# 快速测试脚本
Write-Host "正在测试 Gotux 项目..." -ForegroundColor Cyan
Write-Host ""

# 测试后端
Write-Host "检查后端文件..." -ForegroundColor Yellow
$backendPath = "D:\program\Gotux\backend"
if (Test-Path "$backendPath\main.go") {
    Write-Host "✓ main.go 存在" -ForegroundColor Green
    Write-Host "✓ go.mod 存在: $(Test-Path "$backendPath\go.mod")" -ForegroundColor Green
} else {
    Write-Host "✗ main.go 不存在" -ForegroundColor Red
}

Write-Host ""
Write-Host "检查前端文件..." -ForegroundColor Yellow
$frontendPath = "D:\program\Gotux\frontend"
if (Test-Path "$frontendPath\package.json") {
    Write-Host "✓ package.json 存在" -ForegroundColor Green
    Write-Host "✓ vite.config.js 存在: $(Test-Path "$frontendPath\vite.config.js")" -ForegroundColor Green
} else {
    Write-Host "✗ package.json 不存在" -ForegroundColor Red
}

Write-Host ""
Write-Host "项目结构:" -ForegroundColor Yellow
Get-ChildItem -Path "D:\program\Gotux" -Directory | Format-Table Name
