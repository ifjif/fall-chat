<template>
  <template v-for="convs in [convStore.tempConversations, convStore.conversations]">
    <fl-menu @select="handleSelect" keyProp="id" :activeIndex="active" :isCollapse="collapse" :width="260"
      :data="convs">
      <template #icon="{ item }">
        <div style="position:relative;display:flex;justify-content: center;">
          <Avatar :uri="item.avatar"></Avatar>
          <Badge :max="999" :value="item.unread_count"></Badge>
        </div>
      </template>
      <template #title="{ item }">
        <div>
          <div style="display: flex;gap:2px;align-items: center;justify-content: space-between;">
            <fl-button v-if="item.type === 'group'" size="small" scale="0.7" type="success">
              群聊
            </fl-button>
            <span style="text-overflow: ellipsis;overflow:hidden;flex:1;font-size:14px;">
              {{ item.name }}
            </span>
            <span style="font-size: 10px;">
              {{ parseTime(item.last_msg_time) }}
            </span>
          </div>
          <div style="font-size: 12px;
                white-space: nowrap;
                overflow:hidden;
                text-overflow: ellipsis;
                ">{{ item.last_msg_preview }}</div>

        </div>
      </template>
    </fl-menu>
  </template>
</template>
<script setup>
import { computed, watch } from 'vue'
import { useConversationStore } from '@/stores/conversation_store'
import Avatar from '@/components/common/avatar.vue'
import Badge from '@/components/common/badge.vue'
import { parseTime } from '@/utils/timeinfo'
const props = defineProps({
  collapse: {
    type: Boolean,
    default: false
  }
})
const convStore = useConversationStore()
const active = computed(() => {
  const conv = convStore.getActiveConv()
  return conv.id ?? ""
})


watch(() => active.value, nv => {
  if (nv) {
    convStore.pullMessage(nv)
  }
}, {
  immediate: true
})


const handleSelect = (index, item) => {
  convStore.setActiveConv(item)
}
</script>
