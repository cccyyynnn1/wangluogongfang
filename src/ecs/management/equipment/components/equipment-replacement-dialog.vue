<script setup lang="ts">
  import { flowDeviceListApi, flowDeviceUpdateApi } from '@/api-ecs/equipment'

  import { useTableCopy } from '@/utils'

  import { TableInstance } from 'element-plus'
  import AesEncryptCBC from '~/src/utils/crypto'

  const props = defineProps<{
    showPage: boolean
    currentRow: { id: string }
  }>()

  const visible = ref(false) // 弹框显隐

  const tableRef = ref<TableInstance>() // 表格实例

  const title = ref('添加') // 标题

  const layout = ref('total, sizes, prev, pager, next, jumper')

  const total = ref(0) // 总条数

  let listDate = reactive<object[]>([]) // 表格数据

  const listLoading = ref(false) // 是否加载

  const flowDeviceIds = ref('')

  // 获取型号列表参数
  const queryPage = reactive({
    page: 1,
    limit: 10,
    query: { flag: 1, status: '0' },
  })

  const emit = defineEmits<{
    (e: 'on-close-event'): void
    (e: 'on-reflash'): void
  }>()

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

  // 关闭回调
  const handleClose = () => {
    emit('on-close-event')
  }

  // 提交表单
  const submit = async () => {
    try {
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
              const id = props.currentRow.id
              listLoading.value = true
              const { message } = await flowDeviceUpdateApi(
                { id, flowDeviceIds: flowDeviceIds.value },
                { password: AesEncryptCBC(password.value) }
              )
              ElMessage({ message: message, type: 'success' })
              emit('on-reflash')
              handleClose()
              done()
            } catch (error) {
              instance.confirmButtonLoading = false
            }
          } else {
            done()
          }
        },
      }).catch(() => {})
    } finally {
      listLoading.value = false
    }
  }

  // 初始化数据
  const initData = () => {
    visible.value = props.showPage
    getData()
  }

  // 获取表格数据
  const getData = async () => {
    try {
      listLoading.value = true
      const { data } = await flowDeviceListApi({ ...queryPage, query: JSON.stringify(queryPage.query) })
      listDate = data.list
      total.value = data.total
    } finally {
      listLoading.value = false
    }
  }

  // 每一行点击事件
  const handleCellClick = (row: { id: string }) => {
    flowDeviceIds.value = row.id
    tableRef.value!.setCurrentRow(null)
    tableRef.value!.setCurrentRow(row)
  }

  onMounted(() => {
    initData()
  })
</script>

<script lang="ts">
  export default {
    name: 'EquipmentReplacementDialog',
  }
</script>
<template>
  <div class="equipment-replacement-dialog">
    <el-dialog v-model="visible" :before-close="handleClose" title="流量分析设备" width="800">
      <el-table
        ref="tableRef"
        v-loading="listLoading"
        :border="true"
        :data="listDate"
        highlight-current-row
        row-key="id"
        @cell-click="handleCellClick"
        @cell-contextmenu="useTableCopy"
      >
        <!-- <el-table-column show-overflow-tooltip type="selection" /> -->
        <el-table-column :align="'center'" label="序号" width="55">
          <template #default="{ $index }">
            {{ curIndex + $index }}
          </template>
        </el-table-column>
        <el-table-column :align="'center'" label="名称" prop="name" show-overflow-tooltip width="307" />
        <el-table-column :align="'center'" label="当前IP" prop="customObj" show-overflow-tooltip width="200" />
        <el-table-column :align="'center'" label="版本号" prop="version" show-overflow-tooltip width="200" />
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
      <el-row style="margin-top: 20px">
        <el-col :offset="7" :span="8" style="text-align: right">
          <el-button :loading="listLoading" type="primary" @click="submit">
            {{ listLoading ? '更新中' : '保存' }}
          </el-button>
          <el-button @click="handleClose">取消</el-button>
        </el-col>
      </el-row>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  :deep() {
    .el-scrollbar {
      height: 300px !important;
    }
  }
</style>
