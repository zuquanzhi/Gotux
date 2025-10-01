<template>
  <div class="profile-page">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>个人信息</span>
            </div>
          </template>
          
          <div class="profile-info">
            <el-avatar :size="100" style="margin-bottom: 20px;">
              {{ userStore.userInfo?.username[0].toUpperCase() }}
            </el-avatar>
            
            <div class="info-item">
              <span class="label">用户名：</span>
              <span class="value">{{ userStore.userInfo?.username }}</span>
            </div>
            
            <div class="info-item">
              <span class="label">邮箱：</span>
              <span class="value">{{ userStore.userInfo?.email }}</span>
            </div>
            
            <div class="info-item">
              <span class="label">角色：</span>
              <el-tag v-if="userStore.isAdmin" type="danger">管理员</el-tag>
              <el-tag v-else>普通用户</el-tag>
            </div>
            
            <div class="info-item">
              <span class="label">状态：</span>
              <el-tag type="success">{{ userStore.userInfo?.status }}</el-tag>
            </div>
            
            <div class="info-item">
              <span class="label">注册时间：</span>
              <span class="value">{{ formatDate(userStore.userInfo?.created_at) }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>修改个人信息</span>
            </div>
          </template>
          
          <el-form :model="profileForm" label-width="100px">
            <el-form-item label="邮箱">
              <el-input v-model="profileForm.email" />
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="handleUpdateProfile">更新信息</el-button>
            </el-form-item>
          </el-form>
        </el-card>
        
        <el-card style="margin-top: 20px;">
          <template #header>
            <div class="card-header">
              <span>修改密码</span>
            </div>
          </template>
          
          <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="100px">
            <el-form-item label="原密码" prop="oldPassword">
              <el-input v-model="passwordForm.oldPassword" type="password" show-password />
            </el-form-item>
            
            <el-form-item label="新密码" prop="newPassword">
              <el-input v-model="passwordForm.newPassword" type="password" show-password />
            </el-form-item>
            
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input v-model="passwordForm.confirmPassword" type="password" show-password />
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="handleChangePassword">修改密码</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { updateProfile, changePassword } from '@/api/auth'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const passwordFormRef = ref()

const profileForm = reactive({
  email: ''
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = {
  oldPassword: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' }
  ],
  confirmPassword: [{ required: true, validator: validateConfirmPassword, trigger: 'blur' }]
}

const handleUpdateProfile = async () => {
  try {
    await updateProfile(profileForm)
    await userStore.fetchProfile()
    ElMessage.success('更新成功')
  } catch (error) {
    console.error('Update profile error:', error)
  }
}

const handleChangePassword = async () => {
  try {
    await passwordFormRef.value.validate()
    
    await changePassword(passwordForm.oldPassword, passwordForm.newPassword)
    
    ElMessage.success('密码修改成功')
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
    passwordFormRef.value.resetFields()
  } catch (error) {
    console.error('Change password error:', error)
  }
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  profileForm.email = userStore.userInfo?.email || ''
})
</script>

<style scoped>
.profile-page {
  max-width: 1200px;
  margin: 0 auto;
  animation: fadeInUp 0.5s ease;
}

.card-header {
  font-weight: 600;
  font-size: 18px;
  color: var(--text-primary);
}

.profile-info {
  text-align: center;
}

:deep(.el-avatar) {
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--secondary-color) 100%);
  font-size: 40px;
  font-weight: 600;
  box-shadow: var(--shadow-md);
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.3s ease;
}

.info-item:hover {
  padding-left: 8px;
  background: var(--bg-secondary);
  margin: 0 -20px;
  padding-left: 28px;
  padding-right: 28px;
}

.info-item:last-child {
  border-bottom: none;
}

.label {
  font-weight: 600;
  color: var(--text-secondary);
  font-size: 14px;
}

.value {
  color: var(--text-primary);
  font-size: 14px;
}

:deep(.el-card) {
  transition: all 0.3s ease;
}

:deep(.el-card:hover) {
  box-shadow: var(--shadow-lg) !important;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: var(--text-secondary);
}

:deep(.el-input__inner) {
  border-radius: var(--radius-md);
}

:deep(.el-button) {
  border-radius: var(--radius-md);
  font-weight: 500;
}
</style>
