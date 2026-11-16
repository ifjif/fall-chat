<template>
  <div class="message-container">
    <div class="message-time">
      {{ parseTime(data.send_time) }}
    </div>
    <div :class="['message-body', { 'is-sender': self }]">
      <Avatar :uri="data.sender_avatar"></Avatar>

      <div style="max-width:80%;">
        <div :class="['sender-name', { 'is-sender-name': self }]">
          {{ data.sender_name }}
        </div>

        <div :class="[{ 'is-sender': self }]" style="display:flex;overflow-wrap:break-word;position: relative;">
          <div style="max-width:100%;width:max-content;padding:10px;border-radius:4px;"
            :style="{ backgroundColor: self ? 'beige' : '#ffc9c9' }">
            <fl-loading :loading="data.status === 'sending'"></fl-loading>
            <MessageContent :type="data.msg_type" :content="data.content"></MessageContent>
          </div>

          <div :style="[
            self ?
              { 'padding-right': '10px' }
              :
              { 'padding-left': '10px' }
          ]">
            <fl-icon v-if="data.status === 'success'" size="20" color="green">
              <FlCheckCircle />
            </fl-icon>
            <fl-popover trigger="hover" :showArrow="false" :padding="4">
              <template #reference>
                <fl-icon @click="resend(data)" style="cursor:pointer;" v-if="data.status === 'error'" size="20"
                  color="red">
                  <FlCloseCircle />
                </fl-icon>
              </template>
              重发
            </fl-popover>
          </div>

        </div>
      </div>

    </div>
  </div>
</template>
<script setup>
import { FlCheckCircle, FlCloseCircle } from 'fall-ui/icons'
import Avatar from './common/avatar.vue';
import { useWsStore } from '@/stores/ws_store';
import { useConversationStore } from '@/stores/conversation_store';
import { parseTime } from '@/utils/timeinfo';
import MessageContent from './message-content.vue';
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

.message-body .is-sender {
  flex-direction: row-reverse;
}

.sender-name {
  height: 20px;
  font-size: 14px;
  padding: 4px 0;
  color: gray;
}

.is-sender-name {
  text-align: right;
}
</style>
