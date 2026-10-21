import http from "@/request/http";

export const createGroupApi = (data) => {
  return http.post("/groupInfo/createGroup", data)
}

export const groupsApi = (data) => {
  return http.post("/groupInfo/groups", data)
}

