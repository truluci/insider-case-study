package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	conn *sql.DB
}

func NewConnection() (*Database, error) {
	// SQLite kullanacağız basitlik için
	// Gerçek ortamda PostgreSQL kullanabilirsiniz
	db, err := sql.Open("sqlite3", "./football_league.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{conn: db}, nil
}

func (d *Database) Close() error {
	return d.conn.Close()
}

func (d *Database) GetConnection() *sql.DB {
	return d.conn
}

func (d *Database) CreateTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS teams (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		strength INTEGER NOT NULL DEFAULT 50,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS matches (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		home_team_id INTEGER NOT NULL,
		away_team_id INTEGER NOT NULL,
		week INTEGER NOT NULL,
		home_goals INTEGER DEFAULT 0,
		away_goals INTEGER DEFAULT 0,
		status TEXT DEFAULT 'scheduled',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(home_team_id) REFERENCES teams(id),
		FOREIGN KEY(away_team_id) REFERENCES teams(id)
	);

	CREATE TABLE IF NOT EXISTS team_stats (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		team_id INTEGER NOT NULL UNIQUE,
		played INTEGER DEFAULT 0,
		won INTEGER DEFAULT 0,
		drawn INTEGER DEFAULT 0,
		lost INTEGER DEFAULT 0,
		goals_for INTEGER DEFAULT 0,
		goals_against INTEGER DEFAULT 0,
		goal_diff INTEGER DEFAULT 0,
		points INTEGER DEFAULT 0,
		last_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(team_id) REFERENCES teams(id)
	);

	CREATE TABLE IF NOT EXISTS predictions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		week INTEGER NOT NULL,
		team_id INTEGER NOT NULL,
		position INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(team_id) REFERENCES teams(id)
	);

	CREATE TABLE IF NOT EXISTS tournament_state (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		current_week INTEGER DEFAULT 1,
		total_weeks INTEGER DEFAULT 0,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := d.conn.Exec(schema); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

// SeedDefaultTeams inserts default teams if table is empty
func (d *Database) SeedDefaultTeams() error {
	// Check if teams already exist
	var count int
	err := d.conn.QueryRow("SELECT COUNT(*) FROM teams").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count teams: %w", err)
	}

	// If teams already exist, don't seed
	if count > 0 {
		return nil
	}

	// Default teams
	teams := []struct {
		name     string
		strength int
	}{
		{"Chelsea", 85},
		{"Manchester City", 88},
		{"Arsenal", 78},
		{"Liverpool", 82},
	}

	query := `INSERT INTO teams (name, strength) VALUES (?, ?)`

	for _, team := range teams {
		_, err := d.conn.Exec(query, team.name, team.strength)
		if err != nil {
			return fmt.Errorf("failed to seed team %s: %w", team.name, err)
		}
	}

	return nil
}

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

	// Calculate weeks based on formula
	numTeams := len(teams)
	numWeeks := 0
	if numTeams%2 == 0 {
		// Even: H = 2 * (N - 1)
		numWeeks = 2 * (numTeams - 1)
	} else {
		// Odd: H = 2 * N
		numWeeks = 2 * numTeams
	}

	insertQuery := `INSERT INTO matches (home_team_id, away_team_id, week, status) VALUES (?, ?, ?, 'scheduled')`

	week := 1
	for round := 0; round < numWeeks; round++ {
		// Generate matches for this round
		for i := 0; i < numTeams/2; i++ {
			homeIdx := i
			awayIdx := numTeams - 1 - i

			homeID := teams[homeIdx]
			awayID := teams[awayIdx]

			_, err := d.conn.Exec(insertQuery, homeID, awayID, week)
			if err != nil {
				return fmt.Errorf("failed to schedule match: %w", err)
			}
		}

		// Rotate teams for next round (except first and last)
		if round < numWeeks-1 {
			temp := teams[1]
			for i := 1; i < numTeams-1; i++ {
				teams[i] = teams[i+1]
			}
			teams[numTeams-1] = temp
		}

		week++
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
