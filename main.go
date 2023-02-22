package main

import (
	"fmt"
	"net/http"

	database "github.com/danielabbott/imageboard/database"
	html "github.com/danielabbott/imageboard/html"
)

var db database.DB
var html_templates html.HTMLTemplates

func main() {
	db = database.InitDatabase()
	html_templates = html.LoadTemplates()

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		html.GenerateIndexPage(&db, &html_templates, w, req)
	})

	http.HandleFunc("/image", func(w http.ResponseWriter, req *http.Request) {
		html.GenerateImagePage(&db, &html_templates, w, req)
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			html.Upload(&db, w, req)
		} else {
			http.ServeFile(w, req, "static/upload.html")
		}
	})

	http.HandleFunc("/delete_image", func(w http.ResponseWriter, req *http.Request) {
		html.DeleteImage(&db, w, req)
	})

	http.HandleFunc("/add_tag", func(w http.ResponseWriter, req *http.Request) {
		html.AddTag(&db, w, req)
	})

	http.HandleFunc("/remove_tag", func(w http.ResponseWriter, req *http.Request) {
		html.RemoveTag(&db, w, req)
	})

	http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/i/", http.StripPrefix("/i/", http.FileServer(http.Dir("./images/"))))

	fmt.Println("Starting")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
