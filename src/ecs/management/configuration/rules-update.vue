<script lang="ts">
  export default {
    name: 'RulesUpdate', // 规则更新
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
  import { UploadUserFile } from 'element-plus'
  const uploadIcon = require('@/assets/konwledge/upload.svg')
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const router = useRouter()
  const route = useRoute()
  const queryData = reactive({
    pageNum: 1,
    pageSize: 20,
    enable: '',
    dictType: '',
    dictName: '',
  })

  const editData = reactive<{ file: UploadUserFile[]; remark: string }>({
    file: [],
    remark: '',
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
    $baseConfirm('你确定要删除当前项吗?', null, async () => {
      try {
        const ids = data.map((i) => i.id!)
        const { code } = await deleteDictTypePageApi({ ids, deleteAll: isAll })
        $baseMessage('删除成功', 'success', 'vab-hey-message-success')
        queryData.pageNum = 1
        handleSearch()
      } catch (error) {
        console.log(error)
      }
    })
  }
  const restEditData = () => {
    for (const key in editData) {
      // if (['dictName', 'dictType', 'remark'].includes(key)) {
      //   // @ts-ignore
      //   editData[key] = ''
      // } else if (key === 'enable') {
      //   editData[key] = 1
      // } else {
      //   // @ts-ignore
      //   editData[key] = undefined
      // }
    }
  }
  const submitForm = () => {
    if (!formRef.value) return
    formRef.value.validate().then(async (isValid) => {
      if (isValid) {
        // const { code, msg } = await editDictTypePageApi(editData)
        // $baseMessage(msg, 'success', 'vab-hey-message-success')
        // handleSearch()
        // visible.value = false
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
  const handleShowDialog = () => {
    visible.value = true
  }

  const handleSizeChange = () => {
    queryData.pageNum = 1
    handleSearch()
  }
  const curIndex = computed(() => (queryData.pageNum - 1) * queryData.pageSize + 1)
</script>

<template>
  <div class="rules-update-container">
    <div v-if="!dictionaryListId" class="warp">
      <vab-query-form>
        <vab-query-form-left-panel :span="12">
          <h3>规则更新</h3>
        </vab-query-form-left-panel>
        <vab-query-form-right-panel :span="12">
          <el-button color="#6954F0" :icon="Plus" @click="handleShowDialog">上传</el-button>
        </vab-query-form-right-panel>
      </vab-query-form>
      <el-table v-loading="listLoading" :data="listDate" row-key="id" @selection-change="setSelectRows">
        <!-- <el-table-column fixed="left" show-overflow-tooltip type="selection" width="55" /> -->
        <el-table-column align="center" label="序号" width="55">
          <template #default="{ $index }">
            {{ curIndex + $index }}
          </template>
        </el-table-column>
        <el-table-column label="上传人" prop="dictName" show-overflow-tooltip width="200px">
          <template #header="{ column }">
            <table-filter-tool
              v-model:filter-value="queryData.dictName"
              :column="column"
              filter-type="text"
              :tools="['filter']"
            />
          </template>
        </el-table-column>
        <el-table-column label="文件名称" prop="dictType" show-overflow-tooltip>
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
          :formatter="({ createTime }) => formatNstime(createTime, false)"
          label="创建时间"
          prop="createTime"
          show-overflow-tooltip
          width="200px"
        />
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button size="small" @click="toDictionaryList(row)">版本描述</el-button>
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
        @size-change="handleSizeChange"
      />
      <vab-dialog
        v-model="visible"
        align-center
        :close-on-click-modal="false"
        destroy-on-close
        title="规则上传"
        width="850px"
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
              <el-form-item label="字典名称：" prop="file">
                <el-upload
                  ref="upload"
                  v-model:file-list="editData.file"
                  accept=".zip"
                  action="/v3/ecsPlatform/knowledge/upload"
                  auto-upload
                  class="upload-demo"
                  drag
                >
                  <!--        :http-request="hanldeFileUpload"
        :limit="1"
        :multiple="false"
        :on-error="() => (loading = false)"
        :on-exceed="handleExceed"
        :on-success="handleUploadSuccess"
        :show-file-list="false" -->
                  <el-image :src="uploadIcon" style="width: 48px; height: 44px; margin-bottom: 15px" />
                  <div class="el-upload__text">
                    <span style="color: #4a4759; line-height: 20px">将文件拖到此处，或</span>
                    <em style="line-height: 20px">点击上传</em>
                    <br />
                    <div style="font-size: 13px; color: #bbbdbf; line-height: 14px; margin-top: 9px">
                      只能上传zip文件，且不超过2000M
                    </div>
                  </div>
                </el-upload>
              </el-form-item>
              <el-form-item label="版本描述：" prop="remark">
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
    <!-- <data-rules-update-list v-else :list-id="dictionaryListId" /> -->
  </div>
</template>

<style scoped lang="scss">
  .rules-update-container {
    :deep() {
      .el-scrollbar {
        height: calc(100vh - 170px);
      }
    }
  }
</style>
