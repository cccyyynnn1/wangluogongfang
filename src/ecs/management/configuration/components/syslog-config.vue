<script lang="ts">
  export default {
    name: 'SyslogConfig',
  }
</script>
<script setup lang="ts">
  import type { FormInstance } from 'element-plus'

  import { setSystemConfigApi } from '~/src/api-ecs/system'
  import AesEncryptCBC from '~/src/utils/crypto'

  import { requireRules } from '~/src/utils/rules'

  const props = defineProps<{
    syslog: any
    id?: number
  }>()

  type objType = {
    [key: string]: any
  }
  // 表单数据
  const formData = reactive<objType>({
    logType: 1,
    // logType: ['攻击日志'],
    addr: '',
    post: '',
    enable: true,
    isSave: false,
  })
  const rules = reactive({
    addr: requireRules,
    post: requireRules,
    // logType: requireRules,
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
                  { id: props.id as number, value: JSON.stringify(formData) },
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
    () => props.syslog,
    () => {
      if (props.syslog) {
        const data = JSON.parse(props.syslog)
        for (const key in formData) {
          formData[key] = data[key]
          // if (key == 'logType') formData[key] = data[key] || []
          if (key == 'dataFormat') formData[key] = data[key] || 1
        }
      }
    },
    { immediate: true }
  )
</script>

<template>
  <div class="syslog-config">
    <div class="content">
      <el-row>
        <el-col :span="5" />
        <el-col :span="12">
          <el-form ref="formRef" label-position="right" label-width="150px" :model="formData" :rules="rules">
            <el-form-item label="数据格式类型：" prop="logType">
              <el-radio-group v-model="formData.logType">
                <el-radio :label="0">字符串格式</el-radio>
                <el-radio :label="1">Json格式</el-radio>
              </el-radio-group>
            </el-form-item>
            <!-- <el-form-item label="日志类型：" prop="logType">
              <el-checkbox-group v-model="formData.logType">
                <el-checkbox label="攻击日志" />
                <el-checkbox label="访问日志" />
                <el-checkbox label="DDos日志" />
                <el-checkbox label="DDos日志" />
              </el-checkbox-group>
            </el-form-item> -->
            <el-form-item label="Syslog服务器地址：" prop="addr">
              <el-input v-model="formData.addr" />
            </el-form-item>
            <el-form-item label="Syslog服务器端口：" prop="post">
              <el-input v-model="formData.post" />
            </el-form-item>
            <el-form-item label="是否启用：" prop="enable">
              <el-switch v-model="formData.enable" />
            </el-form-item>
            <el-form-item label="本地储存：" prop="isSave">
              <el-switch v-model="formData.isSave" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submitForm(formRef)">保存</el-button>
            </el-form-item>
          </el-form>
        </el-col>
        <el-col :span="9" />
      </el-row>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .syslog-config {
    .content {
      margin-top: 20px;
    }
  }
</style>
