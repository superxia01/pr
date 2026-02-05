import { test, expect } from '@playwright/test';

test.use({
  launchOptions: {
    channel: 'chrome',
  },
});

test.describe('手动测试', () => {
  test('登录测试', async ({ page }) => {
  // 监听所有 console 消息
  page.on('console', msg => {
    const text = msg.text();
    console.log(`[浏览器] ${text}`);
  });

  // 访问网站
  await page.goto('https://pr.crazyaigc.com');

  console.log('\n' + '='.repeat(60));
  console.log('📱 浏览器已打开');
  console.log('👆 请在浏览器中扫码登录');
  console.log('⏸️ 登录后请手动点击导航菜单测试');
  console.log('⏸️ 浏览器将保持打开，按 Ctrl+C 退出');
  console.log('='.repeat(60) + '\n');

  // 暂停，保持浏览器打开
  await page.pause();
  });
});
