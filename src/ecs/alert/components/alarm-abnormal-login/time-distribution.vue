<script lang="ts">
  export default {
    name: 'TimeDistribution', // 时间分布
  }
</script>

<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  import { AbnormalLandingAlarmDistributionOfTimeNApi } from '~/src/api-ecs/alert'
  import { AbnormalLandingAlarmDistributionOfTimeType } from '~/src/types'
  import numberFormatter from '~/src/utils/number'

  let max = 0
  let min = 0
  const state = reactive({
    initOptions: {
      renderer: 'svg',
    },
    option: {
      color: ['#7FEBDB', '#74F4A4', '#7D75FF', '#5A8FF8'],
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          show: true,
          type: 'cross',
          lineStyle: {
            type: 'dashed',
            width: 1,
          },
        },
      },
      // toolbox: {
      //   show: true,
      //   feature: {
      //     mark: { show: true },
      //     dataView: { show: true, readOnly: false },
      //     restore: { show: true },
      //     saveAsImage: { show: true },
      //   },
      // },
      dataZoom: {
        show: true,
        start: 0,
        end: 100,
        fillerColor: '#6655e7',
        height: '20px',
        handleColor: '#ffffff',
        backgroundColor: '#F2F0FF',
        handleStyle: {
          borderColor: '#6655e7',
        },
        realtime: true,
        showDataShadow: false,
        moveHandleStyle: {
          show: false,
          borderWidth: 0,
          color: '#fff',
          opacity: 0,
        },
        y: 170,
      },
      grid: {
        top: 10,
        left: '2%',
        right: 10,
        bottom: 30,
        containLabel: true,
      },
      xAxis: [
        {
          type: 'time',
          splitNumber: 10,
        },
      ],
      yAxis: [
        {
          type: 'value',
        },
      ],
      animation: false,
      series: [
        {
          name: '',
          type: 'scatter',
          symbolSize: function (value: any) {
            if (value[1]) {
              return (value[1] / max) * 15 > 3 ? (value[1] / max) * 15 : 3
            }
          },
          data: [] as any[],
        },
        {
          name: '',
          type: 'scatter',
          symbolSize: function (value: any) {
            if (value[1]) {
              return (value[1] / max) * 15 > 3 ? (value[1] / max) * 15 : 3
            }
          },
          data: [] as any[],
        },
        {
          name: '',
          type: 'scatter',
          symbolSize: function (value: any) {
            if (value[1]) {
              return (value[1] / max) * 15 > 3 ? (value[1] / max) * 15 : 3
            }
          },
          data: [] as any[],
        },
        {
          name: '',
          type: 'scatter',
          symbolSize: function (value: any) {
            if (value[1]) {
              return (value[1] / max) * 15 > 3 ? (value[1] / max) * 15 : 3
            }
          },
          data: [] as any[],
        },
      ],
    },
  })

  const initData = async (option: AbnormalLandingAlarmDistributionOfTimeType) => {
    const { data } = await AbnormalLandingAlarmDistributionOfTimeNApi({ ...option })
    state.option.series[0].data = []
    state.option.series[1].data = []
    state.option.series[2].data = []
    state.option.series[3].data = []
    min = data[0][1]
    data.forEach((_: any, index: number) => {
      // data[index][1] = +(Math.random() * 150).toFixed(0)
      const num = index % 4

      switch (num) {
        case 0:
          state.option.series[0].data.push(data[index])
          break
        case 1:
          state.option.series[1].data.push(data[index])
          break
        case 2:
          state.option.series[2].data.push(data[index])
          break
        case 3:
          state.option.series[3].data.push(data[index])
          break
        default:
          break
      }
    })
    data.forEach((item: any) => {
      max = max >= item[1] ? max : item[1]
      min = min <= item[1] ? min : item[1]
    })
    // data.forEach((_: any, index: number) => {
    //   if (data[index][1] >= min && data[index][1] < min + (max - min) * 0.25) {
    //     state.option.series[0].data.push(data[index])
    //   } else if (data[index][1] >= min + (max - min) * 0.25 && data[index][1] < min + (max - min) * 0.5) {
    //     state.option.series[1].data.push(data[index])
    //   } else if (data[index][1] >= min + (max - min) * 0.5 && data[index][1] < min + (max - min) * 0.75) {
    //     state.option.series[2].data.push(data[index])
    //   } else if (data[index][1] >= min + (max - min) * 0.75 && data[index][1] <= max) {
    //     state.option.series[3].data.push(data[index])
    //   }
    // })
  }

  defineExpose({
    initData,
  })
</script>

<template>
  <div class="time-distribution">
    <div class="tools">
      <div class="tools-left">时间分布</div>
    </div>
    <vab-chart
      :init-options="state.initOptions"
      :option="state.option"
      style="height: 200px; width: 100%"
      theme="vab-echarts-theme"
    />
  </div>
</template>

<style scoped lang="scss">
  .time-distribution {
    .tools {
      padding: 0 15px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      .tools-left {
        font-weight: 500;
        font-size: 14px;
        color: #303133;
      }
    }
    padding-top: 12px;
    width: 100%;
    height: 240px;
  }
</style>
