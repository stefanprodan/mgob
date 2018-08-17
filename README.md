# mgob

[![Build Status](https://travis-ci.org/stefanprodan/mgob.svg?branch=master)](https://travis-ci.org/stefanprodan/mgob)
[![Docker Image](https://images.microbadger.com/badges/image/stefanprodan/mgob:edge.svg)](https://hub.docker.com/r/stefanprodan/mgob/)

MGOB is a MongoDB backup automation tool built with golang.

#### Features

* schedule backups
* local backups retention
* upload to S3 Object Storage (Minio, AWS, Google Cloud, Azure)
* upload to gcloud storage
* upload to SFTP
* notifications (Email, Slack)
* instrumentation with Prometheus
* http file server for local backups and logs
* distributed as an Alpine Docker image

#### Install

MGOB is available on Docker Hub at [stefanprodan/mgob](https://hub.docker.com/r/stefanprodan/mgob/). 

Supported tags:

* `stefanprodan/mgob:latest` latest stable [release](https://github.com/stefanprodan/mgob/releases)
* `stefanprodan/mgob:edge` master branch latest successful [build](https://travis-ci.org/stefanprodan/mgob)
* `stefanprodan/mgob:0.9` compatible with MongoDB 3.4
* `stefanprodan/mgob:0.10` compatible with MongoDB 3.6

Docker:

```bash
docker run -dp 8090:8090 --name mgob \
    -v "/mogb/config:/config" \
    -v "/mogb/storage:/storage" \
    -v "/mgob/tmp:/tmp" \
    -v "/mgob/data:/data" \
    stefanprodan/mgob \
    -LogLevel=info
```

Kubernetes:

A step by step guide on running MGOB as a StatefulSet with PersistentVolumeClaims can be found [here](https://github.com/stefanprodan/mgob/tree/master/k8s).

#### Configure

Define a backup plan (yaml format) for each database you want to backup inside the `config` dir. 
The yaml file name is being used as the backup plan ID, no white spaces or special characters are allowed. 

_Backup plan_

```yaml
scheduler:
  # run every day at 6:00 and 18:00 UTC
  cron: "0 6,18 */1 * *"
  # number of backups to keep locally
  retention: 14
  # backup operation timeout in minutes
  timeout: 60
target:
  # mongod IP or host name
  host: "172.18.7.21"
  # mongodb port
  port: 27017
  # mongodb database name, leave blank to backup all databases
  database: "test"
  # leave blank if auth is not enabled
  username: "admin"
  password: "secret"
  # add custom params to mongodump (eg. Auth or SSL support), leave blank if not needed
  params: "--ssl --authenticationDatabase admin"
# S3 upload (optional)
s3:
  url: "https://play.minio.io:9000"
  bucket: "backup"
  accessKey: "Q3AM3UQ867SPQQA43P2F"
  secretKey: "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
  # For Minio and AWS use S3v4 for GCP use S3v2 
  api: "S3v4"
# GCloud upload (optional)
gcloud:
  bucket: "backup"
  keyFilePath: /path/to/service-account.json
# Azure blob storage upload (optional)
azure:
  containerName: "backup"
  connectionString: "DefaultEndpointsProtocol=https;AccountName=...;AccountKey=...;EndpointSuffix=core.windows.net"
# SFTP upload (optional)
sftp:
  host: sftp.company.com
  port: 2022
  username: user
  password: secret
  # dir must exist on the SFTP server
  dir: backup
# Email notifications (optional)
smtp:
  server: smtp.company.com
  port: 465
  username: user
  password: secret
  from: mgob@company.com
  to:
    - devops@company.com
    - alerts@company.com
# Slack notifications (optional)
slack:
  url: https://hooks.slack.com/services/xxxx/xxx/xx
  channel: devops-alerts
  username: mgob
  # 'true' to notify only on failures 
  warnOnly: false
```

ReplicaSet example:

```yaml
target:
  host: "mongo-0.mongo.db,mongo-1.mongo.db,mongo-2.mongo.db"
  port: 27017
  database: "test"
```

Sharded cluster with authentication and SSL example:

```yaml
target:
  host: "mongos-0.db,mongos-1.db"
  port: 27017
  database: "test"
  username: "admin"
  password: "secret"
  params: "--ssl --authenticationDatabase admin"
```

#### Web API

* `mgob-host:8090/storage` file server
* `mgob-host:8090/status` backup jobs status
* `mgob-host:8090/metrics` Prometheus endpoint
* `mgob-host:8090/version` mgob version and runtime info
* `mgob-host:8090/debug` pprof endpoint

On demand backup:

* HTTP POST `mgob-host:8090/backup/:planID`

```bash
curl -X POST http://mgob-host:8090/backup/mongo-debug
```

```json
{
  "plan": "mongo-debug",
  "file": "mongo-debug-1494256295.gz",
  "duration": "3.635186255s",
  "size": "455 kB",
  "timestamp": "2017-05-08T15:11:35.940141701Z"
}
```

Scheduler status:

* HTTP GET `mgob-host:8090/status`
* HTTP GET `mgob-host:8090/status/:planID`

```bash
curl -X GET http://mgob-host:8090/status/mongo-debug
```

```json
{
  "plan": "mongo-debug",
  "next_run": "2017-05-13T14:32:00+03:00",
  "last_run": "2017-05-13T11:31:00.000622589Z",
  "last_run_status": "200",
  "last_run_log": "Backup finished in 2.339055539s archive mongo-debug-1494675060.gz size 527 kB"
}
```

#### Logs

View scheduler logs with `docker logs mgob`:

```bash
time="2017-05-05T16:50:55+03:00" level=info msg="Next run at 2017-05-05 16:51:00 +0300 EEST" plan=mongo-dev 
time="2017-05-05T16:50:55+03:00" level=info msg="Next run at 2017-05-05 16:52:00 +0300 EEST" plan=mongo-test 
time="2017-05-05T16:51:00+03:00" level=info msg="Backup started" plan=mongo-dev 
time="2017-05-05T16:51:02+03:00" level=info msg="Backup finished in 2.359901432s archive size 448 kB" plan=mongo-dev 
time="2017-05-05T16:52:00+03:00" level=info msg="Backup started" plan=mongo-test
time="2017-05-05T16:52:02+03:00" level=info msg="S3 upload finished `/storage/mongo-test/mongo-test-1493992320.gz` -> `bktest/mongo-test-1493992320.gz` Total: 1.17 KB, Transferred: 1.17 KB, Speed: 2.96 KB/s " plan=mongo-test 
time="2017-05-05T16:52:02+03:00" level=info msg="Backup finished in 2.855078717s archive size 1.2 kB" plan=mongo-test 
```

The success/fail logs will be sent via SMTP and/or Slack if notifications are enabled.

The mongodump log is stored along with the backup data (gzip archive) in the `storage` dir:

```bash
aleph-mbp:test aleph$ ls -lh storage/mongo-dev
total 4160
-rw-r--r--  1 aleph  staff   410K May  3 17:46 mongo-dev-1493822760.gz
-rw-r--r--  1 aleph  staff   1.9K May  3 17:46 mongo-dev-1493822760.log
-rw-r--r--  1 aleph  staff   410K May  3 17:47 mongo-dev-1493822820.gz
-rw-r--r--  1 aleph  staff   1.5K May  3 17:47 mongo-dev-1493822820.log
```

#### Metrics

Successful backups counter

```bash
mgob_scheduler_backup_total{plan="mongo-dev",status="200"} 8
```

Successful backups duration

```bash
mgob_scheduler_backup_latency{plan="mongo-dev",status="200",quantile="0.5"} 2.149668417
mgob_scheduler_backup_latency{plan="mongo-dev",status="200",quantile="0.9"} 2.39848413
mgob_scheduler_backup_latency{plan="mongo-dev",status="200",quantile="0.99"} 2.39848413
mgob_scheduler_backup_latency_sum{plan="mongo-dev",status="200"} 17.580484907
mgob_scheduler_backup_latency_count{plan="mongo-dev",status="200"} 8
```

Failed jobs count and duration (status 500)

```bash
mgob_scheduler_backup_latency{plan="mongo-test",status="500",quantile="0.5"} 2.4180213
mgob_scheduler_backup_latency{plan="mongo-test",status="500",quantile="0.9"} 2.438254775
mgob_scheduler_backup_latency{plan="mongo-test",status="500",quantile="0.99"} 2.438254775
mgob_scheduler_backup_latency_sum{plan="mongo-test",status="500"} 9.679809477
mgob_scheduler_backup_latency_count{plan="mongo-test",status="500"} 4
```

#### Restore

In order to restore from a local backup you have two options:

Browse `mgob-host:8090/storage` to identify the backup you want to restore. 
Login to your MongoDB server and download the archive using `curl` and restore the backup with `mongorestore` command line.

```bash
curl -o /tmp/mongo-test-1494056760.gz http://mgob-host:8090/storage/mongo-test/mongo-test-1494056760.gz
mongorestore --gzip --archive=/tmp/mongo-test-1494056760.gz --drop
```

You can also restore a backup from within mgob container. 
Exec into mgob, identify the backup you want to restore and use `mongorestore` to connect to your MongoDB server.

```bash
docker exec -it mgob sh
ls /storage/mongo-test
mongorestore --gzip --archive=/storage/mongo-test/mongo-test-1494056760.gz --host mongohost:27017 --drop
```
