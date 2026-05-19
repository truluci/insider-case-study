package models

import "time"

// Prediction represents a user's prediction of the final league table
type Prediction struct {
	ID        int       `json:"id"`
	Week      int       `json:"week"`
	TeamID    int       `json:"team_id"`
	TeamName  string    `json:"team_name"`
	Position  int       `json:"position"` // Predicted position (1-4)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PredictionRepository defines the interface for prediction data access
type PredictionRepository interface {
	Create(prediction *Prediction) (*Prediction, error)
	GetByID(id int) (*Prediction, error)
	GetByWeek(week int) ([]*Prediction, error)
	GetAll() ([]*Prediction, error)
	Update(prediction *Prediction) (*Prediction, error)
	Delete(id int) error
}

// PredictionService defines the interface for prediction business logic
type PredictionService interface {
	CreatePrediction(week, teamID, position int) (*Prediction, error)
	GetPrediction(id int) (*Prediction, error)
	GetPredictionsByWeek(week int) ([]*Prediction, error)
	GetAllPredictions() ([]*Prediction, error)
	UpdatePrediction(prediction *Prediction) (*Prediction, error)
	DeletePrediction(id int) error
}
