<template>
  <div class="upload-page">
    <el-row :gutter="20">
      <!-- 上传区域 -->
      <el-col :span="24">
        <el-card class="upload-card">
          <template #header>
            <div class="card-header">
              <span>📤 上传图片</span>
            </div>
          </template>
          
          <el-upload
            ref="uploadRef"
            class="upload-demo"
            drag
            multiple
            :auto-upload="false"
            :on-change="handleFileChange"
            :file-list="fileList"
            accept="image/*"
            list-type="picture"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                支持 jpg/png/gif/webp 格式，单个文件不超过 10MB
              </div>
            </template>
          </el-upload>
          
          <div class="upload-actions" v-if="fileList.length > 0">
            <el-button type="primary" :loading="uploading" @click="handleUpload">
              <el-icon><upload-filled /></el-icon>
              开始上传 ({{ fileList.length }})
            </el-button>
            <el-button @click="clearFiles">
              <el-icon><delete /></el-icon>
              清空列表
            </el-button>
          </div>
        </el-card>
        
        <!-- 上传结果 -->
        <el-card v-if="uploadedImages.length > 0" style="margin-top: 20px" class="results-card">
          <template #header>
            <div class="card-header">
              <span>✅ 上传成功 ({{ uploadedImages.length }})</span>
              <el-button type="text" @click="uploadedImages = []">清除记录</el-button>
            </div>
          </template>
          
          <el-row :gutter="16">
            <el-col :xs="24" :sm="12" :md="8" v-for="image in uploadedImages" :key="image.id">
              <el-card :body-style="{ padding: '0px' }" shadow="hover" class="result-image-card">
                <el-image
                  :src="`/uploads/${image.file_path}`"
                  fit="cover"
                  style="width: 100%; height: 200px;"
                  :preview-src-list="[`/uploads/${image.file_path}`]"
                  :z-index="9999"
                  :preview-teleported="true"
                />
                <div class="result-image-info">
                  <div class="image-name">{{ image.original_name }}</div>
                  <div class="image-meta">
                    <span class="image-size">{{ formatBytes(image.file_size) }}</span>
                    <span class="image-dimensions">{{ image.width }} × {{ image.height }}</span>
                  </div>
                  <el-button
                    type="primary"
                    size="small"
                    style="margin-top: 10px; width: 100%"
                    @click="showLinks(image.id)"
                  >
                    <el-icon><link /></el-icon>
                    获取链接
                  </el-button>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 上传统计 -->
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :xs="24" :sm="8" v-for="(stat, index) in statsData" :key="index">
        <el-card class="stat-card-inline">
          <div class="stat-content-inline">
            <div class="stat-icon-inline" :class="`icon-${stat.color}`">
              {{ stat.icon }}
            </div>
            <div class="stat-info-inline">
              <div class="stat-label-inline">{{ stat.label }}</div>
              <div class="stat-value-inline" :class="{ success: stat.color === 'green' }">{{ stat.value }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 链接对话框 -->
    <el-dialog v-model="linksDialogVisible" title="图片链接" width="600px">
      <el-form label-width="100px">
        <el-form-item label="URL">
          <el-input v-model="currentLinks.url" readonly>
            <template #append>
              <el-button @click="copyToClipboard(currentLinks.url)">复制</el-button>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item label="Markdown">
          <el-input v-model="currentLinks.markdown" readonly>
            <template #append>
              <el-button @click="copyToClipboard(currentLinks.markdown)">复制</el-button>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item label="HTML">
          <el-input v-model="currentLinks.html" readonly>
            <template #append>
              <el-button @click="copyToClipboard(currentLinks.html)">复制</el-button>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item label="BBCode">
          <el-input v-model="currentLinks.bbcode" readonly>
            <template #append>
              <el-button @click="copyToClipboard(currentLinks.bbcode)">复制</el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { uploadImages, getImageLinks } from '@/api/image'
import { ElMessage } from 'element-plus'
import { UploadFilled, Delete, Link } from '@element-plus/icons-vue'

const router = useRouter()
const uploadRef = ref()
const fileList = ref([])
const uploading = ref(false)
const uploadedImages = ref([])
const linksDialogVisible = ref(false)
const currentLinks = ref({})

// 计算总文件大小
const totalSize = computed(() => {
  return fileList.value.reduce((total, file) => total + (file.size || 0), 0)
})

// 统计数据
const statsData = computed(() => [
  {
    icon: '📋',
    label: '待上传',
    value: fileList.value.length,
    color: 'blue'
  },
  {
    icon: '✅',
    label: '已完成',
    value: uploadedImages.value.length,
    color: 'green'
  },
  {
    icon: '📦',
    label: '总大小',
    value: formatBytes(totalSize.value),
    color: 'purple'
  }
])

const handleFileChange = (file, files) => {
  fileList.value = files
}

const handleUpload = async () => {
  if (fileList.value.length === 0) {
    ElMessage.warning('请选择要上传的文件')
    return
  }
  
  try {
    uploading.value = true
    
    const formData = new FormData()
    fileList.value.forEach(file => {
      formData.append('files', file.raw)
    })
    
    const data = await uploadImages(formData)
    
    uploadedImages.value = data.images
    fileList.value = []
    uploadRef.value.clearFiles()
    
    ElMessage.success(data.message)
    
    if (data.errors && data.errors.length > 0) {
      data.errors.forEach(error => {
        ElMessage.warning(error)
      })
    }
  } catch (error) {
    console.error('Upload error:', error)
  } finally {
    uploading.value = false
  }
}

const clearFiles = () => {
  fileList.value = []
  uploadRef.value.clearFiles()
}

const showLinks = async (imageId) => {
  try {
    const data = await getImageLinks(imageId)
    currentLinks.value = data.links
    linksDialogVisible.value = true
  } catch (error) {
    console.error('Get links error:', error)
  }
}

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

const formatBytes = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}
</script>

