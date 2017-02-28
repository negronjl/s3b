package main

import (
	"fmt"
	"log"
)

func run_test(test_matrix test_matrix) {
	// Initialize statsd
	statsd_client := test_matrix.statsd_client

	// Create agent_id bucket
	timer := statsd_client.NewTiming()
	if test_matrix.connection_object.minioClient.MakeBucket(test_matrix.agent_id,
		test_matrix.connection_object.s3_region) != nil {
		statsd_client.Increment("bucket.create.count_failed")
		if test_matrix.debug {
			log.Printf("Unable to create bucket! [%s]", test_matrix.agent_id)
		}
	} else {
		statsd_client.Increment("bucket.create.count_succeed")
		if test_matrix.debug {
			log.Printf("Created bucket [%s] in [%d]ms", test_matrix.agent_id,
				uint64(timer.Duration().Nanoseconds()/1000/1000))
		}
		timer.Send("bucket.create.time")
	}
	statsd_client.Increment("bucket.create.count_total")

	// iterate over the test_elements
	for _, test_element := range test_matrix.test_elements {
		// Upload the file to S3
		timer = statsd_client.NewTiming()
		_, err := test_matrix.connection_object.minioClient.FPutObject(test_matrix.agent_id, test_element.tag,
			test_element.tmp_filename, "application/octet-stream")
		if err != nil {
			if test_matrix.debug {
				log.Printf("Unable to upload file [%s] to %s/%s",
					test_element.tmp_filename,
					test_matrix.agent_id,
					test_element.tag)
			}
			statsd_client.Increment(fmt.Sprintf("s3.upload.%s.count_failed", test_element.tag))
		} else {
			statsd_client.Increment(fmt.Sprintf("s3.upload.%s.count_succeed", test_element.tag))
			if test_matrix.debug {
				log.Printf("Uploaded file [%s] to %s/%s in [%d]ms",
					test_element.tmp_filename,
					test_matrix.agent_id,
					test_element.tag,
					uint64(timer.Duration().Nanoseconds()/1000/1000))
			}
			timer.Send(fmt.Sprintf("s3.upload.%s.time", test_element.tag))
		}
		statsd_client.Increment(fmt.Sprintf("s3.upload.%s.count_total", test_element.tag))

		// Stat the newly created S3 Object
		timer = statsd_client.NewTiming()
		objInfo, err := test_matrix.connection_object.minioClient.StatObject(test_matrix.agent_id,
			test_element.tag)
		if err != nil {
			if test_matrix.debug {
				log.Printf("Unable to stat file [%s] from bucket [%s]",
					test_element.tag,
					test_matrix.agent_id)
			}
			statsd_client.Increment(fmt.Sprintf("s3.stat.%s.count_failed", test_element.tag))
		} else {
			statsd_client.Increment(fmt.Sprintf("s3.stat.%s.count_succeed", test_element.tag))
			if test_matrix.debug {
				log.Printf("Got status of object [%s] from bucket [%s] in [%d]ms",
					test_element.tag,
					test_matrix.agent_id, uint64(timer.Duration().Nanoseconds()/1000/1000))
			}
			log.Println(objInfo)
			timer.Send(fmt.Sprintf("s3.stat.%s.time", test_element.tag))
		}
		statsd_client.Increment(fmt.Sprintf("s3.stat.%s.count_total", test_element.tag))

		// Download the file from S3
		timer = statsd_client.NewTiming()
		if test_matrix.connection_object.minioClient.FGetObject(test_matrix.agent_id,
			test_element.tag, test_element.tmp_filename) != nil {
			if test_matrix.debug {
				log.Printf("Unable to get object [%s] from bucket [%s]",
					test_element.tag,
					test_matrix.agent_id)
			}
			statsd_client.Increment(fmt.Sprintf("s3.download.%s.count_failed", test_element.tag))
		} else {
			statsd_client.Increment(fmt.Sprintf("s3.download.%s.count_succeed", test_element.tag))
			if test_matrix.debug {
				log.Printf("Downloaded object [%s] to file [%s] in [%d]ms",
					test_element.tag,
					test_element.tmp_filename,
					uint64(timer.Duration().Nanoseconds()/1000/1000))
			}
			timer.Send(fmt.Sprintf("s3.download.%s.time", test_element.tag))
		}
		statsd_client.Increment(fmt.Sprintf("s3.download.%s.count_total", test_element.tag))

		// Delete the file from S3
		timer = statsd_client.NewTiming()
		if test_matrix.connection_object.minioClient.RemoveObject(test_matrix.agent_id,
			test_element.tag) != nil {
			if test_matrix.debug {
				log.Printf("Unable to remove object [%s] from bucket [%s]",
					test_element.tag,
					test_matrix.agent_id)
			}
			statsd_client.Increment(fmt.Sprintf("s3.delete.%s.count_failed", test_element.tag))
		} else {
			statsd_client.Increment(fmt.Sprintf("s3.delete.%s.count_succeed", test_element.tag))
			if test_matrix.debug {
				log.Printf("Deleted object [%s] in [%d]ms",
					test_element.tag,
					uint64(timer.Duration().Nanoseconds()/1000/1000))
			}
			timer.Send(fmt.Sprintf("s3.delete.%s.time", test_element.tag))
		}
		statsd_client.Increment(fmt.Sprintf("s3.delete.%s.count_total", test_element.tag))
	}
	// Delete agent_id bucket
	timer = statsd_client.NewTiming()
	if test_matrix.connection_object.minioClient.RemoveBucket(test_matrix.agent_id) != nil {
		if test_matrix.debug {
			log.Printf("Unable to remove bucket! ", test_matrix.agent_id)
		}
		statsd_client.Increment("bucket.delete.count_failed")
	} else {
		statsd_client.Increment("bucket.delete.count_succeed")
		if test_matrix.debug {
			log.Printf("Deleted bucket [%s] in [%d]ms", test_matrix.agent_id,
				uint64(timer.Duration().Nanoseconds()/1000/1000))
		}
		timer.Send("bucket.delete.time")
	}
	statsd_client.Increment("bucket.delete.count_total")
}
