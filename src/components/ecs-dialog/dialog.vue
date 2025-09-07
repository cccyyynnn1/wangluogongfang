<script lang="ts">
  export default {
    name: 'EcsDialog',
  }
</script>

<script setup lang="ts">
  import { ElDialog } from 'element-plus'
  type MyType = {
    loading: boolean
    modelValue: boolean
    showFooter?: boolean
    width: string | number
    title: string
  }
  defineProps<MyType>()
  const emits = defineEmits<{
    (e: 'cancel'): void
    (e: 'confirm'): void
  }>()

  const handleConfirm = () => {
    emits('confirm')
  }

  const handleClose = () => {
    emits('cancel')
  }
</script>

<template>
  <el-dialog :model-value="modelValue" :title="title" :width="width" @close="() => emits('cancel')">
    <slot name="default"></slot>
    <template #footer v-if="showFooter">
      <span class="dialog-footer">
        <el-button :auto-insert-space="false" :loading="loading" type="primary" @click="handleConfirm">确认</el-button>
        <el-button :auto-insert-space="false" @click="handleClose">取消</el-button>
      </span>
    </template>
  </el-dialog>
</template>
