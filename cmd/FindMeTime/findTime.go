package main

import (
	"math/rand"
)

func FindTimeWorker(tasks []CreateTask, goals []Goal) *FindTimeResponse {
	hoursInAday := [24]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
	blockers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 23}

	var timeSlots []int
	var requiredTime int
	var returnerTasks []ProposedTask
	var returnerGoals []ProposedGoal
	returnerWeek := Week{Days: make(map[string]Day)}

	for _, hidElement := range hoursInAday {
		contained := false
		for _, blockerElement := range blockers {
			if hidElement == blockerElement {
				contained = true
			}
		}
		if !contained {
			timeSlots = append(timeSlots, hidElement)
		}
	}

	for _, task := range tasks {
		requiredTime += task.Duration
	}

	for _, goal := range goals {
		requiredTime += goal.Duration
	}

	if requiredTime > len(timeSlots) {
		return nil
	}

	for _, task := range tasks {
		slotPosition := rand.Intn(len(timeSlots))
		time := timeSlots[slotPosition]
		returnerTasks = append(returnerTasks, ProposedTask{&task, time})
		//Week.Days[""].SortedItems[].StartTime
		/*	TaskId      string
			Title       string
			Description string
			Duration    int
			CreatedOn   string
		*/
		if returnerWeek.Days["03-01-2023"].SortedItems == nil {
			items := []ProposedTask{{&task, time}}
			day := Day{SortedItems: items}
			returnerWeek.Days["03-01-2023"] = day
		} else {
			items := []ProposedTask{{&task, time}}
			returnerWeek.Days["03-01-2023"] = Day{SortedItems: items}
		}
	}

	for _, goal := range goals {
		slotPosition := rand.Intn(len(timeSlots))
		time := timeSlots[slotPosition]
		returnerGoals = append(returnerGoals, ProposedGoal{&goal, time})

	}

	return &FindTimeResponse{ProposedTasks: returnerTasks, ProposedGoals: returnerGoals, Week: returnerWeek}
}
