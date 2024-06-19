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

	type response struct {
		Status       int    `json:"status"`
		Message      string `json:"message"`
		AccessToken  string `json:"access-token"`
		RefreshToken string `json:"resfresh-token"`
	}

	r := response{}

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

	success := wsc.Database.IsValidCredentials(username, password)

	if success {
		fmt.Println("Login successful")
		r.Message = "Login success"
		r.Status = 200
	} else {
		fmt.Println("Login failed")
		r.Message = "Login failed"
		r.Status = 409
	}

	resp.WriteHeader(r.Status)
	resp.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(r)
	if err == nil {
		resp.Write(jsonResponse)
	}
}
