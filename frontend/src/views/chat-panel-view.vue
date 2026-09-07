<template>
  <fl-container style="height:100%;">
    <fl-header height="38" style="background-color: pink;">
      <div style="height:100%;display: flex;align-items: center;padding: 0 10px;">
        {{ conv.name }}
      </div>
    </fl-header>
    <fl-main>
      <fl-scrollbar @resize="handleResize">
        <div :style="{ width: `${messagePanelWidth}px` }">
          <!--
          <MessageComponent :self="false"></MessageComponent>
          <MessageComponent v-for="i in [1, 2, 3, 4, 5, 6, 7]"></MessageComponent>
          <MessageComponent :self="false"></MessageComponent>
          -->
          <MessageComponent :data="msg" v-for="msg in conv.messages"></MessageComponent>
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
import { ref, computed } from 'vue'
import { useWsStore } from '@/stores/ws_store'
import { useUserStore } from '@/stores/user_info'
import { useConversationStore } from '@/stores/conversation_store'

const wsStore = useWsStore()
const userInfoStore = useUserStore()
const convStore = useConversationStore()

// 当前会话
const conv = computed(() => convStore.getActiveConv())

const inputData = ref('')

const messagePanelWidth = ref(0)
const handleResize = ({ clientWidth }) => {
  messagePanelWidth.value = clientWidth
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
  console.log(data)
  wsStore.sendMessage({ type: 'chat', timestamp: Date.now(), data })
  inputData.value = ''
}
</script>
