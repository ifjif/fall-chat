import { defineStore } from "pinia";
import { ref } from "vue";
import { WebSocketService } from "@/service/websocket_service";
import { useConversationStore } from "./conversation_store";

// status [closed, connected, error]
export const useWsStore = defineStore('ws', () => {
  const messages = ref([])
  let status = 'closed'

  function init() {
    const ws = WebSocketService.getInstance()
    ws.on('*', msg => {
      if (msg.type === 'ping') {
        return
      }

      if (msg.type === 'chat_ack') {
        handleChatAck(msg)
        return
      }

      if (msg.type === 'chat_reply') {
        handleChatReply(msg)
        return
      }

      messages.value.push(msg)
    })

    status = 'connected'
  }

  function handleChatReply(msg) {
    console.log(msg)
  }

  function handleChatAck(msg) {
    const ackData = msg.data
    const { server_msg_id, client_msg_id } = ackData
    // 更新会话列表最后一条消息
    const data = {
      conv: ackData.conv,
      to: ackData.to,
      last_msg_id: server_msg_id,
      last_msg_send_time: msg.timestamp || Date.now(),
      client_msg_id
    }

    useConversationStore().updateConvLastMsg(data)
  }

  function sendMessage(msg) {
    if (msg.type === 'chat') {
      // 生成一个全局唯一的临时 ID
      const tempId = `temp_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
      msg.data.client_msg_id = tempId

      const localMessage = {
        id: tempId,
        msg_type: 0,
        content: msg.data.content,
        send_time: msg.timestamp,
        sender_id: msg.data.from,
        status: 'sending',
      }

      useConversationStore().addConversationMsg(msg.data.conv, localMessage)

    }
    WebSocketService.getInstance().send(msg)
  }
  function clearMessage() {
    messages.value = []
  }

  return {
    init,
    sendMessage,
    clearMessage
  }
})
