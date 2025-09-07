<script lang="ts">
  export default {
    name: 'ChartsTotalBL',
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
      data: ['bytesReceive', 'bytesTransmit'],
    },
    series: [
      {
        name: 'bytesReceive',
        itemStyle: {
          opacity: 0,
        },
        data: [] as number[],
        type: 'line',
      },
      {
        name: 'bytesTransmit',
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
      chartOption.series[0].data.push(props.total.bytesReceive / 1024 / 1024)
      chartOption.series[1].data.push(props.total.bytesTransmit / 1024 / 1024)
      // chartOption.series[0].data.push(props.total.bytesReceive)
      // chartOption.series[1].data.push(props.total.bytesTransmit)
    },
    {
      deep: true,
      immediate: true,
    }
  )
</script>

<template>
  <top-bar :mode="'ChartsTotalBL'" title="网卡流量(M)" />
  <vab-chart :option="chartOption" theme="vab-echarts-theme" />
</template>

<style scoped lang="scss">
  .echarts {
    height: calc(100% - 53px);
    width: 100%;
  }
</style>
