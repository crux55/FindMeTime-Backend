package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

func FindTimeWorker(tasks []CreateTask) *FindTimeResponse {
	amountOfDays := 1
	var returnerTasks []ProposedTask
	var returnerGoals []ProposedGoal
	returnerWeek := Week{Days: make(map[string]Day)}
	startDate := nextWeeklyEvent(time.Monday, 0)

	for i := 0; i < amountOfDays; i++ {
		returnerWeek.Days[nextWeeklyEvent(time.Monday, i)] = Day{}
	}
	fmt.Println("about to loop through tasks")
	for taskIndex, task := range tasks {
		dayIndex, proposedDate, time := getDayIndexDateAndTime(task.Title, amountOfDays, returnerWeek, getAvailableTimes(task, returnerWeek), task.Duration)
		if dayIndex == -1 {
			fmt.Println("failed to find time for ", task.Title)
			continue
		}
		putEventIntoWeek(&tasks[taskIndex], proposedDate, time, &returnerWeek)
	}
	fmt.Println("about to return response", returnerTasks)
	for index := range returnerWeek.Days {
		sort.Slice(returnerWeek.Days[index].SortedItems, func(i, j int) bool {
			return returnerWeek.Days[index].SortedItems[i].StartTime < returnerWeek.Days[index].SortedItems[j].StartTime
		})
	}

	return &FindTimeResponse{ProposedTasks: returnerTasks, ProposedGoals: returnerGoals, Week: returnerWeek, StartDate: startDate, EndDate: nextWeeklyEvent(time.Monday, 6)}
}

func getAvailableTimes(task CreateTask, returnerWeek Week) *map[int][]int {
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

	for dayIndex, day := range returnerWeek.Days {
		for _, item := range day.SortedItems {
			di := indexOf(dayIndex, returnerWeek.Days)
			times, _ := allAvailableTimes[di]
			var filteredTimes []int
			for _, t := range times {
				castStartTime, _ := strconv.Atoi(strings.Split(item.StartTime, ":")[0])
				if t < castStartTime || t > castStartTime+item.Duration-1 {
					filteredTimes = append(filteredTimes, t)
				}
			}
			allAvailableTimes[di] = removeDuplicateInt(filteredTimes)
		}
	}
	fmt.Println("Got needed times")
	return &allAvailableTimes
}

func indexOf(_str string, dayArray map[string]Day) int {
	i := 0
	for str := range dayArray {
		if _str == str {
			return i
		}
		i++
	}
	return -1
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
	tries := 0
	fmt.Println("Getting day index date and time")
	for tries < 250 {
		dayIndex, proposedTime, slotPosition := getDayIndexAndTime(allAvailableTimes, amountOfDays, duration)
		proposedDate := Keys(returnerWeek.Days)[dayIndex]
		clash := false
		for _, item := range returnerWeek.Days[proposedDate].SortedItems {
			castTime, _ := strconv.Atoi(strings.Split(item.StartTime, ":")[0])
			castProposedTime, _ := strconv.Atoi(strings.Split(proposedTime, ":")[0])
			if castProposedTime+item.Duration > castTime &&
				castProposedTime < castTime+item.Duration {
				clash = true
				fmt.Println("Found clash")
				break
			}
		}
		if !clash && !contains(returnerWeek.Days[proposedDate].SortedItems, taskName) {
			(*allAvailableTimes)[dayIndex] = append((*allAvailableTimes)[dayIndex][:slotPosition], (*allAvailableTimes)[dayIndex][slotPosition+1:]...)
			return dayIndex, proposedDate, proposedTime
		}
		tries++
		fmt.Println("Increasing times", proposedTime)
	}
	return -1, "", "" //todo error handle
}

func getDayIndexAndTime(allAvailableTimes *map[int][]int, amountOfDays int, duration int) (int, string, int) {
	dayIndex := rand.Intn(amountOfDays)
	time := ""
	var slotPosition int
	for time == "" {

		if len((*allAvailableTimes)[dayIndex]) > 0 {
			availableTimes := (*allAvailableTimes)[dayIndex]
			if len(availableTimes) < duration {
				fmt.Println("loop saftey: getDayIndexAndTime")
				break
			} else if len(availableTimes) == duration {
				time = strconv.Itoa(availableTimes[0]) + ":00"
				return dayIndex, time, 0
			}

			slotPosition := rand.Intn(len(availableTimes)) // imporve this, it doesn't need ot be just random
			fmt.Println(slotPosition)
			if len(availableTimes) >= slotPosition+duration &&
				availableTimes[slotPosition]+duration-1 == availableTimes[slotPosition+duration-1] {
				time = strconv.Itoa(availableTimes[slotPosition]) + ":00"
			}

		} else {
			dayIndex = rand.Intn(amountOfDays)
		}
	}
	return dayIndex, time, slotPosition
}

func putEventIntoWeek(event *CreateTask, date string, time string, week *Week) {
	day, exists := (*week).Days[date]
	if !exists {
		day = Day{}
	}
	day.SortedItems = append(day.SortedItems, ProposedTask{event, time})
	(*week).Days[date] = day
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
