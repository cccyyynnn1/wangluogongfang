<script lang="ts">
  export default {
    name: 'AssetManager',
  }
</script>

<script setup lang="ts">
  import { ref } from 'vue'
  import { Edit } from '@element-plus/icons-vue'
  import { getAssetNewRelatApi, getAssetFirewallPolicyApi, updateAssetInfoApi } from '@/api-ecs/assets-preview'
  import { PreviewAssetsItem } from '@/types/index'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import TraceabilityRelationsDatas from '@/ecs/alert/components/traceability-relations-data.vue'
  import TraceabilityServiceVisit from '@/ecs/alert/components/traceability-service-visit.vue'
  import TraceabilityOnlineDecoding from '@/ecs/alert/components/traceability-online-decoding.vue'
  // import ThreatDistribution from './charts/threat-distribution.vue'
  // import AttackSourceStatistics from './charts/attack-source-statistics.vue'
  import { useMouseInElement } from '@vueuse/core'
  import ThreatView from './charts/threat-view.vue'
  import d3GraphForce2 from '@/components/d3-graph-force2.vue'
  import { ElScrollbar } from 'element-plus'
  import _ from 'lodash'
  import dayjs from 'dayjs'

  const $baseMessage: any = inject('$baseMessage')
  const props = defineProps<{
    previewInfo: PreviewAssetsItem
  }>()
  const inside = ref(null)
  const { isOutside } = useMouseInElement(inside)
  let activitieScrollTop: number[] = []
  const router = useRouter()
  const scrollbarRef = ref<InstanceType<typeof ElScrollbar>>()
  const inFoData = ref<PreviewAssetsItem>(_.cloneDeep(props.previewInfo))
  const baseInfo = ref<PreviewAssetsItem>(_.cloneDeep(props.previewInfo))
  const accessRelationship = reactive({
    links: [] as {
      clientIp: string
      serverIp: string
      source: string
      target: string
    }[],
    nodes: [] as {
      color: string
      icon: string
      id: string
      ip: string
      name: string
    }[],
  })
  const activeName = ref(0)
  const edit = ref('none')
  const activitieClass = ['basic', 'lable', 'fingerprint', 'info', 'charts', 'visiting', 'firewalls']

  const activities = [
    {
      content: '基本信息',
    },
    {
      content: '资产标签',
    },
    {
      content: '流量指纹',
    },
    {
      content: '资产关联信息',
    },
    {
      content: '威胁视图',
    },
    {
      content: '访问关系图',
    },
    {
      content: '防火墙策略',
    },
  ]

  const firewallPolicy = reactive({
    page: 1,
    limit: 20,
    total: 0,
    list: [],
  })
  // 深度挖掘显示
  const detailVisible = ref(false)
  // 深度挖掘数据
  const showInfoData = ref<{ id: string; x: number; y: number; ip: string }>()
  // 深度挖掘查询数据
  const queryData = reactive({
    endDate: '',
    startDate: '',
  })
  const visitingData = reactive({
    endDate: '',
    startDate: '',
  })
  // 导航滚动
  const handleNavClick = (index: number) => {
    scrollbarRef.value?.scrollTo({
      top: activitieScrollTop[index],
      behavior: 'smooth',
    })
  }
  // 资产访问关系
  const handleGetAssetNewRelat = _.debounce(async (ip: string | undefined) => {
    if (!ip) return
    const { data } = await getAssetNewRelatApi({ ip, startTime: visitingData.startDate, endTime: visitingData.endDate })
    accessRelationship.links = (data?.links || []).map((link) => ({
      ...link,
      source: link.clientIp,
      target: link.serverIp,
    }))
    accessRelationship.nodes = (data?.nodes || []).map((node) => ({ ...node, id: node.ip }))
  }, 50)
  // 防火墙策略
  const handleGetAssetFirewallPolicy = async (ip: string | undefined) => {
    if (!ip) return
    const { data } = await getAssetFirewallPolicyApi({
      page: firewallPolicy.page,
      limit: firewallPolicy.limit,
      query: JSON.stringify({ serverIp: ip }),
    })
    firewallPolicy.list = data.list || []
    firewallPolicy.total = data.total || 0
  }

  function getScrollTop() {
    activitieScrollTop = activitieClass.map((i) => {
      const ele = document.querySelector(`.${i}`) as HTMLElement
      return ele?.offsetTop ?? 0
    })
  }

  const scroll = _.debounce(({ scrollTop }: any) => {
    const index = activitieScrollTop.findIndex((i) => scrollTop <= i)
    activeName.value = index
  })

  const handleEditConfirm = async () => {
    const {
      appName,
      appType,
      businessName,
      businessType,
      datasource,
      dbVersion,
      id,
      middlewareVersion,
      netPartation,
      phone,
      protocol,
      sysVersion,
    } = inFoData.value

    const _data = {
      appName,
      appType,
      businessName,
      businessType,
      datasource,
      dbVersion,
      id: id || 0,
      middlewareVersion,
      netPartation,
      phone,
      protocol,
      sysVersion,
    }
    try {
      const { data } = await updateAssetInfoApi(_data)
      baseInfo.value = {
        ...baseInfo.value,
        ..._data,
      }
      $baseMessage('更新成功', 'success', 'vab-hey-message-success')
      handleEditCancel()
    } catch (error) {
      console.error(error)
    }
  }
  const handleEditCancel = () => {
    edit.value = 'none'
  }
  // 详情
  const handleToSite = (hosts: string) => {
    router.push({ name: 'SiteIndex', query: { hosts } })
  }
  /**
   *
   * @param groupId 资产组id  1:流量指纹 2:规则指纹 3:三方指纹
   */
  const getFingerprintLabel = (groupId: 1 | 2 | 3) => {
    const group = props.previewInfo?.labels?.find((item) => item.groupId === groupId)
    return group ? group.labelList : []
  }

  const clickHandler = (params: any) => {
    // 此时params有两种结构，线和点的结构，但同时和检索的访问关系数据有差且，资产标签的数据是不一致，需要做处理。(有source或者target表明是连线，数据用type区分，只在资产标签这里需要，其他照旧)
    detailVisible.value = true
    showInfoData.value = params.source || params.target ? { ...params, type: 'ASSETS_LABELS' } : params
  }
  watch(
    () => detailVisible.value,
    (val) => {
      if (val) {
        const curtimeDate = dayjs()
        queryData.startDate = curtimeDate.subtract(1, 'hour').format('YYYY-MM-DD HH:mm:ss')
        queryData.endDate = curtimeDate.add(10, 'minute').format('YYYY-MM-DD HH:mm:ss')
        // visitingData.startDate = curtimeDate.subtract(6, 'hour').format('YYYY-MM-DD HH:mm:ss')
        // visitingData.endDate = curtimeDate.add(10, 'minute').format('YYYY-MM-DD HH:mm:ss')
      }
    }
  )

  watch(
    () => visitingData,
    () => {
      setTimeout(() => {
        handleGetAssetNewRelat(IP.value)
      }, 0)
    },
    {
      deep: true,
    }
  )

  watch(
    () => isOutside.value,
    () => {
      console.log(isOutside.value)
    }
  )

  const IP = ref()

  onMounted(() => {
    IP.value = baseInfo.value.ipv4 || baseInfo.value.ipv6
    // handleGetAssetNewRelat(IP.value)
    handleGetAssetFirewallPolicy(IP.value)
    nextTick(() => getScrollTop())
  })

  // 查询时间
  const timeDuration = ref('6hours')
  const formateTimeStr = 'YYYY-MM-DD HH:mm:ss'
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

  watchEffect(() => {
    const curtimeDate = dayjs()
    const endDate = [
      '1hours',
      '3hours',
      '6hours',
      '24hours',
      'last-week',
      'last-two-week',
      'last-three-days',
      'last-one-month',
    ].includes(timeDuration.value)
      ? curtimeDate.add(10, 'minute').format(formateTimeStr)
      : curtimeDate.endOf('day').format(formateTimeStr)

    const cur = curtimeDate.startOf('day')

    let startDate = null
    switch (timeDuration.value) {
      case '1hours':
        startDate = curtimeDate.subtract(1, 'hour').format(formateTimeStr)
        break
      case '3hours':
        startDate = curtimeDate.subtract(3, 'hour').format(formateTimeStr)
        break
      case '6hours':
        startDate = curtimeDate.subtract(6, 'hour').format(formateTimeStr)
        break
      case '24hours':
        startDate = curtimeDate.subtract(24, 'hour').format(formateTimeStr)
        break
      case 'today':
        startDate = curtimeDate.startOf('day').format(formateTimeStr)
        break
      case 'yesterday':
        startDate = cur.subtract(1, 'day').format(formateTimeStr)
        break
      case 'last-three-days':
        startDate = curtimeDate.subtract(3, 'day').format(formateTimeStr)
        break
      case 'last-week':
        startDate = curtimeDate.subtract(1, 'week').format(formateTimeStr)
        break
      case 'last-two-week':
        startDate = curtimeDate.subtract(2, 'week').format(formateTimeStr)
        break
      case 'last-one-month':
        startDate = curtimeDate.subtract(1, 'month').format(formateTimeStr)
        break
      default:
        startDate = curtimeDate.startOf('day').format(formateTimeStr)
        break
    }
    visitingData.endDate =
      timeDuration.value === 'yesterday' ? curtimeDate.endOf('day').subtract(1, 'day').format(formateTimeStr) : endDate
    visitingData.startDate = startDate
  })

  watchEffect(() => {
    if (timeDuration.value === 'user-defined') {
      timeDate.value = [dayjs(visitingData.startDate), dayjs(visitingData.endDate)]
    }
  })
  watchEffect(() => {
    if (timeDate.value) {
      const [start, end] = timeDate.value
      visitingData.endDate = dayjs(end).format(formateTimeStr)
      visitingData.startDate = dayjs(start).format(formateTimeStr)
    }
  })
