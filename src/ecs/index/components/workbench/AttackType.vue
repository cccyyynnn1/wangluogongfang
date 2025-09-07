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
          radius: ['35%', '65%'],
          data: [],
          emphasis: {
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
  const router = useRouter()
  const hanldeClick = (e: { name: string }) => {
    router.push({ path: '/alerts/alarm-attack', query: { threatType: e.name } })
  }
  const maxNUm = ref(0)
  watch(
    () => props.threatType,
    () => {
      state.option.series[0].data = props.threatType
      props.threatType.forEach((_: any, index: any) => {
        maxNUm.value = maxNUm.value >= props.threatType[index].value ? maxNUm.value : props.threatType[index].value
      })
    }
  )

  const handleRes = (item: any) => {
    // let res = undefined
    const num = maxNUm.value ? +((item.value / maxNUm.value) * 100).toFixed(0) : 0
    return num
  }
</script>

<script lang="ts">
  export default {
    name: 'AttackType',
  }
</script>
<template>
  <vab-card class="echart" skeleton>
    <template #header>
      <span class="title">攻击类型</span>
    </template>
    <div v-if="threatType?.length > 0" class="content">
      <div v-for="item in threatType" :key="item.name" class="item">
        <el-progress :color="'#FFD44D'" :percentage="handleRes(item)" :stroke-width="12" />
        <div class="item_text">
          <span>{{ item.name }}</span>
          <span style="color: #1e1841">{{ item.value }}</span>
        </div>
      </div>
    </div>
    <el-empty v-if="threatType?.length == 0" description="暂无数据" style="height: 80%; width: 100%" />
    <!-- <vab-chart
      class="echart-content"
      :click="hanldeClick"
      :init-options="state.initOptions"
      :option="state.option"
      theme="vab-echarts-theme"
    /> -->
  </vab-card>
</template>

<style scoped lang="scss">
  :deep() {
    .el-empty {
      --el-empty-image-width: 80px;
    }
  }

  .content {
    height: 100%;
    padding: 15px 24px;
    // display: flex;
    // flex-direction: column;
    // justify-content: space-between;
    overflow-y: auto;

    &::-webkit-scrollbar {
      width: 0;
      height: 0;
    }
    .item {
      height: 50px;
    }
    :deep() {
      .el-progress__text {
        display: none;
      }

      .el-progress-bar__outer {
        border-radius: 4px;
        background-color: #f5f4fe;
      }
    }
    .item_text {
      font-size: 12px;
      color: #88849d;
      margin-top: 8px;
      display: flex;
      // flex-direction: column;
      justify-content: space-between;
    }
  }
</style>
