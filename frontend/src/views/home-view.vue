<template>
  <fl-container style="height:100%;">
    <fl-aside style="background-color: skyblue;" width="48">
      <div style="box-sizing: border-box;padding:10px 0;height:100%;display: flex;flex-direction: column;
      align-items: center;justify-content: space-between;">
        <ul style="list-style: none;margin: 0;padding: 0;
        display:flex;flex-direction: column;align-items: center;
        gap:10px;font-size:24px;
        ">
          <li style="cursor:pointer;" @click="router.push('/home')">
            <Avatar size="40" :uri="avatarUrl"></Avatar>
          </li>
          <li v-for="(item, index) in topActions" :key="index">
            <fl-icon style="cursor:pointer;" @click="item.click">
              <component :is="item.com" />
            </fl-icon>
          </li>
        </ul>

        <ul style="list-style: none;margin: 0;padding: 0;
        display:flex;flex-direction: column;align-items: center;
        gap:10px;font-size: 24px;
        ">
          <li v-for="(item, index) in bottomActions" :key="index">
            <fl-icon style="cursor:pointer;" @click="item.click">
              <component :is="item.com" />
            </fl-icon>
          </li>
        </ul>
      </div>
    </fl-aside>
    <fl-main>
      <LayoutView></LayoutView>
    </fl-main>
  </fl-container>
</template>
<script setup>
import { FlSetting, FlComment, FlTeam, FlAddTeam, FlEarth, FlSmile } from 'fall-ui/icons'
import router from '@/router'
import LayoutView from './layout-view.vue'
import Avatar from '@/components/common/avatar.vue'
import { useUserStore } from '@/stores/user_info'
import { computed } from 'vue'
const userStore = useUserStore()

const avatarUrl = computed(() => userStore.getAvatar())

const topActions = [
  {
    com: FlComment,
    click: () => {
      router.push('/chat')
    }
  },
  {
    com: FlTeam,
    click: () => {
      router.push('/relation')
    }
  },
  {
    com: FlAddTeam,
    click: () => {
      router.push('/searchx')
    }
  },
  {
    com: FlEarth,
    click: () => {
      console.log('earth')
    }
  },
  {
    com: FlSmile,
    click: () => {
      router.push('/emoji')
    }
  }
]
const bottomActions = [
  {
    com: FlSetting,
    click: () => {
      router.push('/setting')
    }
  },
]
</script>
