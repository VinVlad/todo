package internal

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
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

	tasks map[uuid.UUID]Todo
	//nextID int
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
		tasks: make(map[uuid.UUID]Todo)}
	return ts
}

func SaveValues(id uuid.UUID, title string, description string) {

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

func UpdateValues(id uuid.UUID, title string, description string) {

	conn, err := pgx.Connect(Conn)
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Exec("UPDATE tasks SET title = $1, description = $2 WHERE id = $3;", title, description, id)
	if err != nil {
		fmt.Println("Unable to insert data into database:", err)
		return
	}

}

func ReadValues() (*pgx.Rows, error) {

	//return rows, err
	return nil, nil
}

// CreateTodo создаёт новую туду-задачу в хранилище.
func (ts *Storage) CreateTodo(title string, description string) uuid.UUID {
	ts.Lock()
	defer ts.Unlock()

	//ts.nextID++

	//task := Todo{
	//	Title:       title,
	//	Description: description}
	//
	//ts.tasks[ts.nextID] = task
	ID, _ := uuid.NewV4()
	SaveValues(ID, title, description)

	return ID
}

// todo: поправить обращение по ключам. Если в базе записи нет, то всё-равно что-то вернётся)))))
// ChangeTodo изменяет заголовок и/или описание задачи.
func (ts *Storage) ChangeTodo(uuid uuid.UUID, title string, description string) (Todo, error) {
	ts.Lock()
	defer ts.Unlock()

	v := Todo{
		//Id:          v.Id,
		Title:       title,
		Description: description}

	UpdateValues(uuid, title, description)
	return v, nil

}

// GetList отдаёт список всех задач.
func (ts *Storage) GetList() map[uuid.UUID]Todo {

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
		var id uuid.UUID
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

// DeleteTask удаляет задачу по id.
func (ts *Storage) DeleteTask(uuid uuid.UUID) {
	ts.Lock()
	defer ts.Unlock()

	delete(ts.tasks, uuid)
	fmt.Println(ts)
	return

}
