package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

var (
	githash    = "deadbeef"
	buildstamp = "June 6, 2017"
)

// Measurement holds information about a single weather reading
type Measurement struct {
	Name  string
	Value float64
	Units string
}

// This lets me compare directly in the test file.  See:
// https://stackoverflow.com/questions/36091610/comparing-errors-in-go
var (
	errColon          = fmt.Errorf("can't find colon in line, don't know how to split it")
	errIncompleteRead = fmt.Errorf("can't find closing '}' -- incomplete read?")
	measure, value    string
)

// SplitLine spits a string and returns a Measurement struct and err
func SplitLine(s string) (measure Measurement, err error) {
	// FIXME: Account for errors in all this
	m := Measurement{"", 0.0, ""}
	if strings.Index(s, "}") < 0 {
		return m, errIncompleteRead
	}
	if strings.Index(s, ":") < 0 {
		return m, errColon
	}
	// Don't include left bracket itself
	leftBracketPos := strings.Index(s, "{") + 1
	rightBracketPos := strings.Index(s, "}")
	s = s[leftBracketPos:rightBracketPos]
	// log.Printf("[DEBUG] After Trim: %s", s)
	line := strings.Split(s, ":")
	m.Name = line[0]
	// Before parsing, need to remove units:
	// Humidity -> "%"
	// Pressure -> "hP"
	// Prcp -> "NA"
	// Temp -> "C"
	vals := strings.Fields(line[1])
	m.Units = vals[1]
	if m.Value, err = strconv.ParseFloat(vals[0], 64); err != nil {
		err = errors.New("can't figure out value of that measurement")
	}
	return m, err
}

func startupLog() {
	log.Printf("Githash: %s\n", githash)
	log.Printf("Build date: %s\n", buildstamp)
}

func main() {
	startupLog()
	allLoggers := []wxLogger{}

	influx := new(influxLogger)
	if err := influx.init(); err != nil {
		log.Fatalf("[FATAL]: Can't log to InfluxDB: %s", err.Error())
	} else {
		allLoggers = append(allLoggers, influx)
	}

	var wunder wundergroundLogger
	if err := wunder.init(); err != nil {
		log.Printf("[WARN] Can't log to wunderground: %s", err.Error())
	} else {
		// FIXME: Why isn't the init() taking care of this?
		// Arghh: I think I need to be passing around a pointer or something to that init file.
		log.Printf("[DEBUG] wunder: %+v\n", wunder)
		wunder._name = "wunderground"
		allLoggers = append(allLoggers, wunder)

	}

	usbdev := lookUpFromEnvOrDie("USBDEV")
	// Most of the time we'll be logging from Arduino, where I've set
	// baud rate to 9600.  But I'm looking at re-using this logger for
	// the ESP32, which has a hard-coded baud rate of 115200.
	//
	// TODO: Break this out to separate function
	usbBaudRate := 9600
	retVal, exists := os.LookupEnv("USBBAUDRATE")
	if exists {
		if retVal, err := strconv.Atoi(retVal); err == nil {
			usbBaudRate = retVal
		}
	}
	c := &serial.Config{Name: usbdev, Baud: usbBaudRate}
	serialPort, err := serial.OpenPort(c)
	for err != nil {
		log.Printf("[WARN] Can't open serial port, trying to sleep it off: %s\n", err)
		time.Sleep(10 * time.Second)
		serialPort, err = serial.OpenPort(c)
	}
	reader := bufio.NewReader(serialPort)
	log.Println("[INFO] Opened. Next up: looping.")
	for {
		log.Println("[DEBUG] About to read...")
		reply, err := reader.ReadString('}')
		if err != nil {
			log.Println("[WARN] Problem reading: ", err)
			// Sleep for a second
			time.Sleep(1 * time.Second)
			continue
		}
		log.Println(reply)
		measure, err := SplitLine(reply)
		if err != nil {
			log.Printf("[WARN] Could not split line: %s", err)
			// Sleep for a second
			time.Sleep(1 * time.Second)
			continue
		}
		log.Printf("[INFO] Read: %s: %f\n", measure.Name, measure.Value)
		for _, l := range allLoggers {
			if err = l.log(measure); err != nil {
				log.Printf("[WARN] Problem logging to %s: %s", l.name(), err)
			}
		}
	}
}
