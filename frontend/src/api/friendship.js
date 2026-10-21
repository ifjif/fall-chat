import http from "@/request/http"

export const friendshipListApi = () => {
  return http.post("/friendship/list")
}
