package api

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/stefanprodan/mgob/config"
	"net/http"
	"github.com/stefanprodan/mgob/db"
)

type HttpServer struct {
	Config *config.AppConfig
	Stats  *db.StatusStore
}

func (s *HttpServer) Start(version string) {

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	if s.Config.LogLevel == "debug" {
		r.Use(middleware.DefaultLogger)
	}

	r.Mount("/metrics", metricsRouter())
	r.Mount("/debug", middleware.Profiler())

	r.Route("/version", func(r chi.Router) {
		r.Use(appVersionCtx(version))
		r.Get("/", getVersion)
	})

	r.Route("/status", func(r chi.Router) {
		r.Use(statusCtx(s.Stats))
		r.Get("/", getStatus)
	})

	r.Route("/backup", func(r chi.Router) {
		r.Use(configCtx(*s.Config))
		r.Post("/:planID", postBackup)
	})

	r.FileServer("/storage", http.Dir(s.Config.StoragePath))

	logrus.Error(http.ListenAndServe(fmt.Sprintf(":%v", s.Config.Port), r))
}
