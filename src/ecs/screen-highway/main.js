import { createApp } from 'vue'
import App from './components/com/App.vue'

import { createRouter, createWebHashHistory } from 'vue-router'

//rem 适配  1rem = 14px
import './components/module/fontSize.js'
//公共css
import './assets/css/com/lib.css'

import Index from './views/index.vue'
let routes = [
  {
    path: '/',
    name: 'index',
    component: Index,
  },
]
const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

const app = createApp(App)

import * as echarts from 'echarts'
app.config.globalProperties.$echarts = echarts

import tools from './components/module/tools.vue'
app.config.globalProperties.$tools = tools

app.use(router).mount('#app')
