<script lang="ts">
  export default {
    name: 'SourceIP', // 源IP
  }
</script>

<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  import { AbnormalLandingAlarmSrcIpOrAddressIpTopNApi } from '~/src/api-ecs/alert'
  import { AbnormalLandingAlarmSrcIpOrAddressIpTopNType } from '~/src/types/alert'
  import numberFormatter from '~/src/utils/number'

  const emits = defineEmits<{
    (e: 'refresh'): void
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
        formatter: function (params: any) {
          return `${params[0].name}  :  ${params[0].value}`
        },
        confine: true,
      },
      grid: { left: 10, top: 0, bottom: 5, right: 15, containLabel: true },

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
        offset: -5,
        splitLine: { lineStyle: { color: ['#CEEDFF'], type: [5, 8], dashOffset: 3 } },
      },
      yAxis: {
        type: 'category',
        data: [] as number[],
        axisLine: { show: true, lineStyle: { color: '#ccc' } },
        axisTick: { length: 3 },
        splitLine: {
          show: false,
        },
        axisLabel: { show: true, fontSize: 12, color: '#666', margin: 12, padding: 0 },
        inverse: true,
      },
      series: [
        {
          name: '',
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
          data: [],
        },
        {
          name: '数量',
          type: 'bar',
          showBackground: false,
          // backgroundStyle: { color: 'rgba(82, 168, 255, 0.1)', borderRadius: [0, 8, 8, 0] },
          itemStyle: {
            color: 'rgba(82, 168, 255, 0.0)',
            normal: {
              borderRadius: [0, 8, 8, 0],
              color: 'rgba(82, 168, 255, 0.0)',
            },
            opacity: 0,
          },
          label: {
            show: false,
          },
          barGap: '-100%',
          barMaxWidth: 10,
          data: [0],
        },
      ],
    },
  })

  const selectValue = ref<'attackIp' | ' victimIp'>('attackIp')

  const initData = async (option: AbnormalLandingAlarmSrcIpOrAddressIpTopNType) => {
    const { data } = await AbnormalLandingAlarmSrcIpOrAddressIpTopNApi({ ...option })
    const alldata: any[] = data.aggObj || []
    state.option.series[0].data = []
    state.option.series[1].data = []
    state.option.yAxis.data = []
    alldata.forEach((_: any, index: any) => {
      // @ts-ignore
      state.option.yAxis.data.push(alldata[index].name)
      // @ts-ignore
      state.option.series[0].data.push(alldata[index].value)
      // @ts-ignore
      state.option.series[1].data.push(alldata[index].value)
    })
    if (state.option.series[1].data.length == 0) {
      state.option.series[1].data = [0]
    }
  }

  const handleChange = () => {
    emits('refresh')
  }

  defineExpose({
    selectValue: selectValue,
    initData,
  })
</script>

<template>
  <div class="source-ip">
    <div class="tools">
      <div class="tools-left">源IP</div>
      <el-select v-model="selectValue" size="small" style="width: 60px" @change="handleChange">
        <el-option label="源IP" value="attackIp" />
        <el-option label="目的IP" value="victimIp" />
        <!-- <el-option value="world" label="世界" /> -->
      </el-select>
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
  .source-ip {
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
