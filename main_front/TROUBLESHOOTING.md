# Backend Connection Troubleshooting

If you're experiencing issues connecting to the backend server from the frontend, follow these steps:

## 1. Verify Backend Server is Running

Make sure the backend server is running on port 8090:

```bash
cd ../main_back
go run cmd/backend/main.go
```

You should see console output indicating the server is listening.

## 2. Test API Directly

Test the backend API directly with curl or a browser:

```bash
curl http://localhost:8090/cameras
```

You should receive JSON data with camera information.

## 3. CORS Issues

If you encounter CORS errors in the browser console, you have several options:

### Option 1: Use the Vite Development Proxy (Recommended)

The frontend is already configured to proxy API requests via Vite's dev server.
Just make sure your API calls use the `/api` prefix.

### Option 2: Run the CORS Proxy Server

Install the required dependencies and start the CORS proxy:

```bash
npm install
npm run start:proxy
```

Then update the `.env` file to point to the proxy:

```
VITE_API_BASE_URL=http://localhost:8888
```

### Option 3: Update Backend CORS Settings

The backend should already have CORS configured correctly in `internal/api/routers.go`.
Verify that it includes:

```go
corsOptions := handlers.CORS(
    handlers.AllowedOrigins([]string{"*"}),
    handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
    handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
)
```

## 4. Network Issues

Ensure there are no firewall or network restrictions blocking port 8090.

## 5. Backend Error Logs

Check the backend console output for any error messages that might indicate issues with the Go server.

## 6. Restart Both Frontend and Backend

Sometimes a simple restart of both servers can resolve connection issues:

```bash
# Terminal 1: Backend
cd ../main_back
go run cmd/backend/main.go

# Terminal 2: Frontend
npm run dev
```
