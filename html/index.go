package html

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	database "github.com/danielabbott/imageboard/database"
)

type HTMLIndexData struct {
	Title         string
	SearchBoxText string
	NoResults     template.HTML
	Images        template.HTML
	NextButton    template.HTML
}

func atoi(a string) (int32, error) {
	i, err := strconv.ParseInt(a, 10, 32)
	return int32(i), err
}

func GenerateIndexPage(db *database.DB, templates *HTMLTemplates, w http.ResponseWriter, req *http.Request) {
	// 404 for all but index page (net/http seems to think this handler should run for every possible web page)
	if req.URL.Path != "/" {
		w.WriteHeader(404)
		return
	}

	// Query parameters

	// This value is either the 'ser' (serial) field or image id depending on which query is used
	var last_seen_image int32 = -1

	query := req.URL.Query()
	if query.Has("last") {
		x, err := strconv.ParseInt(query.Get("last"), 10, 32)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		last_seen_image = int32(x)
	}

	// Database query

	// min & max comparison values are for the back/foward buttons, they are either image IDs or the image_tag ser field
	image_id_results, lowest_image_value, err :=
		database.GetImages(db, last_seen_image, strings.Split(query.Get("search"), ","))

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	// Generate Html

	no_results_text := ""
	if len(image_id_results) == 0 {
		no_results_text = `<p class="centre">`
		if last_seen_image != -1 {
			no_results_text += `No more results`
		} else {
			no_results_text += `No results`
		}
		no_results_text += "</p>"
	}

	images_html := ""
	for _, id := range image_id_results {
		id_str := strconv.FormatInt(int64(id), 10)
		images_html = images_html +
			`<a href="image?id=` +
			id_str +
			`" class="image_thumb"><div><img src="i/` +
			id_str +
			`" class="image_thumb_img"></div></a>`
	}

	button_html := ""

	if len(image_id_results) == 48 && lowest_image_value > 1 {
		button_html +=
			`<form action="/" method="get">`
		if query.Has("search") {
			button_html += `<input type="hidden" name="search" value="` + query.Get("search") + `">`
		}
		button_html += `<input type="hidden" name="last" value="` + strconv.Itoa(int(lowest_image_value)) + `">` +
			`<button type="submit">Next Page</button>` +
			`</form>`
	} else {
		button_html += `<button disabled>Next Page</button>`
	}

	title := query.Get("search")
	if len(title) == 0 {
		title = "Imageboard"
	}

	templates.index.Execute(w, HTMLIndexData{
		Title:         title,
		SearchBoxText: query.Get("search"),
		NoResults:     template.HTML(no_results_text),
		Images:        template.HTML(images_html),
		NextButton:    template.HTML(button_html),
	})
}
