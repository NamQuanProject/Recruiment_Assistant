import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig(({ mode }) => {
  // Load environment variables based on the current mode (development or production)
  const env = loadEnv(mode, process.cwd(), '');

  return {
    plugins: [react()],
    server: {
      host: '0.0.0.0',
      port: 5173,
      allowedHosts: ['frontend-u3bd.onrender.com'], // Optional, remove if not needed
    },
    define: {
      // Injecting VITE_API_URL from environment variables
      'process.env.VITE_API_URL': JSON.stringify(env.VITE_API_URL),
    },
  };
});
