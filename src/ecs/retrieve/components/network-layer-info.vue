<script setup lang="ts">
  const props = defineProps<{
    showInfoDialog: boolean
    currentRow: any
    tableColumn: any
  }>()

  const visible = ref(false)
  const tableColumn = ref()
  const currentRow = ref()

  // const getData = ()

  onMounted(() => {
    visible.value = props.showInfoDialog
    tableColumn.value = props.tableColumn
    currentRow.value = props.currentRow
    for (let e of tableColumn.value) {
      const key = e.fieldNameEn
      if (key) {
        e[key] = currentRow.value[key]
      }
    }
  })

  const emit = defineEmits<{
    (e: 'on-closeEvent'): void
  }>()

  // 关闭
  const handleClose = () => {
    emit('on-closeEvent')
  }
</script>

<script lang="ts">
  export default {
    name: 'InfoDialog',
  }
</script>
<template>
  <div class="info-dialog">
    <el-dialog v-model="visible" :before-close="handleClose" :fullscreen="true" title="详情">
      <el-descriptions v-if="tableColumn?.length > 0" border :column="2">
        <el-descriptions-item v-for="item in tableColumn" :key="item.id" :label="item.fieldNameCn">
          <div style="padding-right: 30px">{{ item[item.fieldNameEn] }}</div>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss"></style>
