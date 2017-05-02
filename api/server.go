package api

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stefanprodan/mgob/config"
	"github.com/stefanprodan/mgob/scheduler"
	"net/http"
	"time"
)

type HttpServer struct {
	Config *config.AppConfig
	Stats  *scheduler.Stats
}

func (s *HttpServer) Start(version string) {
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/version", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, version)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		js, err := json.MarshalIndent(s.Stats.GetAll(), "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})

	http.Handle("/", http.FileServer(http.Dir(s.Config.StoragePath)))

	logrus.Error(http.ListenAndServe(fmt.Sprintf(":%v", s.Config.Port), http.DefaultServeMux))
}
