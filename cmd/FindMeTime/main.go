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

type CreateTask struct {
	TaskId      string
	Title       string
	Description string
	Duration    int
	CreatedOn   string
}

type ProposedTask struct {
	*CreateTask
	StartTime string
}

type Goal struct {
	*CreateTask
	Frequency int
}

type ProposedGoal struct {
	*Goal
	StartTime int
}

type FindTimeRequestTask struct {
	TaskId []string
}

type FindTimeRequest struct {
	Tasks []string
	Goals []string
}

type FindTimeResponse struct {
	StartDate     string
	EndDate       string
	ProposedTasks []ProposedTask
	ProposedGoals []ProposedGoal
	Week          Week
}

type Week struct {
	Days map[string]Day
}

type Day struct {
	SortedItems []ProposedTask
}

type Tag struct {
	/*
			task_id  VARCHAR ( 50 ) PRIMARY KEY,
		tag_name VARCHAR(50) NOT NULL,
		description VARCHAR(50) NOT NULL,
		owner  INT NOT NULL,
		mon_start TIMESTAMP NOT NULL,
		mon_end TIMESTAMP NOT NULL,
			tue_start TIMESTAMP NOT NULL,
		tue_end TIMESTAMP NOT NULL,
			wed_start TIMESTAMP NOT NULL,
		wed_end TIMESTAMP NOT NULL,
			thu_start TIMESTAMP NOT NULL,
		thu_end TIMESTAMP NOT NULL,
			fri_start TIMESTAMP NOT NULL,
		fri_end TIMESTAMP NOT NULL,
			sat_start TIMESTAMP NOT NULL,
		sat_end TIMESTAMP NOT NULL,
			sun_start TIMESTAMP NOT NULL,
		sun_end TIMESTAMP NOT NULL,
	*/
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var t CreateTask
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
	_, err = db.Query("INSERT INTO Tasks (task_id, title, description, duration, created_on) VALUES ($1, $2, $3, $4, $5);", uuid.New(), t.Title, t.Description, t.Duration, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Task: %+v", t)
}

func CreateGoalHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var g Goal
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userName := "tasker"
	host := "192.168.1.32"

	connStr := "postgresql://" + userName + ":s.o.a.d.@" + host + "/findmetime?sslmode=disable"
	fmt.Print("Before conn")
	fmt.Print(g)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("After conn")
	_, err = db.Query("INSERT INTO Goals (task_id, title, description, duration, created_on, frequency) VALUES ($1, $2, $3, $4, $5, $6);", uuid.New(), g.Title, g.Description, g.Duration, time.Now(), g.Frequency)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Task: %+v", g)
}

func CreateTagHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// var g Goal
	// err := json.NewDecoder(r.Body).Decode(&g)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// userName := "tasker"
	// host := "192.168.1.32"

	// connStr := "postgresql://" + userName + ":s.o.a.d.@" + host + "/findmetime?sslmode=disable"
	// fmt.Print("Before conn")
	// fmt.Print(g)
	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Print("After conn")
	// _, err = db.Query("INSERT INTO Goals (task_id, title, description, duration, created_on, frequency) VALUES ($1, $2, $3, $4, $5, $6);", uuid.New(), g.Title, g.Description, g.Duration, time.Now(), g.Frequency)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Fprintf(w, "Task: %+v", g)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Print(".")
	var tasks []CreateTask

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
		var t CreateTask
		err = rows.Scan(&t.TaskId, &t.Title, &t.Description, &t.Duration, &t.CreatedOn)
		if err != nil {
			fmt.Print("Scan: %v", err)
		}
		tasks = append(tasks, t)
	}

	goaldb, err := sql.Open("postgres", connStr)
	rows, err = goaldb.Query("select task_id, title, description, duration, created_on from Goals")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var t CreateTask
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

	var tasks []CreateTask

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
		var t CreateTask
		err = rows.Scan(&t.TaskId, &t.Title, &t.Description, &t.Duration, &t.CreatedOn)
		if err != nil {
			fmt.Print("Scan: %v", err)
		}
		tasks = append(tasks, t)
	}

	var goals []Goal
	goaldb, goalerr := sql.Open("postgres", connStr)
	goalrows, goalerr := goaldb.Query("select * from Goals where task_id = Any($1)", pq.Array(taskIdList))
	if goalerr != nil {
		log.Fatal(err)
	}
	for goalrows.Next() {
		var g Goal
		var t CreateTask
		err = goalrows.Scan(&t.TaskId, &t.Title, &t.Description, &t.Duration, &t.CreatedOn, &g.Frequency)
		if err != nil {
			fmt.Print("Scan: %v", err)
		}
		g.CreateTask = &t
		goals = append(goals, g)
	}

	response := FindTimeWorker(tasks, goals)

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
        created_on TIMESTAMP NOT NULL
);

CREATE TABLE goals (
        task_id  VARCHAR ( 50 ) PRIMARY KEY,
        title VARCHAR ( 50 ) NOT NULL,
        description VARCHAR ( 50 ) NOT NULL,
        duration INT NOT NULL,
        created_on TIMESTAMP NOT NULL,
		frequency INT NOT NULL
);

CREATE TABLE tags (
	task_id  VARCHAR ( 50 ) PRIMARY KEY,
	tag_name VARCHAR(50) NOT NULL,
	description VARCHAR(50) NOT NULL,
	owner  INT NOT NULL,
	mon_start TIMESTAMP NOT NULL,
	mon_end TIMESTAMP NOT NULL,
		tue_start TIMESTAMP NOT NULL,
	tue_end TIMESTAMP NOT NULL,
		wed_start TIMESTAMP NOT NULL,
	wed_end TIMESTAMP NOT NULL,
		thu_start TIMESTAMP NOT NULL,
	thu_end TIMESTAMP NOT NULL,
		fri_start TIMESTAMP NOT NULL,
	fri_end TIMESTAMP NOT NULL,
		sat_start TIMESTAMP NOT NULL,
	sat_end TIMESTAMP NOT NULL,
		sun_start TIMESTAMP NOT NULL,
	sun_end TIMESTAMP NOT NULL,

)

GRANT ALL ON TABLE tasks TO tasker;
GRANT ALL ON TABLE goals TO tasker;
GRANT ALL ON TABLE tags TO tasker;
*/
