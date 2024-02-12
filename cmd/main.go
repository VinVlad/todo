package main

import (
	"fmt"
	"todo/internal"
)

func main() {
	//server := new(internal.Server)
	//err := server.Run("8000")
	//if err != nil {
	//	log.Fatalf("запуск сервера: %s", err.Error())
	//}

	c := internal.New()
	c.CreateTodo("Задача 1", "Описание 1")
	c.CreateTodo("Задача 2", "Описание 2")
	c.CreateTodo("Задача 3", "Описание 3")
	fmt.Println(c.GetList())

	c.DeleteTask(1)
	fmt.Println(c.GetList())
}
