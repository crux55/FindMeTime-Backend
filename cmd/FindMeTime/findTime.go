package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

func FindTimeWorker(tasks []CreateTask) *FindTimeResponse {
	amountOfDays := 7
	var returnerTasks []ProposedTask
	var returnerGoals []ProposedGoal
	returnerWeek := Week{Days: make(map[string]Day)}
	startDate := nextWeeklyEvent(time.Monday, 0)

	for i := 0; i < amountOfDays; i++ {
		returnerWeek.Days[nextWeeklyEvent(time.Monday, i)] = Day{}
	}
	fmt.Println("about to loop through tasks")
	for index, task := range tasks {
		_, proposedDate, time := getDayIndexDateAndTime(task.Title, amountOfDays, returnerWeek, getAvailableTimes(task), task.Duration)
		putEventIntoWeek(&tasks[index], proposedDate, time, &returnerWeek)
		sort.Slice(returnerWeek.Days[proposedDate].SortedItems, func(i, j int) bool {
			return returnerWeek.Days[proposedDate].SortedItems[i].StartTime < returnerWeek.Days[proposedDate].SortedItems[j].StartTime
		})
	}
	fmt.Println("about to return response", returnerTasks)
	return &FindTimeResponse{ProposedTasks: returnerTasks, ProposedGoals: returnerGoals, Week: returnerWeek, StartDate: startDate, EndDate: nextWeeklyEvent(time.Monday, 6)}
}

func getAvailableTimes(task CreateTask) *map[int][]int {
	allAvailableTimes := make(map[int][]int)
	for _, tag := range task.TagsOnly {
		for _, timeSlot := range tag.TimeSlots {
			var tmpTimes []int
			for n := timeSlot.StartTime; n < timeSlot.EndTime; n++ {
				tmpTimes = append(tmpTimes, n)
			}
			if allAvailableTimes[timeSlot.DayIndex] == nil {
				allAvailableTimes[timeSlot.DayIndex] = make([]int, 0, len(tmpTimes))
			}
			allAvailableTimes[timeSlot.DayIndex] = removeDuplicateInt(append(allAvailableTimes[timeSlot.DayIndex], tmpTimes...))

		}
	}

	for _, notag := range task.TagsNot {
		for _, timeSlot := range notag.TimeSlots {
			times, err := allAvailableTimes[timeSlot.DayIndex]
			if !err {
				break
			}
			var filteredTimes []int
			for _, t := range times {
				if t < timeSlot.StartTime || t > timeSlot.EndTime-1 {
					filteredTimes = append(filteredTimes, t)
				}
			}
			allAvailableTimes[timeSlot.DayIndex] = filteredTimes
		}
	}

	return &allAvailableTimes
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

func getDayIndexDateAndTime(taskName string, amountOfDays int, returnerWeek Week, allAvailableTimes *map[int][]int, duration int) (int, string, string) {
	for {
		dayIndex, time, slotPosition := getDayIndexAndTime(allAvailableTimes, amountOfDays, duration)
		proposedDate := Keys(returnerWeek.Days)[dayIndex]
		if !contains(returnerWeek.Days[proposedDate].SortedItems, taskName) {
			(*allAvailableTimes)[dayIndex] = append((*allAvailableTimes)[dayIndex][:slotPosition], (*allAvailableTimes)[dayIndex][slotPosition+1:]...)
			return dayIndex, proposedDate, time
		}
	}
}

func getDayIndexAndTime(allAvailableTimes *map[int][]int, amountOfDays int, duration int) (int, string, int) {
	dayIndex := rand.Intn(amountOfDays)
	time := ""
	var slotPosition int
	for time == "" {

		if len((*allAvailableTimes)[dayIndex]) > 0 {
			availableTimes := (*allAvailableTimes)[dayIndex]
			slotPosition := rand.Intn(len(availableTimes) + 1)
			if len(availableTimes) > slotPosition+duration-1 &&
				availableTimes[slotPosition]+duration == availableTimes[slotPosition+duration] {
				time = strconv.Itoa(availableTimes[slotPosition]) + ":00"
			}

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
