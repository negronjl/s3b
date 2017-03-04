# s3b
```
NAME:
   s3b - S3/Object Store benchmarking tool

USAGE:
   s3b [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR:
   Juan L. Negron <negronjl@xtremeghost.com>

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --server value, -s value         Target S3 server to test [$S3_SERVER]
   --region value, -r value         S3 server region to use (default: "us-east-1") [$S3_REGION]
   --access-key value, -A value     Access Key to the S3 server [$S3_ACCESS_KEY]
   --secret-key value, -R value     Access Key to the S3 server [$S3_SECRET_KEY]
   --api-signature value, -a value  API Signature version (v2 or v4) (default: "v4") [$S3_API_SIGNATURE]
   --SSL, -S                        Whether or not to use SSL to connect to the S3 server [$S3_SSL]
   --debug, -d                      Print debug HTTP tracing information [$S3_DEBUG]
   --statsd value, -D value         StatsD server to which metrics will be sent (default: "s3b") [$S3_STATSD_HOST]
   --prefix value, -p value         Prefix to use with the StatsD metrics [$S3_STATSD_PREFIX]
   --matrix value, -m value         Comma separated key value pairs of filename=size to use in the testing. [$S3_TEST_MATRIX]
   --matrix-dir value, -M value     Directory containing the files to be used for testing. [$S3_TEST_MATRIX_DIR]
   --help, -h                       show help
   --version, -v                    print the version

```