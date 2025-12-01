<script setup lang="ts">
import { ref } from 'vue'

interface TestResult {
  framework: string
  orm: string
  response: any
  duration: number
  error?: string
}

const frameworks = [
  { name: 'Standard Library', value: 'standard', port: 8081 },
  { name: 'Gin', value: 'gin', port: 8082 },
  { name: 'Fiber', value: 'fiber', port: 8083 },
  { name: 'Echo', value: 'echo', port: 8084 },
  { name: 'Chi', value: 'chi', port: 8085 },
  { name: 'Gorilla Mux', value: 'gorilla', port: 8086 }
]

const orms = [
  { name: 'database/sql', value: 'sql' },
  { name: 'GORM', value: 'gorm' },
  { name: 'SQLx', value: 'sqlx' },
  { name: 'PGX', value: 'pgx' }
]

const endpoints = [
  { name: 'Health Check', path: '/health' },
  { name: 'Simple Test', path: '/api/test/simple' },
  { name: 'Database Test', path: '/api/test/database?limit=10' },
  { name: 'JSON Test', path: '/api/test/json' },
  { name: 'Framework Info', path: '/api/info' }
]

const selectedFramework = ref(frameworks[0])
const selectedOrm = ref(orms[0])
const selectedEndpoint = ref(endpoints[0])
const testResult = ref<TestResult | null>(null)
const loading = ref(false)

const runTest = async () => {
  loading.value = true
  testResult.value = null

  const startTime = performance.now()
  const url = `http://localhost:${selectedFramework.value.port}${selectedEndpoint.value.path}${
    selectedEndpoint.value.path.includes('database') ? `&orm=${selectedOrm.value.value}` : ''
  }`

  try {
    const response = await fetch(url)
    const data = await response.json()
    const duration = performance.now() - startTime

    testResult.value = {
      framework: selectedFramework.value.name,
      orm: selectedOrm.value.name,
      response: data,
      duration
    }
  } catch (error) {
    const duration = performance.now() - startTime
    testResult.value = {
      framework: selectedFramework.value.name,
      orm: selectedOrm.value.name,
      response: null,
      duration,
      error: error instanceof Error ? error.message : 'Unknown error'
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="app">
    <header>
      <h1>🍌 Bananas Framework Tester</h1>
      <p>Vue Client</p>
    </header>

    <div class="controls">
      <div class="control-group">
        <label>Framework</label>
        <select v-model="selectedFramework">
          <option v-for="fw in frameworks" :key="fw.value" :value="fw">
            {{ fw.name }} (:{{ fw.port }})
          </option>
        </select>
      </div>

      <div class="control-group">
        <label>ORM</label>
        <select v-model="selectedOrm">
          <option v-for="orm in orms" :key="orm.value" :value="orm">
            {{ orm.name }}
          </option>
        </select>
      </div>

      <div class="control-group">
        <label>Endpoint</label>
        <select v-model="selectedEndpoint">
          <option v-for="ep in endpoints" :key="ep.path" :value="ep">
            {{ ep.name }}
          </option>
        </select>
      </div>

      <button @click="runTest" :disabled="loading" class="test-button">
        {{ loading ? 'Testing...' : 'Run Test' }}
      </button>
    </div>

    <div v-if="testResult" class="results">
      <h2>Results</h2>
      <div class="result-info">
        <div><strong>Framework:</strong> {{ testResult.framework }}</div>
        <div><strong>ORM:</strong> {{ testResult.orm }}</div>
        <div><strong>Duration:</strong> {{ testResult.duration.toFixed(2) }}ms</div>
      </div>

      <div v-if="testResult.error" class="error">
        <strong>Error:</strong> {{ testResult.error }}
      </div>
      <pre v-else class="response">{{ JSON.stringify(testResult.response, null, 2) }}</pre>
    </div>
  </div>
</template>

<style scoped>
.app {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

header {
  text-align: center;
  margin-bottom: 2rem;
}

header h1 {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
}

header p {
  font-size: 1.2rem;
  color: #42b883;
  font-weight: 600;
}

.controls {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  flex-wrap: wrap;
  align-items: end;
}

.control-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  flex: 1;
  min-width: 200px;
}

.control-group label {
  font-weight: 600;
}

.control-group select {
  padding: 0.5rem;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 4px;
  font-size: 1rem;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  cursor: pointer;
}

.control-group select option {
  background: #242424;
  color: white;
}

.control-group select:focus {
  outline: none;
  border-color: #42b883;
}

.test-button {
  padding: 0.5rem 2rem;
  font-size: 1rem;
  font-weight: 600;
  background-color: #42b883;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.test-button:hover:not(:disabled) {
  background-color: #35a372;
}

.test-button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.results {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 1.5rem;
  background-color: #f9f9f9;
}

.results h2 {
  margin-top: 0;
  margin-bottom: 1rem;
}

.result-info {
  display: flex;
  gap: 2rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.result-info div {
  font-size: 1rem;
}

.error {
  padding: 1rem;
  background-color: #fee;
  border: 1px solid #fcc;
  border-radius: 4px;
  color: #c33;
}

.response {
  background-color: #1e1e1e;
  color: #d4d4d4;
  padding: 1rem;
  border-radius: 4px;
  overflow-x: auto;
  font-size: 0.9rem;
  line-height: 1.5;
}
</style>
