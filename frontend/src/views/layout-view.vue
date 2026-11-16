<template>
  <fl-container style="height:100%;">
    <fl-header height="48" style="line-height: 48px;background-color: beige;">
      <TextEllipsis mode="scroll" :text="notice"></TextEllipsis>
    </fl-header>
    <fl-main>
      <router-view></router-view>
    </fl-main>
  </fl-container>
</template>
<script setup>
import TextEllipsis from '@/components/common/text-ellipsis.vue';
import { ref } from 'vue'

const notice = ref(`
消息瞬息可达，连接无处不在
每一次发送，都是心意的抵达
让沟通没有距离，让交流更有温度
稳定承载亿级消息，流畅守护每段对话
从指尖到云端，消息从不掉线
一行代码，连接世界
消息不丢，连接不断，我们不散
让每一条消息都有归处
用技术缩短人与人的距离
构建实时连接，创造沟通价值
即时，即达，即心
消息在飞，心意在追
连接此刻，沟通未来
你的每一句话，我们都在认真传递
快人一步，聊得尽兴
`)

</script>
