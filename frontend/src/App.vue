<template>
  <div class="container">
    <header class="header">
      <h1>⚽ Football League Simulation</h1>
      <p>Premier League Style Tournament</p>
      <button @click="restartSimulation" class="btn" style="background: #e53e3e; color: white; margin-top: 15px;">🔄 Restart Simulation</button>
    </header>

    <main class="main-content">
      <!-- Navigation Tabs -->
      <nav class="tabs">
        <button 
          v-for="tab in tabs" 
          :key="tab" 
          @click="activeTab = tab"
          :class="['tab-btn', { active: activeTab === tab }]"
        >
          {{ tab }}
        </button>
      </nav>

      <!-- Teams Section -->
      <section v-show="activeTab === 'Teams'" class="section">
        <h2>Teams ({{ teams.length }} takım)</h2>
        <div class="team-form">
          <input v-model="newTeam.name" placeholder="Team Name" type="text">
          <input v-model.number="newTeam.strength" placeholder="Strength (1-100)" type="number" min="1" max="100">
          <button @click="addTeam" class="btn btn-primary" :disabled="tournamentStarted">Add Team & Reschedule</button>
          <span v-if="tournamentStarted" class="warning-text">Tournament started - cannot add teams</span>
        </div>
        
        <div class="teams-grid">
          <div v-for="team in teams" :key="team.id" class="team-card">
            <h3>{{ team.name }}</h3>
            <p>Strength: {{ team.strength }}</p>
          </div>
        </div>
      </section>

      <!-- Matches Section -->
      <section v-show="activeTab === 'Matches'" class="section">
        <div class="week-header">
          <h2>Week {{ currentWeek }} of {{ totalWeeks }}</h2>
          <div class="week-controls">
            <button @click="playWeekMatches" class="btn btn-success" :disabled="currentWeek > totalWeeks">Play Week {{ currentWeek }} Matches</button>
            <button @click="goNextWeek" class="btn btn-info" :disabled="currentWeek >= totalWeeks">Next Week →</button>
          </div>
        </div>

        <div class="match-form">
          <input type="text" placeholder="(Optional) Add custom matches" disabled class="input-disabled">
          <p style="font-size: 0.9rem; color: #666; margin: 10px 0;">Week 1-{{ totalWeeks }} matches are pre-scheduled</p>
        </div>

        <div class="matches-list">
          <div v-if="matches.filter(m => m.week === currentWeek).length === 0" class="no-matches">
            <p>No matches for this week</p>
          </div>
          <div v-for="match in matches.filter(m => m.week === currentWeek)" :key="match.id" class="match-card">
            <div class="match-header">Week {{ match.week }}</div>
            <div class="match-body">
              <span class="team">{{ getTeamName(match.home_team_id) }}</span>
              <span class="score" v-if="match.status === 'completed'">{{ match.home_goals }} - {{ match.away_goals }}</span>
              <span class="score vs" v-else>vs</span>
              <span class="team">{{ getTeamName(match.away_team_id) }}</span>
            </div>
            <div class="match-status" :class="match.status">{{ match.status }}</div>
          </div>
        </div>
      </section>

      <!-- League Table Section -->
      <section v-show="activeTab === 'League'" class="section">
        <h2>League Table</h2>
        <table class="league-table">
          <thead>
            <tr>
              <th>Position</th>
              <th>Team</th>
              <th>P</th>
              <th>W</th>
              <th>D</th>
              <th>L</th>
              <th>GF</th>
              <th>GA</th>
              <th>GD</th>
              <th>Points</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(stat, idx) in leagueStats" :key="idx" :class="['row', `position-${stat.position}`]">
              <td class="position">{{ stat.position }}</td>
              <td class="team-name">{{ stat.team_name }}</td>
              <td>{{ stat.played }}</td>
              <td>{{ stat.won }}</td>
              <td>{{ stat.drawn }}</td>
              <td>{{ stat.lost }}</td>
              <td>{{ stat.goals_for }}</td>
              <td>{{ stat.goals_against }}</td>
              <td>{{ stat.goal_diff }}</td>
              <td class="points"><strong>{{ stat.points }}</strong></td>
            </tr>
          </tbody>
        </table>
      </section>

      <!-- Predictions Section -->
      <section v-show="activeTab === 'Predictions'" class="section">
        <h2>Final Predictions</h2>
        <div class="prediction-form">
          <select v-model.number="newPrediction.teamId">
            <option value="">Select Team</option>
            <option v-for="team in teams" :key="team.id" :value="team.id">{{ team.name }}</option>
          </select>
          <select v-model.number="newPrediction.position">
            <option value="">Predicted Position</option>
            <option value="1">1st Place</option>
            <option value="2">2nd Place</option>
            <option value="3">3rd Place</option>
            <option value="4">4th Place</option>
          </select>
          <button @click="addPrediction" class="btn btn-primary">Add Prediction</button>
        </div>

        <div class="predictions-list">
          <div v-for="pred in predictions" :key="pred.id" class="prediction-card">
            <span class="team-name">{{ pred.team_name }}</span>
            <span class="position-badge">{{ pred.position }}</span>
          </div>
        </div>
      </section>
    </main>

    <footer class="footer">
      <p>&copy; 2024 Football League Simulation | Built with Vue.js & Go</p>
    </footer>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const API_BASE = '/api'

