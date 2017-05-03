# mgob

[![Build Status](https://travis-ci.org/stefanprodan/mgob.svg?branch=master)](https://travis-ci.org/stefanprodan/mgob)

MGOB is a backup manager for MongoDB.

Features:

* schedule backups
* local backups retention
* upload to S3 Object Storage (Minio, AWS, Google Cloud)
* instrumentation with Prometheus
* http file server for local backups and logs
* Alpine Docker image

Install:

```bash
docker run -dp 8090:8090 --name mgob \
    -v "/mgo/config:/config" \
    -v "/mgo/storage:/storage" \
    -v "/mgo/tmp:/tmp" \
    stefanprodan/mgob \
    -LogLevel=info
```

Configure:

At startup MGOB loads the backup plans from the `config` volume.

_Local backup plan_

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

_Local backup with S3 upload plan_

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

Web API:

* `mgob-host:8090/` file server
* `mgob-host:8090/status` backup jobs status
* `mgob-host:8090/metrics` Prometheus endpoint

Logs:

View scheduler logs with `docker logs mgob`:

```bash
time="2017-05-03T17:44:14+03:00" level=info msg="Plan mongo-dev next run at 2017-05-03 17:45:00 +0300 EEST" 
time="2017-05-03T17:44:14+03:00" level=info msg="Plan mongo-test next run at 2017-05-03 17:46:00 +0300 EEST" 
time="2017-05-03T17:45:00+03:00" level=info msg="mongo-dev backup started" 
time="2017-05-03T17:45:02+03:00" level=info msg="mongo-dev backup finished in 2.313856701s archive size 420 kB" 
time="2017-05-03T17:46:00+03:00" level=info msg="mongo-test backup started" 
time="2017-05-03T17:46:03+03:00" level=info msg="S3 upload finshied `/Users/aleph/go/src/github.com/stefanprodan/mgob/test/storage/mongo-test/mongo-test-1493822760.gz` -> `mongo-test/bktest/mongo-test-1493822760.gz` Total: 1.17 KB, Transferred: 1.17 KB, Speed: 3.00 KB/s " 
time="2017-05-03T17:46:03+03:00" level=info msg="mongo-test backup finished in 3.505959481s archive size 1.2 kB" 
```

The mongodump log is stored along with the backup data (gzip archive) in the `storage` dir:

```bash
aleph-mbp:test aleph$ ls -lh storage/mongo-dev
total 4160
-rw-r--r--  1 aleph  staff   410K May  3 17:46 mongo-dev-1493822760.gz
-rw-r--r--  1 aleph  staff   1.9K May  3 17:46 mongo-dev-1493822760.log
-rw-r--r--  1 aleph  staff   410K May  3 17:47 mongo-dev-1493822820.gz
-rw-r--r--  1 aleph  staff   1.5K May  3 17:47 mongo-dev-1493822820.log
```

Metrics:

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