<template>
  <div class="fallchat">
    <div class="body">
      <router-view></router-view>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useWsStore } from './stores/ws_store'
const wsStore = useWsStore()

onMounted(() => {
  wsStore.init()
})
</script>

<style scoped>
.fallchat {
  box-sizing: border-box;
  padding: 20px;
  width: 100%;
  height: 100%;
}

.body {
  height: 100%;
}
</style>
