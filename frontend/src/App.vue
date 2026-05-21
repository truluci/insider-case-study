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

      <!-- Sections -->
      <TeamsSection 
        v-show="activeTab === 'Teams'"
        :teams="teams"
        :tournamentStarted="tournamentStarted"
        @add-team="addTeam"
        @update-team="updateTeam"
        @delete-team="deleteTeam"
      />

      <MatchesSection
        v-show="activeTab === 'Matches'"
        :matches="matches"
        :teams="teams"
        :currentWeek="currentWeek"
        :totalWeeks="totalWeeks"
        @play-week="playWeekMatches"
        @play-all="playAllWeeksMatches"
        @next-week="goNextWeek"
      />

      <LeagueTable
        v-show="activeTab === 'League'"
        :leagueStats="leagueStats"
      />

      <PredictionsSection
        v-show="activeTab === 'Predictions'"
        :teams="teams"
        :predictions="predictions"
        :correctPredictionIds="correctPredictionIds"
        :arePredictionsLocked="arePredictionsLocked"
        @add-prediction="addPrediction"
      />
    </main>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import TeamsSection from './components/TeamsSection.vue'
import MatchesSection from './components/MatchesSection.vue'
import LeagueTable from './components/LeagueTable.vue'
import PredictionsSection from './components/PredictionsSection.vue'

const API_BASE = '/api'

export default {
  name: 'App',
  components: {
    TeamsSection,
    MatchesSection,
    LeagueTable,
    PredictionsSection
  },
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

    const arePredictionsLocked = computed(() => {
      return tournamentStarted.value || (matches.value && matches.value.some(m => m.status === 'completed'))
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

    const addTeam = async (teamData) => {
      if (!teamData.name) {
        alert('Please enter team name')
        return
      }
      try {
        await axios.post(`${API_BASE}/teams`, {
          name: teamData.name,
          strength: teamData.strength
        })
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

    const updateTeam = async (teamData) => {
      try {
        await axios.put(`${API_BASE}/teams/${teamData.id}`, {
          name: teamData.name,
          strength: teamData.strength
        })
        await fetchTeams()
        await fetchMatches()
        await fetchLeagueStats()
        await fetchPredictions()
        alert('Team updated!')
      } catch (error) {
        console.error('Error updating team:', error)
        alert(error.response?.data || 'Failed to update team')
      }
    }

    const deleteTeam = async (teamId) => {
      if (!confirm('Are you sure you want to delete this team?')) return
      try {
        await axios.delete(`${API_BASE}/teams/${teamId}`)
        await fetchTeams()
        await fetchTournamentState()
        await fetchMatches()
        await fetchLeagueStats()
        alert('Team deleted and tournament rescheduled!')
      } catch (error) {
        console.error('Error deleting team:', error)
        if (error.response?.status === 400) {
          alert('Cannot delete teams after tournament has started')
          tournamentStarted.value = true
        } else {
          alert(error.response?.data || 'Failed to delete team')
        }
      }
    }

    let hasAlertedPredictions = false;
    const correctPredictionIds = ref([]);

    const checkPredictions = () => {
      if (!matches.value || matches.value.length === 0) return;
      if (!predictions.value || predictions.value.length === 0) return;
      const allCompleted = matches.value.every(m => m.status === 'completed');
      if (!allCompleted) {
        correctPredictionIds.value = [];
        return;
      }

      const correctIds = [];
      const correctMessages = [];
      predictions.value.forEach(pred => {
        const statIndex = leagueStats.value.findIndex(s => s.team_id === pred.team_id);
        if (statIndex !== -1) {
          const actualPosition = statIndex + 1;
          if (actualPosition === pred.position) {
            correctIds.push(pred.id);
            const teamName = teams.value.find(t => t.id === pred.team_id)?.name || 'Team';
            correctMessages.push(`You correctly predicted ${teamName} at position ${pred.position}!`);
          }
        }
      });

      correctPredictionIds.value = correctIds;

      if (!hasAlertedPredictions && correctMessages.length > 0) {
        alert("🎉 Congratulations!\n\n" + correctMessages.join('\n'));
      }
      hasAlertedPredictions = true;
    }

    const playWeekMatches = async () => {
      try {
        await axios.post(`${API_BASE}/play-week`)
        await fetchMatches()
        await fetchLeagueStats()
        checkPredictions()
      } catch (error) {
        console.error('Error playing matches:', error)
        alert('Failed to play matches')
      }
    }

    const playAllWeeksMatches = async () => {
      try {
        await axios.post(`${API_BASE}/play-all`)
        await fetchTournamentState()
        await fetchMatches()
        await fetchLeagueStats()
        checkPredictions()
      } catch (error) {
        console.error('Error playing all matches:', error)
        alert('Failed to play all matches')
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

    const addPrediction = async (predictionData) => {
      if (arePredictionsLocked.value) {
        alert('Cannot add predictions after the league has started!')
        return
      }

      if (!predictionData.teamId || !predictionData.position) {
        alert('Please fill all fields')
        return
      }

      // Check if this exact prediction already exists
      const exists = predictions.value.find(p => p.team_id === predictionData.teamId && p.position === predictionData.position)
      if (exists) {
        alert('This prediction already exists!')
        return
      }

      try {
        await axios.post(`${API_BASE}/predictions`, {
          week: 1,
          team_id: predictionData.teamId,
          position: predictionData.position
        })
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
        hasAlertedPredictions = false;
        correctPredictionIds.value = [];
        alert('Simulation has been restarted!')
      } catch (error) {
        console.error('Error restarting simulation:', error)
        alert('Failed to restart simulation')
      }
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
      correctPredictionIds,
      arePredictionsLocked,
      currentWeek,
      totalWeeks,
      tournamentStarted,
      addTeam,
      updateTeam,
      deleteTeam,
      playWeekMatches,
      playAllWeeksMatches,
      goNextWeek,
      addPrediction,
      restartSimulation
    }
  }
}
</script>
