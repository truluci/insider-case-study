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