export default {
  name: 'App',
  setup() {
    const activeTab = ref('Teams')
    const tabs = ['Teams', 'Matches', 'League', 'Predictions']
    const teams = ref([])
    const matches = ref([])
    const leagueStats = ref([])
    const predictions = ref([])
    const currentWeek = ref(1)
    const totalWeeks = ref(0)
    const tournamentStarted = ref(false)

    const newTeam = ref({
      name: '',
      strength: 50
    })

    const newMatch = ref({
      homeTeamId: '',
      awayTeamId: '',
      week: 1
    })

    const newPrediction = ref({
      teamId: '',
      position: ''
    })

    const fetchTeams = async () => {
      try {
        const response = await axios.get(`${API_BASE}/teams`)
        teams.value = response.data || []
      } catch (error) {
        console.error('Error fetching teams:', error)
      }
    }

    const fetchMatches = async () => {
      try {
        const response = await axios.get(`${API_BASE}/matches`)
        matches.value = response.data || []
      } catch (error) {
        console.error('Error fetching matches:', error)
      }
    }

    const fetchLeagueStats = async () => {
      try {
        const response = await axios.get(`${API_BASE}/league`)
        leagueStats.value = response.data || []
      } catch (error) {
        console.error('Error fetching league table:', error)
      }
    }

    const fetchPredictions = async () => {
      try {
        const response = await axios.get(`${API_BASE}/predictions`)
        predictions.value = response.data || []
      } catch (error) {
        console.error('Error fetching predictions:', error)
      }
    }

    const fetchTournamentState = async () => {
      try {
        const response = await axios.get(`${API_BASE}/tournament/current-week`)
        currentWeek.value = response.data.current_week
        totalWeeks.value = response.data.total_weeks
        
        // Check if tournament has started by checking if current week > 1
        tournamentStarted.value = currentWeek.value > 1
      } catch (error) {
        console.error('Error fetching tournament state:', error)
      }
    }

    const addTeam = async () => {
      if (!newTeam.value.name) {
        alert('Please enter team name')
        return
      }
      try {
        await axios.post(`${API_BASE}/teams`, {
          name: newTeam.value.name,
          strength: newTeam.value.strength
        })
        newTeam.value = { name: '', strength: 50 }
        await fetchTeams()
        await fetchTournamentState()
        await fetchMatches()
        await fetchLeagueStats()
        alert('Team added and tournament rescheduled!')
      } catch (error) {
        console.error('Error adding team:', error)
        if (error.response?.status === 400) {
          alert('Cannot add teams after tournament has started')
          tournamentStarted.value = true
        } else {
          alert('Failed to add team')
        }
      }
    }

    const addMatch = async () => {
      if (!newMatch.value.homeTeamId || !newMatch.value.awayTeamId) {
        alert('Please select both teams')
        return
      }
      try {
        await axios.post(`${API_BASE}/matches`, {
          home_team_id: newMatch.value.homeTeamId,
          away_team_id: newMatch.value.awayTeamId,
          week: newMatch.value.week
        })
        newMatch.value = { homeTeamId: '', awayTeamId: '', week: 1 }
        await fetchMatches()
      } catch (error) {
        console.error('Error adding match:', error)
        alert('Failed to add match')
      }
    }

    const playWeekMatches = async () => {
      try {
        await axios.post(`${API_BASE}/play-all`)
        await fetchMatches()
        await fetchLeagueStats()
      } catch (error) {
        console.error('Error playing matches:', error)
        alert('Failed to play matches')
      }
    }

    const goNextWeek = async () => {
      try {
        await axios.post(`${API_BASE}/tournament/next-week`)
        await fetchTournamentState()
        await fetchMatches()
        await fetchLeagueStats()
      } catch (error) {
        console.error('Error advancing week:', error)
        alert('Cannot advance to next week')
      }
    }

    const addPrediction = async () => {
      if (!newPrediction.value.teamId || !newPrediction.value.position) {
        alert('Please fill all fields')
        return
      }
      try {
        await axios.post(`${API_BASE}/predictions`, {
          week: 1,
          team_id: newPrediction.value.teamId,
          position: newPrediction.value.position
        })
        newPrediction.value = { teamId: '', position: '' }
        await fetchPredictions()
      } catch (error) {
        console.error('Error adding prediction:', error)
        alert('Failed to add prediction')
      }
    }

    const restartSimulation = async () => {
      if (!confirm('Are you sure you want to completely restart the simulation? All match results and predictions will be lost.')) {
        return
      }
      try {
        await axios.post(`${API_BASE}/restart`)
        await fetchTournamentState()
        await fetchMatches()
        await fetchLeagueStats()
        await fetchPredictions()
        activeTab.value = 'Teams'
        alert('Simulation has been restarted!')
      } catch (error) {
        console.error('Error restarting simulation:', error)
        alert('Failed to restart simulation')
      }
    }

    const getTeamName = (teamId) => {
      const team = teams.value.find(t => t.id === teamId)
      return team ? team.name : 'Unknown'
    }

    onMounted(() => {
      fetchTournamentState()
      fetchTeams()
      fetchMatches()
      fetchLeagueStats()
      fetchPredictions()
    })

    return {
      activeTab,
      tabs,
      teams,
      matches,
      leagueStats,
      predictions,
      currentWeek,
      totalWeeks,
      tournamentStarted,
      newTeam,
      newMatch,
      newPrediction,
      addTeam,
      addMatch,
      playWeekMatches,
      goNextWeek,
      addPrediction,
      getTeamName,
      restartSimulation
    }
  }
}
</script>

