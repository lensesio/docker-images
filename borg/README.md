# Borg #

Borg is an image based on Centos latest that contains the latest release of
[Borg, the deduplicating archiver](https://borgbackup.readthedocs.io).

It is a simple and relatively lightweight image. In Landoop we use Borg for some
of our backup-needs which include backing up docker volumes. This images helps
us with this task.

We usually run it as:

    $ docker run --rm --volumes-from [CONTAINER] -v [BORG_REPO]:[BORG_REPO] landoop/borg create --exclude-caches [BORG_REPO]::[BACKUP NAME] /path/to/backup
