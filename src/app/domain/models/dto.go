package models

type Response struct {
	Message string			`json:"message"`			
	Data interface{}		`json:"data"`
}

type IDResponse struct {
	ID int	`json:"id"`
}

type WordDTO struct {
	Name string `json:"name"`
}

type GroupDTO struct {
	ID int
	Name string
}

type GroupsDTO []Group

type SongsDTO struct {
	Songs []SongDTOResp 
}

type VersesDTO struct {
	Verses []VerseDTOResp
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