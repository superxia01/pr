import { test, expect } from '@playwright/test';

// 使用系统 Chrome
test.use({
  launchOptions: {
    channel: 'chrome', // 使用系统安装的 Chrome
  },
});

test.describe('PR Business 登录测试', () => {
  test('完整登录流程测试', async ({ page, context }) => {
    // 启用调试模式
    await page.setViewportSize({ width: 1280, height: 720 });

    console.log('🌐 正在打开 https://pr.crazyaigc.com ...');

    // 访问首页
    await page.goto('https://pr.crazyaigc.com');

    console.log('✅ 页面已加载');

    // 等待页面加载
    await page.waitForLoadState('networkidle');

    // 截图 - 登录前
    await page.screenshot({ path: 'screenshots/01-login-page.png' });
    console.log('📸 已截图: 01-login-page.png');

    // 检查是否在登录页
    const currentUrl = page.url();
    console.log('📍 当前 URL:', currentUrl);

    if (currentUrl.includes('/login')) {
      console.log('✅ 当前在登录页');

      // 等待微信登录按钮出现
      await page.waitForSelector('text=微信登录', { timeout: 5000 });
      console.log('✅ 找到微信登录按钮');

      // 提示用户扫码
      console.log('\n' + '='.repeat(60));
      console.log('📱 请在浏览器中使用微信扫码登录');
      console.log('='.repeat(60) + '\n');

      // 点击微信登录按钮
      const wechatLoginButton = page.locator('text=微信登录').first();
      await wechatLoginButton.click();

      console.log('✅ 已点击微信登录按钮');

      // 等待跳转（可能跳转到 os.crazyaigc.com 或微信授权页）
      await page.waitForTimeout(3000);

      // 截图 - 点击登录后
      await page.screenshot({ path: 'screenshots/02-after-click-login.png' });
      console.log('📸 已截图: 02-after-click-login.png');

      // 等待登录完成（最长等待 120 秒让用户扫码）
      console.log('⏳ 等待用户扫码登录（最多120秒）...');

      try {
        // 等待 URL 变化（不再是 /login）
        await page.waitForURL(/\/(dashboard|login)/, { timeout: 120000 });

        const newUrl = page.url();
        console.log('📍 登录后 URL:', newUrl);

        if (newUrl.includes('/dashboard')) {
          console.log('✅ 登录成功！已跳转到 Dashboard');

          // 截图 - Dashboard 页面
          await page.waitForTimeout(2000);
          await page.screenshot({ path: 'screenshots/03-dashboard.png' });
          console.log('📸 已截图: 03-dashboard.png');

          // 检查 localStorage 中的 token
          const localStorage = await page.evaluate(() => {
            return {
              accessToken: localStorage.getItem('accessToken'),
              refreshToken: localStorage.getItem('refreshToken'),
              user: localStorage.getItem('user'),
            };
          });

          console.log('\n' + '='.repeat(60));
          console.log('📦 localStorage 内容:');
          console.log('='.repeat(60));
          console.log('Token 存在:', !!localStorage.accessToken);
          console.log('Token 长度:', localStorage.accessToken?.length || 0);
          console.log('RefreshToken 存在:', !!localStorage.refreshToken);
          console.log('User 存在:', !!localStorage.user);
          console.log('User 内容:', localStorage.user ? JSON.parse(localStorage.user) : null);
          console.log('='.repeat(60) + '\n');

          // 点击导航菜单测试
          console.log('🔍 测试点击导航菜单...');

          // 等待导航菜单加载
          await page.waitForTimeout(2000);

          // 尝试点击任意导航链接
          const navLinks = page.locator('a[href^="/"]').all();
          console.log(`找到 ${navLinks.length} 个导航链接`);

          if (navLinks.length > 0) {
            // 点击第一个导航链接（排除 /login 和 /dashboard）
            let clicked = false;
            for (const link of await page.locator('a[href^="/"]').all()) {
              const href = await link.getAttribute('href');
              if (href && href !== '/login' && href !== '/dashboard') {
                console.log(`🖱️ 点击导航链接: ${href}`);

                // 点击前的 URL
                const beforeClickUrl = page.url();

                // 点击链接
                await link.click();

                // 等待导航
                await page.waitForTimeout(3000);

                // 点击后的 URL
                const afterClickUrl = page.url();

                console.log(`点击前 URL: ${beforeClickUrl}`);
                console.log(`点击后 URL: ${afterClickUrl}`);

                // 截图 - 点击导航后
                await page.screenshot({ path: 'screenshots/04-after-nav-click.png' });
                console.log('📸 已截图: 04-after-nav-click.png');

                // 检查是否跳回了登录页
                if (afterClickUrl.includes('/login')) {
                  console.log('❌ 问题确认：点击导航后跳回了登录页！');
                } else {
                  console.log('✅ 正常：导航后保持在页面内');
                }

                clicked = true;
                break;
              }
            }

            if (!clicked) {
              console.log('⚠️ 没有找到可点击的导航链接');
            }
          } else {
            console.log('⚠️ 没有找到导航链接');
          }

          // 检查 Console 日志
          console.log('\n' + '='.repeat(60));
          console.log('📋 浏览器 Console 日志:');
          console.log('='.repeat(60));

          // 监听 console 消息
          page.on('console', msg => {
            const text = msg.text();
            if (text.includes('AuthContext') ||
                text.includes('ProtectedRoute') ||
                text.includes('API') ||
                text.includes('Token') ||
                text.includes('401') ||
                text.includes('error')) {
              console.log(`[浏览器] ${text}`);
            }
          });

          console.log('='.repeat(60) + '\n');

        } else if (newUrl.includes('/login')) {
          console.log('⚠️ 仍在登录页，可能需要等待微信回调');
        }

      } catch (error) {
        console.log('⏰ 等待超时或出错:', error.message);
      }

    } else {
      console.log('⚠️ 当前不在登录页，可能已登录或其他情况');

      // 截图 - 当前页面
      await page.screenshot({ path: 'screenshots/00-current-page.png' });
      console.log('📸 已截图: 00-current-page.png');
    }

    // 暂停，让用户查看
    console.log('\n' + '='.repeat(60));
    console.log('⏸️ 测试完成，浏览器将保持打开 60 秒...');
    console.log('💡 你可以手动测试，查看 Console 和 Network');
    console.log('='.repeat(60) + '\n');

    // 暂停 60 秒让用户手动测试
    await page.waitForTimeout(60000);

    // 暂停执行，保持浏览器打开
    await page.pause();
  });
});
