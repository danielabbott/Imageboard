package html

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"strconv"

	database "github.com/danielabbott/imageboard/database"
)

func Upload(db *database.DB, w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(400)
		return
	}

	// Form data

	const max_file_size = 5 * 1024 * 1024

	req.Body = http.MaxBytesReader(w, req.Body, max_file_size+32*1024)

	title := req.FormValue("title")

	file, file_header, err := req.FormFile("file")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	if len(title) == 0 || file_header.Size > max_file_size {
		w.WriteHeader(400)
		return
	}

	if len(title) > 100 {
		title = title[:100]
	}

	// Database query

	image_id, err := database.NewImage(db, title)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	// Write file

	file_data, err := ioutil.ReadAll(file)

	// TODO convert file to lossy webp (even if file is already webp), provide download link to original

	if err != nil {
		database.RemoveImage(db, image_id)
		w.WriteHeader(500)
		return
	}

	err = ioutil.WriteFile("images/"+strconv.FormatInt(int64(image_id), 10), file_data, fs.FileMode(0777))

	if err != nil {
		database.RemoveImage(db, image_id)
		w.WriteHeader(500)
		return
	}

	// Redirect to the newly created image

	http.Redirect(w, req, "/image?id="+strconv.FormatInt(int64(image_id), 10), 303)
}
