<script setup lang="ts">
  import { FormInstance, TableInstance } from 'element-plus'

  import { requireRules } from '~/src/utils/rules'

  import {
    flowProbeListApi,
    flowKafkaPushAddApi,
    flowKafkaPushUpdataApi,
    kafkaServerConfigApi,
  } from '@/api-ecs/equipment'

  import { Search } from '@element-plus/icons-vue'

  import { useTableCopy } from '@/utils'
  import AesEncryptCBC from '~/src/utils/crypto'

  const props = defineProps<{
    mode: string
    showPage: boolean
    currentRow: object
  }>()

  const serverTableRef = ref<TableInstance>() // 服务配置表单实例

  const probeTableRef = ref<TableInstance>() // 分析链路表单实例

  const formRef = ref<FormInstance>() // 表单实例

  const visible = ref(false) // 弹框显隐

  const kafkaServerConfigVisible = ref(false) // 服务配置popover弹框显隐

  const probeVisible = ref(false) // 分析链路popover弹框显隐

  const title = ref('添加') // 标题

  let serverListDate = reactive<object[]>([]) // 表格数据

  let porbeListDate = reactive<object[]>([]) // 表格数据

  const layout = ref('total, sizes, prev, pager, next, jumper')

  const serverTotals = ref(0) // 总条数

  const porbeTotal = ref(0) // 总条数

  const listLoading = ref(false) // 是否加载

  // 分析链路列表参数
  const porbePayload = reactive({
    page: 1,
    limit: 10,
    query: {
      keyWord: undefined,
    },
    url: 'flowProbe/page',
  })

  // 获取kafka服务列表参数
  const serverPayload = reactive({
    page: 1,
    limit: 20,
    query: {
      keyWord: undefined,
    },
    url: 'kafkaServerConfig/page',
  })

  // 表单数据
  const formData = reactive({
    dialogTitle: 'add', // 添加链路时候的固定值
    moduleType: 'flowKafkaPush',
    key: undefined,
    recordPartition: undefined,
    recordTopic: undefined,
    frequence: '0',
    delayTime: '5',
    kafkaServerConfigId: undefined,
    kafkaServerConfigIdStr: undefined,
    enabled: true,
    probeIds: [],
    probeIdsStr: '',
    flowTable: 'dnsLog',
    note: undefined,
  })

  // 表单数据校验规则
  const rules = reactive({
    recordTopic: requireRules,
    frequence: requireRules,
    delayTime: requireRules,
    enabled: requireRules,
    kafkaServerConfigIdStr: requireRules,
    probeIdsStr: requireRules,
    enableflowTableAnalysis: requireRules,
  })

  // 推送频率option
  const frequenceOption = [
    { value: '0', label: '实时' },
    { value: '10', label: '10秒' },
    { value: '20', label: '30秒' },
    { value: '60', label: '1分钟' },
    { value: '300', label: '5分钟' },
  ]

  // 识别类型option
  const delayTimeOption = [
    { value: '5', label: '5秒' },
    { value: '10', label: '10秒' },
    { value: '30', label: '30秒' },
    { value: '60', label: '1分钟' },
  ]

  // 数据表option
  const flowTableOption = [
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

  const emit = defineEmits<{
    (e: 'on-close-event'): void
    (e: 'on-reflash'): void
  }>()

  const serverIndex = computed(() => (serverPayload.page - 1) * serverPayload.limit + 1)

  const porbeIndex = computed(() => (porbePayload.page - 1) * porbePayload.limit + 1)

  // 关闭回调
  const handleClose = () => {
    probeVisible.value = false
    kafkaServerConfigVisible.value = false
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
                  const { message } = await flowKafkaPushAddApi(
                    {
                      ...formData,
                    },
                    { password: AesEncryptCBC(password.value) }
                  )
                  ElMessage({ message: message, type: 'success' })
                } else {
                  // @ts-ignore
                  const id = props.currentRow.id
                  // @ts-ignore
                  const { message } = await flowKafkaPushUpdataApi(
                    id,
                    {
                      ...formData,
                      id,
                    },
                    { password: AesEncryptCBC(password.value) }
                  )
                  ElMessage({ message: message, type: 'success' })
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
  }

  // 获取分析链路列表
  const getflowProbePageList = async () => {
    try {
      listLoading.value = true
      const { data } = await flowProbeListApi({
        ...porbePayload,
        // @ts-ignore
        query: JSON.stringify(porbePayload.query),
      })
      porbeListDate = data.list
      porbeTotal.value = data.total
    } finally {
      listLoading.value = false
    }
  }

  onMounted(() => {
    initData()
  })

  // 服务配置页容量改变
  const handleServerSizeChange = (val: number) => {
    serverPayload.limit = val
    getkafkaServerConfigList()
  }

  // 服务配置页面改变
  const handleServerChange = (val: number) => {
    serverPayload.page = val
    getkafkaServerConfigList()
  }

  // 分析链路页容量改变
  const handlePorbeSizeChange = (val: number) => {
    serverPayload.limit = val
  }

  // 分析链路页面改变
  const handlePorbeChange = (val: number) => {
    serverPayload.page = val
  }

  // 筛选
  const handleChoose = () => {
    const res = probeTableRef.value?.getSelectionRows()
    const arr: string[] = []
    formData.probeIds = []
    res.forEach((item: any) => {
      arr.push(item.name)
      // @ts-ignore
      formData.probeIds.push(item.id)
    })
    formData.probeIdsStr = arr.join(',')
    probeVisible.value = false
  }

  // 打开popover
  const handleClosePop = (remark: string) => {
    if (remark == 'server') {
      serverTableRef.value?.clearSelection()
      kafkaServerConfigVisible.value = false
    } else {
      probeVisible.value = false
      probeTableRef.value?.clearSelection()
    }
  }

  // 获取kafka服务列表
  const getkafkaServerConfigList = async () => {
    try {
      listLoading.value = true
      const { data } = await kafkaServerConfigApi({
        ...serverPayload,
        query: JSON.stringify(serverPayload.query),
      })
      serverListDate = data.list
      serverTotals.value = data.total
    } finally {
      listLoading.value = false
    }
  }

  // 服务配置选择
  const handleCurrentChange = (val: any) => {
    if (formData.kafkaServerConfigId == val.id) {
      formData.kafkaServerConfigId = undefined
      formData.kafkaServerConfigIdStr = undefined
      serverTableRef.value!.setCurrentRow(false)
    } else {
      formData.kafkaServerConfigId = val.id
      formData.kafkaServerConfigIdStr = val.name
    }
  }

  // 打开
  const openPopEvent = async (remark: string) => {
    if (remark == 'server') {
      kafkaServerConfigVisible.value = true
      await getkafkaServerConfigList()
      const rows = serverListDate.filter((item: any) => {
        return item.id == formData.kafkaServerConfigId
      })
      serverTableRef.value!.setCurrentRow(null)
      serverTableRef.value!.setCurrentRow(rows[0])
    } else {
      await getflowProbePageList()
      probeVisible.value = true
      const rows: any[] = []
      formData.probeIds?.forEach((item: string) => {
        porbeListDate.forEach((td: any) => {
          if (td.id == item) {
            rows.push(td)
          }
        })
      })
      rows.forEach((row) => {
        probeTableRef.value!.toggleRowSelection(row, true)
      })
    }
  }
</script>

<script lang="ts">
  export default {
    name: 'KafkaPushEdit',
  }
</script>
<template>
  <div class="kafka-push-edit">
    <el-dialog v-model="visible" :before-close="handleClose" :title="title" width="800">
      <el-form ref="formRef" label-width="150px" :model="formData" :rules="rules">
        <el-row>
          <el-col :span="12">
            <el-form-item label="信息键" prop="key">
              <el-input v-model="formData.key" maxlength="64" show-word-limit />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="信息分区" prop="recordPartition">
              <el-input v-model="formData.recordPartition" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="信息主题" prop="recordTopic">
              <el-input v-model="formData.recordTopic" maxlength="64" show-word-limit />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="推送频率" prop="frequence">
              <el-select v-model="formData.frequence" placeholder="请选择" style="width: 100%">
                <el-option v-for="item in frequenceOption" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="延迟时间" prop="delayTime">
              <el-select v-model="formData.delayTime" placeholder="请选择" style="width: 100%">
                <el-option v-for="item in delayTimeOption" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="服务配置" prop="kafkaServerConfigIdStr">
              <el-popover placement="bottom" trigger="click" :visible="kafkaServerConfigVisible" width="40%">
                <vab-query-form>
                  <vab-query-form-left-panel :span="4">
                    <el-button type="danger" @click="handleClosePop('server')">关闭</el-button>
                  </vab-query-form-left-panel>
                  <vab-query-form-right-panel :span="20">
                    <el-input
                      v-model="serverPayload.query.keyWord"
                      class="input-with-select"
                      :disabled="listLoading"
                      placeholder="请输入名称"
                      @keyup.enter="getkafkaServerConfigList"
                    >
                      <template #append>
                        <el-button :icon="Search" @click="getkafkaServerConfigList" />
                      </template>
                    </el-input>
                  </vab-query-form-right-panel>
                </vab-query-form>
                <el-table
                  ref="serverTableRef"
                  v-loading="listLoading"
                  :border="true"
                  :data="serverListDate"
                  highlight-current-row
                  row-key="id"
                  @cell-click="handleCurrentChange"
                  @cell-contextmenu="useTableCopy"
                >
                  <el-table-column :align="'center'" fixed="left" label="序号" width="55">
                    <template #default="{ $index }">
                      {{ serverIndex + $index }}
                    </template>
                  </el-table-column>
                  <el-table-column label="名称" property="name" width="120" />
                  <el-table-column label="集群信息" property="addr" width="200" />
                  <el-table-column label="认证方式" property="authTypeStr" />
                </el-table>
                <el-pagination
                  background
                  :current-page="serverPayload.page"
                  :layout="layout"
                  :page-size="serverPayload.limit"
                  :page-sizes="[10, 20, 50, 100]"
                  :total="serverIndex"
                  @current-change="handleServerChange"
                  @size-change="handleServerSizeChange"
                />
                <template #reference>
                  <el-input
                    v-model="formData.kafkaServerConfigIdStr"
                    :disabled="kafkaServerConfigVisible"
                    @click="openPopEvent('server')"
                  />
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
            <el-form-item label="分析链路" prop="probeIdsStr">
              <el-popover placement="bottom" trigger="click" :visible="probeVisible" width="40%">
                <vab-query-form>
                  <vab-query-form-left-panel :span="8">
                    <el-button type="primary" @click="handleChoose">选择</el-button>
                    <el-button type="danger" @click="handleClosePop('probe')">关闭</el-button>
                  </vab-query-form-left-panel>
                  <vab-query-form-right-panel :span="16">
                    <el-input
                      v-model="porbePayload.query.keyWord"
                      class="input-with-select"
                      :disabled="listLoading"
                      placeholder="请输入名称"
                      @keyup.enter="getflowProbePageList"
                    >
                      <template #append>
                        <el-button :icon="Search" @click="getflowProbePageList" />
                      </template>
                    </el-input>
                  </vab-query-form-right-panel>
                </vab-query-form>
                <el-table
                  ref="probeTableRef"
                  v-loading="listLoading"
                  :border="true"
                  :data="porbeListDate"
                  row-key="id"
                  @cell-contextmenu="useTableCopy"
                >
                  <el-table-column show-overflow-tooltip type="selection" />
                  <el-table-column :align="'center'" fixed="left" label="序号" width="55">
                    <template #default="{ $index }">
                      {{ porbeIndex + $index }}
                    </template>
                  </el-table-column>
                  <el-table-column label="名称" property="name" width="120" />
                  <el-table-column label="设备" property="flowDeviceIdStr" width="200" />
                  <el-table-column label="网卡" property="adapterIdStr" :resizable="true" show-overflow-tooltip />
                </el-table>
                <el-pagination
                  background
                  :current-page="porbePayload.page"
                  :layout="layout"
                  :page-size="porbePayload.limit"
                  :page-sizes="[10, 20, 50, 100]"
                  :total="porbeTotal"
                  @current-change="handlePorbeChange"
                  @size-change="handlePorbeSizeChange"
                />
                <template #reference>
                  <el-input v-model="formData.probeIdsStr" :disabled="probeVisible" @click="openPopEvent('probe')" />
                </template>
              </el-popover>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="数据表" prop="flowTable">
              <el-select v-model="formData.flowTable" placeholder="请选择" style="width: 100%">
                <el-option v-for="item in flowTableOption" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="备注" prop="note">
              <el-input v-model="formData.note" maxlength="128" show-word-limit />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <el-row>
        <el-col :offset="19" :span="5" style="text-align: right">
          <el-button type="primary" @click="submitForm(formRef)">保存</el-button>
          <el-button @click="handleClose">取消</el-button>
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
