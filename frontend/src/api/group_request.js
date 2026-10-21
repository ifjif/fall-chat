import http from "@/request/http"

export const listGroupRequestApi = (data) => {
  return http.post("/groupRequest/list", data)
}

export const applyAddGroupApi = (data) => {
  return http.post("/groupRequest/apply", data)
}

export const handleGroupRequestApi = (data) => {
  return http.post("/groupRequest/handle", data)
}
