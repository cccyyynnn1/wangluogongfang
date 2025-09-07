<script lang="ts">
  export default {
    name: 'PayloadLog',
  }
</script>

<script setup lang="ts">
  const props = defineProps<{
    value: any
    type: string
  }>()
</script>

<template>
  <div class="hostTable_title">Payload</div>
  <el-input
    v-copy="value.payload_printable"
    v-fold
    :autosize="{ minRows: 2 }"
    :model-value="value.payload_printable"
    readonly
    resize="none"
    type="textarea"
  />
</template>

<style scoped lang="scss"></style>
