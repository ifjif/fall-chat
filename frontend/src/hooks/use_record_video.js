import { useMediaStream } from "@/webRtc/use_media_stream"
import { FlMessage } from "fall-ui"
import { shallowRef } from 'vue'

export function useRecordVideo(onFinish) {
  let recorder = null
  let localStream = shallowRef(null)

  function end({ blob, duration }) {
    if (onFinish) {
      onFinish({ blob, duration })
    }
  }

  async function start(interval = 500, type) {
    const { stream, recorder: r } = await recordVideo(end, type)
    localStream.value = stream
    recorder = r
    recorder.start(interval)
  }

  function pause() {
    recorder.pause()
  }

  function resume() {
    recorder.resume()
  }

  function stop() {
    recorder.stop()
  }

  return {
    start,
    pause,
    resume,
    stop,
    localStream
  }

}

async function recordVideo(onFinish, type) {
  const { localStream, stop, start } = useMediaStream()
  try {
    await start(true, true, type)
  } catch (e) {
    FlMessage({ type: 'error', bg: true, content: '请求打开摄像头' })
  }
  if (!localStream.value) return

  const recorder = new MediaRecorder(localStream.value, {
    mimeType: MediaRecorder.isTypeSupported("video/webm;codecs=opus")
      ? "video/webm;codecs=opus"
      : "video/webm"
  })

  const chunks = []
  let startTime = null

  recorder.ondataavailable = (e => {
    if (e.data && e.data.size > 0) chunks.push(e.data)
  })

  recorder.onstart = () => {
    startTime = new Date()
  }

  recorder.onstop = () => {
    const duration = (new Date() - startTime) / 1000
    const blob = new Blob(chunks, { type: "video/webm" })
    stop()

    onFinish({ blob, duration })
  }

  return {
    stream: localStream.value,
    recorder
  }
}
