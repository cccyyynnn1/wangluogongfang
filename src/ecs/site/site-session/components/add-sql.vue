<script setup lang="ts">
  import type { FormInstance } from 'element-plus'

  import { ShortcutListType } from '~/src/types'

  import { uuid } from '~/src/utils'

  import { requireRules } from '~/src/utils/rules'

  const props = defineProps<{
    mode: string
    show: boolean
    tableColumn: any
    currentRow?: ShortcutListType
  }>()

  // 表单数据
  const formData = reactive<ShortcutListType>({
    key: '',
    relation: '=',
    value: '',
    label: '',
    enable: true,
    id: '',
  })

  const modes = ref('')

  const lock = ref(false)

  const emit = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
    (e: 'handleUpdateCallback', formData: ShortcutListType): void
  }>()

  const visible = ref(false) // 显隐

  const formRef = ref<FormInstance>() // 表单实例

  const flag = ref()

  let keyOptions = reactive<{ fieldNameCn: string; fieldNameEn: string; fieldType: string; dict: string[] }[]>([])

  // ip
  const relationOptions = reactive<{ value: string; label: string }[]>([
    { value: '=', label: '= 等于' },
    { value: '!=', label: '!= 不等于' },
    { value: '>', label: '> 大于' },
    { value: '<', label: '< 小于' },
    { value: '>=', label: '>= 大于等于' },
    { value: '<=', label: '<= 小于等于' },
    { value: 'exists', label: 'exists 存在' },
    { value: 'not_exists', label: 'not_exists 不存在' },
  ])

  // 平均数
  const relationOptions1 = reactive<{ value: string; label: string }[]>([
    { value: '>', label: '> 大于' },
    { value: '<', label: '< 小于' },
  ])

  // 方向
  const relationOptions2 = reactive<{ value: string; label: string }[]>([{ value: '=', label: '= 等于' }])

  // isn_t
  const relationOptions3 = reactive<{ value: string; label: string }[]>([
    { value: '=', label: '= 等于' },
    { value: '!=', label: '!= 不等于' },
  ])

  // 字符串
  const relationOptions4 = reactive<{ value: string; label: string }[]>([
    { value: '=', label: '= 等于' },
    { value: '!=', label: '!= 不等于' },
    { value: 'like', label: 'like 字符串匹配' },
    { value: 'not_like', label: 'not_like 字符串匹配' },
    { value: 'exists', label: 'exists 存在' },
    { value: 'not_exists', label: 'not_exists 不存在' },
  ])

  const selectOption = ref<string[]>([])

  const handleClose = () => {
    emit('on-closeEvent', false)
  }

  const rules = reactive({
    key: requireRules,
    relation: requireRules,
    value: requireRules,
  })

  const isSelect = ref(false)

  // 回显数据
  const initData = () => {
    visible.value = props.show
    modes.value = props.mode
    keyOptions = props.tableColumn.filter((item: any) => {
      return !item.selectHidden && item.fieldType != 'no_serach' && item.fieldType != 'sum'
    })
    formData.key = keyOptions[0]?.fieldNameEn
    if (modes.value == 'edit') {
      Object.keys(formData).forEach((item: any) => {
        // @ts-ignore
        formData[item] = props.currentRow[item]
      })
    } else {
      lock.value = true
    }
  }

  // 提交
  const submit = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        if (modes.value == 'add') {
          formData.id = uuid()
        }
        const res = keyOptions.find((item: any) => {
          return item.fieldNameEn == formData.key
        })
        if (formData.relation == 'not_exists' || formData.relation == 'exists') {
          formData.value = ''
        }
        formData.label = res!.fieldNameCn
        emit('handleUpdateCallback', formData)
        ElMessage({ message: '操作成功', type: 'success' })
        handleClose()
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  const showRelation = ref(false)

  watch(
    () => formData.relation,
    (newVal) => {
      if (newVal == 'not_exists' || newVal == 'exists') {
        showRelation.value = false
      } else {
        showRelation.value = true
      }
    },
    {
      immediate: true,
    }
  )

  watch(
    () => formData.key,
    (newVal) => {
      const result = keyOptions.find((item: any) => {
        return item.fieldNameEn === newVal
      })
      flag.value = result?.fieldType
      isSelect.value = !!result?.dict.length
      selectOption.value = result?.dict || []
      if (lock.value) {
        switch (flag.value) {
          case 'direction':
            formData.relation = '='
            formData.value = '内对内'
            break
          case 'average':
            formData.relation = '>'
            break
          case 'text':
            formData.relation = '='
            break
          default:
            formData.relation = '='
            break
        }
      }
      lock.value = true
    }
  )

  onMounted(() => {
    initData()
  })
</script>

<script lang="ts">
  export default {
    name: 'AddSql',
  }
</script>
<template>
  <div class="add-sql">
    <el-dialog v-model="visible" :before-close="handleClose" :title="modes == 'add' ? '添加' : '编辑'" width="700px">
      <el-form ref="formRef" label-position="right" label-width="0px" :model="formData" :rules="rules">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item prop="key">
              <el-select v-model="formData.key">
                <el-option
                  v-for="item in keyOptions"
                  :key="item.fieldNameEn"
                  :label="item.fieldNameCn"
                  :value="item.fieldNameEn"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item prop="relation">
              <el-select v-if="flag == 'ip' || flag == 'num'" v-model="formData.relation">
                <el-option v-for="item in relationOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
              <el-select v-else-if="flag == 'average'" v-model="formData.relation">
                <el-option v-for="item in relationOptions1" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
              <el-select v-else-if="flag == 'direction'" v-model="formData.relation">
                <el-option v-for="item in relationOptions2" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
              <el-select v-else-if="flag == 'isn_t'" v-model="formData.relation">
                <el-option v-for="item in relationOptions3" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
              <el-select v-else-if="flag == 'text'" v-model="formData.relation">
                <el-option v-for="item in relationOptions4" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
              <el-select v-else v-model="formData.relation">
                <el-option v-for="item in relationOptions4" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col v-if="showRelation" :span="8">
            <el-form-item v-if="!isSelect" prop="value">
              <el-input v-model="formData.value" />
            </el-form-item>
            <el-form-item v-else prop="value">
              <el-select v-model="formData.value">
                <el-option v-for="item in selectOption" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="submit(formRef)">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss"></style>
