<script lang="ts">
  export default {
    name: 'AlertConfigurationWhiteList',
  }
</script>

<script setup lang="ts">
  // @ts-ignore
  import Whitelist from '@/ecs/alert/components/white-list/index.vue'
  import WhiteTime from './white-time.vue'
  import { Delete, Plus, Search } from '@element-plus/icons-vue'
  import { UploadUserFile, UploadProps, UploadInstance, UploadRawFile, genFileId, FormInstance } from 'element-plus'
  import {
    getWarnWhiteApi,
    getDeleteWarnWhiteApi,
    getUpdateWarnWhiteApi,
    exportWarnWhiteTemplateApi,
    importWarnWhiteTemplateApi,
    exportWarnWhiteApi,
  } from '@/api-ecs/alert'
  import { AlertWhiteList, tableSearch } from '@/types'
  import { formatTime } from '@/utils/time'
  import { downloadFile } from '~/src/utils/download'
  import TableFilterTool from '@/components/table-filter-tool.vue'
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const listDate = ref<AlertWhiteList[]>([])
  const listLoading = ref(false)
  const uploadLoading = ref(false)
  const editVisible = ref(false)
  const showDialog = ref(false)
  const uploadVisible = ref(false)
  const multipleSelection = ref<AlertWhiteList[]>([])
  const queryPage = reactive({
    pageSize: 10,
    pageNum: 1,
    searchStr: '',
  })
  const total = ref(0)
  const curItem = ref()
  // 获取表格序号
  const curIndex = computed(() => (queryPage.pageNum - 1) * queryPage.pageSize + 1)

  const curRow = ref<AlertWhiteList>({
    clientIp: '',
    serverIp: '',
    id: undefined,
    ruleId: '',
  })

  const uplodFormRef = ref<FormInstance>()
  const uplodRef = ref<UploadInstance>()

  // 表单数据
  const uploadFile = reactive<{
    file: UploadUserFile[]
  }>({
    file: [],
  })
  // 表单数据校验规则
  const rules = {
    file: [
      {
        required: true,
        message: '请选择上传文件',
        trigger: 'change',
      },
    ],
  }
  const getWhiteListHandle = async () => {
    const {
      data: { records, total: num },
    } = await getWarnWhiteApi(queryPage)
    total.value = num
    listDate.value = records
  }
  // 多选项改变
  const setSelectRows = (val: any) => {
    multipleSelection.value = val
  }
  const clearData = () => {
    curRow.value = {
      clientIp: '',
      serverIp: '',
      id: undefined,
      ruleId: '',
    }
  }
  const handleEdit = (row?: AlertWhiteList) => {
    if (row) {
      curRow.value = row
    }
    editVisible.value = true
  }
  const handleDelete = async (rows: AlertWhiteList[], deleteAll = false) => {
    $baseConfirm(!deleteAll ? '你确定要删除当前项吗' : '你确定要删除所有数据吗', null, async () => {
      const { msg } = await getDeleteWarnWhiteApi({ ids: rows.map((i) => (i.id ? i.id : 0)), deleteAll })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      queryPage.pageNum = 1
      getWhiteListHandle()
    })
  }
  const save = async () => {
    const { msg } = await getUpdateWarnWhiteApi(curRow.value)
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    await getWhiteListHandle()
    editVisible.value = false
  }

  const handleDownloadExcel = async () => {
    try {
      listLoading.value = true
      const res = await exportWarnWhiteApi({
        ids: multipleSelection.value.map((i) => (i.id ? i.id : 0)),
        searchStr: queryPage.searchStr ? queryPage.searchStr : undefined,
      })
      downloadFile(res, `告警白名单-${formatTime(new Date().getTime())}`)
      ElMessage({ message: '导出成功', type: 'success' })
      //
    } finally {
      listLoading.value = false
    }
  }

  const templateExport = async () => {
    try {
      uploadLoading.value = true
      const res = await exportWarnWhiteTemplateApi()
      downloadFile(res, '白名单模版')
      $baseMessage('下载成功', 'success', 'vab-hey-message-success')
    } finally {
      uploadLoading.value = false
    }
  }
  const templateUpload = async () => {
    uplodFormRef.value!.validate(async () => {
      const { msg } = await importWarnWhiteTemplateApi({ file: uploadFile.file[0].raw! })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      uploadVisible.value = false
      getWhiteListHandle()
    })
  }

  watch(
    () => uploadVisible.value,
    () => {
      if (!uploadVisible.value) {
        uploadFile.file = []
      }
    }
  )
  provide(tableSearch, () => {
    queryPage.pageNum = 1
    getWhiteListHandle()
  })
  const handleExceed: UploadProps['onExceed'] = (files) => {
    uplodRef.value?.clearFiles()
    const file = files[0] as UploadRawFile
    file.uid = genFileId()
    uplodRef.value!.handleStart(file)
  }

  const changeType = (str: any) => {
    const obj = {
      0: '源IP -> 目的IP',
      1: '源IP -> 目的IP:目的端口',
      2: '自定义',
      3: 'Any -> Any',
    }
    // @ts-ignore
    return obj[str]
  }

  onMounted(() => {
    getWhiteListHandle()
  })
