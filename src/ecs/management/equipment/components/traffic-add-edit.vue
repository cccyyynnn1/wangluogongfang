<script setup lang="ts">
  import { FormInstance, TableInstance } from 'element-plus'

  import { requireRules } from '~/src/utils/rules'

  import { getAdapterByIdApi, flowProbeAddApi, flowProbeUpdataApi } from '@/api-ecs/equipment'

  import { Search } from '@element-plus/icons-vue'

  import { useTableCopy } from '@/utils'
  import { pa } from 'element-plus/es/locale'
  import AesEncryptCBC from '~/src/utils/crypto'

  const props = defineProps<{
    mode: string
    showPage: boolean
    currentRow: object
  }>()

  const formRef = ref<FormInstance>() // 表单实例

  const tableRef = ref<TableInstance>() // 表格实例

  const visible = ref(false) // 弹框显隐

  const popVisible = ref(false) // popover弹框显隐

  const title = ref('添加') // 标题

  let listDate = reactive<object[]>([]) // 表格数据

  const layout = ref('total, sizes, prev, pager, next, jumper')

  const totals = ref(0) // 总条数

  const listLoading = ref(false) // 是否加载

  // 获取网卡Id参数
  const adapterPayload = reactive({
    page: 1,
    limit: 10,
    query: {
      keyWord: undefined,
      flowDeviceId: 'a64e8aaa-512c-4801-af82-184c5a3d5513',
      flowDeviceIdStr: 'localhost.localdomain',
    },
    url: 'adapter/page',
  })

  // 表单数据
  const formData = reactive({
    dialogTitle: 'add', // 添加链路时候的固定值
    moduleType: 'flowProbe',
    name: undefined,
    type: '2',
    adapterId: [],
    adapterIdStr: '',
    enabled: true,
    ipLayer: '0',
    ipAddr: undefined,
    macAddr: undefined,
    intranetIp: undefined,
    subType: undefined,
    subValue: undefined,
    enableStorage: true,
    storagePackageCut: undefined,
    enableAnalysis: true,
    analysisPackageCut: undefined,
    sslPackageCut: false,
    removeDupPackage: '0',
    openAnalysis: undefined,
    closeStorage: undefined,
    totalBandwidth: undefined,
    inBandwidth: undefined,
    outBandwidth: undefined,
    flowDirection: '1',
  })

  // 表单数据校验规则
  const rules = reactive({
    name: requireRules,
    logLevelStr: requireRules,
    type: requireRules,
    enabled: requireRules,
    ipLayer: requireRules,
    enableStorage: requireRules,
    enableAnalysis: requireRules,
    sslPackageCut: requireRules,
    removeDupPackage: requireRules,
    // totalBandwidth: requireRules,
    flowDirection: requireRules,
  })

  // 识别IP层级option
  const option = [
    { value: '0', label: '顶层' },
    { value: '1', label: '第一层' },
    { value: '2', label: '第二层' },
    { value: '3', label: '第三层' },
    { value: '4', label: '第四层' },
  ]

  // 识别类型option
  const typeOption = [
    { value: '0', label: 'vlan' },
    { value: '1', label: 'vxlan' },
  ]

  // 开启分析option
  const analyseOption = [
    { value: 'dnsLog', label: 'DNS日志' },
    { value: 'ftpLog', label: 'FTP日志' },
    { value: 'httpsCaLog', label: 'HTTPS证书' },
    { value: 'httpFileLog', label: 'HTTP文件日志' },
    { value: 'httpLog', label: 'HTTP日志' },
    { value: 'httpLoginLog', label: 'HTTP登录日志' },
    { value: 'icmpLog', label: 'ICMP日志' },
    { value: 'msrdpLog', label: 'MSRDP日志' },
    { value: 'smbLog', label: 'SMB日志' },
    { value: 'socksLog', label: 'SOCKS日志' },
    { value: 'sshLog', label: 'SSH日志' },
    { value: 'sslLog', label: 'SSL日志' },
    { value: 'telnetLog', label: 'TELNET日志' },
    { value: 'telnetLoginLog', label: 'TELNET登录日志' },
    { value: 'tftpLog', label: 'TFTP日志' },
    { value: 'sqlLog', label: '异常流量统计' },
    { value: 'emailLog', label: '数据库日志' },
    { value: 'abnormalFlowStat', label: '邮件日志' },
  ]

  // 去重数据包option
  const deWeightOption = [
    { value: '0', label: '不去重' },
    { value: '1', label: '按SEQ去重' },
    { value: '2', label: '按IP去重' },
  ]

  // 关闭存储option
  const storageOption = [
    { value: 'dnsLog', label: 'DNS日志' },
    { value: 'ftpLog', label: 'FTP日志' },
    { value: 'httpFileLog', label: 'HTTP文件日志' },
    { value: 'httpLog', label: 'HTTP日志' },
    { value: 'httpLoginLog', label: 'HTTP登录日志' },
    { value: 'icmpLog', label: 'ICMP日志' },
    { value: 'msrdpLog', label: 'MSRDP日志' },
    { value: 'smbLog', label: 'SMB日志' },
    { value: 'socksLog', label: 'SOCKS日志' },
    { value: 'sshLog', label: 'SSH日志' },
    { value: 'sslLog', label: 'SSL日志' },
    { value: 'telnetLog', label: 'TELNET日志' },
    { value: 'telnetLoginLog', label: 'TELNET登录日志' },
    { value: 'tftpLog', label: 'TFTP日志' },
    { value: 'packetLossStat', label: '丢包统计' },
    { value: 'eventStat', label: '事件统计表' },
    { value: 'flowAlarmLog', label: '告警日志' },
    { value: 'appStatMin', label: '应用统计分表' },
    { value: 'appStatSec', label: '应用统计秒表' },
    { value: 'abnormalFlowStat', label: '异常流量统计' },
    { value: 'dataTransaction', label: '数据传输' },
    { value: 'sqlLog', label: '数据库日志' },
    { value: 'netSegmentStatMin', label: '网段统计分表' },
    { value: 'netSegmentStatSec', label: '网段统计秒表' },
    { value: 'netPerformanceStat', label: '网段性能统计表' },
    { value: 'connectionClose', label: '连接关闭' },
    { value: 'connectionSetup', label: '连接建立' },
    { value: 'emailLog', label: '邮件日志' },
    { value: 'probeStatMin', label: '链路统计分表' },
    { value: 'probeStatSec', label: '链路统计秒表' },
  ]

  // 流量方向option
  const flowOption = [
    { value: '1', label: '双向' },
    { value: '2', label: '出网' },
    { value: '3', label: '进网' },
  ]

  const emit = defineEmits<{
    (e: 'on-close-event'): void
    (e: 'on-reflash'): void
  }>()

  const curIndex = computed(() => (adapterPayload.page - 1) * adapterPayload.limit + 1)

  // 关闭回调
  const handleClose = () => {
    handleClosePop()
    emit('on-close-event')
  }

  // 提交表单
  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
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
                if (props.mode == 'add') {
                  // @ts-ignore
                  const { msg } = await flowProbeAddApi(
                    {
                      ...formData,
                    },
                    { password: AesEncryptCBC(password.value) }
                  )
                  ElMessage({ message: msg, type: 'success' })
                } else {
                  // @ts-ignore
                  const id = props.currentRow.id
                  // @ts-ignore
                  const { msg } = await flowProbeUpdataApi(
                    id,
                    {
                      ...formData,
                      id,
                    },
                    { password: AesEncryptCBC(password.value) }
                  )
                  ElMessage({ message: msg, type: 'success' })
                }
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
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  // 初始化数据
  const initData = () => {
    title.value = props.mode == 'add' ? '添加' : '编辑'
    visible.value = props.showPage
    if (props.mode == 'edit') {
      const data = JSON.parse(JSON.stringify(props.currentRow))
      Object.keys(formData).forEach((item: string) => {
        if (item !== 'dialogTitle' && item !== 'moduleType') {
          // @ts-ignore
          formData[item] = data[item]
        }
      })
    }
    getAdapterById()
  }

  // 获取网卡id
  const getAdapterById = async () => {
    try {
      listLoading.value = true
      const { data } = await getAdapterByIdApi({
        ...adapterPayload,
        query: JSON.stringify(adapterPayload.query),
      })
      listDate = data.list
      totals.value = data.total
    } finally {
      listLoading.value = false
    }
  }

  onMounted(() => {
    initData()
  })

  // 页容量改变
  const handleSizeChange = (val: number) => {
    adapterPayload.limit = val
    getAdapterById()
  }

  // 页面改变
  const handleCurrentChange = (val: number) => {
    adapterPayload.page = val
    getAdapterById()
  }

  // 筛选
  const handleChoose = () => {
    const res = tableRef.value?.getSelectionRows()
    const arr: string[] = []
    formData.adapterId = []
    res.forEach((item: any) => {
      arr.push(item.name)
      // @ts-ignore
      formData.adapterId.push(item.id)
    })
    formData.adapterIdStr = arr.join(',')
    popVisible.value = false
  }

  // 打开popover
  const handleClosePop = () => {
    popVisible.value = false
    tableRef.value?.clearSelection()
  }

  // 打开关闭
  const openPopEvent = () => {
    popVisible.value = true
    const rows: any[] = []
    formData.adapterId.forEach((item: string) => {
      listDate.forEach((td: any) => {
        if (td.id == item) {
          rows.push(td)
        }
      })
    })
    rows.forEach((row) => {
      tableRef.value!.toggleRowSelection(row, true)
    })
  }
