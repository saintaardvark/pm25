#!/bin/bash

set -e

source .secret.sh
while true ; do
        ./wx-logger 2>&1 || true
        sleep 5
done
