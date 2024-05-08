package endpoints

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Quorum-Code/sd-htmx-client/internal/authentication"
)

func (wsc *WSConfig) signupHandler(resp http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/sign-up.html"))
	tmpl.Execute(resp, nil)
}

func (wsc *WSConfig) postSignupHandler(resp http.ResponseWriter, req *http.Request) {
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

	fmt.Println("Username: ", username)
	fmt.Println("Password: ", password)

	user, err := wsc.Database.TryAddUser(username, password)
	if err != nil {
		resp.WriteHeader(http.StatusConflict)
		resp.Write([]byte("failed to create user..."))
		return
	}

	tkn, rtkn, err := authentication.CreateTokens(user.ID)
	if err != nil {
		resp.WriteHeader(409)
		resp.Write([]byte("failed to create auth tokens"))
		return
	}

	response := map[string]string{"access-token": tkn, "refresh-token": rtkn}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		resp.WriteHeader(400)
		fmt.Println("json marshalling failed")
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonResponse)
}
