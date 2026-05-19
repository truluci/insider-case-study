<template>
  <div class="container">
    <header class="header">
      <h1>Football League Simulation</h1>
      <p>Premier League Style Tournament</p>
      <button @click="restartSimulation" class="btn btn-danger">Restart Simulation</button>
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


