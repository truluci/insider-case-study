package services

import (
	"sort"

	"github.com/luci/football-league/pkg/models"
)

type LeagueService struct {
	teamRepo  models.TeamRepository
	matchRepo models.MatchRepository
}

func NewLeagueService(teamRepo models.TeamRepository, matchRepo models.MatchRepository) *LeagueService {
	return &LeagueService{
		teamRepo:  teamRepo,
		matchRepo: matchRepo,
	}
}

// GetCurrentStandings returns the current league standings
func (s *LeagueService) GetCurrentStandings() ([]*models.TeamStats, error) {
	teams, err := s.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}

	matches, err := s.matchRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Calculate stats for each team
	stats := make(map[int]*models.TeamStats)
	for _, team := range teams {
		stats[team.ID] = &models.TeamStats{
			TeamID:   team.ID,
			TeamName: team.Name,
		}
	}

	// Process all completed matches
	for _, match := range matches {
		if match.Status != "completed" {
			continue
		}

		homeStats := stats[match.HomeTeamID]
		awayStats := stats[match.AwayTeamID]

		homeStats.Played++
		awayStats.Played++

		homeStats.GoalsFor += match.HomeGoals
		awayStats.GoalsFor += match.AwayGoals

		homeStats.GoalsAgainst += match.AwayGoals
		awayStats.GoalsAgainst += match.HomeGoals

		homeStats.GoalDiff = homeStats.GoalsFor - homeStats.GoalsAgainst
		awayStats.GoalDiff = awayStats.GoalsFor - awayStats.GoalsAgainst

		// Determine winner
		if match.HomeGoals > match.AwayGoals {
			homeStats.Won++
			homeStats.Points += 3
			awayStats.Lost++
		} else if match.HomeGoals < match.AwayGoals {
			awayStats.Won++
			awayStats.Points += 3
			homeStats.Lost++
		} else {
			homeStats.Drawn++
			homeStats.Points += 1
			awayStats.Drawn++
			awayStats.Points += 1
		}
	}

	// Convert to slice and sort
	result := make([]*models.TeamStats, 0, len(stats))
	for _, stat := range stats {
		result = append(result, stat)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Points != result[j].Points {
			return result[i].Points > result[j].Points
		}
		return result[i].GoalDiff > result[j].GoalDiff
	})

	// Set positions
	for i := range result {
		result[i].Position = i + 1
	}

	return result, nil
}
