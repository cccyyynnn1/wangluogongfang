<script setup lang="ts">
  import AddServiceChain from './service-chain/add-service-chain.vue'

  import TaskDetails from './service-chain/task-details.vue'

  import examine from './service-chain/examine.vue'

  import { useTableCopy } from '@/utils'
  import { deleteChainSsortApi, getChainSsortListApi } from '~/src/api-ecs/assets'
  import { switchOpenChainSsortApi } from '~/src/api-ecs/assets'
  import { switchCloseChainSsortApi } from '~/src/api-ecs/assets'
  import { formatNstime } from '~/src/utils/time'

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  let listDate = reactive<object[]>([{ name: '123', rule: '789' }]) // 表格数据

  let delList: [] = [] // 删除的数组

  const total = ref(0) // 总条数

  const listLoading = ref(false) // 是否加载

  const addPage = ref(false) // 是否显示添加页面

  const showPage = ref(false) // 是否显示详情页面

  const showExamine = ref(false) // 是否显示查看业务链页面

  const MD5 = ref()

  const currentItem = ref()

  // 检索参数
  const queryData = reactive({
    pageNum: 1,
    pageSize: 10,
    site: undefined,
    verifyCode: undefined,
    status: undefined,
    startTime: undefined,
    endTime: undefined,
  })

  // 获取表格序号
  const curIndex = computed(() => (queryData.pageNum - 1) * queryData.pageSize + 1)

  // 添加
  const handleAdd = () => {
    addPage.value = true
  }
  // 删除
  const handleDelete = ({ ...row }) => {
    if (row.row) {
      delList = []
      // @ts-ignore
      delList.push(row.row.id)
    }
    if (delList.length > 0 && !row.deleteAll) {
      $baseConfirm('你确定要删除吗', null, async () => {
        const { msg } = await deleteChainSsortApi({ ids: delList, deleteAll: false })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        delList = []
        await getData()
      })
    }
    if (!row.row && row.deleteAll) {
      $baseConfirm('你确定要删除所有数据吗', null, async () => {
        const { msg } = await deleteChainSsortApi({ ids: [], deleteAll: true })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        await getData()
      })
    }
  }

  // 批量开启
  const handleStart = async () => {
    if (delList.length === 0) return $baseMessage('请选择开启项', 'error', 'vab-hey-message-success')
    const { msg } = await switchOpenChainSsortApi({ ids: delList })
    ElMessage({ message: msg, type: 'success' })
  }

  // 批量结束
  const handleEnd = async () => {
    if (delList.length === 0) return $baseMessage('请选择结束项', 'error', 'vab-hey-message-success')
    const { msg } = await switchCloseChainSsortApi({ ids: delList })
    ElMessage({ message: msg, type: 'success' })
  }

  // 任务详情
  const handleInfo = (row: object) => {
    showPage.value = true
    currentItem.value = row
  }

  // 查看业务链
  const handleCheck = (row: { md5: string }) => {
    showExamine.value = true
    MD5.value = row.md5
  }

  provide('md5Data', MD5)

  // 获取表格数据
  const getData = async () => {
    try {
      listLoading.value = true
      const { data } = await getChainSsortListApi({
        pageNum: queryData.pageNum,
        pageSize: queryData.pageSize,
        site: undefined,
        status: undefined,
        // limit: 1000,
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

  onMounted(() => {
    getData()
  })
</script>

<script lang="ts">
  export default {
    name: 'ServiceChain', // 业务链梳理
  }
</script>
<template>
  <div class="service-chain-container">
    <div v-permissions="['Admin']" class="btns">
      <el-row :gutter="20">
        <el-button type="primary" @click="handleAdd">添加任务</el-button>
        <el-button type="primary" @click="handleStart">批量开启</el-button>
        <el-button type="primary" @click="handleEnd">批量结束</el-button>
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
      <el-table-column align="center" :index="(index) => curIndex + index" label="序号" type="index" width="70" />
      <el-table-column align="center" label="梳理站点域名" prop="site" show-overflow-tooltip />
      <el-table-column align="center" label="任务验证码" prop="verifyCode" show-overflow-tooltip width="150" />
      <!-- <el-table-column align="center" label="执行状态" prop="status" width="150" /> -->
      <el-table-column align="center" label="任务开始时间" prop="startTime" show-overflow-tooltip width="200">
        <template #default="{ row }">
          {{ row.startTime == 0 ? 0 : formatNstime(row.startTime * 1000, false) }}
        </template>
      </el-table-column>
      <el-table-column align="center" label="任务结束时间" prop="" width="200">
        <template #default="{ row }">
          {{ row.updateTime == 0 ? 0 : formatNstime(row.updateTime, false) }}
        </template>
      </el-table-column>
      <el-table-column align="center" fixed="right" label="操作" width="280">
        <template #default="{ row }">
          <el-button class="row_action" size="small" @click="handleInfo(row)">任务详情</el-button>
          <el-button class="row_action" size="small" @click="handleCheck(row)">查看业务链</el-button>
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
      class="known_pagination"
      layout="total, sizes, prev, pager, next, jumper"
      :page-num-sizes="[10, 20, 30]"
      :total="total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />
    <add-service-chain
      v-if="addPage"
      :current-item="currentItem"
      :show-page="addPage"
      @on-close-event="addPage = false"
      @on-reflash="getData"
    />
    <task-details
      v-if="showPage"
      :current-item="currentItem"
      :show-page="showPage"
      @on-close-event="showPage = false"
      @on-reflash="getData"
    />
    <examine v-if="showExamine" :show-page="showExamine" @on-close-event="showExamine = false" />
  </div>
</template>

<style scoped lang="scss">
  .service-chain-container {
    :deep() {
      .el-scrollbar {
        height: calc(100vh - 340px);
      }
    }
  }
  .btns {
    margin-left: 10px;
    margin-bottom: 20px;
  }
</style>
