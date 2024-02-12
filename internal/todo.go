package internal

import (
	"errors"
	"sync"
)

type Todo struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoStore struct {
	sync.Mutex

	tasks  map[int]Todo
	nextID int
}

func New() *TodoStore {
	ts := &TodoStore{
		tasks:  make(map[int]Todo),
		nextID: 0}
	return ts
}

// CreateTodo создаёт новую туду-задачу в хранилище
func (ts *TodoStore) CreateTodo(title string, description string) int {
	ts.Lock()
	defer ts.Unlock()

	task := Todo{
		Id:          ts.nextID,
		Title:       title,
		Description: description}

	ts.tasks[ts.nextID] = task
	ts.nextID++
	return task.Id
}

// ChangeTodo изменяет заголовок и/или описание задачи
func (ts *TodoStore) ChangeTodo(id int, title string, description string) (Todo, error) {
	ts.Lock()
	defer ts.Unlock()

	v, ok := ts.tasks[id]
	if !ok {
		err := errors.New("Поиск таски в базе:")
		return v, err
	} else {
		task := Todo{
			Title:       title,
			Description: description}

		v = task
	}

	return v, nil

}

// GetList отдаёт список всех задач
func (ts *TodoStore) GetList() *TodoStore {
	return ts
}

// DeleteTask удаляет задачу по id
func (ts *TodoStore) DeleteTask(id int) *TodoStore {
	delete(ts.tasks, id)
	return ts

}
