<script lang="ts">
  export default {
    name: 'DepartmemtManagement',
  }
</script>
<script setup lang="ts">
  import { Delete, Plus, Search } from '@element-plus/icons-vue'
  import { getSysDeptApi, saveOrUpdateSysDeptApi, deleteSysDeptApi, getSysDeptByIdApi } from '@/api-ecs/system'
  import { SysDeptItem, tableSearch } from '@/types'
  import { formatNstime } from '@/utils/time'
  import { TableInstance } from 'element-plus'
  import { useTableCopy } from '@/utils'
  import TableFilterTool from '@/components/table-filter-tool.vue'
  type DeptLevelType = 'first' | 'second'

  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')

  const selection = ref<SysDeptItem[]>([]) // tabs选中项

  const tabRef = ref<TableInstance>()

  const deptFormData = reactive({
    deptName: '',
    id: undefined as number | undefined,
    parentId: undefined as number | undefined,
  })

  const listLoading = ref(false) // 是否加载

  const dialogVisible = ref(false) // 弹框显隐

  const dialogIsAdd = ref(false) // 新增还是编辑

  /** first 一级，second 二级 */
  const mode = ref<DeptLevelType>('first') // 添加弹框的分辨标识

  // 表格数据
  const layout = ref('total, sizes, prev, pager, next, jumper')
  const queryPage = reactive({
    total: 0,
    pageNum: 1,
    pageSize: 10,
    title: '',
    deptName: '',
    listDate: [] as SysDeptItem[],
  })
  provide(tableSearch, () => {
    queryPage.pageNum = 1
    queryData()
  })
  // 添加
  const handleAdd = (val: DeptLevelType, row?: SysDeptItem) => {
    dialogVisible.value = true
    mode.value = val
    dialogIsAdd.value = true
    if (row) {
      deptFormData.parentId = mode.value === 'second' ? row.id : row.parentId
      deptFormData.id = undefined
      deptFormData.deptName = ''
    } else {
      deptFormData.parentId = undefined
      deptFormData.id = undefined
      deptFormData.deptName = ''
    }
  }

  const submitDeptHandle = async () => {
    if (!deptFormData.deptName) return $baseMessage('请填写部门名称！', 'error', 'vab-hey-message-error')
    const { deptName, id, parentId } = deptFormData
    const { msg, code } = await saveOrUpdateSysDeptApi({
      deptName,
      id,
      parentId,
    })
    $baseMessage(
      msg,
      code === 50 ? 'error' : 'success',
      code === 50 ? 'vab-hey-message-error' : 'vab-hey-message-success'
    )
    queryPage.pageNum = 1
    dialogVisible.value = false
    queryData()
  }

  // 编辑
  const handleEdit = (row: SysDeptItem) => {
    dialogVisible.value = true
    dialogIsAdd.value = false
    mode.value = row.parentId === 0 ? 'first' : 'second'
    if (row) {
      deptFormData.parentId = mode.value === 'first' ? undefined : row.parentId
      deptFormData.id = row.id
      deptFormData.deptName = row.deptName
    } else {
      deptFormData.parentId = undefined
      deptFormData.id = undefined
      deptFormData.deptName = ''
    }
  }

  // 删除
  const handleDelete = async (rows: SysDeptItem[], isAll = false) => {
    $baseConfirm('你确定要删除当前项吗', null, async () => {
      const ids = rows.map((i) => i.id)
      const { code, msg } = await deleteSysDeptApi({ ids: isAll ? '' : ids.toString(), deleteAll: isAll })
      $baseMessage(
        msg,
        code === 50 ? 'error' : 'success',
        code === 50 ? 'vab-hey-message-error' : 'vab-hey-message-success'
      )
      queryPage.pageNum = 1
      dialogVisible.value = false
      queryData()
    })
  }

  // 查询
  const queryData = async () => {
    listLoading.value = true
    const { pageNum, pageSize, deptName } = queryPage
    const {
      data: { records, total },
    } = await getSysDeptApi({ pageNum, pageSize, deptName })
    queryPage.listDate = records
    queryPage.total = total
    listLoading.value = false
  }

  const loadChildren = (row: SysDeptItem, treeNode: unknown, resolve: (date: SysDeptItem[]) => void) => {
    getSysDeptByIdApi(row.id).then(({ data }) => {
      resolve(data)
    })
  }
  const beforeClose = () => {
    deptFormData.deptName = ''
    deptFormData.id = undefined
    deptFormData.parentId = undefined
  }

  // 多选项改变
  const setSelectRows = (selections: SysDeptItem[]) => {
    selection.value = selections.filter((i) => i.id !== 1)
  }
  // 页容量改变
  const handleSizeChange = (val: number) => {
    queryPage.pageSize = val
    queryData()
  }
  onMounted(() => {
    queryData()
  })
