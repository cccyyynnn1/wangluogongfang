<script setup lang="ts">
  import { FormInstance, UploadUserFile, UploadProps, UploadRawFile, UploadInstance } from 'element-plus'

  import { getToken } from '@/utils/token'

  const $baseMessage: any = inject('$baseMessage')

  const props = defineProps<{
    showUpload: boolean
  }>()

  const upload = ref<UploadInstance>()

  const formRef = ref<FormInstance>() // 表单实例

  const visible = ref(false) // 弹框显隐

  // 表单数据
  const formData = reactive<{
    file: UploadUserFile[]
  }>({
    file: [],
  })

  // 表单数据校验规则
  const rules = reactive({
    file: [
      {
        required: true,
        message: '请选择上传文件',
        trigger: 'change',
      },
    ],
  })

  const emit = defineEmits<{
    (e: 'on-close-event'): void
    (e: 'on-reflash'): void
  }>()

  // 关闭回调
  const handleClose = () => {
    emit('on-close-event')
  }

  const handleExceed: UploadProps['onExceed'] = (files, uploadFiles) => {
    formData.file[0].raw = files[0] as UploadRawFile
    formData.file[0].name = files[0].name
  }

  const handleAvatarSuccess: UploadProps['onSuccess'] = (response, uploadFile) => {
    if (response.code === 20) {
      ElMessage({ message: '上传成功', type: 'success' })
      emit('on-reflash')
      handleClose()
    } else {
      $baseMessage(response.msg, 'error', 'vab-hey-message-error')
    }
  }

  onMounted(() => {
    visible.value = props.showUpload
  })
</script>

<script lang="ts">
  export default {
    name: 'ImportItem',
  }
</script>
<template>
  <div class="import-assets">
    <el-dialog v-model="visible" :before-close="handleClose" title="导入" width="885px">
      <el-form ref="formRef" class="examplesDetailUpload-form" label-width="100px" :model="formData" :rules="rules">
        <el-form-item label="上传文件" prop="file">
          <el-upload
            ref="upload"
            v-model:file-list="formData.file"
            accept="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet, application/vnd.ms-excel"
            action="/v3/ecsPlatform/public/uploadThumbnail"
            :auto-upload="true"
            class="upload-demo"
            drag
            :headers="{ Cookies: `EcsSessionId=${getToken()}` }"
            :limit="1"
            :multiple="false"
            name="filedata"
            :on-exceed="handleExceed"
            :on-success="handleAvatarSuccess"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或
              <em>选择文件</em>
              <br />
              只能上传xls、xlsx格式文件
            </div>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <!-- <el-button :loading="isLoading" type="primary" @click="submitForm(formRef)">开始上传</el-button> -->
        <el-button @click="handleClose">取消</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .upload-demo {
    width: 745px;
  }
  .el-upload__tip {
    font-size: 20px;
    margin: -10px auto 0;
  }
</style>
