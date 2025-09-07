<script lang="ts">
  export default {
    name: 'PcapUpload',
  }
</script>

<script setup lang="ts">
  import { uploadLocalPcapsApi, getRemoteFolderApi, saveRemoteFolderApi } from '@/api-ecs/packet-replay'
  import type { FormInstance, InputInstance } from 'element-plus'
  import { RemoteFolderItem } from '@/types/index'
  import { CircleCloseFilled } from '@element-plus/icons-vue'
  import { cloneDeep } from 'lodash'
  const folderIcon = require('@/assets/index_images/folder.svg')
  const gobackIcon = require('@/assets/index_images/goback.svg')
  const pcapIcon = require('@/assets/index_images/pcap.svg')

  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const props = defineProps<{
    modelValue: boolean
    uploadType: 'online' | 'directory'
  }>()
  const emits = defineEmits<{
    (e: 'update:modelValue', visibility: boolean): void
  }>()

  //在线上传
  const uploading = ref(false)
  const uploadVisible = useVModel(props, 'modelValue', emits)
  const fileList = ref<{ name: string; remark: string; raw: File }[]>([])
  const deleteFile = (arr: any[], index: number) => {
    const deleteItem = arr.splice(index, 1)
  }
  //目录挂载所需要的数据
  // 当前选中的文件
  const currentCheckedFolder = ref<RemoteFolderItem>()
  // 目录下的文件
  const remoteFolders = ref<RemoteFolderItem[]>([])
  const remoteFoldersVisible = ref(false)
  // 目录路径
  const remote_folder_ref = ref<InputInstance>()
  //已选的文件，通过checked状态判断是否选中，在目录挂载上传Pcap展示
  const selectedFolders = ref<Array<RemoteFolderItem & { checked?: boolean }>>([])
  //目录挂载展示列表
  const remoteFolderList = ref<RemoteFolderItem[]>([])
  const remoteFolderForm = ref<FormInstance>()
  const remoteFolderPath = ref('/')
  const remotefoldersLoading = ref(true)
  const remoteFolderFormData = reactive({
    readType: 'smb',
    ip: '',
    port: 445,
    account: '',
    password: '',
    requestId: new Date().getTime(),
  })

  // 在线上传
  const handle2Upload = async () => {
    uploading.value = true
    if (props.uploadType === 'online') {
      if (fileList.value.length === 0) {
        uploading.value = false
        return $baseMessage('请选择Pcap文件', 'error', 'vab-hey-message-error')
      }
      const files = fileList.value.reduce(
        (data, file, index) => {
          data.files.push(file.raw)
          data.remarks[index] = file.remark || ''
          return data
        },
        {
          files: [] as File[],
          remarks: {} as any,
        }
      )
      const { msg } = await uploadLocalPcapsApi(files)
      uploadVisible.value = false
      uploading.value = false
      $baseMessage(msg, 'success', 'vab-hey-message-success')
    } else {
      if (remoteFolderList.value.length === 0) {
        uploading.value = false
        return $baseMessage('请选择Pcap文件', 'error', 'vab-hey-message-error')
      }
      handleSubmitRemoteFolders()
    }
  }
  //目录挂载获取当前目录下文件
  const handleGetRemoteFolders = async () => {
    remotefoldersLoading.value = true
    const { data } = await getRemoteFolderApi({ ...remoteFolderFormData, path: remoteFolderPath.value })
    remoteFolders.value = data || []
    remotefoldersLoading.value = false
  }

  /****
   * @param toSync 是否同步数据到目录挂载已选择列表
   * @description 查询当前目录下文件
   */
  const handle2read = (toSync = false) => {
    remoteFolderForm.value?.validate(async (valid) => {
      if (valid) {
        remoteFoldersVisible.value = true
        handleGetRemoteFolders()
      }
    })
    if (toSync) {
      selectedFolders.value = cloneDeep(remoteFolderList.value)
    }
  }
  //返回上一层
  const handleGoBack = (data: RemoteFolderItem) => {
    remoteFolderPath.value = data.filePath
    handle2read()
  }
  //目录挂载下文件点击
  const handleFolderClick = (data: RemoteFolderItem) => {
    if (data.type === 'folder') {
      remoteFolderPath.value = data.filePath
      handle2read()
      currentCheckedFolder.value = undefined
    } else {
      currentCheckedFolder.value = data
    }
  }
  //点击选择按钮，把当前已选中的数据放入到右侧已选择列表
  const handleToSelected = () => {
    const last = selectedFolders.value.find((i) => i.filePath === currentCheckedFolder.value?.filePath)
    if (!currentCheckedFolder.value || last?.filePath === currentCheckedFolder.value.filePath) return
    selectedFolders.value.push(cloneDeep({ ...currentCheckedFolder.value, checked: true }))
  }
  //把目录读取的已选择文件回显到目录挂载页面
  const handleSaveRemoteFolders = () => {
    const _selectedFolders = selectedFolders.value.filter((i) => i.checked)
    remoteFolderList.value = cloneDeep(_selectedFolders)
    remoteFoldersVisible.value = false
  }
  const handleSubmitRemoteFolders = async () => {
    const dadada = {
      ...remoteFolderFormData,
      path: remoteFolderPath.value,
      remoteFileList: remoteFolderList.value.map(({ fileName, filePath, remark }) => ({ fileName, filePath, remark })),
    }
    const { msg } = await saveRemoteFolderApi(dadada)
    uploadVisible.value = false
    uploading.value = false
    $baseMessage(msg, 'success', 'vab-hey-message-success')
  }

  const handleReadTypeChange = (type: string) => {
    const ports: {
      [key: string]: number
    } = {
      smb: 445,
      sftp: 22,
    }
    remoteFolderFormData.port = ports?.[type] || 21
  }

  watch(
    () => remoteFolderFormData,
    () => {
      remoteFolderFormData.requestId = new Date().getTime()
    },
    {
      deep: true,
    }
  )
