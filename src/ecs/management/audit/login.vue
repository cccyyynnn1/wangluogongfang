<script lang="ts">
  export default {
    name: 'Operation',
  }
</script>

<script setup lang="ts">
  import { useTableCopy } from '@/utils'
  import { getAdminLogPageApi } from '@/api-ecs/login'
  import { LoginLogType, tableSearch } from '@/types/index'
  import { formatNstime, formatTime } from '@/utils/time'
  import { Search } from '@element-plus/icons-vue'
  import TableFilterTool from '@/components/table-filter-tool.vue'
  const listLoading = ref(false) // 是否加载
  const levelOptions = [
    {
      value: 0,
      label: '信息',
    },
    {
      value: 1,
      label: '告警',
    },

    {
      value: 2,
      label: '错误',
    },
  ]
  const multipleSelection = ref<LoginLogType[]>([])
  // 表格数据
  const listDate = ref<LoginLogType[]>([])
  const total = ref(0)
  const layout = ref('total, sizes, prev, pager, next, jumper')
  const queryPage = reactive({
    pageNum: 1,
    pageSize: 20,
    ip: '',
    level: undefined,
  })
  provide(tableSearch, () => {
    getAdminLogPageHandle(true)
  })
  const getAdminLogPageHandle = async (isClear = false) => {
    if (isClear) queryPage.pageNum = 1
    const {
      data: { records, total: _total },
    } = await getAdminLogPageApi(queryPage)
    listDate.value = records
    total.value = _total
  }
  // 多选项改变
  const setSelectRows = (val: any) => {
    multipleSelection.value = val
  }
  function formatJson(filterVal: any, jsonData: any) {
    return jsonData.map((v: any) =>
      filterVal.map((j: any) => {
        return formatExcelData(v, j)
      })
    )
  }

  function formatExcelData(val: any, key: string) {
    switch (key) {
      case 'addTime':
        return formatTime(val[key] * 1000)
      default:
        return val[key] || ''
    }
  }
  const handleDownloadExcel = () => {
    const tHeader = ['名称', '登录IP', '登录路径', '等级', '操作内容', '执行结果', '操作时间']
    const filterVal = ['name', 'ip', 'path', 'proto', 'levelStr', 'msg', 'addTime']
    try {
      listLoading.value = true
      import('@/utils/excel').then((excel) => {
        const list = multipleSelection.value.length ? multipleSelection.value : listDate.value
        const data = formatJson(filterVal, list)
        excel.export_json_to_excel({
          header: tHeader,
          data,
          filename: `登录日志-${formatTime(new Date().getTime())}`,
          autoWidth: true,
          bookType: 'xlsx',
        })
        listLoading.value = false
      })
    } catch (error) {
      console.log(error)
      listLoading.value = false
    }
  }
  onMounted(() => {
    getAdminLogPageHandle()
  })
</script>

<template>
  <div class="operation-container">
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>登录日志</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12">
        <el-button plain style="margin-right: 0 !important" type="primary" @click="handleDownloadExcel">
          导出日志
        </el-button>
      </vab-query-form-right-panel>
    </vab-query-form>
    <el-table
      v-loading="listLoading"
      :data="listDate"
      row-key="id"
      @cell-contextmenu="useTableCopy"
      @selection-change="setSelectRows"
    >
      <el-table-column fixed="left" show-overflow-tooltip type="selection" width="55" />
      <el-table-column label="名称" prop="name" show-overflow-tooltip />
      <el-table-column label="登录IP" prop="ip" show-overflow-tooltip>
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="queryPage.ip"
            :column="column"
            filter-type="text"
            :tools="['filter']"
          />
        </template>
      </el-table-column>
      <el-table-column label="登录路径" prop="path" show-overflow-tooltip />
      <el-table-column label="等级" prop="levelStr" show-overflow-tooltip>
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="queryPage.level"
            :column="column"
            :filter-option="levelOptions"
            filter-type="select"
            :tools="['filter']"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作内容" prop="proto" show-overflow-tooltip />
      <el-table-column label="执行结果" prop="msg" show-overflow-tooltip />
      <!-- 操作时间：秒 -->
      <el-table-column :formatter="({ addTime }) => formatTime(addTime * 1000)" label="操作时间" prop="addTime" />
      <template #empty>
        <el-empty class="vab-data-empty" description="暂无数据" />
      </template>
    </el-table>
    <el-pagination
      v-model:current-page="queryPage.pageNum"
      v-model:page-size="queryPage.pageSize"
      background
      hide-on-single-page
      :layout="layout"
      :page-sizes="[10, 20, 50, 100]"
      :total="total"
      @current-change="() => getAdminLogPageHandle()"
      @size-change="() => getAdminLogPageHandle()"
    />
  </div>
</template>

<style lang="scss" scoped>
  .operation-container {
    position: relative;

    .export {
      float: right;
      margin-bottom: 15px;
    }
    .floatBox {
      float: right;
      margin-bottom: 12px;
      .el-input {
        width: 250px;
        margin-right: 12px;
      }
    }
  }

  :deep() {
    .el-scrollbar {
      height: calc(100vh - 175px) !important;
    }
  }
</style>
