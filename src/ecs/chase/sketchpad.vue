<script lang="ts">
  export default {
    name: 'Sketchpad',
  }
</script>

<script setup lang="ts">
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import SearchSql from '@/components/search-sql/index.vue'
  import UpdataItem from './components/updata-item.vue'
  import DrawerInfo from './components/drawer-info.vue'
  import DrawerFilter from './components/drawer-filter.vue'
  import DrawerSelf from './components/drawer-self.vue'
  import DrawerAlert from './components/drawer-alert.vue'
  import { Search } from '@element-plus/icons-vue'
  import SvgExpand from './components/svg-expand.vue'
  import ChaseGraph from './components/chase-d3-graph.vue'
  import TraceabilityRelationsDatas from '@/ecs/alert/components/traceability-relations-data.vue'
  import TraceabilityServiceVisit from '@/ecs/alert/components/traceability-service-visit.vue'
  import TraceabilityOnlineDecoding from '@/ecs/alert/components/traceability-online-decoding.vue'
  import dayjs from 'dayjs'
  import { chasePathTracingResponse, chasePathTracingResquest, link, node } from '~/src/types'
  import { HuntingPathTracingApi, QueryAssetByIpApi } from '~/src/api-ecs/chase'
  import SketchpadModel from './sketchpad'

  const props = defineProps<{
    type: 'pathTracing' | 'threatHunting'
  }>()
  const activeName = ref('relations')
  //表单数据
  const queryForm = reactive<chasePathTracingResquest>({
    searchSql: '',
    aggregationFields: ['clientIp', 'serverIp', 'serverPort'],
    count: 200,
    indexType: 1,
    startTime: '',
    endTime: '',
  })

  const sqlComponentsRef = ref<InstanceType<typeof SearchSql>>()
  const pageLoading = ref(false)
  const isFullscreen = ref(false)
  const drawerAlert = ref(false)
  const drawerInfo = ref(false)
  const infoData = ref<Record<string, unknown>>({})
  const drawerSelf = ref(false)
  const drawerFilter = ref(false)
  const detailVisible = ref(false)
  const graphRef = ref<InstanceType<typeof ChaseGraph>>()
  const allData = ref<chasePathTracingResponse>({
    linksAll: [],
    linksDefault: [],
    nodesAll: [],
    nodesDefault: [],
  })
  const filterData = reactive({
    ip: true,
    assetName: false,
  })

  const sketchpadData = ref<{ links: link[]; nodes: node[] }>({
    links: [],
    nodes: [],
  })
  const changeData = ref<{ links: link[]; nodes: node[] }>({
    links: [],
    nodes: [],
  })

  // onMounted(() => {
  //   searchhandle()
  // })

  const showInfoData = reactive({
    data: {} as { id: string; x: number; y: number; ip: string },
  })

  const grahpClass = ref<SketchpadModel>()
  function echartClicckHandler(params: any) {
    detailVisible.value = true
    // searchBoxHidden.value = true
    showInfoData.data = params
  }
  const searchhandle = async () => {
    if (sqlComponentsRef.value && !sqlComponentsRef.value?.getValidateSql()) return
    pageLoading.value = true
    getQueryDate()
    try {
      if (props.type == 'pathTracing') {
        const { data } = await HuntingPathTracingApi({
          ...queryForm,
          // startTime: '2024-11-18 15:04:20',
          // endTime: '2024-11-18 16:14:20',
        })
        allData.value = data
        grahpClass.value = new SketchpadModel(allData.value.nodesAll, allData.value.linksAll)

        const { linksDefault, nodesDefault } = grahpClass.value.InitData()
        sketchpadData.value.links = linksDefault
        sketchpadData.value.nodes = nodesDefault
      }
      setTimeout(() => {
        graphRef.value?.initData(sketchpadData.value)
      }, 0)
    } finally {
      pageLoading.value = false
    }
  }

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

  watch(
    () => timeDuration.value,
    (newData, oldData) => {
      searchhandle()
    },
    {
      immediate: true,
    }
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

  const handleFilter = () => {
    drawerFilter.value = !drawerFilter.value
  }

  const handleZoom = (val: 'up' | 'down') => {
    graphRef.value?.HandleZoom(val)
  }

  const handleReset = () => {
    const { linksDefault, nodesDefault } = grahpClass.value!.InitData()
    sketchpadData.value.links = linksDefault
    sketchpadData.value.nodes = nodesDefault
    setTimeout(() => {
      graphRef.value?.initData(sketchpadData.value)
      graphRef.value?.reset()
    }, 0)
  }

  const handleExport = () => {
    sketchpadData.value.links = allData.value.linksAll
    sketchpadData.value.nodes = allData.value.nodesAll
    graphRef.value?.initData(sketchpadData.value)
    graphRef.value?.reset()
  }

  const handleShowInfo = async (nodeData: node) => {
    drawerInfo.value = true
    const { data } = await QueryAssetByIpApi({
      ip: nodeData.ip,
      startTime: '2024-11-18 15:04:20',
      endTime: '2024-11-18 16:14:20',
      status: nodeData.isPath!,
    })
    data['indexIp'] = nodeData.ip
    infoData.value = data
  }

  const handleShowFilter = (data: node) => {
    drawerFilter.value = false
    drawerSelf.value = true
  }

  const handleexpandAll = (ip: string) => {
    changeData.value = {
      links: [],
      nodes: [],
    }
    grahpClass.value?.RecursionGetNextLevelDataByIP(ip, 'all', graphRef.value!.addNode, graphRef.value!.addLink)
  }
  const handleexpandOne = (ip: string) => {
    changeData.value = {
      links: [],
      nodes: [],
    }
    grahpClass.value?.RecursionGetNextLevelDataByIP(ip, 'one', graphRef.value!.addNode, graphRef.value!.addLink)
  }

  const handleDelNode = (ip: string) => {
    const { nodesDel, linksDel } = grahpClass.value!.DelDataByIP(ip)
    graphRef.value!.removeNodesAndLinks(nodesDel, linksDel)
  }

  const handleRemarkColor = (data: node) => {
    grahpClass.value!.ChangeRemarkByNodeInfo(data)
  }

  const handleChange = (mode: keyof typeof filterData) => {
    filterData.assetName = !filterData.assetName
    filterData.ip = !filterData.ip
    const show = filterData.ip == true ? 'ip' : 'assetName'
    graphRef.value!.showText(show)
  }
</script>

<template>
  <div class="search-tool">
    <el-input
      v-model.number="queryForm.count"
      class="search-tool-input"
      :max="1000"
      :min="1"
      placeholder="输入查询条数"
    />
    <SearchSql
      ref="sqlComponentsRef"
      class="search-tool-inputSql"
      :index-type="queryForm.indexType"
      :model-value="queryForm.searchSql"
      placeholder="请输入需要追踪的特征条件"
      @on-change="(str) => (queryForm.searchSql = str)"
      @onSearch="searchhandle"
    />
    <el-button :icon="Search" size="large" type="primary" @click="searchhandle">检索</el-button>
  </div>
  <div v-loading="pageLoading" class="sketchpad" :class="{ 'vab-fullscreen': isFullscreen }">
    <ChaseGraph
      ref="graphRef"
      :graph-data="sketchpadData"
      :is-path="props.type"
      :start-self="drawerSelf"
      @close-filter="drawerSelf = false"
      @del-node="handleDelNode"
      @expand-all="handleexpandAll"
      @expand-one="handleexpandOne"
      @remark-node-color="handleRemarkColor"
      @show-filter="handleShowFilter"
      @show-info="handleShowInfo"
      @show-relate="echartClicckHandler"
    />
    <div class="left_tools" :style="{ left: drawerInfo ? '365px' : '15px' }">
      <el-select v-model="timeDuration" style="margin-right: 6px; height: 30px; width: 120px">
        <el-option v-for="item in timeDuratioOptions" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
      <vab-date-time-picker
        v-if="timeDuration == 'user-defined'"
        v-model="timeDate"
        class="sketchpad_time_picker"
        style="width: 336px"
      />
      <div class="tools_btn has_margin" @click="handleZoom('down')">
        <img class="tools_btn_img" :src="require('@/assets/chase/down.svg')" />
      </div>
      <div class="tools_btn has_margin" @click="handleZoom('up')">
        <img class="tools_btn_img" :src="require('@/assets/chase/up.svg')" />
      </div>
      <div class="tools_btn has_margin" style="cursor: not-allowed">
        <img class="tools_btn_img" :src="require('@/assets/chase/prev.svg')" />
      </div>
      <div class="tools_btn has_margin" style="cursor: not-allowed">
        <img class="tools_btn_img" :src="require('@/assets/chase/next.svg')" />
      </div>
      <div class="tools_btn has_margin" @click="handleReset">
        <img class="tools_btn_img" :src="require('@/assets/chase/refresh.svg')" />
      </div>
      <div class="tools_btn" @click="handleExport">
        <SvgExpand class="tools_btn_img" color="#9793b0" />
      </div>
    </div>
    <div class="right_tools" :style="{ right: drawerFilter || drawerAlert || drawerSelf ? '365px' : '15px' }">
      <div v-if="!isFullscreen" class="tools_btn has_margin" @click="isFullscreen = true">
        <img class="tools_btn_img" :src="require('@/assets/chase/full.svg')" />
      </div>
      <div v-else class="tools_btn has_margin" @click="isFullscreen = false">
        <img class="tools_btn_img" :src="require('@/assets/chase/exit.svg')" />
      </div>
      <div class="tools_btn" @click="handleFilter">
        <img class="tools_btn_img" :src="require('@/assets/chase/filtter.svg')" />
      </div>
    </div>
    <!-- 详情 -->
    <DrawerInfo v-if="drawerInfo" :alldata="infoData" @on-close="drawerInfo = false" />
    <!-- 筛选 -->
    <DrawerFilter
      v-if="drawerFilter"
      :alldata="filterData"
      @on-change-event="handleChange"
      @on-close="drawerFilter = false"
    />
    <!-- 告警列表 -->
    <DrawerAlert v-if="drawerAlert" @on-close="drawerAlert = false" />
    <!-- 自定义连线 -->
    <DrawerSelf v-if="drawerSelf" @on-close="drawerSelf = false" />

    <UpdataItem />
    <!-- 深度挖掘 -->
    <vab-dialog
      v-model="detailVisible"
      align-center
      :close-on-click-modal="false"
      destroy-on-close
      title="深度挖掘"
      width="1375px"
      @close="activeName = 'relations'"
    >
      <!-- @close="activeName = chartType === 'attack-perspective' ? 'alert' : 'relations'" -->
      <el-tabs v-model="activeName" class="demo-tabs">
        <el-tab-pane label="关系" lazy name="relations">
          <traceability-relations-datas
            :node-data="showInfoData.data"
            :search-sql="queryForm.searchSql"
            :time-range="[queryForm.startTime, queryForm.endTime]"
          />
        </el-tab-pane>
        <el-tab-pane label="服务访问" lazy name="service-visit">
          <traceability-service-visit
            :node-data="showInfoData.data"
            :time-range="[queryForm.startTime, queryForm.endTime]"
          />
        </el-tab-pane>
        <el-tab-pane label="数据包分析" lazy name="online-decoding">
          <traceability-online-decoding
            :node-data="showInfoData.data"
            :time-range="[queryForm.startTime, queryForm.endTime]"
          />
        </el-tab-pane>
      </el-tabs>
    </vab-dialog>
  </div>
</template>

<style scoped lang="scss">
  .search-tool {
    position: absolute;
    right: 0;
    display: flex;
    align-items: center;
    top: -46px;

    :deep() {
      --el-input-border: #e5e4ef;
      .el-input__inner {
        height: 32px;
      }
      .el-select {
        .el-input__wrapper {
          height: 30px;
        }
      }

      .el-button {
        height: 34px;
        border-radius: 0 6px 6px 0;
        padding: 0 12px;
      }
      .cm-editor {
        height: 34px;
        border: 1px solid #e5e4ef;
        border-radius: 6px 0 0 6px;
      }
      .cm-content {
        padding: 2px 0;
      }
    }
    .search-tool-input {
      width: 120px;
      border-radius: 6px;
    }

    .search-tool-inputSql {
      width: 280px;
      margin-left: 8px;
      border-radius: 6px;
    }
    .focus {
      width: 600px;
      z-index: 9;
    }
  }

  .sketchpad {
    position: relative;
    background: #fcfcff;

    width: 100%;
    height: 100%;

    .sketchpad-drawer {
      position: absolute;
      height: 100%;
      width: 350px;
      right: 0;
      top: 0;
      padding: 10px 15px;
      background-color: #fff;
      // border-right: 1px solid #e9e6f9;
      border-left: 1px solid #e9e6f9;
      transition: all 0.2s;
    }
    .sketchpad_time_picker {
      position: absolute;
      margin-right: 6px;
      height: 30px;
      width: 310px;
      left: 0;
      top: 35px;
    }

    .left_tools {
      z-index: 9;
      position: absolute;
      top: 15px;
      left: 15px;
      display: flex;
      align-items: center;
    }
    .sketchpad_info {
      position: absolute;
    }

    .tools_btn {
      width: 30px;
      height: 30px;
      background: #ffffff;
      border-radius: 6px;
      border: 1px solid #e6e4f4;
      padding: 7px;

      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      .tools_btn_img {
        width: 100%;
        height: 100%;
      }
    }
    .has_margin {
      margin-right: 6px;
    }

    .right_tools {
      position: absolute;
      z-index: 2;
      top: 15px;
      right: 15px;
      display: flex;
    }
    :deep() {
      .change-drawer-position {
        position: absolute !important;
        z-index: 1 !important;
      }
      .my-title {
        margin: 15px 0 10px 0;

        &::before {
          display: inline-block;
          width: 3px;
          height: 10px;
          margin: 0 6px 0 0px;
          content: '';
          background: #6954f0;
        }
      }
      .tip-item,
      .warp-item {
        display: flex;
        width: 100%;
        margin-top: 8px;
        .left {
          width: 80px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          word-break: break-all;
          word-wrap: break-word;
          font-size: 14px;
          color: #918da5;
        }
        .right {
          width: calc(100% - 80px);
          white-space: normal;
          word-break: break-all;
          word-wrap: break-word;
          font-size: 14px;
          color: #2b2742;
        }
      }
      .warp-item {
        align-items: center;
      }
      .sketchpad-drawer-top {
        display: flex;
        align-items: center;
        justify-content: space-between;
        width: 100%;
        height: 30px;
        .title {
          font-size: 24px;
          color: #1e1842;
        }
        i {
          cursor: pointer;
        }
      }
      .el-input__wrapper {
        height: 30px;
      }
    }
  }
</style>
