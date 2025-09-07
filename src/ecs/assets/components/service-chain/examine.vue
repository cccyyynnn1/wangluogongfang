<script setup lang="ts">
  import PathTracing from '@/ecs/retrieve/path-tracing.vue'
  const props = defineProps<{
    showPage: boolean
  }>()

  const visible = ref(false) // 弹框显隐

  const emit = defineEmits<{
    (e: 'on-close-event'): void
  }>()

  // 关闭回调
  const handleClose = () => {
    emit('on-close-event')
  }

  onMounted(() => {
    visible.value = props.showPage
  })
</script>

<script lang="ts">
  export default {
    name: 'Examine',
  }
</script>
<template>
  <div class="task-details">
    <el-dialog v-model="visible" :before-close="handleClose" title="业务链" width="90%">
      <PathTracing />
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .text {
    color: #0eadff;
  }
</style>
