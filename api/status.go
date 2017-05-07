package api

import (
	"context"
	"github.com/pressly/chi/render"
	"github.com/stefanprodan/mgob/backup"
	"net/http"
)

type appStatus []backup.Result

func (a *appStatus) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func statusCtx(data appStatus) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "app.status", data))
			next.ServeHTTP(w, r)
		})
	}
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("app.status").(appStatus)
	render.JSON(w, r, data)
}
