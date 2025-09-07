<script lang="ts">
  export default {
    name: 'ExtractionRuleEdit',
  }
</script>
<script setup lang="ts">
  import { normalizeSaveOrUpdateType } from '@/types'

  import { requireRules } from '~/src/utils/rules'

  import { useTableCopy } from '@/utils'

  import {
    normalizeGetRowByIdApi,
    normalizeParseLogApi,
    normalizeSaveOrUpdateApi,
    normalizeGroupGetAllApi,
  } from '@/api-ecs/normalize'

  import { customFieldGetByTypeApi } from '@/api-ecs/custom-field'

  import { FormInstance } from 'element-plus'

  const props = defineProps<{
    id: string | number
    mode: string
    currentRow: object
    groupOptions: []
  }>()

  const emit = defineEmits<{
    (e: 'on-save-event', val: boolean): void
    (e: 'on-reflash'): void
  }>()

  // 表格数据
  // let analysFields = ref<object[]>([])

  const mode = ref('')

  const activeName = ref('rule') // tabs选中项

  const listLoading = ref(false) // 是否加载

  const btnDisabled = ref(false) // 按钮置灰

  const newFieldsList = ref() // 新增字段

  // 表单数据
  const formData = reactive<normalizeSaveOrUpdateType>({
    analysisFieldsStr: '',
    analysisType: '0',
    createTime: 0,
    createType: '',
    description: '',
    field: '0',
    groupId: undefined,
    id: undefined,
    indexType: '',
    logSample: '',
    name: '',
    requestHost: '',
    requestUrl: '',
    status: 1,
    updateTime: 0,
    groupSplit: '',
    kvSplit: '',
    analysFields: [],
    newFieldsStr: undefined,
  })

  // 表单校验
  const rules = reactive({
    name: requireRules,
    logSample: requireRules,
    status: requireRules,
    groupId: requireRules,
    requestUrl: requireRules,
    field: requireRules,
  })

  const groupOptions = ref() // 策略分组

  const formRef = ref<FormInstance>() // 表单实例

  // 字段分组
  const fieldOptions = ref([
    { label: '请求负载', value: '0' },
    { label: '响应负载', value: '1' },
    { label: 'url', value: '2' },
    { label: '请求cookie', value: '3' },
  ])

  // 映射字段配置
  const mapFieldOptions = ref([
    {
      fieldNameCn: '',
      fieldNameEn: '',
      disabled: false,
    },
  ])

  //
  const encodeOptions = [
    { value: 'URL', label: 'URL解码' },
    { value: 'base64', label: 'base64解码' },
    { value: 'unicode', label: 'unicode转中文' },
  ]

  onMounted(() => {
    initData()
  })

  // 得到所有分组
  const getAllNormalizeGroup = async () => {
    const { data } = await normalizeGroupGetAllApi()
    groupOptions.value = data
  }

  // 初始化数据
  const initData = async () => {
    mode.value = props.mode
    if (mode.value == 'edit') {
      groupOptions.value = props.groupOptions
      const { data } = await normalizeGetRowByIdApi({ id: props.id })
      Object.keys(data).forEach((item) => {
        if (item == 'status') {
          // @ts-ignore
          formData[item] = Number(data[item])
        } else {
          // @ts-ignore
          formData[item] = data[item]
        }
      })
    } else {
      await getAllNormalizeGroup()
      formData.logSample = JSON.stringify(props.currentRow)
      formData.groupId = groupOptions.value[0].id
      formData.kvSplit = '='
      formData.groupSplit = '&'
      extractionField()
    }
    getMapFidldsOption()
  }

  // 字段提取
  const extractionField = async () => {
    btnDisabled.value = true
    try {
      const { data, msg } = await normalizeParseLogApi({
        logSample: formData.logSample,
        field: formData.field,
        kvSplit: formData.kvSplit,
        groupSplit: formData.groupSplit,
      })
      ElMessage({ message: msg, type: 'success' })
      if (mode.value == 'add') {
        Object.keys(data).forEach((item: string) => {
          // @ts-ignore
          if (item == 'status' && item != 'logSample' && item != 'field' && item != 'groupId') {
            // @ts-ignore
            formData[item] = Number(data[item])
          } else if (item != 'logSample' && item != 'field' && item != 'groupId') {
            // @ts-ignore
            formData[item] = data[item]
          }
        })
      }
      formData.analysFields = await data['analysFields']
      formData.analysisType = await data['analysisType']
      btnDisabled.value = false
      emit('on-reflash')
    } finally {
      btnDisabled.value = false
    }
  }

  // 获取字段映射
  const getMapFidldsOption = async () => {
    const { data } = await customFieldGetByTypeApi({ type: '1' })
    mapFieldOptions.value = data
  }

  // 保存
  const save = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        toCheckNewFields()
        if (newFieldsList.value.length > 0) {
          formData.newFieldsStr = JSON.stringify(newFieldsList.value)
        }
        const arr: object[] = []
        formData.analysFields.forEach((item: any) => {
          if (item.mappingField) {
            const flag = newFieldsList.value.every((e: any) => {
              return e.cnName !== item.mappingField
            })
            if (flag) {
              arr.push(item)
            }
          }
        })
        const { msg } = await normalizeSaveOrUpdateApi({
          ...formData,
          analysFields: undefined,
          analysisFieldsStr: JSON.stringify(arr),
        })
        ElMessage({ message: msg, type: 'success' })
        emit('on-save-event', false)
        emit('on-reflash')
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  // 互斥
  const handleChange = () => {
    const arr = [] as string[]
    formData.analysFields.forEach((item: any) => {
      if (item.mappingField) {
        arr.push(item.mappingField)
      }
    })
    mapFieldOptions.value.forEach((td: any) => {
      td['disabled'] = false
      if (arr.includes(td.fieldNameEn)) {
        td['disabled'] = true
      }
    })
  }

  // 新增字段
  const toCheckNewFields = () => {
    newFieldsList.value = []
    formData.analysFields.forEach((item: any) => {
      const res = mapFieldOptions.value.some((td: any) => {
        return item.mappingField == td.fieldNameEn
      })
      if (!res && item.mappingField) {
        if (item.mappingField) {
          const obj = {
            cnName: '',
            sourceField: '',
          }
          obj.cnName = item.mappingField
          obj.sourceField = item.sourceField
          newFieldsList.value.push(obj)
        }
      }
    })
  }

  // 返回
  const back = () => {
    emit('on-save-event', false)
  }
