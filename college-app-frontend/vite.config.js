import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // Quando o React tentar acessar '/api', será redirecionado para o seu backend Go
      '/api': {
        target: 'http://localhost:8080', // Onde sua API Go está rodando
        changeOrigin: true, // Necessário para hosts virtuais baseados em nome
        rewrite: (path) => path.replace(/^\/api/, ''), // Remove '/api' do path antes de enviar para o Go
      },
    },
  },
})