<script lang="ts">
  export default {
    name: 'ToolBox', // 工具箱
  }
</script>

<script setup lang="ts">
  import { useSettingsStore } from '@/store/modules/settings'
  const settingsStore = useSettingsStore()
</script>

<template>
  <div class="toolBox" @click="() => settingsStore.changeToolboxVisible(true)">
    <img alt="" class="img" src="@/assets/tools/toolbox.png" />
  </div>
</template>

<style scoped lang="scss">
  .img {
    height: 36px;
    width: 36px;
  }
  .toolBox {
    margin-top: 5px;
    border-radius: 5px;
    &:hover {
      color: #7f71eb;
    }
  }
</style>
