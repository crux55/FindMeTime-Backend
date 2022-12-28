import (
	"math/rand"
)

func findTime(findTimeRequest FindTimeRequest) *FindTimeResponse {
	hoursInAday := [24]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
	blockers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 23}

	var timeSlots []int
	var requiredTime int
	var returnerTasks []ProposedTask
	var returnerGoals []ProposedGoal

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

	for _, task := range findTimeRequest.Tasks {
		requiredTime += task.Duration
	}

	for _, goal := range findTimeRequest.Goals {
		requiredTime += goal.Duration
	}

	if requiredTime > len(timeSlots) {
		return nil
	}

	for _, task := range findTimeRequest.Tasks {
		slotPosition := rand.Intn(len(timeSlots))
		time := timeSlots[slotPosition]
		returnerTasks = append(returnerTasks, ProposedTask{&task, time})

	}

	for _, goal := range findTimeRequest.Goals {
		slotPosition := rand.Intn(len(timeSlots))
		time := timeSlots[slotPosition]
		returnerGoals = append(returnerGoals, ProposedGoal{&goal, time})

	}

	return &FindTimeResponse{ProposedTasks: returnerTasks, ProposedGoals: returnerGoals}
}
