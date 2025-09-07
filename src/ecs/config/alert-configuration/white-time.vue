<script lang="ts">
  export default {
    name: 'WhiteTime', // 白名单时间线
  }
</script>

<script setup lang="ts">
  import { formatTime } from '@/utils/time'
  import { DeleteWarntimeApi, getWarnWhiteListApi } from '~/src/api-ecs/alert'
  import { proxyNet } from '@/config/index'

  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')

  const props = defineProps<{
    curData?: any
    isShow: boolean
  }>()
  const uploadVisible = ref(false)
  const listLoading = ref(false)
  const listDate = ref()
  const queryPage = reactive({
    pageSize: 10,
    pageNum: 1,
  })
  const total = ref(0)
  // 获取表格序号
  const curIndex = computed(() => (queryPage.pageNum - 1) * queryPage.pageSize + 1)
  const emit = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
  }>()

  const handleClose = () => {
    emit('on-closeEvent', false)
  }
  const getWhiteListHandle = async () => {
    try {
      listLoading.value = true
      const {
        data: { records, total: num },
      } = await getWarnWhiteListApi(queryPage)
      total.value = num
      listDate.value = records
    } finally {
      listLoading.value = false
    }
  }

  const url = process.env.NODE_ENV == 'development' ? proxyNet : window.location.hostname
  const URL = `https://${url}/download/warnWhiteBk/`
  const handleDownload = (row: { fileName: string }) => {
    window.open(URL + row.fileName)
  }

  const handleDelete = (row: any, deleteAll: boolean) => {
    $baseConfirm(!deleteAll ? '你确定要删除当前项吗' : '你确定要删除所有数据吗', null, async () => {
      const { msg } = await DeleteWarntimeApi({ ids: [row.id], deleteAll })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      getWhiteListHandle()
    })
  }

  onMounted(() => {
    uploadVisible.value = props.isShow
    getWhiteListHandle()
  })
</script>

<template>
  <div class="white-time">
    <el-dialog v-model="uploadVisible" :before-close="handleClose" destroy-on-close title="白名单时间线" width="1160px">
      <el-table v-loading="listLoading" :border="true" :data="listDate" row-key="id">
        <el-table-column :align="'center'" label="序号" width="55">
          <template #default="{ $index }">
            {{ curIndex + $index }}
          </template>
        </el-table-column>
        <el-table-column :align="'center'" label="版本时间" prop="createTime" show-overflow-tooltip>
          <template #default="{ row }">
            {{ formatTime(row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column :align="'center'" label="操作类型" prop="operateType" show-overflow-tooltip />
        <el-table-column :align="'center'" label="操作" width="160">
          <template #default="{ row }">
            <el-button size="small" @click="handleDownload(row)">下载</el-button>
            <el-button size="small" @click="handleDelete(row, false)">删除</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty class="vab-data-empty" description="暂无数据" />
        </template>
      </el-table>
      <template #footer>
        <el-pagination
          v-model:current-page="queryPage.pageNum"
          v-model:page-size="queryPage.pageSize"
          background
          :layout="'total, sizes, prev, pager, next, jumper'"
          :page-sizes="[10, 20, 50]"
          :total="total"
          @current-change="getWhiteListHandle"
          @size-change="getWhiteListHandle"
        />
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss"></style>
