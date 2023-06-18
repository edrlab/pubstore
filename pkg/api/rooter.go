package api

import (
	"github.com/edrlab/pubstore/pkg/stor"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type Api struct {
	stor *stor.Stor
}

func Init(s *stor.Stor) *Api {
	return &Api{stor: s}
}

func (api *Api) Rooter(r *chi.Mux) *chi.Mux {

	validate = validator.New()

	// api Routes CRUD Publication
	// Require Authentication
	r.Group(func(r chi.Router) {
		r.Get("/api/v1/publication/{id}", api.getPublicationHandler)
		r.Post("/api/v1/publication", api.createPublicationHandler)
		// r.Put("/api/v1/publication/{id}", apiV1PublicationPut)
		// r.Delete("/api/v1/publication/{id}", apiV1PublicationDelete)
	})

	// // api Routes CRUD User
	// // Require Authentication
	// r.Group(func(r chi.Router) {
	// 	r.Use(AuthMiddleware)
	// 	r.Get("/api/v1/user/{id}", apiV1UserGet)
	// 	r.Post("/api/v1/user/{id}", apiV1UserPost)
	// 	r.Put("/api/v1/user/{id}", apiV1UserPut)
	// 	r.Delete("/api/v1/user/{id}", apiV1UserDelete)
	// })

	return r
}
