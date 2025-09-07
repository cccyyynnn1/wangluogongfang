<script lang="ts">
  export default {
    name: 'ThreatTrends', // 威胁趋势
  }
</script>

<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  import { reactive } from 'vue'
  import numberFormatter from '~/src/utils/number'

  const props = defineProps<{
    threatLevels: any
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
      legend: {
        data: ['危急', '高危', '中危'],
        left: 'center',
        bottom: 5,
        itemHeight: 10,
        itemGap: 25,
        itemWidth: 18,
      },
      grid: {
        top: 26,
        left: 10,
        right: 10,
        bottom: 15,
        containLabel: true,
      },
      xAxis: [
        {
          type: 'category',
          boundaryGap: true,
          show: false,
          // data: ['数据加载中...'] as string[],
          // axisTick: {
          //   alignWithLabel: true,
          // },
          // axisLabel: {
          //   fontSize: 12,
          // },
          // splitLine: {
          //   show: false,
          // },
        },
        {
          type: 'category',
          boundaryGap: true,
          data: ['数据加载中...'] as string[],
          axisTick: {
            show: false,
          },
          // max: 10,
          position: 'bottom',
          axisLabel: {
            fontSize: 12,
            // padding: [0, 0, 0, -100],
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
          name: '危急',
          type: 'line',
          lineStyle: {
            width: 1.5,
            color: '#ca0a08',
          },
          smooth: true,
          symbol: 'circle',
          showSymbol: false,
          symbolSize: 2,
          symbolColor: '#ca0a08',
          itemStyle: {
            opacity: 1,
            color: '#ca0a08',
            borderWidth: 1,
          },
          data: [0],
        },
        {
          name: '高危',
          type: 'line',
          showAllSymbol: false,
          lineStyle: {
            width: 1.5,
            color: '#ff2927 ',
          },
          showSymbol: false,
          symbolSize: 2,
          symbolColor: '#ff2927',
          smooth: true,
          symbol: 'circle',
          itemStyle: {
            opacity: 1,
            color: '#ff2927',
            borderWidth: 1,
          },
          data: [0],
        },
        {
          name: '中危',
          type: 'line',
          lineStyle: {
            width: 1.5,
            color: '#FFDA36',
          },
          smooth: true,
          symbol: 'circle',
          showSymbol: false,
          symbolSize: 2,
          symbolColor: '#FFDA36',
          itemStyle: {
            opacity: 1,
            color: '#FFDA36',
            borderWidth: 1,
          },
          data: [0],
        },
      ],
    },
  })

  const alldata = ref()
  let code: any[] = []
  watch(
    () => props.threatLevels,
    () => {
      if (props.threatLevels.length == 0) [{ name: '数据加载中...', value: 0 }]
      code = []
      alldata.value = ['', '']
      props.threatLevels.forEach((item: any) => {
        if (item.name == '高危') {
          alldata.value[1] = item
        } else if (item.name == '中危') {
          alldata.value[2] = item
        } else if (item.name == '危急') {
          alldata.value[0] = item
        }
      })
      // code.push(alldata.value[0].threatLevelCode, alldata.value[1].threatLevelCode)
      state.option.series[0].data = alldata.value[0].value.count
      state.option.series[1].data = alldata.value[1].value.count
      state.option.series[2].data = alldata.value[2].value.count
      state.option.xAxis[0].data = alldata.value[0].value.date
      state.option.xAxis[1].data =
        alldata.value[0].value.date.length > 3
          ? [
              alldata.value[0].value.date[0],
              alldata.value[0].value.date[Math.floor(alldata.value[0].value.date.length / 2)],
              alldata.value[0].value.date.pop(),
            ]
          : alldata.value[0].value.date
    }
  )

  const emit = defineEmits<{
    (e: 'on-click-event', val: number): void
  }>()

  // let flag = true
  // const hanldeClick = (e: any) => {
  //   if (!flag) return
  //   const _index = e.seriesName == '高危' ? 0 : 1
  //   const codevalue = code[_index]
  //   flag = false
  //   setTimeout(() => {
  //     flag = true
  //   }, 1000)
  //   emit('on-click-event', codevalue)
  // }
</script>

<template>
  <!-- :click="hanldeClick" -->
  <vab-chart class="threat-trends" :init-options="state.initOptions" :option="state.option" theme="vab-echarts-theme" />
</template>

<style scoped lang="scss">
  .threat-trends {
    width: 100%;
    height: 240px;
  }
</style>
