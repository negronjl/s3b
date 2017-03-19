package main

import (
	"fmt"
	"log"
	"time"
)

func run_test(test_matrix *test_matrix) {

	// Register this agent_id with StatsD
	test_matrix.statsd_object.Increment("agent_id",
		[]string{
			test_matrix.statsd_object.prefix,
			fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
		}, 1)

	// Create agent_id bucket
	timer := time.Now()
	if test_matrix.connection_object.minioClient.MakeBucket(test_matrix.agent_id,
		test_matrix.connection_object.s3_region) != nil {
		test_matrix.statsd_object.Increment("bucket.create.failed",
			[]string{
				test_matrix.statsd_object.prefix,
				fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
				"bucket_create:failed",
			}, 1)
		if test_matrix.debug {
			log.Printf("Unable to create bucket! [%s]", test_matrix.agent_id)
		}
	} else {
		test_matrix.statsd_object.Increment("bucket.create.succeed",
			[]string{
				test_matrix.statsd_object.prefix,
				fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
				"bucket_create:succeed",
			}, 1)
		if test_matrix.debug {
			log.Printf("Created bucket [%s] in [%d]ms", test_matrix.agent_id,
				uint64(time.Since(timer)))
		}
		test_matrix.statsd_object.Timing("bucket.create.time", time.Since(timer),
			[]string{
				test_matrix.statsd_object.prefix,
				fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
				"bucket_create:time",
			}, 1)
	}
	test_matrix.statsd_object.Increment("bucket.create.total",
		[]string{
			test_matrix.statsd_object.prefix,
			fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
			"bucket_create:total",
		}, 1)

	// iterate over the test_elements
	for _, test_element := range test_matrix.test_elements {
		// Calculate the bucket names
		bucket_names := make(map[string]string)
		if test_matrix.statsd_object.datadog {
			bucket_names = map[string]string{
				// Upload
				"s3_upload_failed":  "s3.upload.failed",
				"s3_upload_succeed": "s3.upload.succeed",
				"s3_upload_time":    "s3.upload.time",
				"s3_upload_total":   "s3.upload.total",
				// Stat
				"s3_stat_failed":  "s3.stat.failed",
				"s3_stat_succeed": "s3.stat.succeed",
				"s3_stat_time":    "s3.stat.time",
				"s3_stat_total":   "s3.stat.total",
				// Download
				"s3_download_failed":  "s3.download.failed",
				"s3_download_succeed": "s3.download.succeed",
				"s3_download_time":    "s3.download.time",
				"s3_download_total":   "s3.download.total",
				// Delete
				"s3_delete_failed":  "s3.delete.failed",
				"s3_delete_succeed": "s3.delete.succeed",
				"s3_delete_time":    "s3.delete.time",
				"s3_delete_total":   "s3.delete.total",
			}
		} else {
			bucket_names = map[string]string{
				// Upload
				"s3_upload_failed":  fmt.Sprintf("s3.upload.%s.failed", test_element.tag),
				"s3_upload_succeed": fmt.Sprintf("s3.upload.%s.succeed", test_element.tag),
				"s3_upload_time":    fmt.Sprintf("s3.upload.%s.time", test_element.tag),
				"s3_upload_total":   fmt.Sprintf("s3.upload.%s.total", test_element.tag),
				// Stat
				"s3_stat_failed":  fmt.Sprintf("s3.stat.%s.failed", test_element.tag),
				"s3_stat_succeed": fmt.Sprintf("s3.stat.%s.succeed", test_element.tag),
				"s3_stat_time":    fmt.Sprintf("s3.stat.%s.time", test_element.tag),
				"s3_stat_total":   fmt.Sprintf("s3.stat.%s.total", test_element.tag),
				// Download
				"s3_download_failed":  fmt.Sprintf("s3.download.%s.failed", test_element.tag),
				"s3_download_succeed": fmt.Sprintf("s3.download.%s.succeed", test_element.tag),
				"s3_download_time":    fmt.Sprintf("s3.download.%s.time", test_element.tag),
				"s3_download_total":   fmt.Sprintf("s3.download.%s.total", test_element.tag),
				// Delete
				"s3_delete_failed":  fmt.Sprintf("s3.delete.%s.failed", test_element.tag),
				"s3_delete_succeed": fmt.Sprintf("s3.delete.%s.succeed", test_element.tag),
				"s3_delete_time":    fmt.Sprintf("s3.delete.%s.time", test_element.tag),
				"s3_delete_total":   fmt.Sprintf("s3.delete.%s.total", test_element.tag),
			}
		}

		// Upload the file to S3
		timer = time.Now()
		if test_matrix.debug {
			log.Printf("Uploading [%s] to %s/%s",
				test_element.tmp_filename,
				test_matrix.agent_id,
				test_element.tag)
		}
		_, err := test_matrix.connection_object.minioClient.FPutObject(test_matrix.agent_id, test_element.tag,
			test_element.tmp_filename, "application/octet-stream")
		if err != nil {
			if test_matrix.debug {
				log.Printf("Unable to upload file [%s] to %s/%s",
					test_element.tmp_filename,
					test_matrix.agent_id,
					test_element.tag)
			}
			test_matrix.statsd_object.Increment(bucket_names["s3_upload_failed"],
				[]string{
					test_matrix.statsd_object.prefix,
					fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
					"s3_upload:failed",
					fmt.Sprintf("name:%s", test_element.tag),
					fmt.Sprintf("size:%d", test_element.file_size),
				}, 1)
		} else {
			test_matrix.statsd_object.Increment(bucket_names["s3_upload_succeed"],
				[]string{
					test_matrix.statsd_object.prefix,
					fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
					"s3_upload:succeed",
					fmt.Sprintf("name:%s", test_element.tag),
					fmt.Sprintf("size:%d", test_element.file_size),
				}, 1)
			if test_matrix.debug {
				log.Printf("Uploaded file [%s] to %s/%s in [%d]ms",
					test_element.tmp_filename,
					test_matrix.agent_id,
					test_element.tag,
					uint64(time.Since(timer)))
			}
			test_matrix.statsd_object.Timing(bucket_names["s3_upload_time"],
				time.Since(timer),
				[]string{
					test_matrix.statsd_object.prefix,
					fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
					"s3_upload:time",
					fmt.Sprintf("name:%s", test_element.tag),
					fmt.Sprintf("size:%d", test_element.file_size),
				}, 1)
		}
		test_matrix.statsd_object.Increment(bucket_names["s3_upload_total"],
			[]string{
				test_matrix.statsd_object.prefix,
				fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
				"s3_upload:total",
				fmt.Sprintf("name:%s", test_element.tag),
				fmt.Sprintf("size:%d", test_element.file_size),
			}, 1)

		// Stat the newly created S3 Object
		if test_matrix.debug {
			log.Printf("Getting Status of [%s] from %s/%s",
				test_element.tmp_filename,
				test_matrix.agent_id,
				test_element.tag)
		}
		timer = time.Now()
		objInfo, err := test_matrix.connection_object.minioClient.StatObject(test_matrix.agent_id,
			test_element.tag)
		if err != nil {
			if test_matrix.debug {
				log.Printf("Unable to stat file [%s] from bucket [%s]",
					test_element.tag,
					test_matrix.agent_id)
			}
			test_matrix.statsd_object.Increment(bucket_names["s3_stat_failed"],
				[]string{
					test_matrix.statsd_object.prefix,
					fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
					"s3_stat:failed",
					fmt.Sprintf("name:%s", test_element.tag),
					fmt.Sprintf("size:%d", test_element.file_size),
				}, 1)
		} else {
			test_matrix.statsd_object.Increment(bucket_names["s3_stat_succeed"],
				[]string{
					test_matrix.statsd_object.prefix,
					fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
					"s3_stat:succeed",
					fmt.Sprintf("name:%s", test_element.tag),
					fmt.Sprintf("size:%d", test_element.file_size),
				}, 1)
			if test_matrix.debug {
				log.Printf("Got status of object [%s] from bucket [%s] in [%d]ms",
					test_element.tag,
					test_matrix.agent_id, uint64(time.Since(timer)))
			}
			log.Println(objInfo)
			test_matrix.statsd_object.Timing(bucket_names["s3_stat_time"],
				time.Since(timer),
				[]string{
					test_matrix.statsd_object.prefix,
					fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
					"s3_stat:time",
					fmt.Sprintf("name:%s", test_element.tag),
					fmt.Sprintf("size:%d", test_element.file_size),
				}, 1)
		}
		test_matrix.statsd_object.Increment(bucket_names["s3_stat_total"],
			[]string{
				test_matrix.statsd_object.prefix,
				fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
				"s3_stat:total",
				fmt.Sprintf("name:%s", test_element.tag),
				fmt.Sprintf("size:%d", test_element.file_size),
			}, 1)

		// Download the file from S3
		if test_matrix.debug {
			log.Printf("Downloading [%s] from %s/%s",
				test_element.tmp_filename,
				test_matrix.agent_id,
				test_element.tag)
		}
		timer = time.Now()
		if test_matrix.connection_object.minioClient.FGetObject(test_matrix.agent_id,
			test_element.tag, test_element.tmp_filename) != nil {
			if test_matrix.debug {
				log.Printf("Unable to get object [%s] from bucket [%s]",
					test_element.tag,
					test_matrix.agent_id)
			}
			test_matrix.statsd_object.Increment(bucket_names["s3_download_failed"],
				[]string{
					test_matrix.statsd_object.prefix,
					fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
					"s3_download:failed",
					fmt.Sprintf("name:%s", test_element.tag),
					fmt.Sprintf("size:%d", test_element.file_size),
				}, 1)
		} else {
			test_matrix.statsd_object.Increment(bucket_names["s3_download_succeed"],
				[]string{
					test_matrix.statsd_object.prefix,
					fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
					"s3_download:succeed",
					fmt.Sprintf("name:%s", test_element.tag),
					fmt.Sprintf("size:%d", test_element.file_size),
				}, 1)
			if test_matrix.debug {
				log.Printf("Downloaded object [%s] to file [%s] in [%d]ms",
					test_element.tag,
					test_element.tmp_filename,
					uint64(time.Since(timer)))
			}
			test_matrix.statsd_object.Timing(bucket_names["s3_download_time"],
				time.Since(timer),
				[]string{
					test_matrix.statsd_object.prefix,
					fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
					"s3_download:time",
					fmt.Sprintf("name:%s", test_element.tag),
					fmt.Sprintf("size:%d", test_element.file_size),
				}, 1)
		}
		test_matrix.statsd_object.Increment(bucket_names["s3_download_total"],
			[]string{
				test_matrix.statsd_object.prefix,
				fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
				"s3_download:total",
				fmt.Sprintf("name:%s", test_element.tag),
				fmt.Sprintf("size:%d", test_element.file_size),
			}, 1)

		// Delete the file from S3
		if test_matrix.debug {
			log.Printf("Deleting [%s] from %s/%s",
				test_element.tmp_filename,
				test_matrix.agent_id,
				test_element.tag)
		}

		timer = time.Now()
		if test_matrix.connection_object.minioClient.RemoveObject(test_matrix.agent_id,
			test_element.tag) != nil {
			if test_matrix.debug {
				log.Printf("Unable to remove object [%s] from bucket [%s]",
					test_element.tag,
					test_matrix.agent_id)
			}
			test_matrix.statsd_object.Increment(bucket_names["s3_delete_failed"],
				[]string{
					test_matrix.statsd_object.prefix,
					fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
					"s3_delete:failed",
					fmt.Sprintf("name:%s", test_element.tag),
					fmt.Sprintf("size:%d", test_element.file_size),
				}, 1)
		} else {
			test_matrix.statsd_object.Increment(bucket_names["s3_delete_succeed"],
				[]string{
					test_matrix.statsd_object.prefix,
					fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
					"s3_delete:succeed",
					fmt.Sprintf("name:%s", test_element.tag),
					fmt.Sprintf("size:%d", test_element.file_size),
				}, 1)
			if test_matrix.debug {
				log.Printf("Deleted object [%s] in [%d]ms",
					test_element.tag,
					uint64(time.Since(timer)))
			}
			test_matrix.statsd_object.Timing(bucket_names["s3_delete_time"],
				time.Since(timer),
				[]string{
					test_matrix.statsd_object.prefix,
					fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
					"s3_delete:time",
					fmt.Sprintf("name:%s", test_element.tag),
					fmt.Sprintf("size:%d", test_element.file_size),
				}, 1)
		}
		test_matrix.statsd_object.Increment(bucket_names["s3_delete_total"],
			[]string{
				test_matrix.statsd_object.prefix,
				fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
				"s3_delete:total",
				fmt.Sprintf("name:%s", test_element.tag),
				fmt.Sprintf("size:%d", test_element.file_size),
			}, 1)
	}
	// Delete agent_id bucket
	timer = time.Now()
	if test_matrix.connection_object.minioClient.RemoveBucket(test_matrix.agent_id) != nil {
		if test_matrix.debug {
			log.Printf("Unable to remove bucket! [%s]", test_matrix.agent_id)
		}
		test_matrix.statsd_object.Increment("bucket.delete.failed",
			[]string{
				test_matrix.statsd_object.prefix,
				fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
				"bucket_delete:failed",
			}, 1)
	} else {
		test_matrix.statsd_object.Increment("bucket.delete.succeed",
			[]string{
				test_matrix.statsd_object.prefix,
				fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
				"bucket_delete:succeed",
			}, 1)
		if test_matrix.debug {
			log.Printf("Deleted bucket [%s] in [%d]ms", test_matrix.agent_id,
				uint64(time.Since(timer)))
		}
		test_matrix.statsd_object.Timing("bucket.delete.time", time.Since(timer),
			[]string{
				test_matrix.statsd_object.prefix,
				fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
				"bucket_delete:time",
			}, 1)
	}
	test_matrix.statsd_object.Increment("bucket.delete.total",
		[]string{
			test_matrix.statsd_object.prefix,
			fmt.Sprintf("agent_id:%s", test_matrix.agent_id),
			"bucket_delete:total",
		}, 1)

	//Cleanup
}
