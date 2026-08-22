<template>
  <Upload @change="onchange" accept="image/*">
    <div :style="styles">
      <fl-icon v-if="!imgUrl" :size="size" color="gray">
        <FlPlus />
      </fl-icon>
      <img style="width:100%;height:100%;" v-else :src="imgUrl" />
    </div>
  </Upload>
</template>
<script setup>
import { fileAccessUrlApi } from '@/api/file';
import Upload from './upload.vue';
import { FlPlus } from 'fall-ui/icons'
import { computed, ref, watch } from 'vue'

const props = defineProps({
  size: {
    type: [String, Number],
    default: 100
  },
  url: {
    type: String
  }
})

const emit = defineEmits(['change'])

const styles = computed(() => {
  const style = {
    border: `2px dashed gray`
  }
  style.width = `${props.size}px`
  style.height = `${props.size}px`

  return style
})

const imgUrl = ref('')
watch(() => props.url, async (nv) => {
  imgUrl.value = await fileAccessUrlApi({ file: nv })
})

let imageFile = null
const onchange = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  imageFile = file
  imgUrl.value = await createPreview(imageFile, 300)

  emit('change', imageFile)
}

const createPreview = async (file, maxSize = 300) => {
  // 直接在解码阶段就缩放到目标尺寸，避免解码原图全尺寸位图
  const bitmap = await createImageBitmap(file, {
    resizeWidth: maxSize,
    resizeHeight: maxSize,
    resizeQuality: 'high'
  })
  // 绘制到canvas 导出为blob
  const canvas = document.createElement('canvas')
  canvas.width = bitmap.width
  canvas.height = bitmap.height
  canvas.getContext('2d').drawImage(bitmap, 0, 0)
  // 用完立即释放，防止内存泄漏
  bitmap.close()

  return new Promise((resolve) => {
    canvas.toBlob((blob) => resolve(URL.createObjectURL(blob)), file.type || 'image/jpeg', 0.8)
  })

}
</script>
