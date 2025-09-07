<script lang="ts">
  export default {
    name: 'ModelConfig',
  }
</script>
<script setup lang="ts">
  import { ElInput, type FormInstance } from 'element-plus'
  import { setSystemConfigApi, mailTestSendApi, checkEnableApi } from '~/src/api-ecs/system'
  import AesEncryptCBC from '~/src/utils/crypto'
  import { requireRules } from '~/src/utils/rules'

  const props = defineProps<{
    allData: any
  }>()

  type objType = {
    [key: string]: any
  }

  // 表单数据
  const formData = reactive<objType>({
    address: 'YOUR LLM URL',
    key: '',
    isEnable: false,
  })

  const rules = reactive({
    address: requireRules,
    key: requireRules,
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
                    key: 'module_config',
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
    () => props.allData,
    () => {
      if (props.allData) {
        const data = JSON.parse(props.allData)
        for (const key in formData) {
          formData[key] = data[key]
          if (key == 'isEnable') formData[key] = data[key] || false
        }
      }
    },
    { immediate: true }
  )

  const handleChange = async (val: boolean | number | string) => {
    if (val) {
      const { msg, data } = await checkEnableApi({ key: 'module_config', value: JSON.stringify({ ...formData }) })
      if (data) return ElMessage({ message: msg, type: 'success' })
      formData.isEnable = false
      ElMessage({ message: msg, type: 'error' })
    }
  }
</script>

<template>
  <div class="mail-config">
    <div class="content">
      <el-row>
        <el-col :span="6" />
        <el-col :span="10">
          <el-form ref="formRef" label-position="right" label-width="140px" :model="formData" :rules="rules">
            <el-form-item label="LLM地址：" prop="address">
              <el-input v-model="formData.address" />
            </el-form-item>
            <el-form-item label="API_Key信息：" prop="key">
              <el-input v-model="formData.key" />
            </el-form-item>
            <el-form-item label="是否启用：" prop="isEnable">
              <div class="is-enable">
                <div>
                  <el-switch v-model="formData.isEnable" @change="(e) => handleChange(e)" />
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
