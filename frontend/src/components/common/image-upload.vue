<template>
  <Upload @change="onchange" accept="image/*">
    <div :style="styles">
      <FlIcon v-if="!imgUrl" :size="size" color="gray">
        <FlPlus />
      </FlIcon>
      <img ref="imgRef" style="width:100%;height:100%;" v-show="imgUrl" :src="imgUrl" />
    </div>
  </Upload>
</template>
<script setup>
import { fileAccessUrlApi } from '@/api/file';
import { createImgPreview, getImgSize } from '@/utils/img';
import Upload from './upload.vue';
import { FlPlus } from 'fall-ui/icons'
import { FlIcon } from 'fall-ui';
import { computed, ref, watch } from 'vue'

const props = defineProps({
  size: {
    type: [String, Number],
    default: 100
  },
  url: {
    type: String
  },
  showThumbnail: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['change'])
const imgRef = ref(null)

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

  const imageSize = await getImgSize(file)
  const blob = await createImgPreview(imageFile, 300)
  let showFile = blob
  if (!props.showThumbnail) {
    showFile = file
  }
  const url = URL.createObjectURL(showFile)
  imgRef.value.onload = () => {
    URL.revokeObjectURL(url)
    imgRef.value.onload = null
  }
  imgUrl.value = url

  emit('change', imageFile, imageSize)
}

</script>
