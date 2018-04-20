package main

import (
	"github.com/satori/go.uuid"
	"github.com/urfave/cli"
	"log"
)

func initialize_environment(c *cli.Context) *test_matrix {
	agent_id := uuid.NewV4().String()
	if c.Bool("debug") {
		log.Printf("Agent ID: [%s]", agent_id)
	}
	return initialize_test_matrix(agent_id, initialize_s3_connection(c), initialize_statsd_connection(c), c)
}
