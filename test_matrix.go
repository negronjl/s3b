package main

import (
	"crypto/rand"
	"fmt"
	"gopkg.in/alexcesaro/statsd.v2"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"github.com/urfave/cli"
	"os"
)

func initialize_test_element(agent_id string, tag_size_string string, debug bool) test_element {

	// Variables
	var element_size uint64

	// Parse the string
	element := strings.Split(tag_size_string, "=")
	element_tag := strings.TrimSpace(element[0])

	// See if the file exists
	element_file, err := os.Stat(element_tag)
	if err != nil { // File doesn't exist

		// Parse the element size
		element_size, err := strconv.ParseUint(element[1], 0, 0)
		if err != nil {
			log.Fatalf("Unable to convert [%s] to integer\n", element[1])
		}

		if debug {
			log.Printf("Creating temp file")
		}

		// Create a temporary file
		element_file, err := ioutil.TempFile("", fmt.Sprintf("%s_%s", agent_id, element_tag))
		if err != nil {
			log.Fatalln("Unable to create temp file!")
		} else {
			if debug {
				log.Printf("Created temp file [%s]", element_file.Name())
			}
		}

		// Fill the file with random data
		bytes_written := uint64(0)
		chunk_size := 1000000
		for bytes_written < element_size {
			var byte_array []byte
			bytes_left := element_size - bytes_written

			// Create a byte array
			if bytes_left < uint64(chunk_size) {
				byte_array = make([]byte, bytes_left)
			} else {
				byte_array = make([]byte, chunk_size)
			}

			// Fill the byte array with random data
			_, err = rand.Read(byte_array)
			if err != nil {
				log.Fatalln("Error reading random data")
			}

			// Fill the temp file with random data
			_, err = element_file.Write(byte_array)
			if err != nil {
				log.Fatalln("Could not write temporary file")
			}
			bytes_written += uint64(len(byte_array))
			log.Printf("Wrote %d/%d bytes to: %s", bytes_written, element_size, element_file.Name())

			// Close the temporary file
			err = element_file.Close()
			if err != nil {
				log.Fatalln("Could not close the temporary file object!")
			}
		}
	} else { // File found
		if debug {
			element_size = uint64(element_file.Size())
			log.Printf("Using existing file: %s", element_file.Name(), element_size)
		}
	}

	if debug {
		log.Printf("Test Element: Tag: [%s], Filename: [%s], Size: [%s]", element_tag,
			element_file.Name(), element_size)
	}

	return test_element{
		tag:          element_tag,
		tmp_filename: element_file.Name(),
		file_size:    element_size}
}

func initialize_test_matrix(agent_id string, connection_object s3_connection, c *cli.Context) test_matrix {
	// Debug
	debug := c.Bool("debug")
	if debug {
		log.Printf("Debug is enabled")
	} else {
		log.Printf("Debug disabled.")
	}

	// StatsD host
	statsd_host := c.String("statsd")
	if len(statsd_host) < 1 {
		log.Fatalln("StatsD host not defined!")
	}
	if debug {
		log.Printf("StatsD host: %s", statsd_host)
	}

	// StatsD prefix
	statsd_app_prefix := c.String("prefix")
	if len(statsd_app_prefix) < 1 {
		statsd_app_prefix = "s3b"
		log.Printf("StatsD prefix not defined!  Using standard [%s] prefix", statsd_app_prefix)
	} else {
		if debug {
			log.Printf("StatsD prefix: [%s]", statsd_app_prefix)
		}
	}

	statsd_client, err := statsd.New(statsd.Address(statsd_host),
		statsd.Prefix(statsd_app_prefix), statsd.ErrorHandler(func(err error) {
			log.Fatalf("StatsD Error: %v", err)
		}))
	if err != nil {
		log.Fatal(err)
	}

	// Register this agent_id with StatsD
	statsd_client.Increment("agent_id")

	var test_elements []test_element

	// Matrix Directory
	matrix_dir := c.String("matrix-dir")
	var matrix_string string
	if len(matrix_dir) > 0 { // Matrix directory defined
		if debug {
			log.Println("Processing matrix directory")
		}
		files, err := ioutil.ReadDir(matrix_dir)
		if err != nil {
			log.Fatalf("Could not read matrix-dir: %s", matrix_dir)
		}
		var matrix_elements []string
		for _, file := range files {
			matrix_string += file.Name() + "=" + strconv.FormatInt(file.Size(), 10)
			matrix_elements = append(matrix_elements, matrix_string)
		}
		matrix_string = strings.Join(matrix_elements, ",")
		if debug {
			log.Printf("Calculated matrix-test of [%s] from matrix-dir of [%s]",
				matrix_string,
				matrix_dir)
		}
	} else { // Matrix directory not defined
		// This variable holds a set of key value pairs.
		// The format is meant to be filename=size, filename=size, etc.
		// tag is used as a measurement data point
		// size is the size of the object to be used for testing
		matrix_string = c.String("matrix")
		if len(matrix_string) < 1 {
			log.Fatalln("Test matrix not defined!")
		}

		for _, element := range strings.Split(matrix_string, ",") {
			test_elements = append(test_elements, initialize_test_element(agent_id, element, debug))
		}
	}

	return test_matrix{
		agent_id:          agent_id,
		connection_object: connection_object,
		test_elements:     test_elements,
		statsd_client:     statsd_client,
		debug:             debug}
}
