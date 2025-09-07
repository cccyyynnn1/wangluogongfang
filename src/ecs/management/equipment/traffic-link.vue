<script setup lang="ts">
  import TrafficAddEdit from './components/traffic-add-edit.vue'

  import DeviceList from './components/device-list.vue'

  import { useTableCopy } from '@/utils'

  import { TableInstance } from 'element-plus'

  import { flowProbeDeleteApi, flowProbeListApi } from '@/api-ecs/equipment'
  import AesEncryptCBC from '~/src/utils/crypto'

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  let listDate = reactive<object[]>([]) // 表格数据

  const layout = ref('total, sizes, prev, pager, next, jumper')

  const queryPage = reactive<{ page: number; limit: number; query: { flowDeviceId?: string } }>({
    page: 1,
    limit: 10,
    query: {
      flowDeviceId: undefined,
    },
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
                await flowProbeDeleteApi(delStr, { password: AesEncryptCBC(password.value) })
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

  // 编辑
  const handleEdit = (row: any) => {
    mode.value = 'edit'
    currentRow.value = row
    showPage.value = true
  }

  // 获取数据
  const getData = async () => {
    try {
      listLoading.value = true
      // @ts-ignore
      const { data } = await flowProbeListApi({ ...queryPage, query: JSON.stringify(queryPage.query) })
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

  // 切换右侧菜单
  const changeMenu = (id: string) => {
    queryPage.query.flowDeviceId = id
    getData()
  }

  onMounted(() => {
    getData()
  })
</script>

<script lang="ts">
  export default {
    name: 'TrafficLink',
  }
</script>
<template>
  <div class="traffic-link-container">
    <!-- 顶部按钮 -->
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>流量链路</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12" style="display: none">
        <el-button type="primary" @click="handleAdd">添加</el-button>
        <el-button disabled type="primary" @click="handleExport">启用</el-button>
        <el-button disabled type="primary" @click="handleExport">停用</el-button>
        <el-button disabled type="primary" @click="handleExport">导入</el-button>
        <el-button disabled type="primary" @click="handleExport">导出</el-button>
        <el-button style="margin-right: 0 !important" type="danger" @click="handleDelete">批量删除</el-button>
      </vab-query-form-right-panel>
    </vab-query-form>
    <el-row :gutter="20">
      <el-col :span="4"><device-list @on-change-menu="changeMenu" /></el-col>
      <el-col :span="20">
        <div>
          <!-- 表格 -->
          <div style="width: 100%">
            <el-table
              ref="tableRef"
              v-loading="listLoading"
              :data="listDate"
              row-key="id"
              @cell-contextmenu="useTableCopy"
            >
              <el-table-column show-overflow-tooltip type="selection" />
              <el-table-column fixed="left" label="序号" width="55">
                <template #default="{ $index }">
                  {{ curIndex + $index }}
                </template>
              </el-table-column>
              <el-table-column fixed="left" label="名称" prop="name" show-overflow-tooltip width="100" />
              <el-table-column fixed="left" label="设备" prop="flowDeviceIdStr" show-overflow-tooltip width="100" />
              <el-table-column fixed="left" label="网卡" prop="adapterId" show-overflow-tooltip width="100" />
              <el-table-column label="类型" prop="typeStr" show-overflow-tooltip width="100" />
              <el-table-column label="链路ID" prop="probeId" show-overflow-tooltip width="100" />
              <el-table-column label="状态" prop="enabledStr" show-overflow-tooltip width="100" />
              <el-table-column label="识别IP层级" prop="ipLayerStr" show-overflow-tooltip width="100" />
              <el-table-column label="IP地址" prop="IP地址" show-overflow-tooltip width="100" />
              <el-table-column label="MAC地址" prop="MAC地址" show-overflow-tooltip width="100" />
              <el-table-column label="内网地址" prop="内网地址" show-overflow-tooltip width="100" />
              <el-table-column label="标识类型" prop="标识类型" show-overflow-tooltip width="100" />
              <el-table-column label="标识值" prop="标识值" show-overflow-tooltip width="100" />
              <el-table-column label="启用存储" prop="enableStorageStr" show-overflow-tooltip width="100" />
              <el-table-column label="启用分析" prop="enableAnalysisStr" show-overflow-tooltip width="100" />
              <el-table-column label="存储载包" prop="存储载包" show-overflow-tooltip width="100" />
              <el-table-column label="SSL自动载包" prop="sslPackageCutStr" show-overflow-tooltip width="120" />
              <el-table-column label="去重数据包" prop="removeDupPackageStr" show-overflow-tooltip width="100" />
              <el-table-column label="开启分析" prop="openAnalysisStr" show-overflow-tooltip width="100" />
              <el-table-column label="关闭存储" prop="closeStorageStr" show-overflow-tooltip width="100" />
              <el-table-column label="总宽带" prop="总宽带" show-overflow-tooltip width="100" />
              <el-table-column label="进网宽带" prop="进网宽带" show-overflow-tooltip width="100" />
              <el-table-column label="出网宽带" prop="出网宽带" show-overflow-tooltip width="100" />
              <el-table-column label="流量方向" prop="flowDirectionStr" show-overflow-tooltip width="100" />
              <el-table-column label="创建人" prop="crtUserStr" show-overflow-tooltip width="100" />
              <el-table-column label="创建时间" prop="crtTimeStr" show-overflow-tooltip width="100" />
              <el-table-column label="修改人" prop="updateUserStr" show-overflow-tooltip width="100" />
              <el-table-column label="修改时间" prop="updateTimeStr" show-overflow-tooltip width="100" />
              <!-- <el-table-column fixed="right" label="操作" width="160">
                <template #default="{ row }">
                  <el-button size="small" @click="handleEdit(row)">编辑</el-button>
                  <el-button size="small" @click="handleDelete(row)">删除</el-button>
                </template>
              </el-table-column> -->
              <template #empty>
                <el-empty class="vab-data-empty" description="暂无数据" />
              </template>
            </el-table>
          </div>
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
          <traffic-add-edit
            v-if="showPage"
            :current-row="currentRow"
            :mode="mode"
            :show-page="showPage"
            @on-close-event="showPage = false"
            @on-reflash="getData"
          />
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped lang="scss">
  .traffic-link-container {
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
