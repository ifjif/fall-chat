import { defineStore } from "pinia";
import { ref } from "vue";
import { WebSocketService } from "@/service/websocket_service";
import { useConversationStore } from "./conversation_store";

// status [closed, connected, error]
export const useWsStore = defineStore('ws', () => {
  const ws = WebSocketService.getInstance()
  const messages = ref([])

  ws.on('*', msg => {
    if (msg.type === 'ping') {
      return
    }

    if (msg.type === 'invite') {
      console.log(msg)
    }
  })

  ws.on('chat_ack', msg => {
    handleChatAck(msg)
  })
  ws.on('chat_reply', msg => {
    handleChatReply(msg)
  })
  ws.on('chat_error', msg => {
    handleChatError(msg)
  })


  function handleChatAck(msg) {
    const ackData = msg.data
    const { server_msg_id, client_msg_id } = ackData
    // 更新会话列表最后一条消息
    const data = {
      conv_id: ackData.conv_id,
      to: ackData.to,
      last_msg_id: server_msg_id,
      last_msg_send_time: msg.timestamp || Date.now(),
      last_msg_seq: ackData.msg_seq,
      client_msg_id
    }

    useConversationStore().updateConvLastMsg4ACK(data)
  }

  function handleChatReply(msg) {
    const replyData = msg.data
    const data = {
      conv_id: replyData.conv_id,
      to: replyData.to,
      sender_id: replyData.sender_id,
      content: replyData.content,
      last_msg_id: replyData.msg_id,
      last_msg_send_time: msg.timestamp || Date.now(),
      last_msg_seq: replyData.msg_seq,
      msg_type: replyData.msg_type
    }
    useConversationStore().handleReplyMsg(data)
  }

  function handleChatError(msg) {
    const replyData = msg.data

    const data = {
      conv_id: replyData.data.conv_id,
      msg_id: replyData.data.msg_id,
      timestamp: msg.timestamp
    }

    useConversationStore().handleChatError(data)
  }


  function sendMessage(msg) {
    WebSocketService.getInstance().send(msg)
  }

  function clearMessage() {
    messages.value = []
  }

  function on(type, handler) {
    ws.on(type, handler)
  }

  function onOnce(type, handler) {
    ws.onOnce(type, handler)
  }

  return {
    sendMessage,
    clearMessage,
    on,
    onOnce
  }
})
