package main

import (
	"github.com/joho/godotenv"

	"github.com/Quorum-Code/sd-htmx-client/internal/endpoints"
)

func main() {
	loadEnv()
	endpoints.StartWebServer()
}

func loadEnv() {
	godotenv.Load(".env")
}
