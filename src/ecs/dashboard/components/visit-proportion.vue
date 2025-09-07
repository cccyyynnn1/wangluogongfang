<script lang="ts">
  export default {
    name: 'VisitProportion',
  }
</script>

<script setup lang="ts">
  const props = defineProps<{
    data: {
      normal?: number
      abnormal?: number
    }
  }>()
  const normal = computed(() => {
    const { normal, abnormal } = props.data
    const count = normal! + abnormal!
    if (count == 0) {
      return '未知'
    } else {
      if (isNaN(normal!) || isNaN(abnormal!)) return '未知'
      const num = (normal! / count) * 100
      return num === 100 ? 100 : Number(num.toFixed(2))
    }
  })
  const abnormal = computed(() => {
    if (typeof normal.value === 'string') return '未知'
    return 100 - normal.value
  })
</script>

<template>
  <div class="visit-proportion-container">
    <div class="visit-normal">
      <h4>
        {{ normal === '未知' ? normal : `${normal.toFixed(2)}%` }}
      </h4>
      <h5 class="visit-text">正常访问</h5>
    </div>
    <div class="visit-abnormal">
      <h4>{{ abnormal === '未知' ? abnormal : `${abnormal.toFixed(2)}%` }}</h4>
      <h5 class="visit-text">异常访问</h5>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .visit-proportion-container {
    display: flex;
    justify-content: space-evenly;
    align-items: center;
    .visit-normal,
    .visit-abnormal {
      width: 310px;
      height: 400px;
      text-align: center;
      display: flex;
      justify-content: flex-end;
      flex-direction: column;
      h4,
      h5 {
        margin: 0;
      }
      h4 {
        font-size: 54px;
        color: #0084ff;
      }
      h5 {
        font-size: 30px;
        font-weight: 400;
        color: #8badc0;
        margin: 0 0 30px;
      }
    }
    .visit-abnormal {
      background: url('@/assets/dashboard_images/abnormal.png') no-repeat center;
      h4 {
        color: #fe1a43;
      }
    }
    .visit-normal {
      background: url('@/assets/dashboard_images/normal.png') no-repeat center;
    }
  }
</style>
