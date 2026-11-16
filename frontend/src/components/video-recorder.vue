<template>
  <fl-mask v-if="modelValue" @close="modelValue = false" style="display: flex;flex-direction: column;
  justify-content: center;">
    <div @click.stop style="display:flex;flex-direction: column;align-items: center;">
      <h3>视频录制</h3>
      <p v-show="!started">
        <fl-button type="success" @click="switchType">切换为{{ switchTypeStr }}</fl-button>
      </p>
      <p>
        <video autoplay muted playsinline ref="videoRef"
          style="border: 1px solid black;width:400px;height:300px;"></video>
      </p>

      <div>
        <fl-button :disabled="started" type="success" @click="startRecord">开始</fl-button>
        <fl-button v-show="started" type="warning" @click="pauseOrResumeRecord(state)">{{ stateStr }}</fl-button>
        <fl-button v-show="started" type="error" @click="stopRecord">结束</fl-button>
      </div>
    </div>
  </fl-mask>
</template>
<script setup>
import { useRecordVideo } from '@/hooks/use_record_video';
import { ref, watch, defineEmits, computed } from 'vue'

const emit = defineEmits(['finish'])

const recordTypeStrs = {
  "user": '屏幕',
  "display": '个人'
}
const recordType = ref('user')
const switchTypeStr = computed(() => recordTypeStrs[recordType.value])
const videoRef = ref(null)
const started = ref(false)
const state = ref(0)
const states = ['暂停', '继续']
const stateStr = computed(() => {
  return states[state.value]
})

const modelValue = defineModel({
  type: Boolean,
  default: false
})

function switchType() {
  if (recordType.value === 'user') {
    recordType.value = 'display'
  } else {
    recordType.value = 'user'
  }
}

function finish({ blob, duration }) {
  started.value = false
  state.value = 0
  emit('finish', { blob, duration })
}

const {
  localStream,
  start,
  pause,
  resume,
  stop,
} = useRecordVideo(finish)

watch(() => localStream.value, nv => {
  if (nv) {
    videoRef.value.srcObject = nv
  }
})

const startRecord = () => {
  started.value = true
  start(500, recordType.value)
}

const pauseOrResumeRecord = (s) => {
  switch (s) {
    case 0:
      state.value = 1
      pause()
      return
    case 1:
      state.value = 0
      resume()
      return
  }
}
const stopRecord = () => {
  stop()
}

</script>
