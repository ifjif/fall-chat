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
              <img style="width:50px;height: 50px;border-radius: 50%;" :src="item.avatar" />
            </template>
            <template #title="{ item }">
              <div style="padding:0 20px;display:flex;justify-content: space-between;">
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
                  <fl-button v-if="item.relation_status == 'already_friend'" type="primary">发消息</fl-button>
                  <fl-button v-if="item.relation_status == 'request_pending'" type="primary">等待验证</fl-button>
                  <fl-button v-if="item.relation_status == 'can_add'" type="primary">添加好友</fl-button>
                  <fl-button v-if="item.relation_status == 'already_joined'" type="primary">进入群聊</fl-button>
                  <fl-button v-if="item.relation_status == 'need_approval'" type="primary">申请加入</fl-button>
                  <fl-button v-if="item.relation_status == 'cannot_join'" disabled type="warning">不可加入</fl-button>
                  <fl-button v-if="item.relation_status == 'can_join'" type="primary">加入群聊</fl-button>
                </div>
              </div>
            </template>
          </fl-menu>
        </fl-scrollbar>
      </fl-main>

      <fl-aside width="80">
        <ul style="padding:0;list-style: none;">
          <li><fl-button type="success" @click="showCreateGroupPanel">创建群</fl-button></li>
        </ul>
      </fl-aside>
    </fl-container>
  </fl-container>

  <CreateGroupModal v-model="createGroupPanelVisible"></CreateGroupModal>
</template>
<script setup>
import { fileAccessUrlApi } from '@/api/file'
import { searchApi } from '@/api/user'
import CreateGroupModal from '@/components/modal/create-group-modal.vue'
import { FlSearch } from 'fall-ui/icons'
import { ref } from 'vue'

const menuWidth = ref(0)
const handleResize = ({ clientWidth }) => {
  menuWidth.value = clientWidth
}
const searchName = ref('')

const data = ref([])
let nextCursor = ''

const doSearch = async () => {
  try {
    const result = await searchApi({ name: searchName.value, nextCursor })
    nextCursor = result.next_cursor
    data.value = result?.result ?? []
    data.value.forEach(async (d, i) => {
      data.value[i].avatar = await fileAccessUrlApi({ file: d.avatar })
    })
  } catch (err) { }
}

const createGroupPanelVisible = ref(false)
const showCreateGroupPanel = () => {
  createGroupPanelVisible.value = true
}

</script>
