package handlers

import (
	"encoding/json"
	"net/http"
)

// GetLeagueTable handles GET /api/league
func (h *Handler) GetLeagueTable(w http.ResponseWriter, r *http.Request) {
	table, err := h.leagueService.GetCurrentStandings()
	if err != nil {
		http.Error(w, "Failed to fetch league table", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(table)
}

// GetLeagueTableByWeek handles GET /api/league/week/{week}
func (h *Handler) GetLeagueTableByWeek(w http.ResponseWriter, r *http.Request) {
	week := r.PathValue("week")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "League table for week", "week": week})
}

// GetCurrentWeek handles GET /api/tournament/current-week
func (h *Handler) GetCurrentWeek(w http.ResponseWriter, r *http.Request) {
	currentWeek, err := h.db.GetCurrentWeek()
	if err != nil {
		http.Error(w, "Failed to get current week", http.StatusInternalServerError)
		return
	}

	totalWeeks, _ := h.db.GetTotalWeeks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"current_week": currentWeek,
		"total_weeks":  totalWeeks,
	})
}

// NextWeek handles POST /api/tournament/next-week
func (h *Handler) NextWeek(w http.ResponseWriter, r *http.Request) {
	err := h.db.NextWeek()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentWeek, _ := h.db.GetCurrentWeek()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "Moved to next week",
		"current_week": currentWeek,
	})
}
