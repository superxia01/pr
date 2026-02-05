import { test } from '@playwright/test';

test('调试登录流程', async ({ page }) => {
  // 监听 console
  page.on('console', msg => {
    const text = msg.text();
    if (text.includes('AuthContext') || text.includes('login') || 
        text.includes('Token') || text.includes('localStorage') ||
        text.includes('API') || text.includes('🔍') || text.includes('✅')) {
      console.log(`[浏览器] ${text}`);
    }
  });

  // 访问登录页
  await page.goto('https://pr.crazyaigc.com/login');
  
  console.log('\n📱 浏览器已打开，请扫码登录...');
  console.log('⏸️ 登录后浏览器将保持打开\n');

  // 等待最多5分钟检测登录成功
  try {
    await page.waitForURL(/dashboard/, { timeout: 300000 });
    console.log('\n✅ 检测到登录成功！');
    
    // 等待3秒让页面稳定
    await page.waitForTimeout(3000);
    
    // 检查 localStorage
    const storage = await page.evaluate(() => ({
      token: localStorage.getItem('accessToken'),
      user: localStorage.getItem('user'),
    }));
    
    console.log('\n📦 localStorage 状态:');
    console.log('  Token:', storage.token ? '✅ 存在' : '❌ 不存在');
    console.log('  User:', storage.user ? '✅ 存在' : '❌ 不存在');
    
    // 截图
    await page.screenshot({ path: 'debug-result.png' });
    console.log('\n📸 已保存截图: debug-result.png');
    
    // 保持浏览器打开
    console.log('\n⏸️ 浏览器将保持打开 60 秒供你手动测试...');
    await page.waitForTimeout(60000);
    
  } catch (e) {
    console.log('\n⏰ 等待超时或出错:', e.message);
  }
});
