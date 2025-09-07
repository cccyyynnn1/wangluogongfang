<script lang="ts">
  export default {
    name: 'ConfigManagement',
  }
</script>
<script setup lang="ts">
  import { ElInput, FormInstance, UploadProps, UploadRawFile, UploadUserFile, genFileId } from 'element-plus'
  import { Plus, Close } from '@element-plus/icons-vue'
  import { getToken } from '@/utils/token'
  import { updateSystemDynamic, addSystemDynamic } from '@/api-ecs/system'
  import { useSettingsStore } from '@/store/modules/settings'
  import AesEncryptCBC from '~/src/utils/crypto'
  const settingsStore = useSettingsStore()
  const { systemConfig } = storeToRefs(settingsStore)
  const $baseMessage: any = inject('$baseMessage')
  const submitFormRef = ref<FormInstance>()
  const rules = {
    // name: [{ required: true, message: '请输入系统名称', trigger: 'change' }],
    nameAbb: [{ required: true, message: '请输入系统名称简称', trigger: 'change' }],
    // logo: [{ required: true, message: '请上传系统LOGO', trigger: 'change' }],
    // loginLogo: [{ required: true, message: '请上传登录页LOGO', trigger: 'change' }],
    // loginBg: [{ required: true, message: '请上传登录背景图', trigger: 'change' }],
  }
  /**查询的ID，有ID表示更新 */
  const systemInfo = reactive({
    id: undefined as number | undefined,
    uid: undefined as number | undefined,
    key: '',
    name: '',
    value: '',
  })

  // 表单数据
  const submitForm = reactive({
    name: '',
    nameAbb: '',
    icp: '',
    copyright: '',
    loginLogo: '',
    logginTip: '',
    logo: '',
    loginBg: '',
  })

  const logUrl = ref<UploadUserFile[]>([])
  const loginBgUrl = ref<UploadUserFile[]>([])
  const loginLogoUrl = ref<UploadUserFile[]>([])

  function handleUploadSuccess(response: any) {
    if (response.code === 20) {
      submitForm.logo = response.data
    } else {
      $baseMessage(response.msg, 'error', 'vab-hey-message-error')
    }
  }
  function handleUploadLBSuccess(response: any) {
    if (response.code === 20) {
      submitForm.loginLogo = response.data
    } else {
      $baseMessage(response.msg, 'error', 'vab-hey-message-error')
    }
  }
  function handleBbUploadSuccess(response: any) {
    if (response.code == 20) {
      submitForm.loginBg = response.data
    } else {
      $baseMessage(response.msg, 'error', 'vab-hey-message-error')
    }
  }

  // 保存
  const saveData = async () => {
    if (!submitFormRef.value) return
    await submitFormRef.value.validate(async (valid) => {
      if (valid) {
        systemInfo.value = JSON.stringify(submitForm)
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
              h('div', { style: 'margin: 15px 0 5px 0;color:#55585b' }, '请输入敏感操作密码:(通过验证后方可进行操作)'),
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
                if (systemInfo.id && systemInfo.id > 0) {
                  const { msg } = await updateSystemDynamic(systemInfo, { password: AesEncryptCBC(password.value) })
                  $baseMessage(msg, 'success', 'vab-hey-message-success')
                } else {
                  const { msg } = await addSystemDynamic(submitForm)
                  $baseMessage(msg, 'success', 'vab-hey-message-success')
                }
                settingsStore.getSystemConfig()
                done()
              } catch (error) {
                instance.confirmButtonLoading = false
              }
            } else {
              done()
            }
          },
        }).catch(() => {})
      }
    })
  }

  function resetForm() {
    submitFormRef.value?.resetFields()
  }

  watch(
    () => systemConfig.value,
    () => {
      const {
        id,
        key,
        uid,
        name: varName,
        value: { name, nameAbb, logo, loginBg, icp, copyright, logginTip },
      } = toRaw(systemConfig.value)
      submitForm.copyright = copyright
      submitForm.name = name
      submitForm.nameAbb = nameAbb
      submitForm.logo = logo
      submitForm.icp = icp
      submitForm.loginBg = loginBg
      // submitForm.loginLogo = loginLogo
      submitForm.logginTip = logginTip

      systemInfo.id = id || undefined
      systemInfo.uid = uid || undefined
      systemInfo.key = key || 'system_config'
      systemInfo.name = varName || ''
      // @ts-ignore
      logUrl.value = logo ? [{ name: 'Logo', url: logo }] : undefined

      // loginLogoUrl.value = loginLogo ? [{ name: '登录页LOGO', url: logo }] : undefined
      // @ts-ignore
      loginBgUrl.value = loginBg ? [{ name: '登录背景', url: loginBg }] : undefined
    },
    { immediate: true }
  )

  const upload = ref()
  const uploadBG = ref()
  const handleExceed: UploadProps['onExceed'] = (files) => {
    upload.value!.clearFiles()
    const file = files[0] as UploadRawFile
    file.uid = genFileId()
    upload.value!.handleStart(file)
    upload.value!.submit()
  }
  const handleBGExceed: UploadProps['onExceed'] = (files) => {
    uploadBG.value!.clearFiles()
    const file = files[0] as UploadRawFile
    file.uid = genFileId()
    uploadBG.value!.handleStart(file)
    uploadBG.value!.submit()
  }

  const activeRef = ref(false)
  const activeLogoRef = ref(false)

  const handleLogoDEel = () => {
    uploadBG.value!.clearFiles()
    submitForm.logo = ''
    // @ts-ignore
    logUrl.value = undefined
  }

  const beforeAvatarUpload: UploadProps['beforeUpload'] = (rawFile) => {
    if (rawFile.type !== 'image/jpeg' && rawFile.type !== 'image/gif' && rawFile.type !== 'image/png') {
      ElMessage.error('图片格式不正确!')
      return false
    }
    return true
  }
