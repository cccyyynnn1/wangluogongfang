<script lang="ts">
  export default {
    name: 'DataSwitching',
  }
</script>
<script setup lang="ts">
  import { useTableCopy } from '@/utils'

  import AddEdit from './add-edit.vue'

  import { Delete, Plus } from '@element-plus/icons-vue'

  import { transpondGetGageType } from '~/src/types'

  import { transpondDeleteApi, transpondGetPageApi } from '~/src/api-ecs/transpond'

  import { formatNstime } from '@/utils/time'

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  const listLoading = ref(false) // 是否加载

  const showEdit = ref(false) // 是否显示编辑页

  const mode = ref('add') // 是编辑还是添加

  // 表格数据
  let listDate = reactive<object[]>([])

  const layout = ref('total, sizes, prev, pager, next, jumper')

  let delList: [] = [] // 删除的数组

  const queryPage = reactive<transpondGetGageType>({
    pageSize: 10,
    pageNum: 1,
  })

  // 获取表格序号
  const curIndex = computed(() => (queryPage.pageNum - 1) * queryPage.pageSize + 1)

  const total = ref(0) // 总条数

  const currentData = ref() // 当前行数据

  // 添加
  const handleAdd = () => {
    showEdit.value = true
    mode.value = 'add'
    changeShowEdit(true)
  }

  onMounted(() => {
    getData()
  })

  // 获取数据
  const getData = async () => {
    listLoading.value = true
    const { data } = await transpondGetPageApi({ ...queryPage })
    listDate = data.records
    total.value = data.total
    listLoading.value = false
  }

  // 编辑
  const handleEdit = (row: any) => {
    currentData.value = row
    mode.value = 'edit'
    changeShowEdit(true)
  }

  // 显隐编辑页
  const changeShowEdit = (val: boolean) => {
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
        const { msg } = await transpondDeleteApi({ ids: delList.toString(), deleteAll: false })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        delList = []
        await getData()
      })
    }
    if (!row.row && row.deleteAll) {
      $baseConfirm('你确定要删除所有数据吗', null, async () => {
        const { msg } = await transpondDeleteApi({ ids: '', deleteAll: true })
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
  const handleSizeChange = () => {}

  // 页面改变
  const handleCurrentChange = () => {}
</script>

<template>
  <div class="data-switching-container">
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <el-button :icon="Plus" type="primary" @click="handleAdd()">添加</el-button>
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
      <el-table-column :align="'center'" label="名称" prop="name" show-overflow-tooltip />
      <el-table-column :align="'center'" label="描述" prop="description" show-overflow-tooltip />
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
      v-if="showEdit"
      :current-data="currentData"
      :mode="mode"
      :show-edit="showEdit"
      @on-closeEvent="changeShowEdit"
      @on-reflash="getData"
    />
  </div>
</template>

<style scoped lang="scss">
  :deep() {
    .el-scrollbar {
      height: calc(100vh - 270px) !important;
    }
  }
</style>
