import http from "@/request/http";

export const friendshipListApi = () => {
  return http.post("/social/friendshipList")
}

export const applyAddFriendApi = (data) => {
  return http.post("/social/applyAddFriend", data)
}

export const searchFriendRequestApi = (data) => {
  return http.post("/social/searchFriendRequest", data)
}

export const handleFriendRequestApi = (data) => {
  return http.post('/social/handleFriendRequest', data)
}

export const searchGroupRequestApi = (data) => {
  return http.post("/social/searchGroupRequest", data)
}

export const applyAddGroupApi = (data) => {
  return http.post("/social/applyAddGroup", data)
}

export const handleGroupRequestApi = (data) => {
  return http.post("/social/handleGroupRequest", data)
}
