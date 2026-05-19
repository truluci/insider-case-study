package handlers

import (
	"encoding/json"
	"net/http"
)

// GetTeams handles GET /api/teams
func (h *Handler) GetTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.teamService.GetAllTeams()
	if err != nil {
		http.Error(w, "Failed to fetch teams", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

// CreateTeam handles POST /api/teams
func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	type CreateTeamRequest struct {
		Name     string `json:"name"`
		Strength int    `json:"strength"`
	}

	var req CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if tournament has started - prevent adding teams after matches begin
	hasStarted, err := h.db.HasTournamentStarted()
	if err != nil {
		http.Error(w, "Failed to check tournament status", http.StatusInternalServerError)
		return
	}
	if hasStarted {
		http.Error(w, "Cannot add teams after tournament has started", http.StatusBadRequest)
		return
	}

	team, err := h.teamService.CreateTeam(req.Name, req.Strength)
	if err != nil {
		http.Error(w, "Failed to create team", http.StatusInternalServerError)
		return
	}

	// Initialize team stats with 0 values
	_, err = h.db.GetConnection().Exec(`
		INSERT OR IGNORE INTO team_stats (team_id, played, won, drawn, lost, goals_for, goals_against, goal_diff, points)
		VALUES (?, 0, 0, 0, 0, 0, 0, 0, 0)
	`, team.ID)
	if err != nil {
		http.Error(w, "Failed to initialize team stats", http.StatusInternalServerError)
		return
	}

	// Reschedule matches with auto-calculated weeks
	if err := h.db.ScheduleMatches(); err != nil {
		http.Error(w, "Failed to reschedule matches", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(team)
}
