<script setup lang="ts">
  import ApplicationConfiguration from '@/ecs/assets/components/application-configuration.vue'

  // import InfoDialog from './components/network-layer-info.vue'

  import VabDialog from '@/plugins/VabDialog/index.vue'

  import TraceabilityFieldStatistic from '@/ecs/alert/components/field-statistic.vue'

  import { NetworkLayerParams, RetrieveIndexType } from '@/types/index'

  import { DownloadPcapApi } from '~/src/api-ecs/toolbox'

  import dayjs from 'dayjs'

  import { useCopy, useTableCopy } from '~/src/utils'

  import {
    getEventStatisticsAPI,
    getFlowProbesAPI,
    getPacketDecodeListAPI,
    getTabelListAPI,
    getDownloadQueryAPI,
  } from '~/src/api-ecs/retrieve'

  import { downloadLogPacketFlowSearch } from '~/src/utils/download'

  import type { FormInstance } from 'element-plus'

  import { formatTime } from '~/src/utils/time'

  import { ElTree } from 'element-plus'

  import { useUserStore } from '@/store/modules/user'

  import DecodingPacket from './components/traceability-decoding-packet.vue'

  import DecodingStream from './components/traceability-decoding-stream.vue'

  import { TableColumnItemType } from '~/types/store'

  import { getAllDisPlaysFiledApi } from '~/src/api-ecs/public'

  import { updateDisplayApi } from '~/src/api-ecs/custom-field'
  import { cloneDeep } from 'lodash'
  import { Search } from '@element-plus/icons-vue'

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  const userStore = useUserStore()

  const { getTableColumn } = userStore

  const treeRef = ref<InstanceType<typeof ElTree>>()

  const formRef = ref<FormInstance>()

  const traceabilityFieldRef = ref<InstanceType<typeof TraceabilityFieldStatistic>>()

  // 检索信息表单数据
  const queryForm = reactive<NetworkLayerParams>({
    searchSql: '',
    indexType: 11,
    pageNum: 1,
    pageSize: 10000,
    orderType: 'desc',
    orderField: 'flowBeginTimeNs',
    startTime: '',
    endTime: '',
    scrollId: undefined,
  })

  const query = reactive({
    topField: 'topField',
    timeRange: '',
    flowProbeIds: ['0'],
    searchTable: 'eventStat',
    searchModel: 'all',
    eventFilterList: [],
    keyword: undefined,
  })

  const topValue = ref('100')

  const leftValue = ref('100')

  const dialogValue = ref('100')

  const topOptions = [
    { label: 'Top100', value: '100' },
    { label: 'Top500', value: '500' },
    { label: 'Top1000', value: '1000' },
    { label: 'Top2000', value: '2000' },
    { label: 'Top5000', value: '5000' },
  ]
  const topFieldOptions = [
    { label: '总数据包数', value: 'totalPkts' },
    { label: '客户端数据包数', value: 'clientPkts' },
    { label: '服务端数据包数', value: 'serverPkts' },
    { label: '每秒数据包数', value: 'totalPktps' },
    { label: '客户端每秒数据包数', value: 'clientPktps' },
    { label: '服务端每秒数据包数', value: 'serverPktps' },
    { label: '总字节数', value: 'totalBytes' },
    { label: '客户端字节数', value: 'clientBytes' },
    { label: '服务端字节数', value: 'serverBytes' },
    { label: '比特率', value: 'totalBitps' },
    { label: '客户端比特率', value: 'clientBitps' },
    { label: '服务端比特率', value: 'serverBitps' },
    { label: '总负载', value: 'totalPayload' },
    { label: '客户端负载', value: 'clientPayload' },
    { label: '服务端负载', value: 'serverPayload' },
    { label: '传输效率', value: 'totalTransRate' },
    { label: '客户端传输效率', value: 'clientTransRate' },
    { label: '服务端传输效率', value: 'serverTransRate' },
    { label: '事件总数', value: 'eventCount' },
    { label: '告警数量', value: 'alarmCount' },
    { label: '虚拟网标识', value: 'vlanId' },
    { label: '第二层虚拟网标识', value: 'vlanId2' },
    { label: '虚拟资产标识', value: 'vai' },
  ]
  // 是否加载
  const listLoading = ref(false)

  const decodingVisible = ref(false)

  // 表头字段
  const tableColumn = ref<TableColumnItemType[]>([])

  // 检索返回的数据
  const queryPage = reactive({
    total: 0,
    title: '',
    listDate: [] as object[],
    queryTimes: '',
  })

  const formData = reactive({
    ip: '',
    clientIp: '',
    serverIp: '',
    serverPort: '',
  })
  const eventFilter = ref([])
  const topField = ref('totalBitps')
  const eventFilterList = ref<any[]>([])
  const showFilter = ref(false)
  const packetDecodeList = ref()

  const checkedAll = ref()

  const meta = ref()

  const cur_index = ref()

  const searchValue = ref()

  const itemLoading = ref(false)

  const decodingStreamref = ref()

  const decodingPacketref = ref()

  const activeName = ref('decodingStream')

  const objectList = ref()

  const decodingPackeEvent = () => {
    objectList.value = []
    activeName.value = 'decodingPacket'
    const list = packetDecodeList.value.filter((item: { checked: string }) => {
      return item.checked
    })
    if (list.length > 0) {
      list.forEach((item: any) => {
        currentData.value.clientPort = item.clientPort
        const { clientIp, clientPort, probeId, serverIp, serverPort } = currentData.value
        objectList.value.push({ clientIp, clientPort, probeId, serverIp, serverPort })
      })
      decodingPacketref.value?.getData()
    }
  }

  // 检索
  const handleInput = () => {
    packetDecodeList.value = meta.value.filter((item: { sessionInfo: string }) => {
      if (searchValue.value.trim()) {
        return item.sessionInfo.search(searchValue.value) !== -1
      } else {
        return true
      }
    })
  }

  let decodingStreamFlag = true
  let decodingPacketFlag = true
  const handleActive = (index: number, item: any) => {
    if (cur_index.value == index) return
    decodingStreamFlag = true
    decodingPacketFlag = true
    objectList.value = []
    currentData.value.clientPort = item.clientPort
    cur_index.value = index
    if (activeName.value == 'decodingStream') {
      decodingStreamFlag && decodingStreamref.value.getData()
      decodingStreamFlag = false
    } else {
      decodingPacketFlag && decodingPacketref.value.getData()
      decodingPacketFlag = false
    }
  }

  // 全选
  const checkedAllEvent = () => {
    packetDecodeList.value.forEach((item: { checked: boolean }) => {
      item['checked'] = checkedAll.value
    })
  }

  // 单选
  const checkedEvent = () => {
    const list = packetDecodeList.value.filter((item: { checked: boolean }) => {
      return item['checked'] == true
    })
    checkedAll.value = list.length == packetDecodeList.value.length
  }

  // 获取弹框解码右侧list
  const getPacketDecodeList = async () => {
    const { clientIp, clientPort, probeId, serverIp, serverPort } = currentData.value
    const query = {
      timeRange: '',
      timeStep: '',
      objectList: [],
    }
    query.timeRange = `${queryForm.startTime} - ${queryForm.endTime}`
    query.timeStep = 'minuteStep'
    // @ts-ignore
    query.objectList = [{ clientIp, clientPort, probeId, serverIp, serverPort }]

    itemLoading.value = true
    try {
      const { data } = await getPacketDecodeListAPI({ packetCut: '0', top: dialogValue.value, query })
      packetDecodeList.value = data.data
      packetDecodeList.value.forEach((item: { checked: boolean }) => {
        item['checked'] = false
      })
      meta.value = JSON.parse(JSON.stringify(data.data))
    } finally {
      itemLoading.value = false
    }
  }

  // 当复选框被点击的时候触发
  const checkChangeEvent = (isChecked: boolean, child: any) => {
    remark = 'checked'
    if (isChecked) {
      eventFilterList.value.push(child)
    } else {
      const index = eventFilterList.value.findIndex((item) => item.id === child.id)
      eventFilterList.value.splice(index, 1)
    }
  }

  watch(
    () => eventFilterList.value,
    () => {
      getTableData()
    },
    {
      deep: true,
    }
  )

  const useTableCopyEvent = (row: any, column: any, cell: any, event: any) => {
    const target = column.property
    const value = row[target]
    const arr = [{ label: '复制', callback: useCopy, value: value }]
    useTableCopy(row, column, cell, event, arr)
  }

  const showFiledConfig = ref(false) // 显示表头字段配置

  // 表头字段配置
  const tableHeadConfig = (val: boolean) => {
    showFiledConfig.value = val
  }
  // 用户字段
  const userDisPlaysFiled = ref()

  // 获取用户所有字段
  const getAllDisPlaysFiled = async () => {
    const { data } = await getAllDisPlaysFiledApi()
    userDisPlaysFiled.value = data
    formatColum(24)
  }

  const allColumnData = ref()

  function formatColum(type: RetrieveIndexType) {
    const columnData = getTableColumn(24)
    allColumnData.value = getTableColumn(24)
    const userColumnData = userDisPlaysFiled.value[type] as number[]
    if (userColumnData?.length > 0) {
      tableColumn.value = userColumnData.map((key) =>
        columnData.find((item) => item.id === key)
      ) as TableColumnItemType[]
    } else {
      tableColumn.value = []
    }
  }

  async function configurationHandel(showField: TableColumnItemType[]) {
    tableColumn.value = showField
    const { msg } = await updateDisplayApi({
      indexType: 24,
      displayIds: showField.map((i) => i.id),
    })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    getAllDisPlaysFiled()
  }
  packetDecodeList
  // 自定义时间
  const timeDate = ref()

  const SearchCondition = ref('ip') // 检索条件选项

  const SearchOption = [
    { value: 'ip', label: 'IP地址' },
    { value: 'clientIp', label: '源IP' },
    { value: 'serverIp', label: '目的IP' },
    { value: 'clientServerIp', label: '源+目的IP' },
  ]

  // 时间选项
  const timeDuratioOptions = [
    {
      value: '1hours',
      label: '1小时',
    },
    {
      value: '3hours',
      label: '3小时',
    },
    {
      value: '6ours',
      label: '6小时',
    },
    {
      value: '24hours',
      label: '24小时',
    },
    {
      value: 'today',
      label: '今日',
    },
    {
      value: 'yesterday',
      label: '昨天',
    },
    {
      value: 'last-three-days',
      label: '最近三天',
    },
    {
      value: 'last-week',
      label: '最近一周',
    },
    {
      value: 'last-two-week',
      label: '最近两周',
    },
    {
      value: 'user-defined',
      label: '自定义时间',
    },
  ]

  const treeLoading = ref(false)

  // 获取表格序号
  const curIndex = computed(() => (queryForm.pageNum - 1) * queryForm.pageSize + 1)

  // 查询时间
  const timeDuration = ref('1hours')

  const currentData = ref()

  let remark = ''

  const nodelist = ref()

  const curChangeEvent = (row: any) => {
    nodelist.value = row
    remark = 'node'
    getTableData()
  }

  const getTableData = async () => {
    query.timeRange = `${queryForm.startTime} - ${queryForm.endTime}`
    query.flowProbeIds = flowProbesValue.value as string[]
    query.searchModel = SearchCondition.value
    query.eventFilterList =
      remark == 'node' ? nodelist.value : eventFilterList.value && JSON.parse(JSON.stringify(eventFilterList.value))
    const res = JSON.stringify({ ...formData, ...query })
    listLoading.value = true
    try {
      const { data } = await getTabelListAPI({
        top: topValue.value,
        moduleType: 'eventStat',
        query: res,
        groupBy: 'clientIp,serverIp,serverPort',
      })
      queryPage.listDate = data.list
      queryPage.total = data.total
      queryPage.queryTimes = data.extendObj?.cost
    } finally {
      listLoading.value = false
    }
  }
  const searchEvent = async () => {
    // getTreeData()
    getTableData()
  }

  // 打开详情
  const handleInfo = (row: any) => {
    const { clientIp, clientPort, probeId, serverIp, serverPort } = row
    currentData.value = {
      clientIp,
      clientPort,
      probeId,
      serverIp,
      serverPort,
    }
    currentData.value = row
    decodingVisible.value = true
    getPacketDecodeList()
  }

  watch(
    () => decodingVisible.value,
    () => {
      if (!decodingVisible.value) {
        cur_index.value = undefined
        decodingStreamFlag = true
        decodingPacketFlag = true
        activeName.value = 'decodingStream'
      }
    }
  )

  watch(
    () => activeName.value,
    () => {
      if (!decodingVisible.value || cur_index.value === undefined) return
      if (activeName.value == 'decodingStream') {
        decodingStreamFlag && decodingStreamref.value.getData()
        decodingStreamFlag = false
      } else {
        setTimeout(() => {
          decodingPacketFlag && decodingPacketref.value.getData()
          decodingPacketFlag = false
        }, 0)
      }
    }
  )

  const getQueryDate = () => {
    if (timeDuration.value === 'user-defined') {
      formatUserTime()
      return
    }
    const timeDate = dayjs()
    const endDate = ['1hours', '3hours', '6hours', '24hours', 'last-week', 'last-two-week', 'last-three-days'].includes(
      timeDuration.value
    )
      ? timeDate.add(10, 'minute').format('YYYY-MM-DD HH:mm:ss')
      : timeDate.endOf('day').format('YYYY-MM-DD HH:mm:ss')
    const cur = timeDate.startOf('day')
    let startDate = null
    switch (timeDuration.value) {
      case '1hours':
        startDate = timeDate.subtract(1, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case '3hours':
        startDate = timeDate.subtract(3, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case '6hours':
        startDate = timeDate.subtract(6, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case '24hours':
        startDate = timeDate.subtract(24, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'today':
        startDate = timeDate.startOf('day').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'yesterday':
        startDate = cur.subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'last-three-days':
        startDate = timeDate.subtract(3, 'day').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'last-week':
        startDate = timeDate.subtract(1, 'week').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'last-two-week':
        startDate = timeDate.subtract(2, 'week').format('YYYY-MM-DD HH:mm:ss')
        break
      default:
        startDate = timeDate.startOf('day').format('YYYY-MM-DD HH:mm:ss')
        break
    }
    queryForm.endTime =
      timeDuration.value === 'yesterday'
        ? timeDate.endOf('day').subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss')
        : endDate
    queryForm.startTime = startDate
  }

  const formatDate = (row: any, key: string) => {
    if (!row.key) return
    const time = row[key] / 1000000
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss.SSS')
  }

  const download = async (row: any) => {
    const { getUserId } = useUserStore()
    $baseConfirm('你确定要下载吗', null, async () => {
      try {
        const { clientIp, clientPort, serverIp, serverPort, probeId } = row
        // @ts-ignore
        query.objectList = [{ clientIp, clientPort, probeId, serverIp, serverPort }]
        query.searchModel = 'all'
        const { msg } = await DownloadPcapApi({ uid: getUserId(), query }, 1)
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        // const { data } = await getDownloadQueryAPI({ query: JSON.stringify(query) })
        // downloadLogPacketFlowSearch(data)
      } catch (error) {
        $baseMessage('下载失败', 'error', 'vab-hey-message-error')
      }
    })
  }

  const formatBitps = (row: any, cell: string, unit: 'bitps' | 'pkts') => {
    const unitlist: { bitps: string[]; pkts: string[] } = {
      bitps: ['Gbps', 'Mbps', 'Kbps', 'bps'],
      pkts: ['GB', 'MB', 'KB', 'B'],
    }
    let value = row[cell]
    if (value > 1073000000) {
      value = (value / 1024 / 1024 / 1024).toFixed(2) + unitlist[unit][0]
    } else if (value > 1048000) {
      value = (value / 1024 / 1024).toFixed(2) + unitlist[unit][1]
    } else if (value > 1024) {
      value = (value / 1024).toFixed(2) + unitlist[unit][2]
    } else {
      value + unitlist[unit][3]
    }
    return value
  }

  watch(
    () => timeDuration.value,
    () => {
      getQueryDate()
    },
    {
      immediate: true,
    }
  )

  watch(
    () => SearchCondition.value,
    (newVal, oldVal) => {
      if (newVal != oldVal) {
        queryForm.searchSql = ''
        // changeRules(newVal)
        // formRef.value!.resetFields()
        // formRef.value!.clearValidate()
      }
    },
    { immediate: true }
  )

  const formatUserTime = () => {
    const [startDate, endDate] = timeDate.value
    queryForm.endTime = dayjs(endDate).format('YYYY-MM-DD HH:mm:ss')
    queryForm.startTime = dayjs(startDate).format('YYYY-MM-DD HH:mm:ss')
  }

  watchEffect(() => {
    if (timeDate.value) {
      formatUserTime()
    }
  })

  const handleDownloadExcel = () => {
    const tHeader: string[] = []
    const filterVal: any = []
    try {
      tableColumn.value.forEach(({ fieldNameCn, fieldNameEn }) => {
        tHeader.push(fieldNameCn)
        filterVal.push(fieldNameEn)
      })
      listLoading.value = true
      import('@/utils/excel').then((excel) => {
        const list = queryPage.listDate
        const data = formatJson(filterVal, list)
        excel.export_json_to_excel({
          header: tHeader,
          data,
          filename: `网络层会话-${formatTime(new Date().getTime())}`,
          autoWidth: true,
          bookType: 'xlsx',
        })
        listLoading.value = false
      })
    } catch (error) {
      listLoading.value = false
    }
  }
  function formatJson(filterVal: any, jsonData: any) {
    return jsonData.map((v: any) =>
      filterVal.map((j: any) => {
        return formatExcelData(v, j)
      })
    )
  }

  function formatExcelData(val: any, key: string) {
    switch (key) {
      case 'totalBytes':
      case 'clientBytes':
      case 'serverBytes':
      case 'totalPayload':
      case 'clientPayload':
      case 'serverPayload':
        return formatBitps(val, key, 'pkts')
      case 'totalBitps':
      case 'clientBitps':
      case 'serverBitps':
        return formatBitps(val, key, 'bitps')
      default:
        return val[key] || ''
    }
  }

  // 右侧下拉树
  const flowProbesList = ref()
  const flowProbesValue = ref(['0'])
  const getFlowProbes = async () => {
    const { data } = await getFlowProbesAPI()
    flowProbesList.value = []
    flowProbesList.value.push(data.data)
    flowProbesValue.value = []
    getAllID(flowProbesList.value)
  }

  const getAllID = (data: any) => {
    data?.forEach((item: any) => {
      flowProbesValue.value.push(item.id)
      if (item.children && item.children.length > 0) {
        getAllID(item.children)
      }
    })
  }

  onMounted(async () => {
    await getFlowProbes()
    getTreeData()
    getAllDisPlaysFiled()
  })
  const checkedStatus = (item: any) => {
    if (!eventFilterList.value.length) return false
    return !!eventFilterList.value.find((i: any) => i.id === item.id)
  }
  const eventFilterChange = (index: number, item: any) => {
    if (!item.baseChildren) {
      item.baseChildren = cloneDeep(item.children)
    }
    const str = eventFilter.value[index]
    item.children = item.baseChildren.filter((i: any) => i.title.includes(str))
  }
  const eventSortChange = (item: any) => {
    item.isAscending = !item.isAscending
    item.children.reverse()
  }
  const treeData = ref()
  const getTreeData = async () => {
    query.timeRange = `${queryForm.startTime} - ${queryForm.endTime}`
    query.topField = topField.value
    query.flowProbeIds = flowProbesValue.value as string[]
    query.searchModel = SearchCondition.value
    const keyword = formData.ip || formData.clientIp || formData.serverIp || formData.serverPort || undefined
    // @ts-ignore
    query.keyword = keyword
    const res = JSON.stringify(query)
    treeLoading.value = true
    try {
      const { data } = await getEventStatisticsAPI({
        top: leftValue.value,
        topField: topField.value,
        moduleType: 'flowSearch',
        query: res,
      })

      treeData.value = data.data
    } finally {
      treeLoading.value = false
    }
  }

  watch(
    () => [leftValue.value, topField.value],
    () => {
      if (leftValue.value) {
        getTreeData()
      }
    }
  )

  watch(
    () => dialogValue.value,
    () => {
      if (dialogValue.value) {
        getPacketDecodeList()
      }
    }
  )

  const defaultProps = {
    children: 'children',
    label: 'name',
    value: 'id',
  }
</script>

<script lang="ts">
  export default {
    name: 'NetworkLayer', // 网络层会话
  }
</script>
<template>
  <div class="networkLayer">
    <!-- 头部检索栏 -->
    <vab-query-form style="width: 100%">
      <vab-query-form-left-panel :span="16">
        <el-form class="my-form-top" inline>
          <el-form-item>
            <el-tree-select
              v-model="flowProbesValue"
              :data="flowProbesList"
              default-expand-all
              :expand-on-click-node="false"
              multiple
              node-key="id"
              :props="defaultProps"
              show-checkbox
              style="width: 200px"
            />
            <el-select v-model="timeDuration" style="width: 140px; margin-left: 10px">
              <el-option v-for="item in timeDuratioOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item v-if="timeDuration === 'user-defined'">
            <vab-date-time-picker v-model="timeDate" style="margin-left: 10px" />
          </el-form-item>
        </el-form>
      </vab-query-form-left-panel>
      <vab-query-form-left-panel :span="24">
        <el-form ref="formRef" class="my-form" inline :model="formData" @submit.prevent>
          <el-form-item>
            <el-select v-model="SearchCondition" style="width: 150px">
              <el-option v-for="item in SearchOption" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item class="my-input">
            <el-row :gutter="10">
              <el-col v-if="SearchCondition == 'ip'" :span="24">
                <el-form-item prop="ip">
                  <el-input
                    v-model="formData.ip"
                    placeholder="请输入IP地址，支持IP格式和掩码格式，多个检索条件使用英文逗号分开"
                  />
                </el-form-item>
              </el-col>
              <el-col v-if="SearchCondition == 'clientIp'" :span="24">
                <el-form-item prop="clientIp">
                  <el-input
                    v-model="formData.clientIp"
                    placeholder="源IP地址，支持IP格式和掩码格式，多个使用英文逗号分开"
                  />
                </el-form-item>
              </el-col>
              <el-col v-if="SearchCondition == 'serverIp'" :span="18">
                <el-form-item prop="serverIp">
                  <el-input
                    v-model="formData.serverIp"
                    placeholder="目的IP地址，支持IP格式和掩码格式，多个使用英文逗号分开"
                  />
                </el-form-item>
              </el-col>
              <el-col v-if="SearchCondition == 'serverIp'" :span="6">
                <el-form-item prop="serverPort">
                  <el-input v-model="formData.serverPort" placeholder="目的端口，可为空，为空代表所有" />
                </el-form-item>
              </el-col>
              <el-col v-if="SearchCondition == 'clientServerIp'" :span="9">
                <el-form-item prop="clientIp">
                  <el-input
                    v-model="formData.clientIp"
                    placeholder="源IP地址，支持IP格式和掩码格式，多个使用英文逗号分开"
                  />
                </el-form-item>
              </el-col>
              <el-col v-if="SearchCondition == 'clientServerIp'" :span="9">
                <el-form-item prop="serverIp">
                  <el-input
                    v-model="formData.serverIp"
                    placeholder="目的IP地址，支持IP格式和掩码格式，多个使用英文逗号分开"
                  />
                </el-form-item>
              </el-col>
              <el-col v-if="SearchCondition == 'clientServerIp'" :span="6">
                <el-form-item prop="serverPort">
                  <el-input v-model="formData.serverPort" placeholder="目的端口，可为空，为空代表所有" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form-item>
          <el-form-item>
            <el-button :disabled="listLoading" style="margin: 0 !important" type="primary" @click="() => searchEvent()">
              检索
            </el-button>
          </el-form-item>
        </el-form>
      </vab-query-form-left-panel>
    </vab-query-form>

    <!-- 树以及表格 -->
    <div class="content">
      <div class="left" :class="{ show: showFilter }">
        <el-select v-model="leftValue" class="my-select">
          <el-option v-for="item in topOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-select v-model="topField" class="my-select border" style="width: 128px">
          <el-option v-for="item in topFieldOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-button class="my-select" style="width: 40px" @click="showFilter = !showFilter">
          <el-icon style="font-size: 18px"><Filter /></el-icon>
        </el-button>
        <div v-if="showFilter" v-loading="treeLoading" class="chapter">
          <div v-for="(item, index) in treeData" :key="item.id" class="chapter-item">
            <div class="el-table">
              {{ item.title }}
              <span
                class="caret-wrapper"
                :class="{ descending: !item.isAscending, ascending: item.isAscending }"
                @click="eventSortChange(item)"
              >
                <i class="sort-caret ascending"></i>
                <i class="sort-caret descending"></i>
              </span>
              <span class="overflow" @click="() => (item.overflow = !item.overflow)">
                {{ !item.overflow ? '收起' : '展开' }}
                <el-icon style="vertical-align: -1px">
                  <ArrowUp v-if="!item.overflow" />
                  <ArrowDown v-else />
                </el-icon>
              </span>
            </div>
            <template v-if="!item.overflow">
              <el-input
                v-if="item.title !== '事件类型'"
                v-model="eventFilter[index]"
                clearable
                style="margin-top: 10px; margin-bottom: 4px"
                :suffix-icon="Search"
                @change="eventFilterChange(index, item)"
              />
              <div v-for="child in item.children" :key="child.id" class="event-statistics">
                <el-tooltip :content="child.name" effect="dark" :offset="0" placement="top" :show-after="200">
                  <el-checkbox
                    :checked="checkedStatus(child)"
                    :label="child.name"
                    :value="child.name"
                    @change="(checked) => checkChangeEvent(checked, child)"
                  />
                </el-tooltip>
                <span class="size">
                  {{ child.indexValue }}
                </span>
              </div>
            </template>
          </div>
        </div>
      </div>
      <div class="right">
        <div class="toolbar">
          <el-button type="primary" @click="handleDownloadExcel">导出</el-button>
          <div style="float: right; display: flex; align-items: center">
            <span>总数：{{ queryPage.listDate.length || 0 }}</span>
            &nbsp;&nbsp;
            <span>{{ queryPage.queryTimes || '耗时0秒' }}</span>
            <el-select v-model="topValue" style="margin-left: 10px; width: 100px">
              <el-option v-for="item in topOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
            <el-tooltip content="字段配置" effect="dark" placement="top">
              <el-button style="width: 40px; margin-left: 10px" @click="tableHeadConfig(true)">
                <vab-icon icon="settings-line" style="font-size: 18px; margin: 0" />
              </el-button>
            </el-tooltip>
          </div>
        </div>
      </div>
    </div>
    <el-table v-loading="listLoading" class="my-table" :data="queryPage.listDate" @cell-contextmenu="useTableCopyEvent">
      <el-table-column
        v-if="queryPage.listDate"
        :align="'center'"
        :index="(index) => curIndex + index"
        label="序号"
        type="index"
        width="55"
      />
      <template v-for="item in tableColumn" :key="item.id">
        <el-table-column
          v-if="item.fieldNameCn.includes('节数') || item.fieldNameCn.includes('负载')"
          :label="item.fieldNameCn"
          min-width="170"
          :prop="item.fieldNameEn"
          resizable
          show-overflow-tooltip
        >
          <template #default="{ row }">
            {{ formatBitps(row, item.fieldNameEn, 'pkts') }}
          </template>
        </el-table-column>
        <el-table-column
          v-else-if="item.fieldNameCn.includes('时间')"
          :label="item.fieldNameCn"
          min-width="170"
          :prop="item.fieldNameEn"
          resizable
          show-overflow-tooltip
        >
          <template #default="{ row }">
            {{ formatDate(row, item.fieldNameEn) }}
          </template>
        </el-table-column>
        <el-table-column
          v-else-if="item.fieldNameCn.includes('比特率')"
          :label="item.fieldNameCn"
          min-width="170"
          :prop="item.fieldNameEn"
          resizable
          show-overflow-tooltip
        >
          <template #default="{ row }">
            {{ formatBitps(row, item.fieldNameEn, 'bitps') }}
          </template>
        </el-table-column>
        <el-table-column
          v-else
          :label="item.fieldNameCn"
          min-width="180"
          :prop="item.fieldNameEn"
          resizable
          show-overflow-tooltip
        />
      </template>
      <el-table-column v-if="queryPage.listDate" fixed="right" label="操作" width="100">
        <template #default="{ row }">
          <el-button size="small" @click="handleInfo(row)">详情</el-button>
          <!-- <el-button size="small" @click="download(row)">PCAP保存</el-button> -->
          <el-dropdown style="vertical-align: middle; margin-left: 8px">
            <el-icon>
              <ArrowDown />
            </el-icon>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="download(row)">PCAP保存</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
      <template #empty>
        <el-empty class="vab-data-empty" description="暂无数据" />
      </template>
    </el-table>
    <!-- 表头字段配置 -->
    <application-configuration
      v-model="showFiledConfig"
      :fields="tableColumn"
      :retrieve-index-type="24"
      @handleok="configurationHandel"
    />

    <traceability-field-statistic ref="traceabilityFieldRef" />
    <!-- 弹框 -->
    <div class="dialog">
      <vab-dialog
        v-model="decodingVisible"
        :close-on-click-modal="false"
        destroy-on-close
        show-fullscreen
        title="详情"
        width="1175px"
      >
        <div class="warp">
          <div class="left">
            <el-input
              v-model="searchValue"
              class="w-50 m-2"
              placeholder="输入关键字"
              :suffix-icon="Search"
              @change="handleInput"
            />
            <div class="line-2">
              <el-select v-model="dialogValue" style="width: 150px">
                <el-option v-for="item in topOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
              <el-button type="primary" @click="decodingPackeEvent">解码</el-button>
            </div>
            <div v-loading="itemLoading" class="chapter">
              <div>
                <el-checkbox v-model="checkedAll" label="会话" @change="checkedAllEvent" />
              </div>
              <hr style="opacity: 0.4" />
              <div
                v-for="(item, index) in packetDecodeList"
                :key="index"
                :class="{ item: true, active: index === cur_index }"
              >
                <el-checkbox v-model="item.checked" @change="checkedEvent" />
                &nbsp;
                <span @click="handleActive(index, item)">
                  {{ item.sessionInfo }}
                </span>
              </div>
            </div>
          </div>
          <div class="right">
            <el-tabs v-model="activeName">
              <el-tab-pane label="数据包" lazy name="decodingPacket">
                <decoding-packet
                  ref="decodingPacketref"
                  :current-data="currentData"
                  :object-list="objectList"
                  :time-date="query.timeRange"
                />
              </el-tab-pane>
              <el-tab-pane label="数据流" name="decodingStream">
                <decoding-stream ref="decodingStreamref" :current-data="currentData" :time-date="query.timeRange" />
              </el-tab-pane>
            </el-tabs>
          </div>
        </div>
      </vab-dialog>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .networkLayer {
    height: calc(100vh - 50px);
    padding: 10px 20px;
    .my-table {
      position: relative;
      margin-top: 15px;
      :deep() {
        .el-table__body-wrapper {
          &::-webkit-scrollbar {
            width: 0 !important;
            height: 0;
          }
          .el-scrollbar__bar.is-vertical > div {
            margin-top: 40px;
          }
        }
      }
    }
    :deep() {
      .el-table__body-wrapper {
        position: initial !important;
      }
      .el-scrollbar {
        position: initial !important;
      }
      .el-scrollbar__bar.is-horizontal {
        position: absolute;
        display: block !important;
      }
      .el-tabs__content {
        overflow: initial;
      }
      &::-webkit-scrollbar {
        width: 0 !important;
        height: 0;
      }
    }
    :deep() {
      .el-tree-node__content > .el-tree-node__expand-icon {
        padding: 0;
      }
      .el-tree-node__content > label.el-checkbox {
        margin-right: 0px;
      }
      .my-form-top {
        .el-select-tags-wrapper.has-prefix {
          overflow: hidden;
          height: 31px;
          line-height: 36px;
        }
        .el-select .el-input {
          height: 31px;
        }
      }
      .el-form-item__content {
        position: initial;
      }

      .el-form-item.is-error .el-input__wrapper {
        box-shadow: 0 0 0 1px var(--el-input-border-color, var(--el-border-color)) inset;
      }

      .el-form-item__error {
        top: 33px;
        left: 4px;
      }
    }
    :deep() {
      .left-panel {
        margin-bottom: 0;
      }
      .el-table__body-wrapper {
        .el-scrollbar {
          height: calc(100vh - 225px);
        }
        overflow-y: auto;
        &::-webkit-scrollbar {
          width: 0 !important;
          height: 0;
        }
      }
    }
    .dialog {
      width: 100vw;
      :deep() {
        .el-dialog__body {
          padding: 0 20px !important;
        }
      }
    }
    .warp {
      height: calc(100vh - 90px);
      display: flex;
      .left {
        .active {
          background-color: #f5f7fa;
        }
        height: 100%;
        border-right: 1px solid var(--el-border-color);
        padding: 15px 15px 0;
        background-color: #fff;
        width: 400px;
        .line-2 {
          margin: 20px 0;
          display: flex;
          align-items: center;
          justify-content: space-between;
        }
        .chapter {
          padding: 10px 15px 0;
          width: 100%;
          overflow-y: auto;
          border: 1px solid var(--el-border-color);
          &::-webkit-scrollbar {
            width: 3px;
            height: 0;
          }
          height: calc(100% - 102px);
          .item {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            word-break: break-all;
            word-wrap: break-word;
            &:hover {
              cursor: pointer;
            }
          }
        }
      }
      .right {
        margin-top: 15px;
        width: calc(100% - 400px);
        padding-left: 15px;
        background-color: #fff;
        // overflow-y: auto;
        overflow: hidden;
        &::-webkit-scrollbar {
          width: 3px;
          height: 0;
        }
      }
    }

    .content {
      height: 33px;
      display: flex;
      width: 100%;
      .left {
        height: max-content;
        border: 1px solid var(--el-border-color);
        background-color: #fff;
        width: 300px;
        border-radius: 4px;
        overflow: hidden;
        position: relative;
        z-index: 1991;
        border-bottom: none;
        &.show {
          border-bottom: 1px solid var(--el-border-color);
          .el-button {
            background: var(--el-color-primary);
            color: #fff !important;
          }
        }
        .my-select {
          &.el-button {
            border: none;
            border-bottom: 1px solid var(--el-border-color);
            border-radius: 0;
            height: 33px;
            color: var(--el-input-text-color, var(--el-text-color-regular));
          }
          width: 130px;
          border-bottom: 1px solid var(--el-border-color);
          &.border {
            border-right: 1px solid var(--el-border-color);
            border-left: 1px solid var(--el-border-color);
          }

          :deep() {
            .el-input__wrapper {
              border-radius: 0;
              box-shadow: none;
            }
          }
        }
        .chapter {
          width: 298px;
          height: calc(100vh - 160px);
          overflow-y: auto;
          background: #fff;
          &::-webkit-scrollbar {
            width: 3px;
            height: 0;
          }

          .chapter-item {
            padding-inline: 14px;
            margin-top: 10px;
            .el-table {
              height: 40px;
              background: #f6f5fe;
              border-radius: 4px;
              font-weight: 500;
              font-size: 14px;
              color: #1e1842;
              line-height: 40px;
              padding-inline: 10px;
              margin-bottom: 4px;
              .caret-wrapper {
                margin-left: -5px;
              }
            }
            .event-statistics {
              display: flex;
              height: 30px;
              align-items: center;
              justify-content: space-between;
              font-weight: 400;
              font-size: 14px;
              color: #4a4659;
              .size {
                color: #b7b5be;
                font-size: 13px;
              }
            }
          }
          :deep() {
            .el-tree {
              font-size: 13px;
            }
            .node {
              position: relative;
              width: 100%;
              display: flex;
              align-items: center;
              .node-txt {
                margin-left: 6px;
                max-width: 80px;
                overflow: hidden;
                text-overflow: ellipsis;
                word-break: break-all;
              }
              .node-num {
                position: absolute;
                z-index: 0;
                right: 5px;
                top: 3px;
                max-width: 50px;
                overflow: hidden;
                text-overflow: ellipsis;
                word-break: break-all;
                // vertical-align: text-top;
                height: 18px;
                margin-right: 5px;
                font-size: 12px;
                line-height: 15px;
                text-align: center;
              }
            }
            .overflow {
              font-weight: 500;
              font-size: 12px;
              color: #8275d6;
              float: right;
              cursor: pointer;
            }
            .event-statistics .el-checkbox__label {
              max-width: 155px;
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
            }
            .el-input__suffix-inner {
              .el-icon:last-child {
                order: -1;
              }
            }
          }
        }
      }
      .right {
        flex: 1;
        padding-left: 10px;
        background-color: #fff;
        overflow: hidden;
      }
    }
    .shortcut {
      width: 100%;
      display: flex;
      .word {
        margin-top: 4px;
        width: 100px;
      }
      .items {
        display: flex;
        width: 100%;
        align-items: center;
        flex-wrap: wrap;
        .item {
          :deep(.el-checkbox__label) {
            max-width: 220px;
            overflow: hidden;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            word-break: break-all;
            word-wrap: break-word;
            line-height: 30px;
          }
          // margin-bottom: 5px;
          display: flex;
          align-items: center;
          margin-right: 10px;
          padding: 0 10px;
          border: 1px solid var(--el-border-color);
        }
        :deep() {
          .el-checkbox__label {
            margin: 1px 6px 0 0;
          }
          .el-icon {
            &:hover {
              cursor: pointer;
            }
          }
        }
      }
    }

    .my-form {
      width: 100%;
      display: flex;

      .my-input {
        padding-left: 10px;
        flex: 1;
        :deep() {
          .el-row {
            width: 100%;
          }
          .el-form-item {
            width: 100%;
          }
          .el-form-item__content {
            margin-top: 3px;
          }
        }
      }
    }
  }
</style>
