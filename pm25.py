#!/usr/bin/env python3

import logging
import os
import serial
import time

import click
from influxdb import InfluxDBClient

ser = serial.Serial('/dev/ttyUSB0')


def read_sensor_data(ser):
    """
    ser: opened serial port
    """

    while True:
        data = []
        for i in range(0, 10):
            datum = ser.read()
            data.append(datum)

        pm25 = int.from_bytes(b''.join(data[2:4]), byteorder='little') / 10
        pm10 = int.from_bytes(b''.join(data[4:6]), byteorder='little') / 10

        return({'pm25': pm25, 'pm10': pm10})


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
            "pm25": data['pm25'],
            "pm10": data['pm10'],
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

    DB = "waether"
    host = os.getenv("INFLUX_HOST")
    port = os.getenv("INFLUX_PORT")
    INFLUX_USER = "weather"
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


@click.command(short_help="Import CSV file")
@click.argument("csv_file", type=click.File("r"))
def main(csv_file):
    """
    Main entry point
    """
    logger = logging.getLogger(__name__)
    ser = serial.Serial('/dev/ttyUSB0')
    influx_client = build_influxdb_client()
    while True:
        data = read_sensor_data(ser)
        influx_data = build_influxdb_data(data)
        write_influx_data(influx_data, influx_client)
        logger.info("Logged")
        time.sleep(5)


if __name__ == "__main__":
    log_fmt = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
    logging.basicConfig(level=logging.INFO, format=log_fmt)

    main()
