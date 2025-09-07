<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'

  import { getSiteDataVolume, getAccessStatisticsApi } from '~/src/api-ecs/dashboard'

  import { getAssetsListApi } from '~/src/api-ecs/assets'

  import numberFormatter from '~/src/utils/number'

  const props = defineProps<{
    clientIp: any
  }>()

  const initOptions = {
    renderer: 'svg',
  }

  // 访问统计
  const visitOption = reactive({
    // title: {
    //   // show:props.showBtn ,
    //   text: '访问统计',
    //   textStyle: {
    //     color: '#303133 ',
    //     fontSize: 14,
    //   },
    // },
    tooltip: {
      trigger: 'item',
      formatter: function (data: any) {
        return `${data.marker + data.data.name}: ${numberFormatter.format(data.data.value)}次`
      },
      confine: true,
    },
    series: [
      {
        name: '攻击类型',
        type: 'pie',
        radius: ['30%', '55%'],
        center: ['50%', '50%'],
        avoidLabelOverlap: false,
        padAngle: 2,
        itemStyle: {
          borderRadius: 6,
        },
        data: [],
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)',
          },
        },
        labelLine: {
          show: true,
        },
      },
    ],
  })

  // 站点数据量
  // const siteOption = reactive({
  //   title: {
  //     top: '0%',
  //     left: 'left',
  //     text: '站点数据量',
  //     textStyle: { color: '#333', fontStyle: 'normal', fontSize: 12 },
  //   },
  //   tooltip: {
  //     trigger: 'item',
  //     formatter: function (data: any) {
  //       return `${data.marker + data.data.name}:  ${numberFormatter.format(data.data.value)}`
  //     },
  //   },
  //   series: [
  //     {
  //       name: '站点数据量',
  //       type: 'pie',
  //       radius: ['30%', '55%'],
  //       center: ['50%', '55%'],
  //       data: [],
  //       emphasis: {
  //         itemStyle: {
  //           shadowBlur: 10,
  //           shadowOffsetX: 0,
  //           shadowColor: 'rgba(0, 0, 0, 0.5)',
  //         },
  //       },
  //       labelLine: {
  //         show: true,
  //       },
  //     },
  //   ],
  // })

  // 资产数据量
  // const assetsOption = reactive({
  //   title: {
  //     top: '0%',
  //     left: 'left',
  //     text: '资产数据量',
  //     textStyle: { color: '#333', fontStyle: 'normal', fontSize: 12 },
  //   },
  //   tooltip: {
  //     trigger: 'item',
  //     formatter: function (data: any) {
  //       return `${data.marker + data.data.name}:  ${numberFormatter.format(data.data.value)}`
  //     },
  //   },
  //   series: [
  //     {
  //       name: '资产数据量',
  //       type: 'pie',
  //       radius: ['30%', '55%'],
  //       center: ['50%', '55%'],
  //       data: [
  //         { value: undefined, name: '已知资产' },
  //         { value: undefined, name: '未知资产' },
  //       ],
  //       emphasis: {
  //         itemStyle: {
  //           shadowBlur: 10,
  //           shadowOffsetX: 0,
  //           shadowColor: 'rgba(0, 0, 0, 0.5)',
  //         },
  //       },
  //       labelLine: {
  //         show: true,
  //       },
  //     },
  //   ],
  // })

  watch(
    () => props.clientIp,
    () => {
      visitOption.series[0].data = Array.isArray(props.clientIp) ? props.clientIp : ([] as any)
    }
  )
</script>

<script lang="ts">
  export default {
    name: 'GeneralStatistics',
  }
</script>
<template>
  <vab-card class="echart" skeleton>
    <template #header>
      <span class="title">访问统计</span>
    </template>
    <!-- <el-row :gutter="20">
      <el-col :span="8"> -->
    <vab-chart
      v-if="clientIp?.length > 0"
      class="echart-content"
      :init-options="initOptions"
      :option="visitOption"
      theme="vab-echarts-theme"
    />
    <!-- </el-col>
      <el-col :span="8">
        <vab-chart class="echart-content" :init-options="initOptions" :option="siteOption" theme="vab-echarts-theme" />
      </el-col>
      <el-col :span="8">
        <vab-chart
          class="echart-content"
          :init-options="initOptions"
          :option="assetsOption"
          theme="vab-echarts-theme"
        />
      </el-col>
    </el-row> -->
    <el-empty v-if="clientIp?.length == 0" description="暂无数据" />
  </vab-card>
</template>

<style scoped lang="scss">
  :deep() {
    .el-empty {
      --el-empty-image-width: 80px;
    }
  }
</style>
