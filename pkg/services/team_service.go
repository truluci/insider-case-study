package services

import (
	"github.com/luci/football-league/pkg/models"
)

type TeamService struct {
	repo models.TeamRepository
}

func NewTeamService(repo models.TeamRepository) models.TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) CreateTeam(name string, strength int) (*models.Team, error) {
	team := &models.Team{
		Name:     name,
		Strength: strength,
	}
	return s.repo.Create(team)
}

func (s *TeamService) GetTeam(id int) (*models.Team, error) {
	return s.repo.GetByID(id)
}

func (s *TeamService) GetAllTeams() ([]*models.Team, error) {
	return s.repo.GetAll()
}

func (s *TeamService) UpdateTeam(team *models.Team) (*models.Team, error) {
	return s.repo.Update(team)
}

func (s *TeamService) DeleteTeam(id int) error {
	return s.repo.Delete(id)
}
