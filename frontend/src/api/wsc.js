import http from "@/request/http";

export const wsAccessUrlApi = () => {
  return http.get("/ws/accessUrl")
}
