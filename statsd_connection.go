package main

import (
	"log"

	datadog "github.com/DataDog/datadog-go/statsd"
	"github.com/urfave/cli"
	"gopkg.in/alexcesaro/statsd.v2"
	"time"
)

func initializeStatsdConnection(c *cli.Context) *statsdConnection {
	// Debug information
	debug := c.Bool("debug")

	// DataDog enabled?
	datadogEnabled := c.Bool("datadog")
	if datadogEnabled {
		log.Println("DataDog enabled.")
	}

	// StatsD host
	host := c.String("statsd")
	if len(host) < 1 {
		log.Fatalln("StatsD host not defined!")
	}
	if debug {
		log.Printf("StatsD host: %s", host)
	}

	// StatsD prefix
	prefix := c.String("prefix")
	if len(prefix) < 1 {
		prefix = "s3b"
		log.Printf("StatsD prefix not defined!  Using standard [%s] prefix", prefix)
	} else {
		if debug {
			log.Printf("StatsD prefix: [%s]", prefix)
		}
	}

	// Client
	var statsDClient interface{}
	var err error
	if datadogEnabled {
		statsDClient, err = datadog.New(host)
		if err != nil {
			log.Fatal(err)
		}
		statsDClient.(*datadog.Client).Namespace = prefix + "."
	} else {
		statsDClient, err = statsd.New(statsd.Address(host),
			statsd.Prefix(prefix), statsd.ErrorHandler(func(err error) {
				log.Fatalf("StatsD Error: %v", err)
			}))
		if err != nil {
			log.Fatal(err)
		}
	}

	// Return the new structure
	return &statsdConnection{
		host:    host,
		prefix:  prefix,
		client:  statsDClient,
		datadog: datadogEnabled,
	}
}

func (s statsdConnection) Increment(name string, tags []string, rate float64) error {
	if s.datadog {
		return s.client.(*datadog.Client).Incr(name, tags, rate)
	}
	s.client.(*statsd.Client).Increment(name)
	return nil
}

func (s statsdConnection) Timing(name string, value time.Duration, tags []string, rate float64) error {
	if s.datadog {
		return s.client.(*datadog.Client).Timing(name, value, tags, rate)
	}
	s.client.(*statsd.Client).Timing(name, value)
	return nil
}