</script>

<template>
  <div class="alert-whitelist-container">
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>白名单管理</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12">
        <el-dropdown style="margin-right: 10px">
          <span class="el-dropdown-link">
            <el-button :icon="Delete" type="danger">批量删除</el-button>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleDelete(multipleSelection)">删除选中</el-dropdown-item>
              <el-dropdown-item @click="handleDelete([], true)">删除所有</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button type="primary" @click="showDialog = true">白名单时间线</el-button>
        <el-button plain type="primary" @click="uploadVisible = true">导入</el-button>
        <el-button plain style="margin-right: 0 !important" type="primary" @click="handleDownloadExcel">导出</el-button>
      </vab-query-form-right-panel>
    </vab-query-form>
    <el-table v-loading="listLoading" :data="listDate" row-key="id" @selection-change="setSelectRows">
      <el-table-column show-overflow-tooltip type="selection" />
      <el-table-column label="序号" width="55">
        <template #default="{ $index }">
          {{ curIndex + $index }}
        </template>
      </el-table-column>
      <el-table-column label="名称" prop="ruleName" show-overflow-tooltip>
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="queryPage.searchStr"
            :column="column"
            filter-type="text"
            :tools="['filter']"
          />
        </template>
      </el-table-column>
      <el-table-column label="加白方式" prop="type" show-overflow-tooltip>
        <template #default="{ row }">
          {{ changeType(row.type) }}
        </template>
      </el-table-column>
      <el-table-column label="威胁名称" prop="threatName" show-overflow-tooltip />
      <el-table-column label="添加时间" prop="createTime" show-overflow-tooltip>
        <template #default="{ row }">
          {{ formatTime(row.createTime) }}
        </template>
      </el-table-column>
      <el-table-column label="规则ID" prop="ruleId" show-overflow-tooltip />
      <el-table-column label="操作" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" @click="handleDelete([row])">删除</el-button>
        </template>
      </el-table-column>
      <template #empty>
        <el-empty class="vab-data-empty" description="暂无数据" />
      </template>
    </el-table>
    <el-pagination
      v-model:current-page="queryPage.pageNum"
      v-model:page-size="queryPage.pageSize"
      background
      :layout="'total, sizes, prev, pager, next, jumper'"
      :page-sizes="[10, 20, 50]"
      :total="total"
      @current-change="getWhiteListHandle"
      @size-change="getWhiteListHandle"
    />
    <el-dialog v-if="uploadVisible" v-model="uploadVisible" destroy-on-close title="导入" width="885px">
      <el-form
        ref="uplodFormRef"
        class="examplesDetailUpload-form"
        label-width="100px"
        :model="uploadFile"
        :rules="rules"
      >
        <el-form-item label="上传文件" prop="file">
          <el-upload
            ref="uplodRef"
            v-model:file-list="uploadFile.file"
            accept=".xlsx"
            :auto-upload="false"
            class="upload-demo"
            drag
            :limit="1"
            :on-exceed="handleExceed"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或
              <em>选择文件</em>
              <br />
              只能上传xlsx格式文件
            </div>
          </el-upload>
        </el-form-item>
        <el-form-item label="&nbsp">
          <div class="el-upload__tip">
            <el-button @click="templateExport">Excel导入模版下载</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button :loading="uploadLoading" type="primary" @click="templateUpload">开始上传</el-button>
        <el-button @click="uploadVisible = false">取消</el-button>
      </template>
    </el-dialog>
    <Whitelist
      v-if="editVisible"
      :cur-data="curRow"
      :is-show="editVisible"
      :mode="'edit'"
      @on-closeEvent="editVisible = false"
      @on-reflash="getWhiteListHandle()"
    />
    <WhiteTime v-if="showDialog" :cur-data="curItem" :is-show="showDialog" @on-closeEvent="showDialog = false" />
  </div>
</template>

<style scoped lang="scss">
  .alert-whitelist-action {
    margin-bottom: 12px;
    .right,
    .left {
      display: inline-block;
    }
    .right {
      float: right;
      .el-button {
        margin-left: 12px;
      }
    }
    .el-dropdown {
      vertical-align: initial;
      margin-right: 12px;
    }
  }
  .upload-demo {
    width: 745px;
  }
  .el-upload__tip {
    font-size: 20px;
    margin: -10px auto 0;
  }
  :deep() {
    .el-scrollbar {
      height: calc(100vh - 170px) !important;
    }
  }
</style>
