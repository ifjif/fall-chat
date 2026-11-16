<template>
  <div>设置</div>
  <div>
    <video style="width:400px;height:200px;border: 1px solid black;" ref="localVideo" autoplay muted playsinline />
    <div v-for="peerId in remotePeers" :key="peerId">
      <video style="width:400px;height:200px;border: 1px solid black;" autoplay muted playsinline :ref="(el) => {
        if (el) {
          remoteVideos.set(peerId, el)
          const stream = peers.get(peerId)?.stream
          if (stream) {
            el.srcObject = stream
          }
        }
      }" />
    </div>

    <button @click="joinRoomHandler(roomId, selfId)">加入房间</button>
    <fl-input v-model="targetId"></fl-input>
    <button @click="call(targetId)">呼叫对方</button>
    <button @click="leaveRoom">离开房间</button>

    <button @click="toggleAudio(false)">静音</button>
    <button @click="toggleAudio(true)">开音</button>
    <button @click="toggleVideo(false)">关闭摄像头</button>
    <button @click="toggleVideo(true)">开启摄像头</button>

    <button @click="listDevicesHandle">枚举设备</button>

    <div>
      <label>
        摄像头
        <select v-model="videoDeviceId" @change="onVideoChange">
          <option v-for="d in videoInputs" :key="d.deviceId" :value="d.deviceId">
            {{ d.label || '未知摄像头' }}
          </option>
        </select>
      </label>

      <label>
        麦克风
        <select v-model="audioDeviceId" @change="onAudioChange">
          <option v-for="d in audioInputs" :key="d.deviceId" :value="d.deviceId">
            {{ d.label || '未知麦克风' }}
          </option>
        </select>
      </label>

      <label>
        扬声器
        <select v-model="audioOutputDeviceId" @change="onAudioOutputChange">
          <option v-for="d in audioOutputs" :key="d.deviceId" :value="d.deviceId">
            {{ d.label || '未知扬声器' }}
          </option>
        </select>
      </label>
    </div>

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
const videoInputs = ref([])
const audioInputs = ref([])
const audioOutputs = ref([])
const videoDeviceId = ref('')
const audioDeviceId = ref('')
const audioOutputDeviceId = ref('')

const joinRoomHandler = (roomId, userId) => {
  joinRoom({
    video: true,
    audio: false,
    mediaKind: 'display',
    roomId,
    userId
  })

  /*
  joinRoom(
    roomId,
    userId
  )
  */
}

const onVideoChange = async (e) => {
  console.log('video-change', e)
  switchVideoDevice(videoDeviceId.value)
}
const onAudioChange = (e) => {
  console.log('audio-change', e)
  switchAudioDevice(audioDeviceId.value)
}
const onAudioOutputChange = (e) => {
  console.log('audio-output-change', e)
  switchAudioOutputDevice(audioOutputDeviceId.value)
}

function switchAudioOutputDevice(deviceId) {
  if (!localVideo.value) return
  localVideo.value.setSinkId(deviceId)
}

async function listDevicesHandle() {
  const devices = await listDevices()
  console.log(devices)
  videoInputs.value = devices.videoInputs
  audioInputs.value = devices.audioInputs
  audioOutputs.value = devices.audioOutputs
}

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
