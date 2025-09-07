<script lang="ts">
  export default {
    name: 'EmailAnalysis',
  }
</script>

<script setup lang="ts">
  import dayjs, { Dayjs } from 'dayjs'
  import { useUserStore } from '@/store/modules/user'
  import { getMailAnalysisChartApi, getMailAnalysisListApi, getMailFieldsApi, setMailFieldsApi } from '@/api-ecs/alert'
  import { MailAnalysisChartQuery, MailAnalysisItem } from '@/types'
  import { TableColumnItemType } from '/#/store'
  import { debounce } from 'lodash'
  import type { EChartsType } from 'echarts'
  import { formatNstime } from '@/utils/time'
  import { useScroll } from '@vueuse/core'
  import VabChart from '@/plugins/VabChart/index.vue'
  import AttackHighlight from '@/components/attack-highlight.vue'
  import ApplicationConfiguration from '@/ecs/assets/components/application-configuration.vue'
  import AiMailConfiguration from '@/ecs/assets/components/ai-mail-configuration.vue'
  import JsonPreview from '@/components/json-preview.vue'
  import { timeDuratioOptions as date_options } from '@/data/constant'
  type CenterChartDataType<T> = T extends 'key' ? string[] : number[]

  const $baseMessage: any = inject('$baseMessage')
  const userStore = useUserStore()
  const { getTableColumn } = userStore
  const initOptions = {
    renderer: 'svg',
  }
  const infoType = ref('basic-info')
  const magicIcon = require('@/assets/mofabang.svg')
  // Ai邮件内部邮箱配置
  const emailConfigurationVisible = ref(false)
  const showFiledConfig = ref(false)
  const tableColumn = ref<TableColumnItemType[]>([])
  // 获取中间发件人统计分布图表分页数据
  let CenterChartData: {
    key: string[]
    val: number[]
  } = {
    key: [],
    val: [],
  }
  let CenterChartDataIndex = 0
  let CenterChartPageSzie = 10
  let selectedName = ''
  let domainChartRef = undefined as unknown as EChartsType
  let trendsChartRef = undefined as unknown as EChartsType
  const http_log_ref = ref<InstanceType<typeof AttackHighlight> | null>(null)
  const codeMode = ref()
  const infoVisible = ref()
  const currentInfo = ref()
  const currentInfoIndex = ref(-1)
  const currentInfoLoading = ref(false)
  const _searchStr = ref('')
  const isMini = ref(false)
  const disabledChart = ref(true)
  const timeParty = ref('today')
  const customDate = ref<[Dayjs, Dayjs] | undefined>(undefined)

  const showMore = ref(false)

  const mailLables = ref<{
    [key: string]: {
      count: number
      groupName: string
      id: number
      labelName: string
    }[]
  }>()
  const mimiTags = ref<string[]>([])
  const mailQueryData = reactive({
    //列表查询 包含下面图表查询参数
    pageNum: 1,
    pageSize: 50,
    searchStr: '',
    // 图表查询
    indexType: 33,
    topCount: 100,
    startTime: '',
    endTime: '',
    labels: [] as string[],
    labelRelat: 'or' as 'or' | 'and',
    inOut: false,
    inIn: false,
    outIn: false,
    outOut: false,
    senderMail: '',
    domain: '',
    threatLevel: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10] as number[],
  })
  const mailLogs = reactive({
    list: [] as MailAnalysisItem[],
    total: 0,
    scrollId: '',
    chartLoading: false,
    listLoading: false,
    threatLevelCount: {} as {
      [key: number]: number
    },
  })
  const tableScrollData = ref()
  const arrivedStatus = computed(() => tableScrollData.value?.arrivedState)
  const leftChartOption = reactive({
    animationDuration: 3000,
    title: {
      top: 2,
      left: 10,
      textStyle: {
        color: '#303133',
        fontWeight: 500,
        fontSize: 14,
      },
      text: '邮件趋势',
    },
    legend: {
      data: ['内对内', '内对外', '外对内', '外对外'],
      left: 'center',
      bottom: 5,
      itemHeight: 10,
      itemGap: 25,
      itemWidth: 18,
      selectedMode: 'multiple',
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: {
          backgroundColor: '#6a7985',
        },
      },
      confine: true,
    },
    grid: {
      top: 40,
      left: '3%',
      right: '14%',
      bottom: 25,
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: [] as string[],
      overflow: 'truncate',
      axisLabel: {
        ellipsis: '...',
      },
      splitLine: { lineStyle: { color: ['#CEEDFF'], type: [5, 8], dashOffset: 3 } },
    },
    yAxis: [
      {
        type: 'value',
        splitLine: {
          show: false,
        },
      },
    ],
    series: [] as any,
  })
  const centerChartOption = reactive({
    title: {
      text: '发件人统计分布',
      textStyle: {
        color: '#303133',
        fontWeight: 500,
        fontSize: 14,
      },
      top: 2,
      left: 10,
    },
    grid: {
      top: 30,
      left: '5%',
      bottom: 28,
      containLabel: true,
      show: false,
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow',
      },
    },
    xAxis: {
      type: 'value',
      boundaryGap: false,
      offset: -5,
      splitLine: { lineStyle: { color: ['#CEEDFF'], type: [5, 8], dashOffset: 3 } },
    },
    yAxis: {
      type: 'category',
      data: [] as string[],
      inverse: true,
      axisLabel: {
        width: 80,
        overflow: 'truncate',
        ellipsis: '...',
      },
      splitLine: {
        show: false,
      },
    },
    series: {
      name: '发件人统计分布',
      type: 'bar',
      data: [] as number[],
    },
  })
  const rightChartOption = reactive({
    title: {
      text: '邮件域名分布',
      textStyle: {
        color: '#303133',
        fontWeight: 500,
        fontSize: 14,
      },
      top: 2,
      left: 10,
    },
    tooltip: {
      trigger: 'item',
      formatter: function (val: { value: number }) {
        return val
      },
      textStyle: {
        fontSize: 14,
      },
    },
    legend: {
      type: 'scroll',
      orient: 'vertical',
      height: 190,
      icon: 'circle',
      top: 'center',
      left: '60%',
      itemGap: 11,
      itemHeight: 10,
      itemWidth: 10,
      textStyle: {
        ellipsis: '...',
        width: 120,
        overflow: 'truncate',
      },
      itemStyle: {
        opacity: 1,
      },
      selectorPosition: 'end',
      selectorLabel: {
        color: '#7667ea',
        fontSize: 10,
        borderRadius: 2,
        borderColor: '#EEE',
        rich: {
          fontSize: 6,
        },
      },
    },
    series: {
      name: '邮件域名分布',
      type: 'pie',
      radius: ['50%', '80%'],
      center: ['30%', '55%'],
      selectedMode: 'single',
      padAngle: 5,
      itemStyle: {
        borderRadius: 8,
        borderColor: '#fff',
        borderWidth: 2,
      },
      label: {
        show: false,
        position: 'center',
      },
      emphasis: {
        label: {
          show: true,
          fontSize: 16,
        },
      },
      labelLine: {
        show: false,
      },
      data: [] as { value: number; name: string }[],
    },
  })
  const handleMailReset = async (isClear: boolean) => {
    if (!isClear) {
      Object.assign(mailQueryData, {
        searchStr: _searchStr.value,
      })
    } else {
      selectedName = ''
      CenterChartDataIndex = 0
      mimiTags.value = []
      const nowDate = dayjs()
      Object.assign(mailQueryData, {
        //列表查询 包含下面图表查询参数
        pageNum: 1,
        pageSize: 20,
        searchStr: '',
        // 图表查询
        indexType: 33,
        topCount: 100,
        startTime: nowDate.subtract(1, 'hour').format('YYYY-MM-DD HH:mm:ss'),
        endTime: nowDate.format('YYYY-MM-DD HH:mm:ss'),
        labels: [],
        labelRelat: 'or',
        inOut: false,
        inIn: false,
        outIn: false,
        senderMail: '',
        domain: '',
        threatLevel: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
      })
      timeParty.value = 'today'
    }
  }
  const hanldeGetCenterChartData = (type: 'key' | 'val') => {
    const _values = CenterChartData[type]
    return _values.slice(CenterChartDataIndex, CenterChartDataIndex + CenterChartPageSzie)
  }
  const handleCenterPageNumberChange = (pageNum: number) => {
    CenterChartDataIndex = (pageNum - 1) * CenterChartPageSzie
    centerChartOption.yAxis.data = hanldeGetCenterChartData('key') as CenterChartDataType<'key'>
    centerChartOption.series.data = hanldeGetCenterChartData('val') as CenterChartDataType<'val'>
  }
  const handleGetCharts = debounce(async (charts: MailAnalysisChartQuery) => {
    const { data } = await getMailAnalysisChartApi(charts)
    const { inIn = null, inOut = null, outIn = null, outOut = null } = data.line
    // 邮件趋势
    const leftChartSeries = {
      name: '',
      type: 'line',
      showSymbol: false,
      data: [],
      endLabel: {
        show: true,
        formatter: (params: any) => `${params.seriesName}:${params.data}`,
      },
    }
    // 发件人统计分布
    const leftChart = {
      xAxis: inIn?.date || [],
      series: [
        { ...leftChartSeries, name: '内对内', data: inIn?.count || [] },
        { ...leftChartSeries, name: '内对外', data: inOut?.count || [] },
        { ...leftChartSeries, name: '外对内', data: outIn?.count || [] },
        { ...leftChartSeries, name: '外对外', data: outOut?.count || [] },
      ],
    }
    CenterChartData = Object.entries(data.senderBar).reduce(
      (pre, next) => {
        const [name, value] = next
        pre.key.push(name)
        pre.val.push(value)
        return pre
      },
      {
        key: [] as string[],
        val: [] as number[],
      }
    )
    //邮件域名分布
    const rightChart = Object.entries(data.domainPie).map(([name, value]) => ({ value, name }))

    leftChartOption.series = leftChart.series
    leftChartOption.xAxis.data = leftChart.xAxis
    // 中间发件人统计分布
    handleCenterPageNumberChange(1)

    rightChartOption.series.data = rightChart

    mailLables.value = data.label || {}
    mailLogs.chartLoading = false
    mailLogs.threatLevelCount = data.threatLevel
  }, 800)
  const handleGetList = debounce(
    async (
      listData: MailAnalysisChartQuery & {
        pageNum: number
        pageSize: number
        searchStr: string
      }
    ) => {
      const { data } = await getMailAnalysisListApi({
        ...listData,
        scrollId: mailLogs.scrollId,
        searchStr: _searchStr.value,
      })
      const arr = data.resList || []
      mailLogs.scrollId ? mailLogs.list.push(...arr) : (mailLogs.list = arr)
      mailLogs.total = data.total || 0
      mailLogs.scrollId = data.scrollId || ''
      mailLogs.listLoading = false
    },
    800
  )
  const handleFullTagClick = (tag: string) => {
    const tagIndex = mimiTags.value.findIndex((i) => i === tag)
    if (tagIndex > -1) {
      mimiTags.value.splice(tagIndex, 1)
    } else {
      mimiTags.value.push(tag)
    }
  }
  const handleMiniTagClick = (tag: string) => {
    const tagIndex = mailQueryData.labels.findIndex((i) => i === tag)
    if (tagIndex > -1) {
      mailQueryData.labels.splice(tagIndex, 1)
    } else {
      mailQueryData.labels.push(tag)
    }
  }
  const tableColumnFormat = (row: any, column: any, cellValue: any) => {
    if (column.property !== 'startTimeNs') return cellValue
    return formatNstime(cellValue)
  }
  const leftChartClick = (data: any, chart: EChartsType) => {
    if (!chart) return
    trendsChartRef = chart
    const map = {
      内对内: 'inIn',
      内对外: 'inOut',
      外对内: 'outIn',
      外对外: 'outOut',
    }
    const affectLable = data.name as keyof typeof map
    const affect = map[affectLable]
    const toUnSelect = affectLable === selectedName

    chart.setOption({
      legend: {
        selectedMode: toUnSelect ? 'multiple' : 'single',
      },
    })
    chart.dispatchAction({
      type: toUnSelect ? 'legendAllSelect' : 'legendSelect',
      name: !toUnSelect && affectLable,
    })
    Object.assign(
      mailQueryData,
      toUnSelect
        ? { inIn: false, inOut: false, outIn: false, outOut: false }
        : { inIn: affect === 'inIn', inOut: affect === 'inOut', outIn: affect === 'outIn', outOut: affect === 'outOut' }
    )
    selectedName = affectLable
  }
  const centerChartClick = (data: any) => {
    mailQueryData.senderMail = mailQueryData.senderMail === data?.name ? '' : data?.name
  }
  const rightChartSelectchanged = (data: any, chart: EChartsType) => {
    const {
      fromActionPayload: { dataIndexInside },
      fromAction,
    } = data
    if (dataIndexInside === undefined) {
      chart.setOption(
        {
          series: {
            itemStyle: {
              opacity: 1,
            },
          },
        },
        true
      )
      return
    }
    const domain = rightChartOption.series.data[dataIndexInside].name
    chart.setOption({
      series: {
        itemStyle: {
          opacity: fromAction === 'unselect' ? 1 : 0.5,
        },
      },
    })
    domainChartRef = chart
    mailQueryData.domain = mailQueryData.domain === domain ? '' : domain
  }
  const rightChartLegendselect = (data: any, chart: EChartsType) => {
    const domain = data.name
    const toUnSelect = mailQueryData.domain === domain
    domainChartRef = chart
    chart.dispatchAction({
      type: 'legendSelect',
      name: domain,
    })
    mailQueryData.domain = toUnSelect ? '' : domain
  }
  const handleChange = (value: any) => {
    http_log_ref.value?.setConvertsType(value)
  }
  const handleShowInfo = (row: MailAnalysisItem, index: number) => {
    currentInfoIndex.value = index
    infoVisible.value = true
  }
  const handleKeyDown = (e: any) => {
    if (e.keyCode == 37 && currentInfoIndex.value > 0) {
      handlePrev()
    } else if (e.keyCode == 39 && currentInfoIndex.value < mailLogs.total) {
      handleNext()
    }
  }
  const handlePrev = () => {
    currentInfoIndex.value--
  }
  const handleNext = () => {
    currentInfoIndex.value++
  }

  async function configurationHandel(showField: TableColumnItemType[]) {
    const { msg } = await setMailFieldsApi(showField.map((i) => i.id))
    tableColumn.value = showField
    $baseMessage(msg, 'success', 'vab-hey-message-success')
  }
  const formatDate = () => {
    if (timeParty.value === 'user-defined') {
      if (customDate.value) {
        const [startDate, endDate] = customDate.value
        mailQueryData.endTime = dayjs(endDate).format('YYYY-MM-DD HH:mm:ss')
        mailQueryData.startTime = dayjs(startDate).format('YYYY-MM-DD HH:mm:ss')
      }
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
    ].includes(timeParty.value)
      ? timeDate.add(10, 'minute').format('YYYY-MM-DD HH:mm:ss')
      : timeDate.endOf('day').format('YYYY-MM-DD HH:mm:ss')
    const cur = timeDate.startOf('day')
    let startDate = null
    switch (timeParty.value) {
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
    mailQueryData.endTime =
      timeParty.value === 'yesterday' ? timeDate.endOf('day').subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss') : endDate
    mailQueryData.startTime = startDate
  }
  // 表格标签转换
  const tableLabel = (list: string[]) => {
    const selectedTags = list.filter((item) => mailQueryData.labels.includes(item))
    const otherTags = list.filter((item) => !selectedTags.includes(item))
    return [...selectedTags, ...otherTags].slice(0, 5)
  }
  const handleThreatLevelChange = (threatLevel: number) => {
    // length 11 （0-10）
    if (mailQueryData.threatLevel.length === 11) {
      mailQueryData.threatLevel = [threatLevel]
      return
    }
    if (!mailQueryData.threatLevel.includes(threatLevel)) {
      mailQueryData.threatLevel.push(threatLevel)
      return
    }
    mailQueryData.threatLevel = mailQueryData.threatLevel.filter((i) => i !== threatLevel)
  }
  const setThreatLevel = (isFull = false) => {
    mailQueryData.threatLevel = isFull
      ? [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
      : [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10].filter((i) => !mailQueryData.threatLevel.includes(i))
  }
  watch(
    () => [timeParty.value, customDate.value],
    () => {
      formatDate()
    },
    {
      immediate: true,
    }
  )

  watch(
    () => mimiTags.value,
    () => {
      mailQueryData.labels = [...mimiTags.value]
    },
    {
      deep: true,
    }
  )
  watch(
    () => arrivedStatus.value,
    () => {
      if (arrivedStatus.value.bottom && mailLogs.list.length < mailLogs.total) {
        mailLogs.listLoading = true
        handleGetList(mailQueryData)
      }
    },
    {
      deep: true,
    }
  )
  watch(
    () => mailQueryData,
    () => {
      const { pageNum, pageSize, ...charts } = mailQueryData
      mailLogs.chartLoading = true
      mailLogs.listLoading = true
      mailLogs.scrollId = ''
      domainChartRef?.dispatchAction({
        type: 'unselect',
      })
      handleGetCharts(charts)
      handleGetList(mailQueryData)
      if (tableScrollData.value?.y) {
        tableScrollData.value.y = 0
      }
    },
    {
      deep: true,
      immediate: true,
    }
  )
  watch(
    () => currentInfoIndex.value,
    () => {
      if (currentInfoIndex.value < 0) return
      handleChange('default')
      codeMode.value = ''
      currentInfoLoading.value = true
      currentInfo.value = mailLogs.list[currentInfoIndex.value]
      if (currentInfoIndex.value + 8 > mailLogs.list.length && mailLogs.total > mailLogs.list.length) {
        handleGetList(mailQueryData)
      }
      setTimeout(() => {
        currentInfoLoading.value = false
      }, 500)
    }
  )
  watch(
    () => infoVisible.value,
    () => {
      infoType.value = 'basic-info'
      if (!infoVisible.value) {
        currentInfo.value = undefined
        currentInfoIndex.value = -1
        window.removeEventListener('keyup', handleKeyDown, false)
      }
    }
  )
  onMounted(async () => {
    window.addEventListener('keyup', handleKeyDown)
    const allFields = getTableColumn(33)
    const { data: showFields } = await getMailFieldsApi()
    tableColumn.value = showFields
      .map((item) => allFields.find((fields) => fields.id === item))
      .filter(Boolean) as TableColumnItemType[]

    const dom = document.querySelector('.mail-log-table .el-scrollbar__wrap') as HTMLDivElement
    const _tableScrollData = useScroll(dom)
    tableScrollData.value = _tableScrollData
  })
</script>

<template>
  <div class="email-analysis-container">
    <div class="header">
      <vab-query-form>
        <vab-query-form-left-panel :span="4">
          <h3>
            AI邮件分析
            <div class="email-analysis-total">邮件总数：{{ mailLogs.total }}</div>
          </h3>
        </vab-query-form-left-panel>
        <vab-query-form-right-panel :span="20">
          <el-select v-model="timeParty" style="width: 120px; margin-inline: 22px 10px">
            <el-option v-for="item in date_options" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <vab-date-time-picker
            v-if="timeParty === 'user-defined'"
            v-model="customDate"
            style="width: fit-content; margin-right: 10px"
          />
          <el-input v-model="_searchStr" placeholder="模糊检索" style="width: 240px; margin-right: 10px">
            <template #append>
              <el-button @click="handleMailReset(false)">
                <vab-icon icon="search-line" />
                检索
              </el-button>
            </template>
          </el-input>
          <el-button
            @click="
              () => {
                domainChartRef?.dispatchAction({
                  type: 'unselect',
                })
                leftChartClick({ name: selectedName }, trendsChartRef)
                handleMailReset(true)
              }
            "
          >
            重置页面
          </el-button>
          <!-- 显示表头字段配置 -->
          <el-tooltip class="item" content="字段配置" effect="dark" placement="top">
            <div @click="showFiledConfig = true">
              <vab-icon
                icon="edit-2-line item-icon"
                style="
                  font-size: 20px;
                  color: #b3b9c8;
                  margin-left: 20px;
                  margin-right: 5px;
                  vertical-align: -4px !important;
                "
              />
              <span style="cursor: pointer">编辑</span>
            </div>
          </el-tooltip>
          <el-tooltip class="item" content="邮箱配置" effect="dark" placement="top">
            <div @click="emailConfigurationVisible = true">
              <vab-icon
                icon="settings-3-line"
                style="
                  font-size: 20px;
                  color: #b3b9c8;
                  margin-left: 20px;
                  margin-right: 5px;
                  vertical-align: -4px !important;
                "
              />
              <span style="cursor: pointer">设置</span>
            </div>
          </el-tooltip>
        </vab-query-form-right-panel>
      </vab-query-form>
      <div class="mail-tag-box">
        <div v-if="isMini" class="mini-tag">
          <span
            v-for="tag in mimiTags"
            :key="tag"
            class="mini-tag-itme"
            :class="{ selected: mailQueryData.labels.includes(tag) }"
            @click="handleMiniTagClick(tag)"
          >
            {{ tag }}
          </span>
        </div>
        <div v-else class="all-tag">
          <div v-for="(labels, labelGroup) in mailLables" :key="labelGroup" class="tag-list">
            <div class="tag-title">{{ labelGroup }}:</div>
            <div class="tags">
              <span
                v-for="label in labels"
                :key="label.id"
                class="tag-item"
                :class="{ checked: mailQueryData.labels.includes(label.labelName) }"
                @click="handleFullTagClick(label.labelName)"
              >
                {{ label.labelName }}({{ label.count }})
              </span>
            </div>
          </div>
        </div>
        <div class="threat-level-box">
          <h5>
            风险级别
            <!--   -->
            <div class="selectTools">
              <el-checkbox
                :label="`T0 (${mailLogs.threatLevelCount?.[0] || 0})`"
                :model-value="mailQueryData.threatLevel.includes(0)"
                @change="() => handleThreatLevelChange(0)"
              />
              <el-button :auto-insert-space="false" plain size="small" @click="setThreatLevel(true)">全选</el-button>
              <el-button :auto-insert-space="false" plain size="small" @click="setThreatLevel(false)">反选</el-button>
            </div>
          </h5>
          <div class="threat-level">
            <div
              v-for="num in 10"
              :key="num"
              class="threat-level-item"
              :class="{ [`threat-level-T${num}-active`]: mailQueryData.threatLevel.includes(num) }"
              :data-count="mailLogs.threatLevelCount?.[num] || 0"
              @click="() => handleThreatLevelChange(num)"
            >
              T{{ num }}
            </div>
          </div>
        </div>
        <!-- <div v-if="mimiTags.length" class="preview-tag" @click="isMini = !isMini">
          {{ isMini ? '展开' : '收起' }}
        </div> -->
      </div>
      <div v-loading="mailLogs.chartLoading" class="mail-chart-box">
        <div v-if="disabledChart" class="left-chart">
          <vab-chart
            class="mail-trends"
            :init-options="initOptions"
            :legendselectchanged="leftChartClick"
            :option="leftChartOption"
            theme="vab-echarts-theme"
          />
        </div>
        <div v-if="disabledChart" class="right-chart">
          <vab-chart
            class="mail-trends"
            :click="centerChartClick"
            :init-options="initOptions"
            :option="centerChartOption"
            theme="vab-echarts-theme"
          />
          <el-pagination
            class="centerChartPagination"
            hide-on-single-page
            layout="prev, pager, next,"
            :page-size="CenterChartPageSzie"
            size="small"
            :total="CenterChartData.key.length"
            @current-change="handleCenterPageNumberChange"
          />
          <vab-chart
            class="mail-trends"
            :init-options="initOptions"
            :legendselectchanged="rightChartLegendselect"
            :option="rightChartOption"
            :selectchanged="rightChartSelectchanged"
            theme="vab-echarts-theme"
          />
        </div>
      </div>
    </div>
    <div v-loading="mailLogs.listLoading" class="mail-log-box">
      <el-table class="mail-log-table" :data="mailLogs.list">
        <el-table-column align="center" fixed="left" label="序号" type="index" width="75" />
        <el-table-column
          v-for="column in tableColumn"
          :key="column.id"
          :formatter="tableColumnFormat"
          :label="column.fieldNameCn"
          :prop="column.fieldNameEn"
          :show-overflow-tooltip="column.fieldNameCn !== '标签'"
          :width="`${['邮件摘要', '标签'].includes(column.fieldNameCn) ? 600 : 250}`"
        >
          <template v-if="column.fieldNameCn === '标签'" #default="{ row }">
            <template v-for="(tag, $index) in tableLabel(row.emailTags)" :key="$index">
              <span v-if="tag" class="email-analysis-tag">{{ tag }}</span>
            </template>
          </template>
          <template v-else-if="column.fieldNameCn === '风险级别'" #default="{ row }">
            <div class="email-analysis-level-icon">
              <div
                class="email-analysis-level-text"
                :style="{ left: `${(Math.ceil((row.threatLevel || 1) / 2) - 1) * 30}px` }"
              >
                T{{ row.threatLevel }}
              </div>
              <div
                class="email-analysis-level-item"
                :class="{ active: Math.ceil(row.threatLevel / 2) > 0 }"
                style="--analysis-level: #ffb9b8; border-radius: 8px 0 0 8px"
              />
              <div
                class="email-analysis-level-item"
                :class="{ active: Math.ceil(row.threatLevel / 2) > 1 }"
                style="--analysis-level: #ff9698"
              />
              <div
                class="email-analysis-level-item"
                :class="{ active: Math.ceil(row.threatLevel / 2) > 2 }"
                style="--analysis-level: #ff5e62"
              />
              <div
                class="email-analysis-level-item"
                :class="{ active: Math.ceil(row.threatLevel / 2) > 3 }"
                style="--analysis-level: #f31618"
              />
              <div
                class="email-analysis-level-item"
                :class="{ active: Math.ceil(row.threatLevel / 2) > 4 }"
                style="--analysis-level: #c10003"
              />
            </div>
          </template>
          <template v-else-if="column.fieldNameCn === '邮件摘要'" #default="{ row }">
            {{ row.aiEmailSum?.replaceAll('*', '') || '' }}
          </template>
        </el-table-column>
        <el-table-column align="center" fixed="right" label="操作" width="80">
          <template #default="{ row, $index }">
            <el-button size="small" @click="() => handleShowInfo(row, $index)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <div class="email-analysis-detail">
      <el-drawer v-model="infoVisible" destroy-on-close size="85%" title="邮件分析详情">
        <template #default>
          <div v-loading="currentInfoLoading" class="analysis-detail-box">
            <div class="email-analysis-header">
              <dl class="title">
                <dt><img alt="邮箱" :src="require('@/assets/alert_images/email-icon.png')" /></dt>
                <dd>
                  <h4 class="email-analysis-subject">
                    {{ currentInfo?.subject }}
                  </h4>
                  <p>{{ currentInfo?.mailDate || '-' }}</p>
                </dd>
              </dl>
              <div>
                <el-button :disabled="currentInfoIndex === 0" plain size="small" @click="handlePrev">上一封</el-button>
                <el-button :disabled="currentInfoIndex === mailLogs.total - 1" plain size="small" @click="handleNext">
                  下一封
                </el-button>
              </div>
            </div>
            <div class="email-analysis-content">
              <div class="email-analysis-content-left">
                <div class="email-analysis-threatLevel-box">
                  <div class="email-analysis-level">风险等级：T{{ currentInfo.threatLevel }}</div>
                  <div class="email-analysis-level-icon">
                    <div
                      class="email-analysis-level-text"
                      :style="{ left: `${(Math.ceil((currentInfo.threatLevel || 1) / 2) - 1) * 30}px` }"
                    >
                      T{{ currentInfo.threatLevel }}
                    </div>
                    <div
                      class="email-analysis-level-item"
                      :class="{ active: Math.ceil(currentInfo.threatLevel / 2) > 0 }"
                      style="--analysis-level: #ffb9b8; border-radius: 8px 0 0 8px"
                    />
                    <div
                      class="email-analysis-level-item"
                      :class="{ active: Math.ceil(currentInfo.threatLevel / 2) > 1 }"
                      style="--analysis-level: #ff9698"
                    />
                    <div
                      class="email-analysis-level-item"
                      :class="{ active: Math.ceil(currentInfo.threatLevel / 2) > 2 }"
                      style="--analysis-level: #ff5e62"
                    />
                    <div
                      class="email-analysis-level-item"
                      :class="{ active: Math.ceil(currentInfo.threatLevel / 2) > 3 }"
                      style="--analysis-level: #f31618"
                    />
                    <div
                      class="email-analysis-level-item"
                      :class="{ active: Math.ceil(currentInfo.threatLevel / 2) > 4 }"
                      style="--analysis-level: #c10003"
                    />
                  </div>
                </div>
                <ul class="email-info">
                  <li>
                    <span class="email-info-label">发件人：</span>
                    <div class="email-info-value">
                      <el-tooltip
                        :content="`${currentInfo?.senderName}<${currentInfo?.senderMail}>`"
                        effect="dark"
                        placement="top"
                      >
                        <p>{{ currentInfo?.senderName }}&lt;{{ currentInfo?.senderMail }}&gt;</p>
                      </el-tooltip>
                    </div>
                  </li>
                  <li>
                    <span class="email-info-label">收件人：</span>
                    <div class="email-info-value">
                      <el-tooltip
                        v-for="(sender, index) in currentInfo?.receiverName?.slice(0, 5)"
                        :key="index"
                        :content="`${sender}<${currentInfo?.receiverMail[index]}>`"
                        effect="dark"
                        placement="top"
                      >
                        <p>{{ sender }}&lt;{{ currentInfo?.receiverMail[index] }}&gt;</p>
                      </el-tooltip>
                      <template v-if="showMore">
                        <el-tooltip
                          v-for="(sender, index) in currentInfo?.receiverName?.slice(5)"
                          :key="index"
                          :content="`${sender}<${currentInfo?.receiverMail[5 + index]}>`"
                          effect="dark"
                          placement="top"
                        >
                          <p>{{ sender }}&lt;{{ currentInfo?.receiverMail[5 + index] }}&gt;</p>
                        </el-tooltip>
                      </template>
                      <span
                        v-if="currentInfo?.receiverName?.length - 5 > 0"
                        class="show-more"
                        @click="showMore = !showMore"
                      >
                        <template v-if="!showMore">+另外{{ currentInfo?.receiverName?.length - 5 }}人</template>
                        <template v-else>收起</template>
                      </span>
                    </div>
                  </li>
                  <li>
                    <span class="email-info-label">标签信息：</span>
                    <div class="email-info-value email-analysis-tags">
                      <span v-if="currentInfo?.position" class="email-analysis-tag">{{ currentInfo?.position }}</span>
                      <template v-for="tag in currentInfo?.emailTags" :key="tag">
                        <span v-if="tag" class="email-analysis-tag">{{ tag }}</span>
                      </template>
                    </div>
                  </li>
                </ul>
              </div>
              <div class="email-analysis-content-right">
                <el-button-group>
                  <el-button :type="infoType === 'basic-info' ? 'primary' : 'default'" @click="infoType = 'basic-info'">
                    基本信息
                  </el-button>
                  <el-button :type="infoType === 'raw-data' ? 'primary' : 'default'" @click="infoType = 'raw-data'">
                    原始数据
                  </el-button>
                </el-button-group>
                <template v-if="infoType === 'basic-info'">
                  <div class="email-analysis-info-box">
                    <div class="header">邮件摘要</div>
                    <div
                      class="content"
                      style="
                        max-height: 150px;
                        overflow-y: auto;
                        font-size: 15px;
                        line-height: 26px;
                        line-height: 28px;
                        letter-spacing: 0.3px;
                      "
                    >
                      {{ currentInfo?.aiEmailSum?.replaceAll('*', '') || '暂无数据' }}
                    </div>
                  </div>
                  <div class="email-analysis-info-box" style="flex: 1">
                    <div class="header">
                      邮件内容
                      <div class="fomart-code">
                        <el-tooltip content="魔法棒" effect="dark" placement="top" :show-arrow="false">
                          <img alt="魔法棒" :src="magicIcon" @click="handleChange('magic')" />
                        </el-tooltip>
                        <el-select
                          v-model="codeMode"
                          placeholder="请选择编码模式"
                          style="width: 140px; margin-left: 10px; height: 28px"
                          @change="handleChange"
                        >
                          <el-option key="default" label="恢复原始编码" value="default" />
                          <el-option key="GB2312" label="GB2312" value="gb2312" />
                          <el-option key="GBK" label="GBK" value="gbk" />
                        </el-select>
                      </div>
                    </div>
                    <div class="content">
                      <attack-highlight ref="http_log_ref" :attack-http-msg="currentInfo?.body || '暂无内容'" />
                    </div>
                  </div>
                  <div class="email-analysis-info-box">
                    <div class="header">邮件附件</div>
                    <div class="content attachments">
                      <div v-for="(item, index) in currentInfo?.attachment.split(',')" :key="index" class="attachment">
                        <vab-icon icon="attachment-2" />
                        {{ item }}
                      </div>
                    </div>
                  </div>
                </template>
                <div v-else class="email-analysis-info-box">
                  <div class="header">原始数据</div>
                  <div class="content source_data">
                    <json-preview :json-value="JSON.stringify(currentInfo, null, 4)" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </template>
      </el-drawer>
    </div>
    <!-- 表头字段配置 -->
    <application-configuration
      v-if="showFiledConfig"
      v-model="showFiledConfig"
      :fields="tableColumn"
      :retrieve-index-type="33"
      @handleok="configurationHandel"
    />
    <ai-mail-configuration v-model="emailConfigurationVisible" />
  </div>
</template>

<style scoped lang="scss">
  .email-analysis-container {
    display: flex;
    flex-direction: column;
    overflow-y: auto;
    h3 {
      margin-block: 0;
    }
    .header {
      position: sticky;
      top: 0;
    }
    .email-analysis-total {
      display: inline-block;
      margin-left: 10px;
      padding: 0 11px;
      height: 30px;
      line-height: 30px;
      background: #f4f3fa;
      border-radius: 4px;
      font-weight: 400;
      font-size: 14px;
      color: #605b7d;
    }
    .mail-tag-box {
      margin-top: 12px;
      border: var(--el-border);
      border-radius: 2px;
      padding: 10px 14px 0;
      position: relative;
      display: flex;
      .threat-level-box {
        flex: 1;
        overflow: hidden;
        h5 {
          font-weight: 500;
          font-size: 14px;
          color: #1e1842;
          text-indent: 19px;
          margin-bottom: 20px;
          display: flex;
          align-items: center;
          justify-content: space-between;
          .selectTools {
            display: flex;
            align-items: center;
            .el-checkbox {
              text-indent: 0;
              margin-right: 12px;
            }
            .el-button {
              width: 42px;
              height: 24px;
              background: #ffffff;
              border-radius: 4px;
              & + .el-button {
                margin-left: 4px;
              }
            }
          }
        }
        .threat-level {
          width: 75%;
          margin: 0 auto;
          display: flex;
          flex-wrap: wrap-reverse;
          flex-direction: row-reverse;
          .threat-level-item {
            flex-basis: 20%;
            font-size: 14px;
            color: #4c4242;
            font-weight: 500;
            text-align: center;
            margin-bottom: 35px;
            height: 42px;
            line-height: 32px;
            border-radius: 2px;
            border-right: 2px solid #ffffff;
            background: #ded6d6;
            cursor: pointer;
            position: relative;
            &::after {
              content: attr(data-count);
              position: absolute;
              inset: 0 0 10px 0;
              text-align: center;
              transform: translateY(42px);
              color: #4a4759;
            }
            $colors: (#f5d2d2, #ffbcbc, #ff9494, #ff7676, #ff5151, #ff3535, #ff1516, #ea0201, #de0200, #ca0000);
            @each $c in $colors {
              $i: index($colors, $c);
              &.threat-level-T#{$i}-active {
                background-color: $c;
                color: #fff;
              }
            }
          }
        }
      }
      .all-tag {
        // margin-right: 30px;
        flex: 2;
        overflow: hidden;
        .tag-list {
          line-height: 26px;
          display: flex;
          width: 100%;
          .tag-title {
            width: 70px;
            color: #a9acb3;
          }
          .tags {
            flex: 1;
            word-break: keep-all;
            .tag-item {
              color: #303133;
              cursor: pointer;
              margin-right: 10px;
              &:hover,
              &.checked {
                color: var(--el-color-primary);
                font-weight: 600;
              }
            }
          }
        }
      }
      .mini-tag {
        margin-right: 30px;
        .mini-tag-itme {
          display: inline-block;
          height: 32px;
          padding-inline: 20px;
          margin-right: 12px;
          margin-bottom: 12px;
          font-family: PingFangSC, PingFang SC;
          font-size: 14px;
          font-weight: 400;
          line-height: 32px;
          color: var(--el-button-text-color);
          cursor: pointer;
          user-select: none;
          background: #f3f3fe;
          border-radius: 4px;
          &.selected {
            color: #fff;
            background-color: var(--el-color-primary);
          }
        }
      }
      .preview-tag {
        position: absolute;
        right: 5px;
        top: 5px;
        cursor: pointer;
        background-color: #726aa3;
        border-radius: 2px;
        color: #fff;
        padding-inline: 6px;
      }
    }
    .mail-chart-box {
      border: var(--el-border);
      border-radius: 2px;
      display: flex;
      height: 240px;
      margin-block: 20px;
      overflow: hidden;
      &:empty {
        border: none;
        height: 0;
        margin-bottom: 0;
      }
      .left-chart {
        flex: 1;
        position: relative;
        overflow: hidden;
      }
      .right-chart {
        flex: 2;
        display: flex;
        position: relative;
        overflow: hidden;
        .echarts {
          width: 50%;
        }
        .centerChartPagination {
          position: absolute;
          bottom: 0;
          left: 25%;
          transform: translateX(-50%);
          :deep(.el-input__inner) {
            --el-input-inner-height: 20px;
          }
        }
      }
      .echarts {
        width: 100%;
        height: 100%;
      }
    }
    .mail-log-box {
      flex: 1;
      overflow: auto;
    }
    .email-analysis-level-icon {
      display: flex;
      width: 150px;
      border: 1px solid #fff;
      margin-left: 10px;
      border-radius: 8px;
      margin-top: 20px;
      position: relative;
      .email-analysis-level-text {
        width: 30px;
        height: 16px;
        background: #3c2424;
        border-radius: 2px;
        position: absolute;
        top: -21px;
        left: 2px;
        font-weight: 500;
        font-size: 12px;
        color: #ffffff;
        line-height: 16px;
        text-align: center;
        display: block;
        &::after {
          content: '';
          position: absolute;
          bottom: -8px;
          left: 50%;
          transform: translateX(-50%);
          width: 0;
          height: 0;
          border: 4px solid transparent;
          border-top-color: #3c2424;
        }
      }
      .email-analysis-level-item {
        flex: 1;
        height: 16px;
        background-color: #e4e3e7;
        &:not(:last-child) {
          border-right: 2px solid #fff;
        }
        &:last-child {
          border-radius: 0 8px 8px 0;
        }
        &.active {
          background-color: var(--analysis-level);
        }
      }
    }
    .mt20 {
      margin-top: 20px;
    }
    .email-analysis-detail {
      .analysis-detail-box {
        height: 100%;
        display: flex;
        flex-direction: column;
        .email-analysis-header {
          height: 92px;
          background: #f8f8fd;
          border-radius: 30px 0px 0px 0px;
          padding-inline: 30px;
          display: flex;
          align-items: center;
          justify-content: space-between;
          border-bottom: 1px dashed #d0c8ff;
          .title {
            margin: 0;
            flex: 1;
            display: flex;
            align-items: center;
            overflow: hidden;
            dt {
              width: 54px;
              height: 54px;
              display: block;
              float: left;
              margin-right: 12px;
              img {
                width: 100%;
              }
            }
            dd {
              flex: 1;
              overflow: hidden;
              .email-analysis-subject {
                font-weight: 500;
                font-size: 18px;
                color: #28253a;
                margin-bottom: 4px;
                overflow: hidden;
                text-overflow: ellipsis;
                word-break: keep-all;
                white-space: nowrap;
              }
              p {
                font-weight: 400;
                font-size: 14px;
                color: #84809e;
                margin-bottom: 0;
                width: fit-content;
              }
              margin-bottom: 0;
            }
          }
        }
        .email-analysis-content {
          flex: 1;
          overflow: hidden;
          display: flex;
          padding: 20px 30px;
          .email-analysis-content-left {
            width: 360px;
            border-radius: 6px;
            border: 1px solid #eceaf5;
            overflow-y: auto;
            .show-more {
              font-weight: 400;
              font-size: 13px;
              color: #6954f0;
              cursor: pointer;
              vertical-align: middle;
            }
            .email-analysis-threatLevel-box {
              height: 64px;
              background: url(@/assets/alert_images/mail-default-threatLevel.png) no-repeat center;
              border-radius: 6px;
              background-size: contain;
              display: flex;
              align-items: center;
              justify-content: space-around;
              margin-top: -1px;
              &:has(.email-analysis-level-item.active) {
                background: url(@/assets/alert_images/mail-threatLevel.png) no-repeat center;
                background-size: contain;
              }
              .email-analysis-level {
                font-weight: 500;
                font-size: 13px;
                color: #3c2424;
                margin-left: 30px;
              }
            }
            .email-info {
              padding: 16px;
              margin-bottom: 0;
              height: calc(100% - 68px);
              overflow-y: auto;
              li {
                width: 100%;
                display: flex;
                &:not(:last-child) {
                  margin-bottom: 10px;
                }
                .email-info-label {
                  font-weight: 400;
                  font-size: 15px;
                  color: #84809e;
                }
                .email-info-value {
                  flex: 1;
                  overflow: hidden;
                  font-weight: 400;
                  font-size: 14px;
                  color: #342e58;
                  p {
                    margin-bottom: 0;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    word-break: keep-all;
                    white-space: nowrap;
                  }
                }
              }
            }
          }
          .email-analysis-content-right {
            flex: 1;
            padding-left: 20px;
            display: flex;
            flex-direction: column;
            overflow: hidden;
            .email-analysis-info-box {
              border: var(--el-border);
              border-radius: 6px;
              margin-top: 14px;
              overflow: hidden;
              .header {
                font-weight: 500;
                font-size: 16px;
                color: #352e58;
                line-height: 36px;
                border-bottom: 1px solid #eceaf5;
                padding-inline: 16px;
                background-color: #f6f6fe;
              }
              .content {
                padding: 16px;
                height: calc(100% - 37px);
                &.source_data,
                &.attachments {
                  padding: 0;
                }
                &.attachments {
                  overflow-y: auto;
                  max-height: 210px;
                  padding-block: 10px;
                  .attachment {
                    padding-inline: 16px;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    word-break: keep-all;
                    white-space: nowrap;
                  }
                }
              }
              .fomart-code {
                float: right;
                :deep(.el-input__wrapper) {
                  height: 28px;
                }
                img {
                  width: 16px;
                  height: 16px;
                  cursor: pointer;
                  vertical-align: middle;
                  display: inline-block;
                }
              }
            }
          }
        }
      }

      .email-analysis-tags {
        border-radius: 4px;
        .email-analysis-tag {
          border-color: #d9ecff;
          color: var(--el-color-primary);
          background-color: rgba(102, 85, 231, 0.08);
          white-space: normal;
          display: inline-block;
          vertical-align: text-top;
          margin: 0 4px;
          line-height: 28px;
          padding-inline: 11px;
          border-radius: 4px;
          word-break: break-all;
          margin-bottom: 4px;
          word-wrap: break-word;
        }
      }
      .email-log {
        border: 1px solid var(--el-border-color-lighter);
        border-radius: 3px;
        .top {
          padding: 0 20px;
          height: 40px;
          background: #f5f7fa;
          border-bottom: 1px solid var(--el-border-color-lighter);
          display: flex;
          align-items: center;
          justify-content: space-between;
        }
      }
      :deep() {
        .el-descriptions__label {
          width: 120px;
        }
        .diyy.el-descriptions {
          .el-descriptions__content {
            width: calc((85vw - 260px) / 2);
          }
        }
        .attack-highlight {
          background-color: transparent;
          border-radius: 0;
          border: 0;
          height: 100%;
          .cm-scroller {
            // min-height: 300px;
            .cm-content {
              width: 100%;
              white-space: break-spaces;
              word-wrap: break-word;
            }
          }
        }
      }
    }
    :deep() {
      .v-codemirror .cm-editor {
        width: 100% !important;
        height: calc(100vh - 281px);
        .cm-content {
          background-color: #fff;
        }
      }
      .el-input-group__append {
        color: #fff;
        --el-input-border-color: var(--el-color-primary);
        background-color: var(--el-color-primary);
      }
      .el-table.mail-log-table {
        height: 100%;
        .email-analysis-tag {
          border-color: #d9ecff;
          color: var(--el-color-primary);
          background-color: rgba(102, 85, 231, 0.08);
          white-space: normal;
          display: inline-block;
          vertical-align: text-top;
          margin: 0 4px;
          line-height: 28px;
          padding-inline: 11px;
          border-radius: 4px;
          word-break: break-all;
          word-wrap: break-word;
        }
      }
      .el-drawer {
        border-radius: 30px 0px 0px 30px;
        background: #534b89;
        .el-drawer__header {
          padding: 12px 30px 12px;
          margin-bottom: 0;
          .el-drawer__title {
            font-size: 20px;
            font-weight: 500;
            font-size: 20px;
            line-height: 34px;
            color: #fff;
          }
          .el-drawer__close-btn {
            width: 26px;
            height: 26px;
            border-radius: 50%;
            background: #eeedf9;
            display: flex;
            justify-content: center;
            align-items: center;
          }
        }
        .el-drawer__body {
          border-radius: 30px 0px 0px 30px;
          background: #fff;
          padding: 0;
          overflow-y: hidden;
        }
      }
    }
  }
</style>
