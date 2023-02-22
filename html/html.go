package html

import "html/template"

type HTMLTemplates struct {
	index      *template.Template
	image_page *template.Template
}

func LoadTemplates() HTMLTemplates {
	return HTMLTemplates{
		index:      template.Must(template.ParseFiles("html_templates/index.html")),
		image_page: template.Must(template.ParseFiles("html_templates/image.html")),
	}
}
