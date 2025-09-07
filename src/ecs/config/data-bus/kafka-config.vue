<script lang="ts">
  export default {
    name: 'KafkaConfig',
  }
</script>
<script setup lang="ts">
  import AddKafka from './components/add-kafka.vue'

  import { Plus, Search } from '@element-plus/icons-vue'

  import { kafkaDeleteApi, kafkaGetPageApi } from '@/api-ecs/kafka'

  import { formatNstime } from '@/utils/time'

  import { kafkaSaveOrUpdateType, tableSearch } from '@/types'

  import { useTableCopy } from '@/utils'
  import TableFilterTool from '@/components/table-filter-tool.vue'
  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  const showDrawer = ref(false)

  let listDate = reactive<object[]>([]) // 表格数据

  let delList: [] = [] // 删除的数组

  const layout = ref('total, sizes, prev, pager, next, jumper')

  const queryPage = reactive({
    kafkaName: '',
    pageNum: 1,
    pageSize: 10,
  })

  const total = ref(0) // 总条数

  const mode = ref('') // 标识：添加还是编辑

  const listLoading = ref(false) // 是否加载

  let currentData = reactive<kafkaSaveOrUpdateType>({
    name: '',
    addressPorts: '',
  })
  provide(tableSearch, () => {
    queryPage.pageNum = 1
    getData()
  })
  onMounted(() => {
    getData()
  })

  // 获取数据
  const getData = async () => {
    listLoading.value = true
    const { data } = await kafkaGetPageApi({ ...queryPage })
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
        const { msg } = await kafkaDeleteApi({ ids: delList.toString(), deleteAll: false })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        delList = []
        await getData()
      })
    }
    if (!row.row && row.deleteAll) {
      $baseConfirm('你确定要删除所有数据吗', null, async () => {
        const { msg } = await kafkaDeleteApi({ ids: '', deleteAll: true })
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
        <h3>KAFKA配置</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12">
        <el-button :icon="Plus" type="primary" @click="handleEdit('add')">添加</el-button>
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
      </vab-query-form-right-panel>
    </vab-query-form>
    <el-table
      v-loading="listLoading"
      :data="listDate"
      row-key="id"
      @cell-contextmenu="useTableCopy"
      @selection-change="setSelectRows"
    >
      <el-table-column show-overflow-tooltip type="selection" />
      <el-table-column label="序号" width="55">
        <template #default="{ $index }">
          {{ curIndex + $index }}
        </template>
      </el-table-column>
      <el-table-column label="KAFKA名称" prop="name" show-overflow-tooltip>
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="queryPage.kafkaName"
            :column="column"
            filter-type="text"
            :tools="['filter']"
          />
        </template>
      </el-table-column>
      <el-table-column label="描述" prop="description" show-overflow-tooltip />
      <el-table-column label="IP地址" prop="addressPorts" show-overflow-tooltip />
      <el-table-column label="创建时间" prop="createTime" show-overflow-tooltip>
        <template #default="{ row }">
          {{ row.createTime ? formatNstime(row.createTime, false) : row.createTime }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
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
    <add-kafka
      v-if="showDrawer"
      :current-data="currentData"
      :mode="mode"
      :show-drawer="showDrawer"
      @on-closeEvent="closeEvent"
      @on-reflash="getData"
    />
  </div>
</template>

<style lang="scss" scoped></style>
