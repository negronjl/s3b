package main

import (
	"github.com/minio/minio-go"
)

type s3Connection struct {
	s3Server     string
	s3Region     string
	s3AccessKey  string
	s3SecretKey  string
	apiSignature string
	ssl          bool
	minioClient  *minio.Client
}

type statsdConnection struct {
	host    string
	prefix  string
	client  interface{}
	datadog bool
}

type testElement struct {
	tag         string
	tmpFilename string
	fileSize    uint64
}

type testMatrix struct {
	agentId          string
	connectionObject *s3Connection
	statsdObject     *statsdConnection
	testElements     []testElement
	debug            bool
}
