# 部署说明

WindNav 支持 SQLite 和 PostgreSQL 两种数据库模式。

## SQLite 最简模式

适合个人服务器、NAS 或单机部署。应用容器内同时包含后端服务和前端静态资源，数据保存到 /app/data/windnav.db。

```powershell
docker compose -f deploy/docker-compose.sqlite.yml up --build
```

## PostgreSQL 标准模式

适合更正式的服务化部署。

```powershell
docker compose -f deploy/docker-compose.yml up --build
```

## 关键环境变量

- APP_ENV：production 或 development
- HTTP_ADDR：后端监听地址，默认 :8080
- DB_DRIVER：sqlite 或 postgres
- DB_SQLITE_PATH：SQLite 文件路径
- DATABASE_URL：PostgreSQL 连接串
- JWT_SECRET：登录令牌密钥，生产环境必须修改
- ADMIN_USERNAME：初始化管理员用户名
- ADMIN_PASSWORD：初始化管理员密码
- STATIC_DIR：前端静态资源目录
