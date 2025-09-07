<script lang="ts">
  export default {
    name: 'RoleManagement',
  }
</script>
<script setup lang="ts">
  import AddRole from './add-role.vue'
  import { Plus, Search } from '@element-plus/icons-vue'
  import { getSystemRoleApi, deleteSystemRoleApi, editSystemRoleApi } from '@/api-ecs/system'
  import { SysRole, tableSearch } from '@/types'
  import { useTableCopy } from '@/utils'
  import { formatNstime } from '@/utils/time'
  import TableFilterTool from '@/components/table-filter-tool.vue'
  const activeName = ref('role') // tabs选中项
  const showName = ref('index') // 是否加载
  const listLoading = ref(false) // 是否加载
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  // 表格数据
  const listDate = ref<SysRole[]>([])
  const selection = ref<SysRole[]>([]) // tabs选中项
  const infoData = ref<SysRole | undefined>(undefined)
  const layout = ref('total, sizes, prev, pager, next, jumper')
  const queryPage = reactive({
    total: 0,
    pageNum: 1,
    pageSize: 10,
    roleName: '',
    title: '',
  })
  // 获取表格序号
  const curIndex = computed(() => (queryPage.pageNum - 1) * queryPage.pageSize + 1)
  // 添加
  const handleEdit = (val: string, row?: SysRole) => {
    showName.value = val
    infoData.value = row
  }
  provide(tableSearch, () => {
    queryPage.pageNum = 1
    getSysRole()
  })
  // 删除
  const handleDelete = async (rows: SysRole[], isAll = false) => {
    const names = rows.map((item) => item.userNames).flat(2)
    const text = names.length > 0 ? `当前角色已关联以下用户:\n${names}` : '是否删除当前项？'
    $baseConfirm(
      text,
      null,
      async () => {
        const ids = rows.map((i) => i.id)
        const { code, msg } = await deleteSystemRoleApi({ ids: isAll ? '' : ids.toString(), deleteAll: isAll })
        $baseMessage(
          msg,
          code === 50 ? 'error' : 'success',
          code === 50 ? 'vab-hey-message-error' : 'vab-hey-message-success'
        )
        queryPage.pageNum = 1
        queryPage.roleName = ''
        getSysRole()
      },
      null,
      '强制删除'
    )
  }

  // 多选项改变
  const setSelectRows = (selections: SysRole[]) => {
    selection.value = selections
  }
  // 启、停用
  const handleChange = async (val: number, row: SysRole) => {
    try {
      listLoading.value = true
      const { enable, id, menuIds, roleName } = row
      await editSystemRoleApi({ enable, id, menuIds, roleName })
      $baseMessage(val === 1 ? '已启用' : '已停用', 'success', 'vab-hey-message-success')
      queryPage.pageNum = 1
      queryPage.roleName = ''
      getSysRole()
    } finally {
      listLoading.value = false
    }
  }
  const getSysRole = async () => {
    try {
      const { pageNum, pageSize, roleName } = queryPage
      const {
        data: { records, total },
      } = await getSystemRoleApi({ pageSize, pageNum, roleName })
      listDate.value = records || []
      queryPage.total = total || 0
    } finally {
      listLoading.value = false
    }
  }
  const saveHandler = (load: boolean) => {
    handleEdit('index')
    if (load) {
      queryPage.pageNum = 1
      queryPage.roleName = ''
      getSysRole()
    }
  }
  onMounted(() => {
    getSysRole()
  })
</script>

<template>
  <div class="role-container">
    <template v-if="showName === 'index'">
      <vab-query-form>
        <vab-query-form-left-panel :span="12">
          <h3>角色管理</h3>
        </vab-query-form-left-panel>
        <vab-query-form-right-panel :span="12">
          <el-button :icon="Plus" style="margin-right: 10px" type="primary" @click="handleEdit('add')">添加</el-button>
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

      <el-table
        v-loading="listLoading"
        :data="listDate"
        @cell-contextmenu="useTableCopy"
        @selection-change="setSelectRows"
      >
        <el-table-column show-overflow-tooltip type="selection" />
        <el-table-column :index="(index) => curIndex + index" label="序号" type="index" width="65" />
        <el-table-column label="名称" prop="roleName" show-overflow-tooltip width="200">
          <template #header="{ column }">
            <table-filter-tool
              v-model:filter-value="queryPage.roleName"
              :column="column"
              filter-type="text"
              :tools="['filter']"
            />
          </template>
        </el-table-column>
        <el-table-column label="状态" prop="enable" show-overflow-tooltip width="200">
          <template #default="{ row }">
            <el-switch
              v-model="row.enable"
              active-text="启用"
              :active-value="1"
              inactive-text="停用"
              :inactive-value="0"
              inline-prompt
              style="--el-switch-on-color: #13ce66; --el-switch-off-color: #409eff"
              @change="handleChange(row.enable, row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="创建人" prop="createUserName" show-overflow-tooltip width="200" />
        <el-table-column
          :formatter="({ createTime }) => formatNstime(createTime, false)"
          label="创建时间"
          prop="createTime"
          show-overflow-tooltip
          width="200"
        />
        <el-table-column label="修改人" prop="updateUserName" show-overflow-tooltip />
        <el-table-column
          :formatter="({ updateTime }) => formatNstime(updateTime, false)"
          label="修改时间"
          prop="updateTime"
          show-overflow-tooltip
        />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit('edit', row)">编辑</el-button>
            <el-button size="small" @click="handleDelete([row])">删除</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty class="vab-data-empty" description="暂无数据" />
        </template>
      </el-table>
      <el-col :span="24">
        <el-pagination
          v-model:current-page="queryPage.pageNum"
          v-model:page-size="queryPage.pageSize"
          background
          :layout="layout"
          :page-sizes="[10, 20, 30]"
          :total="queryPage.total"
          @current-change="getSysRole"
          @size-change="getSysRole"
        />
      </el-col>
    </template>
    <add-role v-else :dialog-val="infoData" :show-name="showName" @on-saveData="saveHandler" />
  </div>
</template>

<style lang="scss" scoped>
  :deep() {
    .el-scrollbar {
      height: calc(100vh - 175px) !important;
    }
  }
</style>
