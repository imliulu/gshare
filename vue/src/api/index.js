import axios from 'axios'

const apiClient = axios.create({
  baseURL: process.env.VUE_APP_API_URL+"api/",
  timeout: 1000000,
  headers: {
    'Content-Type': 'application/json'
  }
})

export default apiClient
