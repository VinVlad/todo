package internal

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	"sync"
)

type Todo struct {
	//Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Storage struct {
	sync.Mutex

	tasks  map[int]Todo
	nextID int
}

var Conn = pgx.ConnConfig{
	Host:     "localhost",
	Port:     5432,
	Database: "postgres",
	User:     "postgres",
	Password: "0000",
}

func NewStorage() *Storage {
	ts := &Storage{
		tasks:  make(map[int]Todo),
		nextID: 0}
	return ts
}

func SaveValues(id int, title string, description string) {

	conn, err := pgx.Connect(Conn)
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Exec("INSERT INTO tasks (id, title, description) VALUES ($1, $2, $3)", id, title, description)
	if err != nil {
		fmt.Println("Unable to insert data into database:", err)
		return
	}

}

func ReadValues() (*pgx.Rows, error) {

	//return rows, err
	return nil, nil
}

// CreateTodo создаёт новую туду-задачу в хранилище
func (ts *Storage) CreateTodo(title string, description string) int {
	ts.Lock()
	defer ts.Unlock()

	ts.nextID++

	//task := Todo{
	//	Title:       title,
	//	Description: description}
	//
	//ts.tasks[ts.nextID] = task

	SaveValues(ts.nextID, title, description)

	return ts.nextID
}

// ChangeTodo изменяет заголовок и/или описание задачи
func (ts *Storage) ChangeTodo(id int, title string, description string) (Todo, error) {
	ts.Lock()
	defer ts.Unlock()

	v, ok := ts.tasks[id]
	if !ok {
		err := errors.New("Поиск таски в базе:")
		return v, err
	} else {
		v = Todo{
			//Id:          v.Id,
			Title:       title,
			Description: description}
		ts.tasks[id] = v
	}

	return v, nil

}

// GetList отдаёт список всех задач
func (ts *Storage) GetList() map[int]Todo {

	conn, err := pgx.Connect(Conn)
	if err != nil {
		fmt.Println("Unable to connect to the database:", err)
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT id, title, description FROM tasks")
	if err != nil {
		fmt.Println("Error querying the database:", err)
	}
	defer rows.Close()

	//rows, _ := ReadValues()

	for rows.Next() {
		var id int
		var title string
		var description string
		err := rows.Scan(&id, &title, &description)
		ts.tasks[id] = Todo{
			Title:       title,
			Description: description,
		}
		if err != nil {
			fmt.Println("Error scanning row:", err)
		}

	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
	}

	return ts.tasks
}

// DeleteTask удаляет задачу по id
func (ts *Storage) DeleteTask(id int) {
	ts.Lock()
	defer ts.Unlock()

	delete(ts.tasks, id)
	fmt.Println(ts)
	return

}
