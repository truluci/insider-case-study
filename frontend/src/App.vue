<template>
  <div class="container">
    <div class="notifications-container">
      <div v-for="notif in notifications" :key="notif.id" :class="['notification', notif.type]">
        {{ notif.message }}
      </div>
    </div>
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
        @update-match="updateMatchResult"
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

const API_BASE = import.meta.env.VITE_API_BASE || '/api'

export default {
  name: 'App',
  components: {
    TeamsSection,
    MatchesSection,
    LeagueTable,
    PredictionsSection
  },
  setup() {
    const notifications = ref([])
    let notifId = 0
    const showNotification = (message, type = 'info') => {
      const id = notifId++
      notifications.value.push({ id, message, type })
      setTimeout(() => {
        notifications.value = notifications.value.filter(n => n.id !== id)
      }, 5000)
    }

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
        showNotification('Please enter team name', 'error')
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
        showNotification('Team added and tournament rescheduled!', 'success')
      } catch (error) {
        console.error('Error adding team:', error)
        if (error.response?.status === 400) {
          showNotification('Cannot add teams after tournament has started', 'error')
          tournamentStarted.value = true
        } else {
          showNotification('Failed to add team', 'error')
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
        showNotification('Team updated!', 'success')
      } catch (error) {
        console.error('Error updating team:', error)
        showNotification(error.response?.data || 'Failed to update team', 'error')
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
        showNotification('Team deleted and tournament rescheduled!', 'success')
      } catch (error) {
        console.error('Error deleting team:', error)
        if (error.response?.status === 400) {
          showNotification('Cannot delete teams after tournament has started', 'error')
          tournamentStarted.value = true
        } else {
          showNotification(error.response?.data || 'Failed to delete team', 'error')
        }
      }
    }

    const updateMatchResult = async (matchData) => {
      try {
        await axios.put(`${API_BASE}/matches/${matchData.id}`, {
          home_goals: matchData.homeGoals,
          away_goals: matchData.awayGoals,
          status: 'completed'
        })
        await fetchMatches()
        await fetchLeagueStats()
        checkPredictions()
        showNotification('Match result updated successfully', 'success')
      } catch (error) {
        console.error('Error updating match result:', error)
        showNotification(error.response?.data || 'Failed to update match result', 'error')
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
        showNotification("🎉 Congratulations!\n\n" + correctMessages.join('\n'), 'success');
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
        showNotification('Failed to play matches', 'error')
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
        showNotification('Failed to play all matches', 'error')
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
        showNotification('Cannot advance to next week', 'error')
      }
    }

    const addPrediction = async (predictionData) => {
      if (arePredictionsLocked.value) {
        showNotification('Cannot add predictions after the league has started!', 'error')
        return
      }

      if (!predictionData.teamId || !predictionData.position) {
        showNotification('Please fill all fields', 'error')
        return
      }

      // Check if this exact prediction already exists
      const exists = predictions.value.find(p => p.team_id === predictionData.teamId && p.position === predictionData.position)
      if (exists) {
        showNotification('This prediction already exists!', 'error')
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
        showNotification('Failed to add prediction', 'error')
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
        showNotification('Simulation has been restarted!', 'success')
      } catch (error) {
        console.error('Error restarting simulation:', error)
        showNotification('Failed to restart simulation', 'error')
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
      notifications,
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
      updateMatchResult,
      addPrediction,
      restartSimulation
    }
  }
}
</script>
