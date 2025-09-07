<script setup lang="ts">
  import AddUser from './add-user.vue'
  import EditUser from './edit-user.vue'
  import { useTableCopy } from '@/utils'
  import { Plus } from '@element-plus/icons-vue'
  import TableFilterTool from '@/components/table-filter-tool.vue'
  import { getUserPageApi, getAllRoleOrDeptApi, deleteUserApi, freeUserApi } from '@/api-ecs/system'
  import { formatNstime } from '@/utils/time'
  import { Role, tableSearch } from '@/types'
  import { userTypes } from '@/data/constant'
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const timeDate = ref() // 创建时间
  const isShow = ref(true) // 是否加载
  const layout = 'total, sizes, prev, pager, next, jumper'

  const queryPage = reactive({
    query: {
      pageNum: 1,
      pageSize: 10,
    },
    user: {
      nickName: '',
      mail: '',
      phone: '',
      loginName: '',
      deptIds: undefined as number | undefined,
      roleIds: undefined as number | undefined,
      userType: undefined as number | undefined,
      createSt: undefined as number | undefined,
      createEd: undefined as number | undefined,
    },
    roles: [] as { roleName: string; id: number }[],
    depts: [] as { deptName: string; id: number }[],
    total: 0,
    listDate: [] as Role[],
    listLoading: false,
  })

  const dialogVisible = ref(false) //编辑弹框
  const dialogVal = ref() //编辑相关值
  const selectRows = ref<Role[]>([])

  provide(tableSearch, () => {
    queryPage.query.pageNum = 1
    handleSearch()
  })

  async function getRoleAndDept() {
    try {
      const {
        data: { role, dept },
      } = await getAllRoleOrDeptApi()
      queryPage.roles = role || []
      queryPage.depts = dept || []
    } catch (error) {
      console.log(error)
    }
  }

  // 查询
  const handleSearch = async () => {
    queryPage.listLoading = true
    const { query, user } = queryPage
    try {
      const {
        data: { records, total },
      } = await getUserPageApi({
        ...query,
        ...user,
        deptIds: user.deptIds ? [user.deptIds] : [],
        roleIds: user.roleIds ? [user.roleIds] : [],
      })
      queryPage.listDate = records || []
      queryPage.total = total || 0
    } catch (error) {
      console.log(error)
    }
    queryPage.listLoading = false
  }

  // 重置
  const handleReset = () => {
    queryPage.user = {
      nickName: '',
      mail: '',
      phone: '',
      loginName: '',
      deptIds: undefined,
      roleIds: undefined,
      userType: undefined,
      createSt: undefined as number | undefined,
      createEd: undefined as number | undefined,
    }
    timeDate.value = undefined
  }
  // 获取表格序号
  const curIndex = computed(() => (queryPage.query.pageNum - 1) * queryPage.query.pageSize + 1)
  // 编辑
  const handleEdit = (val: Role) => {
    dialogVal.value = val
    closeEdit(true)
  }
  // 删除
  const handleDelete = (rows: Role[], deleteAll = false) => {
    $baseConfirm(!deleteAll ? '你确定要删除当前项吗' : '你确定要删除所有数据吗', null, async () => {
      const ids = rows.map((i) => i.id)
      const { msg } = await deleteUserApi({ ids: ids.toString(), deleteAll })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      queryPage.query.pageNum = 1
      handleSearch()
    })
  }
  const handleFreeUser = async (rows: Role) => {
    const { msg } = await freeUserApi(rows.id.toString())
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    queryPage.query.pageNum = 1
    handleSearch()
  }
  // 多选项改变
  const setSelectRows = (rows: Role[]) => {
    selectRows.value = rows
  }

  // 页容量改变
  const handleSizeChange = (val: number) => {
    queryPage.query.pageSize = val
    handleSearch()
  }

  // 页面改变
  const handleCurrentChange = () => {
    handleSearch()
  }

  // 关闭编辑弹框
  function closeEdit(visible: boolean) {
    dialogVisible.value = visible
  }

  function closeEditCallback(load: boolean) {
    closeEdit(false)
    if (load) handleSearch()
  }
  watchEffect(() => {
    if (!timeDate.value) {
      queryPage.user.createSt = undefined
      queryPage.user.createEd = undefined
      return
    }
    const [ST, END] = timeDate.value
    queryPage.user.createSt = new Date(ST).getTime()
    queryPage.user.createEd = new Date(END).getTime()
  })

  onMounted(async () => {
    handleSearch()
    getRoleAndDept()
  })
</script>
<script lang="ts">
  export default {
    name: 'Management',
  }
</script>

