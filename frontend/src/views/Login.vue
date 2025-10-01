<template>
  <div class="login-container">
    <div class="login-background">
      <div class="gradient-blob blob-1"></div>
      <div class="gradient-blob blob-2"></div>
      <div class="gradient-blob blob-3"></div>
    </div>
    
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <Logo :size="64" />
          <h2>Gotux</h2>
          <p>简约 · 高效 · 专业的图床管理系统</p>
        </div>
      </template>
      
      <el-form :model="form" :rules="rules" ref="formRef" @submit.prevent="handleLogin">
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            placeholder="用户名"
            size="large"
            :prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            size="large"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleLogin"
            style="width: 100%"
          >
            登录
          </el-button>
        </el-form-item>
        
        <div class="footer-links">
          <router-link to="/register">还没有账号？立即注册</router-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import Logo from '@/components/Logo.vue'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref()
const loading = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    
    await userStore.login(form.username, form.password)
    
    ElMessage.success('登录成功')
    router.push('/')
  } catch (error) {
    console.error('Login error:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 50%, #4a5568 100%);
}

.login-background {
  position: absolute;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.gradient-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.3;
  animation: float 20s infinite ease-in-out;
}

.blob-1 {
  width: 500px;
  height: 500px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  top: -10%;
  left: -10%;
  animation-delay: 0s;
}

.blob-2 {
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  bottom: -10%;
  right: -10%;
  animation-delay: -7s;
}

.blob-3 {
  width: 350px;
  height: 350px;
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: -14s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  33% {
    transform: translate(30px, -50px) scale(1.1);
  }
  66% {
    transform: translate(-20px, 30px) scale(0.9);
  }
}

.login-card {
  width: 440px;
  position: relative;
  z-index: 10;
  backdrop-filter: blur(20px);
  background: rgba(255, 255, 255, 0.95) !important;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25) !important;
  border: 1px solid rgba(255, 255, 255, 0.3) !important;
  animation: fadeInUp 0.6s ease;
}

.card-header {
  text-align: center;
  padding: 20px 0;
}

.card-header h2 {
  margin: 16px 0 8px 0;
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 1px;
}

.card-header p {
  margin: 0;
  color: var(--text-secondary);
  font-size: 13px;
  font-weight: 400;
}

.footer-links {
  text-align: center;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid var(--border-color);
}

.footer-links a {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.3s ease;
}

.footer-links a:hover {
  color: var(--primary-dark);
  transform: translateX(2px);
}

:deep(.el-form-item) {
  margin-bottom: 24px;
}

:deep(.el-input__wrapper) {
  padding: 12px 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06) !important;
}

:deep(.el-input__wrapper:focus-within) {
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1) !important;
}

:deep(.el-button) {
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 0.5px;
}
</style>
