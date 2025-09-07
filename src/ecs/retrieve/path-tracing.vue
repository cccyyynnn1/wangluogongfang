<script lang="ts">
  export default {
    name: 'PathTracing',
  }
</script>
<script setup lang="ts">
  import SearchCharts from './components/search-charts.vue'
</script>

<template>
  <div class="pathTracing-container"><SearchCharts chart-type="path-tracing" /></div>
</template>

<style scoped lang="scss">
  .pathTracing-container {
    padding: 0 !important;
  }
</style>
