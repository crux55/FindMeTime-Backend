package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"database/sql"

	"github.com/google/uuid"
	"github.com/rs/cors" // only needed while CORS is in play

	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
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
	TagsOnly    []Tag
	TagsNot     []Tag
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
	TimeSlots   []TimeSlot
}

type TimeSlot struct {
	DayIndex  int
	StartTime int
	EndTime   int
}

func openDB() (*sql.DB, error) {
	loadConfig, _ := LoadConfig(os.Args[1])
	config := loadConfig.DatabaseConfig
	connStr := "postgresql://" + config.Username + ":" + config.Password + "@" + config.Host + "/" + config.DatabaseName + "?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Print(err)
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
	var tagOnlyIds []string
	var tagNotIds []string
	for _, tagOnly := range t.TagsOnly {
		tagOnlyIds = append(tagOnlyIds, tagOnly.Id)
	}

	for _, tagNot := range t.TagsNot {
		tagNotIds = append(tagNotIds, tagNot.Id)
	}

	db, err := openDB()
	_, err = db.Query("INSERT INTO Tasks (id, title, description, duration, created_on, tags_only, tags_not) VALUES ($1, $2, $3, $4, $5, $6, $7);", uuid.New(), t.Title, t.Description, t.Duration, time.Now(), pq.Array(tagOnlyIds), pq.Array(tagNotIds))
	if err != nil {
		fmt.Print(err)
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
		fmt.Print(err)
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
		fmt.Print(err)
	}
}

func CreateTagHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Print("creating tag")
	var t Tag
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db, err := openDB()
	var timeSlotIds []string
	for _, timeSlot := range t.TimeSlots {
		id := uuid.New()
		_, err = db.Query("INSERT INTO time_slots (id, day_index, start_time, end_time) VALUES ($1, $2, $3, $4);", id, timeSlot.DayIndex, timeSlot.StartTime, timeSlot.EndTime)
		if err != nil {
			fmt.Print(err)
		}
		timeSlotIds = append(timeSlotIds, id.String())
	}
	_, err = db.Query("INSERT INTO tags (id, tag_name, description, time_slots) VALUES ($1, $2, $3, $4);", uuid.New(), t.Name, t.Description, pq.Array(&timeSlotIds))
	if err != nil {
		fmt.Print(err)
	}
	fmt.Fprintf(w, "Task: %+v", t)
}

func GetTagsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var tags []Tag

	db, err := openDB()
	rows, err := db.Query("select * from Tags")
	if err != nil {
		fmt.Print(err)
	}
	for rows.Next() {
		var t Tag
		err = rows.Scan(&t.Id, &t.Name, &t.Description, &t.TimeSlots)
		if err != nil {
			fmt.Print("Scan: %v", err)
		}
		tags = append(tags, t)
	}
	json.NewEncoder(w).Encode(tags)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var tasks []CreateTask

	db, err := openDB()
	rows, err := db.Query("select * from Tasks")
	if err != nil {
		fmt.Print(err)
	}
	for rows.Next() {
		var t CreateTask
		err = rows.Scan(&t.TaskId, &t.Title, &t.Description, &t.Duration, &t.CreatedOn, &t.TagsOnly, &t.TagsNot)
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

	for _, v := range findTimeRequest.Tasks {
		taskIdList = append(taskIdList, v)
	}

	var tasks []CreateTask

	db, err := openDB()
	if err != nil {
		fmt.Print(err)
	}
	rows, err := db.Query("select * from tasks where id = Any($1)", pq.Array(taskIdList))
	if err != nil {
		fmt.Print(err)
	}
	for rows.Next() {
		var t CreateTask
		var to string
		var tn string
		// var toIds []string
		// var tnIds []string
		err = rows.Scan(&t.TaskId, &t.Title, &t.Description, &t.Duration, &t.CreatedOn, &to, &tn)
		// if err != nil {
		// 	fmt.Print("Scan: %v", err)
		// }
		// if len(to) > 0 {
		// 	toIds = strings.Split(to[1:len(to)-1], ",")
		// } else {
		// 	toIds = []string{}
		// }
		// if len(tn) > 0 {
		// 	tnIds = strings.Split(tn[1:len(tn)-1], ",")
		// } else {
		// 	tnIds = []string{}
		// }

		t.TagsOnly = resolveTags(to, db)
		t.TagsNot = resolveTags(tn, db)

		// for _, tagId := range toIds {
		// 	tagId = strings.ReplaceAll(tagId, "\"", "")
		// 	var timeSlots []TimeSlot
		// 	fmt.Println("looping toIds", tagId)

		// 	getTimeSlotIdQuers, _ := db.Prepare("select time_slots from tags where id = $1;")
		// 	timeSlotIds, err := getTimeSlotIdQuers.Query(tagId)

		// 	var timeSlotIdList []string
		// 	for timeSlotIds.Next() {
		// 		var timeSlotId string
		// 		err = timeSlotIds.Scan(&timeSlotId)
		// 		var ts string
		// 		ts = strings.ReplaceAll(timeSlotId, "{", "")
		// 		ts = strings.ReplaceAll(ts, "}", "")
		// 		ts = strings.ReplaceAll(ts, "\"", "")
		// 		timeSlotIdList = append(timeSlotIdList, ts)
		// 		if err != nil {
		// 			fmt.Print("Scan: %v", err)
		// 		}
		// 	}

		// 	timeSlotQuery, _ := db.Prepare("select day_index, start_time, end_time from time_slots where id = Any($1);")
		// 	timeslotrows, err := timeSlotQuery.Query(pq.Array(strings.Split(timeSlotIdList[0], ",")))
		// 	if err != nil {
		// 		fmt.Print("Scan: %v", err)
		// 	}
		// 	for timeslotrows.Next() {
		// 		var timeSlot TimeSlot
		// 		var day_index int
		// 		var start_time int
		// 		var end_time int
		// 		err = timeslotrows.Scan(&day_index, &start_time, &end_time)
		// 		timeSlot.DayIndex = day_index
		// 		timeSlot.StartTime = start_time
		// 		timeSlot.EndTime = end_time
		// 		timeSlots = append(timeSlots, timeSlot)
		// 	}
		// 	t.TagsOnly = []Tag{Tag{TimeSlots: timeSlots}}
		// }

		// for _, tagId := range tnIds {
		// 	var timeSlots []TimeSlot
		// 	timeslotrows, err := db.Query("select * from time_slots where id = Any(select time_slots from tags where id = $1);", string(tagId))
		// 	if err != nil {
		// 		fmt.Print("Scan: %v", err)
		// 	}
		// 	for timeslotrows.Next() {
		// 		var timeSlot TimeSlot
		// 		err = timeslotrows.Scan(&timeSlot.DayIndex, &timeSlot.StartTime, &timeSlot.EndTime)
		// 		timeSlots = append(timeSlots, timeSlot)
		// 	}
		// 	t.TagsNot = []Tag{Tag{TimeSlots: timeSlots}}
		// }

		tasks = append(tasks, t)
	}

	fmt.Print("tasks", tasks)
	response := FindTimeWorker(tasks)

	json.NewEncoder(w).Encode(response)
}

