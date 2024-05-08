package endpoints

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func (wsc *WSConfig) loginHandler(resp http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/login.html"))
	tmpl.Execute(resp, nil)
}

func (wsc *WSConfig) postLoginHandler(resp http.ResponseWriter, req *http.Request) {
	type body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	b := body{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&b)
	if err != nil {
		resp.WriteHeader(401)
		resp.Write([]byte("unable to parse json"))
		fmt.Println("json parsing failed")
		return
	}

	username := b.Username
	password := b.Password

	fmt.Println("Login attempt by: ", username)

	wsc.Database.IsValidCredentials(username, password)

}
