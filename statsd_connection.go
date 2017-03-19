package main

import (
	"log"

	"github.com/urfave/cli"
	datadog "github.com/DataDog/datadog-go/statsd"
	statsd "gopkg.in/alexcesaro/statsd.v2"
	"time"
)

func initialize_statsd_connection(c *cli.Context) *statsd_connection {
	// Debug information
	debug := c.Bool("debug")

	// DataDog enabled?
	datadog_enabled := c.Bool("datadog")
	if datadog_enabled {
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
	var statsd_client interface{}
	var err error
	if datadog_enabled {
		statsd_client, err = datadog.New(host)
		if err != nil {
			log.Fatal(err)
		}
		statsd_client.(*datadog.Client).Namespace = prefix + "."
	} else {
		statsd_client, err = statsd.New(statsd.Address(host),
			statsd.Prefix(prefix), statsd.ErrorHandler(func(err error) {
				log.Fatalf("StatsD Error: %v", err)
			}))
		if err != nil {
			log.Fatal(err)
		}
	}

	// Return the new structure
	return &statsd_connection{
		host:    host,
		prefix:  prefix,
		client:  statsd_client,
		datadog: datadog_enabled,
	}
}

func (s statsd_connection) Increment(name string, tags []string, rate float64) error {
	if s.datadog {
		return s.client.(*datadog.Client).Incr(name, tags, rate)
	} else {
		s.client.(*statsd.Client).Increment(name)
		return nil
	}
}

func (s statsd_connection) Timing(name string, value time.Duration, tags []string, rate float64) error {
	if s.datadog {
		return s.client.(*datadog.Client).Timing(name, value, tags, rate)
	} else {
		s.client.(*statsd.Client).Timing(name, value)
		return nil
	}
}
