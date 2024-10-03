package models

type GroupDTO struct {
	ID int
	Name string
}

type SongDTO struct {
	Group string            `json:"group"`
	Song string             `json:"song"`
	ReleaseDate string   	`json:"releaseDate"`
	Link string             `json:"link"`
	Text string             `json:"text"`
}

type SongDTOResp struct {
	Group Group				`json:"-"`				
	GroupID int             `json:"group_id"`
	Name string             `json:"song"`
	ReleaseDate string   	`json:"releaseDate"`
	Link string             `json:"link"`
}

type VerseDTO struct {
	SongID int              `query:"song_id"`
	Page int                `query:"page"`
}

type VerseDTOResp struct {
	SongID int
	Text string
}