<script lang="ts">
  export default {
    name: 'BackupRestore',
  }
</script>

<script setup lang="ts">
  import {
    getAllBksApi,
    saveManualBkApi,
    setIntervalDaysApi,
    getIntervalDaysApi,
    deleteBkApi,
    rollBkApi,
  } from '@/api-ecs/system'
  import { BackupFileType } from '@/types'
  import { useEcsDialogService } from '@/components/ecs-dialog'
  import SensitiveOperation from '@/components/sensitive-operation.vue'
  import { formatTime } from '@/utils/time'
  import AesEncryptCBC from '~/src/utils/crypto'
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const pageType = ref<'backup' | 'restore'>('backup')
  const backupType = ref<'manual' | 'automation'>('manual')
  const manualFiles = ref<BackupFileType[]>([])
  const automationFiles = ref<BackupFileType[]>([])
  const checkedBackupFileName = ref('')
  const loading = ref(false)
  const backupData = ref<string | number>('')
  let rollIndex = -1
  const handleDelete = (fileName: string, index: number) => {
    rollIndex = index
    $baseConfirm('确认删除当前备份？', null, async () => {
      showDialog(3, { fileNames: [fileName] })
    })
  }
  const handleBackupFileCheckedChange = (checked: boolean, fileName: string) => {
    checkedBackupFileName.value = checked ? fileName : ''
  }

  const handleGetAllBks = async () => {
    loading.value = true
    const { data: allFiles = [] } = await getAllBksApi()
    if (allFiles.length === 0) {
      manualFiles.value = []
      automationFiles.value = []
      loading.value = false
      return
    }
    const _manualFiles = []
    const _automationFiles = []
    for (const file of allFiles) {
      if (file.type === 'auto') {
        _automationFiles.push(file)
      } else {
        _manualFiles.push(file)
      }
    }
    manualFiles.value = _manualFiles.sort((curr, next) => next.createTime - curr.createTime)
    automationFiles.value = _automationFiles.sort((curr, next) => next.createTime - curr.createTime)
    loading.value = false
  }

  const showDialog = (sensitiveAciton: number, params = {}) => {
    /**
     * @param {number} sensitiveAciton 敏感操作
     * @param  sensitiveAciton - 0  手动保存
     * @param  sensitiveAciton - 1  自动保存间隔
     * @param  sensitiveAciton - 2  数据还原
     * @param  sensitiveAciton - 3  删除备份文件
     */
    const actionType = [saveManualBkApi, setIntervalDaysApi, rollBkApi, deleteBkApi]
    const _params: { [key: number]: object } = {
      0: {
        remark: backupData.value,
      },
      1: {
        intervalDays: backupData.value,
      },
      2: {
        fileName: checkedBackupFileName.value,
      },
    }
    if ([0, 1, 2].includes(sensitiveAciton)) {
      params = _params[sensitiveAciton]
    }
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
        return new Promise((resolve, reject) => {
          actionType[sensitiveAciton]({ ...params, password: AesEncryptCBC(password.value) } as any)
            .then(({ msg }) => {
              switch (sensitiveAciton) {
                case 0:
                  backupData.value = ''
                  break
                case 2:
                  checkedBackupFileName.value = ''
                  break
                case 3: {
                  const currData = backupType.value === 'automation' ? automationFiles.value : manualFiles.value
                  currData.splice(rollIndex, 1)
                  rollIndex = -1
                  break
                }
                default:
                  break
              }
              $baseMessage(msg, 'success', 'vab-hey-message-success')
              resolve(true)
            })
            .catch(() => {
              reject(false)
            })
        })
      },
    })
  }
  const handleGetIntervalDays = async () => {
    const { data } = await getIntervalDaysApi()
    backupData.value = data
  }
  watch(backupType, () => {
    if (backupType.value === 'automation') {
      backupData.value = 0
      handleGetIntervalDays()
    } else {
      backupData.value = ''
    }
  })
  watch(pageType, () => {
    if (pageType.value === 'restore') {
      backupType.value = 'automation'
      handleGetAllBks()
    } else {
      backupType.value = 'manual'
    }
  })
</script>

