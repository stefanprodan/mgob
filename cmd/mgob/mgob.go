package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/stefanprodan/mgob"
	"github.com/stefanprodan/mgob/config"
	"github.com/urfave/cli"
)

var (
	appConfig = &config.AppConfig{}
)

func beforeApp(c *cli.Context) error {
	level, err := log.ParseLevel(c.GlobalString("LogLevel"))
	if err != nil {
		log.Fatalf("unable to determine and set log level: %+v", err)
	}
	log.SetLevel(level)

	if c.GlobalBool("JSONLog") {
		// platforms such as Google StackDriver want logs to stdout
		log.SetOutput(os.Stdout)
		log.SetFormatter(&log.JSONFormatter{})
	}

	log.Debug("log level set to ", c.GlobalString("LogLevel"))
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "mgob"
	app.Version = mgob.VERSION
	app.Usage = "mongodb dockerized backup agent"
	app.Action = start
	app.Before = beforeApp
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "ConfigPath,c",
			Usage: "plan yml files dir",
			Value: "/config",
		},
		cli.StringFlag{
			Name:  "StoragePath,s",
			Usage: "backup storage",
			Value: "/storage",
		},
		cli.StringFlag{
			Name:  "TmpPath,t",
			Usage: "temporary backup storage",
			Value: "/tmp",
		},
		cli.StringFlag{
			Name:  "DataPath,d",
			Usage: "db dir",
			Value: "/data",
		},
		cli.IntFlag{
			Name:  "Port,p",
			Usage: "HTTP port to listen on",
			Value: 8090,
		},
		cli.BoolFlag{
			Name:  "JSONLog,j",
			Usage: "logs in JSON format",
		},
		cli.StringFlag{
			Name:  "LogLevel,l",
			Usage: "logging threshold level: debug|info|warn|error|fatal|panic",
			Value: "info",
		},
	}
	app.Run(os.Args)
}

func start(c *cli.Context) error {
	log.Infof("mgob %v", mgob.VERSION)

	appConfig.LogLevel = c.String("LogLevel")
	appConfig.JSONLog = c.Bool("JSONLog")
	appConfig.Port = c.Int("Port")
	appConfig.ConfigPath = c.String("ConfigPath")
	appConfig.StoragePath = c.String("StoragePath")
	appConfig.TmpPath = c.String("TmpPath")
	appConfig.DataPath = c.String("DataPath")

	mgob.Start(appConfig)

	return nil
}
