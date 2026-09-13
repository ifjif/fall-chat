<template>
  <fl-container style="height:100%;">
    <fl-header height="38" style="background-color: pink;">
      <div style="height:100%;display: flex;align-items: center;padding: 0 10px;">
        {{ conv.name }}
      </div>
    </fl-header>
    <fl-main>
      <fl-scrollbar ref="scrollbarRef" @resize="handleResize">
        <div :style="{ width: `${messagePanelWidth}px` }">
          <MessageComponent :self="msg.sender_id === userInfoStore.getId()" :data="msg" v-for="msg in conv.messages">
          </MessageComponent>
        </div>
      </fl-scrollbar>
    </fl-main>
    <fl-footer>
      <fl-textarea v-model="inputData" :rows="6">
        <template #suffixRight>
          <fl-button type="primary" @click=handleSend>发送</fl-button>
        </template>
      </fl-textarea>
    </fl-footer>
  </fl-container>
</template>
<script setup>
import MessageComponent from '@/components/message-component.vue'
import { ref, computed, nextTick } from 'vue'
import { useWsStore } from '@/stores/ws_store'
import { useUserStore } from '@/stores/user_info'
import { useConversationStore } from '@/stores/conversation_store'

const scrollbarRef = ref(null)
const wsStore = useWsStore()
const userInfoStore = useUserStore()
const convStore = useConversationStore()

// 当前会话
const conv = computed(() => convStore.getActiveConv())
let messageContainerHeight = 0

const inputData = ref('')

const messagePanelWidth = ref(0)
const handleResize = ({ clientWidth, clientHeight, scrollWidth, scrollHeight }) => {
  messagePanelWidth.value = clientWidth
  if (messageContainerHeight < scrollHeight) {
    scrollbarRef.value.scrollToBottom()
  }
  messageContainerHeight = scrollHeight
}

const handleSend = () => {
  inputData.value = inputData.value.trim()
  const content = inputData.value
  if (!content) {
    return
  }
  const data = {
    from: userInfoStore.getId(),
    to: conv.value.peer_id,
    conv: conv.value.id,
    content,
    Type: convStore.getActiveConvType()
  }
  wsStore.sendMessage({ type: 'chat', timestamp: Date.now(), data })
  inputData.value = ''
}
</script>
