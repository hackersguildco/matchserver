package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cheersappio/matchserver/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/cheersappio/matchserver/Godeps/_workspace/src/github.com/joho/godotenv"
	_ "github.com/cheersappio/matchserver/models" // init
	_ "github.com/cheersappio/matchserver/utils"
	"github.com/cheersappio/matchserver/ws"
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
