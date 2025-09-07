<script lang="ts">
  export default {
    name: 'AssetOverview',
  }
</script>

<script setup lang="ts">
  import { Plus, Close } from '@element-plus/icons-vue'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import AssetManager from './asset-manager.vue'
  import AssetsLabel from './assets-label.vue'
  import ThreatTrends from './charts/threat-trends.vue'
  import SystemDistribution from './charts/system-distribution.vue'
  import ThreatClassification from './charts/threat-classification.vue'
  import RankingOfAsset from './charts/ranking-of-asset.vue'
  import ThreateningIncident from './charts/threatening-incident.vue'
  import { Search } from '@element-plus/icons-vue'
  import { getAssetsPreviewListApi, getAssetsPreviewSummaryApi, pageChartsApi } from '@/api-ecs/assets-preview'
  import {
    PreviewAssetsItem,
    CMDBfingerprintList,
    NetPartitionSums,
    FingerprintItem,
    PreviewAssetsBodyType,
  } from '@/types/index'

  import { formatNstime } from '@/utils/time'
  import _ from 'lodash'
  import dayjs from 'dayjs'
  import { exportAssetNewApi } from '~/src/api-ecs/assets'
  import { downloadFile } from '~/src/utils/download'
  import numberFormatter from '~/src/utils/number'
  // 不许修改
  let LABEL_KEYS = [] as number[]
  const route = useRoute()
  const activeAssetsStatus = ref<number[]>([])

  const isOfflineStatus = ref<any[]>([])

  const activeAssetsSourseStatus = ref<number[]>([])

  const activeAssetsFormStatus = ref<number[]>([])

  const activeAssetsCausalityStatus = ref<number[]>([])

  const activeAssetsConflictStatus = ref<number[]>([])

  const activeAssetscontainShandowAssetStatus = ref<number[]>([])

  const previewInfo = ref()

  const tableData = ref<PreviewAssetsItem[]>([])
  const list_total = ref(0)
  const multipleSelection = ref<PreviewAssetsItem[]>([])

  const showInfo = ref(false)
  const lableRef = ref<InstanceType<typeof AssetsLabel> | null>(null)
  let queryData = reactive<PreviewAssetsBodyType>({
    pageNum: 1,
    pageSize: 20,
    searchStr: '',
    lastTime: '',
    labelIds: [] as number[],
    dataSourceId: [] as number[],
    otherDataSource: false,
    labelRelat: 'or',
    searchIps: [],
    serverIps: [],
    activeAssetType: 0,
    originImport: null, // 导入
    originFlow: null, // 流量识别
    know: null, // 类别
    conflict: null, // 冲突
    netType: null, // 属性
    isOffline: null, // 属性
    containShandowAsset: null, // 属性
  })
  const timeNow = dayjs()
  const timeParty = ref(6)
  const endTime = timeNow.subtract(0, 'hour').format('YYYY-MM-DD HH:mm:ss')
  let startTime = timeNow.subtract(timeParty.value, 'hour').format('YYYY-MM-DD HH:mm:ss')
  // 所有资产标签
  const allAssetLabel = ref<CMDBfingerprintList[]>([])
  // 数据中心
  const allNetPartitionSums = ref<NetPartitionSums>([])
  // 活跃资产数量
  const activityAssetNum = ref(0)
  // 非活跃资产数量
  const notActivityAssetNum = ref(0)
  // 从标签库选取用来筛选的标签
  const showFields = ref<FingerprintItem[]>([])
  // 其他数据中心总数
  const otherNetPartationCount = ref(0)
  // 流量识别数量
  const flowRecognitionNum = ref(0)
  // 资产导入数量
  const assetImportNum = ref(0)
  // 已知资产数量
  const assetKnowNum = ref(0)
  // 未知资产数量
  const assetUnKnowNum = ref(0)
  // 外网资产
  const assetOutSide = ref(0)
  // 内网资产
  const assetInSide = ref(0)
  // 冲突资产
  const assetConflict = ref(0)
  // 不冲突资产
  const assetNoConflict = ref(0)
  // 存在影子资产
  const assetShandow = ref(0)
  // 不存在影子资产
  const assetNoShandow = ref(0)
  // 离线资产
  const offlineNum = ref(0)
  // 非离线资产
  const notOfflineNum = ref(0)

  const searchStr = ref()

  const rankingOfAssetRef = ref()

  const isShowChart = ref(true) // 显示图标表
  const isTagContentFold = ref(true)
  const tableNode = ref()
  let divHeight: any = undefined

  // 打开详情
  const handleClick = (row: PreviewAssetsItem) => {
    previewInfo.value = row
    showInfo.value = true
  }
  // 表格标签转换
  const tableLabel = (list: CMDBfingerprintList[]) => {
    const curList: Omit<FingerprintItem, 'groupId' | 'isNew'>[] = list.reduce(
      (_labels: Omit<FingerprintItem, 'groupId' | 'isNew'>[], item: CMDBfingerprintList) => {
        const curLabels = item.labelList.map(({ id, labelName }) => ({ id, labelName }))
        return [..._labels, ...curLabels]
      },
      []
    )
    const selectLabels = curList.filter(({ id }) => LABEL_KEYS.includes(id)).map((i) => i.labelName)
    const otherLabels = curList.filter(({ id }) => !LABEL_KEYS.includes(id)).map((i) => i.labelName)
    return [...selectLabels, ...otherLabels].slice(0, 5)
  }
  //打开指纹标签选择
  const openLableDialog = () => {
    lableRef.value?.showDialog(allAssetLabel.value, queryData.labelRelat as 'or' | 'and')
  }
  // 标签选择器
  const handleLableConfirm = ({ labels, labelRelat }: { labels: FingerprintItem[]; labelRelat: 'or' | 'and' }) => {
    const newSelect = labels.map((i) => i.id).filter((id) => queryData.labelIds!.includes(id))
    const hasChange = !_.isEqual(newSelect, queryData.labelIds) || queryData.labelRelat !== labelRelat
    showFields.value = labels
    queryData.labelIds = newSelect
    queryData.labelRelat = labelRelat
    hasChange && handleSearch()
  }
  // 多选项改变
  const setSelectRows = (val: PreviewAssetsItem[]) => {
    multipleSelection.value = val
  }

  const isValidIPv4AndSubnetMask = (ipv4AndSubnetMask: any) => {
    if (!ipv4AndSubnetMask) return false
    const arr = ipv4AndSubnetMask.split(',')
    // 匹配IPv4地址
    const ipv4Regex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
    //  // 匹配子网掩码
    const regex = /^(?:\d{1,2}|\d\d|\d{3})?(?:\.\d{1,2}){3}$/
    const flag = arr.every((item: any) => {
      return ipv4Regex.test(item) || regex.test(item)
    })
    return flag
  }

  // 检索
  const handleSearch = _.debounce(async () => {
    const match = isValidIPv4AndSubnetMask(searchStr.value)
    if (match) {
      queryData.searchStr = ''
      queryData.searchIps = []
      queryData.searchIps = searchStr.value.split(',')
    } else {
      queryData.searchStr = searchStr.value
      queryData.searchIps = []
    }
    const originImport = null
    const originFlow = null
    const {
      data: { records, total },
    } = await getAssetsPreviewListApi({
      ...queryData,
      dataSourceId: queryData.dataSourceId!.filter(Boolean) as number[],
      activeAssetType: [0, 2].includes(activeAssetsStatus.value.length) ? 0 : +activeAssetsStatus.value.toString(),
      originFlow: activeAssetsSourseStatus.value.length == 1 && activeAssetsSourseStatus.value.includes(0) ? 1 : null,
      originImport: activeAssetsSourseStatus.value.length == 1 && activeAssetsSourseStatus.value.includes(1) ? 1 : null,
      know: [0, 2].includes(activeAssetsFormStatus.value.length)
        ? null
        : (+activeAssetsFormStatus.value.toString() as 0 | 1),
      isOffline: [0, 2].includes(isOfflineStatus.value.length)
        ? null
        : +isOfflineStatus.value.toString() == 1
        ? true
        : false,
      netType: [0, 2].includes(activeAssetsCausalityStatus.value.length)
        ? null
        : (+activeAssetsCausalityStatus.value.toString() as 0 | 1),
      conflict: [0, 2].includes(activeAssetsConflictStatus.value.length)
        ? null
        : (+activeAssetsConflictStatus.value.toString() as 0 | 1),
      containShandowAsset: [0, 2].includes(activeAssetscontainShandowAssetStatus.value.length)
        ? null
        : (+activeAssetscontainShandowAssetStatus.value.toString() as 0 | 1),
    })
    LABEL_KEYS = queryData.labelIds!
    tableData.value = records || []
    list_total.value = total || 0
    getPageCharts()
  }, 1000)
  // 获取数据中心和所有指纹标签
  const handleGetAssetsPreviewSummary = async () => {
    const { data } = await getAssetsPreviewSummaryApi()
    allAssetLabel.value = data?.labelVoList || []
    allNetPartitionSums.value = data?.netPartitionSums || []
    notActivityAssetNum.value = data?.notActivityAssetNum || 0
    activityAssetNum.value = data?.activityAssetNum || 0
    otherNetPartationCount.value = data?.otherNetPartationCount || 0
    flowRecognitionNum.value = data?.originFLowNum || 0
    assetImportNum.value = data?.originImportNum || 0
    assetKnowNum.value = data?.assetKnowNum || 0
    assetUnKnowNum.value = data?.assetNotKnowNum || 0
    assetOutSide.value = data?.outNetNum || 0
    assetInSide.value = data?.inNetNum || 0
    showFields.value = data?.topList || []
    assetConflict.value = data?.conflictNum || 0
    assetNoConflict.value = data?.notConflictNum || 0
    assetShandow.value = data?.shandowAsset || 0
    assetNoShandow.value = data?.notShandowAsset || 0
    offlineNum.value = data?.offlineNum || 0
    notOfflineNum.value = data?.notOfflineNum || 0
  }
  // 数据中心改变
  const handleNetPartitionChange = (sourceId: number, index: number) => {
    queryData.otherDataSource = false
    const cur = queryData.dataSourceId![index]
    // @ts-ignore
    queryData.dataSourceId[index] = cur ? null : sourceId
    handleSearch()
  }

  // 其他数据中心改变
  const handleOtherNetPartitionChange = () => {
    queryData.dataSourceId = []
    queryData.otherDataSource = !queryData.otherDataSource
    handleSearch()
  }

  /**
   * @description: 资产状态变化
   * @param {number} status 1-活跃 2-非活跃
   */
  const handleAssetsActiveChange = (status: number) => {
    activeAssetsStatus.value = activeAssetsStatus.value.includes(status)
      ? activeAssetsStatus.value.filter((item) => item !== status)
      : [...activeAssetsStatus.value, status]
    handleSearch()
  }

  /**
   * @description: 资产离线排查
   * @param {number} status 1-离线 2-非活跃
   */
  const handleisOfflineChange = (status: number) => {
    isOfflineStatus.value = isOfflineStatus.value.includes(status)
      ? isOfflineStatus.value.filter((item) => item !== status)
      : [...isOfflineStatus.value, status]
    handleSearch()
  }

  /**
   * @description: 资产来源变化
   */
  const handleAssetsSourseActiveChange = (status: number) => {
    activeAssetsSourseStatus.value = activeAssetsSourseStatus.value.includes(status)
      ? activeAssetsSourseStatus.value.filter((item) => item !== status)
      : [...activeAssetsSourseStatus.value, status]
    handleSearch()
  }

  /**
   * @description: 资产类别变化
   */
  const handleAssetsFormActiveChange = (status: number) => {
    activeAssetsFormStatus.value = activeAssetsFormStatus.value.includes(status)
      ? activeAssetsFormStatus.value.filter((item) => item !== status)
      : [...activeAssetsFormStatus.value, status]
    handleSearch()
  }

  /**
   * @description: 资产属性变化
   */
  const handleAssetsCausalityActiveChange = (status: number) => {
    activeAssetsCausalityStatus.value = activeAssetsCausalityStatus.value.includes(status)
      ? activeAssetsCausalityStatus.value.filter((item) => item !== status)
      : [...activeAssetsCausalityStatus.value, status]
    handleSearch()
  }
  /**
   * @description: 资产冲突变化
   */
  const handleAssetsconflictChange = (status: number) => {
    activeAssetsConflictStatus.value = activeAssetsConflictStatus.value.includes(status)
      ? activeAssetsConflictStatus.value.filter((item) => item !== status)
      : [...activeAssetsConflictStatus.value, status]
    handleSearch()
  }
  /**
   * @description: 影子资产变化
   */
  const handleAssetscontainShandowChange = (status: number) => {
    activeAssetscontainShandowAssetStatus.value = activeAssetscontainShandowAssetStatus.value.includes(status)
      ? activeAssetscontainShandowAssetStatus.value.filter((item) => item !== status)
      : [...activeAssetscontainShandowAssetStatus.value, status]
    handleSearch()
  }
  //
  // 指纹标签删除
  const handleTagClose = (id: number) => {
    const index = showFields.value.findIndex((fidle) => fidle.id === id)
    const selectIndex = queryData.labelIds!.findIndex((fidleId) => fidleId === id)
    if (index > -1) showFields.value.splice(index, 1)
    if (selectIndex > -1) queryData.labelIds!.splice(selectIndex, 1)
    // if (LABEL_Index > -1) LABEL_KEYS.splice(selectIndex, 1)
    handleSearch()
  }
  // 指纹标签选中
  const handleFieldChange = (id: number) => {
    const index = queryData.labelIds!.findIndex((fidleId) => fidleId === id)
    if (index > -1) {
      queryData.labelIds!.splice(index, 1)
    } else {
      queryData.labelIds!.push(id)
    }
    handleSearch()
  }

  const labelRelatChange = (val: string) => {
    queryData.labelRelat = val
    handleSearch()
  }
  onMounted(async () => {
    activeAssetsFormStatus.value = route.query.params === 'unKnown' ? [0] : []
    handleGetAssetsPreviewSummary()
    await handleSearch()
    tableNode.value = document.querySelector('.assets-table')
    divHeight = tableNode.value.offsetHeight
    handleUnFold()
  })

  // 表格数据
  const threatType = ref()
  const sysType = ref()
  const threatLevel = ref()
  const victimIp = ref()
  const threatName = ref()
  const pageChartsParams = reactive<{ threatLevelCode?: number; threatName?: string }>({
    threatLevelCode: undefined,
    threatName: undefined,
  })
  const getPageCharts = async () => {
    const { data } = await pageChartsApi({
      startTime,
      endTime,
      otherDataSource: queryData.otherDataSource,
      activeAssetType: queryData.activeAssetType,
      dataSourceId: queryData.dataSourceId,
      labelIds: queryData.labelIds,
      labelRelat: queryData.labelRelat,
      searchIps: queryData.searchIps,
      searchStr: queryData.searchStr,
      serverIps: queryData.serverIps,
      ...pageChartsParams,
      originFlow: activeAssetsSourseStatus.value.length == 1 && activeAssetsSourseStatus.value.includes(0) ? 1 : null,
      originImport: activeAssetsSourseStatus.value.length == 1 && activeAssetsSourseStatus.value.includes(1) ? 1 : null,
      know: [0, 2].includes(activeAssetsFormStatus.value.length)
        ? null
        : (+activeAssetsFormStatus.value.toString() as 0 | 1),
      netType: [0, 2].includes(activeAssetsCausalityStatus.value.length)
        ? null
        : (+activeAssetsCausalityStatus.value.toString() as 0 | 1),
      isOffline: [0, 2].includes(isOfflineStatus.value.length)
        ? null
        : +isOfflineStatus.value.toString() == 1
        ? true
        : false,
      conflict: [0, 2].includes(activeAssetsConflictStatus.value.length)
        ? null
        : (+activeAssetsConflictStatus.value.toString() as 0 | 1),
      containShandowAsset: [0, 2].includes(activeAssetscontainShandowAssetStatus.value.length)
        ? null
        : (+activeAssetscontainShandowAssetStatus.value.toString() as 0 | 1),
    })
    threatType.value = data.threatType
    threatName.value = data.threatName
    sysType.value = data.sysType
    threatLevel.value = data.threatLevel
    victimIp.value = data.victimIp
    tiemFlag.value = new Date().getTime()
  }

  const handleThreatTrends = (val: number) => {
    pageChartsParams.threatLevelCode = val
    getPageCharts()
  }

  const handleThreatClassification = (val: string) => {
    pageChartsParams.threatName = val
    getPageCharts()
  }

  const handleClickBar = async (val: any) => {
    queryData.serverIps = val
    getPageCharts()
  }

  const handleClickPie = async (val: string) => {
    const threatTypes: string[] = []
    threatTypes.push(val)
    const { data } = await pageChartsApi({
      startTime,
      endTime,
      otherDataSource: queryData.otherDataSource,
      activeAssetType: queryData.activeAssetType,
      dataSourceId: queryData.dataSourceId,
      labelIds: queryData.labelIds,
      labelRelat: queryData.labelRelat,
      searchIps: queryData.searchIps,
      searchStr: queryData.searchStr,
      threatTypes,
      ...pageChartsParams,
      originFlow: activeAssetsSourseStatus.value.length == 1 && activeAssetsSourseStatus.value.includes(0) ? 1 : null,
      originImport: activeAssetsSourseStatus.value.length == 1 && activeAssetsSourseStatus.value.includes(1) ? 1 : null,
      know: [0, 2].includes(activeAssetsFormStatus.value.length)
        ? null
        : (+activeAssetsFormStatus.value.toString() as 0 | 1),
      netType: [0, 2].includes(activeAssetsCausalityStatus.value.length)
        ? null
        : (+activeAssetsCausalityStatus.value.toString() as 0 | 1),
      isOffline: [0, 2].includes(isOfflineStatus.value.length)
        ? null
        : +isOfflineStatus.value.toString() == 1
        ? true
        : false,
      conflict: [0, 2].includes(activeAssetsConflictStatus.value.length)
        ? null
        : (+activeAssetsConflictStatus.value.toString() as 0 | 1),
      containShandowAsset: [0, 2].includes(activeAssetscontainShandowAssetStatus.value.length)
        ? null
        : (+activeAssetscontainShandowAssetStatus.value.toString() as 0 | 1),
    })
    threatType.value = data.threatType
    threatName.value = data.threatName
    sysType.value = data.sysType
    threatLevel.value = data.threatLevel
    victimIp.value = data.victimIp
    tiemFlag.value = new Date().getTime()
  }

  // 展开
  let lock = true
  const tempHeight = 90
  const handleUnFold = () => {
    if (!lock) {
      getPageCharts()
    }
    lock = false
    isShowChart.value = true

    tableNode.value.style.height = isTagContentFold.value ? `${divHeight - 248}px` : `${divHeight - 248 - tempHeight}px`
  }
  // 折叠
  const handleFold = () => {
    isShowChart.value = false
    timeParty.value = 6
    tableNode.value.style.height = isTagContentFold.value ? `${divHeight}px` : `${divHeight - tempHeight}px`
  }

  // 重置页面
  const tiemFlag = ref(new Date().getTime())
  const handleResetPage = () => {
    ;(queryData.pageNum = 1),
      (queryData.pageSize = 10),
      (queryData.searchStr = ''),
      (queryData.lastTime = ''),
      (queryData.labelIds = [] as number[]),
      (queryData.activeAssetType = 0),
      (queryData.dataSourceId = [] as number[]),
      (queryData.conflict = null),
      (queryData.otherDataSource = false),
      (queryData.labelRelat = 'or'),
      (queryData.searchIps = []),
      (queryData.serverIps = []),
      (queryData.originFlow = null),
      (queryData.originImport = null),
      (queryData.know = null),
      (queryData.netType = null),
      (queryData.isOffline = null),
      (queryData.containShandowAsset = null),
      (timeParty.value = 6)
    searchStr.value = undefined
    activeAssetsStatus.value = []
    activeAssetsSourseStatus.value = []
    isOfflineStatus.value = []
    activeAssetsFormStatus.value = []
    activeAssetsCausalityStatus.value = []
    activeAssetsConflictStatus.value = []
    activeAssetscontainShandowAssetStatus.value = []
    pageChartsParams.threatLevelCode = undefined
    pageChartsParams.threatName = undefined
    tiemFlag.value = new Date().getTime()
    rankingOfAssetRef.value.resetKey()
    handleSearch()
  }

  // 重置
  const handleReset = () => {
    queryData.serverIps = []
    pageChartsParams.threatLevelCode = undefined
    pageChartsParams.threatName = undefined
    tiemFlag.value = new Date().getTime()
    timeParty.value = 6
    getPageCharts()
  }
  // 检索IP组
  const handleCheckIPGroup = async (arr: any) => {
    queryData.serverIps = arr
    const {
      data: { records, total },
    } = await getAssetsPreviewListApi({
      ...queryData,
      dataSourceId: queryData.dataSourceId!.filter(Boolean) as number[],
      activeAssetType: [0, 2].includes(activeAssetsStatus.value.length) ? 0 : +activeAssetsStatus.value.toString(),
    })
    tiemFlag.value = new Date().getTime()
    LABEL_KEYS = queryData.labelIds!
    tableData.value = records || []
    list_total.value = total || 0
  }

  watch(
    () => timeParty.value,
    () => {
      startTime = timeNow.subtract(timeParty.value, 'hour').format('YYYY-MM-DD HH:mm:ss')
      queryData.serverIps = []
      pageChartsParams.threatLevelCode = undefined
      pageChartsParams.threatName = undefined
      tiemFlag.value = new Date().getTime()
      // timeParty.value = 6
      getPageCharts()
    }
  )

  // 导出
  const handleExport = async () => {
    const res = await exportAssetNewApi({
      startTime,
      endTime,
      otherDataSource: queryData.otherDataSource,
      activeAssetType: queryData.activeAssetType,
      dataSourceId: queryData.dataSourceId,
      labelIds: queryData.labelIds,
      labelRelat: queryData.labelRelat,
      searchIps: queryData.searchIps,
      searchStr: queryData.searchStr,
      originFlow: activeAssetsSourseStatus.value.length == 1 && activeAssetsSourseStatus.value.includes(0) ? 1 : null,
      originImport: activeAssetsSourseStatus.value.length == 1 && activeAssetsSourseStatus.value.includes(1) ? 1 : null,
      know: [0, 2].includes(activeAssetsFormStatus.value.length)
        ? null
        : (+activeAssetsFormStatus.value.toString() as 0 | 1),
      netType: [0, 2].includes(activeAssetsCausalityStatus.value.length)
        ? null
        : (+activeAssetsCausalityStatus.value.toString() as 0 | 1),
      conflict: [0, 2].includes(activeAssetsConflictStatus.value.length)
        ? null
        : (+activeAssetsConflictStatus.value.toString() as 0 | 1),
      isOffline: [0, 2].includes(isOfflineStatus.value.length)
        ? null
        : +isOfflineStatus.value.toString() == 1
        ? true
        : false,
    })
    downloadFile(res, '资产')
  }

  const changeNetTypeData = (val: number) => {
    if (val == 0) {
      return '内网资产'
    } else if (val == 1) {
      return '外网资产'
    }
    return ''
  }

  watch(
    () => isTagContentFold.value,
    () => {
      if (isTagContentFold.value) {
        tableNode.value.style.height = `${tableNode.value.clientHeight + 90}px`
      } else {
        tableNode.value.style.height = `${tableNode.value.clientHeight - 90}px`
      }
    }
  )

  provide('timePartyProvide', timeParty.value)
  provide('tiemFlag', tiemFlag)
