<script setup lang="ts">
  import AlertAggregationDetailInfo from './alert-detail-aggregation-info.vue'
  import { AlertItem } from '~/src/types'
  import { formatNstime } from '@/utils/time'
  import { useScroll } from '@vueuse/core'
  import { ElTable } from 'element-plus'
  import Traceability from './traceability.vue'
  import { useUserStore } from '@/store/modules/user'
  import { getAllDisPlaysFiledApi } from '@/api-ecs/public'
  import { getSystemConfigApi } from '~/src/api-ecs/system'
  import { TableColumnItemType } from '/#/store'
  import { useCopy } from '@/utils'

  import { searchBySqlApi } from '@/api-ecs/retrieve'
  const alarmIcon = require('@/assets/alarm-icon.svg')
  const props = defineProps<{
    alertDetailVisible: boolean
    selectAlert: any
  }>()
  const emits = defineEmits<{
    (e: 'update:alertDetailVisible', traceabilityVisible: boolean): void
  }>()

  const tableColumn = ref<TableColumnItemType[]>([])
  const userStore = useUserStore()
  const { getTableColumn } = userStore

  const detailVisible = ref(false)
  const traceabilityVisible = ref(false)
  const arrivedInfo = ref()
  //聚合检索必要条件
  const aggregation = reactive({
    scrollId: '',
    pageSize: 100,
    pageNum: 1,
    /** 排序字段 */
    orderField: 'startTimeNs',
    /** 排序方式 */
    orderType: 'decs',
    indexType: 10,
  })
  const aggregation_form = reactive({
    xff: '',
    host: '',
    url: '',
    query: '',
    readStatus: '',
    responseStatusCode: '',
    attackResult: '',
  })
  const aggregationList = ref<AlertItem[]>([])
  const aggregationTotal = ref(0)
  const expandKeys = ref<string[]>([])
  const isLoading = ref(false)
  const infoVal = ref()
  const moduleEnable = ref(false)

  const aggregationTableRef = ref<InstanceType<typeof ElTable> | null>(null)

  const columns = [
    {
      label: '告警时间',
      prop: 'startTimeNs',
      width: 150,
    },
    {
      label: 'URL',
      prop: 'url',
    },
    {
      label: 'Host',
      prop: 'host',
    },
    {
      label: 'Query',
      prop: 'query',
    },
    {
      label: '忽略状态',
      prop: 'ignoreStatus',
      width: 100,
    },
    {
      label: '状态码',
      prop: 'responseStatusCode',
      width: 80,
    },
    {
      label: '攻击结果',
      prop: 'attackResult',
      width: 90,
    },
  ]
  const attackResultType = ['成功', '失败', '未知', '企图']
  function showTraceability(row: any) {
    infoVal.value = row
    traceabilityVisible.value = true
  }
  const queryAlarmInfo = async (isClear = false) => {
    const search_sql = Object.entries(aggregation_form)
      .filter(([key, value]) => {
        return value
      })
      .map(([key, value]) => {
        const mapKey: { [key: string]: string } = {
          readStatus: '读取状态',
          attackResult: '攻击结果',
        }
        const keys = mapKey[key] || key
        return `${keys} = "${value}"`
      })
      .join(' and ')
    isLoading.value = true
    const { startTime, endTime } = props.selectAlert
    const {
      data: { resList, scrollId, total },
      code,
    } = await searchBySqlApi({
      ...aggregation,
      searchSql: search_sql
        ? `告警间隔聚合 = "${props?.selectAlert.alert_aggr_id}" and (${search_sql})`
        : `告警间隔聚合 = "${props?.selectAlert.alert_aggr_id}"`,
      startTime,
      endTime,
    })

    if (isClear) {
      aggregationList.value = resList || []
    } else {
      aggregationList.value.push(...(resList || []))
    }
    aggregation.scrollId = scrollId
    aggregationTotal.value = total
    isLoading.value = false
  }
  const clearHandle = () => {
    aggregation.scrollId = ''
    aggregation.pageNum = 1
  }
  const onSubmit = () => {
    clearHandle()
    queryAlarmInfo(true)
  }
  const expandNextHandle = (row: AlertItem, index: number) => {
    if (index === aggregationList.value.length - 1) return
    const currRow = aggregationList.value[index + 1]
    expandKeys.value = [currRow.id]
    scrollToView()
  }
  const expandPrevHandle = (row: AlertItem, index: number) => {
    if (index === 0) return
    const currRow = aggregationList.value[index - 1]
    expandKeys.value = [currRow.id]
    scrollToView()
  }
  const handleExpandChange = (row: any) => {
    expandKeys.value = expandKeys.value[0] === row.id ? [] : [row.id]
    scrollToView()
  }
  const getModelStatus = async () => {
    const {
      data: { module_config },
    } = await getSystemConfigApi({ keys: 'module_config' })
    moduleEnable.value = JSON.parse(module_config.value).isEnable || false
  }
  const scrollToView = () => {
    nextTick(() => {
      const expandedDom = document.querySelector('.el-table__expanded-cell')?.parentElement as HTMLTableRowElement
      const tableDom = document.querySelector('.my-aggregation-table .el-scrollbar__wrap')
      const expandedTop = expandedDom?.offsetTop - 80
      const tableScrollTop = tableDom!.scrollTop
      const scroll = expandedTop - tableScrollTop
      if (scroll < 100 && scroll > 0) return
      tableDom?.scrollTo({
        top: expandedTop,
        behavior: 'smooth',
      })
    })
  }
  async function formatColum() {
    const alertColumns = getTableColumn(10)
    const { data: userDisPlaysFiled } = await getAllDisPlaysFiledApi()
    const userColumnData = userDisPlaysFiled[10] as number[]
    if (userColumnData?.length > 0) {
      tableColumn.value = userColumnData.map((key) =>
        alertColumns.find((item) => item?.id === key)
      ) as TableColumnItemType[]
    } else {
      tableColumn.value = []
    }
  }
  watch(
    () => props.alertDetailVisible,
    () => {
      detailVisible.value = props.alertDetailVisible
      if (props.alertDetailVisible) {
        queryAlarmInfo(true)
        nextTick(() => {
          const domRef = ref(document.querySelector('.my-aggregation-table .el-scrollbar__wrap') as HTMLDivElement)
          const { arrivedState } = useScroll(domRef.value)
          arrivedInfo.value = arrivedState
        })
      } else {
        aggregation.scrollId = ''
        Object.assign(aggregation_form, {
          xff: '',
          host: '',
          url: '',
          query: '',
          readStatus: '',
          responseStatusCode: '',
          attackResult: '',
        })
      }
    },
    {
      immediate: true,
    }
  )

  watch(
    () => arrivedInfo.value?.bottom,
    () => {
      if (arrivedInfo.value?.bottom && !isLoading.value) {
        if (aggregationTotal.value === aggregationList.value.length) return
        aggregation.pageNum += 1
        queryAlarmInfo()
      }
    }
  )
  onMounted(() => {
    formatColum()
    getModelStatus()
  })
