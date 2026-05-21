package models

import "time"

// Team represents a football team in the league
type Team struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Strength  int       `json:"strength"` // 1-100 scale for team strength
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TeamStats represents the current standing of a team
type TeamStats struct {
	ID            int       `json:"id"`
	TeamID        int       `json:"team_id"`
	TeamName      string    `json:"team_name"`
	Played        int       `json:"played"`
	Won           int       `json:"won"`
	Drawn         int       `json:"drawn"`
	Lost          int       `json:"lost"`
	GoalsFor      int       `json:"goals_for"`
	GoalsAgainst  int       `json:"goals_against"`
	GoalDiff      int       `json:"goal_diff"`
	Points        int       `json:"points"`
	Position      int       `json:"position"`
	WinChance     string    `json:"win_chance"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

// TeamRepository defines the interface for team data access
type TeamRepository interface {
	Create(team *Team) (*Team, error)
	GetByID(id int) (*Team, error)
	GetAll() ([]*Team, error)
	Update(team *Team) (*Team, error)
	Delete(id int) error
}

// TeamService defines the interface for team business logic
type TeamService interface {
	CreateTeam(name string, strength int) (*Team, error)
	GetTeam(id int) (*Team, error)
	GetAllTeams() ([]*Team, error)
	UpdateTeam(team *Team) (*Team, error)
	DeleteTeam(id int) error
}
