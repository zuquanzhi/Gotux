<template>
  <div class="upload-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>上传图片</span>
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
          开始上传 ({{ fileList.length }})
        </el-button>
        <el-button @click="clearFiles">清空</el-button>
      </div>
    </el-card>
    
    <el-card v-if="uploadedImages.length > 0" style="margin-top: 20px">
      <template #header>
        <div class="card-header">
          <span>上传结果</span>
        </div>
      </template>
      
      <el-row :gutter="20">
        <el-col :span="6" v-for="image in uploadedImages" :key="image.id">
          <el-card :body-style="{ padding: '0px' }" shadow="hover">
            <el-image
              :src="`/uploads/${image.file_path}`"
              fit="cover"
              style="width: 100%; height: 200px;"
            />
            <div style="padding: 14px;">
              <div class="image-name">{{ image.original_name }}</div>
              <div class="image-size">{{ formatBytes(image.file_size) }}</div>
              <el-button
                type="primary"
                size="small"
                style="margin-top: 10px; width: 100%"
                @click="showLinks(image.id)"
              >
                获取链接
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>
    
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
import { ref } from 'vue'
import { uploadImages, getImageLinks } from '@/api/image'
import { ElMessage } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'

const uploadRef = ref()
const fileList = ref([])
const uploading = ref(false)
const uploadedImages = ref([])
const linksDialogVisible = ref(false)
const currentLinks = ref({})

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
  max-width: 1200px;
  margin: 0 auto;
  animation: fadeInUp 0.5s ease;
}

.card-header {
  font-weight: 600;
  font-size: 18px;
  color: var(--text-primary);
}

.upload-demo {
  width: 100%;
}

:deep(.el-upload-dragger) {
  border: 2px dashed var(--border-color);
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.03) 0%, rgba(118, 75, 162, 0.03) 100%);
  transition: all 0.3s ease;
  padding: 60px 20px;
}

:deep(.el-upload-dragger:hover) {
  border-color: var(--primary-color);
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.08) 0%, rgba(118, 75, 162, 0.08) 100%);
  transform: scale(1.02);
}

:deep(.el-icon--upload) {
  font-size: 64px;
  color: var(--primary-color);
  margin-bottom: 16px;
}

:deep(.el-upload__text) {
  font-size: 16px;
  color: var(--text-secondary);
}

:deep(.el-upload__text em) {
  color: var(--primary-color);
  font-weight: 600;
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
</style>
