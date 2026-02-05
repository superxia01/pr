# PR Business 部署指南

**最后更新**: 2026-02-05
**版本**: v4.0 (符合 KeenChase 部署标准)

---

## 📋 目录

- [部署架构](#部署架构)
- [前置要求](#前置要求)
- [首次部署](#首次部署)
- [日常部署](#日常部署)
- [服务管理](#服务管理)
- [故障排查](#故障排查)
- [部署 Checklist](#部署-checklist)

---

## 🏗️ 部署架构

### 服务器信息

| 服务器 | IP | 用途 |
|--------|-----|------|
| 上海服务器 | 101.35.120.199 | 应用服务器 (Go + Nginx) |
| 杭州服务器 | 47.110.82.96 | 数据库服务器 (PostgreSQL) |

### 目录结构（标准）

```
/var/www/
├── pr-backend/              # 后端目录
│   ├── pr-business          # 可执行文件
│   ├── .env                 # 环境变量（不提交到 Git）
│   ├── .env.example         # 环境变量模板
│   └── server.log           # 服务日志
│
└── pr-frontend/             # 前端目录
    ├── index.html
    └── assets/
```

### 服务配置

| 组件 | 配置 | 说明 |
|------|------|------|
| **后端服务** | `pr-business-backend.service` | systemd 管理 |
| **前端服务** | Nginx 静态文件 | `/var/www/pr-frontend` |
| **域名** | pr.crazyaigc.com | HTTPS |
| **API** | pr.crazyaigc.com/api/v1 | 反向代理到 :8081 |

---

## ✅ 前置要求

### 1. 本地环境

**前端**：
```bash
cd frontend
node --version  # >= 18.0.0
npm --version   # >= 9.0.0
```

**后端**：
```bash
cd backend
go version      # >= 1.21
```

### 2. SSH 配置

确保 `~/.ssh/config` 中有上海服务器配置：

```ssh
Host shanghai-tencent
    HostName 101.35.120.199
    User ubuntu
    IdentityFile ~/.ssh/xia_mac_shanghai_secure
    ServerAliveInterval 60
```

测试连接：
```bash
ssh shanghai-tencent "echo '连接成功'"
```

### 3. 服务器权限

- sudo 权限（配置 systemd、nginx）
- ubuntu 用户权限

---

## 🚀 首次部署

### Step 1: 创建目录

```bash
ssh shanghai-tencent << 'ENDSSH'
# 创建标准目录
sudo mkdir -p /var/www/pr-backend
sudo mkdir -p /var/www/pr-frontend

# 设置权限
sudo chown -R ubuntu:ubuntu /var/www/pr-*

# 验证
ls -la /var/www/ | grep pr
ENDSSH
```

### Step 2: 创建环境变量

```bash
ssh shanghai-tencent << 'ENDSSH'
sudo tee /var/www/pr-backend/.env << 'EOF'
# ============================================
# 应用配置
# ============================================
APP_ENV=production
APP_PORT=8081
APP_NAME="PR Business"
APP_DEBUG=false

# ============================================
# 数据库配置（通过 SSH 隧道）
# ============================================
DB_HOST=localhost
DB_PORT=5432
DB_USER=nexus_user
DB_PASSWORD=hRJ9NSJApfeyFDraaDgkYowY
DB_NAME=pr_business_db
DB_SSLMODE=disable

# ============================================
# Auth Center 配置
# ============================================
AUTH_CENTER_URL=https://os.crazyaigc.com
AUTH_CENTER_REDIRECT_URI=https://pr.crazyaigc.com/api/v1/auth/callback

# ============================================
# 前端配置
# ============================================
FRONTEND_URL=https://pr.crazyaigc.com

# ============================================
# JWT 配置
# ============================================
JWT_SECRET=151jmeLlr7ZSi9L4KXIhrJ/CfTFBY2PV5CezmfUlLzw=
JWT_ACCESS_TOKEN_EXPIRE=24h
JWT_REFRESH_TOKEN_EXPIRE=168h
EOF

# 设置权限
sudo chmod 600 /var/www/pr-backend/.env
sudo chown ubuntu:ubuntu /var/www/pr-backend/.env
ENDSSH
```

### Step 3: 创建 systemd 服务

```bash
ssh shanghai-tencent << 'ENDSSH'
sudo tee /etc/systemd/system/pr-business-backend.service << 'EOF'
[Unit]
Description=PR Business Backend API
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/var/www/pr-backend
ExecStart=/var/www/pr-backend/pr-business
Restart=always
RestartSec=5
Environment="PORT=8081"
EnvironmentFile=/var/www/pr-backend/.env

# 日志
StandardOutput=append:/var/www/pr-backend/server.log
StandardError=append:/var/www/pr-backend/server.log

[Install]
WantedBy=multi-user.target
EOF

# 重载并启用服务
sudo systemctl daemon-reload
sudo systemctl enable pr-business-backend
ENDSSH
```

### Step 4: 配置 Nginx

```bash
ssh shanghai-tencent << 'ENDSSH'
sudo tee /etc/nginx/sites-available/pr-business << 'EOF'
# PR Business 配置 - V3.5 Vite + React + Go
# 域名: pr.crazyaigc.com

# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name pr.crazyaigc.com;

    # Let's Encrypt 验证
    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    # 其他请求重定向到 HTTPS
    location / {
        return 301 https://$host$request_uri;
    }
}

# HTTPS 主配置
server {
    listen 443 ssl;
    server_name pr.crazyaigc.com;

    # SSL 证书配置
    ssl_certificate /etc/letsencrypt/live/pr.crazyaigc.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/pr.crazyaigc.com/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;

    # 所有 API 请求代理到 Go 后端
    location /api {
        proxy_pass http://127.0.0.1:8081;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # CORS headers
        add_header Access-Control-Allow-Origin https://pr.crazyaigc.com always;
        add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS" always;
        add_header Access-Control-Allow-Headers "Authorization, Content-Type, X-CSRF-Token" always;
        add_header Access-Control-Allow-Credentials true always;

        if ($request_method = 'OPTIONS') {
            add_header Access-Control-Allow-Origin https://pr.crazyaigc.com always;
            add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS" always;
            add_header Access-Control-Allow-Headers "Authorization, Content-Type, X-CSRF-Token" always;
            add_header Access-Control-Allow-Credentials true always;
            add_header Access-Control-Max-Age 1728000;
            add_header Content-Type 'text/plain charset=UTF-8';
            add_header Content-Length 0;
            return 204;
        }
    }

    # 健康检查
    location /health {
        proxy_pass http://127.0.0.1:8081/health;
        access_log off;
    }

    # Vite 静态文件 (SPA)
    location / {
        root /var/www/pr-frontend;
        try_files $uri $uri/ /index.html;

        # 缓存静态资源
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript
               application/x-javascript application/xml+rss
               application/json application/javascript;
}
EOF

# 启用配置
sudo ln -s /etc/nginx/sites-available/pr-business /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重载 Nginx
sudo systemctl reload nginx
ENDSSH
```

### Step 5: 首次部署

```bash
cd /path/to/pr-business

# 部署全部
bash deploy/deploy.sh all
```

---

## 🔄 日常部署

### 使用部署脚本

```bash
cd /path/to/pr-business

# 部署全部（前端 + 后端）
bash deploy/deploy.sh all

# 仅部署前端
bash deploy/deploy.sh frontend

# 仅部署后端
bash deploy/deploy.sh backend
```

### 手动部署（备选）

**前端**：
```bash
cd frontend

# 1. 构建
npm run build

# 2. 部署
rsync -avz --delete \
  --exclude '*.map' \
  --exclude '.DS_Store' \
  dist/ \
  shanghai-tencent:/var/www/pr-frontend/

# 3. 重载 Nginx
ssh shanghai-tencent "sudo systemctl reload nginx"
```

**后端**：
```bash
cd backend

# 1. 交叉编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pr-business-linux

# 2. 上传
scp pr-business-linux shanghai-tencent:/var/www/pr-backend/

# 3. 重启服务
ssh shanghai-tencent << 'ENDSSH'
cd /var/www/pr-backend

# 备份
sudo cp pr-business pr-business.backup.$(date +%Y%m%d_%H%M%S)

# 替换
sudo mv pr-business-linux pr-business
sudo chmod +x pr-business

# 重启
sudo systemctl restart pr-business-backend

# 检查
sudo systemctl status pr-business-backend
ENDSSH
```

---

## 🔧 服务管理

### 后端服务（systemd）

```bash
# 查看状态
ssh shanghai-tencent "sudo systemctl status pr-business-backend"

# 启动
ssh shanghai-tencent "sudo systemctl start pr-business-backend"

# 停止
ssh shanghai-tencent "sudo systemctl stop pr-business-backend"

# 重启
ssh shanghai-tencent "sudo systemctl restart pr-business-backend"

# 查看日志
ssh shanghai-tencent "sudo journalctl -u pr-business-backend -f"

# 查看最近的日志
ssh shanghai-tencent "sudo journalctl -u pr-business-backend -n 100"
```

### 前端服务（Nginx）

```bash
# 测试配置
ssh shanghai-tencent "sudo nginx -t"

# 重载配置（不中断服务）
ssh shanghai-tencent "sudo systemctl reload nginx"

# 重启服务
ssh shanghai-tencent "sudo systemctl restart nginx"

# 查看状态
ssh shanghai-tencent "sudo systemctl status nginx"

# 查看 Nginx 日志
ssh shanghai-tencent "sudo tail -f /var/log/nginx/error.log"
```

### 查看服务器日志

```bash
# 后端日志（journalctl）
ssh shanghai-tencent "sudo journalctl -u pr-business-backend -f"

# 后端日志（文件）
ssh shanghai-tencent "tail -f /var/www/pr-backend/server.log"

# Nginx 访问日志
ssh shanghai-tencent "sudo tail -f /var/log/nginx/access.log"

# Nginx 错误日志
ssh shanghai-tencent "sudo tail -f /var/log/nginx/error.log"
```

---

## 🔍 故障排查

### 问题 1: 后端服务启动失败

```bash
# 检查服务状态
ssh shanghai-tencent "sudo systemctl status pr-business-backend"

# 查看详细日志
ssh shanghai-tencent "sudo journalctl -u pr-business-backend -n 50 --no-pager"

# 检查环境变量
ssh shanghai-tencent "cat /var/www/pr-backend/.env"

# 手动启动测试
ssh shanghai-tencent "cd /var/www/pr-backend && ./pr-business"
```

### 问题 2: 前端无法访问

```bash
# 检查 Nginx 配置
ssh shanghai-tencent "sudo nginx -t"

# 检查前端文件
ssh shanghai-tencent "ls -la /var/www/pr-frontend/"

# 检查 Nginx 错误日志
ssh shanghai-tencent "sudo tail -n 50 /var/log/nginx/error.log"

# 重载 Nginx
ssh shanghai-tencent "sudo systemctl reload nginx"
```

### 问题 3: API 请求失败

```bash
# 测试后端直连
ssh shanghai-tencent "curl -s http://127.0.0.1:8081/api/v1/service-providers"

# 检查数据库连接
ssh shanghai-tencent "sudo systemctl status pg-tunnel"

# 测试数据库连接
ssh shanghai-tencent "psql -h localhost -U nexus_user -d pr_business_db -c 'SELECT 1'"
```

### 问题 4: 目录不符合标准

```bash
# 检查当前目录
ssh shanghai-tencent "ls -la /var/www/ | grep pr"

# 正确应该是：
# /var/www/pr-backend
# /var/www/pr-frontend

# 如果有错误目录（如 pr-business, pr-business-frontend），删除它们
ssh shanghai-tencent "sudo rm -rf /var/www/pr-business /var/www/pr-business-frontend"
```

---

## 📝 部署 Checklist

每次部署前必须检查以下项目：

### 部署前检查

- [ ] **本地测试通过**
  - [ ] 前端构建成功：`cd frontend && npm run build`
  - [ ] 后端编译成功：`cd backend && go build`
  - [ ] 本地功能测试通过

- [ ] **代码已提交**
  - [ ] 重要修改已 commit
  - [ ] commit message 清晰
  - [ ] 必要时已推送到远程

- [ ] **确认部署目标**
  - [ ] 前端？后端？全部？
  - [ ] 生产环境？测试环境？

### 部署中检查

- [ ] **构建阶段**
  - [ ] 前端构建无错误
  - [ ] 后端交叉编译正确（Linux 二进制）
  - [ ] 文件大小合理（前端 < 2MB，后端 < 30MB）

- [ ] **上传阶段**
  - [ ] 文件上传到正确目录：
    - 前端：`/var/www/pr-frontend/`
    - 后端：`/var/www/pr-backend/`
  - [ ] 未上传 `.env` 文件（环境变量与代码分离）
  - [ ] 上传速度正常

### 部署后检查

- [ ] **服务状态**
  - [ ] 后端服务 running：`sudo systemctl status pr-business-backend`
  - [ ] 前端文件存在：`ls -la /var/www/pr-frontend/`
  - [ ] Nginx 配置正确：`sudo nginx -t`

- [ ] **功能测试**
  - [ ] 网站可访问：https://pr.crazyaigc.com
  - [ ] 登录功能正常
  - [ ] API 请求正常（检查浏览器控制台）
  - [ ] 静态资源加载正常（无 404）

- [ ] **性能检查**
  - [ ] 页面加载速度正常
  - [ ] API 响应时间正常
  - [ ] 无明显错误日志

### 回滚准备

如果部署出现问题：

- [ ] **后端回滚**
  ```bash
  ssh shanghai-tencent << 'ENDSSH'
  cd /var/www/pr-backend
  sudo mv pr-business pr-business.failed.$(date +%Y%m%d_%H%M%S)
  sudo mv pr-business.backup.YYYYMMDD_HHMMSS pr-business
  sudo systemctl restart pr-business-backend
  ENDSSH
  ```

- [ ] **前端回滚**
  ```bash
  # 保留最近 3 个版本的备份
  ssh shanghai-tencent "sudo cp -r /var/www/pr-frontend /var/www/pr-frontend.backup.$(date +%Y%m%d_%H%M%S)"

  # 从备份恢复
  rsync -avz --delete shanghai-tencent:/var/www/pr-frontend.backup.YYYYMMDD_HHMMSS/ dist/
  ```

---

## 🎯 快速参考

### 常用命令

```bash
# 部署
bash deploy/deploy.sh all

# 查看日志
ssh shanghai-tencent "sudo journalctl -u pr-business-backend -f"

# 重启服务
ssh shanghai-tencent "sudo systemctl restart pr-business-backend"

# 测试 API
curl https://pr.crazyaigc.com/api/v1/service-providers

# 查看 Nginx 日志
ssh shanghai-tencent "sudo tail -f /var/log/nginx/error.log"
```

### 重要路径

| 类型 | 路径 |
|------|------|
| 后端目录 | `/var/www/pr-backend` |
| 前端目录 | `/var/www/pr-frontend` |
| 环境变量 | `/var/www/pr-backend/.env` |
| 服务配置 | `/etc/systemd/system/pr-business-backend.service` |
| Nginx 配置 | `/etc/nginx/sites-available/pr-business` |

### 端口

| 服务 | 端口 |
|------|------|
| HTTPS | 443 |
| HTTP | 80 |
| 后端 API | 8081 |
| PostgreSQL (via tunnel) | 5432 |

---

## 📚 相关文档

- **[KeenChase 部署标准](../../keenchase-standards/deployment-and-operations.md)** - 通用部署规范
- **[SSH 配置指南](../../keenchase-standards/ssh-setup.md)** - SSH 密钥配置
- **[数据库使用指南](../../keenchase-standards/database-guide.md)** - 数据库连接

---

**文档维护**: 如有疑问或更新需求，请联系技术负责人。