</script>

<template>
  <div class="config-management-container">
    <div class="content">
      <vab-query-form>
        <vab-query-form-left-panel :span="12">
          <h3>系统信息</h3>
        </vab-query-form-left-panel>
      </vab-query-form>
      <el-row>
        <el-col :span="3" />
        <el-col :span="17">
          <el-form ref="submitFormRef" label-position="right" label-width="130px" :model="submitForm" :rules="rules">
            <el-form-item label="主页简称：" prop="nameAbb">
              <el-input v-model.trim="submitForm.nameAbb" placeholder="请输入主页简称" style="width: 50%" />
            </el-form-item>
            <!-- <el-form-item label="主页LOGO：" prop="logo">
              <el-upload
                ref="upload"
                v-model:file-list="logUrl"
                accept="image/jpeg,image/gif,image/png"
                action="/v3/ecsPlatform/public/uploadIndexImg"
                :auto-upload="true"
                class="upload-demo"
                :headers="{ Cookies: `EcsSessionId=${getToken()}` }"
                :limit="1"
                list-type="picture"
                name="file"
                :on-exceed="handleExceed"
                :on-remove="() => (submitForm.logo = '')"
                :on-success="handleUploadSuccess"
              >
                <el-button class="upload-btn" :icon="Plus">上传</el-button>
                <span class="upload-tip">用于系统全局logo展示</span>
              </el-upload>
            </el-form-item> -->
            <el-form-item label="登录页信息：" prop="name">
              <div
                class="login-info"
                :class="{ active: activeRef }"
                :style="{
                  background: `url(${submitForm.loginBg || require('@/assets/login_images/login.png')})   no-repeat`,
                }"
                @mouseenter="activeRef = true"
                @mouseleave="activeRef = false"
              >
                <div class="login-logo">
                  <el-upload
                    ref="upload"
                    v-model:file-list="logUrl"
                    accept="image/jpeg,image/gif,image/png"
                    action="/v3/ecsPlatform/public/uploadIndexImg"
                    :auto-upload="true"
                    :before-upload="beforeAvatarUpload"
                    class="upload-demo"
                    :class="{ logoInputActive: activeRef }"
                    :headers="{ Cookies: `EcsSessionId=${getToken()}` }"
                    :limit="1"
                    list-type="picture"
                    name="file"
                    :on-exceed="handleExceed"
                    :on-remove="() => (submitForm.logo = '')"
                    :on-success="handleUploadSuccess"
                    @mouseenter="activeLogoRef = true"
                    @mouseleave="activeLogoRef = false"
                  >
                    <div class="logo-upload-tip">
                      <el-image
                        v-if="activeLogoRef || !submitForm.logo"
                        class="login-bg-upload"
                        :src="require('@/assets/login_images/plus-upload.svg')"
                      />
                      <!-- <el-icon class="login-bg-upload" style="font-size: 16px">
                        <Plus />
                      </el-icon> -->
                      <div v-if="!submitForm.logo" class="el-upload__text" style="margin-left: 3px">logo</div>
                    </div>
                    <div v-if="submitForm.logo" class="previews">
                      <img alt="" class="bg-logo" :class="{ previewsActive: activeLogoRef }" :src="submitForm.logo" />
                      <div v-if="activeLogoRef" class="selectTag-closeable" @click.stop="handleLogoDEel">
                        <Close />
                      </div>
                    </div>
                  </el-upload>
                  <el-input
                    v-model.trim="submitForm.logginTip"
                    class="login-tip logo-input"
                    :class="{ logoInputActive: activeRef }"
                    placeholder="请输入登录页系统名称"
                    style="margin-left: 6px"
                  />
                </div>
                <el-upload
                  ref="uploadBG"
                  v-model:file-list="loginBgUrl"
                  accept="image/jpeg,image/gif,image/png"
                  action="/v3/ecsPlatform/public/uploadIndexImg"
                  :auto-upload="true"
                  class="login-bg"
                  :class="{ logoInputActive: activeRef }"
                  :headers="{ Cookies: `EcsSessionId=${getToken()}` }"
                  :limit="1"
                  name="file"
                  :on-exceed="handleBGExceed"
                  :on-remove="() => (submitForm.loginBg = '')"
                  :on-success="handleBbUploadSuccess"
                >
                  <div class="up-load-tip">
                    <el-image class="login-bg-upload" :src="require('@/assets/login_images/plus-upload.svg')" />
                    <!-- <el-icon  ><Plus /></el-icon> -->
                    <div class="el-upload__text">上传登录页背景图( 只能上传png/jpg/jpeg)</div>
                  </div>
                </el-upload>
                <el-input
                  v-model.trim="submitForm.name"
                  class="login-mame logo-input"
                  :class="{ logoInputActive: activeRef }"
                  placeholder="请输入登录框标题"
                />
                <el-input
                  v-model.trim="submitForm.icp"
                  class="login-icp logo-input"
                  :class="{ logoInputActive: activeRef }"
                  placeholder="请输入登录页底部信息"
                />
              </div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveData">保存</el-button>
            </el-form-item>
          </el-form>
        </el-col>
        <el-col :span="4" />
      </el-row>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .previews {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .previewsActive {
    opacity: 0.4;
  }
  .selectTag-closeable {
    position: absolute;
    top: 50%;
    left: calc(50% + 26px);
    transform: translate(-50%, -50%);
    z-index: 13;
    width: 12px;
    height: 12px;
    line-height: 12px;
    color: #fff;
    text-align: center;

    background-color: #ff6060;
    border-radius: 100%;
    svg {
      width: 10px;
      height: 10px;
    }
  }
  .bg-logo {
    object-fit: contain;
    width: 28px;
    height: 28px;
    z-index: 11;
    background-size: 100% 100% !important;
  }
  :deep() {
    .login-bg {
      position: absolute;
      z-index: 11;
      left: 3%;
      top: 50%;
      transform: translateY(-50%);
      width: 55%;
      height: 55%;
      border: 1px dashed #fff;
      border-radius: 8px;
      .login-bg-upload {
        height: 36px;
        width: 36px;
      }
      .el-upload {
        height: 100%;
        width: 100%;
        color: #fff !important;
      }
      .el-upload-list {
        display: none;
      }
      .up-load-tip {
        height: 100%;
        width: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        .el-upload__text {
          margin-top: 8px;
          font-size: 15px;
          font-weight: 600;
        }
      }
    }
    .logoInputActive {
      border: 1px dashed #d5c5ff !important;
    }
    .login-logo {
      top: 5%;
      left: 3%;
      position: absolute;
      display: flex;
      align-items: center;
      width: 40%;
      height: 34px;
      z-index: 11;
      .login-bg-upload {
        height: 16px;
        width: 16px;
      }

      .upload-demo {
        border: 1px dashed #fff;
        height: 100%;
        border-radius: 8px;
        width: 72px;
      }
      .el-upload {
        position: relative;
        height: 100%;
        width: 100%;
        color: #fff !important;
      }
      .logo-input {
        width: 80%;
      }
      .el-upload-list {
        display: none;
      }
      .logo-upload-tip {
        z-index: 12;
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        height: 100%;
        width: 100%;
        display: flex;
        // flex-direction: column;
        align-items: center;
        justify-content: center;
        .el-upload__text {
          font-size: 15px;
          font-weight: 600;
        }
      }
    }
    .login-mame {
      position: absolute;
      width: 30%;
      top: 33%;
      right: 6%;
      transform: translateY(-50%);
    }
    .login-icp {
      position: absolute;
      width: 40%;
      left: 50%;
      transform: translateX(-50%);
      bottom: 5%;
    }
    .logo-input {
      z-index: 11;
      border-radius: 8px;
      border: 1px dashed #ffffff;
      .el-input__wrapper {
        background-color: transparent !important;
        border: 0px solid #fff !important;
        box-shadow: initial;
        height: 3% !important;
        .el-input__inner {
          color: #fff !important;
          font-size: 14px;
          font-weight: 600;
          text-align: center;
          &::placeholder {
            color: #fff !important;
          }
        }
      }
    }
    .login-info {
      width: 100%;
      height: 0;
      padding-bottom: 56%;
      object-fit: contain !important;
      background-size: 100% 100% !important;
      position: relative;
      z-index: 10;
      border-radius: 8px;
      border: 2px solid #ffffff;
      overflow: hidden;
      &::after {
        z-index: 9;
        content: ' ';
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        position: absolute;
        opacity: 0.6;
        background: #c7c2e8;
      }
      &::before {
        z-index: 9;
        content: ' ';
        top: 50%;
        transform: translateY(-50%);
        right: 6%;
        height: 40%;
        width: 0;
        padding-right: 30%;

        position: absolute;
        background: url('@/assets/login_images/login-form.png') no-repeat;
        object-fit: contain;
        background-size: 100% 100%;
      }
    }
    .active {
      &::after {
        background: #3a326e;
      }
    }
  }

  :deep(.el-upload-list__item-file-name) {
    color: #3c394f !important;
  }
  .config-management-container {
    .content {
      height: calc(100vh - 80px);
      overflow-y: auto;
      h3 {
        margin-block: 0 0.5em;
      }
      .avatar-uploader .avatar {
        display: block;
        width: 128px;
        height: 128px;
      }
    }
  }
  .upload-tip {
    margin-left: 15px;
    font-weight: 400;
    font-size: 14px;
    color: #b7b5bf;
  }
  .upload-btn {
    width: 86px;
    height: 34px;
    background: #fafaff;
    border-radius: 4px;
    border: 1px dashed #b4a8ff;
    font-weight: 500;
    font-size: 14px;
    color: #6954f0;
  }
  :deep() {
    .avatar-uploader .el-upload {
      position: relative;
      overflow: hidden;
      cursor: pointer;
      border: 1px dashed var(--el-border-color);
      border-radius: 6px;
      transition: var(--el-transition-duration-fast);
    }

    .avatar-uploader .el-upload:hover {
      border-color: var(--el-color-primary);
    }

    .el-icon.avatar-uploader-icon {
      width: 128px;
      height: 128px;
      font-size: 28px;
      color: #8c939d;
      text-align: center;
    }

    .upload-demo {
      width: 50%;
      .el-upload-list__item {
        height: 60px;
        .el-upload-list__item-thumbnail {
          height: 45px;
        }
      }
    }
  }
</style>
