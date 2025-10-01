import request from '@/utils/request'

export function uploadImages(formData) {
  return request.post('/images/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export function getImages(params) {
  return request.get('/images', { params })
}

export function getImageDetail(id) {
  return request.get(`/images/${id}`)
}

export function updateImage(id, data) {
  return request.put(`/images/${id}`, data)
}

export function deleteImage(id) {
  return request.delete(`/images/${id}`)
}

export function batchDeleteImages(imageIds) {
  return request.post('/images/batch-delete', { image_ids: imageIds })
}

export function getImageLinks(id) {
  return request.get(`/images/${id}/links`)
}

// 随机图片 API (不需要认证)
export function getRandomImage(params) {
  // 直接使用 fetch 避免添加 Authorization 头
  const queryString = new URLSearchParams(params).toString()
  const url = `/api/random${queryString ? '?' + queryString : ''}`
  return fetch(url).then(res => res.json())
}

export function getRandomImageUrl(params) {
  const queryString = new URLSearchParams(params).toString()
  return `/api/random/image${queryString ? '?' + queryString : ''}`
}
