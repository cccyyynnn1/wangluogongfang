<script lang="ts">
  export default {
    name: 'EditUser',
  }
</script>
<script setup lang="ts">
  import { editUserApi } from '@/api-ecs/system'
  import { Role } from '@/types'
  import { isPhone, isEmail } from '@/utils/validate'
  import { useSettingsStore } from '@/store/modules/settings'
  import { FormInstance } from 'element-plus'
  import { useUserStore } from '@/store/modules/user'
  import { userTypes } from '@/data/constant'
  import sha1 from 'sha1'
  const settingsStore = useSettingsStore()
  const userStore = useUserStore()
  const $baseMessage: any = inject('$baseMessage')
  const $baseConfirm: any = inject('$baseConfirm')
  const props = defineProps<{
    dialogVisible: boolean
    dialogVal: Role
    roles: { roleName: string; id: number }[]
    depts: { deptName: string; id: number }[]
  }>()

  const emit = defineEmits<{
    (e: 'on-closeEvent', load: boolean): void
  }>()
  const formRef = ref<FormInstance>() // 表单实例
  const formData = ref()
  const visible = ref(false)
  const passwordData = reactive({
    reEnterPassword: '',
    password: '',
    oldPassword: '',
  })

  const nimLength = settingsStore.getSecConfig().minPwdLen || 8
  const maxLength = settingsStore.getSecConfig().maxLength || 32
  const capitalChar = settingsStore.getSecConfig().capitalChar || false // 大写
  const ordinaryChar = settingsStore.getSecConfig().ordinaryChar || false
  const numChar = settingsStore.getSecConfig().numChar || false
  const specialChar = settingsStore.getSecConfig().specialChar || false
  watchEffect(() => {
    const { cloned } = useCloned(props.dialogVal)
    formData.value = cloned.value
  })

  const validatePhone = (rule: any, value: any, callback: any) => {
    if (value && !isPhone(value)) {
      return callback(new Error('请输入合法手机号'))
    }
    return callback()
  }
  const validateMail = (rule: any, value: any, callback: any) => {
    if (value && !isEmail(value)) {
      return callback(new Error('请输入正确邮箱'))
    }
    return callback()
  }
  onMounted(() => {
    visible.value = props.dialogVisible
  })

  // 保存
  const save = async () => {
    const basePassword = 'admin@123'
    formRef.value?.validate(async (valid) => {
      if (passwordData.password !== passwordData.reEnterPassword) {
        return $baseMessage('新密码与确认密码不一致', 'error', 'vab-hey-message-error')
      }
      if (
        passwordData.password &&
        (passwordData.password.length < nimLength || passwordData.password.length > maxLength)
      ) {
        return $baseMessage(`新密码长度${nimLength}-${maxLength}位`, 'error', 'vab-hey-message-error')
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
      if (valid) {
        const { id, nickName, deptId, mail, phone, loginName, roleIds, userType, limitIp } = formData.value
        // @ts-ignore
        let newPassword = null
        // @ts-ignore
        let oldPassword = null
        if (passwordData.password) newPassword = sha1(passwordData.password)
        if (passwordData.oldPassword) oldPassword = sha1(passwordData.oldPassword)
        const { msg, code } = await editUserApi({
          id,
          nickName,
          deptId,
          mail,
          phone,
          limitIp,
          loginName,
          roleIds: userType === 3 ? roleIds : [],
          userType,
          password: newPassword,
          oldPassword: oldPassword,
          confirm: false,
        })
        if (msg == '限制登录IP不包含当前ip地址，如保存则此设备将立即无法使用本平台，是否确认') {
          $baseConfirm(msg, null, async () => {
            const { msg, code } = await editUserApi({
              id,
              nickName,
              deptId,
              mail,
              phone,
              limitIp,
              loginName,
              roleIds: userType === 3 ? roleIds : [],
              userType,
              // @ts-ignore
              password: newPassword,
              // @ts-ignore
              oldPassword: oldPassword,
              confirm: true,
            })
            $baseMessage(msg, 'success', 'vab-hey-message-success')
            if (code === 20) emit('on-closeEvent', true)
          })
        } else {
          $baseMessage(msg, 'success', 'vab-hey-message-success')
          if (code === 20) emit('on-closeEvent', true)
        }
      }
    })
  }
</script>

<template>
  <div class="edit-user">
    <el-dialog
      v-model="visible"
      :before-close="() => emit('on-closeEvent', false)"
      :close-on-click-modal="false"
      title="编辑用户"
      width="35%"
    >
      <div class="content">
        <el-form
          ref="formRef"
          class="login-form"
          label-position="right"
          label-width="120px"
          :model="formData"
          :rules="{
            mail: [{ trigger: 'blur', validator: validateMail }],
            phone: [{ trigger: 'blur', validator: validatePhone }],
            password: [{ message: '登录密码长度8-32位', trigger: 'blur', min: 8, max: 32 }],
          }"
        >
          <div class="item">
            <h4>用户基本信息</h4>
            <el-form-item label="登录名称：" prop="loginName">
              <el-input v-model="formData.loginName" placeholder="32个字符以内" readonly />
            </el-form-item>

            <el-form-item label="用户类型：">
              <el-radio-group v-model="formData.userType">
                <el-radio v-for="user in userTypes" :key="user.lable" :label="user.lable">{{ user.name }}</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item v-if="formData.userType === 3" label="角色名称：" prop="roleIds">
              <el-select
                v-model="formData.roleIds"
                collapse-tags
                collapse-tags-tooltip
                filterable
                multiple
                style="width: 100%"
              >
                <el-option v-for="role in roles" :key="role.roleName" :label="role.roleName" :value="role.id" />
              </el-select>
            </el-form-item>
            <el-form-item v-if="formData?.id === userStore.userId" label="旧密码：">
              <el-input v-model="passwordData.oldPassword" show-password type="password" />
            </el-form-item>
            <el-form-item label="设置密码：">
              <el-input v-model="passwordData.password" show-password type="password" />
            </el-form-item>
            <el-form-item label="再次输入密码：">
              <el-input v-model="passwordData.reEnterPassword" show-password type="password" />
            </el-form-item>
            <el-form-item label="限制登录IP：" prop="limitIp">
              <el-input v-model="formData.limitIp" />
            </el-form-item>
          </div>
          <div class="item">
            <h4>账号信息</h4>
            <el-form-item label="用户名：" prop="nickName">
              <el-input v-model="formData.nickName" />
            </el-form-item>
            <el-form-item label="部门：" prop="deptId">
              <el-select v-model="formData.deptId" filterable style="width: 100%">
                <el-option v-for="dept in depts" :key="dept.id" :label="dept.deptName" :value="dept.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="邮箱：" prop="mail">
              <el-input v-model="formData.mail" />
            </el-form-item>
            <el-form-item label="手机：" prop="phone">
              <el-input v-model="formData.phone" />
            </el-form-item>
          </div>
        </el-form>
      </div>
      <template #footer>
        <el-button type="primary" @click="save">保存</el-button>
        <el-button @click="emit('on-closeEvent', false)">取消</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .content {
    // padding: 0 20px;

    .item {
      padding: 0px 40px 20px 20px;
      background: #f8fbff;
      border-radius: 2px;

      &:nth-child(2) {
        margin-top: 20px;
      }
    }
  }
  .el-select-dropdown__item {
    max-width: 400px;
  }
</style>
