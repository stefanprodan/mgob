package main

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/Sirupsen/logrus"
)

func main() {

	//wait for SIGINT (Ctrl+C) or SIGTERM (docker stop)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan

	logrus.Infof("Shuting down %v signal received", sig)
}
