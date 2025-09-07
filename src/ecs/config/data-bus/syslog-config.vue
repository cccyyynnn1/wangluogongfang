<script lang="ts">
  export default {
    name: 'SyslogConfig', // syslog采集
  }
</script>
<script setup lang="ts">
  import AddSyslog from './components/add-syslog.vue'

  import { Plus, Search } from '@element-plus/icons-vue'

  import { syslogGetPageApi, syslogDeleteApi } from '@/api-ecs/kafka'

  import { formatNstime } from '@/utils/time'

  import { syslogSaveOrUpdateType, warnRuleAnalysFieldsListType } from '@/types'

  import { useTableCopy } from '@/utils'

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  const showDrawer = ref(false)

  let listDate = reactive<object[]>([]) // 表格数据

  let delList: [] = [] // 删除的数组

  const layout = ref('total, sizes, prev, pager, next, jumper')

  const queryPage = reactive({
    ruleName: undefined,
    enable: undefined,
    pageNum: 1,
    pageSize: 10,
  })

  const option = [
    { label: '禁用', value: 0 },
    { label: '启用', value: 1 },
  ]

  const total = ref(0) // 总条数

  const mode = ref('') // 标识：添加还是编辑

  const listLoading = ref(false) // 是否加载

  let currentData = reactive<syslogSaveOrUpdateType>({
    // @ts-ignore
    analysFieldsList: [],
    clientIp: '',
    enable: 0,
    id: undefined,
    logSample: '',
    logType: 0,
    ruleName: '',
    saveIndex: 0,
  })

  onMounted(() => {
    getData()
  })

  // 获取数据
  const getData = async () => {
    listLoading.value = true
    const { data } = await syslogGetPageApi({ ...queryPage })
    listDate = data.records
    total.value = data.total
    listLoading.value = false
  }

  // 添加
  const handleEdit = (val: string, row?: any) => {
    mode.value = val
    currentData = row
    showDrawer.value = true
  }

  // 删除
  const handleDelete = ({ ...row }) => {
    if (row.row) {
      delList = []
      // @ts-ignore
      delList.push(row.row.id)
    }
    if (delList.length > 0 && !row.deleteAll) {
      $baseConfirm('你确定要删除当前项吗', null, async () => {
        const { msg } = await syslogDeleteApi({ ids: delList, deleteAll: false })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        delList = []
        await getData()
      })
    }
    if (!row.row && row.deleteAll) {
      $baseConfirm('你确定要删除所有数据吗', null, async () => {
        const { msg } = await syslogDeleteApi({ ids: [], deleteAll: true })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        await getData()
      })
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

  // 页容量改变
  const handleSizeChange = (val: number) => {
    queryPage.pageSize = val
    getData()
  }

  const curIndex = computed(() => (queryPage.pageNum - 1) * queryPage.pageSize + 1)

  // 页面改变
  const handleCurrentChange = (val: number) => {
    queryPage.pageNum = val
    getData()
  }

  // 关闭抽屉回调
  const closeEvent = () => {
    showDrawer.value = false
  }
</script>

<template>
  <div class="kafka-config-container">
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <el-button :icon="Plus" type="primary" @click="handleEdit('add')">添加</el-button>
        <el-dropdown style="margin-left: 10px; margin-top: -7px">
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
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12">
        <el-form inline :model="queryPage" @submit.prevent>
          <el-form-item>
            <el-select
              v-model="queryPage.enable"
              clearable
              placeholder="请选择状态"
              style="width: 140px; margin-right: 10px"
            >
              <el-option v-for="item in option" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-input v-model.trim="queryPage.ruleName" clearable placeholder="请输入名称" />
          </el-form-item>
          <el-form-item>
            <el-button :disabled="listLoading" :icon="Search" type="primary" @click="getData">检索</el-button>
          </el-form-item>
        </el-form>
      </vab-query-form-right-panel>
    </vab-query-form>

    <el-table
      v-loading="listLoading"
      :border="true"
      :data="listDate"
      row-key="id"
      @cell-contextmenu="useTableCopy"
      @selection-change="setSelectRows"
    >
      <el-table-column show-overflow-tooltip type="selection" />
      <el-table-column :align="'center'" label="序号" width="55">
        <template #default="{ $index }">
          {{ curIndex + $index }}
        </template>
      </el-table-column>
      <el-table-column :align="'center'" label="名称" prop="ruleName" show-overflow-tooltip />
      <el-table-column :align="'center'" label="客户端地址" prop="clientIp" show-overflow-tooltip />
      <!-- <el-table-column :align="'center'" label="解析器" prop="logType" show-overflow-tooltip>
        <template #default="{ row }">
          {{ row.logType == 0 ? '字符串' : 'json' }}
        </template>
      </el-table-column> -->
      <el-table-column :align="'center'" label="状态" prop="enable" show-overflow-tooltip>
        <template #default="{ row }">
          {{ row.enable == 0 ? '禁用' : '启用' }}
        </template>
      </el-table-column>
      <el-table-column :align="'center'" label="创建时间" prop="createTime" show-overflow-tooltip width="180">
        <template #default="{ row }">
          {{ row.createTime ? formatNstime(row.createTime, false) : row.createTime }}
        </template>
      </el-table-column>
      <el-table-column :align="'center'" label="更新时间" prop="updateTime" show-overflow-tooltip width="180">
        <template #default="{ row }">
          {{ row.updateTime ? formatNstime(row.updateTime, false) : row.updateTime }}
        </template>
      </el-table-column>
      <el-table-column :align="'center'" label="操作" width="200">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit('edit', row)">编辑</el-button>
          <el-button size="small" @click="handleDelete({ row })">删除</el-button>
        </template>
      </el-table-column>
      <template #empty>
        <el-empty class="vab-data-empty" description="暂无数据" />
      </template>
    </el-table>
    <el-pagination
      background
      :current-page="queryPage.pageNum"
      :layout="layout"
      :page-size="queryPage.pageSize"
      :page-sizes="[10, 20, 50, 100]"
      :total="total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />
    <add-syslog
      v-if="showDrawer"
      :current-data="currentData"
      :mode="mode"
      :show-drawer="showDrawer"
      @on-closeEvent="closeEvent"
      @on-reflash="getData"
    />
  </div>
</template>

<style lang="scss" scoped>
  .kafka-config-container {
    :deep() {
      .el-scrollbar {
        height: calc(100vh - 270px) !important;
      }
    }
  }
</style>
