<script setup lang="ts">
  import { AssetsIp_Detail } from '../type'

  import { injectStrict } from '@/utils/inject'

  import { getAssetsListApi, deleteAssetsApi, exportAssetsApi } from '@/api-ecs/assets'

  import { getAssetsType, assetsQueryType } from '@/types'

  import dayjs from 'dayjs'

  import { useTableCopy } from '@/utils'

  import { getAllServerIpLabelApi } from '~/src/api-ecs/tagLib'

  import { downloadFile } from '@/utils/download'
  import { formatNstime } from '~/src/utils/time'

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  const props = defineProps<{
    allTags: []
    allSites: []
  }>()

  const { detailVisible, isEdit, currentRow } = injectStrict(AssetsIp_Detail)

  let listDate = reactive<object[]>([]) // 表格数据

  const associatedSite = ref()

  const customerTags = ref()

  // 检索参数
  // const queryData = reactive<assetsQueryType>({
  //   knowOrNot: '0',
  //   name: undefined,
  //   assetsIp: undefined,
  //   lastTime: undefined,
  //   responsible: undefined,
  // })

  let assetsIpMate: string | undefined = ''

  // 检索参数
  const queryData = reactive<getAssetsType>({
    pageNum: 1,
    pageSize: 10,
    knowOrNot: '1',
    name: undefined,
    assetsIp: undefined,
    lastTime: undefined,
    responsible: undefined,
  })

  const total = ref(0) // 总条数

  const listLoading = ref(false) // 是否加载

  const lastTime = ref()

  // 资产最近在线时间下拉选项
  const lastTimeOptions = [
    { label: '一周前', value: 'one-week-ago' },
    { label: '一个月前', value: 'one-months-ago' },
    { label: '三个月前', value: 'three-months-ago' },
    { label: '六个月前', value: 'six-months-ago' },
    { label: '一年前', value: 'one-year-ago' },
  ]

  let delList: [] = [] // 删除的数组

  // 关联站点选项
  const associatedSiteOptions = ref([{ siteName: '', hosts: '', id: '0' }])

  // 客户标签选项
  const customerTagsOptions = ref([{ name: '', value: '', id: '0' }])

  // 获取表格序号
  const curIndex = computed(() => (queryData.pageNum - 1) * queryData.pageSize + 1)

  // 改变页面容量
  function handleSizeChange(params: number) {
    queryData.pageSize = params
    getData()
  }

  // 显示详情页
  const showInfo = (row: any) => {
    isEdit.value = false
    const { assetsIp } = row
    const index = assetsIp.indexOf('/')
    currentRow.value = { ...row, assetsIp: index > 0 ? assetsIp.slice(0, index) : assetsIp }
    detailVisible.value = true
  }

  // 改变页面
  function handleCurrentChange(params: number) {
    queryData.pageNum = params
    getData()
  }

  // 导出
  const exportEvent = async () => {
    listLoading.value = true
    const { knowOrNot, name, assetsIp, lastTime, responsible } = queryData
    try {
      if (delList.length > 0) {
        const res = await exportAssetsApi({ knowOrNot, name, assetsIp, lastTime, responsible, ids: delList })
        downloadFile(res, '未知IP资产')
        ElMessage({ message: '导出成功', type: 'success' })
      } else {
        const res = await exportAssetsApi({ knowOrNot, name, assetsIp, lastTime, responsible })
        downloadFile(res, '未知IP资产')
        ElMessage({ message: '导出成功', type: 'success' })
      }
    } finally {
      listLoading.value = false
    }
  }

  // 多选项改变
  const setSelectRows = (e: any) => {
    delList = []
    e.forEach((item: any) => {
      // @ts-ignore
      delList.push(item.id)
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
      $baseConfirm('你确定要删除吗', null, async () => {
        const { msg } = await deleteAssetsApi({ ids, deleteAll: false, knowOrNot: '1' })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        delList = []
        await getData()
      })
    }
    if (!row.row && row.deleteAll) {
      $baseConfirm('你确定要删除所有数据吗', null, async () => {
        const { msg } = await deleteAssetsApi({ ids: '', deleteAll: true, knowOrNot: '1' })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        await getData()
      })
    }
  }

  // 编辑
  // function handleEdit(row: any) {
  //   isEdit.value = true
  //   detailVisible.value = true
  //   currentRow.value = row
  // }

  // 获取表格数据
  async function getData() {
    // console.log(132)
    listLoading.value = true
    if (queryData.assetsIp) {
      assetsIpMate = JSON.parse(JSON.stringify(queryData.assetsIp))
    } else {
      assetsIpMate = undefined
    }
    changeData()
    const { data } = await getAssetsListApi({
      ...queryData,
    })
    queryData.assetsIp = assetsIpMate
    if (typeof customerTags.value == 'string') {
      // @ts-ignore
      customerTags.value = customerTags.value.split(',')
    }
    if (typeof associatedSite.value == 'string') {
      // @ts-ignore
      associatedSite.value = associatedSite.value.split(',')
    }
    // const res = await JSON.parse(data)
    total.value = data.total
    listDate = data.records
    listLoading.value = false
  }

  // 转化数据
  const changeData = async () => {
    if (customerTags.value?.length > 0) {
      // @ts-ignore
      customerTags.value = customerTags.value.join(',')
      queryData.assetsIp = `${customerTags.value},${queryData.assetsIp}`
    }
    if (associatedSite.value?.length > 0) {
      // @ts-ignore
      associatedSite.value = associatedSite.value.join(',')
      const { data } = await getAllServerIpLabelApi()
      associatedSite.value = []
      data.forEach((item: any) => {
        // @ts-ignore
        associatedSite.value.push(item.value)
      })
      associatedSite.value = associatedSite.value.join(',')
      queryData.assetsIp = `${associatedSite.value},${queryData.assetsIp}`
    }
  }

  watchEffect(() => {
    if (props.allSites || props.allTags) {
      associatedSiteOptions.value = props.allSites
      customerTagsOptions.value = props.allTags
    }
  })

  const formatNetArea = (arr: Array<object>) => {
    if (arr?.length === 0) {
      return null
    } else {
      let str = ''
      arr?.forEach((item: any) => {
        str += item.name
      })
      return str
    }
  }

  onMounted(() => {
    getData()
  })

  // 得到时间
  watch(
    () => lastTime.value,
    () => {
      if (lastTime.value) {
        const timeDate = dayjs()
        let endDate = null
        switch (lastTime.value) {
          case 'one-week-ago':
            endDate = timeDate.subtract(1, 'week').format('YYYY-MM-DD HH:mm:ss')
            break
          case 'one-months-ago':
            endDate = timeDate.subtract(1, 'month').format('YYYY-MM-DD HH:mm:ss')
            break
          case 'three-months-ago':
            endDate = timeDate.subtract(3, 'month').format('YYYY-MM-DD HH:mm:ss')
            break
          case 'six-months-ago':
            endDate = timeDate.subtract(6, 'month').format('YYYY-MM-DD HH:mm:ss')
            break
          case 'one-year-ago':
            endDate = timeDate.subtract(1, 'year').format('YYYY-MM-DD HH:mm:ss')
            break
          default:
            endDate = timeDate.startOf('day').format('YYYY-MM-DD HH:mm:ss')
            break
        }
        queryData.lastTime = `${endDate}`
      } else {
        queryData.lastTime = undefined
      }
    }
  )

  defineExpose({
    getData,
  })
</script>

<script lang="ts">
  export default {
    name: 'AssetsUnknown',
  }
</script>

<template>
  <div class="assets-unknown">
    <el-form label-position="top" label-width="120px" :model="queryData">
      <el-row :gutter="20">
        <el-col :span="24">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="资产IP">
                <el-input v-model="queryData.assetsIp" clearable prop="assetsIp" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="最近在线时间" prop="lastTime">
                <el-select v-model="lastTime" clearable placeholder="请选择" style="width: 100%">
                  <el-option
                    v-for="item in lastTimeOptions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
        </el-col>
        <!-- <el-col :span="16">
                                          <el-row :gutter="20">
                                            <el-col :span="12">
                                              <el-form-item label="关联站点" prop="associatedSite">
                                                <el-select v-model="associatedSite" class="super_height" multiple placeholder="Select">
                                                  <el-option v-for="item in associatedSiteOptions" :key="item.id" :label="item.siteName"
                                                    :value="item.hosts" />
                                                </el-select>
                                              </el-form-item>
                                            </el-col>
                                            <el-col :span="12">
                                              <el-form-item label="客户标签" prop="customerTags">
                                                <el-select v-model="customerTags" class="super_height" multiple placeholder="Select">
                                                  <el-option v-for="item in customerTagsOptions" :key="item.id" :label="item.name" :value="item.value" />
                                                </el-select>
                                              </el-form-item>
                                            </el-col>
                                          </el-row>
                                        </el-col> -->
      </el-row>
      <div v-permissions="['Admin']" class="btns">
        <el-row :gutter="20">
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
      align="center"
      border
      :data="listDate"
      style="width: 100%"
      @cell-contextmenu="useTableCopy"
      @selection-change="setSelectRows"
    >
      <el-table-column show-overflow-tooltip type="selection" />
      <el-table-column align="center" :index="(index) => curIndex + index" label="序号" type="index" width="70" />
      <el-table-column align="center" label="资产IP" prop="assetsIp" />
      <el-table-column align="center" label="网络区域" prop="netArea">
        <template #default="{ row }">
          {{ formatNetArea(row.netArea) }}
        </template>
      </el-table-column>
      <el-table-column align="center" label="在线状态" prop="isOnline">
        <template #default="{ row }">
          {{ row.isOnline ? '在线' : '离线' }}
        </template>
      </el-table-column>
      <el-table-column align="center" label="最近在线时间 " prop="lastTime">
        <template #default="{ row }">
          {{ row.lastTime && formatNstime(row.lastTime, false) }}
        </template>
      </el-table-column>
      <el-table-column align="center" fixed="right" label="操作" width="150">
        <template #default="{ row }">
          <el-button class="row_action" size="small" @click="showInfo(row)">详情</el-button>
          <el-button v-permissions="['Admin']" class="row_action" size="small" @click="handleDelete({ row })">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      v-model:current-page="queryData.pageNum"
      v-model:page-size="queryData.pageSize"
      background
      class="unknown_pagination"
      layout="total, sizes, prev, pager, next, jumper"
      :page-num-sizes="[10, 20, 30]"
      :total="total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />
  </div>
</template>

<style scoped lang="scss">
  .assets-unknown {
    :deep() {
      .el-scrollbar {
        height: calc(100vh - 410px);
      }
    }
    .super_height {
      width: 100%;
      height: 112px;

      :deep() {
        .el-tooltip__trigger,
        .el-input {
          height: 100%;
        }
      }
    }
    .btns {
      margin-left: 10px;
      float: left;
    }
    .search_btn {
      float: right;
      margin-bottom: 20px;
    }
  }
</style>
