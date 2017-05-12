package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/stefanprodan/mgob/api"
	"github.com/stefanprodan/mgob/backup"
	"github.com/stefanprodan/mgob/config"
	"github.com/stefanprodan/mgob/scheduler"
	"os"
	"os/signal"
	"syscall"
)

var version = "undefined"

func main() {
	var appConfig = &config.AppConfig{}
	flag.StringVar(&appConfig.LogLevel, "LogLevel", "debug", "logging threshold level: debug|info|warn|error|fatal|panic")
	flag.IntVar(&appConfig.Port, "Port", 8090, "HTTP port to listen on")
	flag.StringVar(&appConfig.ConfigPath, "ConfigPath", "/config", "plan yml files dir")
	flag.StringVar(&appConfig.StoragePath, "StoragePath", "/storage", "backup storage")
	flag.StringVar(&appConfig.TmpPath, "TmpPath", "/tmp", "temporary backup storage")
	flag.StringVar(&appConfig.DataPath, "DataPath", "/data", "db dir")
	flag.Parse()
	setLogLevel(appConfig.LogLevel)
	logrus.Infof("Starting with config: %+v", appConfig)

	info, err := backup.CheckMongodump()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(info)

	info, err = backup.CheckMinioClient()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(info)

	plans, err := config.LoadPlans(appConfig.ConfigPath)
	if err != nil {
		logrus.Fatal(err)
	}

	stats := scheduler.NewStats(plans)
	sch := scheduler.New(plans, appConfig, stats)
	sch.Start()

	server := &api.HttpServer{
		Config: appConfig,
		Stats:  stats,
	}
	logrus.Infof("Starting HTTP server on port %v", appConfig.Port)
	go server.Start(version)

	//wait for SIGINT (Ctrl+C) or SIGTERM (docker stop)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan

	logrus.Infof("Shuting down %v signal received", sig)
}

func setLogLevel(levelName string) {
	level, err := logrus.ParseLevel(levelName)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetLevel(level)
}
