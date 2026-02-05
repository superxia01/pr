import { test } from '@playwright/test';

test.use({
  launchOptions: {
    channel: 'chrome',
  },
});

test('手动测试登录流程', async ({ page }) => {
  console.log('\n🚀 测试开始...');
  console.log('📍 正在打开 https://pr.crazyaigc.com\n');

  // 访问网站
  await page.goto('https://pr.crazyaigc.com');

  // 监听所有 console 消息
  page.on('console', msg => {
    const text = msg.text();
    if (text.includes('AuthContext') ||
        text.includes('ProtectedRoute') ||
        text.includes('API') ||
        text.includes('Token') ||
        text.includes('401') ||
        text.includes('error') ||
        text.includes('Error')) {
      console.log(`[Console] ${text}`);
    }
  });

  // 监听所有请求
  page.on('request', request => {
    const url = request.url();
    if (url.includes('/api/')) {
      console.log(`[Request] ${request.method()} ${url}`);
    }
  });

  // 监听所有响应
  page.on('response', response => {
    const url = response.url();
    if (url.includes('/api/')) {
      console.log(`[Response] ${response.status()} ${url}`);
      if (response.status() === 401) {
        console.log(`  ⚠️ 401 Unauthorized - Token 可能无效或过期`);
      }
    }
  });

  console.log('✅ 浏览器已打开');
  console.log('📱 请在浏览器中扫码登录...\n');

  // 等待用户登录（最多 5 分钟）
  console.log('⏳ 等待登录（最多 5 分钟）...');

  try {
    // 等待 URL 包含 dashboard
    await page.waitForURL(/dashboard/, { timeout: 300000 });

    console.log('\n✅ 检测到登录成功！跳转到了 Dashboard\n');

    // 等待页面稳定
    await page.waitForTimeout(3000);

    // 检查 localStorage
    const storage = await page.evaluate(() => ({
      accessToken: localStorage.getItem('accessToken'),
      refreshToken: localStorage.getItem('refreshToken'),
      user: localStorage.getItem('user'),
    }));

    console.log('='.repeat(60));
    console.log('📦 localStorage 内容：');
    console.log('='.repeat(60));
    console.log('accessToken:', storage.accessToken ? '✅ 存在' : '❌ 不存在');
    if (storage.accessToken) {
      console.log('  长度:', storage.accessToken.length);
      console.log('  前50字符:', storage.accessToken.substring(0, 50) + '...');
    }
    console.log('refreshToken:', storage.refreshToken ? '✅ 存在' : '❌ 不存在');
    console.log('user:', storage.user ? '✅ 存在' : '❌ 不存在');
    if (storage.user) {
      try {
        const userData = JSON.parse(storage.user);
        console.log('  用户数据:', userData);
      } catch (e) {
        console.log('  (解析失败)');
      }
    }
    console.log('='.repeat(60) + '\n');

    // 截图
    await page.screenshot({ path: 'screenshots/dashboard-after-login.png' });
    console.log('📸 已保存截图: screenshots/dashboard-after-login.png\n');

    console.log('🔍 现在请手动点击导航菜单，观察是否有问题...');
    console.log('⏸️ 浏览器将保持打开，你可以手动测试\n');

    // 暂停，保持浏览器打开
    await page.waitForTimeout(120000); // 等待 2 分钟
    await page.pause();

  } catch (error) {
    console.log('\n❌ 等待登录超时:', error.message);
    console.log('💡 请在浏览器中手动检查是否登录成功\n');
    await page.pause();
  }
});
