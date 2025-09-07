<script lang="ts">
  export default {
    name: 'AttackAlerts',
  }
</script>

<script setup lang="ts">
  import { AlertItem } from '@/types'
  import { formatNstime, formatTime } from '@/utils/time'
  import { levelKey } from '@/ecs/alert/data/index'
  import { getAlertApi } from '@/api-ecs/alert'
  import AlertDetail from '@/ecs/alert/components/alert-detail.vue'
  const props = defineProps<{
    nodeData: any
    timeRange?: [string, string]
  }>()
  const infoVal = ref()
  // 自定义时间
  const timeDate = ref()
  const alertDetailVisible = ref(false)
  const alertInfo = reactive({
    list: [] as AlertItem[],
    total: 0,
    loading: true,
  })
  const arrivedInfo = ref()
  const must_parames = reactive({
    /** 排序字段 */
    orderField: 'startTimeNs',
    /** 排序方式 */
    orderType: 'decs',
    /** 索引类型 */
    indexType: 29,
    /** 当前页码 */
    pageNum: 1,
    /** 当前页数 */
    pageSize: 100,
    /** 开始时间 */
    startTime: '',
    /** 结束时间 */
    endTime: '',
    hasAggOneDocCount: true,
    scrollId: '',
  })
  // 获取表格序号
  const curIndex = computed(() => (must_parames.pageNum - 1) * must_parames.pageSize + 1)

  const obj = { 低危: 1, 中危: 2, 高危: 3, 危急: 4 }
  function getLevel(str: string) {
    // @ts-ignore
    const level = obj[str]
    return levelKey[level]
  }
  function showAlertDetail(row: any) {
    infoVal.value = row
    alertDetailVisible.value = true
  }
  // 查询告警日志
  async function queryData(hasMore = true) {
    alertInfo.loading = true
    const { ip, type, clientIp, serverIp } = props.nodeData
    try {
      const searchSql = type ? `源ip = "${clientIp}" and 目的ip = "${serverIp}"` : `源ip = "${ip}" or 目的ip = "${ip}"`
      const queryData = { ...must_parames, searchSql }
      const {
        data: { total, resList, scrollId },
      } = await getAlertApi(queryData)
      alertInfo.list = hasMore ? [...alertInfo.list, ...resList] : resList
      alertInfo.total = total
      alertInfo.loading = false
      must_parames.scrollId = scrollId
    } catch (error) {
      alertInfo.loading = false
    }
  }
  onMounted(() => {
    const dom = document.querySelector('.alert-relations .el-table__body-wrapper')! as HTMLDivElement
    const table_el = ref(dom)
    const { arrivedState } = useScroll(table_el.value)
    arrivedInfo.value = arrivedState
  })
  watch(
    () => arrivedInfo.value,
    () => {
      if (arrivedInfo.value.bottom) {
        queryData()
      }
    },
    { deep: true }
  )
  watchEffect(() => {
    if (props.timeRange) {
      timeDate.value = props.timeRange
    }
  })
  watchEffect(() => {
    if (timeDate.value) {
      const [startDate, endDate] = timeDate.value
      must_parames.startTime = startDate
      must_parames.endTime = endDate
      queryData(false)
    }
  })
</script>

<template>
  <el-row :gutter="20">
    <el-col :span="10">
      <vab-date-time-picker v-model="timeDate" />
    </el-col>
    <el-col :offset="10" :span="4">
      <el-button
        :loading="alertInfo.loading"
        style="float: right; margin-left: 20px"
        type="primary"
        @click="queryData()"
      >
        检索
      </el-button>
    </el-col>
  </el-row>
  <!-- 告警日志列表 -->
  <el-table
    v-loading="alertInfo.loading"
    class="alert-relations"
    :data="alertInfo.list"
    style="width: 100%; margin-top: 20px; height: 440x"
  >
    <el-table-column :index="curIndex" label="序号" show-overflow-tooltip type="index" width="75" />
    <el-table-column
      :formatter="({ startTimeNs }) => formatNstime(startTimeNs)"
      label="告警时间"
      prop="startTimeNs"
      show-overflow-tooltip
      width="180"
    />
    <el-table-column label="源IP" prop="attackIp" show-overflow-tooltip width="180" />
    <el-table-column label="目的IP" prop="victimIp" show-overflow-tooltip width="180" />
    <el-table-column label="XFF" prop="xff" show-overflow-tooltip width="120" />
    <el-table-column label="威胁类型" prop="threatType" show-overflow-tooltip width="120" />
    <el-table-column label="威胁名称" prop="threatName" show-overflow-tooltip width="180" />
    <el-table-column label="威胁等级" prop="threatLevel" show-overflow-tooltip width="120" />
    <el-table-column label="攻击结果" prop="attackResult" width="120" />
    <el-table-column fixed="right" label="操作" width="80">
      <template #default="{ row }">
        <el-button class="row_action" size="small" @click="showAlertDetail(row)">详情</el-button>
      </template>
    </el-table-column>
  </el-table>
  <alert-detail v-if="alertDetailVisible" v-model:alert-detail-visible="alertDetailVisible" :select-alert="infoVal" />
</template>

<style scoped lang="scss"></style>
