<script setup lang="ts">
  import ApplicationConfiguration from '@/ecs/assets/components/application-configuration.vue'
  import AddSql from '../site/site-session/components/add-sql.vue'
  import VabChart from '@/plugins/VabChart/index.vue'
  import Traceability from './components/traceability.vue'
  import AttackPerspective from './components/attack-perspective.vue'
  import AlertDetail from './components/alert-detail.vue'
  import SearchSql from '@/components/search-sql/index.vue'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import JsonPreview from '@/components/json-preview.vue'
  import dayjs, { Dayjs } from 'dayjs'
  import { TableColumnItemType } from '/#/store'
  import { level_options, levelKey, timeDuratioOptions } from './data/index'
  import { timeDuratioOptions as date_options } from '@/data/constant'
  import { AlertItem, ShortcutListType } from '@/types'
  import { getAlertApi, updateAlertStatusApi, getThreatTypeListApi } from '@/api-ecs/alert'
  import { formatNstime, formatTime } from '@/utils/time'
  import numberFormatte from '@/utils/number'
  import { useTableCopy, useCopy, uuid } from '@/utils'
  import { useUserStore } from '@/store/modules/user'
  import TraceabilityFieldStatistic from '@/ecs/alert/components/field-statistic.vue'
  import { updateDisplayApi } from '@/api-ecs/custom-field'
  import { useScroll } from '@vueuse/core'
  import { getAllDisPlaysFiledApi } from '~/src/api-ecs/public'
  import { getTableCopyData } from '~/src/utils/transition'
  import { getSystemConfigApi } from '~/src/api-ecs/system'
  import numberFormatter from '@/utils/number'
  import { getTopFieldApi } from '~/src/api-ecs/retrieve'
  const userStore = useUserStore()
  const branch_chart = ref()
  const route = useRoute()
  const fieldDetailVisible = ref(false)
  const fieldDetailTitle = ref('')
  const fieldDetailValue = ref('')
  // const threatTypeList = ref<{ dictLabel: string; dictValue: string }[]>([])
  const data = ref()
  const moduleEnable = ref(false)
  const el = ref<HTMLElement | null>(null)
  const traceabilityFieldRef = ref<InstanceType<typeof TraceabilityFieldStatistic>>()
  const { getTableColumn } = userStore
  const $baseMessage: any = inject('$baseMessage')
  // 查询时间
  const timeDuration = ref('1hours')
  // 自定义时间
  const timeDate = ref<[Dayjs, Dayjs] | undefined>(undefined)
  // 自动刷新
  const autoRefresh = ref(false)
  const autoRefreshDuration = ref(30)
  // 定时器ID
  const timerId = ref<NodeJS.Timer>()
  // 快速筛选折叠
  const fold = ref(true)
  // 高级筛选
  const isAdvanced = ref(false)
  const loading = ref(false)
  const sqlComponentsRef = ref<InstanceType<typeof SearchSql>>()
  // 图表箭头切换
  const chartDetail = ref(false)
  // 溯源显示
  const traceabilityVisible = ref(false)
  // 告警详情显示
  const alertDetailVisible = ref(false)
  const multipleSelection = ref<AlertItem[]>([])
  const dataZoomToggle = ref(false)
  // 柱状图实例
  const alertChart_ref = ref()

  const must_parames = reactive({
    /** 排序字段 */
    orderField: 'startTimeNs',
    /** 排序方式 */
    orderType: 'decs',
    /** 索引类型 */
    indexType: 10,
    /** 当前页码 */
    pageNum: 1,
    /** 当前页数 */
    pageSize: 100,
    /** 开始时间 */
    startTime: '',
    /** 结束时间 */
    endTime: '',
  })
  const infoVal = ref()
  const isAttackPerspective = ref(false)
  const alertInfo = reactive({
    list: [] as AlertItem[],
    total: 0,
    pieData: [] as { name: string; value: string | number }[],
    lineData: {
      date: [] as any[],
      count: [] as any[],
    },
    threat_level: [] as any[],
    scrollId: '',
  })

  provide('isFullscreen', isAttackPerspective)
  const attackResultType = ['成功', '失败', '未知', '企图']
  const tableColumn = ref<TableColumnItemType[]>([])

  const userDisPlaysFiled = ref()

  const sqlKey: { [key: string]: string } = {
    attackIp: '源IP',
    sourcePort: '源端口',
    xff: 'XFF',
    url: 'URL',
    host: 'Host',
    victimIp: '目的IP',
    targetPort: '目的端口',
    threatType: '威胁类型',
    threatName: '威胁名称',
    threatLevel: '威胁等级',
    attackResult: '攻击结果',
    readStatus: '已读状态',
  }

  // 告警检索条件
  const queryForm = reactive({
    /** 查询SQL */
    searchSql: '',
    /** 源IP */
    attackIp: '',
    whiteType: 0,
    /** 源端口 */
    sourcePort: null,
    xff: '',
    url: '',
    host: '',
    /** 目的IP */
    victimIp: '',
    /** 目的端口 */
    targetPort: null,
    /** 威胁类型 */
    threatType: [] as string[],
    /** 威胁名称 */
    threatName: '',
    /** 威胁等级 */
    threatLevel: '',
    /** 攻击结果 */
    attackResult: [] as string[],
    /** 攻击次数 */
    attackCount: '',
    /** 已读状态 */
    readStatus: '',
  })

  const modeValue = ref('race')
  // 折线段图配置
  const chartOption = reactive({
    xAxis: {
      type: 'category',
      data: [] as string[],
    },
    yAxis: {
      type: 'value',
      splitNumber: 3,
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
        backgroundColor: '#FEFCFC',
        fillerColor: 'rgba(255,237,236,.5)',
        moveHandleStyle: { color: '#EAD7D7' },
        handleStyle: {
          borderColor: '#EAD7D7',
          color: '#EAD7D7',
        },
        emphasis: {
          moveHandleStyle: { color: '#EAD7D7' },
          handleStyle: {
            borderColor: '#EAD7D7',
          },
        },
        selectedDataBackground: {
          lineStyle: {
            color: '#FF5757',
          },
          areaStyle: {
            color: '#EE0000',
          },
        },
        dataBackground: {
          lineStyle: {
            color: '#FF5757',
          },
          areaStyle: {
            color: '#FFE3E3',
          },
        },
      },
    ],
    series: {
      type: 'line',
      color: '#FF5757',
      areaStyle: {
        color: new VabChart.graphic.LinearGradient(0, 0, 0, 1, [
          {
            offset: 0,
            color: 'rgba(255, 95, 95, .5)',
          },
          {
            offset: 1,
            color: '#fff',
          },
        ]),
      },
      lineStyle: {
        width: 1,
      },
      smooth: 0.6,
      showSymbol: false,
      data: [] as number[],
    },
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
  // 源ip图配置
  const sourceIPEchartOption = reactive({
    title: {
      // show: props.showBtn,
      text: '攻击源IP',
      left: 0,
      textStyle: {
        color: '#303133 ',
        fontSize: 14,
      },
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: {
          backgroundColor: '#6a7985',
        },
      },
      formatter: function (params: any) {
        return `${params[0].name}  :  ${params[0].value}`
      },
      confine: true,
    },
    grid: {
      top: 30,
      left: '0',
      right: 40,
      bottom: 10,
      containLabel: true,
    },
    xAxis: {
      type: 'value',
      boundaryGap: false,
      axisLine: { show: false, lineStyle: { color: '#ccc' } },
      axisTick: { show: false },
      axisLabel: {
        color: '#999',
        formatter: function (value: number) {
          return numberFormatter.format(+value)
        },
      },
      // splitLine: { lineStyle: { color: ['#CEEDFF'], type: [5, 8], dashOffset: 3 } },
    },
    yAxis: {
      type: 'category',
      data: [],
      axisLine: { show: true, lineStyle: { color: '#ccc' } },
      axisTick: { length: 3 },
      splitLine: {
        show: false,
      },
      // axisLabel: { show: true, fontSize: 12, color: '#666', margin: 12, padding: 0 },
      inverse: true,
    },
    series: [
      {
        name: '',
        type: 'bar',
        showBackground: true,
        backgroundStyle: { color: 'rgba(245, 244, 255, 1)', borderRadius: [0, 8, 8, 0] },
        itemStyle: {
          color: '#52A8FF',
          normal: {
            borderRadius: [0, 8, 8, 0],
            color: '#FFC548',
          },
        },
        barMaxWidth: 10,
        data: [],
      },
    ],
  })
  const targetIPEchartOption = reactive({
    title: {
      text: '目的IP',
      textStyle: {
        color: '#303133 ',
        fontSize: 14,
      },
      left: 0,
    },
    tooltip: {
      trigger: 'axis',
      extraCssText: 'z-index:1',
    },
    grid: {
      top: 30,
      left: '0',
      right: 40,
      bottom: 10,
      containLabel: true,
    },

    xAxis: {
      type: 'value',
      boundaryGap: false,
      axisLine: { show: false, lineStyle: { color: '#ccc' } },
      axisTick: { show: false },
      axisLabel: {
        color: '#999',
        formatter: function (value: number) {
          return numberFormatter.format(+value)
        },
      },
      splitLine: { lineStyle: { color: ['#CEEDFF'], type: [5, 8], dashOffset: 3 } },
    },
    yAxis: {
      type: 'category',
      data: [],
      axisLine: { show: true, lineStyle: { color: '#ccc' } },
      axisTick: { length: 3 },
      splitLine: {
        show: false,
      },
      axisLabel: { show: true, fontSize: 12, color: '#666', margin: 12, padding: 0 },
      inverse: true,
    },
    series: [
      {
        name: '',
        type: 'bar',
        showBackground: true,
        backgroundStyle: { color: 'rgba(245, 244, 255, 1)', borderRadius: [0, 8, 8, 0] },
        itemStyle: {
          color: '#858BFF',
          normal: {
            borderRadius: [0, 8, 8, 0],
            color: '#858BFF',
          },
        },
        barMaxWidth: 10,
        data: [],
      },
    ],
  })
  // 面积图配置
  const trendEchartOption = reactive({
    title: {
      text: '攻击趋势图',
      textStyle: {
        color: '#303133 ',
        fontSize: 14,
      },
    },
    tooltip: {
      trigger: 'axis',
      extraCssText: 'z-index:1',
    },
    grid: {
      top: 40,
      left: '0',
      right: 10,
      bottom: 10,
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: [] as any[],
      boundaryGap: false,
    },
    yAxis: [
      {
        type: 'value',
      },
    ],
    series: {
      type: 'line',
      data: [] as any[],
      symbol: 'circle',
      color: '#FF5757',
      areaStyle: {
        color: new VabChart.graphic.LinearGradient(0, 0, 0, 0.7, [
          {
            offset: 0,
            color: 'rgba(255,163, 168, 1)',
          },
          // {
          //   offset: 0.8,
          //   color: 'rgba(255,163, 168, 0.7)',
          // },
          {
            offset: 1,
            color: 'rgba(255,163, 168, 0.05)',
          },
        ]),
      },
      smooth: true,
      yAxisIndex: 0,
      showSymbol: false,
      // areaStyle: {
      //   opacity: 0.2,
      // },
    },
  })

  interface sqlObjType {
    [key: string]: any
  }

  const sqlObj: sqlObjType = {
    attackIp: '源IP',
    sourcePort: '源端口',
    victimIp: '目的IP',
    targetPort: '目的端口',
    host: 'Host',
    threatType: '威胁类型',
    threatName: '威胁名称',
    threatLevel: '威胁等级',
    attackResult: '攻击结果',
    readStatus: '读取状态',
    xff: 'XFF',
    url: 'URL',
  }

  // 饼图配置
  const branchEchartOption = reactive({
    title: {
      text: '告警统计图',
      left: 0,
    },
    tooltip: {
      trigger: 'item',
    },
    legend: {
      bottom: '0',
      left: 'center',
      // selectedMode: false,
      itemStyle: {
        opacity: 1,
      },
    },
    grig: {
      top: 0,
      bottom: 10,
    },
    series: {
      type: 'pie',
      radius: ['35%', '60%'],
      center: ['50%', '35%'],
      selectedMode: 'single',
      minAngle: 10,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2,
      },
      label: {
        show: false,
        position: 'center',
      },
      labelLine: {
        show: false,
      },
      select: {
        itemStyle: {
          shadowColor: 'rgba(0, 0, 0, 0.4)',
          shadowBlur: 10,
          borderWidth: 0.5,
          opacity: 1,
        },
      },
      data: [] as any[],
    },
  })

  const showFiledConfig = ref(false) // 显示表头字段配置

  const shortcutList = ref<ShortcutListType[]>([])

  const shortcutSQL = ref()

  const resetqueryForm = () => {
    for (const key in queryForm) {
      if (['sourcePort', 'targetPort'].includes(key)) {
        // @ts-ignore
        queryForm[key] = null
      } else {
        // @ts-ignore
        queryForm[key] = ''
      }
      // else if (key === 'whiteType')
      //  {}
    }

    for (const shortcut of shortcutList.value) {
      shortcut.enable = false
    }
  }
  // 清空
  const handleEmpty = () => {
    shortcutList.value = []
  }

  const show = ref(false) // 是否显示添加sql页面

  const mode = ref('add') // 是否显示添加sql页面
  const currentRow = ref<ShortcutListType>()

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

  function formatDataTime(row: any, key: string) {
    const time = row[key] / 1000000
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
  }

  // 表头字段配置
  const tableHeadConfig = (val: boolean) => {
    showFiledConfig.value = val
  }

  const resetData = () => {
    alertInfo.scrollId = ''
    alertInfo.list = []
  }

  // const handleGetThreatTypeList = async () => {
  //   const { data } = await getThreatTypeListApi()
  //   threatTypeList.value = data || []
  // }

  const formatSql = () => {
    const arr1 = Object.keys(sqlObj)
    const obj = {} as { [key: string]: any }
    for (const key in queryForm) {
      if (arr1.includes(key)) {
        // @ts-ignore
        if (queryForm[key]) {
          // @ts-ignore
          obj[sqlObj[key]] = queryForm[key]
        }
      }
    }
    if (!obj.威胁类型?.length) delete obj.威胁类型
    if (!obj.攻击结果?.length) delete obj.攻击结果
    const arr = Object.keys(obj)
    if (arr.length > 0) {
      let shortcut = '('
      arr.forEach((td: any, index: number) => {
        if (index == arr.length - 1) {
          // @ts-ignore
          shortcut = `${shortcut}${arr[index]} = ` + `"${obj[arr[index]]}"` + `)`
        } else {
          // @ts-ignore
          shortcut = `${shortcut}${arr[index]} = ` + `"${obj[arr[index]]}"` + ` ` + `and` + ` `
        }
        if (['威胁类型', '攻击结果'].includes(td)) {
          if (obj[td].length === 0) return
          const str = `${arr[index]} = "${obj[td]}"`
          shortcut = shortcut.replace(str, `(${obj[td].map((i: any) => `${td} = "${i}"`).join(' or ')})`)
        }
      })
      shortcutSQL.value = shortcut
    }
  }
  // 或和且互斥功能
  const SQLMode = ref('且')
  const exclusive = ref(false)
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
    const inputSql = queryForm.searchSql ? `( ${queryForm.searchSql} ) and ` : ''
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
      queryData()
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

  // 查询告警日志
  async function queryData(remark = true, updateChart = true) {
    if (sqlComponentsRef.value && !sqlComponentsRef.value?.getValidateSql()) return
    formatDate()
    if (remark) {
      resetData()
    }
    if (!timeDate.value && timeDuration.value === 'user-defined')
      return $baseMessage('请选择时间！', 'error', 'vab-hey-message-error')
    shortcutSQL.value = ''
    const arr = shortcutList.value.filter((item: any) => {
      return item.enable
    })
    formatSql()
    if (arr.length > 0) {
      let shortcut = '('
      arr.forEach((td: any, index: number) => {
        if (index == arr.length - 1) {
          shortcut = `${shortcut} ${arr[index].label} ${arr[index].relation} ${
            arr[index].value || arr[index].value == '0' ? `"${arr[index].value}"` : ''
          })`
        } else {
          shortcut =
            `${shortcut} ${arr[index].label} ${arr[index].relation} ${
              arr[index].value || arr[index].value == '0' ? `"${arr[index].value}"` : ''
            } ` +
            `${SQLMode.value == '且' ? 'and' : 'or'}` +
            ` `
        }
      })
      if (queryForm.searchSql?.trim()) {
        shortcutSQL.value = `(${probeSearchSqlIsIPOrPort()})` + ` ` + `and` + ` ${shortcut}`
      } else {
        if (shortcutSQL.value?.trim()) {
          shortcutSQL.value = `${shortcutSQL.value}` + ` ` + `and` + ` ${shortcut}`
        } else {
          shortcutSQL.value = shortcut
        }
      }
    } else {
      if (queryForm.searchSql?.trim()) {
        if (!shortcutSQL.value) {
          shortcutSQL.value = JSON.parse(JSON.stringify(probeSearchSqlIsIPOrPort()))
        } else {
          shortcutSQL.value =
            `${shortcutSQL.value}` + ` ` + `and` + `(${JSON.parse(JSON.stringify(probeSearchSqlIsIPOrPort()))})`
        }
      }
    }
    loading.value = true
    try {
      const { searchSql } = queryForm
      const queryData = isAdvanced.value
        ? {
            ...must_parames,
            whiteType: queryForm.whiteType ? 1 : 0,
            searchSql: shortcutSQL.value,
            scrollId: alertInfo.scrollId,
          }
        : {
            ...queryForm,
            attackResult: '',
            threatType: '',
            whiteType: queryForm.whiteType ? 1 : 0,
            searchSql: shortcutSQL.value,
            ...must_parames,
            scrollId: alertInfo.scrollId,
          }
      const {
        data: { total, sumaryMap, resList, scrollId },
      } = await getAlertApi(queryData)
      alertInfo.scrollId = scrollId
      alertInfo.list = [...alertInfo.list, ...resList]
      alertInfo.total = total
      if (updateChart) {
        sqlComponentsRef.value?.changeHistories(shortcutSQL.value)
        if (!sumaryMap) return
        const { pillar, pie, threat_level } = JSON.parse(sumaryMap)
        alertInfo.pieData = pie || []
        alertInfo.lineData = {
          date: pillar.date,
          count: pillar.count,
        }
        alertInfo.threat_level = threat_level || []
        if (dataZoomToggle.value) {
          alertChart_ref.value?.chart?.clear()
          alertChart_ref.value?.chart?.setOption(chartOption)
        }
      }
      if (chartDetail.value) {
        initSourceEcharts()
        initTagetEcharts()
      }
      formatColum()
      loading.value = false
      dataZoomToggle.value = false
    } catch (error) {
      loading.value = false
      console.error(error)
    }
  }

  // 获取用户所有字段
  const getAllDisPlaysFiled = async () => {
    const { data } = await getAllDisPlaysFiledApi()
    userDisPlaysFiled.value = data
    formatColum()
  }
  const showTraceabilityField = (column: TableColumnItemType) => {
    const { indexType, startTime, endTime } = must_parames
    const { fieldNameCn, fieldNameEn } = column
    const whiteType = queryForm.whiteType ? 1 : 0
    // const SQL = isAdvanced.value
    //   ? queryForm.searchSql.trim()
    //   : Object.entries(queryForm)
    //       .map((obj) => {
    //         const [key, val] = obj
    //         if (val) return `${sqlKey[key]} ${['威胁名称', '威胁类型'].includes(sqlKey[key]) ? 'like' : '='} "${val}"`
    //         return false
    //       })
    //       .filter(Boolean)
    //       .join(' and ')

    const SQL = shortcutSQL.value
    traceabilityFieldRef.value?.initData({
      title: fieldNameCn,
      query: { indexType, startTime, endTime, aggregationFields: fieldNameEn, searchSql: SQL, whiteType },
      column,
    })
  }
  const obj = { 低危: 1, 中危: 2, 高危: 3, 危急: 4 }
  function getLevel(str: string) {
    // @ts-ignore
    const level = obj[str]
    return levelKey[level]
  }
  const handleDownloadExcel = () => {
    const tHeader: string[] = []
    const filterVal: any = []
    try {
      tableColumn.value.forEach(({ fieldNameCn, fieldNameEn }) => {
        tHeader.push(fieldNameCn)
        filterVal.push(fieldNameEn)
      })
      loading.value = true
      import('@/utils/excel').then((excel) => {
        const list = multipleSelection.value.length ? multipleSelection.value : alertInfo.list
        const data = formatJson(filterVal, list)
        excel.export_json_to_excel({
          header: tHeader,
          data,
          filename: `告警列表-${formatTime(new Date().getTime())}`,
          autoWidth: true,
          bookType: 'xlsx',
        })
        loading.value = false
      })
    } catch (error) {
      console.log(error)
      loading.value = false
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
      case 'warnTime':
      case 'startTimeNs':
        return formatNstime(val[key])
      default:
        return val[key] || ''
    }
  }

  async function configurationHandel(showField: TableColumnItemType[]) {
    tableColumn.value = showField
    const { msg } = await updateDisplayApi({
      indexType: must_parames.indexType,
      displayIds: showField.map((i) => i.id),
    })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    getAllDisPlaysFiled()
  }

  function formatColum() {
    const columnData = getTableColumn(must_parames?.indexType)
    const userColumnData = userDisPlaysFiled.value && (userDisPlaysFiled.value[must_parames?.indexType] as number[])
    if (userColumnData?.length > 0) {
      tableColumn.value = userColumnData.map((key: any) =>
        columnData.find((item) => item.id === key)
      ) as TableColumnItemType[]
    } else {
      tableColumn.value = []
    }
  }

  function showTraceability(row: any) {
    infoVal.value = row
    traceabilityVisible.value = true
  }

  const showPrevBtn = ref(false)

  const showNextBtn = ref(false)

  // 上一条或下一条
  const changeCurrentItemEvent = (val: boolean) => {
    let index = alertInfo.list.findIndex((item: any) => {
      return item.id == infoVal.value.id
    })
    if (val) {
      if (index >= alertInfo.list.length - 5 && alertInfo.list.length! != alertInfo.total) {
        queryData(false, false)
      }
      ++index
      infoVal.value = { ...alertInfo.list[index] }
    } else {
      // o.0 蚌埠了
      if (index === 0) {
        index !== 0 ? (showPrevBtn.value = true) : (showPrevBtn.value = false)
        index !== alertInfo.total - 1 ? (showNextBtn.value = true) : (showNextBtn.value = false)
        return $baseMessage('已是最新告警')
      }
      --index
      infoVal.value = { ...alertInfo.list[index] }
    }
    index !== 0 ? (showPrevBtn.value = true) : (showPrevBtn.value = false)
    index !== alertInfo.total - 1 ? (showNextBtn.value = true) : (showNextBtn.value = false)
  }

  function showAlertDetail(row: any) {
    infoVal.value = row
    alertDetailVisible.value = true
    let index = alertInfo.list.findIndex((item: any) => {
      return item.id == infoVal.value.id
    })
    index !== 0 ? (showPrevBtn.value = true) : (showPrevBtn.value = false)
    index !== alertInfo.total - 1 ? (showNextBtn.value = true) : (showNextBtn.value = false)
  }
  // 多选表格
  function handleSelectionChange(val: AlertItem[]) {
    multipleSelection.value = val
  }

  const colorRef = ref<'red-row' | 'bule-row' | 'green-row' | 'grep-row'>()
  const colorObj = {
    'red-row': 1,
    'bule-row': 2,
    'green-row': 3,
    'grep-row': 4,
  }
  const handleChangeColor = (val: 'red' | 'bule' | 'green' | 'grep') => {
    colorRef.value = `${val}-row`
    updateRowColor(colorObj[colorRef.value])
  }
  /**
   * 修改颜色
   */
  const updateRowColor = async (color: number) => {
    const ids = multipleSelection.value.map((i) => i.id)
    const { msg } = await updateAlertStatusApi({ ...must_parames, color, ids })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    colorRef.value = undefined
    multipleSelection.value = []
    queryData(true, true)
  }

  const tableRowClassName = ({ row, rowIndex }: { row: any; rowIndex: number }) => {
    let color = ''
    for (const key in colorObj) {
      // @ts-ignore
      if (row.color == colorObj[key]) {
        color = key
      }
    }
    return color
  }
  /**
   * type: 标记已读1 标记未读0
   */
  const updateStatus = async (status: '1' | '0') => {
    const ids = multipleSelection.value.map((i) => i.id)
    const { msg } = await updateAlertStatusApi({ ...must_parames, readStatus: status, ids })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    multipleSelection.value = []
  }
  /**
   * type: 标记忽略1 取消忽略0
   */
  const updateIgnoreStatus = async (status: '1' | '0') => {
    const ids = multipleSelection.value.map((i) => i.id)
    const { msg } = await updateAlertStatusApi({ ...must_parames, ignoreStatus: status, ids })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    multipleSelection.value = []
  }

  // 获取表格序号
  const curIndex = computed(() => (must_parames.pageNum - 1) * must_parames.pageSize + 1)
  const formatDate = () => {
    if (timeDuration.value === 'user-defined') {
      formatUserTime()
      return
    }

    const timeDate = dayjs()
    const endDate = [
      '1hours',
      '3hours',
      '6hours',
      '12hours',
      '24hours',
      'last-week',
      'last-two-week',
      'last-three-days',
    ].includes(timeDuration.value)
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
      case '12hours':
        startDate = timeDate.subtract(12, 'hour').format('YYYY-MM-DD HH:mm:ss')
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
    must_parames.endTime =
      timeDuration.value === 'yesterday'
        ? timeDate.endOf('day').subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss')
        : endDate
    must_parames.startTime = startDate
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
  const searchhandle = (val: string) => {
    queryData(true, true)
  }
  const formatSourceData = (val: string) => {
    if (!val) return ''
    let _val = ''
    try {
      _val = JSON.parse(val)
    } catch (error) {
      _val = val
    }
    return _val
  }
  const showFieldHandle = (title: string, val: string) => {
    fieldDetailVisible.value = true
    fieldDetailTitle.value = title
    fieldDetailValue.value = JSON.stringify(formatSourceData(val), null, 4)
  }
  const getModelStatus = async () => {
    const {
      data: { module_config },
    } = await getSystemConfigApi({ keys: 'module_config' })
    moduleEnable.value = JSON.parse(module_config.value).isEnable || false
  }
  watch(
    () => modeValue.value,
    () => {
      if (modeValue.value == 'race') {
        isAdvanced.value = false
        queryForm.searchSql = ''
        queryForm.threatType = []
        isAttackPerspective.value = false
      } else {
        resetqueryForm()
        isAdvanced.value = true
        isAttackPerspective.value = false
        queryData(true, true)
      }
    }
  )

  const handleSwitch = () => {
    queryData(true, true)
  }
  watch(
    () => alertInfo.lineData,
    () => {
      const { date, count } = alertInfo.lineData
      chartOption.xAxis.data = date
      chartOption.series.data = count
      trendEchartOption.xAxis.data = date.slice(-10)
      trendEchartOption.series.data = count.slice(-10)
    }
  )
  watch(
    () => chartDetail.value,
    () => {
      if (!chartDetail.value) {
        queryForm.threatType = []
        nextTick(() => queryData(true, true))
      } else {
        initSourceEcharts()
        initTagetEcharts()
      }
    }
  )

  const initSourceEcharts = async () => {
    const { indexType, startTime, endTime } = must_parames
    const fieldNameEn = 'attackIp'
    const whiteType = queryForm.whiteType ? 1 : 0
    const SQL = shortcutSQL.value
    const query = { indexType, startTime, endTime, aggregationFields: fieldNameEn, searchSql: SQL, whiteType }
    const {
      data: { aggObj },
    } = await getTopFieldApi({
      ...query,
      topCount: 10,
    })
    sourceIPEchartOption.yAxis.data = []
    sourceIPEchartOption.series[0].data = []
    aggObj?.forEach((item: any) => {
      // @ts-ignore
      sourceIPEchartOption.yAxis.data.push(item.name)
      // @ts-ignore
      sourceIPEchartOption.series[0].data.push(item.value)
    })
  }
  const initTagetEcharts = async () => {
    const { indexType, startTime, endTime } = must_parames
    const fieldNameEn = 'victimIp'
    const whiteType = queryForm.whiteType ? 1 : 0
    const SQL = shortcutSQL.value
    const query = { indexType, startTime, endTime, aggregationFields: fieldNameEn, searchSql: SQL, whiteType }
    const {
      data: { aggObj },
    } = await getTopFieldApi({
      ...query,
      topCount: 10,
    })
    targetIPEchartOption.yAxis.data = []
    targetIPEchartOption.series[0].data = []
    aggObj?.forEach((item: any) => {
      // @ts-ignore
      targetIPEchartOption.yAxis.data.push(item.name)
      // @ts-ignore
      targetIPEchartOption.series[0].data.push(item.value)
    })
  }
  watch(
    () => alertInfo.pieData,
    () => {
      branchEchartOption.series.data = alertInfo.pieData
    }
  )
  watch(
    () => isAdvanced.value,
    () => {
      shortcutSQL.value = undefined
      must_parames.pageNum = 1
      must_parames.pageSize = 100
    }
  )
  // watch(
  //   () => branch_chart.value,
  //   () => {
  //     if (branch_chart.value) {
  //       branch_chart.value?.chart.on('selectchanged', function (params: any) {
  //         const {
  //           fromActionPayload: { dataIndexInside, seriesIndex },
  //           fromAction,
  //         } = params
  //         const _option = branch_chart.value?.chart.getOption()
  //         const { series } = _option
  //         if (fromAction === 'unselect') {
  //           series[seriesIndex].itemStyle.opacity = 1
  //           branch_chart.value?.chart.setOption(_option)
  //           queryForm.threatType = []
  //           queryData(true, false)
  //         } else {
  //           series[seriesIndex].itemStyle.opacity = 0.5
  //           branch_chart.value?.chart.setOption(_option)
  //           const type = branchEchartOption.series.data[dataIndexInside]
  //           queryForm.threatType = type?.name ? [type?.name] : []
  //           queryData(true, false)
  //         }
  //       })
  //       branch_chart.value?.chart.on('legendselectchanged', function (params: any) {
  //         const { name } = params
  //         branch_chart.value?.chart.dispatchAction({
  //           type: 'legendSelect',
  //           name,
  //         })
  //         const _option = branch_chart.value?.chart.getOption()
  //         const { series } = _option
  //         series[0].selectedMap = {
  //           [queryForm.threatType[0]]: false,
  //         }
  //         if (queryForm.threatType.includes(name)) {
  //           queryForm.threatType = []
  //           series[0].itemStyle.opacity = 1
  //           series[0].selectedMap[name] = false
  //         } else {
  //           queryForm.threatType = [name]
  //           series[0].itemStyle = {
  //             ...series[0].itemStyle,
  //             opacity: 0.5,
  //           }
  //           series[0].selectedMap = {
  //             ...series[0].selectedMap,
  //             [name]: true,
  //           }
  //         }
  //         branch_chart.value?.chart.setOption(_option)
  //         queryData(true, false)
  //       })
  //     }
  //   }
  // )

  watch(
    () => timeDuration.value,
    () => {
      dataZoomToggle.value = true
      formatDate()
    },
    {
      immediate: true,
    }
  )

  const formatUserTime = () => {
    if (timeDate.value) {
      const [startDate, endDate] = timeDate.value
      must_parames.endTime = dayjs(endDate).format('YYYY-MM-DD HH:mm:ss')
      must_parames.startTime = dayjs(startDate).format('YYYY-MM-DD HH:mm:ss')
    }
  }

  watchEffect(() => {
    if (timeDate.value) {
      formatUserTime()
    }
  })
  watchEffect(() => {
    if (autoRefresh.value) {
      // @ts-ignore
      if (timerId.value) clearInterval(timerId.value)
      timerId.value = setInterval(() => queryData(), autoRefreshDuration.value * 1000)
    } else {
      // @ts-ignore
      clearInterval(timerId.value)
    }
  })
  onMounted(async () => {
    if (route.query.attackIp) {
      queryForm.attackIp = route.query.attackIp as string
      timeDuration.value = '12hours'
    }
    if (route.query.threatType) {
      queryForm.threatType = [route.query.threatType as string]
      timeDuration.value = '12hours'
    }
    getModelStatus()
    // handleGetThreatTypeList()
    setTimeout(() => {
      queryData()
      getAllDisPlaysFiled()
      el.value = document.querySelector('.my-table .el-scrollbar__wrap')
      const { arrivedState } = useScroll(el.value)
      data.value = arrivedState
      isAttackPerspective.value = route.query.params == 'isAttackPerspective' ? true : false
    }, 0)
  })

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

  const useTableCopyEvent = (row: any, column: any, cell: any, event: any) => {
    const arr = getTableCopyData({ tableColumn: tableColumn.value, row, column, mothod: handleUpdateCallback })
    useTableCopy(row, column, cell, event, arr)
  }
  watch(
    () => alertChart_ref.value,
    () => {
      if (alertChart_ref.value) {
        alertChart_ref.value.chart.on('datazoom', () => {
          const { dataZoom, xAxis } = alertChart_ref.value.chart.getModel().option
          const startValue = dataZoom[0].startValue
          const endValue = dataZoom[0].endValue
          const startLable = xAxis[0].data[startValue]
          const endLable = xAxis[0].data[endValue]
          timeDuration.value = 'user-defined'
          timeDate.value = [startLable, endLable]
          must_parames.pageNum = 1
          must_parames.pageSize = 100
          dataZoomToggle.value = false
          queryData(true, false)
        })
      }
    }
  )

  onUnmounted(() => {
    // @ts-ignore
    if (timerId.value) clearInterval(timerId.value)
  })

  watch(
    () => data.value,
    () => {
      if (data.value.bottom) {
        nextTick(() => queryData(false, false))
      }
    },
    { deep: true }
  )
</script>

<script lang="ts">
  export default {
    name: 'TunnelTopics',
  }
</script>

<template>
  <div class="alert-container flex-col">
    <div class="tools">
      <div class="tools-left">
        <el-form v-if="!isAttackPerspective" inline label-position="left">
          <el-form-item style="width: 148px">
            <el-select v-model="timeDuration">
              <el-option v-for="item in date_options" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item v-if="timeDuration === 'user-defined'" style="width: 380px">
            <vab-date-time-picker v-model="timeDate" />
          </el-form-item>
        </el-form>
        <!-- 日志查询方式 -->
        <div class="table_action" style="margin-bottom: 18px">
          <div v-if="!isAttackPerspective" class="action reflesh">
            <label>自动刷新：</label>
            <el-switch v-model="autoRefresh" />
          </div>
          <el-select
            v-if="!isAttackPerspective"
            v-model="autoRefreshDuration"
            class="m-2"
            :disabled="autoRefresh"
            style="width: 100px"
          >
            <el-option v-for="item in timeDuratioOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </div>
      </div>
      <div class="tools-right" style="margin-bottom: 18px">
        <div class="action">
          <label>原始告警：</label>
          <el-switch
            v-model="queryForm.whiteType"
            :active-value="1"
            :inactive-value="0"
            style="margin-right: 20px"
            @change="handleSwitch"
          />
          <el-button-group>
            <el-button
              v-if="modeValue === 'race'"
              type="primary"
              @click="
                () => {
                  modeValue = 'more'
                }
              "
            >
              <el-icon style="font-size: 16px; margin-right: 5px"><Refresh /></el-icon>
              高级筛选
            </el-button>
            <el-button
              v-if="modeValue === 'more'"
              type="primary"
              @click="
                () => {
                  modeValue = 'race'
                }
              "
            >
              <el-icon style="font-size: 16px; margin-right: 5px"><Refresh /></el-icon>
              快速筛选
            </el-button>
            <el-button
              style="margin-left: 20px"
              :type="isAttackPerspective ? 'primary' : 'default'"
              @click="
                () => {
                  isAttackPerspective = true
                }
              "
            >
              攻击透视
            </el-button>
          </el-button-group>
        </div>
        <!-- 显示表头字段配置 -->
        <el-tooltip v-if="!isAttackPerspective" class="item" content="字段配置" effect="dark" placement="top">
          <div @click="tableHeadConfig(true)">
            <vab-icon
              icon="edit-2-line item-icon"
              style="font-size: 20px; color: #b3b9c8; margin-left: 20px; margin-right: 5px"
            />
            <span style="margin-right: 20px; cursor: pointer">编辑</span>
          </div>
        </el-tooltip>
        <div v-if="!isAttackPerspective" class="table_action">
          <el-dropdown trigger="click">
            <el-button>
              更多操作
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <!-- <el-dropdown-item :disabled="multipleSelection.length === 0" @click="download">
              下载Pcap文件
            </el-dropdown-item> -->
                <!-- <el-popover
                  ref="popoverRef"
                  :virtual-ref="buttonRef"
                  trigger="hover"
                  title="With title"
                  virtual-triggering
                ></el-popover> -->
                <el-dropdown-item @click="handleDownloadExcel">导出告警日志</el-dropdown-item>
                <el-dropdown-item :disabled="multipleSelection.length === 0" @click="updateStatus('1')">
                  标记已读
                </el-dropdown-item>
                <el-dropdown-item :disabled="multipleSelection.length === 0" @click="updateStatus('0')">
                  标记未读
                </el-dropdown-item>
                <el-dropdown-item :disabled="multipleSelection.length === 0">
                  <el-popover placement="left" trigger="hover" :width="100">
                    <template #reference>标记颜色</template>
                    <template #default>
                      <div style="height: 100%; width: 100%; padding: 0 20px">
                        <div class="my-circular" @click="handleChangeColor('red')">
                          <div class="circular" style="background: #fff2f2; border: 1px solid #ff7a7a"></div>
                          红色
                        </div>
                        <div class="my-circular" @click="handleChangeColor('bule')">
                          <div class="circular" style="background: #f1f7ff; border: 1px solid #67a6ff"></div>
                          蓝色
                        </div>
                        <div class="my-circular" @click="handleChangeColor('green')">
                          <div class="circular" style="background: #f0fff7; border: 1px solid #79e7a5"></div>
                          绿色
                        </div>
                        <div class="my-circular" @click="handleChangeColor('grep')">
                          <div class="circular" style="background: #f4f4f4; border: 1px solid #c2c2c2"></div>
                          灰色
                        </div>
                      </div>
                    </template>
                  </el-popover>
                </el-dropdown-item>
                <el-dropdown-item :disabled="multipleSelection.length === 0" @click="updateIgnoreStatus('1')">
                  标记忽略
                </el-dropdown-item>
                <el-dropdown-item :disabled="multipleSelection.length === 0" @click="updateIgnoreStatus('0')">
                  取消忽略
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </div>
    <!-- 查询表单 -->
    <hr style="opacity: 0.4; margin-top: -5px; margin-bottom: 15px; transform: scale(1, 0.5)" />
    <el-row>
      <el-col v-if="isAdvanced" :span="24">
        <div v-if="!isAttackPerspective" style="display: flex; align-items: center; padding-bottom: 15px">
          <SearchSql
            ref="sqlComponentsRef"
            :index-type="must_parames.indexType"
            :model-value="queryForm.searchSql"
            @on-change="(str) => (queryForm.searchSql = str)"
            @onSearch="searchhandle"
          />
          <el-button style="margin-left: 20px" type="primary" @click="() => queryData(true, true)">检索</el-button>
        </div>
      </el-col>
      <template v-else>
        <el-col v-if="!isAttackPerspective" :span="21">
          <el-row :gutter="20">
            <el-form inline label-position="top" label-width="110px" style="width: 100%">
              <el-col :span="4">
                <el-form-item label="源IP" prop="attackIp">
                  <el-input v-model="queryForm.attackIp" clearable />
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="源端口" prop="sourcePort">
                  <el-input v-model.number="queryForm.sourcePort" />
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="目的IP" prop="victimIp">
                  <el-input v-model="queryForm.victimIp" clearable />
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="目的端口" prop="targetPort">
                  <el-input v-model.number="queryForm.targetPort" clearable />
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="Host" prop="host">
                  <el-input v-model="queryForm.host" clearable />
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="威胁类型" prop="threatType">
                  <el-select v-model="queryForm.threatType" clearable collapse-tags multiple>
                    <!-- <el-option
                      v-for="item in threatTypeList"
                      :key="item.dictValue"
                      :label="item.dictLabel"
                      :value="item.dictLabel"
                    /> -->
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col v-if="!fold" :span="4">
                <el-form-item label="威胁名称" prop="threatName">
                  <el-input v-model="queryForm.threatName" clearable />
                </el-form-item>
              </el-col>
              <el-col v-if="!fold" :span="4">
                <el-form-item label="威胁等级" prop="threatLevel">
                  <el-select v-model="queryForm.threatLevel">
                    <el-option
                      v-for="item in level_options"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col v-if="!fold" :span="4">
                <el-form-item v-show="!fold" label="攻击结果" prop="attackResult">
                  <el-select v-model="queryForm.attackResult" clearable collapse-tags multiple :value-on-clear="null">
                    <el-option v-for="item in attackResultType" :key="item" :label="item" :value="item" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col v-if="!fold" :span="4">
                <el-form-item label="XFF" prop="xff">
                  <el-input v-model="queryForm.xff" clearable />
                </el-form-item>
              </el-col>
              <el-col v-if="!fold" :span="4">
                <el-form-item label="URL" prop="url">
                  <el-input v-model="queryForm.url" clearable />
                </el-form-item>
              </el-col>
              <el-col v-if="!fold" :span="4">
                <el-form-item label="已读状态" prop="readStatus">
                  <el-select v-model="queryForm.readStatus" clearable>
                    <el-option label="已读" value="1" />
                    <el-option label="未读" value="0" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-form>
          </el-row>
        </el-col>
        <el-col v-if="!isAttackPerspective" :span="3">
          <el-form inline :label-position="isAdvanced ? 'left' : 'top'" label-width="100px">
            <el-form-item class="submit_item" label="&nbsp;" style="margin-right: 0px">
              <el-button plain style="margin-left: 20px" type="danger" @click="resetqueryForm">
                <el-icon><Delete /></el-icon>
              </el-button>
              <div style="display: flex">
                <el-button type="primary" @click="() => queryData(true, true)">检索</el-button>
                <div v-if="!isAdvanced" class="card-header-tag">
                  <img
                    class="form_icon"
                    :class="{ upward: !fold }"
                    src="@/assets/alert_images/down.png"
                    @click="fold = !fold"
                  />
                </div>
              </div>
            </el-form-item>
          </el-form>
        </el-col>
      </template>
    </el-row>
    <div>
      <vab-query-form-left-panel v-if="!isAttackPerspective" :span="24">
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
            <el-button class="items-btn" plain size="small" type="primary" @click="handleUpdate({ remark: 'add' })">
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
    </div>
    <!-- 告警日志操作 -->
    <template v-if="!isAttackPerspective">
      <div>
        <!-- 告警日志图表 -->
        <vab-card class="access" shadow="hover" skeleton>
          <template #header>
            <p>
              共计 - {{ alertInfo.total }}条
              <span v-for="alert of alertInfo.threat_level" :key="alert.name" class="log">
                <span class="risk_status" :class="[getLevel(alert.name)?.class]">
                  {{ getLevel(alert.name)?.label }}
                </span>
                <span>&nbsp;-&nbsp;{{ alert.value }}条</span>
              </span>
            </p>
            <div class="card-header-tag">
              <img
                class="form_icon"
                :class="{ upward: chartDetail }"
                src="@/assets/alert_images/down.png"
                @click="chartDetail = !chartDetail"
              />
            </div>
          </template>

          <template v-if="chartDetail">
            <vab-chart :option="sourceIPEchartOption" theme="vab-echarts-theme" />
            <vab-chart :option="targetIPEchartOption" theme="vab-echarts-theme" />
            <!-- <vab-chart ref="branch_chart" class="trend-echart" :option="branchEchartOption" theme="vab-echarts-theme" /> -->
            <vab-chart :option="trendEchartOption" theme="vab-echarts-theme" />
          </template>
          <template v-else>
            <vab-chart ref="alertChart_ref" class="small" :option="chartOption" theme="vab-echarts-theme" />
          </template>
        </vab-card>
      </div>
      <div class="flex-col" style="height: 85vh">
        <!-- 告警日志列表 -->
        <el-table
          v-loading="loading"
          :border="true"
          class="my-table"
          :data="alertInfo.list"
          element-loading-text="Loading..."
          :row-class-name="tableRowClassName"
          @cell-contextmenu="useTableCopyEvent"
          @selection-change="handleSelectionChange"
        >
          <el-table-column align="center" type="selection" width="55" />
          <el-table-column align="center" :index="curIndex" label="序号" type="index" width="75" />
          <el-table-column
            v-for="item in tableColumn"
            :key="item.id"
            align="center"
            :label="item.fieldNameCn"
            :prop="item.fieldNameEn"
            :show-overflow-tooltip="!['sourceData'].includes(item.fieldNameEn)"
            :width="
              [
                'threatName',
                'startTimeNs',
                'server_country_code',
                'client_country_code',
                'server_country',
                'abnormal_city',
                'abnormal_country',
                'abnormal_incident',
                'abnormal_province',
                'abnormal_time',
              ].includes(item.fieldNameEn)
                ? 250
                : 120
            "
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
            <template #default="{ row }" v-if="item.fieldNameEn == 'startTimeNs'">
              {{ formatDataTime(row, item.fieldNameEn) }}
            </template>
            <template #default="{ row }" v-else-if="['sourceData'].includes(item.fieldNameEn)">
              <div class="tabEllipsis" @click="() => showFieldHandle(item.fieldNameCn, row[item.fieldNameEn])">
                {{ row[item.fieldNameEn] }}
              </div>
            </template>
            <template #default="{ row }" v-else-if="item.fieldNameEn == 'threatLevel'">
              <span :class="['alert_tag', getThreatLevel(row.threatLevel)]">{{ row.threatLevel }}</span>
            </template>
            <template #default="{ row }" v-else-if="item.fieldNameEn == 'attackResult'">
              <span
                :class="[
                  'attackResult',
                  `attackResult-${attackResultType.findIndex((type) => type === row?.attackResult)}`,
                ]"
              >
                <el-icon><WarnTriangleFilled /></el-icon>
                {{ row.attackResult }}
              </span>
            </template>
          </el-table-column>

          <el-table-column align="center" fixed="right" label="操作" width="150">
            <template #default="{ row }">
              <el-button class="row_action" size="small" @click="showAlertDetail(row)">详情</el-button>
              <el-button class="row_action" size="small" @click="showTraceability(row)">溯源</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </template>
    <AttackPerspective v-else />

    <vab-dialog v-model="fieldDetailVisible" destroy-on-close :title="fieldDetailTitle" width="1375px">
      <div class="mask">
        <json-preview :json-value="fieldDetailValue" />
      </div>
    </vab-dialog>

    <!-- 表头字段配置 -->
    <application-configuration
      v-if="showFiledConfig"
      v-model="showFiledConfig"
      :fields="tableColumn"
      :retrieve-index-type="must_parames.indexType"
      @handleok="configurationHandel"
    />
    <traceability v-model:traceabilityVisible="traceabilityVisible" :select-alert="infoVal" />
    <alert-detail
      v-model:alert-detail-visible="alertDetailVisible"
      :module-enable="moduleEnable"
      :next="showNextBtn"
      :prev="showPrevBtn"
      :select-alert="infoVal"
      @on-skip-event="changeCurrentItemEvent"
    />
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
  </div>
