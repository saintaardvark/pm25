package main

import (
	"fmt"
	"log"
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
	pt, err := client.NewPoint(abbrevName, i.tags, fields, time.Now())
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
	measureAbbrevs := map[string]string{
		"Humd":    "humidity",
		"Prcp":    "precipitation",
		"PrcpMtr": "precipitation_meter",
		"Pres":    "pressure",
		"Temp":    "temperature",
	}
	return measureAbbrevs[name]
}
