<script setup lang="ts">
  import LoadBalancingAdd from './components/load-balancing-add.vue'

  import { useTableCopy } from '@/utils'

  import { TableInstance } from 'element-plus'

  import { loadBalanceDeleteApi, loadBalanceApi } from '@/api-ecs/equipment'
  import AesEncryptCBC from '~/src/utils/crypto'

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  let listDate = reactive<object[]>([]) // 表格数据

  const layout = ref('total, sizes, prev, pager, next, jumper')

  const queryPage = reactive({
    page: 1,
    limit: 10,
    query: {},
  })

  const total = ref(0) // 总条数

  let delStr = ''

  const listLoading = ref(false) // 是否加载

  const tableRef = ref<TableInstance>() // 表格实例

  const showPage = ref(false) // 是否显示编辑页

  const currentRow = ref()

  const mode = ref('')

  // 添加
  const handleAdd = () => {
    mode.value = 'add'
    showPage.value = true
  }

  // 导出
  const handleExport = () => {
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
            if (row.id) {
              delStr = row.id
            } else {
              const res = tableRef.value?.getSelectionRows()
              const arr: string[] = []
              res.forEach((item: any) => {
                arr.push(item.id)
              })
              delStr = arr.join(',')
            }
            if (delStr) {
              $baseConfirm('你确定要删除吗', null, async () => {
                await loadBalanceDeleteApi(delStr, { password: AesEncryptCBC(password.value) })
                $baseMessage('删除成功', 'success', 'vab-hey-message-success')
                await getData()
              })
            }
            done()
          } catch (error) {
            instance.confirmButtonLoading = false
          }
        } else {
          done()
        }
      },
    }).catch(() => {})
  }

  // 获取数据
  const getData = async () => {
    try {
      listLoading.value = true
      const { data } = await loadBalanceApi({ ...queryPage, query: JSON.stringify(queryPage.query) })
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
    name: 'LoadBalancing', // 负载均衡管理,
  }
</script>
<template>
  <div class="load-balancing-container">
    <!-- 顶部按钮 -->
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>负载均衡管理</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12" style="display: none">
        <el-button type="primary" @click="handleAdd">添加</el-button>
        <el-button disabled type="primary" @click="handleExport">导入</el-button>
        <el-button disabled type="primary" @click="handleExport">导出</el-button>
        <el-button click="handleDelete" style="margin-right: 0 !important" type="danger">批量删除</el-button>
      </vab-query-form-right-panel>
    </vab-query-form>

    <!-- 表格 -->
    <el-table ref="tableRef" v-loading="listLoading" :data="listDate" row-key="id" @cell-contextmenu="useTableCopy">
      <el-table-column show-overflow-tooltip type="selection" />
      <el-table-column label="序号" width="55">
        <template #default="{ $index }">
          {{ curIndex + $index }}
        </template>
      </el-table-column>
      <el-table-column label="型号" prop="modelIdStr" show-overflow-tooltip width="100" />
      <el-table-column label="获取方式" prop="modeStr" show-overflow-tooltip width="100" />
      <el-table-column label="状态" prop="enabledStr" show-overflow-tooltip width="100" />
      <el-table-column label="名称" prop="name" show-overflow-tooltip width="100" />
      <el-table-column label="主机名" prop="hostName" show-overflow-tooltip width="100" />
      <el-table-column label="硬件型号" prop="hadrwareModel" show-overflow-tooltip width="100" />
      <el-table-column label="管理端口IP" prop="managementIp" show-overflow-tooltip width="100" />
      <el-table-column label="硬件序列号" prop="hardwareSerialNum" show-overflow-tooltip width="100" />
      <el-table-column label="激活秘钥" prop="activeKey" show-overflow-tooltip width="100" />
      <el-table-column label="激活日期" prop="activeDateStr" show-overflow-tooltip width="100" />
      <el-table-column label="连接超时" prop="连接超时" show-overflow-tooltip width="100" />
      <el-table-column label="查询超时" prop="查询超时" show-overflow-tooltip width="100" />
      <el-table-column label="负载均衡端口" prop="负载均衡端口" show-overflow-tooltip width="120" />
      <el-table-column label="采集类型" prop="采集类型" show-overflow-tooltip width="100" />
      <el-table-column label="采集周期" prop="采集周期" show-overflow-tooltip width="100" />
      <el-table-column label="采集开始时间" prop="采集开始时间" show-overflow-tooltip width="120" />
      <el-table-column label="最后采集时间" prop="最后采集时间" show-overflow-tooltip width="120" />
      <el-table-column label="采集次数" prop="采集次数" show-overflow-tooltip width="100" />
      <el-table-column label="备注" prop="备注" show-overflow-tooltip width="100" />
      <el-table-column label="创建人" prop="crtUserStr" show-overflow-tooltip width="100" />
      <el-table-column label="创建时间" prop="crtTimeStr" show-overflow-tooltip width="100" />
      <el-table-column label="修改人" prop="updateUserStr" show-overflow-tooltip width="100" />
      <el-table-column label="修改时间" prop="updateTimeStr" show-overflow-tooltip width="100" />
      <!-- <el-table-column fixed="right" label="操作" width="100">
        <template #default="{ row }">
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
    <!-- 编辑页 -->
    <load-balancing-add
      v-if="showPage"
      :current-row="currentRow"
      :mode="mode"
      :show-page="showPage"
      @on-close-event="showPage = false"
      @on-reflash="getData"
    />
  </div>
</template>

<style scoped lang="scss">
  .load-balancing-container {
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
