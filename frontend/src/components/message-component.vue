<template>
  <div class="message-container">
    <div class="message-time">
      {{ parseTime(data.send_time) }}
    </div>
    <div :class="['message-body', { 'is-sender': self }]">
      <Avatar :uri="data.sender_avatar"></Avatar>

      <div style="max-width:80%;">
        <div :class="['sender-name', { 'is-sender-name': self }]">{{ data.sender_name }}</div>
        <div
          style="overflow-wrap:break-word;position: relative;padding:10px;border-radius:4px;background-color: beige;">
          <fl-loading :loading="data.status === 'sending'"></fl-loading>
          {{ data.content }}
        </div>
      </div>

      <div style="margin-top: 24px;">
        <fl-icon v-if="data.status === 'success'" size="20" color="green">
          <FlCheckCircle />
        </fl-icon>
        <fl-popover trigger="hover" :showArrow="false" :padding="4">
          <template #reference>
            <fl-icon @click="resend(data)" style="cursor:pointer;" v-if="data.status === 'error'" size="20" color="red">
              <FlCloseCircle />
            </fl-icon>
          </template>
          重发
        </fl-popover>
      </div>

    </div>
  </div>
</template>
<script setup>
import { FlCheckCircle, FlCloseCircle } from 'fall-ui/icons'
import Avatar from './avatar.vue';
import { useWsStore } from '@/stores/ws_store';
import { useConversationStore } from '@/stores/conversation_store';
import { parseTime } from '@/utils/timeinfo';
const wsStore = useWsStore()
const convStore = useConversationStore()
const props = defineProps({
  self: {
    type: Boolean,
    default: true
  },
  data: {
    type: Object
  }

})

const resend = (msg) => {
  const data = {
    from: msg.sender_id,
    to: msg.peer_id,
    conv: msg.conv_id,
    content: msg.content,
    conv_type: msg.conv_type,
    msg_type: msg.msg_type
  }
  wsStore.sendMessage(convStore.generateMsg(data))
}

</script>
<style scoped>
.message-container {
  padding: 10px;
}

.message-time {
  padding: 4px 0;
  text-align: center;
  font-size: 14px;
  color: gray;
}

.message-body {
  display: flex;
  gap: 10px;
}

.message-body.is-sender {
  flex-direction: row-reverse;
}

.sender-name {
  font-size: 14px;
  padding: 4px 0;
  color: gray;
}

.is-sender-name {
  text-align: right;
}
</style>
