<script setup lang="ts">
  import { FormInstance, TableInstance, UploadProps, UploadUserFile } from 'element-plus'

  import { Plus } from '@element-plus/icons-vue'

  import { requireRules } from '~/src/utils/rules'

  import { loadBalanceAddApi, loadBalanceModelListApi } from '@/api-ecs/equipment'

  import { Search } from '@element-plus/icons-vue'

  import { useTableCopy } from '@/utils'
  import AesEncryptCBC from '~/src/utils/crypto'

  const props = defineProps<{
    mode: string
    showPage: boolean
    currentRow: object
  }>()

  const loadBalanceTableRef = ref<TableInstance>() // 型号表单实例

  const uploadInputValue = ref('')

  const formRef = ref<FormInstance>() // 表单实例

  const visible = ref(false) // 弹框显隐

  const loadBalanceConfigVisible = ref(false) // 型号popover弹框显隐

  const title = ref('添加') // 标题

  let loadBalanceListDate = reactive<object[]>([]) // 表格数据

  const layout = ref('total, sizes, prev, pager, next, jumper')

  const loadBalanceTotals = ref(0) // 总条数

  const listLoading = ref(false) // 是否加载

  // 获取型号列表参数
  const loadBalancePayload = reactive({
    page: 1,
    limit: 10,
    query: {
      keyWord: undefined,
      flowDeviceId: undefined,
    },
    url: 'loadBalanceConfig/page',
  })

  // 表单数据
  const formData = reactive({
    dialogTitle: 'add', // 添加链路时候的固定值
    moduleType: 'loadBalance',
    modelId: undefined,
    modelIdStr: undefined,
    mode: '2',
    policyContentFile: undefined, // 文件
    enabled: true,
    configName: undefined,
    note: undefined,
    loadBalancePort: '22',
    retryCount: '3',
    retryWait: '60',
  })

  // 表单数据校验规则
  const rules = reactive({
    modelIdStr: requireRules,
    recmodeordTopic: requireRules,
    policyContentFile: requireRules,
    enabled: requireRules,
  })

  const emit = defineEmits<{
    (e: 'on-close-event'): void
    (e: 'on-reflash'): void
  }>()

  const loadBalanceIndex = computed(() => (loadBalancePayload.page - 1) * loadBalancePayload.limit + 1)

  // 关闭回调
  const handleClose = () => {
    emit('on-close-event')
  }

  // 提交表单
  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        try {
          listLoading.value = true
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
                h(
                  'div',
                  { style: 'margin: 15px 0 5px 0;color:#55585b' },
                  '请输入敏感操作密码:(通过验证后方可进行操作)'
                ),
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
                  if (props.mode == 'add') {
                    const { msg } = await loadBalanceAddApi(
                      {
                        data: JSON.stringify({ ...formData, policyContentFile: undefined }),
                        // @ts-ignore
                        policyContentFile: formData.policyContentFile[0].raw,
                      },
                      { password: AesEncryptCBC(password.value) }
                    )
                    ElMessage({ message: msg, type: 'success' })
                  }
                  done()
                } catch (error) {
                  instance.confirmButtonLoading = false
                }
              } else {
                done()
              }
            },
          }).catch(() => {})

          //  else {listLoading
          //   // @ts-ignore
          //   const id = props.currentRow.id
          //   // @ts-ignore
          //   const { message } = await loadBalanceUpdataApi(id, {
          //     ...formData,
          //     id,
          //   })
          //   ElMessage({ message: message, type: 'success' })
          // }
          emit('on-reflash')
          handleClose()
        } finally {
          listLoading.value = false
        }
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  // 初始化数据
  const initData = () => {
    title.value = props.mode == 'add' ? '添加' : '编辑'
    visible.value = props.showPage
    // if (props.mode == 'edit') {
    //   const data = JSON.parse(JSON.stringify(props.currentRow))
    //   Object.keys(formData).forEach((item: string) => {
    //     if (item !== 'dialogTitle' && item !== 'moduleType') {
    //       // @ts-ignore
    //       formData[item] = data[item]
    //     }
    //   })
    // }
  }

  onMounted(() => {
    initData()
  })

  // 型号页容量改变
  const handleloadBalanceSizeChange = (val: number) => {
    loadBalancePayload.limit = val
    getloadBalanceConfigList()
  }

  // 型号页面改变
  const handleloadBalanceChange = (val: number) => {
    loadBalancePayload.page = val
    getloadBalanceConfigList()
  }

  // 打开popover
  const handleClosePop = () => {
    loadBalanceTableRef.value?.clearSelection()
    loadBalanceConfigVisible.value = false
  }

  // 获取kafka服务列表
  const getloadBalanceConfigList = async () => {
    try {
      listLoading.value = true
      const { data } = await loadBalanceModelListApi({
        ...loadBalancePayload,
        // @ts-ignore
        query: JSON.stringify(loadBalancePayload.query),
      })
      loadBalanceListDate = data.list
      loadBalanceTotals.value = data.list
    } finally {
      listLoading.value = false
    }
  }

  // 型号选择
  const handleCurrentChange = (val: any) => {
    if (formData.modelId == val.id) {
      formData.modelId = undefined
      formData.modelIdStr = undefined
      loadBalanceTableRef.value!.setCurrentRow(false)
    } else {
      formData.modelId = val.id
      formData.modelIdStr = val.name
    }
  }

  // 打开
  const openPopEvent = async () => {
    loadBalanceConfigVisible.value = true
    await getloadBalanceConfigList()
    const rows = loadBalanceListDate.filter((item: any) => {
      return item.id == formData.modelId
    })
    loadBalanceTableRef.value!.setCurrentRow(rows[0])
  }

  // 大小超过
  const handleExceed: UploadProps['onExceed'] = (files) => {
    console.log('handleExceed')
  }

  // 上传成功
  const handleChange: UploadProps['onChange'] = (uploadFile, uploadFiles) => {
    uploadInputValue.value = uploadFile.name
  }
