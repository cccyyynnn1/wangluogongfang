<script lang="ts">
  export default {
    name: 'AttackCharacterization',
  }
</script>

<script setup lang="ts">
  import { useEcsDialogService } from '@/components/ecs-dialog'
  import AttackHighlightRule from './attack-highlight-rule.vue'
  import { UploadUserFile, UploadProps, UploadInstance, UploadRawFile, genFileId, FormInstance } from 'element-plus'
  import SensitiveOperation from '@/components/sensitive-operation.vue'
  import {
    getAttackHighlightConfigApi,
    getAttackHighlightWhiteApi,
    exportHighLightApi,
    exportHighLightTemplateApi,
    importHighLightApi,
    deleteAttackHighLightApi,
  } from '@/api-ecs/attack-characterization'
  import AesEncryptCBC from '~/src/utils/crypto'
  import { AttackCharacterizationType } from '@/types/index'
  import { downloadFile } from '~/src/utils/download'
  type ScopeType = 'requestHeader' | 'requestPayload' | 'responseHeader' | 'responsePayload' | 'all'

  const uploadIcon = require('@/assets/konwledge/upload.svg')
  const uploadVisible = ref(false)
  const uploadLoading = ref(false)
  const $baseMessage: any = inject('$baseMessage')
  const uplodFormRef = ref<FormInstance>()
  const uplodRef = ref<UploadInstance>()
  // 表单数据
  const uploadFile = reactive<{
    file: UploadUserFile[]
  }>({
    file: [],
  })
  const rules = {
    file: [
      {
        required: true,
        message: '请选择上传文件',
        trigger: 'change',
      },
    ],
  }

  const attackConfigQuery = reactive({
    pageNum: 1,
    pageSize: 50,
    total: 0,
  })
  const attackWhiteQuery = reactive({ pageNum: 1, pageSize: 50, total: 0 })
  const attackConfigList = ref<AttackCharacterizationType[]>([])
  const attackWhiteList = ref<AttackCharacterizationType[]>([])
  const attackCharacterizationChecked = reactive<{ configIds: number[]; whiteIds: number[] }>({
    configIds: [],
    whiteIds: [],
  })
  const handleExceed: UploadProps['onExceed'] = (files) => {
    uplodRef.value?.clearFiles()
    const file = files[0] as UploadRawFile
    file.uid = genFileId()
    uplodRef.value!.handleStart(file)
  }
  // 导出攻击特征和特征白名单
  const handleExportAttackCharacterization = async () => {
    const data = await exportHighLightApi()
    downloadFile(data, '')
    $baseMessage('下载成功', 'success', 'vab-hey-message-success')
  }
  // 模板导入
  const templateUpload = async () => {
    uplodFormRef.value!.validate(async () => {
      const { msg } = await importHighLightApi({ file: uploadFile.file[0].raw! })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      uploadVisible.value = false
      handleGetList()
    })
  }
  // 模板导出
  const templateExport = async () => {
    try {
      uploadLoading.value = true
      const res = await exportHighLightTemplateApi()
      downloadFile(res, '攻击特征模版')
      $baseMessage('下载成功', 'success', 'vab-hey-message-success')
    } finally {
      uploadLoading.value = false
    }
  }

  const handleAttackDelete = (data: { configIds: number[]; whiteIds: number[] }, isAll = false) => {
    const password = ref('')
    const { destroy } = useEcsDialogService({
      title: '提示',
      content: () =>
        h(SensitiveOperation, {
          modelValue: password.value,
          'onUpdate:modelValue': (val: string) => {
            password.value = val
          },
        }),
      async confirm() {
        if (!password.value.length) return Promise.reject()
        const params = {
          ...data,
          deleteAll: isAll,
          password: AesEncryptCBC(password.value),
        }
        return new Promise((resolve, reject) => {
          deleteAttackHighLightApi(params)
            .then(({ msg }) => {
              $baseMessage(msg, 'success', 'vab-hey-message-success')
              handleGetList()
              resolve(true)
            })
            .catch(() => {
              resolve(true)
            })
          resolve(params)
        })
      },
    })
  }

  const handleAttackSelectionChange = (
    rows: AttackCharacterizationType[],
    type: 'highLightConfig' | 'highLightWhite'
  ) => {
    const isd = rows.map(({ id }) => id as number)
    if (type === 'highLightConfig') {
      attackCharacterizationChecked.configIds = isd
    } else {
      attackCharacterizationChecked.whiteIds = isd
    }
  }
  const handleEditRule = (
    row: Partial<AttackCharacterizationType> & { type?: 'highLightConfig' | 'highLightWhite' }
  ) => {
    const { destroy } = useEcsDialogService({
      title: row.id ? '编辑' : '添加',
      showFooter: false,
      content: () => h(AttackHighlightRule, { attackHighlightData: row }),
      confirm() {
        handleGetList()
      },
    })
  }
  const handleGetAttackConfig = async () => {
    const { pageNum, pageSize } = attackConfigQuery
    const { data } = await getAttackHighlightConfigApi({ pageNum, pageSize })
    attackConfigList.value = data.records || []
    attackConfigQuery.total = data.total || 0
  }
  const handleGetAttackWhite = async () => {
    const { pageNum, pageSize } = attackWhiteQuery
    const { data } = await getAttackHighlightWhiteApi({ pageNum, pageSize })
    attackWhiteList.value = data.records || []
    attackWhiteQuery.total = data.total || 0
  }
  const handleGetList = () => {
    handleGetAttackWhite()
    handleGetAttackConfig()
  }

  const formatScope = (row: AttackCharacterizationType) => {
    const formatObj: { [key in ScopeType]: string } = {
      requestHeader: '请求头',
      requestPayload: '请求体',
      responseHeader: '响应头',
      responsePayload: '响应体',
      all: '全部',
    }
    return formatObj[row.scope]
  }
  watchEffect(() => {
    if (!uploadVisible.value) {
      uploadFile.file = []
    }
  })
  onMounted(() => {
    handleGetList()
  })
