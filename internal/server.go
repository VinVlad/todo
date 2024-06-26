package internal

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"log"
	"mime"
	"net/http"
	"time"
)

type Server struct {
	HttpServer *http.Server
	Storage    *Storage
}

func NewServer() *Server {
	store := NewStorage()
	return &Server{Storage: store}
}

func (s *Server) Run(port string) {
	s.HttpServer = &http.Server{
		Addr:           "localhost:" + port,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
}

// renderJSON формирует JSON.
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// CreateTaskHandler хэндлер для создания новой карточки.
func (s *Server) CreateTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling task create at %s\n", req.URL.Path)

	type RequestTask struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	type ResponseId struct {
		Id uuid.UUID `json:"id"`
	}

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rt RequestTask
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := s.Storage.CreateTodo(rt.Title, rt.Description)
	renderJSON(w, ResponseId{Id: id})
}

// ChangeTaskHandler хэндлер для обновления существующей карточки.
func (s *Server) ChangeTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling task change at %s\n", req.URL.Path)

	type RequestTask struct {
		Id          uuid.UUID `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
	}

	type ResponseId struct {
		Id          uuid.UUID `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
	}

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rt RequestTask
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	td, err := s.Storage.ChangeTodo(rt.Id, rt.Title, rt.Description)
	if err != nil {
		fmt.Println("ошибка изменения карточки:")
	}

	renderJSON(w, ResponseId{Id: rt.Id, Title: td.Title, Description: td.Description})
}

// GetTasksListHandler хэндлер для получения всех карточке.
func (s *Server) GetTasksListHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all tasks  at %s\n", req.URL.Path)

	type ResponseId struct {
		Id          uuid.UUID `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
	}

	td := s.Storage.GetList()

	renderJSON(w, td)
}

// DeleteTaskHandler хэндлер для удаления карточки по uuid.
func (s *Server) DeleteTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete task  at %s\n", req.URL.Path)

	type RequestTask struct {
		Id uuid.UUID `json:"id"`
	}

	type ResponseId struct {
		Message string `json:"message"`
	}

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rt RequestTask
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := s.Storage.DeleteTask(rt.Id)

	renderJSON(w, ResponseId{Message: message})

	//id, _ := strconv.Atoi(mux.Vars(req)["id"])
	//s.Storage.DeleteTask(uuid.UUID)
}
