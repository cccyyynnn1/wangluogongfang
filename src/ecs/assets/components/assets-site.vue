<script setup lang="ts">
  import EditSite from '@/ecs/site/edit-site.vue'
  import Application from './assets-site/application.vue'
  import ImportAssets from './import-assets.vue'
  import ApiInterface from '@/ecs/site/interface/index.vue'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import { siteGetPageType } from '~/src/types'
  import { exportSiteAssetsApi, siteGetPageApi, siteSessionDeleteApi } from '~/src/api-ecs/site'
  import { useTableCopy } from '@/utils'
  import { downloadFile } from '~/src/utils/download'
  import _ from 'lodash'

  const props = defineProps<{
    module?: number
  }>()

  let listDate = reactive<object[]>([]) // 表格数据

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  // 数据
  const queryData = reactive<siteGetPageType>({
    pageNum: 1,
    pageSize: 20,
    hosts: undefined,
    siteName: undefined,
    user: undefined,
    module: 1,
  })

  const siteDetailVisible = ref(false)

  const listLoading = ref(false)

  const showPage = ref(false) // 站点应用开关

  const total = ref(0)

  const sessionData = ref() // 站点id

  const mode = ref('edit') // 编辑站点还是新增站点

  const showEditSite = ref(false) // 新增和编辑站点页面

  const showApiInterface = ref(false) // API接口

  const sonTitle = ref('编辑站点')

  const rowData = ref()

  // 获取表格序号
  const curIndex = computed(() => (queryData.pageNum - 1) * queryData.pageSize + 1)

  let delList: string[] = [] // 选中项的数组

  let selectList: [] = []

  const showUpload = ref(false) // 是否显示站点资产导入页面

  const remark = ref('site') // 标识

  // 多选项改变
  const setSelectRows = (e: any) => {
    delList = []
    selectList = []
    e.forEach((item: { id: string }) => {
      // @ts-ignore
      selectList.push(item)
      delList.push(item.id)
    })
  }

  // 添加
  const handleAdd = () => {
    mode.value = 'add'
    sonTitle.value = '新增站点'
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
      // if (delList.length > 0) {
      const ids = delList.join(',')
      if (queryData.module == 1) {
        const res = await exportSiteAssetsApi({ ids })
        downloadFile(res, '站点资产')
      }
      ElMessage({ message: '导出成功', type: 'success' })
      // }
    } finally {
      listLoading.value = false
    }
  }

  onMounted(() => {
    getData()
  })

  // 获取页面数据
  const getData = _.debounce(async () => {
    queryData.module = 1
    listLoading.value = true
    const { data } = await siteGetPageApi({ ...queryData })
    listDate = data.records
    total.value = data.total
    listLoading.value = false
  }, 10)

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
      const ids = delList.join(',')
      const flag = selectList.some((item: { hasChildren: boolean }) => {
        return item.hasChildren
      })
      const message = flag ? '删除项中存在API，确定删除' : '您确定要删除所选项吗'
      $baseConfirm(message, null, async () => {
        const { msg } = await siteSessionDeleteApi({ ids, module: queryData.module, deleteAll: false })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        delList = []
        await getData()
      })
    }
    if (!row.row && row.deleteAll) {
      $baseConfirm('你确定要删除所有数据吗', null, async () => {
        const { msg } = await siteSessionDeleteApi({ ids: '', module: queryData.module, deleteAll: true })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        await getData()
      })
    }
  }

  // 站点应用
  const handleApiApplication = (row: any) => {
    showPage.value = true
    sessionData.value = row
  }

  const handleApiInterface = (row: any) => {
    showApiInterface.value = true
    rowData.value = row
  }

  // 编辑
  const handleEdit = (row: any) => {
    rowData.value = row
    mode.value = 'edit'
    sonTitle.value = '编辑站点'
    showEditSite.value = true
  }

  // 编辑和新增站点
  const addSite = (val: boolean, res?: string) => {
    showEditSite.value = val
    mode.value = res || ''
    getData()
  }

  const formatHost = (str: string) => {
    let res = ''
    try {
      const obj = JSON.parse(str)
      const keys = Object.keys(obj)

      keys.forEach((item, index) => {
        if (index !== keys.length - 1) {
          res += obj[item].selected.length > 0 ? `${obj[item].selected.join(',')}<br />` : `${item}<br />`
        } else {
          res += obj[item].selected.length > 0 ? `${obj[item].selected.join(',')}` : `${item}`
        }
      })
    } catch {
      res = str
    }
    return res
  }
