import http from "@/request/http";

export const convSingleApi = (conv_id) => {
  return http.get("/conversation/single", { conv_id })
}

export const convListApi = () => {
  return http.get("/conversation/list")
}
