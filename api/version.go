package api

import (
	"context"
	"github.com/pressly/chi/render"
	"net/http"
	"runtime"
	"strconv"
)

type appVersion map[string]string

func (a *appVersion) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func appVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := appVersion{
				"mgob_version": version,
				"repository":   "github.com/stefanprodan/mgob",
				"go_version":   runtime.Version(),
				"os":           runtime.GOOS,
				"arch":         runtime.GOARCH,
				"max_procs":    strconv.FormatInt(int64(runtime.GOMAXPROCS(0)), 10),
				"goroutines":   strconv.FormatInt(int64(runtime.NumGoroutine()), 10),
				"cpu_count":    strconv.FormatInt(int64(runtime.NumCPU()), 10),
			}
			r = r.WithContext(context.WithValue(r.Context(), "app.version", data))
			next.ServeHTTP(w, r)
		})
	}
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("app.version").(appVersion)
	render.JSON(w, r, data)
}
