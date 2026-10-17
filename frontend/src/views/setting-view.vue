<template>
  <div>设置</div>
  <div>
    <video style="width:200px;height:200px;border: 1px solid black;" ref="localVideo" autoplay muted playsinline />
    <div v-for="peerId in remotePeers" :key="peerId">
      <video style="width:200px;height:200px;border: 1px solid black;" autoplay muted playsinline :ref="(el) => {
        if (el) {
          remoteVideos.set(peerId, el)
          const stream = peers.get(peerId)?.stream
          if (stream) {
            el.srcObject = stream
          }
        }
      }" />
    </div>

    <button @click="joinRoom(roomId, selfId)">加入房间</button>
    <fl-input v-model="targetId"></fl-input>
    <button @click="call(targetId)">呼叫对方</button>
    <button @click="leaveRoom">离开房间</button>

    <button @click="toggleAudio(false)">静音</button>
    <button @click="toggleVideo(false)">关闭摄像头</button>

    <button @click="listDevices().then(console.log)">枚举设备</button>

  </div>
</template>
<script setup>
import { useUserStore } from '@/stores/user_info'
import { useWebRtc } from '@/webRtc/use_web_rtc'
import { watch, ref, computed } from 'vue'

// lishi 7aeabd8c21cbdbca
// zs 7aeabdf9ad0bdbca
const targetId = ref('')

const selfId = computed(() => useUserStore().getId())
const roomId = ref('room-1')
const localVideo = ref(null)
const remoteVideos = ref(new Map())

const {
  localStream,
  peers,
  roomId: currentRoom,
  remotePeers,
  joinRoom,
  call,
  leaveRoom,
  toggleAudio,
  toggleVideo,
  switchAudioDevice,
  switchVideoDevice,
  listDevices
} = useWebRtc()

watch(() => localStream.value, () => {
  if (localVideo.value && localStream.value) {
    localVideo.value.srcObject = localStream.value
  }
})

function getRemoteVideoEl(peerId) {
  return remoteVideos.value.get(peerId)
}
</script>
