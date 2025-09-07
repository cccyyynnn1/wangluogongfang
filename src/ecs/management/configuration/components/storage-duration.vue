<script lang="ts">
  export default {
    name: 'StorageDuration',
  }
</script>
<script setup lang="ts">
  import type { FormInstance } from 'element-plus'

  import { setSystemConfigApi } from '~/src/api-ecs/system'
  import AesEncryptCBC from '~/src/utils/crypto'

  import { requireRules } from '~/src/utils/rules'

  const props = defineProps<{
    ecsDataTime: any
    id?: number
  }>()

  type objType = {
    [key: string]: any
  }
  // 表单数据
  const formData = reactive<objType>({
    duration: 30,
  })
  const rules = reactive({
    duration: requireRules,
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
    () => props.ecsDataTime,
    () => {
      if (props.ecsDataTime) {
        const data = JSON.parse(props.ecsDataTime)
        for (const key in formData) {
          formData[key] = data[key]
          // // if (key == 'logType') formData[key] = data[key] || []
          // if (key == 'dataFormat') formData[key] = data[key] || 1
        }
      }
    },
    { immediate: true }
  )
</script>

<template>
  <div class="ecs-data-time">
    <div class="content">
      <el-row>
        <el-col :span="5" />
        <el-col :span="12">
          <el-form ref="formRef" label-position="right" label-width="150px" :model="formData" :rules="rules">
            <el-form-item label="数据存储天数：" prop="duration">
              <el-input-number v-model="formData.duration" :min="180" />
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
  .ecs-data-time {
    .content {
      margin-top: 20px;
    }
  }
</style>
