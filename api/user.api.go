package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/hiennguyen9874/stockk-go/controllers"
	"github.com/hiennguyen9874/stockk-go/repository"
	"gorm.io/gorm"
)

func UserRoute(router *chi.Mux, db *gorm.DB) {
	repo := repository.UserRepo{}

	userCtrler := controllers.NewUserHandler(db, repo)

	router.Group(func(r chi.Router) {
		r.Get("/user", userCtrler.GetUsers)
		r.Post("/user", userCtrler.CreateUser)
		r.Route("/user/{id}", func(r chi.Router) {
			r.Get("/", userCtrler.GetUser)
		})
	})
}
