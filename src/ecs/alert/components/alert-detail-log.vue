<script lang="ts">
  export default {
    name: 'AlertDetailLog',
  }
  type LogTemplateType = 'Payload' | 'Http' | 'ThuApi'
</script>

<script setup lang="ts">
  import { Component } from 'vue'
  import PayloadLog from './payload-log.vue'
  import HttpLog from './http-log.vue'
  import ThuAlertApi from './thu-http-log.vue'

  const isPayload = ['cloud-insight']
  const isHttp = ['Skyeye—http', 'Threatbook—http', 'TUH']
  const isQinghuaApi = ['thu_api_alert']

  const props = defineProps<{
    logVal: any
  }>()

  const logTemplate: { [key in LogTemplateType]: Component } = {
    Payload: PayloadLog,
    Http: HttpLog,
    ThuApi: ThuAlertApi,
  }
  const componentTyps = Object.keys(logTemplate)
  const templateType = computed(() => {
    const index = [isPayload, isHttp, isQinghuaApi].findIndex((type) => type.includes(props.logVal.trade_name))
    const logComponent = componentTyps[index] as LogTemplateType
    return index > -1 ? logComponent : ''
  })
</script>

<template>
  <div class="logTemplate">
    <component :is="logTemplate[templateType]" v-if="templateType" :type="logVal.trade_name" :value="logVal" />
  </div>
</template>

<style scoped lang="scss">
  .logTemplate {
    :deep() {
      .hostTable_title {
        height: 14px;
        line-height: 14px;
        position: relative;
        text-indent: 0.8em;
        font-weight: 700;
        color: #303133;
        &::before {
          content: ' ';
          display: inline-block;
          position: absolute;
          width: 3px;
          height: 100%;
          left: 0;
          background: var(--el-color-primary);
          margin-right: 5px;
        }
      }
      .my-item {
        background: #f8fbff;
        padding: 20px;
        line-height: 20px;
        overflow: hidden;
        word-break: break-all;
        word-wrap: break-word;
        box-shadow: 0 0 0 1px var(--el-input-border-color, var(--el-border-color)) inset;
        border-radius: var(--el-input-border-radius, var(--el-border-radius-base));
        transition: var(--el-transition-box-shadow);
        white-space: pre-wrap;
        margin: 16px 0 26px;
        &:hover {
          cursor: text;
        }
      }
    }
  }
</style>
