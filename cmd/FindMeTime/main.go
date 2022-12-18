package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"database/sql"

	"github.com/julienschmidt/httprouter"
)

type CreateTask struct {
	title    string
	duration int
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var t CreateTask
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	connStr := "postgresql://<username>:<password>@<database_ip>/todos?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Query("INSERT INTO Tasks (title, duration) VALUES ($1, $2);", t.title, t.duration)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Task: %+v", t)
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	log.Fatal(http.ListenAndServe(":8080", router))
}
