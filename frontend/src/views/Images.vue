<template>
  <div class="images-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>图片管理</span>
          <div class="header-actions">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索图片"
              style="width: 200px; margin-right: 10px;"
              clearable
              @clear="fetchImages"
              @keyup.enter="fetchImages"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-button type="primary" @click="fetchImages">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button
              type="danger"
              :disabled="selectedImages.length === 0"
              @click="handleBatchDelete"
            >
              <el-icon><Delete /></el-icon>
              批量删除
            </el-button>
          </div>
        </div>
      </template>
      
      <el-row :gutter="20" v-loading="loading">
        <el-col :xs="24" :sm="12" :md="8" :lg="6" v-for="image in images" :key="image.id">
          <el-card :body-style="{ padding: '0px' }" shadow="hover" class="image-card">
            <el-checkbox
              v-model="selectedImages"
              :label="image.id"
              class="image-checkbox"
            />
            <el-image
              :src="`/uploads/${image.file_path}`"
              :preview-src-list="[`/uploads/${image.file_path}`]"
              fit="cover"
              style="width: 100%; height: 200px; cursor: pointer;"
              @click="viewImage(image)"
            />
            <div style="padding: 14px;">
              <div class="image-name" :title="image.original_name">
                {{ image.original_name }}
              </div>
              <div class="image-info">
                <span>{{ formatBytes(image.file_size) }}</span>
                <span>{{ image.width }} x {{ image.height }}</span>
              </div>
              <div class="image-stats">
                <el-icon><View /></el-icon>
                <span>{{ image.stats?.view_count || 0 }}</span>
              </div>
              <div class="image-actions">
                <el-button size="small" @click="showLinks(image.id)">
                  <el-icon><Link /></el-icon>
                  链接
                </el-button>
                <el-button size="small" type="primary" @click="editImage(image)">
                  <el-icon><Edit /></el-icon>
                  编辑
                </el-button>
                <el-button size="small" type="danger" @click="deleteImage(image.id)">
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
      
      <el-empty v-if="!loading && images.length === 0" description="暂无图片" />
      
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[12, 24, 48, 96]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchImages"
          @current-change="fetchImages"
        />
      </div>
    </el-card>
    
    <!-- 编辑对话框 -->
    <el-dialog v-model="editDialogVisible" title="编辑图片信息" width="500px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="描述">
          <el-input v-model="editForm.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="editForm.tags" placeholder="用逗号分隔多个标签" />
        </el-form-item>
        <el-form-item label="公开">
          <el-switch v-model="editForm.is_public" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleEdit">保存</el-button>
      </template>
    </el-dialog>
    
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
import { ref, onMounted } from 'vue'
import { getImages as fetchImagesApi, updateImage as updateImageApi, deleteImage as deleteImageApi, batchDeleteImages, getImageLinks } from '@/api/image'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Delete, View, Link, Edit } from '@element-plus/icons-vue'

const loading = ref(false)
const images = ref([])
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)
const searchKeyword = ref('')
const selectedImages = ref([])
const editDialogVisible = ref(false)
const linksDialogVisible = ref(false)
const currentLinks = ref({})
const editForm = ref({
  id: null,
  description: '',
  tags: '',
  is_public: true
})

const fetchImages = async () => {
  try {
    loading.value = true
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }
    
    const data = await fetchImagesApi(params)
    images.value = data.images
    total.value = data.total
  } catch (error) {
    console.error('Fetch images error:', error)
  } finally {
    loading.value = false
  }
}

const editImage = (image) => {
  editForm.value = {
    id: image.id,
    description: image.description || '',
    tags: image.tags || '',
    is_public: image.is_public
  }
  editDialogVisible.value = true
}

const handleEdit = async () => {
  try {
    await updateImageApi(editForm.value.id, {
      description: editForm.value.description,
      tags: editForm.value.tags,
      is_public: editForm.value.is_public
    })
    
    ElMessage.success('更新成功')
    editDialogVisible.value = false
    fetchImages()
  } catch (error) {
    console.error('Update image error:', error)
  }
}

const deleteImage = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这张图片吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await deleteImageApi(id)
    ElMessage.success('删除成功')
    fetchImages()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete image error:', error)
    }
  }
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedImages.value.length} 张图片吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await batchDeleteImages(selectedImages.value)
    ElMessage.success('批量删除成功')
    selectedImages.value = []
    fetchImages()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Batch delete error:', error)
    }
  }
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

const viewImage = (image) => {
  // 点击图片时使用 Element Plus 的预览功能
}

const formatBytes = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

onMounted(() => {
  fetchImages()
})
</script>

<style scoped>
.images-page {
  animation: fadeInUp 0.5s ease;
}

@media (max-width: 768px) {
  .image-card {
    margin-bottom: 16px;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .header-actions {
    width: 100%;
    flex-wrap: wrap;
  }
  
  .header-actions .el-input {
    width: 100% !important;
    margin-right: 0 !important;
    margin-bottom: 8px;
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 18px;
  color: var(--text-primary);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.image-card {
  position: relative;
  margin-bottom: 24px;
  transition: all 0.3s ease;
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.image-card:hover {
  transform: translateY(-6px);
  box-shadow: var(--shadow-lg) !important;
}

.image-checkbox {
  position: absolute;
  top: 12px;
  left: 12px;
  z-index: 10;
  background: rgba(255, 255, 255, 0.95);
  padding: 4px 8px;
  border-radius: var(--radius-md);
  backdrop-filter: blur(10px);
}

:deep(.el-checkbox__label) {
  display: none;
}

:deep(.el-image) {
  position: relative;
  overflow: hidden;
}

:deep(.el-image::after) {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(180deg, rgba(0,0,0,0) 0%, rgba(0,0,0,0.1) 100%);
  pointer-events: none;
}

.image-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 10px;
}

.image-info {
  font-size: 12px;
  color: var(--text-tertiary);
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
  padding: 6px 10px;
  background: var(--bg-secondary);
  border-radius: var(--radius-sm);
}

.image-stats {
  font-size: 12px;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 12px;
  padding: 4px 8px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  border-radius: var(--radius-sm);
  width: fit-content;
}

.image-stats .el-icon {
  color: var(--primary-color);
}

.image-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.image-actions .el-button {
  flex: 1;
  min-width: 70px;
}

.pagination {
  margin-top: 32px;
  display: flex;
  justify-content: center;
}

:deep(.el-pagination) {
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: var(--radius-lg);
}

:deep(.el-dialog__header) {
  padding: 24px;
  border-bottom: 1px solid var(--border-color);
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
}

:deep(.el-dialog__title) {
  font-weight: 600;
  font-size: 18px;
  color: var(--text-primary);
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-empty) {
  padding: 80px 0;
}

:deep(.el-empty__image) {
  width: 160px;
}

:deep(.el-empty__description) {
  color: var(--text-tertiary);
  font-size: 14px;
}
</style>