<style scoped>
.container {
  max-width: 1200px;
  margin: 0 auto;
  background: white;
  border-radius: 10px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  overflow: hidden;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 40px 20px;
  text-align: center;
}

.header h1 {
  font-size: 2.5rem;
  margin-bottom: 10px;
}

.header p {
  font-size: 1.1rem;
  opacity: 0.9;
}

.main-content {
  padding: 30px;
}

.tabs {
  display: flex;
  gap: 10px;
  margin-bottom: 30px;
  border-bottom: 2px solid #eee;
}

.tab-btn {
  padding: 12px 24px;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 500;
  color: #666;
  border-bottom: 3px solid transparent;
  transition: all 0.3s ease;
}

.tab-btn:hover {
  color: #667eea;
}

.tab-btn.active {
  color: #667eea;
  border-bottom-color: #667eea;
}

.section {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.section h2 {
  margin-bottom: 20px;
  color: #333;
}

/* Forms */
.team-form, .match-form, .prediction-form {
  display: flex;
  gap: 10px;
  margin-bottom: 30px;
  flex-wrap: wrap;
}

input, select {
  padding: 10px 15px;
  border: 1px solid #ddd;
  border-radius: 5px;
  font-size: 1rem;
}

input:focus, select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

/* Buttons */
.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-primary {
  background: #667eea;
  color: white;
}

.btn-primary:hover {
  background: #5568d3;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.btn-success {
  background: #48bb78;
  color: white;
}

.btn-success:hover {
  background: #38a169;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(72, 187, 120, 0.4);
}

.btn-info {
  background: #3182ce;
  color: white;
}

.btn-info:hover {
  background: #2c5aa0;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(49, 130, 206, 0.4);
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.warning-text {
  color: #e53e3e;
  font-size: 0.9rem;
  font-weight: 600;
  margin-left: 10px;
}

.week-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  gap: 20px;
  flex-wrap: wrap;
}

.week-header h2 {
  margin: 0;
  color: #667eea;
  font-size: 1.8rem;
}

.week-controls {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.input-disabled {
  opacity: 0.6;
  cursor: not-allowed;
  background-color: #f5f5f5;
}

/* Teams Grid */
.teams-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
}

.team-card {
  background: #f7fafc;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.team-card:hover {
  border-color: #667eea;
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.1);
  transform: translateY(-2px);
}

