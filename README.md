# ⚽ Football League Simulation

Futbol ligi simülasyonu projesi. 4 takımın ligi oynayıp sonuçlarını tahmin etme simülasyonu.

**Teknolojiler**: Go 1.23 + Vue.js 3  
**Mimari**: Clean Architecture + Interface-Based Design

---

## 🚀 Hızlı Başlangıç

### Gereksinimler
- Go 1.23+
- Node.js 18+
- Docker & Docker Compose (opsiyonel)

### Kurulum ve Çalıştırma

**Seçenek 1: Makefile kullan (önerilen)**
```bash
make install  # Tüm bağımlılıkları kur
make dev      # Dev ortamında çalıştır (backend + frontend)
```

**Seçenek 2: Manual kurulum**

Backend'i başlat:
```bash
go mod tidy
go run ./cmd/main.go
```
Server `http://localhost:8080` adresinde başlayacak.

Frontend'i başlat (yeni terminal):
```bash
cd frontend
npm install
npm run dev
```
Frontend `http://localhost:3000` adresinde çalışacak.

**Seçenek 3: Docker ile**
```bash
make docker
```

---

## 📁 Proje Yapısı

```
insider-case-study/
│
├── cmd/                        # Backend entry point
│   └── main.go
│
├── pkg/                        # Uygulamanın çekirdeği
│   ├── models/                 # Domain models ve interfaces
│   │   ├── team.go            # Team entity ve interface'ler
│   │   ├── match.go           # Match entity ve interface'ler
│   │   └── prediction.go      # Prediction entity ve interface'ler
│   │
│   ├── services/               # Business logic katmanı
│   │   ├── team_service.go    # Team işletme logic'i
│   │   ├── match_service.go   # Match simülasyonu ve logic'i
│   │   ├── league_service.go  # Lig tablosu hesapları
│   │   ├── prediction_service.go
│   │   ├── team_repository.go      # Data access layer
│   │   ├── match_repository.go
│   │   └── prediction_repository.go
│   │
│   ├── handlers/               # HTTP API katmanı
│   │   ├── handler.go         # Handler initialize
│   │   ├── team_handlers.go
│   │   ├── match_handlers.go
│   │   ├── league_handlers.go
│   │   └── prediction_handlers.go
│   │
│   ├── database/               # Database katmanı
│   │   └── database.go        # Connection ve schema
│   │
│   └── utils/                  # Yardımcı fonksiyonlar
│
├── config/
│   └── config.go              # Configuration yönetimi
│
├── frontend/                   # Vue.js SPA
│   ├── src/
│   │   ├── App.vue            # Ana component (Teams, Matches, League, Predictions)
│   │   └── main.js
│   ├── index.html
│   ├── vite.config.js
│   ├── package.json
│   ├── Dockerfile             # Frontend container
│   └── README.md
│
├── sql/
│   └── schema.sql             # PostgreSQL schema
│
├── go.mod                      # Go module definition
├── go.sum                      # Dependency hashes
│
├── docker-compose.yml         # Multi-container orchestration
├── Dockerfile                 # Backend container
├── Makefile                   # Build commands
│
├── .env                       # Environment variables (git ignored)
├── .env.example               # Template
├── .gitignore                 # Git ignore rules
│
└── README.md                  # Bu dosya
```

---

## 🏗️ Mimari Tasarım

### Clean Architecture Katmanları

```
┌──────────────────────────────────────┐
│  HTTP Handlers (main.go routes)     │ ← Request/Response
├──────────────────────────────────────┤
│  Services (Business Logic)           │ ← İş mantığı
├──────────────────────────────────────┤
│  Repositories (Data Access)          │ ← Database işlemleri
├──────────────────────────────────────┤
│  Database/Storage Layer (SQLite)     │ ← Persistence
└──────────────────────────────────────┘
```

### Interface-Based Design

Tüm bağımlılıklar interface'ler üzerinden tanımlanmıştır:

