<script lang="ts">
  export default {
    name: 'KnowledgeIndex', // 知识
  }
</script>

<script setup lang="ts">
  import VabDialog from '@/plugins/VabDialog/index.vue'

  import { Plus } from '@element-plus/icons-vue'

  import { FormInstance, UploadUserFile, UploadProps, UploadRawFile } from 'element-plus'

  import {
    knowledgeDeleteApi,
    knowledgeGetAllApi,
    knowledgeUploadApi,
    knowledgeFileDeleteApi,
    knowledgeFileSaveApi,
  } from '@/api-ecs/knowledge'

  import numberFormatter from '~/src/utils/number'
  import { formatNstime } from '~/src/utils/time'

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  const deleteIcon = require('@/assets/konwledge/shanchu.svg')
  const uploadIcon = require('@/assets/konwledge/upload.svg')
  const docIcon = require('@/assets/konwledge/DOC.svg')
  const mdIcon = require('@/assets/konwledge/MD.svg')
  const txtIcon = require('@/assets/konwledge/TXT.svg')
  const loading = ref(false)
  const num = ref(0)

  const visible = ref(false)

  const formRef = ref<FormInstance>() // 表单实例

  // 表单数据
  const formData = reactive<{
    file: UploadUserFile[]
    webFiles: {
      name: string
      value: string
    }[]
  }>({
    file: [],
    webFiles: [],
  })

  const list = ref<{ id: number; suffix: string; name: string; numberOfByte: number; createTime: number }[]>([])

  const showDialogDel = ref(false)

  const uploadFiles = ref<string[]>([])

  const curID = ref<number>()

  onMounted(() => {
    getPage()
  })

  const handleShowDialog = () => {
    visible.value = true
  }

  const getPage = async () => {
    const { data } = await knowledgeGetAllApi()
    list.value = data.gather
    num.value = data.sumSize
  }

  const handleExceed: UploadProps['onExceed'] = async (files, uploadFiles) => {
    formData.file[0].raw = files[0] as UploadRawFile
    formData.file[0].name = files[0].name
    try {
      const { data }: any = await hanldeFileUpload()
      formData.webFiles.push(data)
    } finally {
      loading.value = false
    }
  }
  const hanldeFileUpload = () => {
    loading.value = true
    return knowledgeUploadApi({ file: formData.file[0].raw })
  }
  const handleUploadSuccess = (response: any) => {
    formData.webFiles.push(response.data)
    loading.value = false
  }
  const hanldeFileDelete = async (path: string) => {
    await knowledgeFileDeleteApi(path)
    formData.webFiles = formData.webFiles.filter((i) => i.value !== path)
  }
  const submitForm = async () => {
    const files = formData.webFiles.map((i) => i.value)
    loading.value = true
    try {
      const { msg } = await knowledgeFileSaveApi([...files])
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      getPage()
      visible.value = false
    } finally {
      loading.value = false
    }
  }

  const handleDel = async (id: number) => {
    $baseConfirm('你确定要删除该项吗？', null, async () => {
      const { msg } = await knowledgeDeleteApi(JSON.stringify({ ids: [id] }))
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      getPage()
    })
  }

  const formate = (name: string) => {
    if (name) {
      const arr = name.split('.')
      let res = undefined
      switch (arr[arr.length - 1]) {
        case 'txt':
          res = txtIcon
          break
        case 'md':
          res = mdIcon
          break
        case 'docx':
          res = docIcon
          break
        case 'doc':
          res = docIcon
          break

        default:
          break
      }
      return res
    } else {
      return undefined
    }
  }

  const handleClose = () => {
    formData.file = []
    formData.webFiles = []
  }
</script>

