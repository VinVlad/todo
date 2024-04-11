package internal

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
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

var Ctx = context.Background()

func NewConnection() *pgx.ConnConfig {
	var config, _ = pgx.ParseConfig("")
	config.Host = "localhost"
	config.Port = 5432
	config.User = "postgres"
	config.Password = "0000"
	config.Database = "postgres"
	return config
}

var config *pgx.ConnConfig = NewConnection()

func NewStorage() *Storage {
	ts := &Storage{
		tasks: make(map[uuid.UUID]Todo)}
	return ts
}

// isExist проверяет есть ли запись в базе с таким uuid.
func isExist(ctx context.Context, id uuid.UUID) (bool, error) {

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return false, err
	}
	defer func() { _ = conn.Close(ctx) }()

	exist := conn.QueryRow(ctx, "Select EXISTS (SELECT 1 FROM tasks  WHERE id = $1) ;", id)
	var exists bool
	if err := exist.Scan(&exists); err != nil {
		return false, err
	}
	return exists, err

}

// SaveValues сохраняет запись в базу.
func SaveValues(
	ctx context.Context,
	id uuid.UUID,
	title string,
	description string,
) {

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer func() { _ = conn.Close(ctx) }()

	_, err = conn.Exec(ctx, "INSERT INTO tasks (id, title, description) VALUES ($1, $2, $3)", id, title, description)
	if err != nil {
		fmt.Println("Unable to insert data into database:", err)
		return
	}

}

// UpdateValues изменяет имеющуюся запись в базе по uuid.
func UpdateValues(
	ctx context.Context,
	id uuid.UUID,
	title string,
	description string,
) {

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer func() { _ = conn.Close(ctx) }()

	_, err = conn.Exec(ctx, "UPDATE tasks SET title = $1, description = $2 WHERE id = $3;", title, description, id)
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
func (ts *Storage) CreateTodo(
	title string,
	description string,
) uuid.UUID {
	ts.Lock()
	defer ts.Unlock()

	ctx, cancel := context.WithCancel(Ctx)
	defer cancel()

	//ts.nextID++
	//task := Todo{
	//	Title:       title,
	//	Description: description}
	//
	//ts.tasks[ts.nextID] = task
	id, _ := uuid.NewV4()
	SaveValues(ctx, id, title, description)

	return id
}

// TODO: расширить ошибки из функций save и update
// ChangeTodo изменяет заголовок и/или описание задачи.
func (ts *Storage) ChangeTodo(
	id uuid.UUID,
	title string,
	description string,
) (Todo, error) {
	ts.Lock()
	defer ts.Unlock()

	ctx, cancel := context.WithCancel(Ctx)
	defer cancel()

	v := Todo{
		//Id:          v.Id,
		Title:       title,
		Description: description}

	exists, err := isExist(ctx, id)
	if err != nil {
		return Todo{}, err
	}
	if exists {
		UpdateValues(ctx, id, title, description)
	} else {
		SaveValues(ctx, id, title, description)
	}

	return v, nil
}

// GetList отдаёт список всех задач.
func (ts *Storage) GetList() map[uuid.UUID]Todo {
	ctx, cancel := context.WithCancel(Ctx)
	defer cancel()

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		fmt.Println("Unable to connect to the database:", err)
	}
	defer func() { _ = conn.Close(ctx) }()

	rows, err := conn.Query(ctx, "SELECT id, title, description FROM tasks")
	if err != nil {
		fmt.Println("Error querying the database:", err)
	}
	defer rows.Close()

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
func (ts *Storage) DeleteTask(id uuid.UUID) (string, error) {
	ts.Lock()
	defer ts.Unlock()

	ctx, cancel := context.WithCancel(Ctx)
	defer cancel()

	//delete(ts.tasks, uuid)
	//fmt.Println(ts)

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return "", err
	}
	defer func() { _ = conn.Close(ctx) }()

	_, err = conn.Exec(ctx, "DELETE FROM tasks WHERE id = $1;", id)
	if err != nil {
		fmt.Println("Удаление записи:", err)
		return "", err
	}

	return "Карточка удалена", nil

}
