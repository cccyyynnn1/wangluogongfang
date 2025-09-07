<script lang="ts">
  export default {
    name: 'ChartsTotalY',
  }
</script>
<script setup lang="ts">
  // import { EChartsOption } from 'echarts'

  import VabChart from '@/plugins/VabChart/index.vue'

  import TopBar from './top-bar.vue'

  import { thirteenBitTimestamp } from '~/src/utils/index'

  import { calculateDiffTime } from '~/src/utils/time'

  type totalType = {
    ldle: number
    uptime: number | string
  }

  const props = defineProps<{
    total: totalType
  }>()

  const chartOption = reactive({
    // tooltip: {
    //   formatter: '{a} <br/>{b} : {c}%'
    // },
    title: {
      text: `总体运行时长：400天20小时50分钟`,
      top: '80%',
      left: '50%',
      textAlign: 'center',
      textStyle: {
        fontSize: 13,
        fontWeight: 'normal',
      },
    },
    series: [
      {
        name: 'Pressure',
        startAngle: 200,
        endAngle: -20,
        type: 'gauge',
        center: ['50%', '50%'],
        progress: {
          show: true,
          width: 20,
        },
        axisLine: {
          lineStyle: {
            width: 20, //柱子的宽度
          },
        },
        color: undefined as undefined | string,
        detail: {
          show: false,
        },
        axisTick: {
          show: false,
        },
        axisLabel: {
          show: false,
        },
        splitLine: {
          show: false,
        },
        pointer: {
          offsetCenter: [0, -10],
          width: 8, //指针的宽度
          length: '40%', //指针长度，按照半圆半径的百分比
        },
        data: [
          {
            value: 50,
            name: '正常',
          },
        ],
      },
    ],
  })

  const timer = ref('总体状态:2023年10月11日')

  watch(
    () => props.total,
    () => {
      timer.value = `总体状态:${thirteenBitTimestamp(Date.now())}`
      // @ts-ignore
      chartOption!.title!.text = `总体运行时长：${calculateDiffTime(Math.ceil(props.total.uptime))}`

      const num1 = props.total.ldle == 0 ? 100 : props.total.ldle

      const num: number = 100 - num1
      // @ts-ignore
      chartOption.series[0].data[0].value = num as number
      // @ts-ignore
      chartOption.series[0].data[0].name = `${num.toFixed(2)}%`
      // @ts-ignore
      chartOption.series[0].color = num > 50 || num == 0 ? undefined : undefined
      // visitOption.series[0].data = Array.isArray(props.clientIp) ? props.clientIp : ([] as any)
    },
    {
      deep: true,
      immediate: true,
    }
  )
</script>

<template>
  <top-bar :mode="'ChartsTotalY'" :title="timer" />
  <vab-chart :option="chartOption" theme="vab-echarts-theme" />
</template>

<style scoped lang="scss">
  .echarts {
    height: calc(100% - 53px);
    width: 100%;
  }
</style>
