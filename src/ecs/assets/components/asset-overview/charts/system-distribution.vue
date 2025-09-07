<script lang="ts">
  export default {
    name: 'SystemDistribution', // 系统分布
  }
</script>

<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  // @ts-ignore
  import Pagination from './pagination'
  import { reactive } from 'vue'
  import numberFormatter from '~/src/utils/number'

  const props = defineProps<{
    sysTypes: any
  }>()

  const state = reactive({
    initOptions: {
      renderer: 'svg',
    },
    option: {
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
      grid: { left: 0, top: 12, bottom: 10, right: 10, containLabel: true },

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
        data: ['暂无数据'] as string[],
        axisLine: { show: true, lineStyle: { color: '#ccc' } },
        axisTick: { length: 3 },
        splitLine: {
          show: false,
        },
        axisLabel: {
          show: true,
          fontSize: 12,
          color: '#666',
          margin: 12,
          padding: 0,
        },
        inverse: true,
      },
      series: [
        {
          name: '数量',
          type: 'bar',
          showBackground: true,
          backgroundStyle: { color: 'rgba(82, 168, 255, 0.1)', borderRadius: [0, 8, 8, 0] },
          itemStyle: {
            color: '#52A8FF',
            normal: {
              borderRadius: [0, 8, 8, 0],
              color: '#4F68FF',
            },
          },
          barMaxWidth: 10,
          data: [0] as number[],
        },
      ],
    },
  })
  const currentPage1 = ref(1)
  const pageSize1 = ref(5)
  const alldata = ref()
  const totals = ref(7)
  watch(
    () => props.sysTypes,
    () => {
      alldata.value = props.sysTypes.length > 0 ? props.sysTypes : [{ name: '数据加载中...', value: 0 }]
      totals.value = alldata.value.length
      getData()
    }
  )

  const getData = () => {
    state.option.yAxis.data = []
    state.option.series[0].data = []
    const startNum = pageSize1.value * (currentPage1.value - 1)
    const endNum =
      pageSize1.value * currentPage1.value > totals.value ? totals.value : pageSize1.value * currentPage1.value
    alldata.value.forEach((_: any, index: any) => {
      if (index >= startNum && index < endNum) {
        state.option.yAxis.data.push(alldata.value[index].name)
        state.option.series[0].data.push(alldata.value[index].value)
      }
    })
  }

  watch(
    () => currentPage1.value,
    () => {
      getData()
    }
  )
  const handleCurrentChange = (val: number) => {
    currentPage1.value = val
  }
</script>

<template>
  <div class="system-distribution">
    <vab-chart
      :init-options="state.initOptions"
      :option="state.option"
      style="height: 210px; width: 100%"
      theme="vab-echarts-theme"
    />
    <pagination
      :current-page="currentPage1"
      :page-size="pageSize1"
      :total="totals"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<style scoped lang="scss">
  .system-distribution {
    width: 100%;
    height: 240px;
  }
</style>
