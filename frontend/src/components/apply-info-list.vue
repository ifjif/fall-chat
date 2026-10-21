<template>
  <div style="height:100%;display: flex;flex-direction: column;">
    <FlTabs v-model="activeTab" onlyHeader>
      <FlTabPane v-for="item in statusList" :key="item.name" :name="item.name" :label="item.label">
      </FlTabPane>
    </FlTabs>
    <div style="flex:1;min-height: 0;">
      <FlScrollBar @resize="handleResize">
        <FlMenu :data="menuData" :width="menuWidth">
          <template #icon="{ item }">
            <Avatar :uri="item.avatar"></Avatar>
          </template>
          <template #title="{ item }">
            <div>
              <div>{{ item.nickname }}</div>
              <ul style="color:gray;font-size:12px;list-style: none;">
                <li>
                  {{ item.message }}
                </li>
                <li>
                  来自：
                  <template v-if="item.apply_type === 'search'">搜索</template>
                </li>
                <li>
                  状态：
                  <template v-if="item.status === 'pending'">
                    <FlButton type="warning" scale="0.7" size="small">待处理</FlButton>
                  </template>
                  <template v-else-if="item.status === 'approved'">
                    <FlButton type="success" scale="0.7" size="small">已同意</FlButton>
                  </template>
                  <template v-else-if="item.status === 'rejected'">
                    <FlButton type="error" scale="0.7" size="small">已拒绝</FlButton>
                  </template>
                </li>
                <li>
                  时间：{{ parseTime(item.created_at) }}
                </li>
              </ul>
              <div v-if="!item.sender && item.status === 'pending'">
                <fl-button type="success" size="small" @click="handleApply(1, item)">同意</fl-button>
                <fl-button type="error" size="small" @click="handleApply(2, item)">拒绝</fl-button>
              </div>
            </div>
          </template>
        </FlMenu>
      </FlScrollBar>
    </div>
  </div>

</template>

<script setup>
import { listFriendRequestApi, handleFriendRequestApi } from '@/api/friendship_request'
import { FlMenu, FlScrollBar, FlButton, FlTabs, FlTabPane } from 'fall-ui';
import { onMounted, ref, watch } from 'vue'
import Avatar from './avatar.vue';
import { parseTime } from '@/utils/timeinfo';
const activeTab = ref('all')

const menuData = ref([])
const statusValue = {
  all: [],
  pending: [0],
  approved: [1],
  rejected: [2],
  expired: [3],
}
const statusList = ref([
  {
    name: 'all',
    label: '全部',
  },
  {
    name: 'pending',
    label: '待处理'
  },
  {
    name: 'approved',
    label: '已通过',
  },
  {
    name: 'rejected',
    label: '已拒绝',
  },
  {
    name: 'expired',
    label: '已过期',
  },
])

watch(() => activeTab.value, nv => {
  fetchApplyInfo(statusValue[nv])
})

const menuWidth = ref(0)
const handleResize = ({ clientWidth, scrollWidth, clientHeight, scrollHeight }) => {
  menuWidth.value = clientWidth
}

const fetchApplyInfo = async (status = []) => {
  const data = {
    filter_status: status
  }
  const result = await listFriendRequestApi(data)
  menuData.value = result
}

const handleApply = async (action, item) => {
  const data = { action, id: item.id }
  await handleFriendRequestApi(data)
  fetchApplyInfo()
}

onMounted(() => {
  fetchApplyInfo()
})
</script>
