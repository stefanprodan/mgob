package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/stefanprodan/mgob/api"
	"github.com/stefanprodan/mgob/config"
	_ "github.com/stefanprodan/mgob/mongodump"
	"os"
	"os/signal"
	"syscall"
	"github.com/stefanprodan/mgob/scheduler"
)

func main() {

	var appConfig = &config.AppConfig{}
	flag.StringVar(&appConfig.LogLevel, "LogLevel", "debug", "logging threshold level: debug|info|warn|error|fatal|panic")
	flag.IntVar(&appConfig.Port, "Port", 8090, "HTTP port to listen on")
	flag.StringVar(&appConfig.ConfigPath, "ConfigPath", "/Users/aleph/go/src/github.com/stefanprodan/mgob/test/config", "plan yml files dir")
	flag.StringVar(&appConfig.StoragePath, "StoragePath", "/Users/aleph/go/src/github.com/stefanprodan/mgob/test/storage", "backup storage")
	flag.StringVar(&appConfig.TmpPath, "TmpPath", "/Users/aleph/go/src/github.com/stefanprodan/mgob/test/tmp", "temporary backup storage")

	setLogLevel(appConfig.LogLevel)

	server := &api.HttpServer{
		Config: appConfig,
	}
	logrus.Infof("Starting HTTP server on port %v", appConfig.Port)
	go server.Start()

	plans, err := config.LoadPlans(appConfig.ConfigPath)

	if err != nil {
		logrus.Fatal(err)
	}

	//err = mongodump.Run(plans[0], appConfig)
	//if err != nil {
	//	logrus.Fatal(err)
	//}
	//logrus.Info("done")

	scheduler.Start(plans, appConfig)

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
