<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'

  import { getAttackSourceTypeApi } from '~/src/api-ecs/dashboard'

  import numberFormatter from '@/utils/number'

  const props = defineProps<{
    threatType: any
  }>()

  const state = reactive({
    initOptions: {
      renderer: 'svg',
    },
    option: {
      tooltip: {
        trigger: 'item',
        formatter: function (val: { value: number }) {
          return numberFormatter.format(+val.value)
        },
      },
      series: [
        {
          name: '攻击类型',
          type: 'pie',
          data: [],
          radius: ['40%', '70%'],
          minAngle: 10,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2,
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 20,
              fontWeight: 'bold',
            },
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)',
            },
          },
        },
      ],
    },
  })

  // onMounted(() => {
  //   initData()
  // })

  // // 初始化数据
  // const initData = async () => {
  //   const { data } = await getAttackSourceTypeApi()
  //   state.option.series[0].data = data.threatType
  // }

  watch(
    () => props.threatType,
    () => {
      state.option.series[0].data = props.threatType
    }
  )
</script>

<script lang="ts">
  export default {
    name: 'AttackType',
  }
</script>
<template>
  <vab-card class="echart" shadow="hover" skeleton>
    <template #header>
      <span class="title">攻击类型</span>
    </template>
    <vab-chart
      class="echart-content"
      :init-options="state.initOptions"
      :option="state.option"
      theme="vab-echarts-theme"
    />
  </vab-card>
</template>

<style scoped lang="scss"></style>
