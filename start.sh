#!/bin/bash

echo "==================================="
echo "    Gotux 图床管理系统启动"
echo "==================================="
echo ""

# 检查 Go 环境
echo "检查 Go 环境..."
if ! command -v go &> /dev/null; then
    echo "错误: 未找到 Go 环境，请先安装 Go 1.21+"
    exit 1
fi
echo "✓ $(go version)"
echo ""

# 检查 Node 环境
echo "检查 Node 环境..."
if ! command -v node &> /dev/null; then
    echo "错误: 未找到 Node.js 环境，请先安装 Node.js 18+"
    exit 1
fi
echo "✓ Node $(node --version)"
echo ""

# 启动后端
echo "启动后端服务..."
cd backend

# 检查并创建 .env 文件
if [ ! -f ".env" ]; then
    echo "创建 .env 配置文件..."
    cp .env.example .env
fi

echo "正在启动后端服务 (端口 8080)..."
go run main.go &
BACKEND_PID=$!

cd ..
sleep 3

# 启动前端
echo ""
echo "启动前端服务..."
cd frontend

# 检查 node_modules
if [ ! -d "node_modules" ]; then
    echo "安装前端依赖..."
    npm install
fi

echo "正在启动前端服务 (端口 5173)..."
npm run dev &
FRONTEND_PID=$!

cd ..

echo ""
echo "==================================="
echo "启动完成！"
echo "==================================="
echo ""
echo "访问地址:"
echo "  前端: http://localhost:5173"
echo "  后端: http://localhost:8080"
echo ""
echo "默认管理员账户:"
echo "  用户名: admin"
echo "  密码: admin123"
echo ""
echo "⚠️  请在首次登录后立即修改默认密码！"
echo ""
echo "按 Ctrl+C 停止服务"

# 等待中断信号
trap "kill $BACKEND_PID $FRONTEND_PID; exit" INT TERM
wait
