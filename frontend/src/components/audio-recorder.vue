<template>
  <fl-mask v-if="modelValue" @close="modelValue = false" style="display:flex;
    flex-direction: column;justify-content:center;">
    <div @click.stop style="display:flex;justify-content: center;">
      <div>
        <fl-button type="success" :disabled="started" @click="startRecordAudio">开始</fl-button>
        <fl-button v-if="started" type="warning" @click="handlePauseOrResume(state)">{{ stateStr }}</fl-button>
        <fl-button v-if="started" type="error" @click="stopRecordAudio">结束</fl-button>
      </div>
    </div>
  </fl-mask>
</template>

<script setup>
import { useRecordAudio } from '@/hooks/use_record_audio';
import { ref, defineModel, defineEmits, computed } from 'vue'

// 遮罩层关闭，不销毁当前组件
// 状态依旧
const modelValue = defineModel({
  type: Boolean,
  default: false
})

const emit = defineEmits(['finish'])

const started = ref(false)
const state = ref(0)
const states = ['暂停', '继续']
const stateStr = computed(() => states[state.value])

function finishRecord({ blob, duration }) {
  state.value = 0
  emit('finish', { blob, duration })
}

const {
  start,
  stop,
  pause,
  resume
} = useRecordAudio(finishRecord)

const startRecordAudio = async () => {
  start()
  started.value = true
}
const handlePauseOrResume = (s) => {
  switch (s) {
    case 0:
      state.value = 1
      pause()
      break
    case 1:
      state.value = 0
      resume()
  }
}
const stopRecordAudio = () => {
  stop()
  started.value = false
}
</script>
