package main

import (
	"gojek/library/database"
	"gojek/library/internal/router"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := router.NewRouter(db)
	log.Fatal(http.ListenAndServe(":8080", r))
}
