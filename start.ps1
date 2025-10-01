# Gotux 启动脚本

Write-Host "===================================" -ForegroundColor Green
Write-Host "    Gotux 图床管理系统启动" -ForegroundColor Green
Write-Host "===================================" -ForegroundColor Green
Write-Host ""

# 检查 Go 环境
Write-Host "检查 Go 环境..." -ForegroundColor Yellow
if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "错误: 未找到 Go 环境，请先安装 Go 1.21+" -ForegroundColor Red
    exit 1
}
$goVersion = go version
Write-Host "✓ $goVersion" -ForegroundColor Green
Write-Host ""

# 检查 Node 环境
Write-Host "检查 Node 环境..." -ForegroundColor Yellow
if (!(Get-Command node -ErrorAction SilentlyContinue)) {
    Write-Host "错误: 未找到 Node.js 环境，请先安装 Node.js 18+" -ForegroundColor Red
    exit 1
}
$nodeVersion = node --version
Write-Host "✓ Node $nodeVersion" -ForegroundColor Green
Write-Host ""

# 启动后端
Write-Host "启动后端服务..." -ForegroundColor Yellow
Set-Location backend

# 检查并创建 .env 文件
if (!(Test-Path ".env")) {
    Write-Host "创建 .env 配置文件..." -ForegroundColor Yellow
    Copy-Item ".env.example" ".env"
}

Write-Host "正在启动后端服务 (端口 8080)..." -ForegroundColor Cyan
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$PWD'; go run main.go"

Set-Location ..
Start-Sleep -Seconds 3

# 启动前端
Write-Host ""
Write-Host "启动前端服务..." -ForegroundColor Yellow
Set-Location frontend

# 检查 node_modules
if (!(Test-Path "node_modules")) {
    Write-Host "安装前端依赖..." -ForegroundColor Yellow
    npm install
}

Write-Host "正在启动前端服务 (端口 5173)..." -ForegroundColor Cyan
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$PWD'; npm run dev"

Set-Location ..

Write-Host ""
Write-Host "===================================" -ForegroundColor Green
Write-Host "启动完成！" -ForegroundColor Green
Write-Host "===================================" -ForegroundColor Green
Write-Host ""
Write-Host "访问地址:" -ForegroundColor Yellow
Write-Host "  前端: http://localhost:5173" -ForegroundColor Cyan
Write-Host "  后端: http://localhost:8080" -ForegroundColor Cyan
Write-Host ""
Write-Host "默认管理员账户:" -ForegroundColor Yellow
Write-Host "  用户名: admin" -ForegroundColor Cyan
Write-Host "  密码: admin123" -ForegroundColor Cyan
Write-Host ""
Write-Host "⚠️  请在首次登录后立即修改默认密码！" -ForegroundColor Red
Write-Host ""
