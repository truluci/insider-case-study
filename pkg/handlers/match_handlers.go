package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetMatches handles GET /api/matches
func (h *Handler) GetMatches(w http.ResponseWriter, r *http.Request) {
	currentWeek, _ := h.db.GetCurrentWeek()

	matches, err := h.matchService.GetMatchesByWeek(currentWeek)
	if err != nil {
		http.Error(w, "Failed to fetch matches", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

// CreateMatch handles POST /api/matches
func (h *Handler) CreateMatch(w http.ResponseWriter, r *http.Request) {
	type CreateMatchRequest struct {
		HomeTeamID int `json:"home_team_id"`
		AwayTeamID int `json:"away_team_id"`
		Week       int `json:"week"`
	}

	var req CreateMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	match, err := h.matchService.CreateMatch(req.HomeTeamID, req.AwayTeamID, req.Week)
	if err != nil {
		http.Error(w, "Failed to create match", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(match)
}

// UpdateMatch handles PUT /api/matches/{id}
func (h *Handler) UpdateMatch(w http.ResponseWriter, r *http.Request) {
	type UpdateMatchRequest struct {
		HomeGoals int    `json:"home_goals"`
		AwayGoals int    `json:"away_goals"`
		Status    string `json:"status"`
	}

	var req UpdateMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// ID'yi path'ten al
	id := r.PathValue("id")

	// Match'i update et
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Match updated", "id": id})
}

// PlayAllMatches handles POST /api/play-all
func (h *Handler) PlayAllMatches(w http.ResponseWriter, r *http.Request) {
	currentWeek, err := h.db.GetCurrentWeek()
	if err != nil {
		http.Error(w, "Failed to get current week", http.StatusInternalServerError)
		return
	}

	if err := h.matchService.PlayMatchesForWeek(currentWeek); err != nil {
		http.Error(w, "Failed to play matches", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf("Week %d matches played successfully", currentWeek),
		"week":    currentWeek,
	})
}
