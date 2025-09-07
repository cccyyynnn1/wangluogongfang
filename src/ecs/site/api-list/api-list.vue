<script lang="ts">
  export default {
    name: 'ApiList', // API列表
  }
</script>
<script setup lang="ts">
  import TraceabilityFieldStatistic from '@/ecs/alert/components/field-statistic.vue'
  import { useTableCopy, useCopy } from '@/utils'
  import { AlertItem } from '@/types'
  import { TableColumnItemType } from '/#/store'
  import dayjs from 'dayjs'
  import { timeDuratioOptions as date_options } from '@/data/constant'
  import AlertDetail from '../../alert/components/alert-detail.vue'
  import { level_options, levelKey, timeDuratioOptions } from '../../alert/data/index'
  import { getAlertApi } from '@/api-ecs/alert'
  import { useScroll, useMouseInElement, useElementBounding } from '@vueuse/core'
  import { getAllDisPlaysFiledApi } from '~/src/api-ecs/public'
  import { useUserStore } from '@/store/modules/user'
  import { getSiteApiListApi, getListBySiteApiApi, exportResponseDataApi } from '~/src/api-ecs/site'
  import { getSystemConfigApi } from '~/src/api-ecs/system'
  import { uuid } from '~/src/utils'
  import inputSearch from '~/library/components/VabInputSearch/index.vue'
  import _lodash from 'lodash'
  import html2canvas from 'html2canvas'
  const userStore = useUserStore()
  const traceabilityFieldRef = ref<InstanceType<typeof TraceabilityFieldStatistic>>()
  const multipleSelection = ref<AlertItem[]>([])
  const $baseMessage: any = inject('$baseMessage')
  const { getTableColumn } = userStore
  const props = defineProps<{
    hosts: string | null
    siteName: string | null
  }>()
  const isAdvanced = ref(false)
  const loading = ref(false)
  const isLoading = ref(false)
  const hosts = ref()
  const dataList = ref()
  const itemIndex = ref(0)
  const alternativeType = ref(false) // 单选或多选，默认单选
  const showToggle = ref(false) // 打开选择
  const activeName = ref('first')
  const activeArr = ref(['query', 0])
  const activeResp = ref([0, 0])
  const hoverValueArr = ref<any[]>([])
  const scaleResNum = ref(1)
  const scaleRespNum = ref(1)
  const canDragRes = ref(false)
  const canDragResp = ref(false)
  const canDragView = ref(false)
  const currentData = ref<{ host: string; id: string; method: string; type: number; path: string }>({
    host: '',
    id: '',
    path: '',
    method: '',
    type: 0,
  })
  const queryPData = ref()
  const headerData = ref()
  const responseData = ref()
  const rightJson = ref()
  const treeData = ref()
  const tagsIndex = ref()
  const prexArr = ref<any[]>([])
  const showMap = ref(true)
  const checkedNodesArr = ref<{ fullUrl: string }[]>([])
  const ctxZoom = ref(0)
  interface ColorObj {
    [key: string]: any
  }
  const moduleEnable = ref(false)
  const showPrevBtn = ref(false)
  const showNextBtn = ref(false)

  let nodes: HTMLElement[] | undefined = undefined

  const data = ref()

  const el = ref<HTMLElement | null>(null)
  // 告警检索条件
  const queryForm = reactive({
    /** 查询SQL */
    searchSql: '',
    /** 源IP */
    attackIp: '',
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
    threatType: '',
    /** 威胁名称 */
    threatName: '',
    /** 威胁等级 */
    threatLevel: '',
    /** 攻击结果 */
    attackResult: '',
    /** 攻击次数 */
    attackCount: '',
    /** 已读状态 */
    readStatus: '',
  })
  // 查询时间
  const timeDuration = ref('1hours')
  // 高级筛选
  // const isAdvanced = ref(false)

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
    whiteType: 0,
  })

  // 告警详情显示
  const alertDetailVisible = ref(false)

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

  const infoVal = ref()

  const siteSearch = ref()

  const userDisPlaysFiled = ref()

  // 自定义时间
  const timeDate = ref()
  const showBottom = ref<number[]>([])
  const tableColumn = ref<TableColumnItemType[]>([])
  // 获取表格序号
  const curIndex = computed(() => (must_parames.pageNum - 1) * must_parames.pageSize + 1)

  // 多选表格
  function handleSelectionChange(val: AlertItem[]) {
    multipleSelection.value = val
  }

  const showTraceabilityField = (fieldCn: string, fieldEn: string) => {
    let str = ''
    const arr = props.hosts?.split(',')
    arr?.forEach((_: any, index: any) => {
      if (index == 0) {
        str += `host = "${arr[index]}"`
      } else {
        str += ` or host = ${arr[index]}`
      }
    })
    if (checkedNodesArr.value && checkedNodesArr?.value?.length > 0) {
      checkedNodesArr.value.forEach((_: any, index: any) => {
        if (index == 0) {
          str = `(${str}) and (url like "${checkedNodesArr.value[index].fullUrl}"`
        } else {
          str += ` or url like "${checkedNodesArr.value[index].fullUrl}"`
        }
      })
      str += ')'
    }
    const { indexType, startTime, endTime, whiteType } = must_parames
    traceabilityFieldRef.value?.initData({
      type: 'alert',
      title: fieldCn,
      query: { indexType, startTime, endTime, aggregationFields: fieldEn, searchSql: str, whiteType: whiteType },
    })
  }
  const obj = { 低危: 1, 中危: 2, 高危: 3, 危急: 4 }
  function getLevel(str: string) {
    // @ts-ignore
    const level = obj[str]
    return levelKey[level]
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

  // 获取用户所有字段
  const getAllDisPlaysFiled = async () => {
    const { data } = await getAllDisPlaysFiledApi()
    userDisPlaysFiled.value = data
    formatColum()
  }

  function formatColum() {
    const columnData = getTableColumn(must_parames.indexType)
    const userColumnData = userDisPlaysFiled.value[must_parames.indexType] as number[]
    if (userColumnData?.length > 0) {
      tableColumn.value = userColumnData.map((key) =>
        columnData.find((item) => item.id === key)
      ) as TableColumnItemType[]
    } else {
      tableColumn.value = []
    }
  }

  const resetData = () => {
    alertInfo.scrollId = ''
    alertInfo.list = []
  }

  // 查询告警日志
  async function queryData(remark = true, updateChart = true) {
    if (!queryForm.host) return
    if (remark) {
      resetData()
    }
    if (!timeDate.value && timeDuration.value === 'user-defined')
      return $baseMessage('请选择时间！', 'error', 'vab-hey-message-error')
    loading.value = true
    if (timeDuration.value !== 'user-defined') {
      formatDate()
    }
    try {
      const { searchSql } = queryForm
      let str = ''
      const arr = props.hosts?.split(',')
      arr?.forEach((_: any, index: any) => {
        if (index == 0) {
          if (arr[index].endsWith(':*')) {
            str += `host like "${arr[index].slice(0, -2)}"`
          } else {
            str += `host = "${arr[index]}"`
          }
        } else {
          if (arr[index].endsWith(':*')) {
            str += ` or host like "${arr[index].slice(0, -2)}"`
          } else {
            str += ` or host = "${arr[index]}"`
          }
        }
      })
      if (checkedNodesArr.value && checkedNodesArr?.value?.length > 0) {
        checkedNodesArr.value.forEach((_: any, index: any) => {
          if (index == 0) {
            str = `(${str}) and (url like "${checkedNodesArr.value[index].fullUrl}"`
          } else {
            str += ` or url like "${checkedNodesArr.value[index].fullUrl}"`
          }
        })
        str += ')'
      }
      const queryData = isAdvanced.value
        ? {
            ...must_parames,
            threatType: queryForm.threatType,
            searchSql: `${searchSql.trim()} and (${str})`,
            scrollId: alertInfo.scrollId,
          }
        : { ...queryForm, searchSql: str, ...must_parames, scrollId: alertInfo.scrollId }
      const {
        data: { total, sumaryMap, resList, scrollId },
      } = await getAlertApi(queryData)
      alertInfo.scrollId = scrollId
      if (resList.length > 0) {
        const id = resList[0].id
        const flag = alertInfo.list.some((item: any) => {
          return item.id == id
        })
        if (flag) return false
        alertInfo.list = [...alertInfo.list, ...resList]
      }

      alertInfo.total = total
      if (updateChart && sumaryMap) {
        const { pillar, pie, threat_level } = JSON.parse(sumaryMap)
        alertInfo.pieData = pie || []
        alertInfo.lineData = {
          date: pillar.date,
          count: pillar.count,
        }
        alertInfo.threat_level = threat_level || []
      }
      formatColum()
      loading.value = false
    } catch (error) {
      loading.value = false
      console.error(error)
    }
  }

  const formatDate = () => {
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
    must_parames.endTime =
      timeDuration.value === 'yesterday'
        ? timeDate.endOf('day').subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss')
        : endDate
    must_parames.startTime = startDate
  }

  const tagsList = ref<any[]>([])

  let resViewClickX = 0
  let resViewClickY = 0
  let resViewEndX = 0
  let resViewEndY = 0
  let resViewOffsetX = 0
  let resViewOffsetY = 0

  let respViewClickX = 0
  let respViewClickY = 0
  let respViewEndX = 0
  let respViewEndY = 0
  let respViewOffsetX = 0
  let respViewOffsetY = 0

  const initResViewPosition = () => {
    resViewClickX = 0
    resViewClickY = 0
    resViewEndX = 0
    resViewEndY = 0
    resViewOffsetX = 0
    resViewOffsetY = 0
  }
  const initRespViewPosition = () => {
    respViewClickX = 0
    respViewClickY = 0
    respViewEndX = 0
    respViewEndY = 0
    respViewOffsetX = 0
    respViewOffsetY = 0
  }

  let offsetLeftRes = 0
  let offsetTopRes = 0

  let offsetLeftResp = 0
  let offsetTopResp = 0
  const ininitialisationViewNodePositionRes = () => {
    offsetLeftRes = 0
    offsetTopRes = 0
  }
  const ininitialisationViewNodePositionResp = () => {
    offsetLeftResp = 0
    offsetTopResp = 0
  }

  // 思维导图
  let resPosX = 0
  let resPosY = 0
  let resCurX = 0
  let resCurY = 0
  let resOffsetX = 0
  let resOffsety = 0
  let respPosX = 0
  let respPosY = 0
  let respCurX = 0
  let respCurY = 0
  let respOffsetX = 0
  let respOffsety = 0
  const ininitialisationDragData = () => {
    ininitialisationResDragData()
    ininitialisationRespDragData()
    ininitialisationViewNodePositionRes()
    ininitialisationViewNodePositionResp()
    initResViewPosition()
    initRespViewPosition()
  }

  const ininitialisationResDragData = () => {
    resPosX = 0
    resPosY = 0
    resCurX = 0
    resCurY = 0
    resOffsetX = 0
    resOffsety = 0
  }

  const ininitialisationRespDragData = () => {
    respPosX = 0
    respPosY = 0
    respCurX = 0
    respCurY = 0
    respOffsetX = 0
    respOffsety = 0
  }

  const cleaeData = () => {
    showBottom.value = []
    queryPData.value = undefined
    headerData.value = undefined
    responseData.value = undefined
    siteSearch.value = undefined
    rightJson.value = undefined
    tagsList.value = []
    prexArr.value = []
    currentData.value.host = ''
    currentData.value.method = ''
    currentData.value.path = ''
    alternativeType.value = false
    // length.value = undefined
    tagsIndex.value = undefined
    dataList.value = undefined
    showToggle.value = false
    checkedNodesArr.value = []
    activeName.value = 'first'
    activeArr.value = ['query', 0]
    activeResp.value = [0, 0]
    scaleResNum.value = 1
    scaleRespNum.value = 1
    ininitialisationDragData()
  }

  const resetDataWhenClickItem = () => {
    showBottom.value = []
    queryPData.value = undefined
    headerData.value = undefined
    responseData.value = undefined
    siteSearch.value = undefined
    rightJson.value = undefined
    tagsList.value = []
    prexArr.value = []
    currentData.value.host = ''
    currentData.value.method = ''
    currentData.value.path = ''
    alternativeType.value = false
    tagsIndex.value = undefined
    dataList.value = undefined
    checkedNodesArr.value = []
    activeName.value = 'first'
    activeArr.value = ['query', 0]
    activeResp.value = [0, 0]
    scaleResNum.value = 1
    scaleRespNum.value = 1
    ininitialisationDragData()
  }

  const itemClick = (data: any) => {
    if (data.type == 1) return
    if (currentData.value['host'] !== data.fullUrl) {
      resetDataWhenClickItem()
    }
    currentData.value['host'] = data.fullUrl
    currentData.value['id'] = data.id
    currentData.value['method'] = data.method
    currentData.value['type'] = data.type
    currentData.value['path'] = data.name
    setTimeout(() => {
      initialisationWheelScrollEvent()
      ininitialisationDragData()
    }, 0)
  }

  const formatParameters = (parameters: any) => {
    if (!parameters) return
    queryPData.value = {
      type: '',
      properties: [],
    }
    headerData.value = {
      type: '',
      properties: [],
    }
    parameters.forEach((item: any) => {
      const res = item.schema
      const obj1 = {
        name: '',
        in: '',
        type: '',
        connection: [],
      }
      obj1.in = item.in
      obj1.connection = []
      obj1.name = item.name
      const oneOf = res.oneOf
      const tpyeArr: string[] = []
      oneOf?.forEach((element: any) => {
        const obj2 = {
          one_of: '',
          value: undefined,
          description: '',
        }
        tpyeArr.push(element.type)
        obj2.one_of = element.type
        obj2['value'] = element.enum
        obj2.description = element.description
        // @ts-ignore
        obj1.connection.push(obj2)
      })
      obj1.type = tpyeArr.join(',')
      if (item.in == 'query' || item.in == 'path') {
        queryPData.value.properties.push(obj1)
      } else {
        headerData.value.properties.push(obj1)
      }
    })
  }

  const formatResponses = (responses: any) => {
    if (!responses) return
    responseData.value = []
    for (const key in responses) {
      const obj = {
        code: '',
        content: '',
        properties: [],
      }
      obj.code = key
      const { content } = responses[key]
      for (const k in content) {
        obj.content = k
        const res = content[k].properties
        obj.properties = []
        for (const t in res) {
          const obj1 = {
            name: '',
            connection: [],
            type: '',
          }
          obj1.connection = []
          // obj1.name = `"${t}"`
          obj1.name = t
          const tpyeArr: string[] = []
          const { description, type, enum: value } = res[t]
          const obj2 = {
            one_of: '',
            value: undefined,
            description,
          }
          tpyeArr.push(type)
          obj2.one_of = type
          obj2['value'] = value
          obj1.type = tpyeArr.join(',')
          // @ts-ignore
          obj1.connection.push(obj2)
          // @ts-ignore
          obj.properties.push(obj1)
        }
      }
      if (obj.code) {
        responseData.value.push(obj)
        // formatRightData(responseData.value[0])
      }
    }
  }

  // 获取右侧数据
  const getAlldata = async () => {
    if (currentData.value.type) return
    const { data } = await getSiteApiListApi({ url: currentData.value.host })
    dataList.value = []
    tagsList.value = []
    for (const key in data) {
      dataList.value.push(data[key])
      tagsList.value.push(key)
    }
    if (dataList.value.length > 0) {
      tagsClick(0)
    }
  }

  const tagsClick = (index: number) => {
    if (index == tagsIndex.value) return
    showBottom.value = []
    queryPData.value = undefined
    headerData.value = undefined
    responseData.value = undefined
    rightJson.value = undefined
    activeName.value = 'first'
    activeArr.value = ['query', 0]
    activeResp.value = [0, 0]
    const { parameters, responses } = dataList.value[index]
    formatParameters(parameters)
    formatResponses(responses)
    initActiveItem()
    tagsIndex.value = index
  }

  const initActiveItem = () => {
    let _index = 0
    if (queryPData.value.properties.length == 0) {
      _index = headerData.value.properties.findIndex((item: { connection: any[] }) => {
        return item.connection.length
      })
      activeArr.value[0] = 'head'
    } else {
      _index = queryPData.value.properties.findIndex((item: { connection: any[] }) => {
        return item.connection.length
      })
    }
    activeArr.value[1] = _index == -1 ? 0 : _index
    let activeResp0 = 0
    let activeResp1 = 0
    activeResp0 = responseData.value.findIndex((item: any) => {
      return item.properties.length > 0
    })
    if (activeResp0 != -1) {
      activeResp1 = recursiveGetActiveResp(activeResp0) || 0
      activeResp.value = [activeResp0, activeResp1]
    }
    initWheelScroll()
  }

  const recursiveGetActiveResp = (idx: number): number => {
    let activeRespOne = 0
    if (idx != -1) {
      activeRespOne = responseData.value[idx].properties.findIndex((item: any) => {
        return item.connection.length > 0
      })
    }
    if (activeRespOne == -1 && idx + 1 < responseData.value.length) {
      return recursiveGetActiveResp(idx + 1)
    } else {
      return activeRespOne
    }
  }

  watch(
    () => currentData.value?.host,
    async () => {
      if (currentData.value?.host) {
        getAlldata()
        queryForm.host = currentData.value?.host
        queryData()
      } else {
        cleaeData()
      }
    }
  )

  watch(
    () => timeDuration.value,
    () => {
      if (timeDuration.value !== 'user-defined') {
        timeDate.value = undefined
      }
      formatDate()
      queryData()
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

  watchEffect(() => {
    if (timeDate.value) {
      formatUserTime()
    }
  })

  // 递归得到目录
  type objType = {
    [key: string]: string[]
  }
  const getCataloguesByRecurrence = (arr: string[], prex: string) => {
    const newOBJ: objType = {}
    if (arr.length < 2) return prexArr.value.push(prex + arr[0])
    arr.forEach((item: any) => {
      const arrUrl = item.split('/')
      if (arrUrl[0] && arrUrl.length > 1) {
        const KeyStr = `${prex}/${arrUrl[0]}`
        if (!newOBJ[KeyStr]) newOBJ[KeyStr] = []
        newOBJ[KeyStr].push(arrUrl.slice(1).join('/') || '/')
      } else if (arrUrl[0] && arrUrl.length == 1) {
        prexArr.value.push(`${prex}/`)
      }
    })
    for (const key in newOBJ) {
      if (newOBJ[key].length > 1) {
        prexArr.value.push(key)
        getCataloguesByRecurrence(newOBJ[key], key)
      }
    }
    prexArr.value = Array.from(new Set(prexArr.value))
  }

  const formatDataTree = (dataTree: any[]) => {
    // 1：变成{ _ : [...]}数据结构
    const dataTreeToObj: string[] = []
    dataTree.forEach((item: any) => {
      dataTreeToObj.push(item.url)
    })
    // 2: 递归得到
    if (dataTreeToObj.length == 0) return
    getCataloguesByRecurrence(dataTreeToObj, '')
    prexArr.value = prexArr.value.reverse()
    prexArr.value.forEach((item, index) => {
      if (prexArr.value[index] == '/') {
        prexArr.value[index] = '/'
      } else {
        prexArr.value[index] = item.substring(1)
      }
      if (prexArr.value[index][prexArr.value[index].length - 1] != '/') {
        prexArr.value[index] += '/'
      }
    })
    prexArr.value = Array.from(new Set(prexArr.value))
  }

  // 将dataTree改成目录结构
  const getTreeByPrexArr = (dataTree: any) => {
    const arr: any[] = []
    prexArr.value.forEach((item: any) => {
      if (item !== '/') {
        const obj = {
          type: 1,
          id: uuid(),
          name: item,
          fullUrl: item,
          url: item,
          children: 0,
          checked: false,
          hasWarn: false,
          multiple: false,
          // showToggle: false,
          overflow: true,
          export: true,
          disabled: 'flex',
        }
        arr.push(obj)
        let warn = false
        dataTree.forEach((td: any) => {
          if (td.url.slice(0, item.length) == item && !td.id) {
            td['type'] = 0
            td['id'] = uuid()
            td['host'] = td['fullUrl'] = td.url
            td['method'] = undefined
            td['checked'] = false
            td['overflow'] = true
            td['disabled'] = 'flex'
            td['multiple'] = false
            // td['showToggle'] = false
            td['name'] = td.url.slice(item.length) || '/'
            td['pUrl'] = item.slice(0)
            if (!warn) {
              warn = td.hasWarn
            }
            arr.forEach((i) => {
              if (i.fullUrl == item) {
                i.children += 1
                if (!i.hasWarn) {
                  i.hasWarn = warn
                }
              }
            })
            arr.push(td)
          }
        })
      }
    })
    dataTree.forEach((td: any) => {
      if (!td.id) {
        td['type'] = 0
        td['id'] = uuid()
        td['host'] = td['fullUrl'] = td['name'] = td.url
        td['method'] = undefined
        td['overflow'] = true
        td['checked'] = false
        td['multiple'] = false
        td['disabled'] = 'flex'
        arr.push(td)
      }
    })
    treeData.value = arr.filter((item) => {
      return (item.type == 1 && item.children > 0) || item.type == 0
    })
  }

  // 获取左侧树
  const getSiteApiList = async () => {
    if (hosts.value) {
      let arr: any[] = []
      // try {
      //   const obj = JSON.parse(hosts.value)
      //   console.log('🚀 ~ getSiteApiList ~ obj:', hosts.value)
      //   for (const key in obj) {
      //     const selectArr = obj[key].selected
      //     if (selectArr.length > 0) {
      //       selectArr.forEach((item: any) => {
      //         arr.push(`${key}:${item}`)
      //       })
      //     } else {
      //       arr.push(`${key}`)
      //     }
      //   }
      // } catch (err) {
      //   arr = hosts.value.split(',')
      // }
      try {
        isLoading.value = true
        const { data } = await getListBySiteApiApi({
          startTime: must_parames.startTime,
          endTime: must_parames.endTime,
          hosts: hosts.value,
        })
        const dataTree = await data?.tree
        formatDataTree(dataTree)
        if (prexArr.value.length < 1) return
        // 将dataTree改成目录结构
        getTreeByPrexArr(dataTree)
        setTimeout(() => {
          checkOverflow()
        }, 0)
        getCurrent()
      } finally {
        isLoading.value = false
      }
    }
  }

  interface Tree {
    [key: string]: any
  }

  const getCurrent = () => {
    const data = treeData.value.find((item: { type: number }) => {
      return item.type == 0
    })
    currentData.value['host'] = data.fullUrl
    currentData.value['id'] = data.id
    currentData.value['method'] = data.method
    currentData.value['type'] = data.type
    currentData.value['path'] = data.name
  }

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
  const getModelStatus = async () => {
    const {
      data: { module_config },
    } = await getSystemConfigApi({ keys: 'module_config' })
    moduleEnable.value = JSON.parse(module_config.value).isEnable || false
  }
  watch(
    () => props.hosts,
    () => {
      hosts.value = props.hosts
      cleaeData()
      treeData.value = undefined
      showMap.value = true
      getSiteApiList()
      itemIndex.value = 0
      resetData()
    },
    { immediate: true }
  )

  const route = useRoute()
  onMounted(async () => {
    getModelStatus()
    await getAllDisPlaysFiled()
    // queryData()
    initDragNode()
    // @ts-ignore
    nodes = document.querySelectorAll('.el-table__body-wrapper')
    el.value = nodes![0] as HTMLElement
    const { arrivedState } = useScroll(el.value)
    data.value = arrivedState
    if (route.query.timeDuration) {
      timeDuration.value = route.query.timeDuration as string
    }
  })

  const defaultProps = {
    children: 'children',
    label: 'name',
    disabled: 'hasChildren',
  }

  watch(
    () => data.value,
    () => {
      if (data.value.bottom) {
        queryData(false, false)
      }
    },
    { deep: true }
  )

  const handleCopy = () => {
    useCopy(JSON.stringify(rightJson.value))
  }

  function formatDataTime(row: any, key: string) {
    const time = row[key] / 1000000
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
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
  const listLoading = ref(false)
  // 导出
  // const exportEvent = async () => {
  //   listLoading.value = true
  //   try {
  //     let arr = []
  //     arr = hosts.value.split(',')
  //     if (arr.length == 0) return (listLoading.value = false)
  //     const res = await exportResponseDataApi({ hosts: arr })
  //     const name = props.siteName
  //     downloadFile(res, `${name}`)
  //     ElMessage({ message: '导出成功', type: 'success' })
  //   } finally {
  //     listLoading.value = false
  //   }
  // }

  const handleToggle = () => {
    try {
      listLoading.value = true
      showToggle.value = !showToggle.value
      alternativeType.value = false
      treeData.value?.forEach((item: { checked: boolean; multiple: boolean }) => {
        item.checked = false
        item.multiple = false
      })
      if (!showToggle.value) queryData()
    } finally {
      setTimeout(() => {
        listLoading.value = false
      }, 300)
    }
  }

  const handleCheckedNodes = _lodash.debounce(() => {
    checkedNodesArr.value = []
    treeData.value.forEach((item: any) => {
      if (item.type == 0 && item.checked && item.disabled == 'flex') {
        checkedNodesArr.value.push(item)
      }
    })
    if (checkedNodesArr.value.length > 0) queryData()
  }, 1000)

  // 检查文字是否超出div
  const checkOverflow = () => {
    const nodes = ref()
    nodes.value = document.querySelectorAll('.checkOverflow')
    const overArr: number[] = []
    nodes.value.forEach((_: any, index: number) => {
      if (nodes.value[index].scrollWidth > nodes.value[index].offsetWidth) {
        overArr.push(index)
      }
    })
    if (overArr.length > 0) {
      overArr.forEach((item) => {
        treeData.value[item].overflow = false
      })
    }
  }

  // 子站点展开或闭合
  const toggleExpand = (idx: number) => {
    treeData.value[idx].export = !treeData.value[idx].export
    const child: any[] = []
    const pUrl = treeData.value[idx].name
    treeData.value.forEach((_: any, index: number) => {
      if (treeData.value[index].pUrl == pUrl && treeData.value[index].type == 0) {
        child.push(index)
        treeData.value[index].disabled = treeData.value[idx].export ? 'flex' : 'none'
      }
    })
  }

  // 切换复选框
  const handleChange = (idx: number) => {
    if (treeData.value[idx].type == 1) {
      const pUrl = treeData.value[idx].name
      treeData.value.forEach((_: any, index: number) => {
        if (treeData.value[index].pUrl == pUrl && treeData.value[index].type == 0) {
          treeData.value[index].checked = !treeData.value[idx].checked
        }
      })
      handleCheckedNodes()
    } else {
      setTimeout(() => {
        const pUrl = treeData.value[idx].pUrl
        const arr = ref<boolean[]>([])
        let _index = 9999
        treeData.value.forEach((_: any, index: number) => {
          if (treeData.value[index].name == pUrl) _index = index
          if (treeData.value[index].pUrl == pUrl && treeData.value[index].type == 0) {
            arr.value.push(treeData.value[index].checked)
          }
        })
        const flag = arr.value.every((item) => {
          return item
        })
        treeData.value[_index].checked = flag
        handleCheckedNodes()
      }, 0)
    }
  }

  // 输入框搜索
  const handleSearch = () => {
    if (!treeData.value) return
    isLoading.value = true
    setTimeout(() => {
      const urlArr: string[] = []
      if (treeData.value.length > 0 && siteSearch.value?.trim()) {
        treeData.value.forEach((item: { fullUrl: string; disabled: string; pUrl: string }) => {
          if (
            !urlArr.includes(item.pUrl) &&
            item.pUrl &&
            item.fullUrl.toUpperCase().includes(siteSearch.value?.toUpperCase())
          ) {
            urlArr.push(item.pUrl)
          }
          if (item.fullUrl.toUpperCase().includes(siteSearch.value?.toUpperCase())) {
            urlArr.push(item.fullUrl)
          }
        })
        treeData.value.forEach((item: { fullUrl: string; disabled: string; type: number }) => {
          item.disabled = urlArr.includes(item.fullUrl) ? 'flex' : 'none'
        })
      }
      isLoading.value = false
    }, 300)
  }

  // 筛选告警单选多选切换
  const handleAlternativeTypeChange = () => {
    treeData.value.forEach((item: { multiple: boolean; checked: boolean }) => {
      item.multiple = alternativeType.value
      item.checked = false
    })
  }

  const handleRadioClick = (idx: number) => {
    treeData.value.forEach((item: { checked: boolean }) => {
      item.checked = false
    })
    treeData.value[idx].checked = true
    handleCheckedNodes()
  }

  const handleToggleActive = (remark: string, idx: any, e: any) => {
    const num = activeArr.value[1] == idx && activeArr.value[0] == remark ? -1 : idx
    activeArr.value = [remark, num]
    const viewContent = document.querySelector('.content .el-tabs__content')
    if (!viewContent) return
    const { x, y } = viewContent.getBoundingClientRect()
    const { clientX, clientY } = e
    const obj = {
      offsetX: clientX - x,
      offsetY: clientY - y,
    }
    getThumbnailInfo(obj)
  }

  const handleToggleActiveResp = (_idx: any, idx: any, e: any) => {
    const num = activeResp.value[1] == idx && activeResp.value[0] == _idx ? -1 : idx
    activeResp.value = [_idx, num]
    const viewContent = document.querySelector('.content .el-tabs__content')
    if (!viewContent) return
    const { x, y } = viewContent.getBoundingClientRect()
    const { clientX, clientY } = e
    const obj = {
      offsetX: clientX - x,
      offsetY: clientY - y,
    }
    getThumbnailInfo(obj)
  }

  const handleMouseEnter = (title: any, subheading?: string, _index?: number) => {
    hoverValueArr.value[0] = title
    hoverValueArr.value[1] = subheading
    hoverValueArr.value[2] = _index
    hoverValueArr.value = hoverValueArr.value.filter((item) => {
      return item !== undefined
    })
  }

  const handleMouseLeave = () => {
    hoverValueArr.value = []
  }

  // 复制功能
  const handleCopyEvent = () => {
    if (hoverValueArr.value.length == 1 && hoverValueArr.value[0] == 'query') {
      formatDataWhenHoverValueArrLengthIsOne(queryPData.value)
    } else if (hoverValueArr.value.length == 1 && hoverValueArr.value[0] == 'head') {
      formatDataWhenHoverValueArrLengthIsOne(headerData.value)
    } else if (
      hoverValueArr.value.length == 1 &&
      (hoverValueArr.value[0] !== 'query' || hoverValueArr.value[0] !== 'head')
    ) {
      const item = responseData.value.find((item: { code: string }) => {
        return item.code == hoverValueArr.value[0]
      })
      formatDataWhenHoverValueArrLengthIsOne(item)
    } else if (hoverValueArr.value.length == 2 && hoverValueArr.value[0] == 'query') {
      formatDataWhenHoverValueArrLengthIsTwo(queryPData.value)
    } else if (hoverValueArr.value.length == 2 && hoverValueArr.value[0] == 'head') {
      formatDataWhenHoverValueArrLengthIsTwo(headerData.value)
    } else if (
      hoverValueArr.value.length == 2 &&
      (hoverValueArr.value[0] !== 'query' || hoverValueArr.value[0] !== 'head')
    ) {
      const item = responseData.value.find((item: { code: string }) => {
        return item.code == hoverValueArr.value[0]
      })
      formatDataWhenHoverValueArrLengthIsTwo(item)
    } else if (hoverValueArr.value.length == 3 && hoverValueArr.value[0] == 'query') {
      formatDataWhenHoverValueArrLengthIsThree(queryPData.value)
    } else if (hoverValueArr.value.length == 3 && hoverValueArr.value[0] == 'head') {
      formatDataWhenHoverValueArrLengthIsThree(headerData.value)
    } else if (
      hoverValueArr.value.length == 3 &&
      (hoverValueArr.value[0] !== 'query' || hoverValueArr.value[0] !== 'head')
    ) {
      const item = responseData.value.find((item: { code: string }) => {
        return item.code == hoverValueArr.value[0]
      })
      formatDataWhenHoverValueArrLengthIsThree(item)
    }
  }

  const formatDataWhenHoverValueArrLengthIsOne = (rightData: any) => {
    const { properties } = rightData
    const myObject = {}
    properties.forEach((item: any) => {
      const obj = {
        key: '',
        value: '',
      }
      obj.key = item.name
      let str = ''
      item.connection.forEach((td: any) => {
        if (td.value) {
          const str1 = td.value.join(',')
          if (str1) {
            str = `${str + str1};`
          }
        }
      })
      str = str.substring(0, str.length - 1)
      obj.value = str
      // @ts-ignore
      myObject[item.name] = str
    })
    rightJson.value = myObject
    handleCopy()
  }

  const formatDataWhenHoverValueArrLengthIsTwo = (rightData: any) => {
    const { properties } = rightData
    const obj = {
      key: '',
      value: '',
    }
    let str = ''
    obj.key = hoverValueArr.value[1]
    const item = properties.find((item: { name: string }) => {
      return item.name == hoverValueArr.value[1]
    })
    item.connection.forEach((td: any) => {
      if (td.value) {
        const str1 = td.value.join(',')
        if (str1) {
          str = `${str + str1};`
        }
      }
    })
    str = str.substring(0, str.length - 1)
    obj.value = str
    rightJson.value = obj
    handleCopy()
  }

  const formatDataWhenHoverValueArrLengthIsThree = (rightData: any) => {
    const { properties } = rightData
    const index = hoverValueArr.value[2]
    const obj = {
      type: '',
      value: '',
      description: '',
    }
    let str = ''
    const item = properties.find((item: { name: string }) => {
      return item.name == hoverValueArr.value[1]
    })
    const { one_of, value, description } = item.connection[index]
    const str1 = value.join(',')
    str = `${str + str1};`
    str = str.substring(0, str.length - 1)
    obj.value = str
    obj.type = one_of
    obj.description = description
    rightJson.value = obj
    handleCopy()
  }

  // 上下，左右拖动功能
  const horizontalMove = ref(false)
  const verticalMove = ref(false)
  const target = ref(null)
  const inside = ref(null)
  const topHeight = ref<number | string>()
  const leftWidth = ref<number | string>(360)

  const { elementY, elementX } = useMouseInElement(target)

  const { isOutside } = useMouseInElement(inside)

  const { height, width } = useElementBounding(inside)

  const handleDown = (e: any) => {
    e.preventDefault()
    horizontalMove.value = true
  }
  const handlVerticalNodeDown = (e: any) => {
    e.preventDefault()
    verticalMove.value = true
  }

  const handleUp = (e: any) => {
    horizontalMove.value = false
    canDragRes.value = false
    canDragResp.value = false
    canDragView.value = false
  }
  const handleVerticalNodeUp = (e: any) => {
    verticalMove.value = false
    canDragRes.value = false
    canDragResp.value = false
  }

  const horizontalNode = ref()
  const verticalNode = ref()
  const initDragNode = () => {
    horizontalNode.value = document.querySelector('.divider-horizontal')
    horizontalNode.value?.addEventListener('mousedown', handleDown)

    verticalNode.value = document.querySelector('.divider-vertical')
    verticalNode.value?.addEventListener('mousedown', handlVerticalNodeDown)

    window.addEventListener('mouseup', handleVerticalNodeUp)
    window.addEventListener('mouseup', handleUp)
  }

  onUnmounted(() => {
    horizontalNode.value?.removeEventListener('mousedown', handleDown)
    verticalNode.value?.removeEventListener('mousedown', handlVerticalNodeDown)
    window.removeEventListener('mouseup', handleUp)
    window.removeEventListener('mouseup', handleVerticalNodeUp)
  })

  watch(
    () => isOutside.value,
    () => {
      if (horizontalMove.value && isOutside.value) {
        horizontalMove.value = false
      }
      if (verticalMove.value && isOutside.value) {
        verticalMove.value = false
      }
      // if (canDragRes.value && isOutside.value) {
      //   canDragRes.value = false
      // }
      // if (canDragResp.value && isOutside.value) {
      //   canDragResp.value = false
      // }
      if (canDragView.value && isOutside.value) {
        canDragView.value = false
      }
    }
  )

  const getleftWidth = _lodash.throttle(() => {
    if (verticalMove.value) {
      if (elementX.value >= width.value * 0.7) {
        leftWidth.value = width.value * 0.7
      } else if (elementX.value < 360) {
        leftWidth.value = 360
      } else {
        leftWidth.value = elementX.value
      }
      getViewingAreaHeightAndWidthZoom()
      return leftWidth.value
    }
  }, 50)

  const getTopHeight = _lodash.throttle(() => {
    if (horizontalMove.value && nodes) {
      const node = nodes![0] as HTMLElement
      if (elementY.value >= height.value * 0.7) {
        topHeight.value = height.value * 0.7
      } else if (elementY.value < height.value * 0.3) {
        topHeight.value = height.value * 0.3
      } else {
        topHeight.value = elementY.value
      }
      const val = height.value - topHeight.value - 3
      node.style.setProperty('max-height', `${val - 100}px`)
      node.style.setProperty('min-height', `${val - 100}px`)
      getViewingAreaHeightAndWidthZoom()
      return topHeight.value
    }
  }, 50)

  watch(
    () => elementX.value,
    () => {
      getleftWidth()
    }
  )
  watch(
    () => elementY.value,
    () => {
      getTopHeight()
    }
  )

  // 缩放功能
  const resNode = ref()
  const respNode = ref()
  const initWheelScroll = () => {
    setTimeout(() => {
      resNode.value = document.querySelector('.request-content')
      resNode.value?.addEventListener('wheel', resNodeWheelScrollEvent)
      resNode.value?.addEventListener('mousedown', resNodeMousedown, false)
      resNode.value?.addEventListener('mouseup', resNodeMouseup, false)
      resNode.value?.addEventListener('mousemove', resNodeMousemove, false)

      respNode.value = document.querySelector('.response-content')
      respNode.value?.addEventListener('wheel', respNodeWheelScrollEvent)
      respNode.value?.addEventListener('mousedown', respNodeMousedown, false)
      respNode.value?.addEventListener('mouseup', respNodeMouseup, false)
      respNode.value?.addEventListener('mousemove', respNodeMousemove, false)

      getThumbnailInfo()
    }, 0)
  }

  const resNodeWheelScrollEvent = _lodash.throttle((event: any) => {
    // if (ctxZoom.value < 0.05) return
    if (event.deltaY < 0) {
      // 向上滚动，放大页面
      scaleResNum.value = scaleResNum.value >= 2 ? 2 : scaleResNum.value + 0.1
    } else {
      // 向下滚动，缩小页面
      scaleResNum.value = scaleResNum.value - 0.1 <= 0.3 ? 0.3 : scaleResNum.value - 0.1
    }
    resNode.value.style.transform = `translate3d(${resCurX}px, ${resCurY}px, 0) scale(${scaleResNum.value})`
    getViewingAreaHeightAndWidthZoom()
  }, 50)

  const respNodeWheelScrollEvent = _lodash.throttle((event: any) => {
    // if (ctxZoom.value < 0.05) return
    // scaleRespNum
    if (event.deltaY < 0) {
      // 向上滚动，放大页面
      scaleRespNum.value = scaleRespNum.value >= 2 ? 2 : scaleRespNum.value + 0.1
    } else {
      // 向下滚动，缩小页面
      scaleRespNum.value = scaleRespNum.value - 0.1 <= 0.3 ? 0.3 : scaleRespNum.value - 0.1
    }
    respNode.value.style.transform = `translate3d(${respCurX}px, ${respCurY}px, 0) scale(${scaleRespNum.value})`
    getViewingAreaHeightAndWidthZoom()
  }, 50)

  const initialisationWheelScrollEvent = () => {
    if (resNode.value) resNode.value.style.transform = `scale(${scaleResNum.value})`
    if (respNode.value) respNode.value.style.transform = `scale(${scaleRespNum.value})`
    // if (resNode.value || respNode.value) getViewingAreaHeightAndWidthZoom()
  }

  onUnmounted(() => {
    resNode.value?.removeEventListener('wheel', resNodeWheelScrollEvent)
    resNode.value?.removeEventListener('mousedown', resNodeMousedown)
    resNode.value?.removeEventListener('mouseup', resNodeMouseup)
    resNode.value?.removeEventListener('mousemove', resNodeMousemove)
    respNode.value?.removeEventListener('wheel', respNodeWheelScrollEvent)
    respNode.value?.removeEventListener('mousedown', respNodeMousedown)
    respNode.value?.removeEventListener('mouseup', respNodeMouseup)
    respNode.value?.removeEventListener('mousemove', respNodeMousemove)
    viewNode.value?.removeEventListener('mousedown', viewMousedown, false)
    viewNode.value?.removeEventListener('mouseup', viewMouseUP, false)
    viewNode.value?.removeEventListener('mousemove', viewNodeMousemove, false)
  })

  // 思维导图拖拽功能
  const resNodeMousedown = (e: any) => {
    resPosX = e.clientX - resOffsetX
    resPosY = e.clientY - resOffsety
    resViewClickX = e.clientX - resViewOffsetX
    resViewClickY = e.clientY - resViewOffsetY
    canDragRes.value = true
  }
  const resNodeMouseup = () => {
    resPosX = resCurX
    resPosY = resCurY
    resViewClickX = resViewEndX
    resViewClickY = resViewEndY
    canDragRes.value = false
  }
  const resNodeMousemove = _lodash.throttle((e: any) => {
    if (canDragRes.value) {
      const num = activeName.value == 'first' ? scaleResNum.value : scaleRespNum.value
      e.preventDefault()
      resCurX = e.clientX - resPosX
      resCurY = e.clientY - resPosY
      resOffsetX = resCurX
      resOffsety = resCurY
      resViewEndX = e.clientX - resViewClickX
      resViewEndY = e.clientY - resViewClickY
      resViewOffsetX = resViewEndX
      resViewOffsetY = resViewEndY
      setTranslate(resCurX, resCurY, resNode.value)
      setViewTranslate(resViewEndX / num, resViewEndY / num)
    }
  }, 10)

  const respNodeMousedown = (e: any) => {
    respPosX = e.clientX - respOffsetX
    respPosY = e.clientY - respOffsety
    respViewClickX = e.clientX - respViewOffsetX
    respViewClickY = e.clientY - respViewOffsetY
    canDragResp.value = true
  }
  const respNodeMouseup = () => {
    respPosX = respCurX
    respPosY = respCurY
    respViewClickX = respViewEndX
    respViewClickY = respViewEndY
    canDragResp.value = false
  }
  const respNodeMousemove = _lodash.throttle((e: any) => {
    const num = activeName.value == 'first' ? scaleResNum.value : scaleRespNum.value
    if (canDragResp.value) {
      respCurX = e.clientX - respPosX
      respCurY = e.clientY - respPosY
      respOffsetX = respCurX
      respOffsety = respCurY
      respViewEndX = e.clientX - respViewClickX
      respViewEndY = e.clientY - respViewClickY
      respViewOffsetX = respViewEndX
      respViewOffsetY = respViewEndY
      setTranslate(respCurX, respCurY, respNode.value)
      setViewTranslate(respViewEndX / num, respViewEndY / num)
    }
  }, 10)

  const setTranslate = (xPos: number, yPos: number, el: HTMLDivElement) => {
    const num = activeName.value == 'first' ? scaleResNum.value : scaleRespNum.value
    el.style.transform = `translate3d(${xPos}px, ${yPos}px, 0) scale(${num})`
    // getViewingAreaHeightAndWidthZoom()
  }
  const setViewTranslate = (xPos: number, yPos: number) => {
    const offsetX = activeName.value == 'first' ? offsetLeftRes : offsetLeftResp
    const offsetY = activeName.value == 'first' ? offsetTopRes : offsetTopResp
    // const num = activeName.value == 'first' ? scaleResNum.value : scaleRespNum.value
    const viewNode = document.querySelector('.viewing-area') as HTMLDivElement
    if (!viewNode) return
    viewNode.style.transform = `translate3d(${-xPos + offsetX}px, ${-yPos + offsetY}px, 0)`
  }

  const Loading = ref(false)
  // 缩略图功能
  const getThumbnailInfo = async (obj?: any) => {
    const htmlContent = activeName.value == 'first' ? resNode.value : respNode.value
    const Y = activeName.value == 'first' ? resViewEndY : respViewEndY
    const X = activeName.value == 'first' ? resViewEndX : respViewEndX
    const viewContent = document.querySelector('.content .el-tabs__content')
    if (!viewContent) return
    const contentHeight = viewContent.getBoundingClientRect().height - 50
    const contentWidth = viewContent.getBoundingClientRect().width - 100
    if (Y > contentHeight || -Y > contentHeight || X > contentWidth || -X > contentWidth) {
      Loading.value = true
    }
    await htmlToCanvas(htmlContent, obj)
  }

  const htmlToCanvas = _lodash.throttle(async (htmlContent: HTMLElement, obj?: any) => {
    setTimeout(async () => {
      if (!htmlContent) return
      const { width, height } = getMindMapHeightAndWeight()
      const num = activeName.value == 'first' ? scaleResNum.value : scaleRespNum.value
      const canvas = await html2canvas(htmlContent, { scale: 1 / num, width, height })
      if (canvas) {
        const canvasWidth = canvas.getAttribute('width')!
        const canvasHeight = canvas.getAttribute('height')!
        canvas.setAttribute('style', `width: ${canvasWidth};height: ${canvasHeight}`)
        const node = document.querySelector('.canvas-map') as HTMLDivElement
        if (node) {
          const nodeWidth = canvasWidth
          const nodeHeight = canvasHeight
          node.style.width = `${nodeWidth}px`
          node.style.height = `${nodeHeight}px`
          const zoom = 180 / +nodeWidth < 92 / +nodeHeight ? 180 / +nodeWidth : 92 / +nodeHeight

          ctxZoom.value = zoom
          node.style.transform = `translate3d(-50%, -50%, 0) scale(${zoom})`
          const viewing = document.querySelector('.viewing-area')
          viewing?.setAttribute('style', 'display: none')
          setTimeout(() => {
            const child = document.getElementsByTagName('canvas')[0]
            child && node.removeChild(child)
            node?.appendChild(canvas)
            viewing?.removeAttribute('style')
            getViewingAreaHeightAndWidthZoom()
            initViewNode()
            if (obj) {
              if (!viewing) return
              const { x, y } = viewing.getBoundingClientRect()
              const newObj = {
                clientX: (obj.offsetX * zoom) / num + x,
                clientY: (obj.offsetY * zoom) / num + y,
              }
              const Y = activeName.value == 'first' ? resViewEndY : respViewEndY
              const X = activeName.value == 'first' ? resViewEndX : respViewEndX
              const viewContent = document.querySelector('.content .el-tabs__content')
              if (!viewContent) return
              const contentHeight = viewContent.getBoundingClientRect().height - 50
              const contentWidth = viewContent.getBoundingClientRect().width - 100
              if (Y > contentHeight || -Y > contentHeight || X > contentWidth || -X > contentWidth) {
                handleMapClick(newObj)
              }
            }
            Loading.value = false
          }, 0)
        }
      }
    }, 0)
  })

  // 获取思维导图宽高
  const getMindMapHeightAndWeight = () => {
    const htmlContent = activeName.value == 'first' ? resNode.value : respNode.value
    if (htmlContent) {
      const column1 = htmlContent.querySelectorAll('.column-1-title')
      const column2 = htmlContent.querySelectorAll('.column-1-item')
      let width1 = 0
      let height1 = 0
      let width2 = 0
      let height2 = 0
      column1?.forEach((item: any) => {
        width1 = item.getBoundingClientRect().width > width1 ? item.getBoundingClientRect().width : width1
        height1 = item.getBoundingClientRect().height
      })
      column2?.forEach((item: any) => {
        width2 = item.getBoundingClientRect().width > width2 ? item.getBoundingClientRect().width : width2
        height2 += item.getBoundingClientRect().height
      })
      const height3 = (column1.length - column2.length) * height1
      const width = +(width1 + width2 + 200).toFixed(0)
      const height = +(height3 + height2 + 45).toFixed(0)

      return { width, height }
    } else {
      return { width: 0, height: 0 }
    }
  }

  // 视觉区域比例
  const getViewingAreaHeightAndWidthZoom = () => {
    // 影响因素：缩略图思维导图大小；思维导图缩放比；视口大小
    // // 1.缩略图思维导图大小
    // const canvasMap = document.querySelector('.canvas-map canvas')
    // if (!canvasMap) return
    // const width = +canvasMap.getAttribute('width')!
    // const height = +canvasMap.getAttribute('height')!
    // 2.思维导图缩放比
    const zoom = activeName.value == 'first' ? scaleResNum.value : scaleRespNum.value
    // 3.视口大小
    const viewContent = document.querySelector('.content .el-tabs__content')
    if (!viewContent) return
    const contentHeight = viewContent.getBoundingClientRect().height
    const contentWidth = viewContent.getBoundingClientRect().width
    // 在缩略图中canvas的大小在思维导图缩放过程中大小是不变的，所以变化的视觉区域
    // 思维导图放大，可以看成视觉区域缩小
    const viewAreaHeight = contentHeight / zoom
    // const viewAreaHeight = contentHeight / zoom > height ? height : contentHeight / zoom
    // const viewAreaWidth = contentWidth / zoom > width ? width : contentWidth / zoom
    const viewAreaWidth = contentWidth / zoom
    const viewing = document.querySelector('.viewing-area')
    viewing?.setAttribute('style', `height:${viewAreaHeight}px;width:${viewAreaWidth}px`)

    getViewNodePosition()
  }

  const getScaleNum = () => {
    return (activeName.value == 'first' ? scaleResNum.value * 100 : scaleRespNum.value * 100).toFixed(0)
  }

  const initPosition = () => {
    if (activeName.value == 'first') {
      ininitialisationResDragData()
      ininitialisationViewNodePositionRes()
      initResViewPosition()
      scaleResNum.value = 1
      resNode.value.style.transform = `translate3d(0px, 0px, 0) scale(${1})`
    } else {
      ininitialisationRespDragData()
      ininitialisationViewNodePositionResp()
      initRespViewPosition()
      scaleRespNum.value = 1
      respNode.value.style.transform = `translate3d(0px, 0px, 0) scale(${1})`
    }
    getViewingAreaHeightAndWidthZoom()
    getThumbnailInfo()
  }

  const viewNode = ref()
  const toggleMapShow = () => {
    showMap.value = !showMap.value
    if (showMap.value) {
      const htmlContent = activeName.value == 'first' ? resNode.value : respNode.value
      htmlToCanvas(htmlContent)
      initViewNode()
    }
  }

  const initViewNode = () => {
    setTimeout(() => {
      viewNode.value = document.querySelector('.viewing-area')
      viewNode.value?.addEventListener('mousedown', viewMousedown, false)
      viewNode.value?.addEventListener('mouseup', viewMouseUP, false)
      viewNode.value?.addEventListener('mousemove', viewNodeMousemove, false)
    }, 0)
  }

  const viewMousedown = (e: any) => {
    e.preventDefault()
    e.stopPropagation()
    const num = activeName.value == 'first' ? scaleResNum.value : scaleRespNum.value
    if (activeName.value == 'first') {
      resViewClickX = e.clientX + resViewOffsetX * ctxZoom.value
      resViewClickY = e.clientY + resViewOffsetY * ctxZoom.value
      resPosX = e.clientX + (resOffsetX * ctxZoom.value) / num
      resPosY = e.clientY + (resOffsety * ctxZoom.value) / num
    } else {
      respViewClickX = e.clientX + respViewOffsetX * ctxZoom.value
      respViewClickY = e.clientY + respViewOffsetY * ctxZoom.value
      respPosX = e.clientX + (respOffsetX * ctxZoom.value) / num
      respPosY = e.clientY + (respOffsety * ctxZoom.value) / num
    }

    canDragView.value = true
  }

  const viewMouseUP = () => {
    if (activeName.value == 'first') {
      resViewClickX = resViewEndX
      resViewClickY = resViewEndY
      resPosX = resCurX
      resPosY = resCurY
    } else {
      respViewClickX = respViewEndX
      respViewClickY = respViewEndY
      respPosX = respCurX
      respPosY = respCurY
    }
    canDragView.value = false
  }

  const viewNodeMousemove = _lodash.throttle((e: any) => {
    if (canDragView.value) {
      // 首先要初始化位置
      const viewNode = document.querySelector('.viewing-area') as HTMLDivElement
      if (!viewNode) return
      // const { width, height } = getMindMapHeightAndWeight()
      // const zoom = 180 / width < 92 / height ? 180 / width : 92 / height
      const num = activeName.value == 'first' ? scaleResNum.value : scaleRespNum.value
      // Request tab签
      if (activeName.value == 'first') {
        resCurX = 0 - ((e.clientX - resPosX) / ctxZoom.value) * num
        resCurY = 0 - ((e.clientY - resPosY) / ctxZoom.value) * num
        resOffsetX = resCurX
        resOffsety = resCurY

        resViewEndX = 0 - (e.clientX - resViewClickX) / ctxZoom.value
        resViewEndY = 0 - (e.clientY - resViewClickY) / ctxZoom.value
        resViewOffsetX = resViewEndX
        resViewOffsetY = resViewEndY

        setTranslate(resCurX, resCurY, resNode.value)
        setViewTranslate(resViewEndX, resViewEndY)
      } else {
        respCurX = 0 - ((e.clientX - respPosX) / ctxZoom.value) * num
        respCurY = 0 - ((e.clientY - respPosY) / ctxZoom.value) * num
        respOffsetX = respCurX
        respOffsety = respCurY

        respViewEndX = 0 - (e.clientX - respViewClickX) / ctxZoom.value
        respViewEndY = 0 - (e.clientY - respViewClickY) / ctxZoom.value
        respViewOffsetX = respViewEndX
        respViewOffsetY = respViewEndY

        setTranslate(respCurX, respCurY, respNode.value)
        setViewTranslate(respViewEndX, respViewEndY)
      }
    }
  }, 50)

  const getViewNodePosition = () => {
    const viewNode = document.querySelector('.viewing-area') as HTMLDivElement
    if (!viewNode) return
    const { width, height } = viewNode.getBoundingClientRect()
    // 思维导图大小
    const num = activeName.value == 'first' ? scaleResNum.value : scaleRespNum.value
    // 计算偏移
    if (activeName.value == 'first') {
      offsetLeftRes = num >= 1 ? (width - width / num) / ctxZoom.value : (width - width / num) * num * 2
      offsetTopRes = num >= 1 ? (height - height / num) / ctxZoom.value : (height - height / num) * num * 2
      setViewTranslate(respViewEndX, respViewEndY)
    } else {
      offsetLeftResp = num >= 1 ? (width - width / num) / ctxZoom.value : (width - width / num) * num * 2
      offsetTopResp = num >= 1 ? (height - height / num) / ctxZoom.value : (height - height / num) * num * 2
      setViewTranslate(respViewEndX, respViewEndY)
    }
  }

  const handleScale = (flag: boolean) => {
    // if (ctxZoom.value < 0.05) return
    if (activeName.value == 'first') {
      if (flag) {
        // 向上滚动，放大页面
        scaleResNum.value = scaleResNum.value >= 2 ? 2 : scaleResNum.value + 0.1
      } else {
        // 向下滚动，缩小页面
        scaleResNum.value = scaleResNum.value - 0.1 <= 0.3 ? 0.3 : scaleResNum.value - 0.1
      }
      resNode.value.style.transform = `translate3d(${resCurX}px, ${resCurY}px, 0) scale(${scaleResNum.value})`
      getViewingAreaHeightAndWidthZoom()
    } else {
      if (flag) {
        // 向上滚动，放大页面
        scaleRespNum.value = scaleRespNum.value >= 2 ? 2 : scaleRespNum.value + 0.1
      } else {
        // 向下滚动，缩小页面
        scaleRespNum.value = scaleRespNum.value - 0.1 <= 0.3 ? 0.3 : scaleRespNum.value - 0.1
      }
      respNode.value.style.transform = `translate3d(${respCurX}px, ${respCurY}px, 0) scale(${scaleRespNum.value})`
      getViewingAreaHeightAndWidthZoom()
    }
  }

  // 点击画布位移
  const handleMapClick = (e: any) => {
    const canvasMap = document.querySelector('.canvas-map canvas')
    if (!canvasMap) return
    const { x, y } = canvasMap.getBoundingClientRect()
    const { width, height } = viewNode.value.getBoundingClientRect()
    const { clientX, clientY } = e
    const num = activeName.value == 'first' ? scaleResNum.value : scaleRespNum.value
    if (activeName.value == 'first') {
      resViewEndX = -((clientX - x - (width / 2) * num) / ctxZoom.value)
      resViewEndY = -((clientY - y - (height / 2) * num) / ctxZoom.value)
      resViewOffsetX = resViewEndX
      resViewOffsetY = resViewEndY
      // setViewTranslate(resViewEndX * num, resViewEndY * num)
      const viewNode = document.querySelector('.viewing-area') as HTMLDivElement
      if (!viewNode) return
      viewNode.style.transform = `translate3d(${-resViewEndX + offsetLeftRes}px, ${-resViewEndY + offsetTopRes}px, 0)`

      resCurX = resViewEndX * num
      resCurY = resViewEndY * num
      resOffsetX = resCurX
      resOffsety = resCurY
      setTranslate(resCurX, resCurY, resNode.value)
    } else {
      respViewEndX = -((clientX - x - width / 2) / ctxZoom.value)
      respViewEndY = -((clientY - y - height / 2) / ctxZoom.value)
      respViewOffsetX = respViewEndX
      respViewOffsetY = respViewEndY
      setViewTranslate(resViewEndX, resViewEndY)

      respCurX = resViewEndX * num
      respCurY = resViewEndY * num
      respOffsetX = respCurX
      respOffsety = respCurY
      setTranslate(respCurX, respCurY, respNode.value)
    }
  }

  // 切换tab
  const handleTabChange = (name: any) => {
    getThumbnailInfo()
  }

  const isFullscreen = ref(false)
  // 全屏的开关
  const clickFullScreen = () => {
    isFullscreen.value = !isFullscreen.value
  }
</script>
<template>
  <div ref="inside" class="api-list" :class="{ 'vab-fullscreen': isFullscreen }">
    <!-- <el-tooltip class="item" :content="isFullscreen ? '关闭全屏' : '全屏'" effect="dark" placement="top">
      <vab-icon
        class="icon"
        :icon="isFullscreen ? 'fullscreen-exit-fill' : 'fullscreen-fill'"
        @click="clickFullScreen"
      />
    </el-tooltip> -->
    <div ref="target" class="resize-top" :style="{ height: topHeight + 'px' }">
      <div class="left" style="height: 100%" :style="{ width: leftWidth + 'px' }">
        <div class="my-search">
          <input-search
            v-model="siteSearch"
            clearable
            placeholder="搜索"
            style="width: 160px"
            @on-search="handleSearch"
          />
          <div class="alternativeType">
            <el-checkbox v-if="showToggle" v-model="alternativeType" @change="handleAlternativeTypeChange" />
            <span v-if="showToggle" style="margin: 0 10px 0 0; color: #2b2742">多选</span>
            <el-button
              :disabled="listLoading || (treeData && treeData.length == 0)"
              size="small"
              type="primary"
              @click="handleToggle"
            >
              {{ showToggle ? '关闭筛选' : '筛选告警' }}
            </el-button>
          </div>
        </div>
        <div v-loading="isLoading" class="my-tree">
          <template v-for="(item, index) in treeData" :key="item.id">
            <div v-if="item.disabled == 'flex'" class="node node-item">
              <div
                v-if="item.type == 1"
                class="nodeIcon"
                :style="{ transform: item.export ? 'rotate(90deg)' : 'none' }"
                @click="toggleExpand(index)"
              >
                <el-icon style="color: #b2add4"><CaretRight /></el-icon>
              </div>
              <el-checkbox
                v-if="showToggle && item.multiple"
                v-model="item.checked"
                :style="{
                  marginLeft: item.type == 1 || item.name == item.fullUrl ? '0px' : '36px',
                  marginRight: '6px',
                  width: '12px',
                  height: '12px',
                }"
                @click="handleChange(index)"
              />
              <div
                v-if="item.type == 0 && !showToggle"
                :style="{
                  marginLeft: item.type == 1 || item.name == item.fullUrl ? '0px' : '36px',
                  marginRight: '6px',
                }"
              ></div>
              <div
                v-if="item.type == 0 && showToggle && !item.multiple"
                class="my-radio"
                :class="{ radio_active: item.checked }"
                :style="{
                  marginLeft: item.type == 1 || item.name == item.fullUrl ? '0px' : '36px',
                  marginRight: '4px',
                }"
                @click="handleRadioClick(index)"
              ></div>
              <el-image
                v-if="item.type == 1"
                :src="require('@/assets/site_images/site-icon.svg')"
                style="width: 14px; height: 12px; margin: 0 2px 0 6px"
              />
              <div class="circle" :style="{ background: item.hasWarn ? '#df5352' : '#4FB769' }"></div>
              <span
                class="node-txt checkOverflow"
                :class="{ nodeActive: item.id == currentData.id, isHover: item.type != 1 }"
                :style="{
                  maxWidth: item.type == 1 ? 'calc(100% - 91px)' : 'calc(100% - 100px)',
                }"
                @click="itemClick(item)"
              >
                <el-tooltip v-if="!item.overflow" :content="item.name" effect="dark" placement="top">
                  {{ item.name }}
                </el-tooltip>
                <span v-else>{{ item.name }}</span>
              </span>
              <span v-if="item.children >= 1" class="node-num">
                {{ item.children > 99 ? '99+' : item.children }}
              </span>
            </div>
          </template>
        </div>
      </div>
      <div
        class="divider-vertical"
        style="width: 2px; height: 100%; background-color: #e7e4fb; cursor: ew-resize"
      ></div>
      <div class="right">
        <div class="content">
          <div class="content-host">
            <div class="content-host-text">
              <span v-if="tagsList.length > 0">
                <span
                  v-for="(item, index) in tagsList"
                  :key="index"
                  class="tag"
                  :class="{ tagActive: itemIndex == index ? true : false }"
                  @click="tagsClick(index)"
                >
                  {{ item && item.toUpperCase() }}
                </span>
              </span>
              <span v-if="currentData && currentData.host" class="urlName">
                {{ currentData && currentData.host ? currentData.host : '/' }}
              </span>
            </div>
            <!-- <el-button
              v-if="currentData?.host"
              :loading="listLoading"
              size="small"
              style="width: 48px; height: 26px; background: #6954f0; border-radius: 4px"
              type="primary"
              @click="exportEvent"
            >
              导出
            </el-button> -->
          </div>
          <div v-if="currentData && currentData.host" class="content-view">
            <el-tabs v-model="activeName" v-loading="Loading" class="demo-tabs" @tab-change="handleTabChange">
              <el-tab-pane label="Request" name="first">
                <div v-if="queryPData || headerData" class="column request-content" style="cursor: move">
                  <div class="column-1">
                    <div class="column-1-title" @mouseenter="handleMouseEnter('query')" @mouseleave="handleMouseLeave">
                      Query Parameters
                      <el-icon
                        v-if="
                          hoverValueArr[0] == 'query' &&
                          !hoverValueArr[1] &&
                          !hoverValueArr[2] &&
                          queryPData &&
                          queryPData.properties.length > 0
                        "
                        class="documentCopy"
                        style="margin-left: 5px; cursor: pointer"
                        @click.stop="handleCopyEvent"
                      >
                        <DocumentCopy />
                      </el-icon>
                      <div
                        v-if="(queryPData && queryPData.properties.length == 0) || (dataList && dataList.length == 0)"
                        class="mark"
                      ></div>
                    </div>
                    <div v-if="queryPData && queryPData.properties.length > 0" class="column-1-item">
                      <div v-for="(item, index) in queryPData.properties" :key="index" class="column-2">
                        <div
                          class="column-2-title"
                          @click.stop="(e) => handleToggleActive('query', index, e)"
                          @mouseenter="handleMouseEnter('query', item.name)"
                          @mouseleave="handleMouseLeave"
                        >
                          {{ item.name }}
                          <el-icon
                            v-if="
                              hoverValueArr[0] == 'query' &&
                              hoverValueArr.length == 2 &&
                              hoverValueArr[1] == item.name &&
                              !hoverValueArr[2] &&
                              item.connection.length > 0
                            "
                            class="documentCopy"
                            style="margin-left: 5px; cursor: pointer"
                            @click.stop="handleCopyEvent"
                          >
                            <DocumentCopy />
                          </el-icon>
                          <div
                            v-if="index == activeArr[1] && activeArr[0] == 'query' && item.connection.length > 0"
                            class="btn"
                            @click.stop="(e) => handleToggleActive('query', index, e)"
                          >
                            -
                          </div>
                        </div>
                        <div class="column-3">
                          <div v-for="(td, e) in item.connection" :key="e" class="column-3-item">
                            <div
                              v-if="index != activeArr[1] || activeArr[0] !== 'query'"
                              class="btn-num"
                              @click.stop="(e) => handleToggleActive('query', index, e)"
                            >
                              {{ item.connection.length }}
                            </div>
                            <div
                              v-if="index == activeArr[1] && activeArr[0] == 'query'"
                              class="column-3-item-title"
                              @mouseenter="handleMouseEnter('query', item.name, e)"
                              @mouseleave="handleMouseLeave"
                            >
                              <div class="content-warp">
                                <div class="content-warp-title">Type:</div>
                                <div class="content-value">
                                  <div
                                    style="
                                      display: flex;
                                      align-items: center;
                                      justify-content: space-between;
                                      width: 100%;
                                    "
                                  >
                                    {{ td.one_of }}
                                    <el-icon
                                      v-if="
                                        hoverValueArr[0] == 'query' &&
                                        hoverValueArr[1] == item.name &&
                                        hoverValueArr[2] == e
                                      "
                                      style="cursor: pointer"
                                      @click.stop="handleCopyEvent"
                                    >
                                      <DocumentCopy />
                                    </el-icon>
                                  </div>
                                </div>
                              </div>
                              <div class="content-warp">
                                <div class="content-warp-title">Value:</div>
                                <div class="content-value content-val-sp">
                                  <div v-for="(t, _index) in td.value" :key="_index" class="value">{{ `"${t}"` }}</div>
                                </div>
                              </div>
                              <div class="content-warp">
                                <div class="content-warp-title">Description:</div>
                                <div class="content-value">
                                  <div>{{ (td && td.description) || '-' }}</div>
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="column-1">
                    <div class="column-1-title" @mouseenter="handleMouseEnter('head')" @mouseleave="handleMouseLeave">
                      Header Parameters
                      <el-icon
                        v-if="
                          hoverValueArr[0] == 'head' &&
                          !hoverValueArr[1] &&
                          !hoverValueArr[2] &&
                          headerData &&
                          headerData.properties.length > 0
                        "
                        class="documentCopy"
                        style="margin-left: 5px; cursor: pointer"
                        @click.stop="handleCopyEvent"
                      >
                        <DocumentCopy />
                      </el-icon>
                      <div
                        v-if="(headerData && headerData.properties.length == 0) || (dataList && dataList.length == 0)"
                        class="mark"
                      ></div>
                    </div>
                    <div v-if="headerData && headerData.properties.length > 0" class="column-1-item">
                      <div v-for="(item, index) in headerData.properties" :key="index" class="column-2">
                        <div
                          class="column-2-title"
                          @click.stop="(e) => handleToggleActive('head', index, e)"
                          @mouseenter="handleMouseEnter('head', item.name)"
                          @mouseleave="handleMouseLeave"
                        >
                          {{ item.name }}
                          <el-icon
                            v-if="
                              hoverValueArr[0] == 'head' &&
                              hoverValueArr.length == 2 &&
                              hoverValueArr[1] == item.name &&
                              !hoverValueArr[2] &&
                              item.connection.length > 0
                            "
                            class="documentCopy"
                            style="margin-left: 5px; cursor: pointer"
                            @click.stop="handleCopyEvent"
                          >
                            <DocumentCopy />
                          </el-icon>
                          <div
                            v-if="index == activeArr[1] && activeArr[0] == 'head' && item.connection.length > 0"
                            class="btn"
                            @click.stop="(e) => handleToggleActive('head', index, e)"
                          >
                            -
                          </div>
                        </div>
                        <div class="column-3">
                          <div v-for="(td, e) in item.connection" :key="e" class="column-3-item">
                            <div
                              v-if="index != activeArr[1] || activeArr[0] !== 'head'"
                              class="btn-num"
                              @click.stop="(e) => handleToggleActive('head', index, e)"
                            >
                              {{ item.connection.length }}
                            </div>
                            <div
                              v-if="index == activeArr[1] && activeArr[0] == 'head'"
                              class="column-3-item-title"
                              @mouseenter="handleMouseEnter('head', item.name, e)"
                              @mouseleave="handleMouseLeave"
                            >
                              <div class="content-warp">
                                <div class="content-warp-title">Type:</div>
                                <div class="content-value">
                                  <div
                                    style="
                                      display: flex;
                                      align-items: center;
                                      justify-content: space-between;
                                      width: 100%;
                                    "
                                  >
                                    {{ td.one_of }}
                                    <el-icon
                                      v-if="
                                        hoverValueArr[0] == 'head' &&
                                        hoverValueArr[1] == item.name &&
                                        hoverValueArr[2] == e
                                      "
                                      style="cursor: pointer"
                                      @click.stop="handleCopyEvent"
                                    >
                                      <DocumentCopy />
                                    </el-icon>
                                  </div>
                                </div>
                              </div>
                              <div class="content-warp">
                                <div class="content-warp-title">Value:</div>
                                <div class="content-value content-val-sp">
                                  <div v-for="(t, _index) in td.value" :key="_index" class="value">
                                    {{ `"${t}" ` }}
                                  </div>
                                </div>
                              </div>
                              <div class="content-warp">
                                <div class="content-warp-title">Description:</div>
                                <div class="content-value">
                                  <div>{{ (td && td.description) || '-' }}</div>
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </el-tab-pane>
              <el-tab-pane label="Response" name="second">
                <div v-if="responseData" class="column response-content" style="cursor: move">
                  <div
                    v-for="(warp, _index) in responseData"
                    :key="_index"
                    class="column-1"
                    :style="{ borderLeft: responseData.length == 1 ? 'none' : '1px solid #e7e4fb' }"
                  >
                    <div
                      class="column-1-title"
                      :class="{ myMark: responseData.length == 1 }"
                      @mouseenter="handleMouseEnter(warp.code)"
                      @mouseleave="handleMouseLeave"
                    >
                      {{ warp.code }}
                      <el-icon
                        v-if="
                          hoverValueArr[0] == warp.code &&
                          !hoverValueArr[1] &&
                          !hoverValueArr[2] &&
                          warp.properties.length > 0
                        "
                        class="documentCopy"
                        style="margin-left: 5px; cursor: pointer"
                        @click.stop="handleCopyEvent"
                      >
                        <DocumentCopy />
                      </el-icon>
                      <div v-if="warp.properties.length == 0" class="mark"></div>
                    </div>
                    <div v-if="warp.properties.length > 0" class="column-1-item">
                      <div v-for="(item, index) in warp.properties" :key="index" class="column-2">
                        <div
                          class="column-2-title"
                          warp.code
                          @click.stop="(e) => handleToggleActiveResp(_index, index, e)"
                          @mouseenter="handleMouseEnter(warp.code, item.name)"
                          @mouseleave="handleMouseLeave"
                        >
                          {{ item.name }}
                          <el-icon
                            v-if="
                              hoverValueArr[0] == warp.code &&
                              hoverValueArr.length == 2 &&
                              hoverValueArr[1] == item.name &&
                              !hoverValueArr[2] &&
                              item.connection.length > 0
                            "
                            class="documentCopy"
                            style="margin-left: 5px; cursor: pointer"
                            @click.stop="handleCopyEvent"
                          >
                            <DocumentCopy />
                          </el-icon>
                          <div
                            v-if="index == activeResp[1] && activeResp[0] == _index && item.connection.length > 0"
                            class="btn"
                            @click.stop="(e) => handleToggleActiveResp(_index, index, e)"
                          >
                            -
                          </div>
                        </div>
                        <div class="column-3">
                          <div v-for="(td, e) in item.connection" :key="e" class="column-3-item">
                            <div
                              v-if="index != activeResp[1] || activeResp[0] !== _index"
                              class="btn-num"
                              @click.stop="(e) => handleToggleActiveResp(_index, index, e)"
                            >
                              {{ item.connection.length }}
                            </div>
                            <div
                              v-if="index == activeResp[1] && activeResp[0] == _index"
                              class="column-3-item-title"
                              @mouseenter="handleMouseEnter(warp.code, item.name, e)"
                              @mouseleave="handleMouseLeave"
                            >
                              <div class="content-warp">
                                <div class="content-warp-title">Type:</div>
                                <div class="content-value">
                                  <div
                                    style="
                                      display: flex;
                                      align-items: center;
                                      justify-content: space-between;
                                      width: 100%;
                                    "
                                  >
                                    {{ td.one_of }}
                                    <el-icon
                                      v-if="
                                        hoverValueArr[0] == warp.code &&
                                        hoverValueArr[1] == item.name &&
                                        hoverValueArr[2] == e
                                      "
                                      style="cursor: pointer"
                                      @click.stop="handleCopyEvent"
                                    >
                                      <DocumentCopy />
                                    </el-icon>
                                  </div>
                                </div>
                              </div>
                              <div class="content-warp">
                                <div class="content-warp-title">Value:</div>
                                <div class="content-value content-val-sp">
                                  <div v-for="(t, _index) in td.value" :key="_index" class="value">{{ `"${t}"` }}</div>
                                </div>
                              </div>
                              <div class="content-warp">
                                <div class="content-warp-title">Description:</div>
                                <div class="content-value">
                                  <div>{{ (td && td.description) || '-' }}</div>
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </el-tab-pane>
            </el-tabs>
            <div v-if="queryPData || headerData" class="thumbnail">
              <div class="tools">
                <el-image
                  alt=""
                  :src="require('@/assets/site_images/map.svg')"
                  style="height: 14px; width: 14px; cursor: pointer"
                  @click="toggleMapShow"
                />
                <div class="control-keys">
                  <el-button
                    link
                    style="font-size: 18px; font-weight: 600; color: #938eaf; margin-bottom: 2px"
                    @click="handleScale(false)"
                  >
                    -
                  </el-button>
                  <div>{{ getScaleNum() }}%</div>
                  <el-button
                    link
                    style="font-size: 18px; font-weight: 600; color: #938eaf; margin-bottom: 2px"
                    @click="handleScale(true)"
                  >
                    +
                  </el-button>
                </div>
                <el-image
                  alt=""
                  :src="require('@/assets/site_images/position.svg')"
                  style="height: 14px; width: 14px; cursor: pointer"
                  @click="initPosition"
                />
              </div>
              <div v-if="showMap" class="canvas">
                <div class="canvas-map" @click="handleMapClick">
                  <div class="viewing-area" @click.stop></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div
      class="divider-horizontal"
      style="height: 3px; width: 100%; background-color: #e7e4fb; cursor: ns-resize"
    ></div>
    <div class="footer" style="width: 100%">
      <div class="time">
        <div style="font-size: 16px; color: #303338; line-height: 22px">告警数据</div>
        <el-form inline>
          <el-form-item>
            <el-select v-model="timeDuration" style="width: 180px">
              <el-option v-for="item in date_options" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item v-if="timeDuration === 'user-defined'" style="width: 380px">
            <vab-date-time-picker v-model="timeDate" style="margin-left: 20px" />
          </el-form-item>
        </el-form>
      </div>
      <el-table
        v-loading="loading"
        class="my-table"
        :data="alertInfo.list"
        style="width: 100%"
        @cell-contextmenu="useTableCopy"
        @selection-change="handleSelectionChange"
      >
        <el-table-column align="center" type="selection" width="55" />
        <el-table-column align="center" :index="curIndex" label="序号" type="index" width="55" />
        <el-table-column
          v-for="item in tableColumn"
          :key="item.id"
          align="center"
          :label="item.fieldNameCn"
          :prop="item.fieldNameEn"
          show-overflow-tooltip
          :width="['threatName', 'warnTime'].includes(item.fieldNameEn) ? 250 : 150"
        >
          <template #header>
            <span style="cursor: pointer">
              {{ item.fieldNameCn }}
              <vab-icon
                v-if="item.supportAgg"
                icon="filter-line"
                style="font-size: 12px; margin-left: 6px; vertical-align: -1px !important"
                @click="showTraceabilityField(item.fieldNameCn, item.fieldNameEn)"
              />
            </span>
          </template>
          <template #default="{ row }" v-if="item.fieldNameEn == 'startTimeNs'">
            {{ formatDataTime(row, item.fieldNameEn) }}
          </template>
          <template #default="{ row }" v-else-if="item.fieldNameEn == 'threatLevel'">
            <span :class="['alert_tag', getThreatLevel(row.threatLevel)]">
              <el-icon><WarnTriangleFilled /></el-icon>
              {{ row.threatLevel }}
            </span>
          </template>
          <template #default="{ row }" v-else-if="item.fieldNameEn == 'warnTime'">
            {{ formatDataTime(row, item.fieldNameEn) }}
          </template>
        </el-table-column>

        <el-table-column align="center" fixed="right" label="操作" width="150">
          <template #default="{ row }">
            <el-button class="row_action" size="small" @click="showAlertDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <alert-detail
      v-model:alert-detail-visible="alertDetailVisible"
      :module-enable="moduleEnable"
      :next="showNextBtn"
      :prev="showPrevBtn"
      :select-alert="infoVal"
      @on-skip-event="changeCurrentItemEvent"
    />
    <traceability-field-statistic ref="traceabilityFieldRef" />
  </div>
