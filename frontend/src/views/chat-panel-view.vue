<template>
  <fl-container style="height:100%;">
    <fl-header height="38" style="background-color: pink;">
      <div style="height:100%;display: flex;align-items: center;padding: 0 10px;">
        张三
      </div>
    </fl-header>
    <fl-main>
      <fl-scrollbar @resize="handleResize">
        <div :style="{ width: `${messagePanelWidth}px` }">
          <MessageComponent :self="false"></MessageComponent>
          <MessageComponent v-for="i in [1, 2, 3, 4, 5, 6, 7]"></MessageComponent>
          <MessageComponent :self="false"></MessageComponent>
        </div>
      </fl-scrollbar>
    </fl-main>
    <fl-footer>
      <fl-textarea :rows="6">
        <template #suffixRight>
          <fl-button type="primary">发送</fl-button>
        </template>
      </fl-textarea>
    </fl-footer>
  </fl-container>
</template>
<script setup>
import MessageComponent from '@/components/message-component.vue'
import { ref } from 'vue'

const messagePanelWidth = ref(0)
const handleResize = ({ clientWidth }) => {
  messagePanelWidth.value = clientWidth

}
</script>
