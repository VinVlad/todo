package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"todo/internal"
)

func main() {

	router := mux.NewRouter()

	server := internal.NewServer()

	router.HandleFunc("/task", server.CreateTaskHandler).Methods(http.MethodPut)
	router.HandleFunc("/task", server.ChangeTaskHandler).Methods(http.MethodPost)
	router.HandleFunc("/task", server.GetTasksListHandler).Methods(http.MethodGet)
	router.HandleFunc("/task", server.DeleteTaskHandler).Methods(http.MethodDelete)

	server.Run("8000")

	log.Fatal(http.ListenAndServe(server.HttpServer.Addr, router))

}
