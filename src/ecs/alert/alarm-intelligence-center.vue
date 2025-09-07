<script lang="ts">
  export default {
    name: 'AlarmIntelligenceCenter',
  }
</script>

<script setup lang="ts">
  import { InfoCustomOptions } from '~/src/types'
  import IntelligenceCenterCloud from './components/intelligence-center-cloud.vue'
  import IntelligenceCenterCustom from './components/intelligence-center-custom.vue'
  import IntelligenceCenterWhite from './components/intelligence-center-white.vue'
  import { getInfoCustomOptionsApi } from '@/api-ecs/alert'
  const infoCustomOptions = reactive<{
    [key in 'threatLevel' | 'iocType' | 'reliable' | 'threatType']: InfoCustomOptions
  }>({
    threatLevel: [],
    iocType: [],
    reliable: [],
    threatType: [],
  })
  const activeName = ref('cloud')
  // const setBg = () => {
  //   document.querySelector('.vab-app-main')?.setAttribute('style', 'background-color:#010123; transition:none')
  //   document.querySelector('.vab-app-main section')?.setAttribute('style', 'background-color:#010123;transition:none')
  //   document.querySelector('#app')?.setAttribute('style', 'background-color:#010123;transition:none')
  // }

  const getInfoCustomOptionsHandle = async () => {
    const { data } = await getInfoCustomOptionsApi()
    infoCustomOptions.iocType = data.info_ioc_type || []
    infoCustomOptions.reliable = data.info_reliable || []
    infoCustomOptions.threatType = data.info_threat_type || []
    infoCustomOptions.threatLevel = data.alertThreatLevel_code_cn || []
  }
  // watchEffect(() => {
  //   if (activeName.value === 'cloud') {
  //     setBg()
  //   } else {
  //     document.querySelector('.vab-app-main')?.setAttribute('style', 'transition:none')
  //     document.querySelector('.vab-app-main section')?.setAttribute('style', 'transition:none')
  //     document.querySelector('#app')?.setAttribute('style', 'transition:none')
  //   }
  // })
  onMounted(() => {
    // setBg()
    getInfoCustomOptionsHandle()
  })
  // onBeforeRouteLeave(() => {
  //   document.querySelector('.vab-app-main')?.removeAttribute('style')
  //   document.querySelector('.vab-app-main section')?.removeAttribute('style')
  //   document.querySelector('#app')?.removeAttribute('style')
  // })
</script>

<template>
  <div class="alarm-intelligence-center" :class="{ hasBg: activeName === 'cloud' }">
    <el-tabs v-model="activeName" class="alarm-intelligence-center-tabs">
      <el-tab-pane label="情报检索" name="cloud" />
      <el-tab-pane label="自定义情报" name="custom" />
      <el-tab-pane label="情报白名单" name="white" />
    </el-tabs>
    <intelligence-center-cloud v-if="activeName === 'cloud'" :info-custom-options="infoCustomOptions" />
    <intelligence-center-custom v-else-if="activeName === 'custom'" :info-custom-options="infoCustomOptions" />
    <intelligence-center-white v-else :info-custom-options="infoCustomOptions" />
  </div>
</template>

<style scoped lang="scss">
  .alarm-intelligence-center {
    height: calc(100vh);
    padding-top: 10px;
    display: flex;
    flex-direction: column;
    position: relative;
    // inset: 0 -10px;
    &.hasBg {
      &::after {
        content: ' ';
        display: block;
        position: fixed;
        width: 1px;
        height: 100vh;
        background-color: #eeeef0 !important;
        opacity: 0.13;
        left: auto;
        top: 0;
      }
      background: url('@/assets/alert_images/alarm-intelligence.png') no-repeat center;
      background-size: cover;
      .alarm-intelligence-center-tabs {
        --el-text-color-primary: #ffffffa9;
        --el-color-primary: #fff;
      }
      :deep() {
        .el-tabs__nav-wrap::after {
          background-color: #eeeef0 !important;
          opacity: 0.13;
        }
      }
    }
    .alarm-intelligence-center-tabs {
      padding-left: 35px;
      --el-text-color-primary: #303133;
      --el-color-primary: #6954f0;
      :deep() {
        .ioc-dialog.el-dialog,
        .ioc-whitelistdialog.el-dialog {
          .ioc-form {
            margin-right: 50px;
            margin-top: 10px;
            .el-input,
            .el-select :deep(.el-input),
            .el-textarea {
              width: 536px;
            }
            .el-form-item {
              margin-bottom: 16px;
            }
          }
          .dialog-footer {
            margin-right: 50px;
            margin-bottom: 10px;
          }
        }
      }
    }
  }
</style>
