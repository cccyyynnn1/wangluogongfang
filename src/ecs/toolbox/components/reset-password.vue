<script lang="ts">
  export default {
    name: 'ResetPassword',
  }
</script>

<script setup lang="ts">
  import { useUserStore } from '@/store/modules/user'
  import { resetPasswordApi } from '@/api-ecs/login'
  import { useSettingsStore } from '@/store/modules/settings'
  import router, { resetRouter } from '@/router'
  import sha1 from 'sha1'
  const settingsStore = useSettingsStore()
  const $baseMessage: any = inject('$baseMessage')
  const userStore = useUserStore()
  const passwordData = reactive({
    oldPassword: '',
    password: '',
    enterPassword: '',
  })
  const nimLength = settingsStore.getSecConfig()?.minPwdLen || 8
  const maxLength = settingsStore.getSecConfig()?.maxLength || 32
  const capitalChar = settingsStore.getSecConfig()?.capitalChar || false // 大写
  const ordinaryChar = settingsStore.getSecConfig()?.ordinaryChar || false
  const numChar = settingsStore.getSecConfig()?.numChar || false
  const specialChar = settingsStore.getSecConfig()?.specialChar || false
  const basePassword = 'admin@123'
  const resetPassword = async () => {
    if (passwordData.password !== passwordData.enterPassword) {
      return $baseMessage('新密码与确认密码不一致', 'error', 'vab-hey-message-error')
    }
    if (
      passwordData.password &&
      (passwordData.password.length < nimLength || passwordData.password.length > maxLength)
    ) {
      return $baseMessage(`新密码要求长度${nimLength}-${maxLength}位`, 'error', 'vab-hey-message-error')
    }
    if (passwordData.password.toLowerCase() === basePassword) {
      return $baseMessage(`当前密码包含${basePassword}`, 'error', 'vab-hey-message-error')
    }
    if (capitalChar) {
      const reg = new RegExp(/^(?=.*[A-Z])/)
      const res = reg.exec(passwordData.password)
      if (!res) {
        return $baseMessage(`必须包含大写`, 'error', 'vab-hey-message-error')
      }
    }
    if (ordinaryChar) {
      const reg = new RegExp(/^(?=.*[a-z])/)
      const res = reg.exec(passwordData.password)
      if (!res) {
        return $baseMessage(`必须包含小写`, 'error', 'vab-hey-message-error')
      }
    }
    if (numChar) {
      const reg = new RegExp(/^(?=.*[0-9])/)
      const res = reg.exec(passwordData.password)
      if (!res) {
        return $baseMessage(`必须包含数字`, 'error', 'vab-hey-message-error')
      }
    }
    if (specialChar) {
      const reg = new RegExp(/^(?=.*[./~!@#$%^&*)(_+}{|:?><])/)
      const res = reg.exec(passwordData.password)
      if (!res) {
        return $baseMessage(`必须包含特殊字符`, 'error', 'vab-hey-message-error')
      }
    }
    try {
      const { msg, code } = await resetPasswordApi({
        password: sha1(passwordData.password),
        oldPassword: sha1(passwordData.oldPassword),
      })
      $baseMessage('修改密码成功', 'success', 'vab-hey-message-success')
      settingsStore.changeToolboxVisible(false)
      userStore.logout()
      // router.replace('/login')
    } catch (error) {
      console.log(error)
    }
  }
</script>

<template>
  <div class="reset-password">
    <!-- <el-alert
      v-if="userStore.expire"
      :closable="false"
      effect="dark"
      title="账号长时间不登录，需重置密码后重新登录！"
      type="warning"
    /> -->
    <el-form class="reset-password-form" label-position="right" label-width="120px">
      <el-form-item label="旧密码：">
        <el-input v-model="passwordData.oldPassword" show-password type="password" />
      </el-form-item>
      <el-form-item label="设置密码：">
        <el-input v-model="passwordData.password" show-password type="password" />
      </el-form-item>
      <el-form-item label="再次输入密码：">
        <el-input v-model="passwordData.enterPassword" show-password type="password" />
      </el-form-item>
      <el-form-item label="">
        <el-button type="primary" @click="resetPassword">保存</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<style scoped lang="scss">
  .reset-password {
    padding: 0 25px;
    .reset-password-form {
      margin-bottom: 25px;
      :deep() {
        .el-form-item__content {
          justify-content: end;
        }
        .el-button--primary {
          background: var(--el-color-primary);
          border: 0;
          &:hover {
            opacity: 0.9;
          }
        }
      }
    }
  }
</style>
