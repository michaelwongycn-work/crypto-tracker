package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/michaelwongycn/crypto-tracker/controller"
	"github.com/michaelwongycn/crypto-tracker/handler/middleware"
)

type handler struct {
	timeout    time.Duration
	controller controller.Controller
	cors       *cors.Cors
}

func NewHandler(timeout time.Duration, controller controller.Controller) *handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return &handler{
		timeout:    timeout,
		controller: controller,
		cors:       c,
	}
}

func (h *handler) StartRoute() *http.Server {
	r := chi.NewRouter()

	r.Use(h.cors.Handler)
	r.Get("/ping", h.controller.Ping)

	r.Post("/login", h.controller.Login)
	r.Post("/register", h.controller.Register)
	r.Post("/logout", h.controller.Logout)
	r.Post("/refresh-token", h.controller.RefreshToken)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate)

		r.Get("/crypto", h.controller.ShowUserAsset)
		r.Post("/crypto", h.controller.InsertUserAsset)
		r.Delete("/crypto", h.controller.DeleteUserAsset)
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", 2000),
		WriteTimeout: h.timeout * time.Second,
		ReadTimeout:  h.timeout * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("listen: %s", err)
		}
	}()

	return srv
}
