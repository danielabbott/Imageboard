package html

import (
	"fmt"
	"net/http"
	"strconv"

	database "github.com/danielabbott/imageboard/database"
)

func RemoveTag(db *database.DB, w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(400)
		return
	}

	// Form data

	tag := req.FormValue("tag")
	image_id, err := strconv.ParseInt(req.FormValue("image_id"), 10, 32)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// Database query

	database.RemoveTag(db, database.ImageID(image_id), tag)

	http.Redirect(w, req, "/image?id="+strconv.FormatInt(int64(image_id), 10), 303)
}

func AddTag(db *database.DB, w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(400)
		return
	}

	// Form data

	tag := req.FormValue("tag")
	if len(tag) > 80 {
		w.WriteHeader(400)
		return
	}

	image_id, err := strconv.ParseInt(req.FormValue("image_id"), 10, 32)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// Database query

	err = database.AddTag(db, database.ImageID(image_id), tag)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	http.Redirect(w, req, "/image?id="+strconv.FormatInt(int64(image_id), 10), 303)
}
