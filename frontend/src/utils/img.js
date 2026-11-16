// 等比例 得到thumbnail
export async function createImgPreview(file, maxSize = 300) {
  // 直接在解码阶段就缩放到目标尺寸，避免解码原图全尺寸位图
  /*
  const bitmap = await createImageBitmap(file, {
    resizeWidth: maxSize,
    resizeHeight: maxSize,
    resizeQuality: 'high'
  })
  */
  const bitmap = await createImageBitmap(file)
  const scale = Math.min(1, maxSize / Math.max(bitmap.width, bitmap.height))
  const w = Math.round(bitmap.width * scale)
  const h = Math.round(bitmap.height * scale)
  // 绘制到canvas 导出为blob
  const canvas = document.createElement('canvas')
  //canvas.width = bitmap.width
  //canvas.height = bitmap.height
  canvas.width = w
  canvas.height = h
  canvas.getContext('2d').drawImage(bitmap, 0, 0, w, h)
  // 用完立即释放，防止内存泄漏
  bitmap.close()

  return new Promise((resolve) => {
    canvas.toBlob((blob) => resolve(blob), file.type || 'image/jpeg', 0.8)
  })

}

// 创建 File Object 通过 Blob
export function createFileUsingBlob(blob, name, type) {
  return new File([blob], name, { type })
}

// 采用 createImageBitmap 得到图片的原始 width,height
export async function getImgSize(file) {
  return createImageBitmap(file).then(bitmap => {
    const size = { width: bitmap.width, height: bitmap.height }
    bitmap.close()
    return size
  })
}

// 采用 Image 得到图片的原始 width, height, 当 createImageBitmap不支持时
export async function getImgSizeFallback(file) {
  return new Promise((resolve, reject) => {
    const url = URL.createObjectURL(file)
    const img = new Image()
    img.onload = () => {
      resolve({ width: img.naturalWidth, height: img.naturalHeight })
      URL.revokeObjectURL(url)
    }
    img.onerror = reject
    img.src = url
  })
}
