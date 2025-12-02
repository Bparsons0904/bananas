<script lang="ts">
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

  let selectedFramework = frameworks[0]
  let selectedOrm = orms[0]
  let selectedEndpoint = endpoints[0]
  let testResult: TestResult | null = null
  let loading = false

  async function runTest() {
    loading = true
    testResult = null

    const startTime = performance.now()
    const url = `http://localhost:${selectedFramework.port}${selectedEndpoint.path}${
      selectedEndpoint.path.includes('database') ? `&orm=${selectedOrm.value}` : ''
    }`

    try {
      const response = await fetch(url)
      const data = await response.json()
      const duration = performance.now() - startTime

      testResult = {
        framework: selectedFramework.name,
        orm: selectedOrm.name,
        response: data,
        duration
      }
    } catch (error) {
      const duration = performance.now() - startTime
      testResult = {
        framework: selectedFramework.name,
        orm: selectedOrm.name,
        response: null,
        duration,
        error: error instanceof Error ? error.message : 'Unknown error'
      }
    } finally {
      loading = false
    }
  }
</script>

<div class="app">
  <header>
    <h1>🍌 Bananas Framework Tester</h1>
    <p>Svelte Client</p>
  </header>

  <div class="controls">
    <div class="control-group">
      <label for="framework">Framework</label>
      <select id="framework" bind:value={selectedFramework}>
        {#each frameworks as fw}
          <option value={fw}>
            {fw.name} (:{fw.port})
          </option>
        {/each}
      </select>
    </div>

    <div class="control-group">
      <label for="orm">ORM</label>
      <select id="orm" bind:value={selectedOrm}>
        {#each orms as orm}
          <option value={orm}>
            {orm.name}
          </option>
        {/each}
      </select>
    </div>

    <div class="control-group">
      <label for="endpoint">Endpoint</label>
      <select id="endpoint" bind:value={selectedEndpoint}>
        {#each endpoints as ep}
          <option value={ep}>
            {ep.name}
          </option>
        {/each}
      </select>
    </div>

    <button on:click={runTest} disabled={loading} class="test-button">
      {loading ? 'Testing...' : 'Run Test'}
    </button>
  </div>

  {#if testResult}
    <div class="results">
      <h2>Results</h2>
      <div class="result-info">
        <div><strong>Framework:</strong> {testResult.framework}</div>
        <div><strong>ORM:</strong> {testResult.orm}</div>
        <div><strong>Duration:</strong> {testResult.duration.toFixed(2)}ms</div>
      </div>

      {#if testResult.error}
        <div class="error">
          <strong>Error:</strong> {testResult.error}
        </div>
      {:else}
        <pre class="response">{JSON.stringify(testResult.response, null, 2)}</pre>
      {/if}
    </div>
  {/if}
</div>

<style>
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
    color: #ff3e00;
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
    width: 200px;
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
    border-color: #ff3e00;
  }

  .test-button {
    padding: 0.5rem 2rem;
    font-size: 1rem;
    font-weight: 600;
    background-color: #ff3e00;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .test-button:hover:not(:disabled) {
    background-color: #cc3200;
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
    color: #333;
  }

  .results h2 {
    margin-top: 0;
    margin-bottom: 1rem;
    color: #333;
  }

  .result-info {
    display: flex;
    gap: 2rem;
    margin-bottom: 1rem;
    flex-wrap: wrap;
  }

  .result-info div {
    font-size: 1rem;
    color: #333;
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
