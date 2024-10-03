package controller

import (
	"time"
	"test_effective_mobile/app/domain/models"
	"test_effective_mobile/app/domain"
	"test_effective_mobile/app/repository"
)

type Controller struct {
	Repo   *repository.Repository
	Domain *domain.Domain
}

func (controller *Controller) CreateGroup(name string) (group *models.Group, err error) {
	group = controller.Domain.CreateGroup(name)
	_, err = controller.Repo.SaveGroup(group)
	return group, err
}

func (controller *Controller) CreateSong(name string, releasedate *time.Time, link string, group *models.Group) (id int, err error) {
	song := controller.Domain.CreateSong(name, releasedate, link, group)
	id, err = controller.Repo.SaveSong(song)
	return id, err
}

func (controller *Controller) CreateVerse(text string, idSong int) (id int, err error) {
	verse := controller.Domain.CreateVerse(text, idSong)
	id, err = controller.Repo.SaveVerse(verse)
	return id, err
}