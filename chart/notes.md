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

To install the chart, clone the Git repository first and from the repository's root directory, run:

    $ helm install --namespace my-kubernetes-ns ./chart

This will install the chart on your cluster.