# Kubernetes MongoDB Backup Operator

This is a step by step guide on setting up 
MGOB to automate MongoDB backups on Google Kubernetes Engine.

Requirements:

* GKE cluster minimum version v1.8
* kubctl admin config

#### Create a MongoDB ReplicaSet

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

Create a 3 nodes Replica Set, each replica provisioned with a 1Gi SSD disk:

```bash
$ kubectl apply -f ./mongo-rs.yaml 
service "mongo" created
statefulset "mongo" created
```

The above command creates a Headless Service and a Stateful Set for the Mongo Replica Set and a Service Account for the Mongo sidecar.
Each POD contains a MongoDB instance and a [mongo-k8s-sidecar](https://github.com/cvallance/mongo-k8s-sidecar), 
the sidecar will initialize the Replica Set and will add the rs members as soon as the PODs are up.

You can monitor the rs initialization by looking at the sidecar logs:

```bash
$ kubectl -n db logs mongo-0 mongo-sidecar
Using mongo port: 27017
Starting up mongo-k8s-sidecar
The cluster domain 'cluster.local' was successfully verified.
Pod has been elected for replica set initialization
initReplSet 10.52.2.127:27017
```

Inspect the newly created cluster with kubectl:

```bash
$ kubectl -n db get all
NAME                 DESIRED   CURRENT   AGE
statefulsets/mongo   3         3         8m

NAME         READY     STATUS    RESTARTS   AGE
po/mongo-0   2/2       Running   0          8m
po/mongo-1   2/2       Running   0          7m
po/mongo-2   2/2       Running   0          6m

NAME        TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)     AGE
svc/mongo   ClusterIP   None         <none>        27017/TCP   11m
```

You can run a temporary mongo pod inside the `db` namespace to create a `test` database and insert some data:

```bash
$ kubectl -n db run -it --rm --restart=Never mongo-cli --image=mongo --command -- sh
mongo "mongodb://mongo-0.mongo,mongo-1.mongo,mongo-2.mongo:27017/dbname_?"
rs0:PRIMARY> use test
rs0:PRIMARY> db.inventory.insert({item: "one", val: "test" })
WriteResult({ "nInserted" : 1 })
```