</script>

<script lang="ts">
  export default {
    name: 'AssetsSite',
  }
</script>

<template>
  <div class="assets-site">
    <el-form label-position="right" label-width="auto" :model="queryData">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-form-item label="站点名称" prop="siteName">
            <el-input v-model="queryData.siteName" clearable />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="HOST" prop="hosts">
            <el-input v-model="queryData.hosts" clearable />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="负责人" prop="user">
            <el-input v-model="queryData.user" clearable />
          </el-form-item>
        </el-col>
      </el-row>
      <div class="btns">
        <el-row :gutter="20">
          <el-button type="primary" @click="handleAdd">添加</el-button>
          <el-button type="primary" @click="importEvent">导入</el-button>
          <el-button v-permissions="['Admin']" :loading="listLoading" type="primary" @click="exportEvent">
            导出
          </el-button>
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
        <el-button :disabled="listLoading" :loading="listLoading" type="primary" @click="getData">检索</el-button>
      </div>
    </el-form>
    <el-table
      v-loading="listLoading"
      class="my-table"
      :data="listDate"
      style="width: 100%"
      @cell-contextmenu="useTableCopy"
      @selection-change="setSelectRows"
    >
      <el-table-column show-overflow-tooltip type="selection" />
      <el-table-column :index="(index) => curIndex + index" label="序号" type="index" width="70" />
      <el-table-column label="站点名称" prop="siteName" show-overflow-tooltip />
      <el-table-column label="HOST" prop="hosts" show-overflow-tooltip>
        <template #default="{ row }">
          <span style="overflow-wrap: break-word" v-html="formatHost(row.hosts)"></span>
        </template>
      </el-table-column>
      <el-table-column label="负责人" prop="user" show-overflow-tooltip width="200" />
      <el-table-column label="联系电话" prop="phone" show-overflow-tooltip width="200" />
      <el-table-column v-permissions="['Admin']" fixed="right" label="操作" width="300">
        <template #default="{ row }">
          <el-button class="row_action" size="small" @click="handleApiApplication(row)">站点应用</el-button>
          <el-button class="row_action" size="small" @click="handleApiInterface(row)">API接口</el-button>
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
      :page-sizes="[20, 30, 50]"
      :total="total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />
    <!-- 编辑新增站点 -->
    <EditSite
      :current-item="rowData"
      :mode="mode"
      :show-edit-site="showEditSite"
      source="assets"
      :title="sonTitle"
      @on-closeEvent="addSite"
    />
    <vab-dialog v-model="showApiInterface" destroy-on-close title="API接口" width="1250px">
      <api-interface :id="rowData.id" />
    </vab-dialog>

    <Application v-if="showPage" :session-data="sessionData" :show-page="showPage" @on-close-event="showPage = false" />

    <!-- 站点导入资产 -->
    <ImportAssets
      v-if="showUpload"
      :mode="remark"
      :show-upload="showUpload"
      @on-close-event="showUpload = false"
      @on-reflash="getData"
    />
  </div>
</template>

<style scoped lang="scss">
  .assets-site {
    .my-table {
      border-bottom: 1px solid var(--el-border-color-lighter);
    }
    padding: 10px 20px;
    :deep() {
      .el-form-item__label-wrap {
        margin-right: 0 !important;
      }
      .el-dialog__footer:empty {
        display: none;
      }
      .el-scrollbar {
        height: calc(100vh - 230px);
      }
    }
    .super_height {
      height: 112px;
    }

    .btns {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 20px;
      margin-left: 10px;
    }
  }
</style>
