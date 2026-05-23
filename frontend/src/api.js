import axios from 'axios'

const API_BASE = import.meta.env.VITE_API_BASE || '/api'

const apiClient = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json'
  }
})

export default apiClient