</script>

<template>
  <div class="extraction-rule-edit-container">
    <el-tabs v-model="activeName">
      <el-tab-pane label="提取规则" name="rule">
        <el-button type="primary" @click="back">返回</el-button>
        <el-row style="margin-top: 20px">
          <el-col :span="2" />
          <el-col :span="17">
            <el-form
              ref="formRef"
              class="login-form"
              label-position="right"
              label-width="140px"
              :model="formData"
              :rules="rules"
            >
              <el-form-item label="日志样本" prop="logSample">
                <el-input v-model="formData.logSample" resize="none" :rows="4" type="textarea" />
              </el-form-item>
              <el-row>
                <el-col :span="12">
                  <el-form-item label="策略名称" prop="name">
                    <el-input v-model="formData.name" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="状态" prop="status">
                    <el-radio-group v-model="formData.status">
                      <el-radio :label="1">启用</el-radio>
                      <el-radio :label="0">关闭</el-radio>
                    </el-radio-group>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="策略分组" prop="groupId">
                    <el-select v-model="formData.groupId" class="m-2" style="width: 100%">
                      <el-option v-for="item in groupOptions" :key="item.value" :label="item.name" :value="item.id" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="URL" prop="requestUrl">
                    <el-input v-model="formData.requestUrl" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Host" prop="requestHost">
                    <el-input v-model="formData.requestHost" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="策略描述" prop="description">
                    <el-input v-model="formData.description" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="提取字段" prop="field">
                    <el-select v-model="formData.field" class="m-2" style="width: 100%">
                      <el-option
                        v-for="item in fieldOptions"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                      />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="组分隔符" prop="groupSplit">
                    <el-input v-model="formData.groupSplit" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="键值对分隔符" prop="kvSplit">
                    <el-input v-model="formData.kvSplit" />
                  </el-form-item>
                </el-col>
              </el-row>

              <el-form-item label="">
                <el-button :disabled="btnDisabled" type="primary" @click="extractionField">字段提取</el-button>
              </el-form-item>
              <el-form-item label="提取结果" prop="analysFields">
                <el-table
                  v-loading="listLoading"
                  :border="true"
                  :data="formData.analysFields"
                  @cell-contextmenu="useTableCopy"
                >
                  <el-table-column :align="'center'" label="字段" prop="sourceField" show-overflow-tooltip />
                  <el-table-column :align="'center'" label="映射字段" prop="mappingField" show-overflow-tooltip>
                    <template #default="{ row }">
                      <el-select
                        v-model="row.mappingField"
                        :allow-create="true"
                        class="m-2"
                        :clearable="true"
                        filterable
                        style="width: 100%"
                        @change="handleChange"
                      >
                        <el-option
                          v-for="item in mapFieldOptions"
                          :key="item.fieldNameEn"
                          :disabled="item.disabled"
                          :label="item.fieldNameCn"
                          placeholder="请选择"
                          :value="item.fieldNameEn"
                        />
                      </el-select>
                    </template>
                  </el-table-column>
                  <el-table-column :align="'center'" label="提取值" prop="sourceValue" show-overflow-tooltip />
                  <el-table-column :align="'center'" label="解码" prop="encode" show-overflow-tooltip>
                    <template #default="{ row }">
                      <el-select
                        v-model="row.encode"
                        :allow-create="true"
                        class="m-2"
                        :clearable="true"
                        filterable
                        style="width: 100%"
                      >
                        <el-option
                          v-for="item in encodeOptions"
                          :key="item.value"
                          :label="item.label"
                          placeholder="请选择"
                          :value="item.value"
                        />
                      </el-select>
                    </template>
                  </el-table-column>
                  <template #empty>
                    <el-empty class="vab-data-empty" description="暂无数据" />
                  </template>
                </el-table>
              </el-form-item>
              <el-form-item label="">
                <el-button type="primary" @click="save(formRef)">保存</el-button>
              </el-form-item>
            </el-form>
          </el-col>
          <el-col :span="5" />
        </el-row>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped lang="scss">
  .extraction-rule-edit-container {
    .btn {
      width: 100%;
      display: flex;
      justify-content: center;
      margin-bottom: 20px;
    }
  }
</style>
