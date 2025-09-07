<script lang="ts">
  export default {
    name: 'Pagination', // 自定义分页
  }
</script>

<script setup lang="ts">
  import { onMounted, ref, watch } from 'vue'

  const props = defineProps<{
    total: number
    currentPage: number
    pageSize: number
  }>()
  const emit = defineEmits<{
    (e: 'current-change', val: number): void
  }>()
  const pageTotal = ref()
  const num = ref(1)

  const prevDis = ref(true)
  const nextDis = ref(true)
  const timeParty: any = inject('tiemFlag') as any

  onMounted(() => {
    num.value = props.currentPage
    pageTotal.value = Math.ceil(props.total / props.pageSize)
  })

  const handleClick = (flag: 'up' | 'down') => {
    if (flag == 'up') {
      if (num.value == pageTotal.value) return
      num.value = num.value + 1 > pageTotal.value ? pageTotal.value : num.value + 1
    } else {
      if (num.value == 1) return
      num.value = num.value - 1 == 1 ? 1 : num.value - 1
    }
  }

  const handleChange = (value: any) => {
    num.value = value
  }

  watch(
    () => timeParty?.value,
    () => {
      num.value = 1
    }
  )

  watch(
    () => props.total,
    () => {
      pageTotal.value = Math.ceil(props.total / props.pageSize)
    },
    { immediate: true }
  )

  watch(
    () => num.value,
    () => {
      prevDis.value = num.value == 1 ? false : true
      nextDis.value = num.value == pageTotal.value ? false : true
      emit('current-change', num.value)
    },
    { immediate: true }
  )
</script>

<template>
  <div class="x-pagination">
    <el-icon class="x-btn" :style="{ color: prevDis ? '#000' : '#ddd' }" @click="handleClick('down')">
      <ArrowLeft />
    </el-icon>
    &nbsp;
    <el-input-number
      v-model="num"
      :controls="false"
      :max="pageTotal"
      :min="1"
      size="small"
      style="width: 45px"
      @change="handleChange"
    />
    &nbsp; / &nbsp; {{ pageTotal }}
    &nbsp;
    <el-icon class="x-btn" :style="{ color: nextDis ? '#000' : '#ddd' }" @click="handleClick('up')">
      <ArrowRight />
    </el-icon>
  </div>
</template>

<style scoped lang="scss">
  .x-pagination {
    display: flex;
    align-items: center;
    justify-content: center;
    .x-btn {
      &:hover {
        cursor: pointer;
      }
    }
  }
</style>
