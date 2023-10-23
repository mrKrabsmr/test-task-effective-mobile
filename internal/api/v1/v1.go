package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/mrKrabsmr/commerce-edu-api/internal/apps/persons"
)

func ConfigureRoutes(router chi.Router) {
	userApp := persons.NewApp()

	router.Route("/v1", func(r chi.Router) {
		userApp.ConfigureRoutes(r)
	})
}
