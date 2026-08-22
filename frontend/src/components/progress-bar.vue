<template>
  <div class="progressContainer">
    <div class="progressBody" :style="{ width: `${width}px`, borderColor: bdColor }">
      <div class="progressContent" :style="{ backgroundColor: color, width: `${progress}%` }">
      </div>
    </div>
    <span>{{ progress }}%</span>
  </div>
</template>

<script setup>

const props = defineProps({
  progress: {
    type: Number
  },
  width: {
    type: [Number, String],
    default: 200
  },
  color: {
    type: String,
    default: 'green'
  },
  bdColor: {
    type: String,
    default: 'gray'
  }
})
</script>

<style scoped>
.progressContainer {
  display: flex;
  align-items: center;
  gap: 4px;
}

.progressBody {
  height: 8px;
  border-width: 1px;
  border-style: solid;
  border-radius: 8px;
}

.progressContent {
  height: 100%;
  border-radius: 8px;
}
</style>
