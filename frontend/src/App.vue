<template>
  <el-config-provider :locale="zhCn">
    <router-view />
  </el-config-provider>
</template>

<script setup>
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

:root {
  /* 主色调 - 优雅的紫色渐变 */
  --primary-color: #667eea;
  --primary-dark: #5568d3;
  --primary-light: #7c8ff0;
  --secondary-color: #764ba2;
  
  /* 中性色 - 更柔和的灰色调 */
  --bg-color: #f8f9fc;
  --bg-secondary: #ffffff;
  --text-primary: #1a202c;
  --text-secondary: #4a5568;
  --text-tertiary: #718096;
  --border-color: #e2e8f0;
  
  /* 功能色 */
  --success-color: #48bb78;
  --warning-color: #ed8936;
  --danger-color: #f56565;
  --info-color: #4299e1;
  
  /* 阴影 */
  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  
  /* 圆角 */
  --radius-sm: 6px;
  --radius-md: 10px;
  --radius-lg: 16px;
  --radius-xl: 24px;
}

body {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background: var(--bg-color);
  color: var(--text-primary);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  line-height: 1.6;
}

#app {
  width: 100%;
  min-height: 100vh;
}

/* 滚动条美化 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: var(--bg-color);
}

::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a0aec0;
}

/* Element Plus 样式覆盖 */
.el-card {
  border-radius: var(--radius-lg) !important;
  box-shadow: var(--shadow-md) !important;
  border: 1px solid var(--border-color) !important;
}

.el-button {
  border-radius: var(--radius-md) !important;
  font-weight: 500 !important;
  transition: all 0.3s ease !important;
}

.el-button--primary {
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--secondary-color) 100%) !important;
  border: none !important;
}

.el-button--primary:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg) !important;
}

.el-input__wrapper {
  border-radius: var(--radius-md) !important;
  box-shadow: var(--shadow-sm) !important;
  transition: all 0.3s ease !important;
}

.el-input__wrapper:hover {
  box-shadow: var(--shadow-md) !important;
}

.el-table {
  border-radius: var(--radius-lg) !important;
}

/* 动画 */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fade-in-up {
  animation: fadeInUp 0.5s ease;
}
</style>
