/// <reference types="vitest" />

import path from 'path';
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    setupFiles: './src/setupTests.ts',
    environment: 'jsdom',
    reporters: 'html',
    includeTaskLocation: true,
  },
  resolve: {
    alias: [{ find: 'src', replacement: path.resolve(__dirname, './src') }],
  },
});
