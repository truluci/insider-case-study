# ⚽ Football League Simulation

A full-stack football league simulation application built with Go (backend) and Vue 3 (frontend). It allows you to create teams, simulate matches week by week (or all at once), track league standings, and even predict the champion with dynamic win probabilities.

## ✨ Features
- **Team Management:** Add teams and define their core strength.
- **Match Simulation:** Automatically generates a schedule and simulates match results based on team strength.
- **Dynamic League Table:** Live tracking of Points, Goal Difference, and Championship Win %.
- **Predictions:** Lock in your prediction for the final standings before the league starts!
- **Neo-Brutalist UI:** A sleek, custom Vue 3 frontend interface.

## 🚀 Quick Start

### Requirements
- Go 1.23+
- Node.js 18+

### Setup & Run
The easiest way to start both the backend and frontend is using the provided Makefile:

```bash
# 1. Install dependencies
make install

# 2. Run the development environment (Backend + Frontend)
make dev
```

The application will be available at:
- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080

**Note for deployed version:** The backend is deployed on Render and the frontend on Vercel — when opening the site from Vercel it may take a few moments to connect to the backend; please refresh the page a couple of times to ensure it becomes fully functional.

### Manual Setup
If you don't have `make` installed:

**Backend:**
```bash
go mod tidy
go run ./cmd/main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```
