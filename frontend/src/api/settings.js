import request from '@/utils/request'

/**
 * 获取用户设置
 */
export const getSettings = () => {
  return request({
    url: '/user/settings',
    method: 'get'
  })
}

/**
 * 更新用户设置
 */
export const updateSettings = (data) => {
  return request({
    url: '/user/settings',
    method: 'put',
    data
  })
}