</script>

<template>
  <div class="remote-folders-box">
    <el-dialog v-model="uploadVisible" :title="uploadType === 'online' ? '在线上传' : '目录挂载'" width="870px">
      <template v-if="uploadType === 'online'">
        <el-upload
          v-model:file-list="fileList"
          accept=".pcap,.cap,.pcapng"
          :auto-upload="false"
          class="upload"
          drag
          :limit="3"
          multiple
          :show-file-list="false"
        >
          <div class="el-upload__text">
            <em>
              <el-icon><CirclePlus /></el-icon>
              上传数据包
            </em>
          </div>
        </el-upload>
        <ul class="upload-file-list">
          <li v-for="(file, index) in fileList" :key="file.name" class="file-item">
            <div class="handle">
              <div class="name">
                <vab-icon icon="file-text-fill" style="color: #ffc13d; margin-inline: 10px 7px" />
                {{ file.name }}
              </div>
              <el-input
                v-model="file.remark"
                clearable
                placeholder="添加备注"
                style="height: 40px; margin-top: -1px; margin-right: -1px"
              />
              <el-icon class="delete-icon" @click="deleteFile(fileList, index)"><Delete /></el-icon>
            </div>
          </li>
        </ul>
      </template>
      <template v-else-if="uploadType === 'directory'">
        <div style="min-height: 450px">
          <el-form
            ref="remoteFolderForm"
            class="directory-upload-form"
            label-width="105px"
            :model="remoteFolderFormData"
            :rules="{
              account: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
              password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
              ip: [{ required: true, message: '请输入IP', trigger: 'blur' }],
              port: [{ required: remoteFolderFormData.readType !== 'smb', message: '请输入端口', trigger: 'blur' }],
            }"
          >
            <el-form-item label="读取方式：" prop="readType">
              <el-select v-model="remoteFolderFormData.readType" @change="handleReadTypeChange">
                <el-option label="SMB读取" value="smb" />
                <el-option label="SFTP读取" value="sftp" />
                <el-option label="FTP读取" value="ftp" />
                <el-option label="FTPS读取" value="ftps" />
              </el-select>
            </el-form-item>
            <el-form-item class="abreast-item" label="IP：" prop="ip">
              <el-input v-model="remoteFolderFormData.ip" clearable />
            </el-form-item>
            <el-form-item class="abreast-item" label="端口：" prop="port">
              <el-input v-model="remoteFolderFormData.port" clearable />
            </el-form-item>
            <el-form-item class="abreast-item" label="目录：">
              <el-input v-model="remoteFolderPath" clearable />
            </el-form-item>
            <el-form-item class="abreast-item" label="用户名：" prop="account">
              <el-input v-model="remoteFolderFormData.account" clearable />
            </el-form-item>
            <el-form-item class="abreast-item" label="密码：" prop="password">
              <el-input v-model="remoteFolderFormData.password" clearable show-password type="password" />
            </el-form-item>
            <el-form-item label="&nbsp;">
              <el-button type="primary" @click="handle2read(true)">确认读取</el-button>
            </el-form-item>
          </el-form>
          <ul v-if="remoteFolderList.length" class="upload-file-list directory">
            <li v-for="(file, index) in remoteFolderList" :key="file.fileName" class="file-item">
              <div class="handle">
                <div class="name">
                  <vab-icon icon="file-text-fill" style="color: #ffc13d; margin-inline: 10px 7px" />
                  {{ file.fileName }}
                </div>
                <el-input
                  v-model="file.remark"
                  clearable
                  placeholder="添加备注"
                  style="height: 40px; margin-top: -1px; margin-right: -1px"
                />
                <el-icon class="delete-icon" @click="deleteFile(remoteFolderList, index)"><Delete /></el-icon>
              </div>
            </li>
          </ul>
        </div>
      </template>
      <template #footer>
        <div class="dialog-footer">
          <el-button :auto-insert-space="false" :loading="uploading" type="primary" @click="handle2Upload">
            确定
          </el-button>
          <el-button :auto-insert-space="false" @click="uploadVisible = false">取消</el-button>
        </div>
      </template>
      <!-- 目录读取弹框 -->
      <el-dialog v-model="remoteFoldersVisible" class="remotefolders-dialog" :show-close="false" width="1100px">
        <div v-loading="remotefoldersLoading" class="remotefolders">
          <div class="remotefolders-left">
            <div class="fixed">
              <el-button
                :auto-insert-space="false"
                link
                style="text-indent: 5px; margin-left: 8px"
                @click="handleGoBack(remoteFolders[0])"
              >
                <img alt="返回上一级" class="folder-icon" :src="gobackIcon" width="14" />
                返回上一级
              </el-button>
            </div>
            <template v-if="remoteFolders.length > 1">
              <div
                v-for="(item, index) in remoteFolders.slice(1)"
                :key="index"
                class="remotefolders-item left-item"
                :class="{ checked: currentCheckedFolder?.filePath === item.filePath }"
                @click="handleFolderClick(item)"
              >
                <img alt="文件" class="folder-icon" :src="item.type === 'folder' ? folderIcon : pcapIcon" width="14" />
                <span class="file-name">{{ item.fileName }}</span>
                <span class="file-size">{{ item.fileSize }}</span>
              </div>
            </template>
            <el-empty v-else class="vab-data-empty" description="暂无数据" />
          </div>
          <div class="remotefolders-center">
            <el-button :auto-insert-space="false" type="primary" @click="handleToSelected">选择</el-button>
          </div>
          <div class="remotefolders-right">
            <div class="fixed"><h4>已选文件：</h4></div>
            <div v-for="(item, index) in selectedFolders" :key="index" class="remotefolders-item right-item">
              <el-checkbox v-model="item.checked" :label="item.fileName" />
            </div>
          </div>
        </div>
        <template #header>
          <div class="my-header" @click="remote_folder_ref?.focus()">
            <el-input
              ref="remote_folder_ref"
              v-model="remoteFolderPath"
              class="remote-folder-path"
              clearable
              @keydown.enter="handleGetRemoteFolders"
            />
          </div>
        </template>
        <template #footer>
          <div class="dialog-footer">
            <el-button :auto-insert-space="false" type="primary" @click="handleSaveRemoteFolders">确定</el-button>
            <el-button :auto-insert-space="false" @click="remoteFoldersVisible = false">取消</el-button>
          </div>
        </template>
      </el-dialog>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .upload {
    width: 790px;
    margin: 0 auto 20px;
    :deep() {
      .el-upload-dragger {
        width: 790px;
        height: 180px;
        background: #fafaff;
      }
      .el-upload__text {
        margin-top: 35px;
        font-weight: 400;
        font-size: 15px;
        .el-icon {
          vertical-align: -3px;
        }
      }
    }
  }
  .upload-file-list {
    padding-inline-start: 0;
    &.directory {
      width: 795px;
      background: #fafaff;
      border-radius: 4px;
      min-height: 240px;
      padding: 20px 0;
      margin-left: 15px;
      .handle {
        margin-inline: 15px 30px !important;
      }
    }
    .file-item {
      height: 40px;
      line-height: 40px;
      display: flex;
      margin-bottom: 10px;
      &:hover {
        .delete-icon {
          display: block;
        }
      }
      .handle {
        width: 100%;
        margin-inline: 20px;
        display: flex;
        border: 1px solid var(--el-border-color);
        background-color: #fff;
        position: relative;
        div {
          flex: 1;
        }
      }
      .delete-icon {
        display: none;
        position: absolute;
        right: -20px;
        top: 11px;
        color: var(--el-color-primary);
        cursor: pointer;
      }
    }
  }

  .directory-upload-form {
    margin-right: 20px;
    :deep() {
      .el-form-item {
        &.abreast-item {
          display: inline-flex;
          .el-select,
          .el-input {
            width: 300px;
          }
        }
        .el-select {
          width: 100%;
        }
        .el-form-item__label {
          padding-right: 2px;
        }
      }
    }
  }
  .remote-folders-box {
    :deep() {
      .el-dialog.remotefolders-dialog {
        .el-dialog__header {
          margin-right: 0;
          padding-block: 10px 0;
        }
        .el-dialog__body {
          margin-bottom: 10px;
          border-bottom: 1px solid var(--el-border-color) !important;
          padding: 0 !important;
        }
      }
      .remotefolders {
        height: 500px;
        display: flex;
        &-left,
        &-right {
          flex: 1;
          padding: 0 10px;
          overflow: auto;
          &::-webkit-scrollbar {
            width: 8px;
            height: 8px;
          }
          .fixed {
            position: sticky;
            top: 0;
            background-color: #fff;
            line-height: 32px;
            z-index: 999;
            h4 {
              margin-bottom: 0;
            }
          }
          .remotefolders-item {
            text-indent: 10px;
            line-height: 32px;
            font-size: 14px;
            height: 32px;
            margin-bottom: 2px;
            &.right-item {
              .el-checkbox__label {
                padding-left: 0;
                font-weight: normal;
                color: var(--el-checkbox-text-color);
              }
            }
            &.left-item {
              border-radius: 6px;
              cursor: pointer;
              display: flex;
              align-items: center;
              .file-name {
                display: inline-block;
                height: 32px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
                flex: 1;
              }
              .file-size {
                align-self: flex-end;
                margin-right: 10px;
              }
              .folder-icon {
                width: 14px;
                height: 14px;
                margin-inline: 10px 0px;
              }
              &:hover:not(.checked) {
                background-color: #f3f1fe;
              }
              &.checked {
                background-color: var(--el-color-primary);
                color: #fff;
              }
            }
          }
        }
        &-center {
          width: 100px;
          text-align: center;
          border: 1px solid var(--el-border-color);
          border-top: none;
          border-bottom: none;
          .el-button {
            margin-top: 239px;
          }
        }
      }
      .remote-folder-path {
        width: 100%;
        .el-input__wrapper {
          box-shadow: none !important;
          padding-inline: 0 !important;
          .el-input__inner {
            font-weight: bold;
            font-size: 18px;
            color: #303133;
          }
        }
      }
    }
  }
</style>
