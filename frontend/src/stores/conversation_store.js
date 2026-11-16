import { convListApi, convSingleApi } from "@/api/conversation"
import { FindConvMembBriefJob } from "@/job/find_conv_memb_brief_job"
import { delayTaskMananger } from "@/utils/task_manager"
import { defineStore } from "pinia"
import { ref } from 'vue'
import { useUserStore } from "./user_info"
import { listImMessageApi, patchListImMessageApi } from "@/api/im_message"
import { getPreview } from "@/utils/im_message"
import { convMembCommitReadApi } from "@/api/conversation_member"

export const useConversationStore = defineStore('conversation', () => {
  // 
  // 当不处于活跃会话时，新消息到达,unread_count会增加, 显示在用户会话头像上
  // 当处于活跃会话是，新消息到达，unread_count,position.count会增加会,
  //  position.count在聊天框底部展示有多少条新消息，点击可到达底部
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
  //      position: {
  //        next: // 下次拉取的游标
  //        has_next: // 是否还能拉
  //        first: // 首次加载,
  //        top: // 滚动到顶部是
  //        cout: // 新消息,但未滚动到底部的数量
  //      }
  //
  //  --会话中消息--  最前的seq小，最后的seq大
  //      messages:[
  //        {
  //          id: 
  //          msg_type: 
  //          conv_type:
  //          content: 
  //          send_time: 
  //          sender_id: 
  //          seq: 
  //          status: 
  //          sender_name:
  //          sender_avatar:
  //          conv_id:
  //          peer_id;       // 发送消息时，失败后重发需要它。拉取消息不需要它
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

  function findNewSeq() {
    const conv = activeconversation.value
    let last_msg_seq = conv.last_msg_seq
    let next = conv.position.next
    let has_next = conv.position.has_next

    // 当前存在的消息中，最新的seq
    let latest_msg_seq = conv.messages?.[conv.messages.length - 1].seq || 0

    console.log(latest_msg_seq, next, has_next)

    if (last_msg_seq === latest_msg_seq) {
      return {
        hasNew: false,
        next,
        has_next
      }
    }
    const count = last_msg_seq - latest_msg_seq // 还未拉取的数量
    const first = latest_msg_seq === 0 || count > 40 // 拉最新

    if (first) {
      next = conv.position.next = ""
      has_next = conv.position.has_next = true
      conv.messages = []
    } else {
      last_msg_seq = last_msg_seq + 1
    }

    return {
      hasNew: true, // 是否有新消息
      next, // next_cursor
      has_next, // 是否还可以拉历史消息
      last_msg_seq, // 最新消息seq
      latest_msg_seq, // 已拉取的消息中最大seq
    }
  }

  // 拉历史消息，或最新消息
  // 如果都没有，不进行
  async function pullMessage(id, top = false) {
    const {
      hasNew,
      next,
      has_next,
      last_msg_seq,
      latest_msg_seq,
    } = findNewSeq()
    // 拉顶部, 但没有了
    if (top && !has_next) return

    const data = {
      conv_id: id,
      sort_key: "seq_desc",
      next,
      max_seq: last_msg_seq,
      min_seq: latest_msg_seq
    }

    console.log(data)

    const canCommitRead = next === "" || hasNew
    let result = null
    top = top || next === ""
    // 拉取会话消息
    if (top) {
      result = await listImMessageApi(data)
      if (!result) {
        return
      }
      activeconversation.value.position.next = result.next
      activeconversation.value.position.has_next = result.has_next
    } else if (hasNew) {
      // 补新消息
      const value = await patchListImMessageApi(data)
      result = {
        data: value
      }
      if (!result.data) {
        result.data = result.data.reverse()
      }
    }

    if (result?.data) {
      for (const msg of result.data) {
        addConversationMsg(id, msg, top)
      }
    }

    if (canCommitRead) {
      // 首次/patch，进行读提交
      readCommit()
    }
  }

  // 提交一个 会话成员的读
  async function readCommit() {
    const conv = getActiveConv()
    const msg = getLatestMessage()
    if (!msg || conv.last_read_seq === msg.seq) return
    const data = {
      conv_id: conv.id,
      msg_id: msg.id,
      msg_seq: msg.seq
    }

    try {
      await convMembCommitReadApi(data)
      const unreadCount = conv.unread_count
      conv.last_read_seq = msg.seq
      conv.unread_count = 0
      console.log(unreadCount)
    } catch {
    }
  }

  // 获取活跃会话中最新的消息
  function getLatestMessage() {
    const conv = getActiveConv()
    const len = conv.messages.length ?? 0
    return len > 0 ? conv.messages[len - 1] : null
  }


  // 从运程拉取用户所有会话
  async function pullConversation() {
    const result = await convListApi()
    if (result) {
      conversations.value = result
    }
  }

  // 设置活跃会话
  //
  //   首次，拉最新的, seq是当前往后的
  //   没有新消息，seq可以继续拉
  //
  //   有新消息，如果直接加入到本地，seq可继续使用
  //
  //   有新消息，没有直接加入到本地，需要新把新消息拉完，seq才可以继续使用
  //      这个要注意，如果太多，不能一次性拉
  //
  //
  async function setActiveConv(data) {
    // 不需要设置next,has_next , 由请求时进行判断
    if (!data.position) {
      data.position = {
        count: 0,
      }
    } else {
      data.position.count = 0
    }
    activeconversation.value = data
  }

  // 得到活跃会话
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
    if (index >= conversations.value.length) return
    if (index === 0) {
      setActiveConv(conversations.value[0])
      return
    }
    const [target] = conversations.value.splice(index, 1)
    conversations.value.unshift(target)
    setActiveConv(target)
  }

  // 移动临时会话到头
  function moveTempConvToFront(index) {
    if (index >= tempConversations.value.length) return
    if (index === 0) {
      setActiveConv(tempConversations.value[0])
      return
    }
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
  //    id: 
  //    msg_type: 
  //    content: 
  //    send_time: 
  //    sender_id: 
  //    status: 
  //    conv_type: 
  //    conv_id: 
  //    peer_id: 
  // }
  function addConversationMsg(id, msg, top) {
    if (top) {
      addConversationMsg2Head(id, msg)
    } else {
      addConversationMsg2Tail(id, msg)
    }
    dispatch('messageAppend')
  }

  // true 表示添加到head
  function addMsg2Head(messages, msg) {
    messages.unshift(msg)
    return true
  }
  function addMsg2Tail(messages, msg) {
    messages.push(msg)
    return false
  }

  function addConversationMsg2Tail(id, msg) {
    let result = addMsgToTempConv(id, msg, addMsg2Tail)
    if (result === -1) {
      addMsgToRealConv(id, msg, addMsg2Tail)
    }
  }

  function addConversationMsg2Head(id, msg) {
    let result = addMsgToTempConv(id, msg, addMsg2Head)
    if (result === -1) {
      addMsgToRealConv(id, msg, addMsg2Head)
    }
  }

  // 首次会话 自己发送消息
  function addMsgToTempConv(id, msg, addMsgFn) {
    const index = findConversationByTemp(id)
    if (index !== -1) {
      msg.sender_name = useUserStore().getNickname()
      msg.sender_avatar = useUserStore().getAvatar()

      if (!tempConversations.value[index].messages) tempConversations.value[index].messages = []
      addMsgFn(tempConversations.value[index].messages, msg)
    }

    return index
  }

  // 为正式会话添加消息
  // 只有活跃的会话才进行
  function addMsgToRealConv(id, msg, addMsgFn) {
    const index = findConversation(id)
    if (index !== -1) {
      if (!conversations.value[index].messages) conversations.value[index].messages = []
      let lastIndex = conversations.value[index].messages.length
      const head = addMsgFn(conversations.value[index].messages, msg)
      if (head) {
        lastIndex = 0
      }
      let conv = conversations.value[index]
      let mm = conv.messages[lastIndex]
      addMessageUserInfo(conv, mm)
    }

    return index
  }

  // job 完成回调
  function findConvMembBriefJobCallback(result, conv, messages) {
    if (!result) return
    conv.members = { ...conv.members, ...result }

    for (const message of messages) {
      // 回填，如果没有，再次找
      addMessageUserInfo(conv, message)
    }
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
        seq: data.last_msg_seq,
        status: 'success'
      })

      // 获取预览
      const preview = getPreview(targetMessage.msg_type, targetMessage.content)

      // 更新会话最后一条消息信息
      updateConvLastMsg({
        conv_id: data.conv_id,
        last_msg_id: data.last_msg_id,
        last_msg_preview: preview,
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
  //    seq
  //    success
  // }
  function updateMessageIdAndStatus(message, data) {
    message.id = data.last_msg_id
    message.send_time = data.last_msg_send_time
    message.status = data.status
    message.seq = data.seq
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
  // last_read_seq
  // unread_count
  //
  // isSender[true: 更新会话last_msg_seq]
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
    const msg = {
      id: data.last_msg_id,
      msg_type: data.msg_type,
      content: data.content,
      send_time: data.last_msg_send_time,
      sender_id: data.sender_id,
      seq: data.last_msg_seq,
      status: 'success',
    }

    const preview = getPreview(msg.msg_type, data.content)
    const updateConv = {
      conv_id: data.conv_id,
      last_msg_id: data.last_msg_id,
      last_msg_preview: preview,
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
      activeconversation.value.position.first = false
      activeconversation.value.position.top = false
      // 聊天框底部展示 新到达数量
      activeconversation.value.position.count += 1
      addMsgToRealConv(data.conv_id, msg, addMsg2Tail)
      // 会话列表中展示
      updateConvUnreadCount(data.conv_id, updateConv.last_msg_seq, false, index)
    }
  }

  //  处理错误
  //  data = {
  //    conv_id: 
  //    msg_id: 
  //    timestamp: 
  //  }
  function handleChatError(data) {
    let conv
    let index = findConversationByTemp(data.conv_id)
    if (index !== -1) {
      conv = tempConversations.value[index]
    } else {
      index = findConversation(data.conv_id)
    }

    if (!conv && index !== -1) {
      conv = conversations.value[index]
    }

    const msg_index = conv.messages.findIndex(m => m.id === data.msg_id)
    if (msg_index === -1) return
    const targetMessage = conv.messages[msg_index]

    updateMessageIdAndStatus(targetMessage, {
      last_msg_id: data.msg_id,
      last_msg_send_time: data.timestamp,
      status: 'error'
    })

  }


  // data{
  //  from
  //  to
  //  conv
  //  content,
  //  conv_type
  //  msg_type
  // }
  function generateMsg(data) {
    // 生成一个全局唯一的临时 ID
    const tempId = `temp_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`
    data.client_msg_id = tempId

    const timestamp = Date.now()
    const localMessage = {
      id: tempId,
      msg_type: data.msg_type,
      content: data.content,
      send_time: timestamp,
      sender_id: data.from,
      status: 'sending',
      conv_type: data.conv_type,
      conv_id: data.conv,
      peer_id: data.to
    }

    addConversationMsg(data.conv, localMessage)

    return {
      type: 'chat',
      timestamp,
      data
    }
  }

  return {
    conversations,
    tempConversations,
    on,
    todoConcat,
    pullConversation,
    setActiveConv,
    getActiveConv,
    getActiveConvType,
    updateConvLastMsg4ACK,
    addConversationMsg,
    generateMsg,
    handleReplyMsg,
    handleChatError,
    findNewSeq,
    pullMessage,
    readCommit,
  }
})
