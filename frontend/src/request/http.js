import request from './request'

// config用来覆盖 默认intance的配置
// 如：某个接口请求需要更长， timeout=30000
class Http {
  get(url, params = {}, config = {}) {
    return request.get(url, { params, ...config })
  }

  post(url, data = {}, config = {}) {
    return request.post(url, data, config)
  }

  put(url, data = {}, config = {}) {
    return request.put(url, data, config)
  }

  delete(url, data = {}, config = {}) {
    return request.delete(url, { data, ...config })
  }
}

export default new Http()
