package main

import (
	"fmt"
	"log"
	"time"
)

func runTest(testMatrix *testMatrix) {

	// Register this agentId with StatsD
	testMatrix.statsdObject.Increment("agentId",
		[]string{
			testMatrix.statsdObject.prefix,
			fmt.Sprintf("agentId:%s", testMatrix.agentId),
		}, 1)

	// Create agentId bucket
	timer := time.Now()
	if testMatrix.connectionObject.minioClient.MakeBucket(testMatrix.agentId,
		testMatrix.connectionObject.s3Region) != nil {
		testMatrix.statsdObject.Increment("bucket.create.failed",
			[]string{
				testMatrix.statsdObject.prefix,
				fmt.Sprintf("agentId:%s", testMatrix.agentId),
				"bucket_create:failed",
			}, 1)
		if testMatrix.debug {
			log.Printf("Unable to create bucket! [%s]", testMatrix.agentId)
		}
	} else {
		testMatrix.statsdObject.Increment("bucket.create.succeed",
			[]string{
				testMatrix.statsdObject.prefix,
				fmt.Sprintf("agentId:%s", testMatrix.agentId),
				"bucket_create:succeed",
			}, 1)
		if testMatrix.debug {
			log.Printf("Created bucket [%s] in [%d]ms", testMatrix.agentId,
				uint64(time.Since(timer)))
		}
		testMatrix.statsdObject.Timing("bucket.create.time", time.Since(timer),
			[]string{
				testMatrix.statsdObject.prefix,
				fmt.Sprintf("agentId:%s", testMatrix.agentId),
				"bucket_create:time",
			}, 1)
	}
	testMatrix.statsdObject.Increment("bucket.create.total",
		[]string{
			testMatrix.statsdObject.prefix,
			fmt.Sprintf("agentId:%s", testMatrix.agentId),
			"bucket_create:total",
		}, 1)

	// iterate over the testElements
	for _, testElement := range testMatrix.testElements {
		// Calculate the bucket names
		bucketNames := make(map[string]string)
		if testMatrix.statsdObject.datadog {
			bucketNames = map[string]string{
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
			bucketNames = map[string]string{
				// Upload
				"s3_upload_failed":  fmt.Sprintf("s3.upload.%s.failed", testElement.tag),
				"s3_upload_succeed": fmt.Sprintf("s3.upload.%s.succeed", testElement.tag),
				"s3_upload_time":    fmt.Sprintf("s3.upload.%s.time", testElement.tag),
				"s3_upload_total":   fmt.Sprintf("s3.upload.%s.total", testElement.tag),
				// Stat
				"s3_stat_failed":  fmt.Sprintf("s3.stat.%s.failed", testElement.tag),
				"s3_stat_succeed": fmt.Sprintf("s3.stat.%s.succeed", testElement.tag),
				"s3_stat_time":    fmt.Sprintf("s3.stat.%s.time", testElement.tag),
				"s3_stat_total":   fmt.Sprintf("s3.stat.%s.total", testElement.tag),
				// Download
				"s3_download_failed":  fmt.Sprintf("s3.download.%s.failed", testElement.tag),
				"s3_download_succeed": fmt.Sprintf("s3.download.%s.succeed", testElement.tag),
				"s3_download_time":    fmt.Sprintf("s3.download.%s.time", testElement.tag),
				"s3_download_total":   fmt.Sprintf("s3.download.%s.total", testElement.tag),
				// Delete
				"s3_delete_failed":  fmt.Sprintf("s3.delete.%s.failed", testElement.tag),
				"s3_delete_succeed": fmt.Sprintf("s3.delete.%s.succeed", testElement.tag),
				"s3_delete_time":    fmt.Sprintf("s3.delete.%s.time", testElement.tag),
				"s3_delete_total":   fmt.Sprintf("s3.delete.%s.total", testElement.tag),
			}
		}

		// Upload the file to S3
		timer = time.Now()
		if testMatrix.debug {
			log.Printf("Uploading [%s] to %s/%s",
				testElement.tmpFilename,
				testMatrix.agentId,
				testElement.tag)
		}
		_, err := testMatrix.connectionObject.minioClient.FPutObject(testMatrix.agentId, testElement.tag,
			testElement.tmpFilename, "application/octet-stream")
		if err != nil {
			if testMatrix.debug {
				log.Printf("Unable to upload file [%s] to %s/%s",
					testElement.tmpFilename,
					testMatrix.agentId,
					testElement.tag)
			}
			testMatrix.statsdObject.Increment(bucketNames["s3_upload_failed"],
				[]string{
					testMatrix.statsdObject.prefix,
					fmt.Sprintf("agentId:%s", testMatrix.agentId),
					"s3_upload:failed",
					fmt.Sprintf("name:%s", testElement.tag),
					fmt.Sprintf("size:%d", testElement.fileSize),
				}, 1)
		} else {
			testMatrix.statsdObject.Increment(bucketNames["s3_upload_succeed"],
				[]string{
					testMatrix.statsdObject.prefix,
					fmt.Sprintf("agentId:%s", testMatrix.agentId),
					"s3_upload:succeed",
					fmt.Sprintf("name:%s", testElement.tag),
					fmt.Sprintf("size:%d", testElement.fileSize),
				}, 1)
			if testMatrix.debug {
				log.Printf("Uploaded file [%s] to %s/%s in [%d]ms",
					testElement.tmpFilename,
					testMatrix.agentId,
					testElement.tag,
					uint64(time.Since(timer)))
			}
			testMatrix.statsdObject.Timing(bucketNames["s3_upload_time"],
				time.Since(timer),
				[]string{
					testMatrix.statsdObject.prefix,
					fmt.Sprintf("agentId:%s", testMatrix.agentId),
					"s3_upload:time",
					fmt.Sprintf("name:%s", testElement.tag),
					fmt.Sprintf("size:%d", testElement.fileSize),
				}, 1)
		}
		testMatrix.statsdObject.Increment(bucketNames["s3_upload_total"],
			[]string{
				testMatrix.statsdObject.prefix,
				fmt.Sprintf("agentId:%s", testMatrix.agentId),
				"s3_upload:total",
				fmt.Sprintf("name:%s", testElement.tag),
				fmt.Sprintf("size:%d", testElement.fileSize),
			}, 1)

		// Stat the newly created S3 Object
		if testMatrix.debug {
			log.Printf("Getting Status of [%s] from %s/%s",
				testElement.tmpFilename,
				testMatrix.agentId,
				testElement.tag)
		}
		timer = time.Now()
		objInfo, err := testMatrix.connectionObject.minioClient.StatObject(testMatrix.agentId,
			testElement.tag)
		if err != nil {
			if testMatrix.debug {
				log.Printf("Unable to stat file [%s] from bucket [%s]",
					testElement.tag,
					testMatrix.agentId)
			}
			testMatrix.statsdObject.Increment(bucketNames["s3_stat_failed"],
				[]string{
					testMatrix.statsdObject.prefix,
					fmt.Sprintf("agentId:%s", testMatrix.agentId),
					"s3_stat:failed",
					fmt.Sprintf("name:%s", testElement.tag),
					fmt.Sprintf("size:%d", testElement.fileSize),
				}, 1)
		} else {
			testMatrix.statsdObject.Increment(bucketNames["s3_stat_succeed"],
				[]string{
					testMatrix.statsdObject.prefix,
					fmt.Sprintf("agentId:%s", testMatrix.agentId),
					"s3_stat:succeed",
					fmt.Sprintf("name:%s", testElement.tag),
					fmt.Sprintf("size:%d", testElement.fileSize),
				}, 1)
			if testMatrix.debug {
				log.Printf("Got status of object [%s] from bucket [%s] in [%d]ms",
					testElement.tag,
					testMatrix.agentId, uint64(time.Since(timer)))
			}
			log.Println(objInfo)
			testMatrix.statsdObject.Timing(bucketNames["s3_stat_time"],
				time.Since(timer),
				[]string{
					testMatrix.statsdObject.prefix,
					fmt.Sprintf("agentId:%s", testMatrix.agentId),
					"s3_stat:time",
					fmt.Sprintf("name:%s", testElement.tag),
					fmt.Sprintf("size:%d", testElement.fileSize),
				}, 1)
		}
		testMatrix.statsdObject.Increment(bucketNames["s3_stat_total"],
			[]string{
				testMatrix.statsdObject.prefix,
				fmt.Sprintf("agentId:%s", testMatrix.agentId),
				"s3_stat:total",
				fmt.Sprintf("name:%s", testElement.tag),
				fmt.Sprintf("size:%d", testElement.fileSize),
			}, 1)

		// Download the file from S3
		if testMatrix.debug {
			log.Printf("Downloading [%s] from %s/%s",
				testElement.tmpFilename,
				testMatrix.agentId,
				testElement.tag)
		}
		timer = time.Now()
		if testMatrix.connectionObject.minioClient.FGetObject(testMatrix.agentId,
			testElement.tag, testElement.tmpFilename) != nil {
			if testMatrix.debug {
				log.Printf("Unable to get object [%s] from bucket [%s]",
					testElement.tag,
					testMatrix.agentId)
			}
			testMatrix.statsdObject.Increment(bucketNames["s3_download_failed"],
				[]string{
					testMatrix.statsdObject.prefix,
					fmt.Sprintf("agentId:%s", testMatrix.agentId),
					"s3_download:failed",
					fmt.Sprintf("name:%s", testElement.tag),
					fmt.Sprintf("size:%d", testElement.fileSize),
				}, 1)
		} else {
			testMatrix.statsdObject.Increment(bucketNames["s3_download_succeed"],
				[]string{
					testMatrix.statsdObject.prefix,
					fmt.Sprintf("agentId:%s", testMatrix.agentId),
					"s3_download:succeed",
					fmt.Sprintf("name:%s", testElement.tag),
					fmt.Sprintf("size:%d", testElement.fileSize),
				}, 1)
			if testMatrix.debug {
				log.Printf("Downloaded object [%s] to file [%s] in [%d]ms",
					testElement.tag,
					testElement.tmpFilename,
					uint64(time.Since(timer)))
			}
			testMatrix.statsdObject.Timing(bucketNames["s3_download_time"],
				time.Since(timer),
				[]string{
					testMatrix.statsdObject.prefix,
					fmt.Sprintf("agentId:%s", testMatrix.agentId),
					"s3_download:time",
					fmt.Sprintf("name:%s", testElement.tag),
					fmt.Sprintf("size:%d", testElement.fileSize),
				}, 1)
		}
		testMatrix.statsdObject.Increment(bucketNames["s3_download_total"],
			[]string{
				testMatrix.statsdObject.prefix,
				fmt.Sprintf("agentId:%s", testMatrix.agentId),
				"s3_download:total",
				fmt.Sprintf("name:%s", testElement.tag),
				fmt.Sprintf("size:%d", testElement.fileSize),
			}, 1)

		// Delete the file from S3
		if testMatrix.debug {
			log.Printf("Deleting [%s] from %s/%s",
				testElement.tmpFilename,
				testMatrix.agentId,
				testElement.tag)
		}

		timer = time.Now()
		if testMatrix.connectionObject.minioClient.RemoveObject(testMatrix.agentId,
			testElement.tag) != nil {
			if testMatrix.debug {
				log.Printf("Unable to remove object [%s] from bucket [%s]",
					testElement.tag,
					testMatrix.agentId)
			}
			testMatrix.statsdObject.Increment(bucketNames["s3_delete_failed"],
				[]string{
					testMatrix.statsdObject.prefix,
					fmt.Sprintf("agentId:%s", testMatrix.agentId),
					"s3_delete:failed",
					fmt.Sprintf("name:%s", testElement.tag),
					fmt.Sprintf("size:%d", testElement.fileSize),
				}, 1)
		} else {
			testMatrix.statsdObject.Increment(bucketNames["s3_delete_succeed"],
				[]string{
					testMatrix.statsdObject.prefix,
					fmt.Sprintf("agentId:%s", testMatrix.agentId),
					"s3_delete:succeed",
					fmt.Sprintf("name:%s", testElement.tag),
					fmt.Sprintf("size:%d", testElement.fileSize),
				}, 1)
			if testMatrix.debug {
				log.Printf("Deleted object [%s] in [%d]ms",
					testElement.tag,
					uint64(time.Since(timer)))
			}
			testMatrix.statsdObject.Timing(bucketNames["s3_delete_time"],
				time.Since(timer),
				[]string{
					testMatrix.statsdObject.prefix,
					fmt.Sprintf("agentId:%s", testMatrix.agentId),
					"s3_delete:time",
					fmt.Sprintf("name:%s", testElement.tag),
					fmt.Sprintf("size:%d", testElement.fileSize),
				}, 1)
		}
		testMatrix.statsdObject.Increment(bucketNames["s3_delete_total"],
			[]string{
				testMatrix.statsdObject.prefix,
				fmt.Sprintf("agentId:%s", testMatrix.agentId),
				"s3_delete:total",
				fmt.Sprintf("name:%s", testElement.tag),
				fmt.Sprintf("size:%d", testElement.fileSize),
			}, 1)
	}
	// Delete agentId bucket
	timer = time.Now()
	if testMatrix.connectionObject.minioClient.RemoveBucket(testMatrix.agentId) != nil {
		if testMatrix.debug {
			log.Printf("Unable to remove bucket! [%s]", testMatrix.agentId)
		}
		testMatrix.statsdObject.Increment("bucket.delete.failed",
			[]string{
				testMatrix.statsdObject.prefix,
				fmt.Sprintf("agentId:%s", testMatrix.agentId),
				"bucket_delete:failed",
			}, 1)
	} else {
		testMatrix.statsdObject.Increment("bucket.delete.succeed",
			[]string{
				testMatrix.statsdObject.prefix,
				fmt.Sprintf("agentId:%s", testMatrix.agentId),
				"bucket_delete:succeed",
			}, 1)
		if testMatrix.debug {
			log.Printf("Deleted bucket [%s] in [%d]ms", testMatrix.agentId,
				uint64(time.Since(timer)))
		}
		testMatrix.statsdObject.Timing("bucket.delete.time", time.Since(timer),
			[]string{
				testMatrix.statsdObject.prefix,
				fmt.Sprintf("agentId:%s", testMatrix.agentId),
				"bucket_delete:time",
			}, 1)
	}
	testMatrix.statsdObject.Increment("bucket.delete.total",
		[]string{
			testMatrix.statsdObject.prefix,
			fmt.Sprintf("agentId:%s", testMatrix.agentId),
			"bucket_delete:total",
		}, 1)

	//Cleanup
}
