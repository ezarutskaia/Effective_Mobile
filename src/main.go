package main

import (
	"fmt"
	"test_effective_mobile/app/domain"
	"test_effective_mobile/app/domain/models"
	"test_effective_mobile/app/repository"
	"test_effective_mobile/app/interfaces"
	"test_effective_mobile/app/controller"
	"test_effective_mobile/app"
)

func main() {
	fmt.Println("go started")

	db := repository.Engine()
	(*db).AutoMigrate(&models.Group{}, &models.Song{}, &models.Verse{})
	
	app := &app.App{
		Domain: &domain.Domain{},
		Repository: &repository.Repository{
			DB: db,
		},
		Interfaces: &interfaces.Interfaces{},
	}

	app.Interfaces.HttpServer.HandleHttpRequest(&controller.Controller{
		Repo: app.Repository,
		Domain: app.Domain,
	})
}