package interfaces

import (
	"time"
	"gorm.io/gorm"
	"strings"
	"strconv"
	"net/http"
	"os"
	log "github.com/sirupsen/logrus"
	"test_effective_mobile/app/controller"
	"test_effective_mobile/app/domain/models"
	"github.com/labstack/echo/v4"
)

type Data map[string]interface{}

type HttpServer struct{}

func (server HttpServer) Response (c echo.Context, message string, data map[string]interface{}) (error) {
	return c.JSON(http.StatusOK, map[string]interface{}{
        "message": message,
		"data": data,
    })
}

// AddGroup		 godoc
// @Summary      Add group
// @Accept       json
// @Produce      json
// @Router       /group/add [post]
// @Param        name   body   string  true  "Group name"
// @Success      200  {object}	models.Response{data=models.IDResponse}
func (server HttpServer) AddGroup(controller *controller.Controller) (func(c echo.Context)  (err error)) {
	return func(c echo.Context) (err error) {
		log.Debug("Handle add group")

		name := new(models.WordDTO)
		if err := c.Bind(name); err != nil {
			log.Info("Error group name")
			return server.Response(c, "Data reading error", Data{"id": nil})
		}
		
		var group *models.Group
		groupExist, err := controller.Repo.GetGroup(name.Name)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				groupNew, err := controller.CreateGroup(name.Name)
				if err != nil {
					log.Info("Cannot create group")
					return server.Response(c, "data recording error", Data{"id": nil})
				}
				group = groupNew
			} else {
				log.Info("database is not available")
				return server.Response(c, "database error", Data{"id": nil})
			}
		} else {
			group = groupExist
		}
	
		return server.Response(c, "note was created", Data{"id": group.ID})
	}
}

// AddSong		 godoc
// @Summary      Add group
// @Accept       json
// @Produce      json
// @Router       /song/add [post]
// @Param        group   body   string  true  "Group name"
// @Param        song   body   string  true  "Song name"
// @Param        link   body   string  false  "Link"
// @Param        text   body   string  false  "Song text. Versers are divided by /n/n"
// @Success      200  {object}	models.Response{data=models.IDResponse}
func (server HttpServer) AddSong(controller *controller.Controller) (func(c echo.Context)  (err error)) {
	return func(c echo.Context) (err error) {
		log.Debug("Handle add song")
		songDTO := new(models.SongDTO)
		
		if err := c.Bind(songDTO); err != nil {
			log.Info("The data for song is incorrect")
			return server.Response(c, "data reading error", Data{"id": nil})
		}
		
		var releaseDate *time.Time
		releaseDate = nil
		if songDTO.ReleaseDate != "" {
			releaseDateTemp, err := time.Parse("2006-01-02", songDTO.ReleaseDate)
			if err != nil {
				log.Info("The date is incorrect")
				return server.Response(c, "invalid date format", Data{"id": nil})
			} else {
				releaseDate = &releaseDateTemp
			}
		}

		var group *models.Group
		groupExist, err := controller.Repo.GetGroup(songDTO.Group)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				groupNew, err := controller.CreateGroup(songDTO.Group)
				if err != nil {
					log.Info("Cannot create group")
					return server.Response(c, "data recording error", Data{"id": nil})
				}
				group = groupNew
			} else {
				log.Info("database is not available")
				return server.Response(c, "database error", Data{"id": nil})
			}
		} else {
			group = groupExist
		}

		id, err := controller.CreateSong(songDTO.Song, releaseDate, songDTO.Link, group)
		if err != nil {
			log.Info("database is not available")
			return server.Response(c, "data recording error", Data{"id": nil})
		}

		if songDTO.Text != "" {
			verses := strings.Split(songDTO.Text, "\n\n")
			for _, verse := range verses {
				_, err = controller.CreateVerse(verse, id)
				if err != nil {
					log.Info("database is not available")
					return server.Response(c, "database error", Data{"id": nil})
				}
			}
		}
		
		return server.Response(c, "note was created", Data{"id": id})
	}
}

