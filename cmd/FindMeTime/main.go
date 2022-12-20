package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"database/sql"

	"github.com/google/uuid"
	"github.com/rs/cors" // only needed while CORS is in play

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

type CreateTask struct {
	title       string
	description string
	duration    string
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Print(".")
	var t CreateTask
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userName := "tasker"
	host := "192.168.1.33"

	connStr := "postgresql://" + userName + ":s.o.a.d.@" + host + "/findmetime?sslmode=disable"
	fmt.Print("Before conn")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("After conn")
	_, err = db.Query("INSERT INTO Tasks (task_id, title, description, duration, created_on) VALUES ($1, $2, $3, $4, $5);", uuid.New(), t.title, t.description, t.duration, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Task: %+v", t)
}

func main() {
	router := httprouter.New()

	router.POST("/api/v1/task/create", CreateTaskHandler)
	handler := cors.Default().Handler(router)
	http.ListenAndServe(":8080", handler)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

/*
\c findmetime;
drop table tasks;
CREATE TABLE tasks (
        task_id  VARCHAR ( 50 ) PRIMARY KEY,
        title VARCHAR ( 50 ) NOT NULL,
        description VARCHAR ( 50 ) NOT NULL,
        duration VARCHAR ( 255 ) NOT NULL,
        created_on TIMESTAMP NOT NULL
); GRANT ALL ON TABLE tasks TO tasker;
*/
