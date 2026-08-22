import http from '../request/http'

export const passRestCodeApi = (data) => {
  return http.post('/email/passRestCode', data)
}
