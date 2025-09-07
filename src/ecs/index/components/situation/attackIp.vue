<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'

  import wordcloud from 'echarts-wordcloud'

  wordcloud

  const props = defineProps<{
    attackIp: any
  }>()

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
          rotationRange: [0, 0],
          rotationStep: 45,
          gridSize: 8,
          shape: 'circle',
          width: '95%',
          height: '95%',
          textStyle: {
            color: function () {
              const color = `rgb(${[
                Math.round(Math.random() * 255),
                Math.round(Math.random() * 255),
                Math.round(Math.random() * 255),
              ].join(',')})`
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

  watch(
    () => props.attackIp,
    () => {
      keywords.value = props.attackIp.length > 0 ? props.attackIp : [{ name: '暂无数据', value: 0 }]
    }
  )
</script>

<script lang="ts">
  export default {
    name: 'AttackSource',
  }
</script>
<template>
  <vab-card class="echart" shadow="hover" skeleton>
    <template #header>
      <span class="title">攻击源</span>
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
