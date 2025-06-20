// CORS Proxy for development
// Run this script with Node.js if you encounter CORS issues
// Usage: node cors-proxy.js

const express = require('express');
const cors = require('cors');
const { createProxyMiddleware } = require('http-proxy-middleware');

const app = express();
const port = 8888;

// Enable CORS for all routes
app.use(cors());

// Proxy all requests to the backend
app.use('/', createProxyMiddleware({
  target: 'http://localhost:8090', // Target backend URL
  changeOrigin: true,
  onProxyRes: (proxyRes) => {
    proxyRes.headers['Access-Control-Allow-Origin'] = '*';
    proxyRes.headers['Access-Control-Allow-Methods'] = 'GET, POST, PUT, DELETE, OPTIONS';
    proxyRes.headers['Access-Control-Allow-Headers'] = 'Content-Type, Authorization';
  }
}));

app.listen(port, () => {
  console.log(`CORS Proxy running at http://localhost:${port}`);
  console.log(`Proxying requests to http://localhost:8090`);
  console.log(`Update your frontend to use http://localhost:${port} instead of directly accessing the backend`);
});
