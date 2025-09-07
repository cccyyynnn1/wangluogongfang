<script lang="ts">
  export default {
    name: 'ChartsTotalL',
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
    series: [
      {
        name: 'Memory Used Utilization',
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
      // const used = props.total.total?.slice(0, -2) - props.total.free?.slice(0, -2)
      chartOption.series[0].data.push(props.total.rate)
    },
    {
      deep: true,
      immediate: true,
    }
  )
</script>

<template>
  <top-bar :mode="'ChartsTotalL'" title="内存使用率(%)" />
  <vab-chart :option="chartOption" theme="vab-echarts-theme" />
  <ul class="using">
    <li>
      {{ total.total }}
      <br />
      <span>内存总量</span>
    </li>
    <li>
      {{ total.free }}
      <br />
      <span>剩余内存量</span>
    </li>
  </ul>
</template>

<style scoped lang="scss">
  .echarts {
    height: calc(100% - 53px);
    width: 100%;
  }
  .using {
    display: flex;
    padding-inline: 20px;
    margin-top: -50px;
    li {
      width: 50%;
      text-align: center;
    }
  }
</style>
