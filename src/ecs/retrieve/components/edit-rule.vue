<script lang="ts">
  export default {
    name: 'EditRule', // 白名单
  }
</script>

<script setup lang="ts">
  import { onMounted, ref } from 'vue'

  import type { FormInstance } from 'element-plus'

  import { requireRules } from '~/src/utils/rules'

  import { GetConfigFieldsApi, WarnWhiteSaveOrUpdatesApi } from '~/src/api-ecs/alert'

  import { useUserStore } from '@/store/modules/user'

  import { IndexTypeTpye, LevelRuleSaveUpdateModel } from '~/src/types'

  import { LevelRuleSaveUpdateAPI } from '~/src/api-ecs/retrieve'

  const userStore = useUserStore()

  const { getTableColumn, getAllIndexType, username } = userStore

  const tagType = ref<IndexTypeTpye[]>(getAllIndexType())

  const $baseConfirm: any = inject('$baseConfirm')

  const props = defineProps<{
    isShow: boolean
    curData: any
    remark: string
  }>()

  const visible = ref(false) // 显隐

  const isLoading = ref(false) // 加载

  const mode = ref('')

  const title = ref('添加规则')

  const levelOption = [
    { label: '普通', value: '普通' },
    { label: '一般', value: '一般' },
    { label: '严重', value: '严重' },
  ]

  const relatOption = [
    { label: '!=', value: '!=' },
    { label: '=', value: '=' },
  ]

  // 表单数据
  const formData = reactive<LevelRuleSaveUpdateModel>({
    ruleName: '',
    cnd: '',
    id: undefined,
    indexType: 1,
    level: '普通',
  })

  const cndList = ref()

  const fieldRules = reactive<{ [key: string]: any }>({})

  let fieldList1: { fieldNameCn: string; fieldNameEn: string; disabled?: boolean }[] = []

  let fieldData3 = reactive<{ [key: string]: { fields: string; relat: string; value: string } }>({})

  const fieldVlaue3 = ref<string[]>([])

  const rules = reactive({
    ruleName: requireRules,
  })

  const currentdData = ref()

  const formRef = ref<FormInstance>()

  const fieldRef = ref<FormInstance>()

  const emit = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
    (e: 'on-reflash'): void
  }>()

  const handleClose = () => {
    emit('on-closeEvent', false)
  }

  const submitForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.validate(async (valid) => {
      if (valid) {
        formetData(fieldData3, fieldVlaue3.value)
      } else {
        return false
      }
    })
  }

  const formetData = async (resdata: any, strArr: string[]) => {
    cndList.value = []
    for (const key in resdata) {
      if (strArr.includes(key)) {
        const obj = {
          fields: '',
          relat: '=',
          value: [],
        }
        // obj.fields = resdata[key].fields
        obj.relat = resdata[key].relat
        obj.value = resdata[key].value
        cndList.value.push(obj)
      }
    }
    cndList.value.forEach((item: any, index: any) => {
      cndList.value[index].fields = fieldVlaue3.value[index]
    })
    const flag = cndList.value.some((item: any) => {
      return !item.value
    })
    isLoading.value = true
    try {
      if (flag) return ElMessage({ message: '提交失败，表单中存在空白项目', type: 'error' })
      formData.cnd = JSON.stringify(cndList.value)
      const { msg } = await LevelRuleSaveUpdateAPI({
        ...formData,
      })
      ElMessage({ message: msg, type: 'success' })
      handleClose()
      emit('on-reflash')
    } finally {
      isLoading.value = false
    }
  }

  const initData = () => {
    mode.value = props.remark
    if (mode.value == 'edit') {
      for (const key in formData) {
        if (key == 'cnd') {
          fieldVlaue3.value = []
          const res = JSON.parse(currentdData.value[key])
          res.forEach((item: any) => {
            fieldVlaue3.value.push(item.fields)
            fieldData3[item.fields] = { fields: item.fields, relat: item.relat, value: item.value }
          })
        } else if (key == 'level') {
          formData[key] = currentdData.value[key]
        } else {
          // @ts-ignore
          formData[key] = currentdData.value[key]
        }
      }
    }
  }

  onMounted(async () => {
    visible.value = props.isShow
    currentdData.value = props.curData
    setTimeout(() => {
      initData()
    }, 10)
  })

  const handlerChange3 = async (index: number, str: string) => {
    filterHas()
  }

  const filterHas = () => {
    fieldList1.forEach((td: any) => {
      if (fieldVlaue3.value.includes(td.fieldNameEn)) {
        td['disabled'] = true
      } else {
        td['disabled'] = false
      }
    })
  }

  const handlerAddRow = () => {
    if (Object.keys(fieldVlaue3.value).length >= fieldList1.length)
      return ElMessage({ message: '已是最大长度，不能再添加', type: 'warning' })
    const arr: any[] = []
    fieldList1.filter((item: any) => {
      if (!fieldVlaue3.value.includes(item.fieldNameEn)) {
        arr.push(item)
      }
    })
    fieldVlaue3.value.push(arr[0].fieldNameEn)
    fieldData3[arr[0].fieldNameEn] = { fields: arr[0].fieldNameEn, relat: '=', value: '' }
    filterHas()
  }

  const handlerDelete = (str: string) => {
    $baseConfirm('你确定要删除该项吗？', null, async () => {
      fieldData3[str] = { fields: '', relat: '=', value: '' }
      fieldVlaue3.value = fieldVlaue3.value.filter((item: any) => {
        return item != str
      })
      filterHas()
    })
  }

  const resetData = () => {
    fieldList1 = []
    fieldVlaue3.value = []
  }

  watch(
    () => formData.indexType,
    () => {
      resetData()
      const res = getTableColumn(formData.indexType)
      fieldList1 = res.filter((item: any) => {
        return item.isDeleted == 2
      })
      res.forEach((item: any) => {
        const field = {
          fieldNameCn: '',
          fieldNameEn: '',
        }
        field.fieldNameCn = item.fieldNameCn
        field.fieldNameEn = item.fieldNameEn
        fieldRules[item.fieldNameEn] = requireRules
        fieldData3[item.fieldNameEn] = {
          fields: field.fieldNameEn,
          relat: '=',
          value: '',
        }
      })
      fieldVlaue3.value.push(fieldList1[0].fieldNameEn)
      initData()
      filterHas()
    },
    { immediate: true }
  )
