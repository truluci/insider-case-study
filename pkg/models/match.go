package models

import "time"

// Match represents a football match
type Match struct {
	ID         int       `json:"id"`
	HomeTeamID int       `json:"home_team_id"`
	AwayTeamID int       `json:"away_team_id"`
	Week       int       `json:"week"`
	HomeGoals  int       `json:"home_goals"`
	AwayGoals  int       `json:"away_goals"`
	Status     string    `json:"status"` // "scheduled", "completed", "played"
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// MatchDetail contains match with team names
type MatchDetail struct {
	ID           int       `json:"id"`
	Week         int       `json:"week"`
	HomeTeamID   int       `json:"home_team_id"`
	HomeTeamName string    `json:"home_team_name"`
	AwayTeamID   int       `json:"away_team_id"`
	AwayTeamName string    `json:"away_team_name"`
	HomeGoals    int       `json:"home_goals"`
	AwayGoals    int       `json:"away_goals"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// MatchRepository defines the interface for match data access
type MatchRepository interface {
	Create(match *Match) (*Match, error)
	GetByID(id int) (*Match, error)
	GetAll() ([]*Match, error)
	GetByWeek(week int) ([]*Match, error)
	Update(match *Match) (*Match, error)
	Delete(id int) error
}

// MatchService defines the interface for match business logic
type MatchService interface {
	CreateMatch(homeTeamID, awayTeamID, week int) (*Match, error)
	GetMatch(id int) (*Match, error)
	GetAllMatches() ([]*Match, error)
	GetMatchesByWeek(week int) ([]*Match, error)
	UpdateMatch(match *Match) (*Match, error)
	DeleteMatch(id int) error
	PlayMatch(id int) (*Match, error)
	PlayAllMatches() error
	PlayMatchesForWeek(week int) error
}
