package services

import (
	"database/sql"
	"time"

	"github.com/luci/football-league/pkg/models"
)

type PredictionRepository struct {
	db *sql.DB
}

func NewPredictionRepository(db *sql.DB) models.PredictionRepository {
	return &PredictionRepository{db: db}
}

func (r *PredictionRepository) Create(prediction *models.Prediction) (*models.Prediction, error) {
	now := time.Now()
	query := `INSERT INTO predictions (week, team_id, position, created_at, updated_at) VALUES (?, ?, ?, ?, ?) RETURNING id`

	err := r.db.QueryRow(query, prediction.Week, prediction.TeamID, prediction.Position, now, now).Scan(&prediction.ID)
	if err != nil {
		return nil, err
	}

	prediction.CreatedAt = now
	prediction.UpdatedAt = now
	return prediction, nil
}

func (r *PredictionRepository) GetByID(id int) (*models.Prediction, error) {
	query := `SELECT id, week, team_id, position, created_at, updated_at FROM predictions WHERE id = ?`

	prediction := &models.Prediction{}
	err := r.db.QueryRow(query, id).Scan(&prediction.ID, &prediction.Week, &prediction.TeamID, &prediction.Position, &prediction.CreatedAt, &prediction.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return prediction, nil
}

func (r *PredictionRepository) GetByWeek(week int) ([]*models.Prediction, error) {
	query := `SELECT id, week, team_id, position, created_at, updated_at FROM predictions WHERE week = ?`

	rows, err := r.db.Query(query, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	predictions := make([]*models.Prediction, 0)
	for rows.Next() {
		prediction := &models.Prediction{}
		err := rows.Scan(&prediction.ID, &prediction.Week, &prediction.TeamID, &prediction.Position, &prediction.CreatedAt, &prediction.UpdatedAt)
		if err != nil {
			return nil, err
		}
		predictions = append(predictions, prediction)
	}

	return predictions, rows.Err()
}

func (r *PredictionRepository) GetAll() ([]*models.Prediction, error) {
	query := `SELECT id, week, team_id, position, created_at, updated_at FROM predictions ORDER BY week, position`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	predictions := make([]*models.Prediction, 0)
	for rows.Next() {
		prediction := &models.Prediction{}
		err := rows.Scan(&prediction.ID, &prediction.Week, &prediction.TeamID, &prediction.Position, &prediction.CreatedAt, &prediction.UpdatedAt)
		if err != nil {
			return nil, err
		}
		predictions = append(predictions, prediction)
	}

	return predictions, rows.Err()
}

func (r *PredictionRepository) Update(prediction *models.Prediction) (*models.Prediction, error) {
	now := time.Now()
	query := `UPDATE predictions SET week = ?, team_id = ?, position = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.Exec(query, prediction.Week, prediction.TeamID, prediction.Position, now, prediction.ID)
	if err != nil {
		return nil, err
	}

	prediction.UpdatedAt = now
	return prediction, nil
}

func (r *PredictionRepository) Delete(id int) error {
	query := `DELETE FROM predictions WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
