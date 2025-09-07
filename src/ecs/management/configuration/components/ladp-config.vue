<script lang="ts">
  export default {
    name: 'LadpConfig', // LDAP配置
  }
</script>

<script setup lang="ts">
  import type { FormInstance } from 'element-plus'

  import { setSystemConfigApi } from '~/src/api-ecs/system'

  import { requireRules } from '~/src/utils/rules'

  import sha1 from 'sha1'
  import AesEncryptCBC from '~/src/utils/crypto'

  const props = defineProps<{
    ladp: any
    id?: number
  }>()

  type objType = {
    [key: string]: any
  }
  // 表单数据
  const formData = reactive<objType>({
    enable: 1,
    serverType: 0,
    version: 0,
    serverIP: '',
    serverPort: '',
    account: '',
    password: '',
    domain: '',
    overtime: '',
    chase: 0,
    baseDN: '',
  })
  const rules = reactive({
    enable: requireRules,
    serverType: requireRules,
    version: requireRules,
    serverIP: requireRules,
    serverPort: requireRules,
    account: requireRules,
    password: requireRules,
    domain: requireRules,
    overtime: requireRules,
    chase: requireRules,
    baseDN: requireRules,
  })

  const formRef = ref<FormInstance>()

  const serverTypeOption = [
    { label: 'Microsoft AD', value: 0 },
    { label: 'OpenLADP', value: 1 },
  ]

  const versionOption = [{ label: '3.0', value: 0 }]

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
                    value: JSON.stringify({ ...formData, password: sha1(formData.password) }),
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
    () => props.ladp,
    () => {
      if (props.ladp) {
        const data = JSON.parse(props.ladp)
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
  <div class="ladp-config">
    <div class="content">
      <el-row>
        <el-col :span="5" />
        <el-col :span="12">
          <el-form ref="formRef" label-position="right" label-width="150px" :model="formData" :rules="rules">
            <el-form-item label="启用认证：" prop="enable">
              <el-radio-group v-model="formData.enable" class="ml-4">
                <el-radio :label="1" size="large">启用</el-radio>
                <el-radio :label="0" size="large">停用</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="服务器类型：" prop="serverType">
              <el-select v-model="formData.serverType" style="width: 100%">
                <el-option v-for="item in serverTypeOption" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="协议版本：" prop="version">
              <el-select v-model="formData.version" style="width: 100%">
                <el-option v-for="item in versionOption" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="服务器IP：" prop="serverIP">
              <el-input v-model="formData.serverIP" maxlength="39" show-word-limit />
            </el-form-item>
            <el-form-item label="服务器端口：" prop="serverPort">
              <el-input v-model="formData.serverPort" />
            </el-form-item>
            <el-form-item label="登录用户：" prop="account">
              <el-input v-model="formData.account" maxlength="32" show-word-limit />
            </el-form-item>
            <el-form-item label="登录密码：" prop="password">
              <el-input v-model="formData.password" type="password" />
            </el-form-item>
            <el-form-item label="域名：" prop="domain">
              <el-input v-model="formData.domain" maxlength="32" show-word-limit />
            </el-form-item>
            <el-form-item label="超时时间：" prop="overtime">
              <el-input v-model="formData.overtime">
                <template #append>毫秒</template>
              </el-input>
            </el-form-item>
            <el-form-item label="Chase referrals：" prop="chase">
              <el-switch v-model="formData.chase" />
            </el-form-item>
            <el-form-item label="Base DN：" prop="baseDN">
              <el-input v-model="formData.baseDN" />
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
  .ladp-config {
    .content {
      margin-top: 20px;
    }
  }
</style>
