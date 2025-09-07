<script lang="ts">
  export default {
    name: 'Download', //
  }
</script>

<script setup lang="ts">
  import { useSettingsStore } from '@/store/modules/settings'
  const settingsStore = useSettingsStore()

  const props = defineProps<{
    data: any
  }>()
</script>

<template>
  <div class="download" @click="() => settingsStore.changeDownloadVisible(true, props.data)">
    <!-- <vab-icon icon="download-2-line" /> -->
    <img alt="" class="img" src="@/assets/tools/download.png" />
  </div>
</template>

<style scoped lang="scss">
  .img {
    height: 36px;
    width: 36px;
  }
  .download {
    margin-top: 5px;
    border-radius: 5px;
    &:hover {
      color: #7f71eb;
      // background-color: #f0eefd !important;
    }
  }
</style>
