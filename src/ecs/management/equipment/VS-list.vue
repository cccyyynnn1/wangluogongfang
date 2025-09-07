<script setup lang="ts">
  import LoadBalancingList from './components/load-balancing-list.vue'

  import { useTableCopy } from '@/utils'

  import { TableInstance } from 'element-plus'

  import { loadBalanceVsApi } from '@/api-ecs/equipment'

  let listDate = reactive<object[]>([]) // 表格数据

  const layout = ref('total, sizes, prev, pager, next, jumper')

  const queryPage = reactive<{ page: number; limit: number; query: { loadBalanceId?: string } }>({
    page: 1,
    limit: 10,
    query: {
      loadBalanceId: undefined,
    },
  })

  const total = ref(0) // 总条数

  const listLoading = ref(false) // 是否加载

  const tableRef = ref<TableInstance>() // 表格实例

  // 导出
  const handleExport = () => {
    ElMessage({ message: '功能完善中...', type: 'warning' })
  }

  // 获取数据
  const getData = async () => {
    try {
      listLoading.value = true
      // @ts-ignore
      const { data } = await loadBalanceVsApi({ ...queryPage, query: JSON.stringify(queryPage.query) })
      listDate = data.list
      total.value = data.total
    } finally {
      listLoading.value = false
    }
  }

  // 页容量改变
  const handleSizeChange = (val: number) => {
    queryPage.limit = val
    getData()
  }

  const curIndex = computed(() => (queryPage.page - 1) * queryPage.limit + 1)

  // 页面改变
  const handleCurrentChange = (val: number) => {
    queryPage.page = val
    getData()
  }

  // 切换右侧菜单
  const changeMenu = (id: string) => {
    queryPage.query.loadBalanceId = id
    getData()
  }

  onMounted(() => {
    getData()
  })
</script>
<script lang="ts">
  export default {
    name: 'VSList', // VS列表
  }
</script>
<template>
  <div class="VS-list-container">
    <!-- 顶部按钮 -->
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>VS列表</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12" style="display: none">
        <el-button disabled style="margin-right: 0 !important" type="primary" @click="handleExport">导出</el-button>
      </vab-query-form-right-panel>
    </vab-query-form>
    <el-row :gutter="20">
      <el-col :span="4"><load-balancing-list @on-change-menu="changeMenu" /></el-col>
      <el-col :span="20">
        <div>
          <!-- 表格 -->
          <el-table
            ref="tableRef"
            v-loading="listLoading"
            :data="listDate"
            row-key="id"
            @cell-contextmenu="useTableCopy"
          >
            <el-table-column show-overflow-tooltip type="selection" />
            <el-table-column label="序号" width="55">
              <template #default="{ $index }">
                {{ curIndex + $index }}
              </template>
            </el-table-column>
            <el-table-column label="负载均衡" prop="loadBalanceIdStr" show-overflow-tooltip width="100" />
            <el-table-column label="VS名称" prop="vsName" show-overflow-tooltip width="100" />
            <el-table-column label="VS地址" prop="vsIp" show-overflow-tooltip width="100" />
            <el-table-column label="VS端口" prop="vsPort" show-overflow-tooltip width="100" />
            <el-table-column label="会话保持" prop="persist" show-overflow-tooltip width="100" />
            <el-table-column label="PROFILE名称" prop="profileName" show-overflow-tooltip width="120" />
            <el-table-column label="IRULES名称" prop="irulesName" show-overflow-tooltip width="120" />
            <el-table-column label="POOL名称" prop="poolNameStr" show-overflow-tooltip width="100" />
            <el-table-column label="SNAT名称" prop="snatName" show-overflow-tooltip width="100" />
            <el-table-column label="SNAT地址" prop="snatAddr" show-overflow-tooltip width="100" />
            <el-table-column label="XFORWARDEDFOR" prop="xForwardedForStr" show-overflow-tooltip width="160" />
            <el-table-column label="转换地址" prop="translateAddrStr" show-overflow-tooltip width="100" />
            <el-table-column label="转换端口" prop="translatePortStr" show-overflow-tooltip width="100" />
            <el-table-column label="创建人" prop="crtUserStr" show-overflow-tooltip width="100" />
            <el-table-column label="创建时间" prop="crtTimeStr" show-overflow-tooltip width="100" />
            <el-table-column label="修改人" prop="updateUserStr" show-overflow-tooltip width="100" />
            <el-table-column label="修改时间" prop="updateTimeStr" show-overflow-tooltip width="100" />
            <template #empty>
              <el-empty class="vab-data-empty" description="暂无数据" />
            </template>
          </el-table>
          <el-pagination
            background
            :current-page="queryPage.page"
            :layout="layout"
            :page-size="queryPage.limit"
            :page-sizes="[10, 20, 50, 100]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
          />
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped lang="scss">
  .VS-list-container {
    h3 {
      margin-block: 0 0.5em;
    }
    :deep() {
      .el-scrollbar {
        height: calc(100vh - 180px) !important;
      }
    }
  }
</style>
