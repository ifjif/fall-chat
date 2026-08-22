import http from "../request/http";

export const fileFetchApi = (data, config) => {
  return http.get(`/file/fetch`, data, config)
}

export const fileUploadApi = (data, config) => {
  return http.post("/file/upload", data, config)
}

export const fileUploadChunkApi = (data, config) => {
  return http.post("/file/chunk", data, config)
}

export const fileUploadMergeApi = (data) => {
  return http.post("/file/merge", data)
}

export const fileInstanceCheckApi = (data) => {
  return http.post("/file/instanceCheck", data)
}

export const fileChunkStatusApi = (data) => {
  return http.post("/file/chunkStatus", data)
}

export const fileAccessUrlApi = (data) => {
  return http.get("/file/accessUrl", data)
}
