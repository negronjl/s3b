package main

import (
	"github.com/minio/minio-go"
	"gopkg.in/alexcesaro/statsd.v2"
)

type s3_connection struct {
	s3_server     string
	s3_region     string
	s3_access_key string
	s3_secret_key string
	api_signature string
	ssl           bool
	minioClient   *minio.Client
}

type test_element struct {
	tag          string
	tmp_filename string
	file_size    uint64
}

type test_matrix struct {
	agent_id          string
	connection_object *s3_connection
	test_elements     []test_element
	statsd_client     *statsd.Client
	debug             bool
}
