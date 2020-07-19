# arduino-wx-go-logger

Does what it says on the tin.  Meant for use with the arduino-wx
codebase.  I've split off this logger to its own repo to try and help
with vendoring, god help me.

# Usage

This looks for a few different environment variables:

- INFLUXDB_USER: username for InfluxDB

- INFLUXDB_PASS: password for INFLUXDB_USER

- INFLUXDB_ADDR: hostname for InfluxDB server

- INFLUXDB_DB: database on the InfluxDB server

- USBDEV: the USB serial device to read from

- NODE: string to describe what node this is ("node1" is what I've
  used in the past)

- LOCATION: string for the location

- LOC_LAT: string for latitude

- LOC_LONG: string for longitude

Suggested values, unless I've changed something:

```
export INFLUXDB_USER=weather
export INFLUXDB_PASS=nice try
export INFLUXDB_ADDR=https://influxdb:8086
export INFLUXDB_DB=weather
export USBDEV=/dev/ttyACM0						# (or /dev/ttyUSB0)
export NODE=node1
export LOCATION=BBY
export LOC_LAT="0.123"
export LOC_LONG="-0.123"
```

There are no default values for these anymore; this is left to deployment.

# Docker

`make docker-build` will build the container.

`make docker` will build the container and start a shell for testing.
