<script setup lang="ts">
  import type { FormInstance } from 'element-plus'

  import { getDataSourceApi, updateNetworkPartitionSetApi } from '~/src/api-ecs/assets'

  import { requireRules } from '~/src/utils/rules'

  import VabDialog from '@/plugins/VabDialog/index.vue'

  import { formatNstime } from '@/utils/time'
  import { getDictTypePageListApi, editDictTypePageListApi, deleteDictTypePageListApi } from '@/api-ecs/data-dictionary'
  import { DictTypePageListItem } from '@/types'
  import { Plus } from '@element-plus/icons-vue'
  import AesEncryptCBC from '~/src/utils/crypto'
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')

  type itemType = {
    name: string
    rule: string
    dataSource: string
    region: string
    netType: number
  }

  const props = defineProps<{
    showPage: boolean
    currentItem: any
    modes: string
    options: any
    cityOptions: any
  }>()

  const option = ref()

  const emit = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
    (e: 'on-reflash'): void
  }>()

  const visible = ref(false)

  const visible1 = ref(false)

  // 表单数据
  const formData = reactive<itemType>({
    name: '',
    rule: '',
    dataSource: '',
    region: '',
    netType: 0,
  })

  // 表单数据校验
  const rules = reactive({
    name: requireRules,
    rule: requireRules,
    dataSource: requireRules,
    region: requireRules,
    netType: requireRules,
  })

  const formRef = ref<FormInstance>()

  const mode = ref('')

  // 提交
  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        const { msg } = await updateNetworkPartitionSetApi({
          ...formData,
          id: mode.value == 'edit' ? props.currentItem.id : null,
        })
        ElMessage({ message: msg, type: 'success' })
        emit('on-reflash')
        handleClose()
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  // 关闭
  const handleClose = () => {
    emit('on-closeEvent', false)
  }

  // 初始化数据
  const initData = () => {
    visible.value = props.showPage
    mode.value = props.modes
    option.value = props.options
    formData.region = props.cityOptions[0]?.en || ''
    if (mode.value === 'edit') {
      const data = JSON.parse(JSON.stringify(props.currentItem))
      Object.keys(formData).forEach((item) => {
        // @ts-ignore
        formData[item as keyof typeof formData] = data[item]
      })
    }
  }

  const handleCloseVab = () => {
    visible1.value = false
  }

  const handelOpen = () => {
    handleSearch()
    visible1.value = true
  }

  const getDataSource = async () => {
    const { data } = await getDataSourceApi()
    option.value = data.datasource || []
  }

  onMounted(() => {
    initData()
  })

  const queryData = reactive({
    pageNum: 1,
    pageSize: 20,
  })
  const editData = reactive<DictTypePageListItem>({
    dictLabel: '',
    dictValue: '',
    enable: 1,
    remark: '',
    id: undefined,
  })
  const formRef1 = ref<FormInstance>()
  const visible2 = ref(false)
  const selection = ref<DictTypePageListItem[]>([])
  const listLoading = ref(false)
  // 表格数据
  const listDate = ref<DictTypePageListItem[]>([])
  // 表格数据
  const listTotal = ref(0)
  const restEditData = () => {
    for (const key in editData) {
      if (['dictName', 'dictType', 'remark'].includes(key)) {
        // @ts-ignore
        editData[key] = ''
      } else if (key === 'enable') {
        editData[key] = 1
      } else {
        // @ts-ignore
        editData[key] = undefined
      }
    }
  }
  const submitForm1 = () => {
    if (!formRef1.value) return
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
            if (!formRef.value) return
            formRef1.value?.validate().then(async (isValid) => {
              if (isValid) {
                const { code, msg } = await editDictTypePageListApi(
                  { ...editData, dictTypeId: '6' },
                  { password: AesEncryptCBC(password.value) }
                )
                $baseMessage(msg, 'success', 'vab-hey-message-success')
                handleSearch()
                getDataSource()
                visible2.value = false
              }
            })
            done()
          } catch (error) {
            instance.confirmButtonLoading = false
          }
        } else {
          done()
        }
      },
    }).catch(() => {
      visible2.value = false
    })
  }
  // 查询
  const handleSearch = async () => {
    listLoading.value = true
    try {
      const {
        data: { records, total },
      } = await getDictTypePageListApi({ ...queryData, dictTypeId: '6' })
      listDate.value = Array.isArray(records) ? records : []
      listTotal.value = total
    } finally {
      listLoading.value = false
    }
  }
  // 多选项改变
  const setSelectRows = (selections: any[]) => {
    selection.value = selections
  }
  const editDataDictionaryHandle = (data?: DictTypePageListItem) => {
    if (data) {
      for (const key in data) {
        // @ts-ignore
        editData[key] = data[key]
      }
    }
    visible2.value = true
  }
  const deleteDataDictionaryHandle = (data: DictTypePageListItem[], isAll = false) => {
    $baseConfirm('你确定要删除当前项吗?', null, async () => {
      try {
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
                const ids = data.map((i) => i.id!)
                const { code } = await deleteDictTypePageListApi(
                  { ids, deleteAll: isAll, dictTypeId: '6' },
                  { password: AesEncryptCBC(password.value) }
                )
                $baseMessage('删除成功', 'success', 'vab-hey-message-success')
                queryData.pageNum = 1
                handleSearch()
                getDataSource()
                done()
              } catch (error) {
                instance.confirmButtonLoading = false
              }
            } else {
              done()
            }
          },
        }).catch(() => {
          visible.value = false
        })
      } catch (error) {
        console.log(error)
      }
    })
  }

  const handleReturn = () => {
    visible1.value = false
  }
