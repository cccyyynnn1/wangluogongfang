<script setup lang="ts">
  import { useTableCopy } from '@/utils'

  import DeviceList from './components/device-list.vue'

  import { networkListApi, adapterStartApi, adapterStopApi } from '@/api-ecs/equipment'
  import AesEncryptCBC from '~/src/utils/crypto'

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

  const listLoading = ref(false) // 是否加载

  // 启、停用
  const password = ref('')
  const handleChange = async (val: boolean, id: string) => {
    try {
      if (val) {
        const res = await adapterStartApi(id, { password: AesEncryptCBC(password.value) })
        ElMessage({ message: '启用成功', type: 'success' })
      } else {
        const res = await adapterStopApi(id, { password: AesEncryptCBC(password.value) })
        ElMessage({ message: '停用成功', type: 'success' })
      }
    } finally {
      getData()
      password.value = ''
    }
  }

  // 切换右侧菜单
  const changeMenu = (id: string) => {
    queryPage.query.flowDeviceId = id
    getData()
  }

  // 删除
  const handleDelete = () => {
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

  // 获取数据
  const getData = async () => {
    try {
      listLoading.value = true
      // @ts-ignore
      const { data } = await networkListApi({ ...queryPage, query: JSON.stringify(queryPage.query) })
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
            return true
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
</script>

<script lang="ts">
  export default {
    name: 'NetworkInformation', // 设备网卡信息
  }
</script>
<template>
  <div class="network-information-container">
    <!-- 顶部按钮 -->
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>设备网卡信息</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12" style="display: none">
        <el-button disabled type="primary" @click="handleChange(true, '0')">启用</el-button>
        <el-button disabled type="primary" @click="handleChange(false, '0')">停用</el-button>
        <el-button disabled style="margin-right: 0 !important" type="danger" @click="handleDelete">批量删除</el-button>
      </vab-query-form-right-panel>
    </vab-query-form>
    <el-row :gutter="20">
      <el-col :span="4"><device-list @on-change-menu="changeMenu" /></el-col>
      <el-col :span="20">
        <div>
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
            <el-table-column fixed="left" label="名称" prop="adapterName" show-overflow-tooltip width="100" />
            <el-table-column fixed="left" label="设备" prop="flowDeviceIdStr" show-overflow-tooltip width="100" />
            <el-table-column label="IPv4地址" prop="ipV4" show-overflow-tooltip width="100" />
            <el-table-column label="IPv6地址" prop="ipV6" show-overflow-tooltip width="100" />
            <el-table-column label="MAC地址" prop="mac" show-overflow-tooltip width="100" />
            <el-table-column label="掩码" prop="netmask" show-overflow-tooltip width="100" />
            <el-table-column label="广播地址" prop="broadaddr" show-overflow-tooltip width="100" />
            <el-table-column label="关联Agent" prop="关联Agent" show-overflow-tooltip width="100" />
            <el-table-column label="状态" prop="enabledStr" show-overflow-tooltip width="100" />
            <!-- <el-table-column label="启用抓包" prop="capture" show-overflow-tooltip width="160">
              <template #default="{ row }">
                <el-switch
                  v-model="row.capture"
                  :before-change="handlBeforeChange"
                  @change="handleChange(row.capture, row.id)"
                />
              </template>
            </el-table-column> -->
            <el-table-column label="描述" prop="描述" show-overflow-tooltip width="100" />
            <el-table-column label="创建时间" prop="crtTimeStr" show-overflow-tooltip width="100" />
            <el-table-column label="更新时间" prop="lastTimeStr" show-overflow-tooltip width="100" />
            <el-table-column label="是否与抓包线程匹配" prop="是否与抓包线程匹配" show-overflow-tooltip width="160" />
            <!-- <el-table-column fixed="right" label="操作" width="80">
              <template #default>
                <el-button size="small" @click="handleDelete">删除</el-button>
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
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped lang="scss">
  .network-information-container {
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
