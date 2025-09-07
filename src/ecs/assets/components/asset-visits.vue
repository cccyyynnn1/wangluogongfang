<script lang="ts">
  export default {
    name: 'AssetVisits',
  }
</script>

<script setup lang="ts">
  import SearchCharts from '../../retrieve/components/search-charts.vue'
</script>

<template>
  <div class="assetVisits-container"><SearchCharts chart-type="asset-visits" /></div>
</template>

<style scoped lang="scss">
  .assetVisits-container {
    padding: 10px 20px !important;
  }
</style>
