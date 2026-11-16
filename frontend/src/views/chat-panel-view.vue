<template>
  <fl-container style="height:100%;">
    <fl-header height="38" style="background-color: pink;">
      <div style="height:100%;display: flex;align-items: center;padding: 0 10px;">
        {{ conv.name }}
      </div>
    </fl-header>
    <fl-main>
      <div style="height:100%;position:relative;">
        <fl-virtual-list :buffer="5" :estimatedHeight="120" v-show="convMessage.length > 0" :items="convMessage"
          @bottom="handleBottom" @top="handleTop" ref="scrollbarRef" @resize="handleResize">
          <template #default="{ item }">
            <div :style="{ width: `${messagePanelWidth}px` }">
              <MessageComponent :self="item.sender_id === userInfoStore.getId()" :data="item">
              </MessageComponent>
            </div>
          </template>
        </fl-virtual-list>
        <!--
        <fl-scrollbar @bottom="handleBottom" @top="handleTop" ref="scrollbarRef" @resize="handleResize">
          <div :style="{ width: `${messagePanelWidth}px` }">
            <MessageComponent :self="msg.sender_id === userInfoStore.getId()" :data="msg" v-for="msg in conv.messages">
            </MessageComponent>
          </div>
        </fl-scrollbar>
        -->

        <p style="text-align: center;width:100%;border-radius: 4px;bottom:0;position:absolute;font-size:14px;"
          v-if="conv.position && conv.position.count !== 0">
          <span @click="viewNewMessage"
            style="cursor:pointer;border-radius: 4px;padding:8px;background-color: rgba(0, 250, 0, 0.3);">
            有 {{
              conv.position.count
            }}
            条新消息
          </span>
        </p>
      </div>
    </fl-main>
    <fl-footer>
      <fl-textarea @blur="blurEvent" @focus="focusEvent" v-model="inputData" :rows="rows">
        <template #prefix>
          <div style="width:100%;">
            <ul style="display:flex;list-style: none;padding:0;margin: 0;gap: 10px;">
              <li style="cursor:pointer;">
                <FlPopover tag="span" :padding="0" :showArrow="false">
                  <template #reference>
                    <fl-icon :size="23">
                      <FlMeh></FlMeh>
                    </fl-icon>
                  </template>
                  <EmojiPanel @select="handleSelectEmoji"></EmojiPanel>
                </FlPopover>
              </li>
              <li>
                <Upload @change="handleChange">
                  <fl-icon style="cursor:pointer;" :size="23">
                    <FlFolder>
                    </FlFolder>
                  </fl-icon>
                </Upload>
              </li>
              <li style="cursor:pointer;">
                <fl-icon @click="handleAudio" :size="23">
                  <FlAudio></FlAudio>
                </fl-icon>
              </li>
              <li style="cursor:pointer;">
                <fl-icon @click="handleVideo" :size="23">
                  <FlVideo></FlVideo>
                </fl-icon>
              </li>
            </ul>
            <div v-if="fileList.length > 0">
              <fl-scrollbar wheel="x">
                <div style="display:flex;gap: 16px;">
                  <div style="height:50px;position:relative;" v-for="(f, index) in fileList">
                    <img v-if="[2, 3].includes(f.msgType)" style="height:100%;"
                      :style="{ 'aspect-ratio': `${f.width}/${f.height}` }" :src="f.localUrl">
                    <span v-else-if="f.msgType === 4">{{ f.name }}</span>
                    <audio v-else-if="f.msgType === 5" controls :src="f.localUrl"
                      :style="{ 'aspect-ratio': `${f.width}/${f.height}` }"></audio>
                    <fl-icon @click="deleteFile(index)" style="cursor:pointer;
                    position:absolute;top:0;right:-16px;
                    " :size="20">
                      <FlCloseCircle />
                    </fl-icon>
                  </div>
                </div>
              </fl-scrollbar>
            </div>
          </div>
        </template>
        <template #suffixRight>
          <fl-button type="primary" @click=handleSend>发送</fl-button>
        </template>
      </fl-textarea>
    </fl-footer>
  </fl-container>

  <AudioRecorder v-model="audioMask" @finish="finishAudioRecord"></AudioRecorder>
  <VideoRecorder v-model="videoMask" @finish="finishVideoRecord"></VideoRecorder>
</template>
<script setup>
import MessageComponent from '@/components/message-component.vue'
import { ref, computed, nextTick } from 'vue'
import { useWsStore } from '@/stores/ws_store'
import { useUserStore } from '@/stores/user_info'
import { useConversationStore } from '@/stores/conversation_store'
import { FlCloseCircle, FlMeh, FlFolder, FlAudio, FlVideo } from 'fall-ui/icons'
import Upload from '@/components/common/upload.vue'
import { getAudioMsgName, getAudioMsgType, getMediaMsgContent, getVideoMsgName, handleFileMsgUpload, handleImgMsgUpload, handleVideoMsgUpload, isImgFile, isVideoFile } from '@/utils/im_message'
import { createFileUsingBlob } from '@/utils/img'
import AudioRecorder from "@/components/audio-recorder.vue"
import VideoRecorder from '@/components/video-recorder.vue'
import { FlPopover } from 'fall-ui'
import EmojiPanel from '@/components/emoji-panel.vue'
import { decodeEmojiSequence } from '@/utils/emoji'

