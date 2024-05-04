package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env")

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/verify-login", verifyLoginHandler)
	http.HandleFunc("/create-account", createAccountHandler)
	http.HandleFunc("/logged-in", loggedInHandler)
	http.HandleFunc("/", handler)
	http.HandleFunc("/events", eventHandler)
	http.HandleFunc("/events/redirect", redirectHandler)

	http.HandleFunc("POST /api/sign-up", PostApiSignup)

	// startFileServer()
	http.ListenAndServe(":8080", nil)
}

func PostApiSignup(resp http.ResponseWriter, req *http.Request) {
	response := map[string]string{"token": "some token"}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		resp.WriteHeader(400)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonResponse)
}

func redirectHandler(resp http.ResponseWriter, req *http.Request) {
	time.Sleep(3 * time.Second)

	resp.Header().Set("Content-Type", "text/event-stream")
	resp.Header().Set("Cache-Control", "no-cache")
	resp.Header().Set("Connection", "keep-alive")

	fmt.Fprintf(resp, "event: redirect\n")
	data, _ := json.Marshal(map[string]string{"redirectTo": "/logged-in"})
	fmt.Fprintf(resp, "data: %s\n\n", data)

	resp.(http.Flusher).Flush()
}

func loggedInHandler(resp http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("passed.html"))
	tmpl.Execute(resp, nil)
}

func createAccountHandler(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		resp.WriteHeader(400)
		resp.Write([]byte("Error loading form"))
		return
	}

	username := req.PostFormValue("username")
	password := req.PostFormValue("password")

	fmt.Println(username)
	fmt.Println(password)

	jwtSecret := os.Getenv("JWT_SECRET")
	fmt.Println(jwtSecret)

	sda := jwt.RegisteredClaims{Issuer: "sd-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour).UTC()),
		Subject:   "no-subject"}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, sda)

	signedToken, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		resp.WriteHeader(400)
		return
	}

	fmt.Println(token)
	fmt.Println(signedToken)

	tmpl := template.Must(template.ParseFiles("verify-login.html"))
	tmpl.Execute(resp, nil)
}

func verifyLoginHandler(resp http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("verify-login.html"))
	tmpl.Execute(resp, nil)
}

func loginHandler(resp http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("login.html"))
	tmpl.Execute(resp, nil)
}

func startFileServer() {
	d := "./static"
	fs := http.FileServer(http.Dir(d))

	http.Handle("/static/", fs)

	// _, err := os.Stat(d)
	// if os.IsNotExist(err) {
	// 	fmt.Printf("Directory '%s' not found.\n", d)
	// 	return
	// }

	// fileServer := http.FileServer(http.Dir("public"))

	// http.Handle("/public/", http.StripPrefix("/static/", fileServer))
}

func handler(resp http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		resp.WriteHeader(404)
		resp.Write([]byte("Page not found..."))
		return
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(resp, nil)
}

func eventHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/event-stream")
	resp.Header().Set("Cache-Control", "no-cache")
	resp.Header().Set("Connection", "keep-alive")

	dataCh := make(chan string)

	_, cancel := context.WithCancel(req.Context())
	defer cancel()

	go func() {
		for data := range dataCh {
			fmt.Fprintf(resp, "data: %s\n\n", data)
			resp.(http.Flusher).Flush()
		}
	}()

	x := 0
	for {
		dataCh <- strconv.Itoa(x)
		time.Sleep(1 * time.Second)
		x++
	}
}

type siteConfig struct {
	connectionCount int
}

func (sc siteConfig) handleNewConnection(resp http.ResponseWriter, req *http.Request) {

}
