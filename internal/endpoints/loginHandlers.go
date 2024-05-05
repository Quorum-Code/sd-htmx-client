package endpoints

import (
	"html/template"
	"net/http"
)

func (wsc *WSConfig) loginHandler(resp http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("login.html"))
	tmpl.Execute(resp, nil)
}

func (wsc *WSConfig) postLoginHandler(resp http.ResponseWriter, req *http.Request) {

}