</script>

<template>
  <div class="edit-rule">
    <el-dialog v-model="visible" :before-close="handleClose" destroy-on-close :title="title" width="848px">
      <el-form ref="formRef" label-position="right" label-width="120px" :model="formData" :rules="rules">
        <el-form-item label="名称：" prop="ruleName">
          <el-input v-model="formData.ruleName" placeholder="请输入名称" style="width: 655px" />
        </el-form-item>
        <el-form-item label="选择等级：" prop="saveIndex">
          <el-select v-model="formData.level" style="width: 655px">
            <el-option v-for="item in levelOption" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="选择协议：" prop="saveIndex">
          <el-select v-model="formData.indexType" style="width: 655px">
            <el-option v-for="item in tagType" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="解析内容：">
          <div class="wrap">
            <div class="item_3">
              <el-form
                ref="fieldRef"
                label-position="right"
                label-width="140px"
                :model="fieldData3"
                :rules="fieldRules"
              >
                <el-form-item v-for="(item, index) in fieldVlaue3" :key="item">
                  <template #label>
                    <el-select
                      v-model="fieldVlaue3[index]"
                      style="width: 100%"
                      @change="handlerChange3(index, fieldData3[fieldVlaue3[index]].fields)"
                    >
                      <el-option
                        v-for="item in fieldList1"
                        :key="item.fieldNameEn"
                        :disabled="item.disabled"
                        :label="item.fieldNameCn"
                        :value="item.fieldNameEn"
                      />
                    </el-select>
                  </template>
                  <el-select v-model="fieldData3[fieldVlaue3[index]].relat" style="width: 30%; margin-right: 10px">
                    <el-option v-for="item in relatOption" :key="item.value" :label="item.label" :value="item.value" />
                  </el-select>
                  <el-input v-model="fieldData3[fieldVlaue3[index]].value" style="width: 30%" />
                  <el-icon
                    v-if="fieldVlaue3.length > 1"
                    class="my_icon"
                    color="#448ef7"
                    style="margin: 0 0 16px 18px"
                    @click="handlerDelete(fieldVlaue3[index])"
                  >
                    <Delete />
                  </el-icon>
                </el-form-item>

                <el-form-item>
                  <el-button type="primary" @click="handlerAddRow">添加</el-button>
                </el-form-item>
              </el-form>
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button :loading="isLoading" type="primary" @click="submitForm(formRef)">确认</el-button>
        <el-button @click="handleClose">取消</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .wrap {
    width: 655px;
    height: 325px;
    padding: 25px;
    background: #f5f7fa;
    border-radius: 4px;
    overflow-y: auto;
    &::-webkit-scrollbar {
      width: 0;
      height: 0;
    }

    :deep() {
      .el-input {
        margin-bottom: 15px;
      }
      .my_icon {
        cursor: pointer;
      }
    }
  }
</style>
