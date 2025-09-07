<script lang="ts">
  export default {
    name: 'EditWindow',
  }
</script>
<script setup lang="ts">
  import VueDraggable from 'vuedraggable'

  // import { searchBySqlApi } from '~/src/api-ecs/retrieve'

  import { requireRules } from '@/utils/rules'

  import type { FormInstance } from 'element-plus'

  import { filedDataType, MonitoringItem, IndexTypeTpye } from '@/types'

  import {
    siteGetAllPageApi,
    getApiInterfaceApiByIdApi,
    siteSearchApi,
    getBySiteIdApi,
    getAllDefaultFieldsApi,
  } from '~/src/api-ecs/site'

  import { ElMessage } from 'element-plus'

  import { useUserStore } from '@/store/modules/user'

  // import { usePubilcStore } from '@/store/modules/public'

  import { getFlowProbesAPI } from '~/src/api-ecs/retrieve'

  const { getTableColumn, getAllIndexINcouldsNetType } = useUserStore()

  // const publicStore = usePubilcStore()

  // const { getFlowProbesList, getFlowProbesIDList } = publicStore

  const indexTypeOption = ref<IndexTypeTpye[]>(getAllIndexINcouldsNetType())

  const props = defineProps<{
    showpage: boolean
    currentItem: MonitoringItem
    mode: string
    title: string
  }>()

  const emit = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
    (e: 'on-reflash', val: MonitoringItem): void
  }>()

  const visible = ref(false)

  const formRef = ref<FormInstance>()

  const modes = ref('')

  const titles = ref('')

  const allHttpField = ref<filedDataType[]>([]) // http字段
  const allDefaultField = ref<any>({}) // http字段

  const currentItems = ref()

  // const flowProbesValue = ref(['0'])
  // 表单数据
  const formData = reactive({
    name: '名称',
    siteApiId: undefined,
    siteSessionId: '',
    startTime: '',
    endTime: '',
    searchSql: '',
    indexType: 1,
    pageNum: 1,
    pageSize: 10,
    orderField: 'requestTimeNs',
    orderType: 'desc',
    displayFields: [] as filedDataType[],
    module: '0',
    ip: '',
    clientIp: '',
    serverIp: '',
    serverPort: '',
    topValue: '',
    leftValue: undefined,
    searchModel: '',
    flowProbeIds: ['0'],
    // id: undefined,
  })

  const SearchOption = [
    { value: 'ip', label: 'IP地址' },
    { value: 'clientIp', label: '源IP' },
    { value: 'serverIp', label: '目的IP' },
    { value: 'clientServerIp', label: '源+目的IP' },
  ]

  const topOptions = [
    { label: 'Top100', value: '100' },
    { label: 'Top500', value: '500' },
    { label: 'Top1000', value: '1000' },
    { label: 'Top2000', value: '2000' },
    { label: 'Top5000', value: '5000' },
  ]

  const siteOption = ref()

  const siteApiIdOption = ref()

  // 表单数据校验
  const rules = reactive({
    name: requireRules,
    siteSessionId: requireRules,
    searchSql: requireRules,
    displayFields: requireRules,
  })

  // 获取站点
  const getAllSite = async () => {
    const { data } = await siteGetAllPageApi()
    siteOption.value = data
    // 判断formData.siteSessionId的值是否在siteOption.value中
    if (!siteOption.value.some((item: any) => item.id === formData.siteSessionId)) {
      formData.siteSessionId = ''
    }
  }

  // 获取默认字段
  const getDefaultFields = async () => {
    const { data } = await getAllDefaultFieldsApi()
    allDefaultField.value = data
  }

  const getSonSite = async () => {
    const { data } = await getApiInterfaceApiByIdApi(formData.siteSessionId)
    siteApiIdOption.value = data
  }

  // 得到所有字段
  const getHttpField = async () => {
    if (formData.indexType == 22) {
      allHttpField.value = getTableColumn(1)?.filter((item: any) => {
        return item.fieldNameCn
      })
    } else {
      allHttpField.value = getTableColumn(formData.indexType)?.filter((item: any) => {
        return item.fieldNameCn
      })
    }
    if (formData.indexType == 1) {
      if (formData.siteSessionId) {
        allHttpField.value = getTableColumn(30).filter((item: any) => {
          return item.fieldNameCn
        })
        formData.displayFields = []
        if (currentItems.value.hasError) return
        const {
          data: { displayFieldsArr },
        } = await getBySiteIdApi({ id: formData.siteSessionId, apiId: formData.siteApiId })
        // formData.displayFields =
        displayFieldsArr.forEach((item: number) => {
          allHttpField.value.forEach((td) => {
            if (item == td.id) {
              formData.displayFields.push(td)
            }
          })
        })
      } else {
        const arr: number[] = []
        currentItems.value.data?.displayFields &&
          currentItems.value.data?.displayFields((item: { id: number }) => {
            arr.push(item.id)
          })
        const ids = arr.length > 0 ? arr : (allDefaultField.value[formData.indexType] as number[])
        formData.displayFields = allHttpField.value.filter((item: { isDisplay: number; id: number | string }) => {
          return ids?.includes(+item.id)
        })
        allHttpField.value = allHttpField.value.filter((item: { isDisplay: number; id: number | string }) => {
          return !ids?.includes(+item.id)
        })
      }
    } else {
      const arr: number[] = []
      currentItems.value.data?.displayFields &&
        currentItems.value.data?.displayFields((item: { id: number }) => {
          arr.push(item.id)
        })
      let ids = formData.indexType == 22 ? allDefaultField.value[1] : allDefaultField.value[formData.indexType]
      ids = arr.length > 0 ? arr : ids
      formData.displayFields = allHttpField.value.filter((item: { isDisplay: number; id: number | string }) => {
        return ids?.includes(+item.id)
      })
      allHttpField.value = allHttpField.value.filter((item: { isDisplay: number; id: number | string }) => {
        return !ids?.includes(+item.id)
      })
    }
  }

  // 提交
  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        currentItems.value.data = formData
        if (formData.module == '0') {
          currentItems.value.data['url'] = 'siteSearchApi'
        } else if (formData.indexType == 24) {
          currentItems.value.data['url'] = 'getTabelListAPI'
        } else {
          currentItems.value.data['url'] = 'searchBySqlApi'
        }
        emit('on-reflash', currentItems.value)
        ElMessage({ message: '操作成功', type: 'success' })
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
    emit('on-closeEvent', false)
  }

  const initData = async () => {
    visible.value = props.showpage
    modes.value = props.mode
    titles.value = props.title
    currentItems.value = JSON.parse(JSON.stringify(props.currentItem))
    Object.keys(currentItems?.value?.data)?.forEach((td) => {
      if (td == 'siteSessionId') {
        formData.siteSessionId = currentItems?.value?.data['siteSessionId'] || ''
      } else {
        // @ts-ignore
        formData[td] = currentItems?.value?.data[td]
      }
    })
    getAllSite()
    setTimeout(() => {
      if (currentItems?.value?.data?.siteApiId) {
        formData.siteApiId = currentItems?.value?.data?.siteApiId || ''
      }
    }, 0)
  }

  watch(
    () => formData.module,
    (newVal) => {
      formData.indexType = 1
      resetForm(formRef.value)
      formData.module = newVal
      formData.indexType = currentItems.value?.data?.indexType || 1
      setTimeout(() => {
        formRef.value!.clearValidate()
      }, 0)
      if (newVal == '1') {
        getHttpField()
        formData.indexType = 22
      } else if (newVal == '2') {
        getHttpField()
      }
    },
    { immediate: true }
  )

  watch(
    () => formData.indexType,
    (newVal, oldVal) => {
      if (newVal != oldVal && formData.indexType !== 22) {
        getHttpField()
      }
      if (formData.indexType == 1 || formData.indexType == 22) {
        formData.orderField = 'requestTimeNs'
      } else {
        if (formData.indexType == 24) {
          formData.searchModel = 'ip'
          formData.topValue = '100'
        }
        formData.orderField = 'startTimeNs'
      }
    }
  )

  watch(
    () => formData.siteSessionId,
    (newVal, oldVal) => {
      if (newVal != oldVal) {
        getSonSite()
        getHttpField()
        formData.siteApiId = undefined
      }
    }
  )

  watch(
    () => formData.siteApiId,
    (newVal, oldVal) => {
      if (newVal != oldVal) {
        getHttpField()
      }
    }
  )

  const getAllDefaultFields = async () => {
    const res = await getAllDefaultFieldsApi()
    allDefaultField.value = res.data
  }

  const defaultProps = {
    children: 'children',
    label: 'name',
    value: 'id',
  }

  // 右侧下拉树
  const flowProbesList = ref()
  const flowProbesValue = ref(['0'])
  const getFlowProbes = async () => {
    const { data } = await getFlowProbesAPI()
    flowProbesList.value = []
    flowProbesList.value.push(data.data)
    flowProbesValue.value = []
    getAllID(flowProbesList.value)
    formData.flowProbeIds = flowProbesValue.value as string[]
  }

  const getAllID = (data: any) => {
    data?.forEach((item: any) => {
      flowProbesValue.value.push(item.id)
      if (item.children && item.children.length > 0) {
        getAllID(item.children)
      }
    })
  }

  onMounted(async () => {
    getAllDefaultFields()
    await initData()
    await getFlowProbes()
  })
