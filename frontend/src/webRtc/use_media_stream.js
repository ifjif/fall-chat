import { ref, shallowRef } from 'vue'

export function useMediaStream() {
  const localStream = shallowRef(null)
  const error = ref(null)
  const videoDeviceId = ref('')
  const audioDeviceId = ref('')

  async function listDevices() {
    const devices = await navigator.mediaDevices.enumerateDevices()
    return {
      videoInputs: devices.filter(d => d.kind === 'videoinput'),
      audioInputs: devices.filter(d => d.kind === 'audioinput'),
      audioOuputs: devices.filter(d => d.kind === 'audiooutput')
    }
  }

  async function start(video = true, audio = true) {
    try {
      const constraints = {
        video: video
          ? {
            deviceId: videoDeviceId.value
              ? { exact: videoDeviceId.value }
              : undefined
          }
          : false,
        audio: audio
          ? {
            deviceId: audioDeviceId.value
              ? { exact: audioDeviceId.value }
              : undefined
          }
          : false
      }
      // 同一电脑无法同时获取
      // const stream = await navigator.mediaDevices.getUserMedia(constraints)

      // 手机浏览器进行屏幕共享限制较大
      const stream = await navigator.mediaDevices.getDisplayMedia(constraints)
      localStream.value = stream
      error.value = null
      return stream
    } catch (e) {
      error.value = e
      throw e
    }
  }

  async function switchVideoDevice(deviceId) {
    const wasEnabled = isVideoEnabled()
    videoDeviceId.value = deviceId
    await restartStream(wasEnabled, true)
  }

  async function switchAudioDevice(deviceId) {
    const wasEnabled = isAudioEnabled()
    audioDeviceId.value = deviceId
    await restartStream(true, wasEnabled)
  }

  async function restartStream(video = true, audio = true) {
    stop()
    return start(video, audio)
  }

  function stop() {
    localStream.value?.getTracks().forEach(t => t.stop())
    localStream.value = null
  }

  function isVideoEnabled() {
    return localStream.value?.getVideoTracks().some(t => t.enabled) ?? false
  }

  function isAudioEnabled() {
    return localStream.value?.getAudioTracks().some(t => t.enabled) ?? false
  }

  function toggleVideo(enabled) {
    localStream.value?.getVideoTracks().forEach(t => t.enabled = enabled)
  }

  function toggleAudio(enabled) {
    localStream.value?.getAudioTracks().forEach(t => t.enabled = enabled)
  }

  return {
    localStream,
    error,
    start,
    stop,
    restartStream,
    toggleVideo,
    toggleAudio,
    switchVideoDevice,
    switchAudioDevice,
    listDevices,
    isVideoEnabled,
    isAudioEnabled
  }
}