</script>

<template>
  <div class="departmemt-container">
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>部门管理</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12">
        <el-button :icon="Plus" style="margin-right: 10px" type="primary" @click="handleAdd('first')">添加</el-button>
        <el-dropdown>
          <span class="el-dropdown-link">
            <el-button type="danger">
              批量删除
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleDelete(selection)">删除选中</el-dropdown-item>
              <el-dropdown-item @click="handleDelete(selection, true)">删除所有</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </vab-query-form-right-panel>
    </vab-query-form>
    <div v-loading="listLoading" style="min-height: 280px">
      <el-table
        ref="tabRef"
        :data="queryPage.listDate"
        row-key="id"
        @cell-contextmenu="useTableCopy"
        @select="setSelectRows"
      >
        <el-table-column show-overflow-tooltip type="selection" />
        <el-table-column :align="'center'" label="部门名称" prop="deptName" show-overflow-tooltip>
          <template #header="{ column }">
            <table-filter-tool
              v-model:filter-value="queryPage.deptName"
              :column="column"
              filter-type="text"
              :tools="['filter']"
            />
          </template>
        </el-table-column>
        <el-table-column
          :align="'center'"
          :formatter="(row) => row.employeeCount || 0"
          label="人员数"
          prop="employeeCount"
          show-overflow-tooltip
        />
        <el-table-column
          :align="'center'"
          :formatter="(row) => formatNstime(row.createTime, false)"
          label="创建时间"
          prop="createTime"
          show-overflow-tooltip
        />
        <el-table-column :align="'center'" label="操作" width="400">
          <template #default="{ row }">
            <el-button v-if="row.parentId === 0" size="small" @click="handleAdd('second', row)">添加二级部门</el-button>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button v-if="row.id != 1" size="small" @click="handleDelete([row])">删除</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty class="vab-data-empty" description="暂无数据" />
        </template>
      </el-table>
    </div>
    <el-col :span="24">
      <el-pagination
        v-model:current-page="queryPage.pageNum"
        v-model:page-size="queryPage.pageSize"
        background
        :layout="layout"
        :page-sizes="[10, 20, 50, 100]"
        :total="queryPage.total"
        @current-change="queryData"
        @size-change="handleSizeChange"
      />
    </el-col>
    <el-dialog
      v-model="dialogVisible"
      :title="
        mode === 'first' ? `${dialogIsAdd ? '新增' : '编辑'}一级部门` : `${dialogIsAdd ? '新增' : '编辑'}二级部门`
      "
      width="30%"
      @close="beforeClose"
    >
      <el-form class="login-form" label-position="right" label-width="100px">
        <el-form-item :label="mode === 'first' ? '一级部门名称' : '二级部门名称'" prop="deptName">
          <el-input v-model="deptFormData.deptName" placeholder="部门名称在10个字以内" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="submitDeptHandle">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style lang="scss" scoped>
  .departmemt-container {
    :deep() {
      th.el-table-column--selection .cell {
        display: none;
      }
    }
  }
  :deep() {
    .el-scrollbar {
      height: calc(100vh - 175px) !important;
    }
  }
</style>
