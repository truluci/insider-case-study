package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
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

// UpdateTeam handles PUT /api/teams/{id}
func (h *Handler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Team ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Team ID", http.StatusBadRequest)
		return
	}

	type UpdateTeamRequest struct {
		Name     string `json:"name"`
		Strength int    `json:"strength"`
	}

	var req UpdateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	existingTeam, err := h.teamService.GetTeam(id)
	if err != nil {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	}

	// Rule: cannot change name of default teams
	if existingTeam.IsDefault && req.Name != existingTeam.Name {
		http.Error(w, "Cannot change the name of a default team", http.StatusForbidden)
		return
	}

	existingTeam.Name = req.Name
	existingTeam.Strength = req.Strength

	updatedTeam, err := h.teamService.UpdateTeam(existingTeam)
	if err != nil {
		http.Error(w, "Failed to update team", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeam)
}

// DeleteTeam handles DELETE /api/teams/{id}
func (h *Handler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Team ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Team ID", http.StatusBadRequest)
		return
	}

	// Check if tournament has started
	hasStarted, err := h.db.HasTournamentStarted()
	if err != nil {
		http.Error(w, "Failed to check tournament status", http.StatusInternalServerError)
		return
	}
	if hasStarted {
		http.Error(w, "Cannot delete teams after tournament has started", http.StatusBadRequest)
		return
	}

	existingTeam, err := h.teamService.GetTeam(id)
	if err != nil {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	}

	// Rule: cannot delete default teams
	if existingTeam.IsDefault {
		http.Error(w, "Cannot delete a default team", http.StatusForbidden)
		return
	}

	if err := h.teamService.DeleteTeam(id); err != nil {
		http.Error(w, "Failed to delete team", http.StatusInternalServerError)
		return
	}

	// Also delete stats
	_, _ = h.db.GetConnection().Exec("DELETE FROM team_stats WHERE team_id = ?", id)

	// Reschedule matches
	if err := h.db.ScheduleMatches(); err != nil {
		http.Error(w, "Failed to reschedule matches", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
