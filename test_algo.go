package main

import (
	"fmt"
)

func main() {
	teams := []int{1, 2, 3, 4, 5}

	if len(teams)%2 != 0 {
		teams = append(teams, -1)
	}
	
	numTeams := len(teams)
	numWeeksHalf := numTeams - 1
	numWeeks := numWeeksHalf * 2

	for round := 0; round < numWeeks; round++ {
		week := round + 1
		isSecondHalf := round >= numWeeksHalf
		
		fmt.Printf("Week %d:\n", week)
		for i := 0; i < numTeams/2; i++ {
			homeIdx := i
			awayIdx := numTeams - 1 - i

			homeID := teams[homeIdx]
			awayID := teams[awayIdx]

			if homeID == -1 || awayID == -1 {
				continue
			}

			if i == 0 && round%2 == 1 {
				homeID, awayID = awayID, homeID
			}

			if isSecondHalf {
				homeID, awayID = awayID, homeID
			}

			fmt.Printf("  %d vs %d\n", homeID, awayID)
		}

		temp := teams[1]
		for i := 1; i < numTeams-1; i++ {
			teams[i] = teams[i+1]
		}
		teams[numTeams-1] = temp
	}
}
