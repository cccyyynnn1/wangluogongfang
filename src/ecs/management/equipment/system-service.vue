<script setup lang="ts">
  import { useTableCopy } from '@/utils'

  import { TableInstance } from 'element-plus'

  import { sysServiceListApi, sysServiceStopApi, sysServiceStartApi } from '@/api-ecs/equipment'
  import AesEncryptCBC from '~/src/utils/crypto'

  let listDate = reactive<any[]>([]) // 表格数据

  const layout = ref('total, sizes, prev, pager, next, jumper')

  const queryPage = reactive({
    page: 1,
    limit: 10,
    query: {},
  })

  const total = ref(0) // 总条数

  const listLoading = ref(false) // 是否加载

  const tableRef = ref<TableInstance>() // 表格实例

  // 导出
  const handleExport = () => {
    ElMessage({ message: '功能完善中...', type: 'warning' })
  }

  // 获取数据
  const getData = async () => {
    try {
      listLoading.value = true
      const { data } = await sysServiceListApi({ ...queryPage, query: JSON.stringify(queryPage.query) })
      listDate = data.list
      total.value = data.total
    } finally {
      listLoading.value = false
    }
  }

  const curIndex = computed(() => (queryPage.page - 1) * queryPage.limit + 1)

  // 启、停用
  const password = ref('')
  const handleChange = async (val: boolean, id: string) => {
    try {
      if (val) {
        const { msg } = await sysServiceStartApi(id, { password: AesEncryptCBC(password.value) })
        ElMessage({ message: msg, type: 'success' })
      } else {
        const { msg } = await sysServiceStopApi(id, { password: AesEncryptCBC(password.value) })
        ElMessage({ message: msg, type: 'success' })
      }
    } finally {
      // setTimeout(() => {
      getData()
      // }, 1000)
      password.value = ''
    }
  }

  const handlBeforeChange = () => {
    return ElMessageBox({
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
            done()
            return true
          } catch (error) {
            instance.confirmButtonLoading = false
            return false
          }
        } else {
          done()
          return false
        }
      },
    }).catch(() => {
      return false
    })
  }

  onMounted(() => {
    getData()
  })
</script>

<script lang="ts">
  export default {
    name: 'SystemService', // 系统服务
  }
</script>
<template>
  <div class="system-service-container">
    <!-- 顶部按钮 -->
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>系统服务</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12" style="display: none">
        <el-button disabled type="primary" @click="handleExport">启用</el-button>
        <el-button disabled style="margin-right: 0 !important" type="primary" @click="handleExport">停用</el-button>
      </vab-query-form-right-panel>
    </vab-query-form>

    <!-- 表格 -->
    <el-table ref="tableRef" v-loading="listLoading" :data="listDate" row-key="id" @cell-contextmenu="useTableCopy">
      <el-table-column show-overflow-tooltip type="selection" />
      <el-table-column fixed="left" label="序号" width="55">
        <template #default="{ $index }">
          {{ curIndex + $index }}
        </template>
      </el-table-column>
      <el-table-column fixed="left" label="服务ID" prop="id" show-overflow-tooltip width="230" />
      <el-table-column fixed="left" label="服务名称" prop="name" show-overflow-tooltip width="180" />
      <el-table-column label="类型" prop="typeStr" show-overflow-tooltip width="100" />
      <!-- <el-table-column label="状态" prop="enabled" show-overflow-tooltip width="160">
        <template #default="{ row }">
          <el-switch
            v-model="row.enabled"
            :before-change="handlBeforeChange"
            @change="handleChange(row.enabled, row.serviceId)"
          />
        </template>
      </el-table-column> -->
      <el-table-column label="创建线程数" prop="threadCount" show-overflow-tooltip width="100" />
      <el-table-column label="运行线程数" prop="runThreadCount" show-overflow-tooltip width="100" />
      <el-table-column label="待处理数据" prop="dataCount" show-overflow-tooltip width="100" />
      <el-table-column label="启动时间" prop="startTime" show-overflow-tooltip width="180" />
      <el-table-column label="停止时间" prop="stopTimeStr" show-overflow-tooltip width="180" />
      <el-table-column label="运行时长" prop="runTimeStr" show-overflow-tooltip width="100" />
      <el-table-column label="服务描述" prop="info" show-overflow-tooltip width="400" />
      <template #empty>
        <el-empty class="vab-data-empty" description="暂无数据" />
      </template>
    </el-table>
  </div>
</template>

<style scoped lang="scss">
  .system-service-container {
    h3 {
      margin-block: 0 0.5em;
    }
    :deep() {
      .el-scrollbar {
        height: calc(100vh - 140px) !important;
      }
    }
  }
</style>
