import { shallowRef, ref } from "vue";

export function useMediaStream() {
  const localStream = shallowRef(null)
  const error = ref(null)
  const videoDeviceId = ref('')
  const audioDeviceId = ref('')

  // 找出所有可用设备
  async function listDevices() {
    const devices = await navigator.mediaDevices.enumerateDevices()
    return {
      videoInputs: devices.filter(d => d.kind === 'videoinput'),
      audioInputs: devices.filter(d => d.kind === 'audioinput'),
      audioOutputs: devices.filter(d => d.kind === 'audiooutput')
    }
  }

  // 针对单一 kind 取新track, 不动已有流
  async function acquireTrack(kind, deviceId) {
    const constraints = getDeviceConstraints(kind, deviceId)
    const tmp = await navigator.mediaDevices.getUserMedia(constraints)
    return kind === 'video' ? tmp.getVideoTracks()[0] : tmp.getAudioTracks()[0]
  }

  function getDeviceConstraints(kind, deviceId) {
    switch (kind) {
      case 'video':
        return {
          video: {
            deviceId: deviceId ? { exact: deviceId } : undefined
          }
        }

      case 'audio':
        return {
          audio: {
            deviceId: deviceId ? { exact: deviceId } : undefined
          }
        }
    }

    return {}
  }

  async function createMediaStream(constraints, kind) {
    switch (kind) {
      case 'user':
        return await navigator.mediaDevices.getUserMedia(constraints)
      case 'display':
        return await navigator.mediaDevices.getDisplayMedia(constraints)
    }

    throw new Error('不支持的media类型: ' + kind)
  }

  // kind = "user" | "display"
  async function start(video = true, audio = true, kind = "user") {
    const constraints = {
      video: video ? getDeviceConstraints('video', videoDeviceId.value)
        : false,
      audio: audio ? getDeviceConstraints('audio', audioDeviceId.value)
        : false
    }

    try {
      const stream = await createMediaStream(constraints, kind)
      localStream.value = stream
      error.value = null
    } catch (e) {
      error.value = e
      throw e
    }
  }

  // 切换设备，仅 kind="user" 支持
  // @return {removedTrack, addedTrack}
  // 交给上次进行 replaceTrack
  // 如果还没有开启，直接start
  function switchDevice(mediaKind, deviceKind, deviceId) {
    if (mediaKind !== 'user') {
      error.value = new Error('不支持的媒体种类设备切换:' + kind)
      throw error.value
    }
    if (!['video', 'audio'].includes(deviceKind)) {
      error.value = new Error('不支持的设备种类切换:' + kind)
      throw error.value
    }

    switch (deviceKind) {
      case 'video': return switchVideoDevice(mediaKind, deviceKind, deviceId,)
      case 'audio': return switchAudioDevice(deviceId, mediaKind)
    }
  }

  async function switchVideoDevice(mediaKind, deviceKind, deviceId) {
    videoDeviceId.value = deviceId
    const stream = localStream.value
    if (!stream || stream.getVideoTracks().length === 0) {
      await start(true, !!audioDeviceId.value, mediaKind)
      return null
    }

    const addedTrack = await acquireTrack(deviceKind, deviceId)
    const removedTrack = stream.getVideoTracks()[0]
    // 组装新 stream(新视频 + 原音频) 
    localStream.value = new MediaStream([addedTrack, ...stream.getAudioTracks()])
    if (removedTrack) removedTrack.stop()
    return { removedTrack, addedTrack }
  }

  async function switchAudioDevice(mediaKind, deviceKind, deviceId) {
    audioDeviceId.value = deviceId
    const stream = localStream.value
    if (!stream || stream.getAudioTracks().length === 0) {
      await start(!!videoDeviceId.value, true, mediaKind)
      return null
    }

    const addedTrack = await acquireTrack(deviceKind, deviceId)
    const removedTrack = stream.getAudioTracks()[0]
    localStream.value = new MediaStream([addedTrack, ...stream.getVideoTracks()])
    if (removedTrack) removedTrack.stop()

    return { removedTrack, addedTrack }
  }

  async function restartStream(video, audio, mediaKind) {
    stop()
    return start(video, audio, mediaKind)
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
    localStream.value.getVideoTracks().forEach(t => t.enabled = enabled)
  }
  function toggleAudio(enabled) {
    localStream.value.getAudioTracks().forEach(t => t.enabled = enabled)
  }

  return {
    localStream,
    error,
    start,
    stop,
    listDevices,
    restartStream,
    switchDevice,
    toggleVideo,
    toggleAudio,
    isVideoEnabled,
    isAudioEnabled
  }

}