</script>

<script lang="ts">
  export default {
    name: 'LoadBalancingAdd',
  }
</script>
<template>
  <div class="load-balancing-add">
    <el-dialog v-model="visible" :before-close="handleClose" :title="title" width="800">
      <el-form ref="formRef" label-width="150px" :model="formData" :rules="rules">
        <el-row>
          <el-col :span="12">
            <el-form-item label="型号" prop="modelIdStr">
              <el-popover placement="bottom" trigger="click" :visible="loadBalanceConfigVisible" width="40%">
                <vab-query-form>
                  <vab-query-form-left-panel :span="4">
                    <el-button type="danger" @click="handleClosePop">关闭</el-button>
                  </vab-query-form-left-panel>
                  <vab-query-form-right-panel :span="20">
                    <el-input
                      v-model="loadBalancePayload.query.keyWord"
                      class="input-with-select"
                      :disabled="listLoading"
                      placeholder="请输入名称"
                      @keyup.enter="getloadBalanceConfigList"
                    >
                      <template #append>
                        <el-button :icon="Search" @click="getloadBalanceConfigList" />
                      </template>
                    </el-input>
                  </vab-query-form-right-panel>
                </vab-query-form>
                <el-table
                  ref="loadBalanceTableRef"
                  v-loading="listLoading"
                  :border="true"
                  :data="loadBalanceListDate"
                  highlight-current-row
                  row-key="id"
                  @cell-click="handleCurrentChange"
                  @cell-contextmenu="useTableCopy"
                >
                  <el-table-column :align="'center'" label="序号" width="55">
                    <template #default="{ $index }">
                      {{ loadBalanceIndex + $index }}
                    </template>
                  </el-table-column>
                  <el-table-column label="型号名称" property="name" />
                </el-table>
                <el-pagination
                  background
                  :current-page="loadBalancePayload.page"
                  :layout="layout"
                  :page-size="loadBalancePayload.limit"
                  :page-sizes="[10, 20, 50, 100]"
                  :total="loadBalanceIndex"
                  @current-change="handleloadBalanceChange"
                  @size-change="handleloadBalanceSizeChange"
                />
                <template #reference>
                  <el-input v-model="formData.modelIdStr" :disabled="loadBalanceConfigVisible" @click="openPopEvent" />
                </template>
              </el-popover>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="获取方式" prop="mode">
              <el-radio v-model="formData.mode" label="1">直线</el-radio>
              <el-radio v-model="formData.mode" label="2">离线</el-radio>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="配置文件" prop="policyContentFile">
              <el-upload
                v-model:file-list="formData.policyContentFile"
                :auto-upload="false"
                class="upload-demo"
                :limit="1"
                :name="'1'"
                :on-change="handleChange"
                :on-exceed="handleExceed"
              >
                <el-input v-model="uploadInputValue">
                  <template #append>
                    <el-icon><Plus /></el-icon>
                  </template>
                </el-input>
              </el-upload>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="enabled">
              <el-radio v-model="formData.enabled" :label="true">启用</el-radio>
              <el-radio v-model="formData.enabled" :label="false">停用</el-radio>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="名称" prop="configName">
              <el-input v-model="formData.configName" maxlength="128" show-word-limit />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="备注" prop="note">
              <el-input v-model="formData.note" maxlength="512" show-word-limit />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <el-row>
        <el-col :offset="11" :span="2" style="text-align: right">
          <el-button :loading="listLoading" type="primary" @click="submitForm(formRef)">保存</el-button>
        </el-col>
      </el-row>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  :deep() {
    .el-scrollbar {
      height: 200px !important;
    }
    .el-upload-list.el-upload-list--text {
      display: none;
    }
  }
</style>
