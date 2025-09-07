<script lang="ts">
  export default {
    name: 'ThreatDistribution', // 疑似威胁分布
  }
</script>

<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  import { reactive, watch } from 'vue'
  import numberFormatter from '~/src/utils/number'

  const props = defineProps<{
    distributions: any
  }>()

  const state = reactive({
    initOptions: {
      renderer: 'svg',
    },
    option: {
      title: {
        text: '疑似威胁分布',
        textStyle: {
          color: '#303133 ',
          fontSize: 14,
        },
      },
      tooltip: {
        trigger: 'item',
        formatter: function (val: { value: number }) {
          return numberFormatter.format(+val.value)
        },
        textStyle: {
          fontSize: 14,
        },
      },
      legend: {
        orient: 'vertical',
        height: 160,
        icon: 'circle',
        top: 'center',
        left: '40%',
        itemGap: 11,
        itemHeight: 10,
        itemWidth: 10,
        textStyle: {
          ellipsis: '...',
          width: 120,
          overflow: 'truncate',
        },
      },
      series: [
        {
          name: '威胁分布',
          type: 'pie',
          radius: ['40%', '60%'],
          center: ['20%', '46%'],
          avoidLabelOverlap: false,
          padAngle: 5,
          itemStyle: {
            borderRadius: 8,
            borderColor: '#fff',
            borderWidth: 2,
          },
          label: {
            show: false,
            position: 'center',
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 16,
              // fontWeight: 'bold',
            },
          },
          labelLine: {
            show: false,
          },
          data: [
            { value: 0, name: '数据加载中... ' },
            { value: 0, name: '数据加载中...' },
          ],
        },
      ],
    },
  })

  watch(
    () => props.distributions,
    () => {
      const aggObj = props.distributions.length > 0 ? props.distributions : [{ name: '暂无数据', value: 0 }]
      const alldata = aggObj.filter((item: any) => {
        return !!item.name
      })
      state.option.series[0].data = []
      for (const iterator of alldata) {
        state.option.series[0].data.push(iterator)
      }
      // getData()
    }
  )
</script>

<template>
  <vab-chart
    class="threat-distribution"
    :init-options="state.initOptions"
    :option="state.option"
    theme="vab-echarts-theme"
  />
</template>

<style scoped lang="scss">
  .threat-distribution {
    padding-top: 15px;
    height: 220px;
    width: 100%;
  }
</style>
