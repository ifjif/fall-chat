<template>
  <img :style="avatarStyle" :src="url" />
</template>

<script setup>
import { computed, watch, ref } from 'vue'
import { fileAccessUrlApi } from '@/api/file'
const props = defineProps({
  uri: {
    type: String
  },
  size: {
    type: [Number, String],
    default: 50
  },
  radius: {
    type: [Number, String],
    default: 50
  }
})

const avatarStyle = computed(() => {
  const styles = {
    width: `${props.size}px`,
    height: `${props.size}px`,
    borderRadius: `${props.radius}%`
  }

  return styles
})
const url = ref('')

watch(() => props.uri, nv => {
  fileAccessUrlApi({ file: nv })
    .then(v => url.value = v)
}, {
  immediate: true
})

</script>
