import request from '@/utils/request'

export function login(username, password) {
  return request.post('/auth/login', { username, password })
}

export function register(username, email, password) {
  return request.post('/auth/register', { username, email, password })
}

export function getProfile() {
  return request.get('/user/profile')
}

export function updateProfile(data) {
  return request.put('/user/profile', data)
}

export function changePassword(oldPassword, newPassword) {
  return request.post('/user/change-password', {
    old_password: oldPassword,
    new_password: newPassword
  })
}

export function getStats() {
  return request.get('/user/stats')
}
