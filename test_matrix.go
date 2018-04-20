package main

import (
	"crypto/rand"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func initializeTestElement(agentId string, tagSizeString string, debug bool) testElement {

	// Variables
	elementSize := uint64(0)
	elementFilename := ""

	// Parse the string
	element := strings.Split(tagSizeString, "=")
	elementTag := strings.TrimSpace(element[0])

	// See if the file exists
	elementFile, err := os.Stat(elementTag)
	if err != nil { // File doesn't exist

		// Parse the element size
		elementSize, err := strconv.ParseUint(element[1], 0, 0)
		if err != nil {
			log.Fatalf("Unable to convert [%s] to integer\n", element[1])
		}

		if debug {
			log.Println("Creating temp file")
		}

		// Create a temporary file
		elementFile, err := ioutil.TempFile("", fmt.Sprintf("%s_%s", agentId, elementTag))
		if err != nil {
			log.Fatalln("Unable to create temp file!")
		} else {
			if debug {
				log.Printf("Created temp file [%s]", elementFile.Name())
			}
		}

		// Fill the file with random data
		bytesWritten := uint64(0)
		chunkSize := 1000000
		for bytesWritten < elementSize {
			var byteArray []byte
			bytesLeft := elementSize - bytesWritten

			// Create a byte array
			if bytesLeft < uint64(chunkSize) {
				byteArray = make([]byte, bytesLeft)
			} else {
				byteArray = make([]byte, chunkSize)
			}

			// Fill the byte array with random data
			_, err = rand.Read(byteArray)
			if err != nil {
				log.Fatalln("Error reading random data")
			}

			// Fill the temp file with random data
			_, err = elementFile.Write(byteArray)
			if err != nil {
				log.Fatalln("Could not write temporary file")
			}
			bytesWritten += uint64(len(byteArray))
			log.Printf("Wrote %d/%d bytes to: %s", bytesWritten, elementSize, elementFile.Name())
		}
		// Close the temporary file
		err = elementFile.Close()
		if err != nil {
			log.Fatalln("Could not close the temporary file object!")
		}
		if debug {
			log.Println("Closed temp file")
		}

		// Update outer elementFilename
		elementFilename = elementFile.Name()
	} else { // File found
		if debug {
			elementSize = uint64(elementFile.Size())
			log.Printf("Using existing file: %s (size: %d)", elementFile.Name(), elementSize)
		}

		// Update outer elementFilename
		elementFilename = elementFile.Name()

		if debug {
			log.Printf("Test Element: Tag: [%s], Filename: [%s], Size: [%d]", elementTag,
				elementFilename, elementSize)
		}

	}

	return testElement{
		tag:         elementTag,
		tmpFilename: elementFilename,
		fileSize:    elementSize}
}

func initializeTestMatrix(agentId string,
	connectionObject *s3Connection,
	statsDObject *statsdConnection,
	c *cli.Context) *testMatrix {

	// Debug
	debug := c.Bool("debug")

	testElements := make([]testElement, 0)

	// Matrix Directory
	matrixDir := c.String("matrix-dir")
	var matrixString string
	if len(matrixDir) > 0 { // Matrix directory defined
		if debug {
			log.Println("Processing matrix directory")
		}
		files, err := ioutil.ReadDir(matrixDir)
		if err != nil {
			log.Fatalf("Could not read matrix-dir: %s", matrixDir)
		}
		var matrixElements []string
		for _, file := range files {
			elementString := file.Name() + "=" + strconv.FormatInt(file.Size(), 10)
			matrixElements = append(matrixElements, elementString)

		}
		matrixString = strings.Join(matrixElements, ",")
		if debug {
			log.Printf("Calculated matrix-test of [%s] from matrix-dir of [%s]",
				matrixString,
				matrixDir)
		}
	} else { // Matrix directory not defined
		// This variable holds a set of key value pairs.
		// The format is meant to be filename=size, filename=size, etc.
		// tag is used as a measurement data point
		// size is the size of the object to be used for testing
		matrixString = c.String("matrix")
		if len(matrixString) < 1 {
			log.Fatalln("Test matrix not defined!")
		}

		for _, element := range strings.Split(matrixString, ",") {
			testElements = append(testElements, initializeTestElement(agentId, element, debug))
		}
	}

	return &testMatrix{
		agentId:          agentId,
		connectionObject: connectionObject,
		testElements:     testElements,
		statsdObject:     statsDObject,
		debug:            debug}
}
