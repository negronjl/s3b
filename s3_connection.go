package main

import (
	"log"

	"github.com/minio/minio-go"
	"github.com/urfave/cli"
)

// Validates that the required environment variables
// have been set.
// Returns an s3_connection
// On error, it log.Fatal an exits(1)
func initialize_s3_connection(c *cli.Context) *s3_connection {

	// Debug information
	debug := c.Bool("debug")
	if debug {
		log.Println("Debug enabled")
	}

	// S3 Server
	s3_server := c.String("server")
	if len(s3_server) < 1 {
		cli.ShowAppHelp(c)
		log.Fatalln("S3 server not defined!")
	} else {
		if debug {
			log.Printf("S3 Server: %s", s3_server)
		}
	}

	// Most S3 servers need a region
	s3_region := c.String("region")
	if len(s3_region) < 1 {
		cli.ShowAppHelp(c)
		log.Fatalln("S3 region not defined!")
	} else {
		if debug {
			log.Printf("S3 Region: %s", s3_region)
		}
	}

	// S3 Access Key
	s3_access_key := c.String("access-key")
	if len(s3_access_key) < 1 {
		cli.ShowAppHelp(c)
		log.Fatalln("S3 Access Key not defined!")
	} else {
		if debug {
			log.Printf("S3 Access Key: %s", s3_access_key)
		}
	}

	// S3 Secret Key
	s3_secret_key := c.String("secret-key")
	if len(s3_secret_key) < 1 {
		cli.ShowAppHelp(c)
		log.Fatalln("S3 Secret Key not defined!")
	} else {
		if debug {
			log.Printf("S3 Secret Key: %s", s3_secret_key)
		}
	}

	// API Signature
	api_signature := c.String("api-signature")
	if debug {
		log.Printf("Using S3 Signature API: %s", api_signature)
	}

	// Are we connecting over SSL?
	ssl := c.Bool("SSL")
	if ssl {
		if debug {
			log.Printf("SSL enabled")
		}
	} else {
		if debug {
			log.Printf("SSL disabled")
		}
	}

	// Initialize minio client object.
	minioClient := new(minio.Client)
	var err error
	if api_signature == "v2" {
		minioClient, err = minio.NewV2(s3_server, s3_access_key, s3_secret_key, ssl)
	} else if api_signature == "v4" {
		minioClient, err = minio.NewV4(s3_server, s3_access_key, s3_secret_key, ssl)
	} else {
		minioClient, err = minio.New(s3_server, s3_access_key, s3_secret_key, ssl)
	}
	if err != nil {
		log.Fatalln("Error connecting to S3 server.\n", err)
	} else {
		if debug {
			log.Printf("S3 client initialized")
		}
	}

	return &s3_connection{
		s3_server:     s3_server,
		s3_region:     s3_region,
		s3_access_key: s3_access_key,
		s3_secret_key: s3_secret_key,
		api_signature: api_signature,
		ssl:           ssl,
		minioClient:   minioClient}
}
