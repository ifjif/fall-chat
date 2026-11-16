import { getBlobAudioDuration } from "@/utils/audio"
import { useMediaStream } from "@/webRtc/use_media_stream"
import { FlMessage } from "fall-ui"

// @param1:  ({ blob, duration }) => void
export function useRecordAudio(onFinish) {
  let recorder = null

  async function end({ blob, duration }) {
    try {
      duration = await getBlobAudioDuration(blob)
    } catch (e) {
      console.log('AudioContext解码失败，使用手动计时时长', e)
    }

    onFinish?.({ blob, duration })
  }

  const start = async (interval = 500) => {
    recorder = await audioRecorder(end)
    recorder.start(interval)
  }

  const stop = () => {
    recorder.stop()
  }
  const pause = () => {
    recorder.pause()
  }
  const resume = () => {
    recorder.resume()
  }

  return {
    start,
    stop,
    pause,
    resume
  }

}

export async function audioRecorder(onFinish) {
  const { start, localStream, stop } = useMediaStream()
  try {
    await start(false, true)
  } catch (e) {
    FlMessage({ type: 'error', bg: true, content: '请允许使用麦克风' })
  }
  if (!localStream.value) return

  const recorder = new MediaRecorder(localStream.value, {
    mimeType: MediaRecorder.isTypeSupported('audio/webm;codecs=opus')
      ? 'audio/webm;codecs=opus'
      : 'audio/webm'
  })


  const chunks = []
  let startTime = null

  recorder.ondataavailable = (e) => {
    if (e.data && e.data.size > 0) chunks.push(e.data)
  }

  recorder.onstart = () => {
    startTime = new Date()
  }

  recorder.onstop = () => {
    const duration = (new Date() - startTime) / 1000
    const blob = new Blob(chunks, { type: 'audio/webm' })
    stop()
    onFinish({ blob, duration })
  }

  // 暂停
  //recorder.pause()

  // 继续
  //recorder.resume()
  //
  return recorder

}
