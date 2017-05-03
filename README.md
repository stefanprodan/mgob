# mgob

[![Build Status](https://travis-ci.org/stefanprodan/mgob.svg?branch=master)](https://travis-ci.org/stefanprodan/mgob)

MGOB is a backup manager for MongoDB.

Features:
* schedule backups
* local backups retention
* upload to S3 Object Storage (Minio, AWS, Google Cloud)
* instrumentation with Prometheus
* Alpine Docker image

Install:

```bash
docker run -dp 8090:8090 --name mgob \
    -v "/mgo/config:/config" \
    -v "/mgo/storage:/storage" \
    -v "/mgo/tmp:/tmp" \
    stefanprodan/mgo \
    -LogLevel=info
```

Configure:

At startup MGOB loads the backup plans from the `config` volume.

***Local backup plan***

```yaml
target:
  host: "172.18.7.21"
  port: 27017
  database: "test" 
scheduler:
  cron: "*/1 * * * *"
  # number of backups to keep
  retention: 7
  # backup operation timeout in seconds
  timeout: 60
```

***Local backup with S3 upload plan***

```yaml
target:
  host: "172.18.7.21"
  port: 27017
  database: "test" 
  username: "admin"
  password: "secret"
scheduler:
  cron: "0 6,18 */1 * *"
  retention: 5
  timeout: 30
s3:
  url: "https://play.minio.io:9000"
  bucket: "backup"
  accessKey: "Q3AM3UQ867SPQQA43P2F"
  secretKey: "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
  api: "S3v4"
```