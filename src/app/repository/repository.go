package repository

import (
	"time"
	"gorm.io/gorm"
	"test_effective_mobile/app/domain/models"
)

type Repository struct {
	DB *gorm.DB
}

func (repo *Repository) SaveGroup(group *models.Group) (id int, err error) {
	result := (*repo.DB).Create(group)
	return group.ID, result.Error
}

func (repo *Repository) SaveSong(song *models.Song) (id int, err error) {
	result := (*repo.DB).Create(song)
	return song.ID, result.Error
}

func (repo *Repository) SaveVerse(verse *models.Verse) (id int, err error) {
	result := (*repo.DB).Create(verse)
	return verse.ID, result.Error
}

func (repo *Repository) GetGroup(name string) (group *models.Group, err error) {
	result := (*repo.DB).Where("name = ?", name).First(&group)
	return group, result.Error
}

func (repo *Repository) GetGroupName(id int) (name string, err error) {
	var group *models.Group
	result := (*repo.DB).Where("id = ?", id).First(&group)
	return group.Name, result.Error
}

func (repo *Repository) ListGroups() (groups []*models.GroupDTO, err error) {
	result := (*repo.DB).Model(&models.Group{}).Find(&groups)
	return groups, result.Error
}

func (repo *Repository) GroupBySymbol(abc string) (groups []*models.GroupDTO, err error) {
	result := (*repo.DB).Model(&models.Group{}).Where("name ILIKE ?", "%"+abc+"%").Find(&groups)
	return groups, result.Error
}

func (repo *Repository) ListSongs() (songs []*models.SongDTOResp, err error) {
	result := (*repo.DB).Model(&models.Song{}).Find(&songs)
	return songs, result.Error
}

func (repo *Repository) FilterSongs(where map[string]interface{}) (songs []*models.SongDTOResp, err error) {
	result := (*repo.DB).Model(&models.Song{}).Preload("Group").Where(where).Find(&songs)
	return songs, result.Error
}

func (repo *Repository) ListVerses(songID int) (verses []*models.VerseDTOResp, err error) {
	result := (*repo.DB).Model(&models.Verse{}).Where("song_id = ?", songID).Find(&verses)
	return verses, result.Error
}

func (repo *Repository) PageVerses(songID int, page int) (verses []*models.VerseDTOResp, err error) {
	limit := 2
	offset := (page - 1) * limit

	result := (*repo.DB).Model(&models.Verse{}).Where("song_id = ?", songID).Limit(limit).Offset(offset).Find(&verses)
	return verses, result.Error
}

func (repo *Repository) UpdateSong(songID int, songDTO *models.SongDTO) error {
    updateFields := make(map[string]interface{})

    if songDTO.Group != "" {
		group, err := repo.GetGroup(songDTO.Group)
		if err == nil {
			updateFields["group_id"] = group.ID
		}
    }
    if songDTO.Song != "" {
        updateFields["name"] = songDTO.Song
    }
    if songDTO.ReleaseDate != "" {
		releaseDate, err := time.Parse("2006-01-02", songDTO.ReleaseDate)
		if err == nil {
			updateFields["release_date"] = releaseDate 
		}
    }
	if songDTO.Link != "" {
		updateFields["link"] = songDTO.Link
	}

    if len(updateFields) > 0 {
        result := (*repo.DB).Model(&models.Song{}).Where("id = ?", songID).Updates(updateFields)
		return result.Error
    }

    return nil
}

func (repo *Repository) DeleteVerses(songID int) error {
	result := (*repo.DB).Where("song_id = ?", songID).Delete(&models.Verse{})
	return result.Error
}

func (repo *Repository) DeleteSong(songID int) error {
	result := (*repo.DB).Where("id = ?", songID).Delete(&models.Song{})
	return result.Error
}