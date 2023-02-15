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

type User struct {
	ID       string
	UserName string
}

type CreateTask struct {
	TaskId      string
	Title       string
	Description string
	Duration    int
	CreatedOn   string
	Owner       User
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
	Id          string
	Name        string
	Description string
	Owner       User
	Mon_start   string
	Mon_end     string
	Tues_start  string
	Tues_end    string
	Wed_start   string
	Wed_end     string
	Thur_start  string
	Thur_end    string
	Fri_start   string
	Fri_end     string
	Sat_start   string
	Sat_end     string
	Sun_start   string
	Sun_end     string
}

func openDB() (*sql.DB, error) {
	userName := "tasker"
	host := "192.168.1.32"
	pass := "s.o.a.d."
	database := "findmetime"

	connStr := "postgresql://" + userName + ":" + pass + "@" + host + "/" + database + "?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db, err
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var t CreateTask
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db, err := openDB()
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

	db, err := openDB()
	_, err = db.Query("INSERT INTO Goals (task_id, title, description, duration, created_on, frequency) VALUES ($1, $2, $3, $4, $5, $6);", uuid.New(), g.Title, g.Description, g.Duration, time.Now(), g.Frequency)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Task: %+v", g)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Print(".")
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	db, err := openDB()
	_, err = db.Query("INSERT INTO USERS (id, username) VALUES ($1, $2)", user.ID, user.UserName)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTagHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var t Tag
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db, err := openDB()
	_, err = db.Query("INSERT INTO Tags (task_id, tag_name, description, owner, mon_start, mon_end, tue_start, tue_end, wed_start, wed_end, thu_start, thu_end, fri_start, fri_end, sat_start, sat_end, sun_start, sun_end) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18);", uuid.New(), t.Name, t.Description, t.Owner.ID, t.Mon_start, t.Mon_end, t.Tues_start, t.Tues_end, t.Wed_start, t.Wed_end, t.Thur_start, t.Thur_end, t.Fri_start, t.Fri_end, t.Sat_start, t.Sat_end, t.Sun_start, t.Sun_end)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Task: %+v", t)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var tasks []CreateTask

	db, err := openDB()
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

	goaldb, err := openDB()
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

	db, err := openDB()
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
	goaldb, goalerr := openDB()
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
	router.POST("/api/v1/users/", CreateUserHandler)
	router.POST("/api/v1/tags/", CreateTagHandler)
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
drop table tags;
drop table users;

CREATE TABLE users (
	id  VARCHAR (20) PRIMARY KEY,
	username VARCHAR(20) NOT NULL
);

CREATE TABLE tasks (
	task_id  VARCHAR (20) PRIMARY KEY,
	title VARCHAR (20) NOT NULL,
	description VARCHAR (20) NOT NULL,
	duration INT NOT NULL,
	created_on TIMESTAMP NOT NULL,
	owner  VARCHAR (20) NOT NULL
);

CREATE TABLE goals (
	task_id  VARCHAR (20) PRIMARY KEY,
	title VARCHAR (20) NOT NULL,
	description VARCHAR (20) NOT NULL,
	duration INT NOT NULL,
	created_on TIMESTAMP NOT NULL,
	frequency INT NOT NULL,
	owner  VARCHAR (20) NOT NULL
);

CREATE TABLE tags (
	task_id  VARCHAR (20) PRIMARY KEY,
	tag_name VARCHAR(20) NOT NULL,
	description VARCHAR(20) NOT NULL,
	owner  VARCHAR (20) NOT NULL,
	mon_start TIMESTAMP ,
	mon_end TIMESTAMP ,
		tue_start TIMESTAMP ,
	tue_end TIMESTAMP ,
		wed_start TIMESTAMP ,
	wed_end TIMESTAMP ,
		thu_start TIMESTAMP ,
	thu_end TIMESTAMP ,
		fri_start TIMESTAMP ,
	fri_end TIMESTAMP ,
		sat_start TIMESTAMP ,
	sat_end TIMESTAMP ,
		sun_start TIMESTAMP ,
	sun_end TIMESTAMP

);
GRANT ALL ON TABLE users TO tasker;
GRANT ALL ON TABLE tasks TO tasker;
GRANT ALL ON TABLE goals TO tasker;
GRANT ALL ON TABLE tags TO tasker;
*/
