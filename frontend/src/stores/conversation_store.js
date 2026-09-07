import { defineStore } from "pinia"
import { ref } from 'vue'

export const useConversationStore = defineStore('conversation', () => {
  const conversations = ref([])
  const tempConversations = ref([])
  const activeconversation = ref({})

  function setActiveConv(data) {
    activeconversation.value = data
  }
  function getActiveConv() {
    return activeconversation.value
  }
  function isActiveConv(id) {
    return activeconversation.value.id === id
  }
  function getActiveConvType() {
    if (activeconversation.value.type === 'peer') {
      return 0
    } else if (activeconversation.value.type === 'group') {
      return 1
    }
  }
  function addTempConversation(data) {
    tempConversations.value.unshift(data)
    setActiveConv(data)
  }

  function findConversation(id) {
    return conversations.value.findIndex(conv => conv.id === id)
  }

  function findConversationByTemp(peer_id) {
    return tempConversations.value.findIndex(conv => conv.peer_id == peer_id)
  }

  function moveConvToFront(index) {
    if (index === 0 || index >= conversations.value.length) return
    const [target] = conversations.value.splice(index, 1)
    conversations.value.unshift(target)
    setActiveConv(target)
  }
  function moveTempConvToFront(index) {
    if (index === 0 || index >= tempConversations.value.length) return
    const [target] = tempConversations.value.splice(index, 1)
    tempConversations.value.unshift(target)
    setActiveConv(target)
  }

  function todoConcat(data) {
    let index = findConversation(data.id)
    if (index !== -1) {
      moveConvToFront(index)
      return
    }
    index = findConversationByTemp(data.id)
    if (index !== -1) {
      moveTempConvToFront(index)
      return
    }
    addTempConversation(data)
  }

  function addConversationMsg(id, msg) {
    let index = findConversationByTemp(id)
    if (index !== -1) {
      if (!tempConversations.value[index].messages) tempConversations.value[index].messages = []
      tempConversations.value[index].messages.push(msg)
      return
    }
    index = findConversation(id)
    if (index !== -1) {
      if (!conversations.value[index].messages) conversations.value[index].messages = []
      conversations.value[index].messages.push(msg)
    }
  }


  function updateConvLastMsg(data) {
    let index = findConversationByTemp(data.to)
    if (index !== -1) {
      const [target] = tempConversations.value.splice(index)
      const isActive = isActiveConv(target.id)
      // 更新消息状态
      const msg_index = target.messages.findIndex(m => m.id === data.client_msg_id)
      target.messages[msg_index].id = data.last_msg_id
      target.messages[msg_index].send_time = data.last_msg_send_time
      target.messages[msg_index].status = 'success'

      const last_msg_preview = target.messages[msg_index].content
      const last_msg_sender_id = target.messages[msg_index].sender_id

      // 更新最后一条消息
      target.id = data.conv
      target.last_msg_id = data.last_msg_id
      target.last_msg_preview = last_msg_preview
      target.last_msg_sender_id = last_msg_sender_id
      target.last_msg_send_time = data.last_msg_send_time
      conversations.value.unshift(target)
      if (isActive) {
        setActiveConv(target)
      }
      return
    }
    index = findConversation(data.conv)
    if (index !== -1) {
      // 更新消息状态
      const msg_index = conversations.value[index].messages.findIndex(m => m.id === data.client_msg_id)
      conversations.value[index].messages[msg_index].id = data.last_msg_id
      conversations.value[index].messages[msg_index].send_time = data.last_msg_send_time
      conversations.value[index].messages[msg_index].status = 'success'

      const last_msg_preview = conversations.value[index].messages[msg_index].content
      const last_msg_sender_id = conversations.value[index].messages[msg_index].sender_id

      // 更新最后一条消息
      conversations.value[index].last_msg_id = data.last_msg_id
      conversations.value[index].last_msg_preview = last_msg_preview
      conversations.value[index].last_msg_sender_id = last_msg_sender_id
      conversations.value[index].last_msg_send_time = data.last_msg_send_time
    }
  }

  return {
    conversations,
    tempConversations,
    addTempConversation,
    todoConcat,
    setActiveConv,
    getActiveConv,
    getActiveConvType,
    updateConvLastMsg,
    addConversationMsg,
  }
})
