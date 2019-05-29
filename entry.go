package mgob

import (
	"os"
	"os/signal"
	"path"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/stefanprodan/mgob/api"
	"github.com/stefanprodan/mgob/backup"
	"github.com/stefanprodan/mgob/config"
	"github.com/stefanprodan/mgob/db"
	"github.com/stefanprodan/mgob/scheduler"
)

var appConfig = &config.AppConfig{}

func Start(appConfig *config.AppConfig) {
	log.Infof("starting with config: %+v", appConfig)

	info, err := backup.CheckMongodump()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(info)

	info, err = backup.CheckMinioClient()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(info)

	info, err = backup.CheckGCloudClient()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(info)

	info, err = backup.CheckAzureClient()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(info)

	plans, err := config.LoadPlans(appConfig.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	store, err := db.Open(path.Join(appConfig.DataPath, "mgob.db"))
	if err != nil {
		log.Fatal(err)
	}
	statusStore, err := db.NewStatusStore(store)
	if err != nil {
		log.Fatal(err)
	}
	sch := scheduler.New(plans, appConfig, statusStore)
	sch.Start()

	server := &api.HttpServer{
		Config: appConfig,
		Stats:  statusStore,
	}
	log.Infof("starting http server on port %v", appConfig.Port)
	go server.Start(VERSION)

	// wait for SIGINT (Ctrl+C) or SIGTERM (docker stop)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan

	log.Infof("shutting down %v signal received", sig)
}
