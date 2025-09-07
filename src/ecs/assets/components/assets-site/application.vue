<script setup lang="ts">
  import { useTableCopy } from '@/utils'

  import AddApplication from './add-application.vue'

  import { deleteApplicationApi, getAppPageListApi } from '@/api-ecs/assets'

  import { ApplicationType } from '@/types'

  import ImportAssets from '../import-assets.vue'

  import { downloadFile } from '~/src/utils/download'

  import { exportApplicationApi } from '@/api-ecs/assets'

  let listDate = reactive<object[]>([]) // 表格数据

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  // 数据
  const queryData = reactive({
    pageNum: 1,
    pageSize: 10,
    siteId: 0,
  })

  const props = defineProps<{
    showPage: boolean
    sessionData: any
  }>()

  const currentItem = reactive<ApplicationType>({
    appName: '',
    appIp: '',
    siteId: undefined,
    id: undefined,
  })

  const sessionData = ref()

  const title = ref()

  const listLoading = ref(false)

  const total = ref(0)

  const showEditSite = ref(false) // 新增和编辑站点页面

  // 获取表格序号
  const curIndex = computed(() => (queryData.pageNum - 1) * queryData.pageSize + 1)

  let delList: number[] = [] // 选中项的数组

  const showUpload = ref(false) // 是否显示站点资产导入页面

  const visible = ref(false) // 弹框显隐

  let selectList: [] = []

  const mode = ref('edit') // 编辑站点还是新增站点

  const emit = defineEmits<{
    (e: 'on-close-event'): void
    (e: 'on-reflash'): void
  }>()

  // 多选项改变
  const setSelectRows = (e: any) => {
    delList = []
    selectList = []
    e.forEach((item: { id: number }) => {
      // @ts-ignore
      selectList.push(item)
      delList.push(item.id)
    })
  }

  // 添加
  const handleAdd = () => {
    mode.value = 'add'
    showEditSite.value = true
  }

  // 导入
  const importEvent = () => {
    showUpload.value = true
  }

  // 导出
  const exportEvent = async () => {
    listLoading.value = true
    try {
      const res = await exportApplicationApi({ siteId: currentItem.siteId as number, ids: delList })
      downloadFile(res, '站点应用')
      ElMessage({ message: '导出成功', type: 'success' })
    } finally {
      listLoading.value = false
    }
  }

  // 关闭回调
  const handleClose = () => {
    emit('on-close-event')
  }

  const getData = async () => {
    if (!sessionData.value.id) return
    queryData.siteId = sessionData.value.id
    listLoading.value = true
    try {
      const { data } = await getAppPageListApi({ ...queryData })
      listDate = data.records
      total.value = data.total
    } finally {
      listLoading.value = false
    }
  }

  // 改变页面容量
  function handleSizeChange(params: number) {
    queryData.pageSize = params
    getData()
  }

  // 改变页面
  function handleCurrentChange(params: number) {
    queryData.pageNum = params
    getData()
  }

  // 删除
  const handleDelete = ({ ...row }) => {
    if (row.row) {
      delList = []
      // @ts-ignore
      delList.push(row.row.id)
    }
    if (delList.length > 0 && !row.deleteAll) {
      // const ids = delList.join(',')
      $baseConfirm('你确定要删除选择数据吗', null, async () => {
        const { msg } = await deleteApplicationApi({
          siteId: currentItem.siteId as number,
          ids: delList,
          deleteAll: false,
        })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        delList = []
        await getData()
      })
    }
    if (!row.row && row.deleteAll) {
      $baseConfirm('你确定要删除所有数据吗', null, async () => {
        const { msg } = await deleteApplicationApi({ siteId: currentItem.siteId as number, deleteAll: true })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        await getData()
      })
    }
  }

  // 编辑
  const handleEdit = (row: any) => {
    for (const key in currentItem) {
      // @ts-ignore
      currentItem[key as keyof ApplicationType] = row[key as keyof ApplicationType]
    }
    mode.value = 'edit'
    showEditSite.value = true
  }

  const handleReflash = () => {
    showEditSite.value = false
    getData()
  }

  onMounted(() => {
    visible.value = props.showPage
    sessionData.value = props.sessionData
    title.value = props.sessionData?.siteName
    currentItem.siteId = sessionData.value.id
    getData()
  })
</script>

<script lang="ts">
  export default {
    name: 'Application',
  }
</script>
<template>
  <div class="application">
    <el-dialog v-model="visible" :before-close="handleClose" :title="title" width="1180px">
      <div class="btns">
        <el-row :gutter="20">
          <el-button type="primary" @click="handleAdd">添加</el-button>
          <el-button type="primary" @click="importEvent">导入</el-button>
          <el-button :loading="listLoading" type="primary" @click="exportEvent">导出</el-button>
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
        </el-row>
      </div>
      <el-table
        v-loading="listLoading"
        align="center"
        border
        :data="listDate"
        style="width: 100%"
        @cell-contextmenu="useTableCopy"
        @selection-change="setSelectRows"
      >
        <el-table-column show-overflow-tooltip type="selection" />
        <el-table-column align="center" :index="(index) => curIndex + index" label="序号" type="index" width="80" />
        <el-table-column align="center" label="应用名称" prop="appName" show-overflow-tooltip />
        <el-table-column align="center" label="应用IP" prop="appIp" show-overflow-tooltip />
        <el-table-column v-permissions="['Admin']" align="center" fixed="right" label="操作" width="150">
          <template #default="{ row }">
            <el-button class="row_action" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button class="row_action" size="small" @click="handleDelete({ row })">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="queryData.pageNum"
        v-model:page-size="queryData.pageSize"
        background
        class="site_pagination"
        layout="total, sizes, prev, pager, next, jumper"
        :page-sizes="[10, 20, 30]"
        :total="total"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      />
    </el-dialog>
    <AddApplication
      v-if="showEditSite"
      :current-item="currentItem"
      :mode="mode"
      :show-edit-site="showEditSite"
      @on-closeEvent="showEditSite = false"
      @on-reflash="handleReflash"
    />
    <ImportAssets
      v-if="showUpload"
      mode="app"
      :show-upload="showUpload"
      :site-id="currentItem.siteId"
      @on-close-event="showUpload = false"
      @on-reflash="getData"
    />
  </div>
</template>

<style scoped lang="scss">
  .btns {
    display: flex;
    align-items: center;
    margin-left: 10px;
    margin-bottom: 20px;
    justify-content: space-between;
  }
</style>
