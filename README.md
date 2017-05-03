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

Metrics:

Successful backups counter

```bash
# TYPE mgob_scheduler_backup_total counter
mgob_scheduler_backup_total{plan="mongo-dev",status="200"} 8
```

Successful backups duration

```bash
# HELP mgob_scheduler_backup_latency Backup duration in seconds.
# TYPE mgob_scheduler_backup_latency summary
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