**Models/Interfaces** (pkg/models/)
```go
type TeamService interface {
    CreateTeam(name string, strength int) (*Team, error)
    GetTeam(id int) (*Team, error)
    GetAllTeams() ([]*Team, error)
}

type MatchService interface {
    PlayMatch(id int) (*Match, error)
    PlayAllMatches() error
}

type TeamRepository interface {
    Create(team *Team) (*Team, error)
    GetAll() ([]*Team, error)
}
```

**Implementation** (pkg/services/)
```go
type TeamServiceImpl struct {
    repo models.TeamRepository
}

type MatchServiceImpl struct {
    matchRepo models.MatchRepository
    teamRepo  models.TeamRepository
}
```

---

## 📊 API Endpoints

### Teams
```
GET    /api/teams              - Tüm takımları listele
POST   /api/teams              - Yeni takım oluştur
```
**Request Body (POST):**
```json
{
  "name": "Chelsea",
  "strength": 85
}
```

### Matches
```
GET    /api/matches            - Tüm maçları listele
POST   /api/matches            - Yeni maç oluştur
PUT    /api/matches/{id}       - Maç sonucunu güncelle
POST   /api/play-all           - Tüm maçları otomatik oyna
```
**Request Body (POST):**
```json
{
  "home_team_id": 1,
  "away_team_id": 2,
  "week": 1
}
```

### League Table
```
GET    /api/league             - Güncel lig tablosu
GET    /api/league/week/{week} - Belirli haftanın lig tablosu
```

### Predictions
```
GET    /api/predictions        - Tüm tahminleri listele
```

---

## 🎮 League Kuralları (Premier League Standardı)

| Kural | Değer |
|-------|-------|
| Takım Sayısı | 4 |
| Galibiyet Puanı | 3 |
| Beraberlik Puanı | 1 |
| Kayıp Puanı | 0 |

**Sıralama Kriterleri:**
1. Toplam Puan (P) ⬇
2. Gol Farkı (GD) ⬇
3. Attığı Gol (GF) ⬇

**Match Result Simulation:**
- Takım gücü (strength) temelinde algoritm
- Random factor ile realizm
- Home team advantage (optional future)

---

## 🎮 Simulation Flow

1. **Takımları Oluştur** 
   - "Teams" sekmesine git
   - Her takım için isim ve güç (1-100) gir
   - "Add Team" butonuna tıkla

2. **Maçları Planla**
   - "Matches" sekmesine git
   - Home/Away takımları seç
   - Hafta seçimi yap
   - "Add Match" butonuna tıkla

3. **Otomatik Oyna**
   - "Play All Matches" butonuna tıkla
   - Sistem tüm scheduled maçları oyna

4. **Lig Tablosunu Görüntüle**
   - "League" sekmesine git
   - Güncel standings'i göster

5. **Tahmin Et**
   - "Predictions" sekmesine git
   - Her takım için final pozisyonunu tahmin et

---

## 🔧 Environment Konfigürasyonu

`.env` dosyasını düzenle:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=football_league

# Server Configuration
SERVER_PORT=8080

