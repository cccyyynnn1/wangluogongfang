<script lang="ts" setup>
  import { useSettingsStore } from '@/store/modules/settings'
  const settingsStore = useSettingsStore()
  const { systemConfig } = storeToRefs(settingsStore)
</script>

<template>
  <footer class="vab-footer">
    {{ systemConfig.value.copyright }}
    <vab-icon icon="copyright-line" />
    <span class="text">v3.0.5-RC</span>
  </footer>
</template>

<style lang="scss" scoped>
  .vab-footer {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 40px;
    color: rgba(0, 0, 0, 0.45);
    background: var(--el-color-white);
    border-top: 1px dashed $base-border-color;
    position: absolute;
    bottom: 0;
    left: 10px;
    right: 10px;
    z-index: 100;
    i {
      margin: 0 5px;
    }
    .text {
      color: #fff;
      opacity: 0.8;
    }
  }
</style>
