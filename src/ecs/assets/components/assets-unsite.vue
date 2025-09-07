<script setup lang="ts">
  import EditSite from '@/ecs/site/edit-unsite.vue'
  import Application from './assets-site/application.vue'
  import ImportAssets from './import-assets.vue'
  import ApiInterface from '@/ecs/site/interface/index.vue'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import { injectStrict } from '@/utils/inject'
  import FieldStatistic from './field-statistic.vue'
  import { SiteUnknowPageType, siteGetPageType } from '~/src/types'
  import {
    exportSiteAssetsApi,
    siteSessionDeleteApi,
    exportUnknowAssetsApi,
    getSiteUnknowPageApi,
  } from '~/src/api-ecs/site'
  import { useTableCopy } from '@/utils'
  import { downloadFile } from '~/src/utils/download'
  import _ from 'lodash'

  let listDate = reactive<object[]>([]) // 表格数据

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  // 数据
  const queryData = reactive<SiteUnknowPageType>({
    pageNum: 1,
    pageSize: 20,
    host: '',
    clientIp: '',
    serverIp: '',
    position: '',
  })

  const listLoading = ref(false)

  const showPage = ref(false) // 站点应用开关

  const total = ref(0)

  const sessionData = ref() // 站点id

  const mode = ref('edit') // 编辑站点还是新增站点

  const showEditSite = ref(false) // 新增和编辑站点页面

  const showApiInterface = ref(false) // API接口

  const sonTitle = ref('编辑未知站点')

  const rowData = ref()

  // 获取表格序号
  const curIndex = computed(() => (queryData.pageNum - 1) * queryData.pageSize + 1)

  let delList: number[] = [] // 选中项的数组

  let selectList: [] = []

  const showUpload = ref(false) // 是否显示站点资产导入页面

  const remark = ref('site') // 标识

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
      // if (delList.length > 0) {
      const res = await exportUnknowAssetsApi({
        ids: delList,
        clientIp: queryData.clientIp,
        host: queryData.host,
        serverIp: queryData.serverIp,
        position: queryData.position,
      })
      downloadFile(res, '未知站点资产')
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
    listLoading.value = true
    const { data } = await getSiteUnknowPageApi({ ...queryData })
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

  // 编辑
  const handleEdit = (row: any) => {
    rowData.value = row
    mode.value = 'edit'
    showEditSite.value = true
  }

  // 编辑和新增站点
  const addSite = (val: boolean, res?: string) => {
    showEditSite.value = val
    mode.value = res || ''
    getData()
  }

  const tableLabel = (labels: string[]) => {
    return labels.join(' 、')
    // const curList: Omit<FingerprintItem, 'groupId' | 'isNew'>[] = list.reduce(
    //   (_labels: Omit<FingerprintItem, 'groupId' | 'isNew'>[], item: CMDBfingerprintList) => {
    //     const curLabels = item.labelList.map(({ id, labelName }) => ({ id, labelName }))
    //     return [..._labels, ...curLabels]
    //   },
    //   []
    // )
    // const selectLabels = curList.filter(({ id }) => LABEL_KEYS.includes(id)).map((i) => i.labelName)
    // const otherLabels = curList.filter(({ id }) => !LABEL_KEYS.includes(id)).map((i) => i.labelName)
    // return [...selectLabels, ...otherLabels].slice(0, 5)
  }

  const showFieldStatistic = ref(false)
  const fieldList = ref<string[]>([])
  const field_title = ref('')
  const curtableColumn = ref<{ fieldNameCn: string; fieldNameEn: string }[] | []>([])
  const handleFieldStatistic = (column: { property: string; label: string }) => {
    // showFieldStatistic.value
    showFieldStatistic.value = true
    curtableColumn.value = [{ fieldNameCn: column.label, fieldNameEn: column.property }]
    fieldList.value[0] = column.property
    field_title.value = column.label || ''
  }
</script>

<script lang="ts">
  export default {
    name: 'AssetsUNSite',
  }
</script>

