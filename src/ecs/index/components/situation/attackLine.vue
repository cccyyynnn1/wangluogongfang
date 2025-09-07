<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'

  import numberFormatter from '~/src/utils/number'

  const props = defineProps<{
    attackLine?: object
  }>()

  const initOptions = {
    renderer: 'svg',
  }

  const assetsOption = reactive({
    initOptions: {
      renderer: 'svg',
    },
    option: {
      // color: ['#ABF8CA', '#92DCF8', '#AEACFF', '#E6B8FF', '#FEE6C4'],
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
      grid: {
        top: '4%',
        left: '3%',
        right: '4%',
        bottom: '4%',
        containLabel: true,
      },
      xAxis: [
        {
          type: 'category',
          boundaryGap: false,
          data: [] as string[],
          axisTick: {
            alignWithLabel: false,
          },
          axisLabel: {
            align: 'left',
          },
          // boundaryGap: ['5%', '95%'],
          splitLine: {
            show: false,
          },
        },
      ],
      yAxis: [
        {
          axisLabel: {
            formatter: function (value: number) {
              return numberFormatter.format(+value)
            },
          },
          axisTick: {
            show: false,
          },
          type: 'value',
          axisLine: {
            show: false,
          },
          splitLine: {
            lineStyle: {
              type: 'dashed',
              width: 1,
              opacity: '1',
            },
          },
        },
      ],
      series: [
        {
          name: '攻击趋势',
          type: 'line',
          symbol: 'circle',
          smooth: true,
          yAxisIndex: 0,
          stack: 'Total',
          lineStyle: {
            width: 1,
            color: 'rgb(255, 84, 84)',
          },
          showSymbol: false,
          areaStyle: {
            opacity: 0.72,
            color: new VabChart.graphic.LinearGradient(0, 0, 0, 1, [
              {
                offset: 0,
                color: 'rgb(255, 197, 200)',
              },
              {
                offset: 1,
                color: 'rgb(255, 255, 255)',
              },
            ]),
          },
          emphasis: {
            focus: 'series',
          },
          data: [] as string[],
        },
      ],
    },
  })

  watch(
    () => props.attackLine,
    () => {
      // @ts-ignore
      assetsOption.option.series[0].data = props.attackLine.count
      // @ts-ignore
      assetsOption.option.xAxis[0].data = props.attackLine.date
    }
  )
</script>

<script lang="ts">
  export default {
    name: 'AttackLine',
  }
</script>
<template>
  <vab-card class="echart" shadow="hover" skeleton>
    <template #header>
      <span class="title">攻击趋势</span>
    </template>
    <vab-chart
      class="echart-content"
      :init-options="assetsOption.initOptions"
      :option="assetsOption.option"
      theme="vab-echarts-theme"
    />
  </vab-card>
</template>

<style scoped lang="scss"></style>
