package main

import (
	"math/rand"
	"sort"
	"strconv"
	"time"
)

func FindTimeWorker(tasks []CreateTask, goals []Goal, timeSlotsOnly []TimeSlot, timeSlotsNot []TimeSlot) *FindTimeResponse {
	// hoursInAday := [24]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
	// blockers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 23}
	amountOfDays := 7
	// var timeSlots []int
	// var requiredTime int
	var returnerTasks []ProposedTask
	var returnerGoals []ProposedGoal
	returnerWeek := Week{Days: make(map[string]Day)}
	startDate := nextWeeklyEvent(time.Monday, 0)
	allAvailableTimes := make(map[int][]int)

	// for _, hidElement := range hoursInAday {
	// 	contained := false
	// 	for _, blockerElement := range blockers {
	// 		if hidElement == blockerElement {
	// 			contained = true
	// 		}
	// 	}
	// 	if !contained {
	// 		timeSlots = append(timeSlots, hidElement)
	// 	}
	// }

	for _, timeSlot := range timeSlotsOnly {
		var tmpTimes []int
		// tmpTimes[tag[timeSlot.DayIndex]] = make([]int, timeSlot.EndTime - timeSlot.StartTime)
		for n := timeSlot.StartTime; n < timeSlot.EndTime; n++ {
			tmpTimes = append(tmpTimes, n)
		}
		if allAvailableTimes[timeSlot.DayIndex] == nil {
			allAvailableTimes[timeSlot.DayIndex] = make([]int, len(tmpTimes))
		}
		allAvailableTimes[timeSlot.DayIndex] = removeDuplicateInt(append(allAvailableTimes[timeSlot.DayIndex], tmpTimes...))

	}

	for _, tsn := range timeSlotsNot {
		for n := tsn.StartTime; n < tsn.EndTime; n++ {
			j := 0
			for _, allowedTime := range allAvailableTimes[tsn.DayIndex] {
				if allowedTime == n {
					allAvailableTimes[tsn.DayIndex][j] = allowedTime
					j++
				}
			}
			allAvailableTimes[tsn.DayIndex] = allAvailableTimes[tsn.DayIndex][:j]
		}
	}

	// returnerWeek.Days[nextWeeklyEvent(time.Monday, i)] = Day{}
	// allAvailableTimes[i] = make([]int, len(timeSlots))
	// copy(allAvailableTimes[i], timeSlots)

	// for _, _task := range tasks {
	// 	requiredTime += _task.Duration
	// }

	// for _, goal := range goals {
	// 	requiredTime += goal.Duration
	// }

	// if requiredTime > len(timeSlots) {
	// 	return nil
	// }

	for index, _ := range tasks {
		_, proposedDate, time := getDayIndexDateAndTime(tasks[index].Title, amountOfDays, returnerWeek, &allAvailableTimes)
		putEventIntoWeek(&tasks[index], proposedDate, time, &returnerWeek)
		sort.Slice(returnerWeek.Days[proposedDate].SortedItems, func(i, j int) bool {
			return returnerWeek.Days[proposedDate].SortedItems[i].StartTime < returnerWeek.Days[proposedDate].SortedItems[j].StartTime
		})
	}

	for goalIndex, goal := range goals {
		returnerGoals = append(returnerGoals, ProposedGoal{&goals[goalIndex], 0})
		for frequencyIndex := 0; frequencyIndex < goal.Frequency; frequencyIndex++ {
			_, proposedDate, time := getDayIndexDateAndTime(goal.Title, amountOfDays, returnerWeek, &allAvailableTimes)
			putEventIntoWeek(goals[goalIndex].CreateTask, proposedDate, time, &returnerWeek)
			sort.Slice(returnerWeek.Days[proposedDate].SortedItems, func(i, j int) bool {
				return returnerWeek.Days[proposedDate].SortedItems[i].StartTime < returnerWeek.Days[proposedDate].SortedItems[j].StartTime
			})
		}
	}

	return &FindTimeResponse{ProposedTasks: returnerTasks, ProposedGoals: returnerGoals, Week: returnerWeek, StartDate: startDate, EndDate: nextWeeklyEvent(time.Monday, 6)}
}

func nextWeeklyEvent(weekday time.Weekday, mod int) string {
	days := int((7 + (weekday - time.Now().Weekday())) % 7)
	y, m, d := time.Now().AddDate(0, 0, days+mod).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Now().Location()).Format("20060102")
}

func Keys(m map[string]Day) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}

func contains(tasks []ProposedTask, title string) bool {
	for _, task := range tasks {
		if task.CreateTask.Title == title {
			return true
		}
	}
	return false
}

func getDayIndexDateAndTime(taskName string, amountOfDays int, returnerWeek Week, allAvailableTimes *map[int][]int) (int, string, string) {
	for {
		dayIndex, time, slotPosition := getDayIndexAndTime(allAvailableTimes, amountOfDays)
		proposedDate := Keys(returnerWeek.Days)[dayIndex]
		if !contains(returnerWeek.Days[proposedDate].SortedItems, taskName) {
			(*allAvailableTimes)[dayIndex] = append((*allAvailableTimes)[dayIndex][:slotPosition], (*allAvailableTimes)[dayIndex][slotPosition+1:]...)
			return dayIndex, proposedDate, time
		}
	}
}

func getDayIndexAndTime(allAvailableTimes *map[int][]int, amountOfDays int) (int, string, int) {
	dayIndex := rand.Intn(amountOfDays)
	time := ""
	var slotPosition int
	for time == "" {

		if len((*allAvailableTimes)[dayIndex]) > 0 {
			availableTimes := (*allAvailableTimes)[dayIndex]
			slotPosition := rand.Intn(len(availableTimes))
			time = strconv.Itoa(availableTimes[slotPosition]) + ":00"
		} else {
			dayIndex = rand.Intn(amountOfDays)
		}
	}
	return dayIndex, time, slotPosition
}

func putEventIntoWeek(event *CreateTask, date string, time string, week *Week) {
	if week.Days[date].SortedItems == nil {
		items := []ProposedTask{{event, time}}
		day := Day{SortedItems: items}
		week.Days[date] = day
	} else {
		items := ProposedTask{event, time}
		newDay := Day{SortedItems: append(week.Days[date].SortedItems, items)}
		week.Days[date] = newDay
	}
}

func removeDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
