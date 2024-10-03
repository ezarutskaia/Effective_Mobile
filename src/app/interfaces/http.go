package interfaces

import (
	"fmt"
	"time"
	"gorm.io/gorm"
	"strings"
	// "sync"
	// "math/rand"
	 "strconv"
	 "net/http"
	// "context"
	// "errors"
	// "encoding/json"
	"test_effective_mobile/app/controller"
	"test_effective_mobile/app/domain/models"
	"github.com/labstack/echo/v4"
)

type Data map[string]interface{}

type HttpServer struct{}

type WordRequest struct {
	Name string `json:"name"`
}

func (server HttpServer) Response (c echo.Context, message string, data map[string]interface{}) (error) {
	return c.JSON(http.StatusOK, map[string]interface{}{
        "message": message,
		"data": data,
    })
}

func (server HttpServer) HandleHttpRequest(controller *controller.Controller) {
	fmt.Println("HTTP server have started.")
	e := echo.New()

	// group/add

	e.POST("/group/add", func(c echo.Context) (err error) {
		name := new(WordRequest)
		if err := c.Bind(name); err != nil {
			return server.Response(c, "data reading error", Data{"id": ""})
		}

		id, err := controller.CreateGroup(name.Name)
		if err != nil {
			return server.Response(c, "data recording error", Data{"id": ""})
		}

		return server.Response(c, "note was created", Data{"id": id})
	})

	// song/add

	e.POST("/song/add", func(c echo.Context) (err error) {
		songDTO := new(models.SongDTO)
		
		if err := c.Bind(songDTO); err != nil {
			return server.Response(c, "data reading error", Data{"id": ""})
		}
		
		var releaseDate *time.Time
		releaseDate = nil
		if songDTO.ReleaseDate != "" {
			releaseDateTemp, err := time.Parse("2006-01-02", songDTO.ReleaseDate)
			if err != nil {
				return server.Response(c, "invalid date format", Data{"id": ""})
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
					return server.Response(c, "data recording error", Data{"id": ""})
				}
				group = groupNew
			} else {
				return server.Response(c, "database error", Data{"id": ""})
			}
		} else {
			group = groupExist
		}

		id, err := controller.CreateSong(songDTO.Song, releaseDate, songDTO.Link, group)
		if err != nil {
			return server.Response(c, "data recording error", Data{"id": ""})
		}

		if songDTO.Text != "" {
			verses := strings.Split(songDTO.Text, "\n\n")
			for _, verse := range verses {
				_, err = controller.CreateVerse(verse, id)
				if err != nil {
					return server.Response(c, "database error", Data{"id": ""})
				}
			}
		}
		
		return server.Response(c, "note was created", Data{"id": id})
	})

	// group/info

	e.GET("/group/info", func(c echo.Context) (err error) {
		symbols := new(WordRequest)
		if err := c.Bind(symbols); err != nil {
			return server.Response(c, "data reading error", Data{"groups": ""})
		}
		
		if symbols.Name == "" {
			groups, err := controller.Repo.ListGroups()
			if err != nil {
				return server.Response(c, "database error", Data{"groups": ""})
			}
			return server.Response(c, "success", Data{"groups": groups})
		}

		groups, err := controller.Repo.GroupBySymbol(symbols.Name)
		if err != nil {
			return server.Response(c, "database error", Data{"groups": ""})
		}
		
		return server.Response(c, "success", Data{"groups": groups})
	})

	// song/info

	e.GET("/song/info", func(c echo.Context) (err error) {
		where := make(map[string]interface{})

		songDTO := new(models.SongDTO)
		if err := c.Bind(songDTO); err != nil {
			return server.Response(c, "data reading error", Data{"songs": ""})
		}

		if songDTO.Group != "" {
			group, err := controller.Repo.GetGroup(songDTO.Group)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return server.Response(c, "no such group", Data{"songs": ""})
				}
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
        		return server.Response(c, "invalid date format", Data{"songs": ""})  
    		}
			where["release_date"] = releaseDate
		}

		if len(where) == 0 {
			songs, err := controller.Repo.ListSongs()
			if err != nil {
				return server.Response(c, "database error", Data{"songs": ""})
			}
			return server.Response(c, "success", Data{"groups": songs})
		}

		songs, err := controller.Repo.FilterSongs(where)
		if err != nil {
			return server.Response(c, "database error", Data{"songs": ""})
		}

		return server.Response(c, "success", Data{"groups": songs})
	})

	// verse/info

	e.GET("/verse/info", func(c echo.Context) (err error) {
		verseDTO := new(models.VerseDTO)
		if err := c.Bind(verseDTO); err != nil {
			return server.Response(c, "data reading error", Data{"verses": ""})
		}

		var verses []*models.VerseDTOResp
		if verseDTO.Page == 0 {
			verses, err = controller.Repo.ListVerses(verseDTO.SongID)
			if err != nil {
				return server.Response(c, "database error", Data{"verses": ""})
			}
			return server.Response(c, "success", Data{"verses": verses})
		}
		
		verses, err = controller.Repo.PageVerses(verseDTO.SongID, verseDTO.Page)
		fmt.Printf("%v", verses)
		if err != nil {
			return server.Response(c, "database error", Data{"verses": ""})
		}
		return server.Response(c, "success", Data{"verse": verses})
	})

	// song/update

	e.PUT("/song/update", func(c echo.Context) (err error) {
		id := c.QueryParam("id")

		songDTO := new(models.SongDTO)
		if err := c.Bind(songDTO); err != nil {
			return server.Response(c, "data reading error", Data{"id": ""})
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			return server.Response(c, "invalid id", Data{"id": ""})
		}

		err = controller.Repo.UpdateSong(idInt, songDTO)
		if err != nil {
			return server.Response(c, "update error", Data{"id": ""})
		}

		if songDTO.Text != "" {
			err = controller.Repo.DeleteVerses(idInt)
			if err != nil {
				return server.Response(c, "delete error", Data{"id": ""})
			}
			verses := strings.Split(songDTO.Text, "\n\n")
			for _, verse := range verses {
				_, err = controller.CreateVerse(verse, idInt)
				if err != nil {
					return server.Response(c, "database error", Data{"id": ""})
				}
			}
		}

		return server.Response(c, "update successful", Data{"id": id})
	})

	// song/delete

	e.DELETE("/song/delete", func(c echo.Context) (err error) {
		id := c.QueryParam("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			return server.Response(c, "invalid id", Data{"songs_id": ""})
		}

		err = controller.Repo.DeleteSong(idInt)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return server.Response(c, "no such song", Data{"songs_id": idInt})
			}
			return server.Response(c, "database error", Data{"songs": ""})
		}

		err = controller.Repo.DeleteVerses(idInt)
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return server.Response(c, "database error", Data{"songs": ""})
				}
			}

		return server.Response(c, "song was deleted", Data{"songs_id": idInt})
	})

	e.Logger.Fatal(e.Start("0.0.0.0:6050"))
}