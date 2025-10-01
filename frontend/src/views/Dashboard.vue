<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #409EFF">
              <el-icon :size="30"><Picture /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.image_count }}</div>
              <div class="stat-label">图片总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #67C23A">
              <el-icon :size="30"><FolderOpened /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ formatBytes(stats.storage_used) }}</div>
              <div class="stat-label">已用空间</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #E6A23C">
              <el-icon :size="30"><View /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.total_views }}</div>
              <div class="stat-label">总浏览量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #F56C6C">
              <el-icon :size="30"><User /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ userStore.userInfo?.username }}</div>
              <div class="stat-label">当前用户</div>
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
  padding: 10px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 5px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
}

.quick-actions {
  display: flex;
  gap: 10px;
}
</style>
