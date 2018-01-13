# Kubernetes MongoDB Backup Operator

This is a step by step guide on setting up 
MGOB to automate MongoDB backups on Google Kubernetes Engine.

Requirements:

* GKE cluster minimum version v1.8
* kubctl admin config

Clone the mgob repository:

```bash
$ git clone https://github.com/stefanprodan/mgob.git
$ cd mgob/k8s
```

Create a cluster admin user:

```bash
kubectl create clusterrolebinding "cluster-admin-$(whoami)" \
    --clusterrole=cluster-admin \
    --user="$(gcloud config get-value core/account)"
```

### Create a MongoDB RS with Stateful Sets

Create the `db` namespace:

```bash
$ kubectl apply -f ./namespace.yaml 
namespace "db" created
```

Create the `ssd` and `hdd` storage classes:

```bash
$ kubectl apply -f ./storage.yaml 
storageclass "ssd" created
storageclass "hdd" created
```

Create the `startup-script` _Daemon Set_ to disable hugepage on all hosts:

```bash
$ kubectl apply -f ./mongo-ds.yaml 
daemonset "startup-script" created
```

Create a 3 nodes _Replica Set_, each replica provisioned with a 1Gi SSD disk:

```bash
$ kubectl apply -f ./mongo-rs.yaml 
service "mongo" created
statefulset "mongo" created
clusterrole "default" configured
serviceaccount "default" configured
clusterrolebinding "system:serviceaccount:db:default" configured
```

The above command creates a _Headless Service_ and a _Stateful Set_ for the Mongo _Replica Set_ and a _Service Account_ for the Mongo sidecar.
Each pod contains a Mongo instance and a sidecar.
The sidecar will initialize the _Replica Set_ and will add the rs members as soon as the pods are up.
You can safely scale up or down the _Stateful Set_ replicas, the sidecar will add or remove rs members.

You can monitor the rs initialization by looking at the sidecar logs:

```bash
$ kubectl -n db logs mongo-0 mongo-sidecar
Using mongo port: 27017
Starting up mongo-k8s-sidecar
The cluster domain 'cluster.local' was successfully verified.
Pod has been elected for replica set initialization
initReplSet 10.52.2.127:27017
```

Inspect the newly created cluster with `kubectl`:

```bash
$ kubectl -n db get pods --selector=role=mongo
NAME         READY     STATUS    RESTARTS   AGE
po/mongo-0   2/2       Running   0          8m
po/mongo-1   2/2       Running   0          7m
po/mongo-2   2/2       Running   0          6m
```

Connect to the container running in `mongo-0` pod, create a `test` database and insert some data:

```bash
$ kubectl -n db exec -it mongo-0 -c mongod mongo
rs0:PRIMARY> use test
rs0:PRIMARY> db.inventory.insert({item: "one", val: "two" })
WriteResult({ "nInserted" : 1 })
```

Each MongoDB replica has its own DNS address as in `<pod-name>.<service-name>.<namespace>`.
If you need to access the _Replica Set_ from another namespace use the following connection url:

```
mongodb://mongo-0.mongo.db,mongo-1.mongo.db,mongo-2.mongo.db:27017/dbname_?
```

Test the connectivity by creating a temporary pod in the default namespace:

```
$ kubectl run -it --rm --restart=Never mongo-cli --image=mongo --command -- /bin/bash
root@mongo-cli:/# mongo "mongodb://mongo-0.mongo.db,mongo-1.mongo.db,mongo-2.mongo.db:27017/test"
rs0:PRIMARY> db.getCollectionNames()
[ "inventory" ]
```

