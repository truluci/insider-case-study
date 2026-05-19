package services

import (
	"math/rand"

	"github.com/luci/football-league/pkg/models"
)

type MatchService struct {
	matchRepo models.MatchRepository
	teamRepo  models.TeamRepository
}

func NewMatchService(matchRepo models.MatchRepository, teamRepo models.TeamRepository) models.MatchService {
	return &MatchService{
		matchRepo: matchRepo,
		teamRepo:  teamRepo,
	}
}

func (s *MatchService) CreateMatch(homeTeamID, awayTeamID, week int) (*models.Match, error) {
	match := &models.Match{
		HomeTeamID: homeTeamID,
		AwayTeamID: awayTeamID,
		Week:       week,
		Status:     "scheduled",
	}
	return s.matchRepo.Create(match)
}

func (s *MatchService) GetMatch(id int) (*models.Match, error) {
	return s.matchRepo.GetByID(id)
}

func (s *MatchService) GetAllMatches() ([]*models.Match, error) {
	return s.matchRepo.GetAll()
}

func (s *MatchService) GetMatchesByWeek(week int) ([]*models.Match, error) {
	return s.matchRepo.GetByWeek(week)
}

func (s *MatchService) UpdateMatch(match *models.Match) (*models.Match, error) {
	return s.matchRepo.Update(match)
}

func (s *MatchService) DeleteMatch(id int) error {
	return s.matchRepo.Delete(id)
}

// PlayMatch plays a single match and generates result based on team strengths
func (s *MatchService) PlayMatch(id int) (*models.Match, error) {
	match, err := s.matchRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Get team strengths
	homeTeam, _ := s.teamRepo.GetByID(match.HomeTeamID)
	awayTeam, _ := s.teamRepo.GetByID(match.AwayTeamID)

	// Generate result based on team strengths
	homeGoals := s.generateGoals(homeTeam.Strength, awayTeam.Strength)
	awayGoals := s.generateGoals(awayTeam.Strength, homeTeam.Strength)

	match.HomeGoals = homeGoals
	match.AwayGoals = awayGoals
	match.Status = "completed"

	return s.matchRepo.Update(match)
}

// PlayAllMatches plays all scheduled matches
func (s *MatchService) PlayAllMatches() error {
	matches, err := s.matchRepo.GetAll()
	if err != nil {
		return err
	}

	for _, match := range matches {
		if match.Status == "scheduled" {
			s.PlayMatch(match.ID)
		}
	}
	return nil
}

// PlayMatchesForWeek plays all matches for a specific week
func (s *MatchService) PlayMatchesForWeek(week int) error {
	matches, err := s.matchRepo.GetByWeek(week)
	if err != nil {
		return err
	}

	for _, match := range matches {
		if match.Status == "scheduled" {
			s.PlayMatch(match.ID)
		}
	}
	return nil
}

// generateGoals generates goals based on team strengths
func (s *MatchService) generateGoals(teamStrength, opponentStrength int) int {
	// Strength difference affects goal probability
	diff := teamStrength - opponentStrength
	baseGoals := rand.Intn(3) // 0, 1, or 2

	// Add some randomness based on strength
	if diff > 0 {
		baseGoals += rand.Intn(2)
	}

	return baseGoals
}
