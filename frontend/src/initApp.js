import { refreshApi } from "./api/auth"
import { FlMessage } from "fall-ui"
import { useUserStore } from "./stores/user_info"
import { useWsStore } from './stores/ws_store'
import router from "./router"
import { useConversationStore } from "./stores/conversation_store"

export const initApp = async () => {
  const userStore = useUserStore()
  try {
    await refreshApi()
    await userStore.loginSuccess()
    useWsStore()
    useConversationStore().pullConversation()
  } catch (err) {
    localStorage.removeItem('token')
    FlMessage({ type: 'warning', content: '登录已过期，请重新登录' })
    router.push('/login')
  }
}
