// vite.config.ts
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
      allowedHosts: ['frontend-eqtg.onrender.com'], // hoặc bỏ nếu không cần giới hạn host
    },
    define: {
      'process.env': env, // nếu bạn muốn dùng process.env.VITE_API_URL
    },
  };
});
