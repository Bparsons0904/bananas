# Bananas Frontend Clients

This directory contains a **unified testing interface** with multiple frontend framework implementations.

## 🌟 Main Interface

**Access the main testing UI at:** `http://localhost:5172`

The main interface provides:
- **Framework Switcher**: Toggle between React, Solid, and Angular implementations in real-time
- **Unified Experience**: One URL to access all frontend frameworks
- **Keyboard Shortcuts**: Alt+1 (React), Alt+2 (Solid), Alt+3 (Angular)
- **Persistent Preference**: Remembers your last selected framework

## Available Frameworks

- **Main Launcher** - Port 5172 (⭐ Start here!)
- **React** - Port 5173
- **Solid** - Port 5174
- **Angular** - Port 5175

## Running the Clients

### Recommended: Use Tilt (starts everything)

```bash
tilt up
# Then visit http://localhost:5172
```

### Individual Frameworks

```bash
# Main launcher
cd frontend && npm run dev

# React
cd react && npm run dev

# Solid
cd solid && npm run dev

# Angular
cd angular && npm start
```

## Features

Each client provides:
- Framework selector (Standard, Gin, Fiber, Echo, Chi, Gorilla Mux)
- ORM selector (database/sql, GORM, SQLx, PGX)
- Endpoint selector (Health, Simple Test, Database Test, JSON Test, Framework Info)
- Real-time request duration measurement
- JSON response display
- Error handling

## Architecture

All clients share the same functionality but use framework-specific patterns:

- **React**: Uses hooks (useState) and async/await
- **Solid**: Uses signals (createSignal) for reactive state
- **Angular**: Uses signals and standalone components

## Backend Integration

The frontends connect to the Go backend servers running on ports 8081-8086.
Ensure the backend services are running before testing.
