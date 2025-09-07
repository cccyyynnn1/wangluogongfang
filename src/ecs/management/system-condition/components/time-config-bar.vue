<script lang="ts">
  export default {
    name: 'TimeConfigBar',
  }
</script>
<script setup lang="ts">
  defineProps<{
    isTrue: boolean
  }>()
  const remark = ref('0') // 选择标识

  // 时间选择数据
  const dateList = reactive([
    { id: '0', label: '实时' },
    // { id: '1', label: '1小时' },
    // { id: '2', label: '2小时' },
    // { id: '3', label: '6小时' },
    // { id: '4', label: '1天' },
    // { id: '5', label: '3天' },
    // { id: '6', label: '7天' },
    // { id: '7', label: '自定义' },
  ])

  const value1 = ref('')

  const value2 = ref('')

  const disabledDate = (time: Date) => {
    return time.getTime() < Date.now()
  }

  const clickEvent = (id: string) => {
    remark.value = id
  }
</script>

<template>
  <div class="bar" style="display: flex">
    <div
      v-for="item in dateList"
      :key="item.id"
      :class="{ btn: true, active: remark == item.id }"
      @click="clickEvent(item.id)"
    >
      {{ item.label }}
    </div>
    <div v-if="isTrue" class="div">
      <el-date-picker v-model="value1" :disabled="remark !== '7'" placeholder="开始日期" type="date" />
      <el-date-picker
        v-model="value2"
        :disabled="remark !== '7'"
        :disabled-date="disabledDate"
        placeholder="结束日期"
        type="date"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">
  .btn {
    font-size: 12px;
    padding: 4px 8px;
    color: var(--el-color-primary);
    background: #f5f4ff;
    border-radius: 4px;
  }
</style>
