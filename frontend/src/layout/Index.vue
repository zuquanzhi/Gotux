<template>
  <div class="layout">
    <el-container>
      <el-aside width="260px">
        <div class="logo-container">
          <Logo :size="48" />
          <div class="logo-text">
            <h2>Gotux</h2>
            <p>图床管理系统</p>
          </div>
        </div>
        <el-menu
          :default-active="$route.path"
          router
          :background-color="'transparent'"
          text-color="#e2e8f0"
          active-text-color="#ffffff"
          class="custom-menu"
        >
          <el-menu-item index="/dashboard">
            <el-icon><DataAnalysis /></el-icon>
            <span>仪表盘</span>
          </el-menu-item>
          <el-menu-item index="/upload">
            <el-icon><Upload /></el-icon>
            <span>上传图片</span>
          </el-menu-item>
          <el-menu-item index="/images">
            <el-icon><Picture /></el-icon>
            <span>图片管理</span>
          </el-menu-item>
          <el-sub-menu v-if="userStore.isAdmin" index="/admin">
            <template #title>
              <el-icon><Setting /></el-icon>
              <span>系统管理</span>
            </template>
            <el-menu-item index="/admin/users">用户管理</el-menu-item>
            <el-menu-item index="/admin/images">图片管理</el-menu-item>
          </el-sub-menu>
        </el-menu>
      </el-aside>
      <el-container>
        <el-header>
          <div class="header-content">
            <div class="breadcrumb">
              <h3>{{ currentTitle }}</h3>
            </div>
            <div class="user-info">
              <el-dropdown @command="handleCommand">
                <span class="user-dropdown">
                  <el-avatar :size="32">{{ userStore.userInfo?.username[0].toUpperCase() }}</el-avatar>
                  <span class="username">{{ userStore.userInfo?.username }}</span>
                  <el-icon><ArrowDown /></el-icon>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                    <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </el-header>
        <el-main>
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import Logo from '@/components/Logo.vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const currentTitle = computed(() => route.meta.title || '首页')

const handleCommand = (command) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/login')
    ElMessage.success('已退出登录')
  } else if (command === 'profile') {
    router.push('/profile')
  }
}
</script>

<style scoped>
.layout {
  height: 100vh;
  background: var(--bg-color);
}

.el-container {
  height: 100%;
}

.el-aside {
  background: linear-gradient(180deg, #1a202c 0%, #2d3748 100%);
  color: #fff;
  height: 100vh;
  overflow-y: auto;
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
}

.logo-container {
  height: 100px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
}

.logo-text {
  margin-top: 12px;
  text-align: center;
}

.logo-text h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #f093fb 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 0.5px;
}

.logo-text p {
  margin: 4px 0 0 0;
  font-size: 12px;
  color: #a0aec0;
  font-weight: 300;
}

.custom-menu {
  border: none;
  padding: 12px 8px;
}

.custom-menu :deep(.el-menu-item) {
  margin: 6px 0;
  border-radius: 10px;
  transition: all 0.3s ease;
  font-weight: 500;
}

.custom-menu :deep(.el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.1) !important;
  transform: translateX(4px);
}

.custom-menu :deep(.el-menu-item.is-active) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.custom-menu :deep(.el-sub-menu__title) {
  border-radius: 10px;
  margin: 6px 0;
  font-weight: 500;
  transition: all 0.3s ease;
}

.custom-menu :deep(.el-sub-menu__title:hover) {
  background: rgba(255, 255, 255, 0.1) !important;
}

.el-header {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  display: flex;
  align-items: center;
  padding: 0 32px;
  border-bottom: 1px solid var(--border-color);
}

.header-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.breadcrumb h3 {
  margin: 0;
  color: var(--text-primary);
  font-size: 20px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  padding: 8px 16px;
  border-radius: 12px;
  transition: all 0.3s ease;
  background: var(--bg-color);
  border: 1px solid var(--border-color);
}

.user-dropdown:hover {
  background: #fff;
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.username {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.el-main {
  background: var(--bg-color);
  padding: 24px;
  overflow-y: auto;
}

/* 滚动条美化 - 侧边栏 */
.el-aside::-webkit-scrollbar {
  width: 6px;
}

.el-aside::-webkit-scrollbar-track {
  background: transparent;
}

.el-aside::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 3px;
}

.el-aside::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>
