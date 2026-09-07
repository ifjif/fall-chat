import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/home'
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/login'
  },
  {
    path: '/login',
    component: () => import('../views/login-view.vue')
  },
  {
    path: '/home',
    component: () => import('../views/home-view.vue'),
    redirect: '/user-home',
    children: [
      {
        path: '/user-home',
        component: () => import('../views/user-home.vue')
      },
      {
        path: '/chat',
        component: () => import('../views/chat-view.vue')
      },
      {
        path: '/searchX',
        component: () => import('../views/searchx-view.vue')
      },
      {
        path: '/setting',
        component: () => import('../views/setting-view.vue')
      },
      {
        path: '/relation',
        component: () => import('../views/relation-view.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to, from) => {
  if (to.fullPath !== '/login') {
    // 检查
    const token = localStorage.getItem('token')
    if (!token) {
      return '/login'
    }
  }

  if (to.fullPath === '/login') {
    const token = localStorage.getItem('token')
    if (token) {
      return '/home'
    }
  }

})

export default router
