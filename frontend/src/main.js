import { initApp } from './initApp'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './asserts/app.css'

import App from './App.vue'
import router from './router'
import fallui from 'fall-ui'
import 'fall-ui/themes/index.css'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)
app.use(router)
app.use(fallui)

app.mount('#app')

initApp()
