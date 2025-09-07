<script lang="ts">
  export default {
    name: 'SafeConfig', // 安全配置
  }
</script>

<script setup lang="ts">
  import type { FormInstance } from 'element-plus'

  import { SecConfigUpdateApi, getSecConfigApi } from '~/src/api-ecs/system'

  import { SecConfigUpdateModel } from '~/src/types'

  import { requireRules } from '~/src/utils/rules'

  import { useSettingsStore } from '@/store/modules/settings'
  import AesEncryptCBC from '~/src/utils/crypto'

  const settingsStore = useSettingsStore()

  // 表单数据
  const formData = reactive<SecConfigUpdateModel>({
    capitalChar: false,
    changeDays: 0,
    changePwd: false,
    expireTime: 0,
    ipWhilteList: '',
    loginFailLock: false,
    loginIpLimit: false,
    loginTimes: 0,
    minPwdLen: 0,
    numChar: false,
    ordinaryChar: false,
    overtimeConfig: false,
    specialChar: false,
    unlockTime: 0,
    maxPwdLen: 32,
  })

  const rules = reactive({
    capitalChar: requireRules,
    changeDays: requireRules,
    changePwd: requireRules,
    expireTime: requireRules,
    // ipWhilteList: requireRules,
    loginFailLock: requireRules,
    loginIpLimit: requireRules,
    loginTimes: requireRules,
    minPwdLen: requireRules,
    numChar: requireRules,
    ordinaryChar: requireRules,
    overtimeConfig: requireRules,
    specialChar: requireRules,
    unlockTime: requireRules,
  })

  const formRef = ref<FormInstance>()

  const submitForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.validate(async (valid) => {
      if (valid) {
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
                const { msg } = await SecConfigUpdateApi({ ...formData }, { password: AesEncryptCBC(password.value) })
                settingsStore.updataSecConfig(formData)
                ElMessage({ message: msg, type: 'success' })
                done()
              } catch (error) {
                instance.confirmButtonLoading = false
              }
            } else {
              done()
            }
          },
        }).catch(() => {})
      } else {
        console.log('error submit!')
        return false
      }
    })
  }

  const init = async () => {
    const { data } = await getSecConfigApi()
    for (const key in formData) {
      // @ts-ignore
      formData[key] = data[key]
    }
  }

  onMounted(() => {
    init()
  })
</script>

<template>
  <div class="safe-config">
    <div class="content">
      <el-row>
        <el-col :span="5" />
        <el-col :span="12">
          <el-form ref="formRef" label-position="right" label-width="180px" :model="formData" :rules="rules">
            <el-form-item label="密码长度：" prop="minPwdLen">
              <el-input-number v-model="formData.minPwdLen" :max="32" :min="8" type="number" />
            </el-form-item>
            <el-form-item label="是否包含大写：" prop="capitalChar">
              <el-radio-group v-model="formData.capitalChar" class="ml-4">
                <el-radio :label="true" size="large">是</el-radio>
                <el-radio :label="false" size="large">否</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="是否包含小写：" prop="ordinaryChar">
              <el-radio-group v-model="formData.ordinaryChar" class="ml-4">
                <el-radio :label="true" size="large">是</el-radio>
                <el-radio :label="false" size="large">否</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="是否包含数字：" prop="numChar">
              <el-radio-group v-model="formData.numChar" class="ml-4">
                <el-radio :label="true" size="large">是</el-radio>
                <el-radio :label="false" size="large">否</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="是否包含特殊字符：" prop="specialChar">
              <el-radio-group v-model="formData.specialChar" class="ml-4">
                <el-radio :label="true" size="large">是</el-radio>
                <el-radio :label="false" size="large">否</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="登录失败是否锁定：" prop="loginFailLock">
              <el-radio-group v-model="formData.loginFailLock" class="ml-4">
                <el-radio :label="true" size="large">是</el-radio>
                <el-radio :label="false" size="large">否</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="登录次数：" prop="loginTimes">
              <el-input-number v-model="formData.loginTimes" :min="1" type="number" />
            </el-form-item>
            <el-form-item label="自动解锁时间：" prop="unlockTime">
              <el-input-number v-model="formData.unlockTime" type="number" />
              <span class="unit">分钟</span>
            </el-form-item>
            <el-form-item label="超时配置：" prop="overtimeConfig">
              <el-radio-group v-model="formData.overtimeConfig" class="ml-4">
                <el-radio :label="true" size="large">是</el-radio>
                <el-radio :label="false" size="large">否</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="令牌失效时间：" prop="expireTime">
              <el-input-number v-model="formData.expireTime" :min="1000" type="number" />
              <span class="unit">分钟</span>
            </el-form-item>
            <el-form-item label="修改密码配置：" prop="changePwd">
              <el-radio-group v-model="formData.changePwd" class="ml-4">
                <el-radio :label="true" size="large">是</el-radio>
                <el-radio :label="false" size="large">否</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="密码更换周期：" prop="changeDays">
              <el-input-number v-model="formData.changeDays" :min="1" type="number" />
              <span class="unit">天</span>
            </el-form-item>
            <el-form-item label="限制IP登录：" prop="loginIpLimit">
              <el-radio-group v-model="formData.loginIpLimit" class="ml-4">
                <el-radio :label="true" size="large">是</el-radio>
                <el-radio :label="false" size="large">否</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="ip白名单：" prop="ipWhilteList">
              <el-input v-model="formData.ipWhilteList" />
            </el-form-item>
            <el-form-item>
              <!-- <el-button :disabled="true" type="primary">测试</el-button> -->
              <el-button type="primary" @click="submitForm(formRef)">保存</el-button>
            </el-form-item>
          </el-form>
        </el-col>
        <el-col :span="9" />
      </el-row>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .safe-config {
    .content {
      padding-bottom: 20px;
      margin-top: 20px;
      height: calc(100vh - 200px);
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
      .unit {
        margin-left: 10px;
        color: var(--el-text-color-placeholder);
      }
    }
  }
</style>
