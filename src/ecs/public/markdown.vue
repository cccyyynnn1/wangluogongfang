<script setup lang="ts">
  import 'github-markdown-css'

  import { useSettingsStore } from '@/store/modules/settings'

  import { useFavicon } from '@vueuse/core'

  const marked = require('marked')

  const icon = useFavicon()

  const settingsStore = useSettingsStore()

  const { logo } = storeToRefs(settingsStore)

  const dataList = ref()

  const allTitle = ref()

  const allData = ref()

  const activeIndex = ref(-1)

  const input = ref('')

  const jsonData = ref()

  const initData = async () => {
    const url = `/download/manual/user.md`
    try {
      const res = await axios.get(url)
      jsonData.value = marked.marked(res.data)
      const article = document.getElementById('article')
      article!.innerHTML = jsonData.value
      dataList.value = article?.childNodes
      allTitle.value = []
      allData.value = []
      dataList.value?.forEach((item: { nodeName: string }) => {
        if (
          item.nodeName == 'H1' ||
          item.nodeName == 'H2' ||
          item.nodeName == 'H3' ||
          item.nodeName == 'H4' ||
          item.nodeName == 'H5' ||
          item.nodeName == 'H6'
        ) {
          allTitle.value.push(item)
          allData.value.push(item)
        }
      })
    } catch {
      jsonData.value = ''
      ElMessage({ message: '加载文档失败...', type: 'error' })
    }
  }

  const handleActive = (row: any, index: number) => {
    const elment = document.querySelector('.content') as HTMLBodyElement
    const top = row.offsetTop
    scrollTopAnimation({ element: elment!, time: 0.3, offsetTop: top })
    activeIndex.value = index
  }

  const scrollTopAnimation = ({ element, time, offsetTop }: { element: Element; time: number; offsetTop: number }) => {
    if (!time) {
      element.scrollTop = offsetTop
      return offsetTop
    }
    const spacingTime = 20 // 设置循环的间隔时间  值越小消耗性能越高
    let spacingInex = (time * 1000) / spacingTime // 计算循环的次数
    let nowTop = element.scrollTop // 获取当前滚动条位置
    let everTop = (offsetTop - nowTop) / spacingInex // 计算每次滑动的距离
    let scrollTimer = setInterval(() => {
      if (spacingInex > 0) {
        spacingInex--
        element.scrollTop += everTop
      } else {
        clearInterval(scrollTimer) // 清除计时器
      }
    }, spacingTime)
  }

  const handleInput = () => {
    allTitle.value = allData.value.filter((item: { outerText: string }) => {
      if (input.value.trim()) {
        return item.outerText.search(input.value) !== -1
      } else {
        return true
      }
    })
  }

  onMounted(() => {
    initData()
    icon.value = logo.value
  })
</script>

<script lang="ts">
  export default {
    name: 'Markdown',
  }
</script>
<template>
  <div class="help">
    <div class="menuList">
      <el-input v-model="input" @change="handleInput" style="width: 250px; margin-left: 20px"></el-input>
      <div
        v-for="(item, index) in allTitle"
        :key="index"
        :class="[item.nodeName, 'menu', index === activeIndex ? 'active' : null]"
        @click="handleActive(item, index)"
      >
        {{ item.innerText }}
      </div>
    </div>
    <div class="content">
      <div id="article" class="markdown-body"></div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .active {
    background-color: #eeeeee;
  }
  .help {
    display: flex;
    width: 100%;
    height: 100vh;
  }
  .menuList {
    min-width: 300px;
    max-width: 300px;
    height: 100vh;
    background-color: #fafafa;
    padding: 20px 0;
    overflow-y: auto;
    .menu {
      // height: 30px;
      line-height: 30px;
      &:hover {
        cursor: pointer;
      }
    }
    .H1 {
      // font-size: 18px;
      padding-left: 20px !important;
    }
    .H2 {
      // font-size: 18px;
      padding-left: 40px;
    }
    .H3 {
      // font-size: 16px;
      padding-left: 50px;
    }
    .H4 {
      // font-size: 16px;
      padding-left: 60px;
    }
    .H5 {
      // font-size: 14px;
      padding-left: 70px;
    }
    .H6 {
      // font-size: 14px;
      padding-left: 80px;
    }
  }
  .content {
    background-color: #fafafa;
    overflow: auto !important;
    flex: 1;
    height: 100%;
    border-color: #fff;
    border-left: 1px solid var(--el-border-color);
  }
  .markdown-body {
    background-color: #fff;
    box-sizing: border-box;
    min-width: 200px;
    max-width: 1280px;
    margin: 0 auto;
    padding: 0 40px;
    pre {
      background-color: #f8f8f8;
      padding: 8px;
      border-radius: 3px;
      border: 1px solid var(--el-border-color);
    }
  }
</style>
