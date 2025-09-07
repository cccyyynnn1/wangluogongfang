<script lang="ts">
  export default {
    name: 'ExtractionRule',
  }
</script>
<script setup lang="ts">
  import ExtractionRuleEdit from './extraction-rule-edit.vue'

  import { Delete } from '@element-plus/icons-vue'

  import { normalizeGroupGetAllApi, normalizeGetPageApi, normalizeDeleteApi } from '~/src/api-ecs/normalize'

  import { formatNstime } from '@/utils/time'

  import { normalizeSaveOrUpdateType } from '~/src/types'

  import { useTableCopy } from '@/utils'

  const route = useRoute()

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  // 表格数据
  let listDate = reactive<object[]>([])

  const layout = ref('total, sizes, prev, pager, next, jumper')

  const total = ref(0) // 总条数

  const queryPage = reactive({
    groupId: undefined,
    pageNum: 1,
    pageSize: 10,
  })

  // 获取表格序号
  const curIndex = computed(() => (queryPage.pageNum - 1) * queryPage.pageSize + 1)

  const options = ref() // 全部分组选项

  const currentId = ref() // 全部分组选项

  const currentRow = ref<any>({})

  const mode = ref('edit')

  const listLoading = ref(false) // 是否加载

  const showEdit = ref(false) // 是否显示编辑页

  let delList: [] = [] // 删除的数组

  onMounted(() => {
    getAllNormalizeGroup()
  })

  onMounted(() => {
    getData()
    getDataBySkip()
  })

  // 如果从站点和资产进来
  const getDataBySkip = () => {
    if (route.query.allData) {
      mode.value = 'add'
      showEdit.value = true
      currentRow.value = JSON.parse(localStorage.getItem('fieldEextraction') as string)
    } else {
      mode.value = 'edit'
    }
  }

  // 获取数据
  const getData = async () => {
    listLoading.value = true
    const { data } = await normalizeGetPageApi({ ...queryPage })
    listDate = data.records
    total.value = data.total
    listLoading.value = false
  }

  // 得到所有分组
  const getAllNormalizeGroup = async () => {
    const { data } = await normalizeGroupGetAllApi()
    options.value = data
  }

  // 编辑
  const handleEdit = (val: boolean, row?: normalizeSaveOrUpdateType) => {
    if (row) {
      currentRow.value = row
      currentId.value = row.id
    }
    mode.value = 'edit'
    showEdit.value = val
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
        const { msg } = await normalizeDeleteApi({ ids: delList.toString(), deleteAll: false })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        delList = []
        await getData()
      })
    }
    if (!row.row && row.deleteAll) {
      $baseConfirm('你确定要删除所有数据吗', null, async () => {
        const { msg } = await normalizeDeleteApi({ ids: '', deleteAll: true })
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
</script>

<template>
  <div class="extraction-rule-container">
    <div v-if="!showEdit">
      <vab-query-form>
        <vab-query-form-left-panel :span="12">
          <h3>提取规则</h3>
        </vab-query-form-left-panel>
        <vab-query-form-right-panel :span="12">
          <el-form inline :model="queryPage" @submit.prevent>
            <el-form-item>
              <el-select v-model="queryPage.groupId" class="m-2" style="width: 100%" @change="getData">
                <el-option v-for="item in options" :key="item.value" :label="item.name" :value="item.id" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-dropdown>
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
            </el-form-item>
          </el-form>
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
        <el-table-column label="策略名称" prop="name" show-overflow-tooltip />
        <el-table-column label="状态" prop="status" show-overflow-tooltip>
          <template #default="{ row }">
            <span v-if="row.status === '1'">启用</span>
            <span v-else>关闭</span>
          </template>
        </el-table-column>
        <el-table-column label="分组" prop="groupId" show-overflow-tooltip />
        <el-table-column label="requestUrl" prop="requestUrl" show-overflow-tooltip />
        <el-table-column label="requestHost" prop="requestHost" show-overflow-tooltip />
        <el-table-column label="策略描述" prop="strategyDescription" show-overflow-tooltip />
        <el-table-column label="创建时间" prop="createTime" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.createTime ? formatNstime(row.createTime, false) : row.createTime }}
          </template>
        </el-table-column>
        <el-table-column label="更新时间" prop="updateTime" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.updateTime ? formatNstime(row.updateTime, false) : row.createTime }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(true, row)">编辑</el-button>
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
    </div>

    <extraction-rule-edit
      v-else
      :id="currentId"
      :current-row="currentRow"
      :group-options="options"
      :mode="mode"
      @on-reflash="getData"
      @on-save-event="handleEdit"
    />
  </div>
</template>

<style scoped lang="scss">
  :deep() {
    .el-scrollbar {
      height: calc(100vh - 170px) !important;
    }
  }
</style>
