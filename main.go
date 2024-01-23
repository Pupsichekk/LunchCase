package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *sql.DB
}

func main() {
	connectionString := "user=postgres password=Qqwerasdf1 dbname=lunchcase sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	apiCfg := apiConfig{
		DB: db,
	}

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to PostgreSQL succesfully")

	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Printf("No variable named port exists")
	}

	http.HandleFunc("/user", corsMiddleware(apiCfg.handlerUser))
	log.Printf("Server starting on port %s", port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Printf("Unable to start a server on port %s: %v", port, err)
	}

}