const scrollbarRef = ref(null)
const wsStore = useWsStore()
const userInfoStore = useUserStore()
const convStore = useConversationStore()
let sendCommit = false // 发送提交时不进行 readCommit

// 当前会话
const conv = computed(() => {
  nextTick(() => {
    scrollbarRef.value.scrollToBottom()
  })
  return convStore.getActiveConv()
})

const convMessage = computed(() => {
  return [...(conv.value.messages ?? [])]
})


const rows = ref(2)
const inputData = ref('')

const blurEvent = (e) => {
  rows.value = 2
}
const focusEvent = (e) => {
  rows.value = 5
}

const messagePanelWidth = ref(0)
// 首次进入，滚动到最下面
// 向上拉，滚动到之前的top
// 新消息加入，不进行滚动
const handleResize = ({ clientWidth, clientHeight, scrollWidth, scrollHeight }) => {
  messagePanelWidth.value = clientWidth
}

const handleSelectEmoji = (item) => {
  const emojiType = item.media_type
  if (emojiType === 'unicode') {
    const u = decodeEmojiSequence(item.font_name)
    inputData.value += u
    return
  }
  console.log(item)

  const c = {
    url: item.media_url,
    thumb_url: item.media_url,
    width: item.width,
    height: item.height,
    size: item.file_size,
    type: '',
    name: Date.now(),
    text: ''
  }

  const content = JSON.stringify(c)

  sendContent(content, 2)
}


const viewNewMessage = () => {
  scrollbarRef.value.scrollToBottom()
}

// 滚动到顶部拉历史消息
const handleTop = async () => {
  console.log('top event')
  const anchro = scrollbarRef.value.getAnchor()
  await convStore.pullMessage(conv.value.id, true)
  await nextTick()
  scrollbarRef.value.restoreAnchor(anchro)
}

// 滚动到底部 有新消息的 提交读
// 发消息，滚动到底部不处理
const handleBottom = () => {
  console.log('底部')
  if (sendCommit) {
    console.log('发送提交')
    conv.value.position.count = 0
    sendCommit = false
    return
  }
  if (conv.value.position.count) {
    conv.value.position.count = 0
    convStore.readCommit()
  }
}

// 仅支持单张图片/单条视频
const fileList = ref([])
const handleChange = async (e) => {
  const file = e.target.files[0]
  let fileObj = null
  if (isImgFile(file.type)) {
    fileObj = await handleImgMsgUpload(file)
  } else if (isVideoFile(file.type)) {
    fileObj = await handleVideoMsgUpload(file)
  } else {
    fileObj = handleFileMsgUpload(file)
  }

  console.log(fileObj)
  deleteFile(0)
  if (fileObj) {
    fileList.value.push(fileObj)
  }
}

const audioMask = ref(false)
// 音频处理
const handleAudio = () => {
  audioMask.value = true
}
// 音频录制结束
const finishAudioRecord = async ({ blob, duration }) => {
  const localUrl = URL.createObjectURL(blob)
  const name = getAudioMsgName()
  const f = {
    file: createFileUsingBlob(blob, name, blob.type),
    localUrl,
    revokeUrl: () => URL.revokeObjectURL(localUrl),
    name,
    type: blob.type,
    size: blob.size,
    duration,
    width: 300,
    height: 50,
    msgType: getAudioMsgType(),
  }
  console.log(f)

  deleteFile(0)
  fileList.value.push(f)
  audioMask.value = false
}

const videoMask = ref(false)
// 视频处理
const handleVideo = () => {
  videoMask.value = true
}

const finishVideoRecord = async ({ blob, duration }) => {
  const name = getVideoMsgName()
  const file = createFileUsingBlob(blob, name, blob.type)
  const f = await handleVideoMsgUpload(file)
  if (!Number.isFinite(f.duration)) {
    f.duration = duration
  }

  console.log(f)
  deleteFile(0)
  fileList.value.push(f)
  videoMask.value = false
}

const deleteFile = (index) => {
  if (fileList.value.length <= 0) return
  const f = fileList.value[index]
  if (!f) return
  if (f.revokeUrl) {
    f.revokeUrl()
  }
  fileList.value.splice(index, 1)
}


// 进行消息发送
const handleSend = async () => {
  let msgType = 1
  inputData.value = inputData.value.trim()
  let content = inputData.value
  if (!content && fileList.value.length === 0) {
    return
  }

  if (fileList.value.length > 0) {
    const f = fileList.value[0]
    const c = await getMediaMsgContent(f)

    if (c) {
      msgType = f.msgType
      c.text = content
      content = JSON.stringify(c)
    }
  }

  sendContent(content, msgType)
  deleteFile(0)
}

// 发送滚动到底部时不进行 readCommit
const sendContent = (content, msgType) => {
  const data = {
    from: userInfoStore.getId(),
    to: conv.value.peer_id,
    conv: conv.value.id,
    content,
    conv_type: convStore.getActiveConvType(),
    msg_type: msgType
  }
  console.log(data)
  sendCommit = true
  const msg = convStore.generateMsg(data)
  wsStore.sendMessage(msg)
  inputData.value = ''

  nextTick(() => {
    scrollbarRef.value.scrollToBottom()
  })
}
</script>
