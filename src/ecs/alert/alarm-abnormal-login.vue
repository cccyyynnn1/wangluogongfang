<script lang="ts">
  export default {
    name: 'AbnormalLogin',
  }
</script>

<script setup lang="ts">
  import ApplicationConfiguration from '@/ecs/assets/components/application-configuration.vue'
  import AddSql from '../site/site-session/components/add-sql.vue'
  import TraceabilityFieldStatistic from '@/ecs/alert/components/field-statistic.vue'
  import { getAllDisPlaysFiledApi } from '@/api-ecs/public'
  import { TableColumnItemType } from '/#/store'
  import { useUserStore } from '@/store/modules/user'
  import AlertDetail from './components/alert-detail.vue'
  import { getSystemConfigApi } from '~/src/api-ecs/system'
  import { ElDivider } from 'element-plus'
  import { level_options, levelKey, timeDuratioOptions } from './data/index'
  import { timeDuratioOptions as date_options } from '@/data/constant'
  import { CloseBold } from '@element-plus/icons-vue'
  import AbnormalLoginRules from './components/abnormal-login-rules.vue'
  import { downloadFile } from '~/src/utils/download'
  import { formatNstime, formatTime } from '@/utils/time'
  import SearchSql from '@/components/search-sql/index.vue'
  import JsonPreview from '@/components/json-preview.vue'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import {
    AbnormalLandingAlarmDeleteByIdApi,
    AbnormalLandingAlarmQueryAllApi,
    AbnormalLandingAlarmUpdateStatusApi,
    AbnormalLandingAlarmImportApi,
    AbnormalLandingAlarmExportTemplateApi,
    AbnormalLandingAlarmExportApi,
    updateAlertStatusApi,
    AbnormalLandingAlarmSearchApi,
    AbnormalLandingAlarmQueryOneApi,
  } from '@/api-ecs/alert'
  import IntelligenceCenterWhiteUpload from './components/intelligence-center-white/intelligence-center-white-upload.vue'
  import dayjs from 'dayjs'
  import { AbnormalLandingAlarmSearchType, AlertItem, ShortcutListType, DowmloadType } from '~/src/types'
  import RegionalDistribution from './components/alarm-abnormal-login/regional-distribution.vue'
  import SourceIP from './components/alarm-abnormal-login/source-ip.vue'
  import TimeDistribution from './components/alarm-abnormal-login/time-distribution.vue'
  import { updateDisplayApi } from '~/src/api-ecs/custom-field'
  import { getTableCopyData } from '~/src/utils/transition'
  import { useTableCopy } from '~/src/utils'
  import { DownloadLogApi } from '~/src/api-ecs/toolbox'
  import _ from 'lodash'
  const multipleSelection = ref<AlertItem[]>([])
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const spacer = h(ElDivider, { direction: 'vertical' })
  const userStore = useUserStore()
  const { getTableColumn } = userStore
  const traceabilityFieldRef = ref<InstanceType<typeof TraceabilityFieldStatistic>>()
  const ruleVisible = ref(false)
  const ruleUploadVisible = ref(false)
  const tableColumn = ref<TableColumnItemType[]>([])
  const alarmIndexType = 34
  const alarmInfoVal = ref()
  const alarmInfoIndex = ref(-1)
  const alertDetailVisible = ref(false)
  const showTableSelection = ref(false)
  const attackResultType = ['成功', '失败', '未知', '企图']
  // const loading = ref(false)
  let endTimeList: any[] = []
  // 查询时间
  const timeDuration = ref('1hours')
  // 自定义时间
  const timeDate = ref()
  const shortcutList = ref<ShortcutListType[]>([])
  // 自动刷新
  const autoRefresh = ref(false)
  const autoRefreshDuration = ref(30)
  // 定时器ID
  const timerId = ref<NodeJS.Timer>()
  const showFiledConfig = ref(false) // 显示表头字段配置
  const rulesTooltipVisible = ref(false)
  const alarmTableRef = ref()
  const virtualRef = ref()
  const sqlComponentsRef = ref<InstanceType<typeof SearchSql>>()
  const selection = ref([])
  const listLoading = ref(true)
  const regionalDistributionRef = ref<InstanceType<typeof RegionalDistribution>>()
  const sourceIPRef = ref<InstanceType<typeof SourceIP>>()
  const timeDistributionRef = ref<InstanceType<typeof TimeDistribution>>()
  // 表格数据
  // const listDate = ref(dadaddadaad)
  // 表格数据
  const listTotal = ref(0)
  const showPrevBtn = ref(true)
  const showNextBtn = ref(true)
  // 高级筛选
  const isAdvanced = ref(true)
  // 告警检索条件
  const must_parames = reactive<AbnormalLandingAlarmSearchType>({
    /** 查询SQL */
    searchSql: '',
    /** 源IP */
    attackIp: '',
    // whiteType: 0,
    /** 源端口 */
    sourcePort: undefined,
    xff: '',
    url: '',
    host: '',
    /** 目的IP */
    victimIp: '',
    /** 目的端口 */
    targetPort: undefined,
    // /** 威胁类型 */
    threatType: '异常登录告警',
    whiteType: 0,
    // /** 威胁名称 */
    threatName: '',
    // /** 威胁等级 */
    // threatLevel: '',
    /** 攻击结果 */
    attackResult: [] as string[],
    /** 排序字段 */
    orderField: 'startTimeNs',
    /** 排序方式 */
    orderType: 'decs',
    /** 索引类型 */
    indexType: 34,
    /** 当前页码 */
    pageNum: 1,
    /** 当前页数 */
    pageSize: 20,
    /** 开始时间 */
    startTime: '',
    /** 结束时间 */
    endTime: '',
    // hasAggOneDocCount: true,
    // topCount: 100,
  })
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
  // 获取表格序号
  const curIndex = computed(() => (must_parames.pageNum - 1) * must_parames.pageSize + 1)
  const sortChange = async () => {
    must_parames.pageNum = 1
    queryData(true)
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
    const inputSql = must_parames.searchSql ? `( ${must_parames.searchSql} ) ` : ''
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
      initData()
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
    if (!must_parames.searchSql?.trim()) return ''
    if (ipv4Regex.test(must_parames.searchSql) || ipv6Regex.test(must_parames.searchSql)) {
      return `(源IP = '${must_parames.searchSql}' or 目的IP = '${must_parames.searchSql}')`
    }
    if (portRegex.test(must_parames.searchSql)) {
      return `(源端口 = '${must_parames.searchSql}' or 目的端口 = '${must_parames.searchSql}')`
    }
    return must_parames.searchSql
  }

  // 规则列表
  const searchName = ref('')
  const selectAll = ref(false)
  const ruleList = ref<any[]>([])
  const abnormalLandingAlarmQueryAll = async () => {
    const { data } = await AbnormalLandingAlarmQueryAllApi({ name: searchName.value })
    ruleList.value = data
    if (ruleList.value.length > 0) {
      ruleList.value.forEach((_: any, index: number) => {
        ruleList.value[index]['checked'] = false
      })
      handleCheckboxItemChange()
    }
  }

  const handelDel = (val: number) => {
    const ids: number[] = []
    if (val == -1) {
      ruleList.value.forEach((_: any, index: number) => {
        if (ruleList.value[index].checked) {
          ids.push(ruleList.value[index].id)
        }
      })
    } else {
      ids.push(val)
    }
    if (ids.length > 0) {
      $baseConfirm('你确定要删除选中项吗', null, async () => {
        const { msg } = await AbnormalLandingAlarmDeleteByIdApi({ ids })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        await abnormalLandingAlarmQueryAll()
      })
    }
  }

  const handleUpdateStatus = async (val: number) => {
    const ids: number[] = []
    ruleList.value.forEach((_: any, index: number) => {
      if (ruleList.value[index].checked) {
        ids.push(ruleList.value[index].id)
      }
    })
    if (ids.length > 0) {
      const { msg } = await AbnormalLandingAlarmUpdateStatusApi({
        ids,
        status: val,
      })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      await abnormalLandingAlarmQueryAll()
    }
  }

  const curData = ref()
  const handleEdit = (row: any) => {
    curData.value = row
    ruleVisible.value = true
  }

  const handleSelectAll = () => {
    if (ruleList.value.length > 0) {
      ruleList.value.forEach((_: any, index: number) => {
        ruleList.value[index]['checked'] = selectAll.value
      })
    }
  }

  const handleCheckboxItemChange = () => {
    const flag = ruleList.value.every((item: any) => {
      return item.checked
    })
    selectAll.value = flag
  }

  const handleShow = () => {
    abnormalLandingAlarmQueryAll()
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
  const showAlertDetail = async (row: any, index: number) => {
    alarmInfoIndex.value = index
    alarmInfoVal.value = row
    alertDetailVisible.value = true
  }
  const moduleEnable = ref(false)
  const getModelStatus = async () => {
    const {
      data: { module_config },
    } = await getSystemConfigApi({ keys: 'module_config' })
    moduleEnable.value = JSON.parse(module_config.value).isEnable || false
  }
  const changeCurrentItemEvent = async (isAdd: boolean) => {
    const index = isAdd
      ? Math.min(alarmInfoIndex.value + 1, alertInfo.list.length)
      : Math.max(alarmInfoIndex.value - 1, 0)
    alarmInfoIndex.value = index
  }
  // 获取用户所有字段
  const getAllDisPlaysFiled = async () => {
    listLoading.value = true
    const { data } = await getAllDisPlaysFiledApi()
    const alarmDisplays = data[alarmIndexType] as number[]
    const alarmAllFields = getTableColumn(alarmIndexType)
    const fields = alarmDisplays.map((alarm) => alarmAllFields.find((item) => item.id === alarm))
    tableColumn.value = fields as TableColumnItemType[]
    listLoading.value = false
  }
  watch(
    () => alarmInfoIndex.value,
    (index) => {
      if (alarmInfoIndex.value === 0) {
        showPrevBtn.value = false
        showNextBtn.value = true
      }
      if (alarmInfoIndex.value === alertInfo.list.length - 1) {
        showNextBtn.value = false
        showPrevBtn.value = true
      }
      if (alarmInfoIndex.value > alertInfo.list.length - 5) {
        console.log('getData')
      }
      alarmInfoVal.value = alertInfo.list[index]
    }
  )

  const handleAdd = () => {
    curData.value = null
    ruleVisible.value = true
  }

  const handleExport = async () => {
    const ids: number[] = []
    ruleList.value.forEach((_: any, index: number) => {
      if (ruleList.value[index].checked) {
        ids.push(ruleList.value[index].id)
      }
    })
    if (ids.length > 0) {
      const res = await AbnormalLandingAlarmExportApi({ ids })
      downloadFile(res, '规则列表')
    }
  }

  const lastPageBtnNode = ref()
  const getLastPageBtnNode = () => {
    const nodes = document.querySelectorAll('.el-pager li')
    if (nodes.length > 5) {
      lastPageBtnNode.value = nodes[nodes.length - 1] as HTMLElement
      lastPageBtnNode.value.setAttribute('style', 'display: none')
      if (must_parames.pageNum >= Math.ceil(listTotal.value / must_parames.pageSize) - 5) {
        lastPageBtnNode.value?.setAttribute('style', 'display: block')
      } else {
        lastPageBtnNode.value?.setAttribute('style', 'display: none')
      }
      // currentChange()
    }
  }

  const currentChangeFalg = ref(false)
  const currentChangeInitFalg = ref(true)
  const currentChange = () => {
    if (must_parames.pageNum >= Math.ceil(listTotal.value / must_parames.pageSize) - 5) {
      lastPageBtnNode.value?.setAttribute('style', 'display: block')
    } else {
      lastPageBtnNode.value?.setAttribute('style', 'display: none')
    }
    // if (currentChangeInitFalg.value) return (currentChangeInitFalg.value = false)
    queryData(true)
  }

  onMounted(() => {
    getModelStatus()
    setTimeout(() => {
      initData()
      getAllDisPlaysFiled()
    }, 0)
  })

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

  watch(
    () => timeDuration.value,
    () => {
      formatDate()
      endTimeList = []
    },
    {
      immediate: true,
    }
  )

  const formatUserTime = () => {
    const [startDate, endDate] = timeDate.value
    must_parames.endTime = dayjs(endDate).format('YYYY-MM-DD HH:mm:ss')
    must_parames.startTime = dayjs(startDate).format('YYYY-MM-DD HH:mm:ss')
  }

  const initData = async () => {
    await queryData()
    if (sqlComponentsRef.value && !sqlComponentsRef.value?.getValidateSql()) return
    sourceIPRefDataChange()
    timeDistributionRefDataChange()
    regionalDistributionRefDataChange()
    setTimeout(() => {
      getLastPageBtnNode()
    }, 0)
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
      timerId.value = setInterval(() => {
        must_parames.pageNum = 1
        queryData()
      }, autoRefreshDuration.value * 1000)
    } else {
      // @ts-ignore
      clearInterval(timerId.value)
    }
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
    // threatType: '威胁类型',
    attackResult: '攻击结果',
    xff: 'XFF',
    url: 'URL',
  }
  const formatSql = () => {
    const arr1 = Object.keys(sqlObj)
    const obj = {} as { [key: string]: any }
    for (const key in must_parames) {
      if (arr1.includes(key)) {
        // @ts-ignore
        if (must_parames[key]) {
          // @ts-ignore
          obj[sqlObj[key]] = must_parames[key]
        }
      }
    }
    // if (!obj.威胁类型?.length) delete obj.威胁类型
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

  const shortcutSQL = ref('')
  const queryData = async (remark = false) => {
    formatDate()
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
      if (must_parames.searchSql?.trim()) {
        shortcutSQL.value = `(${probeSearchSqlIsIPOrPort()})` + ` ` + `and` + ` ${shortcut}`
      } else {
        if (shortcutSQL.value?.trim()) {
          shortcutSQL.value = `${shortcutSQL.value}` + ` ` + `and` + ` ${shortcut}`
        } else {
          shortcutSQL.value = shortcut
        }
      }
    } else {
      if (must_parames.searchSql?.trim()) {
        if (!shortcutSQL.value) {
          shortcutSQL.value = JSON.parse(JSON.stringify(probeSearchSqlIsIPOrPort()))
        } else {
          shortcutSQL.value =
            `${shortcutSQL.value}` + ` ` + `and` + ` (${JSON.parse(JSON.stringify(probeSearchSqlIsIPOrPort()))})`
        }
      }
    }
    shortcutSQL.value += shortcutSQL.value.trim() ? ` and (威胁类型 = "异常登录告警")` : `(威胁类型 = "异常登录告警")`
    listLoading.value = true
    try {
      const _index = Math.floor(must_parames.pageNum / (10000 / must_parames.pageSize))
      const curTimeObj = endTimeList[_index]
      let _value: any = undefined
      let _page = 0
      if (curTimeObj) {
        const _name = Object.keys(curTimeObj)[0]
        _value = curTimeObj[_name]
        if (_name && +_name != 1) {
          _page = 10000 / must_parames.pageSize - (+_name % (10000 / must_parames.pageSize))
        }
      }
      const {
        data: { total, resList },
      } = await AbnormalLandingAlarmSearchApi({
        ...must_parames,
        attackResult: must_parames.attackResult.join(','),
        searchSql: shortcutSQL.value,
        pageNum:
          must_parames.pageNum % (10000 / must_parames.pageSize) == 0
            ? 10000 / must_parames.pageSize
            : (must_parames.pageNum % (10000 / must_parames.pageSize)) + _page,
        total: remark ? listTotal.value : undefined,
        endTime: _value || must_parames.endTime,
      })
      alertInfo.list = resList
      listTotal.value = total
      recordingTime()
      sqlComponentsRef.value?.changeHistories(shortcutSQL.value)
      listLoading.value = false
    } catch (error) {
      listLoading.value = false
      console.error(error)
    }
  }

  const recordingTime = () => {
    if (must_parames.pageNum == 1) {
      const item = alertInfo.list[0]
      const time = item?.startTimeNs
      if (time) {
        endTimeList[0] = { [must_parames.pageNum]: must_parames.endTime }
      }
    }
    if (10000 / must_parames.pageSize - (must_parames.pageNum % (10000 / must_parames.pageSize)) < 5) {
      const _index = Math.ceil(must_parames.pageNum / (10000 / must_parames.pageSize))
      if (_index) {
        const item = alertInfo.list[alertInfo.list.length - 1]
        const time = item?.startTimeNs
        const pre_index = must_parames.pageNum
        if (!currentChangeFalg.value) {
          endTimeList[_index] = { [pre_index]: dayjs(+time / 1000000).format('YYYY-MM-DD HH:mm:ss') }
          currentChangeFalg.value = true
        }
      }
    } else {
      currentChangeFalg.value = false
    }
  }

  const regionalDistributionRefDataChange = () => {
    regionalDistributionRef.value?.initData({
      searchSql: shortcutSQL.value.trim() ? `${shortcutSQL.value}` : `(威胁类型 = "异常登录告警")`,
      startTime: must_parames.startTime,
      // startTime: '2024-07-18 08:38:35',
      endTime: must_parames.endTime,
      // endTime: '2024-09-19 09:48:35',
      indexType: must_parames.indexType,
      topCount: 100,
      whiteType: 0,
      aggregationFields:
        regionalDistributionRef.value?.selectValue == 'client_province'
          ? [regionalDistributionRef.value?.selectValue, 'client_longitude', 'client_latitude']
          : [regionalDistributionRef.value?.selectValue],
    })
  }

  const sourceIPRefDataChange = () => {
    sourceIPRef.value?.initData({
      searchSql: shortcutSQL.value.trim() ? `${shortcutSQL.value}` : `(威胁类型 = "异常登录告警")`,
      startTime: must_parames.startTime,
      // startTime: '2024-07-18 08:38:35',
      endTime: must_parames.endTime,
      // endTime: '2024-09-19 09:48:35',
      indexType: must_parames.indexType,
      topCount: 100,
      whiteType: 0,
      aggregationFields: sourceIPRef.value?.selectValue,
    })
  }

  const timeDistributionRefDataChange = () => {
    timeDistributionRef.value?.initData({
      ...must_parames,
      searchSql: shortcutSQL.value,
      attackResult: must_parames.attackResult.join(','),
    })
  }

  const handleDownloadExcel = async (isAll = false) => {
    if (isAll) {
      const obj: DowmloadType = {
        downloadCnd: shortcutSQL.value,
        edTime: dayjs(must_parames.endTime).valueOf(),
        indexType: must_parames.indexType,
        stTime: dayjs(must_parames.startTime).valueOf(),
        type: 0,
        dataType: 1,
      }
      try {
        const { msg } = await DownloadLogApi({ ...obj }, 1)
        $baseMessage(msg, 'success', 'vab-hey-message-success')
      } catch (error) {
        $baseMessage('下载失败', 'error', 'vab-hey-message-error')
      }
    } else {
      const tHeader: string[] = []
      const filterVal: any = []
      try {
        tableColumn.value.forEach(({ fieldNameCn, fieldNameEn }) => {
          tHeader.push(fieldNameCn)
          filterVal.push(fieldNameEn)
        })
        listLoading.value = true
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
          listLoading.value = false
        })
      } catch (error) {
        console.log(error)
        listLoading.value = false
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
      case 'warnTime':
      case 'startTimeNs':
        return formatNstime(val[key])
      default:
        return val[key] || ''
    }
  }

  const handleHiddenSelection = () => {
    showTableSelection.value = false
    multipleSelection.value = []
    alarmTableRef.value?.clearSelection()
  }
  /**
   * type: 标记忽略1 取消忽略0
   */
  const updateIgnoreStatus = async (status: '1' | '0') => {
    const ids = multipleSelection.value.map((i) => i.id)
    const { msg } = await updateAlertStatusApi({ ...must_parames, ignoreStatus: +status, ids })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    multipleSelection.value = []
    queryData(true)
  }

  /**
   * type: 标记已读1 标记未读0
   */
  const updateStatus = async (status: '1' | '0') => {
    const ids = multipleSelection.value.map((i) => i.id)
    const { msg } = await updateAlertStatusApi({ ...must_parames, readStatus: +status, ids })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    multipleSelection.value = []
    queryData(true)
  }

  const searchhandle = (val: string) => {
    initData()
  }
  // 清空
  const handleEmpty = () => {
    shortcutList.value = []
  }

  let BASICDATA: AlertItem[] = []
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

  const fieldDetailVisible = ref(false)
  const fieldDetailTitle = ref('')
  const fieldDetailValue = ref('')
  const showFieldHandle = (title: string, val: string) => {
    fieldDetailVisible.value = true
    fieldDetailTitle.value = title
    fieldDetailValue.value = JSON.stringify(formatSourceData(val), null, 4)
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

  // 表头字段配置
  const tableHeadConfig = (val: boolean) => {
    showFiledConfig.value = val
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
    queryData(true)
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

  const useTableCopyEvent = (row: any, column: any, cell: any, event: any) => {
    const arr = getTableCopyData({ tableColumn: tableColumn.value, row, column, mothod: handleUpdateCallback })
    useTableCopy(row, column, cell, event, arr)
  }
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

  // 多选表格
  function handleSelectionChange(val: AlertItem[]) {
    multipleSelection.value = val
  }

  function formatDataTime(row: any, key: string) {
    const time = row[key] / 1000000
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss.SSS')
  }

  function formatDataTimeDec(time: any) {
    if (!time) return '-'
    const times = time
    return dayjs(times).format('YYYY-MM-DD HH:mm:ss')
  }

  const showTraceabilityField = (column: TableColumnItemType) => {
    const { indexType, startTime, endTime } = must_parames
    const { fieldNameCn, fieldNameEn } = column
    const whiteType = 0
    const SQL = shortcutSQL.value
    traceabilityFieldRef.value?.initData({
      title: fieldNameCn,
      query: { indexType, startTime, endTime, aggregationFields: fieldNameEn, searchSql: SQL, whiteType },
      column,
    })
  }

  const router = useRouter()
  const fullFlowSurveyHandle = async (row: any) => {
    alarmInfoVal.value = row
    const { clientIp, serverIp, sourcePort, targetPort, startTimeNs, protocol } = alarmInfoVal.value
    const time = startTimeNs / 1000000
    const sqlStr = `源ip = "${clientIp}" and 源端口 = "${sourcePort}" and 目的ip = "${serverIp}" and 目的端口 = "${targetPort}"`
    const indexType = 1
    const start = formatTime(dayjs(time).subtract(30, 'minute'))
    const end = formatTime(dayjs(time).add(30, 'minute'))
    const resolveRouter = router.resolve({
      path: '/retrieve/index',
      query: {
        info: encodeURIComponent(
          JSON.stringify({
            sql: sqlStr,
            indexType,
            timeRanges: [start, end],
            protocol,
          })
        ),
      },
    })
    window.open(resolveRouter.href, '_blank')
  }

  watch(
    () => isAdvanced.value,
    () => {
      must_parames.searchSql = ''
    }
  )

  const curId = ref(-1)
  const popperValue = ref({
    incident: '',
  })
  const infoObj = {
    httpLogin: 'HTTP登录',
    sql: '数据库登录',
    mail: '邮箱登录',
    msrdp: 'RDP登录',
    sshLog: 'SSH登录',
  }
  const getConfidence = (code: string) => {
    return ['低', '中', '高'].findIndex((i) => code === i) + 1
  }
  const queryRuleIdOrThreatName = async (type: string, value: number | string, row: any) => {
    if (curId.value == row.id) return
    popperValue.value['incident'] = row.abnormal_incident
    curId.value = row.id
    // const { data } = await AbnormalLandingAlarmQueryOneApi({ [type]: type == 'id' ? +value : value })
    // console.log(row, 'rwo')
    // console.log(data, 'data')
    // popperValue.value = data
    // if (typeof data == 'string') return $baseMessage(data, 'error', 'vab-hey-message-error')
    getRuleIdOrThreatNameRemark(row)
  }

  const getRuleIdOrThreatNameRemark = (row: any) => {
    //
    // 目的IP：相等；源IP地理位置：不相等；时间：不相等；适配事件：相等。
    // const { abnormal_incident } = row
    // const { country, province, city, startTime, endTime, timeStatus } = popperValue.value
    // if (abnormal_country != country || abnormal_province != province || abnormal_city != city)
    //   popperValue.value['addressFlag'] = true
    // const time = `${startTime || '00:00'}-${endTime || '23:59'}(${timeStatus || '每天'})`
    // if (abnormal_time != time) popperValue.value['timeFlag'] = true
    // const incidentArr = abnormal_incident?.split(',') || []
    // console.log(incidentArr)
    // incidentArr.forEach((item: string) => {
    //   const regExp = new RegExp(item, 'g')
    //   popperValue.value.incident = popperValue.value.incident.replaceAll(
    //     regExp,
    //     `<span class="abnormal-login-popper-in-table-tip">${item}</span>`
    //   )
    // })
    if (row.protocol) {
      // @ts-ignore
      const item = infoObj[row.protocol]
      const regExp = new RegExp(item, 'g')
      popperValue.value.incident = popperValue.value.incident.replaceAll(
        regExp,
        `<span class="abnormal-login-popper-in-table-tip">${item}</span>`
      )
    }
  }

  const scrollRef = ref<HTMLElement | null>(null)
  const { isScrolling } = useScroll(scrollRef)
  const scrollFalg = ref(false)
  watch(
    () => isScrolling.value,
    () => {
      scrollFalg.value = isScrolling.value
    }
  )

  const isWrapScrollref = ref()
  nextTick(() => {
    isWrapScrollref.value = document.querySelector('.alarm-table .el-scrollbar__wrap') as HTMLElement
    if (isWrapScrollref.value) {
      isWrapScrollref.value.addEventListener('scroll', handelScrollEvent)
    }
  })

  const handelScrollEvent = _.debounce(() => {
    if (!scrollFalg.value) {
      scrollFalg.value = true
    }
    setTimeout(() => {
      scrollFalg.value = false
    }, 100)
  }, 200)

  onUnmounted(() => {
    isWrapScrollref.value.removeEventListener('scroll', handelScrollEvent)
  })
</script>

<template>
  <div class="abnormal-login-container">
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3 class="title">异常登录告警</h3>
        <el-form inline label-position="left" style="margin-left: 15px">
          <el-form-item style="width: 148px">
            <el-select v-model="timeDuration">
              <el-option v-for="item in date_options" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item v-if="timeDuration === 'user-defined'" style="width: 300px">
            <vab-date-time-picker v-model="timeDate" style="margin-left: 5px" />
          </el-form-item>
        </el-form>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12">
        <div class="auto-refresh">
          <label>自动刷新：</label>
          <el-switch v-model="autoRefresh" />
        </div>
        <!-- </div> -->
        <el-select
          v-model="autoRefreshDuration"
          class="auto-refresh-selset"
          :disabled="autoRefresh"
          style="width: 85px; margin-right: 10px"
        >
          <el-option v-for="item in timeDuratioOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-button type="primary" @click="handleAdd">添加规则</el-button>
        <el-button ref="virtualRef" style="margin-right: 10px !important" type="primary">规则列表</el-button>
        <el-dropdown trigger="click">
          <el-button @click="showTableSelection = true">
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
              <el-dropdown-item :disabled="multipleSelection.length === 0" @click="handleDownloadExcel(false)">
                导出选择日志
              </el-dropdown-item>
              <el-dropdown-item @click="handleDownloadExcel(true)">导出全部日志</el-dropdown-item>
              <el-dropdown-item @click="handleHiddenSelection">隐藏多选框</el-dropdown-item>
              <el-dropdown-item :disabled="multipleSelection.length === 0" divided @click="updateStatus('1')">
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
      </vab-query-form-right-panel>
      <div class="search-sql">
        <div class="search-value">
          <SearchSql
            v-if="isAdvanced"
            ref="sqlComponentsRef"
            :index-type="must_parames.indexType"
            :model-value="must_parames.searchSql"
            @on-change="(str) => (must_parames.searchSql = str)"
            @onSearch="searchhandle"
          />
        </div>
        <div class="search-tool">
          <!-- <el-row style="width: 100%"> -->
          <el-form inline :label-position="isAdvanced ? 'left' : 'top'">
            <el-form-item class="submit_item" label="&nbsp;" style="margin-right: 0px; margin-bottom: 0px !important">
              <el-button type="primary" @click="initData">检索</el-button>
              <el-button @click="tableHeadConfig(true)">
                <el-icon><Setting /></el-icon>
              </el-button>
            </el-form-item>
            <el-form-item
              v-if="!isAdvanced"
              class="submit_item"
              label="&nbsp;"
              style="margin-right: 0px; margin-bottom: 0px !important"
            >
              <div style="height: 46px; width: 200px"></div>
            </el-form-item>
          </el-form>
          <!-- </el-row> -->
        </div>
      </div>
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
    <div
      ref="scrollRef"
      class="content-warp"
      :style="{ height: isAdvanced ? 'calc(100vh - 206px)' : 'calc(100vh - 310px)' }"
    >
      <div class="charts-warp">
        <el-row>
          <el-col :span="6">
            <RegionalDistribution ref="regionalDistributionRef" @refresh="regionalDistributionRefDataChange" />
          </el-col>
          <el-col :span="6"><SourceIP ref="sourceIPRef" @refresh="sourceIPRefDataChange" /></el-col>
          <el-col :span="12"><TimeDistribution ref="timeDistributionRef" /></el-col>
        </el-row>
      </div>
      <el-table
        ref="alarmTableRef"
        v-loading="listLoading"
        class="my-table alarm-table"
        :data="alertInfo.list"
        element-loading-text="Loading..."
        :row-class-name="tableRowClassName"
        @cell-contextmenu="useTableCopyEvent"
        @selection-change="handleSelectionChange"
      >
        <el-table-column v-if="showTableSelection" type="selection" width="55" />
        <el-table-column :index="curIndex" label="序号" type="index" width="75" />
        <el-table-column
          v-for="item in tableColumn"
          :key="item.id"
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
              'abnormal_dst_ip',
            ].includes(item.fieldNameEn)
              ? 220
              : 145
          "
        >
          <template #header>
            <span :style="{ cursor: 'pointer' }">
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
          <template #default="{ row }" v-else-if="['ruleId'].includes(item.fieldNameEn)">
            <div :style="{ display: 'flex', alignItems: 'center', width: '100%' }">
              <span>{{ row[item.fieldNameEn] }}</span>
              <el-popover
                :disabled="scrollFalg"
                :persistent="false"
                placement="right"
                popper-class="abnormal-login-popper-in-table"
                :show-arrow="false"
                trigger="click"
                :width="400"
              >
                <template #reference>
                  <img
                    v-if="row[item.fieldNameEn]"
                    alt=""
                    src="@/assets/alert_images/abnormal_landing_alarm/detail.png"
                    :style="{ width: '20px', height: '20px', cursor: 'pointer', marginLeft: '4px' }"
                    @click="queryRuleIdOrThreatName('id', row[item.fieldNameEn], row)"
                  />
                </template>
                <div class="content">
                  <div class="title-head">
                    <div v-if="row.ruleStatus == 'enable'" class="tag">启用</div>
                    <div v-else-if="row.ruleStatus == 'disable'" class="tag disablement">停用</div>
                    <div :style="{ fontWeight: 500, fontSize: '15px', color: '#1e1842' }">
                      规则名称：
                      <span>{{ row.abnormal_name }}</span>
                      <span v-if="row.ruleStatus == 'not_exist'">（已删除）</span>
                    </div>
                  </div>
                  <div class="rule-info">
                    <div class="rule-info-ul">
                      <div class="rule-info-li">
                        <img
                          :src="require('@/assets/alert_images/abnormal_landing_alarm/ip.svg')"
                          :style="{ height: '14px', width: '14px' }"
                        />
                        <div class="info info-checked">
                          <div style="display: flex">
                            <div style="width: 70px">目的IP：</div>
                            <div class="abnormal-login-popper-in-table-tip">{{ row.abnormal_dst_ip }}</div>
                          </div>
                        </div>
                      </div>
                      <div class="rule-info-li">
                        <img
                          :src="require('@/assets/alert_images/abnormal_landing_alarm/address.svg')"
                          :style="{ height: '14px', width: '14px' }"
                        />
                        <div class="info">
                          <div style="display: flex">
                            <div style="width: 70px">地址：</div>
                            <div
                              :class="{
                                'abnormal-login-popper-in-table-tip':
                                  !!row.abnormal_city || !!row.abnormal_province || !!row.abnormal_country,
                              }"
                            >
                              {{
                                `${row.abnormal_country}${row.abnormal_province ? `-${row.abnormal_province}` : ''}` +
                                  `${row.abnormal_city ? `-${row.abnormal_city}` : ''}` || '-'
                              }}
                            </div>
                          </div>
                        </div>
                      </div>
                      <div class="rule-info-li">
                        <el-image
                          :src="require('@/assets/alert_images/abnormal_landing_alarm/time.svg')"
                          :style="{ height: '14px', width: '14px' }"
                        />
                        <!-- formatTimeStatus(item.timeStatus) -->
                        <div class="info info-checked">
                          <div style="display: flex">
                            <div style="width: 70px">时间：</div>
                            <div :class="{ 'abnormal-login-popper-in-table-tip': !!row.abnormal_time }">
                              <!-- {{
                                `${popperValue.startTime || '00:00'}-${popperValue.endTime || '23:59'}（${
                                  popperValue.timeStatus || '每天'
                                }）` || '00:00-23:59'
                              }} -->
                              {{ row.abnormal_time || '00:00-23:45(每天)' }}
                            </div>
                          </div>
                        </div>
                      </div>
                      <div class="rule-info-li rule-info-li-sp">
                        <div>
                          <img
                            :src="require('@/assets/alert_images/abnormal_landing_alarm/pack.svg')"
                            :style="{ height: '14px', width: '14px', marginTop: '-4px' }"
                          />
                        </div>
                        <div class="info">
                          <div style="display: flex">
                            <div style="width: 80px">适配事件：</div>
                            <div v-html="popperValue.incident || '-'"></div>
                          </div>
                        </div>
                      </div>
                      <div class="rule-info-li">
                        <el-image
                          :src="require('@/assets/alert_images/abnormal_landing_alarm/time.svg')"
                          :style="{ height: '14px', width: '14px' }"
                        />
                        <!-- formatTimeStatus(item.timeStatus) -->
                        <div class="info info-checked">
                          <div style="display: flex">
                            <div style="width: 70px">创建时间：</div>
                            <div>{{ formatDataTimeDec(+row.ruleUpdateTime) }}</div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </el-popover>
            </div>
          </template>
          <template #default="{ row }" v-else-if="['threatName'].includes(item.fieldNameEn)">
            <div style="display: flex">
              <div class="threatName-row">
                {{ row[item.fieldNameEn] }}
              </div>
              <el-popover
                :disabled="scrollFalg"
                :persistent="false"
                placement="right"
                popper-class="abnormal-login-popper-in-table"
                :show-arrow="false"
                trigger="click"
                :width="400"
              >
                <template #reference>
                  <img
                    v-if="row[item.fieldNameEn]"
                    alt=""
                    src="@/assets/alert_images/abnormal_landing_alarm/detail.png"
                    style="width: 20px; height: 20px; cursor: pointer; margin-left: 4px"
                    @click="queryRuleIdOrThreatName('name', row[item.fieldNameEn], row)"
                  />
                </template>
                <div v-if="popperValue" class="content">
                  <div class="title-head">
                    <!-- <div v-if="popperValue?.status == 1" class="tag">启用</div> -->
                    <!-- <div v-else-if="popperValue?.status == 0" class="tag disablement">停用</div> -->
                    <div style="font-weight: 500; font-size: 15px; color: #1e1842">
                      规则名称：{{ row.abnormal_name }}
                    </div>
                  </div>
                  <div class="rule-info">
                    <div class="rule-info-ul">
                      <div class="rule-info-li">
                        <img
                          :src="require('@/assets/alert_images/abnormal_landing_alarm/ip.svg')"
                          :style="{ height: '14px', width: '14px' }"
                        />
                        <div class="info info-checked">
                          <div style="display: flex">
                            <div style="width: 70px">目的IP：</div>
                            <div>{{ row.abnormal_dst_ip }}</div>
                          </div>
                        </div>
                      </div>
                      <div class="rule-info-li">
                        <img
                          :src="require('@/assets/alert_images/abnormal_landing_alarm/address.svg')"
                          :style="{ height: '14px', width: '14px' }"
                        />
                        <div class="info">
                          <div style="display: flex">
                            <div style="width: 70px">地址：</div>
                            <div
                              :class="{
                                'abnormal-login-popper-in-table-tip':
                                  !!row.abnormal_city || !!row.abnormal_province || !!row.abnormal_country,
                              }"
                            >
                              {{
                                `${row.abnormal_country}${row.abnormal_province ? `-${row.abnormal_province}` : ''}` +
                                  `${row.abnormal_city ? `-${row.abnormal_city}` : ''}` || '-'
                              }}
                            </div>
                          </div>
                        </div>
                      </div>
                      <div class="rule-info-li">
                        <el-image
                          :src="require('@/assets/alert_images/abnormal_landing_alarm/time.svg')"
                          :style="{ height: '14px', width: '14px' }"
                        />
                        <!-- formatTimeStatus(item.timeStatus) -->
                        <div class="info info-checked">
                          <div style="display: flex">
                            <div style="width: 70px">时间：</div>
                            <div>
                              {{ row.abnormal_time || '00:00-23:45(每天)' }}
                            </div>
                          </div>
                        </div>
                      </div>
                      <div class="rule-info-li rule-info-li-sp">
                        <div>
                          <img
                            :src="require('@/assets/alert_images/abnormal_landing_alarm/pack.svg')"
                            :style="{ height: '14px', width: '14px', marginTop: '-4px' }"
                          />
                        </div>
                        <div class="info">
                          <div style="display: flex">
                            <div style="width: 80px">适配事件：</div>
                            <div v-html="popperValue.incident || '-'"></div>
                          </div>
                        </div>
                      </div>
                      <div class="rule-info-li">
                        <el-image
                          :src="require('@/assets/alert_images/abnormal_landing_alarm/time.svg')"
                          :style="{ height: '14px', width: '14px' }"
                        />
                        <!-- formatTimeStatus(item.timeStatus) -->
                        <div class="info info-checked">
                          <div style="display: flex">
                            <div style="width: 70px">创建时间：</div>
                            <div>{{ formatDataTimeDec(+row.ruleUpdateTime) }}</div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </el-popover>
            </div>
          </template>
          <template #default="{ row }" v-else-if="['sourceData'].includes(item.fieldNameEn)">
            <div class="tabEllipsis" @click="() => showFieldHandle(item.fieldNameCn, row[item.fieldNameEn])">
              {{ row[item.fieldNameEn] }}
            </div>
          </template>
          <template #default="{ row }" v-else-if="item.fieldNameEn == 'threatLevel'">
            <span :class="['alert_tag', getThreatLevel(row.threatLevel)]">
              <el-icon><WarnTriangleFilled /></el-icon>
              {{ row.threatLevel }}
            </span>
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

          <template #default="{ row }" v-else-if="item.fieldNameCn == '置信度'">
            <el-rate
              :colors="['#67C23A', '#67C23A', '#67C23A']"
              disabled
              disabled-void-color="#C7C6D4"
              :max="3"
              :model-value="getConfidence(row.confidence)"
            />
          </template>
        </el-table-column>

        <el-table-column fixed="right" label="操作" width="80">
          <template #default="{ row }">
            <el-button class="row_action" size="small" @click="fullFlowSurveyHandle(row)">调查</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <el-pagination
      v-model:current-page="must_parames.pageNum"
      v-model:page-size="must_parames.pageSize"
      background
      class="known_pagination"
      layout="total, sizes, prev, pager, next"
      :page-sizes="[20, 50, 100]"
      :total="listTotal"
      @current-change="currentChange"
      @size-change="sortChange"
    />

    <!-- 表头字段配置 -->
    <application-configuration
      v-if="showFiledConfig"
      v-model="showFiledConfig"
      :fields="tableColumn"
      :retrieve-index-type="must_parames.indexType"
      @handleok="configurationHandel"
    />

    <vab-dialog v-model="fieldDetailVisible" destroy-on-close :title="fieldDetailTitle" width="1375px">
      <div class="mask">
        <json-preview :json-value="fieldDetailValue" />
      </div>
    </vab-dialog>

    <alert-detail
      v-model:alert-detail-visible="alertDetailVisible"
      :module-enable="moduleEnable"
      :next="showNextBtn"
      :prev="showPrevBtn"
      :select-alert="alarmInfoVal"
      @on-skip-event="changeCurrentItemEvent"
    />
    <abnormal-login-rules v-model:visible="ruleVisible" :rule-data="curData" />
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
    <el-popover
      placement="top"
      popper-class="abnormal-login-popper"
      :show-arrow="false"
      :teleported="false"
      trigger="click"
      :virtual-ref="virtualRef"
      virtual-triggering
      width="420"
      @show="handleShow"
    >
      <template #default>
        <div class="rule-title rule-title-head">
          <span>规则列表</span>
          <el-button :auto-insert-space="false" :icon="CloseBold" link type="primary" @click="handleAdd">
            添加
          </el-button>
        </div>
        <div class="rule-list">
          <vab-input-search
            v-model="searchName"
            clearable
            placeholder="请输入规则名称"
            :style="{ marginInline: '8px', width: '380px', marginBottom: '8px' }"
            @on-search="abnormalLandingAlarmQueryAll"
          />
          <el-collapse accordion class="collapse-items">
            <el-collapse-item v-for="item in ruleList" :key="item.id">
              <template #title>
                <div class="title-head">
                  <el-checkbox
                    v-model="item.checked"
                    size="large"
                    :style="{ marginInline: '15px 5px' }"
                    @change="handleCheckboxItemChange"
                    @click.stop="() => {}"
                  />
                  <div v-if="item.status == 1" class="tag">启用</div>
                  <div v-else-if="item.status == 0" class="tag disablement">停用</div>
                  <div>{{ item.name }}</div>
                </div>
              </template>
              <div class="rule-info">
                <div class="rule-info-ul">
                  <div class="rule-info-li">
                    <img
                      :src="require('@/assets/alert_images/abnormal_landing_alarm/ip.svg')"
                      :style="{ height: '14px', width: '14px' }"
                    />
                    <div class="info info-checked">{{ item.serverIp }}</div>
                  </div>
                  <div class="rule-info-li">
                    <img
                      :src="require('@/assets/alert_images/abnormal_landing_alarm/address.svg')"
                      :style="{ height: '14px', width: '14px' }"
                    />
                    <div class="info">{{ item.city || item.province || item.country || '-' }}</div>
                  </div>
                  <div class="rule-info-li">
                    <el-image
                      :src="require('@/assets/alert_images/abnormal_landing_alarm/time.svg')"
                      :style="{ height: '14px', width: '14px' }"
                    />
                    <!-- formatTimeStatus(item.timeStatus) -->
                    <div class="info info-checked">
                      {{
                        `${item.startTime || '00:00'}-${item.endTime || '23:59'}（${item.timeStatus || '每天'}）` ||
                        '00:00-23:59'
                      }}
                    </div>
                  </div>
                  <!-- <div class="rule-info-li">
                    <el-image
                      :src="require('@/assets/alert_images/abnormal_landing_alarm/assets.svg')"
                      :style="{ height: '14px', width: '14px' }"
                    />
                    <div class="info">{{ item.assetArray || '-' }}</div>
                  </div> -->
                  <div class="rule-info-li">
                    <img
                      :src="require('@/assets/alert_images/abnormal_landing_alarm/pack.svg')"
                      :style="{ height: '14px', width: '14px' }"
                    />
                    <div class="info">{{ item.incident || '-' }}</div>
                  </div>
                </div>
                <el-space class="my-spacer" :spacer="spacer">
                  <span @click="handleEdit(item)">编辑</span>
                  <span @click="handelDel(item.id)">删除</span>
                </el-space>
              </div>
            </el-collapse-item>
          </el-collapse>
        </div>
        <div class="rule-footer">
          <div class="left">
            <el-checkbox v-model="selectAll" label="全选" size="large" @change="handleSelectAll" @click.stop="" />
          </div>
          <div class="right">
            <el-button class="img-btn" @click="ruleUploadVisible = true">
              <el-tooltip class="item" content="导入" effect="dark" placement="top">
                <img
                  :src="require('@/assets/alert_images/abnormal_landing_alarm/export.svg')"
                  :style="{ height: '14px', width: '14px' }"
                />
              </el-tooltip>
            </el-button>
            <el-button class="img-btn" @click="handleExport">
              <el-tooltip class="item" content="导出" effect="dark" placement="top">
                <img
                  :src="require('@/assets/alert_images/abnormal_landing_alarm/import.svg')"
                  :style="{ height: '14px', width: '14px' }"
                />
              </el-tooltip>
            </el-button>
            <el-button class="img-btn" @click="handelDel(-1)">
              <el-tooltip class="item" content="删除" effect="dark" placement="top">
                <img
                  :src="require('@/assets/alert_images/abnormal_landing_alarm/del.svg')"
                  :style="{ height: '14px', width: '14px' }"
                />
              </el-tooltip>
            </el-button>
            <el-button style="background-color: #6655e7; color: #fff" type="primary" @click="handleUpdateStatus(0)">
              停用
            </el-button>
            <el-button style="background-color: #6655e7; color: #fff" type="primary" @click="handleUpdateStatus(1)">
              启用
            </el-button>
          </div>
        </div>
      </template>
    </el-popover>
    <IntelligenceCenterWhiteUpload
      v-model:visible="ruleUploadVisible"
      :download-temp="AbnormalLandingAlarmExportTemplateApi"
      :import-assets-fnc="AbnormalLandingAlarmImportApi"
      :title="'规则导入'"
    />
  </div>