</script>

<template>
  <div class="asset-heat-map-container">
    <vab-query-form style="width: 100%">
      <vab-query-form-left-panel :span="8">
        <span style="font-size: 20px; color: #303133; margin-top: -4px; font-weight: 500">资产总览</span>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="16">
        <div style="display: flex; align-items: center">
          <el-input
            v-model="searchStr"
            class="my-input"
            clearable
            placeholder="模糊检索"
            style="margin-left: 10px; height: 32px"
          />
          <el-button
            :auto-insert-space="false"
            :icon="Search"
            style="margin-left: -2px; height: 32px"
            type="primary"
            @click="handleSearch"
          >
            检索
          </el-button>
          <el-button :auto-insert-space="false" style="margin-left: 10px; height: 32px" @click="handleResetPage">
            重置页面
          </el-button>
          <el-button :auto-insert-space="false" style="margin-left: 10px; height: 32px" @click="handleExport">
            导出
          </el-button>
        </div>
      </vab-query-form-right-panel>
      <vab-query-form-left-panel :span="24">
        <div class="casts">
          <div class="parametric">
            <span class="parametric-title">资产来源:</span>
            <span
              class="assets-tag-item sapn-btn"
              :class="{ xselected: activeAssetsSourseStatus.includes(0) }"
              @click="handleAssetsSourseActiveChange(0)"
            >
              流量识别({{ flowRecognitionNum }})
            </span>
            &nbsp;
            <span
              class="assets-tag-item sapn-btn"
              :class="{ xselected: activeAssetsSourseStatus.includes(1) }"
              @click="handleAssetsSourseActiveChange(1)"
            >
              资产导入({{ assetImportNum }})
            </span>
          </div>
          <div class="parametric">
            <span class="parametric-title">资产类别:</span>
            <span
              class="assets-tag-item sapn-btn"
              :class="{ xselected: activeAssetsFormStatus.includes(1) }"
              @click="handleAssetsFormActiveChange(1)"
            >
              已知资产({{ assetKnowNum }})
            </span>
            &nbsp;
            <span
              class="assets-tag-item sapn-btn"
              :class="{ xselected: activeAssetsFormStatus.includes(0) }"
              @click="handleAssetsFormActiveChange(0)"
            >
              未知资产({{ assetUnKnowNum }})
            </span>
          </div>
          <div class="parametric" style="width: 200px">
            <span class="parametric-title">资产冲突:</span>
            <span
              class="assets-tag-item sapn-btn"
              :class="{ xselected: activeAssetsConflictStatus.includes(1) }"
              @click="handleAssetsconflictChange(1)"
            >
              是({{ assetConflict }})
            </span>
            &nbsp;
            <span
              class="assets-tag-item sapn-btn"
              :class="{ xselected: activeAssetsConflictStatus.includes(0) }"
              @click="handleAssetsconflictChange(0)"
            >
              否({{ assetNoConflict }})
            </span>
          </div>
          <div class="parametric" style="width: 200px; z-index: 10">
            <span class="parametric-title">影子资产:</span>
            <span
              class="assets-tag-item sapn-btn"
              :class="{ xselected: activeAssetscontainShandowAssetStatus.includes(1) }"
              @click="handleAssetscontainShandowChange(1)"
            >
              是({{ assetShandow }})
            </span>
            &nbsp;
            <span
              class="assets-tag-item sapn-btn"
              :class="{ xselected: activeAssetscontainShandowAssetStatus.includes(0) }"
              @click="handleAssetscontainShandowChange(0)"
            >
              否({{ assetNoShandow }})
            </span>
          </div>
          <div class="parametric">
            <span class="parametric-title">资产离线排查:</span>
            <span
              class="assets-tag-item sapn-btn"
              :class="{ xselected: isOfflineStatus.includes(1) }"
              @click="handleisOfflineChange(1)"
            >
              离线资产({{ offlineNum }})
            </span>
            &nbsp;
            <span
              class="assets-tag-item sapn-btn"
              :class="{ xselected: isOfflineStatus.includes(0) }"
              @click="handleisOfflineChange(0)"
            >
              非离线资产({{ notOfflineNum }})
            </span>
          </div>
        </div>
      </vab-query-form-left-panel>
      <vab-query-form-left-panel :span="24">
        <div class="casts">
          <div class="parametric">
            <span class="parametric-title">资产属性:</span>
            <span
              class="assets-tag-item sapn-btn"
              :class="{ xselected: activeAssetsCausalityStatus.includes(0) }"
              @click="handleAssetsCausalityActiveChange(0)"
            >
              内网资产({{ assetInSide }})
            </span>
            &nbsp;
            <span
              class="assets-tag-item sapn-btn"
              :class="{ xselected: activeAssetsCausalityStatus.includes(1) }"
              @click="handleAssetsCausalityActiveChange(1)"
            >
              外网资产({{ assetOutSide }})
            </span>
          </div>
          <div class="parametric">
            <span class="parametric-title">资产类型:</span>
            <span
              class="assets-tag-item sapn-btn"
              :class="{ xselected: activeAssetsStatus.includes(1) }"
              @click="handleAssetsActiveChange(1)"
            >
              活跃资产({{ activityAssetNum }})
            </span>
            &nbsp;
            <span
              class="assets-tag-item sapn-btn"
              :class="{ xselected: activeAssetsStatus.includes(2) }"
              @click="handleAssetsActiveChange(2)"
            >
              非活跃资产({{ notActivityAssetNum }})
            </span>
          </div>

          <div class="center">
            <span class="parametric-title">数据中心:</span>
            <span
              v-for="(netPartition, index) in allNetPartitionSums"
              :key="netPartition.dataSourceId"
              class="assets-tag-item sapn-btn"
              :class="{ xselected: queryData.dataSourceId[index] }"
              style="margin-left: 4px"
              @click="() => handleNetPartitionChange(netPartition.dataSourceId, index)"
            >
              {{ netPartition.dataSourceName }}({{ netPartition.count }}) &nbsp;
            </span>
            <span
              class="assets-tag-item sapn-btn"
              :class="{ xselected: queryData.otherDataSource }"
              @click="handleOtherNetPartitionChange"
            >
              其他({{ otherNetPartationCount }})
            </span>
          </div>
        </div>
      </vab-query-form-left-panel>
    </vab-query-form>
    <!-- 过滤筛选查询区 -->
    <el-row class="x-content" :gutter="20" style="padding-top: 6px">
      <el-col :span="24" style="position: relative">
        <div class="tag-content assetLabel" :style="{ height: isTagContentFold ? '60px' : '150px' }">
          <div class="warp">
            <div
              v-for="field of showFields"
              :key="field.id"
              class="tag-item selectTag"
              :class="{ selected: queryData.labelIds?.includes(field.id) }"
              @click.self="handleFieldChange(field.id)"
            >
              {{ field.labelName }}
              <div class="selectTag-closeable" @click.stop="handleTagClose(field.id)">
                <Close />
              </div>
            </div>
          </div>
          <div style="position: absolute; top: 50%; transform: translateY(-50%); height: 60px; right: 0; width: 102px">
            <ul class="labelRelat">
              <li :class="{ active: queryData.labelRelat == 'or' }" @click="labelRelatChange('or')">或</li>
              <li :class="{ active: queryData.labelRelat == 'and' }" @click="labelRelatChange('and')">且</li>
              <div style="margin-left: 6px">
                <div>
                  <el-icon
                    v-if="isTagContentFold"
                    style="width: 30px; height: 30px; border-radius: 0px 4px 0px 0px; border: 1px solid #d8dce6"
                    @click="isTagContentFold = false"
                  >
                    <ArrowDown />
                  </el-icon>
                  <el-icon
                    v-else
                    style="width: 30px; height: 30px; border-radius: 0px 4px 0px 0px; border: 1px solid #d8dce6"
                    @click="isTagContentFold = true"
                  >
                    <ArrowUp />
                  </el-icon>
                </div>
                <div>
                  <el-icon
                    style="
                      width: 30px;
                      height: 30px;
                      border-radius: 0px 0px 4px 0px;
                      border: 1px solid #d8dce6;
                      border-top: none;
                    "
                    @click="openLableDialog"
                  >
                    <Setting />
                  </el-icon>
                </div>
              </div>
            </ul>
          </div>
        </div>
      </el-col>
    </el-row>
    <!-- 图表 -->
    <div v-if="!isShowChart" class="x-fold" @click="handleUnFold">
      <el-icon color="#fff"><ArrowDown /></el-icon>
    </div>

    <div v-if="isShowChart" class="assets-chart">
      <el-select v-model="timeParty" class="my-x-time" size="small">
        <el-option label="6小时" :value="6" />
        <el-option label="12小时" :value="12" />
        <el-option label="24小时" :value="24" />
        <el-option label="72小时" :value="72" />
      </el-select>
      <el-row :gutter="20" style="height: 100%">
        <el-col :span="9" style="border-right: 1px solid var(--el-border-color)">
          <el-tabs style="height: 100%" :tab-position="'left'">
            <el-tab-pane label="威胁趋势">
              <ThreatTrends :threat-levels="threatLevel" @on-click-event="handleThreatTrends" />
            </el-tab-pane>
            <el-tab-pane label="系统分布"><SystemDistribution :sys-types="sysType" /></el-tab-pane>
            <el-tab-pane label="威胁告警">
              <ThreatClassification :threat-types="threatName" @on-click-event="handleThreatClassification" />
            </el-tab-pane>
          </el-tabs>
        </el-col>
        <el-col :span="15">
          <el-row :gutter="20" style="height: 100%">
            <el-col :span="12">
              <RankingOfAsset
                ref="rankingOfAssetRef"
                :show-btn="true"
                :victim-ips="victimIp"
                @on-check-ips="handleCheckIPGroup"
                @on-click-bar="handleClickBar"
                @on-reset="handleReset"
              />
            </el-col>
            <el-col :span="12">
              <ThreateningIncident :threat-types="threatType" @on-click-pie="handleClickPie" />
            </el-col>
          </el-row>
        </el-col>
      </el-row>
      <div class="x-fold x-fold-up" @click="handleFold">
        <el-icon color="#fff"><ArrowUp /></el-icon>
      </div>
    </div>

    <el-table class="assets-table" :data="tableData" @selection-change="setSelectRows">
      <!-- <el-table-column align="center" type="selection" width="50" /> -->
      <el-table-column align="center" label="序号" show-overflow-tooltip type="index" width="55" />
      <el-table-column align="center" label="名称" prop="appName" show-overflow-tooltip width="155" />
      <el-table-column
        align="center"
        :formatter="({ ipv4, ipv6 }) => ipv4 || ipv6"
        label="资产IP"
        show-overflow-tooltip
        width="180"
      />
      <el-table-column align="center" label="业务名称" prop="businessName" show-overflow-tooltip width="155" />
      <el-table-column align="center" label="资产属性" prop="netTypeStr" width="155" />
      <el-table-column
        align="center"
        :formatter="({ updateTime }) => formatNstime(updateTime, false)"
        label="最近在线时间"
        prop="updateTime"
        show-overflow-tooltip
        width="180"
      />
      <el-table-column align="center" label="资产标签 " prop="labels" show-overflow-tooltip>
        <template #default="{ row }">
          <div class="labelBox foldBox">
            <div v-for="(item, index) in tableLabel(row.labels)" :key="index" :alt="item" class="table-label">
              {{ item }}
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column align="center" fixed="right" label="操作" width="80">
        <template #default="{ row }">
          <el-button size="small" @click="handleClick(row)">详情</el-button>
        </template>
      </el-table-column>

      <template #empty><el-empty /></template>
    </el-table>
    <el-pagination
      v-model:current-page="queryData.pageNum"
      v-model:page-size="queryData.pageSize"
      background
      layout="total, sizes, prev, pager, next, jumper"
      :page-sizes="[20, 30, 50, 100]"
      :total="list_total"
      @current-change="handleSearch"
      @size-change="handleSearch"
    />
    <vab-dialog v-model="showInfo" destroy-on-close show-fullscreen title="资产管理">
      <asset-manager :preview-info="previewInfo" />
    </vab-dialog>
    <assets-label ref="lableRef" :selected-fieelds="showFields" @confirm="handleLableConfirm" />
  </div>
