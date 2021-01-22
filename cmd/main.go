package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/namphung1998/auth-service-go/internal"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	handler := internal.Handler{}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Post("/greet", handler.HandleGreet())

	http.ListenAndServe(":3090", r)
}
