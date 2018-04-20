package main

import (
	"github.com/minio/minio-go"
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

type statsd_connection struct {
	host    string
	prefix  string
	client  interface{}
	datadog bool
}

type test_element struct {
	tag          string
	tmp_filename string
	file_size    uint64
}

type test_matrix struct {
	agent_id          string
	connection_object *s3_connection
	statsd_object     *statsd_connection
	test_elements     []test_element
	debug             bool
}
