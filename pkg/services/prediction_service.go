package services

import (
	"github.com/luci/football-league/pkg/models"
)

type PredictionService struct {
	repo models.PredictionRepository
}

func NewPredictionService(repo models.PredictionRepository) models.PredictionService {
	return &PredictionService{repo: repo}
}

func (s *PredictionService) CreatePrediction(week, teamID, position int) (*models.Prediction, error) {
	prediction := &models.Prediction{
		Week:     week,
		TeamID:   teamID,
		Position: position,
	}
	return s.repo.Create(prediction)
}

func (s *PredictionService) GetPrediction(id int) (*models.Prediction, error) {
	return s.repo.GetByID(id)
}

func (s *PredictionService) GetPredictionsByWeek(week int) ([]*models.Prediction, error) {
	return s.repo.GetByWeek(week)
}

func (s *PredictionService) GetAllPredictions() ([]*models.Prediction, error) {
	return s.repo.GetAll()
}

func (s *PredictionService) UpdatePrediction(prediction *models.Prediction) (*models.Prediction, error) {
	return s.repo.Update(prediction)
}

func (s *PredictionService) DeletePrediction(id int) error {
	return s.repo.Delete(id)
}
