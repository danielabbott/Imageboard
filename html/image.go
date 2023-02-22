package html

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

	database "github.com/danielabbott/imageboard/database"
)

type TagData struct {
	Name        string
	URL         template.URL
	NameEscaped string
}

type HTMLImagePageData struct {
	Title   string
	Src     string
	ImageID int
	Tags    []TagData
}

func GenerateImagePage(db *database.DB, templates *HTMLTemplates, w http.ResponseWriter, req *http.Request) {
	// Query parameters

	query := req.URL.Query()
	if !query.Has("id") {
		w.WriteHeader(400)
		return
	}

	image_id, err := strconv.ParseInt(query.Get("id"), 10, 32)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	// Databse query

	title, err := database.GetImageTitle(db, database.ImageID(image_id))

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	tag_names, err := database.GetImageTags(db, database.ImageID(image_id))

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	tags := []TagData{}
	for _, name := range tag_names {
		tags = append(tags, TagData{
			Name:        name,
			URL:         template.URL("/?search=" + url.QueryEscape(name)),
			NameEscaped: url.QueryEscape(name),
		})
	}

	// Generate Html

	templates.image_page.Execute(w, HTMLImagePageData{
		Title:   title,
		Src:     "/i/" + strconv.FormatInt(int64(image_id), 10),
		ImageID: int(image_id),
		Tags:    tags,
	})
}

func DeleteImage(db *database.DB, w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(400)
		return
	}

	// Form data

	id, err := strconv.ParseInt(req.FormValue("id"), 10, 32)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// Database query

	database.DeleteImage(db, database.ImageID(id))

	// Redirect to home page

	http.Redirect(w, req, "/", 303)
}
