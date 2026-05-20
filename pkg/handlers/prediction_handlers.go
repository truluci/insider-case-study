package handlers

import (
	"encoding/json"
	"net/http"
)

// GetPredictions handles GET /api/predictions
func (h *Handler) GetPredictions(w http.ResponseWriter, r *http.Request) {
	predictions, err := h.predictionService.GetAllPredictions()
	if err != nil {
		http.Error(w, "Failed to fetch predictions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(predictions)
}

// CreatePrediction handles POST /api/predictions
func (h *Handler) CreatePrediction(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Week     int `json:"week"`
		TeamID   int `json:"team_id"`
		Position int `json:"position"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	prediction, err := h.predictionService.CreatePrediction(req.Week, req.TeamID, req.Position)
	if err != nil {
		http.Error(w, "Failed to create prediction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(prediction)
}
