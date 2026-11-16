<template>
  <div :style="[avatarStyle, { position: 'relative' }]">
    <FlLoading :style="[avatarStyle, { 'border': '1px solid black' }]" :loading="loading"></FlLoading>
    <img v-if="!loading" :style="avatarStyle" :src="url" />
  </div>
</template>

<script setup>
import { computed, watch, ref } from 'vue'
import { fileAccessUrlApi } from '@/api/file'
import { FlLoading } from 'fall-ui'
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
const loading = ref(true)

watch(() => props.uri, nv => {
  if (nv) {
    fileAccessUrlApi({ file: nv })
      .then(v => url.value = v)
      .finally(() => loading.value = false)
  }
}, {
  immediate: true
})

</script>
