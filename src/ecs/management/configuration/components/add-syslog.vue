<script lang="ts">
  export default {
    name: 'AddSyslog',
  }
</script>
<script setup lang="ts">
  import { requireRules } from '~/src/utils/rules'

  import { FormInstance } from 'element-plus'

  import { Delete, Plus } from '@element-plus/icons-vue'

  import { syslogSaveOrUpdateApi, parseLogApi, getWarnRuleDictApi, getByWarnRuleIdApi } from '~/src/api-ecs/kafka'

  import { uuid } from '~/src/utils'

  import { syslogSaveOrUpdateType, warnRuleAnalysFieldsListType } from '@/types'

  import { useUserStore } from '@/store/modules/user'
  import AesEncryptCBC from '~/src/utils/crypto'

  // import { Editor } from '@wangeditor/editor-for-vue'
  // import { en } from 'element-plus/es/locale'

  const userStore = useUserStore()

  const { getTableColumn } = userStore

  // 表单数据
  const formData = reactive<syslogSaveOrUpdateType>({
    parseList: [
      {
        parseType: 0,
        matchRegex: '',
        // matchRegex: 'srcmac=([^\\f]+)\\sproto=(\\w+)[^\\f]+dst=(\\d{1,3}.\\d{1,3}.\\d{1,3}.\\d{1,3}:\\d{1,5})',
        matchLog: '',
        fieldMappings: [],
        // fieldMappingsVo: [],
        groupSplit: '',
        kvSplit: '',
        appendFields: [],
        matchJsonFields: [],
      },
    ] as warnRuleAnalysFieldsListType[],
    clientIp: '',
    enable: 1,
    id: undefined,
    logSample: '',
    ruleName: '',
    saveIndex: 0,
    prefixRegex: '',
    // prefixRegex: 'srcmac=00:05:1b:ad:0b:6b dstmac=00:50:56:b4:88:b3',
  })
  // \{.*\}
  // srcmac=(\S+)[^\f]+\sproto=(\w+)[^\f]+dst=(\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}:\d{1,5})

  const formRef = ref<FormInstance>() // 表单实例

  const saveIndexOption = [
    { label: '告警', value: 0 },
    { label: '审计日志', value: 1 },
    { label: '用户日志', value: 2 },
  ]

  const logTypeOption = [
    { label: '正则', value: 0 },
    { label: 'json', value: 1 },
    { label: 'kv', value: 2 },
  ]

  const listLoading = ref(false) // 是否加载

  // 表单数据校验
  const rules = reactive({
    ruleName: requireRules,
    saveIndex: requireRules,
    // parseType: requireRules,
    enable: requireRules,
    logSample: requireRules,
    prefixRegex: requireRules,
    parseList: requireRules,
  })

  const dateWarningarnRuleDict = ref()

  const defalutWarningarnRuleDict = ref()

  // 映射字段配置
  const mapFieldOptions = ref()

  const props = defineProps<{
    showDrawer: boolean
    mode: string
    currentData?: syslogSaveOrUpdateType
  }>()

  const emit = defineEmits<{
    (e: 'on-closeEvent'): void
    (e: 'on-reflash'): void
  }>()

  const nodeRef = ref()

  const disabled = ref(false)

  const btnDisabled = ref(false)

  const drawer = ref<boolean>()

  const onBlur = () => {
    // formData.logSample = nodeRef.value.innerHTML
    formData.logSample = nodeRef.value.innerText
  }

  // 获取字典
  const getWarnRuleDict = async () => {
    const { data } = await getWarnRuleDictApi()
    defalutWarningarnRuleDict.value = data.defaultDict
    dateWarningarnRuleDict.value = data.dateDict
  }

  // 解析类型改变
  const handleChangeParseType = (index: number) => {
    formData.parseList[index].fieldMappings = []
  }

  // 正则回显高亮
  const handleHighlight = () => {
    const search = nodeRef.value.innerText
    if (search <= 0) return
    let metaJson = JSON.parse(JSON.stringify(search))
    let indexArr: number[] = []
    formData.parseList.forEach((item) => {
      item.fieldMappings?.forEach((td) => {
        indexArr.push(td.startIndex)
        indexArr.push(td.endIndex)
      })
    })
    indexArr = indexArr.filter((item) => {
      return item == 0 ? true : !!item
    })
    if (indexArr.length <= 1) return
    const satrtStr = '<span style="background-color: yellow;">'
    const endStr = '</span>'
    let str = ''
    const strArr = []
    const str1 = metaJson.substring(0, indexArr[0])
    strArr.push(str1)
    for (let index = 1; index < indexArr.length; index++) {
      if (index % 2 != 0) {
        const res = metaJson.substring(indexArr[index - 1], indexArr[index])
        const reg = satrtStr + res + endStr
        strArr.push(reg)
        if (index == indexArr.length - 1) {
          const res = metaJson.substring(indexArr[index], metaJson.length)
          strArr.push(res)
        }
      } else {
        const res = metaJson.substring(indexArr[index - 1], indexArr[index])
        strArr.push(res)
      }
    }
    strArr.forEach((item) => {
      str += item
    })
    // nodeRef.value.innerHTML = str
    nodeRef.value.innerText = str
  }

  // 去掉标签
  // const toClearElement = () => {
  //   if (formData.logSample.length <= 0) return
  //   let metaJson = JSON.parse(JSON.stringify(formData.logSample))
  //   const satrtStr = '<span style="background-color: yellow;">'
  //   const p = '<p>'
  //   const p1 = '</p>'
  //   const br = '<br />'
  //   const endStr = '</span>'
  //   const nbsp = ' &nbsp;'
  //   const reg = new RegExp(endStr, 'gi')
  //   const reg1 = new RegExp(p, 'gi')
  //   const reg2 = new RegExp(p1, 'gi')
  //   const reg3 = new RegExp(satrtStr, 'gi')
  //   const reg4 = new RegExp(br, 'gi')
  //   const reg5 = new RegExp(nbsp, 'gi')
  //   let str = metaJson.replaceAll(reg, '')
  //   str = str.replaceAll(reg1, '')
  //   str = str.replaceAll(reg2, '')
  //   str = str.replaceAll(reg3, '')
  //   str = str.replaceAll(reg4, '')
  //   str = str.replaceAll(reg5, '')
  //   formData.logSample = str
  // }

  // 解析
  const handleParse = async () => {
    try {
      disabled.value = true
      const regexArr = []
      formData.parseList.forEach((item: any) => {
        item.fieldMappingsVo = null
        if (item.matchRegex) regexArr.push(true)
      })
      if (regexArr.length != formData.parseList.length)
        return ElMessage({ message: '匹配正则不能为空', type: 'warning' })
      // toClearElement()
      const { data, msg } = await parseLogApi({
        logSample: nodeRef.value.innerText,
        prefixRegex: formData.prefixRegex,
        saveIndex: formData.saveIndex,
        parseList: formData.parseList,
      })
      formData.parseList = data.parseVoList
      formData.parseList.forEach((item: any) => {
        // @ts-ignore
        item.fieldMappings = JSON.parse(JSON.stringify(item.fieldMappingsVo))
      })
      ElMessage({ message: msg, type: 'success' })
      btnDisabled.value = true
      handleChange()
      // handleHighlight()
    } finally {
      disabled.value = false
    }
  }

  // 互斥
  const handleChange = () => {
    const arr = [] as string[]
    formData.parseList?.forEach((item: any) => {
      if (item.fieldMappings?.length > 0) {
        item.fieldMappings.forEach((td: any) => {
          if (td.mappingField) {
            arr.push(td.mappingField)
          }
          if (td.children.length > 0) {
            td.children.forEach((warp: any) => {
              arr.push(warp.mappingField)
            })
          }
        })
      }
    })
    mapFieldOptions.value?.forEach((td: any) => {
      td['disabled'] = false
      if (arr.includes(td.fieldNameEn)) {
        td['disabled'] = true
      }
    })
  }

  // 添加子解析项
  const handleAddChild = (row: { uuid: number }, index: number) => {
    const child = {
      dictId: null,
      endIndex: null,
      groupIndex: null,
      jsonKey: '',
      level: 2,
      mappingField: '',
      mappingValue: '',
      sourceValue: '',
      uuid: uuid(),
      startIndex: null,
      children: [],
    }
    const td = formData.parseList[index].fieldMappings?.find((item) => {
      return item.uuid == row.uuid
    })
    child.jsonKey = td.jsonKey
    child.sourceValue = td.sourceValue
    formData.parseList[index].fieldMappings?.forEach((item) => {
      if (item.uuid == row.uuid) {
        item.children.push(child)
      }
    })
  }

  // 追加匹配字段
  const handleAddMergeFiled = (index: number) => {
    formData.parseList[index].matchJsonFields?.push({
      matchUuid: uuid(),
    })
  }

  // 追加原始字段
  const handleAddFiled = (index: number) => {
    formData.parseList[index].appendFields?.push({
      appendUuid: uuid(),
    })
  }

  // 添加
  const handleAdd = () => {
    const parseList: warnRuleAnalysFieldsListType = {
      parseType: 0,
      matchRegex: '',
      matchLog: '',
      fieldMappings: [],
      // fieldMappingsVo?: undefined,
      groupSplit: '',
      kvSplit: '',
      appendFields: [],
      matchJsonFields: [],
    }
    formData.parseList.push(parseList)
  }

  // 删除子解析项
  const handleDeleteChild = (row: { uuid: number }, index: number) => {
    formData.parseList[index].fieldMappings?.forEach((item) => {
      if (item.children.length > 0) {
        item.children = item.children.filter((td: any) => {
          return td.uuid != row.uuid
        })
      }
    })
  }

  // 删除解析器
  const handleDelete = (index: number) => {
    formData.parseList = formData.parseList.filter((item: any, idx: any) => {
      return idx != index
    })
  }

  //  删除匹配字段
  const handleDeleteMergeFiled = (row: { matchUuid: number }, index: number) => {
    formData.parseList[index].matchJsonFields = formData.parseList[index]?.matchJsonFields?.filter((item: any) => {
      return item.matchUuid !== row.matchUuid
    })
  }

  //  删除原始字段
  const handleDeleteFiled = (row: { appendUuid: number }, index: number) => {
    formData.parseList[index].appendFields = formData.parseList[index]?.appendFields?.filter((item: any) => {
      return item.appendUuid !== row.appendUuid
    })
  }

  // 回显数据
  const initData = async () => {
    drawer.value = props.showDrawer
    if (props.mode == 'edit') {
      btnDisabled.value = true
      const id = JSON.parse(JSON.stringify(props.currentData))?.id
      if (id) {
        const { data } = await getByWarnRuleIdApi({ id })
        Object.keys(formData).forEach((item) => {
          // @ts-ignore
          formData[item] = data[item]
        })
        formData.parseList = data.parseVoList
        formData.parseList.forEach((item: any) => {
          // @ts-ignore
          item.fieldMappings = item.fieldMappingsVo
        })
        // nodeRef.value.innerHTML = formData.logSample
        nodeRef.value.innerText = formData.logSample
      }
    }
  }

  // 保存
  const confirmClick = async (formEl: FormInstance | undefined) => {
    if (formData.parseList.length == 0) ElMessage({ message: '表格不能为空', type: 'warning' })
    const falgArr = []
    const regexArr = []
    formData.parseList.forEach((item: any) => {
      let matchFlag = false
      let appendFlag = false
      item.fieldMappingsVo = null
      if (item.matchRegex) regexArr.push(true)
      if (item?.appendFields?.length > 0) {
        appendFlag = item?.appendFields.some((td: any) => {
          return !td.appendKey || !td.appendValue
        })
        if (appendFlag) falgArr.push(appendFlag)
      }
      if (item?.matchJsonFields?.length > 0) {
        matchFlag = item?.matchJsonFields.some((td: any) => {
          return !td.matchKey || !td.matchKey
        })
        if (matchFlag) falgArr.push(matchFlag)
      }
    })
    if (regexArr.length != formData.parseList.length) return ElMessage({ message: '匹配正则不能为空', type: 'warning' })
    if (falgArr.length > 0) return ElMessage({ message: '表格存在空值', type: 'warning' })
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        // toClearElement()
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
                const { msg } = await syslogSaveOrUpdateApi(
                  {
                    ...formData,
                    logSample: nodeRef.value.innerText,
                  },
                  { password: AesEncryptCBC(password.value) }
                )
                emit('on-reflash')
                ElMessage({ message: msg, type: 'success' })
                closeEvent()
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
        console.log('error submit!', fields)
      }
    })
  }

  const closeEvent = () => {
    emit('on-closeEvent')
  }

  onMounted(async () => {
    await initData()
    mapFieldOptions.value = getTableColumn(10)
    getWarnRuleDict()
    handleChange()
    // setTimeout(() => {
    //   handleHighlight()
    // }, 0)
  })
