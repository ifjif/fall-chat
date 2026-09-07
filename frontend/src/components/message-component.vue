<template>
  <div class="message-container">
    <div class="message-time">
      {{ data.send_time }}
    </div>
    <div :class="['message-body', { 'is-sender': self }]">
      <img src="../asserts/avatar1.png" style="width:50px;height:50px;border-radius: 50%;" />
      <div style="position: relative;padding:10px;border-radius:4px;background-color: beige;">
        <fl-loading :loading="data.status === 'sending'"></fl-loading>
        {{ data.content }}
      </div>
      <fl-icon v-if="data.status === 'success'" size="20" color="green">
        <FlCheckCircle />
      </fl-icon>
      <fl-icon v-if="data.status === 'error'" size="20" color="red">
        <FlCloseCircle />
      </fl-icon>
    </div>
  </div>
</template>
<script setup>
import { FlCheckCircle, FlCloseCircle } from 'fall-ui/icons'
const props = defineProps({
  self: {
    type: Boolean,
    default: true
  },
  data: {
    type: Object
  }
})

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
</style>
