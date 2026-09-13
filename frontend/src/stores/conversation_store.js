import { convSingleApi } from "@/api/conversation"
import { FindConvMembBriefJob } from "@/job/find_conv_memb_brief_job"
import { delayTaskMananger } from "@/utils/task_manager"
import { defineStore } from "pinia"
import { nextTick, ref, watch } from 'vue'
import { useUserStore } from "./user_info"

export const useConversationStore = defineStore('conversation', () => {
  // 
  // conv {
  //  --会话基本信息--
  //      id: "7a895f6ce1cbdbca",
  //      name: "我是李四",
  //      avatar: "files/80fea4a632bb159ef84be63a4ae25a2d/avatar55.webp",
  //      type: "peer",
  //      peer_id: "7aeabd8c21cbdbca",
  //      last_msg_id: 90799988350652416,
  //      last_msg_preview: "hello,hello,hello,hello,hello",
  //      last_msg_sender_id: 81476227113684992,
  //      last_msg_time: "2026-09-08T05:26:44.205Z",
  //      last_msg_seq: 2,
  //      last_read_seq: 2,
  //      unread_count: 0
  //
  //  --会话中消息--
  //      messages:[
  //        {
  //          id: 
  //          msg_type: 
  //          content: 
  //          send_time: 
  //          sender_id: 
  //          seq: 
  //          status: 
  //          sender_name:
  //          sender_avatar:
  //        }
  //      ]
  //
  //  --会话中成员信息--
  //      members: {
  //        id : {
  //          id: '',
  //          nickname:'', 
  //          in_group_nickname:''
  //          avatar:'',
  //        }
  //      }
  // 
  const conversations = ref([])
  const tempConversations = ref([])
  const activeconversation = ref({})
  const events = {}

  function on(eventName, callback) {
    events[eventName] = callback
  }

  function dispatch(eventName) {
    const callback = events[eventName]
    if (callback) {
      callback()
    }
  }

  // 设置活跃会话
  function setActiveConv(data) {
    activeconversation.value = data
  }

  // 得到获取会话
  function getActiveConv() {
    return activeconversation.value
  }

  // 判断指定会话是否是活跃会话
  function isActiveConv(id) {
    return activeconversation.value.id === id
  }

  // 获取活跃会话的类型
  function getActiveConvType() {
    if (activeconversation.value.type === 'peer') {
      return 0
    } else if (activeconversation.value.type === 'group') {
      return 1
    }
  }

  // 添加会话到临时会话
  function addTempConversation(data) {
    tempConversations.value.unshift(data)
    setActiveConv(data)
  }

  // 从正式会话中找会话
  // 可能是从别的页面来的
  function findConversation(id) {
    return conversations.value.findIndex(conv => conv.id === id || conv.peer_id === id)
  }

  // 从远程中找会话
  async function findConversationByRemote(id) {
    const conv = await convSingleApi(id)
    // 加入实时会话
    if (conv) {
      conversations.value.push(conv)
    }
  }

  // 从临时会话中找会话
  function findConversationByTemp(peer_id) {
    return tempConversations.value.findIndex(conv => conv.peer_id == peer_id)
  }

  // 移动正式会话到头
  function moveConvToFront(index) {
    if (index === 0 || index >= conversations.value.length) return
    const [target] = conversations.value.splice(index, 1)
    conversations.value.unshift(target)
    setActiveConv(target)
  }

  // 移动临时会话到头
  function moveTempConvToFront(index) {
    if (index === 0 || index >= tempConversations.value.length) return
    const [target] = tempConversations.value.splice(index, 1)
    tempConversations.value.unshift(target)
    setActiveConv(target)
  }

  // 从其它页面进入 会话页面
  //
  // data {
  //    id
  //    peer_id
  //    name
  //    avatar
  //    type
  // }
  //
  // 不会进行远程查询，不存在，建立临时会话
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


  // msg {
  //    id
  //    msg_type
  //    content 
  //    send_time 
  //    sender_id
  //    status 
  // }
  function addConversationMsg(id, msg) {
    let result = addMsgToTempConv(id, msg)
    if (result === -1) {
      addMsgToRealConv(id, msg)
    }
    dispatch('messageAppend')
  }

  // 首次会话 自己发送消息
  function addMsgToTempConv(id, msg) {
    msg.sender_name = useUserStore().getNickname()
    msg.sender_avatar = useUserStore().getAvatar()

    const index = findConversationByTemp(id)
    if (index !== -1) {
      if (!tempConversations.value[index].messages) tempConversations.value[index].messages = []
      tempConversations.value[index].messages.push(msg)
    }

    return index
  }

  // 为正式会话添加消息
  // 只有活跃的会话才进行
  function addMsgToRealConv(id, msg) {
    const index = findConversation(id)
    if (index !== -1) {
      if (!conversations.value[index].messages) conversations.value[index].messages = []
      let lastIndex = conversations.value[index].messages.length
      conversations.value[index].messages.push(msg)

      let conv = conversations.value[index]
      let mm = conv.messages[lastIndex]
      addMessageUserInfo(conv, mm)

    }

    return index
  }

  // job 完成回调
  function findConvMembBriefJobCallback(result, conv, messages) {
    if (!result) return
    for (const message of messages) {
      const senderID = message.sender_id
      const member = result[senderID]
      message.sender_name = member.in_group_nickname || member.nickname
      message.sender_avatar = member.avatar
    }

    conv.members = { ...conv.members, ...result }
  }

  // 为消息添加用户信息
  function addMessageUserInfo(conv, message) {
    let senderID = message.sender_id
    const member = conv.members?.[senderID]
    if (member) {
      message.sender_name = member.in_group_nickname || member.nickname
      message.sender_avatar = member.avatar
      return
    }

    let job = new FindConvMembBriefJob(conv, message, findConvMembBriefJobCallback)
    delayTaskMananger.submit(job)
  }


  // ACK 响应后更新消息
  function updateConvLastMsg4ACK(data) {
    let index = findConversationByTemp(data.to)
    if (index !== -1) {
      const [target] = tempConversations.value.splice(index)
      const isActive = isActiveConv(target.id)

      target.id = data.conv_id
      conversations.value.unshift(target)

      if (isActive) {
        setActiveConv(target)
      }
    }

    index = findConversation(data.conv_id)
    if (index !== -1) {
      const msg_index = conversations.value[index].messages.findIndex(m => m.id === data.client_msg_id)
      const targetMessage = conversations.value[index].messages[msg_index]

      // 更新消息状态
      updateMessageIdAndStatus(targetMessage, {
        last_msg_id: data.last_msg_id,
        last_msg_send_time: data.last_msg_send_time,
        status: 'success'
      })

      // 更新会话最后一条消息信息
      updateConvLastMsg({
        conv_id: data.conv_id,
        last_msg_id: data.last_msg_id,
        last_msg_preview: targetMessage.content,
        last_msg_time: data.last_msg_send_time,
        last_msg_sender_id: targetMessage.sender_id,
        last_msg_seq: data.last_msg_seq
      })

      // 更新未读数
      updateConvUnreadCount(data.conv_id, data.last_msg_seq, true, index)
    }

  }
  // 更新消息状态
  //
  // data {
  //    last_msg_id
  //    last_msg_send_time
  //    success
  // }
  function updateMessageIdAndStatus(message, data) {
    message.id = data.last_msg_id
    message.send_time = data.last_msg_send_time
    message.status = data.status
  }

  // 更新会话最后一条消息信息
  // 
  // data {
  //     conv_id
  //     last_msg_id
  //     last_msg_preview
  //     last_msg_sender_id
  //     last_msg_time
  //     last_msg_seq
  //  }
  function updateConvLastMsg(data) {
    const index = findConversation(data.conv_id)
    if (index !== -1) {
      const conv = conversations.value[index]
      // 更新最后一条消息
      conv.last_msg_id = data.last_msg_id
      conv.last_msg_preview = data.last_msg_preview
      conv.last_msg_sender_id = data.last_msg_sender_id
      conv.last_msg_time = data.last_msg_time
      conv.last_msg_seq = data.last_msg_seq
    }

    return index
  }

  // 更新未读数
  // 
  // last_read_seq
  // unread_count
  //
  // unread_count = (new_seq - last_read_seq)
  function updateConvUnreadCount(id, last_msg_seq, isSender, index) {
    if (index < 0) index = findConversation(id)
    let last_read_seq = conversations.value[index].last_read_seq
    if (isSender) {
      conversations.value[index].last_read_seq = last_msg_seq
      last_read_seq = last_msg_seq
    }
    let unreadCount = last_msg_seq - last_read_seq
    conversations.value[index].unread_count = unreadCount
  }

  // 消息回复处理
  async function handleReplyMsg(data) {
    // 先临时
    const msg = {
      id: data.last_msg_id,
      msg_type: 0,
      content: data.content,
      send_time: data.last_msg_send_time,
      sender_id: data.sender_id,
      seq: data.last_msg_seq,
      status: 'success',
    }

    const updateConv = {
      conv_id: data.conv_id,
      last_msg_id: data.last_msg_id,
      last_msg_preview: data.content,
      last_msg_sender_id: data.sender_id,
      last_msg_time: data.last_msg_send_time,
      last_msg_seq: data.last_msg_seq
    }

    // 找临时
    let index = findConversationByTemp(data.sender_id)
    if (index !== -1) {
      // 添加到实时，并设置id
      const [target] = tempConversations.value.splice(index, 1)
      target.id = data.conv_id
      conversations.value.push(target)
    }

    // 进行更新
    index = updateConvLastMsg(updateConv)

    // 远程拉会话并进行更新
    if (index === -1) {
      await findConversationByRemote(data.conv_id)
      index = updateConvLastMsg(updateConv)
    }

    // 更新未读数
    updateConvUnreadCount(data.conv_id, updateConv.last_msg_seq, false, index)

    // 判断当前会话是否是活跃的，是 追加消息
    if (isActiveConv(data.conv_id)) {
      addMsgToRealConv(data.conv_id, msg)
      // 活跃清零
      // todo 如果未在底部，可以寄一个临时的
      updateConvUnreadCount(data.conv_id, updateConv.last_msg_seq, true, index)
    }
  }

  return {
    conversations,
    tempConversations,
    on,
    todoConcat,
    setActiveConv,
    getActiveConv,
    getActiveConvType,
    updateConvLastMsg4ACK,
    addConversationMsg,
    handleReplyMsg,
  }
})
