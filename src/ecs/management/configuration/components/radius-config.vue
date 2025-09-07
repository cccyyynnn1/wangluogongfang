<script lang="ts">
  export default {
    name: 'RadiusConfig', // radius配置
  }
</script>
<script setup lang="ts">
  import type { FormInstance } from 'element-plus'

  import { setSystemConfigApi } from '~/src/api-ecs/system'
  import AesEncryptCBC from '~/src/utils/crypto'

  import { requireRules } from '~/src/utils/rules'

  const props = defineProps<{
    radius: any
    id?: number
  }>()

  type objType = {
    [key: string]: any
  }
  // 表单数据
  const formData = reactive<objType>({
    enable: 1,
    serverIP: '',
    authenticationPort: '',
    billingPort: '',
    secret: '',
    encryption: 0,
    repetition: '',
    overtime: '',
  })
  const rules = reactive({
    enable: requireRules,
    serverIP: requireRules,
    authenticationPort: requireRules,
    billingPort: requireRules,
    secret: requireRules,
    encryption: requireRules,
    repetition: requireRules,
    overtime: requireRules,
  })

  const formRef = ref<FormInstance>()

  const encryptionOption = [
    { label: 'pap', value: 0 },
    { label: 'chap', value: 1 },
  ]

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
    () => props.radius,
    () => {
      if (props.radius) {
        const data = JSON.parse(props.radius)
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
  <div class="radius-config">
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
            <el-form-item label="服务器IP：" prop="serverIP">
              <el-input v-model="formData.serverIP" maxlength="39" show-word-limit />
            </el-form-item>
            <el-form-item label="认证端口：" prop="authenticationPort">
              <el-input v-model="formData.authenticationPort" />
            </el-form-item>
            <el-form-item label="计费端口：" prop="billingPort">
              <el-input v-model="formData.billingPort" />
            </el-form-item>
            <el-form-item label="密钥：" prop="secret">
              <el-input v-model="formData.secret" maxlength="32" show-word-limit />
            </el-form-item>
            <el-form-item label="加密方式：" prop="encryption">
              <el-select v-model="formData.encryption" style="width: 100%">
                <el-option v-for="item in encryptionOption" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="重复次数：" prop="repetition">
              <el-input v-model="formData.repetition" />
            </el-form-item>
            <el-form-item label="超时时间：" prop="overtime">
              <el-input v-model="formData.overtime">
                <template #append>毫秒</template>
              </el-input>
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
  .radius-config {
    .content {
      margin-top: 20px;
    }
  }
</style>
