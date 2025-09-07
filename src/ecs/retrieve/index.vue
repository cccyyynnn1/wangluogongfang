<script lang="ts">
  export default {
    name: 'Retrieve',
  }
</script>

<script setup lang="ts">
  import ApplicationConfiguration from '@/ecs/assets/components/application-configuration.vue'
  import NetworkLayer from './network-layer.vue'
  import LevelRule from './components/level-rule.vue'
  import Publish from './components/publish.vue'
  import AddSql from '../site/site-session/components/add-sql.vue'
  import PathTracing from './path-tracing.vue'
  import VabChart from '@/plugins/VabChart/index.vue'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import DbDetail from './components/detail-db.vue'
  import SmbDetail from './components/detail-smb.vue'
  import DnsDetail from './components/detail-dns.vue'
  import MailDetail from './components/detail-mail.vue'
  import FtpDetail from './components/detail-ftp.vue'
  import DetailFile from './components/detail-file.vue'
  import AlertDetail from '@/ecs/alert/components/alert-detail.vue'
  import DetailOther from '@/ecs/retrieve/components/detail-other.vue'
  import TraceabilityFieldStatistic from '@/ecs/alert/components/field-statistic.vue'
  import SessionInfo from '@/ecs/site/site-session/session-info.vue'
  import Download from '~/library/components/VabColumnBar/download.vue'
  import ToolBox from '~/library/components/VabColumnBar/toolbox.vue'
  import Favorites from './favorites.vue'
  import RetrieveHistory from './components/history.vue'
  import SearchSql from '@/components/search-sql/index.vue'
  import RetrieveWorkSpace from './components/work-space.vue'
  import RetrieveTraceability from './traceability.vue'
  import JsonPreview from '@/components/json-preview.vue'
  import { useTableCopy } from '@/utils'
  import {
    searchBySqlApi,
    publishRuleApi,
    UpdateHistoryApi,
    toShareApi,
    findHistoryApi,
    getAllUserNameAPI,
    getIndexCountAPI,
  } from '@/api-ecs/retrieve'
  import { getAllDisPlaysFiledApi } from '@/api-ecs/public'
  import { getSystemConfigApi } from '~/src/api-ecs/system'
  import { updateDisplayApi } from '@/api-ecs/custom-field'
  import {
    SearchBySqlParams,
    RetrieveIndexType,
    ShortcutListType,
    IndexTypeTpye,
    RetrieveHistoryItem,
    DowmloadType,
  } from '@/types/index'
  import { TableColumnItemType } from '/#/store'
  import { formatNstime, formatTime } from '@/utils/time'
  import numberFormatte from '@/utils/number'
  import dayjs from 'dayjs'
  import { useUserStore } from '@/store/modules/user'
  import { h } from 'vue'
  import { ElMessageBox, ElInput, ElSelect } from 'element-plus'
  import { useScroll } from '@vueuse/core'
  import { getTableCopyData } from '~/src/utils/transition'
  import { DownloadLogApi } from '~/src/api-ecs/toolbox'
  import EscUp from '@/components/up-and-down/esc-up.vue'
  import EscDown from '@/components/up-and-down/esc-down.vue'
  import numberFormatter from '~/src/utils/number'
  import IframePage from './iframe-page.vue'
  import _ from 'lodash'

  let abortController = new AbortController()
  const sqlComponentsRef = ref<InstanceType<typeof SearchSql>>()
  const $baseMessage: any = inject('$baseMessage')
  const userStore = useUserStore()
  const { getTableColumn, getAllIndexType, username } = userStore
  const tagType = ref<IndexTypeTpye[]>(getAllIndexType())
  const route = useRoute()
  const showFiledConfig = ref(false) // 显示表头字段配置
  const showPublish = ref(false) // 显示表头字段配置
  const traceabilityFieldRef = ref<InstanceType<typeof TraceabilityFieldStatistic>>()
  const isFirst = ref(true)
  const el = ref<HTMLElement | null>(null)

  const fieldDetailVisible = ref(false)
  const fieldDetailTitle = ref('')
  const fieldDetailValue = ref('')
  const showTableSelection = ref(false)
  const retrieveTableRef = ref()
  const dataZoomToggle = ref(false)

  // 检索信息表单数据
  const queryForm = reactive<SearchBySqlParams>({
    searchSql: '',
    indexType: 0,
    pageNum: 1,
    pageSize: 100,
    orderType: 'desc',
    orderField: 'requestTimeNs',
    startTime: '',
    endTime: '',
    scrollId: undefined,
    inputSql: '',
    filterSqlArr: '',
    sqlRelat: '',
  })

  const userDisPlaysFiled = ref()

  const activeName = ref<string>('applicationLayer')

  const DNSDict = ref()
  const nameList = ref([])
  // 获取表格序号
  const curIndex = computed(() => (queryForm.pageNum - 1) * queryForm.pageSize + 1)
  // 查询时间
  const timeDuration = ref('1hours')
  // 自定义时间
  const timeDate = ref()
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
      value: '6hours',
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
  // 柱状图实例
  const application_ref = ref()
  // 图表配置
  const retrieveChartOption = reactive({
    xAxis: {
      type: 'category',
      data: [] as string[],
    },
    yAxis: {
      type: 'value',
      splitNumber: 2,
      axisLabel: {
        formatter: function (value: number) {
          return numberFormatte.format(+value)
        },
      },
    },
    dataZoom: [
      {
        show: true,
        realtime: false,
        bottom: 8,
        xAxisIndex: [0],
        throttle: 200,
        brushSelect: true,
        backgroundColor: 'rgba(243, 241, 254, .5)',
        fillerColor: 'rgba(240,238,254,.1)',
        moveHandleStyle: { color: '#BFBADB' },
        handleStyle: {
          borderColor: '#BFBADB',
          color: '#BFBADB',
        },
        emphasis: {
          moveHandleStyle: { color: '#BFBADB' },
          handleStyle: {
            borderColor: '#BFBADB',
          },
        },
        selectedDataBackground: {
          areaStyle: {
            color: '#6954f0',
          },
          lineStyle: {
            color: '#6954f0',
          },
        },
        dataBackground: {
          lineStyle: {
            color: '#6954f0',
          },
          areaStyle: {
            color: 'rgba(243, 241, 254, .5)',
          },
        },
      },
    ],
    series: [
      {
        type: 'line',
        color: '#6954f0',
        lineStyle: {
          width: 1,
        },
        areaStyle: {
          color: new VabChart.graphic.LinearGradient(0, 0, 0, 1, [
            {
              offset: 0,
              color: '#D1DCFF',
            },
            {
              offset: 1,
              color: '#fff',
            },
          ]),
        },
        smooth: 0.6,
        showSymbol: false,
        data: [] as number[],
      },
    ],
    grid: {
      top: 10,
      left: '2%',
      right: 10,
      bottom: 0,
      height: 60,
      containLabel: true,
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow',
      },
    },
  })
  // 是否加载
  const listLoading = ref(false)
  const showLevelRule = ref(false) //规则列表
  const showHistory = ref(false)
  const work_space_ref = ref<InstanceType<typeof RetrieveWorkSpace>>()
  const sqlColums = ref<TableColumnItemType[]>([])
  const SecurityId = 9999999
  const spaceId = computed(() => work_space_ref.value?.spaceId || SecurityId)
  const favoritesDisable = ref(true)
  const hasHistory = ref(false)
  const historyData = ref()

  const indexTypeCount = ref()

  // 详情弹窗Visible 除HTTP检索使用
  const detailVisible = ref(false)
  // HTTP检索单独使用！
  const sessionInfoVisible = ref(false)
  const detailType = ref(1)
  const detailVal = ref()
  // 表格数据
  const layout = ref('sizes, prev, pager, next, jumper')
  // 检索返回的数据
  const queryPage = reactive({
    total: 0,
    title: '',
    listDate: [] as object[],
    queryTimes: '0.0',
  })
  // 选择图表时间
  const subChartSearchData = ref<string[]>([])
  // 选择图表查询的总数
  const subChartSearchTotal = ref(0)
  const moduleEnable = ref(false)
  const show = ref(false) // 是否显示添加sql页面

  const mode = ref('add') // 是否显示添加sql页面
  const multipleSelection = ref<any[]>([])
  const currentRow = ref<ShortcutListType>()

  const shortcutList = ref<ShortcutListType[]>([])

  const shortcutSQL = ref()

  const noTakeEffect = [
    'requestTimeNs',
    'responseTimeNs',
    'responseTime',
    'requestContentLength',
    'requestId',
    'responseContentLength',
    'probeIds',
    'deviceId',
    'startTimeNs',
    'responseCode',
    'kafkakey',
    'date',
    'id',
    'clientAssets',
    'serverAssets',
    'clientArea',
    'serverArea',
    'direction',
  ]

  const detailDom: { [key: number]: any } = {
    2: DnsDetail,
    3: DbDetail,
    4: FtpDetail,
    5: SmbDetail,
    6: MailDetail,
    10: AlertDetail,
    12: DetailOther,
    13: DetailOther,
    14: DetailOther,
    15: DetailOther,
    16: DetailOther,
    17: DetailOther,
    18: DetailOther,
    19: DetailOther,
    20: DetailOther,
    21: DetailOther,
    25: DetailOther,
    26: DetailOther,
    27: DetailFile,
    28: DetailOther,
    32: DetailOther,
  }

  const showPrevBtn = ref(false)

  const showNextBtn = ref(false)

  // 或和且互斥功能
  const SQLMode = ref('且')
  const exclusive = ref(false)
  watch(
    () => exclusive.value,
    () => {
      if (exclusive.value) {
        shortcutList.value.forEach((_: any, index: any) => {
          shortcutList.value[index].enable = false
        })
      }
    }
  )

  const handleHiddenSelection = () => {
    showTableSelection.value = false
    multipleSelection.value = []
    retrieveTableRef.value?.clearSelection()
  }
  const handleGetIndexCount = async (obj: SearchBySqlParams) => {
    const { data } = await getIndexCountAPI(obj)
    indexTypeCount.value = data.reduce((a, b) => Object.assign(a, b), {})
  }

  const handleCheckboxChange = (id: string) => {
    if (exclusive.value) {
      shortcutList.value.forEach((_: any, index: any) => {
        if (shortcutList.value[index].id != id) {
          shortcutList.value[index].enable = false
        }
      })
      searchBySqlHandel()
    }
  }
  const handleMerge2Sql = () => {
    const checkedFilter = []
    const otherFilter = []
    for (const element of shortcutList.value) {
      if (element.enable) {
        checkedFilter.push(element)
        continue
      }
      otherFilter.push(element)
    }
    if (!checkedFilter.length) return
    const inputSql = queryForm.searchSql ? `( ${queryForm.searchSql} ) ` : ''
    const checkedStr = checkedFilter
      .map(
        ({ label, relation, value }) =>
          `${label} ${relation.toLowerCase().endsWith('exists') ? `${relation}` : `${relation} "${value}"`}`
      )
      .join(SQLMode.value === '且' ? ' and ' : ' or ')
    const mergeSql = `${inputSql ? `${inputSql} and ( ${checkedStr} )` : checkedStr}`
    sqlComponentsRef.value?.changeSql(mergeSql)
    shortcutList.value = otherFilter
  }
  const probeSearchSqlIsIPOrPort = () => {
    // IPV4
    const ipv4Regex =
      /^(?!0)(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])$/
    // IPV6
    const ipv6Regex =
      /^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$/
    // 端口
    const portRegex = /^(?:[1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$/
    if (!queryForm.searchSql?.trim()) return ''
    if (ipv4Regex.test(queryForm.searchSql) || ipv6Regex.test(queryForm.searchSql)) {
      return `源IP = "${queryForm.searchSql}" or 目的IP = "${queryForm.searchSql}"`
    }
    if (portRegex.test(queryForm.searchSql)) {
      return `源端口 = "${queryForm.searchSql}" or 目的端口 = "${queryForm.searchSql}"`
    }
    return queryForm.searchSql
  }

  // 上一条或下一条
  const changeCurrentItemEvent = (val: boolean) => {
    let index = queryPage.listDate.findIndex((item: any) => {
      return item.id == detailVal.value.id
    })
    detailVal.value = undefined
    if (val) {
      if (index >= queryPage.listDate.length - 5 && queryPage.listDate.length! != queryPage.total) {
        if (subChartSearchData.value.length > 0) {
          nextTick(async () => await getSubChartDataHandel(false))
        } else {
          nextTick(async () => await searchBySql(false))
        }
      }
      ++index
      detailVal.value = { ...queryPage.listDate[index], indexType: detailType.value }
    } else {
      --index
      detailVal.value = { ...queryPage.listDate[index], indexType: detailType.value }
    }
    index !== 0 ? (showPrevBtn.value = true) : (showPrevBtn.value = false)
    index !== queryPage.total - 1 ? (showNextBtn.value = true) : (showNextBtn.value = false)
    if (detailVisible.value == true) {
      reloadingEvent()
    }
    detailDomRef.value?.getCode && detailDomRef.value?.getCode()
    detailDomRef.value?.getCode && detailDomRef.value?.getAssetName()
  }

  const tableColumn = ref<TableColumnItemType[]>([])
  // sql更新回调
  const handleUpdateCallback = (res: ShortcutListType) => {
    const flag = shortcutList.value.findIndex((item: any) => {
      return item?.id == res.id
    })
    if (flag < 0) {
      shortcutList.value.push(res)
    } else {
      shortcutList.value[flag] = JSON.parse(JSON.stringify(res))
    }
  }

  // 关闭
  const closeSql = (val: boolean) => {
    show.value = false
  }

  // 清空
  const handleEmpty = () => {
    shortcutList.value = []
  }

  // 添加修改sql
  const handleUpdate = ({ remark, data }: { remark: string; data?: ShortcutListType }) => {
    mode.value = remark
    show.value = true
    currentRow.value = data
  }

  // 删除sql
  const handleClear = (id: string) => {
    const _index = shortcutList.value.findIndex((item: any) => {
      return item.id == id
    })
    shortcutList.value.splice(_index, 1)
  }

  // 打开详情
  const handleInfo = (row: any) => {
    detailVal.value = { ...row, indexType: detailType.value }
    const index = queryPage.listDate.findIndex((item: any) => {
      return item.id == detailVal.value.id
    })
    index !== 0 ? (showPrevBtn.value = true) : (showPrevBtn.value = false)
    index !== queryPage.total - 1 ? (showNextBtn.value = true) : (showNextBtn.value = false)
    if (detailType.value === 1) {
      return (sessionInfoVisible.value = true)
    }
    detailVisible.value = true
  }

  const resetData = () => {
    queryForm.scrollId = undefined
    queryPage.listDate = []
  }

  // 多选项改变
  const setSelectRows = (val: any) => {
    multipleSelection.value = val
  }

  // 页容量改变
  // const handleSizeChange = (val: number) => {
  //   searchBySql(false)
  // }

  // // 页面改变
  // const handleCurrentChange = () => {
  //   searchBySql(false)
  // }

  // 表头字段配置
  const tableHeadConfig = (val: boolean) => {
    showFiledConfig.value = val
  }

  // 发布
  const publishEvent = (val: boolean) => {
    showPublish.value = val
  }

  const searchBySqlHandel = () => {
    queryForm.pageNum = 1
    isFirst.value = false
    resetData()
    searchBySql()
  }
  // 检索信息
  async function searchBySql(updateChart = true) {
    if (sqlComponentsRef.value && !sqlComponentsRef.value?.getValidateSql()) return
    if (dataZoomToggle.value) {
      subChartSearchData.value = []
      subChartSearchTotal.value = 0
      queryForm.pageNum = 1
      queryForm.pageSize = 100
    }
    shortcutSQL.value = ''
    const arr = shortcutList.value.filter((item: any) => {
      return item.enable
    })
    if (arr.length > 0) {
      let shortcut = '('
      arr.forEach((td: any, index: number) => {
        if (index == arr.length - 1) {
          shortcut = `${shortcut} ${arr[index].label} ${arr[index].relation} ${
            // @ts-ignore
            arr[index].value || arr[index].value == 0 || arr[index].value == '0' ? `"${arr[index].value}"` : ''
          })`
        } else {
          shortcut =
            `${shortcut} ${arr[index].label} ${arr[index].relation} ${
              // @ts-ignore
              arr[index].value || arr[index].value == 0 || arr[index].value == '0' ? `"${arr[index].value}"` : ''
            } ` +
            `${SQLMode.value == '且' ? 'and' : 'or'}` +
            ` `
        }
      })
      if (queryForm.searchSql?.trim()) {
        shortcutSQL.value = `(${probeSearchSqlIsIPOrPort()})` + ` ` + `and` + ` ${shortcut}`
      } else {
        shortcutSQL.value = shortcut
      }
    } else {
      if (queryForm.searchSql?.trim()) {
        shortcutSQL.value = JSON.parse(JSON.stringify(probeSearchSqlIsIPOrPort()))
      }
    }
    listLoading.value = true
    getQueryDate()
    if (updateChart) {
      resetData()
    }
    try {
      let filterSqlArr = ''
      if (shortcutList.value.length > 0) {
        filterSqlArr = JSON.stringify(shortcutList.value)
      }
      let sqlRelat = ''
      sqlRelat = JSON.stringify({ SQLMode: SQLMode.value, exclusive: exclusive.value })
      if (shortcutSQL.value) {
        handleGetIndexCount({
          ...queryForm,
          searchSql: shortcutSQL.value,
          workspaceId: spaceId.value,
          inputSql: queryForm.searchSql,
          filterSqlArr,
        })
      } else {
        indexTypeCount.value = undefined
      }
      const {
        data: { queryTime, resList, sumaryMap, total, scrollId },
        code,
      } = await searchBySqlApi(
        {
          ...queryForm,
          workspaceId: spaceId.value,
          searchSql: shortcutSQL.value,
          inputSql: queryForm.searchSql,
          filterSqlArr,
          sqlRelat,
        },
        abortController.signal
      )
      const my_arr = resList || []
      queryForm.scrollId = scrollId
      queryPage.listDate.push(...my_arr)
      queryPage.queryTimes = (queryTime / 1000 + 0.1).toFixed(1)
      queryPage.total = total
      if (updateChart) {
        sqlComponentsRef.value?.changeHistories(shortcutSQL.value)
        if (!sumaryMap) return
        const { count, date } = JSON.parse(sumaryMap)
        retrieveChartOption.xAxis.data = date || []
        retrieveChartOption.series[0].data = count || []
      }
      // chartCurData.value = queryPage.listDate.length > 0 && queryPage.listDate[0]
      detailType.value = queryForm.indexType
      formatColum(queryForm.indexType)
      if (dataZoomToggle.value) {
        application_ref.value?.chart?.clear()
        application_ref.value?.chart?.setOption(retrieveChartOption)
      }
    } finally {
      favoritesDisable.value = false
      listLoading.value = false
      dataZoomToggle.value = false
    }
  }
  const updateCollect = async (row: RetrieveHistoryItem) => {
    const { code } = await UpdateHistoryApi({ ...row, collectStatus: !row.collectStatus })
    $baseMessage(row.collectStatus ? '已取消收藏' : '收藏成功', 'success', 'vab-hey-message-success')
    hasHistory.value = false
  }
  const changeDNSData = () => {
    queryPage.listDate.forEach((item: any) => {
      Object.keys(item).forEach((td: any) => {
        if (td == 'requestType' || td == 'requestClass') {
          const res = item[td]
          item[td] = DNSDict.value[td][res]
        }
      })
    })
  }
  const getThreatLevel = (type: string) => {
    const obj: {
      [key: string]: string
    } = {
      高危: 'high',
      中危: 'mid',
      低危: 'low',
      危急: 'critical',
    }
    return obj[type] || 'default'
  }
  const showTraceabilityField = (column: TableColumnItemType) => {
    const { fieldNameCn, fieldNameEn } = column
    const { searchSql, indexType, startTime, endTime } = queryForm
    traceabilityFieldRef.value?.initData({
      title: fieldNameCn,
      query: {
        searchSql: shortcutSQL.value,
        indexType,
        startTime,
        endTime,
        aggregationFields: fieldNameEn,
        whiteType: undefined,
      },
      column,
    })
  }
  // 获取用户所有字段
  const getAllDisPlaysFiled = async () => {
    const { data } = await getAllDisPlaysFiledApi()
    userDisPlaysFiled.value = data
    formatColum(queryForm.indexType)
  }

  function formatDate(row: any, key: string) {
    if (!row[key]) return row[key]
    const time = row[key] / 1000000
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss.SSS')
  }

  async function formatColum(type: RetrieveIndexType) {
    const columnData = await getTableColumn(queryForm.indexType)
    sqlColums.value = columnData
    if (userDisPlaysFiled.value) {
      const userColumnData = userDisPlaysFiled.value[type] as number[]
      if (userColumnData?.length > 0) {
        tableColumn.value = userColumnData.map((key: any) =>
          columnData?.find((item) => item?.id === key)
        ) as TableColumnItemType[]
      } else {
        tableColumn.value = []
      }
    }
  }

  async function configurationHandel(showField: TableColumnItemType[]) {
    tableColumn.value = showField
    const { msg } = await updateDisplayApi({
      indexType: queryForm.indexType,
      displayIds: showField.map((i) => i.id),
    })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    getAllDisPlaysFiled()
  }

  function changeCellStyle(title: string) {
    if (['客户端口', '服务端口', '请求方式', '状态码', '源端口', '目的端口'].includes(title)) {
      return '120'
    } else if (['客户端IP', '服务端IP', 'HOST', 'XFF'].includes(title)) {
      return '120'
    } else if (['请求时间', '响应时长'].includes(title)) {
      return '160'
    } else {
      return '180'
    }
  }
  async function trafficModelHandle(data: { searchTime: number; ruleName: string }) {
    const _shortcutSql = shortcutList.value
      .filter((i) => i.enable)
      .map(({ value, label, relation }) => `${label} ${relation} ${value}`)
      .join(' and ')
    const sql = [queryForm.searchSql ? `( ${queryForm.searchSql} )` : '', _shortcutSql].filter(Boolean).join(' and ')
    if (!data.ruleName) return $baseMessage('请填写规则名称', 'error', 'vab-hey-message-error')
    if (!data.searchTime) return $baseMessage('请填写健康时长', 'error', 'vab-hey-message-error')
    if (!sql) return $baseMessage('请填写查询SQL', 'error', 'vab-hey-message-error')
    const query = {
      ...data,
      ...queryForm,
      searchSql: sql,
      displayFields: tableColumn.value.map((i) => i.fieldNameEn),
    }
    try {
      const { msg } = await publishRuleApi(query)
      $baseMessage(msg, 'success', 'vab-hey-message-success')
    } catch (error) {
      $baseMessage('发布失败', 'error', 'vab-hey-message-error')
    }
    publishEvent(false)
  }

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
  // const getDnsDict = async () => {
  //   const { data } = await getDnsDictApi()
  //   DNSDict.value = data
  // }

  const useTableCopyEvent = (row: any, column: any, cell: any, event: any) => {
    if (
      (column.property == 'serverAssets' && row.serverAssets == undefined) ||
      (column.property == 'clientAssets' && row.clientAssets == undefined)
    )
      return
    const arr = getTableCopyData({ tableColumn: tableColumn.value, row, column, mothod: handleUpdateCallback })
    useTableCopy(row, column, cell, event, arr)
  }

  const reloadSearchHandle = (data: RetrieveHistoryItem, hasTime: boolean) => {
    let { searchSql } = data
    const { searchEdTime, searchStTime, indexType, inputSql, filterSqlArr, sqlRelat } = data
    if (inputSql) {
      searchSql = inputSql
    } else {
      if (filterSqlArr) {
        searchSql = ''
      }
    }
    if (sqlRelat) {
      const res = JSON.parse(sqlRelat)
      SQLMode.value = res.SQLMode
      exclusive.value = res.exclusive
    }
    if (filterSqlArr) {
      shortcutList.value = JSON.parse(filterSqlArr)
    } else {
      shortcutList.value = []
    }
    queryForm.pageNum = 1
    queryForm.pageSize = 100
    queryForm.startTime = hasTime ? searchStTime : ''
    queryForm.endTime = hasTime ? searchEdTime : ''
    queryForm.searchSql = searchSql
    timeDuration.value = hasTime ? 'user-defined' : '1hours'
    timeDate.value = hasTime ? [searchStTime, searchEdTime] : null
    sqlComponentsRef.value?.changeSql(queryForm.searchSql)
    if (queryForm.indexType === indexType) {
      searchBySql()
    } else {
      queryForm.indexType = indexType
    }
  }
  const shareHandle = () => {
    if (favoritesDisable.value) return
    getAllUserNameAPI().then((res) => {
      nameList.value = res.data || []
    })
    const rNames = ref<string[]>([])
    const remarks = ref('')
    let filterSqlArr = ''
    if (shortcutList.value.length > 0) {
      filterSqlArr = JSON.stringify(shortcutList.value)
    }
    let sqlRelat = ''
    sqlRelat = JSON.stringify({ SQLMode: SQLMode.value, exclusive: exclusive.value })
    ElMessageBox({
      title: '分享给',
      confirmButtonText: '分享',
      customClass: 'sharedMessageBox',
      // closeOnClickModal: false,
      message: () =>
        h('div', null, [
          h(
            ElSelect,
            {
              modelValue: rNames.value,
              placeholder: '请选择被分享人',
              multiple: false,
              class: 'abc',
              style: { width: '100%' },
              'onUpdate:modelValue': (val: string) => {
                rNames.value = []
                rNames.value.push(val)
              },
            },
            () =>
              nameList.value.map((item: string) => {
                return h(ElSelect.Option, {
                  key: item,
                  label: item,
                  value: item,
                })
              })
          ),
          h(ElInput, {
            modelValue: remarks.value,
            type: 'textarea',
            placeholder: '请输入分享备注',
            rows: 5,
            style: 'margin-top: 10px',
            resize: 'none',
            'onUpdate:modelValue': (val: string) => {
              remarks.value = val
            },
          }),
        ]),
      beforeClose: async (action, instance, done) => {
        if (action === 'confirm') {
          const { indexType, endTime, startTime } = queryForm
          instance.confirmButtonLoading = true
          const { code } = await toShareApi({
            rNames: rNames.value.length === 0 ? '所有人' : rNames.value.join(','),
            sName: username,
            searchStTime: startTime,
            searchSql: shortcutSQL.value,
            inputSql: queryForm.searchSql,
            filterSqlArr,
            sqlRelat,
            searchEdTime: endTime,
            indexType,
            workspaceId: spaceId.value,
            remarks: remarks.value,
          })
          $baseMessage('分享成功', 'success', 'vab-hey-message-success')
          done()
        } else {
          done()
        }
      },
    }).catch(() => {})
  }
  const favoritesHandle = async () => {
    if (favoritesDisable.value || !shortcutSQL.value.length) return
    if (hasHistory.value) {
      updateCollect(historyData.value)
      return
    }
    findHistoryApi({
      searchSql: shortcutSQL.value,
      indexType: queryForm.indexType,
      collectStatus: true,
      workspaceId: spaceId.value,
    }).then((res: any) => {
      if (res.data) {
        hasHistory.value = true
        historyData.value = res.data
      } else {
        hasHistory.value = false
      }
    })
    const { indexType, endTime, startTime } = queryForm
    let filterSqlArr = ''
    if (shortcutList.value.length > 0) {
      filterSqlArr = JSON.stringify(shortcutList.value)
    }
    let sqlRelat = ''
    sqlRelat = JSON.stringify({ SQLMode: SQLMode.value, exclusive: exclusive.value })
    const remarks = ref('')
    ElMessageBox({
      title: '收藏',
      confirmButtonText: '确定',
      message: () =>
        h(ElInput, {
          modelValue: remarks.value,
          type: 'textarea',
          placeholder: '请输入收藏备注',
          rows: 5,
          resize: 'none',
          'onUpdate:modelValue': (val: string) => {
            remarks.value = val
          },
        }),
      beforeClose: async (action, instance, done) => {
        if (action === 'confirm') {
          const { code } = await UpdateHistoryApi({
            collectStatus: true,
            searchStTime: startTime,
            searchSql: shortcutSQL.value,
            inputSql: queryForm.searchSql,
            filterSqlArr,
            sqlRelat,
            searchEdTime: endTime,
            indexType,
            workspaceId: spaceId.value,
            remarks: remarks.value,
          })
          $baseMessage('收藏成功', 'success', 'vab-hey-message-success')
          done()
        } else {
          done()
        }
      },
    }).catch(() => {})
  }
  // 选择图表时间查询数据
  const getSubChartDataHandel = async (remark = true) => {
    if (remark) {
      resetData()
    }
    try {
      listLoading.value = true
      const {
        data: { resList, total, scrollId },
      } = await searchBySqlApi({
        ...queryForm,
        searchSql: shortcutSQL.value,
        workspaceId: spaceId.value,
        startTime: subChartSearchData.value[0],
        endTime: subChartSearchData.value[1],
      })
      queryForm.scrollId = scrollId
      queryPage.total = total
      queryPage.listDate.push(...resList)
      // chartCurData.value = queryPage.listDate.length > 0 && queryPage.listDate[0]
      formatColum(queryForm.indexType)
    } finally {
      listLoading.value = false
      dataZoomToggle.value = false
    }
  }

  // 日志
  const handleDownloadExcel = async () => {
    if (multipleSelection.value.length > 0) {
      const tHeader: string[] = []
      const filterVal: any = []
      try {
        tableColumn.value.forEach(({ fieldNameCn, fieldNameEn }) => {
          tHeader.push(fieldNameCn)
          filterVal.push(fieldNameEn)
        })
        listLoading.value = true
        import('@/utils/excel').then((excel) => {
          const list = multipleSelection.value.length ? multipleSelection.value : queryPage.listDate
          const data = formatJson(filterVal, list)
          excel.export_json_to_excel({
            header: tHeader,
            data,
            filename: `${formatTime(new Date().getTime())}-日志导出/数据导出`,
            autoWidth: true,
            bookType: 'xlsx',
          })
          listLoading.value = false
        })
      } catch (error) {
        console.log(error)
        listLoading.value = false
      }
    } else {
      getQueryDate()
      const obj: DowmloadType = {
        downloadCnd: shortcutSQL.value,
        edTime: dayjs(queryForm.endTime).valueOf(),
        indexType: queryForm.indexType,
        stTime: dayjs(queryForm.startTime).valueOf(),
        type: 0,
        dataType: 1,
      }
      try {
        const { msg } = await DownloadLogApi({ ...obj }, 1)
        $baseMessage(msg, 'success', 'vab-hey-message-success')
      } catch (error) {
        $baseMessage('下载失败', 'error', 'vab-hey-message-error')
      }
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
      case 'startTimeNs':
      case 'requestTimeNs':
        return formatNstime(val[key])
      default:
        return val[key] || ''
    }
  }
  const showFieldHandle = (title: string, val: string) => {
    fieldDetailVisible.value = true
    fieldDetailTitle.value = title
    fieldDetailValue.value = val
  }
  const searchhandle = (val: string) => {
    searchBySqlHandel()
  }
  const getConfidence = (code: string) => {
    return ['低', '中', '高'].findIndex((i) => code === i) + 1
  }
  watch(
    () => timeDuration.value,
    (newData, oldData) => {
      dataZoomToggle.value = true
      getQueryDate()
    },
    {
      immediate: true,
    }
  )
  watch(
    () => userStore.indexTypeList,
    () => {
      tagType.value = userStore.indexTypeList ? userStore.indexTypeList : []
    },
    {
      immediate: true,
    }
  )
  watch(
    () => queryForm.indexType,
    () => {
      dataZoomToggle.value = true
      if (queryForm.indexType !== 1) {
        queryForm.orderField = 'startTimeNs'
      } else {
        queryForm.orderField = 'requestTimeNs'
      }
      abortController?.abort()
      abortController = new AbortController()
      searchBySqlHandel()
    }
  )
  watch(
    () => [shortcutList.value, queryForm.searchSql],
    () => {
      favoritesDisable.value = true
    },
    { deep: true }
  )

  const formatUserTime = () => {
    if (timeDate.value) {
      const [startDate, endDate] = timeDate.value
      queryForm.endTime = dayjs(endDate).format('YYYY-MM-DD HH:mm:ss')
      queryForm.startTime = dayjs(startDate).format('YYYY-MM-DD HH:mm:ss')
    }
  }

  watchEffect(() => {
    if (timeDate.value) {
      formatUserTime()
    }
  })
  const getModelStatus = async () => {
    const {
      data: { module_config },
    } = await getSystemConfigApi({ keys: 'module_config' })
    moduleEnable.value = JSON.parse(module_config.value).isEnable || false
  }
  const data = ref()
  const reloading = ref(false)
  const reloadingEvent = () => {
    reloading.value = true
    setTimeout(() => {
      reloading.value = false
    }, 400)
  }
  const detailDomRef = ref()
  const handleKeyDown = (e: any) => {
    if (e.keyCode == 37 && showPrevBtn.value == true) {
      changeCurrentItemEvent(false)
    } else if (e.keyCode == 39 && showNextBtn.value == true) {
      changeCurrentItemEvent(true)
    }
  }

  const remark = ref<'history' | 'favorites' | 'share' | 'shareto'>('history')
  const infoObj = {
    httpLogin: 1,
    sql: 3,
    mail: 6,
    msrdp: 16,
    sshLog: 18,
  }

  onMounted(() => {
    getModelStatus()
    const query = route.query
    window.addEventListener('keyup', handleKeyDown)
    if (query.info) {
      const res = JSON.parse(decodeURIComponent(query.info as string)) as any
      queryForm.searchSql = res.sql
      queryForm.indexType = res.indexType
      sqlComponentsRef.value?.changeSql(queryForm.searchSql)
      if (res?.timeRanges) {
        if (res.protocol) {
          // @ts-ignore
          queryForm.indexType = infoObj[res.protocol]
        }
        const [start, end] = res.timeRanges
        timeDate.value = [dayjs(start), dayjs(end)]
        timeDuration.value = 'user-defined'
      }
      // if (res.indexType === 1) {
      //   searchBySql()
      // }
    } else {
      queryForm.indexType = 1
    }
    application_ref.value.chart.on('datazoom', () => {
      const { dataZoom, xAxis } = application_ref.value.chart.getModel().option
      const startValue = dataZoom[0].startValue
      const endValue = dataZoom[0].endValue
      const startLable = xAxis[0].data[startValue]
      const endLable = xAxis[0].data[endValue]
      timeDuration.value = 'user-defined'
      timeDate.value = [startLable, endLable]
      subChartSearchData.value = [startLable, endLable]
      queryForm.pageNum = 1
      queryForm.pageSize = 100
      dataZoomToggle.value = false
      getSubChartDataHandel()
    })
    // getDnsDict()
    getAllDisPlaysFiled()

    el.value = document.querySelector('.my-table .el-scrollbar__wrap')
    const { arrivedState } = useScroll(el.value)
    data.value = arrivedState
    activeName.value = (route.query.params as string) || 'applicationLayer'

    if (query.to) {
      showHistory.value = true
      remark.value = 'shareto'
    }
  })

  watch(
    () => route.query,
    () => {
      if (route.query.infos) {
        const res = JSON.parse(decodeURIComponent(route.query.infos as string)) as any
        queryForm.workspaceId = res.workspaceId
        const sqlRelat = res.sqlRelat
        const filterSqlArr = res.filterSqlArr
        const inputSql = res.inputSql
        let searchSql = res.searchSql
        if (inputSql) {
          searchSql = inputSql
        } else {
          if (filterSqlArr) {
            searchSql = ''
          }
        }
        if (sqlRelat) {
          const res = JSON.parse(sqlRelat)
          SQLMode.value = res.SQLMode
          exclusive.value = res.exclusive
        }
        if (filterSqlArr) {
          shortcutList.value = JSON.parse(filterSqlArr)
        } else {
          shortcutList.value = []
        }
        queryForm.searchSql = searchSql
        if (res?.timeRanges) {
          const [start, end] = res.timeRanges
          timeDate.value = [dayjs(start), dayjs(end)]
          timeDuration.value = 'user-defined'
        }
        queryForm.indexType = res.indexType
      }
    },
    { deep: true, immediate: true }
  )

  watch(
    () => data.value,
    () => {
      if (data.value.bottom && !listLoading.value) {
        if (queryPage.total === queryPage.listDate.length) return
        if (subChartSearchData.value.length > 0) {
          nextTick(() => getSubChartDataHandel(false))
        } else {
          nextTick(() => searchBySql(false))
        }
      }
    },
    { deep: true }
  )

  const getLevel = (type: number) => {
    const obj: {
      [key: string]: string
    } = {
      严重: 'high',
      一般: 'mid',
      普通: 'low',
    }
    return obj[type] || 'default'
  }

  const handlerExChange = (str: number) => {
    if (str == 0) {
      return '普通'
    } else if (str == 1) {
      return '一般'
    } else if (str == 2) {
      return '严重'
    }
  }
</script>

<template>
  <div class="retrieve-container">
    <el-tabs v-model="activeName" style="height: 100%">
      <el-tab-pane class="flex-col" label="应用层会话" lazy name="applicationLayer" style="height: 100%">
        <!-- 头部检索栏 -->
        <vab-query-form style="width: 100%">
          <vab-query-form-left-panel :span="16">
            <el-form inline>
              <el-form-item>
                <el-select v-model="queryForm.indexType">
                  <el-option
                    v-for="item in tagType"
                    :key="item.value"
                    :label="item.label"
                    style="padding-inline: 10px 15px; width: 250px"
                    :value="item.value"
                  >
                    {{ item.label }}
                    <span v-if="indexTypeCount" style="float: right">
                      {{
                        indexTypeCount?.[item.value] > 0
                          ? numberFormatter.format(indexTypeCount?.[item.value])
                          : '暂无数据'
                      }}
                    </span>
                  </el-option>
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-select v-model="timeDuration" style="margin: 0 10px">
                  <el-option
                    v-for="item in timeDuratioOptions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
              </el-form-item>
              <el-form-item v-if="timeDuration === 'user-defined'">
                <vab-date-time-picker v-model="timeDate" />
              </el-form-item>
            </el-form>
          </vab-query-form-left-panel>
          <vab-query-form-right-panel :span="8">
            <!-- 右侧操作按钮  v-permissions="['Admin']" -->
            <span>总用时{{ queryPage.queryTimes }}秒</span>
            <el-divider direction="vertical" />
            <el-tooltip class="item" content="回溯记录" effect="dark" placement="top">
              <vab-icon icon="time-line" style="font-size: 20px; color: #b3b9c8" @click="showHistory = true" />
            </el-tooltip>
            <el-divider direction="vertical" />
            <el-tooltip class="item" content="分享" effect="dark" placement="top">
              <vab-icon icon="share-line" style="font-size: 20px; color: #b3b9c8" @click="shareHandle" />
            </el-tooltip>
            <el-divider direction="vertical" />
            <!-- <div v-if="queryForm.indexType === 1" v-permissions="['Admin']" style="display: inline-block">
              <el-tooltip content="发布" effect="dark" placement="top">
                <el-icon :size="20" style="vertical-align: middle; color: #b3b9c8" @click="publishEvent(true)">
                  <Promotion />
                </el-icon>
              </el-tooltip>
              <el-divider direction="vertical" />
            </div> -->
            <!-- 显示表头字段配置 -->
            <el-tooltip class="item" content="字段配置" effect="dark" placement="top">
              <vab-icon icon="edit-2-line" style="font-size: 20px; color: #b3b9c8" @click="tableHeadConfig(true)" />
            </el-tooltip>
            <el-divider direction="vertical" />
            <el-tooltip class="item" content="等级规则" effect="dark" placement="top">
              <vab-icon icon="article-line" style="font-size: 20px; color: #b3b9c8" @click="showLevelRule = true" />
            </el-tooltip>
          </vab-query-form-right-panel>
          <vab-query-form-right-panel :span="24">
            <el-form class="my-form" inline :model="queryForm" @submit.prevent>
              <el-form-item class="my-input">
                <SearchSql
                  ref="sqlComponentsRef"
                  :index-type="queryForm.indexType"
                  :model-value="queryForm.searchSql"
                  @on-change="(str) => (queryForm.searchSql = str)"
                  @onSearch="searchhandle"
                />
                <!-- 收藏 -->
                <div class="favoritesBtn" :class="{ hasHistory }" @click="favoritesHandle">
                  <vab-icon :icon="hasHistory ? 'star-fill' : 'star-line'" />
                </div>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="searchBySqlHandel">检索</el-button>
                <el-dropdown trigger="click">
                  <el-button type="default" @click="showTableSelection = true">
                    导出日志
                    <el-icon class="el-icon--right"><arrow-down /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item :disabled="multipleSelection.length === 0" @click="handleDownloadExcel">
                        导出选中
                      </el-dropdown-item>
                      <el-dropdown-item @click="handleDownloadExcel">导出所有</el-dropdown-item>
                      <el-dropdown-item @click="handleHiddenSelection">隐藏多选框</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </el-form-item>
            </el-form>
          </vab-query-form-right-panel>
          <vab-query-form-left-panel :span="24">
            <div class="shortcut">
              <div class="word">过滤条件：</div>
              <el-radio-group v-model="SQLMode" class="SQLMode" size="small">
                <el-radio-button :disabled="exclusive" label="或" />
                <el-radio-button :disabled="exclusive" label="且" />
              </el-radio-group>
              <div class="containment">
                <el-checkbox v-model="exclusive" />
                &nbsp;&nbsp;互斥
              </div>
              <div class="items">
                <div v-for="item in shortcutList" :key="item.id" class="item">
                  <el-checkbox v-model="item.enable" @change="handleCheckboxChange(item.id)">
                    <span style="my-span">{{ item.label }} {{ item.relation }} {{ `"${item.value ?? ''}"` }}</span>
                  </el-checkbox>
                  <el-popover placement="bottom" trigger="hover" width="160">
                    <template #reference>
                      <el-icon><ArrowDown /></el-icon>
                    </template>
                    <template #default>
                      <el-button link @click="handleUpdate({ remark: 'edit', data: item })">
                        <el-icon><Edit /></el-icon>
                        &nbsp;编辑
                      </el-button>
                      <el-button link @click="handleClear(item.id)">
                        <el-icon><Close /></el-icon>
                        &nbsp;删除
                      </el-button>
                    </template>
                  </el-popover>
                </div>
                <el-button
                  class="items-btn"
                  :disabled="listLoading"
                  plain
                  size="small"
                  type="primary"
                  @click="handleUpdate({ remark: 'add' })"
                >
                  <el-icon><Plus /></el-icon>
                  添加
                </el-button>
                <el-button
                  class="items-btn"
                  :disabled="shortcutList.length === 0"
                  plain
                  size="small"
                  type="primary"
                  @click="handleMerge2Sql"
                >
                  <el-icon>
                    <svg
                      id="图层_1"
                      version="1.1"
                      viewBox="0 0 12 12"
                      x="0px"
                      xml:space="preserve"
                      xmlns="http://www.w3.org/2000/svg"
                      xmlns:xlink="http://www.w3.org/1999/xlink"
                      y="0px"
                    >
                      <g id="检索优化1" transform="translate(-987, -124)">
                        <g id="编组-23" transform="translate(366, 116)">
                          <g id="被TA看过" transform="translate(614, 1)">
                            <g id="编组-11" transform="translate(8, 4)">
                              <g id="hebing-2" transform="translate(0, 4)">
                                <g id="形状">
                                  <path
                                    class="st0"
                                    d="M3.2,6.3c-0.1,0-0.2,0-0.2-0.1l0,0C2.9,6,2.9,5.8,3,5.7l0.5-0.5H1.1V4.5h2.3L3,4C2.9,3.9,2.9,3.7,3,3.6
							l0,0c0.1-0.1,0.1-0.1,0.2-0.1s0.2,0,0.2,0.1l1.3,1.3L3.5,6.2C3.4,6.2,3.3,6.3,3.2,6.3 M6.8,6.3c-0.1,0-0.2,0-0.2-0.1L5.3,4.9
							l1.2-1.3c0.1-0.1,0.1-0.1,0.2-0.1c0.1,0,0.2,0,0.2,0.1l0,0C7.1,3.7,7.1,3.9,7,4L6.5,4.5h2.3v0.7H6.5L7,5.7
							C7.1,5.8,7.1,6,7,6.2l0,0v0C6.9,6.2,6.9,6.3,6.8,6.3 M0.8,10C0.3,10,0,9.6,0,9.1V0.9C0,0.4,0.4,0,0.8,0h2.8
							c0.5,0,0.8,0.4,0.8,0.9v0.9H3.8V0.9c0-0.1,0-0.1,0-0.1c0,0-0.1-0.1-0.1-0.1H0.8c0,0-0.1,0-0.1,0.1c0,0-0.1,0.1,0,0.1v8.2
							c0,0.1,0.1,0.2,0.2,0.2h2.8c0,0,0.1,0,0.1-0.1c0,0,0.1-0.1,0-0.1V8.2h0.6v0.9c0,0.5-0.4,0.9-0.8,0.9L0.8,10L0.8,10L0.8,10z
							 M6.4,10c-0.5,0-0.8-0.4-0.8-0.9V8.3h0.6v0.9c0,0.1,0.1,0.2,0.2,0.2h2.8c0.1,0,0.2-0.1,0.2-0.2V0.9c0-0.1-0.1-0.2-0.2-0.2H6.4
							c-0.1,0-0.2,0.1-0.2,0.2v0.9H5.6V0.9C5.6,0.4,5.9,0,6.4,0h2.8C9.6,0,10,0.4,10,0.9v8.2C10,9.6,9.6,10,9.2,10L6.4,10L6.4,10z"
                                  />
                                </g>
                              </g>
                            </g>
                          </g>
                        </g>
                      </g>
                    </svg>
                  </el-icon>
                  归并
                </el-button>
                <el-button
                  class="items-btn"
                  :disabled="shortcutList.length === 0"
                  plain
                  size="small"
                  type="danger"
                  @click="handleEmpty"
                >
                  <el-icon><Delete /></el-icon>
                  清空
                </el-button>
              </div>
            </div>
          </vab-query-form-left-panel>
        </vab-query-form>
        <div class="flex-col" style="height: 0">
          <!-- 图表数据 -->
          <div class="chart">
            <div class="panel">
              <div class="num">
                {{ numberFormatte.format(queryPage.total) }}
              </div>
              <div style="color: #a9acb3">查询结果总数</div>
            </div>
            <vab-chart
              ref="application_ref"
              class="target-echart2"
              :option="retrieveChartOption"
              theme="vab-echarts-theme"
            />
          </div>
          <!-- 表格 -->
          <el-table
            ref="retrieveTableRef"
            v-loading="listLoading"
            class="my-table"
            :data="queryPage.listDate"
            @cell-contextmenu="useTableCopyEvent"
            @selection-change="setSelectRows"
          >
            <el-table-column v-if="showTableSelection" fixed="left" type="selection" width="55" />
            <el-table-column
              :align="'center'"
              fixed="left"
              :index="(index) => curIndex + index"
              label="序号"
              type="index"
              width="65"
            />
            <template v-for="item in tableColumn" :key="item.id">
              <el-table-column
                :label="item.fieldNameCn"
                :min-width="changeCellStyle(item.fieldNameCn)"
                :prop="item.fieldNameEn"
                resizable
                :show-overflow-tooltip="!['requestPayload', 'responsePayload'].includes(item.fieldNameEn)"
              >
                <template #header>
                  <span style="cursor: pointer">
                    {{ item.fieldNameCn }}
                    <el-image
                      v-if="item.supportAgg"
                      class="table-filter"
                      :src="require('@/assets/tongji-3.svg')"
                      @click="showTraceabilityField(item)"
                    />
                  </span>
                </template>
                <template #default="{ row }" v-if="['requestPayload', 'responsePayload'].includes(item.fieldNameEn)">
                  <span class="tabEllipsis" @click="() => showFieldHandle(item.fieldNameCn, row[item.fieldNameEn])">
                    {{ row[item.fieldNameEn] }}
                  </span>
                </template>
                <template #default="{ row }" v-else-if="item.fieldNameCn == '置信度'">
                  <el-rate
                    :colors="['#67C23A', '#67C23A', '#67C23A']"
                    disabled
                    disabled-void-color="#C7C6D4"
                    :max="3"
                    :model-value="getConfidence(row.confidence)"
                  />
                </template>
                <template #default="{ row }" v-else-if="item.fieldNameCn === '威胁等级'">
                  <span v-if="row.threatLevel" :class="['threatLevel_tag', getThreatLevel(row.threatLevel)]">
                    <el-icon><WarnTriangleFilled /></el-icon>
                    {{ row.threatLevel }}
                  </span>
                </template>
                <template
                  v-else-if="
                    (item.fieldNameCn.includes('时间') || item.fieldNameCn.includes('时长')) &&
                    item.fieldNameCn !== '时间服务器' &&
                    item.fieldNameCn !== '种子创建时间' &&
                    item.fieldNameCn !== 'IP租用时间' &&
                    item.fieldNameCn !== '响应时长' &&
                    !item.fieldNameCn.includes('时间戳')
                  "
                  #default="{ row }"
                >
                  {{ formatDate(row, item.fieldNameEn) }}
                </template>
                <template v-else-if="item.fieldNameCn.includes('等级')" #default="{ row }">
                  <span v-if="row.ruleLevel" :class="['alert_tag', getLevel(row.ruleLevel)]">
                    {{ row.ruleLevel }}
                  </span>
                </template>
              </el-table-column>
            </template>
            <el-table-column :align="'center'" fixed="right" label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" @click="handleInfo(row)">详情</el-button>
              </template>
            </el-table-column>
            <template #empty>
              <el-empty class="vab-data-empty" description="暂无数据" />
            </template>
          </el-table>
        </div>
      </el-tab-pane>
      <el-tab-pane label="数据收藏夹" lazy name="favorites">
        <favorites v-if="activeName === 'favorites'" :space-id="spaceId" />
      </el-tab-pane>
      <el-tab-pane label="溯源图" lazy name="traceability">
        <div class="traceability">
          <retrieve-traceability v-if="activeName === 'traceability'" :space-id="spaceId" />
        </div>
      </el-tab-pane>
      <el-tab-pane label="路径追踪" lazy name="pathTracing">
        <PathTracing />
      </el-tab-pane>
      <!-- <el-tab-pane label="测试" lazy name="pathTrac">
        <IframePage />
      </el-tab-pane> -->
      <!-- <el-tab-pane label="场景模型" lazy name="trafficModel">
        <TrafficModel />
      </el-tab-pane> -->
    </el-tabs>
    <!-- 表头字段配置 -->
    <application-configuration
      v-model="showFiledConfig"
      :fields="tableColumn"
      :retrieve-index-type="queryForm.indexType"
      @handleok="configurationHandel"
    />
    <!-- 发布 -->
    <publish
      v-if="showPublish"
      :show-publish="showPublish"
      @on-closeEvent="publishEvent"
      @on-submit="trafficModelHandle"
    />
    <!-- //告警索引 -->
    <component
      :is="detailDom[detailType]"
      v-if="detailType === 10"
      v-model:alert-detail-visible="detailVisible"
      v-loading="reloading"
      :info-val="detailVal"
      :module-enable="moduleEnable"
      :next="showNextBtn"
      :prev="showPrevBtn"
      :select-alert="detailVal"
      :workspace-id="spaceId"
      @on-skip-event="changeCurrentItemEvent"
    />

    <!-- 检索详情 -->
    <vab-dialog v-if="detailType !== 10" v-model="detailVisible" :destroy-on-close="true" title="详情" width="1125px">
      <div v-if="detailVisible" v-dialogBackTop>
        <component
          :is="detailDom[detailType] || DetailOther"
          ref="detailDomRef"
          v-loading="reloading"
          :info-val="detailVal"
          :workspace-id="spaceId"
        />
      </div>
      <Download class="download" :data="detailVal" />
      <ToolBox class="toolBox" />
      <EscUp v-if="showPrevBtn" class="ecs-up" @on-click-up="changeCurrentItemEvent(false)" />
      <EscDown v-if="showNextBtn" class="ecs-down" @on-click-down="changeCurrentItemEvent(true)" />
    </vab-dialog>
    <session-info
      v-if="sessionInfoVisible"
      :info-data="detailVal"
      :next="showNextBtn"
      :prev="showPrevBtn"
      :show-session-info="sessionInfoVisible"
      :workspace-id="spaceId"
      @on-close-event="sessionInfoVisible = false"
      @on-skip-event="changeCurrentItemEvent"
    />
    <!-- :dns-dict="DNSDict" -->
    <!-- 字段统计 -->
    <traceability-field-statistic
      ref="traceabilityFieldRef"
      show-context-menu
      @change-shortcut="handleUpdateCallback"
    />
    <!-- sql字段编辑添加 -->
    <add-sql
      v-if="show"
      :current-row="currentRow"
      :mode="mode"
      :show="show"
      :table-column="tableColumn"
      @handleUpdateCallback="handleUpdateCallback"
      @on-closeEvent="closeSql"
    />
    <!-- 检索历史 -->
    <retrieve-history
      v-if="showHistory"
      v-model="showHistory"
      :active="remark"
      :space-id="spaceId"
      @reloadSearch="reloadSearchHandle"
    />
    <!-- 收藏夹 工作空间 -->
    <retrieve-work-space ref="work_space_ref" />
    <!-- 规则列表页 -->
    <LevelRule v-if="showLevelRule" :is-show="showLevelRule" @on-closeEvent="showLevelRule = false" />
    <vab-dialog v-model="fieldDetailVisible" destroy-on-close :title="fieldDetailTitle" width="1325px">
      <div class="mask">
        <json-preview :json-value="fieldDetailValue" />
      </div>
    </vab-dialog>
  </div>
</template>

<style scoped lang="scss">
  $criticalColor: #ff0202;
  $lowColor: #1b81fe;
  $midColor: #f1b04d;
  $highColor: #fa6d15;
  $defaultColor: #909399;
  .ecs-up {
    position: fixed;
    left: calc((100vw - 1300px) / 2 - 50px);
    top: 50%;
    transform: translateY(-50%);
  }
  .ecs-down {
    position: fixed;
    right: calc((100vw - 1300px) / 2 - 40px);
    top: 50%;
    transform: translateY(-50%);
  }
  .download {
    position: fixed;
    left: 15px;
    bottom: 105px;
  }
  .toolBox {
    position: fixed;
    left: 15px;
    bottom: 65px;
  }
  .alert_tag {
    display: inline-block;
    width: 55px;
    border-radius: 5px;
    text-align: center;
    color: #fff;
    height: 30px;
    line-height: 30px;
    text-align: center;
    &.critical {
      background-color: rgba($criticalColor, 0.8);
    }
    &.low {
      background-color: rgba($lowColor, 0.8);
    }
    &.mid {
      background-color: rgba($midColor, 0.8);
    }
    &.high {
      background-color: rgba($highColor, 0.8);
    }
    &.default {
      background-color: #fff;
      color: inherit;
    }
  }
  .threatLevel_tag {
    display: inline-block;
    width: 55px;
    padding: 2px 3px;
    border-radius: 5px;
    color: #fff;
    font-size: 14px;
    .el-icon {
      margin-right: -4px;
      vertical-align: -3px;
      font-size: 17px;
    }
    &.critical {
      background-color: #ca0a08;
    }
    &.low {
      background-color: #ffbe36;
    }
    &.mid {
      background-color: #ff7212;
    }
    &.high {
      background-color: #ff2927;
    }
    &.default {
      background-color: #9d9aba;
      color: inherit;
    }
  }
  .tabEllipsis {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    width: clamp(120px, 90% - 10px, 90%);
    color: var(--el-color-primary);
    cursor: pointer;
    &:hover {
      text-decoration: underline;
    }
  }
  .mask {
    max-height: 600px;
    overflow-y: auto;
    border: 1px solid rgb(220, 223, 230);
    background-color: #f8f7ff;
  }
  .my-table {
    position: relative;
    flex: 1;
    :deep() {
      .el-table__row {
        height: 55px;
      }
      .el-table__body-wrapper {
        max-height: calc(100vh - 400px);
        min-height: calc(100vh - 400px);
        overflow-y: auto;
        &::-webkit-scrollbar {
          width: 0 !important;
          height: 0;
        }
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
    .cm-scroller .cm-content {
      padding-right: 40px;
    }
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
  }

  .retrieve-container {
    height: calc(100vh - 25px);
    overflow: hidden;
    margin-bottom: 30px;
    position: relative;
    .button {
      color: #0d88fe;
      margin-left: 30px;

      &:hover {
        cursor: pointer;
      }
    }

    .model {
      margin-left: 10px;
      position: relative;

      &::after {
        position: absolute;
        content: '';
        right: -19px;
        top: 2;
        height: 12px;
        border: 1px solid #0d88fe;
        opacity: 0.2;
      }
    }

    .btn {
      position: absolute;
      z-index: 99;
      top: 30px;
      right: 20px;
      font-size: 15px;
      color: #a9acb3;

      &:hover {
        cursor: pointer;
      }
    }

    .my-form {
      width: 100%;
      display: flex;
      .el-form-item {
        width: 190px;
        .el-button.el-button--primary {
          margin-left: 10px !important;
        }
      }
      .my-input {
        flex: 1;
        position: relative;
        .retrieve-sql {
          :deep(.cm-editor) {
            padding-right: 40px;
          }
        }
        .favoritesBtn {
          position: absolute;
          right: 3px;
          top: 2px;
          cursor: pointer;
          width: 40px;
          height: 36px;
          font-size: 18px;
          background-color: #fff;
          &::before {
            box-shadow: inset -10px 0 10px -10px rgba(0, 0, 0, 0.25);
            content: '';
            position: absolute;
            top: 0;
            width: 10px;
            left: -10px;
            bottom: -1px;
            overflow-x: hidden;
            overflow-y: hidden;
            touch-action: none;
            pointer-events: none;
          }
          i {
            margin-left: 12px;
            margin-top: 3px;
          }
          &.hasHistory i {
            color: #e6a23c;
          }
        }
      }
    }
  }

  :deep() {
    .favorites-action {
      position: absolute;
      right: 20px;
      top: 10px;
    }
    .vab-query-form[data-v-23d83642] .right-panel {
      margin-bottom: 0;
    }
    .el-checkbox {
      height: 25px;
    }
    .el-tooltip__trigger {
      cursor: pointer !important;
    }
  }
  .shortcut {
    width: 100%;
    display: flex;
    // align-items: center;
    .word {
      width: 70px;
      margin-top: 8px;
      line-height: 28px;
    }
    .SQLMode {
      :deep() {
        width: 73px;
        align-items: inherit;
        margin-top: 8px;
        height: 28px;
        .el-radio-button__inner {
          padding: 5px 9px;
          height: 28px;
          display: flex;
          align-items: center;
          font-size: 13px;
        }
      }
    }
    .containment {
      width: 70px;
      height: 28px;
      background: #ffffff;
      border-radius: 2.5px;
      border: 1px solid var(--el-border-color);
      margin: 8px 10px 0 0px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 12px;
      font-weight: 400;
      font-size: 13px;
      color: #1e1842;
    }
    .items-btn {
      height: 28px;
      margin-top: 8px;
      .el-icon {
        margin-right: 2px;
      }
    }
    .items {
      display: flex;
      width: calc(100% - 223px);
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
        display: flex;
        align-items: center;
        margin-right: 10px;
        margin-top: 8px;
        padding: 0 10px;
        border-radius: 2.5px;
        height: 28px;
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

  .chart {
    width: 100%;
    height: 117px;
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    justify-content: center;

    .panel {
      width: 120px;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      background-color: #f7fbff;
      flex-direction: column;

      .num {
        font-size: 25px;
        font-weight: 500;
        color: #303133;
      }
    }

    .target-echart2 {
      flex: 1;
      height: 100%;
    }
  }
  .traceability {
    background-color: #0d88fe;
    height: calc(100vh - 110px);
  }
  .table_action {
    margin-bottom: 20px;
    display: flex;
    justify-content: flex-end;

    .reflesh {
      margin-right: 30px;
    }

    .action {
      margin-left: 30px;
      display: flex;
      align-items: center;
    }

    .el-button {
      margin-left: 0;
    }

    .el-dropdown {
      float: right;
    }
  }

  .flex-col {
    display: flex;
    flex-direction: column;
    flex: 1;
  }

  ::v-deep .el-tabs {
    flex: 1;
    display: flex;
    flex-direction: column;
    // height: 100%;
  }
  ::v-deep .el-tabs__content {
    flex: 1;
  }
</style>
