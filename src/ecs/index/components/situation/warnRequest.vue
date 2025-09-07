<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'

  import numberFormatter from '~/src/utils/number'

  const props = defineProps<{
    warnRequest?: object
  }>()

  const initOptions = {
    renderer: 'svg',
  }

  // 资产数据量
  const assetsOption = reactive({
    color: ['#469CFF', '#FCD338'],
    tooltip: {
      trigger: 'item',
      formatter: function (data: any) {
        return `${data.marker + data.data.name}:  ${numberFormatter.format(data.data.value)}`
      },
    },
    legend: {
      bottom: '0',
      left: 'center',
    },
    series: [
      {
        // name: '资产数据量',
        type: 'pie',
        radius: ['0', '50%'],
        center: ['50%', '45%'],
        // itemStyle: {
        // borderRadius: 10,
        // borderWidth: 2,
        // color: function (params:any) {
        //   const colorList = ['#76AAFF','#FF797F'];
        //   return colorList[params.dataIndex]
        // },
        // },
        data: [
          { value: undefined, name: '正常访问' },
          { value: undefined, name: '异常访问' },
        ],
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)',
          },
        },
        labelLine: {
          show: false,
        },
      },
    ],
  })

  watch(
    () => props.warnRequest,
    () => {
      // @ts-ignore
      const sum = props.warnRequest['1'] + props.warnRequest['0']
      // @ts-ignore
      const normal = props.warnRequest['1'] || 0
      // @ts-ignore
      const unnormal = props.warnRequest['0'] || 0
      assetsOption.series[0].data[0].value = normal
      assetsOption.series[0].data[1].value = unnormal
      assetsOption.series[0].data[0].name = `正常访问 ${Number((normal / sum).toFixed(2)) * 100}%`
      assetsOption.series[0].data[1].name = `异常访问 ${Number((unnormal / sum).toFixed(2)) * 100}%`
    }
  )
</script>

<script lang="ts">
  export default {
    name: 'WarnRequest',
  }
</script>
<template>
  <vab-card class="echart" shadow="hover" skeleton>
    <template #header>
      <span class="title">异常访问占比</span>
    </template>
    <vab-chart class="echart-content" :init-options="initOptions" :option="assetsOption" theme="vab-echarts-theme" />
  </vab-card>
</template>

<style scoped lang="scss"></style>
