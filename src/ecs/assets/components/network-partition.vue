<script setup lang="ts">
  import ImportAssets from './import-assets.vue'

  import UpdateNetwork from './update-network.vue'

  import { NetworkPartitionType } from '@/types'

  import {
    getNetworkPartitionListApi,
    deleteNetworkPartitionApi,
    exportNetworkPartitionApi,
    getDataSourceApi,
  } from '@/api-ecs/assets'

  import { downloadFile } from '~/src/utils/download'

  import { useTableCopy } from '@/utils'

  const option = ref<any[]>([])
  const cityOption = ref<
    {
      cn: string
      en: string
    }[]
  >([])

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  let listDate = reactive<object[]>([]) // 表格数据

  let delList: string[] = [] // 删除的数组

  const total = ref(0) // 总条数

  const listLoading = ref(false) // 是否加载

  const showPage = ref(false) // 是否显示资产添加页面

  const mode = ref()

  const currentItem = ref()

  const showUpload = ref(false) // 是否显示资产导入页面
  // 检索参数
  const queryData = reactive<NetworkPartitionType>({
    pageNum: 1,
    pageSize: 20,
    dataSources: [],
    netTypes: [],
    name: undefined,
    rule: undefined,
  })

  // 获取表格序号
  const curIndex = computed(() => (queryData.pageNum - 1) * queryData.pageSize + 1)

  // 添加
  const handleAdd = () => {
    mode.value = 'add'
    showPage.value = true
  }

  // 导入
  const importEvent = () => {
    showUpload.value = true
  }

  // 导出
  const exportEvent = async () => {
    listLoading.value = true
    const { name, rule } = queryData
    try {
      if (delList.length > 0) {
        const res = await exportNetworkPartitionApi({ ids: delList, name, rule })
        downloadFile(res, '网络分区')
        ElMessage({ message: '导出成功', type: 'success' })
      } else {
        const res = await exportNetworkPartitionApi({ name, rule })
        downloadFile(res, '网络分区')
        ElMessage({ message: '导出成功', type: 'success' })
      }
    } finally {
      listLoading.value = false
    }
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
      $baseConfirm('你确定要删除吗', null, async () => {
        const { msg } = await deleteNetworkPartitionApi({ ids, deleteAll: false })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        delList = []
        await getData()
      })
    }
    if (!row.row && row.deleteAll) {
      $baseConfirm('你确定要删除所有数据吗', null, async () => {
        const { msg } = await deleteNetworkPartitionApi({ ids: '', deleteAll: true })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        await getData()
      })
    }
  }

  // 编辑
  const handleEdit = (row: any) => {
    mode.value = 'edit'
    currentItem.value = row
    showPage.value = true
  }

  // 获取表格数据
  const getData = async () => {
    try {
      listLoading.value = true
      const { data } = await getNetworkPartitionListApi({
        ...queryData,
        dataSources: queryData.dataSources,
        netTypes: queryData.netTypes,
      })
      total.value = data.total
      listDate = data.records
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

  // 多选项改变
  const setSelectRows = (e: any) => {
    delList = []
    e.forEach((item: any) => {
      // @ts-ignore
      delList.push(item.id)
    })
  }

  const changeData = (data: string) => {
    const res = option.value.find((item: any) => {
      return item.dictValue == data
    })
    return res?.dictLabel
  }

  const changeNetTypeData = (val: number) => {
    if (val == 0) {
      return '内网资产'
    } else if (val == 1) {
      return '外网资产'
    }
    return ''
  }

  const getDataSource = async () => {
    const { data } = await getDataSourceApi()
    option.value = data.datasource || []
    cityOption.value = data.region || []
  }

  onMounted(() => {
    getData()
    getDataSource()
  })
</script>

<script lang="ts">
  export default {
    name: 'NetworkPartition', // 网络分区
  }
</script>
<template>
  <div class="network-partition">
    <el-form label-position="right" label-width="auto" :model="queryData">
      <el-row :gutter="20">
        <el-col :span="24">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-form-item label="数据中心归属" prop="dataSource">
                <el-select
                  v-model="queryData.dataSources"
                  clearable
                  collapse-tags
                  multiple
                  placeholder="请选择数据中心归属"
                  style="width: 100%"
                >
                  <el-option
                    v-for="item in option"
                    :key="item.dictValue"
                    :label="item.dictLabel"
                    :value="item.dictValue"
                  />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="资产属性" prop="netType">
                <el-select
                  v-model="queryData.netTypes"
                  clearable
                  multiple
                  placeholder="请选择资产属性"
                  style="width: 100%"
                >
                  <el-option label="内网资产" :value="0" />
                  <el-option label="外网资产" :value="1" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="分区名称" prop="name">
                <el-input v-model="queryData.name" clearable />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="分区规则" prop="rule">
                <el-input v-model="queryData.rule" clearable />
              </el-form-item>
            </el-col>
          </el-row>
        </el-col>
      </el-row>
      <div v-permissions="['Admin']" class="btns">
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
      <el-button class="search_btn" :loading="listLoading" type="primary" @click="getData">检索</el-button>
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
      <el-table-column label="数据中心归属" prop="dataSource" show-overflow-tooltip>
        <template #default="{ row }">
          {{ changeData(row.dataSource) }}
        </template>
      </el-table-column>
      <el-table-column label="分区名称" prop="name" show-overflow-tooltip />
      <el-table-column label="分区规则" prop="rule" show-overflow-tooltip />
      <el-table-column label="资产属性" prop="netType" show-overflow-tooltip>
        <template #default="{ row }">
          {{ changeNetTypeData(row.netType) }}
        </template>
      </el-table-column>
      <el-table-column fixed="right" label="操作" width="150">
        <template #default="{ row }">
          <el-button class="row_action" size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button
            v-if="!row.isDefault"
            v-permissions="['Admin']"
            class="row_action"
            size="small"
            @click="handleDelete({ row })"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      v-model:current-page="queryData.pageNum"
      v-model:page-size="queryData.pageSize"
      background
      class="known_pagination"
      layout="total, sizes, prev, pager, next, jumper"
      :page-sizes="[20, 30, 50]"
      :total="total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />
    <!-- 导入资产 -->
    <ImportAssets
      v-if="showUpload"
      mode="net"
      :show-upload="showUpload"
      @on-close-event="showUpload = false"
      @on-reflash="getData"
    />
    <!-- 新增和编辑 -->
    <update-network
      v-if="showPage"
      :city-options="cityOption"
      :current-item="currentItem"
      :modes="mode"
      :options="option"
      :show-page="showPage"
      @on-close-event="showPage = false"
      @on-reflash="getData"
    />
  </div>
</template>

<style scoped lang="scss">
  .network-partition {
    padding: 10px 20px;
    :deep() {
      .el-form-item__label-wrap {
        margin-right: 0 !important;
      }
      .el-scrollbar {
        height: calc(100vh - 230px);
      }
    }
  }

  .my-table {
    border-bottom: 1px solid var(--el-border-color-lighter);
  }

  .btns {
    margin-left: 10px;
    float: left;
  }
  .search_btn {
    float: right;
    margin-bottom: 20px;
  }
</style>
