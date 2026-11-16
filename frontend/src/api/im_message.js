import http from "@/request/http";

export const listImMessageApi = (data) => {
  return http.post("/imMessage/list", data)
}

export const patchListImMessageApi = (data) => {
  return http.post("/imMessage/patchList", data)
}
