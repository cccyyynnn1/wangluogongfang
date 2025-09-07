<script lang="ts">
  export default {
    name: 'WhiteList', // 白名单
  }
</script>

<script setup lang="ts">
  import { onMounted, ref } from 'vue'

  import type { FormInstance } from 'element-plus'

  import { requireRules } from '~/src/utils/rules'

  import { GetConfigFieldsApi, WarnWhiteSaveOrUpdatesApi } from '~/src/api-ecs/alert'

  const $baseConfirm: any = inject('$baseConfirm')

  const props = defineProps<{
    isShow: boolean
    curData: any
    mode?: string
  }>()

  const visible = ref(false) // 显隐

  const isLoading = ref(false) // 加载

  const title = ref('添加白名单')

  // 表单数据
  const formData = reactive({
    ruleName: '',
    type: 0,
    threatName: '',
    id: undefined,
  })

  const cndList = ref()

  const fieldList1: { cn: string; en: string; disabled?: boolean }[] = []

  const fieldList2: { cn: string; en: string; disabled?: boolean }[] = []

  const fieldRules = reactive<{ [key: string]: any }>({})

  const fieldData1 = reactive<{ [key: string]: string[] }>({})

  const fieldVlaue1 = ref<string[]>(['attackIp', 'victimIp'])

  const fieldData2 = reactive<{ [key: string]: string[] }>({})

  const fieldVlaue2 = ref<string[]>(['attackIp', 'victimIp', 'targetPort'])

  const fieldData3 = reactive<{ [key: string]: string[] }>({})

  const fieldVlaue3 = ref<string[]>([])

  const fieldData4 = reactive<{ [key: string]: string[] }>({})

  const fieldVlaue4 = ref<string[]>([])

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
        if (formData.type == 0) {
          formetData(fieldData1, fieldVlaue1.value)
        } else if (formData.type == 1) {
          formetData(fieldData2, fieldVlaue2.value)
        } else if (formData.type == 2) {
          formetData(fieldData3, fieldVlaue3.value)
        } else if (formData.type == 3) {
          formetData(fieldData4, fieldVlaue4.value)
        }
      } else {
        console.log('error submit!')
        return false
      }
    })
  }

  const formetData = async (resdata: any, strArr: string[]) => {
    cndList.value = []
    for (const key in resdata) {
      if (strArr.includes(key)) {
        const obj = {
          field: '',
          values: [],
        }
        obj.field = key
        obj.values = resdata[key]
        cndList.value.push(obj)
      }
    }
    const flag = cndList.value.some((item: any) => {
      return item.values.length == 0 || (item.values.length == 1 && item.values[0] == '')
    })
    isLoading.value = true
    try {
      if (flag) return ElMessage({ message: '提交失败，表单中存在空白项目', type: 'error' })
      // if (!currentdData.value.ruleId) return ElMessage({ message: '提交失败，规则id为空', type: 'error' })
      const { data, msg } = await WarnWhiteSaveOrUpdatesApi({
        ...formData,
        ruleId: currentdData.value.ruleId,
        cndList: cndList.value,
      })
      ElMessage({ message: msg, type: 'success' })
      handleClose()
      emit('on-reflash')
    } finally {
      isLoading.value = false
    }
  }

  const GetConfigFields = async () => {
    const { data } = await GetConfigFieldsApi()
    if (data.length == 0) return ElMessage({ message: '获取默认字段失败', type: 'error' })
    cndList.value = []
    data.forEach((item: any) => {
      const obj = {
        field: '',
        values: [],
      }
      const field = {
        cn: '',
        en: '',
      }
      obj.field = item.fieldNameEn
      field.cn = item.fieldNameCn
      field.en = item.fieldNameEn
      fieldRules[item.fieldNameEn] = requireRules
      fieldData1[item.fieldNameEn] = ['']
      fieldData2[item.fieldNameEn] = ['']
      fieldData3[item.fieldNameEn] = ['']
      fieldData4[item.fieldNameEn] = ['']
      fieldList1.push(field)
      fieldList2.push(field)
      // cndList.value.push(obj)
    })
    fieldVlaue3.value.push(fieldList1[0].en)
    // fieldVlaue4.value.push(fieldList2[0].en)
  }

  onMounted(async () => {
    visible.value = props.isShow
    currentdData.value = props.curData
    formData.threatName = currentdData.value.threatName
    await GetConfigFields()
    if (formData.type == 0 && props.mode != 'edit') {
      fieldData1.attackIp[0] = currentdData.value.attackIp
      fieldData1.victimIp[0] = currentdData.value.victimIp
    } else if (formData.type == 1 && props.mode != 'edit') {
      fieldData2.attackIp[0] = currentdData.value.attackIp
      fieldData2.victimIp[0] = currentdData.value.victimIp
      fieldData2.targetPort[0] = currentdData.value.targetPort
    }
    if (props.mode == 'edit') {
      formData.ruleName = currentdData.value.ruleName
      formData.type = currentdData.value.type
      formData.id = currentdData.value.id
      setTimeout(() => {
        ininForm()
      }, 0)
    }
    filterHas()
    filterHas2()
  })

  const ininForm = () => {
    const obj = {
      0: fieldData1,
      1: fieldData2,
      2: fieldData3,
      3: fieldData4,
    }
    const filedsValue = {
      0: fieldVlaue1.value,
      1: fieldVlaue2.value,
      2: fieldVlaue3.value,
      3: fieldVlaue4.value,
    }
    // @ts-ignore
    const field = filedsValue[formData.type]
    currentdData.value.cndList.forEach((item: any) => {
      // @ts-ignore
      obj[formData.type][item.field] = item.values
      if (!field?.includes(item.field)) {
        field.push(item.field)
      }
    })
  }

  const handlerChange3 = (index: number, str: string) => {
    fieldData3[str] = ['']
    filterHas()
  }

  const handlerChange4 = (index: number, str: string) => {
    fieldData4[str] = ['']
    filterHas2()
  }

  const filterHas = () => {
    fieldList1.forEach((td: any) => {
      if (fieldVlaue3.value.includes(td.en)) {
        td['disabled'] = true
      } else {
        td['disabled'] = false
      }
    })
  }

  const filterHas2 = () => {
    fieldList2.forEach((td: any) => {
      if (fieldVlaue4.value.includes(td.en)) {
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
      if (!fieldVlaue3.value.includes(item.en)) {
        arr.push(item)
      }
    })
    fieldVlaue3.value.push(arr[0].en)
    fieldData3[arr[0].en] = ['']
    filterHas()
  }

  const handlerAddRow2 = () => {
    if (Object.keys(fieldVlaue4.value).length >= fieldList2.length)
      return ElMessage({ message: '已是最大长度，不能再添加', type: 'warning' })
    const arr: any[] = []
    fieldList2.filter((item: any) => {
      if (!fieldVlaue4.value.includes(item.en)) {
        arr.push(item)
      }
    })
    fieldVlaue4.value.push(arr[0].en)
    fieldData4[arr[0].en] = ['']
    filterHas2()
  }

  const handlerAddLine = (value: string) => {
    fieldData3[value].push('')
  }

  const handlerAddLine2 = (value: string) => {
    fieldData4[value].push('')
  }

  const handlerDelete = (str: string, index: number) => {
    if (index == 0) {
      $baseConfirm('你确定要删除该项吗？', null, async () => {
        fieldData3[str] = ['']
        fieldVlaue3.value = fieldVlaue3.value.filter((item: any) => {
          return item != str
        })
        filterHas()
      })
    } else {
      $baseConfirm('你确定要删除该行吗？', null, async () => {
        fieldData3[str].splice(index, 1)
      })
    }
    filterHas()
  }

  watch(
    () => formData.type,
    () => {
      if (formData.type == 2) {
        filterHas()
      } else if (formData.type == 3) {
        filterHas2()
      } else if (formData.type == 1 && props.mode != 'edit') {
        fieldData2.attackIp[0] = currentdData.value.attackIp
        fieldData2.victimIp[0] = currentdData.value.victimIp
        fieldData2.targetPort[0] = currentdData.value.targetPort
      }
    }
  )

  const handlerDelete2 = (str: string, index: number) => {
    if (index == 0) {
      $baseConfirm('你确定要删除该项吗？', null, async () => {
        fieldData4[str] = ['']
        fieldVlaue4.value = fieldVlaue4.value.filter((item: any) => {
          return item != str
        })
        filterHas2()
      })
    } else {
      $baseConfirm('你确定要删除该行吗？', null, async () => {
        fieldData4[str].splice(index, 1)
      })
    }
    filterHas2()
  }
</script>

<template>
  <div class="white-list">
    <el-dialog
      v-model="visible"
      append-to-body
      :before-close="handleClose"
      destroy-on-close
      :title="title"
      width="848px"
    >
      <el-form ref="formRef" label-position="right" label-width="120px" :model="formData" :rules="rules">
        <el-form-item label="名称：" prop="ruleName">
          <el-input v-model="formData.ruleName" placeholder="请输入名称" style="width: 655px" />
        </el-form-item>
        <el-form-item label="添加方式：" prop="type">
          <el-radio-group v-model="formData.type">
            <el-radio :label="0" size="large">源IP -> 目的IP</el-radio>
            <el-radio :label="1" size="large">源IP -> 目的IP:目的端口</el-radio>
            <el-radio :label="2" size="large">自定义</el-radio>
            <el-radio :label="3" size="large">Any -> Any</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <div class="wrap">
            <div v-if="formData.type == 0" class="item_1">
              <el-form
                ref="fieldRef"
                label-position="right"
                label-width="100px"
                :model="fieldData3"
                :rules="fieldRules"
              >
                <el-form-item label="源IP：" prop="attackIp">
                  <template v-for="(item, index) in fieldData1.attackIp" :key="index">
                    <el-input v-model="fieldData1.attackIp[index]" style="width: 60%" />
                  </template>
                </el-form-item>
                <el-form-item label="目的IP：" prop="victimIp">
                  <template v-for="(item, index) in fieldData1.victimIp" :key="index">
                    <el-input v-model="fieldData1.victimIp[index]" style="width: 60%" />
                  </template>
                </el-form-item>
              </el-form>
            </div>
            <div v-else-if="formData.type == 1" class="item_2">
              <el-form
                ref="fieldRef"
                label-position="right"
                label-width="100px"
                :model="fieldData2"
                :rules="fieldRules"
              >
                <el-form-item label="源IP：" prop="attackIp">
                  <template v-for="(item, index) in fieldData2.attackIp" :key="index">
                    <el-input v-model="fieldData2.attackIp[index]" style="width: 60%" />
                  </template>
                </el-form-item>
                <el-form-item label="目的IP：" prop="victimIp">
                  <template v-for="(item, index) in fieldData2.victimIp" :key="index">
                    <el-input v-model="fieldData2.victimIp[index]" style="width: 60%" />
                  </template>
                </el-form-item>
                <el-form-item label="目的端口：" prop="targetPort">
                  <template v-for="(item, index) in fieldData2.targetPort" :key="index">
                    <el-input v-model="fieldData2.targetPort[index]" style="width: 60%" />
                  </template>
                </el-form-item>
              </el-form>
            </div>
            <div v-else-if="formData.type == 2" class="item_3">
              <el-form
                v-if="fieldVlaue3.length > 0"
                ref="fieldRef"
                label-position="right"
                label-width="140px"
                :model="fieldData3"
                :rules="fieldRules"
              >
                <el-form-item v-for="(item, index) in fieldVlaue3" :key="index" :prop="item">
                  <template #label>
                    <el-select
                      v-model="fieldVlaue3[index]"
                      style="width: 100%"
                      @change="handlerChange3(index, fieldVlaue3[index])"
                    >
                      <el-option
                        v-for="item in fieldList1"
                        :key="item.en"
                        :disabled="item.disabled"
                        :label="item.cn"
                        :value="item.en"
                      />
                    </el-select>
                  </template>
                  <template v-for="(it, idx) in fieldData3[fieldVlaue3[index]]" :key="idx">
                    <el-input v-model="fieldData3[fieldVlaue3[index]][idx]" style="width: 60%" />
                    <el-icon
                      v-if="1 + idx == fieldData3[fieldVlaue3[index]].length"
                      class="my_icon"
                      color="#448ef7"
                      style="margin: 0 0 16px 18px"
                      @click="handlerAddLine(fieldVlaue3[index])"
                    >
                      <Plus />
                    </el-icon>
                    <el-icon
                      v-if="
                        fieldVlaue3.length != 1 ||
                        (fieldData3[fieldVlaue3[index]].length > 1 && fieldVlaue3.length == 1 && idx != 0)
                      "
                      class="my_icon"
                      color="#448ef7"
                      style="margin: 0 0 16px 18px"
                      @click="handlerDelete(fieldVlaue3[index], idx)"
                    >
                      <Delete />
                    </el-icon>
                  </template>
                </el-form-item>

                <el-form-item label=" ">
                  <el-button type="primary" @click="handlerAddRow">添加</el-button>
                </el-form-item>
              </el-form>
            </div>
            <div v-else-if="formData.type == 3" class="item_4">
              <h5>排除以下加白</h5>
              <el-form
                ref="fieldRef"
                label-position="right"
                label-width="140px"
                :model="fieldData4"
                :rules="fieldRules"
              >
                <el-form-item v-for="(item, index) in fieldVlaue4" :key="index" :prop="item">
                  <template #label>
                    <el-select
                      v-model="fieldVlaue4[index]"
                      style="width: 100%"
                      @change="handlerChange4(index, fieldVlaue4[index])"
                    >
                      <el-option
                        v-for="item in fieldList2"
                        :key="item.en"
                        :disabled="item.disabled"
                        :label="item.cn"
                        :value="item.en"
                      />
                    </el-select>
                  </template>
                  <template v-for="(it, idx) in fieldData4[fieldVlaue4[index]]" :key="idx">
                    <el-input v-model="fieldData4[fieldVlaue4[index]][idx]" style="width: 60%" />
                    <el-icon
                      v-if="1 + idx == fieldData4[fieldVlaue4[index]].length"
                      class="my_icon"
                      color="#448ef7"
                      style="margin: 0 0 16px 18px"
                      @click="handlerAddLine2(fieldVlaue4[index])"
                    >
                      <Plus />
                    </el-icon>
                    <el-icon
                      v-if="
                        fieldVlaue4.length != 1 ||
                        (fieldData4[fieldVlaue4[index]].length > 1 && fieldVlaue4.length == 1 && idx != 0)
                      "
                      class="my_icon"
                      color="#448ef7"
                      style="margin: 0 0 16px 18px"
                      @click="handlerDelete2(fieldVlaue4[index], idx)"
                    >
                      <Delete />
                    </el-icon>
                  </template>
                </el-form-item>

                <el-form-item label=" ">
                  <el-button type="primary" @click="handlerAddRow2">添加</el-button>
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