</script>

<template>
  <div class="attack-characterization-container">
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>高亮特征配置</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12">
        <el-button type="primary" @click="handleEditRule({})">添加</el-button>
        <el-button type="primary" @click="uploadVisible = true">导入</el-button>
        <el-button type="primary" @click="handleExportAttackCharacterization">导出</el-button>
        <el-dropdown style="margin-left: 10px">
          <span class="el-dropdown-link">
            <el-button type="danger">
              批量删除
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleAttackDelete(attackCharacterizationChecked)">删除选中</el-dropdown-item>
              <el-dropdown-item @click="handleAttackDelete(attackCharacterizationChecked, true)">
                删除所有
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </vab-query-form-right-panel>
    </vab-query-form>
    <div class="attack-characterization-config">
      <div class="attack-characterization-library">
        <div class="config-header">
          <vab-icon icon="file-settings-line" />
          高亮特征库
        </div>
        <el-table
          :data="attackConfigList"
          size="small"
          @selection-change="(rows) => handleAttackSelectionChange(rows, 'highLightConfig')"
        >
          <el-table-column type="selection" width="50" />
          <el-table-column label="序号" type="index" width="80" />
          <el-table-column label="高亮字段" prop="content" />
          <el-table-column :formatter="formatScope" label="作用域" prop="scope" />
          <el-table-column label="备注" prop="remark" />
          <el-table-column label="操作" width="150">
            <template #default="{ row }">
              <el-button plain size="small" @click="handleEditRule({ ...row, type: 'highLightConfig' })">
                编辑
              </el-button>
              <el-button plain size="small" @click="handleAttackDelete({ configIds: [row.id], whiteIds: [] })">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination
          v-model:current-page="attackConfigQuery.pageNum"
          v-model:page-size="attackConfigQuery.pageSize"
          layout="total, prev, pager, next"
          :page-sizes="[50, 100, 150]"
          :total="attackConfigQuery.total"
          @current-change="handleGetAttackConfig"
          @size-change="handleGetAttackConfig"
        />
      </div>
      <div class="attack-characterization-white">
        <div class="config-header">
          <vab-icon icon="file-lock-line" />
          特征库白名单
        </div>
        <el-table
          :data="attackWhiteList"
          size="small"
          @selection-change="(rows) => handleAttackSelectionChange(rows, 'highLightWhite')"
        >
          <el-table-column type="selection" width="50" />
          <el-table-column label="序号" type="index" width="80" />
          <el-table-column label="高亮字段" prop="content" />
          <el-table-column :formatter="formatScope" label="作用域" prop="scope" />
          <el-table-column label="备注" prop="remark" />
          <el-table-column label="操作" width="150">
            <template #default="{ row }">
              <el-button plain size="small" @click="handleEditRule({ ...row, type: 'highLightWhite' })">编辑</el-button>
              <el-button plain size="small" @click="handleAttackDelete({ configIds: [], whiteIds: [row.id] })">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination
          v-model:current-page="attackWhiteQuery.pageNum"
          v-model:page-size="attackWhiteQuery.pageSize"
          layout="total, prev, pager, next"
          :page-sizes="[50, 100, 150]"
          :total="attackWhiteQuery.total"
          @current-change="handleGetAttackWhite"
          @size-change="handleGetAttackWhite"
        />
      </div>
    </div>
    <el-dialog v-if="uploadVisible" v-model="uploadVisible" destroy-on-close title="导入" width="850px">
      <el-form
        ref="uplodFormRef"
        class="examplesDetailUpload-form"
        label-position="top"
        label-width="100px"
        :model="uploadFile"
        :rules="rules"
      >
        <el-form-item label="上传文件" prop="file">
          <el-upload
            ref="uplodRef"
            v-model:file-list="uploadFile.file"
            accept=".xlsx"
            :auto-upload="false"
            class="upload-demo"
            drag
            :limit="1"
            :on-exceed="handleExceed"
          >
            <el-icon class="el-icon--upload">
              <el-image :src="uploadIcon" style="width: 48px; height: 44px" />
            </el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或
              <em>选择文件</em>
              <br />
              只能上传xlsx格式文件
            </div>
          </el-upload>
        </el-form-item>
        <el-form-item label="&nbsp">
          <div class="el-upload__tip">
            <el-button :loading="uploadLoading" @click="templateExport">Excel导入模版下载</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button :loading="uploadLoading" type="primary" @click="templateUpload">开始上传</el-button>
        <el-button @click="uploadVisible = false">取消</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .attack-characterization-container {
    display: flex;
    flex-direction: column;
    h3 {
      margin-block: 0;
    }
    .attack-characterization-config {
      flex: 1;
      overflow: hidden;
      margin-top: 10px;
      border: 1px solid #ebeef5;
      display: flex;
      & > div {
        flex: 1;
        overflow: hidden;
        .config-header {
          height: 40px;
          font-weight: 500;
          font-size: 15px;
          color: #645f84;
          line-height: 40px;
          padding-left: 20px;
        }
        .el-pagination {
          padding-right: 10px;
          justify-content: center;
        }
        :deep() {
          .el-table.el-table--small {
            th {
              background-color: var(--el-table-tr-bg-color) !important;
            }
            .el-table__cell {
              padding-block: 8px;
              .el-button.is-plain {
                width: 44px;
                border-radius: 4px;
                & + .el-button {
                  margin-left: 8px;
                }
              }
              &.el-table-column--selection > .cell {
                margin-left: 7px;
              }
            }
          }
        }
        &.attack-characterization-library {
          .config-header {
            background: #e4e3ec;
          }
          .el-table {
            height: calc(100% - 92px);
          }
        }
        &.attack-characterization-white {
          border-left: 1px solid #ebeef5;
          .config-header {
            background: #f4f2ff;
          }
          .el-table {
            height: calc(100% - 92px);
          }
        }
      }
    }
    .upload-demo {
      width: 100%;
      :deep() {
        .el-upload-dragger {
          background: #fafaff;
          border-radius: 8px;
          padding-block: 55px;
        }
        .el-icon--upload {
          margin-bottom: 0;
        }
      }
    }
    .el-upload__tip {
      font-size: 20px;
      margin: -10px auto 0;
    }
  }
</style>
