<template>
  <template v-for="convs in [convStore.tempConversations, convStore.conversations]">
    <fl-menu @select="handleSelect" keyProp="id" :activeIndex="active" :isCollapse="false" :width="260" :data="convs">
      <template #icon="{ item }">
        <Avatar :uri="item.avatar"></Avatar>
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
              {{ item.last_msg_send_time }}
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
import { computed } from 'vue'
import { useConversationStore } from '@/stores/conversation_store'
import Avatar from '@/components/avatar.vue'
const convStore = useConversationStore()
const active = computed(() => convStore.getActiveConv().id)

const handleSelect = (index, item) => {
  convStore.setActiveConv(item)
}
</script>
