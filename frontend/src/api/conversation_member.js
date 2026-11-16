import http from "@/request/http";

export const convMembListBriefBatchApi = (data) => {
  return http.post("/conversationMember/listBriefBatch", data)
}

export const convMembCommitReadApi = (data) => {
  return http.post("/conversationMember/commit/read", data)
}
