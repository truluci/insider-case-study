package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luci/football-league/pkg/database"
	"github.com/luci/football-league/pkg/handlers"
)

func main() {
	// Database bağlantısını başlat
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	// Tabloları oluştur
	if err := db.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Default takımları ekle
	if err := db.SeedDefaultTeams(); err != nil {
		log.Fatalf("Failed to seed default teams: %v", err)
	}

	// Otomatik maç schedule'ı oluştur (formula: even=2*(n-1), odd=2*n)
	if err := db.ScheduleMatches(); err != nil {
		log.Fatalf("Failed to schedule matches: %v", err)
	}

	// Handler'ları initialize et
	h := handlers.NewHandler(db.GetConnection(), db)

	// Router'ı setup et
	router := setupRouter(h)

	// Server'ı başlat
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful shutdown için signal handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\nShutdown signal received: %v\n", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}
	}()

	log.Println("Server starting on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func setupRouter(h *handlers.Handler) http.Handler {
	mux := http.NewServeMux()

	// Teams endpoints
	mux.HandleFunc("GET /api/teams", h.GetTeams)
	mux.HandleFunc("POST /api/teams", h.CreateTeam)

	// Matches endpoints
	mux.HandleFunc("GET /api/matches", h.GetMatches)
	mux.HandleFunc("POST /api/matches", h.CreateMatch)
	mux.HandleFunc("PUT /api/matches/{id}", h.UpdateMatch)

	// League table endpoints
	mux.HandleFunc("GET /api/league", h.GetLeagueTable)
	mux.HandleFunc("GET /api/league/week/{week}", h.GetLeagueTableByWeek)

	// Tournament endpoints
	mux.HandleFunc("GET /api/tournament/current-week", h.GetCurrentWeek)
	mux.HandleFunc("POST /api/tournament/next-week", h.NextWeek)

	// Predictions endpoints
	mux.HandleFunc("GET /api/predictions", h.GetPredictions)

	// Auto play endpoints
	mux.HandleFunc("POST /api/play-all", h.PlayAllMatches)

	return mux
}
