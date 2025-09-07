<script lang="ts">
  export default {
    name: 'SSOConfig',
  }
</script>

<script setup lang="ts">
  import type { FormInstance } from 'element-plus'

  import { setSystemConfigApi } from '~/src/api-ecs/system'
  import AesEncryptCBC from '~/src/utils/crypto'

  import { requireRules } from '~/src/utils/rules'

  const props = defineProps<{
    sso: any
    id?: number
  }>()

  type objType = {
    [key: string]: any
  }
  // 表单数据
  const formData = reactive<objType>({
    enable: 1,
    authentication: 1,
    authenticationAddr: '',
    APPID: 0,
    APPSecret: '',
    // callback: '',
    // importType: 1,
    // LDPURL: '',
    accountDefaultRole: '',
  })
  const rules = reactive({
    enable: requireRules,
    authentication: requireRules,
    // callback: requireRules,
    // importType: requireRules,
    // LDPURL: requireRules,
  })

  const formRef = ref<FormInstance>()

  const encryptionOption = [{ label: 'pap', value: 0 }]

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
                const { msg } = await setSystemConfigApi(
                  {
                    id: props.id as number,
                    value: JSON.stringify(formData),
                  },
                  { password: AesEncryptCBC(password.value) }
                )
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

  watch(
    () => props.sso,
    () => {
      if (props.sso) {
        const data = JSON.parse(props.sso)
        for (const key in formData) {
          formData[key] = data[key]
          if (key == 'enable') {
            formData[key] = data[key] == 1 ? 1 : 0
          }
          // if (key == 'logType') formData[key] = data[key] || []
          // if (key == 'dataFormat') formData[key] = data[key] || 1
        }
      }
    },
    { immediate: true }
  )
</script>

<template>
  <div class="sso-config">
    <div class="content">
      <el-row>
        <el-col :span="5" />
        <el-col :span="12">
          <el-form ref="formRef" label-position="right" label-width="180px" :model="formData" :rules="rules">
            <el-form-item label="启用认证：" prop="enable">
              <el-radio-group v-model="formData.enable" class="ml-4">
                <el-radio :label="1" size="large">启用</el-radio>
                <el-radio :label="0" size="large">停用</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="认证方式：" prop="authentication">
              <el-radio-group v-model="formData.authentication" class="ml-4">
                <el-radio :label="1" size="large">OAUTH2.0</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="认证地址：" prop="authenticationAddr">
              <el-input v-model="formData.authenticationAddr" maxlength="128" show-word-limit />
            </el-form-item>
            <el-form-item label="应用ID：" prop="APPID">
              <el-input v-model="formData.APPID" maxlength="128" show-word-limit />
            </el-form-item>
            <el-form-item label="应用密钥：" prop="APPSecret">
              <el-input v-model="formData.APPSecret" maxlength="128" show-word-limit />
            </el-form-item>

            <!-- <el-form-item label="回调URL：" prop="callback">
              <el-input v-model="formData.callback" />
            </el-form-item>
            <el-form-item label="元数据导入方式：" prop="importType">
              <el-radio-group v-model="formData.importType" class="ml-4">
                <el-radio :label="0" size="large">导入URL</el-radio>
                <el-radio :label="1" size="large">自定义</el-radio>
                <el-radio :label="2" size="large">上传认证信息</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="LDP服务元数据URL：" prop="LDPURL">
              <el-input v-model="formData.LDPURL" />
            </el-form-item> -->
            <el-form-item label="用户默认角色：" prop="accountDefaultRole">
              <el-input v-model="formData.accountDefaultRole" />
            </el-form-item>
            <el-form-item>
              <el-button :disabled="true" type="primary">测试</el-button>
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
  .sso-config {
    .content {
      margin-top: 20px;
    }
  }
</style>
