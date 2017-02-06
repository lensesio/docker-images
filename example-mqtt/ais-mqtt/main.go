/*
A simple program to decode AIS messages, select Class A Position
reports, encode them to Avro, serialize them to binary and send
them to MQTT. It has a simple worker pool to speed-up the slowest
part of the process (waiting for delivery confirmation from MQTT)
and the number of messages sent is configurable (iterates over
the same source until it reaches the quota).
*/
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	ais "github.com/andmarios/aislib"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/linkedin/goavro"
)

var (
	mqttServerAddr   = flag.String("server", "localhost:1883", "MQTT server endpoint")
	originName       = flag.String("origin", "/ais", "origin queue or topic")
	logFilename      = flag.String("log", "", "file to write output (and logs), stdout if left empty")
	schemaFilename   = flag.String("schema", "/classAPositionReportSchema.json", "file containing the avro schema to use")
	nmeaDataFilename = flag.String("nmea", "/opt/go/src/landoop/ais-mqtt/nmea-sample",
		"file containing nmea sentences (size isn't important but better to have a few thousand Class A sentences)")
	testMessages = flag.Int("messages", 100000, "number of messages to send to kafka")
	numWorkers   = flag.Int("workers", 64, "number of workers to encode messages to avro and send them to MQTT")
)

var workerWg sync.WaitGroup

func main() {
	flag.Parse()

	// If asked, write log to a file.
	if len(*logFilename) > 0 {
		logFile, err := os.OpenFile(*logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error opening log file: %v\n", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	// Load Schema File
	tempBuf, err := ioutil.ReadFile(*schemaFilename)
	if err != nil {
		log.Fatalf("Error opening schema file: %v\n", err)
	}
	schema := string(tempBuf)

	// The msgBus will deliver the decoded messages to the workers.
	msgBus := make(chan ais.ClassAPositionReport, 1024*256)

	// Spawn our workers. They encode the messages to avro and send them to MQTT.
	workerWg.Add(*numWorkers)
	for i := 1; i <= *numWorkers; i++ {
		go mqttWorker(msgBus, schema)
	}

	// This is aislib specific code. It creates an AIS router where we send
	// AIS sentences and the router decodes their type and sends them back
	// in a type - payload struct.
	aisSentences := make(chan string, 1024*128)
	classifiedSentences := make(chan ais.Message, 1024*128)
	failedSentences := make(chan ais.FailedSentence, 1024*128)
	go ais.Router(aisSentences, classifiedSentences, failedSentences)

	// Open the file that contains raw NMEA sentences. It was recorded from
	// a live AIS server. It is trivial in golang to use another source for
	// this, like a TCP connection or stdin.
	nmeaData, err := os.Open(*nmeaDataFilename)
	if err != nil {
		log.Fatalln(err)
	}
	defer nmeaData.Close()

	// This is a function that takes the decoded messages that are Class A
	// Position Reports and sends them to msgBus for processing. It is self
	// feeding, as long as it needs more messages to reach the target we set
	// it will re-supply the AIS router with sentences from the NMEA file,
	// rolling over if needed. It is fairly quick and not the bottleneck of
	// the code.
	workerWg.Add(1)
	go func() {
		defer workerWg.Done()

		var message ais.Message
		// To keep track of messages decoded and send to workers.
		numMessages := 0
	routerLoop:
		for {
			select {
			case message = <-classifiedSentences:
				switch message.Type {
				// If message is class A position report
				case 1, 2, 3:
					// Decode message
					t, _ := ais.DecodeClassAPositionReport(message.Payload)
					msgBus <- t
					numMessages++
					if numMessages == *testMessages {
						close(msgBus)
						log.Printf("Finished processing input. %d Class A "+
							"messages decoded. Waiting for workers to finish"+
							" encoding and sending records to MQTT.\n",
							numMessages)
						break routerLoop
					}
				// Special message for auto-feeding. There isn't an AIS with type 255
				case 255:
					go func() {
						in := bufio.NewScanner(nmeaData)
						in.Split(bufio.ScanLines)

						for in.Scan() {
							aisSentences <- in.Text()
						}
						// Roll to the start of the file
						nmeaData.Seek(0, 0)
						// This message will restart this process
						classifiedSentences <- ais.Message{Type: 255}
					}()
				default:
				}
			case <-failedSentences:
			}
		}
	}()

	log.Println("Starting to proccess messages.")
	// First entry will start our self-feeding functio above.
	classifiedSentences <- ais.Message{Type: 255}

	workerWg.Wait()
	log.Println("All messages processed and send to MQTT.")
}

// mqttWorker receives ClassAPositionReport messages over a channel, encodes
// them to avro and sends them to mqtt
func mqttWorker(msgBus chan ais.ClassAPositionReport, schema string) {
	defer workerWg.Done()

	// Create Avro Codec
	codec, err := goavro.NewCodec(schema)
	if err != nil {
		log.Fatalln(err)
	}

	// Create an Avro Record
	record, err := goavro.NewRecord(goavro.RecordSchema(schema))
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Create a unique string for MQTT client id
	t := time.Now().UnixNano()
	uniqueText := fmt.Sprintf("%d", t)

	// Set connection options and connect to ActiveMQ
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://" + *mqttServerAddr)
	opts.SetClientID("ais-source" + uniqueText)
	opts.SetCleanSession(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)

	tempBuf := new(bytes.Buffer)
	for msg := range msgBus {
		// Set ais messages data into the avro record
		classA2Record(msg, record)
		// Encode the avro record to an avro binary message
		tempBuf.Reset()
		if err = codec.Encode(tempBuf, record); err != nil {
			log.Println(err)
		}
		// Publish the binary message to MQTT. Ask to deliver it at least once (QoS 1)
		token := client.Publish(*originName, byte(1), false, tempBuf.Bytes())
		token.Wait()
	}

}

// class2ARecord will take a ClassAPositionReport and set its values into an
// avro record
func classA2Record(m ais.ClassAPositionReport, r *goavro.Record) {
	r.Set("Type", int32(m.Type))
	r.Set("Repeat", int32(m.Repeat))
	r.Set("MMSI", int64(m.MMSI))
	r.Set("Speed", float32(m.Speed))
	r.Set("Accuracy", m.Accuracy)
	r.Set("Longitude", float64(m.Lon))
	r.Set("Latitude", float64(m.Lat))
	r.Set("Course", float32(m.Course))
	r.Set("Heading", int32(m.Heading))
	r.Set("Second", int32(m.Second))
	r.Set("RAIM", m.RAIM)
	r.Set("Radio", int64(m.Radio))
	r.Set("Status", int32(m.Status))
	r.Set("Turn", float32(m.Turn))
	r.Set("Maneuver", int32(m.Maneuver))
	// Due to the way the messages are processed we can guarantee that
	// messages are almost processed in order and if not, then they
	// are processed with errors in the millisecond level. Class A Position
	// Reports are usually at a frequency from 3 min to 2 secs for very fast
	// moving vehicles, so millisecond accuracy in ordering is enough.
	r.Set("Timestamp", time.Now().UnixNano())
}
