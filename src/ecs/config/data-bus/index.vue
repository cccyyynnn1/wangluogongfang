<script lang="ts">
  export default {
    name: 'DataBus',
  }
</script>
<script setup lang="ts">
  import KafkaConfig from './kafka-config.vue'
</script>

<template>
  <kafka-config />
</template>

<style scoped lang="scss">
  :deep() {
    .el-scrollbar {
      height: calc(100vh - 170px) !important;
    }
  }
</style>
