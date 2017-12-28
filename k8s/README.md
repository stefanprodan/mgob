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

### Create a MongoDB Replica Set with Stateful Sets

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
Each POD contains a MongoDB instance and a [mongo-k8s-sidecar](https://github.com/cvallance/mongo-k8s-sidecar), 
the sidecar will initialize the _Replica Set_ and will add the rs members as soon as the PODs are up.

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

You can run a temporary _mongo-cli_ pod inside the `db` namespace to create a `test` database and insert some data:

```bash
$ kubectl -n db run -it --rm --restart=Never mongo-cli --image=mongo --command -- sh
mongo "mongodb://mongo-0.mongo,mongo-1.mongo,mongo-2.mongo:27017/dbname_?"
rs0:PRIMARY> use test
rs0:PRIMARY> db.inventory.insert({item: "one", val: "two" })
WriteResult({ "nInserted" : 1 })
```

Each MongoDB replica has its own DNS address as in `<pod-name>.<service-name>.<namespace>`.
If you need to access the _Replica Set_ from another namespace use the following connection url:

```bash
mongodb://mongo-0.mongo.db,mongo-1.mongo.db,mongo-2.mongo.db:27017/dbname_?
```
