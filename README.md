# windnav

微风导航，简单轻快的自建导航页。

## 技术栈

- 后端：Go、Gin、GORM
- 数据库：SQLite 或 PostgreSQL
- 前端：Vue 3、Vite、TypeScript、Tailwind CSS
- 部署：单容器 SQLite 或 Docker Compose PostgreSQL

## 本地开发

启动后端：

```powershell
Set-Location backend
go run ./cmd/server
```

启动前端：

```powershell
Set-Location frontend
npm install
npm run dev
```

默认管理员：

- 用户名：admin
- 密码：admin123456

前端开发地址通常是 http://localhost:5173，后端地址是 http://localhost:8080。

## 数据库模式

默认使用 SQLite：

```powershell
$env:DB_DRIVER="sqlite"
$env:DB_SQLITE_PATH="./data/windnav.db"
go run ./cmd/server
```

使用 PostgreSQL：

```powershell
$env:DB_DRIVER="postgres"
$env:DATABASE_URL="postgres://windnav:windnav_password@localhost:5432/windnav?sslmode=disable"
go run ./cmd/server
```

## 部署

SQLite 最简部署：

```powershell
docker compose -f deploy/docker-compose.sqlite.yml up --build
```

PostgreSQL 标准部署：

```powershell
docker compose -f deploy/docker-compose.yml up --build
```

生产环境必须修改 JWT_SECRET 和 ADMIN_PASSWORD。
