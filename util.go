package main

import (
	"log"
	"os"
)

// lookUpFromEnvOrDie looks up and returns an environment value, and dies
// if it can't find it
func lookUpFromEnvOrDie(envVar string) string {
	retVal, exists := os.LookupEnv(envVar)
	if exists == false {
		log.Fatalf("[FATAL] Can't proceed without environment var %s !", envVar)
	}
	return retVal
}
