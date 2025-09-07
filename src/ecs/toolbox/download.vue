<script lang="ts">
  export default {
    name: 'EcsDownload',
  }
</script>
<script setup lang="ts">
  import { useSettingsStore } from '@/store/modules/settings'

  import { DownloadLogDelApi, DownloadLogPageApi, DownloadLogApi, DownloadPcapApi } from '~/src/api-ecs/toolbox'

  import { DowmloadType } from '~/src/types'

  import { useTableCopy } from '~/src/utils'

  import { useUserStore } from '@/store/modules/user'

  import { proxyNet } from '@/config/index'

  const settingsStore = useSettingsStore()

  const $baseMessage: any = inject('$baseMessage')

  const $baseConfirm: any = inject('$baseConfirm')

  interface Props {
    modelValue: boolean
    selectAlert?: any
    fileName?: string
  }

  const emits = defineEmits<{
    (e: 'update:modelValue', visible: boolean): void
  }>()

  const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
    selectAlert: null,
  })

  const visible = useVModel(props, 'modelValue', emits)

  // const multipleSelection = ref<number[]>([])

  const queryForm = reactive({
    pageNum: 1,
    pageSize: 10,
    dataType: 0,
    searchStr: '',
  })

  // 获取表格序号
  const curIndex = computed(() => (queryForm.pageNum - 1) * queryForm.pageSize + 1)

  // 检索返回的数据
  const queryPage = reactive({
    total: 0,
    loading: true,
    listDate: [] as object[],
  })

  const getlogPageHandle = async () => {
    queryPage.loading = true
    try {
      const {
        data: { records, total },
      } = await DownloadLogPageApi({ ...queryForm, dataType: queryForm.dataType == 0 ? null : queryForm.dataType })
      queryPage.listDate = records || []
      queryPage.total = total || 0
    } finally {
      queryPage.loading = false
    }
  }

  const url = process.env.NODE_ENV == 'development' ? proxyNet : window.location.hostname
  const URL = `https://${url}/download/exportLogs/`
  const handleDownload = (row: { filePath: string }) => {
    window.open(URL + row.filePath)
  }

  let delList: string[] = [] // 选中项的数组
  let selectList: [] = []

  const exportEvent = () => {
    if (delList.length == 0) return
    const arr: string[] = []
    queryPage.listDate.forEach((item: any) => {
      if (delList.includes(item.id) && item.status) {
        arr.push(item.filePath)
      }
    })
    arr.forEach((item: any) => {
      window.open(URL + item)
    })
  }

  // 删除
  const handleDelete = ({ ...row }) => {
    if (row.row) {
      delList = []
      // @ts-ignore
      delList.push(row.row.id)
    }
    if (delList.length > 0 && !row.deleteAll) {
      const ids = delList.join(',')
      const flag = selectList.some((item: { hasChildren: boolean }) => {
        return item.hasChildren
      })
      const message = flag ? '删除项中存在API，确定删除' : '您确定要删除所选项吗'
      $baseConfirm(message, null, async () => {
        const { msg } = await DownloadLogDelApi({ ids: delList, deleteAll: false })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        delList = []
        await getlogPageHandle()
      })
    }
    if (!row.row && row.deleteAll) {
      $baseConfirm('你确定要删除所有数据吗', null, async () => {
        const { msg } = await DownloadLogDelApi({ deleteAll: true })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        await getlogPageHandle()
      })
    }
  }

  const handleReLoadload = async (row: {
    indexType: number
    id: number
    type: number
    downloadCnd: string
    dataType: number
  }) => {
    if (row.type == 0) {
      const obj: DowmloadType = {
        id: row.id, //失败任务id
        indexType: row.indexType, // 索引类型
        status: 0,
        dataType: row.dataType,
      }
      try {
        const { msg } = await DownloadLogApi({ ...obj }, row.dataType)
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        getlogPageHandle()
      } catch (error) {
        $baseMessage('下载失败', 'error', 'vab-hey-message-error')
      }
    } else {
      const { getUserId } = useUserStore()
      try {
        const { msg } = await DownloadPcapApi({ uid: getUserId(), query: JSON.parse(row.downloadCnd) }, row.dataType)
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        getlogPageHandle()
      } catch (error: any) {
        $baseMessage('下载失败', 'error', 'vab-hey-message-error')
      }
    }
  }

  function handleSelectionChange(val: any[]) {
    delList = val.map((i) => i.id)
  }

  const useTableCopyEvent = (row: any, column: any, cell: any, event: any) => {
    if (!row.downloadCnd) return
    const arr: any[] = []
    useTableCopy(row, column, cell, event, arr)
  }

  onMounted(() => {
    if (props.fileName) {
      queryForm.searchStr = props.fileName
    }
    getlogPageHandle()
  })
