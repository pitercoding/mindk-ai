/// <reference types="vitest/config" />
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': import.meta.dirname + '/src',
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./src/test/setup.ts'],
    css: false,
    // Run each test file in an isolated VM context. The default pools fail on
    // Windows when the terminal's working directory has a lower-case drive
    // letter ("d:\..."), because Vitest ends up loaded twice and every file
    // fails with "Cannot read properties of undefined (reading 'config')".
    pool: 'vmThreads',
  },
})
