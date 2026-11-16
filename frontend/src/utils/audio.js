export async function getBlobAudioDuration(blob) {
  const arrayBuffer = await blob.arrayBuffer()
  const audioContext = new AudioContext()
  const audioBuffer = await audioContext.decodeAudioData(arrayBuffer)
  return audioBuffer.duration
}

