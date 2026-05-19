-- Teams tablosu
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    strength INTEGER NOT NULL DEFAULT 50 CHECK (strength >= 1 AND strength <= 100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Matches tablosu
CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    home_team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    away_team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    week INTEGER NOT NULL CHECK (week >= 1),
    home_goals INTEGER DEFAULT 0 CHECK (home_goals >= 0),
    away_goals INTEGER DEFAULT 0 CHECK (away_goals >= 0),
    status VARCHAR(50) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'completed', 'played')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(home_team_id, away_team_id, week)
);

-- Team istatistikleri
CREATE TABLE IF NOT EXISTS team_stats (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL UNIQUE REFERENCES teams(id) ON DELETE CASCADE,
    played INTEGER DEFAULT 0 CHECK (played >= 0),
    won INTEGER DEFAULT 0 CHECK (won >= 0),
    drawn INTEGER DEFAULT 0 CHECK (drawn >= 0),
    lost INTEGER DEFAULT 0 CHECK (lost >= 0),
    goals_for INTEGER DEFAULT 0 CHECK (goals_for >= 0),
    goals_against INTEGER DEFAULT 0 CHECK (goals_against >= 0),
    goal_diff INTEGER DEFAULT 0,
    points INTEGER DEFAULT 0 CHECK (points >= 0),
    last_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Predictions tablosu
CREATE TABLE IF NOT EXISTS predictions (
    id SERIAL PRIMARY KEY,
    week INTEGER NOT NULL CHECK (week >= 1),
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    position INTEGER NOT NULL CHECK (position >= 1 AND position <= 4),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_matches_week ON matches(week);
CREATE INDEX idx_matches_home_team ON matches(home_team_id);
CREATE INDEX idx_matches_away_team ON matches(away_team_id);
CREATE INDEX idx_predictions_week ON predictions(week);
CREATE INDEX idx_predictions_team ON predictions(team_id);
