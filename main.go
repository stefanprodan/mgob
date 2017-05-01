package main

import (
	"github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"github.com/stefanprodan/mgob/config"
	"flag"
	"github.com/stefanprodan/mgob/api"
)

func main() {

	var config = &confing.AppConfig{}
	flag.StringVar(&config.LogLevel, "LogLevel", "debug", "logging threshold level: debug|info|warn|error|fatal|panic")
	flag.IntVar(&config.Port, "Port", 8090, "HTTP port to listen on")

	setLogLevel(config.LogLevel)

	server := &api.HttpServer{
		Config: config,
	}
	logrus.Infof("Starting HTTP server on port %v", config.Port)
	go server.Start()

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