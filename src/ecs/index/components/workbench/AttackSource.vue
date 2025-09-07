<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'

  import wordcloud from 'echarts-wordcloud'

  wordcloud

  const props = defineProps<{
    attackIp: any
  }>()

  const colorList = [
    '#F39EC2',
    '#5BD7A6',
    '#5A8FF8',
    '#6F5EF9',
    '#5FB1FB',
    '#5D7092',
    '#9860C0',
    '#EDC045',
    '#8AD8FF',
    '#4EA1A1',
  ]

  const keywords = ref<{ name: string; value: number }[]>([])

  const state = reactive({
    initOptions: {
      renderer: 'svg',
    },
    option: {
      tooltip: {
        show: true,
        trigger: 'item',
        confine: true,
      },
      series: [
        {
          type: 'wordCloud',
          // maskImage: maskImage,
          sizeRange: [12, 40],
          rotationRange: [0, 10],
          rotationStep: 45,
          gridSize: 18,
          shape: 'circle',
          width: '100%',
          height: '100%',
          textStyle: {
            color: function (val: any) {
              const color = colorList[val?.dataIndex % 10] || 0
              return color
            },
            fontWeight: 'normal',
            emphasis: {
              shadowBlur: 10,
              shadowColor: '#333',
            },
          },
          data: keywords,
        },
      ],
    },
  })
  const router = useRouter()
  const hanldeClick = (e: { name: string }) => {
    if (!e.name || e.name == '暂无数据') return
    router.push({ path: '/alerts/alarm-attack', query: { attackIp: e.name } })
  }

  watch(
    () => props.attackIp,
    () => {
      keywords.value = props.attackIp?.length > 0 ? props.attackIp : []
    }
  )
</script>

<script lang="ts">
  export default {
    name: 'AttackSource',
  }
</script>
<template>
  <vab-card class="echart" skeleton>
    <template #header>
      <span class="title">攻击源</span>
    </template>
    <vab-chart
      v-if="attackIp?.length > 0"
      class="echart-content"
      :click="hanldeClick"
      :init-options="state.initOptions"
      :option="state.option"
      theme="vab-echarts-theme"
    />
    <el-empty v-if="attackIp?.length == 0" description="暂无数据" />
  </vab-card>
</template>

<style scoped lang="scss">
  :deep() {
    .el-empty {
      --el-empty-image-width: 80px;
    }
  }
</style>
