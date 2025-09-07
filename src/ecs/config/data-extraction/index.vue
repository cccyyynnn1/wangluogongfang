<script lang="ts">
  export default {
    name: 'DataExtraction',
  }
</script>
<script setup lang="ts">
  import AddEdit from './add-edit.vue'

  import { Delete, Plus, Search } from '@element-plus/icons-vue'

  import { customFieldDeleteApi, customFieldGetPageApi, customFieldSaveOrUpdateApi } from '~/src/api-ecs/custom-field'

  import { customFieldGetGageType, customFieldSaveOrUpdateType, IndexTypeTpye, tableSearch } from '~/src/types'
  import TableFilterTool from '@/components/table-filter-tool.vue'
  import { formatNstime } from '@/utils/time'

  import { useTableCopy } from '@/utils'

  import { useUserStore } from '@/store/modules/user'

  const { getAllIndexType } = useUserStore()

  const tagType = ref<IndexTypeTpye[]>(getAllIndexType())

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  const showDrawer = ref(false) // 打开弹框

  const title = ref('添加') // 弹框标题

  const layout = ref('total, sizes, prev, pager, next, jumper')

  let listDate = reactive<object[]>([]) // 表格数据

  let delList: [] = [] // 删除的数组

  const total = ref(0) // 总条数

  const queryPage = reactive<customFieldGetGageType>({
    pageSize: 10,
    pageNum: 1,
    indexType: 1,
    fieldNameCn: '',
  })

  // 获取表格序号
  const curIndex = computed(() => (queryPage.pageNum - 1) * queryPage.pageSize + 1)

  let currentData = reactive<customFieldSaveOrUpdateType>({
    fieldNameCn: '',
    type: 1,
    isDisplay: '0',
    id: null,
  })

  const listLoading = ref(false) // 是否加载
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
    const { data } = await customFieldGetPageApi({ ...queryPage })
    listDate = data.records
    total.value = data.total
    listLoading.value = false
  }

  // 添加
  const handleAdd = () => {
    showDrawer.value = true
    title.value = '添加'
  }

  // 编辑
  const handleEdit = (row: any) => {
    currentData = row
    showDrawer.value = true
    title.value = '编辑'
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
        const { msg } = await customFieldDeleteApi({ ids: delList.toString(), deleteAll: false })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        delList = []
        await getData()
      })
    }
    if (!row.row && row.deleteAll) {
      $baseConfirm('你确定要删除所有数据吗', null, async () => {
        const { msg } = await customFieldDeleteApi({ ids: '', deleteAll: true })
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

  // 页面改变
  const handleCurrentChange = (val: number) => {
    queryPage.pageNum = val
    getData()
  }

  //
  const formatIsDisplay = (val: number): boolean => {
    return val == 1 ? true : false
  }

  // 保存
  const closeEvent = () => {
    showDrawer.value = false
  }

  // 转换数据
  const formatDate = (val: number | string) => {
    let value = ''
    tagType.value.find((item) => {
      if (item.value == val) {
        value = item.label
      }
    })
    return value
  }
  const displayChange = async (val: string | number | boolean, row: any) => {
    listLoading.value = true
    const { fieldNameCn, id, type } = row
    const { msg } = await customFieldSaveOrUpdateApi({ fieldNameCn, id, type, isDisplay: val === 1 ? true : false })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    getData()
  }
</script>

<template>
  <div class="data-extraction-container">
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>自定义字段</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12">
        <el-button :icon="Plus" type="primary" @click="handleAdd">添加</el-button>
        <el-dropdown style="margin-left: 10px">
          <el-button type="danger">
            批量删除
            <el-icon class="el-icon--right"><arrow-down /></el-icon>
          </el-button>
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
      <el-table-column :align="'center'" label="序号" width="55">
        <template #default="{ $index }">
          {{ curIndex + $index }}
        </template>
      </el-table-column>
      <el-table-column :align="'center'" label="字段名" prop="fieldNameCn" show-overflow-tooltip>
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="queryPage.fieldNameCn"
            :column="column"
            filter-type="text"
            :tools="['filter']"
          />
        </template>
      </el-table-column>
      <el-table-column :align="'center'" label="索引类型" prop="type" show-overflow-tooltip>
        <template #default="{ row }">
          {{ formatDate(row.type) }}
        </template>
      </el-table-column>
      <el-table-column :align="'center'" label="创建时间" prop="createTime" show-overflow-tooltip>
        <template #default="{ row }">
          {{ row.createTime ? formatNstime(row.createTime, false) : row.createTime }}
        </template>
      </el-table-column>
      <el-table-column :align="'center'" label="操作" width="200">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
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
      :page-sizes="[5, 10, 20, 50]"
      :total="total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />
    <add-edit
      v-if="showDrawer"
      :current-data="currentData"
      :show-drawer="showDrawer"
      :title="title"
      @on-closeEvent="closeEvent"
      @on-reflash="getData"
    />
  </div>
</template>

<style scoped lang="scss">
  :deep() {
    .el-scrollbar {
      height: calc(100vh - 170px) !important;
    }
  }
  .data-extraction-container {
    height: 100%;
  }
</style>