</script>

<script lang="ts">
  export default {
    name: 'AlertDetailAggregation',
  }
</script>

<template>
  <div class="alarm-aggregation-detail">
    <el-drawer
      v-model="detailVisible"
      size="85%"
      @close="
        () => {
          clearHandle()
          emits('update:alertDetailVisible', false)
        }
      "
    >
      <template #header="{ titleId, titleClass }">
        <h2 :id="titleId" :class="titleClass">威胁告警详情</h2>
      </template>
      <template #default>
        <div v-loading="isLoading" class="alarm-aggregation-detail-content">
          <div class="alarm-aggregation-detail-left">
            <dl class="alarm-aggregation-detail-header">
              <dt>
                <el-image :src="alarmIcon" style="width: 42px; height: 52px; margin-bottom: 4px" />
                <a class="alarm-btn">{{ selectAlert?.threatLevel || '-' }}</a>
              </dt>
              <dd>
                <label>威胁名称</label>
                <div @click="useCopy(selectAlert?.threatName)">{{ selectAlert?.threatName || '-' }}</div>
              </dd>
              <dd>
                <label>威胁类型</label>
                <div @click="useCopy(selectAlert?.threatType)">{{ selectAlert?.threatType || '-' }}</div>
              </dd>
              <dd>
                <label>源IP</label>
                <div style="display: ruby" @click="useCopy(selectAlert?.clientIp)">
                  <span>{{ selectAlert?.clientIp || '-' }}</span>
                  <span v-if="selectAlert?.client_country" style="white-space: nowrap">
                    <span>（</span>
                    <span>{{ selectAlert?.client_country }}</span>
                    <img
                      v-if="selectAlert?.client_country_code"
                      :src="require(`@/assets/flag/${selectAlert?.client_country_code.toLowerCase()}.png`)"
                      style="width: 20px; height: 15px; margin: 0 0 3px 4px"
                    />
                    <span>）</span>
                  </span>
                </div>
              </dd>
              <dd>
                <label>目的IP</label>
                <div style="display: ruby" @click="useCopy(selectAlert?.serverIp)">
                  <span>{{ selectAlert?.victimIp || '-' }}</span>
                  <span v-if="selectAlert?.server_country" style="white-space: nowrap">
                    <span>（</span>
                    <span>{{ selectAlert?.server_country }}</span>
                    <img
                      v-if="selectAlert?.server_country_code"
                      :src="require(`@/assets/flag/${selectAlert?.server_country_code.toLowerCase()}.png`)"
                      style="width: 20px; height: 15px; margin: 0 0 3px 4px"
                    />
                    <span>）</span>
                  </span>
                </div>
              </dd>
              <dd>
                <label>目的端口</label>
                <div @click="useCopy(selectAlert?.targetPort)">{{ selectAlert?.targetPort || '-' }}</div>
              </dd>
              <dd>
                <label>攻击阶段</label>
                <div @click="useCopy(selectAlert?.killchain)">{{ selectAlert?.killchain || '-' }}</div>
              </dd>
              <dd>
                <label>攻击次数</label>
                <div @click="useCopy(selectAlert?.attackCount)">{{ selectAlert?.attackCount || '-' }}</div>
              </dd>
              <dd>
                <label>方向</label>
                <div @click="useCopy(selectAlert?.position)">{{ selectAlert?.position || '-' }}</div>
              </dd>
              <dd v-if="selectAlert.vulnHarm">
                <label>风险危害</label>
                <el-tooltip
                  :content="selectAlert.vulnHarm"
                  placement="top"
                  popper-class="selectAlert-attackDesc"
                  :show-after="1000"
                >
                  <div class="word-content" @click="useCopy(selectAlert.vulnHarm)">{{ selectAlert.vulnHarm }}</div>
                </el-tooltip>
              </dd>
              <dd v-if="selectAlert.detailInfo">
                <label>威胁详情</label>
                <el-tooltip
                  :content="selectAlert.detailInfo"
                  placement="top"
                  popper-class="selectAlert-attackDesc"
                  :show-after="1000"
                >
                  <div class="word-content" @click="useCopy(selectAlert.detailInfo)">
                    {{ selectAlert.detailInfo }}
                  </div>
                </el-tooltip>
              </dd>
              <dd v-if="selectAlert.attackDesc">
                <label>威胁描述</label>
                <el-tooltip
                  :content="selectAlert.attackDesc"
                  placement="top"
                  popper-class="selectAlert-attackDesc"
                  :show-after="1000"
                >
                  <div class="word-content" @click="useCopy(selectAlert.attackDesc)">{{ selectAlert.attackDesc }}</div>
                </el-tooltip>
              </dd>
              <dd v-if="selectAlert.bulletin">
                <label>解决方案</label>
                <el-tooltip
                  :content="selectAlert.bulletin"
                  placement="top"
                  popper-class="selectAlert-attackDesc"
                  :show-after="1000"
                >
                  <div class="word-content" @click="useCopy(selectAlert.bulletin)">{{ selectAlert.bulletin }}</div>
                </el-tooltip>
              </dd>
            </dl>
          </div>
          <div class="alarm-aggregation-detail-right">
            <el-form
              class="alarm-aggregation-detail-form"
              :inline="true"
              :label-width="70"
              :model="aggregation_form"
              style="width: 100%"
            >
              <el-form-item label="XFF">
                <el-input v-model="aggregation_form.xff" clearable />
              </el-form-item>
              <el-form-item label="Host">
                <el-input v-model="aggregation_form.host" clearable />
              </el-form-item>
              <el-form-item label="URL">
                <el-input v-model="aggregation_form.url" clearable />
              </el-form-item>
              <el-form-item label="Query">
                <el-input v-model="aggregation_form.query" clearable />
              </el-form-item>
              <el-form-item label="读取状态">
                <el-select v-model="aggregation_form.readStatus" clearable>
                  <el-option label="已读" value="已读" />
                  <el-option label="未读" value="未读" />
                </el-select>
              </el-form-item>
              <el-form-item label="攻击结果">
                <el-select v-model="aggregation_form.attackResult" clearable>
                  <el-option v-for="item in attackResultType" :key="item" :label="item" :value="item" />
                </el-select>
              </el-form-item>
              <el-form-item class="search-btn" label="&nbsp;">
                <el-button type="primary" @click="onSubmit">检索</el-button>
              </el-form-item>
            </el-form>
            <div class="alarm-aggregation-detail-table">
              <el-table
                ref="aggregationTableRef"
                border
                class="my-aggregation-table"
                :data="aggregationList"
                :expand-row-keys="expandKeys"
                height="100%"
                row-key="id"
                @expand-change="handleExpandChange"
              >
                <el-table-column type="expand">
                  <template #default="{ row, $index }">
                    <alert-aggregation-detail-info
                      :alert-info="row"
                      :detail-total-index="aggregationList.length - 1"
                      :info-index="$index"
                      :module-enable="moduleEnable"
                      :table-column="tableColumn"
                      @next="expandNextHandle"
                      @prev="expandPrevHandle"
                    />
                  </template>
                </el-table-column>
                <el-table-column align="center" label="序号" type="index" width="55" />
                <el-table-column
                  v-for="column in columns"
                  :key="column.prop"
                  align="left"
                  :label="column.label"
                  :prop="column.prop"
                  show-overflow-tooltip
                  :width="column.width"
                >
                  <template v-if="column.label === '告警时间'" #default="{ row }">
                    {{ formatNstime(row[column.prop]) }}
                  </template>
                  <template #default="{ row }" v-else-if="column.prop == 'attackResult'">
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
              </el-table>
            </div>
          </div>
        </div>
      </template>
    </el-drawer>
    <traceability v-model:traceabilityVisible="traceabilityVisible" :select-alert="infoVal" />
  </div>
