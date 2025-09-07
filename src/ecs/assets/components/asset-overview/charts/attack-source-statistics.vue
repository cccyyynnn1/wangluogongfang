<script lang="ts">
  export default {
    name: 'AttackSourceStatistics', // 攻击源统计
  }
</script>
<script setup lang="ts">
  // @ts-ignore
  import VabChart from '@/plugins/VabChart/index.vue'
  // @ts-ignore
  import Pagination from './pagination'
  import { reactive, ref, watch } from 'vue'
  // @ts-ignore
  import numberFormatter from '~/src/utils/number'

  const props = defineProps<{
    attackSources: any
  }>()

  const state = reactive({
    initOptions: {
      renderer: 'svg',
    },
    option: {
      title: {
        text: '攻击源统计',
        textStyle: {
          color: '#303133 ',
          fontSize: 14,
        },
      },
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
      grid: { left: 0, top: 30, bottom: 5, right: 15, containLabel: true },

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
        data: ['数据加载中...'],
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
          name: '数量',
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
          // , 7, 6, 3
          data: [0],
        },
      ],
    },
  })
  const currentPage1 = ref(1)
  const pageSize1 = ref(7)
  const alldata = ref()
  const totals = ref(12)

  watch(
    () => props.attackSources,
    () => {
      alldata.value = props.attackSources.length > 0 ? props.attackSources : [{ name: '数据加载中...', value: 0 }]
      totals.value = alldata.value.length
      getData()
    }
  )

  const getData = () => {
    state.option.yAxis.data = []
    state.option.series[0].data = []
    const startNum = pageSize1.value * (currentPage1.value - 1)
    const endNum =
      pageSize1.value * currentPage1.value > totals.value ? totals.value : pageSize1.value * currentPage1.value
    alldata.value.forEach((_: any, index: any) => {
      if (index >= startNum && index < endNum) {
        state.option.yAxis.data.push(alldata.value[index].name)
        state.option.series[0].data.push(alldata.value[index].value)
      }
    })
  }

  watch(
    () => currentPage1.value,
    () => {
      getData()
    }
  )

  const handleCurrentChange = (val: number) => {
    currentPage1.value = val
  }
</script>

<template>
  <div class="system-distribution">
    <vab-chart
      :init-options="state.initOptions"
      :option="state.option"
      style="height: 185px; width: 100%"
      theme="vab-echarts-theme"
    />
    <pagination
      :current-page="currentPage1"
      :page-size="pageSize1"
      :total="totals"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<style scoped lang="scss">
  .system-distribution {
    padding-top: 15px;
    width: 100%;
    height: 220px;
  }
</style>
