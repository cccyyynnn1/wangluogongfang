<script setup lang="ts">
  import { FormInstance } from 'element-plus'

  import { addChainSsortApi } from '~/src/api-ecs/assets'

  import { requireRules } from '~/src/utils/rules'

  const props = defineProps<{
    showPage: boolean
  }>()

  const formRef = ref<FormInstance>() // 表单实例

  const option = [
    { value: 'https://', label: 'https://' },
    { value: 'http://', label: 'http://' },
  ]

  const visible = ref(false) // 弹框显隐

  // 表单数据
  const formData = reactive({
    name: 'http://',
    assetsIp: undefined,
  })

  // 表单数据校验规则
  const rules = reactive({
    assetsIp: requireRules,
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
        const { msg } = await addChainSsortApi({ site: formData.name + formData.assetsIp })
        emit('on-reflash')
        ElMessage({ message: msg, type: 'success' })
        handleClose()
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  onMounted(() => {
    visible.value = props.showPage
  })
</script>

<script lang="ts">
  export default {
    name: 'AddServiceChain',
  }
</script>
<template>
  <div class="add-service-chain">
    <el-dialog v-model="visible" :before-close="handleClose" title="添加任务" width="680px">
      <el-form ref="formRef" :model="formData" :rules="rules">
        <el-row :gutter="20">
          <el-col :span="7">
            <el-form-item label="站点域名" prop="name">
              <el-select v-model="formData.name" style="width: 100%">
                <el-option v-for="item in option" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="17">
            <el-form-item prop="assetsIp">
              <el-input v-model="formData.assetsIp" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col style="text-align: right">
            <el-button type="primary" @click="submitForm(formRef)">确认</el-button>
            <el-button @click="handleClose">取消</el-button>
          </el-col>
        </el-row>
      </el-form>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss"></style>