</template>

<style scoped lang="scss">
  $criticalColor: #ff4340;
  $lowColor: #1b81fe;
  $midColor: #f1b04d;
  $highColor: #fa6d15;
  $defaultColor: #9d9aba;

  .shortcut {
    width: 100%;
    display: flex;
    // align-items: center;
    .word {
      width: 70px;
      margin-bottom: 8px;
      line-height: 28px;
    }
    .SQLMode {
      :deep() {
        width: 72px;
        align-items: inherit;
        margin-bottom: 8px;
        height: 28px;
        .el-radio-button__inner {
          padding: 5px 9px !important;
          height: 28px !important;
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
      margin: 0px 10px 8px 0px;
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
      margin-bottom: 8px;
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
        margin-bottom: 8px;
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

  :deep() {
    .el-table .red-row {
      background-color: #fff3f2 !important;
    }
    .el-table .bule-row {
      background-color: #f1f7ff !important;
    }
    .el-table .green-row {
      background-color: #f1fff2 !important;
    }
    .el-table .grep-row {
      background-color: #f4f4f4 !important;
    }
    .attackResult {
      .el-icon {
        font-size: 17px;
        vertical-align: -3px;
        margin-right: -4px;
      }
      &-0 {
        color: $criticalColor;
      }
      &-1 {
        color: $defaultColor;
      }
      &-2 {
        color: #ffa515;
      }
      &-3 {
        color: #6954f0;
      }
    }
  }

  .my-circular {
    position: relative;
    height: 24px;
    line-height: 24px;
    &:hover {
      cursor: pointer;
    }
    .circular {
      position: absolute;
      width: 12px;
      height: 12px;
      border-radius: 50%;
      left: -18px;
      top: 6px;
    }
  }
  .alert-container {
    // overflow: auto;
    // height: calc(100vh - 35px);
    height: calc(100vh - 10px);
    overflow: hidden;

    &::-webkit-scrollbar {
      width: 0;
      height: 0;
    }
  }
  .tabEllipsis {
    overflow: hidden;
    white-space: nowrap !important;
    text-overflow: ellipsis;
    -o-text-overflow: ellipsis;
    word-break: break-all;
    word-wrap: break-word;
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
  .alert_tag {
    display: inline-block;
    width: 55px;
    padding: 2px 3px;
    border-radius: 5px;
    color: #fff;
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
      background-color: rgba($defaultColor, 0.8);
      color: inherit;
    }
  }
  .item-icon {
    &:hover {
      cursor: pointer;
    }
  }
  .my-table {
    position: relative;
    // flex: 1;
    height: calc(100vh - 460px);
    :deep() {
      .el-table__body-wrapper {
        // overflow-y: auto;
        &::-webkit-scrollbar {
          width: 0 !important;
          height: 0;
        }
        .el-scrollbar__bar.is-vertical > div {
          margin-top: 40px;
        }
        // -moz-binding: url(‘ellipsis.xml#ellipsis’);
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
  }

  :deep() {
    .el-card__header {
      padding: 20px 20 14px;

      p {
        margin: 0;
      }
    }

    .el-card__body {
      padding-top: 0;
      display: flex;

      .echarts {
        flex: 1;
        width: auto;
        height: 275px;
        &.small {
          height: 120px;
        }
      }
    }

    .tools {
      margin-top: 12px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      .tools-left {
        flex: 1;
        display: flex;
        align-items: center;
      }
      .tools-right {
        min-width: 285px;
        display: flex;
        align-items: center;
      }
    }

    .table_action {
      display: flex;
      justify-content: flex-end;

      .reflesh {
        margin-right: 20px;
      }

      .action {
        margin-left: 10px;
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

    .risk_status,
    .el-table__cell .el-tag {
      &.critical {
        color: $criticalColor;
      }

      &.low {
        color: $lowColor;
      }

      &.mid {
        color: $midColor;
      }

      &.high {
        color: $highColor;
      }
    }

    .el-table__cell .el-tag {
      border: 0;
      background-color: rgba($criticalColor, 0.3);

      &.critical {
        background-color: rgba($criticalColor, 0.3);
      }

      &.low {
        background-color: rgba($lowColor, 0.3);
      }

      &.mid {
        background-color: rgba($midColor, 0.3);
      }

      &.high {
        background-color: rgba($highColor, 0.3);
      }
    }
    .el-form-item {
      width: 100%;
      .el-form-item__content {
        .el-input,
        .el-select {
          width: 100%;
        }
      }
    }
    .submit_item {
      .el-form-item__content {
        justify-content: space-between;
      }
    }
  }

  .log {
    margin: 0 10px;
    color: #6c6d6e;
    font-size: 14px;
  }

  .card-header-tag {
    width: 32px;
    height: 32px;
    border: 1px solid #e4e7ed;
    border-radius: 2px;
    display: flex;
    align-items: center;
    top: 10px !important;
    margin-left: 10px;
    cursor: pointer;
  }

  .form_icon {
    width: 20px;
    height: 20px;
    margin: 0 auto;
    transition: transform 0.2s;

    &.upward {
      transform: rotate(180deg);
    }
  }

  .chart_detail {
    width: 50%;
    display: inline-block;

    h5 {
      margin: 10px 0 10px;
      font-weight: 500;
      color: #303133;
    }
  }
  .flex-col {
    display: flex;
    flex-direction: column;
    flex: 1;
  }
</style>
