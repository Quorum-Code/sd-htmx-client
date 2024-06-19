package endpoints

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Quorum-Code/sd-htmx-client/internal/authentication"
)

func (wsc *WSConfig) indexHandler(resp http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		tmpl := template.Must(template.ParseFiles("html/page-not-found.html"))
		tmpl.Execute(resp, nil)
		return
	}

	_, _ = authentication.GetUserFromHeader(req)

	tmpl := template.Must(template.ParseFiles("html/index.html"))
	tmpl.Execute(resp, nil)
}

func (wsc *WSConfig) PostIndexHandler(resp http.ResponseWriter, req *http.Request) {
	auth, err := authentication.RequestToToken(req)
	if err != nil {
		return
	}

	fmt.Println(auth.Token)
	fmt.Println(auth.Claim)

	d := map[string]string{
		"response": "got"}

	jsonData, err := json.Marshal(d)
	if err != nil {
		fmt.Println("failed marshaling json")
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonData)
}
