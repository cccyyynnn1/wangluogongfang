<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  const props = defineProps<{
    visible: boolean
    title: string
  }>()
  const emits = defineEmits<{
    (e: 'update:visible', visible: boolean): void
  }>()
  const site_chart = ref<InstanceType<typeof VabChart>>()
  const detailVisible = ref(false)

  function changeVisible() {
    emits('update:visible', false)
  }

  function download() {
    if (site_chart.value) {
      try {
        const dataUrl = site_chart.value.getDataURL({
          type: 'png',
          backgroundColor: '#F3F9FF',
        })
        const downloadBtn = document.createElement('a')
        downloadBtn.download = `${props.title}-站点资产`
        downloadBtn.style.display = 'none'
        downloadBtn.href = dataUrl
        downloadBtn.click()
      } catch (error) {
        console.error(error)
      }
    }
  }

  watchEffect(() => {
    detailVisible.value = props.visible
  })
</script>

<script lang="ts">
  export default {
    name: 'AssetsSiteDetail',
  }
</script>

<template>
  <div class="site-detail">
    <el-dialog
      v-model="detailVisible"
      :close-on-click-modal="false"
      destroy-on-close
      :title="props.title"
      width="1120px"
      @close="changeVisible"
    >
      <div class="site-action">
        <span>姓名：{{ '李长生' }}</span>
        <span>联系电话：{{ 11011011011 }}</span>
        <el-button @click="download">
          <vab-icon icon="download-2-line" style="font-size: 14px" />
          下载到本地
        </el-button>
      </div>
      <div class="siteChart">
        <vab-chart ref="site_chart" class="charts" theme="vab-echarts-theme" />
      </div>
      <template #footer>
        <el-button type="primary" @click="changeVisible">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .siteChart {
    width: 100%;
    height: 448px;
    background: #f3f9ff;
    margin-top: 20px;

    .charts {
      width: 100%;
      height: 448px;
    }
  }

  .site-action {
    display: flex;
    align-items: center;

    span {
      margin-right: 20px;
    }

    .el-button {
      margin-left: auto;
    }
  }
</style>