func stripArrayChars(str string) string {
	str = strings.ReplaceAll(str, "{", "")
	str = strings.ReplaceAll(str, "}", "")
	str = strings.ReplaceAll(str, "\"", "")
	return str
}

func resolveTags(tagIdStr string, db *sql.DB) []Tag {
	if len(tagIdStr) == 0 {
		return []Tag{}
	}
	var timeSlots []TimeSlot
	tagIds := strings.Split(stripArrayChars(tagIdStr), ",")
	for _, tagId := range tagIds {
		fmt.Println("looping tag ids", tagId)

		getTimeSlotIdQuers, _ := db.Prepare("select time_slots from tags where id = $1;")
		timeSlotIds, err := getTimeSlotIdQuers.Query(tagId)

		var timeSlotIdList []string
		for timeSlotIds.Next() {
			var timeSlotId string
			err = timeSlotIds.Scan(&timeSlotId)
			timeSlotIdList = append(timeSlotIdList, stripArrayChars(timeSlotId)) //pretty sure this isn't needed
			if err != nil {
				fmt.Print("Scan: %v", err)
			}
		}

		timeSlotQuery, _ := db.Prepare("select day_index, start_time, end_time from time_slots where id = Any($1);")
		timeslotrows, err := timeSlotQuery.Query(pq.Array(strings.Split(timeSlotIdList[0], ",")))
		if err != nil {
			fmt.Print("Scan: %v", err)
		}
		for timeslotrows.Next() {
			var timeSlot TimeSlot
			var day_index int
			var start_time int
			var end_time int
			err = timeslotrows.Scan(&day_index, &start_time, &end_time)
			timeSlot.DayIndex = day_index
			timeSlot.StartTime = start_time
			timeSlot.EndTime = end_time
			timeSlots = append(timeSlots, timeSlot)
		}
	}
	return []Tag{Tag{TimeSlots: timeSlots}}
}

func main() {
	router := httprouter.New()

	router.POST("/api/v1/task/create", CreateTaskHandler)
	router.POST("/api/v1/goal/create", CreateGoalHandler)
	router.GET("/api/v1/task/all", GetTasksHandler)
	router.POST("/api/v1/findtime", FindTime)
	router.POST("/api/v1/users", CreateUserHandler)
	router.POST("/api/v1/tags", CreateTagHandler)
	router.GET("/api/v1/tags", GetTagsHandler)
	handler := cors.Default().Handler(router)
	fmt.Print("Started....")
	http.ListenAndServe(":8080", handler)
}
