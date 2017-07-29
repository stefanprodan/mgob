package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stefanprodan/mgob/db"
)

type appStatus []*db.Status

func (a *appStatus) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func statusCtx(store *db.StatusStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, err := store.GetAll()
			if err != nil {
				render.Status(r, 500)
				render.JSON(w, r, map[string]string{"error": err.Error()})
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "app.status", appStatus(data)))
			next.ServeHTTP(w, r)
		})
	}
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("app.status").(appStatus)
	render.JSON(w, r, data)
}

func getPlanStatus(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("app.status").(appStatus)
	planID := chi.URLParam(r, "planID")
	for _, s := range data {
		if s.Plan == planID {
			render.JSON(w, r, s)
			return
		}
	}

	render.Status(r, 404)
	render.JSON(w, r, map[string]string{"error": "Plan not found"})
}
