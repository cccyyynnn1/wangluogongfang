<script lang="ts">
  export default {
    name: 'EscUp', //
  }
</script>

<script setup lang="ts">
  const emit = defineEmits<{
    (e: 'on-click-up'): void
  }>()

  const handlerClick = () => {
    emit('on-click-up')
  }
</script>

<template>
  <div class="esc-up" @click="handlerClick">
    <el-tooltip content="上一条" effect="light" placement="top">
      <vab-icon icon="arrow-left-s-line" />
    </el-tooltip>
  </div>
</template>

<style scoped lang="scss">
  .esc-up {
    background-color: #fff;
    height: 40px;
    width: 40px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    &:hover {
      cursor: pointer;
    }
  }
</style>
