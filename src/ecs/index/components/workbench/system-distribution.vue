<script lang="ts">
  export default {
    name: 'SystemDistribution', // 系统分布
  }
</script>

<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  // @ts-ignore
  import Pagination from '@/ecs/assets/components/asset-overview/charts/pagination.vue'
  import { reactive } from 'vue'
  import numberFormatter from '~/src/utils/number'

  const props = defineProps<{
    sysTypes: any
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
      grid: { left: 0, top: 12, bottom: 10, right: 10, containLabel: true },

      xAxis: {
        type: 'value',
        boundaryGap: false,
        axisLine: { show: false, lineStyle: { color: '#ccc' } },
        axisTick: { show: false },
        position: 'top',
        axisLabel: {
          color: '#999',
          formatter: function (value: number) {
            return numberFormatter.format(+value)
          },
        },
        splitLine: { lineStyle: { color: ['#CEEDFF'], type: [5, 8], dashOffset: 3 } },
      },
      yAxis: {
        type: 'category',
        data: ['暂无数据'] as string[],
        axisLine: { show: true, lineStyle: { color: '#ccc' } },
        axisTick: { length: 3 },
        splitLine: {
          show: false,
        },
        axisLabel: {
          show: true,
          fontSize: 12,
          color: '#666',
          margin: 12,
          padding: 0,
        },
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
          data: [] as number[],
        },
      ],
    },
  })
  // const currentPage1 = ref(1)
  const pageSize1 = ref(5)
  const alldata = ref()
  const totals = ref(7)
  watch(
    () => props.sysTypes,
    () => {
      alldata.value = props.sysTypes?.length > 0 ? props.sysTypes : []
      totals.value = alldata.value?.length
      getData()
    }
  )

  const curData = ref()
  const maxNUm = ref(0)
  const getData = () => {
    state.option.yAxis.data = []
    state.option.series[0].data = []
    curData.value = []
    // const startNum = pageSize1.value * (currentPage1.value - 1)
    // const endNum =
    //   pageSize1.value * currentPage1.value > totals.value ? totals.value : pageSize1.value * currentPage1.value
    alldata.value.forEach((_: any, index: any) => {
      maxNUm.value = maxNUm.value >= alldata.value[index].value ? maxNUm.value : alldata.value[index].value
      // if (index >= startNum && index < endNum) {
      const obj = {
        name: alldata.value[index].name,
        value: alldata.value[index].value,
      }
      curData.value.push(obj)
      state.option.yAxis.data.push(alldata.value[index].name)
      state.option.series[0].data.push(alldata.value[index].value)
      // }
    })
  }

  // watch(
  //   () => currentPage1.value,
  //   () => {
  //     getData()
  //   }
  // )
  // const handleCurrentChange = (val: number) => {
  //   currentPage1.value = val
  // }
  const handleRes = (item: any) => {
    // let res = undefined
    const num = maxNUm.value ? +((item.value / maxNUm.value) * 100).toFixed(0) : 0
    return num
  }
</script>

<template>
  <vab-card class="echart system-distribution" skeleton>
    <template #header>
      <span class="title">系统分布</span>
      <!-- <pagination
        :current-page="currentPage1"
        :page-size="pageSize1"
        :total="totals"
        @current-change="handleCurrentChange"
      /> -->
    </template>
    <div v-if="curData?.length > 0" class="content">
      <div v-for="item in curData" :key="item.name" class="item">
        <el-progress :color="'#7CDFFF'" :percentage="handleRes(item)" :stroke-width="12" />
        <div class="item_text">
          <span>{{ item.name }}</span>
          <span style="color: #1e1841">{{ item.value }}</span>
        </div>
      </div>
    </div>
    <el-empty v-if="curData?.length == 0" description="暂无数据" />
    <!-- <vab-chart
      :init-options="state.initOptions"
      :option="state.option"
      style="height: 210px; width: 100%"
      theme="vab-echarts-theme"
    /> -->
  </vab-card>
</template>

<style scoped lang="scss">
  .system-distribution {
    :deep() {
      .el-empty {
        --el-empty-image-width: 80px;
      }
    }
    width: 100%;
    .content {
      height: 100%;
      padding: 15px 24px;
      // display: flex;
      // flex-direction: column;
      // justify-content: space-between;
      overflow-y: auto;

      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
      .item {
        height: 50px;
      }
      :deep() {
        .el-progress__text {
          display: none;
        }
        .el-progress-bar__outer {
          border-radius: 4px;
          background-color: #f5f4fe;
        }
      }
      .item_text {
        font-size: 12px;
        color: #88849d;
        margin-top: 8px;
        display: flex;
        // flex-direction: column;
        justify-content: space-between;
      }
    }
  }
</style>
