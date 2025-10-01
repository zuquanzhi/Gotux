<template>
  <div class="profile-page">
    <el-row :gutter="20">
      <el-col :xs="24" :sm="24" :md="8" :lg="6">
        <el-card class="sticky-card">
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
      
      <el-col :xs="24" :sm="24" :md="16" :lg="18">
        <!-- 基本信息 -->
        <el-card>
          <template #header>
            <div class="card-header">
              <span>修改个人信息</span>
            </div>
          </template>
          
          <el-form :model="profileForm" label-width="120px">
            <el-form-item label="邮箱">
              <el-input v-model="profileForm.email" />
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="handleUpdateProfile">更新信息</el-button>
            </el-form-item>
          </el-form>
        </el-card>
        
        <!-- 链接设置 -->
        <el-card style="margin-top: 20px;">
          <template #header>
            <div class="card-header">
              <span>🔗 链接设置</span>
            </div>
          </template>
          
          <el-form :model="settingsForm" label-width="120px">
            <el-form-item label="自定义域名">
              <el-input 
                v-model="settingsForm.custom_domain" 
                placeholder="例如: https://img.example.com"
              >
                <template #prepend>
                  <el-icon><Link /></el-icon>
                </template>
              </el-input>
              <template #extra>
                <span class="form-tip">配置后，图片链接将使用此域名替换默认域名</span>
              </template>
            </el-form-item>
            
            <el-form-item label="默认链接格式">
              <el-select v-model="settingsForm.default_link_format" style="width: 100%;">
                <el-option label="URL 直链" value="url" />
                <el-option label="Markdown 格式" value="markdown" />
                <el-option label="HTML 格式" value="html" />
                <el-option label="BBCode 格式" value="bbcode" />
              </el-select>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="handleUpdateSettings">保存链接设置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
        
        <!-- 图片处理设置 -->
        <el-card style="margin-top: 20px;">
          <template #header>
            <div class="card-header">
              <span>🎨 图片处理设置</span>
            </div>
          </template>
          
          <el-form :model="settingsForm" label-width="120px">
            <el-form-item label="自动压缩">
              <el-switch v-model="settingsForm.compress_image" />
              <template #extra>
                <span class="form-tip">上传时自动压缩图片以节省空间</span>
              </template>
            </el-form-item>
            
            <el-form-item label="压缩质量" v-if="settingsForm.compress_image">
              <el-slider 
                v-model="settingsForm.compress_quality" 
                :min="1" 
                :max="100"
                show-input
              />
              <template #extra>
                <span class="form-tip">质量越高，文件越大 (建议 70-90)</span>
              </template>
            </el-form-item>
            
            <el-form-item label="启用水印">
              <el-switch v-model="settingsForm.enable_watermark" />
            </el-form-item>
            
            <el-form-item label="水印文字" v-if="settingsForm.enable_watermark">
              <el-input 
                v-model="settingsForm.watermark_text" 
                placeholder="输入水印文字"
              />
            </el-form-item>
            
            <el-form-item label="水印位置" v-if="settingsForm.enable_watermark">
              <el-select v-model="settingsForm.watermark_position" style="width: 100%;">
                <el-option label="左上角" value="top-left" />
                <el-option label="顶部居中" value="top-center" />
                <el-option label="右上角" value="top-right" />
                <el-option label="左侧居中" value="middle-left" />
                <el-option label="正中心" value="center" />
                <el-option label="右侧居中" value="middle-right" />
                <el-option label="左下角" value="bottom-left" />
                <el-option label="底部居中" value="bottom-center" />
                <el-option label="右下角" value="bottom-right" />
              </el-select>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="handleUpdateSettings">保存图片设置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
        
        <!-- 上传限制设置 -->
        <el-card style="margin-top: 20px;">
          <template #header>
            <div class="card-header">
              <span>📁 上传限制设置</span>
            </div>
          </template>
          
          <el-form :model="settingsForm" label-width="120px">
            <el-form-item label="单文件大小限制">
              <el-input-number 
                v-model="maxImageSizeMB" 
                :min="0.1" 
                :max="50"
                :step="0.5"
                :precision="1"
              />
              <span style="margin-left: 10px;">MB</span>
              <template #extra>
                <span class="form-tip">限制单个图片文件的最大大小</span>
              </template>
            </el-form-item>
            
            <el-form-item label="允许的格式">
              <el-select 
                v-model="allowedTypesArray" 
                multiple 
                style="width: 100%;"
                placeholder="选择允许的图片格式"
              >
                <el-option label="JPG" value="jpg" />
                <el-option label="JPEG" value="jpeg" />
                <el-option label="PNG" value="png" />
                <el-option label="GIF" value="gif" />
                <el-option label="WebP" value="webp" />
                <el-option label="BMP" value="bmp" />
                <el-option label="SVG" value="svg" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="存储配额">
              <el-progress 
                :percentage="storagePercentage" 
                :color="storageColor"
                :stroke-width="20"
              >
                <span class="storage-text">
                  {{ formatBytes(settingsForm.used_storage) }} / {{ formatBytes(settingsForm.storage_quota) }}
                </span>
              </el-progress>
            </el-form-item>
            
            <el-form-item label="图片审核">
              <el-switch v-model="settingsForm.enable_image_review" />
              <template #extra>
                <span class="form-tip">上传图片后需要管理员审核才能公开</span>
              </template>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="handleUpdateSettings">保存上传设置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
        
        <!-- 修改密码 -->
        <el-card style="margin-top: 20px;">
          <template #header>
            <div class="card-header">
              <span>🔒 修改密码</span>
            </div>
          </template>
          
          <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="120px">
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
import { ref, reactive, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { updateProfile, changePassword } from '@/api/auth'
import { getSettings, updateSettings } from '@/api/settings'
import { ElMessage } from 'element-plus'
import { Link } from '@element-plus/icons-vue'

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

const settingsForm = reactive({
  custom_domain: '',
  default_link_format: 'url',
  enable_watermark: false,
  watermark_text: '',
  watermark_position: 'bottom-right',
  compress_image: false,
  compress_quality: 80,
  max_image_size: 10485760, // 10MB in bytes
  allowed_image_types: 'jpg,jpeg,png,gif,webp',
  enable_image_review: false,
  storage_quota: 1073741824, // 1GB
  used_storage: 0
})

// 计算属性：MB 转换
const maxImageSizeMB = computed({
  get: () => (settingsForm.max_image_size / 1024 / 1024).toFixed(1),
  set: (val) => {
    settingsForm.max_image_size = Math.round(val * 1024 * 1024)
  }
})

// 计算属性：允许的类型数组
const allowedTypesArray = computed({
  get: () => settingsForm.allowed_image_types ? settingsForm.allowed_image_types.split(',') : [],
  set: (val) => {
    settingsForm.allowed_image_types = val.join(',')
  }
})

// 计算属性：存储使用百分比
const storagePercentage = computed(() => {
  if (settingsForm.storage_quota === 0) return 0
  return Math.min((settingsForm.used_storage / settingsForm.storage_quota) * 100, 100)
})

// 计算属性：存储进度条颜色
const storageColor = computed(() => {
  const percentage = storagePercentage.value
  if (percentage < 50) return '#67c23a'
  if (percentage < 80) return '#e6a23c'
  return '#f56c6c'
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

const handleUpdateSettings = async () => {
  try {
    await updateSettings(settingsForm)
    ElMessage.success('设置保存成功')
    await loadSettings()
  } catch (error) {
    console.error('Update settings error:', error)
    ElMessage.error('保存设置失败')
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

const loadSettings = async () => {
  try {
    const data = await getSettings()
    Object.assign(settingsForm, data.settings)
  } catch (error) {
    console.error('Load settings error:', error)
    ElMessage.error('加载设置失败')
  }
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('zh-CN')
}

const formatBytes = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

onMounted(() => {
  profileForm.email = userStore.userInfo?.email || ''
  loadSettings()
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

.form-tip {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 4px;
  display: block;
}

:deep(.el-form-item__extra) {
  margin-top: 4px;
}

:deep(.el-progress__text) {
  font-size: 12px !important;
  font-weight: 600;
}

.storage-text {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

:deep(.el-slider) {
  margin-right: 20px;
}

:deep(.el-select) {
  width: 100%;
}

:deep(.el-input-number) {
  width: 150px;
}

:deep(.el-switch) {
  --el-switch-on-color: var(--primary-color);
}

:deep(.el-progress) {
  line-height: 1.5;
}

:deep(.el-card__header) {
  border-bottom: 1px solid var(--border-color);
}

.card-header span {
  display: flex;
  align-items: center;
  gap: 8px;
}

.sticky-card {
  position: sticky;
  top: 20px;
}

@media (max-width: 768px) {
  .sticky-card {
    position: relative;
    top: 0;
    margin-bottom: 20px;
  }
}
</style>
