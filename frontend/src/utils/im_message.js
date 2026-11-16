import { createFileUsingBlob, createImgPreview, getImgSize } from './img'
import { uploadFile } from './file_upload'
import { getVideoInfo } from './video'
// 消息类型 1:文本 2:图片 3:视频 4:文件 5:语音 6:语音通话 7:视频通话
export function getMsgTypeUsingFileType(type) {
  if (type.includes("image")) return 2
  if (type.includes("video")) return 3
  if (type.includes("audio")) return 5

  return 4
}

export function getAudioMsgName() {
  return `audio_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`
}
export function getVideoMsgName() {
  return `video_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`
}
export function getAudioMsgType() {
  return 5
}
export function getVideoMsgType() {
  return 3
}

export function isImgFile(type) {
  return type.includes("image")
}
export function isVideoFile(type) {
  return type.includes("video")
}


// 获取预览
// 消息类型 1:文本 2:图片 3:视频 4:文件 5:语音 6:语音通话 7:视频通话
export function getPreview(type, content) {
  const max = Math.min(content.length, 100)
  switch (type) {
    case 1:
      return content.substring(0, max)
    case 2:
      return "[图片]"
    case 3:
      return "[视频]"
    case 4:
      return "[文件]"
    case 5:
      return "[语音]"
    case 6:
      return "[语音通话]"
    case 7:
      return "[视频通话]"
    default:
      return "[未知]"
  }
}


// 图片消息
export const handleImgMsgUpload = async (file) => {
  const { width, height } = await getImgSize(file)
  const blob = await createImgPreview(file)
  const localUrl = URL.createObjectURL(blob)

  return {
    file,
    thumbFile: createFileUsingBlob(blob, `thumb_${file.name}`, file.type),
    localUrl,
    revokeUrl: () => URL.revokeObjectURL(localUrl),
    width,
    height,
    name: file.name,
    size: file.size,
    type: file.type,
    msgType: getMsgTypeUsingFileType(file.type)
  }
}

// 视频消息
export const handleVideoMsgUpload = async (file) => {
  const {
    duration,
    height,
    width,
    coverBlob,
    localUrl,
    coverType
  } = await getVideoInfo(file)

  return {
    file,
    coverFile: createFileUsingBlob(coverBlob, `cover_${file.name}`, coverType),
    duration,
    height,
    width,
    localUrl,
    name: file.name,
    size: file.size,
    type: file.type,
    msgType: getMsgTypeUsingFileType(file.type)
  }
}

// 文件消息
export const handleFileMsgUpload = (file) => {
  return {
    file,
    name: file.name,
    size: file.size,
    type: file.type,
    msgType: getMsgTypeUsingFileType(file.type)
  }
}

export async function getMediaMsgContent(f) {
  switch (f.msgType) {
    case 2:
      return await getImgFileContet(f)
    case 3:
      return await getVideoContent(f)
    case 4:
      return await getFileContent(f)
    case 5:
      return await getAudioContent(f)
  }
}

// 上传图片消息文件，返回带url的content
export const getImgFileContet = async (f) => {
  // 上传缩略图
  const p1 = new Promise((resolve, reject) => {
    uploadFile(f.thumbFile, ({ fileUrl }) => {
      resolve(fileUrl)
    })
  })

  // 上传原图
  const p2 = new Promise((resolve, reject) => {
    uploadFile(f.file, ({ fileUrl }) => {
      resolve(fileUrl)
    })
  })

  const [thumb_url, url] = await Promise.all([p1, p2])

  const content = {
    url,
    thumb_url,
    width: f.width,
    height: f.height,
    size: f.size,
    type: f.type,
    name: f.name,
    text: ''
  }

  return content
}

// 上传其他文件
export const getFileContent = async (f) => {
  const url = await new Promise((resolve, reject) => {
    uploadFile(f.file, ({ fileUrl }) => {
      resolve(fileUrl)
    })
  })

  return {
    url,
    name: f.name,
    size: f.size,
    type: f.type,
    text: ''
  }
}

// 上传视频消息文件，返回带url的内容
export const getVideoContent = async (f) => {
  const p1 = new Promise((resolve, reject) => {
    uploadFile(f.coverFile, ({ fileUrl }) => {
      resolve(fileUrl)
    })
  })

  const p2 = new Promise((resolve, reject) => {
    uploadFile(f.file, ({ fileUrl }) => {
      resolve(fileUrl)
    })
  })

  const [cover_url, url,] = await Promise.all([p1, p2])

  return {
    url,
    cover_url,
    width: f.width,
    height: f.height,
    size: f.size,
    type: f.type,
    name: f.name,
    duration: f.duration,
    text: ''
  }
}

// 上传 音频消息文件，返回带url的内容
export async function getAudioContent(f) {
  const url = await new Promise((resolve, reject) => {
    uploadFile(f.file, ({ fileUrl }) => {
      resolve(fileUrl)
    })
  })

  return {
    url,
    width: f.width,
    height: f.height,
    size: f.size,
    type: f.type,
    name: f.name,
    duration: f.duration,
    text: ''
  }
}
