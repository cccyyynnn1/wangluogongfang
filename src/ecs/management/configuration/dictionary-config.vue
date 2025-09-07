<script lang="ts">
  export default {
    name: 'DataDictionaryConfig',
  }
</script>

<script setup lang="ts">
  import { formatNstime } from '@/utils/time'
  import { getDictTypePageApi, editDictTypePageApi, deleteDictTypePageApi } from '@/api-ecs/data-dictionary'
  import { DictTypePageItem, tableSearch } from '@/types'
  import { requireRules } from '~/src/utils/rules'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import { FormInstance } from 'ant-design-vue'
  import { Plus } from '@element-plus/icons-vue'
  import DataDictionaryConfigList from './components/dictionary-config-list.vue'
  import TableFilterTool from '@/components/table-filter-tool.vue'
  import AesEncryptCBC from '~/src/utils/crypto'
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const router = useRouter()
  const route = useRoute()
  const option = [
    {
      label: '全部',
      value: '',
    },
    {
      label: '正常',
      value: 1,
    },
    {
      label: '停用',
      value: 0,
    },
  ]
  const queryData = reactive({
    pageNum: 1,
    pageSize: 20,
    enable: '',
    dictType: '',
    dictName: '',
  })

  const editData = reactive<DictTypePageItem>({
    dictName: '',
    dictType: '',
    enable: 1,
    remark: '',
    id: undefined,
  })

  const dictionaryListId = ref('')
  const formRef = ref<FormInstance>()
  const visible = ref(false)
  const selection = ref<DictTypePageItem[]>([])
  const listLoading = ref(false)
  // 表格数据
  const listDate = ref<DictTypePageItem[]>([])
  // 表格数据
  const listTotal = ref(0)
  provide(tableSearch, () => {
    queryData.pageNum = 1
    handleSearch()
  })
  // 查询
  const handleSearch = async () => {
    listLoading.value = true
    try {
      const {
        data: { records, total },
      } = await getDictTypePageApi(queryData)
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
  const editDataDictionaryHandle = (data?: DictTypePageItem) => {
    if (data) {
      for (const key in data) {
        // @ts-ignore
        editData[key] = data[key]
      }
    }

    visible.value = true
  }
  const deleteDataDictionaryHandle = (data: DictTypePageItem[], isAll = false) => {
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
            $baseConfirm('你确定要删除当前项吗?', null, async () => {
              try {
                const ids = data.map((i) => i.id!)
                const { code } = await deleteDictTypePageApi(
                  { ids, deleteAll: isAll },
                  { password: AesEncryptCBC(password.value) }
                )
                $baseMessage('删除成功', 'success', 'vab-hey-message-success')
                queryData.pageNum = 1
                handleSearch()
              } catch (error) {
                console.log(error)
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
    }).catch(() => {})
  }
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
  const submitForm = () => {
    if (!formRef.value) return
    formRef.value.validate().then(async (isValid) => {
      if (isValid) {
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
                const { code, msg } = await editDictTypePageApi(editData, { password: AesEncryptCBC(password.value) })
                $baseMessage(msg, 'success', 'vab-hey-message-success')
                handleSearch()
                visible.value = false
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
      }
    })
  }
  const toDictionaryList = (data: DictTypePageItem) => {
    router.push({ name: `DataDictionaryConfig`, params: { id: data.id } }) // ->  dataDictionaryConfig/${data.id}
  }
  watch(
    () => queryData.enable,
    () => {
      handleSearch()
    }
  )
  watch(
    () => route.params,
    (data) => {
      if (data.id) {
        dictionaryListId.value = data.id === ':id' ? '' : (data.id as string)
      }
    },
    {
      immediate: true,
    }
  )
  onMounted(() => {
    handleSearch()
  })
</script>

<template>
  <div class="dictionary-config-container">
    <div v-if="!dictionaryListId" class="warp">
      <vab-query-form>
        <vab-query-form-left-panel :span="12">
          <h3>字典配置</h3>
        </vab-query-form-left-panel>
        <vab-query-form-right-panel :span="12">
          <el-button :icon="Plus" style="margin-right: 10px" type="primary" @click="editDataDictionaryHandle()">
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
        </vab-query-form-right-panel>
      </vab-query-form>
      <el-table v-loading="listLoading" :data="listDate" row-key="id" @selection-change="setSelectRows">
        <el-table-column fixed="left" show-overflow-tooltip type="selection" width="55" />
        <el-table-column label="字典名称" prop="dictName" show-overflow-tooltip>
          <template #header="{ column }">
            <table-filter-tool
              v-model:filter-value="queryData.dictName"
              :column="column"
              filter-type="text"
              :tools="['filter']"
            />
          </template>
        </el-table-column>
        <el-table-column label="字典类型" prop="dictType" show-overflow-tooltip>
          <template #header="{ column }">
            <table-filter-tool
              v-model:filter-value="queryData.dictType"
              :column="column"
              filter-type="text"
              :tools="['filter']"
            />
          </template>
        </el-table-column>
        <el-table-column
          :formatter="(row) => (row.isDefault == 1 ? '内置' : '自定义')"
          label="字典键值"
          prop="isDefault"
          show-overflow-tooltip
        />
        <el-table-column label="状态" prop="enable">
          <template #header="{ column }">
            <table-filter-tool
              v-model:filter-value="queryData.enable"
              :column="column"
              :filter-option="option"
              filter-type="select"
              :tools="['filter']"
            />
          </template>
          <template #default="{ row }">
            <el-button plain size="small" :type="row.enable == 1 ? 'primary' : 'info'">
              {{ row.enable == 1 ? '正常' : '停用' }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="备注" prop="remark" show-overflow-tooltip />
        <el-table-column
          :formatter="({ createTime }) => formatNstime(createTime, false)"
          label="创建时间"
          prop="createTime"
          show-overflow-tooltip
          width="200px"
        />
        <el-table-column label="操作" width="220">
          <template #default="{ row }">
            <el-button size="small" @click="editDataDictionaryHandle(row)">编辑</el-button>
            <el-button size="small" @click="toDictionaryList(row)">列表</el-button>
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
      <vab-dialog
        v-model="visible"
        align-center
        :close-on-click-modal="false"
        destroy-on-close
        :title="editData?.id ? '编辑字典类型' : '新增字典类型'"
        width="800px"
        @close="restEditData"
      >
        <el-row>
          <el-col :offset="1" :span="20">
            <el-form
              ref="formRef"
              label-position="right"
              label-width="150px"
              :model="editData"
              :rules="{
                dictName: requireRules,
                dictType: requireRules,
              }"
            >
              <el-form-item label="字典名称：" prop="dictName">
                <el-input v-model="editData.dictName" />
              </el-form-item>
              <el-form-item label="字典类型：" prop="dictType">
                <el-input v-model="editData.dictType" />
              </el-form-item>
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
                <el-button type="primary" @click="submitForm">保存</el-button>
              </el-form-item>
            </el-form>
          </el-col>
        </el-row>
      </vab-dialog>
    </div>
    <data-dictionary-config-list v-else :list-id="dictionaryListId" />
  </div>
</template>

<style scoped lang="scss">
  .dictionary-config-container {
    :deep() {
      .el-scrollbar {
        height: calc(100vh - 170px);
      }
    }
  }
</style>