<style scoped>
.upload-page {
  animation: fadeInUp 0.5s ease;
}

.card-header {
  font-weight: 600;
  font-size: 18px;
  color: var(--text-primary);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.upload-demo {
  width: 100%;
}

.upload-card {
  min-height: 500px;
  display: flex;
  flex-direction: column;
}

.upload-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.results-card {
  animation: fadeInUp 0.5s ease;
}

.result-image-card {
  margin-bottom: 16px;
  transition: all 0.3s ease;
}

.result-image-card:hover {
  transform: translateY(-4px);
}

.result-image-info {
  padding: 14px;
}

.image-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 8px;
}

:deep(.el-upload-dragger) {
  border: 2px dashed var(--border-color);
  border-radius: var(--radius-lg);
  background: 
    linear-gradient(135deg, rgba(102, 126, 234, 0.03) 0%, rgba(118, 75, 162, 0.03) 100%),
    radial-gradient(circle at 20% 30%, rgba(102, 126, 234, 0.05) 0%, transparent 50%),
    radial-gradient(circle at 80% 70%, rgba(118, 75, 162, 0.05) 0%, transparent 50%);
  transition: all 0.3s ease;
  padding: 80px 40px;
  position: relative;
  overflow: hidden;
}

:deep(.el-upload-dragger::before) {
  content: '';
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  border-radius: var(--radius-lg);
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: -1;
}

:deep(.el-upload-dragger:hover) {
  border-color: transparent;
  background: 
    linear-gradient(135deg, rgba(102, 126, 234, 0.08) 0%, rgba(118, 75, 162, 0.08) 100%),
    radial-gradient(circle at 20% 30%, rgba(102, 126, 234, 0.1) 0%, transparent 50%),
    radial-gradient(circle at 80% 70%, rgba(118, 75, 162, 0.1) 0%, transparent 50%);
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
}

:deep(.el-upload-dragger:hover::before) {
  opacity: 0.1;
}

:deep(.el-icon--upload) {
  font-size: 72px;
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--secondary-color) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 20px;
  filter: drop-shadow(0 2px 8px rgba(102, 126, 234, 0.3));
}

:deep(.el-upload__text) {
  font-size: 17px;
  color: var(--text-secondary);
  font-weight: 500;
  letter-spacing: 0.3px;
}

:deep(.el-upload__text em) {
  color: var(--primary-color);
  font-weight: 700;
  font-style: normal;
}

:deep(.el-upload__tip) {
  margin-top: 16px;
  font-size: 13px;
  color: var(--text-tertiary);
}

.upload-actions {
  margin-top: 32px;
  text-align: center;
  display: flex;
  gap: 16px;
  justify-content: center;
}

.upload-actions .el-button {
  min-width: 140px;
}

.image-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.image-size {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 4px;
}

:deep(.el-card) {
  transition: all 0.3s ease;
}

:deep(.el-card:hover) {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg) !important;
}

:deep(.el-image) {
  border-radius: var(--radius-md);
}



/* 横向统计卡片 */
.stat-card-inline {
  border: 1px solid var(--border-color);
  transition: all 0.3s ease;
}

.stat-card-inline:hover {
  border-color: var(--primary-color);
  box-shadow: var(--shadow-lg);
  transform: translateY(-2px);
}

.stat-content-inline {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px 0;
}

.stat-icon-inline {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
}

.stat-icon-inline.icon-blue {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
}

.stat-icon-inline.icon-green {
  background: linear-gradient(135deg, rgba(72, 187, 120, 0.1) 0%, rgba(56, 161, 105, 0.1) 100%);
}

.stat-icon-inline.icon-purple {
  background: linear-gradient(135deg, rgba(159, 122, 234, 0.1) 0%, rgba(128, 90, 213, 0.1) 100%);
}

.stat-info-inline {
  flex: 1;
  min-width: 0;
}

.stat-label-inline {
  font-size: 13px;
  color: var(--text-tertiary);
  margin-bottom: 4px;
  font-weight: 500;
}

.stat-value-inline {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.2;
}

.stat-value-inline.success {
  color: #67c23a;
}

/* 响应式 */
@media (max-width: 768px) {
  .stat-content-inline {
    flex-direction: column;
    text-align: center;
  }
  
  .stat-icon-inline {
    width: 56px;
    height: 56px;
    font-size: 28px;
  }
  
  .stat-value-inline {
    font-size: 28px;
  }
}
</style>
