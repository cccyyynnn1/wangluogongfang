<script lang="ts">
  export default {
    name: 'LogDataOutput',
  }
</script>

<script setup lang="ts">
  import SyslogOutput from './components/syslog-output.vue'
</script>

<template>
  <div class="logDataOutput-container">
    <syslog-output />
  </div>
</template>

<style scoped lang="scss">
  .logDataOutput-container {
    height: calc(100% - 30px);
    :deep() {
      h3 {
        margin-block: 0 0.5em;
      }
      .el-table {
        height: calc(100% - 60px);
      }
    }
  }
</style>