</template>

<style scoped lang="scss">
  .asset-heat-map-container {
    .assets-chart {
      position: relative;
      .my-x-time {
        position: absolute;
        top: 15px;
        left: 13px;
        width: 73px;
        z-index: 99;
      }
    }
    position: relative;
    // margin-top: 5px;

    .casts {
      display: flex;
      align-items: center;
    }
    .parametric {
      line-height: 32px;
      height: 32px;
      width: 305px;
      overflow: hidden;
    }
    .center {
      width: calc(100% - 610px);
      height: 32px;
      line-height: 32px;
      overflow-y: auto;
      // white-space: nowrap;
      word-break: break-all;
      word-wrap: break-word;

      &::-webkit-scrollbar {
        width: 0px;
        height: 2px;
      }
    }
    .parametric-title {
      color: #a9acb3;
    }
    .sapn-btn {
      &:hover {
        cursor: pointer;
      }
    }
    .assets-tag-item {
      color: #303133;
      height: 32px;
      &.xselected {
        color: var(--el-color-primary) !important;
        font-weight: 600;
      }
      &.selected {
        --el-button-bg-color: var(--el-color-primary-light-9) !important;
        --el-button-border-color: var(--el-color-primary-light-7) !important;
        --el-button-text-color: var(--el-color-primary) !important;
        font-weight: 600;
      }
    }
    .is-plain.assets-tag-item:hover,
    .is-plain.assets-tag-item:focus {
      color: var(--el-button-text-color);
      border-color: var(--el-button-border-color);
      background-color: var(--el-button-bg-color);
    }
    :deep() {
      .el-input__wrapper:not(:focus) {
        width: 110px !important;
        z-index: 0;
      }
      .el-input__wrapper.is-focus {
        width: 300px !important;
        z-index: 99;
      }

      .left-panel {
        display: flex !important;
        align-items: center !important;
        margin: 0 !important;
        line-height: 40px;
      }
      .el-checkbox__label {
        color: #303133;
        font-weight: 400;
      }
    }

    .tag-content {
      width: 100%;
      height: 60px;
      padding: 15px 0 15px 20px;
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
      border: 1px solid #d8dce6;
      border-radius: 2px;
      .warp {
        width: calc(100% - 100px);
      }
      .tag-item {
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
        &.selectTag {
          position: relative;
          .selectTag-closeable {
            position: absolute;
            top: -5px;
            right: -5px;
            width: 12px;
            height: 12px;
            line-height: 12px;
            color: #fff;
            text-align: center;
            background-color: #ff6060;
            border-radius: 100%;
            svg {
              width: 10px;
              height: 10px;
            }
          }
          &.selected {
            color: #fff;
            background-color: var(--el-color-primary);
          }
        }
      }
      &.diy {
        display: flex;
        // flex-direction: column;
        align-items: center;
        justify-content: space-evenly;
        padding: 0;
        .tag-item {
          margin-bottom: 0;
        }
      }
      &.assetLabel {
        &:hover {
          .selectTag-closeable {
            display: block;
          }
        }
        .selectTag-closeable {
          display: none;
        }
      }
    }

    .x-fold {
      float: right;
      background-color: #726aa3;
      height: 20px;
      width: 20px;
      border-radius: 2px;
      display: flex;
      align-items: center;
      justify-content: center;
      .x-icon {
        font-size: 12px;
      }
      &:hover {
        cursor: pointer;
      }
    }

    .assets-chart {
      margin-top: 14px;
      border: 1px solid var(--el-border-color);
      width: 100%;
      height: 240px;
      border-radius: 2px;
      position: relative;
      .x-fold-up {
        position: absolute;
        right: 0px;
        top: 0px;
      }

      :deep() {
        .el-tabs__active-bar {
          display: none;
        }
        .el-tabs__item {
          height: 40px;
        }
        .el-tabs__nav {
          margin-top: 57px;
        }
      }
    }

    .assets-table {
      margin-top: 14px;
      height: calc(100vh - 270px);
      .labelBox {
        display: flex;
        &.foldBox {
          display: block;
          .table-label {
            display: inline-block;
          }
        }
        .table-label {
          padding: 0 10px;
          height: 30px;
          border-radius: 4px;
          line-height: 30px;
          background: #f1f5fb;
          margin-right: 10px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }

      :deep() {
        // .cell {
        //   display: flex;
        // }
        .el-table__row.hover-row {
          td.el-table__cell {
            background-color: #fcfdff !important;
          }
        }
      }
    }

    :deep() {
      .labelRelat {
        display: flex;
        align-items: center;
        padding-left: 0;
        margin-bottom: 8px;
        cursor: pointer;
        li {
          &:not(.tips) {
            width: 28px;
            height: 24px;
            font-size: 13px;
            line-height: 24px;
            text-align: center;
            border-radius: 4px;
            cursor: pointer;
          }

          &.active {
            background-color: var(--el-color-primary);
            color: #fff;
          }
          &.tips {
            margin-left: 8px;
            color: #a9acb3;
          }
        }
      }
      .el-form--inline .el-form-item {
        margin-right: 0;
      }
      .el-table__empty-block {
        position: absolute;
        top: 50%;
        transform: translateY(-50%);
      }
    }
  }
</style>
