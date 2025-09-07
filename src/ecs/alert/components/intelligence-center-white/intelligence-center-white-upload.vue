<script lang="ts">
  export default {
    name: 'IntelligenceCenterWhiteUpload', // 情报白名单导入
  }
</script>

<script setup lang="ts">
  import { useVModel } from '@vueuse/core'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import { UploadUserFile, UploadProps, UploadRawFile, UploadInstance } from 'element-plus'
  import { ResponseData } from '~/src/types'
  import { downloadFile } from '~/src/utils/download'
  // import download from '@/assets/alert_images/download.svg'

  const props = defineProps<{
    visible: boolean
    title: string
    downloadTemp: any
    importAssetsFnc: any
  }>()

  const emit = defineEmits<{
    (e: 'update:visible', visible: boolean): void
    (e: 'reflash'): void
  }>()

  const dialogVisible = useVModel(props, 'visible', emit)

  const uploadIcon = require('@/assets/konwledge/upload.svg')

  const upload = ref<UploadInstance>()

  const $baseMessage: any = inject('$baseMessage')

  const deleteIcon = require('@/assets/konwledge/shanchu.svg')
  const docIcon = require('@/assets/alert_images/Xlsx.svg')

  // 表单数据
  const formData = reactive<{
    file: UploadUserFile[]
  }>({
    file: [],
  })

  const loading = ref(false)

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
  const submitForm = async () => {
    try {
      loading.value = true
      const { msg } = await props.importAssetsFnc({ file: formData.file[0].raw })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      emit('reflash')
      emit('update:visible', false)
    } finally {
      loading.value = false
    }
  }

  const handleTempDownLoad = async () => {
    const res = await props.downloadTemp()
    downloadFile(res, `${props.title}模版`)
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
    <el-button color="#6954F0" style="float: right; margin-bottom: 15px" type="primary" @click="handleTempDownLoad">
      <!-- Generator: Adobe Illustrator 24.0.0, SVG Export Plug-In . SVG Version: 6.00 Build 0)  -->
      <svg
        id="图层_1"
        style="height: 12px; width: 12px; fill: #ffffff; margin-right: 2px"
        version="1.1"
        viewBox="0 0 14 14"
        x="0px"
        xml:space="preserve"
        xmlns="http://www.w3.org/2000/svg"
        xmlns:xlink="http://www.w3.org/1999/xlink"
        y="0px"
      >
        <title>xiazai-4</title>
        <g id="情报中心">
          <g id="情报中心-自定义情报-导入" transform="translate(-1267, -365)">
            <g id="侧边弹窗" transform="translate(536, 265)">
              <g id="下载模版" transform="translate(718, 90)">
                <g id="编组-14" transform="translate(13, 7)">
                  <g id="编组-15" transform="translate(0, 3)">
                    <g id="xiazai-4" transform="translate(0, 0)">
                      <path
                        id="路径"
                        class="st0"
                        d="M14,11.6c0,0.8-0.7,1.4-1.5,1.4H1.5c-0.8,0-1.5-0.6-1.5-1.4V8.3c0-0.4,0.3-0.7,0.7-0.7
								s0.7,0.3,0.7,0.7v2.5c0,0.5,0.4,0.9,0.9,0.9h9.3c0.5,0,0.9-0.4,0.9-0.9V8.3c0-0.4,0.3-0.7,0.7-0.7c0.4,0,0.7,0.3,0.7,0.7
								V11.6z"
                      />
                      <path
                        id="路径_1_"
                        class="st0"
                        d="M7.2,0L7.2,0c0.5,0,0.8,0.2,0.8,0.7v6.7c0,0.5-0.3,0.7-0.8,0.7l0,0
								c-0.5,0-0.8-0.2-0.8-0.7V0.7C6.4,0.2,6.7,0,7.2,0z"
                      />
                      <path
                        id="路径_2_"
                        class="st0"
                        d="M4.2,5.5c0.1-0.1,0.3-0.2,0.5-0.2c0.2,0,0.4,0.1,0.5,0.2l2.7,2.7c0.3,0.3,0.3,0.8,0,1
								s-0.8,0.3-1,0L4.2,6.5C4.1,6.4,4,6.2,4,6S4.1,5.6,4.2,5.5L4.2,5.5z"
                      />
                      <path
                        id="路径_3_"
                        class="st0"
                        d="M10.6,5.5c0.1,0.1,0.2,0.3,0.2,0.5s-0.1,0.4-0.2,0.5L7.9,9.2c-0.3,0.3-0.8,0.3-1,0
								c-0.3-0.3-0.3-0.8,0-1l2.7-2.7c0.1-0.1,0.3-0.2,0.5-0.2S10.4,5.4,10.6,5.5L10.6,5.5z"
                      />
                    </g>
                  </g>
                </g>
              </g>
            </g>
          </g>
        </g>
      </svg>
      下载模板
    </el-button>
    <!-- :on-error="() => (loading = false)"
      :on-exceed="handleExceed"
      
      :show-file-list="false"
      :on-preview="handlePreview"
      v-model:file-list="formData.file" -->
    <el-upload
      ref="upload"
      v-model:file-list="formData.file"
      accept=".xls,.xlsx"
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
          目前仅支上传xls、xlsx文件类型
        </div>
      </div>
    </el-upload>
    <div style="height: 60px; margin-top: 20px">
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
    <div class="dialog-footer" style="margin-top: 100px; text-align: right">
      <!--  -->
      <el-button color="#6954F0" :disabled="loading || !formData.file.length" type="primary" @click="submitForm">
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
