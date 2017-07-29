package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func metricsRouter() http.Handler {
	promHandler := func(next http.Handler) http.Handler { return promhttp.Handler() }
	emptyHandler := func(w http.ResponseWriter, r *http.Request) {}
	r := chi.NewRouter()
	r.Use(promHandler)
	r.Get("/", emptyHandler)
	return r
}
