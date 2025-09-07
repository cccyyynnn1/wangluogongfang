<script lang="ts">
  export default {
    name: 'UpdateConfig',
  }
</script>

<script setup lang="ts">
  import { getToken } from '@/utils/token'
  import type { UploadProps, UploadInstance } from 'element-plus'
  import { checkVerificationApi } from '~/src/api-ecs/system'
  import AesEncryptCBC from '~/src/utils/crypto'
  const $baseMessage: any = inject('$baseMessage')
  const upload = ref<UploadInstance>()

  const loading = ref(false)
  function handleError() {
    loading.value = false
    upload.value!.clearFiles()
  }
  function handleSuccess($val: any) {
    if ($val.code === 50) {
      $baseMessage($val.msg, 'error', 'vab-hey-message-error')
      upload.value!.clearFiles()
    } else {
      $baseMessage($val.msg, 'success', 'vab-hey-message-success')
    }
    loading.value = false
  }
  const beforeAvatarUpload: UploadProps['beforeUpload'] = async (rawFile) => {
    if (rawFile.size / 1024 / 1024 > 2000) {
      $baseMessage('文件大小不能超出2000M', 'error', 'vab-hey-message-error')
      return false
    }
    loading.value = true
    return true
  }
</script>

<template>
  <div v-loading="loading" class="update-container" element-loading-text="系统更新中...">
    <el-upload
      ref="upload"
      accept=".zip,.ZIP"
      action="/v3/eht/update/upload"
      :auto-upload="true"
      :before-upload="beforeAvatarUpload"
      drag
      element-loading-background="rgba(0,0,0,0.1)"
      element-loading-spinner="el-icon-loading"
      element-loading-text="拼命上传中"
      :headers="{ Cookies: `EcsSessionId=${getToken()}` }"
      :limit="1"
      multiple
      name="filedata"
      :on-error="handleError"
      :on-success="handleSuccess"
      style="width: 70%; margin: 80px auto 0"
    >
      <i class="el-icon-upload"></i>
      <div class="el-upload__text">
        将更新文件拖到此处，或
        <em>点击上传</em>
      </div>
      <div class="el-upload__tip">只能上传zip文件，且不超过2000M</div>
    </el-upload>
    <el-alert :closable="false">
      <h3>更新说明</h3>
      <p>更新文件须按照格式加密压缩</p>
      <p>更新文件时禁止进行其他操作</p>
    </el-alert>
  </div>
</template>

<style lang="scss" scoped>
  .el-tab-pane,
  .el-alert {
    width: 70%;
    margin: 50px auto !important;
  }
  :deep(.el-upload-dragger) {
    padding: 80px 0;
  }
</style>
