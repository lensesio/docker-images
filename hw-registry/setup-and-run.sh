#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

exec /usr/bin/supervisord -c /etc/supervisord.conf
