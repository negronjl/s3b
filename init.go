package main

import (
	"github.com/satori/go.uuid"
	"github.com/urfave/cli"
	"log"
)

func initializeEnvironment(c *cli.Context) *testMatrix {
	agentId := uuid.NewV4().String()
	if c.Bool("debug") {
		log.Printf("Agent ID: [%s]", agentId)
	}
	return initializeTestMatrix(agentId, initializeS3Connection(c), initializeStatsdConnection(c), c)
}
