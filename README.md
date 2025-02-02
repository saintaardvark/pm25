# pm25

A simple logger for an SDS011 air quality monitor, based on [this
article](https://www.raspberrypi.org/blog/monitor-air-quality-with-a-raspberry-pi/).

## Usage

- Install python3-serial:

```
sudo apt-get install -y python3-serial
```

- Create `.secret.sh`:

```
export INFLUX_DB=pm25
export INFLUX_USER=pm25
export INFLUX_PASS=my_secret_password
export INFLUX_HOST=127.0.0.1
export INFLUX_PORT=8086
# Default is /dev/ttyUSB0, but can be overridden here:
export SERIAL_PORT=/dev/sds011
```

- Run `make setup`

- Run:

```
source ./.secret.sh
./pm25.py
```

# Resources

- [Environment Canada monitoring station in Burnaby][0]

# License

GPL v3.

# Future options

- https://github.com/ikalchev/py-sds011
- https://github.com/menschel/sds011

[0]: https://aqicn.org/city/british-comlumbia/burnaby-south