</template>
<style lang="scss">
  .abnormal-login-popper-in-table {
    background: #ffffff;
    box-shadow: 0px 0px 4px 4px rgba(77, 81, 86, 0.05);
    border-radius: 8px;
    border: 1px solid #f1f0ff;
    padding: 0px !important;
    .abnormal-login-popper-in-table-tip {
      color: #ff3a3a;
      background: #fff0f1;
    }
    .title-head {
      width: 400px;
      height: 46px;
      padding-left: 10px;
      background: #f7f6ff;
      border-radius: 8px 8px 0px 0px;
      border: 1px solid #f2f0ff;
      display: flex;
      align-items: center;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      word-break: break-all;
      word-wrap: break-word;
      .tag {
        width: 36px;
        height: 22px;
        background: #4abf57;
        border-radius: 4px;
        margin: 0 6px 0 4px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #fff;
      }
      .disablement {
        background: #9690b7;
      }
    }
    .rule-info {
      color: #7d7990;
      .rule-info-ul {
        // height: 116px;
        padding: 10px 15px;
      }

      .rule-info-li {
        display: flex;
        align-items: center;
        margin-bottom: 4px;
      }
      .rule-info-li-sp {
        align-items: baseline;
      }
      .info {
        margin-left: 6px;
        line-height: 27px;
        width: 360px;
        color: #494758;
        // overflow: hidden;
        // text-overflow: ellipsis;
        // white-space: nowrap;
        // word-break: break-all;
        // word-wrap: break-word;
      }
    }
  }
  .abnormal-login-popper {
    padding: 0px !important;
    border-radius: 10px !important;
    .el-button + .el-button {
      margin-left: 8px !important;
    }
  }
