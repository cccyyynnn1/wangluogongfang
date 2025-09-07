<script lang="ts">
  export default {
    name: 'LogDetails',
  }
</script>

<script setup lang="ts">
  import VabJsonViewer from 'vue-json-viewer'
  import { AuditSystemLog } from '@/types'
  import { formatNstime } from '@/utils/time'
  import { downloadSourceLog } from '@/utils/download'
  const activeName = ref('logDetails') // tabs选中项

  const props = defineProps<{
    isShow: boolean
    infoVal: AuditSystemLog
  }>()

  const emit = defineEmits<{
    (e: 'on-saveData', val: boolean): void
  }>()
  const formatSourceData = (val: string) => {
    if (!val) return ''
    let _val = ''
    try {
      _val = JSON.parse(val)
    } catch (error) {
      _val = val
    }
    return _val
  }
  const saveData = () => {
    emit('on-saveData', true)
  }
  const downloadLog = () => {
    const name = `${props.infoVal.logTypeStr}-${formatNstime(props.infoVal.timeStamp)}`
    downloadSourceLog(props.infoVal.sourceData, name)
  }
</script>

<template>
  <div class="log-details">
    <el-tabs v-model="activeName">
      <el-tab-pane label="操作日志" name="logDetails">
        <el-button class="btn" type="primary" @click="saveData">返回</el-button>
        <div class="item">发生时间：{{ formatNstime(infoVal.timeStamp) }}</div>
        <div class="item">日志类型：{{ infoVal.logTypeStr }}</div>
        <div class="log item">
          <div class="word">原始日志：</div>
          <div class="warp">
            <vab-json-viewer :expand-depth="5" preview-mode sort :value="formatSourceData(infoVal.sourceData)" />
          </div>
        </div>
        <div class="download">
          <el-button type="primary" @click="downloadLog">下载</el-button>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>
<style lang="scss" scoped>
  // :deep() {
  //   .el-scrollbar {
  //     height: calc(100vh - 220px) !important;
  //   }
  // }
  :deep() {
    .el-table__body-wrapper {
      max-height: 380px;
      min-height: 380px;
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0 !important;
        height: 0;
      }
    }
  }
  .log-details {
    .btn {
      margin: 10px 0 0 0px;
    }

    .item {
      height: 20px;
      margin-top: 15px;
      line-height: 20px;
    }

    .log {
      display: flex;
      height: 500px;

      .warp {
        flex: 1;
        height: 100%;
        background: #f5f7fa;
        .jv-container {
          background-color: inherit;
        }
      }
    }

    .download {
      position: relative;
      height: 30px;
      margin-top: 20px;

      .el-button {
        position: absolute;
        bottom: 0px;
        left: 50%;
        transform: translateX(-50%);
      }
    }
  }
</style>
