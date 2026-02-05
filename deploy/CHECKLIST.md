# 部署 Checklist (Deployment Checklist)

**目的**: 确保每次部署都符合 KeenChase V4.0 标准，避免遗漏和错误

**使用方法**: 每次部署时逐项检查，完成一项打勾一项

---

## 🚀 部署前检查 (Pre-Deployment)

### 本地环境
- [ ] **前端构建测试**
  ```bash
  cd frontend && npm run build
  ```
  预期结果：构建成功，无错误，dist/ 目录大小 < 2MB

- [ ] **后端编译测试**
  ```bash
  cd backend && go build
  ```
  预期结果：编译成功，生成可执行文件

- [ ] **本地功能测试**
  - [ ] 登录功能正常
  - [ ] 主要页面访问正常
  - [ ] API 调用正常

### 代码管理
- [ ] **代码已提交**
  ```bash
  git status
  ```
  预期结果：无未提交的重要修改

- [ ] **Commit Message 清晰**
  - [ ] 描述了修改内容
  - [ ] 遵循 commit 规范

### 确认部署范围
- [ ] **部署目标明确**
  - [ ] 前端 only？
  - [ ] 后端 only？
  - [ ] 全部？

- [ ] **部署环境确认**
  - [ ] 生产环境 (pr.crazyaigc.com)
  - [ ] 测试环境

---

## 📦 部署中检查 (During Deployment)

### 前端部署
- [ ] **构建步骤**
  ```bash
  cd frontend
  npm run build
  ```
  - [ ] 构建时间 < 30 秒
  - [ ] 无 TypeScript 错误
  - [ ] dist/ 目录生成

- [ ] **上传步骤**
  ```bash
  rsync -avz --delete --exclude '*.map' dist/ shanghai-tencent:/var/www/pr-frontend/
  ```
  - [ ] 上传到 **正确目录**: `/var/www/pr-frontend/`
  - [ ] ❌ 不是 `/var/www/pr-business-frontend/`
  - [ ] 文件数量正确（通常 6-10 个文件）
  - [ ] 总大小 < 1MB

- [ ] **Nginx 重载**
  ```bash
  ssh shanghai-tencent "sudo systemctl reload nginx"
  ```
  - [ ] 无错误输出

### 后端部署
- [ ] **交叉编译**
  ```bash
  cd backend
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pr-business-linux
  ```
  - [ ] 编译成功
  - [ ] 文件大小约 20-30MB
  - [ ] `file pr-business-linux` 显示 ELF 64-bit

- [ ] **上传步骤**
  ```bash
  scp pr-business-linux shanghai-tencent:/var/www/pr-backend/
  ```
  - [ ] 上传到 **正确目录**: `/var/www/pr-backend/`
  - [ ] ❌ 不是 `/var/www/pr-business/`
  - [ ] 上传速度正常

- [ ] **服务重启**
  ```bash
  ssh shanghai-tencent << 'ENDSSH'
  cd /var/www/pr-backend
  sudo cp pr-business pr-business.backup.$(date +%Y%m%d_%H%M%S)
  sudo mv pr-business-linux pr-business
  sudo chmod +x pr-business
  sudo systemctl restart pr-business-backend
  ENDSSH
  ```
  - [ ] 备份创建成功
  - [ ] 服务重启成功
  - [ ] **❌ 未上传 .env 文件**（环境变量与代码分离）

---

## ✅ 部署后检查 (Post-Deployment)

### 服务状态检查
- [ ] **后端服务状态**
  ```bash
  ssh shanghai-tencent "sudo systemctl status pr-business-backend"
  ```
  - [ ] 状态：`active (running)`
  - [ ] 内存占用正常（< 50MB）
  - [ ] 无重启失败记录

- [ ] **前端文件检查**
  ```bash
  ssh shanghai-tencent "ls -la /var/www/pr-frontend/"
  ```
  - [ ] `index.html` 存在
  - [ ] `assets/` 目录存在
  - [ ] 文件时间戳是最新的

- [ ] **Nginx 配置检查**
  ```bash
  ssh shanghai-tencent "sudo nginx -t"
  ```
  - [ ] 配置测试通过

### 功能验证
- [ ] **网站可访问**
  ```bash
  curl -I https://pr.crazyaigc.com
  ```
  预期结果：`HTTP/1.1 200 OK`

- [ ] **登录功能测试**
  - [ ] 打开 https://pr.crazyaigc.com
  - [ ] 点击登录
  - [ ] 微信扫码登录正常
  - [ ] 登录后显示正确的用户信息

- [ ] **API 测试**
  - [ ] 浏览器开发者工具无 API 错误
  - [ ] Network 面板 API 请求正常
  - [ ] 无 401/403/500 错误

- [ ] **页面功能测试**
  - [ ] Dashboard 正常显示
  - [ ] 商家列表/服务商列表正常
  - [ ] 数据加载正常

### 日志检查
- [ ] **后端日志检查**
  ```bash
  ssh shanghai-tencent "sudo journalctl -u pr-business-backend -n 50 --no-pager"
  ```
  - [ ] 无严重错误
  - [ ] 无数据库连接错误
  - [ ] 启动日志正常

- [ ] **Nginx 日志检查**
  ```bash
  ssh shanghai-tencent "sudo tail -n 20 /var/log/nginx/error.log"
  ```
  - [ ] 无新的错误

