<script lang="ts">
  export default {
    name: 'RankingOfAsset ', // 资产威胁统计排行
  }
</script>

<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  // @ts-ignore
  import Pagination from './pagination'
  import { reactive, ref } from 'vue'
  import numberFormatter from '~/src/utils/number'

  const props = defineProps<{
    victimIps: any
    showBtn: boolean
  }>()

  const btnDisable = ref(false)

  const state = reactive({
    initOptions: {
      renderer: 'svg',
    },
    option: {
      title: {
        show: props.showBtn,
        text: '资产威胁统计排行',
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
        formatter: function (params: any) {
          return `${params[0].name}  :  ${params[0].value}`
        },
        confine: true,
      },
      grid: { left: 10, top: props.showBtn ? 30 : 0, bottom: 5, right: 15, containLabel: true },

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
        data: ['暂无数据...'],
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
          data: [0],
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
  const currentPage1 = ref(1)
  const pageSize1 = ref(7)
  const alldata = ref()
  const totals = ref(12)
  // const isShow = ref(true)

  // 聚合
  let polymerizationObj: any = {}
  let polymerizationFalg = false
  const handlePolymerization = () => {
    if (polymerizationFalg) return
    polymerizationFalg = true
    polymerizationObj = {}
    alldata.value.forEach((item: any) => {
      const arr = item.name.split('.')
      arr.pop()
      const str = arr.join('.')
      if (!Object.keys(polymerizationObj).includes(str)) {
        polymerizationObj[str] = []
      }
      polymerizationObj[str].push(item.name)
    })
    objSort()
    getData()
  }

  const objSort = () => {
    let arr = Object.keys(polymerizationObj)
    const dict = {}
    for (const iterator of arr) {
      // @ts-ignore
      dict[iterator] = polymerizationObj[iterator].length
    }
    const sort_values = Object.values(dict).sort((a: any, b: any) => b - a)
    //  console.log(sort_values)
    const sort_dict = {}
    for (let i of sort_values) {
      for (let key in dict) {
        //  console.log(i,key)
        // @ts-ignore
        if (dict[key] === i) {
          // @ts-ignore
          sort_dict[key] = i
          // @ts-ignore
          delete dict[key]
        }
      }
    }
    const polymerizationObj1 = {}
    for (const key in sort_dict) {
      // @ts-ignore
      polymerizationObj1[key] = polymerizationObj[key]
    }
    polymerizationObj = polymerizationObj1
  }

  const getData = () => {
    state.option.yAxis.data = []
    state.option.series[0].data = []
    state.option.series[1].data = []
    let max = 0
    let keyArr: string[] = []
    // let polymerizationArr: string[] = []
    if (polymerizationFalg) {
      keyArr = Object.keys(polymerizationObj)
      totals.value = Object.keys(polymerizationObj).length
    }
    const startNum = pageSize1.value * (currentPage1.value - 1)
    const endNum =
      pageSize1.value * currentPage1.value > totals.value ? totals.value : pageSize1.value * currentPage1.value
    if (polymerizationFalg) {
      max = polymerizationObj[keyArr[0]].length
      for (let index = startNum; index < endNum; index++) {
        // max = max > polymerizationObj[keyArr[index]].length ? max : polymerizationObj[keyArr[index]].length
        state.option.yAxis.data.push(`${keyArr[index]}.1/24`)
        state.option.series[0].data.push(polymerizationObj[keyArr[index]].length)
      }
      for (let index = startNum; index < endNum; index++) {
        state.option.series[1].data.push(max)
      }
    } else {
      alldata.value.forEach((_: any, index: any) => {
        if (index >= startNum && index < endNum) {
          state.option.yAxis.data.push(alldata.value[index].name)
          state.option.series[0].data.push(alldata.value[index].value)
          state.option.series[1].data.push(alldata.value[index].value)
        }
      })
    }
  }
  watch(
    () => props.victimIps,
    () => {
      alldata.value = props.victimIps.length > 0 ? props.victimIps : [{ name: '暂无数据...', value: 0 }]
      totals.value = alldata.value.length
      getData()
    }
  )

  watchEffect(() => {
    if (state.option.series[0].data.length === 1 || polymerizationFalg) {
      btnDisable.value = true
    } else {
      btnDisable.value = false
    }
  })

  watch(
    () => currentPage1.value,
    () => {
      getData()
    }
  )

  const handleCurrentChange = (val: number) => {
    currentPage1.value = val
  }

  const emit = defineEmits<{
    (e: 'on-click-bar', val: any): void
    (e: 'on-check-ips', val: any): void
    (e: 'on-reset'): void
  }>()

  let flag = true
  const hanldeClick = (e: any) => {
    if (!flag) return
    flag = false
    setTimeout(() => {
      flag = true
    }, 1000)
    let arr = []

    polymerizationFalg ? (arr = hanldeClickPolymerization(e.name)) : arr.push(e.name)
    polymerizationFalg = false
    emit('on-click-bar', arr)
  }

  const hanldeClickPolymerization = (key: string) => {
    return polymerizationObj[key.slice(0, -5)]
  }

  const resetKey = () => {
    polymerizationObj = {}
    polymerizationFalg = false
  }
  const handleReset = () => {
    resetKey()
    emit('on-reset')
  }
  // 检查当前IP组
  const handleCheckIPGroup = () => {
    if (state.option.yAxis.data.length == 0) return
    const arr: string[] = []
    alldata.value.forEach((item: any) => {
      arr.push(item.name)
    })
    resetKey()
    emit('on-check-ips', [...new Set(arr)])
  }

  defineExpose({
    resetKey,
    btnDisable,
  })
</script>

<template>
  <div class="ranking-of-asset">
    <vab-chart
      :click="hanldeClick"
      :init-options="state.initOptions"
      :option="state.option"
      style="height: 195px; width: 100%"
      theme="vab-echarts-theme"
    />
    <pagination
      :current-page="currentPage1"
      :page-size="pageSize1"
      :total="totals"
      @current-change="handleCurrentChange"
    />
    <div v-if="showBtn" class="contry-btn">
      <el-button class="btn" :disabled="btnDisable" link type="primary" @click="handlePolymerization">聚合</el-button>
      <el-button class="btn" :disabled="polymerizationFalg" link type="primary" @click="handleCheckIPGroup">
        检索IP组
      </el-button>
      <el-button class="btn" link type="primary" @click="handleReset">重置</el-button>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .ranking-of-asset {
    padding-top: 15px;
    width: 100%;
    height: 240px;
    position: relative;
    .contry-btn {
      position: absolute;
      top: 15px;
      z-index: 999;
      right: 10px;

      .btn {
        font-size: 12px;
        font-weight: 400;
      }
      // &:hover {
      //   cursor: pointer;
      // }
    }
  }
</style>