<template>
  <div class="assets-site">
    <el-form label-position="right" label-width="auto" :model="queryData">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-form-item label="HOST" prop="host">
            <el-input v-model="queryData.host" clearable />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="源IP" prop="clientIp">
            <el-input v-model="queryData.clientIp" clearable />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="目的IP" prop="serverIp">
            <el-input v-model="queryData.serverIp" clearable />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="方向" prop="position">
            <el-select v-model="queryData.position" clearable style="width: 100%">
              <el-option label="内对内" value="内对内" />
              <el-option label="内对外" value="内对外" />
              <el-option label="外对内" value="外对内" />
              <el-option label="外对外" value="外对外" />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <div class="btns">
        <el-row :gutter="20">
          <el-button v-permissions="['Admin']" :loading="listLoading" type="primary" @click="exportEvent">
            导出
          </el-button>
          <!-- <el-dropdown style="margin-left: 10px">
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
          </el-dropdown> -->
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
      <el-table-column label="HOST" prop="host" show-overflow-tooltip>
        <template #header="{ column }">
          {{ column.label }}
          <el-image
            v-if="column.supportAgg"
            class="table-filter"
            :src="require('@/assets/tongji-3.svg')"
            @click="handleFieldStatistic(column)"
          />
        </template>
      </el-table-column>
      <el-table-column label="源IP" prop="clientIp" show-overflow-tooltip width="140" />
      <el-table-column label="目的IP" prop="serverIp" show-overflow-tooltip width="140">
        <template #header="{ column }">
          {{ column.label }}
          <vab-icon
            v-if="column.supportAgg"
            icon="filter-line"
            style="font-size: 12px; position: absolute; margin-left: 6px; cursor: pointer"
            @click="handleFieldStatistic(column)"
          />
        </template>
      </el-table-column>
      <el-table-column label="方向" prop="position" show-overflow-tooltip width="100" />
      <el-table-column label="Title" prop="titles" show-overflow-tooltip>
        <!-- <template #default="{ row }">
          <div v-if="row.titles.length > 0" class="labelBox foldBox">
            <span
              v-for="(item, index) in tableLabel(row.titles).split(' 、')"
              :key="index"
              :alt="item"
              class="table-label"
            >
              {{ item }}
            </span>
          </div>
        </template> -->
      </el-table-column>
      <el-table-column label="标签" prop="labels" show-overflow-tooltip>
        <!-- <template #default="{ row }">
          <div v-if="row.labels.length > 0" class="labelBox foldBox">
            <span
              v-for="(item, index) in tableLabel(row.labels).split(' 、')"
              :key="index"
              :alt="item"
              class="table-label"
            >
              {{ item }}
            </span>
          </div>
        </template> -->
      </el-table-column>
      <el-table-column v-permissions="['Admin']" fixed="right" label="操作" width="100">
        <template #default="{ row }">
          <el-button class="row_action" size="small" @click="handleEdit(row)">编辑</el-button>
          <!-- <el-button class="row_action" size="small" @click="handleDelete({ row })">删除</el-button> -->
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
      source="un-konw-site"
      :title="sonTitle"
      @on-closeEvent="addSite"
    />
    <!-- 字段统计 -->

    <field-statistic
      v-if="showFieldStatistic"
      :aggregations-payload="{ ...queryData }"
      :field-list="fieldList"
      show-context-menu
      :show-field-statistic="showFieldStatistic"
      :table-column="curtableColumn"
      :title="field_title"
      @on-closeEvent="showFieldStatistic = false"
    />
  </div>
</template>

<style scoped lang="scss">
  .labelBox {
    display: flex;
    &.foldBox {
      display: block;
      .table-label {
        display: inline-block;
      }
    }
  }
  .table-label {
    padding: 0 10px;
    height: 30px;
    border-radius: 4px;
    line-height: 30px;
    background: #f1f5fb;
    margin-right: 10px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .my-table {
    border-bottom: 1px solid var(--el-border-color-lighter);
  }
  .assets-site {
    padding: 10px 20px;
    :deep() {
      .el-dialog__footer:empty {
        display: none;
      }
      .el-form-item__label-wrap {
        margin-right: 0 !important;
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
