<script lang="ts">
  export default {
    name: 'CpuUsageRatesBL',
  }
</script>
<script setup lang="ts">
  // import { EChartsOption } from 'echarts'

  import VabChart from '@/plugins/VabChart/index.vue'

  import TopBar from './top-bar.vue'
  import { formatTimeToHSM } from '~/src/utils/time'

  const props = defineProps<{
    total: any
  }>()

  const chartOption = reactive({
    xAxis: {
      type: 'category',
      data: [] as string[],
    },
    yAxis: {
      type: 'value',
    },
    grid: {
      left: 50,
      top: 30,
      bottom: 90,
    },
    legend: {
      type: 'plain',
      bottom: 25,
      itemWidth: 10,
      itemHeight: 10,
      itemGap: 14,
      data: ['IO使用率'],
    },
    series: [
      {
        name: 'IO使用率',
        itemStyle: {
          opacity: 0,
        },
        data: [] as number[],
        type: 'line',
      },
    ],
  })
  const timer = ref('00:00:00')

  watch(
    () => props.total,
    () => {
      timer.value = `${formatTimeToHSM(Date.now())}`
      chartOption.xAxis.data.push(timer.value)
      const util = props.total?.reduce((l: any, r: any) => l + r.util, 0)
      chartOption.series[0].data.push(util / props.total.length)
    },
    {
      deep: true,
      immediate: true,
    }
  )
</script>

<template>
  <top-bar :mode="'CpuUsageRatesBL'" title="IO使用率(%)" />
  <vab-chart :option="chartOption" theme="vab-echarts-theme" />
</template>

<style scoped lang="scss">
  .echarts {
    height: calc(100% - 53px);
    width: 100%;
  }
</style>