</script>
<template>
  <el-dialog v-model="visible" title="下载" width="1180px">
    <div class="btn">
      <el-row :gutter="20">
        <el-button type="primary" @click="exportEvent">批量导出</el-button>
        <el-dropdown style="margin-left: 10px">
          <span class="el-dropdown-link">
            <el-button type="danger">
              批量删除
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleDelete">删除选中</el-dropdown-item>
              <el-dropdown-item @click="(e) => handleDelete({ row: false, deleteAll: true })">
                删除所有
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-select v-model="queryForm.dataType" style="margin-left: 10px; width: 120px" @change="getlogPageHandle">
          <el-option key="0" :label="'调查和告警'" :value="0" />
          <el-option key="1" :label="'调查'" :value="1" />
          <el-option key="2" :label="'告警'" :value="2" />
        </el-select>
      </el-row>
      <div class="right">
        <el-input v-model="queryForm.searchStr" clearable placeholder="请输入" style="width: 250px" />
        &nbsp;
        <el-button type="primary" @click="getlogPageHandle">检索</el-button>
      </div>
    </div>
    <div v-if="props.selectAlert" class="download-info">
      <el-descriptions border :column="4">
        <el-descriptions-item :label="props.selectAlert?.attackIp ? '源IP' : '客户端IP'">
          <div style="padding-right: 30px">{{ props.selectAlert?.attackIp || props.selectAlert?.clientIp }}</div>
        </el-descriptions-item>
        <el-descriptions-item :label="props.selectAlert?.sourcePort ? '源端口' : '客户端端口'">
          <div style="padding-right: 30px">
            {{ props.selectAlert?.sourcePort || props.selectAlert?.clientPort }}
          </div>
        </el-descriptions-item>
        <el-descriptions-item :label="props.selectAlert?.victimIp ? '目的IP' : '服务端IP'">
          <div style="padding-right: 30px">
            {{ props.selectAlert?.victimIp || props.selectAlert?.serverIp }}
          </div>
        </el-descriptions-item>
        <el-descriptions-item :label="props.selectAlert?.targetPort ? '目的端口' : '服务端端口'">
          <div style="padding-right: 30px">
            {{ props.selectAlert?.targetPort || props.selectAlert?.serverPort }}
          </div>
        </el-descriptions-item>
      </el-descriptions>
    </div>
    <el-table
      v-loading="queryPage.loading"
      align="center"
      border
      :data="queryPage.listDate"
      style="width: 100%"
      @cell-contextmenu="useTableCopyEvent"
      @selection-change="handleSelectionChange"
    >
      <el-table-column show-overflow-tooltip type="selection" />
      <el-table-column align="center" :index="(index) => curIndex + index" label="序号" type="index" width="70" />
      <el-table-column align="center" label="文件名" prop="filePath" show-overflow-tooltip width="240" />
      <el-table-column align="center" label="备注（下载数据量<=3万）" prop="downloadCnd" show-overflow-tooltip />
      <el-table-column align="center" label="数据类型" prop="type" show-overflow-tooltip width="120">
        <template #default="{ row }">
          {{ row.type == 0 ? '日志' : 'Pcap' }}
        </template>
      </el-table-column>
      <el-table-column align="center" label="状态" prop="status" show-overflow-tooltip width="120">
        <template #default="{ row }">
          <el-tag v-if="row.status == 0 || row.status == 1" class="ml-2" type="warning">进行中</el-tag>
          <el-tag v-else-if="row.status == 2" class="ml-2" type="success">完成</el-tag>
          <el-tag v-else-if="row.status == 3" class="ml-2" type="danger">失败</el-tag>
        </template>
      </el-table-column>
      <el-table-column align="center" fixed="right" label="操作" width="200">
        <template #default="{ row }">
          <el-button v-if="row.status == 2" class="row_action" size="small" @click="handleDownload(row)">
            导出
          </el-button>
          <el-button class="row_action" size="small" @click="handleDelete({ row })">
            {{ row.status == 2 ? '删除' : '取消' }}
          </el-button>
          <el-button v-if="row.status == 3" class="row_action" size="small" @click="handleReLoadload(row)">
            重新下载
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      v-model:current-page="queryForm.pageNum"
      v-model:page-size="queryForm.pageSize"
      background
      layout="sizes, prev, pager, next, jumper"
      :page-sizes="[10, 20, 30, 40, 50]"
      style="margin-bottom: 20px"
      :total="queryPage.total > 1000 ? 1000 : queryPage.total"
      @current-change="getlogPageHandle"
      @size-change="getlogPageHandle"
    />
  </el-dialog>
</template>

<style scoped lang="scss">
  .download-info {
    margin-bottom: 15px;
  }
  .btn {
    margin-left: 10px;
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  :deep() {
    .el-scrollbar {
      max-height: 416px;
      overflow-y: auto;
    }
    &::-webkit-scrollbar {
      width: 0 !important;
      height: 0;
    }
  }
</style>
