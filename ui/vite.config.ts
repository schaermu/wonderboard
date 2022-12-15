import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    emptyOutDir: false
  },
  server: {
    proxy: {
      '/api': 'http://localhost:3000/api'
    }
  },
  plugins: [svelte()],
})
