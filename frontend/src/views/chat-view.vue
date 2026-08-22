<template>
  <fl-container style="height:100%;">
    <fl-aside style="background-color: white;">
      <fl-scrollbar>
        <ChatListView></ChatListView>
      </fl-scrollbar>
    </fl-aside>
    <fl-main>
      <ChatPanelView></ChatPanelView>
    </fl-main>
    <fl-aside width="200" style="background-color: green;">
      侧边
    </fl-aside>
  </fl-container>
</template>
<script setup>
import ChatListView from './chat-list-view.vue';
import ChatPanelView from './chat-panel-view.vue';

</script>
