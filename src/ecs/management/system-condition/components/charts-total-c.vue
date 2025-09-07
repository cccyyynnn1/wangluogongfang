<script lang="ts">
  export default {
    name: 'ChartsTotalC',
  }
</script>
<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'

  import TopBar from './top-bar.vue'

  const props = defineProps<{
    total: any
  }>()

  const chartOption = reactive({
    title: [
      {
        text: '使用率',
        textStyle: {
          color: '#36CBCB',
          fontSize: 18,
          fontWeight: 'normal',
        },
        itemGap: 20,
        left: 'center',
        top: '29%',
      },
      {
        text: '45%',
        textStyle: {
          color: '#36CBCB',
          fontSize: 18,
          fontWeight: 'normal',
        },
        itemGap: 20,
        left: 'center',
        top: '42%',
      },
    ],
    series: [
      {
        center: ['50%', '40%'],
        name: 'Access From',
        type: 'pie',
        radius: ['38%', '60%'],
        avoidLabelOverlap: false,
        label: {
          show: false,
          position: 'center',
        },
        emphasis: {
          scale: false,
          label: {
            show: false,
            fontSize: 40,
            fontWeight: 'bold',
          },
        },
        labelLine: {
          show: false,
        },
        data: [
          {
            value: 55,
            name: '总容量',
            itemStyle: {
              color: '#E6EBF8',
            },
          },
          {
            value: 45,
            name: '剩余容量',
            itemStyle: {
              color: '#36CBCB',
            },
          },
        ],
      },
    ],
  })
  watch(
    () => props.total,
    () => {
      if (props.total?.size) {
        const used = props.total?.used?.slice(0, -2)
        const unused = props.total?.size?.slice(0, -2) - used
        chartOption.series[0].data[0].value = Math.ceil(used)
        chartOption.series[0].data[1].value = Math.ceil(unused)
        chartOption.title[1].text = Math.ceil(props.total?.usePercent) + '%'
      }
    },
    {
      immediate: true,
      deep: true,
    }
  )
</script>

<template>
  <top-bar :mode="'ChartsTotalC'" title="磁盘使用率(%)" />
  <vab-chart :option="chartOption" theme="vab-echarts-theme" />
  <ul class="using">
    <li>
      {{ total?.size }}
      <br />
      <span>磁盘总量</span>
    </li>
    <li>
      {{ total?.used }}
      <br />
      <span>使用容量</span>
    </li>
    <li>
      {{ total.avail }}
      <br />
      <span>剩余容量</span>
    </li>
  </ul>
</template>

<style scoped lang="scss">
  .echarts {
    height: calc(100% - 58px);
    width: 100%;
  }
  .using {
    display: flex;
    padding-inline: 20px;
    margin-top: -50px;
    li {
      width: 33%;
      text-align: center;
    }
  }
</style>
