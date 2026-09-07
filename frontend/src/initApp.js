import { refreshApi } from "./api/auth"
import { FlMessage } from "fall-ui"
import { useUserStore } from "./stores/user_info"
import router from "./router"

export const initApp = async () => {
  const userStore = useUserStore()
  try {
    await refreshApi()
    userStore.loginSuccess()
  } catch (err) {
    localStorage.removeItem('token')
    FlMessage({ type: 'warning', content: '登录已过期，请重新登录' })
    router.push('/login')
  }
}
