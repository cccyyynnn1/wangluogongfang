<script setup lang="ts">
  import CollectingDeviceEdit from './components/collecting-device-edit.vue'

  import { useTableCopy } from '@/utils'

  import { equipmentListApi, restartEquipmentApi } from '@/api-ecs/equipment'
  import AesEncryptCBC from '~/src/utils/crypto'

  let listDate = reactive<object[]>([]) // 表格数据

  const layout = ref('total, sizes, prev, pager, next, jumper')

  const queryPage = reactive({
    page: 1,
    limit: 10,
    query: {},
  })

  const total = ref(0) // 总条数

  const listLoading = ref(false) // 是否加载

  const showPage = ref(false) // 是否显示编辑页

  const currentRow = ref()

  const moreOptions = [
    { value: 'analyse', label: '分析' },
    { value: 'restart', label: '重启服务' },
    { value: 'inquiry', label: '信息查询' },
  ]

  // 导出
  const handleExport = () => {}

  // 操作下拉框改变
  const handleChange = async (e: string, row: { id: string }) => {
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
            let res = undefined
            switch (e) {
              case 'analyse':
                ElMessage({ message: '功能完善中...', type: 'warning' })
                break
              case 'restart':
                res = await restartEquipmentApi(row.id, { password: AesEncryptCBC(password.value) })
                // @ts-ignore
                ElMessage({ message: res.msg, type: 'success' })
                break

              case 'inquiry':
                ElMessage({ message: '功能完善中...', type: 'warning' })
                break
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

  // 删除
  const handleDelete = () => {
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
  }

  // 编辑
  const handleEdit = (row: any) => {
    currentRow.value = row
    showPage.value = true
  }

  // 获取数据
  const getData = async () => {
    try {
      listLoading.value = true
      const { data } = await equipmentListApi({ ...queryPage, query: JSON.stringify(queryPage.query) })
      listDate = data.list
      total.value = data.total
    } finally {
      listLoading.value = false
    }
  }

  // 多选项改变
  const setSelectRows = (e: any) => {}

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
    name: 'CollectingDevice', // 采集设备
  }
</script>
<template>
  <div class="collecting-device-container">
    <!-- 顶部按钮 -->
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>采集设备</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12" style="display: none">
        <el-button disabled type="primary" @click="handleExport">导出</el-button>
        <el-button disabled type="danger" @click="handleDelete">批量删除</el-button>
      </vab-query-form-right-panel>
    </vab-query-form>

    <!-- 表格 -->
    <el-table
      v-loading="listLoading"
      :data="listDate"
      row-key="id"
      @cell-contextmenu="useTableCopy"
      @selection-change="setSelectRows"
    >
      <el-table-column show-overflow-tooltip type="selection" />
      <el-table-column fixed="left" label="序号" width="55">
        <template #default="{ $index }">
          {{ curIndex + $index }}
        </template>
      </el-table-column>
      <el-table-column fixed="left" label="名称" prop="name" show-overflow-tooltip width="130" />
      <el-table-column fixed="left" label="当前IP" prop="ip" show-overflow-tooltip width="130" />
      <el-table-column label="序号" prop="sn" show-overflow-tooltip width="100" />
      <el-table-column label="版本" prop="version" show-overflow-tooltip width="100" />
      <el-table-column label="版本标识" prop="flag" show-overflow-tooltip width="100" />
      <el-table-column label="状态" prop="statusStr" show-overflow-tooltip width="100" />
      <el-table-column label="启动时间" prop="startTimeStr" show-overflow-tooltip width="100" />
      <el-table-column label="启动压缩存储" prop="zipStorageStr" show-overflow-tooltip width="110" />
      <el-table-column label="日志级别" prop="logLevelStr" show-overflow-tooltip width="100" />
      <el-table-column label="事件统计间隔" prop="eventStatSecondStr" show-overflow-tooltip width="110" />
      <el-table-column label="数据传输合并间隔" prop="dataTransactionSecondStr" show-overflow-tooltip width="140" />
      <el-table-column label="数据库内存限制" prop="dbMemoryLimitGbStr" show-overflow-tooltip width="130" />
      <el-table-column label="是否已激活" prop="是否已激活" show-overflow-tooltip width="100" />
      <el-table-column label="服务到期时间" prop="服务到期时间" show-overflow-tooltip width="110" />
      <el-table-column label="使用到期时间" prop="使用到期时间" show-overflow-tooltip width="110" />
      <el-table-column label="授权链路数" prop="授权链路数" show-overflow-tooltip width="100" />
      <el-table-column label="授权会话数" prop="授权会话数" show-overflow-tooltip width="100" />
      <el-table-column label="授权流量大小" prop="授权流量大小" show-overflow-tooltip width="110" />
      <el-table-column label="是否需要重启服务" prop="needRestartServiceStr" show-overflow-tooltip width="140" />
      <el-table-column label="操作系统" prop="system" show-overflow-tooltip width="100" />
      <el-table-column label="架构" prop="arch" show-overflow-tooltip width="100" />
      <el-table-column label="CPU个数" prop="cpuCount" show-overflow-tooltip width="100" />
      <el-table-column label="CPU核心数" prop="cpuCores" show-overflow-tooltip width="100" />
      <el-table-column label="CPU主频" prop="cpuMhzStr" show-overflow-tooltip width="100" />
      <el-table-column label="内存总量" prop="memoryStr" show-overflow-tooltip width="100" />
      <el-table-column label="硬盘总量" prop="diskStr" show-overflow-tooltip width="100" />
      <el-table-column label="创建时间" prop="crtTimeStr" show-overflow-tooltip width="100" />
      <el-table-column label="更新时间" prop="lastTimeStr" show-overflow-tooltip width="100" />
      <el-table-column label="同步系统时间" prop="systemTimeStr" show-overflow-tooltip width="110" />
      <el-table-column label="是否开启numa" prop="isNumaStr" show-overflow-tooltip width="120" />
      <!-- <el-table-column fixed="right" label="操作" width="220">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" @click="handleDelete">删除</el-button>
          <el-select
            placeholder="更多"
            size="small"
            style="width: 60px; margin-left: 20px"
            @change="handleChange($event, row)"
          >
            <el-option v-for="item in moreOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
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
    <CollectingDeviceEdit
      v-if="showPage"
      :current-row="currentRow"
      :show-page="showPage"
      @on-close-event="showPage = false"
      @on-reflash="getData"
    />
  </div>
</template>

<style scoped lang="scss">
  .collecting-device-container {
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
