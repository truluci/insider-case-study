package database

import (
	"fmt"
)

// ResetMatches deletes all existing matches and tournament state
func (d *Database) ResetMatches() error {
	_, err := d.conn.Exec("DELETE FROM matches")
	if err != nil {
		return fmt.Errorf("failed to reset matches: %w", err)
	}
	_, err = d.conn.Exec("DELETE FROM tournament_state")
	if err != nil {
		return fmt.Errorf("failed to reset tournament state: %w", err)
	}
	return nil
}

// ScheduleMatches creates round-robin matches for all teams
// Weeks calculation:
// - Even teams (4, 18): H = 2 * (N - 1)
// - Odd teams (5, 7): H = 2 * N
func (d *Database) ScheduleMatches() error {
	// Get all teams
	rows, err := d.conn.Query("SELECT id FROM teams ORDER BY id")
	if err != nil {
		return fmt.Errorf("failed to get teams: %w", err)
	}
	defer rows.Close()

	var teams []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		teams = append(teams, id)
	}

	if len(teams) == 0 {
		return nil
	}

	// Reset existing matches
	if err := d.ResetMatches(); err != nil {
		return err
	}

	if len(teams)%2 != 0 {
		teams = append(teams, -1) // dummy team for odd number of teams
	}

	numTeams := len(teams)
	numWeeksHalf := numTeams - 1
	numWeeks := numWeeksHalf * 2

	insertQuery := `INSERT INTO matches (home_team_id, away_team_id, week, status) VALUES (?, ?, ?, 'scheduled')`

	for round := 0; round < numWeeks; round++ {
		week := round + 1
		isSecondHalf := round >= numWeeksHalf
		
		for i := 0; i < numTeams/2; i++ {
			homeIdx := i
			awayIdx := numTeams - 1 - i

			homeID := teams[homeIdx]
			awayID := teams[awayIdx]

			if homeID == -1 || awayID == -1 {
				continue // skip matches with the dummy team (bye week)
			}

			// Swap home/away for the second half of the season
			if isSecondHalf {
				homeID, awayID = awayID, homeID
			}

			_, err := d.conn.Exec(insertQuery, homeID, awayID, week)
			if err != nil {
				return fmt.Errorf("failed to schedule match: %w", err)
			}
		}

		// Rotate teams for next round (except first team)
		temp := teams[1]
		for i := 1; i < numTeams-1; i++ {
			teams[i] = teams[i+1]
		}
		teams[numTeams-1] = temp
	}

	// Initialize tournament_state
	_, err = d.conn.Exec(`INSERT OR REPLACE INTO tournament_state (id, current_week, total_weeks) VALUES (1, 1, ?)`, numWeeks)
	if err != nil {
		return fmt.Errorf("failed to init tournament state: %w", err)
	}

	return nil
}

// GetCurrentWeek returns the current tournament week
func (d *Database) GetCurrentWeek() (int, error) {
	var week int
	err := d.conn.QueryRow("SELECT COALESCE(current_week, 1) FROM tournament_state WHERE id = 1").Scan(&week)
	if err != nil {
		return 1, nil // Default to week 1
	}
	return week, nil
}

// GetTotalWeeks returns total weeks in tournament
func (d *Database) GetTotalWeeks() (int, error) {
	var weeks int
	err := d.conn.QueryRow("SELECT COALESCE(total_weeks, 0) FROM tournament_state WHERE id = 1").Scan(&weeks)
	if err != nil {
		return 0, nil
	}
	return weeks, nil
}

// NextWeek advances to the next week
func (d *Database) NextWeek() error {
	currentWeek, _ := d.GetCurrentWeek()
	totalWeeks, _ := d.GetTotalWeeks()

	if currentWeek >= totalWeeks {
		return fmt.Errorf("tournament already finished")
	}

	_, err := d.conn.Exec("UPDATE tournament_state SET current_week = ?, updated_at = CURRENT_TIMESTAMP WHERE id = 1", currentWeek+1)
	return err
}

// HasTournamentStarted checks if any match has been played
func (d *Database) HasTournamentStarted() (bool, error) {
	var count int
	err := d.conn.QueryRow("SELECT COUNT(*) FROM matches WHERE status = 'completed'").Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
