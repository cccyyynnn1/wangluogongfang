<script lang="ts">
  export default {
    name: 'AttackPerspective',
  }
</script>

<script setup lang="ts">
  import SearchCharts from '../../retrieve/components/search-charts.vue'
</script>

<template>
  <div class="attackPerspective-container"><SearchCharts chart-type="attack-perspective" /></div>
</template>

<style scoped lang="scss">
  .attackPerspective-container {
    padding: 0 !important;
  }
</style>
