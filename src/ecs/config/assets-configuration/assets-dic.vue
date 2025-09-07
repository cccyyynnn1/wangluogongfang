<script lang="ts">
  export default {
    name: 'AssetsDic', //
  }
</script>

<script setup lang="ts">
  import { formatNstime, formatTime } from '@/utils/time'
  import { getAllAssetDictApi, updateAssetDictApi } from '~/src/api-ecs/custom-field'
  import { AssetLabelDict } from '@/types'
  import { requireRules } from '~/src/utils/rules'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import { FormInstance } from 'ant-design-vue'
  const $baseMessage: any = inject('$baseMessage')

  const queryData = reactive({
    pageNum: 1,
    pageSize: 20,
    searchStr: '',
  })

  const editData = reactive<AssetLabelDict>({
    sourceValue: '',
    transferValue: '',
    remark: '',
    id: undefined,
  })

  // const dictionaryListId = ref('')
  const formRef = ref<FormInstance>()
  const visible = ref(false)
  const selection = ref<AssetLabelDict[]>([])
  const listLoading = ref(false)
  // 表格数据
  const listDate = ref<AssetLabelDict[]>([])
  // 表格数据
  const listTotal = ref(0)
  const allData = ref()

  // 多选项改变
  const setSelectRows = (selections: any[]) => {
    selection.value = selections
  }
  const editDataDictionaryHandle = (data?: AssetLabelDict) => {
    if (data) {
      for (const key in editData) {
        // @ts-ignore
        editData[key] = data[key]
      }
    }

    visible.value = true
  }
  const restEditData = () => {
    for (const key in editData) {
      if (['sourceValue', 'transferValue', 'remark'].includes(key)) {
        // @ts-ignore
        editData[key] = ''
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
        const { code, msg } = await updateAssetDictApi(editData)
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        handleSearch()
        visible.value = false
      }
    })
  }

  let order: undefined | 'descending' | 'ascending' = undefined
  type propTypt = undefined | 'labelCount' | 'sourceValueLength' | 'createTime'
  let prop: propTypt = undefined
  const sortChangeEvent = (e: any) => {
    order = e.order
    prop = e.prop
    sortChange()
  }

  const sortChange = () => {
    listDate.value = []
    allData.value.sort(function (a: any, b: any) {
      return order == 'ascending' ? a[prop!] - b[prop!] : b[prop!] - a[prop!]
    })
    const count =
      (queryData.pageNum - 1) * queryData.pageSize + queryData.pageSize < allData.value.length
        ? (queryData.pageNum - 1) * queryData.pageSize + queryData.pageSize
        : allData.value.length
    for (let index = (queryData.pageNum - 1) * queryData.pageSize; index < count; index++) {
      listDate.value.push(allData.value[index])
    }
  }
  // 导出
  const handleDownloadExcel = () => {
    const tHeader = ['原始标签', '转义标签', '引用次数', '原始字符长度', '备注', '创建时间']
    const filterVal = ['sourceValue', 'transferValue', 'labelCount', 'sourceValueLength', 'remark', 'createTime']
    try {
      import('@/utils/excel').then((excel) => {
        const list: any[] = selection.value.length ? selection.value : listDate.value
        const data = list.map((item) =>
          filterVal.map((valKey) => {
            const val = item[valKey]
            return valKey === 'createTime' ? formatNstime(val, false) : val
          })
        )
        excel.export_json_to_excel({
          header: tHeader,
          data,
          filename: `资产字典-${formatTime(new Date().getTime())}`,
          autoWidth: true,
          bookType: 'xlsx',
        })
      })
    } catch (error) {
      console.log(error)
    }
  }

  // 查询
  const handleSearch = async () => {
    listLoading.value = true
    try {
      listDate.value = []
      const res = await getAllAssetDictApi(queryData)
      const count =
        (queryData.pageNum - 1) * queryData.pageSize + queryData.pageSize < res.data.length
          ? (queryData.pageNum - 1) * queryData.pageSize + queryData.pageSize
          : res.data.length
      for (let index = (queryData.pageNum - 1) * queryData.pageSize; index < count; index++) {
        listDate.value.push(res.data[index])
      }
      listTotal.value = res.data.length
      allData.value = res.data
      sortChange()
    } finally {
      listLoading.value = false
    }
  }

  onMounted(() => {
    handleSearch()
  })
</script>

<template>
  <div class="assets-dic-container">
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>资产字典</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12">
        <div style="display: flex">
          <el-input
            v-model="queryData.searchStr"
            clearable
            placeholder="模糊检索"
            style="margin: 0 12px 0 0; width: 350px"
          />
          <el-button type="primary" @click="handleSearch">检索</el-button>
          <el-button :loading="listLoading" type="primary" @click="handleDownloadExcel">导出</el-button>
        </div>
      </vab-query-form-right-panel>
    </vab-query-form>
    <el-row :gutter="20">
      <!-- 字典列表 -->
      <el-col :span="24">
        <el-table
          v-loading="listLoading"
          :data="listDate"
          row-key="uid"
          @selection-change="setSelectRows"
          @sort-change="sortChangeEvent"
        >
          <el-table-column fixed="left" show-overflow-tooltip type="selection" width="55" />
          <el-table-column label="原始标签" prop="sourceValue" show-overflow-tooltip />
          <el-table-column label="转义标签" prop="transferValue" show-overflow-tooltip />
          <el-table-column label="引用次数" prop="labelCount" show-overflow-tooltip sortable />
          <el-table-column label="原始字符长度" prop="sourceValueLength" show-overflow-tooltip sortable />
          <el-table-column label="备注" prop="remark" show-overflow-tooltip />
          <el-table-column
            :formatter="({ createTime }) => formatNstime(createTime, false)"
            label="创建时间"
            prop="createTime"
            show-overflow-tooltip
            sortable
            width="200px"
          />
          <el-table-column label="操作" width="80">
            <template #default="{ row }">
              <el-button size="small" @click="editDataDictionaryHandle(row)">编辑</el-button>
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
          :page-sizes="[10, 20, 50]"
          :total="listTotal"
          @current-change="sortChange"
          @size-change="sortChange"
        />
      </el-col>
    </el-row>
    <vab-dialog
      v-model="visible"
      align-center
      :close-on-click-modal="false"
      destroy-on-close
      :title="editData?.id ? '编辑标签名称' : '新增标签名称'"
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
              transferValue: requireRules,
            }"
          >
            <el-form-item label="转义标签：" prop="transferValue">
              <el-input v-model="editData.transferValue" />
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
</template>

<style scoped lang="scss">
  .assets-dic-container {
    :deep() {
      // .el-dropdown {
      //   margin-top: -24px !important;
      // }
      .el-scrollbar {
        height: calc(100vh - 170px);
      }
    }
  }
</style>
