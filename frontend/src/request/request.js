import axios from 'axios'
import router from '@/router'
import { refreshApi } from '@/api/auth'
import { FlMessage } from 'fall-ui'

let isRefreshing = false // 是否正在刷新token
let requestQueue = []   // 带等待刷新完成的请求


const instance = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  }
})

instance.interceptors.request.use((config) => {
  const token = localStorage.getItem("token")
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  return config
}, (error) => {
  return Promise.reject(error)
})

instance.interceptors.response.use((response) => {
  const _data = response.data
  if (_data instanceof Blob) {
    return _data
  }
  const { code, message, data } = _data
  if (message) {
    FlMessage({ type: 'success', content: message })
  }

  return data
}, (error) => {
  const response = error.response
  const status = response.status
  let content = ''

  if (response) {
    const _data = response.data
    const isBlob = _data instanceof Blob
    const { code, _, message } = _data
    // access 进行刷新
    if (status === 401 && (code === 2001 || isBlob) && !error.config._retry) {
      return refreshToken(error.config)
    }

    content = message
  } else if (error.message.includes('timeout')) {
    content = '请求超时，请稍后尝试'
  } else {
    content = '网络异常，请稍后尝试'
  }

  if (content) {
    FlMessage({ type: 'error', content })
  }


  return Promise.reject(error)
})

// access_token 过期进行刷新
const refreshToken = async (originRequestConfig) => {
  if (isRefreshing) {
    return new Promise((resolve) => {
      requestQueue.push(() => {
        const token = localStorage.getItem('token')
        // 用新token重试请求
        originRequestConfig.headers.Authorization = `Bearer ${token}`
        resolve(instance(originRequestConfig))
      });
    })
  }

  // 标记已被处理
  originRequestConfig._retry = true
  isRefreshing = true

  // 调用 refres
  try {
    const data = await refreshApi({}, {
      withCredentials: true // 自动携带cookie
    })

    const newToken = data.access_token
    localStorage.setItem("token", newToken)

    // 执行 等待队列
    requestQueue.forEach(rq => rq())
    requestQueue = []

    // 当前请求重新执行
    originRequestConfig.headers.Authorization = `Bearer ${newToken}`
    return instance(originRequestConfig)
  } catch (err) {
    // refresh 失败，则refresh token 失效，重新login
    requestQueue = []
    FlMessage({ type: 'warning', content: '登录已过期，请重新登录' })
    localStorage.removeItem('token')
    router.push("/login")
    return Promise.reject(err)
  } finally {
    isRefreshing = false
  }

}


export default instance
