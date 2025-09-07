<script lang="ts">
  export default {
    name: 'SiteSession',
  }
</script>
<script setup lang="ts">
  import SessionInfo from './session-info.vue'

  import FieldStatistic from './field-statistic.vue'

  import AddSql from './components/add-sql.vue'

  import SiteHistory from './components/history.vue'

  // import { EChartsOption } from 'echarts'

  import dayjs from 'dayjs'

  import VabChart from '@/plugins/VabChart/index.vue'
  import SearchSql from '@/components/search-sql/index.vue'
  import { getAllTagApi, siteSearchApi, getBySiteIdApi } from '~/src/api-ecs/site'

  import { siteSearchType, filedDataType, ShortcutListType, SiteHistoryItem } from '@/types'

  import numberFormatte from '@/utils/number'

  import { useUserStore } from '@/store/modules/user'

  import { useTableCopy, useCopy, uuid } from '@/utils'

  import type { FormInstance } from 'element-plus'

  import { useScroll } from '@vueuse/core'

  import { getTableCopyData } from '~/src/utils/transition'
  import { UpdateHistoryApi } from '~/src/api-ecs/site'
  const route = useRoute()

  const $baseMessage: any = inject('$baseMessage')

  const data = ref()

  const el = ref<HTMLElement | null>(null)

  const formRef = ref<FormInstance>()
  const sqlComponentsRef = ref<InstanceType<typeof SearchSql>>()
  const { getTableColumn } = useUserStore()
  const props = defineProps<{
    siteSessionId: number | string
    apiId: number | string
  }>()
  const arrowDown = ref(true) // 展开

  const retrieveChart = ref()
  const dataZoomToggle = ref(false)
  const curtableColumn = ref()
  // 表单数据
  const queryForm = reactive<siteSearchType>({
    startTime: '',
    siteSessionId: '',
    siteApiId: '',
    endTime: '',
    pageNum: 1,
    pageSize: 100,
    cookie: '',
    searchSql: '',
    orderField: 'requestTimeNs',
    orderType: 'desc',
    indexType: '1',
    xff: '',
    title: '',
    responseStatusCode: '',
    requestMethod: '',
    url: '',
    responsePayload: '',
    host: '',
    serverIp: '',
    userAgent: '',
    clientIp: '',
    requestPayload: '',
    scrollId: '',
    inputSql: '',
    filterSqlArr: '',
    sqlRelat: '',
  })

  // 图表配置
  const retrieveChartOption = reactive({
    xAxis: {
      type: 'category',
      data: null,
    },
    yAxis: {
      type: 'value',
      splitNumber: 5,
      axisLabel: {
        formatter: function (value: number) {
          return numberFormatte.format(+value)
        },
      },
    },
    series: [
      {
        data: null,
        type: 'line',
        color: '#6954f0',
        lineStyle: {
          width: 1,
        },
        areaStyle: {
          color: new VabChart.graphic.LinearGradient(0, 0, 0, 1, [
            {
              offset: 0,
              color: '#6954f0',
            },
            {
              offset: 1,
              color: '#fff',
            },
          ]),
        },
        smooth: 0.6,
        showSymbol: false,
      },
    ],
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
    grid: {
      top: 10,
      left: '2%',
      right: 10,
      bottom: 0,
      height: 66,
      containLabel: true,
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow',
      },
    },
  })

  const listLoading = ref(false) // 是否加载

  const favoritesDisable = ref(true)

  const showHistory = ref(false)

  const show = ref(false) // 是否显示添加sql页面

  const mode = ref('add') // 是否显示添加sql页面

  const currentRow = ref<ShortcutListType>()

  const rules = reactive({}) // 表单校验规则

  const queryPage = reactive({
    total: 0,
    listDate: [] as object[],
  })

  const shortcutList = ref<ShortcutListType[]>([])

  // 获取表格序号
  const curIndex = computed(() => (queryForm.pageNum - 1) * queryForm.pageSize + 1)

  const tableColumn = ref<filedDataType[]>([]) // 表头

  const showSessionInfo = ref(false) // 是否显示会话信息页面

  const showFieldStatistic = ref(false) // 显示字段统计页面

  const timeDate = ref() // 自定义时间

  const allTags = ref()

  const shortcutSQL = ref()

  const timeQuantum = ref('1hours') // 时间段最近一小时

  const timeDuration = ref('30s') // 时间长度30s

  const switchValue = ref(false) // 是否自动刷新

  const infoData = ref() // 详情信息

  const timeQuantumOptions = [
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

  const timeDuratioOptions = [
    {
      value: '30s',
      label: '30秒',
    },
    {
      value: '1min',
      label: '1分钟',
    },
    {
      value: '5min',
      label: '5分钟',
    },
    {
      value: '10min',
      label: '10分钟',
    },
  ]

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

  const cleanData = () => {
    const createFormData = () => ({
      startTime: '',
      siteSessionId: '',
      siteApiId: '',
      endTime: '',
      pageNum: 1,
      pageSize: 100,
      cookie: '',
      searchSql: '',
      orderField: 'requestTimeNs',
      orderType: 'desc',
      indexType: '1',
      xff: '',
      title: '',
      responseStatusCode: '',
      requestMethod: '',
      url: '',
      responsePayload: '',
      host: '',
      serverIp: '',
      userAgent: '',
      clientIp: '',
      requestPayload: '',
      scrollId: '',
    })
    Object.assign(queryForm, createFormData())
    dataZoomToggle.value = true
    SQLMode.value = '且'
    exclusive.value = false
  }

  // 添加修改sql
  const handleUpdate = ({ remark, data }: { remark: string; data?: ShortcutListType }) => {
    mode.value = remark
    show.value = true
    currentRow.value = data
  }

  const getTableColumnItem = async () => {
    const {
      data: { displayFieldsArr },
    } = await getBySiteIdApi({ id: props.siteSessionId, apiId: props.apiId })
    if (displayFieldsArr) {
      const allField = getTableColumn(props?.apiId ? 31 : 30)
      emit('setDisplayFieldsArr', displayFieldsArr)
      tableColumn.value = displayFieldsArr
        .map((i: number) => {
          return allField.find((item) => item.id === i)
        })
        .filter(Boolean)
    }
  }

  // 自动刷新
  let timer: any
  const autoReflash = (time: string) => {
    if (switchValue.value) {
      const res = changeTime(time)
      timer = setInterval(async () => {
        if (!listLoading.value) {
          clearInterval(timer)
          formatDayDate()
          await siteSessionSearch()
          autoReflash(time)
        }
      }, res)
    } else {
      clearInterval(timer)
    }
  }

  onUnmounted(() => {
    clearInterval(timer)
  })

  const changeTime = (time: string) => {
    let timer = 0
    switch (time) {
      case '30s':
        // todo
        timer = 30 * 1000
        break
      case '1min':
        timer = 60 * 1000
        break
      case '5min':
        timer = 5 * 60 * 1000
        break
      case '10min':
        timer = 10 * 60 * 1000
        break
    }
    return timer
  }

  // 删除sql
  const handleClear = (id: string) => {
    const _index = shortcutList.value.findIndex((item: any) => {
      return item.id == id
    })
    shortcutList.value.splice(_index, 1)
  }

  // 打开
  const unfold = () => {
    arrowDown.value = !arrowDown.value
    formRef.value?.resetFields()
  }

  const resetData = () => {
    queryForm.scrollId = ''
    queryPage.listDate = []
  }

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
  const handleCheckboxChange = (id: string) => {
    if (exclusive.value) {
      shortcutList.value.forEach((_: any, index: any) => {
        if (shortcutList.value[index].id != id) {
          shortcutList.value[index].enable = false
        }
      })
      siteSessionSearch()
    }
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
  // 站点检索
  const siteSessionSearch = async (remark = true) => {
    if (sqlComponentsRef.value && !sqlComponentsRef.value?.getValidateSql()) return
    if (remark) {
      resetData()
    }
    if (!queryForm.siteSessionId) return
    shortcutSQL.value = ''
    const arr = shortcutList.value.filter((item: any) => {
      return item.enable
    })
    if (arrowDown.value && arr.length > 0) {
      let shortcut = '('
      arr.forEach((td: any, index: number) => {
        if (index == arr.length - 1) {
          shortcut = `${shortcut}${arr[index].label} ${arr[index].relation} ${
            // @ts-ignore
            arr[index].value || arr[index].value == 0 || arr[index].value == '0' ? `"${arr[index].value}"` : ''
          })`
        } else {
          shortcut =
            `${shortcut}${arr[index].label} ${arr[index].relation} ${
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
      if (queryForm.searchSql?.trim() && arrowDown.value) {
        shortcutSQL.value = JSON.parse(JSON.stringify(probeSearchSqlIsIPOrPort()))
      }
    }
    listLoading.value = true
    formatDayDate()
    try {
      let filterSqlArr = ''
      if (shortcutList.value.length > 0) {
        filterSqlArr = JSON.stringify(shortcutList.value)
      }
      let sqlRelat = ''
      sqlRelat = JSON.stringify({ SQLMode: SQLMode.value, exclusive: exclusive.value })
      const {
        data: { resList, sumaryMap, total, scrollId },
      } = await siteSearchApi({
        ...queryForm,
        searchSql: shortcutSQL.value,
        inputSql: queryForm.searchSql,
        filterSqlArr,
        sqlRelat,
      })
      await getTableColumnItem()
      queryForm.scrollId = scrollId
      if (!queryForm.scrollId) {
        queryPage.listDate = []
      }
      if (resList.length > 0) {
        const id = resList[0].id
        const flag = queryPage.listDate.some((item: any) => {
          return item.id == id
        })
        if (flag) return false
        queryPage.listDate = [...queryPage.listDate, ...resList]
      }
      queryPage.total = total
      if (sumaryMap) {
        const { count, date } = JSON.parse(sumaryMap)
        retrieveChartOption.xAxis.data = date || []
        retrieveChartOption.series[0].data = count || []
      }

      if (dataZoomToggle.value) {
        retrieveChart.value?.chart?.clear()
        retrieveChart.value?.chart?.setOption(retrieveChartOption)
      }
      dataZoomToggle.value = false
      if (remark) {
        sqlComponentsRef.value?.changeHistories(shortcutSQL.value)
      }
    } finally {
      favoritesDisable.value = false
      listLoading.value = false
    }
  }

  const reloadSearchHandle = (data: SiteHistoryItem, hasTime: boolean) => {
    let { searchSql } = data
    const { searchEdTime, searchStTime, indexType, inputSql, filterSqlArr, siteId, apiId, sqlRelat } = data
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
    queryForm.pageSize = 10
    queryForm.startTime = hasTime ? searchStTime : ''
    queryForm.endTime = hasTime ? searchEdTime : ''
    queryForm.searchSql = searchSql
    queryForm.siteSessionId = siteId
    queryForm.siteApiId = apiId
    timeQuantum.value = hasTime ? 'user-defined' : timeQuantum.value
    timeDate.value = hasTime ? [searchStTime, searchEdTime] : null
    dataZoomToggle.value = true
    if (queryForm.indexType === indexType) {
      siteSessionSearch()
    } else {
      queryForm.indexType = indexType
    }
  }

  // 查询
  // const queryData = () => { }

  const favoritesHandle = async () => {
    const { indexType, endTime, startTime } = queryForm
    let filterSqlArr = ''
    if (shortcutList.value.length > 0) {
      filterSqlArr = JSON.stringify(shortcutList.value)
    }
    const { code } = await UpdateHistoryApi({
      collectStatus: true,
      searchStTime: startTime,
      searchSql: shortcutSQL.value,
      searchEdTime: endTime,
      indexType,
      inputSql: queryForm.searchSql,
      filterSqlArr,
      siteId: queryForm.siteSessionId,
      apiId: queryForm.siteApiId,
      // workspaceId: spaceId.value,
    })
    $baseMessage('收藏成功', 'success', 'vab-hey-message-success')
  }

  const showPrevBtn = ref(false)

  const showNextBtn = ref(false)

  // 上一条或下一条
  const changeCurrentItemEvent = (val: boolean) => {
    let index = queryPage.listDate.findIndex((item: any) => {
      return item.id == infoData.value.id
    })
    infoData.value = undefined
    if (val) {
      if (index >= queryPage.listDate.length - 5 && queryPage.listDate.length! != queryPage.total) {
        siteSessionSearch(false)
      }
      ++index
      infoData.value = { ...queryPage.listDate[index] }
    } else {
      --index
      infoData.value = { ...queryPage.listDate[index] }
    }
    index !== 0 ? (showPrevBtn.value = true) : (showPrevBtn.value = false)
    index !== queryPage.total - 1 ? (showNextBtn.value = true) : (showNextBtn.value = false)
  }

  // 详情
  const handleInfo = (val: boolean, row?: any) => {
    infoData.value = row
    showSessionInfo.value = true
    const index = queryPage.listDate.findIndex((item: any) => {
      return item.id == infoData.value.id
    })
    index !== 0 ? (showPrevBtn.value = true) : (showPrevBtn.value = false)
    index !== queryPage.total - 1 ? (showNextBtn.value = true) : (showNextBtn.value = false)
  }

  // 会话字段统计
  const fieldList: string[] = []
  let field_title = ''
  const handleFieldStatistic = async (val: boolean, row?: string, title?: string, column?: any) => {
    showFieldStatistic.value = val
    curtableColumn.value = column
    if (row) {
      fieldList[0] = row
    }
    field_title = title || ''
  }

  // 获取所有标签
  const getAllTag = async () => {
    const { data } = await getAllTagApi()
    allTags.value = data
  }

  // 页容量改变
  const handleSizeChange = (val: number) => {
    queryForm.pageSize = val
    siteSessionSearch()
  }

  // 页面改变
  const handleCurrentChange = (val: number) => {
    queryForm.pageNum = val
    siteSessionSearch()
  }

  function changeCellStyle(title: string) {
    if (['客户端口', '服务端口', '请求方式', '状态码'].includes(title)) {
      return '125'
    } else if (['客户端IP', '服务端IP', 'HOST', 'XFF'].includes(title)) {
      return '120'
    } else if (['请求时间', '响应时长'].includes(title)) {
      return '160'
    } else {
      return '150'
    }
  }
  const searchhandle = (val: string) => {
    siteSessionSearch()
  }
  const formatDayDate = () => {
    if (timeQuantum.value === 'user-defined') {
      formatUserTime()
      return
    }
    const timeDate = dayjs()
    const endDate = timeDate.format('YYYY-MM-DD HH:mm:ss')
    const cur = timeDate.startOf('day')
    let startDate = null
    switch (timeQuantum.value) {
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
      timeQuantum.value === 'yesterday'
        ? timeDate.endOf('day').subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss')
        : endDate
    queryForm.startTime = startDate
  }

  // 得到时间
  watchEffect(() => {
    formatDayDate()
  })

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

  // 得到自定义时间
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

  // id改变更新页面
  watch(
    () => [props.siteSessionId, props.apiId],
    () => {
      if (props.siteSessionId || props.apiId) {
        handleEmpty()
        cleanData()
        queryForm.siteSessionId = props.siteSessionId
        queryForm.siteApiId = props.apiId
        siteSessionSearch()
      }
    },
    {
      immediate: true,
    }
  )

  watch(
    () => switchValue.value,
    () => {
      autoReflash(timeDuration.value)
    }
  )
  type obj = {
    siteId: string | number
    apiId?: string | number
  }
  const emit = defineEmits<{
    (e: 'changeActive', data: obj): void
    (e: 'setDisplayFieldsArr', data: Array<any>): void
  }>()
  watch(
    () => queryForm.siteSessionId,
    () => {
      emit('changeActive', {
        siteId: queryForm.siteSessionId,
        apiId: queryForm.siteApiId,
      })
    }
  )
  watch(
    () => queryForm.siteApiId,
    () => {
      emit('changeActive', {
        siteId: queryForm.siteSessionId,
        apiId: queryForm.siteApiId,
      })
    }
  )
  watch(
    () => retrieveChart.value,
    () => {
      retrieveChart.value.chart.on('datazoom', () => {
        const { dataZoom, xAxis } = retrieveChart.value.chart.getModel().option
        const startValue = dataZoom[0].startValue
        const endValue = dataZoom[0].endValue
        const startLable = xAxis[0].data[startValue]
        const endLable = xAxis[0].data[endValue]
        timeQuantum.value = 'user-defined'
        timeDate.value = [startLable, endLable]
        queryForm.pageNum = 1
        queryForm.pageSize = 100
        dataZoomToggle.value = false
        siteSessionSearch()
      })
    }
  )
  onMounted(() => {
    // getAllTag()
    el.value = document.querySelector('.my-table .el-table__body-wrapper')
    const { arrivedState } = useScroll(el.value)
    data.value = arrivedState
    if (route.query.timeDuration) {
      timeQuantum.value = route.query.timeDuration as string
    }
  })

  function formatDate(row: any, key: string) {
    const time = row[key] / 1000000
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss.SSS')
  }

  const useTableCopyEvent = (row: any, column: any, cell: any, event: any) => {
    const arr = getTableCopyData({ tableColumn: tableColumn.value, row, column, mothod: handleUpdateCallback })
    useTableCopy(row, column, cell, event, arr)
  }

  watch(
    () => [shortcutList.value, queryForm.searchSql],
    () => {
      favoritesDisable.value = true
    },
    { deep: true }
  )

  watch(
    () => data.value,
    () => {
      if (data.value.bottom) {
        if (queryPage.listDate.length === queryPage.total) return
        siteSessionSearch(false)
      }
    },
    { deep: true }
  )
  defineExpose({
    siteSessionSearch,
    showHistory,
  })
  // class="word">过滤条件：</div>
</script>

<template>
  <div class="site-session">
    <!-- 检索栏 -->
    <vab-query-form style="width: 100%">
      <vab-query-form-left-panel :span="16">
        <el-select v-model="timeQuantum" class="m-2">
          <el-option v-for="item in timeQuantumOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <vab-date-time-picker v-if="timeQuantum === 'user-defined'" v-model="timeDate" style="width: 300px" />
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="8">
        <el-select v-model="timeDuration" class="m-2" :disabled="switchValue">
          <el-option v-for="item in timeDuratioOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-switch v-model="switchValue" active-text="自动刷新" />
      </vab-query-form-right-panel>
      <vab-query-form-right-panel :span="24">
        <el-form class="my-form" inline :model="queryForm" @submit.prevent>
          <el-form-item class="my-input">
            <search-sql
              v-if="arrowDown"
              ref="sqlComponentsRef"
              :model-value="queryForm.searchSql"
              :site-data="{
                siteId: +queryForm.siteSessionId,
                apiId: queryForm.siteApiId ? +queryForm.siteApiId : undefined,
              }"
              @on-change="(str) => (queryForm.searchSql = str)"
              @onSearch="searchhandle"
            />
            <el-input v-else :column-list="tableColumn" :disabled="!arrowDown" />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              @click="
                () => {
                  dataZoomToggle = true
                  siteSessionSearch()
                }
              "
            >
              检索
            </el-button>
            <!-- 收藏 -->
            <el-button :disabled="favoritesDisable || !shortcutSQL.length" type="primary" @click="favoritesHandle">
              <vab-icon icon="star-line" style="font-size: 18px" />
            </el-button>
            <!-- <el-button @click="unfold">
              <el-tooltip class="item" :content="arrowDown ? '展开' : '折叠'" effect="dark" placement="top">
                <vab-icon v-if="arrowDown" icon="arrow-down-s-line" style="font-size: 20px; color: #ccc" />
                <vab-icon v-else icon="arrow-up-s-line" style="font-size: 20px; color: #ccc" />
              </el-tooltip>
            </el-button> -->
          </el-form-item>
        </el-form>
      </vab-query-form-right-panel>
      <vab-query-form-left-panel v-if="arrowDown" :span="24">
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
    <!-- 图表或表单 -->
    <div v-if="arrowDown" class="chart">
      <div class="panel">
        <div class="num">
          {{ numberFormatte.format(queryPage.total) }}
        </div>
        <div style="color: #a9acb3">查询结果总数</div>
      </div>
      <vab-chart ref="retrieveChart" class="target-echart2" :option="retrieveChartOption" theme="vab-echarts-theme" />
    </div>

    <el-form
      v-else
      ref="formRef"
      class="login-form"
      label-position="top"
      :model="queryForm"
      :rules="rules"
      style="width: 100%; display: flex"
    >
      <el-row :gutter="10">
        <el-col :span="8">
          <el-row :gutter="10">
            <el-col :span="12">
              <el-form-item label="源IP" prop="clientIp">
                <el-input v-model="queryForm.clientIp" clearable />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="目的IP" prop="serverIp">
                <el-input v-model="queryForm.serverIp" clearable />
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="URL" prop="url">
            <el-input v-model="queryForm.url" clearable />
          </el-form-item>
          <el-form-item label="标题" prop="title">
            <el-input v-model="queryForm.title" clearable />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-row :gutter="10">
            <el-col :span="12">
              <el-form-item label="XFF" prop="xff">
                <el-input v-model="queryForm.xff" clearable />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="HOST" prop="host">
                <el-input v-model="queryForm.host" clearable />
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="请求payload关键字" prop="requestPayload">
            <el-input v-model="queryForm.requestPayload" clearable />
          </el-form-item>
          <el-form-item label="User agent" prop="userAgent">
            <el-input v-model="queryForm.userAgent" clearable />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-row :gutter="10">
            <el-col :span="12">
              <el-form-item label="请求方式" prop="requestMethod">
                <el-input v-model="queryForm.requestMethod" clearable />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="状态码" prop="responseStatusCode">
                <el-input v-model="queryForm.responseStatusCode" clearable />
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="响应payload关键字" prop="responsePayload">
            <el-input v-model="queryForm.responsePayload" clearable />
          </el-form-item>
          <el-form-item label="cookie" prop="cookie">
            <el-input v-model="queryForm.cookie" clearable />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
    <!-- 表格 -->
    <el-table
      v-loading="listLoading"
      class="my-table"
      :data="queryPage.listDate"
      style="margin-top: 20px; width: 100%"
      @cell-contextmenu="useTableCopyEvent"
    >
      <el-table-column v-if="!listLoading" label="序号" width="65">
        <template #default="{ $index }">
          <span>{{ curIndex + $index }}</span>
        </template>
      </el-table-column>
      <template v-for="item in tableColumn" :key="item?.id">
        <el-table-column
          v-if="item.fieldNameCn.includes('时间')"
          min-width="160"
          :prop="item.fieldNameEn"
          :resizable="true"
          show-overflow-tooltip
        >
          <template #header>
            {{ item.fieldNameCn }}

            <el-image
              v-if="item.supportAgg"
              class="table-filter"
              :src="require('@/assets/tongji-3.svg')"
              style="margin-left: 2px"
              @click="handleFieldStatistic(true, item.fieldNameEn, item.fieldNameCn, item)"
            />
          </template>
          <template #default="{ row }">
            {{ formatDate(row, item.fieldNameEn) }}
          </template>
        </el-table-column>
        <el-table-column
          v-else
          :min-width="changeCellStyle(item.fieldNameCn)"
          :prop="item.fieldNameEn"
          :resizable="true"
          show-overflow-tooltip
        >
          <template #header>
            {{ item.fieldNameCn }}

            <el-image
              v-if="item.supportAgg"
              class="table-filter"
              :src="require('@/assets/tongji-3.svg')"
              style="margin-left: 2px"
              @click="handleFieldStatistic(true, item.fieldNameEn, item.fieldNameCn, item)"
            />
          </template>
        </el-table-column>
        <!-- <el-table-column v-else align="center" :label="item.fieldNameCn" min-width="100" :prop="item.fieldNameEn"
                                                                                                                                                                                    :resizable="false" show-overflow-tooltip /> -->
      </template>
      <el-table-column v-if="!listLoading" fixed="right" label="操作" width="100">
        <template #default="{ row }">
          <el-button size="small" @click="handleInfo(true, row)">详情</el-button>
        </template>
      </el-table-column>
      <template #empty>
        <el-empty class="vab-data-empty" description="暂无数据" />
      </template>
    </el-table>
  </div>
  <!-- 会话详情 -->
  <session-info
    v-if="showSessionInfo"
    :info-data="infoData"
    :next="showNextBtn"
    :prev="showPrevBtn"
    :show-session-info="showSessionInfo"
    @on-closeEvent="showSessionInfo = false"
    @onSkipEvent="changeCurrentItemEvent"
  />
  <!-- 会话字段统计 -->
  <field-statistic
    v-if="showFieldStatistic"
    :aggregations-payload="{ ...queryForm, searchSql: shortcutSQL }"
    :field-list="fieldList"
    show-context-menu
    :show-field-statistic="showFieldStatistic"
    :table-column="curtableColumn"
    :title="field_title"
    @change-shortcut="handleUpdateCallback"
    @on-closeEvent="handleFieldStatistic"
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
  <site-history
    v-if="showHistory"
    v-model="showHistory"
    :site-id="queryForm.siteSessionId"
    @reloadSearch="reloadSearchHandle"
  />
</template>

<style scoped lang="scss">
  .my-table {
    position: relative;
    :deep() {
      .el-table__inner-wrapper {
        width: 100%;
      }
      .el-table__body-wrapper {
        position: initial !important;
        width: 100%;
      }
      .el-scrollbar {
        position: initial !important;
      }
      .el-scrollbar__bar.is-horizontal {
        position: absolute;
        display: block !important;
      }
    }
    :deep() {
      .el-table__body-wrapper {
        max-height: calc(100vh - 400px) !important;
        min-height: calc(100vh - 400px) !important;
        overflow-y: auto;
        &::-webkit-scrollbar {
          width: 0 !important;
          height: 0;
        }
      }
    }
  }
  :deep() {
    .super_height {
      width: 100%;
      height: 192px;

      .el-tooltip__trigger,
      .el-input {
        height: 100%;
      }
    }
  }

  .site-session {
    margin-top: 15px;
    width: calc(100%);
    padding-bottom: 20px;
  }

  .m-2 {
    margin-right: 20px;
  }

  .form-tags {
    margin-left: 10px;
    width: 360px;
  }

  .tags {
    padding: 10px;
    height: 192px;
    overflow: auto;

    &::-webkit-scrollbar {
      width: 0px;
      height: 0px;
    }

    border: 1px solid var(--el-border-color);
  }

  :deep(.el-tag) {
    margin-left: 10px;
  }

  .field-statistic {
    &:hover {
      cursor: pointer;
    }
  }

  :deep() {
    .vab-query-form[data-v-23d83642] .right-panel {
      margin-bottom: 0;
    }
    .el-checkbox {
      height: 25px;
    }
  }
  .my-form {
    width: 100%;
    display: flex;

    .my-input {
      flex: 1;
      width: calc(100% - 175px);
      // overflow-x: auto;
      :deep() {
        .el-form-item__content {
          width: 100%;
        }
      }
    }
    .my-action {
      width: 210px;
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
    height: 6.5vw;
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
</style>
