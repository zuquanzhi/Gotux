<template>
  <div class="dashboard">
    <div class="welcome-banner">
      <div class="welcome-content">
        <h1>欢迎回来, {{ userStore.userInfo?.username }}! 👋</h1>
        <p>让我们开始管理您的图片资源</p>
      </div>
    </div>

    <el-row :gutter="24" class="stats-row">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card-wrapper">
          <div class="stat-card gradient-blue">
            <div class="stat-icon-container">
              <div class="stat-icon">
                <el-icon :size="32"><Picture /></el-icon>
              </div>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.image_count }}</div>
              <div class="stat-label">图片总数</div>
            </div>
            <div class="stat-trend">
              <span class="trend-up">↑ 活跃</span>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card-wrapper">
          <div class="stat-card gradient-green">
            <div class="stat-icon-container">
              <div class="stat-icon">
                <el-icon :size="32"><FolderOpened /></el-icon>
              </div>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ formatBytes(stats.storage_used) }}</div>
              <div class="stat-label">已用空间</div>
            </div>
            <div class="stat-trend">
              <span class="trend-info">存储中</span>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card-wrapper">
          <div class="stat-card gradient-orange">
            <div class="stat-icon-container">
              <div class="stat-icon">
                <el-icon :size="32"><View /></el-icon>
              </div>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.total_views }}</div>
              <div class="stat-label">总浏览量</div>
            </div>
            <div class="stat-trend">
              <span class="trend-up">↑ 增长</span>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card-wrapper">
          <div class="stat-card gradient-purple">
            <div class="stat-icon-container">
              <div class="stat-icon">
                <el-icon :size="32"><User /></el-icon>
              </div>
            </div>
            <div class="stat-content">
              <div class="stat-value username-value">{{ userStore.userInfo?.username }}</div>
              <div class="stat-label">当前用户</div>
            </div>
            <div class="stat-trend">
              <span class="trend-online">在线</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>快速操作</span>
            </div>
          </template>
          <div class="quick-actions">
            <el-button type="primary" :icon="Upload" @click="goToUpload">
              上传图片
            </el-button>
            <el-button type="success" :icon="Picture" @click="goToImages">
              查看图片
            </el-button>
            <el-button type="info" :icon="User" @click="goToProfile">
              个人中心
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最近上传</span>
            </div>
          </template>
          <el-table :data="recentImages" style="width: 100%" v-loading="loading">
            <el-table-column prop="original_name" label="文件名" show-overflow-tooltip />
            <el-table-column label="预览" width="100">
              <template #default="{ row }">
                <el-image
                  :src="`/uploads/${row.file_path}`"
                  :preview-src-list="[`/uploads/${row.file_path}`]"
                  fit="cover"
                  style="width: 60px; height: 60px; border-radius: 4px;"
                />
              </template>
            </el-table-column>
            <el-table-column prop="file_size" label="大小" width="120">
              <template #default="{ row }">
                {{ formatBytes(row.file_size) }}
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="上传时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getStats } from '@/api/auth'
import { getImages } from '@/api/image'
import { Upload, Picture, User, View, FolderOpened } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const stats = ref({
  image_count: 0,
  storage_used: 0,
  total_views: 0
})
const recentImages = ref([])
const loading = ref(false)

const formatBytes = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('zh-CN')
}

const goToUpload = () => router.push('/upload')
const goToImages = () => router.push('/images')
const goToProfile = () => router.push('/profile')

const fetchData = async () => {
  try {
    loading.value = true
    const [statsData, imagesData] = await Promise.all([
      getStats(),
      getImages({ page: 1, page_size: 5 })
    ])
    stats.value = statsData
    recentImages.value = imagesData.images
  } catch (error) {
    console.error('Fetch data error:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.dashboard {
  max-width: 1400px;
  margin: 0 auto;
}

.welcome-banner {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: var(--radius-xl);
  padding: 40px;
  margin-bottom: 32px;
  color: white;
  box-shadow: var(--shadow-xl);
  position: relative;
  overflow: hidden;
}

.welcome-banner::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -10%;
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, transparent 70%);
  border-radius: 50%;
}

.welcome-content h1 {
  font-size: 32px;
  font-weight: 700;
  margin: 0 0 8px 0;
  position: relative;
  z-index: 1;
}

.welcome-content p {
  font-size: 16px;
  opacity: 0.9;
  margin: 0;
  position: relative;
  z-index: 1;
}

.stats-row {
  margin-bottom: 32px;
}

.stat-card-wrapper {
  transition: all 0.3s ease;
  height: 100%;
}

.stat-card-wrapper:hover {
  transform: translateY(-8px);
}

.stat-card-wrapper :deep(.el-card__body) {
  padding: 0;
}

.stat-card {
  padding: 24px;
  position: relative;
  overflow: hidden;
  border-radius: var(--radius-lg);
  min-height: 140px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  right: 0;
  width: 100px;
  height: 100px;
  opacity: 0.1;
  border-radius: 50%;
}

.gradient-blue {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.gradient-green {
  background: linear-gradient(135deg, #48bb78 0%, #38a169 100%);
  color: white;
}

.gradient-orange {
  background: linear-gradient(135deg, #ed8936 0%, #dd6b20 100%);
  color: white;
}

.gradient-purple {
  background: linear-gradient(135deg, #9f7aea 0%, #805ad5 100%);
  color: white;
}

.stat-icon-container {
  margin-bottom: 16px;
}

.stat-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: var(--radius-md);
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
}

.stat-content {
  margin-bottom: 12px;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 4px;
  line-height: 1.2;
}

.username-value {
  font-size: 24px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
  font-weight: 500;
}

.stat-trend {
  display: flex;
  align-items: center;
}

.stat-trend span {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
}

.trend-up {
  color: rgba(255, 255, 255, 0.95);
}

.trend-info {
  color: rgba(255, 255, 255, 0.9);
}

.trend-online {
  color: rgba(255, 255, 255, 0.95);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 18px;
  color: var(--text-primary);
}

.quick-actions {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.quick-actions .el-button {
  flex: 1;
  min-width: 140px;
}

:deep(.el-table) {
  font-size: 14px;
}

:deep(.el-table th) {
  background: var(--bg-color);
  color: var(--text-secondary);
  font-weight: 600;
}
</style>
