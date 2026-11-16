<template>
  <fl-container style="height:100%;">
    <fl-aside style="background-color: white;">
      <div style="height:100%;display:flex;flex-direction: column;">
        <div style="padding:8px 0;;width:68px;text-align: center;">
          <fl-button type="primary" size="small" @click="handleCollapse">
            {{ collapseText[collapse] }}
          </fl-button>
        </div>
        <div v-if="false">
          <fl-input>
            <template #append>
              <fl-button>搜索</fl-button>
            </template>
          </fl-input>
        </div>
        <div style="flex:1;min-height: 0;">
          <fl-scrollbar>
            <ChatListView :collapse="collapse"></ChatListView>
          </fl-scrollbar>
        </div>
      </div>
    </fl-aside>
    <fl-main>
      <ChatPanelView></ChatPanelView>
    </fl-main>
    <fl-aside :width="width" style="background-color: green;">
      侧边
    </fl-aside>
  </fl-container>
</template>
<script setup>
import ChatListView from './chat-list-view.vue';
import ChatPanelView from './chat-panel-view.vue';
import { ref } from 'vue'
const collapseText = ref({
  false: '折叠',
  true: '展开'
})
const collapse = ref(false)
const width = ref('40')

const handleCollapse = () => {
  collapse.value = !collapse.value
}

</script>
