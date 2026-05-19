# Makefile

.PHONY: help install dev build run clean

help:
	@echo "Football League Simulation - Available commands:"
	@echo "  make install    - Install dependencies (Go & Node)"
	@echo "  make dev        - Run in development mode"
	@echo "  make build      - Build backend and frontend"
	@echo "  make run        - Run backend server"
	@echo "  make frontend   - Run frontend development server"
	@echo "  make docker     - Run with Docker Compose"
	@echo "  make clean      - Clean build files"

install:
	@echo "Installing Go dependencies..."
	cd . && go mod tidy
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

dev: 
	@echo "Cleaning up old processes..."
	@lsof -i :8080 | grep -v COMMAND | awk '{print $$2}' | xargs kill -9 2>/dev/null || true
	@lsof -i :3000 | grep -v COMMAND | awk '{print $$2}' | xargs kill -9 2>/dev/null || true
	@sleep 1
	@echo "Removing old database..."
	@rm -f football_league.db
	@echo "Starting development environment..."
	@echo "Backend starting on http://localhost:8080"
	@echo "Frontend starting on http://localhost:3000"
	@echo ""
	go run ./cmd/main.go &
	cd frontend && npm run dev

build:
	@echo "Building backend..."
	go build -o bin/main ./cmd
	@echo "Building frontend..."
	cd frontend && npm run build

run:
	@echo "Running server on :8080"
	./bin/main

frontend:
	@echo "Running frontend dev server on :3000"
	cd frontend && npm run dev

docker:
	docker-compose up --build

clean:
	rm -rf bin/
	rm -rf frontend/dist/
	rm -f *.db

.DEFAULT_GOAL := help
