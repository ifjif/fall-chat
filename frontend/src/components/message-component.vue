<template>
  <div class="message-container">
    <div class="message-time">
      2026/01/01 01:01:01
    </div>
    <div :class="['message-body', { 'is-sender': self }]">
      <img src="../asserts/avatar1.png" style="width:50px;height:50px;border-radius: 50%;" />
      <div style="background-color: beige;">这是内容</div>
    </div>
  </div>
</template>
<script setup>
const props = defineProps({
  self: {
    type: Boolean,
    default: true
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
