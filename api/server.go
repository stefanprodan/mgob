package api

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stefanprodan/mgob/config"
	"github.com/stefanprodan/mgob/scheduler"
	"net/http"
	"runtime"
	"strconv"
)

type HttpServer struct {
	Config *config.AppConfig
	Stats  *scheduler.Stats
}

func (s *HttpServer) Start(version string) {
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/version", func(w http.ResponseWriter, req *http.Request) {
		info := map[string]string{
			"mgob_version": version,
			"repository":   "github.com/stefanprodan/mgob",
			"go_version":   runtime.Version(),
			"os":           runtime.GOOS,
			"arch":         runtime.GOARCH,
			"max_procs":    strconv.FormatInt(int64(runtime.GOMAXPROCS(0)), 10),
			"goroutines":   strconv.FormatInt(int64(runtime.NumGoroutine()), 10),
			"cpu_count":    strconv.FormatInt(int64(runtime.NumCPU()), 10),
		}

		js, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
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
