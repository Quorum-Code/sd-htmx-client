package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handler)
	http.HandleFunc("/events", eventHandler)

	// startFileServer()
	http.ListenAndServe(":8080", nil)
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
