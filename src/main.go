package main

import (
	"SCTI/database"
	"SCTI/fileserver"
	"SCTI/middleware"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = database.OpenDatabase()
	if err != nil {
		log.Printf("Error connecting to postgres database\n%v", err)
	}
	defer database.CloseDatabase()

	fileserver.RunFileServer()

	fmt.Println("Server Started")
	mux := http.NewServeMux()
	LoadRoutes(mux)

	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.EndpointLogging(mux),
	}

	log.Fatal(server.ListenAndServe())
}
