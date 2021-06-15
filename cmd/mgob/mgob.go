package main

import (
	"os"
	"os/signal"
	"path"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/stefanprodan/mgob/pkg/api"
	"github.com/stefanprodan/mgob/pkg/backup"
	"github.com/stefanprodan/mgob/pkg/config"
	"github.com/stefanprodan/mgob/pkg/db"
	"github.com/stefanprodan/mgob/pkg/scheduler"
)

var (
	appConfig = &config.AppConfig{}
	version   = "v1.3.0-dev"
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
	app.Version = version
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
		cli.StringFlag{
			Name:  "Host,h",
			Usage: "HTTP host to listen on",
			Value: "",
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
	log.Infof("mgob %v", version)

	appConfig.LogLevel = c.String("LogLevel")
	appConfig.JSONLog = c.Bool("JSONLog")
	appConfig.Port = c.Int("Port")
	appConfig.ConfigPath = c.String("ConfigPath")
	appConfig.StoragePath = c.String("StoragePath")
	appConfig.TmpPath = c.String("TmpPath")
	appConfig.DataPath = c.String("DataPath")
	appConfig.Version = version

	log.Infof("starting with config: %+v", appConfig)

	appConfig.UseAwsCli = true
	appConfig.HasGpg = true

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

	info, err = backup.CheckAWSClient()
	if err != nil {
		log.Warn(err)
		appConfig.UseAwsCli = false
	}
	log.Info(info)

	info, err = backup.CheckGpg()
	if err != nil {
		log.Warn(err)
		appConfig.HasGpg = false
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
	go server.Start(appConfig.Version)

	// wait for SIGINT (Ctrl+C) or SIGTERM (docker stop)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan

	log.Infof("shutting down %v signal received", sig)

	return nil
}
