import SparkMD5 from "spark-md5"

const CHUNK_SIZE = 5 * 1024 * 1024

self.onmessage = (e) => {
  const { file, taskId } = e.data
  const largeFile = isLargeFile(file)
  const filename = file.name
  splitChunks(file, async ({ id, chunks, total }) => {
    self.postMessage({ taskId, id, chunks, total, largeFile, filename })
  })
}

function isLargeFile(file) {
  return file.size > CHUNK_SIZE
}


function splitChunks(file, callback) {
  const total = Math.ceil(file.size / CHUNK_SIZE)
  const chunks = []
  const spark = new SparkMD5.ArrayBuffer()
  const reader = new FileReader()
  let currentIndex = 0

  const loadChunk = () => {
    const start = currentIndex * CHUNK_SIZE
    const end = start + CHUNK_SIZE
    const slice = file.slice(start, end)

    reader.readAsArrayBuffer(slice)
  }

  reader.onload = (e) => {
    const chunk = e.target.result
    spark.append(chunk)
    chunks.push({ index: currentIndex, chunk })
    currentIndex++

    if (currentIndex >= total) {
      const id = spark.end()
      callback({ id, chunks, total })
      return
    }

    loadChunk()
  }

  loadChunk()
}
