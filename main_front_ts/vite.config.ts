import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // Ensure requests starting with /api are forwarded to your backend
      '/api': {
        target: 'http://localhost:8090', // Replace with your backend server address if different
        changeOrigin: true, // Needed for virtual hosted sites
        rewrite: (path) => path.replace(/^\/api/, ''), // Remove the /api prefix when forwarding
      },
    },
  },
})
