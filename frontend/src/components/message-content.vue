<template>
  <div @click="previewHandler">
    <div v-if="type === 1">
      {{ realContent }}
    </div>
    <div style="width:400px;" v-else-if="type === 2">
      <img :src="realUrl" style="width:100%;"
        :style="{ 'aspect-ratio': `${realContent.width}/${realContent.height}` }" />
      <Teleport to="body">
        <fl-mask @close="isPreview = false" v-if="isPreview">
          <img :src="realContentUrl" style=" height:80%;"
            :style="{ 'aspect-ratio': `${realContent.width}/${realContent.height}` }" />
        </fl-mask>
      </Teleport>
      <span> {{ realContent.text }}</span>

      <!--
    <p>url: {{ realContent.url }}</p>
    <p>thumb_url: {{ realContent.thumb_url }}</p>
    <p>size: {{ realContent.size }}</p>
    <p>width: {{ realContent.width }}</p>
    <p>height: {{ realContent.height }}</p>
    <p>name: {{ realContent.name }}</p>
    <p>type: {{ realContent.type }}</p>
    -->
    </div>
    <div style="width:400px;" v-else-if="type === 3">
      <div v-if="!isPreview" style="position: relative;">
        <img :src="realUrl" style="width:100%;"
          :style="{ 'aspect-ratio': `${realContent.width}/${realContent.height}` }" />
        <p style="position:absolute;width:100%;height:50px;top:50%;
      margin-top:-25px;
      text-align: center;">
          <fl-icon :size="50" style="border-radius: 50%; background-color: rgba(255, 255, 255, 0.7);">
            <FlPlayCircle />
          </fl-icon>
        </p>
        <p style="margin:0;display:flex;justify-content: right;position:absolute;width:100%;bottom:10px;">
          <span style="background-color: rgba(255, 255, 255, 0.7);padding:4px;">
            {{ realContent.duration }}
          </span>
        </p>
      </div>
      <video controls autoplay v-else-if="isPreview" :src="realContentUrl" style="width:100%;"
        :style="{ 'aspect-ratio': `${realContent.width}/${realContent.height}` }"></video>
      <span> {{ realContent.text }}</span>

      <!--
    <p>url: {{ realContent.url }}</p>
    <p>cover_url: {{ realContent.cover_url }}</p>
    <p>duration: {{ realContent.duration }}</p>
    <p>size: {{ realContent.size }}</p>
    <p>width: {{ realContent.width }}</p>
    <p>height: {{ realContent.height }}</p>
    <p>name: {{ realContent.name }}</p>
    <p>type: {{ realContent.type }}</p>
    <p>text: {{ realContent.text }}</p>
    -->
    </div>

    <div style="width:300px" v-else-if="type === 5">
      <audio style="width:100%;" controls :src="realUrl"
        :style="{ 'aspect-ratio': `${realContent.width}/${realContent.height}` }">
      </audio>
      <span>{{ realContent.text }}</span>

      <!--
    <p>url: {{ realContent.url }}</p>
    <p>name: {{ realContent.name }}</p>
    <p>size: {{ realContent.size }}</p>
    <p>duration: {{ realContent.duration }}</p>
    <p>type: {{ realContent.type }}</p>
    <p>width: {{ realContent.width }}</p>
    <p>height: {{ realContent.height }}</p>
    <p>text: {{ realContent.text }}</p>
    -->
    </div>
    <div v-else-if="type === 4">
      <p style="margin:0;">
        {{ realContent.name }}
        <fl-icon :size="24">
          <FlFile />
        </fl-icon>
      </p>
      <p style="font-size:14px;">大小:{{ realContent.size }}</p>
      <span>{{ realContent.text }}</span>
      <!--
    <p>url: {{ realContent.url }}</p>
    <p>name: {{ realContent.name }}</p>
    <p>type: {{ realContent.type }}</p>
    <p>size: {{ realContent.size }}</p>
    <p>text: {{ realContent.text }}</p>
    -->
    </div>
  </div>
</template>

<script setup>
import { fileAccessUrlApi } from '@/api/file';
import { computed, watch, ref } from 'vue';
import { FlPlayCircle, FlFile } from 'fall-ui/icons'

const props = defineProps({
  type: {
    type: Number,
    default: 1
  },
  content: {
    type: String,
    default: ""
  }
})

const realContent = computed(() => {
  if (props.type === 1) {
    return props.content
  }

  if ([2, 3, 4, 5].includes(props.type)) {
    return JSON.parse(props.content)
  }
})
const realUrl = ref('')
const realContentUrl = ref('')

watch(() => realContent.value, async nv => {
  let req_url = ''
  switch (props.type) {
    case 2:
      req_url = nv.thumb_url
      break
    case 3:
      req_url = nv.cover_url
      break
    case 5:
      req_url = nv.url
  }

  if (!req_url) return

  const url = await fileAccessUrlApi({ file: req_url })
  realUrl.value = url
}, {
  immediate: true
})

const isPreview = ref(false)

const previewHandler = async () => {
  if (![2, 3].includes(props.type)) return
  const reqUrl = realContent.value.url
  if (!reqUrl) return
  const url = await fileAccessUrlApi({ file: reqUrl })
  realContentUrl.value = url
  isPreview.value = true
}
</script>
