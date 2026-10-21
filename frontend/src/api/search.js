import http from "@/request/http"

export const searchUserOrGroupApi = (data) => {
  return http.post("/search/userOrGroup", data)
}
