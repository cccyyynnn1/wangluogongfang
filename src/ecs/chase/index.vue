<script lang="ts">
  export default {
    name: 'Chase', // 狩猎
  }
</script>

<script setup lang="ts">
  import PathTracing from './path-tracing.vue'
  import ThreatHunting from './threat-hunting.vue'

  const activeName = ref('path') // tabs选中项
</script>

<template>
  <div class="Chase-container">
    <el-tabs v-model="activeName" style="height: 100%; width: 100%">
      <el-tab-pane label="路径追踪" name="path">
        <PathTracing />
      </el-tab-pane>
      <el-tab-pane label="威胁狩猎" name="hunt">
        <ThreatHunting />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped lang="scss">
  .Chase-container {
    height: 100vh;
    overflow: hidden;
  }
  :deep() {
    .el-tabs {
      --el-tabs-header-height: 32px;
    }
    .el-tabs__header {
      margin: 0 0 12px;
    }
    // 隐藏tab下划线
    .el-tabs__nav-wrap {
      &::after {
        display: none;
      }
    }
    .el-tabs__content {
      overflow: initial;
      border: 1px solid #e9e6f9;
      // border-radius: 8px;
      height: calc(100% - 45px);
      width: 100%;
      .el-tab-pane {
        overflow: initial;
        height: 100%;
      }
    }
  }
</style>