// GroupInfo	 godoc
// @Summary      Get group info
// @Accept       json
// @Produce      json
// @Router       /group/info [post]
// @Param        name   body   string  true  "Group name"
// @Success      200  {object}	models.Response{data=models.GroupsDTO}
func (server HttpServer) GroupInfo(controller *controller.Controller) (func(c echo.Context)  (err error)) {
	return func(c echo.Context) (err error) {
		log.Debug("Handle group info")
		symbols := new(models.WordDTO)
		if err := c.Bind(symbols); err != nil {
			log.Info("The symbols for group is incorrect")
			return server.Response(c, "data reading error", Data{"groups": ""})
		}
		
		if symbols.Name == "" {
			groups, err := controller.Repo.ListGroups()
			if err != nil {
				log.Info("database is not available")
				return server.Response(c, "database error", Data{"groups": ""})
			}
			return server.Response(c, "success", Data{"groups": groups})
		}

		groups, err := controller.Repo.GroupBySymbol(symbols.Name)
		if err != nil {
			log.Info("database is not available")
			return server.Response(c, "database error", Data{"groups": ""})
		}
		
		return server.Response(c, "success", Data{"groups": groups})
	}
}

// SongInfo	 	 godoc
// @Summary      Get song info
// @Accept       json
// @Produce      json
// @Router       /song/info [post]
// @Param        group   body   string  true  "Group name"
// @Param        song   body   string  true  "Song name" 
// @Param        releaseDat   body   string  false  "Date release"
// @Param        link   body   string  false  "HTTP Link"
// @Success      200  {object}	models.Response{data=models.SongsDTO}
func (server HttpServer) SongInfo(controller *controller.Controller) (func(c echo.Context)  (err error)) {
	return func(c echo.Context) (err error) {
		log.Debug("Handle song info")
		where := make(map[string]interface{})

		songDTO := new(models.SongDTO)
		if err := c.Bind(songDTO); err != nil {
			log.Info("The data for song is incorrect")
			return server.Response(c, "data reading error", Data{"songs": ""})
		}

		if songDTO.Group != "" {
			group, err := controller.Repo.GetGroup(songDTO.Group)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					log.Info("Cannot find group")
					return server.Response(c, "no such group", Data{"songs": ""})
				}
				log.Info("database is not available")
				return server.Response(c, "database error", Data{"songs": ""})
			}
			where["group_id"] = group.ID
		}
		if songDTO.Song != "" {
			where["name"] = songDTO.Song
		}
		if songDTO.ReleaseDate != "" {
			releaseDate, err := time.Parse("2006-01-02", songDTO.ReleaseDate)
    		if err != nil {
				log.Info("The date is incorrect")
        		return server.Response(c, "invalid date format", Data{"songs": ""})  
    		}
			where["release_date"] = releaseDate
		}

		if len(where) == 0 {
			songs, err := controller.Repo.ListSongs()
			if err != nil {
				log.Info("database is not available")
				return server.Response(c, "database error", Data{"songs": ""})
			}
			return server.Response(c, "success", Data{"songs": songs})
		}

		songs, err := controller.Repo.FilterSongs(where)
		if err != nil {
			log.Info("database is not available")
			return server.Response(c, "database error", Data{"songs": ""})
		}

		return server.Response(c, "success", Data{"songs": songs})
	}
}

// VerseInfo	 godoc
// @Summary      Get song verse
// @Accept       json
// @Produce      json
// @Router       /verse/info [get]
// @Param        song_id   query   int  true  "Song ID"
// @Param        page   query   int  false  "Page" 
// @Success      200  {object}	models.Response{data=models.VersesDTO}
func (server HttpServer) VerseInfo(controller *controller.Controller) (func(c echo.Context)  (err error)) {
	return func(c echo.Context) (err error) {
		log.Debug("Handle verse info")
		verseDTO := new(models.VerseDTO)
		if err := c.Bind(verseDTO); err != nil {
			log.Info("The data for verse is incorrect")
			return server.Response(c, "data reading error", Data{"verses": ""})
		}
		
		var verses []*models.VerseDTOResp
		if verseDTO.Page == 0 {
			verses, err = controller.Repo.ListVerses(verseDTO.SongID)
			if err != nil {
				log.Info("database is not available")
				return server.Response(c, "database error", Data{"verses": ""})
			}
			return server.Response(c, "Success", Data{"verses": verses})
		}
		
		verses, err = controller.Repo.PageVerses(verseDTO.SongID, verseDTO.Page)	
		if err != nil {
			log.Info("database is not available")
			return server.Response(c, "database error", Data{"verses": ""})
		}
		return server.Response(c, "success", Data{"verse": verses})
	}
}

