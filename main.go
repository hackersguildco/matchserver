package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/hackersguildco/matchserver/models" // init
	_ "github.com/hackersguildco/matchserver/utils"
	"github.com/hackersguildco/matchserver/ws"
	"github.com/joho/godotenv"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ws/{username}", ws.ServeWS)
	http.Handle("/", router)
	port := os.Getenv("PORT")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func init() {
	if err := godotenv.Load(".env_dev"); err != nil {
		log.Fatal("Error loading .env_dev file")
	}
}
