import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: '.',
  fullyParallel: false,
  use: {
    headless: false, // 显示浏览器窗口
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },
});
