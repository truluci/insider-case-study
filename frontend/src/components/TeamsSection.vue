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
        <div v-if="editingTeamId === team.id" class="edit-team-form" style="display: flex; flex-direction: column; gap: 0.5rem;">
          <input v-model="editForm.name" :disabled="team.is_default" placeholder="Team Name" type="text" style="width: 100%; box-sizing: border-box;">
          <input v-model.number="editForm.strength" placeholder="Strength (1-100)" type="number" min="1" max="100" style="width: 100%; box-sizing: border-box;">
          <div style="display: flex; gap: 0.5rem; margin-top: 0.5rem;">
            <button @click="saveEdit" class="btn btn-primary" style="flex: 1;">Save</button>
            <button @click="cancelEdit" class="btn" style="flex: 1;">Cancel</button>
          </div>
        </div>
        <div v-else>
          <h3>{{ team.name }} <span v-if="team.is_default" style="font-size: 0.7rem; background: var(--secondary-color); color: white; padding: 2px 6px; border-radius: 4px; vertical-align: middle;">Default</span></h3>
          <p>Strength: {{ team.strength }}</p>
          <div class="team-actions" style="margin-top: 1rem; display: flex; gap: 0.5rem;">
            <button @click="startEdit(team)" class="btn btn-small" style="background: transparent;">Edit</button>
            <button v-if="!team.is_default" @click="handleDelete(team.id)" class="btn btn-small btn-danger" :disabled="tournamentStarted" style="background: transparent;">Delete</button>
          </div>
        </div>
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
  emits: ['add-team', 'update-team', 'delete-team'],
  setup(props, { emit }) {
    const localNewTeam = ref({ name: '', strength: 50 })
    const editingTeamId = ref(null)
    const editForm = ref({ id: null, name: '', strength: 50 })

    const handleAddTeam = () => {
      emit('add-team', { name: localNewTeam.value.name, strength: localNewTeam.value.strength })
      localNewTeam.value = { name: '', strength: 50 }
    }

    const startEdit = (team) => {
      editingTeamId.value = team.id
      editForm.value = { id: team.id, name: team.name, strength: team.strength }
    }

    const saveEdit = () => {
      emit('update-team', { ...editForm.value })
      editingTeamId.value = null
    }

    const cancelEdit = () => {
      editingTeamId.value = null
    }

    const handleDelete = (id) => {
      emit('delete-team', id)
    }

    return {
      localNewTeam,
      editingTeamId,
      editForm,
      handleAddTeam,
      startEdit,
      saveEdit,
      cancelEdit,
      handleDelete
    }
  }
}
</script>
