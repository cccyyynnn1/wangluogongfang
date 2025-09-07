<script lang="ts">
  export default {
    name: 'ThreateningIncident', // 威胁分类分布
  }
</script>

<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  import { reactive } from 'vue'
  import numberFormatter from '~/src/utils/number'

  const props = defineProps<{
    threatTypes: any
  }>()

  const state = reactive({
    initOptions: {
      renderer: 'svg',
    },
    option: {
      title: {
        text: '威胁分类分布',
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
        type: 'scroll',
        height: 190,
        icon: 'circle',
        top: 0,
        right: '5%',
        itemGap: 11,
        itemHeight: 10,
        itemWidth: 10,
        data: [] as any[],
        selected: {},
        textStyle: {
          ellipsis: '...',
          width: 150,
          overflow: 'truncate',
        },
        pageButtonPosition: 'end',
        pageIconSize: [10, 8],
        pageIconColor: '#6954f0',
        selectorLabel: {
          color: '#7667ea',
          fontSize: 10,
          borderRadius: 2,
          borderColor: '#EEE',
          rich: {
            fontSize: 6,
          },
        },
        selector: [
          {
            // 全选
            type: 'all',
            // 可以是任意你喜欢的标题
            title: '全选',
          },
          {
            // 反选
            type: 'inverse',
            // 可以是任意你喜欢的标题
            title: '反选',
          },
        ],
      },
      series: [
        {
          name: '威胁事件分布',
          type: 'pie',
          radius: ['55%', '75%'],
          center: ['30%', '50%'],
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
    () => props.threatTypes,
    () => {
      const aggObj = props.threatTypes.length > 0 ? props.threatTypes : [{ name: '暂无数据', value: 0 }]
      const alldata = aggObj.filter((item: any) => {
        return !!item.name
      })
      state.option.series[0].data = []
      state.option.legend.data = []
      for (const iterator of alldata) {
        state.option.series[0].data.push(iterator)
        state.option.legend.data.push(iterator.name)
      }
      reset()
    }
  )

  const reset = () => {
    for (const key of state.option.legend.data) {
      // @ts-ignore
      state.option.legend.selected[key] = true
    }
  }

  const emit = defineEmits<{
    (e: 'on-click-pie', val: string): void
  }>()

  let flag = true
  const hanldeClick = (e: any) => {
    if (!flag) return
    flag = false
    setTimeout(() => {
      flag = true
    }, 1000)
    emit('on-click-pie', e.data?.name)
  }
</script>

<template>
  <vab-chart
    class="threatening-incident"
    :click="hanldeClick"
    :init-options="state.initOptions"
    :option="state.option"
    theme="vab-echarts-theme"
  />
</template>

<style scoped lang="scss">
  .threatening-incident {
    padding-top: 15px;
    height: 240px;
    width: 100%;
  }
</style>
