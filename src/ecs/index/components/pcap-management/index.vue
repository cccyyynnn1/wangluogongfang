<script lang="ts">
  export default {
    name: 'PcapManagement',
  }
</script>

<script setup lang="ts">
  import { getLocalPcapsApi, editLocalPcapsStatusApi, deleteLocalPcapsApi } from '@/api-ecs/packet-replay'
  const props = defineProps<{
    modelValue: boolean
  }>()
  const emits = defineEmits<{
    (e: 'update:modelValue', visibility: boolean): void
  }>()
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const queryForm = reactive({
    pageNum: 1,
    pageSize: 10,
    searchStr: '',
  })
  const pcapManagementVisible = useVModel(props, 'modelValue', emits)
  const pcapList = ref()
  const pcapTotal = ref(0)
  const handleGetLocalPcaps = async () => {
    const { data } = await getLocalPcapsApi(queryForm)
    pcapList.value = data.records || []
    pcapTotal.value = data.total || 0
  }

  const handleEditNote = async (row: any) => {
    const { id, remark } = row
    const { msg } = await editLocalPcapsStatusApi({ id, remark })
    row.edit = false
    $baseMessage(msg, 'success', 'vab-hey-message-success')
  }
  const handleDelete = (row: any) => {
    $baseConfirm('确认删除当前Pcap文件么？', null, async () => {
      const { msg } = await deleteLocalPcapsApi({ ids: [row.id] })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      handleGetLocalPcaps()
    })
  }
  watchEffect(() => {
    if (pcapManagementVisible.value) {
      handleGetLocalPcaps()
    }
  })
</script>

<template>
  <div class="pcap-manggement">
    <el-dialog v-model="pcapManagementVisible" title="Pcap文件管理" width="1100px">
      <el-table border class="packet-management-table" :data="pcapList" style="height: 548px">
        <!-- <el-table-column type="selection" width="55" /> -->
        <el-table-column label="序号" type="index" width="60" />
        <el-table-column label="Pcap名称" prop="originName" show-overflow-tooltip />
        <el-table-column label="备注" show-overflow-tooltip>
          <template #default="{ row }">
            <template v-if="row.edit">
              <el-input v-model="row.remark" clearable @keydown.enter="handleEditNote(row)" />
            </template>
            <template v-else>
              {{ row.remark }}
            </template>
          </template>
        </el-table-column>
        <el-table-column align="center" fixed="right" label="操作" width="140">
          <template #default="{ row }">
            <el-button v-if="row.edit" class="row_action" size="small" type="primary" @click="handleEditNote(row)">
              完成
            </el-button>
            <el-button v-else class="row_action" size="small" @click="row.edit = true">编辑</el-button>
            <el-button class="row_action" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="queryForm.pageNum"
        v-model:page-size="queryForm.pageSize"
        background
        layout="sizes, prev, pager, next, jumper"
        :page-sizes="[10, 20, 30, 40, 50, 100]"
        :total="pcapTotal"
      />
    </el-dialog>
  </div>
</template>

<style scoped lang="scss"></style>
