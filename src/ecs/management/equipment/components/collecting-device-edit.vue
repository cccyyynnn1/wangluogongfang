<script setup lang="ts">
  import { FormInstance } from 'element-plus'

  import { requireRules, requireAndNumberRules } from '~/src/utils/rules'

  import { updateEquipmentApi } from '@/api-ecs/equipment'
  import AesEncryptCBC from '~/src/utils/crypto'

  const props = defineProps<{
    showPage: boolean
    currentRow: object
  }>()

  const formRef = ref<FormInstance>() // 表单实例

  const visible = ref(false) // 弹框显隐

  // 表单数据
  const formData = reactive({
    id: '',
    dataTransactionSecond: 5, // 数据传输合并间隔
    dbMemoryLimitGb: 50, // 数据库内存限制
    eventStatSecond: 10, // 事件统计间隔
    isNuma: true, // 是否开启numa
    needRestartService: false, // 是否需要重启服务
    systemTime: true, // 是否同步系统时间
    zipStorage: true, // 启⽤压缩存储
    logLevelStr: 'DEBUG', // ⽇志级别
    dialogTitle: 'edit',
    moduleType: 'flowDevice',
  })

  // 表单数据校验规则
  const rules = reactive({
    dataTransactionSecond: requireAndNumberRules,
    dbMemoryLimitGb: requireAndNumberRules,
    eventStatSecond: requireAndNumberRules,
    logLevelStr: requireRules,
    zipStorage: requireRules,
  })

  const emit = defineEmits<{
    (e: 'on-close-event'): void
    (e: 'on-reflash'): void
  }>()

  // 关闭回调
  const handleClose = () => {
    emit('on-close-event')
  }

  // 提交表单
  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
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
                const { msg } = await updateEquipmentApi(
                  formData.id,
                  { ...formData },
                  { password: AesEncryptCBC(password.value) }
                )
                emit('on-reflash')
                ElMessage({ message: msg, type: 'success' })
                handleClose()
                done()
              } catch (error) {
                instance.confirmButtonLoading = false
              }
            } else {
              done()
            }
          },
        }).catch(() => {})
        // @ts-ignore
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  // 初始化数据
  const initData = () => {
    visible.value = props.showPage
    const data = JSON.parse(JSON.stringify(props.currentRow))
    Object.keys(formData).forEach((item: string) => {
      if (item !== 'dialogTitle' && item !== 'moduleType') {
        // @ts-ignore
        formData[item] = data[item]
      }
    })
  }

  onMounted(() => {
    initData()
  })
</script>

<script lang="ts">
  export default {
    name: 'CollectingDeviceEdit',
  }
</script>
<template>
  <div class="collecting-device-edit">
    <el-dialog v-model="visible" :before-close="handleClose" title="编辑" width="600px">
      <el-form ref="formRef" label-width="150px" :model="formData" :rules="rules">
        <el-form-item label="启用压缩存储" prop="zipStorage">
          <el-radio v-model="formData.zipStorage" :label="true">启用</el-radio>
          <el-radio v-model="formData.zipStorage" :label="false">停用</el-radio>
        </el-form-item>
        <el-form-item label="日志级别" prop="logLevelStr">
          <el-input v-model="formData.logLevelStr" :disabled="true" />
        </el-form-item>
        <el-form-item label="事件统计间隔" prop="eventStatSecond">
          <el-input v-model="formData.eventStatSecond">
            <template #suffix>秒</template>
          </el-input>
        </el-form-item>
        <el-form-item label="数据传输合并间隔" prop="dataTransactionSecond">
          <el-input v-model="formData.dataTransactionSecond">
            <template #suffix>秒</template>
          </el-input>
        </el-form-item>
        <el-form-item label="数据库内存限制" prop="dbMemoryLimitGb">
          <el-input v-model="formData.dbMemoryLimitGb">
            <template #suffix>GB</template>
          </el-input>
        </el-form-item>
        <el-form-item label="是否同步系统时间" prop="systemTime">
          <el-radio v-model="formData.systemTime" :label="true">是</el-radio>
          <el-radio v-model="formData.systemTime" :label="false">否</el-radio>
        </el-form-item>
        <el-form-item label="是否开启numa" prop="isNuma">
          <el-radio v-model="formData.isNuma" :label="true">是</el-radio>
          <el-radio v-model="formData.isNuma" :label="false">否</el-radio>
        </el-form-item>
        <el-row>
          <el-col :offset="10" :span="4" style="text-align: right">
            <el-button type="primary" @click="submitForm(formRef)">保存</el-button>
          </el-col>
        </el-row>
      </el-form>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss"></style>
