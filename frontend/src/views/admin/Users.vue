<template>
  <div class="admin-users">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
        </div>
      </template>
      
      <el-table :data="users" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="role" label="角色" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.role === 'admin'" type="danger">管理员</el-tag>
            <el-tag v-else>普通用户</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.status === 'active'" type="success">激活</el-tag>
            <el-tag v-else type="danger">禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="storage_quota" label="存储配额" width="150">
          <template #default="{ row }">
            <span v-if="row.storage_quota === 0">无限制</span>
            <span v-else>{{ formatBytes(row.storage_quota) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'active' && row.role !== 'admin'"
              size="small"
              type="warning"
              @click="updateStatus(row.id, 'disabled')"
            >
              禁用
            </el-button>
            <el-button
              v-if="row.status === 'disabled'"
              size="small"
              type="success"
              @click="updateStatus(row.id, 'active')"
            >
              激活
            </el-button>
            <el-button
              v-if="row.role !== 'admin'"
              size="small"
              type="primary"
              @click="showQuotaDialog(row)"
            >
              配额
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
          @size-change="fetchUsers"
          @current-change="fetchUsers"
        />
      </div>
    </el-card>

    <!-- 配额设置对话框 -->
    <el-dialog v-model="quotaDialogVisible" title="设置存储配额" width="500px">
      <el-form :model="quotaForm" label-width="100px">
        <el-form-item label="用户">
          <el-input v-model="quotaForm.username" disabled />
        </el-form-item>
        <el-form-item label="当前使用">
          <el-input :value="formatBytes(quotaForm.storage_used || 0)" disabled />
        </el-form-item>
        <el-form-item label="存储配额">
          <el-input
            v-model.number="quotaForm.storage_quota_mb"
            type="number"
            placeholder="输入配额大小(MB),0表示无限制"
          >
            <template #append>MB</template>
          </el-input>
          <div style="margin-top: 8px; font-size: 12px; color: #909399;">
            设置为 0 表示无限制。1 GB = 1024 MB
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="quotaDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="updateQuota">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAllUsers, updateUserStatus, updateUserQuota } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const users = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const quotaDialogVisible = ref(false)
const quotaForm = ref({
  user_id: 0,
  username: '',
  storage_used: 0,
  storage_quota_mb: 0
})

const fetchUsers = async () => {
  try {
    loading.value = true
    const data = await getAllUsers({
      page: currentPage.value,
      page_size: pageSize.value
    })
    users.value = data.users
    total.value = data.total
  } catch (error) {
    console.error('Fetch users error:', error)
  } finally {
    loading.value = false
  }
}

const updateStatus = async (userId, status) => {
  try {
    const action = status === 'disabled' ? '禁用' : '激活'
    await ElMessageBox.confirm(`确定要${action}此用户吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await updateUserStatus(userId, status)
    ElMessage.success(`${action}成功`)
    fetchUsers()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Update status error:', error)
    }
  }
}

const showQuotaDialog = (user) => {
  quotaForm.value = {
    user_id: user.id,
    username: user.username,
    storage_used: user.storage_used || 0,
    storage_quota_mb: user.storage_quota ? Math.round(user.storage_quota / 1024 / 1024) : 0
  }
  quotaDialogVisible.value = true
}

const updateQuota = async () => {
  try {
    const quotaBytes = quotaForm.value.storage_quota_mb * 1024 * 1024
    
    await updateUserQuota(quotaForm.value.user_id, quotaBytes)
    ElMessage.success('配额设置成功')
    quotaDialogVisible.value = false
    fetchUsers()
  } catch (error) {
    console.error('Update quota error:', error)
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
  fetchUsers()
})
</script>

<style scoped>
.admin-users {
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
