import { Component, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

interface Framework {
  name: string;
  value: string;
  port: number;
}

interface Orm {
  name: string;
  value: string;
}

interface Endpoint {
  name: string;
  path: string;
}

interface TestResult {
  framework: string;
  orm: string;
  response: any;
  duration: number;
  error?: string;
}

@Component({
  selector: 'app-root',
  imports: [CommonModule, FormsModule],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  frameworks: Framework[] = [
    { name: 'Standard Library', value: 'standard', port: 8081 },
    { name: 'Gin', value: 'gin', port: 8082 },
    { name: 'Fiber', value: 'fiber', port: 8083 },
    { name: 'Echo', value: 'echo', port: 8084 },
    { name: 'Chi', value: 'chi', port: 8085 },
    { name: 'Gorilla Mux', value: 'gorilla', port: 8086 }
  ];

  orms: Orm[] = [
    { name: 'database/sql', value: 'sql' },
    { name: 'GORM', value: 'gorm' },
    { name: 'SQLx', value: 'sqlx' },
    { name: 'PGX', value: 'pgx' }
  ];

  endpoints: Endpoint[] = [
    { name: 'Health Check', path: '/health' },
    { name: 'Simple Test', path: '/api/test/simple' },
    { name: 'Database Test', path: '/api/test/database?limit=10' },
    { name: 'JSON Test', path: '/api/test/json' },
    { name: 'Framework Info', path: '/api/info' }
  ];

  selectedFramework = signal<Framework>(this.frameworks[0]);
  selectedOrm = signal<Orm>(this.orms[0]);
  selectedEndpoint = signal<Endpoint>(this.endpoints[0]);
  testResult = signal<TestResult | null>(null);
  loading = signal(false);

  async runTest() {
    this.loading.set(true);
    this.testResult.set(null);

    const startTime = performance.now();
    const framework = this.selectedFramework();
    const orm = this.selectedOrm();
    const endpoint = this.selectedEndpoint();

    const url = `http://localhost:${framework.port}${endpoint.path}${
      endpoint.path.includes('database') ? `&orm=${orm.value}` : ''
    }`;

    try {
      const response = await fetch(url);
      const data = await response.json();
      const duration = performance.now() - startTime;

      this.testResult.set({
        framework: framework.name,
        orm: orm.name,
        response: data,
        duration
      });
    } catch (error) {
      const duration = performance.now() - startTime;
      this.testResult.set({
        framework: framework.name,
        orm: orm.name,
        response: null,
        duration,
        error: error instanceof Error ? error.message : 'Unknown error'
      });
    } finally {
      this.loading.set(false);
    }
  }

  onFrameworkChange(value: string) {
    const framework = this.frameworks.find(f => f.value === value);
    if (framework) {
      this.selectedFramework.set(framework);
    }
  }

  onOrmChange(value: string) {
    const orm = this.orms.find(o => o.value === value);
    if (orm) {
      this.selectedOrm.set(orm);
    }
  }

  onEndpointChange(path: string) {
    const endpoint = this.endpoints.find(e => e.path === path);
    if (endpoint) {
      this.selectedEndpoint.set(endpoint);
    }
  }
}
