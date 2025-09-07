<script lang="ts">
  export default {
    name: 'UploadDate', // 日历导入
  }
</script>

<script setup lang="ts">
  import { useVModel } from '@vueuse/core'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import { UploadUserFile, UploadProps, UploadRawFile, UploadInstance, FormInstance } from 'element-plus'
  import { ResponseData } from '~/src/types'
  import { downloadFile } from '~/src/utils/download'
  // import download from '@/assets/alert_images/download.svg'
  import { requireRules } from '~/src/utils/rules'
  import AesEncryptCBC from '~/src/utils/crypto'

  const props = defineProps<{
    visible: boolean
    title: string
    importAssetsFnc: any
  }>()

  const emit = defineEmits<{
    (e: 'update:visible', visible: boolean): void
    (e: 'reflash'): void
  }>()

  const dialogVisible = useVModel(props, 'visible', emit)

  const yearOption: number[] = []
  for (let index = 2023; index < 2100; index++) {
    yearOption.push(index)
  }

  const uploadIcon = require('@/assets/konwledge/upload.svg')

  const upload = ref<UploadInstance>()

  const $baseMessage: any = inject('$baseMessage')

  const deleteIcon = require('@/assets/konwledge/shanchu.svg')
  const docIcon = require('@/assets/alert_images/Xlsx.svg')

  // 表单数据
  const formData = reactive<{
    file: UploadUserFile[]
    annual: number
  }>({
    file: [],
    annual: 2024,
  })

  const loading = ref(false)
  const formRef = ref<FormInstance>()
  const rules = reactive({
    file: requireRules,
    annual: requireRules,
  })

  const handleExceed: UploadProps['onExceed'] = (files, uploadFiles) => {
    if (files.length > 1) return $baseMessage('暂不支持多个文件同时上传', 'info', 'vab-hey-message-success')
    formData.file[0].raw = files[0] as UploadRawFile
    formData.file[0].name = files[0].name
  }
  const handleUploadSuccess = (response: any) => {
    loading.value = false
  }
  const hanldeFileDelete = async () => {
    formData.file = []
  }

  const submitForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.validate(async (valid) => {
      if (valid) {
        try {
          loading.value = true
          const password = ref('')
          ElMessageBox({
            title: '提示',
            showCancelButton: true,
            customClass: 'need-password-message-box',
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            customStyle: {
              maxWidth: '500px',
            },
            message: () =>
              h('div', null, [
                h(
                  'div',
                  { style: 'margin: 15px 0 5px 0;color:#55585b' },
                  '请输入敏感操作密码:(通过验证后方可进行操作)'
                ),
                h(ElInput, {
                  type: 'password',
                  modelValue: password.value,
                  placeholder: '请输入敏感操作密码',
                  showPassword: true,
                  style: 'margin-block: 10px',
                  'onUpdate:modelValue': (val: string) => {
                    password.value = val
                  },
                }),
              ]),

            beforeClose: async (action, instance, done) => {
              if (action === 'confirm') {
                instance.confirmButtonLoading = true
                try {
                  const { msg } = await props.importAssetsFnc(
                    { annual: formData.annual, password: AesEncryptCBC(password.value) },
                    { file: formData.file[0].raw }
                  )
                  $baseMessage(msg, 'success', 'vab-hey-message-success')
                  emit('reflash')
                  emit('update:visible', false)
                  done()
                } catch (error) {
                  instance.confirmButtonLoading = false
                }
              } else {
                done()
              }
            },
          }).catch(() => {})
        } finally {
          loading.value = false
        }
      } else {
        console.log('error submit!')
        return false
      }
    })
  }

  const handleClose = () => {
    formData.file = []
    loading.value = false
  }
</script>

<template>
  <vab-dialog
    v-model="dialogVisible"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    destroy-on-close
    :loading="loading"
    :title="title"
    width="850px"
    @close="handleClose"
  >
    <!-- :on-error="() => (loading = false)"
      :on-exceed="handleExceed"
      
      :show-file-list="false"
      :on-preview="handlePreview"
      v-model:file-list="formData.file" -->
    <el-form ref="formRef" label-position="top" :model="formData" :rules="rules">
      <el-form-item label="上传文件" prop="file">
        <el-upload
          ref="upload"
          v-model:file-list="formData.file"
          accept=".text,.txt"
          action="/v3/ecsPlatform/infoWhite/importInfoWhite"
          class="upload-demo"
          drag
          :limit="1"
          :multiple="false"
          :on-error="() => (loading = false)"
          :on-exceed="handleExceed"
          :on-success="handleUploadSuccess"
          :show-file-list="false"
        >
          <el-image :src="uploadIcon" style="width: 48px; height: 44px; margin-bottom: 15px" />
          <div class="el-upload__text">
            <span style="color: #4a4759; line-height: 20px">将文件拖到此处，或</span>
            <em style="line-height: 20px">点击上传</em>
            <br />
            <div style="font-size: 13px; color: #bbbdbf; line-height: 14px; margin-top: 9px">
              目前仅text,txt支上传文件类型
            </div>
          </div>
        </el-upload>
        <div v-if="formData.file.length > 0" style="height: 60px; margin-top: 20px; width: 100%">
          <div v-for="(webfile, index) in formData.file" :key="index" class="file">
            <div class="file-left">
              <el-image :src="docIcon" style="width: 24px; height: 28px" />
              <span style="margin-left: 8px">{{ webfile.name }}</span>
            </div>
            <div class="del" @click="() => hanldeFileDelete()">
              <el-image :src="deleteIcon" style="width: 14px; height: 14px; cursor: pointer" />
            </div>
          </div>
        </div>
      </el-form-item>
      <el-form-item label="选择日历年份" prop="annual">
        <el-select v-model="formData.annual" style="width: 100%">
          <el-option v-for="(item, index) in yearOption" :key="index" :label="item" :value="item" />
        </el-select>
      </el-form-item>
    </el-form>
    <div class="dialog-footer" style="margin-top: 100px; text-align: right">
      <!--  -->
      <el-button
        color="#6954F0"
        :disabled="loading || !formData.file.length"
        type="primary"
        @click="submitForm(formRef)"
      >
        确定
      </el-button>
      <el-button @click="dialogVisible = false">取消</el-button>
    </div>
  </vab-dialog>
</template>

<style scoped lang="scss">
  .upload-demo {
    width: 100%;
    // margin-top: 12px;
    :deep() {
      .el-upload-dragger {
        width: 100%;
      }
    }
  }
  :deep() {
    .el-upload-dragger {
      background: #fafaff;
      height: 242px;
      display: flex;
      border-radius: 8px;
      align-items: center;
      justify-content: center;
      flex-direction: column;
    }
    .el-dialog__footer:empty {
      display: none;
    }
  }
  .file {
    width: calc(100% + 2px);
    margin-left: -1px;
    height: 60px;
    // margin-top: 20px;
    background: #ffffff;
    border-radius: 8px;
    border: 1px solid #eae7f5;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 19px;

    .file-left {
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .del {
      width: 24px;
      height: 24px;
      background: #f4f3fc;
      border-radius: 4px;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }
</style>
