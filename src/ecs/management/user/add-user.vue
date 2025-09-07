<script lang="ts">
  export default {
    name: 'AddUser',
  }
</script>

<script setup lang="ts">
  import { FormInstance, FormItemRule, FormRules } from 'element-plus'
  import { AddUserParam } from '@/types'
  import { formatNstime, formatTime } from '@/utils/time'
  import { useSettingsStore } from '@/store/modules/settings'
  import { addUserApi } from '@/api-ecs/system'
  import { userTypes } from '@/data/constant'
  import { isEmail, isPhone } from '@/utils/validate'
  import sha1 from 'sha1'
  import dayjs from 'dayjs'
  const settingsStore = useSettingsStore()
  const nimLength = settingsStore.getSecConfig().minPwdLen || 8
  const maxLength = settingsStore.getSecConfig().maxLength || 32
  const capitalChar = settingsStore.getSecConfig().capitalChar || false // 大写
  const ordinaryChar = settingsStore.getSecConfig().ordinaryChar || false
  const numChar = settingsStore.getSecConfig().numChar || false
  const specialChar = settingsStore.getSecConfig().specialChar || false
  const $baseMessage: any = inject('$baseMessage')
  const props = defineProps<{
    isShow: boolean
    roles: { roleName: string; id: number }[]
    depts: { deptName: string; id: number }[]
  }>()

  const emits = defineEmits<{
    (e: 'on-cancel', load: boolean): void
  }>()

  const activeName = ref('addUser') // tabs选中项
  const rules = {
    nickName: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    loginName: [{ required: true, message: '请输入登录名', trigger: 'blur' }],
    deptId: [{ required: true, message: '请选择部门', trigger: 'blur' }],
    roleIds: [{ required: true, message: '请选择角色', trigger: 'blur' }],
    password: [
      { required: true, message: '请输入登录密码', trigger: 'blur' },
      { message: `登录密码长度${nimLength}-${maxLength}位`, trigger: 'blur', min: nimLength, max: maxLength },
      { pattern: capitalChar ? /^(?=.*[A-Z])/ : '', message: capitalChar ? '必须包含大写' : '包含大写' },
      { pattern: ordinaryChar ? /^(?=.*[a-z])/ : '', message: ordinaryChar ? '必须包含小写' : '包含小写' },
      { pattern: numChar ? /^(?=.*[0-9])/ : '', message: numChar ? '必须包含数字' : '包含数字' },
      {
        pattern: specialChar ? /^(?=.*[./~!@#$%^&*)(_+}{|:?><])/ : '',
        message: specialChar ? '必须包含特殊符号' : '特殊符号',
      },
    ],
    phone: [{ trigger: 'blur', validator: validatePhone }],
    mail: [{ trigger: 'blur', validator: validateMail }],
  }
  const queryForm = reactive({
    timeDate: undefined,
    confirmon: '',
    showtime: false,
  })
  const formData = reactive<AddUserParam>({
    nickName: '',
    loginName: '',
    description: '',
    password: '',
    effectiveSt: '',
    effectiveEd: '',
    userType: 1,
    roleIds: [],
    deptId: props.depts[0].id,
    phone: '',
    mail: '',
    address: '',
    ipWhiteList: '',
    limitIp: '',
    twoFactorAuth: true,
    visibleEncryptData: true,
    visibleBuiltinData: true,
  }) // 表单数据
  const formRef = ref<FormInstance>() // 表单实例

  // 保存
  const saveData = () => {
    if ((!formData.password && !queryForm.confirmon) || formData.password !== queryForm.confirmon)
      return $baseMessage('两次密码不一致，请重新输入', 'error', 'vab-hey-message-error')
    if (!formData.effectiveEd) return $baseMessage('请配置有效时间', 'error', 'vab-hey-message-error')
    if (formData.userType === 3 && formData.roleIds.length === 0)
      return $baseMessage('请选择角色', 'error', 'vab-hey-message-error')
    formRef.value?.validate(async (valid) => {
      if (valid) {
        const { msg } = await addUserApi({
          ...formData,
          password: sha1(formData.password),
          roleIds: formData.userType === 3 ? formData.roleIds : [],
        })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        emits('on-cancel', true)
      }
    })
  }
  function validatePhone(rule: any, value: any, callback: any) {
    if (value && !isPhone(value)) {
      return callback(new Error('请输入合法手机号'))
    }
    return callback()
  }
  function validateMail(rule: any, value: any, callback: any) {
    if (value && !isEmail(value)) {
      return callback(new Error('请输入正确邮箱'))
    }
    return callback()
  }
  function setTimeDateHandle() {
    const nowdate = dayjs()
    queryForm.timeDate = [formatTime(nowdate), formatTime(nowdate.add(99, 'year'))] as any
  }
  watchEffect(() => {
    if (queryForm.timeDate) {
      const [ST, END] = queryForm.timeDate as any[]
      formData.effectiveSt = ST
      formData.effectiveEd = END
    }
  })
</script>

<template>
  <div class="add-user">
    <el-tabs v-model="activeName">
      <el-tab-pane label="添加用户" name="addUser">
        <el-row>
          <el-col :span="4" />
          <el-col :span="12">
            <el-form
              ref="formRef"
              class="login-form"
              label-position="right"
              label-width="220px"
              :model="formData"
              :rules="rules"
            >
              <el-form-item label="登录名称" prop="loginName">
                <el-input v-model="formData.loginName" />
              </el-form-item>
              <el-form-item label="真实姓名" prop="nickName">
                <el-input v-model="formData.nickName" />
              </el-form-item>
              <el-form-item label="用户描述" prop="description">
                <el-input v-model="formData.description" />
              </el-form-item>
              <el-form-item label="密码" prop="password">
                <el-input v-model="formData.password" show-password type="password" />
              </el-form-item>
              <el-form-item label="确认密码" required>
                <el-input v-model="queryForm.confirmon" show-password type="password" />
              </el-form-item>
              <el-form-item label="有效时间配置" required>
                <el-date-picker
                  v-model="queryForm.timeDate"
                  end-placeholder="结束日期"
                  start-placeholder="开始日期"
                  type="daterange"
                />
                <el-button text type="primary" @click="setTimeDateHandle">永久有效</el-button>
              </el-form-item>
              <el-form-item label="部门" prop="deptId">
                <el-select v-model="formData.deptId" filterable style="width: 100%">
                  <el-option v-for="dept in depts" :key="dept.id" :label="dept.deptName" :value="dept.id" />
                </el-select>
              </el-form-item>
              <el-form-item label="用户类型" prop="userType">
                <el-radio-group v-model="formData.userType">
                  <el-radio v-for="user in userTypes" :key="user.lable" :label="user.lable">{{ user.name }}</el-radio>
                </el-radio-group>
              </el-form-item>
              <el-form-item v-if="formData.userType === 3" label="角色" prop="roleIds">
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
              <el-form-item label="手机号" prop="phone">
                <el-input v-model="formData.phone" />
              </el-form-item>
              <el-form-item label="电子邮箱" prop="mail">
                <el-input v-model="formData.mail" />
              </el-form-item>
              <el-form-item label="限制登录IP" prop="limitIp">
                <el-input v-model="formData.limitIp" />
              </el-form-item>
              <el-form-item class="btn">
                <el-button type="primary" @click="saveData">保存</el-button>
                <el-button @click="emits('on-cancel', false)">取消</el-button>
              </el-form-item>
            </el-form>
          </el-col>
          <el-col :span="7" />
        </el-row>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style lang="scss" scoped>
  :deep() {
    .el-tabs__content {
      height: calc(100vh - 100px);
      overflow-y: auto;
    }
  }
  .add-user {
  }
  .el-select-dropdown__item {
    max-width: 400px;
  }
</style>