# Environment Type
ENVIRONMENT=development  # production|development
```

Development için SQLite otomatik kullanılır.

---

## 🐳 Docker ile Çalıştırma

```bash
docker-compose up --build
```

**Services:**
- **app**: Backend Go server (8080)
- **db**: PostgreSQL (5432)
- **frontend**: Vue.js dev server (3000)

---

## 📝 Veritabanı Şeması

### Teams Table
```sql
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    strength INTEGER (1-100),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### Matches Table
```sql
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    home_team_id INTEGER REFERENCES teams(id),
    away_team_id INTEGER REFERENCES teams(id),
    week INTEGER,
    home_goals INTEGER,
    away_goals INTEGER,
    status VARCHAR(50),  -- scheduled|completed|played
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### Team Stats Table
```sql
CREATE TABLE team_stats (
    id SERIAL PRIMARY KEY,
    team_id INTEGER UNIQUE REFERENCES teams(id),
    played INTEGER,
    won INTEGER,
    drawn INTEGER,
    lost INTEGER,
    goals_for INTEGER,
    goals_against INTEGER,
    goal_diff INTEGER,
    points INTEGER,
    last_updated_at TIMESTAMP
);
```

### Predictions Table
```sql
CREATE TABLE predictions (
    id SERIAL PRIMARY KEY,
    week INTEGER,
    team_id INTEGER REFERENCES teams(id),
    position INTEGER (1-4),  -- Final predicted position
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

Detaylı schema: **sql/schema.sql**

---

## 📖 Özet Komutlar

```bash
# Kurulum
make install          # Tüm bağımlılıkları kur

# Geliştirme
make dev              # Backend + Frontend dev mode
make run              # Sadece backend
make frontend         # Sadece frontend

# Build & Deploy
make build            # Production build (backend + frontend)
make docker           # Docker Compose ile çalıştır

# Temizleme
make clean            # Build dosyalarını ve DB'yi sil
```

---

## ✨ Modular Code Prensipleri

✅ **Interface-Based Design**: Tüm servisleri interface'ler aracılığıyla tanımla  
✅ **Repository Pattern**: Data access logic ayrıştırılmış  
✅ **Service Layer**: Business logic ayrıştırılmış  
✅ **Dependency Injection**: Handler'lara dependencies inject edilir  
✅ **Single Responsibility**: Her dosyanın tek sorumluluğu  
✅ **Separation of Concerns**: Her katının net görevi

---

## 🔄 Veri Akışı

```
Frontend (Vue.js)
       ↓
HTTP Request (Axios)
       ↓
Handler (pkg/handlers/)
       ↓
Service (pkg/services/)
       ↓
Repository (pkg/services/*_repository.go)
       ↓
Database (SQLite/PostgreSQL)
       ↓
Response back to Frontend
```

---

## 🚧 Gelecek Geliştirmeler

- [ ] Websocket - Real-time match updates
- [ ] Unit Tests - 80%+ coverage
- [ ] Integration Tests
- [ ] Swagger API Dokumentasyonu
- [ ] Authentication/Authorization
- [ ] Advanced Analytics Dashboard
- [ ] Match editing ve recalculation
- [ ] Statik file serving (binary embed)
- [ ] Middleware (logging, CORS, rate limiting)
- [ ] Graceful shutdown
- [ ] Health checks

---

## 💾 SQL Queries Örneği

**Get League Table with Stats:**
```sql
SELECT 
    t.id, t.name,
    ts.played, ts.won, ts.drawn, ts.lost,
    ts.goals_for, ts.goals_against, ts.goal_diff,
    ts.points,
    ROW_NUMBER() OVER (ORDER BY ts.points DESC, ts.goal_diff DESC) as position
FROM teams t
JOIN team_stats ts ON t.id = ts.team_id
ORDER BY position;
```

**Get Matches with Team Names:**
```sql
SELECT 
    m.id, m.week,
    ht.name as home_team, at.name as away_team,
    m.home_goals, m.away_goals, m.status
FROM matches m
JOIN teams ht ON m.home_team_id = ht.id
JOIN teams at ON m.away_team_id = at.id
ORDER BY m.week, m.id;
```

---

## 🤝 Kod Kalitesi

- **Go Fmt**: Kod formatlama
- **Go Vet**: Potansiyel hataları kontrol
- **Interfaces**: Loosely coupled design
- **Error Handling**: Tüm hataları handle et
- **Comments**: Export edilen fonksiyonları document et

---

## 📄 Lisans

MIT

---

## 👨‍💻 Geliştirici Notları

- SQLite development için kullanılır, production'da PostgreSQL önerilir
- `main.go` modülü başlatma ve routing'i yönetir
- Handler'lar HTTP request/response döngüsünü yönetir
- Service'ler business logic'i yönetir
- Repository'ler database işlemlerini yönetir
- Model'ler domain objelerini ve interface'leri tanımlar

**Yeni Feature Ekleme Checklist:**
- [ ] Model/Interface (pkg/models/)
- [ ] Repository (pkg/services/*_repository.go)
- [ ] Service (pkg/services/*_service.go)
- [ ] Handler (pkg/handlers/*_handlers.go)
- [ ] Route (cmd/main.go setupRouter)
- [ ] Frontend Component (frontend/src/App.vue)

---

**Başlangıç Tarihi**: May 2024  
**Son Güncelleme**: May 2024
