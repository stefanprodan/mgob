package api

import (
	"github.com/pressly/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func metricsRouter() http.Handler {
	promHandler := func(next http.Handler) http.Handler { return promhttp.Handler() }
	emptyHandler := func(w http.ResponseWriter, r *http.Request) {}
	r := chi.NewRouter()
	r.Use(promHandler)
	r.Get("/", emptyHandler)
	return r
}
