import http from "@/request/http";

export const createEmojiPackApi = (data) => {
  return http.post("/emoji/createPack", data)
}

export const listEmojiPackItemApi = (data) => {
  return http.post("/emoji/listItem", data)
}

export const listEmojiPackApi = (data) => {
  return http.post("/emoji/listPack", data)
}

export const emojiPackDetail = (packId) => {
  return http.get(`/emoji/pack/${packId}`)
}
export const emojiPublish = (data) => {
  return http.post("/emoji/publish", data)
}
export const emojiPublishVerifyCount = (data) => {
  return http.post("/emoji/publish/verifyCount", data)
}

export const listEmojiPackOwnApi = (data) => {
  return http.post("/emoji/listPack/own", data)
}
