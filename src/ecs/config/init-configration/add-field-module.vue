<script lang="ts">
  export default {
    name: 'AddFieldModule',
  }
</script>

<script setup lang="ts">
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import { getComparisonFieldsTemplateApi } from '@/api-ecs/custom-field'
  import FieldModules from './field-modules.vue'
  import { tableSearch, InitialisationItem } from '@/types'
  const handleTableSearch = inject(tableSearch)
  const props = withDefaults(
    defineProps<{
      modelValue: boolean
    }>(),
    {
      modelValue: false,
    }
  )
  const emits = defineEmits<{
    (e: 'update:modelValue', visible: boolean): void
  }>()

  const templateData = ref()
  const displayStatus = ref('default')
  const fieldInfovVisible = ref(false)
  const templateId = ref()
  const queryForm = reactive({
    templateName: '',
    remark: '',
  })

  const templateList = ref<InitialisationItem[]>([])
  const dialogVisible = useVModel(props, 'modelValue', emits)
  const handleShowInfo = (display: 'preview' | 'add') => {
    templateData.value = Object.assign({ id: templateId.value }, display === 'add' ? queryForm : {})
    displayStatus.value = display
    fieldInfovVisible.value = true
    handleTableSearch?.()
  }
  const handleGetAllFieldsTemplate = async () => {
    const { data } = await getComparisonFieldsTemplateApi()
    templateList.value = data || []
    templateId.value = data[0] ? data[0]?.id : undefined
  }
  const handelReload = () => {
    dialogVisible.value = false
    handleTableSearch?.()
  }
  watch(
    () => dialogVisible.value,
    () => {
      if (dialogVisible.value) {
        queryForm.templateName = ''
        queryForm.remark = ''
        handleGetAllFieldsTemplate()
      }
    }
  )
  onMounted(() => handleGetAllFieldsTemplate())
</script>

<template>
  <vab-dialog v-model="dialogVisible" destroy-on-close title="新建模版" width="680px">
    <el-form
      label-position="top"
      label-width="100px"
      :model="queryForm"
      :rules="{
        templateName: [{ required: true, trigger: 'blur', message: '请输入模版名称' }],
      }"
      style="margin: 0 20px"
    >
      <el-form-item label="名称:" prop="templateName">
        <el-input v-model="queryForm.templateName" clearable />
      </el-form-item>
      <el-form-item label="备注" prop="remark">
        <el-input v-model="queryForm.remark" clearable />
      </el-form-item>
      <el-form-item label="引用模版基准" prop="template">
        <el-select v-model="templateId">
          <el-option
            v-for="template in templateList"
            :key="template.id"
            :label="template.templateName"
            :value="template.id"
          />
        </el-select>
        <span class="overview" @click="() => handleShowInfo('preview')">预览</span>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button :disabled="!queryForm.templateName || !templateId" type="primary" @click="() => handleShowInfo('add')">
        下一步
      </el-button>
      <el-button @click="dialogVisible = false">取消</el-button>
    </template>
    <field-modules
      v-model:visible="fieldInfovVisible"
      :display="displayStatus"
      :template-data="templateData"
      @close-handle="handelReload"
    />
  </vab-dialog>
</template>

<style scoped lang="scss">
  .overview {
    cursor: pointer;
    margin-left: 10px;
    color: var(--el-color-primary);
    font-weight: 400;
    font-size: 14px;
  }
</style>
