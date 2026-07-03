# 本地运行与调试

本文档说明如何在本地启动 WindNav，适用于日常开发和功能调试。

## 推荐模式

本地调试建议使用 SQLite。它不依赖外部数据库，启动后端时会自动创建本地数据库文件、表结构和默认管理员账号。

## 环境要求

- Go
- Node.js
- npm

可在项目根目录执行以下命令确认环境：

```powershell
go version
node --version
npm --version
```

## 启动后端

打开第一个终端，在项目根目录执行：

```powershell
Set-Location backend
$env:DB_DRIVER="sqlite"
$env:DB_SQLITE_PATH="./data/windnav.db"
$env:JWT_SECRET="dev-secret"
$env:ADMIN_USERNAME="admin"
$env:ADMIN_PASSWORD="admin123456"
go run ./cmd/server
```

后端默认监听：

```text
http://localhost:8080
```

健康检查地址：

```text
http://localhost:8080/healthz
```

## 启动前端

打开第二个终端，在项目根目录执行：

```powershell
Set-Location frontend
npm install
npm run dev
```

前端默认监听：

```text
http://localhost:5173
```

Vite 已在本地开发模式下配置代理，前端请求 `/api` 会自动转发到后端：

```text
http://localhost:8080
```

## 访问应用

公开导航页：

```text
http://localhost:5173
```

后台入口：

```text
http://localhost:5173/admin
```

默认管理员账号：

```text
admin / admin123456
```

## 重置本地数据

如果需要重新初始化本地数据：

1. 停止后端服务。
2. 删除 SQLite 数据库文件：

```text
backend/data/windnav.db
```

3. 重新启动后端：

```powershell
go run ./cmd/server
```

后端会重新创建数据库、表结构和默认管理员账号。

## 常见问题

### 前端接口请求失败

确认后端正在运行，并检查健康检查地址是否可访问：

```text
http://localhost:8080/healthz
```

### 后台登录失败

确认启动后端时设置的管理员环境变量：

```powershell
$env:ADMIN_USERNAME="admin"
$env:ADMIN_PASSWORD="admin123456"
```

如果数据库已经存在，修改环境变量不会覆盖已有管理员密码。需要删除本地数据库文件后重新启动后端。

### 端口被占用

后端端口可通过环境变量修改：

```powershell
$env:HTTP_ADDR=":8081"
go run ./cmd/server
```

如果修改后端端口，也需要同步调整前端开发代理配置。
