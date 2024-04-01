package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"todo/internal"
)

func main() {

	router := mux.NewRouter()
	//todo не работает игнор слэша
	router.StrictSlash(true)

	server := internal.NewServer()

	router.HandleFunc("/task/", server.CreateTaskHandler).Methods("POST")
	router.HandleFunc("/task/change/", server.ChangeTaskHandler).Methods("POST")
	router.HandleFunc("/task/all", server.GetTasksListHandler).Methods("GET")
	router.HandleFunc("/task/{id:[0-9]+}/", server.DeleteTaskHandler).Methods("DELETE")

	server.Run("8000")
	log.Fatal(http.ListenAndServe(server.HttpServer.Addr, router))

}