<template>
  <div class="user-management-container">
    <template v-if="isShow">
      <vab-query-form>
        <vab-query-form-left-panel :span="12">
          <h3>用户管理</h3>
        </vab-query-form-left-panel>
        <vab-query-form-right-panel :span="12">
          <el-button :icon="Plus" type="primary" @click="isShow = false">添加</el-button>
          <el-dropdown style="margin-left: 10px">
            <span class="el-dropdown-link">
              <el-button type="danger">
                批量删除
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleDelete(selectRows)">删除选中</el-dropdown-item>
                <el-dropdown-item @click="(e) => handleDelete([], true)">删除所有</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </vab-query-form-right-panel>
      </vab-query-form>
      <el-table
        v-loading="queryPage.listLoading"
        :data="queryPage.listDate"
        row-key="id"
        @cell-contextmenu="useTableCopy"
        @selection-change="setSelectRows"
      >
        <el-table-column show-overflow-tooltip type="selection" />
        <el-table-column align="center" :index="curIndex" label="序号" type="index" width="55" />
        <el-table-column
          :formatter="({ createTime }) => formatNstime(createTime, false)"
          label="创建时间"
          prop="createTime"
          show-overflow-tooltip
          width="200"
        >
          <template #header="{ column }">
            <table-filter-tool
              v-model:filter-value="timeDate"
              :column="column"
              filter-type="date-range"
              :tools="['filter']"
            />
          </template>
        </el-table-column>
        <el-table-column label="登录名称" prop="loginName" show-overflow-tooltip width="200">
          <template #header="{ column }">
            <table-filter-tool
              v-model:filter-value="queryPage.user.loginName"
              :column="column"
              filter-type="text"
              :tools="['filter']"
            />
          </template>
        </el-table-column>
        <el-table-column label="姓名" prop="nickName" show-overflow-tooltip>
          <template #header="{ column }">
            <table-filter-tool
              v-model:filter-value="queryPage.user.nickName"
              :column="column"
              filter-type="text"
              :tools="['filter']"
            />
          </template>
        </el-table-column>
        <el-table-column label="角色" prop="roleName" show-overflow-tooltip>
          <template #header="{ column }">
            <table-filter-tool
              v-model:filter-value="queryPage.user.roleIds"
              :column="column"
              :filter-option="queryPage.roles.map((item) => ({ label: item.roleName, value: item.id }))"
              filter-type="select"
              :tools="['filter']"
            />
          </template>
        </el-table-column>
        <el-table-column
          :formatter="({ userType }) => (userType === 1 ? '管理员' : userType === 2 ? '审计员' : '普通用户')"
          label="用户类型"
          prop="userType"
          show-overflow-tooltip
        >
          <template #header="{ column }">
            <table-filter-tool
              v-model:filter-value="queryPage.user.userType"
              :column="column"
              :filter-option="userTypes.map((item) => ({ label: item.name, value: item.lable }))"
              filter-type="select"
              :tools="['filter']"
            />
          </template>
        </el-table-column>
        <el-table-column label="部门" prop="deptName" show-overflow-tooltip>
          <template #header="{ column }">
            <table-filter-tool
              v-model:filter-value="queryPage.user.deptIds"
              :column="column"
              :filter-option="queryPage.depts.map((item) => ({ label: item.deptName, value: item.id }))"
              filter-type="select"
              :tools="['filter']"
            />
          </template>
        </el-table-column>
        <el-table-column label="邮箱" prop="mail" show-overflow-tooltip />
        <el-table-column label="手机" prop="phone" show-overflow-tooltip />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button v-if="row?.status === 1" size="small" @click="handleFreeUser(row)">释放</el-button>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleDelete([row])">删除</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty class="vab-data-empty" description="暂无数据" />
        </template>
      </el-table>
      <el-pagination
        v-model:current-page="queryPage.query.pageNum"
        v-model:page-size="queryPage.query.pageSize"
        background
        :layout="layout"
        :page-sizes="[10, 20, 50, 100]"
        :total="queryPage.total"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      />
    </template>
    <!-- 子页面添加用户 -->
    <add-user
      v-else
      :depts="queryPage.depts"
      :is-show="isShow"
      :roles="queryPage.roles"
      @on-cancel="
        (load) => {
          isShow = true
          if (load) handleSearch()
        }
      "
    />
    <!-- 子页面编辑用户 -->
    <edit-user
      v-if="dialogVisible"
      :depts="queryPage.depts"
      :dialog-val="dialogVal"
      :dialog-visible="dialogVisible"
      :roles="queryPage.roles"
      @on-close-event="closeEditCallback"
    />
  </div>
</template>

<style lang="scss" scoped>
  .user-management-container {
    :deep() {
      .el-scrollbar {
        height: calc(100vh - 180px) !important;
      }
    }
  }
</style>
