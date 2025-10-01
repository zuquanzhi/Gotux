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
