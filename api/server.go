package api

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stefanprodan/mgob/config"
	"net/http"
)

type HttpServer struct {
	Config *config.AppConfig
}

func (s *HttpServer) Start() {
	http.Handle("/metrics", promhttp.Handler())

	logrus.Error(http.ListenAndServe(fmt.Sprintf(":%v", s.Config.Port), http.DefaultServeMux))
}
