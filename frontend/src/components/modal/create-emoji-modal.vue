<template>
  <fl-modal title="创建表情包" @ok="handleOk" v-model="modelValue">
    <template #content>
      <fl-form :model="formModel">
        <fl-form-item label="封面">
          <ImageUpload @change="imageChange" :size="80"></ImageUpload>
        </fl-form-item>
        <fl-form-item label="名字">
          <fl-input v-model="formModel.name"></fl-input>
        </fl-form-item>
        <fl-form-item label="收费">
          <fl-radio-group v-model="formModel.priceType" size="small">
            <fl-radio :value="1" label="免费"></fl-radio>
            <fl-radio :value="2" label="收费"></fl-radio>
          </fl-radio-group>
        </fl-form-item>
        <fl-form-item v-show="formModel.priceType === 2" label="价格">
          <fl-input v-model="formModel.price"></fl-input>
        </fl-form-item>
      </fl-form>
    </template>
  </fl-modal>
</template>
<script setup>
import { uploadFile } from '@/utils/file_upload'
import ImageUpload from '../common/image-upload.vue'
import { ref } from 'vue'
import { createEmojiPackApi } from '@/api/emoji'

const modelValue = defineModel({
  type: Boolean,
  default: false
})

const formModel = ref({
  coverUrl: '',
  name: '',
  priceType: 1,
  price: '0',
  source: 2
})

let coverFile = null
const imageChange = (file) => {
  coverFile = file
}

const uploadCover = async () => {
  if (!coverFile) {
    return ""
  }
  const p1 = new Promise((resolve, reject) => {
    uploadFile(coverFile, ({ fileUrl }) => {
      resolve(fileUrl)
    })
  })

  return await p1
}

const handleOk = async () => {
  formModel.value.coverUrl = await uploadCover()
  console.log(formModel.value)

  const data = {
    name: formModel.value.name,
    source: formModel.value.source,
    price_type: formModel.value.priceType,
    price: parseFloat(formModel.value.price),
    cover_url: formModel.value.coverUrl
  }

  await createEmojiPackApi(data)
}
</script>
