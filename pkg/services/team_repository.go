package services

import (
	"database/sql"
	"time"

	"github.com/luci/football-league/pkg/models"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) models.TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(team *models.Team) (*models.Team, error) {
	now := time.Now()
	query := `INSERT INTO teams (name, strength, is_default, created_at, updated_at) VALUES (?, ?, ?, ?, ?) RETURNING id`

	err := r.db.QueryRow(query, team.Name, team.Strength, team.IsDefault, now, now).Scan(&team.ID)
	if err != nil {
		return nil, err
	}

	team.CreatedAt = now
	team.UpdatedAt = now
	return team, nil
}

func (r *TeamRepository) GetByID(id int) (*models.Team, error) {
	query := `SELECT id, name, strength, is_default, created_at, updated_at FROM teams WHERE id = ?`

	team := &models.Team{}
	err := r.db.QueryRow(query, id).Scan(&team.ID, &team.Name, &team.Strength, &team.IsDefault, &team.CreatedAt, &team.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (r *TeamRepository) GetAll() ([]*models.Team, error) {
	query := `SELECT id, name, strength, is_default, created_at, updated_at FROM teams ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teams := make([]*models.Team, 0)
	for rows.Next() {
		team := &models.Team{}
		err := rows.Scan(&team.ID, &team.Name, &team.Strength, &team.IsDefault, &team.CreatedAt, &team.UpdatedAt)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, rows.Err()
}

func (r *TeamRepository) Update(team *models.Team) (*models.Team, error) {
	now := time.Now()
	query := `UPDATE teams SET name = ?, strength = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.Exec(query, team.Name, team.Strength, now, team.ID)
	if err != nil {
		return nil, err
	}

	team.UpdatedAt = now
	return team, nil
}

func (r *TeamRepository) Delete(id int) error {
	query := `DELETE FROM teams WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
