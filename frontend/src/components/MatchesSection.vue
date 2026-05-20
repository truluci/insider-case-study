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
</template>

<script>
export default {
  name: 'MatchesSection',
  props: {
    matches: { type: Array, required: true },
    teams: { type: Array, required: true },
    currentWeek: { type: Number, required: true },
    totalWeeks: { type: Number, required: true }
  },
  emits: ['play-week', 'play-all', 'next-week'],
  setup(props) {
    const getTeamName = (teamId) => {
      const team = props.teams.find(t => t.id === teamId)
      return team ? team.name : 'Unknown'
    }
    return { getTeamName }
  }
}
</script>
