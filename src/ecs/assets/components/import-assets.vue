<script setup lang="ts">
  import { FormInstance, UploadUserFile, UploadProps, UploadRawFile, UploadInstance } from 'element-plus'

  import {
    importApplicationApi,
    templateAssetsApi,
    importAssetsApi,
    templateNetworkPartitionApi,
    importNetworkPartitionApi,
    templateApplicationApi,
  } from '~/src/api-ecs/assets'

  import { importSiteAssetsApi, templateSiteAssetsApi } from '~/src/api-ecs/site'

  import { downloadFile } from '~/src/utils/download'

  import { proxyNet } from '@/config/index'

  import { genFileId } from 'element-plus'
  import LogicFlow from '@logicflow/core'
  import { number } from 'echarts'

  const props = defineProps<{
    showUpload: boolean
    mode: string
    siteId?: number
  }>()

  const upload = ref<UploadInstance>()

  const formRef = ref<FormInstance>() // 表单实例

  const visible = ref(false) // 弹框显隐

  const remark = ref<string>('site')

  const isLoading = ref(false) // 按钮动画

  // 表单数据
  const formData = reactive<{
    file: UploadUserFile[]
    type?: number
  }>({
    file: [],
    type: undefined,
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

  const urlObj = {
    site: templateSiteAssetsApi,
    know: templateAssetsApi,
    net: templateNetworkPartitionApi,
    app: templateApplicationApi,
  }

  const exportObj = {
    site: importSiteAssetsApi,
    know: importAssetsApi,
    net: importNetworkPartitionApi,
    app: importApplicationApi,
  }

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

  // 导入模板
  const templateAssets = async () => {
    try {
      isLoading.value = true
      // @ts-ignore
      const url = urlObj[remark.value]
      const res = await url()
      downloadFile(res, '导入模板')
      ElMessage({ message: '下载成功', type: 'success' })
    } finally {
      isLoading.value = false
    }
  }

  // 提交表单
  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        isLoading.value = true
        try {
          // @ts-ignore
          const url = exportObj[remark.value]
          const { msg, data } = await url({ file: formData.file[0].raw, type: formData.type, siteId: props.siteId })
          ElMessage({ message: msg, type: 'success' })
          isLoading.value = false
          if (remark.value == 'know') {
            const url = process.env.NODE_ENV == 'development' ? proxyNet : window.location.hostname
            const URL = `https://${url}/download/`
            window.open(URL + data)
          }
          emit('on-reflash')
          handleClose()
        } catch (error: any) {
          isLoading.value = false
        }
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  onMounted(() => {
    visible.value = props.showUpload
    remark.value = props.mode
    if (remark.value == 'know') {
      formData.type = 0
    }
  })
</script>

<script lang="ts">
  export default {
    name: 'ImportAssets',
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
            accept=".xlsx"
            :auto-upload="false"
            class="upload-demo"
            drag
            :limit="1"
            :multiple="false"
            :on-exceed="handleExceed"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或
              <em>选择文件</em>
              <br />
              只能上传xlsx格式文件
            </div>
          </el-upload>
        </el-form-item>
        <el-form-item v-if="remark == 'know'" label="资产属性">
          <el-radio-group v-model="formData.type" class="ml-4">
            <el-radio :label="0">内网资产</el-radio>
            <el-radio :label="1">外网资产</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="&nbsp">
          <div class="el-upload__tip">
            <el-button :loading="isLoading" @click="templateAssets">Excel导入模板下载</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button :loading="isLoading" type="primary" @click="submitForm(formRef)">开始上传</el-button>
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