### 性能检查
- [ ] **页面加载速度**
  - [ ] 首屏加载 < 2 秒
  - [ ] API 响应 < 500ms

- [ ] **静态资源**
  - [ ] 无 404 错误
  - [ ] JS/CSS 加载正常

---

## 🔥 回滚检查 (Rollback Checklist)

**如果部署出现问题，立即执行以下步骤：**

### 确认回滚需求
- [ ] **问题确认**
  - [ ] 服务无法启动？
  - [ ] 功能异常？
  - [ ] 性能严重下降？

### 后端回滚
- [ ] **停止部署**
  ```bash
  ssh shanghai-tencent "sudo systemctl stop pr-business-backend"
  ```

- [ ] **恢复备份**
  ```bash
  ssh shanghai-tencent << 'ENDSSH'
  cd /var/www/pr-backend
  sudo mv pr-business pr-business.failed.$(date +%Y%m%d_%H%M%S)
  sudo mv pr-business.backup.YYYYMMDD_HHMMSS pr-business
  sudo chmod +x pr-business
  ENDSSH
  ```

- [ ] **重启服务**
  ```bash
  ssh shanghai-tencent "sudo systemctl start pr-business-backend"
  ```

- [ ] **验证回滚**
  ```bash
  ssh shanghai-tencent "sudo systemctl status pr-business-backend"
  curl https://pr.crazyaigc.com/api/v1/service-providers
  ```

### 前端回滚
- [ ] **恢复备份**
  ```bash
  rsync -avz --delete \
    shanghai-tencent:/var/www/pr-frontend.backup.YYYYMMDD_HHMMSS/ \
    dist/

  rsync -avz --delete dist/ shanghai-tencent:/var/www/pr-frontend/
  ```

- [ ] **清除浏览器缓存**
  - [ ] Chrome: Cmd+Shift+R (Mac) / Ctrl+Shift+R (Windows)
  - [ ] 或使用无痕模式测试

- [ ] **验证回滚**
  ```bash
  curl -I https://pr.crazyaigc.com
  ```

---

## 🚨 常见错误检查

### 目录错误
- [ ] **检查服务器目录**
  ```bash
  ssh shanghai-tencent "ls -la /var/www/ | grep pr"
  ```
  正确输出应该只有：
  - `pr-backend`
  - `pr-frontend`

  ❌ 如果出现以下目录，说明部署错误：
  - `pr-business` → 删除
  - `pr-business-frontend` → 删除

  删除命令：
  ```bash
  ssh shanghai-tencent "sudo rm -rf /var/www/pr-business /var/www/pr-business-frontend"
  ```

### 环境变量错误
- [ ] **确认 .env 文件存在**
  ```bash
  ssh shanghai-tencent "ls -la /var/www/pr-backend/.env"
  ```
  - [ ] 文件存在
  - [ ] 权限正确（600）
  - [ ] 所有者正确（ubuntu:ubuntu）

- [ ] **确认 .env 未被部署覆盖**
  - [ ] 检查配置值是否正确
  - [ ] 特别是数据库密码、JWT_SECRET 等

### 服务名称错误
- [ ] **确认服务名**
  ```bash
  ssh shanghai-tencent "sudo systemctl list-units | grep pr-business"
  ```
  正确输出：`pr-business-backend.service`

  ❌ 如果是 `pr-business.service`，说明配置错误

### Nginx 配置错误
- [ ] **检查 Nginx root 路径**
  ```bash
  ssh shanghai-tencent "sudo cat /etc/nginx/sites-enabled/pr-business | grep 'root '"
  ```
  正确输出：`root /var/www/pr-frontend;`

  ❌ 如果是 `/var/www/pr-business-frontend`，说明配置错误

---

## 📋 部署记录模板

每次部署后记录：

```markdown
## 部署记录 - YYYY-MM-DD HH:mm

**部署人**: [你的名字]
**部署类型**: [前端/后端/全部]
**Commit**: [commit hash]

### 部署内容
- [ ] 前端修改：[简述]
- [ ] 后端修改：[简述]

### 部署过程
- [ ] 构建时间：X 秒
- [ ] 部署时间：X 秒
- [ ] 遇到问题：[无 / 具体问题]

### 部署结果
- [ ] ✅ 成功 / ❌ 失败 / ⚠️ 部分失败

### 验证结果
- [ ] 服务状态：正常 / 异常
- [ ] 功能测试：通过 / 失败
- [ ] 性能测试：正常 / 异常

### 备注
[其他需要记录的信息]
```

---

## 🎯 快速检查命令

```bash
# 一键检查所有关键状态
ssh shanghai-tencent << 'ENDSSH'
echo "=== 服务状态 ==="
sudo systemctl status pr-business-backend --no-pager | head -10

echo ""
echo "=== 目录结构 ==="
ls -la /var/www/ | grep pr

echo ""
echo "=== 前端文件 ==="
ls -lh /var/www/pr-frontend/

echo ""
echo "=== Nginx 配置 ==="
sudo nginx -t

echo ""
echo "=== 后端日志（最近 10 行）==="
sudo journalctl -u pr-business-backend -n 10 --no-pager
ENDSSH
```

---

## 📞 联系方式

**遇到问题？**
- 技术负责人：[联系方式]
- 紧急联系：[联系方式]

---

**最后更新**: 2026-02-05
**维护人**: DevOps Team
