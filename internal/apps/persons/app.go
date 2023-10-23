package persons

import (
	"github.com/go-chi/chi/v5"
)

type App struct {
	controller *controller
}

func NewApp() *App {
	return &App{
		controller: newController(),
	}
}

func (a *App) ConfigureRoutes(router chi.Router) {
	router.Get("/persons/", a.controller.list)
	router.Post("/persons/", a.controller.create)
	router.Put("/persons/{id}/", a.controller.update)
	router.Patch("/persons/{id}/", a.controller.update)
	router.Delete("/persons/{id}/", a.controller.delete)
}
