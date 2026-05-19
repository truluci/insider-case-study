package services

import (
	"database/sql"
	"time"

	"github.com/luci/football-league/pkg/models"
)

type MatchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) models.MatchRepository {
	return &MatchRepository{db: db}
}

func (r *MatchRepository) Create(match *models.Match) (*models.Match, error) {
	now := time.Now()
	query := `INSERT INTO matches (home_team_id, away_team_id, week, home_goals, away_goals, status, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`

	err := r.db.QueryRow(query, match.HomeTeamID, match.AwayTeamID, match.Week, match.HomeGoals, match.AwayGoals, match.Status, now, now).Scan(&match.ID)
	if err != nil {
		return nil, err
	}

	match.CreatedAt = now
	match.UpdatedAt = now
	return match, nil
}

func (r *MatchRepository) GetByID(id int) (*models.Match, error) {
	query := `SELECT id, home_team_id, away_team_id, week, home_goals, away_goals, status, created_at, updated_at FROM matches WHERE id = ?`

	match := &models.Match{}
	err := r.db.QueryRow(query, id).Scan(&match.ID, &match.HomeTeamID, &match.AwayTeamID, &match.Week, &match.HomeGoals, &match.AwayGoals, &match.Status, &match.CreatedAt, &match.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return match, nil
}

func (r *MatchRepository) GetAll() ([]*models.Match, error) {
	query := `SELECT id, home_team_id, away_team_id, week, home_goals, away_goals, status, created_at, updated_at FROM matches ORDER BY week, id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]*models.Match, 0)
	for rows.Next() {
		match := &models.Match{}
		err := rows.Scan(&match.ID, &match.HomeTeamID, &match.AwayTeamID, &match.Week, &match.HomeGoals, &match.AwayGoals, &match.Status, &match.CreatedAt, &match.UpdatedAt)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return matches, rows.Err()
}

func (r *MatchRepository) GetByWeek(week int) ([]*models.Match, error) {
	query := `SELECT id, home_team_id, away_team_id, week, home_goals, away_goals, status, created_at, updated_at FROM matches WHERE week = ? ORDER BY id`

	rows, err := r.db.Query(query, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]*models.Match, 0)
	for rows.Next() {
		match := &models.Match{}
		err := rows.Scan(&match.ID, &match.HomeTeamID, &match.AwayTeamID, &match.Week, &match.HomeGoals, &match.AwayGoals, &match.Status, &match.CreatedAt, &match.UpdatedAt)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return matches, rows.Err()
}

func (r *MatchRepository) Update(match *models.Match) (*models.Match, error) {
	now := time.Now()
	query := `UPDATE matches SET home_goals = ?, away_goals = ?, status = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.Exec(query, match.HomeGoals, match.AwayGoals, match.Status, now, match.ID)
	if err != nil {
		return nil, err
	}

	match.UpdatedAt = now
	return match, nil
}

func (r *MatchRepository) Delete(id int) error {
	query := `DELETE FROM matches WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
