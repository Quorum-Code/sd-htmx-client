package endpoints

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/Quorum-Code/sd-htmx-client/internal/database"
)

type WSConfig struct {
	Database *database.Database
}

func StartWebServer() {
	fmt.Println("Starting web server...")

	wsc := WSConfig{}
	db, err := database.InitDatabase()
	if err != nil {
		fmt.Println("failed to start db...", err)
		return
	}
	wsc.Database = db

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", wsc.indexHandler)
	http.HandleFunc("/sign-up", wsc.signupHandler)
	http.HandleFunc("/login", wsc.loginHandler)

	http.HandleFunc("/events", wsc.eventHandler)
	http.HandleFunc("/events/redirect", wsc.redirectHandler)

	http.HandleFunc("POST /api/index", wsc.PostIndexHandler)
	http.HandleFunc("POST /api/sign-up", wsc.postSignupHandler)
	http.HandleFunc("POST /api/login", wsc.postLoginHandler)

	fmt.Println("Starting web server...")

	port := ":8080"

	printLocalHost(port)

	http.ListenAndServe(port, nil)

	fmt.Println("Closing web server...")
}

func printLocalHost(port string) {
	name, err := os.Hostname()
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return
	}

	fmt.Println("Possible connections...")

	for _, a := range addrs {
		fmt.Printf(" - %s%s\n", a, port)
	}

	fmt.Printf(" - %s%s\n", "localhost", port)
}