<template>
  <div v-loading="loading" class="backup-restore-container">
    <el-button-group class="header">
      <el-button
        :auto-insert-space="false"
        :type="pageType === 'backup' ? 'primary' : 'default'"
        @click="pageType = 'backup'"
      >
        备份
      </el-button>
      <el-button
        :auto-insert-space="false"
        :type="pageType === 'restore' ? 'primary' : 'default'"
        @click="pageType = 'restore'"
      >
        还原
      </el-button>
    </el-button-group>
    <div v-if="pageType === 'backup'" class="backup-content">
      <el-form label-position="right" label-width="110px" @submit.self.prevent>
        <el-form-item label="备份类型:">
          <el-radio-group v-model="backupType">
            <el-radio label="manual">手动</el-radio>
            <el-radio label="automation">自动</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="backupType === 'automation'" label="备份间隔时间:">
          <el-input-number
            v-model="backupData"
            class="automation-input"
            :max="30"
            :min="1"
            placeholder="请输入备份间隔时间"
            style="width: 100px"
          />
          <span class="day">天</span>
          <span class="tips">温馨提示：默认保存最近30次备份文件，历史会被删除</span>
        </el-form-item>
        <el-form-item v-else label="备注:">
          <el-input v-model="backupData" placeholder="请输入备注" style="width: 300px" />
        </el-form-item>
        <el-form-item label="&nbsp;">
          <el-button :auto-insert-space="false" type="primary" @click="showDialog(backupType === 'manual' ? 0 : 1)">
            确认
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div v-else class="restore-content">
      <div class="header">
        选择文件：
        <span class="backup-type" :class="{ checked: backupType === 'automation' }" @click="backupType = 'automation'">
          自动备份文件
        </span>
        <el-divider direction="vertical" />
        <span class="backup-type" :class="{ checked: backupType === 'manual' }" @click="backupType = 'manual'">
          手动备份文件
        </span>
        <el-button
          v-if="checkedBackupFileName"
          :auto-insert-space="false"
          style="float: right"
          type="primary"
          @click="showDialog(2)"
        >
          还原
        </el-button>
      </div>
      <div class="backup-histories">
        <el-row :gutter="9">
          <el-col
            v-for="(file, index) of backupType === 'automation' ? automationFiles : manualFiles"
            :key="file.createTime"
            :span="6"
          >
            <div class="backup-history" :class="{ checked: file.fileName === checkedBackupFileName }">
              <div class="backup-info">
                <el-checkbox
                  v-model="checkedBackupFileName"
                  :label="file.fileName"
                  :true-label="file.fileName"
                  @change="handleBackupFileCheckedChange($event, file.fileName)"
                />
                <span class="date">{{ formatTime(file.createTime) }}</span>
                <span class="remark">{{ file.type === 'auto' ? '路径：' : '备注：' }}{{ file.remark || '空' }}</span>
                <el-icon class="backup-delete" @click="handleDelete(file.fileName, index)"><Delete /></el-icon>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .backup-restore-container {
    padding-bottom: 0 !important;
    & > .header {
      .el-button {
        width: 80px;
      }
    }
    .backup-content {
      margin-top: 16px;
      height: calc(100vh - 80px);
      background: #f9f8ff;
      border-radius: 8px;
      padding: 20px;
      padding-bottom: 0;
      .el-form {
        width: 70%;
        margin: 30px auto 0;
        .automation-input {
          :deep(.el-input__inner) {
            text-align: center;
          }
        }
        .day {
          font-weight: 400;
          font-size: 14px;
          color: #b7b5bf;
          margin-left: 5px;
        }
        .tips {
          font-weight: 400;
          font-size: 14px;
          color: #4a4759;
          margin-left: 30px;
        }
      }
      .el-radio {
        margin-right: 20px;
        :deep(.el-radio__label) {
          padding-left: 4px;
        }
      }
    }
    .restore-content {
      .header {
        line-height: 32px;
        margin-top: 10px;
        font-weight: 400;
        font-size: 14px;
        color: #908e9e;
        .backup-type {
          font-weight: 500;
          font-size: 15px;
          color: #7d7990;
          cursor: pointer;
          &.checked {
            color: var(--el-color-primary);
          }
        }
      }
      .backup-histories {
        margin-top: 10px;
        height: calc(100vh - 116px);
        background: #f9f8ff;
        border-radius: 8px;
        padding: 20px;
        padding-bottom: 0;
        margin-inline: 0 !important;
        .backup-history {
          width: 100%;
          height: 106px;
          border-radius: 6px;
          border: 1px solid #e5e4ef;
          padding: 10px 16px;
          position: relative;
          background-color: #fff;
          margin-bottom: 8px;
          cursor: pointer;
          cursor: pointer;
          .backup-info {
            display: flex;
            flex-direction: column;
            .el-checkbox {
              width: calc(100% - 20px);
              :deep(.el-checkbox__label) {
                flex: 1;
                overflow: hidden;
                word-break: keep-all;
                text-overflow: ellipsis;
                font-weight: 500;
                font-size: 15px;
                color: #494758;
              }
            }
            .date,
            .remark {
              font-weight: 400;
              font-size: 14px;
              color: #7d7990;
              text-indent: 24px;
              line-height: 24px;
              overflow: hidden;
              word-break: keep-all;
              text-overflow: ellipsis;
              white-space: nowrap;
            }
          }
          .backup-delete {
            position: absolute;
            top: 21px;
            right: 15px;
          }
          &.checked {
            border-color: var(--el-color-primary);
          }
        }
      }
    }
  }
</style>