.team-card h3 {
  color: #667eea;
  margin-bottom: 10px;
}

.team-card p {
  color: #666;
}

/* Matches List */
.matches-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 15px;
}

.match-card {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  overflow: hidden;
  background: white;
}

.match-header {
  background: #edf2f7;
  padding: 10px 15px;
  font-weight: 600;
  color: #667eea;
}

.match-body {
  padding: 15px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.team {
  font-weight: 600;
  color: #333;
  flex: 1;
}

.score {
  margin: 0 15px;
  font-size: 1.2rem;
  font-weight: 700;
  color: #667eea;
  min-width: 60px;
  text-align: center;
}

.vs {
  color: #999;
  font-size: 0.9rem;
}

.match-status {
  padding: 5px 15px;
  text-align: center;
  font-size: 0.85rem;
  font-weight: 600;
  background: #f0f0f0;
  color: #666;
  text-transform: capitalize;
}

.match-status.completed {
  background: #c6f6d5;
  color: #22543d;
}

.match-status.scheduled {
  background: #bee3f8;
  color: #2c5282;
}

.no-matches {
  text-align: center;
  padding: 40px 20px;
  color: #999;
  background: #f9f9f9;
  border-radius: 8px;
  border: 1px dashed #ddd;
}

.no-matches p {
  margin: 0;
  font-size: 1.1rem;
}

/* League Table */
.league-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

.league-table th {
  background: #667eea;
  color: white;
  padding: 12px;
  text-align: left;
  font-weight: 600;
}

.league-table td {
  padding: 12px;
  border-bottom: 1px solid #e2e8f0;
}

.league-table tbody tr {
  transition: background-color 0.2s ease;
}

.league-table tbody tr:hover {
  background: #f7fafc;
}

.league-table .position {
  font-weight: 700;
  color: #667eea;
  min-width: 50px;
}

.league-table .team-name {
  font-weight: 600;
}

.league-table .points {
  color: #667eea;
  font-weight: 700;
}

/* Predictions */
.predictions-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
}

.prediction-card {
  background: #edf2f7;
  border: 2px solid #cbd5e0;
  border-radius: 8px;
  padding: 15px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.prediction-card .team-name {
  font-weight: 600;
  color: #333;
}

.position-badge {
  background: #667eea;
  color: white;
  border-radius: 50%;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
}

/* Footer */
.footer {
  background: #f7fafc;
  padding: 20px;
  text-align: center;
  color: #666;
  border-top: 1px solid #e2e8f0;
}

/* Responsive */
@media (max-width: 768px) {
  .header h1 {
    font-size: 1.8rem;
  }

  .teams-grid,
  .matches-list,
  .predictions-list {
    grid-template-columns: 1fr;
  }

  .tabs {
    flex-wrap: wrap;
  }

  .team-form, .match-form, .prediction-form {
    flex-direction: column;
  }

  .team-form input,
  .match-form select,
  .prediction-form select,
  .btn {
    width: 100%;
  }

  .league-table {
    font-size: 0.9rem;
  }

  .league-table th,
  .league-table td {
    padding: 8px;
  }
}
</style>
