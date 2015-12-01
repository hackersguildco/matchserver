package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cheersapp/matchserver/models"
	"github.com/cheersapp/matchserver/utils"
	"github.com/cheersapp/matchserver/ws"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	initEnv()
	utils.InitLog()
	models.InitDB()
	go ws.System.Run()
	router := mux.NewRouter()
	router.HandleFunc("/ws/{username}", ws.ServeWS)
	http.Handle("/", router)
	port := os.Getenv("PORT")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func initEnv() {
	if err := godotenv.Load(".env_dev"); err != nil {
		log.Fatal("Error loading .env_dev file")
	}
}
