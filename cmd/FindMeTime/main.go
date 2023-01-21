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
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Task struct {
	TaskId      string
	Title       string
	Description string
	Duration    int
	CreatedOn   string
	Frequency   int
}

type ProposedTask struct {
	*Task
	StartTime string
}

type FindTimeRequestTask struct {
	TaskId []string
}

type FindTimeRequest struct {
	Tasks []string
}

type FindTimeResponse struct {
	StartDate     string
	EndDate       string
	ProposedTasks []ProposedTask
	Week          Week
}

type Week struct {
	Days map[string]Day
}

type Day struct {
	SortedItems []ProposedTask
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var t Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userName := "tasker"
	host := "192.168.1.32"

	connStr := "postgresql://" + userName + ":s.o.a.d.@" + host + "/findmetime?sslmode=disable"
	fmt.Print("Before conn")
	fmt.Print(t)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("After conn")
	_, err = db.Query("INSERT INTO Tasks (task_id, title, description, duration, created_on, frequency) VALUES ($1, $2, $3, $4, $5);", uuid.New(), t.Title, t.Description, t.Duration, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Task: %+v", t)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Print(".")
	var tasks []Task

	userName := "tasker"
	host := "192.168.1.32"

	connStr := "postgresql://" + userName + ":s.o.a.d.@" + host + "/findmetime?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("select * from Tasks")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var t Task
		err = rows.Scan(&t.TaskId, &t.Title, &t.Description, &t.Duration, &t.CreatedOn)
		if err != nil {
			fmt.Print("Scan: %v", err)
		}
		tasks = append(tasks, t)
	}

	json.NewEncoder(w).Encode(tasks)
}

func FindTime(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Print("in find time")
	var findTimeRequest FindTimeRequest
	err := json.NewDecoder(r.Body).Decode(&findTimeRequest)
	if err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var taskIdList []string
	var goalIdList []string

	for _, v := range findTimeRequest.Tasks {
		taskIdList = append(taskIdList, v)
	}

	for _, v := range findTimeRequest.Goals {
		goalIdList = append(goalIdList, v)
	}

	var tasks []Task

	userName := "tasker"
	host := "192.168.1.32"

	connStr := "postgresql://" + userName + ":s.o.a.d.@" + host + "/findmetime?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("select * from Tasks where task_id = Any($1)", pq.Array(taskIdList))
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var t Task
		err = rows.Scan(&t.TaskId, &t.Title, &t.Description, &t.Duration, &t.CreatedOn, &t.Frequency)
		if err != nil {
			fmt.Print("Scan: %v", err)
		}
		tasks = append(tasks, t)
	}

	response := FindTimeWorker(tasks)

	json.NewEncoder(w).Encode(response)
}

func main() {
	router := httprouter.New()

	router.POST("/api/v1/task/create", CreateTaskHandler)
	router.POST("/api/v1/goal/create", CreateGoalHandler)
	router.GET("/api/v1/task/all", GetTasksHandler)
	router.POST("/api/v1/findtime", FindTime)
	handler := cors.Default().Handler(router)
	fmt.Print("Started....")
	http.ListenAndServe(":8080", handler)
}

// func enableCors(w *http.ResponseWriter) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "*")
// 	fmt.Print(".")
// }

/*
\c findmetime;
drop table tasks;
drop table goals;
CREATE TABLE tasks (
        task_id  VARCHAR ( 50 ) PRIMARY KEY,
        title VARCHAR ( 50 ) NOT NULL,
        description VARCHAR ( 50 ) NOT NULL,
        duration INT NOT NULL,
        created_on TIMESTAMP NOT NULL,
		frequency INT NOT NULL
);

CREATE TABLE goals (
        task_id  VARCHAR ( 50 ) PRIMARY KEY,
        title VARCHAR ( 50 ) NOT NULL,
        description VARCHAR ( 50 ) NOT NULL,
        duration INT NOT NULL,
        created_on TIMESTAMP NOT NULL,
		frequency INT NOT NULL
);

GRANT ALL ON TABLE tasks TO tasker;
GRANT ALL ON TABLE goals TO tasker;
*/
