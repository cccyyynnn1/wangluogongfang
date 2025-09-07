<script lang="ts">
  export default {
    name: 'ThreatView', // 威胁视图
  }
</script>

<script setup lang="ts">
  // @ts-ignore
  import VabChart from '@/plugins/VabChart/index.vue'
  // @ts-ignore
  import Pagination from './pagination'
  import { reactive, ref } from 'vue'
  // @ts-ignore
  // import numberFormatter from '~/src/utils/number'
  import { detailChartsApi } from '@/api-ecs/assets-preview'
  import dayjs from 'dayjs'

  const props = defineProps<{
    serverIp: any
  }>()

  onMounted(() => {})

  const isShow = ref(true)

  const colorArr = ['#8BB3FF', '#FFAAAA', '#FFD26E']

  const state = reactive({
    initOptions: {
      renderer: 'svg',
    },
    option: {
      tooltip: {
        trigger: 'item',
        triggerOn: 'mousemove',
        formatter: (params: any) => {
          const lableArr = ['目的IP', '攻击IP', '威胁名称']
          const _index = colorArr.findIndex((item) => {
            return item == params.color
          })
          const label1 = lableArr[_index] ? `（${lableArr[_index]}）` : ''
          const str = `${params.name} ${label1} :  ${params.value}（数量）`
          return str
        },
      },
      // 10.209.0.188
      // 10.99.19.33
      series: [
        {
          type: 'sankey',
          data: [] as any[],
          links: [] as any[],
          emphasis: {
            focus: 'adjacency',
          },
          nodeAlign: 'left',
          levels: [
            {
              depth: 0,
              itemStyle: {
                color: colorArr[0],
              },
              lineStyle: {
                color: colorArr[0],
                opacity: 0.6,
              },
            },
            {
              depth: 1,
              itemStyle: {
                color: colorArr[1],
              },
              lineStyle: {
                color: colorArr[1],
                opacity: 0.6,
              },
            },
            {
              depth: 2,
              itemStyle: {
                color: colorArr[2],
              },
              lineStyle: {
                color: colorArr[1],
                opacity: 0.6,
              },
            },
          ],
          top: 10,
          lineStyle: {
            curveness: 0.5,
          },
          draggable: false,
          nodeGap: 20,
          nodeWidth: 20,
        },
      ],
    },
  })

  const node = ref()
  watch(
    () => props.serverIp,
    () => {
      if (props.serverIp) {
        getPageCharts()
      }
    }
  )

  // // 表格数据
  const timeParty: any = inject('timePartyProvide') as any
  const timeNow = dayjs()
  const endTime = timeNow.subtract(0, 'hour').format('YYYY-MM-DD HH:mm:ss')
  const startTime = timeNow.subtract((timeParty.value as number) || 6, 'hour').format('YYYY-MM-DD HH:mm:ss')
  const getPageCharts = async () => {
    state.option.series[0].data = []
    state.option.series[0].links = []

    const { data } = await detailChartsApi({ startTime, endTime, serverIps: props.serverIp })
    let total = 0
    const linkArr: string[] = []
    const nodeOBJArr: any[] = []
    for (const iterator of data.victimIp_attackIp_threatName) {
      total += iterator.value
    }
    for (const iterator of data.victimIp_attackIp_threatName) {
      const arr = iterator.name.split('~')
      arr.forEach((item: any) => {
        if (!linkArr.includes(item)) {
          linkArr.push(item)
        }
      })
      for (let index = 2; index >= 0; index--) {
        if (index == 2) {
          const obj = {
            source: '',
            target: '',
            value: 0,
          }
          obj.value = iterator.value
          obj.source = arr[index]
          obj.target = arr[index - 1]
          nodeOBJArr.push(obj)
        }
        if (index == 1) {
          const flag = nodeOBJArr.findIndex((item: any) => {
            return item.source == arr[index]
          })
          if (flag >= 0) {
            nodeOBJArr[flag].value += iterator.value
          } else {
            const obj = {
              source: '',
              target: '',
              value: 0,
            }
            obj.value = iterator.value
            obj.source = arr[index]
            obj.target = arr[index - 1]
            // obj.value = iterator.value / total
            nodeOBJArr.push(obj)
          }
        }
      }
      // console.log(nodeOBJArr)
    }
    linkArr.forEach((item) => {
      const obj = {
        name: '',
      }
      obj.name = item
      state.option.series[0].data.push(obj)
    })
    state.option.series[0].links = nodeOBJArr
    isShow.value = nodeOBJArr.length == 0 ? false : true
    setTimeout(() => {
      node.value = document.querySelector('.threat-view')
      if (node.value) {
        const targetDOMHeight = node.value?.querySelector('g')?.getBoundingClientRect()?.height
        node.value.style.height = `${targetDOMHeight + 400}px`
      }
    }, 0)
  }
</script>

<template>
  <vab-chart
    v-if="isShow"
    class="threat-view"
    :init-options="state.initOptions"
    :option="state.option"
    style="min-height: 630px; width: 100%"
    theme="vab-echarts-theme"
  />
  <el-empty v-else description="暂无数据" />
</template>

<style scoped lang="scss"></style>
