<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'

  import numberFormatter from '~/src/utils/number'

  const props = defineProps<{
    index?: object
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
            alignWithLabel: true,
          },
          axisLabel: {
            align: 'left',
          },
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
          name: '应用数据量趋势',
          type: 'line',
          symbol: 'circle',
          smooth: true,
          yAxisIndex: 0,
          stack: 'Total',
          lineStyle: {
            width: 1,
            color: 'rgb(79, 148, 255)',
          },
          showSymbol: false,
          areaStyle: {
            opacity: 0.72,
            color: new VabChart.graphic.LinearGradient(0, 0, 0, 1, [
              {
                offset: 0,
                color: 'rgb(158, 198, 255)',
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

  // const getRandomColor = () => {
  //   const color = `rgb(${[
  //     Math.round(Math.random() * 255),
  //     Math.round(Math.random() * 255),
  //     Math.round(Math.random() * 255),
  //   ].join(',')})`
  //   return color
  // }

  // const getSeries = () => {
  //   assetsOption.option.series = []
  //   // @ts-ignore
  //   Object.keys(props.index).forEach((item: any) => {
  //     const meta = {
  //       name: '0',
  //       type: 'line',
  //       symbol: 'circle',
  //       smooth: true,
  //       yAxisIndex: 0,
  //       lineStyle: {
  //         width: 2,
  //         color: undefined,
  //       },
  //       showSymbol: false,
  //       areaStyle: {
  //         opacity: 0.8,
  //         color: undefined,
  //       },
  //       emphasis: {
  //         focus: 'series',
  //       },
  //       data: [] as string[],
  //     }
  //     const color = getRandomColor()
  //     meta.name = item
  //     // @ts-ignore
  //     meta.lineStyle.color = color
  //     meta.areaStyle.color = new VabChart.graphic.LinearGradient(0, 0, 0, 1, [
  //       {
  //         offset: 0,
  //         color: color,
  //       },
  //       {
  //         offset: 1,
  //         color: 'rgb(255, 255, 255)',
  //       },
  //     ])
  //     // @ts-ignore
  //     meta.data = props.index[item].count
  //     // @ts-ignore
  //     assetsOption.option.series.push(meta)
  //   })
  // }

  watch(
    () => props.index,
    () => {
      // getSeries()
      // @ts-ignore
      assetsOption.option.xAxis[0].data = props.index.date
      // @ts-ignore
      assetsOption.option.series[0].data = props.index.count
      // console.log(props.index)
    }
  )
</script>

<script lang="ts">
  export default {
    name: 'Index',
  }
</script>
<template>
  <vab-card class="echart" shadow="hover" skeleton>
    <template #header>
      <span class="title">应用数据量趋势</span>
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
