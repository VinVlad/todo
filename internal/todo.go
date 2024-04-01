package internal

import (
	"errors"
	"fmt"
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

func NewStorage() *Storage {
	ts := &Storage{
		tasks:  make(map[int]Todo),
		nextID: 0}
	return ts
}

// CreateTodo создаёт новую туду-задачу в хранилище
func (ts *Storage) CreateTodo(title string, description string) int {
	ts.Lock()
	defer ts.Unlock()

	ts.nextID++

	task := Todo{
		Title:       title,
		Description: description}

	ts.tasks[ts.nextID] = task

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
