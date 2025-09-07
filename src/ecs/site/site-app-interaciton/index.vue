<script lang="ts">
  export default {
    name: 'SiteAppInteraction',
  }
</script>

<script setup lang="ts">
  import LargeGraph from '@/components/large-graph.vue'
  import { getSiteAppInterActionApi } from '@/api-ecs/site'
  const props = defineProps<{
    siteId?: number | string
  }>()
  const site_id = ref()
  const isFullscreen = ref(false)
  const isLight = ref(true)
  // 全屏的开关
  const clickFullScreen = () => {
    isFullscreen.value = !isFullscreen.value
  }
  const clickChangeColor = () => {
    isLight.value = !isLight.value
  }
  watch(
    props,
    () => {
      site_id.value = props.siteId
    },
    {
      immediate: true,
    }
  )
</script>

<template>
  <div class="siteAppInteraction" :class="{ 'vab-fullscreen': isFullscreen, isLight: isLight, isDark: !isLight }">
    <div class="btnn">
      <div class="item">
        <vab-icon
          class="icon"
          :icon="isFullscreen ? 'fullscreen-exit-fill' : 'fullscreen-fill'"
          @click="clickFullScreen"
        />
        {{ isFullscreen ? '取消全屏' : '全屏' }}
      </div>
      <div class="item item_color">
        <el-icon @click="clickChangeColor"><Switch /></el-icon>
        切换颜色
      </div>
    </div>
    <LargeGraph :is-full="isFullscreen" :site-id="site_id" />
  </div>
</template>

<style scoped lang="scss">
  .siteAppInteraction {
    position: relative;
    width: 100%;
    height: calc(100vh - 100px);
    padding-top: 20px;
    &.isDark {
      background-color: rgb(43, 47, 51);
      &.vab-fullscreen {
        background-color: rgb(43, 47, 51);
      }
      :deep() {
        .g6-component-contextmenu {
          background-color: #363b40 !important;
          border-color: #363b40 !important;
          color: hsla(0, 0%, 100%, 0.85) !important;
          box-shadow: 0 5px 18px 0 rgba(0, 0, 0, 0.6) !important;
        }
      }
    }
    &.isLight {
      background-color: #e7f1fa;
      &.vab-fullscreen {
        background-color: #e7f1fa;
      }
      // .g6-component-contextmenu {
      // }
    }
    .btnn {
      position: absolute;
      top: 50%;
      left: 20px;
      transform: translateY(-50%);
      z-index: 1;
      display: flex;
      flex-direction: column;
      cursor: pointer;
      padding: 18px 9px;
      background: #ffffff;
      border-radius: 4px;
      border: 1px solid #dbe6ef;
      .item {
        display: flex;
        flex-direction: column;
        align-items: center;
        color: #0d88fe;
        font-size: 12px;
        &.item_color {
          color: #00b776;
          margin-top: 10px;
        }
      }
      i {
        font-size: 17px;
        color: #0d88fe;
        width: 40px;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: #e5f3ff;
        margin-bottom: 2px;
        border-radius: 4px;
      }
      .el-icon {
        background: #ecfff5;
        color: #00b776;
      }
    }
  }
</style>