</template>

<style scoped lang="scss">
  $itemHeight: 48px;
  $criticalColor: #d60705;
  $lowColor: #f8ae0a;
  $midColor: #ff751f;
  $highColor: #ff3c3a;
  $defaultColor: #909399;
  $topHeight: calc((100vh - 85px) * 0.7);
  $img: 30px;
  .divider-vertical,
  .divider-horizontal {
    position: relative;
    &::after {
      position: absolute;
      content: '';
      top: 50%;
      transform: translateY(-50%);
      left: -8px;
      background: url('@/assets//site_images/drag.svg');
      width: 16px;
      height: 30px;
      z-index: 10;
    }
  }
  .divider-horizontal {
    &::after {
      left: 50%;
      transform: translateX(-50%) rotate(90deg);
      top: -14px;
    }
  }
  .my-table {
    position: relative;
    :deep() {
      .el-table__inner-wrapper {
        width: 100%;
      }
      height: 100%;
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
        max-height: calc(100vh - 100px - $topHeight - 72px);
        min-height: calc(100vh - 100px - $topHeight - 72px);
        // max-height: calc(100% - 100px);
        // min-height: calc(100% - 100px);
        overflow-y: auto;
        // .el-scrollbar__bar.is-horizontal {
        //   bottom: 10px;
        // }
        &::-webkit-scrollbar {
          width: 0px !important;
          height: 0px !important;
        }
      }
    }
  }
  .alert_tag {
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
      background-color: rgba($criticalColor, 1);
    }
    &.low {
      background-color: rgba($lowColor, 1);
    }
    &.mid {
      background-color: rgba($midColor, 1);
    }
    &.high {
      background-color: rgba($highColor, 1);
    }
    &.default {
      background-color: rgba($defaultColor, 1);
      color: #fff;
    }
  }

  .api-list {
    width: 100%;
    height: calc(100vh - 70px);
    border: 1px solid #e7e4fb;
    border-radius: 10px;
    overflow: hidden;
    // 表情盒子样式
    -webkit-user-select: none;
    -moz-user-select: none;
    -o-user-select: none;
    user-select: none;
  }
  .resize-top {
    position: relative;
    display: flex;
    height: $topHeight;
    .left {
      position: relative;
      overflow: hidden;

      .my-tree {
        height: calc(100%);
        padding-top: 41px;
        overflow-y: auto;
        &::-webkit-scrollbar {
          width: 4px;
          height: 20px !important;
        }
      }
      .my-search {
        position: absolute;
        top: 0;
        left: 0;
        z-index: 999;
        display: flex;
        align-items: center;
        justify-content: space-between;
        height: 41px;
        background: #f2f2ff;
        width: 100%;
        padding: 0 15px;
      }
      .alternativeType {
        display: flex;
        align-items: center;
        :deep(.el-checkbox) {
          margin-right: 3px;
        }
      }
      border-right: 1px solid var(--el-border-color);
      :deep(.el-tree-node__content) {
        position: relative;
      }
      :deep(.el-button) {
        padding: 8px !important;
      }
      .node {
        position: relative;
        height: 34px;
        line-height: 34px;
        display: flex;
        width: 100%;
        padding: 15px;
        background-color: #f8f8ff;
        align-items: center;
        color: #303133;
        &:nth-of-type(odd) {
          background-color: #fff;
        }
        .isHover {
          cursor: pointer;
          &:hover {
            color: rgba(102, 85, 231, 0.8) !important;
          }
        }
        .nodeActive {
          color: rgba(102, 85, 231, 1) !important;
        }
        .my-radio {
          width: 14px;
          height: 14px;
          border-radius: 50%;
          border: 1px solid var(--el-border-color);
          &:hover {
            cursor: pointer;
            border: 1px solid rgba(102, 85, 231, 0.8);
          }
        }
        .radio_active {
          background-color: rgba(102, 85, 231);
          position: relative;
          &::after {
            content: '';
            position: absolute;
            width: 6px;
            height: 6px;
            top: 3px;
            left: 3px;
            border-radius: 50%;
            background-color: #fff;
          }
        }
        .nodeIcon {
          margin-right: 6px;
          margin-top: 2px;
          cursor: pointer;
        }
        .circle {
          width: 8px;
          height: 8px;
          background: #df5352;
          border-radius: 50%;
          margin: 0 4px 0 4px;
        }
        .node-txt {
          margin-left: 4px;
          max-width: calc(100% - 70px);
          overflow: hidden;
          white-space: nowrap;
          text-overflow: ellipsis;
          word-break: break-all;
        }
        .node-num {
          position: absolute;
          z-index: 0;
          right: 5px;
          top: 8px;
          // vertical-align: text-top;
          padding: 2px 5px;
          height: 18px;
          margin-right: 5px;
          font-size: 12px;
          line-height: 22px;
          color: #303133;
          border-radius: 4px;
          background-color: #ecedef;
          line-height: 15px;
          text-align: center;
        }
      }

      .items {
        height: calc(100% - 55px);
        overflow-y: auto;

        &::-webkit-scrollbar {
          width: 3px;
          height: 0;
        }
        .item {
          margin-top: 5px;
          padding: 10px;
          width: 224px;
          // background-color: #f3faff;
          display: flex;
          .get {
            background-color: #4dac61;
          }
          .info {
            font-size: 12px;
            line-height: 18px;
            overflow: hidden;
            text-align: center;
            text-overflow: ellipsis;
            word-break: break-all;
            white-space: 3;
            text-align: left;
          }
          &:hover {
            cursor: pointer;
            background-color: #f3faff;
            .info {
              color: #2493fe;
            }
          }
        }
        .itemActive {
          background-color: #f3faff;
          .info {
            color: #2493fe;
          }
        }
      }
    }
    .right {
      // width: 100;
      margin-top: 15;
      flex: 1;
      display: flex;
      overflow-y: auto;

      .max-title {
        max-width: 180px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        word-break: break-all;
        word-wrap: break-word;
      }

      &::-webkit-scrollbar {
        width: 3px;
        height: 0;
      }

      .content {
        box-sizing: border-box;
        height: 100%;
        position: relative;
        border-right: 1px solid var(--el-border-color);
        width: 100%;
        overflow: hidden;

        .content-warp {
          display: flex;
          align-items: center;
          flex-wrap: wrap;
          overflow: hidden;
          border-bottom: 1px solid #e7e4fb;
          background: #f9f9ff;
          &:nth-last-of-type() {
            border-bottom: none;
          }
          .content-warp-title {
            border-radius: 4px 0px 0px 4px;
            line-height: 36px;
            height: 100%;
            width: 100px;
            display: flex;
            justify-content: center;
            align-items: center;
          }
          .content-value {
            background-color: #fff;
            min-height: 36px;
            border-left: 1px solid #e7e4fb;
          }
          .content-val-sp {
            padding: 10px 20px;
          }
        }

        .one-of {
          height: 26px;
          line-height: 26px;
          padding: 0px 12px;
          margin-left: 10px;
          font-size: 13px;
          background: #5d7cc5;
          border-radius: 2px;
          color: #fff;
        }
        .content-value {
          flex: 1;
          display: flex;
          flex-wrap: wrap;
          padding: 0 20px;
          &.description {
            line-height: 36px;
          }
        }
        .value {
          max-width: 100%;
          line-height: 24px;
          font-size: 13px;
          margin-right: 4px;
          word-break: break-all;
          white-space: 3;
          text-align: left;
        }
        :deep() {
          .el-tabs {
            height: 100%;
          }
          .el-tabs__content {
            height: calc(100% - 40px);
            overflow: hidden;
          }
          .el-tab-pane {
            height: 100%;
          }
          .el-tabs__nav-wrap {
            &::after {
              display: inline;
            }
          }
          .el-tabs__nav {
            margin-left: 20px;
          }
        }
        .content-view {
          padding-top: 41px;
          height: 100%;
          position: relative;
          .thumbnail {
            position: absolute;
            padding: 0 10px;
            z-index: 11;
            right: 15px;
            bottom: 15px;
            width: 179px;
            height: 36px;
            background: #ffffff;
            display: flex;
            align-items: center;
            box-shadow: 0px 0px 3px 3px rgba(77, 81, 86, 0.05);
            border-radius: 6px;
            border: 1px solid #e9e8fa;
            .tools {
              display: flex;
              align-items: center;
              .control-keys {
                height: 14px;
                width: 107px;
                // padding: 0 10px;
                margin: 0 10px;
                display: flex;
                align-items: center;
                justify-content: space-around;
                border-left: 1px solid #e9e8fa;
                border-right: 1px solid #e9e8fa;
              }
            }
            .canvas {
              position: absolute;
              right: -1px;
              bottom: 38px;
              width: 180px;
              height: 92px;
              overflow: hidden;
              background: #ffffff;
              box-shadow: 0px 0px 3px 3px rgba(77, 81, 86, 0.05);
              border-radius: 6px;
              border: 1px solid #e9e8fa;
              padding: 2px;
              .canvas-map {
                position: absolute;
                top: 50%;
                left: 50%;
                canvas {
                  opacity: 0.2;
                }
              }
              .viewing-area {
                position: absolute;
                z-index: 999;
                top: 0px;
                left: 0px;
                // width: 50%;
                // height: 50%;
                // transform: translate(-50%, -50%);
                border-radius: 6%;
                background: #6954f0;
                opacity: 0.3;
                cursor: move;
              }
            }
          }
        }
        .content-host {
          position: absolute;
          top: 0;
          left: 0;
          z-index: 9;
          height: 41px;
          background: #f2f2ff;
          width: 100%;
          padding: 0 15px;
          text-align: left;
          font-size: 18px;
          display: flex;
          align-items: center;
          justify-content: space-between;
          font-weight: 500;
          color: #303338;

          .content-host-text {
            width: calc(100% - 60px);
            overflow: hidden;
            word-break: break-all;
            white-space: nowrap;
            text-overflow: ellipsis;
          }
          .tag {
            vertical-align: middle;
            padding: 2px 5px;
            height: 22px;
            margin-right: 5px;
            font-size: 13px;
            color: #8e87b6;
            border-radius: 3px;
            line-height: 19px;
            text-align: center;
            background: #ffffff;
            border-radius: 4px;
            border: 1px solid #aea7d3;
            &:hover {
              cursor: pointer;
            }
          }
          .tagActive {
            background: #59547b;
            color: #ffffff;
          }
          .urlName {
            vertical-align: bottom;
            line-height: 26px;
          }
        }
      }
    }
    .column {
      padding: 20px 0 20px 55px;
      height: 100%;
      display: flex;
      // align-items: center;
      justify-content: flex-start;
      flex-direction: column;
    }
    .column-1,
    .column-2 {
      display: flex;
      align-items: center;
    }
    .column-1 {
      border-left: 1px solid #e7e4fb;
      position: relative;
      &:nth-child(1) {
        &::after {
          position: absolute;
          z-index: 1;
          content: '';
          width: 3px;
          top: 0px;
          left: -1px;
          height: 50%;
          background-color: #fff;
        }
      }
      &:nth-last-child(1) {
        &::after {
          position: absolute;
          z-index: 1;
          content: '';
          width: 3px;
          bottom: -1px;
          left: -1px;
          height: 50%;
          background-color: #fff;
        }
      }
    }
    .column-2 {
      border-left: 1px solid #e7e4fb;
      position: relative;
      &:nth-child(1) {
        &::after {
          position: absolute;
          z-index: 1;
          content: '';
          width: 3px;
          top: 0;
          left: -1px;
          height: 50%;
          background-color: #fff;
        }
      }
      &:last-child {
        &::after {
          position: absolute;
          z-index: 1;
          content: '';
          width: 3px;
          bottom: -1px;
          left: -1px;
          height: 50%;
          background-color: #fff;
        }
        &::before {
          position: absolute;
          z-index: 1;
          content: '';
          width: 3px;
          bottom: -1px;
          left: -1px;
          height: 50%;
          background-color: #fff;
        }
      }
    }
    // border: 1px solid #E7E4FB;
    .column-1-title {
      margin: 30px;
      // width: 150px;
      height: 30px;
      padding: 15px 30px 15px 15px;
      line-height: 30px;
      text-align: center;
      display: flex;
      align-items: center;
      background: #ffffff;
      border-radius: 4px;
      border: 1px solid #4b4668;
      font-size: 14px;
      color: #4b4668;
      position: relative;
      white-space: nowrap;

      .documentCopy {
        position: absolute;
        right: 8px;
        top: 8px;
      }
      .mark {
        position: absolute;
        width: 31px;
        height: 2px;
        right: -32px;
        top: 14px;
        background-color: #fff;
        z-index: 3;
      }
      &::after {
        position: absolute;
        z-index: 1;
        content: '';
        width: 31px;
        top: 15px;
        left: -32px;
        height: 1px;
        background-color: #e7e4fb;
      }
      &::before {
        position: absolute;
        z-index: 1;
        content: '';
        width: 31px;
        top: 15px;
        right: -32px;
        height: 1px;
        background-color: #e7e4fb;
      }
    }
    .myMark {
      &::after {
        width: 0px;
      }
    }
    .column-2-title {
      margin: 10px 30px;
      // overflow: hidden;
      height: 30px;
      line-height: 30px;
      text-align: center;
      background: #4b4669;
      border-radius: 4px;
      border: 1px solid #5b547f;
      font-size: 14px;
      color: #ffffff;
      padding: 0 30px 0 15px;
      position: relative;
      display: flex;
      align-items: center;
      white-space: nowrap;
      cursor: pointer;
      .documentCopy {
        position: absolute;
        right: 8px;
        top: 7px;
      }
      .btn {
        background: #ffeeee;
        border-radius: 6px;
        border: 1px solid #ff4f4f;
        position: absolute;
        height: 16px;
        width: 20px;
        display: flex;
        color: #ff4f4f;
        align-items: center;
        justify-content: center;
        right: -26px;
        top: 6px;
        cursor: pointer;
        &::after {
          position: absolute;
          z-index: 1;
          content: '';
          width: 5px;
          top: 7px;
          left: -6px;
          height: 1px;
          background-color: #e7e4fb;
        }
        &::before {
          position: absolute;
          z-index: 1;
          content: '';
          width: 5px;
          top: 7px;
          right: -6px;
          height: 1px;
          background-color: #e7e4fb;
        }
      }
      &::after {
        position: absolute;
        z-index: 1;
        content: '';
        width: 30px;
        top: 14px;
        left: -31px;
        height: 1px;
        background-color: #e7e4fb;
      }
      // &::before {
      //   position: absolute;
      //   z-index: 1;
      //   content: '';
      //   width: 30px;
      //   top: 14px;
      //   right: -31px;
      //   height: 1px;
      //   background-color: #e7e4fb;
      // }
    }
    .column-3 {
      margin-right: 30px;
      border-left: 1px solid #e7e4fb;
      position: relative;
    }
    .column-3-item {
      // border-left: 1px solid #e7e4fb;
      position: relative;
      .btn-num {
        background: #ffeeee;
        border-radius: 6px;
        border: 1px solid #ff4f4f;
        position: absolute;
        height: 16px;
        min-width: 20px;
        display: flex;
        color: #ff4f4f;
        align-items: center;
        justify-content: center;
        right: 4px;
        top: -8px;
        cursor: pointer;
      }
      &:nth-child(1) {
        &::after {
          position: absolute;
          z-index: 1;
          content: '';
          width: 3px;
          top: -10px;
          left: -1px;
          height: calc(50% + 10px);
          background-color: #fff;
        }
      }
      &:last-child {
        &::before {
          position: absolute;
          z-index: 1;
          content: '';
          width: 3px;
          bottom: -11px;
          left: -1px;
          height: calc(50% + 10px);
          background-color: #fff;
        }
      }
    }
    .column-3-item-title {
      margin: 10px 0 10px 30px;
      overflow: hidden;
      width: 620px;
      line-height: 30px;
      background: #ffffff;
      border-radius: 4px;
      border: 1px solid #dfdbf0;
      &::after {
        position: absolute;
        z-index: 1;
        content: '';
        width: 30px;
        top: 50%;
        left: 0px;
        height: 1px;
        background-color: #e7e4fb;
      }
    }
  }
  .footer {
    z-index: 999;
    // height: calc(100% - $topHeight - 50px);
    .time {
      height: 50px;
      width: 100%;
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 0 15px;
      .el-form-item {
        margin-bottom: 0;
        margin-right: 0;
      }
    }
  }
</style>
