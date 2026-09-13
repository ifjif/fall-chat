import http from "@/request/http";

export const convMembListBriefBatchApi = (data) => {
  return http.post("/conversationMember/listBriefBatch", data)
}
