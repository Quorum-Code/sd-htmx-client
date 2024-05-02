package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/verify-login", verifyLoginHandler)
	http.HandleFunc("/create-account", createAccountHandler)
	http.HandleFunc("/logged-in", loggedInHandler)
	http.HandleFunc("/", handler)
	http.HandleFunc("/events", eventHandler)
	http.HandleFunc("/events/redirect", redirectHandler)

	// startFileServer()
	http.ListenAndServe(":8080", nil)
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

func authJs(resp http.ResponseWriter, req http.Request) {

}

func handler(resp http.ResponseWriter, req *http.Request) {
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
