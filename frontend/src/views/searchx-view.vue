<template>
  <fl-container style="height:100%;">
    <fl-header style="padding:0 10px;">
      <fl-input v-model="searchName" :prefix-icon="FlSearch">
        <template #append>
          <fl-button @click="doSearch">搜索</fl-button>
        </template>
      </fl-input>
    </fl-header>
    <fl-container>

      <fl-main>
        <fl-scrollbar @resize="handleResize">
          <fl-menu :data="data" :width="menuWidth">
            <template #icon="{ item }">
              <Avatar :uri="item.avatar"></Avatar>
            </template>
            <template #title="{ item }">
              <div style="display:flex;justify-content: space-between;">
                <div>
                  <fl-button v-if="item.type === 'group'" size="small" scale="0.8" type="success">群聊</fl-button>
                  <span>
                    {{ item.name }}
                    <template v-if="item.nickname">({{ item.nickname }})</template>
                  </span>
                  <div style="font-size: 14px;color: gray;">
                    {{ item.description }}
                  </div>
                </div>
                <div>
                  <fl-button v-if="item.relation_status == 'already_friend'" type="success">发消息</fl-button>
                  <fl-button v-else-if="item.relation_status == 'request_pending'" disabled
                    type="warning">等待验证</fl-button>
                  <fl-button v-else-if="item.relation_status == 'can_add'" type="primary"
                    @click="addFriend(item)">添加好友</fl-button>
                  <fl-button v-else-if="item.relation_status == 'already_joined'" type="success">进入群聊</fl-button>
                  <fl-button v-else-if="item.relation_status == 'need_approval'" type="primary"
                    @click="addGroup(item)">申请加入</fl-button>
                  <fl-button v-else-if="item.relation_status == 'cannot_join'" disabled type="warning">不可加入</fl-button>
                  <fl-button v-else-if="item.relation_status == 'can_join'" type="success">加入群聊</fl-button>
                </div>
              </div>
            </template>
          </fl-menu>
        </fl-scrollbar>
      </fl-main>

      <fl-aside width="140">
        <ul style="padding:0;list-style: none;">
          <li><fl-button type="primary" text @click="showCreateGroupPanel">创建群</fl-button></li>
          <li><fl-button type="success" text @click="showApplyInfo">好友申请信息</fl-button></li>
          <li><fl-button type="warning" text @click="showGroupApplyInfo">群申请信息</fl-button></li>
        </ul>
      </fl-aside>
    </fl-container>
  </fl-container>

  <CreateGroupModal v-model="createGroupPanelVisible"></CreateGroupModal>
</template>
<script setup>
import { applyAddFriendApi } from '@/api/friendship_request'
import { applyAddGroupApi } from '@/api/group_request'
import { searchUserOrGroupApi } from '@/api/search'
import CreateGroupModal from '@/components/modal/create-group-modal.vue'
import { FlDrawer, FlMessageBox, FlTextarea } from 'fall-ui'
import { FlSearch } from 'fall-ui/icons'
import { ref, h } from 'vue'
import ApplyInfoList from '@/components/apply-info-list.vue'
import GroupApplyInfoList from '@/components/group-apply-info-list.vue'
import Avatar from '@/components/common/avatar.vue'

const menuWidth = ref(0)
const handleResize = ({ clientWidth }) => {
  menuWidth.value = clientWidth
}
const searchName = ref('')

const data = ref([])
let nextCursor = ''

const doSearch = async () => {
  if (!searchName.value.trim()) {
    searchName.value = ''
    return
  }
  try {
    const result = await searchUserOrGroupApi({ name: searchName.value, nextCursor })
    nextCursor = result.next_cursor
    data.value = result?.result ?? []
  } catch (err) { }
}

const createGroupPanelVisible = ref(false)
const showCreateGroupPanel = () => {
  createGroupPanelVisible.value = true
}

const greeting = ref('')
const addFriend = (item) => {
  const fmb = FlMessageBox({
    type: 'info',
    title: `添加 ${item.nickname} 为好友`,
    content: h(FlTextarea, {
      placeholder: '打招呼',
      maxlength: 255,
      count: true,
      // modelValue: greeting.value  (不需要，会覆盖内部默认定义的，然后丢失响应式，也就是内部监听modelValue会失效)
      'onUpdate:modelValue': (e) => greeting.value = e,
    }),
    onOk: () => { fmb.close(); doAddFriend(item) }
  })
}

const doAddFriend = async (item) => {
  greeting.value = greeting.value.trim()
  const data = {
    target_id: item.id,
    greeting: greeting.value,
    apply_type: 0, // search
  }

  await applyAddFriendApi(data)
}

const addGroup = (item) => {
  const fmb = FlMessageBox({
    type: 'info',
    title: `申请进入群：${item.name}`,
    content: h(FlTextarea, {
      count: true,
      maxlength: 255,
      'onUpdate:modelValue': (e) => greeting.value = e
    }),
    onOk: () => { fmb.close(); doAddGroup(item) }
  })
}

const doAddGroup = async (item) => {
  greeting.value = greeting.value.trim()
  const data = {
    target_id: item.id,
    greeting: greeting.value,
    invited_id: '',
    apply_type: 0,
  }

  await applyAddGroupApi(data)
}

const showApplyInfo = () => {
  FlDrawer({
    placement: 'left',
    fixedScreen: true,
    title: '好友申请信息',
    content: h(ApplyInfoList)
  })
}

const showGroupApplyInfo = () => {
  FlDrawer({
    title: '群申请信息',
    placement: 'left',
    fixedScreen: true,
    content: h(GroupApplyInfoList)
  })
}
</script>