The [mongo-k8s-sidecar](https://github.com/cvallance/mongo-k8s-sidecar) deals with ReplicaSet provisioning only. 
if you want to run a sharded cluster on GKE, take a look at [pkdone/gke-mongodb-shards-demo](https://github.com/pkdone/gke-mongodb-shards-demo). 

### Create a MongoDB Backup agent with Stateful Sets

First let's create two databases `test1` and `test2`:

```bash
$ kubectl -n db exec -it mongo-0 -c mongod mongo
rs0:PRIMARY> use test1
rs0:PRIMARY> db.inventory.insert({item: "one", val: "two" })
WriteResult({ "nInserted" : 1 })
rs0:PRIMARY> use test2
rs0:PRIMARY> db.inventory.insert({item: "one", val: "two" })
WriteResult({ "nInserted" : 1 })
```

Create a ConfigMap to schedule backups every minute for `test1` and every two minutes for `test2`:

```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  labels:
    role: backup
  name: mgob-config
  namespace: db
data:
  test1.yml: |
    target:
      host: "mongo-0.mongo.db,mongo-1.mongo.db,mongo-2.mongo.db"
      port: 27017
      database: "test1"
    scheduler:
      cron: "*/1 * * * *"
      retention: 5
      timeout: 60
  test2.yml: |
    target:
      host: "mongo-0.mongo.db,mongo-1.mongo.db,mongo-2.mongo.db"
      port: 27017
      database: "test2"
    scheduler:
      cron: "*/2 * * * *"
      retention: 10
      timeout: 60
```

Apply the config:

```bash
kubectl apply -f ./mgob-cfg.yaml
```

Deploy mgob _Headless Service_ and _Stateful Set_ with two disks, 3Gi for the long term backup storage 
and 1Gi for the temporary storage of the running backups:

```bash
kubectl apply -f ./mgob-dep.yaml
```

To monitor the backups you can stream the mgob logs:

```bash
$ kubectl -n db logs -f mgob-0 
msg="Backup started" plan=test1 
msg="Backup finished in 261.76829ms archive test1-1514491560.gz size 307 B" plan=test1 
msg="Next run at 2017-12-28 20:07:00 +0000 UTC" plan=test1 
msg="Backup started" plan=test2
msg="Backup finished in 266.635088ms archive test2-1514491560.gz size 313 B" plan=test2 
msg="Next run at 2017-12-28 20:08:00 +0000 UTC" plan=test2 
```

Or you can `curl` the mgob API:

```bash
kubectl -n db exec -it mgob-0 -- curl mgob-0.mgob.db:8090/status
```

Let's run an on demand backup for `test2` database:

```bash
kubectl -n db exec -it mgob-0 -- curl -XPOST mgob-0.mgob.db:8090/backup/test2
{"plan":"test2","file":"test2-1514492080.gz","duration":"61.109042ms","size":"313 B","timestamp":"2017-12-28T20:14:40.604057546Z"}
```

You can restore a backup from within mgob container. 
Exec into mgob and identify the backup you want to restore, the backups are in `/storage/<plan-name>`.

```bash
$ kubectl -n db exec -it mgob-0 /bin/bash
ls -lh /storage/test1
-rw-r--r--    1 root     root         307 Dec 28 20:23 test1-1514492580.gz
-rw-r--r--    1 root     root         162 Dec 28 20:23 test1-1514492580.log
-rw-r--r--    1 root     root         307 Dec 28 20:24 test1-1514492640.gz
-rw-r--r--    1 root     root         162 Dec 28 20:24 test1-1514492640.log
```

Use `mongorestore` to connect to your MongoDB server and restore a backup:

```bash
$ kubectl -n db exec -it mgob-0 /bin/bash
mongorestore --gzip --archive=/storage/test1/test1-1514492640.gz --host mongo-0.mongo.db:27017 --drop
```

### Monitoring and alerting

For each backup plan you can configure alerting via email or Slack:

```yaml
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

Mgob exposes Prometheus metrics on the `/metrics` endpoint. 

Successful/failed backups counter:

```
mgob_scheduler_backup_total{plan="test1",status="200"} 8
mgob_scheduler_backup_total{plan="test2",status="500"} 2
```

Backup duration:

```
mgob_scheduler_backup_latency{plan="test1",status="200",quantile="0.5"} 2.149668417
mgob_scheduler_backup_latency{plan="test1",status="200",quantile="0.9"} 2.39848413
mgob_scheduler_backup_latency{plan="test1",status="200",quantile="0.99"} 2.39848413
```

### Backup to GCP Storage Bucket

For long term backup storage you could use a GCP Bucket since is a cheaper option than keeping all 
backups on disk.

First you need to create an GCP service account key from the `API & Services` page. Download the JSON file 
and rename it to `service-account.json`. 

Store the JSON file as a secret in the `db` namespace:

```bash
kubectl -n db create secret generic gcp-key --from-file=service-account.json=service-account.json
```

From the GCP web UI, navigate to _Storage_ and create a regional bucket named `mgob`. 
If the bucket name is taken you'll need to change it in the `mgob-gstore-cfg.yaml` file:

```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  labels:
    role: mongo-backup
  name: mgob-gstore-config
  namespace: db
data:
  test.yml: |
    target:
      host: "mongo-0.mongo.db,mongo-1.mongo.db,mongo-2.mongo.db"
      port: 27017
      database: "test"
    scheduler:
      cron: "*/1 * * * *"
      retention: 1
      timeout: 60
    gcloud:
      bucket: "mgob"
      keyFilePath: /etc/mgob/service-account.json
```

Apply the config:

```bash
kubectl apply -f ./mgob-gstore-cfg.yaml
```

Deploy mgob with the `gcp-key` secret map to a volume:

```bash
kubectl apply -f ./mgob-gstore-dep.yaml
```

After one minute the backup will be uploaded to the GCP bucket:

```bash
$ kubectl -n db logs -f mgob-0 
msg="Google Cloud SDK 181.0.0 bq 2.0.27 core 2017.11.28 gsutil 4.28"
msg="Backup started" plan=test
msg="GCloud upload finished Copying file:///storage/test/test-1514544660.gz"
```

