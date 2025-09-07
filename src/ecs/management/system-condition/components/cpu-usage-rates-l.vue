<script lang="ts">
  export default {
    name: 'CpuUsageRatesL',
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
      data: ['CPU Total'],
    },
    series: [
      {
        name: 'CPU Total',
        data: [] as number[],
        type: 'line',
        itemStyle: {
          opacity: 0,
        },
      },
    ],
  })
  const timer = ref('00:00:00')
  watch(
    () => props.total,
    () => {
      timer.value = `${formatTimeToHSM(Date.now())}`
      chartOption.xAxis.data.push(timer.value)
      chartOption.series[0].data.push(props.total.usage)
    },
    {
      deep: true,
      immediate: true,
    }
  )
</script>

<template>
  <top-bar :mode="'CpuUsageRatesL'" title="CPU使用率(%)" />
  <vab-chart :option="chartOption" theme="vab-echarts-theme" />
</template>

<style scoped lang="scss">
  .echarts {
    height: calc(100% - 53px);
    width: 100%;
  }
</style>
