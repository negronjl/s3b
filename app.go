package main

import (
	"time"

	"github.com/urfave/cli"
)

func s3b_flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "server, s",
			Usage:  "Target S3 server to test",
			EnvVar: "S3_SERVER",
		},
		cli.StringFlag{
			Name:   "region, r",
			Usage:  "S3 server region to use",
			EnvVar: "S3_REGION",
			Value: "us-east-1",
		},
		cli.StringFlag{
			Name:   "access-key, A",
			Usage:  "Access Key to the S3 server",
			EnvVar: "S3_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "secret-key, R",
			Usage:  "Access Key to the S3 server",
			EnvVar: "S3_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "api-signature, a",
			Usage:  "API Signature version (v2 or v4)",
			EnvVar: "S3_API_SIGNATURE",
			Value: "v4",
		},
		cli.BoolFlag{
			Name:   "SSL, S",
			Usage:  "Whether or not to use SSL to connect to the S3 server",
			EnvVar: "S3_SSL",
		},
		cli.BoolFlag{
			Name:   "debug, d",
			Usage:  "Print debug HTTP tracing information",
			EnvVar: "S3_DEBUG",
		},
		cli.StringFlag{
			Name:   "statsd, D",
			Usage:  "StatsD server to which metrics will be sent",
			EnvVar: "S3_STATSD_HOST",
			Value: "s3b",
		},
		cli.StringFlag{
			Name:   "prefix, p",
			Usage:  "Prefix to use with the StatsD metrics",
			EnvVar: "S3_STATSD_PREFIX",
		},
		cli.StringFlag{
			Name:   "matrix, m",
			Usage:  "Comma separated key value pairs of filename=size to use in the testing.",
			EnvVar: "S3_TEST_MATRIX",
		},
		cli.StringFlag{
			Name:   "matrix-dir, M",
			Usage:  "Directory containing the files to be used for testing.",
			EnvVar: "S3_TEST_MATRIX_DIR",
		},
	}
}

func s3b_action(c *cli.Context) error {
	//run_test(initialize_environment(c))
	for {
		run_test(initialize_environment(c))
	}
	return nil
}

func s3b_app() *cli.App {
	app := cli.NewApp()
	app.Name = "s3b"
	app.Usage = "S3/Object Store benchmarking tool"
	app.Version = "0.0.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name:  "Juan L. Negron",
			Email: "negronjl@xtremeghost.com",
		},
	}
	app.Flags = s3b_flags()
	app.Action = s3b_action
	return app
}
