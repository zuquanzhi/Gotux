import request from '@/utils/request'

export function getAllUsers(params) {
  return request.get('/admin/users', { params })
}

export function updateUserStatus(id, status) {
  return request.put(`/admin/users/${id}/status`, { status })
}

export function updateUserQuota(id, storageQuota) {
  return request.put(`/admin/users/${id}/quota`, { storage_quota: storageQuota })
}

export function getAllImages(params) {
  return request.get('/admin/images', { params })
}

export function getSystemStats() {
  return request.get('/admin/stats')
}
