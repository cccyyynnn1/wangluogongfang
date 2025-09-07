<script lang="ts">
  export default {
    name: 'Interface',
  }
</script>
<script setup lang="ts">
  import { getApiInterfaceApi, deleteApiInterfaceApi } from '~/src/api-ecs/site'
  import { useTableCopy } from '@/utils'
  import ApiEdit from './api-edit.vue'
  import { siteApiInterfaceItemType } from '@/types'
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const router = useRouter()
  const props = defineProps<{
    id: number
  }>()
  const showApiEdit = ref(false) // 是否显示编辑弹窗
  const curApiData = ref<undefined | siteApiInterfaceItemType>(undefined)
  // 表格数据
  const layout = ref('total, sizes, prev, pager, next, jumper')
  const queryForm = reactive({
    pageNum: 1,
    pageSize: 10,
    apiName: '',
    apiUrl: '',
    searchSql: '',
  })
  const queryPage = reactive({
    total: 0,
    apiList: [] as siteApiInterfaceItemType[],
    loading: true,
    setSelectRows: [] as siteApiInterfaceItemType[],
  })

  const dataLabel = ref() // 数据标签
  const warnLabel = ref() // 告警标签

  onMounted(() => {
    queryData()
  })

  // 获取告警标签
  const getApiInterface = async () => {
    queryPage.loading = true
    const {
      data: { total, records },
    } = await getApiInterfaceApi({ ...queryForm, sessionId: props.id })
    queryPage.total = total
    queryPage.apiList = records
    queryPage.loading = false
  }
  // 获取表格序号
  const curIndex = computed(() => (queryForm.pageNum - 1) * queryForm.pageSize + 1)
  // 查询
  const queryData = () => {
    queryForm.pageNum = 1
    getApiInterface()
  }

  // 详情
  const handleInfo = (val: siteApiInterfaceItemType) => {
    router.push({ name: 'SiteIndex', query: { apiId: val.id, sessionId: props.id } })
  }

  // 编辑
  const handleEdit = (val: siteApiInterfaceItemType) => {
    showApiEdit.value = true
    curApiData.value = val
  }

  // 添加
  const handleAddApi = () => {
    curApiData.value = undefined
    showApiEdit.value = true
  }
  // 删除
  const handleDeleteApi = (selects: siteApiInterfaceItemType[]) => {
    $baseConfirm('你确定要删除当前项吗', null, async () => {
      const ids = selects.map((i) => i.id)
      const { msg } = await deleteApiInterfaceApi(ids.toString())
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      queryData()
    })
  }
  // 多选项改变
  const setSelectRows = (rows: siteApiInterfaceItemType[]) => {
    queryPage.setSelectRows = rows
  }
</script>

<template>
  <div class="interface">
    <el-form label-position="top" style="width: 100%; display: flex">
      <el-row :gutter="10" style="width: 100%; display: flex">
        <el-col :span="4">
          <el-form-item label="接口名称">
            <el-input v-model="queryForm.apiName" clearable />
          </el-form-item>
        </el-col>
        <el-col :span="18">
          <el-form-item label="接口地址">
            <el-input v-model="queryForm.apiUrl" clearable />
          </el-form-item>
        </el-col>
        <el-col :span="2">
          <el-form-item label="&nbsp">
            <el-button style="margin-left: 36px" type="primary" @click="queryData">检索</el-button>
          </el-form-item>
        </el-col>
        <el-col v-if="false" :span="8">
          <el-form-item class="form-tags" label="数据标签" prop="roleName">
            <div class="tags">
              <el-tag v-for="(item, index) in dataLabel" :key="index">
                {{ item[1] }}
              </el-tag>
            </div>
          </el-form-item>
        </el-col>
        <el-col v-if="false" :span="8">
          <el-form-item class="form-tags" label="告警标签" prop="roleName">
            <div class="tags">
              <el-tag v-for="(item, index) in warnLabel" :key="index">
                {{ item[1] }}
              </el-tag>
            </div>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
    <el-space style="margin-bottom: 20px">
      <el-button type="primary" @click="handleAddApi">添加</el-button>
      <!-- <el-button type="primary" @click="queryData">接口发现</el-button> -->
      <el-button
        :disabled="queryPage.setSelectRows.length === 0"
        type="danger"
        @click="handleDeleteApi(queryPage.setSelectRows)"
      >
        批量删除
      </el-button>
    </el-space>
    <!-- 表格 -->
    <el-table
      v-loading="queryPage.loading"
      class="my-table"
      :data="queryPage.apiList"
      row-key="id"
      @cell-contextmenu="useTableCopy"
      @selection-change="setSelectRows"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column :index="curIndex" label="序号" type="index" width="80" />
      <el-table-column label="接口名称" prop="apiName" show-overflow-tooltip width="300" />
      <el-table-column label="接口地址" prop="apiUrl" show-overflow-tooltip />
      <el-table-column label="接口是否展示" width="120">
        <template #default="{ row }">
          <span class="apiDisplayTag" :class="{ isDisplay: row.isDisplay === 1 }">
            {{ row.isDisplay === 1 ? '是' : '否' }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="240">
        <template #default="{ row }">
          <el-button size="small" @click="handleInfo(row)">会话</el-button>
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" @click="handleDeleteApi([row])">删除</el-button>
        </template>
      </el-table-column>
      <template #empty>
        <el-empty description="暂无数据" />
      </template>
    </el-table>

    <!-- pagination -->
    <el-pagination
      v-model:current-page="queryForm.pageNum"
      v-model:page-size="queryForm.pageSize"
      background
      :layout="layout"
      :page-sizes="[10, 20, 50]"
      :total="queryPage.total"
      @current-change="getApiInterface"
      @size-change="queryData"
    />
  </div>
  <!-- 编辑api站点  -->
  <api-edit
    v-if="showApiEdit"
    :api-data="curApiData"
    :session-id="props.id"
    :show-api-edit="showApiEdit"
    @on-closeEvent="showApiEdit = false"
    @on-reflash="queryData"
  />
</template>

<style scoped lang="scss">
  .interface {
    height: 100%;
    width: 100%;
    :deep() {
      .el-table__body-wrapper,
      .el-scrollbar {
        height: 420px !important;
      }
    }
  }
  .apiDisplayTag {
    display: inline-block;
    padding: 2px 12px;
    border-radius: 4px;
    background-color: #f9f8ff;
    border: 1px solid #d6d3e7;
    color: #4a4759;
    &.isDisplay {
      background-color: #f4f3ff;
      border: 1px solid #bdb2ff;
      color: #6954f0;
    }
  }

  .form-tags {
    width: 100%;
  }

  .tags {
    height: 112px;
    overflow-y: auto;
    padding: 10px;

    &::-webkit-scrollbar {
      width: 0;
      height: 0;
    }

    border: 1px solid var(--el-border-color);
  }
</style>