</script>

<template>
  <div class="asset-manager-content">
    <el-timeline class="navBarBox">
      <el-timeline-item
        v-for="(activity, index) in activities"
        :key="index"
        center
        :class="{ isActive: activeName === index }"
        @click="() => handleNavClick(index)"
      >
        <div class="navBar">{{ activity.content }}</div>
      </el-timeline-item>
    </el-timeline>
    <div class="asset-info-content">
      <el-scrollbar ref="scrollbarRef" always @scroll="scroll">
        <el-descriptions border class="basic" :column="3" title="基本信息">
          <template #extra>
            <div class="extra" @click="() => (edit = 'basic')">
              <Edit />
              编辑
            </div>
          </template>
          <el-descriptions-item label="业务名称">
            <el-select
              v-if="edit === 'basic'"
              v-model="inFoData.businessName"
              allow-create
              filterable
              style="width: 100%"
            >
              <el-option
                v-for="option in inFoData.historyList"
                :key="option.businessName"
                :label="option.businessName"
                :value="option.businessName"
              />
            </el-select>
            <div v-else v-copy="baseInfo.businessName">{{ baseInfo.businessName }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="应用名称">
            <el-input v-if="edit === 'basic'" v-model="inFoData.appName" />
            <div v-else v-copy="baseInfo.appName">{{ baseInfo.appName }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="应用类型">
            <el-input v-if="edit === 'basic'" v-model="inFoData.appType" />
            <div v-else v-copy="baseInfo.appType">{{ baseInfo.appType }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="数据中心">
            <el-input v-if="edit === 'basic'" v-model="inFoData.datasource" />
            <div v-else v-copy="baseInfo.datasource">{{ baseInfo.datasource }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="业务IPV4">
            <div v-copy="baseInfo.ipv4">{{ baseInfo.ipv4 }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="业务IPV6">
            <div v-copy="baseInfo.ipv6">{{ baseInfo.ipv6 }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="网络分区">
            <el-input v-if="edit === 'basic'" v-model="inFoData.netPartation" />
            <div v-else v-copy="baseInfo.netPartation">{{ baseInfo.netPartation }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="系统版本">
            <el-input v-if="edit === 'basic'" v-model="inFoData.sysVersion" />
            <div v-else v-copy="baseInfo.sysVersion">{{ baseInfo.sysVersion }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="中间库版本">
            <el-input v-if="edit === 'basic'" v-model="inFoData.middlewareVersion" />
            <div v-else v-copy="baseInfo.middlewareVersion">{{ baseInfo.middlewareVersion }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="数据库版本">
            <el-input v-if="edit === 'basic'" v-model="inFoData.dbVersion" />
            <div v-else v-copy="baseInfo.dbVersion">{{ baseInfo.dbVersion }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="联系方式">
            <el-input v-if="edit === 'basic'" v-model="inFoData.phone" />
            <div v-else v-copy="baseInfo.phone">{{ baseInfo.phone }}</div>
          </el-descriptions-item>
        </el-descriptions>
        <div v-if="edit === 'basic'" class="edit-bar">
          <el-button type="primary" @click="handleEditConfirm">保存</el-button>
          <el-button @click="handleEditCancel">取消</el-button>
        </div>
        <el-descriptions border class="lable single" :column="3" title="资产标签">
          <el-descriptions-item>
            <div style="height: 30px"></div>
          </el-descriptions-item>
        </el-descriptions>
        <el-descriptions border class="fingerprint single" :column="3" title="流量指纹">
          <el-descriptions-item>
            <div v-for="label in getFingerprintLabel(1)" :key="label.id" class="fingerprint-label">
              {{ label.labelName }}
            </div>
          </el-descriptions-item>
        </el-descriptions>
        <el-descriptions border class="info" :column="1" title="资产关联信息">
          <el-descriptions-item label="服务端口">
            <div>{{ baseInfo?.portList.join(' , ') }}</div>
          </el-descriptions-item>
          <div v-for="(item, index) in baseInfo?.hostVoList" :key="index">
            <el-descriptions-item class="ceshi" />
            <el-descriptions-item label="HOST信息">
              <div>
                {{ item?.host }}
                <span v-if="!item.know">（未加入已知站点）</span>
              </div>
            </el-descriptions-item>
            <el-descriptions-item label="资产详情">
              <div style="display: flex; flex-wrap: wrap">
                <div>
                  {{ item?.businessInfo.join(' , ') }}
                </div>
              </div>
            </el-descriptions-item>
            <el-descriptions-item class-name="noPading flex" label="API接口">
              <div style="display: flex; flex-wrap: wrap">
                <div v-for="(td, index) in item?.apis" :key="`url-${index}`" class="api-item">
                  <div :style="{ cursor: item.know ? 'pointer' : 'default', width: '100%' }">
                    <!-- <el-button v-if="item.know" link type="primary" @click="handleToSite(baseInfo.host)">
                      {{ td }}
                    </el-button> -->
                    <a v-if="item.know" target="_blank">
                      <span style="text-decoration: underline" @click="handleToSite(baseInfo.host)">{{ td }}</span>
                    </a>
                    <span v-else>{{ td }}</span>
                  </div>
                </div>
              </div>
            </el-descriptions-item>
          </div>
        </el-descriptions>
        <el-descriptions border class="charts single" :column="1" title="威胁视图">
          <el-descriptions-item>
            <div class="my-x-charts">
              <ThreatView :server-ip="IP" />
            </div>
            <!-- <el-row :gutter="20" style="height: 240px">
              <el-col class="right-part" :span="11">
                <ThreatDistribution :distributions="threatType" />
              </el-col>
              <el-col :span="13">
                <AttackSourceStatistics :attack-sources="attackIp" />
              </el-col>
            </el-row> -->
          </el-descriptions-item>
        </el-descriptions>
        <!-- <el-descriptions border class="fingerprint" :column="1" title="流量指纹">
          <el-descriptions-item label="系统指纹">
            <div v-for="label in getFingerprintLabel(2)" :key="label.id" class="fingerprint-label">
              {{ label.labelName }}
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="服务指纹">
            <div v-for="label in getFingerprintLabel(3)" :key="label.id" class="fingerprint-label">
              {{ label.labelName }}
            </div>
          </el-descriptions-item>
        </el-descriptions> -->
        <el-descriptions border class="visiting single noPadding" :column="1" title="访问关系">
          <el-descriptions-item>
            <div style="height: 454px; background: #f3f9ff; position: relative">
              <div class="search-left">
                <el-select
                  ref="inside"
                  v-model="timeDuration"
                  popper-class="visiting-select"
                  style="margin-right: 6px; height: 30px; width: 120px"
                >
                  <el-option
                    v-for="item in timeDuratioOptions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
                <vab-date-time-picker v-if="timeDuration == 'user-defined'" v-model="timeDate" style="width: 336px" />
              </div>
              <d3-graph-force2 :disable-keys="['wheel']" :graph-data="accessRelationship" @click="clickHandler" />
            </div>
          </el-descriptions-item>
        </el-descriptions>
        <el-descriptions border class="firewalls single noBorder" :column="1" title="防火墙策略" />
        <div style="width: 100%; height: 545px">
          <el-table :data="firewallPolicy.list" style="width: 100%; margin-right: -10px; height: 100%">
            <el-table-column label="防火墙" prop="firewallIdStr" width="150" />
            <el-table-column label="策略ID" prop="policyName" width="280" />
            <el-table-column label="优先级" prop="num" />
            <el-table-column label="源区域" prop="fromZoneIdsStr" width="250" />
            <el-table-column label="目的区域" prop="toZoneIdsStr" width="250" />
            <el-table-column label="源名称" prop="clientName" />
            <el-table-column label="目的名称" prop="serverName" width="100" />
            <el-table-column label="源IP" prop="clientIp" width="180" />
            <el-table-column label="目的IP" prop="serverIp" width="180" />
            <el-table-column label="有效期" prop="expirationDate" width="180" />
            <el-table-column label="状态" prop="enabledStr" />
            <el-table-column label="是否过期" prop="overdueStr" width="100" />
            <template #empty><el-empty /></template>
          </el-table>
          <el-pagination
            v-model:current-page="firewallPolicy.page"
            v-model:page-size="firewallPolicy.limit"
            background
            hide-on-single-page
            layout="total, sizes, prev, pager, next, jumper"
            :page-sizes="[10, 20, 30, 50, 100]"
            style="margin: 20px"
            :total="firewallPolicy.total"
          />
        </div>
      </el-scrollbar>
    </div>
    <vab-dialog
      v-model="detailVisible"
      align-center
      :close-on-click-modal="false"
      destroy-on-close
      title="深度挖掘"
      width="1375px"
    >
      <el-tabs :model-value="'relations'">
        <el-tab-pane label="关系" lazy name="relations">
          <traceability-relations-datas
            :node-data="showInfoData"
            :search-sql="''"
            :time-range="[queryData.startDate, queryData.endDate]"
          />
        </el-tab-pane>
        <el-tab-pane label="服务访问" lazy name="service-visit">
          <traceability-service-visit
            :node-data="showInfoData"
            :time-range="[queryData.startDate, queryData.endDate]"
          />
        </el-tab-pane>
        <el-tab-pane label="数据包分析" lazy name="online-decoding">
          <traceability-online-decoding
            :node-data="showInfoData"
            :time-range="[queryData.startDate, queryData.endDate]"
          />
        </el-tab-pane>
      </el-tabs>
    </vab-dialog>
  </div>
</template>
<style>
  .visiting-select {
    z-index: 2000;
  }
</style>
<style scoped lang="scss">
  .asset-manager-content {
    :deep() {
      .info {
        tbody {
          tr {
            &:nth-child(4n - 2) {
              .el-descriptions__cell {
                background-color: #fff;
                border: none !important;
                height: 20px;
              }
            }
            &:nth-child(4n - 1) {
              .el-descriptions__cell {
                background: #f5f7fa;
                color: #333333;
                font-weight: 500;
              }
            }
          }
        }
      }
    }
    display: flex;
    height: calc(100vh - 120px);
    padding-bottom: 20px;

    .my-x-charts {
      height: 80vh;
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
    }
    .asset-info-content {
      flex: 1;
      padding-right: 10px;
      overflow-y: auto;
      .edit-bar {
        margin: -20px 0 40px 0;
        text-align: right;
      }
      .extra {
        display: flex;
        margin-right: 20px;
        font-size: 14px;
        font-weight: 400;
        color: #a9acb3;
        cursor: pointer;
        svg {
          width: 15px;
          margin-right: 2px;
        }
      }
      .host-info {
        margin-top: 10px;
      }
      .tabView {
        display: flex;
        align-items: center;
        min-width: 50%;

        div {
          flex: 1;
          line-height: 38px;
          text-align: left;
          padding-left: 11px;
          width: 100%;
          height: 38px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          word-break: break-all;
          word-wrap: break-word;
          &:not(:nth-child(2n)) {
            border-right: 1px solid #e1e3ec;
          }
          // &:not(:last-child) {
          //   border-right: 1px solid #e1e3ec;
          // }
        }
      }
      .api-item {
        flex: 1;
        text-align: left;
        padding-left: 11px;
        min-width: 50%;
        &:not(:nth-child(2n)) {
          border-right: 1px solid #e1e3ec;
        }
        &:nth-child(n) {
          border-top: 1px solid #e1e3ec;
        }
        &:last-child {
          border-right: 0px solid #e1e3ec;
        }
        div {
          line-height: 38px;
          height: 38px;
          width: 100%;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          word-break: break-all;
          word-wrap: break-word;
        }
        &:nth-child(1) {
          border-top: 0px solid #e1e3ec;
        }
        &:nth-child(2) {
          border-top: 0px solid #e1e3ec;
        }
      }
      .fingerprint-label {
        color: var(--el-color-primary);
        background-color: var(--el-color-primary-light-9);
        display: inline-block;
        padding: 0 20px;
        height: 36px;
        line-height: 36px;
        border-radius: 4px;
        margin: 5px 10px 5px 0;
      }
      :deep() {
        .el-descriptions {
          font-family: PingFangSC, PingFang SC;
          font-size: 14px;
          font-weight: 400;
          .search-left {
            position: absolute;
            z-index: 2012;
            display: flex;
            top: 15px;
            left: 15px;
          }
          &:not(:last-child) {
            margin-bottom: 40px;
          }
          &.single {
            .el-descriptions__label {
              display: none;
            }
          }
          &.noBorder {
            margin-bottom: 0px;
            .el-descriptions__table {
              border: none;
            }
          }
          &.noPadding {
            .el-descriptions__content {
              padding: 0 !important;
            }
          }
          .el-descriptions__header {
            height: 42px;
            padding-left: 16px;
            margin-bottom: 0;
            background: var(--el-descriptions-item-bordered-label-background);
            border-inline: var(--el-descriptions-table-border);
            border-top: var(--el-descriptions-table-border);
          }
          .el-descriptions__title {
            font-family: PingFangSC, PingFang SC;
            font-size: 16px;
            font-size: 16px;
            font-weight: 500;
            color: #303133;
            &::before {
              top: 6px;
              background-color: var(--el-color-primary);
            }
          }
          .el-descriptions__label {
            width: 166px;
            color: #303133;
            text-indent: 17px;
          }
          .el-descriptions__content {
            color: #606266;
            &.noPading {
              min-height: 38px;
              padding: 0;
            }
          }
        }
      }
    }
    .navBarBox {
      width: 180px;
      padding-left: 30px;
    }
    :deep() {
      .el-timeline-item__wrapper {
        cursor: pointer;
      }
      .el-timeline-item__tail {
        border-left-style: dashed;
      }
      .el-timeline-item__node {
        top: 7px;
      }
      .el-timeline-item__timestamp {
        display: none;
      }
      .el-timeline-item__content {
        width: fit-content;
        font-family: PingFangSC, PingFang SC;
        font-size: 14px;
        font-weight: 500;
        color: #606266;
      }
      .navBar {
        position: relative;
        height: 30px;
        padding-inline: 14px;
        line-height: 30px;
        border-radius: 4px;
        transition: all 0.3s;
        &::before {
          position: absolute;
          top: 11px;
          left: -8px;
          z-index: 1;
          display: block;
          content: '';
          border: 4px solid transparent;
          border-right-color: var(--el-color-primary);
          opacity: 0;
          transition: all 0.3s;
        }
      }
      .isActive {
        .el-timeline-item__node {
          background: var(--el-color-white);
          border-color: var(--el-color-primary);
          border-style: solid;
          border-width: 3px;
        }
        .navBar {
          color: #ffffff;
          background-color: var(--el-color-primary);
          &::before {
            opacity: 1;
          }
        }
      }
    }
  }
  :deep() {
    .el-dialog__footer {
      display: none !important;
    }
    .basic .el-descriptions__content {
      width: 400px;
    }
  }
  :deep() {
    .el-table__body-wrapper {
      max-height: 409px;
      min-height: 409px;
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0 !important;
        height: 0;
      }
    }
  }
</style>
