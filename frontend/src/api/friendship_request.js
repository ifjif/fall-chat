import http from "@/request/http"

export const listFriendRequestApi = (data) => {
  return http.post("/friendRequest/list", data)
}

export const applyAddFriendApi = (data) => {
  return http.post("/friendRequest/apply", data)
}

export const handleFriendRequestApi = (data) => {
  return http.post('/friendRequest/handle', data)
}

