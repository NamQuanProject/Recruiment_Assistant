// vite.config.ts
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',
    port: 5173, // hoặc cổng bạn muốn
    allowedHosts: ['frontend-c5z1.onrender.com'], // 👈 thêm dòng này
  },
});
