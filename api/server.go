package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"github.com/stefanprodan/mgob/config"
	"fmt"
)

type HttpServer struct {
	Config *confing.AppConfig
}

func (s *HttpServer) Start() {
	http.Handle("/metrics", promhttp.Handler())

	logrus.Error(http.ListenAndServe(fmt.Sprintf(":%v", s.Config.Port), http.DefaultServeMux))
}
