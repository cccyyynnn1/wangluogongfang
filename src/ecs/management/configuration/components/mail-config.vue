<script lang="ts">
  export default {
    name: 'MailConfig',
  }
</script>
<script setup lang="ts">
  import type { FormInstance } from 'element-plus'
  import { setSystemConfigApi, mailTestSendApi } from '~/src/api-ecs/system'
  import { requireRules } from '~/src/utils/rules'
  import sha1 from 'sha1'
  import AesEncryptCBC from '~/src/utils/crypto'

  const props = defineProps<{
    mail: any
    id?: number
  }>()

  type objType = {
    [key: string]: any
  }

  // 表单数据
  const formData = reactive<objType>({
    address: '',
    port: '',
    account: '',
    password: '',
    isEnable: true,
  })

  const recipients = ref()

  const dialogVisible = ref(false)

  const sendMailEvent = async () => {
    if (!recipients.value) return
    const { msg } = await mailTestSendApi({ mailAddr: recipients.value })
    ElMessage({ message: msg, type: 'success' })
    dialogVisible.value = false
  }

  const sendMail = () => {
    dialogVisible.value = true
  }

  const rules = reactive({
    address: requireRules,
    port: requireRules,
    account: requireRules,
    password: requireRules,
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
                const { msg } = await setSystemConfigApi(
                  {
                    id: props.id as number,
                    value: JSON.stringify({ ...formData }),
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
    () => props.mail,
    () => {
      if (props.mail) {
        const data = JSON.parse(props.mail)
        for (const key in formData) {
          formData[key] = data[key]
          if (key == 'isEnable') formData[key] = data[key] || false
        }
      }
    },
    { immediate: true }
  )
</script>

<template>
  <div class="mail-config">
    <div class="content">
      <el-row>
        <el-col :span="6" />
        <el-col :span="10">
          <el-form ref="formRef" label-position="right" label-width="120px" :model="formData" :rules="rules">
            <el-form-item label="服务器地址：" prop="address">
              <el-input v-model="formData.address" />
            </el-form-item>
            <el-form-item label="端口：" prop="port">
              <el-input v-model="formData.port" />
            </el-form-item>
            <el-form-item label="发送人账户：" prop="account">
              <el-input v-model="formData.account" />
            </el-form-item>
            <el-form-item label="密码：" prop="password">
              <el-input v-model="formData.password" />
            </el-form-item>
            <el-form-item label="是否启用：" prop="isEnable">
              <div class="is-enable">
                <div>
                  <el-switch v-model="formData.isEnable" />
                </div>
                <div class="btn">
                  <el-button link type="primary" @click="sendMail">发送测试邮件</el-button>
                </div>
              </div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submitForm(formRef)">保存</el-button>
            </el-form-item>
          </el-form>
        </el-col>
        <el-col :span="8" />
      </el-row>
    </div>
    <el-dialog v-model="dialogVisible" title="添加收件人" width="30%">
      <el-form ref="formRef" label-position="right" label-width="100px" :model="formData">
        <el-form-item label="收件人：">
          <el-input v-model="recipients" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="sendMailEvent">确认</el-button>
          <el-button @click="dialogVisible = false">取消</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style lang="scss" scoped>
  .mail-config {
    .content {
      margin-top: 20px;

      .is-enable {
        display: flex;
        align-items: center;
        justify-content: space-between;
        width: 100%;
      }
    }
  }
</style>
