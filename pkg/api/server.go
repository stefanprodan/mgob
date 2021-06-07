package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/stefanprodan/mgob/pkg/config"
	"github.com/stefanprodan/mgob/pkg/db"
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
		r.Get("/{planID}", getPlanStatus)
	})

	r.Route("/backup", func(r chi.Router) {
		r.Use(configCtx(*s.Config))
		r.Post("/{planID}", postBackup)
	})

	FileServer(r, "/storage", http.Dir(s.Config.StoragePath))

	log.Error(http.ListenAndServe(fmt.Sprintf("%s:%v", s.Config.Host, s.Config.Port), r))
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
