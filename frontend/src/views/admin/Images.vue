<template>
  <div class="admin-images">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>所有图片</span>
        </div>
      </template>
      
      <el-table :data="images" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
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
        <el-table-column prop="original_name" label="文件名" show-overflow-tooltip />
        <el-table-column prop="user.username" label="上传者" width="120" />
        <el-table-column prop="file_size" label="大小" width="100">
          <template #default="{ row }">
            {{ formatBytes(row.file_size) }}
          </template>
        </el-table-column>
        <el-table-column prop="stats.view_count" label="浏览量" width="100">
          <template #default="{ row }">
            {{ row.stats?.view_count || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="上传时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button
              size="small"
              type="danger"
              @click="deleteImage(row.id)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchImages"
          @current-change="fetchImages"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAllImages } from '@/api/admin'
import { deleteImage as deleteImageApi } from '@/api/image'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const images = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const fetchImages = async () => {
  try {
    loading.value = true
    const data = await getAllImages({
      page: currentPage.value,
      page_size: pageSize.value
    })
    images.value = data.images
    total.value = data.total
  } catch (error) {
    console.error('Fetch images error:', error)
  } finally {
    loading.value = false
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

onMounted(() => {
  fetchImages()
})
</script>

<style scoped>
.admin-images {
  max-width: 1400px;
  margin: 0 auto;
}

.card-header {
  font-weight: bold;
  font-size: 16px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