</script>

<template>
  <div class="edit-window">
    <el-dialog
      v-model="visible"
      :before-close="handleClose"
      :style="{ minHeight: '300px' }"
      :title="titles"
      width="800px"
    >
      <el-form ref="formRef" label-position="right" label-width="110px" :model="formData" :rules="rules">
        <el-form-item label="类型：" prop="module">
          <el-radio-group v-model="formData.module">
            <el-radio label="0">站点接口</el-radio>
            <el-radio label="1">字段提取</el-radio>
            <el-radio label="2">自定义</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="窗口名称：" prop="name">
          <el-input v-model="formData.name" placeholder="请输入窗口名称" />
        </el-form-item>
        <el-row v-if="formData.module == '0'" :gutter="20">
          <el-col :span="24">
            <el-form-item label="站点：" prop="siteSessionId">
              <el-select v-model="formData.siteSessionId" style="width: 50%">
                <el-option v-for="item in siteOption" :key="item.id" :label="item.siteName" :value="item.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="子站点：" prop="siteApiId">
              <el-select v-model="formData.siteApiId" style="width: 50%">
                <el-option v-for="item in siteApiIdOption" :key="item.id" :label="item.apiName" :value="item.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item v-if="formData.module == '2'" label="数据类型：" prop="indexType" style="width: 100%">
          <el-select v-model="formData.indexType" style="width: 100%">
            <el-option v-for="item in indexTypeOption" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="formData.module == '2' && formData.indexType == 24"
          class="my-form-top"
          label="选择链路："
          prop="indexType"
        >
          <el-form-item style="width: 100%">
            <el-tree-select
              v-model="formData.flowProbeIds"
              :data="flowProbesList"
              default-expand-all
              :expand-on-click-node="false"
              multiple
              node-key="id"
              :props="defaultProps"
              show-checkbox
              style="width: 100%"
            />
          </el-form-item>
        </el-form-item>
        <el-form-item v-if="formData.module == '2' && formData.indexType == 24" label="TOP：" prop="indexType">
          <el-select v-model="formData.topValue" style="width: 100%">
            <el-option v-for="item in topOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item v-if="formData.module == '2' && formData.indexType == 24" prop="searchModel">
          <template #label>
            <el-select v-model="formData.searchModel" style="width: 140px">
              <el-option v-for="item in SearchOption" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </template>
          <el-row style="width: 100%">
            <el-col v-if="formData.searchModel == 'ip'" :span="24">
              <el-form-item prop="ip">
                <el-input
                  v-model="formData.ip"
                  placeholder="请输入IP地址，支持IP格式和掩码格式，多个检索条件使用英文逗号分开"
                />
              </el-form-item>
            </el-col>
            <el-col v-if="formData.searchModel == 'clientIp'" :span="24">
              <el-form-item prop="clientIp">
                <el-input
                  v-model="formData.clientIp"
                  placeholder="源IP地址，支持IP格式和掩码格式，多个使用英文逗号分开"
                />
              </el-form-item>
            </el-col>
            <el-col v-if="formData.searchModel == 'serverIp'" :span="24" style="margin-bottom: 10px">
              <el-form-item prop="serverIp">
                <el-input
                  v-model="formData.serverIp"
                  placeholder="目的IP地址，支持IP格式和掩码格式，多个使用英文逗号分开"
                />
              </el-form-item>
            </el-col>
            <el-col v-if="formData.searchModel == 'serverIp'" :span="24">
              <el-form-item prop="serverPort">
                <el-input v-model="formData.serverPort" placeholder="目的端口，可为空，为空代表所有" />
              </el-form-item>
            </el-col>
            <el-col v-if="formData.searchModel == 'clientServerIp'" :span="24" style="margin-bottom: 10px">
              <el-form-item prop="clientIp">
                <el-input
                  v-model="formData.clientIp"
                  placeholder="源IP地址，支持IP格式和掩码格式，多个使用英文逗号分开"
                />
              </el-form-item>
            </el-col>
            <el-col v-if="formData.searchModel == 'clientServerIp'" :span="24" style="margin-bottom: 10px">
              <el-form-item prop="serverIp">
                <el-input
                  v-model="formData.serverIp"
                  placeholder="目的IP地址，支持IP格式和掩码格式，多个使用英文逗号分开"
                />
              </el-form-item>
            </el-col>
            <el-col v-if="formData.searchModel == 'clientServerIp'" :span="24">
              <el-form-item prop="serverPort">
                <el-input v-model="formData.serverPort" placeholder="目的端口，可为空，为空代表所有" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form-item>

        <el-form-item v-if="formData.module != '0' && formData.indexType != 24" label="检索语句：" prop="searchSql">
          <el-input v-model="formData.searchSql" placeholder="请输入检索语句" resize="none" :rows="4" type="textarea" />
        </el-form-item>
        <el-form-item v-if="formData.module != '0'" prop="displayFields">
          <template #label>
            <div style="text-align: right">
              <span>展示字段：</span>
              <br />
              <span style="color: #a9acb3; font-size: 13px">(拖动排序)&nbsp;</span>
            </div>
          </template>
          <div class="space">
            <vue-draggable
              animation="300"
              chosen-class="chosenClass"
              ghost-class="ghost"
              group="my-group"
              item-key="id"
              :list="formData.displayFields"
            >
              <template #item="{ element }">
                <el-tag>{{ element.fieldNameCn }}</el-tag>
              </template>
            </vue-draggable>
          </div>
        </el-form-item>
        <el-form-item v-if="formData.module != '0'" label="可选字段：">
          <div class="space">
            <vue-draggable
              animation="300"
              chosen-class="chosenClass"
              ghost-class="ghost"
              group="my-group"
              item-key="id"
              :list="allHttpField"
              :sort="false"
            >
              <template #item="{ element }">
                <el-tag type="info">{{ element.fieldNameCn }}</el-tag>
              </template>
            </vue-draggable>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="submitForm(formRef)">确认</el-button>
          <!-- <el-button @click="handleClose">取消</el-button> -->
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
    div {
      height: calc(100% - 20px);
      min-height: 130px;
    }
  }

  :deep() {
    .el-tree-node__content > .el-tree-node__expand-icon {
      padding: 0;
    }
    .el-tree-node__content > label.el-checkbox {
      margin-right: 0px;
    }
    .my-form-top {
      .el-select-tags-wrapper.has-prefix {
        overflow: hidden;
        height: 31px;
        line-height: 36px;
      }
      .el-select .el-input {
        height: 31px;
      }
    }
    // .el-form-item__content {
    //   position: initial;
    // }

    // .el-form-item.is-error .el-input__wrapper {
    //   box-shadow: 0 0 0 1px var(--el-input-border-color, var(--el-border-color)) inset;
    // }

    // .el-form-item__error {
    //   top: 33px;
    //   left: 4px;
    // }
  }

  :deep(.el-tag) {
    margin-left: 10px;
  }
</style>
