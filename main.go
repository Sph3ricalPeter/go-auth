package main

import (
	"fmt"
	"net/http"

	"github.com/Sph3ricalPeter/go-auth/auth"
	"github.com/Sph3ricalPeter/go-auth/config"
	"github.com/Sph3ricalPeter/go-auth/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func handleAlive(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func handleTestJwtAuth(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("JWT is valid"))
}

func main() {
	s := storage.NewPgxStorage()
	_ = s.CreateUser("admin", "admin")

	jwtAuth := auth.NewJwtAuth(s)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/alive", handleAlive)
	r.Post("/login", jwtAuth.HandleLogin)
	r.Post("/logout", jwtAuth.HandleLogout)
	r.Post("/refresh", jwtAuth.HandleRefresh)

	r.Group(func(r chi.Router) {
		r.Use(jwtAuth.JwtAuthHandler)
		r.Get("/test", handleTestJwtAuth)
	})

	_ = http.ListenAndServe(fmt.Sprintf("%s:%s", config.Host, config.Port), r)
}