// UpdateSong	 godoc
// @Summary      Update song
// @Accept       json
// @Produce      json
// @Router       /song/update [get]
// @Param        group   body   string  true  "Group name"
// @Param        song   body   string  true  "Song name" 
// @Param        releaseDat   body   string  false  "Date release"
// @Param        link   body   string  false  "HTTP Link"
// @Success      200  {object}	models.Response{data=models.IDResponse}
func (server HttpServer) UpdateSong(controller *controller.Controller) (func(c echo.Context)  (err error)) {
	return func(c echo.Context) (err error) {
		log.Debug("Handle song update")
		id := c.QueryParam("id")

		songDTO := new(models.SongDTO)
		if err := c.Bind(songDTO); err != nil {
			log.Info("The data for song is incorrect")
			return server.Response(c, "data reading error", Data{"id": nil})
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			log.Info("id is incorrect")
			return server.Response(c, "invalid id", Data{"id": nil})
		}

		err = controller.Repo.UpdateSong(idInt, songDTO)
		if err != nil {
			log.Info("cannot update song")
			return server.Response(c, "update error", Data{"id": nil})
		}

		if songDTO.Text != "" {
			err = controller.Repo.DeleteVerses(idInt)
			if err != nil {
				log.Info("canot delete verse")
				return server.Response(c, "delete error", Data{"id": nil})
			}
			verses := strings.Split(songDTO.Text, "\n\n")
			for _, verse := range verses {
				_, err = controller.CreateVerse(verse, idInt)
				if err != nil {
					log.Info("database is not available")
					return server.Response(c, "database error", Data{"id": nil})
				}
			}
		}

		return server.Response(c, "update successful", Data{"id": id})
	}
}


// DeleteSong	 godoc
// @Summary      Delete song
// @Accept       json
// @Produce      json
// @Router       /song/delete [delete]
// @Param        song_id   query   int  true  "Song ID"
// @Success      200  {object}	models.Response{data=models.IDResponse}
func (server HttpServer) DeleteSong(controller *controller.Controller) (func(c echo.Context)  (err error)) {
	return func(c echo.Context) (err error) {
		log.Debug("Handle delete song")
		id := c.QueryParam("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			log.Info("id is incorrect")
			return server.Response(c, "invalid id", Data{"id": ""})
		}

		err = controller.Repo.DeleteSong(idInt)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Info("Cannot find song")
				return server.Response(c, "no such song", Data{"id": idInt})
			}
			log.Info("database is not available")
			return server.Response(c, "database error", Data{"id": ""})
		}

		err = controller.Repo.DeleteVerses(idInt)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Info("Cannot find verse")
				return server.Response(c, "no such verse", Data{"id": idInt})
			}
			log.Info("database is not available")
			return server.Response(c, "database error", Data{"id": ""})
		}

		return server.Response(c, "song was deleted", Data{"id": idInt})
	}
}

func (server HttpServer) HandleHttpRequest(controller *controller.Controller) {
	log.Info("HTTP server have started.")
	e := echo.New()
	
	e.POST("/group/add", server.AddGroup(controller))
	e.POST("/song/add", server.AddSong(controller))
	e.POST("/group/info", server.GroupInfo(controller))
	e.POST("/song/info", server.SongInfo(controller))
	e.GET("/verse/info", server.VerseInfo(controller))
	e.PUT("/song/update", server.UpdateSong(controller))
	e.DELETE("/song/delete", server.DeleteSong(controller))
	
	e.Logger.Fatal(e.Start("0.0.0.0:" + os.Getenv("port")))
}