</script>

<template>
  <div class="add-kafka">
    <el-drawer v-model="drawer" :close-on-click-modal="false" size="70%" @close="closeEvent">
      <template #header>
        <h4>{{ mode == 'edit' ? '编辑' : '新建' }}</h4>
      </template>
      <template #default>
        <el-form ref="formRef" label-position="right" label-width="160px" :model="formData" :rules="rules">
          <div class="my-title"><span>基本信息</span></div>
          <el-form-item label="名称：" prop="ruleName">
            <el-input v-model="formData.ruleName" />
          </el-form-item>
          <el-form-item label="客户端地址：" prop="clientIp">
            <el-input v-model="formData.clientIp" placeholder="请输入客户端地址" />
          </el-form-item>
          <el-form-item label="选择入库：" prop="saveIndex">
            <el-select v-model="formData.saveIndex" style="width: 100%">
              <el-option v-for="item in saveIndexOption" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="是否启用：" prop="enable">
            <el-radio v-model="formData.enable" :label="1">启用</el-radio>
            <el-radio v-model="formData.enable" :label="0">停用</el-radio>
          </el-form-item>
          <el-form-item label="原始日志：" prop="logSample">
            <div
              ref="nodeRef"
              class="my-textarea"
              contenteditable="true"
              placeholder="请输入内容"
              style="overflow: auto; height: 120px; width: 100%; padding: 5px; border: 1px solid var(--el-border-color)"
              @blur="onBlur"
            ></div>
          </el-form-item>
          <el-form-item label="前置过滤：" prop="prefixRegex">
            <el-input v-model="formData.prefixRegex" placeholder="请输入前置过滤条件" />
          </el-form-item>
          <el-form-item label="" prop="prefixRegex">
            <el-button :disabled="disabled" type="primary" @click="handleParse">日志解析</el-button>
          </el-form-item>
          <div class="my-title">
            <el-button :icon="Plus" type="primary" @click="handleAdd">添加解析器</el-button>
          </div>
          <div v-for="(td, index) in formData.parseList" :key="index" class="parses">
            <el-button v-if="index !== 0" class="close" :icon="Delete" type="warning" @click="handleDelete(index)" />
            <el-form-item label="内容解析器：">
              <el-select v-model="td.parseType" style="width: 100%" @change="handleChangeParseType(index)">
                <el-option v-for="item in logTypeOption" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="匹配正则：" prop="parseList">
              <el-input v-model="td.matchRegex" />
            </el-form-item>
            <el-form-item label="匹配日志：">
              <el-input v-model="td.matchLog" />
            </el-form-item>
            <el-form-item v-if="td.parseType == 2" label="组分隔符：" prop="parseList">
              <el-input v-model="td.groupSplit" />
            </el-form-item>
            <el-form-item v-if="td.parseType == 2" label="键值对分隔符：" prop="parseList">
              <el-input v-model="td.kvSplit" />
            </el-form-item>
            <el-form-item label="内容解析：">
              <el-table
                v-loading="listLoading"
                :border="true"
                :data="td.fieldMappings"
                default-expand-all
                row-key="uuid"
              >
                <el-table-column
                  v-if="td.parseType == 1"
                  :align="'center'"
                  label="解析KEY"
                  prop="jsonKey"
                  show-overflow-tooltip
                />
                <el-table-column
                  v-else-if="td.parseType == 2"
                  :align="'center'"
                  label="解析KEY"
                  prop="kvKey"
                  show-overflow-tooltip
                />
                <el-table-column :align="'center'" label="映射字段" prop="mappingField" show-overflow-tooltip>
                  <template #default="{ row }">
                    <el-select
                      v-model="row.mappingField"
                      :allow-create="true"
                      class="m-2"
                      :clearable="true"
                      style="width: 80%"
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
                <el-table-column :align="'center'" label="字典表" prop="dictId" show-overflow-tooltip>
                  <template #default="{ row }">
                    <el-select
                      v-model="row.dictId"
                      :allow-create="true"
                      class="m-2"
                      :clearable="true"
                      style="width: 80%"
                    >
                      <template v-if="row.mappingField == 'startTimeNs'">
                        <el-option
                          v-for="item in dateWarningarnRuleDict"
                          :key="item.id"
                          :label="item.dictName"
                          placeholder="请选择"
                          :value="item.id"
                        />
                      </template>
                      <template v-else>
                        <el-option
                          v-for="item in defalutWarningarnRuleDict"
                          :key="item.id"
                          :label="item.dictName"
                          placeholder="请选择"
                          :value="item.id"
                        />
                      </template>
                    </el-select>
                  </template>
                </el-table-column>
                <el-table-column :align="'center'" label="提取值" prop="sourceValue" show-overflow-tooltip />
                <el-table-column :align="'center'" label="映射值" prop="mappingValue" show-overflow-tooltip />
                <el-table-column :align="'center'" label="操作">
                  <template #default="{ row }">
                    <el-button v-if="row.level != 2" size="small" @click="handleAddChild(row, index)">
                      添加子解析项
                    </el-button>
                    <el-button v-else size="small" @click="handleDeleteChild(row, index)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-form-item>
            <el-form-item :data="td.appendFields" label="原始日志追加字段：">
              <el-table v-loading="listLoading" :border="true" :data="td.appendFields">
                <el-table-column :align="'center'" label="字段名" prop="appendKey" show-overflow-tooltip>
                  <template #default="{ row }">
                    <el-input v-model="row.appendKey" :clearable="true" />
                  </template>
                </el-table-column>
                <el-table-column :align="'center'" label="值" prop="appendValue" show-overflow-tooltip>
                  <template #default="{ row }">
                    <el-input v-model="row.appendValue" :clearable="true" />
                  </template>
                </el-table-column>
                <el-table-column :align="'center'">
                  <template #header>
                    <el-button :icon="Plus" type="primary" @click="handleAddFiled(index)" />
                  </template>

                  <template #default="{ row }">
                    <el-button size="small" @click="handleDeleteFiled(row, index)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-form-item>
            <el-form-item v-if="td.parseType == 1" label="匹配字段：">
              <el-table v-loading="listLoading" :border="true" :data="td.matchJsonFields">
                <el-table-column :align="'center'" label="字段名" prop="matchKey" show-overflow-tooltip>
                  <template #default="{ row }">
                    <el-input v-model="row.matchKey" :clearable="true" />
                  </template>
                </el-table-column>

                <el-table-column :align="'center'" label="值" prop="matchValue" show-overflow-tooltip>
                  <template #default="{ row }">
                    <el-input v-model="row.matchValue" :clearable="true" />
                  </template>
                </el-table-column>
                <el-table-column :align="'center'">
                  <template #header>
                    <el-button :icon="Plus" type="primary" @click="handleAddMergeFiled(index)" />
                  </template>
                  <template #default="{ row }">
                    <el-button size="small" @click="handleDeleteMergeFiled(row, index)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-form-item>
          </div>

          <el-form-item style="margin-top: 20px">
            <el-button type="primary" @click="confirmClick(formRef)">保存</el-button>
            <el-button @click="closeEvent">取消</el-button>
          </el-form-item>
        </el-form>
      </template>
    </el-drawer>
  </div>
</template>

<style lang="scss" scoped>
  .add-kafka {
    padding-right: 40px;
    :deep() {
      .el-table__body-wrapper .el-scrollbar {
        height: initial !important;
      }
    }

    .parses {
      padding: 20px 0 20px;
      border-top: 1px solid var(--el-border-color);
      border-bottom: 1px solid var(--el-border-color);
      margin-bottom: 15px;
      position: relative;
      .close {
        position: absolute;
        right: 0px;
        bottom: 3px;
      }
    }

    .my-title {
      margin-bottom: 30px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      // &::before {
      //   display: inline-block;
      //   width: 3px;
      //   height: 10px;
      //   margin: 0 6px 0 30px;
      //   content: '';
      //   background: #1890ff;
      // }
    }
  }
</style>
