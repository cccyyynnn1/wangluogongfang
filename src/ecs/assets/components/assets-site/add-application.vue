<script lang="ts">
  export default {
    name: 'AddApplication',
  }
</script>
<script setup lang="ts">
  import { requireRules } from '@/utils/rules'

  import type { FormInstance } from 'element-plus'

  import { ApplicationType } from '@/types'

  import { ElMessage } from 'element-plus'

  import { siteAppSaveUpdateApi } from '~/src/api-ecs/assets'

  const props = defineProps<{
    showEditSite: boolean
    currentItem: ApplicationType
    mode: string
  }>()

  const emit = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
    (e: 'on-reflash'): void
  }>()

  const visible = ref(false)

  const formRef = ref<FormInstance>()

  const mode = ref('')

  const title = ref('')

  // 表单数据
  const formData = reactive<ApplicationType>({
    appName: '',
    appIp: '',
    siteId: undefined,
    id: undefined,
  })

  // 表单数据校验
  const rules = reactive({
    siteName: requireRules,
    appIp: requireRules,
  })

  // 提交
  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        const { msg } = await siteAppSaveUpdateApi({
          ...formData,
        })
        ElMessage({ message: msg, type: 'success' })
        emit('on-reflash')
        handleClose()
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  // 重置表单
  const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
  }

  // 关闭
  const handleClose = () => {
    resetForm(formRef.value)
    emit('on-closeEvent', false)
  }

  const initData = () => {
    visible.value = props.showEditSite
    mode.value = props.mode
    if (props.showEditSite) {
      formData.siteId = props.currentItem.siteId
      if (mode.value == 'edit') {
        title.value = '编辑'
        for (const key in formData) {
          // @ts-ignore
          formData[key as keyof typeof formData] = props.currentItem[key as keyof typeof formData]
        }
      } else {
        title.value = '新增'
      }
    }
  }

  onMounted(() => {
    initData()
  })
</script>

<template>
  <div class="edit-site">
    <el-dialog v-model="visible" :before-close="handleClose" :title="title" width="600px">
      <el-form ref="formRef" label-position="right" label-width="110px" :model="formData" :rules="rules">
        <el-form-item label="应用名称：" prop="appName">
          <el-input v-model="formData.appName" />
        </el-form-item>
        <el-form-item label="应用IP：" prop="appIp">
          <el-input v-model="formData.appIp" resize="none" :rows="4" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="submitForm(formRef)">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .space {
    width: 100%;
    min-height: 150px;
    padding: 10px;
    border: 1px solid var(--el-border-color);
    border-radius: 2px;
  }

  :deep(.el-tag) {
    margin-left: 10px;
  }
</style>
