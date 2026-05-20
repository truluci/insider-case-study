<template>
  <section class="section">
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
          <th>Win %</th>
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
          <td class="prediction-perc">{{ stat.win_chance === 'N' ? 'N' : '%' + stat.win_chance }}</td>
        </tr>
      </tbody>
    </table>
  </section>
</template>

<script>
export default {
  name: 'LeagueTable',
  props: {
    leagueStats: {
      type: Array,
      required: true
    }
  }
}
</script>
