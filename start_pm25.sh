#!/bin/bash

set -e

cd /home/pi/dev/src/pm25
source .venv/bin/activate
source .secret.sh
while true ; do
	./pm25.py || true
	sleep 5
done
