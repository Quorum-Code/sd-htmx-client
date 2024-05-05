package endpoints

import (
	"html/template"
	"net/http"
)

func (wsc *WSConfig) indexHandler(resp http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		tmpl := template.Must(template.ParseFiles("html/page-not-found.html"))
		tmpl.Execute(resp, nil)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/index.html"))
	tmpl.Execute(resp, nil)
}
