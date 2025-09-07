<script setup lang="ts">
  import LoadBalancingList from './components/load-balancing-list.vue'

  import { useTableCopy } from '@/utils'

  import { TableInstance } from 'element-plus'

  import { loadBalancePoolApi } from '@/api-ecs/equipment'

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
      const { data } = await loadBalancePoolApi({ ...queryPage, query: JSON.stringify(queryPage.query) })
      listDate = data.list
      total.value = data.total
    } finally {
      listLoading.value = false
    }
  }

  // 切换右侧菜单
  const changeMenu = (id: string) => {
    queryPage.query.loadBalanceId = id
    getData()
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

  onMounted(() => {
    getData()
  })
</script>
<script lang="ts">
  export default {
    name: 'POOLList', // POOL列表
  }
</script>
<template>
  <div class="POOL-list-container">
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>POOL列表</h3>
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
            <el-table-column label="POOL名称" prop="poolName" show-overflow-tooltip width="100" />
            <el-table-column label="负载均衡算法" prop="algorithm" show-overflow-tooltip width="120" />
            <el-table-column label="健康检查算法" prop="healthCheck" show-overflow-tooltip width="120" />
            <el-table-column label="POOL地址端口" prop="poolIpPort" show-overflow-tooltip width="120" />
            <el-table-column label="VS名称" prop="vsNameStr" show-overflow-tooltip width="100" />
            <el-table-column label="VS地址" prop="vsIp" show-overflow-tooltip width="100" />
            <el-table-column label="VS端口" prop="vsPort" show-overflow-tooltip width="100" />
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
  .POOL-list-container {
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
