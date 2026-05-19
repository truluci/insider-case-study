<template>
  <section class="section">
    <h2>Final Predictions</h2>
    <div class="prediction-form">
      <select v-model.number="localPrediction.teamId">
        <option value="">Select Team</option>
        <option v-for="team in teams" :key="team.id" :value="team.id">{{ team.name }}</option>
      </select>
      <select v-model.number="localPrediction.position">
        <option value="">Predicted Position</option>
        <option value="1">1st Place</option>
        <option value="2">2nd Place</option>
        <option value="3">3rd Place</option>
        <option value="4">4th Place</option>
      </select>
      <button @click="handleAddPrediction" class="btn btn-primary">Add Prediction</button>
    </div>

    <div class="predictions-list">
      <div v-for="pred in predictions" :key="pred.id" class="prediction-card">
        <span class="team-name">{{ pred.team_name }}</span>
        <span class="position-badge">{{ pred.position }}</span>
      </div>
    </div>
  </section>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'PredictionsSection',
  props: {
    teams: { type: Array, required: true },
    predictions: { type: Array, required: true }
  },
  emits: ['add-prediction'],
  setup(props, { emit }) {
    const localPrediction = ref({ teamId: '', position: '' })

    const handleAddPrediction = () => {
      emit('add-prediction', { teamId: localPrediction.value.teamId, position: localPrediction.value.position })
      localPrediction.value = { teamId: '', position: '' }
    }

    return {
      localPrediction,
      handleAddPrediction
    }
  }
}
</script>