</script>

<script lang="ts">
  export default {
    name: 'UpdateNetwork', // 新增或修改网络分区
  }
</script>
<template>
  <div class="update-network">
    <el-dialog v-model="visible" :before-close="handleClose" :title="mode == 'add' ? '添加' : '编辑'" width="720px">
      <el-form ref="formRef" label-width="120px" :model="formData" :rules="rules">
        <el-form-item label="数据中心名称" prop="dataSource">
          <el-select v-model="formData.dataSource" class="select">
            <el-option v-for="item in option" :key="item.dictValue" :label="item.dictLabel" :value="item.dictValue" />
          </el-select>
          <el-button style="margin-left: 15px" type="primary" @click="handelOpen">配置</el-button>
          <!-- <el-input v-model="formData.dataSource" /> -->
        </el-form-item>
        <el-form-item label="分区名称" prop="name">
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item label="分区规则" prop="rule">
          <el-input v-model="formData.rule" />
          <div class="tips">
            支持IPv4或IPv6，但IPv4和IPv6不能同时存在；支持单个IP地址（192.168.0.1）；
            <br />
            支持IP地址范围（192.168.1.5-192.168.1.9）；支持IP地址掩码（192.168.0.0/16）；
            <br />
            多个用英文逗号分隔。
          </div>
        </el-form-item>
        <el-form-item label="资产属性">
          <el-radio-group v-model="formData.netType" class="ml-4">
            <el-radio :label="0">内网资产</el-radio>
            <el-radio :label="1">外网资产</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="分区所属地" prop="region">
          <el-select v-model="formData.region" style="width: 100%">
            <el-option v-for="item in cityOptions || []" :key="item.en" :label="item.cn" :value="item.en" />
          </el-select>
          <!-- <el-input v-model="formData.dataSource" /> -->
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="submitForm(formRef)">确认</el-button>
          <el-button @click="handleClose">取消</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
  <vab-dialog v-model="visible1" destroy-on-close title="数据中心管理" width="1080px" @close="handleCloseVab">
    <div class="dictionary-list-container">
      <el-row :gutter="20">
        <!-- 字典操作 -->
        <el-col :span="24">
          <el-button
            :icon="Plus"
            style="margin-bottom: 0 !important; margin-right: 10px"
            type="primary"
            @click="editDataDictionaryHandle()"
          >
            添加
          </el-button>
          <el-dropdown>
            <span class="el-dropdown-link">
              <el-button type="danger">
                批量删除
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="deleteDataDictionaryHandle(selection)">删除选中</el-dropdown-item>
                <el-dropdown-item @click="deleteDataDictionaryHandle(selection, true)">删除所有</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-button style="float: right" type="primary" @click="handleReturn">返回</el-button>
        </el-col>
        <!-- 字典列表 -->
        <el-col :span="24" style="margin-top: 20px">
          <el-table
            v-loading="listLoading"
            :border="true"
            :data="listDate"
            row-key="id"
            @selection-change="setSelectRows"
          >
            <el-table-column fixed="left" show-overflow-tooltip type="selection" width="55" />
            <el-table-column :align="'center'" label="字典标签" prop="dictLabel" show-overflow-tooltip />
            <el-table-column :align="'center'" label="字典键值" prop="dictValue" show-overflow-tooltip />
            <el-table-column
              :align="'center'"
              :formatter="(row) => (row.isDefault == 1 ? '内置' : '自定义')"
              label="字典键值"
              prop="isDefault"
              show-overflow-tooltip
            />
            <el-table-column :align="'center'" label="状态" prop="enable">
              <template #default="{ row }">
                <el-button plain size="small" :type="row.enable == 1 ? 'primary' : 'info'">
                  {{ row.enable == 1 ? '正常' : '停用' }}
                </el-button>
              </template>
            </el-table-column>
            <el-table-column :align="'center'" label="备注" prop="remark" show-overflow-tooltip />
            <el-table-column
              :align="'center'"
              :formatter="({ createTime }) => formatNstime(createTime, false)"
              label="创建时间"
              prop="createTime"
              show-overflow-tooltip
              width="200px"
            />
            <el-table-column :align="'center'" label="操作" width="150">
              <template #default="{ row }">
                <el-button size="small" @click="editDataDictionaryHandle(row)">编辑</el-button>
                <el-button v-if="row.isDefault != 1" size="small" @click="deleteDataDictionaryHandle([row])">
                  删除
                </el-button>
              </template>
            </el-table-column>
            <template #empty>
              <el-empty class="vab-data-empty" description="暂无数据" />
            </template>
          </el-table>
          <el-pagination
            v-model:current-page="queryData.pageNum"
            v-model:page-size="queryData.pageSize"
            background
            class="known_pagination"
            layout="total, sizes, prev, pager, next, jumper"
            :page-sizes="[10, 20, 30]"
            :total="listTotal"
            @current-change="handleSearch"
            @size-change="handleSearch"
          />
        </el-col>
      </el-row>
      <vab-dialog
        v-model="visible2"
        align-center
        :close-on-click-modal="false"
        destroy-on-close
        :title="editData?.id ? '编辑标签' : '新增标签'"
        width="800px"
        @close="restEditData"
      >
        <el-row>
          <el-col :offset="1" :span="20">
            <el-form
              ref="formRef1"
              label-position="right"
              label-width="150px"
              :model="editData"
              :rules="{
                dictLabel: requireRules,
                dictValue: requireRules,
              }"
            >
              <el-form-item label="字典标签：" prop="dictLabel">
                <el-input v-model="editData.dictLabel" />
              </el-form-item>
              <!-- <el-form-item label="字典键值：" prop="dictValue">
                <el-input v-model="editData.dictValue" />
              </el-form-item> -->
              <el-form-item label="字典状态：" prop="enable">
                <el-select v-model="editData.enable">
                  <el-option label="正常" :value="1" />
                  <el-option label="停用" :value="0" />
                </el-select>
              </el-form-item>
              <el-form-item label="备注：" prop="remark">
                <el-input v-model="editData.remark" :autosize="{ minRows: 5 }" resize="none" type="textarea" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="submitForm1">保存</el-button>
              </el-form-item>
            </el-form>
          </el-col>
        </el-row>
      </vab-dialog>
    </div>
  </vab-dialog>
</template>

<style scoped lang="scss">
  .tips {
    font-size: 12px;
    color: #888888;
    line-height: 20px;
  }

  .select {
    width: calc(100% - 80px);
  }

  .add {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 20px;
  }

  .dictionary-list-container {
    :deep() {
      .el-scrollbar {
        height: calc(100vh - 500px) !important;
      }
    }
  }
</style>
