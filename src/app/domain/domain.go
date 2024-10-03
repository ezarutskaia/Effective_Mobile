package domain

import (
	"time"
	"test_effective_mobile/app/domain/models"
)

type Domain struct {}

func (domain *Domain) CreateGroup(group string) (*models.Group){
	return &models.Group{Name: group}
}

func (domain *Domain) CreateSong(name string, releasedate *time.Time, link string, group *models.Group) (*models.Song) {
	return &models.Song{Name: name, ReleaseDate: releasedate, Link: link,GroupID: group.ID}
}

func (domain *Domain) CreateVerse(text string, idSong int) (*models.Verse) {
	return &models.Verse{Text: text, SongID: idSong}
}