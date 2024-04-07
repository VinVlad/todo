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

	router.HandleFunc("/task/", server.CreateTaskHandler).Methods("PUT")
	router.HandleFunc("/task/", server.ChangeTaskHandler).Methods("POST")
	router.HandleFunc("/task/", server.GetTasksListHandler).Methods("GET")
	router.HandleFunc("/task/", server.DeleteTaskHandler).Methods("DELETE")

	//fmt.Println(internal.Config.Config)

	server.Run("8000")

	log.Fatal(http.ListenAndServe(server.HttpServer.Addr, router))

}
