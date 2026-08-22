<template>
  <div class="fallchat">
    <div class="body">
      <router-view></router-view>
    </div>
  </div>
</template>

<script setup>
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
