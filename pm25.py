#!/usr/bin/env python3

# pm25.py: simple InfluxDB logger for SDS011 air quality meter.
# Copyright (C) 2021 Hugh Brown

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.


import logging
import os
import serial
import time

from sds011 import SDS011
from influxdb import InfluxDBClient


DEFAULT_SERIAL_PORT = "/dev/ttyUSB0"


def read_sensor_data(ser):
    """
    ser: opened serial port
    """

    while True:
        data = []
        for i in range(0, 10):
            datum = ser.read()
            data.append(datum)

        pm25 = int.from_bytes(b"".join(data[2:4]), byteorder="little") / 10
        pm10 = int.from_bytes(b"".join(data[4:6]), byteorder="little") / 10

        return {"pm25": pm25, "pm10": pm10}


def build_influxdb_data(data):
    """
    Build influxdb data out of, uh, data
    """
    logger = logging.getLogger(__name__)
    logger.info("Building influxdb data...")

    influx_data = []
    logger.debug(data)
    measurement = {
        "measurement": "pm",
        "fields": {
            "pm25": data["pm25"],
            "pm10": data["pm10"],
        },
        "tags": {
            "location": "NEWWEST",
        },
    }
    influx_data.append(measurement)

    logger.debug(influx_data)
    return influx_data


def build_influxdb_client():
    """
    Build and return InfluxDB client
    """
    # Setup influx client
    logger = logging.getLogger(__name__)

    DB = os.getenv("INFLUX_DB")
    host = os.getenv("INFLUX_HOST")
    port = os.getenv("INFLUX_PORT")
    INFLUX_USER = os.getenv("INFLUX_USER")
    INFLUX_PASS = os.getenv("INFLUX_PASS")

    influx_client = InfluxDBClient(
        host=host,
        port=port,
        username=INFLUX_USER,
        password=INFLUX_PASS,
        database=DB,
        ssl=True,
        verify_ssl=True,
    )
    logger.info("Connected to InfluxDB version {}".format(influx_client.ping()))
    return influx_client


def write_influx_data(influx_data, influx_client):
    """
    Write influx_data to database
    """
    logger = logging.getLogger(__name__)
    logger.info("Writing data to influxdb...")

    influx_client.write_points(influx_data, time_precision="s")


def main():
    """
    Main entry point
    """
    logger = logging.getLogger(__name__)
    serial_port = os.getenv("SERIAL_PORT", DEFAULT_SERIAL_PORT)
    sds_client = SDS011(serial_port)
    sds_client.set_work_period(work_time=2)

    influx_client = build_influxdb_client()
    while True:
        # time.sleep(110)
        # data = read_sensor_data(ser)
        data = {}
        data["pm25"], data["pm10"] = sds_client.query()
        influx_data = build_influxdb_data(data)
        write_influx_data(influx_data, influx_client)
        logger.debug(influx_data)
        time.sleep(120)


if __name__ == "__main__":
    log_fmt = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
    logging.basicConfig(level=logging.DEBUG, format=log_fmt)

    main()
