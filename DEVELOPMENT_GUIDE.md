# Development Guide – ONVIF Camera Management System

## Project Structure

- **main_back/**: Go backend (API, camera logic, config)
- **main_front/**: React frontend (UI, config, API calls)

## Quick Start

```cmd
# From project root
start_system.bat
```

### Manual Start

**Backend:**
```cmd
cd main_back
go run cmd/backend/main.go
```

**Frontend:**
```cmd
cd main_front
npm install
npm run dev
```

- Frontend: http://localhost:5173
- Backend API: http://localhost:8090

## Development Workflow

### Frontend (main_front)
- Components: `src/components/`
- Pages: `src/pages/`
- API: `src/services/api.js`
- Styles: `src/App.css`
- Main entry: `src/App.jsx`

**To add features:**
- New UI: add to `components/` or `pages/`
- New API call: update `services/api.js`
- Update routing: edit `App.jsx`

### Backend (main_back)
- Entry: `cmd/backend/main.go`
- API: `internal/api/`
- Camera logic: `internal/camera/`
- Config: `config/cameras.json`

**To add features:**
- New endpoint: add to `internal/api/`
- Camera logic: extend `internal/camera/`
- Update config: edit `config/cameras.json`

## Testing

- **Frontend:** Use browser dev tools, React component tests, or `/debug` route.
- **Backend:** Use curl/Postman for API, Go tests for logic.

## Debugging

- **Frontend:** Check browser console, use `console.log`.
- **Backend:** Check terminal output, use `go run ... 2>&1 | tee backend.log`.

## Best Practices

- Use environment variables for secrets (especially in production).
- Keep dependencies updated (`npm update`, `go get -u`).
- Use `.gitignore` to avoid committing build artifacts and secrets.
- Document new features and endpoints.

## Maintenance

- Regularly update dependencies.
- Monitor for security issues.
- Backup configuration files.

---

_Last updated: June 2025_
