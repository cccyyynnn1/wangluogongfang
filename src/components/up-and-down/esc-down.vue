<script lang="ts">
  export default {
    name: 'EscDown', //
  }
</script>

<script setup lang="ts">
  const emit = defineEmits<{
    (e: 'on-click-down'): void
  }>()

  const handlerClick = () => {
    emit('on-click-down')
  }
</script>

<template>
  <div class="esc-down" @click="handlerClick">
    <el-tooltip content="下一条" effect="light" placement="top">
      <vab-icon icon="arrow-right-s-line" />
    </el-tooltip>
  </div>
</template>

<style scoped lang="scss">
  .esc-down {
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
