package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/phanirithvij/todo-cli"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var usr, _ = user.Current()
var hdir = usr.HomeDir
var repositoryFilePath = filepath.Join(hdir, ".todo-cli.json")

func loadTasksFromRepositoryFile() (todos []*todo.Task, doneTodos []*todo.Task, latestTaskID int) {
	f, err := os.Open(repositoryFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return todos, doneTodos, latestTaskID
		}
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	var t []*todo.Task
	if err = json.NewDecoder(f).Decode(&t); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for _, v := range t {
		if v.IsDone {
			doneTodos = append(doneTodos, v)
			continue
		}
		todos = append(todos, v)

		if v.ID >= latestTaskID {
			latestTaskID = v.ID
		}
	}

	return todos, doneTodos, latestTaskID
}

func (m model) saveTasks() {
	f, err := os.OpenFile(repositoryFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Println(err)
		}
		f, err = os.Create(repositoryFilePath)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
	defer f.Close()
	if err := f.Truncate(0); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	todos := append(m.tasks, m.doneTasks...)
	data, _ := json.Marshal(todos)

	_, err = f.Write(data)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
