<template>
  <div v-if="info">
    <p>作者:{{ info.owner_id }}</p>
    <p>名字:{{ info.name }}</p>
    <p>创建时间:{{ parseTime(info.created_at) }}</p>
    <p>状态:
      <span v-if="info.status === 'added'">已发布</span>
      <span v-if="info.status === 'draft'">草稿</span>
    </p>
  </div>
  <p style="color:red;">
    注意：页面离开后，数据会丢失，请及时发布<br>
    最多推送 50 个
  </p>
  <div v-if="info?.status === 'draft'">
    <fl-button @click="addImage" size="small" type="primary">添加图片</fl-button>
    <fl-button @click="addUnicode" size="small" type="success">添加字体表情</fl-button>
    <fl-button @click="publishHandler" size="small" type="warning">发布</fl-button>
  </div>
  <div>
    <Grid cols="auto-fill">
      <template v-for="item in packItems" :key="item.id">
        <span style="text-align: center;font-size:50px;" v-if="item.media_type === 'unicode'">
          {{ decodeEmojiSequence(item.font_name) }}
        </span>
        <Image :style="{ 'aspect-ratio': `${item.width}/${item.height}` }" v-if="item.media_type === 'dynamic' ||
          item.media_type === 'static'" :load="item.revoke_media_url" width="auto" height="80px" :uri="item.media_url"
          :remoteUrl="getRemoteFileUrl"></Image>
      </template>
    </Grid>
  </div>
</template>
<script setup>
import { emojiPackDetail, emojiPublish, emojiPublishVerifyCount, listEmojiPackItemApi } from '@/api/emoji';
import ImageUpload from '@/components/common/image-upload.vue';
import Image from '@/components/common/image.vue';
import { codePointsToHex, decodeEmojiSequence, detectEmojiInputType } from '@/utils/emoji';
import { getRemoteFileUrl } from '@/utils/file';
import { uploadFile } from '@/utils/file_upload';
import { parseTime } from '@/utils/timeinfo';
import { FlInput, FlMessageBox } from 'fall-ui';
import { h, computed, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import Grid from '@/components/common/grid.vue';
const route = useRoute()
const packId = computed(() => {
  return route.query.id
})

const info = ref(null)
const itemNext = ref("")
const packItems = ref([])

const unicodeValue = ref('')

const addImage = () => {
  let finalFile = null
  let fileWH = {}
  const fb = FlMessageBox({
    title: '添加图片',
    content: h(ImageUpload, {
      showThumbnail: false,
      onChange: (file, wh) => {
        finalFile = file
        fileWH = wh
      }
    }),
    onOk: () => {
      fb.close()
      generateImageItem(finalFile, fileWH)
    }
  })
}
const generateImageItem = (file, fileWH) => {
  let media_type_number = 1
  let media_type = 'static'
  if (file.type.includes('gif')) {
    media_type_number = 2
    media_type = 'dynamic'
  }

  let media_url = URL.createObjectURL(file)
  const item = {
    id: Date.now().toString(),
    media_type_number,
    media_type,
    font_name: '',
    media_file: file,
    media_url,
    revoke_media_url: () => URL.revokeObjectURL(media_url),
    thumb_url: "",
    width: fileWH.width,
    height: fileWH.height,
    file_size: file.size,
    sort_no: 0
  }

  packItems.value.push(item)
}

const addUnicode = () => {
  FlMessageBox({
    title: '添加unicode emoji',
    content: h(FlInput, {
      'onUpdate:modelValue': (v) => unicodeValue.value = v
    }),
    onOk: () => {
      const v = detectEmojiInputType(unicodeValue.value)
      if (v.type === 'unknown') return
      generateUnicodeItem(v)
    }
  })
}

const generateUnicodeItem = (info) => {
  const codepoints = info.value
  const item = {
    id: Date.now().toString(),
    media_type_number: 4,
    media_type: 'unicode',
    font_name: codePointsToHex(codepoints),
    media_url: "",
    thumb_url: "",
    width: 0,
    height: 0,
    file_size: 0,
    sort_no: 0
  }
  packItems.value.push(item)
}

const getPackDetail = async () => {
  info.value = await emojiPackDetail(packId.value)
  console.log(info.value)
}

const getPackItem = async () => {
  const data = {
    pack_id: packId.value,
    next: itemNext.value,
    sort_key: "sort_no_asc"
  }
  const result = await listEmojiPackItemApi(data)
  if (result) {
    itemNext.value = result.next
  }
  if (result.data) {
    packItems.value.push(...result.data);
  }

  console.log(packItems.value)
}

const publishHandler = async () => {
  const item_count = packItems.value.length ?? 0
  const pack_id = packId.value

  try {
    await emojiPublishVerifyCount({
      pack_id,
      item_count
    })
  } catch {
    return
  }

  const ps = []
  const items = []
  packItems.value.forEach((item, index) => {
    const p1 = new Promise(async (resolve, reject) => {
      await handleItem(item)
      const realItem = {
        media_type: item.media_type_number,
        font_name: item.font_name,
        media_url: item.media_remote_uri ?? "",
        thumb_url: item.thumb_remote_uri ?? "",
        width: item.width,
        height: item.height,
        file_size: item.file_size,
        sort_no: index
      }
      items.push(realItem)
      resolve()
    })
    ps.push(p1)
  })

  await Promise.all(ps)
  const data = {
    id: packId.value,
    items
  }
  console.log(data)
  await emojiPublish(data)
}

const handleItem = async (item) => {
  const ps = []
  const idxs = []
  if (item.media_url) {
    const p = new Promise((resolve, reject) => {
      uploadFile(item.media_file, ({ fileUrl }) => {
        resolve(fileUrl)
      })
    })
    ps.push(p)
    idxs.push(1)
  }

  if (item.thumb_url) {
    const p = new Promise((resolve, reject) => {
      uploadFile(item.thumb_file, ({ fileUrl }) => {
        resolve(fileUrl)
      })
    })
    ps.push(p)
    idxs.push(2)
  }

  const result = await Promise.all(ps)

  for (const idx of idxs) {
    switch (idx) {
      case 1:
        item.media_remote_uri = result[0]
        break
      case 2:
        item.thumb_remote_uri = result[1] ?? result[0]
    }
  }
}

onMounted(() => {
  getPackDetail()
  getPackItem()
})
</script>
