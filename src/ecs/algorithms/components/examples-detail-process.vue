<script setup lang="ts">
  import { getContainersProgressApi } from '~/src/api-ecs/algorithms'

  import { useTableCopy } from '@/utils'

  const props = defineProps<{
    currentItem: any
  }>()

  const listDate = ref<object[]>([]) // 表格数据

  const listLoading = ref(false) // 是否加载

  const columnList = ref<string[]>([])

  onMounted(() => {
    getData()
  })

  // 获取数据
  const getData = async () => {
    listDate.value = []
    columnList.value = []
    listLoading.value = true
    // @ts-ignore
    const { Processes, Titles } = await getContainersProgressApi({ id: props.currentItem.Id })

    Titles.forEach((item: string, index: number) => {
      columnList.value.push(Titles[index])
    })
    Processes.forEach((item: string[]) => {
      const obj = {}
      item.forEach((td: string, index: number) => {
        // @ts-ignore
        obj[Titles[index]] = item[index]
      })
      listDate.value.push(obj)
    })
    listLoading.value = false
  }
</script>

<script lang="ts">
  export default {
    name: 'ExamplesDetailProcess',
  }
</script>

<template>
  <el-table
    v-loading="listLoading"
    align="center"
    :border="true"
    :data="listDate"
    style="width: 100%"
    @cell-contextmenu="useTableCopy"
  >
    <el-table-column v-for="item in columnList" :key="item" align="center" :label="item" :prop="item" />
  </el-table>
</template>

<style scoped lang="scss"></style>
