export async function getVideoInfo(file) {
  return new Promise((resolve, reject) => {
    const time = 1
    const video = document.createElement("video")
    const url = URL.createObjectURL(file)
    video.src = url
    video.muted = true;
    video.playsInline = true;

    const info = {}

    video.onloadedmetadata = (e) => {
      info.duration = video.duration
      info.height = video.videoHeight
      info.width = video.videoWidth
      video.currentTime = time
    }

    video.onseeked = () => {
      const canvas = document.createElement("canvas")
      canvas.width = info.width
      canvas.height = info.height
      const ctx = canvas.getContext('2d')
      ctx.drawImage(video, 0, 0, canvas.width, canvas.height)

      canvas.toBlob((blob) => {
        URL.revokeObjectURL(url)
        info.coverType = 'image/jpeg'
        info.coverBlob = blob
        info.localUrl = canvas.toDataURL(info.coverType, 0.9)
        resolve(info)
      }, info.coverType, 0.9)
    }

    video.onerror = () => {
      URL.revokeObjectURL(url)
      reject(new Error('视频加载失败'))
    }
  })
}
