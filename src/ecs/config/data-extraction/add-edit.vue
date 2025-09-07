<script lang="ts">
  export default {
    name: 'AddEdit',
  }
</script>
<script setup lang="ts">
  import { requireRules } from '~/src/utils/rules'
  import type { FormInstance } from 'element-plus'

  import { customFieldSaveOrUpdateApi } from '~/src/api-ecs/custom-field'

  import { customFieldSaveOrUpdateType, IndexTypeTpye } from '@/types'

  import { useUserStore } from '@/store/modules/user'

  const { getAllIndexType } = useUserStore()

  const tagType = ref<IndexTypeTpye[]>(getAllIndexType())

  const props = defineProps<{
    showDrawer: boolean
    title: string
    currentData?: customFieldSaveOrUpdateType
  }>()

  const emit = defineEmits<{
    (e: 'on-closeEvent'): void
    (e: 'on-reflash'): void
  }>()

  const visible = ref(false)

  const formRef = ref<FormInstance>() // 表单实例

  const title = ref('')

  // 表单数据
  const formData = reactive<customFieldSaveOrUpdateType>({
    fieldNameCn: '',
    type: 1,
    isDisplay: '0',
    id: null,
  })

  const currentId = ref(0)

  // 表单数据校验规则
  const rules = reactive({
    fieldNameCn: requireRules,
    type: requireRules,
  })

  onMounted(() => {
    initData()
  })

  // 回显数据
  const initData = () => {
    visible.value = props.showDrawer
    title.value = props.title
    if (title.value == '编辑') {
      const data = JSON.parse(JSON.stringify(props.currentData))
      currentId.value = data.id
      Object.keys(formData).forEach((item) => {
        if (item == 'type') {
          formData[item] = Number(data[item])
        } else if (item == 'isDisplay') {
          // @ts-ignore
          formData[item] = data[item] == 1 ? true : false
        } else {
          // @ts-ignore
          formData[item] = data[item]
        }
      })
    }
    // console.log(formData, 'formData')
  }

  // 提交
  const submit = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        const { msg } = await customFieldSaveOrUpdateApi({ ...formData })
        ElMessage({ message: msg, type: 'success' })
        emit('on-reflash')
        handleClose()
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  const handleClose = () => {
    emit('on-closeEvent')
  }
</script>

<template>
  <div class="add-edit">
    <el-dialog v-model="visible" :before-close="handleClose" :title="title" width="550px">
      <el-form ref="formRef" label-position="right" label-width="100px" :model="formData" :rules="rules">
        <el-form-item label="字段名称" prop="fieldNameCn">
          <el-input v-model="formData.fieldNameCn" maxlength="14" show-word-limit />
        </el-form-item>
        <el-form-item label="索引类型" prop="type">
          <el-select v-model="formData.type" class="m-2" disabled style="width: 100%">
            <el-option v-for="item in tagType" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="submit(formRef)">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .add-edit {
    :deep() {
      .el-descriptions__content {
        width: 350px;
      }
    }
  }
</style>
