<template>
  <section class="section">
    <h2>Final Predictions</h2>
    <div class="prediction-form" v-if="!arePredictionsLocked">
      <select v-model.number="localPrediction.teamId">
        <option value="">Select Team</option>
        <option v-for="team in teams" :key="team.id" :value="team.id">{{ team.name }}</option>
      </select>
      <select v-model.number="localPrediction.position">
        <option value="">Predicted Position</option>
        <option v-for="n in teams.length" :key="n" :value="n">{{ n }}{{ n === 1 ? 'st' : n === 2 ? 'nd' : n === 3 ? 'rd' : 'th' }} Place</option>
      </select>
      <button @click="handleAddPrediction" class="btn btn-primary">Add Prediction</button>
    </div>
    <div class="warning-text" v-else style="margin-bottom: 20px;">
      <p>The league has started! Predictions are locked.</p>
    </div>

    <div class="predictions-list">
      <div v-for="pred in predictions" :key="pred.id" :class="['prediction-card', { 'correct-prediction': correctPredictionIds && correctPredictionIds.includes(pred.id) }]">
        <span class="team-name">{{ getTeamName(pred.team_id) }}</span>
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
    predictions: { type: Array, required: true },
    correctPredictionIds: { type: Array, default: () => [] },
    arePredictionsLocked: { type: Boolean, default: false }
  },
  emits: ['add-prediction'],
  setup(props, { emit }) {
    const localPrediction = ref({ teamId: '', position: '' })

    const handleAddPrediction = () => {
      emit('add-prediction', { teamId: localPrediction.value.teamId, position: localPrediction.value.position })
      localPrediction.value = { teamId: '', position: '' }
    }

    const getTeamName = (teamId) => {
      const team = props.teams.find(t => t.id === teamId)
      return team ? team.name : 'Unknown'
    }

    return {
      localPrediction,
      handleAddPrediction,
      getTeamName
    }
  }
}
</script>
