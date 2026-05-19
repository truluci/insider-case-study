package handlers

import (
	"database/sql"

	"github.com/luci/football-league/pkg/database"
	"github.com/luci/football-league/pkg/models"
	"github.com/luci/football-league/pkg/services"
)

// Handler struct holds all service dependencies
type Handler struct {
	teamService       models.TeamService
	matchService      models.MatchService
	predictionService models.PredictionService
	leagueService     *services.LeagueService
	db                *database.Database
}

// NewHandler creates a new handler with dependencies
func NewHandler(dbConn *sql.DB, db *database.Database) *Handler {
	// Initialize repositories
	teamRepo := services.NewTeamRepository(dbConn)
	matchRepo := services.NewMatchRepository(dbConn)
	predictionRepo := services.NewPredictionRepository(dbConn)

	// Initialize services
	teamService := services.NewTeamService(teamRepo)
	matchService := services.NewMatchService(matchRepo, teamRepo)
	predictionService := services.NewPredictionService(predictionRepo)
	leagueService := services.NewLeagueService(teamRepo, matchRepo)

	return &Handler{
		teamService:       teamService,
		matchService:      matchService,
		predictionService: predictionService,
		leagueService:     leagueService,
		db:                db,
	}
}
