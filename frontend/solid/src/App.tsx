import { createSignal } from 'solid-js'
import './App.css'

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

function App() {
  const [selectedFramework, setSelectedFramework] = createSignal(frameworks[0])
  const [selectedOrm, setSelectedOrm] = createSignal(orms[0])
  const [selectedEndpoint, setSelectedEndpoint] = createSignal(endpoints[0])
  const [testResult, setTestResult] = createSignal<TestResult | null>(null)
  const [loading, setLoading] = createSignal(false)

  const runTest = async () => {
    setLoading(true)
    setTestResult(null)

    const startTime = performance.now()
    const url = `http://localhost:${selectedFramework().port}${selectedEndpoint().path}${
      selectedEndpoint().path.includes('database') ? `&orm=${selectedOrm().value}` : ''
    }`

    try {
      const response = await fetch(url)
      const data = await response.json()
      const duration = performance.now() - startTime

      setTestResult({
        framework: selectedFramework().name,
        orm: selectedOrm().name,
        response: data,
        duration
      })
    } catch (error) {
      const duration = performance.now() - startTime
      setTestResult({
        framework: selectedFramework().name,
        orm: selectedOrm().name,
        response: null,
        duration,
        error: error instanceof Error ? error.message : 'Unknown error'
      })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div class="app">
      <header>
        <h1>🍌 Bananas Framework Tester</h1>
        <p>Solid Client</p>
      </header>

      <div class="controls">
        <div class="control-group">
          <label>Framework</label>
          <select
            value={selectedFramework().value}
            onChange={(e) => setSelectedFramework(frameworks.find(f => f.value === e.currentTarget.value)!)}
          >
            {frameworks.map(fw => (
              <option value={fw.value}>
                {fw.name} (:{fw.port})
              </option>
            ))}
          </select>
        </div>

        <div class="control-group">
          <label>ORM</label>
          <select
            value={selectedOrm().value}
            onChange={(e) => setSelectedOrm(orms.find(o => o.value === e.currentTarget.value)!)}
          >
            {orms.map(orm => (
              <option value={orm.value}>
                {orm.name}
              </option>
            ))}
          </select>
        </div>

        <div class="control-group">
          <label>Endpoint</label>
          <select
            value={selectedEndpoint().path}
            onChange={(e) => setSelectedEndpoint(endpoints.find(ep => ep.path === e.currentTarget.value)!)}
          >
            {endpoints.map(ep => (
              <option value={ep.path}>
                {ep.name}
              </option>
            ))}
          </select>
        </div>

        <button onClick={runTest} disabled={loading()} class="test-button">
          {loading() ? 'Testing...' : 'Run Test'}
        </button>
      </div>

      {testResult() && (
        <div class="results">
          <h2>Results</h2>
          <div class="result-info">
            <div><strong>Framework:</strong> {testResult()!.framework}</div>
            <div><strong>ORM:</strong> {testResult()!.orm}</div>
            <div><strong>Duration:</strong> {testResult()!.duration.toFixed(2)}ms</div>
          </div>

          {testResult()!.error ? (
            <div class="error">
              <strong>Error:</strong> {testResult()!.error}
            </div>
          ) : (
            <pre class="response">
              {JSON.stringify(testResult()!.response, null, 2)}
            </pre>
          )}
        </div>
      )}
    </div>
  )
}

export default App
