package main

import (
	"log"

	"github.com/minio/minio-go"
	"github.com/urfave/cli"
)

// Validates that the required environment variables
// have been set.
// Returns an s3Connection
// On error, it log.Fatal an exits(1)
func initializeS3Connection(c *cli.Context) *s3Connection {

	// Debug information
	debug := c.Bool("debug")
	if debug {
		log.Println("Debug enabled")
	}

	// S3 Server
	s3Server := c.String("server")
	if len(s3Server) < 1 {
		cli.ShowAppHelp(c)
		log.Fatalln("S3 server not defined!")
	} else {
		if debug {
			log.Printf("S3 Server: %s", s3Server)
		}
	}

	// Most S3 servers need a region
	s3Region := c.String("region")
	if len(s3Region) < 1 {
		cli.ShowAppHelp(c)
		log.Fatalln("S3 region not defined!")
	} else {
		if debug {
			log.Printf("S3 Region: %s", s3Region)
		}
	}

	// S3 Access Key
	s3AccessKey := c.String("access-key")
	if len(s3AccessKey) < 1 {
		cli.ShowAppHelp(c)
		log.Fatalln("S3 Access Key not defined!")
	} else {
		if debug {
			log.Printf("S3 Access Key: %s", s3AccessKey)
		}
	}

	// S3 Secret Key
	s3SecretKey := c.String("secret-key")
	if len(s3SecretKey) < 1 {
		cli.ShowAppHelp(c)
		log.Fatalln("S3 Secret Key not defined!")
	} else {
		if debug {
			log.Printf("S3 Secret Key: %s", s3SecretKey)
		}
	}

	// API Signature
	apiSignature := c.String("api-signature")
	if debug {
		log.Printf("Using S3 Signature API: %s", apiSignature)
	}

	// Are we connecting over SSL?
	ssl := c.Bool("SSL")
	if ssl {
		if debug {
			log.Println("SSL enabled")
		}
	} else {
		if debug {
			log.Println("SSL disabled")
		}
	}

	// Initialize minio client object.
	minioClient := new(minio.Client)
	var err error
	switch apiSignature {
	case "v2":
		minioClient, err = minio.NewV2(s3Server, s3AccessKey, s3SecretKey, ssl)
	case "v4":
		minioClient, err = minio.NewV4(s3Server, s3AccessKey, s3SecretKey, ssl)
	default:
		minioClient, err = minio.New(s3Server, s3AccessKey, s3SecretKey, ssl)
	}
	if err != nil {
		log.Fatalln("Error connecting to S3 server.\n", err)
	} else {
		if debug {
			log.Println("S3 client initialized")
		}
	}

	return &s3Connection{
		s3Server:     s3Server,
		s3Region:     s3Region,
		s3AccessKey:  s3AccessKey,
		s3SecretKey:  s3SecretKey,
		apiSignature: apiSignature,
		ssl:          ssl,
		minioClient:  minioClient}
}
