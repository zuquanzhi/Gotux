<template>
  <div class="dashboard">
    <!-- 统计概览 -->
    <div class="stats-overview">
      <div class="stat-card" v-for="(stat, index) in statsData" :key="index">
        <div class="stat-icon" :class="`icon-${stat.color}`">
          <component :is="stat.icon" />
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stat.value }}</div>
          <div class="stat-label">{{ stat.label }}</div>
        </div>
      </div>
    </div>
    
    <!-- 最近上传 -->
    <el-row :gutter="20" style="margin-top: 32px">
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
                  :z-index="9999"
                  :preview-teleported="true"
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
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getStats } from '@/api/auth'
import { getImages } from '@/api/image'
import { Picture, FolderOpened, View, Calendar } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const stats = ref({
  image_count: 0,
  storage_used: 0,
  total_views: 0
})
const recentImages = ref([])
const loading = ref(false)

const statsData = computed(() => [
  {
    icon: Picture,
    value: stats.value.image_count,
    label: '图片总数',
    color: 'blue'
  },
  {
    icon: FolderOpened,
    value: formatBytes(stats.value.storage_used),
    label: '存储空间',
    color: 'green'
  },
  {
    icon: View,
    value: stats.value.total_views,
    label: '总浏览量',
    color: 'purple'
  },
  {
    icon: Calendar,
    value: recentImages.value.length,
    label: '最近上传',
    color: 'orange'
  }
])

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
  animation: fadeInUp 0.5s ease;
}

/* 统计概览 */
.stats-overview {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 32px;
}

.stat-card {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 20px;
  transition: all 0.3s ease;
  cursor: default;
}

.stat-card:hover {
  border-color: var(--primary-color);
  box-shadow: var(--shadow-lg);
  transform: translateY(-2px);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
}

.stat-icon.icon-blue {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  color: #667eea;
}

.stat-icon.icon-green {
  background: linear-gradient(135deg, rgba(72, 187, 120, 0.1) 0%, rgba(56, 161, 105, 0.1) 100%);
  color: #48bb78;
}

.stat-icon.icon-purple {
  background: linear-gradient(135deg, rgba(159, 122, 234, 0.1) 0%, rgba(128, 90, 213, 0.1) 100%);
  color: #9f7aea;
}

.stat-icon.icon-orange {
  background: linear-gradient(135deg, rgba(237, 137, 54, 0.1) 0%, rgba(221, 107, 32, 0.1) 100%);
  color: #ed8936;
}

.stat-info {
  flex: 1;
  min-width: 0;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 4px;
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: var(--text-tertiary);
  font-weight: 500;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 18px;
  color: var(--text-primary);
}

:deep(.el-table) {
  font-size: 14px;
}

:deep(.el-table th) {
  background: var(--bg-secondary);
  color: var(--text-secondary);
  font-weight: 600;
}

:deep(.el-table td), :deep(.el-table th) {
  padding: 12px 0;
}

:deep(.el-card) {
  border: 1px solid var(--border-color);
  transition: all 0.3s ease;
}

:deep(.el-card:hover) {
  box-shadow: var(--shadow-md);
}

/* 响应式布局 */
@media (max-width: 1200px) {
  .stats-overview {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-overview {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .stat-card {
    padding: 20px;
  }
  
  .stat-icon {
    width: 48px;
    height: 48px;
    font-size: 20px;
  }
  
  .stat-value {
    font-size: 24px;
  }
  
  .stat-label {
    font-size: 12px;
  }
}
</style>
