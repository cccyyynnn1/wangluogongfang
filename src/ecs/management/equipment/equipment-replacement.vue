<script setup lang="ts">
  import EquipmentReplacementDialog from './components/equipment-replacement-dialog.vue'

  import { useTableCopy } from '@/utils'

  import { TableInstance } from 'element-plus'

  import { flowDeviceUpdateDeleteApi, flowDeviceUpdateListApi } from '@/api-ecs/equipment'

  // const $baseConfirm: any = inject('$baseConfirm')

  // const $baseMessage: any = inject('$baseMessage')

  let listDate = reactive<object[]>([]) // 表格数据

  const showPage = ref(false) // 是否显示编辑页

  const currentRow = ref()

  const layout = ref('total, sizes, prev, pager, next, jumper')

  const queryPage = reactive({
    page: 1,
    limit: 20,
    query: {},
  })

  const total = ref(0) // 总条数

  // let delStr = ''

  const listLoading = ref(false) // 是否加载

  const tableRef = ref<TableInstance>() // 表格实例

  // 更新
  const handleUpdate = (row: any) => {
    currentRow.value = row
    showPage.value = true
  }

  // 上传
  const handleUpload = () => {
    ElMessage({ message: '功能完善中...', type: 'warning' })
  }

  // 删除
  const handleDelete = (row?: any) => {
    const password = ref('')
    ElMessageBox({
      title: '提示',
      showCancelButton: true,
      customClass: 'need-password-message-box',
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      customStyle: {
        maxWidth: '500px',
      },
      message: () =>
        h('div', null, [
          h('div', { style: 'margin: 15px 0 5px 0;color:#55585b' }, '请输入敏感操作密码:(通过验证后方可进行操作)'),
          h(ElInput, {
            type: 'password',
            modelValue: password.value,
            placeholder: '请输入敏感操作密码',
            showPassword: true,
            style: 'margin-block: 10px',
            'onUpdate:modelValue': (val: string) => {
              password.value = val
            },
          }),
        ]),

      beforeClose: async (action, instance, done) => {
        if (action === 'confirm') {
          instance.confirmButtonLoading = true
          try {
            ElMessage({ message: '功能完善中...', type: 'warning' })
            done()
          } catch (error) {
            instance.confirmButtonLoading = false
          }
        } else {
          done()
        }
      },
    }).catch(() => {})
    // if (row.id) {
    //   delStr = row.id
    // } else {
    //   const res = tableRef.value?.getSelectionRows()
    //   const arr: string[] = []
    //   res.forEach((item: any) => {
    //     arr.push(item.id)
    //   })
    //   delStr = arr.join(',')
    // }
    // if (delStr) {
    //   $baseConfirm('你确定要删除吗', null, async () => {
    //     await flowDeviceUpdateDeleteApi(delStr)
    //     $baseMessage('删除成功', 'success', 'vab-hey-message-success')
    //     await getData()
    //   })
    // }
  }

  // 获取数据
  const getData = async () => {
    try {
      listLoading.value = true
      // @ts-ignore
      const { data } = await flowDeviceUpdateListApi({ ...queryPage, query: JSON.stringify(queryPage.query) })
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

  onMounted(() => {
    getData()
  })
</script>

<script lang="ts">
  export default {
    name: 'EquipmentReplacement', // 流量设备更新
  }
</script>
<template>
  <div class="equipment-replacement-container">
    <!-- 顶部按钮 -->
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>流量设备更新</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12" style="display: none">
        <el-button disabled type="primary" @click="handleUpload">上传文件</el-button>
        <el-button type="danger" @click="handleDelete">批量删除</el-button>
      </vab-query-form-right-panel>
    </vab-query-form>

    <!-- 表格 -->
    <el-table
      ref="tableRef"
      v-loading="listLoading"
      class="my-table"
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
      <el-table-column label="文件名" prop="fileName" show-overflow-tooltip width="180" />
      <el-table-column label="文件大小" prop="fileSizeStr" show-overflow-tooltip width="100" />
      <el-table-column label="版本号" prop="version" show-overflow-tooltip width="100" />
      <el-table-column label="版本类型" prop="type" show-overflow-tooltip width="100" />
      <el-table-column label="版本标识" prop="flag" show-overflow-tooltip width="100" />
      <el-table-column label="打包时间" prop="pkgTimeStr" show-overflow-tooltip width="160" />
      <el-table-column label="指纹" prop="fingerPoint" show-overflow-tooltip width="280" />
      <el-table-column label="创建人" prop="crtUserStr" show-overflow-tooltip width="100" />
      <el-table-column label="创建时间" prop="crtTimeStr" show-overflow-tooltip width="160" />
      <el-table-column label="更新人" prop="updateUserStr" show-overflow-tooltip width="100" />
      <el-table-column label="更新时间" prop="updateTimeStr" show-overflow-tooltip width="160" />
      <!-- <el-table-column fixed="right" label="操作" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="handleUpdate(row)">更新</el-button>
          <el-button size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column> -->
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
    <equipment-replacement-dialog
      v-if="showPage"
      :current-row="currentRow"
      :show-page="showPage"
      @on-close-event="showPage = false"
      @on-reflash="getData"
    />
  </div>
</template>

<style scoped lang="scss">
  .equipment-replacement-container {
    h3 {
      margin-block: 0 0.5em;
    }
    .my-table {
      :deep() {
        .el-scrollbar {
          height: calc(100vh - 180px) !important;
        }
      }
    }
  }
</style>
