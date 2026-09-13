import { fileAccessUrlApi } from "@/api/file";
import { userProfileApi } from "@/api/user";
import { defineStore } from "pinia";
import { ref } from 'vue'

export const useUserStore = defineStore("user", () => {
  const users = ref([{ avatar: '', nickname: '' }])
  const userInfo = ref(null)

  async function loginSuccess() {
    const data = await userProfileApi()
    console.log('初始成功', data)
    setInfo(data)
  }

  function setInfo(data) {
    userInfo.value = {
      id: data.id,
      username: data.username,
      nickname: data.nickname,
      phone: data.phone,
      avatar: data.avatar,
      joinTime: data.created_at
    }
  }

  function getId() {
    return userInfo.value.id
  }

  function getAvatar() {
    return userInfo.value?.avatar
  }
  function getNickname() {
    return userInfo.value?.nickname
  }

  function updateInfo(data) {
    setInfo(data)
  }

  async function fetchUserAvatarUrl() {
    return await fileAccessUrlApi({ file: userInfo.value?.avatar })
  }

  return {
    userInfo,
    getId,
    getAvatar,
    getNickname,
    loginSuccess,
    updateInfo,
    fetchUserAvatarUrl
  }
})
