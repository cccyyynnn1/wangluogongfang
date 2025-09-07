<script lang="ts">
  export default {
    name: 'SensitiveOperation',
  }
</script>

<script setup lang="ts">
  import { useVModel } from '@vueuse/core'
  const props = defineProps<{
    modelValue: string
  }>()
  const emit = defineEmits<{
    (e: 'update:modelValue', modelValue: string): void
  }>()
  const password = useVModel(props, 'modelValue', emit)
</script>

<template>
  <el-form label-position="top" @submit.self.prevent>
    <el-form-item label="请输入敏感操作密码：(通过验证后方可进行操作)">
      <el-input v-model="password" placeholder="请输入敏感操作密码" show-password type="password" />
    </el-form-item>
  </el-form>
</template>
<style scoped lang="scss"></style>
