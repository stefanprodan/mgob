Usage
=====
This script assumes a user has been created in the mongoDB instance with sufficient read privileges to create the
backups.

The password set when creating the user below is referred to in the `values.yaml`-file for each plan. Make sure that
those match, or your backups will fail.

    db.createUser({
      user: "mongodb-backup",
      pwd: "backup-user-pwd",
      roles: [
        { role: "backup", db: "admin" }
      ]
    });
    
Secondly, it assumes that [Helm](https://github.com/kubernetes/helm) is installed on your cluster and that you have
a properly configured [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed.

To install the chart, first clone the Git repository. Secondly, edit the `values.yaml`-file to define your backup
plans. Documentation for those can be found in the default repository `readme.md`, easily accessible as the
[mgob repository start page](https://github.com/stefanprodan/mgob). When the `values.yaml`-file properly represents the 
plan(s) you want to create, simply run:

    $ helm install --namespace my-kubernetes-ns ./chart

This will install the chart on your cluster.