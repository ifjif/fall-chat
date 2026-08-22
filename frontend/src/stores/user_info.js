import { fileAccessUrlApi } from "@/api/file";
import { userProfileApi } from "@/api/user";
import { defineStore } from "pinia";
import { ref, watch } from 'vue'

export const useUserStore = defineStore("user", () => {
  const userInfo = ref(null)
  const fullAvatarUrl = ref('')

  watch(() => userInfo.value?.avatar, async () => {
    fullAvatarUrl.value = await fetchUserAvatarUrl()
  })

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

  function updateInfo(data) {
    setInfo(data)
  }

  async function fetchUserAvatarUrl() {
    return await fileAccessUrlApi({ file: userInfo.value?.avatar })
  }

  return {
    userInfo,
    fullAvatarUrl,
    loginSuccess,
    updateInfo,
    fetchUserAvatarUrl
  }
})