</script>

<script lang="ts">
  export default {
    name: 'TrafficAddEdit',
  }
</script>
<template>
  <div class="traffic-add-edit">
    <el-dialog v-model="visible" :before-close="handleClose" :fullscreen="true" :title="title">
      <el-form ref="formRef" label-width="150px" :model="formData" :rules="rules">
        <el-row>
          <el-col :span="12">
            <el-form-item label="名称" prop="name">
              <el-input v-model="formData.name" maxlength="128" show-word-limit />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="类型" prop="type">
              <el-radio v-model="formData.type" label="1">普通链路</el-radio>
              <el-radio v-model="formData.type" label="2">Agent聚合链路</el-radio>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="网卡" prop="adapterId">
              <el-popover placement="bottom" trigger="click" :visible="popVisible" width="40%">
                <vab-query-form>
                  <vab-query-form-left-panel :span="8">
                    <el-button type="primary" @click="handleChoose">选择</el-button>
                    <el-button type="danger" @click="handleClosePop">关闭</el-button>
                  </vab-query-form-left-panel>
                  <vab-query-form-right-panel :span="16">
                    <el-input
                      v-model="adapterPayload.query.keyWord"
                      class="input-with-select"
                      :disabled="listLoading"
                      placeholder="请输入名称"
                      @keyup.enter="getAdapterById"
                    >
                      <template #append>
                        <el-button :icon="Search" @click="getAdapterById" />
                      </template>
                    </el-input>
                  </vab-query-form-right-panel>
                </vab-query-form>
                <el-table
                  ref="tableRef"
                  v-loading="listLoading"
                  :border="true"
                  :data="listDate"
                  row-key="id"
                  @cell-contextmenu="useTableCopy"
                >
                  <el-table-column show-overflow-tooltip type="selection" />
                  <el-table-column :align="'center'" fixed="left" label="序号" width="55">
                    <template #default="{ $index }">
                      {{ curIndex + $index }}
                    </template>
                  </el-table-column>
                  <el-table-column label="名称" property="name" width="120" />
                  <el-table-column label="设备" property="flowDeviceIdStr" width="200" />
                  <el-table-column label="IPv4地址" property="address" />
                </el-table>
                <el-pagination
                  background
                  :current-page="adapterPayload.page"
                  :layout="layout"
                  :page-size="adapterPayload.limit"
                  :page-sizes="[10, 20, 50, 100]"
                  :total="totals"
                  @current-change="handleCurrentChange"
                  @size-change="handleSizeChange"
                />
                <template #reference>
                  <el-input v-model="formData.adapterIdStr" :disabled="popVisible" @click="openPopEvent" />
                </template>
              </el-popover>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="enabled">
              <el-radio v-model="formData.enabled" :label="true">启用</el-radio>
              <el-radio v-model="formData.enabled" :label="false">停用</el-radio>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="识别IP层级" prop="ipLayer">
              <el-select v-model="formData.ipLayer" placeholder="请选择" style="width: 100%">
                <el-option v-for="item in option" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="IP地址" prop="ipAddr">
              <el-input v-model="formData.ipAddr" maxlength="512" show-word-limit />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="MAC地址" prop="macAddr">
              <el-input v-model="formData.macAddr" maxlength="512" show-word-limit />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="内网地址" prop="intranetIp">
              <el-input v-model="formData.intranetIp" maxlength="512" show-word-limit />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="标识类型" prop="subType">
              <el-select v-model="formData.subType" placeholder="请选择" style="width: 100%">
                <el-option v-for="item in typeOption" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="标识值" prop="subValue">
              <el-input v-model="formData.subValue" maxlength="128" show-word-limit />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="启用存储" prop="enableStorage">
              <el-radio v-model="formData.enableStorage" :label="true">启用</el-radio>
              <el-radio v-model="formData.enableStorage" :label="false">停用</el-radio>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="存储裁包" prop="storagePackageCut">
              <el-input v-model="formData.storagePackageCut">
                <template #suffix>Byte</template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="启用分析" prop="enableAnalysis">
              <el-radio v-model="formData.enableAnalysis" :label="true">启用</el-radio>
              <el-radio v-model="formData.enableAnalysis" :label="false">停用</el-radio>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分析裁包" prop="analysisPackageCut">
              <el-input v-model="formData.analysisPackageCut">
                <template #suffix>Byte</template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="SSL自动裁包" prop="sslPackageCut">
              <el-radio v-model="formData.sslPackageCut" :label="true">启用</el-radio>
              <el-radio v-model="formData.sslPackageCut" :label="false">停用</el-radio>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="去重数据包" prop="removeDupPackage">
              <el-select v-model="formData.removeDupPackage" placeholder="请选择" style="width: 100%">
                <el-option v-for="item in deWeightOption" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="开启分析" prop="openAnalysis">
              <el-select v-model="formData.openAnalysis" multiple placeholder="请选择" style="width: 100%">
                <el-option v-for="item in analyseOption" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="关闭存储" prop="closeStorage">
              <el-select v-model="formData.closeStorage" multiple placeholder="请选择" style="width: 100%">
                <el-option v-for="item in storageOption" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="总带宽" prop="totalBandwidth">
              <el-input v-model="formData.totalBandwidth">
                <template #suffix>Mbps</template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="进网带宽" prop="inBandwidth">
              <el-input v-model="formData.inBandwidth">
                <template #suffix>Mbps</template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="出网带宽" prop="outBandwidth">
              <el-input v-model="formData.outBandwidth">
                <template #suffix>Mbps</template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="流量方向" prop="flowDirection">
              <el-select v-model="formData.flowDirection" placeholder="请选择" style="width: 100%">
                <el-option v-for="item in flowOption" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <el-row>
        <el-col :offset="11" :span="2" style="text-align: right">
          <el-button type="primary" @click="submitForm(formRef)">保存</el-button>
        </el-col>
      </el-row>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  :deep() {
    .el-scrollbar {
      height: 200px !important;
    }
  }
</style>
