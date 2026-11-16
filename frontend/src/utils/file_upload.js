import { fileChunkStatusApi, fileInstanceCheckApi, fileUploadApi, fileUploadChunkApi, fileUploadMergeApi } from "../api/file"
import { Scheduler } from "../scheduler/scheduler"

const worker = new Worker(
  new URL('../file_worker.js', import.meta.url),
  { type: 'module' }
)

const taskCallbacks = new Map()

// finish -> 完成时回调 ({fileUrl}) => void
// process -> 处理+上传 总进度 ({total}) => void  total为总chunks数，每次上传完一个chunk进行调用
export const uploadFile = (file, finish, process) => {
  const taskId = Date.now() + "_" + Math.random().toString(36).slice(2)
  taskCallbacks.set(taskId, { finish, process })

  worker.postMessage({ file, taskId })

  return {
    taskId,
    abort: () => {
      worker.postMessage({ type: 'abort', taskId })
      taskCallbacks.delete(taskId)
    }
  }
}


worker.onmessage = (e) => {
  const { taskId, id, chunks, total, largeFile, filename } = e.data
  const callbacks = taskCallbacks.get(taskId)
  if (!callbacks) return

  upload(id, chunks, filename, largeFile, callbacks)
    .finally(() => {
      taskCallbacks.delete(taskId)
    })
}

const upload = async (id, chunks, filename, largeFile, callbacks) => {
  // 秒传
  const fileUrl = await fileExists(id)
  if (fileUrl) {
    callbacks.finish?.({ fileUrl })
    return
  }

  // 查看已上传分片
  const beforeChunks = await fileChunkStatusApi({ file_id: id })
  const uploadChunks = filterChunks(chunks, beforeChunks)

  // 单文件或分片上传
  await doUpload(id, uploadChunks, largeFile, filename, callbacks)

  // 大文件合并
  if (largeFile) {
    const indices = chunks.map(data => data.index)
    const data = {
      file_id: id,
      filename: filename,
      total_chunk: chunks.length,
      chunks: indices
    }
    const fileUrl = await fileUploadMergeApi(data)
    callbacks.finish?.({ fileUrl })
  }
}

const doUpload = (id, chunks, largeFile, filename, callbacks) => {
  const total = chunks.length
  const scheduler = new Scheduler(3)

  let tasks = []
  if (largeFile) {
    tasks = chunks.map(data => scheduler.add(() => doUploadChunks(id, data, total, callbacks)))
  } else {
    tasks = chunks.map(data => scheduler.add(() => doUploadFile(id, data, filename, callbacks)))
  }

  return Promise.all(tasks)
}

const doUploadChunks = async (id, data, total, callbacks) => {
  const formData = new FormData()
  formData.append("file_id", id)
  formData.append("index", data.index)
  const file = new Blob([data.chunk])
  formData.append("chunk", file)

  try {
    await fileUploadChunkApi(formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    callbacks.process?.({ total })
  } catch (e) {
  }
}

const doUploadFile = async (id, data, filename, callbacks) => {
  const formData = new FormData()
  formData.append("file_id", id)
  const file = new File([data.chunk], filename)
  formData.append("file", file)

  try {
    const fileUrl = await fileUploadApi(formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    callbacks.finish?.({ fileUrl })
  } catch (e) { }
}


const filterChunks = (chunks, beforeChunks) => {
  return chunks.filter(data => !beforeChunks.some(bc => data.index == bc))
}

const fileExists = async (id) => {
  try {
    return await fileInstanceCheckApi({ file_id: id })
  } catch (e) {
  }
}