<template>
  <div class="knowledge-container">
    <div class="title">
      <div class="left">
        <span class="word">知识库</span>
        <span class="count">共计：{{ numberFormatter.format(num) }}字数</span>
      </div>
      <div class="right">
        <el-button color="#6954F0" :icon="Plus" @click="handleShowDialog">上传知识</el-button>
      </div>
    </div>
    <div class="content">
      <el-row :gutter="15" style="overflow: hidden">
        <el-col
          v-for="item in list"
          :key="item.id"
          :span="6"
          @mouseenter="curID = item.id"
          @mouseleave="curID = undefined"
        >
          <div class="item">
            <div class="left">
              <el-image :src="formate(item.name)" style="width: 38px; height: 46px" />
            </div>
            <div class="middle">
              <div class="title">
                <el-tooltip class="box-item" :content="item.name" effect="dark" placement="top-start">
                  <span>{{ item.name }}</span>
                </el-tooltip>
              </div>
              <div style="margin: 6px 0 2px 0">字数：{{ item.numberOfByte }}</div>
              <div class="updateTime">更新时间：{{ formatNstime(item.createTime, false) }}</div>
            </div>
            <div v-if="item.id == curID" class="right" @click="handleDel(item.id)">
              <el-image :src="deleteIcon" style="width: 14px; height: 14px; cursor: pointer" />
            </div>
          </div>
        </el-col>
      </el-row>

      <el-empty
        v-if="list.length == 0"
        description="暂无数据"
        style="height: 100%; cursor: pointer"
        @click="handleShowDialog"
      />
    </div>
    <vab-dialog
      v-model="visible"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      destroy-on-close
      :loading="loading"
      title="数据中心管理"
      width="850px"
      @close="handleClose"
    >
      <el-upload
        ref="upload"
        v-model:file-list="formData.file"
        accept=".tar,.doc,.txt,.docx,.md"
        action="/v3/ecsPlatform/knowledge/upload"
        auto-upload
        class="upload-demo"
        drag
        :http-request="hanldeFileUpload"
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
            目前仅支持 Doc、Docx、Txt、Markdown文件类型
          </div>
        </div>
      </el-upload>

      <div v-for="(webfile, index) in formData.webFiles" :key="index" class="file">
        <div class="file-left">
          <el-image :src="formate(webfile.value)" style="width: 24px; height: 28px" />
          <span style="margin-left: 8px">{{ webfile.name }}</span>
        </div>
        <div class="del" @click="() => hanldeFileDelete(webfile.value)">
          <el-image :src="deleteIcon" style="width: 14px; height: 14px; cursor: pointer" />
        </div>
      </div>
      <div class="dialog-footer" style="margin-top: 160px; text-align: right">
        <el-button color="#6954F0" :disabled="loading || !formData.webFiles.length" type="primary" @click="submitForm">
          确定
        </el-button>
        <el-button @click="visible = false">取消</el-button>
      </div>
    </vab-dialog>
  </div>
</template>

<style scoped lang="scss">
  .knowledge-container {
    // padding-top: 20px !important;
    padding-bottom: 20px;

    .updateTime {
      overflow: hidden;
      height: 24px;
      text-overflow: ellipsis;
      white-space: nowrap;
      word-break: break-all;
      word-wrap: break-word;
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

    .upload-demo {
      height: 242px;
    }

    .file {
      width: calc(100% + 2px);
      margin-left: -1px;
      height: 60px;
      margin-top: 20px;
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

    .title {
      width: 100%;
      display: flex;
      align-items: center;
      justify-content: space-between;
    }

    .left {
      display: flex;
      align-items: center;
    }

    .word {
      font-weight: 500;
      font-size: 20px;
      color: #303133;
      line-height: 30px;
    }

    .count {
      margin-left: 10px;
      padding: 0 11px;
      height: 30px;
      background: #f4f3fa;
      border-radius: 4px;
      font-weight: 400;
      font-size: 14px;
      color: #605b7d;
      line-height: 20px;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .content {
      margin-top: 20px;
      height: calc(100vh - 115px);
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
      .item {
        // height: 134px;
        background: #ffffff;
        border-radius: 12px;
        border: 1px solid #eae7f5;
        margin-bottom: 15px;
        padding: 20px 20px;
        display: flex;
        flex-direction: row;
        align-items: flex-start;

        .middle {
          margin: 0 14px;
          width: calc(100% - 80px);
          font-size: 14px;
          color: #9d9baa;
          line-height: 24px;
          .title {
            font-weight: 500;
            font-size: 16px;
            color: #342e58;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            word-break: break-all;
            word-wrap: break-word;
            width: 100%;
            :deep() {
              .el-only-child__content {
                width: calc(100%);
                white-space: nowrap;
                overflow: hidden;
                text-overflow: ellipsis;
                word-break: break-all;
                word-wrap: break-word;
              }
            }
          }
        }
        .right {
          width: 14px;
          margin-top: 5px;
        }
      }
    }
  }
</style>
