package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

type influxTags map[string]string

type influxLogger struct {
	ic       client.Client
	_name    string
	database string
	tags     influxTags
}

func (i *influxLogger) init() error {
	var err error

	log.Println("[INFO] Initializing InfluxDB values...")

	i.database = lookUpFromEnvOrDie("INFLUXDB_DB")

	log.Println("[INFO] Setting up InfluxDB client...")
	i.ic, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     lookUpFromEnvOrDie("INFLUXDB_ADDR"),
		Username: lookUpFromEnvOrDie("INFLUXDB_USER"),
		Password: lookUpFromEnvOrDie("INFLUXDB_PASS"),
		Timeout:  (time.Duration(10) * time.Second),
	})
	if err != nil {
		log.Println("[WARN] Could not set up new InfluxDB client")
		return err
	}
	log.Println("[INFO] Setting up InfluxDB tags...")
	i.tags = influxTags{
		"arduino":  lookUpFromEnvOrDie("NODE"),
		"location": lookUpFromEnvOrDie("LOCATION"),
		"lat":      lookUpFromEnvOrDie("LOC_LAT"),
		"long":     lookUpFromEnvOrDie("LOC_LONG"),
	}

	// log.Printf("[DEBUG] Right after setting, i.ic == |%v|\n", i.ic)
	i._name = "InfluxDB"
	return nil
}

// logToInfluxdb sends a measurement to an InfluxDB server
func (i *influxLogger) log(m Measurement) error {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  i.database,
		Precision: "s",
	})
	if err != nil {
		return fmt.Errorf("Error creating InfluxDB batch: %s", err)
	}

	// Create a point and add to batch
	fields := map[string]interface{}{
		m.Name: m.Value,
	}

	abbrevName := GetAbbrev(m.Name)

	// The soil temperature probes are the first time I actually have
	// multiple instances of a thing.  I'm going to handle this with a
	// special case :grimace: ; if I get more, I'll need to refactor
	// here.

	pt, err := client.NewPoint(abbrevName, i.tags, fields, time.Now())
	if abbrevName == "soil_temp" {
		// Format: "soil_temp_3"
		whichProbe := m.Name[strings.LastIndex(m.Name, "_")+1:]
		tags := i.tags
		tags["probe"] = whichProbe
		pt, err = client.NewPoint(abbrevName, tags, fields, time.Now())
	}

	if err != nil {
		return fmt.Errorf("error in client.NewPoint: %s", err)
	}
	bp.AddPoint(pt)
	// Write the batch
	if err := i.ic.Write(bp); err != nil {
		return fmt.Errorf("error writing to Influxdb: %s", err)
	}
	return nil
}

func (i *influxLogger) name() string {
	return i._name
}

// GetAbbrev returns abbreviation for a particular measurement
func GetAbbrev(name string) string {
	log.Printf("[DEBUG] Trying to find match for |%s|\n", name)
	// See comment in log() up above for details about the soil
	// temperature probes.
	measureAbbrevs := map[string]string{
		"Humd":        "humidity",
		"Prcp":        "precipitation",
		"PrcpMtr":     "precipitation_meter",
		"Pres":        "pressure",
		"Temp":        "temperature",
		"soil_temp_1": "soil_temp",
		"soil_temp_2": "soil_temp",
		"soil_temp_3": "soil_temp",
	}
	return measureAbbrevs[name]
}
