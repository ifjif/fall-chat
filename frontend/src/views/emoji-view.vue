<template>
  <div style="height:100%;display:flex;flex-direction: column;">
    <div style="display: flex;align-items: center;margin-bottom: 4px;background-color: white;">
      <fl-menu :activeIndex="activeMenuIdx" @select="handleSelect" mode="horizontal" :data="menuData">
      </fl-menu>
      <div style="min-width:60px;">
        <fl-button size="small" type="success" @click="createPackHandler">创建</fl-button>
      </div>
    </div>
    <div>
      <template v-if="activeMenuIdx === 'all'">
        <fl-radio-group v-model="priceType">
          <fl-radio :value="1" label="免费"></fl-radio>
          <fl-radio :value="2" label="付费"></fl-radio>
        </fl-radio-group>
      </template>
      <template v-if="activeMenuIdx === 'my'">
        <fl-radio-group v-model="sourceType">
          <fl-radio :value="2" label="我创建"></fl-radio>
          <fl-radio :value="3" label="已购买"></fl-radio>
        </fl-radio-group>
        <fl-radio-group v-model="packStatus">
          <fl-radio :value="0" label="已下架"></fl-radio>
          <fl-radio v-show="sourceType === 2" :value="1" label="待发布"></fl-radio>
          <fl-radio :value="2" label="已发布"></fl-radio>
          <fl-radio v-show="sourceType === 2" :value="3" label="待审核"></fl-radio>
        </fl-radio-group>
      </template>
    </div>
    <fl-scrollbar @resize="handleResize">
      <Grid :min-item-width="140" cols="auto-fill" :style="{ width: `${realWidth}px` }">
        <div @click="toEmojiDetail(item)" v-for="item in emojiPacks" :key="item.id"
          style="height:240px;border: 1px solid sandybrown;">
          <p> {{ item.name }}</p>
          <Image :height="100" :uri="item.cover_url" :remote-url="getCoverUrl"></Image>
          <p style="font-size:14px;">{{ parseTime(item.created_at) }}</p>
          <p style="font-size:14px;">
            <fl-button v-if="item.price_type === 'free'" size="small" type="success">免费</fl-button>
            <fl-button v-if="item.price_type === 'payment'" size="small" type="warning">付费</fl-button>
            <span style="margin-left: 4px;" v-if="item.price_type === 'payment'">{{ item.price }}元</span>
          </p>
        </div>
      </Grid>
    </fl-scrollbar>
  </div>

  <CreateEmojiModal v-model="showCreatePackModal"></CreateEmojiModal>

</template>
<script setup>
import { listEmojiPackApi, listEmojiPackOwnApi } from '@/api/emoji';
import { fileAccessUrlApi } from '@/api/file';
import Image from '@/components/common/image.vue';
import CreateEmojiModal from '@/components/modal/create-emoji-modal.vue';
import { parseTime } from '@/utils/timeinfo';
import { onMounted, ref, watch } from 'vue';
import router from '@/router';
import Grid from '@/components/common/grid.vue';

const menuData = ref([
  {
    title: '全部',
    index: 'all',
  },
  {
    title: '我的',
    index: 'my',
  },
])
const activeMenuIdx = ref('all')
const handleSelect = (idx, item) => {
  activeMenuIdx.value = idx
}

const emojiPacks = ref([])
const next = ref('')
const priceType = ref(0)
const sourceType = ref("")
const packStatus = ref("")

watch(() => [sourceType.value, packStatus.value], () => {
  next.value = ""
  emojiPacks.value = []
  getEmojiOwnPack()
})
watch(() => priceType.value, () => {
  next.value = ""
  emojiPacks.value = []
  getEmojiPack()
})
watch(() => activeMenuIdx.value, (nv) => {
  next.value = ""
  emojiPacks.value = []
  if (nv === 'my') getEmojiOwnPack()
  else if (nv === 'all') getEmojiPack()
})

const realWidth = ref(0)
const handleResize = ({ clientWidth }) => {
  realWidth.value = clientWidth
}

const toEmojiDetail = (item) => {
  console.log(item)
  router.push({
    path: '/emoji-detail',
    query: {
      id: item.id
    }
  })
}

const getEmojiOwnPack = async () => {
  const data = {
    source_type: sourceType.value || 0
  }
  if (packStatus.value !== "") data['status'] = packStatus.value

  const result = await listEmojiPackOwnApi(data)
  if (result) {
    next.value = result.next
  }

  const appendData = result.data ?? []
  emojiPacks.value.push(...appendData)
}

const getEmojiPack = async () => {
  const data = {
    price_type: priceType.value || 0,
    status: 2,
    next: next.value,
    sort_key: 'created_desc'
  }
  const result = await listEmojiPackApi(data)
  if (!result) {
    next.value = result.next
  }
  const appendData = result.data ?? []
  emojiPacks.value.push(...appendData)

  console.log(emojiPacks.value)
}

const showCreatePackModal = ref(false)
const createPackHandler = () => {
  showCreatePackModal.value = true
}

const getCoverUrl = async (uri) => {
  return await fileAccessUrlApi({ file: uri })
}


onMounted(() => {
  getEmojiPack()
})
</script>
