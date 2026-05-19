<template>
  <section class="section">
    <h2>Teams ({{ teams.length }} takım)</h2>
    <div class="team-form">
      <input v-model="localNewTeam.name" placeholder="Team Name" type="text">
      <input v-model.number="localNewTeam.strength" placeholder="Strength (1-100)" type="number" min="1" max="100">
      <button @click="handleAddTeam" class="btn btn-primary" :disabled="tournamentStarted">Add Team & Reschedule</button>
      <span v-if="tournamentStarted" class="warning-text">Tournament started - cannot add teams</span>
    </div>
    
    <div class="teams-grid">
      <div v-for="team in teams" :key="team.id" class="team-card">
        <h3>{{ team.name }}</h3>
        <p>Strength: {{ team.strength }}</p>
      </div>
    </div>
  </section>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'TeamsSection',
  props: {
    teams: {
      type: Array,
      required: true
    },
    tournamentStarted: {
      type: Boolean,
      required: true
    }
  },
  emits: ['add-team'],
  setup(props, { emit }) {
    const localNewTeam = ref({ name: '', strength: 50 })

    const handleAddTeam = () => {
      emit('add-team', { name: localNewTeam.value.name, strength: localNewTeam.value.strength })
      localNewTeam.value = { name: '', strength: 50 }
    }

    return {
      localNewTeam,
      handleAddTeam
    }
  }
}
</script>