</template>

<style scoped lang="scss">
  .alarm-aggregation-detail {
    .alarm-aggregation-detail-header {
      overflow: hidden;
      display: flex;
      flex-direction: column;
      margin-top: 28px;
      .word-content {
        // line-height: 24px;
        overflow: hidden;
        // max-height: 48px;
        text-overflow: ellipsis;
        word-break: break-all;
        word-wrap: break-word;
        display: -webkit-box;
        // background-clip: text;
        appearance: none;
        line-clamp: 2;
        -webkit-line-clamp: 2;
        /* autoprefixer: ignore next */
        -webkit-box-orient: vertical;
      }
      dt {
        display: flex;
        flex-direction: column;
        align-items: center;
        margin-bottom: 30px;
        .alarm-btn {
          width: 56px;
          height: 24px;
          line-height: 20px;
          overflow: hidden;
          text-align: center;
          font-weight: 500;
          font-size: 14px;
          color: #ff4340;
          margin-top: 2px;
          border-radius: 12px;
          border: 2px solid transparent;
          background-image: linear-gradient(#fff0f0, #fff0f0),
            linear-gradient(166deg, rgba(255, 182, 182, 1), rgba(255, 184, 169, 1), rgba(255, 98, 98, 1));
          background-origin: border-box;
          background-clip: content-box, border-box;
        }
        .alarm-aggregation-detail-header-copy {
          font-size: 12px;
          color: #ffffff;
          width: 40px;
          height: 22px;
          border-radius: 4px;
          border: 1px solid #f8f7ff;
          display: inline-flex;
          position: absolute;
          top: 7px;
          right: 5px;
          text-indent: 0;
          text-indent: 8px;
          line-height: 20px;
        }
      }
      dd {
        display: flex;
        margin-left: 24px;
        text-align: left;
        label {
          font-weight: 500;
          font-size: 15px;
          color: #342e58;
          display: block;
          width: 80px;
          line-height: 26px;
        }
        div {
          flex: 1;
          overflow: hidden;
          display: -webkit-box;
          -webkit-box-orient: vertical;
          -webkit-line-clamp: 3;
          word-break: break-all;
          text-overflow: ellipsis;
          font-weight: 400;
          font-size: 15px;
          color: #4a4759;
          line-height: 26px;
          margin-right: 20px;
          cursor: pointer;
        }
      }
    }
    .alarm-aggregation-detail-content {
      display: flex;
      height: 100%;
      .alarm-aggregation-detail-left {
        width: 316px;
      }
      .alarm-aggregation-detail-right {
        flex: 1;
        background: #ffff;
        padding: 20px 30px;
        display: flex;
        flex-direction: column;
      }
      .alarm-aggregation-detail-table {
        width: 100%;
        flex: 1;
        margin-bottom: 13px;
        overflow-y: auto;
        .my-aggregation-table {
          position: relative;
          flex: 1;
          :deep() {
            .el-table__body-wrapper {
              .el-scrollbar__bar.is-vertical > div {
                margin-top: 40px;
              }
            }
            tr.el-table__row:has(
                > .el-table__expand-column .cell .el-table__expand-icon.el-table__expand-icon--expanded
              ) {
              position: sticky;
              top: 0;
              z-index: 999;
              .el-table__expand-column {
                border-left: var(--el-table-border);
              }
            }
          }
        }
      }
    }

    :deep() {
      .el-drawer {
        border-radius: 30px 0px 0px 30px;
        background: #534b89;
        .el-drawer__header {
          padding: 10px 30px 10px;
          margin-bottom: 0;
          .el-drawer__title {
            font-size: 20px;
            font-weight: 500;
            font-size: 20px;
            line-height: 34px;
            color: #fff;
          }
          .el-drawer__close-btn {
            width: 22px;
            height: 22px;
            border-radius: 50%;
            background: #eeedf9;
            display: flex;
            justify-content: center;
            align-items: center;
          }
        }
        .el-drawer__body {
          border-radius: 30px 0px 0px 30px;
          background: #f6f6fa;
          padding: 0;
        }
      }
      .alarm-aggregation-detail-form {
        .el-form-item {
          width: 24%;
          margin-right: 1%;
          &:nth-of-type(4) {
            width: 25%;
            margin-right: 0;
          }
          &.search-btn {
            width: 50% !important;
            margin-right: 0;
            .el-form-item__content {
              justify-content: end;
            }
          }
        }
      }
      .attackResult {
        .el-icon {
          font-size: 17px;
          vertical-align: -3px;
          margin-right: -4px;
        }
        &-0 {
          color: #ff4340;
        }
        &-1 {
          color: #9d9aba;
        }
        &-2 {
          color: #ffa515;
        }
        &-3 {
          color: #6954f0;
        }
      }
    }
  }
</style>
