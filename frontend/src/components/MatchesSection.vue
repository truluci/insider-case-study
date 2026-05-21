<template>
  <section class="section">
    <div class="week-header">
      <h2>Week {{ currentWeek }} of {{ totalWeeks }}</h2>
      <div class="week-controls">
        <button @click="$emit('play-week')" class="btn btn-success" :disabled="currentWeek > totalWeeks">Play Week {{ currentWeek }} Matches</button>
        <button @click="$emit('play-all')" class="btn btn-primary" :disabled="currentWeek > totalWeeks">Play All Matches</button>
        <button @click="$emit('next-week')" class="btn btn-info" :disabled="currentWeek >= totalWeeks">Next Week →</button>
      </div>
    </div>


    <div class="matches-list">
      <div v-if="matches.filter(m => m.week === currentWeek).length === 0" class="no-matches">
        <p>No matches for this week</p>
      </div>
      <div v-for="match in matches.filter(m => m.week === currentWeek)" :key="match.id" class="match-card">
        <div class="match-header">Week {{ match.week }}</div>
        
        <div v-if="editingMatchId === match.id" class="edit-match-form" style="display: flex; flex-direction: column; gap: 0.5rem; margin-bottom: 0.5rem;">
          <div style="display: flex; justify-content: space-between; align-items: center;">
            <span class="team">{{ getTeamName(match.home_team_id) }}</span>
            <input v-model.number="editForm.homeGoals" type="number" min="0" style="width: 50px; text-align: center; padding: 2px;">
            <span class="score vs">-</span>
            <input v-model.number="editForm.awayGoals" type="number" min="0" style="width: 50px; text-align: center; padding: 2px;">
            <span class="team" style="text-align: right;">{{ getTeamName(match.away_team_id) }}</span>
          </div>
          <div style="display: flex; gap: 0.5rem; margin-top: 0.5rem; justify-content: center;">
            <button @click="saveEdit" class="btn btn-primary btn-small" style="padding: 4px 8px; font-size: 0.9rem;">Save</button>
            <button @click="cancelEdit" class="btn btn-small" style="padding: 4px 8px; font-size: 0.9rem;">Cancel</button>
          </div>
        </div>

        <div v-else style="display: flex; flex-direction: column; height: 100%;">
          <div class="match-body">
            <span class="team">{{ getTeamName(match.home_team_id) }}</span>
            <span class="score" v-if="match.status === 'completed'">{{ match.home_goals }} - {{ match.away_goals }}</span>
            <span class="score vs" v-else>vs</span>
            <span class="team">{{ getTeamName(match.away_team_id) }}</span>
          </div>
          <div class="match-status" :class="match.status" style="display: flex; justify-content: space-between; align-items: center; margin-top: auto;">
            <span>{{ match.status }}</span>
            <button v-if="match.status === 'completed'" @click="startEdit(match)" class="btn btn-small" style="background: transparent; padding: 2px 8px; font-size: 0.8rem; border: 1px solid var(--border-color);">Edit</button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'MatchesSection',
  props: {
    matches: { type: Array, required: true },
    teams: { type: Array, required: true },
    currentWeek: { type: Number, required: true },
    totalWeeks: { type: Number, required: true }
  },
  emits: ['play-week', 'play-all', 'next-week', 'update-match'],
  setup(props, { emit }) {
    const editingMatchId = ref(null)
    const editForm = ref({ homeGoals: 0, awayGoals: 0 })

    const startEdit = (match) => {
      editingMatchId.value = match.id
      editForm.value = { homeGoals: match.home_goals, awayGoals: match.away_goals }
    }

    const saveEdit = () => {
      emit('update-match', { id: editingMatchId.value, ...editForm.value })
      editingMatchId.value = null
    }

    const cancelEdit = () => {
      editingMatchId.value = null
    }

    const getTeamName = (teamId) => {
      const team = props.teams.find(t => t.id === teamId)
      return team ? team.name : 'Unknown'
    }
    return { 
      getTeamName,
      editingMatchId,
      editForm,
      startEdit,
      saveEdit,
      cancelEdit
    }
  }
}
</script>
