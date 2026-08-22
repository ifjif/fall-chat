import http from '../request/http.js'

export const refreshApi = (params, config) => {
  return http.get("/auth/refresh", params, config)
}
