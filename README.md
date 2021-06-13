# pm25

A simple logger for an SDS011 air quality monitor, based on [this
article](https://www.raspberrypi.org/blog/monitor-air-quality-with-a-raspberry-pi/).

## Usage

- Create `.secret.sh`:

```
export INFLUX_PASS=my_secret_password
export INFLUX_HOST=127.0.0.1
export INFLUX_PORT=8086
```

- Run `make setup`

- Run:

```
source ./.secret.sh
./pm25.py
```
