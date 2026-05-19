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
