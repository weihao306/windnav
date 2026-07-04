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

## 生成随机测试数据

确认后端已按上文启动后，可在项目根目录执行以下 PowerShell 脚本。脚本会登录默认管理员账号，创建或复用一组分类和标签，并随机新增 36 个站点。

```powershell
$ErrorActionPreference = "Stop"
$baseUrl = "http://localhost:8080"

$login = Invoke-RestMethod `
  -Method Post `
  -Uri "$baseUrl/api/auth/login" `
  -ContentType "application/json" `
  -Body (@{ username = "admin"; password = "admin123456" } | ConvertTo-Json)

$token = $login.data.accessToken
if (-not $token) { $token = $login.data.token }
if (-not $token) { throw "登录成功但响应中未找到 token" }

$headers = @{ Authorization = "Bearer $token" }

function Invoke-WindNavJson($Method, $Path, $Body) {
  Invoke-RestMethod `
    -Method $Method `
    -Uri "$baseUrl$Path" `
    -Headers $headers `
    -ContentType "application/json" `
    -Body ($Body | ConvertTo-Json -Depth 10)
}

$categorySeeds = @(
  @{ name = "开发工具"; slug = "dev-tools"; description = "提升研发效率的工具集合"; icon = "code"; color = "#2563eb"; sortOrder = 10; isVisible = $true },
  @{ name = "设计灵感"; slug = "design-inspiration"; description = "界面、图标与视觉灵感站点"; icon = "palette"; color = "#db2777"; sortOrder = 20; isVisible = $true },
  @{ name = "AI 工具"; slug = "ai-tools"; description = "AI 写作、绘图、编程与效率工具"; icon = "sparkles"; color = "#7c3aed"; sortOrder = 30; isVisible = $true },
  @{ name = "学习资源"; slug = "learning"; description = "课程、文档、教程与知识库"; icon = "book"; color = "#16a34a"; sortOrder = 40; isVisible = $true },
  @{ name = "效率办公"; slug = "productivity"; description = "协作、笔记、自动化与项目管理"; icon = "briefcase"; color = "#ea580c"; sortOrder = 50; isVisible = $true },
  @{ name = "云服务"; slug = "cloud-services"; description = "部署、存储、监控与基础设施"; icon = "cloud"; color = "#0891b2"; sortOrder = 60; isVisible = $true }
)

$tagSeeds = @(
  @{ name = "免费"; slug = "free"; color = "#22c55e" },
  @{ name = "开源"; slug = "open-source"; color = "#0ea5e9" },
  @{ name = "中文"; slug = "chinese"; color = "#f97316" },
  @{ name = "推荐"; slug = "recommended"; color = "#ef4444" },
  @{ name = "教程"; slug = "tutorial"; color = "#8b5cf6" },
  @{ name = "SaaS"; slug = "saas"; color = "#14b8a6" },
  @{ name = "API"; slug = "api"; color = "#6366f1" },
  @{ name = "团队协作"; slug = "collaboration"; color = "#ec4899" }
)

$categories = @()
foreach ($seed in $categorySeeds) {
  try {
    $created = Invoke-WindNavJson "Post" "/api/admin/categories" $seed
    $categories += $created.data
  } catch {
    $existing = Invoke-RestMethod -Method Get -Uri "$baseUrl/api/admin/categories" -Headers $headers
    $categories += @($existing.data | Where-Object { $_.slug -eq $seed.slug })[0]
  }
}

$tags = @()
foreach ($seed in $tagSeeds) {
  try {
    $created = Invoke-WindNavJson "Post" "/api/admin/tags" $seed
    $tags += $created.data
  } catch {
    $existing = Invoke-RestMethod -Method Get -Uri "$baseUrl/api/admin/tags" -Headers $headers
    $tags += @($existing.data | Where-Object { $_.slug -eq $seed.slug })[0]
  }
}

$siteNames = @("Atlas", "Nova", "Pixel", "Stack", "Flow", "Spark", "Orbit", "Craft", "Mint", "Vector", "Nimbus", "Rocket", "Bridge", "Matrix", "Beacon", "Studio", "Lab", "Hub", "Pilot", "Forge", "Garden", "Desk", "Board", "Wave")
$domains = @("example.com", "windnav.dev", "toolbox.local", "demo.site", "nav.test", "workspace.app")
$descriptions = @("适合日常高频使用的精选工具。", "提供清晰易用的工作流与模板。", "可用于快速验证导航站展示效果。", "包含丰富案例和实用资源。", "面向个人与团队的效率平台。", "帮助你发现更优质的互联网资源。")
$sitesCreated = 0

for ($i = 1; $i -le 36; $i++) {
  $category = $categories[($i - 1) % $categories.Count]
  $name = $siteNames | Get-Random
  $domain = $domains | Get-Random
  $selectedTags = @($tags | Get-Random -Count (Get-Random -Minimum 1 -Maximum 4))

  $body = @{
    categoryId = [uint32]$category.id
    title = "$name $i"
    url = "https://$($name.ToLower())-$i.$domain"
    description = $descriptions | Get-Random
    iconUrl = ""
    fallbackIcon = $name.Substring(0, 1).ToUpper()
    sortOrder = $i * 10
    isPinned = ($i % 7 -eq 0)
    isVisible = $true
    tagIds = @($selectedTags | ForEach-Object { [uint32]$_.id })
  }

  Invoke-WindNavJson "Post" "/api/admin/sites" $body | Out-Null
  $sitesCreated++
}

Write-Host "已创建随机站点：$sitesCreated"
```

如需重新生成一批全新的测试数据，可先按下方步骤重置本地数据库，再重新执行脚本。

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