</style>

<style scoped lang="scss">
  $criticalColor: #ff4340;
  $lowColor: #45be3e;
  $midColor: #ffb407;
  $highColor: #ff7c06;
  $defaultColor: #9d9aba;
  $criticalBgColor: #d60705;
  $lowBgColor: #f8ae0a;
  $midBgColor: #ff751f;
  $highBgColor: #ff3c3a;
  $defaultBgColor: #9d9aba;
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
  .abnormal-login-container {
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
    .threatName-row {
      max-width: 170px;
      overflow-x: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      word-break: break-all;
      word-wrap: break-word;
    }
    .content-warp {
      height: calc(100vh - 210px);
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
    }
    .charts-warp {
      width: 100%;
      height: 250px;
      border-radius: 8px;
      border: 1px solid #e9e6f9;
    }
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
      .cm-editor {
        height: 34px;
        border-radius: 2.5px;
        .cm-content {
          padding: 2px 0;
        }
      }
    }
    .search-sql {
      width: 100%;
      display: flex;
      align-items: center;
      margin-bottom: 10px;
      .search-value {
        width: calc(100% - 135px);
      }
      .search-tool {
        width: 135px;
        display: flex;
        align-items: center;
      }
    }
    .submit_item {
      .el-form-item__content {
        justify-content: space-between;
      }
    }
    .auto-refresh {
      width: 144px;
      height: 32px;
      padding: 0 12px;
      background: #ffffff;
      border-radius: 2.5px 0px 0px 2.5px;
      border: 1px solid #e4e2f8;
      border-right: 0;
    }
    .auto-refresh-selset {
      :deep() {
        .el-input__wrapper {
          border-radius: 0 2.5px 2.5px 0;
        }
      }
    }
    .my-spacer {
      :deep() {
        > span {
          border-right: 1px solid var(--el-collapse-border-color);
        }
      }
    }
    .title-head {
      display: flex;
      align-items: center;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      word-break: break-all;
      word-wrap: break-word;
      .tag {
        width: 36px;
        height: 22px;
        background: #4abf57;
        border-radius: 4px;
        margin: 0 6px 0 4px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #fff;
      }
      .disablement {
        background: #9690b7;
      }
    }
    .rule-info {
      color: #7d7990;
      .rule-info-ul {
        height: 116px;
        padding: 10px 15px;
      }

      .rule-info-li {
        display: flex;
        align-items: center;
        margin-bottom: 4px;
      }
      .info {
        margin-left: 6px;
        line-height: 20px;
        width: 360px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        word-break: break-all;
        word-wrap: break-word;
      }
    }
    .collapse-items {
      height: calc(100vh - 250px);
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
      :deep() {
        .el-collapse-item__header {
          height: 40px;
        }
      }
    }
    .title {
      margin-block: 0 0.5em;
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
    .alert_tag {
      display: inline-block;
      width: 54px;
      height: 24px;
      line-height: 24px;
      text-align: center;
      background: #9d9aba;
      border-radius: 4px;
      font-size: 13px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
      &.critical {
        background-color: rgba($criticalBgColor, 1);
      }
      &.low {
        background-color: rgba($lowBgColor, 1);
      }
      &.mid {
        background-color: rgba($midBgColor, 1);
      }
      &.high {
        background-color: rgba($highBgColor, 1);
      }
      &.default {
        background-color: rgba($defaultBgColor, 1);
        color: #fff;
      }
    }
    :deep() {
      .my-table {
        margin-top: 15px;
        height: calc(100vh - 470px);
      }
      .el-popper.abnormal-login-popper {
        left: -24px !important;
        padding-inline: 0;
        padding-top: 0;
        .rule-footer {
          height: 60px;
          display: flex;
          justify-content: space-between;
          padding: 0 20px 0 37px;
          border-top: 1px solid var(--el-border-color-lighter);
          .el-button {
            border-radius: 4px !important;
          }
          .left {
            display: flex;
            align-items: center;
          }
          .right {
            display: flex;
            align-items: center;
            .img-btn {
              padding: 8px !important;
            }
          }
          // padding-top: 12px;
        }
        .rule-title {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding-inline: 20px;
          height: 48px;
          border-bottom: 1px solid var(--el-border-color-lighter);
          align-items: center;
          margin-bottom: 14px;
          .el-icon {
            transform: rotate(45deg);
          }
          & > span {
            font-weight: 500;
            font-size: 16px;
            color: #1e1842;
          }
        }
        .rule-title-head {
          height: 54px;
        }
        .rule-list {
          // height: 54px !important;
          padding-inline: 12px;
          .el-collapse {
            border-top: none;
            border-bottom: none;
            margin-inline: 8px;
            .el-collapse-item {
              margin-bottom: 8px;
              ul {
                height: 116px;
                padding: 0 15px;
              }
              .el-space {
                width: 100%;
                justify-content: space-evenly;
                border-top: 1px solid var(--el-collapse-border-color);
                height: 34px;
                line-height: 34px;
                .el-divider--vertical {
                  height: 110%;
                }
                span {
                  cursor: pointer;
                }
              }
              .el-collapse-item__header {
                background: #f9f8ff;
                border: 1px solid var(--el-collapse-border-color);
                border-radius: 6px;
              }
              .el-collapse-item__wrap {
                border: 1px solid var(--el-collapse-border-color);
                border-top: none;
                border-radius: 0 0 6px 6px;
                .el-collapse-item__content {
                  padding-bottom: 0;
                }
              }
              &.is-active {
                .el-collapse-item__header {
                  border-radius: 6px 6px 0 0;
                }
              }
            }
          }
        }
      }
    }
  }
</style>
