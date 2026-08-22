<template>
  <fl-menu :isCollapse="false" :width="260" :data="menuData">
    <template #icon="{ item }">
      <img style="width:50px;height:50px;border-radius: 50%;" src="../asserts/avatar1.png" />
    </template>
    <template #title="{ item }">
      <div style="padding-left: 20px;">
        <div style="display: flex;gap:2px;align-items: center;justify-content: space-between;">
          <fl-button v-if="item.type === 2" size="small" scale="0.7" type="success">
            群聊
          </fl-button>
          <span style="text-overflow: ellipsis;overflow:hidden;flex:1;font-size:14px;">
            {{ item.name }}
          </span>
          <span style="font-size: 10px;">
            {{ item.lastTime }}
          </span>
        </div>
        <div style="font-size: 12px;
                white-space: nowrap;
                overflow:hidden;
                text-overflow: ellipsis;
                ">{{ item.lastMessage }}fdsfsdfsdffdjsfsdfs</div>

      </div>
    </template>
  </fl-menu>
</template>
<script setup>
const menuData = [
  {
    type: 2,
    img: '../asserts/avatar1.png',
    name: '张三',
    lastMessage: '这是最后一条消息',
    lastTime: '2026/01/02 12:12:12'
  },
  {
    img: '../asserts/avatar1.png',
    name: '张三',
    lastMessage: '这是最后一条消息',
    lastTime: '2026/01/02 12:12:12'
  },
  {
    img: '../asserts/avatar1.png',
    name: '张三',
    lastMessage: '这是最后一条消息',
    lastTime: '2026/01/02 12:12:12'
  },
  {
    img: '../asserts/avatar1.png',
    name: '张三',
    lastMessage: '这是最后一条消息',
    lastTime: '2026/01/02 12:12:12'
  },
  {
    img: '../asserts/avatar1.png',
    name: '张三',
    lastMessage: '这是最后一条消息',
    lastTime: '2026/01/02 12:12:12'
  },
  {
    img: '../asserts/avatar1.png',
    name: '张三',
    lastMessage: '这是最后一条消息',
    lastTime: '2026/01/02 12:12:12'
  },
  {
    img: '../asserts/avatar1.png',
    name: '张三',
    lastMessage: '这是最后一条消息',
    lastTime: '2026/01/02 12:12:12'
  },
  {
    img: '../asserts/avatar1.png',
    name: '张三',
    lastMessage: '这是最后一条消息',
    lastTime: '2026/01/02 12:12:12'
  },
  {
    img: '../asserts/avatar1.png',
    name: '张三',
    lastMessage: '这是最后一条消息',
    lastTime: '2026/01/02 12:12:12'
  },
]
</script>
