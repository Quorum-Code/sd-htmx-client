package endpoints

import (
	"fmt"
	"net/http"

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

	// http.HandleFunc("/verify-login", verifyLoginHandler)
	// http.HandleFunc("/create-account", createAccountHandler)
	// http.HandleFunc("/logged-in", loggedInHandler)

	http.HandleFunc("/events", wsc.eventHandler)
	http.HandleFunc("/events/redirect", wsc.redirectHandler)

	http.HandleFunc("POST /api/sign-up", wsc.postSignupHandler)
	http.HandleFunc("POST /api/login", wsc.postLoginHandler)

	fmt.Println("Web server live...")

	http.ListenAndServe(":8080", nil)
}
