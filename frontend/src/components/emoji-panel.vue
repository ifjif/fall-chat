<template>
  <div style="display:flex;
  flex-direction: column;
  padding:4px;
  background-color:white;height:300px;width:100%;
  box-shadow: gray 0px 0px 10px 2px;
  ">
    <div style="flex:1;min-height:0;height:200px;">
      <FlScrollBar>
        <Grid style="width:400px;" :minItemWidth="50" cols="auto-fill">
          <div @click="selectEmoji(item)" v-for="item in packItems" :key="item.id"
            style="text-align: center;cursor:pointer;line-height: 60px;text-align: center;width:100%;height:100%;">
            <span style="font-size:30px;" v-if="item.media_type === 'unicode'">
              {{ decodeEmojiSequence(item.font_name) }}
            </span>

            <Image :style="{ 'aspect-ratio': `${item.width}/${item.height}`, 'max-width': '100%' }" v-if="item.media_type === 'dynamic' ||
              item.media_type === 'static'" :load="item.revoke_media_url" width="auto" height="60px"
              :uri="item.media_url" :remoteUrl="getRemoteFileUrl"></Image>
          </div>
        </Grid>
      </FlScrollBar>
    </div>
    <div style="border-top: 1px solid gray;min-height:20px;max-height:100px;">
      <FlScrollBar>
        <FlRow :gutter="20" v-for="items in packRows" style="padding:4px;width:400px;font-size:16px;">
          <FlCol style="cursor:pointer;" :style="[item.id === packId ? { 'background-color': 'skyblue' } : '']"
            @click="selectPack(item)" :span="packColSpan" v-for="item in items" :key="item.id">
            <TextEllipsis mode="scroll">
              <span>{{ item.name }}</span>
            </TextEllipsis>
          </FlCol>
        </FlRow>
      </FlScrollBar>
    </div>
  </div>
</template>
<script setup>
import { FlScrollBar, FlRow, FlCol } from "fall-ui";
import TextEllipsis from "./common/text-ellipsis.vue";
import Image from "./common/image.vue";
import { onMounted, ref, computed, watch } from 'vue'
import { listEmojiPackItemApi, listEmojiPackOwnApi } from "@/api/emoji";
import { decodeEmojiSequence } from "@/utils/emoji";
import { getRemoteFileUrl } from '@/utils/file';
import Grid from "./common/grid.vue";

const emit = defineEmits(['select'])

const emojiPacks = ref([])
const packItems = ref([])
const packColSpan = ref(4)
const packRows = computed(() => {
  const maxCol = 24 / packColSpan.value
  const maxRow = Math.ceil(emojiPacks.value.length / maxCol)
  const a = []
  for (let i = 0; i < maxRow; i++) {
    a.push(emojiPacks.value.slice(i * maxCol, i * maxCol + maxCol))
  }

  return a
})
const packId = ref('')
watch(() => packId.value, () => {
  packItems.value = []
  itemNext.value = ''
  getPackItem()
})
const next = ref('')
const itemNext = ref('')

const selectEmoji = (item) => {
  emit('select', item)
}

const selectPack = (item) => {
  packId.value = item.id
}

const getEmojiOwnPack = async () => {
  const data = {
    status: 2
  }

  const result = await listEmojiPackOwnApi(data)
  if (result) {
    next.value = result.next
  }

  const appendData = result.data ?? []
  emojiPacks.value.push(...appendData)
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
}

onMounted(() => {
  getEmojiOwnPack()
  if (packId.value) {
    getPackItem()
  }
})
</script>
