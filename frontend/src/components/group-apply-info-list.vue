<template>
  <div style="height:100%;display: flex;flex-direction: column;">
    <FlTabs v-model="activeStatus" onlyHeader>
      <FlTabPane v-for="item in statusList" :label="item.label" :name="item.name" :key="item.name"></FlTabPane>
    </FlTabs>
    <div style="flex:1;min-height: 0;">
      <FlScrollBar @resize="handleResize">
        <FlMenu :data="groupApplyData" :width="menuWidth">
          <template #icon="{ item }">
            <Avatar :uri="item.from_user_avatar"></Avatar>
          </template>
          <template #title="{ item }">
            <div>
              <div>{{ item.from_user_name }}</div>
              <ul style="font-size: 12px;color: gray;list-style: none;">
                <li>
                  {{ item.message }}
                </li>
                <li>
                  群：{{ item.group_name }}
                </li>
                <li v-if="item.invite_user_id">
                  邀请人：{{ item.invite_user_name }}
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
                <li v-if="item.handle_by_user_id">
                  处理人：{{ item.handle_by_user_name }}
                </li>
                <li>
                  时间：
                  {{ parseTime(item.created_at) }}
                </li>
              </ul>
              <div v-if="!item.sender && item.status === 'pending'">
                <FlButton size="small" type="success" @click="handleRequest(item, 1)">同意</FlButton>
                <FlButton size="small" type="error" @click="handleRequest(item, 2)">拒绝</FlButton>
              </div>
            </div>
          </template>
        </FlMenu>
        <LoadingMore :loading="loading" :more="hasNext" :cursor="nextCursor" @click="nextPage"></LoadingMore>
      </FlScrollBar>
    </div>
  </div>
</template>
<script setup>
import { onMounted, ref, watch } from 'vue'
import { FlScrollBar, FlTabs, FlTabPane, FlMenu, FlButton } from 'fall-ui';
import { listGroupRequestApi, handleGroupRequestApi } from '@/api/group_request'
import Avatar from './common/avatar.vue';
import { parseTime } from '@/utils/timeinfo';
import LoadingMore from './common/loading-more.vue';

const activeStatus = ref('all')

const menuWidth = ref(0)
const handleResize = ({ clientWidth }) => {
  menuWidth.value = clientWidth
}

const statusValue = {
  all: [],
  pending: [0],
  approved: [1],
  rejected: [2],
  expired: [3]
}

const statusList = [
  {
    name: 'all',
    label: "全部",
  },
  {
    name: 'pending',
    label: "待处理",
  },
  {
    name: 'approved',
    label: "已同意",
  },
  {
    name: 'rejected',
    label: "已拒绝",
  },
  {
    name: 'expired',
    label: "已过期",
  },
]

const groupApplyData = ref([])
const nextCursor = ref("")
const loading = ref(false)
const hasNext = ref(false)

watch(() => activeStatus.value, nv => {
  nextCursor.value = ""
  groupApplyData.value = []
  fetchApplyInfo(statusValue[nv])
})

const fetchApplyInfo = async (status = []) => {
  const data = {
    status,
    next: nextCursor.value,
    sort_key: "created_desc"
  }

  const result = await listGroupRequestApi(data)
  if (result && result.data) {
    groupApplyData.value.push(...result.data)
    nextCursor.value = result.next
    hasNext.value = result.has_next
  }
}

const nextPage = () => {
  loading.value = true
  fetchApplyInfo(statusValue[activeStatus.value]).finally(() => {
    loading.value = false
  })
}

const handleRequest = async (item, action) => {
  const data = {
    id: item.id,
    action
  }

  await handleGroupRequestApi(data)
  fetchApplyInfo()
}

onMounted(() => {
  fetchApplyInfo()
})
</script>
