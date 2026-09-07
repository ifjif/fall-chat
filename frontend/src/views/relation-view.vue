<template>
  <fl-row>
    <fl-col :span="18">
      <fl-container>
        <fl-header>
          <p style="text-align:center;color:gray;font-size:14px;">好友</p>
        </fl-header>
        <fl-main>
          <fl-menu @select="handleSelect" :data="friends" width="100%">
            <template #icon="{ item }">
              <Avatar :uri="item.avatar"></Avatar>
            </template>
            <template #title="{ item }">
              <div>
                <div>
                  <template v-if="item.remark">
                    {{ item.remark }}
                  </template>
                  <template v-else>
                    {{ item.nickname }}
                  </template>
                </div>
                <div style="font-size:12px;color:gray;">来自：
                  <template v-if="item.source === 'search'">搜索</template>
                </div>
              </div>
            </template>
          </fl-menu>
        </fl-main>
      </fl-container>
      <fl-container>
        <fl-header>
          <p style="text-align:center;color:gray;font-size:14px;">群</p>
        </fl-header>
        <fl-main>
          <fl-menu @select="selectGroup" :data="groups" width="100%">
            <template #icon="{ item }">
              <Avatar :uri="item.avatar"></Avatar>
            </template>
            <template #title="{ item }">
              <div style="padding-left: 20px;">
                <div>
                  {{ item.name }}
                </div>
                <div style="font-size:12px;color:gray;">描述：
                  {{ item.description }}
                </div>
              </div>
            </template>
          </fl-menu>
        </fl-main>
      </fl-container>
    </fl-col>

    <fl-col :span="6">
      <div v-if="selectFriend">
        <h3>好友详情</h3>
        <p>昵称：{{ selectFriend.nickname ?? selectFriend.name }}</p>
        <fl-button type="success" @click="handleConcat">发送消息</fl-button>
      </div>
    </fl-col>
  </fl-row>
</template>

<script setup>
import { groupsApi } from '@/api/group'
import { friendshipListApi } from '@/api/social'
import router from '@/router'
import { onMounted, ref } from 'vue'
import { useConversationStore } from '@/stores/conversation_store'
import Avatar from '@/components/avatar.vue'

const convStore = useConversationStore()

const friends = ref([])
const groups = ref([])

const fetchFriends = async () => {
  const data = await friendshipListApi()
  friends.value = data
}

const fetchGrops = async () => {
  const data = {
    filter: []
  }
  const result = await groupsApi(data)
  groups.value = result ?? []
}

const selectFriend = ref(null)
const selectType = ref('')
const handleSelect = (index, item) => {
  selectFriend.value = item
  selectType.value = 'peer'
}
const selectGroup = (index, item) => {
  selectFriend.value = item
  selectType.value = 'group'
}

const handleConcat = () => {
  const data = {
    id: selectFriend.value.id,
    peer_id: selectFriend.value.id,
    name: selectFriend.value.nickname ?? selectFriend.value.name,
    avatar: selectFriend.value.avatar,
    type: selectType.value
  }
  convStore.todoConcat(data)
  router.push('/chat')
}

onMounted(() => {
  fetchFriends()
  fetchGrops()
})
</script>
