<script lang="ts">
  export default {
    name: 'FieldInitialisation', // 初始化配置
  }
</script>

<script setup lang="ts">
  import { getFieldsTemplateApi, usingFieldsTemplateApi, deleteFieldsTemplateApi } from '~/src/api-ecs/custom-field'
  import TableFilterTool from '@/components/table-filter-tool.vue'
  import { Plus } from '@element-plus/icons-vue'
  import { formatNstime } from '@/utils/time'
  import { InitialisationItem, tableSearch } from '@/types/index'
  import FieldDeleteHistory from './delete-history.vue'
  import AddFieldModule from './add-field-module.vue'
  import FieldModules from './field-modules.vue'
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const showAddFieldDialog = ref(false) // 打开编辑页
  const fieldInfovVisible = ref(false)
  const listDate = ref<InitialisationItem[]>([])
  const listLoading = ref(false) // 是否加载
  const total = ref(0)

  // 当前模版下的字段详情
  const templateData = ref()
  const showHistory = ref(false)
  const displayStatus = ref('default')
  // 获取表格序号
  const curIndex = computed(() => (queryPage.pageNum - 1) * queryPage.pageSize + 1)
  const queryPage = reactive({
    templateName: '',
    pageNum: 1,
    pageSize: 10,
  })

  provide(tableSearch, () => {
    queryPage.pageNum = 1
    getAllTemplate()
  })
  const getAllTemplate = async () => {
    listLoading.value = true
    try {
      const { data } = await getFieldsTemplateApi(queryPage)
      listDate.value = data.records || []
      total.value = data.total
    } finally {
      listLoading.value = false
    }
  }
  const handleFieldStatusChange = async (status: any, row: InitialisationItem) => {
    try {
      const { msg } = await usingFieldsTemplateApi({ status, id: row.id })
      getAllTemplate()
      $baseMessage(msg, 'success', 'vab-hey-message-success')
    } catch (error) {
      row.status = status === 0 ? 1 : 0
    }
  }

  const handleDeleteTemplate = async (id: number) => {
    $baseConfirm('你确定要删除当前项吗？', null, async () => {
      const { msg } = await deleteFieldsTemplateApi([id])
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      getAllTemplate()
    })
  }

  const hanldleShowHistory = (show: boolean) => {
    showHistory.value = show
  }
  const handleShowFielInfo = (template_data: any, display = 'default') => {
    displayStatus.value = template_data.templateName === '系统默认模版' ? 'preview' : display
    templateData.value = template_data
    fieldInfovVisible.value = true
  }
  onMounted(() => {
    getAllTemplate()
  })
</script>

<template>
  <div class="field-initialisation-container">
    <template v-if="!showHistory">
      <vab-query-form>
        <vab-query-form-left-panel :span="12">
          <h3>字段初始化</h3>
        </vab-query-form-left-panel>
        <vab-query-form-right-panel :span="12">
          <el-button plain @click="() => hanldleShowHistory(true)">历史删除</el-button>
          <el-button
            :icon="Plus"
            style="margin-right: 0 !important"
            type="primary"
            @click="() => (showAddFieldDialog = true)"
          >
            新建自定义
          </el-button>
        </vab-query-form-right-panel>
      </vab-query-form>
      <el-table v-loading="listLoading" class="my-table" :data="listDate" row-key="id">
        <el-table-column align="left" :index="curIndex" label="序号" type="index" width="75" />
        <el-table-column label="名称" prop="templateName" show-overflow-tooltip>
          <template #header="{ column }">
            <table-filter-tool
              v-model:filter-value="queryPage.templateName"
              :column="column"
              filter-type="text"
              :tools="['filter']"
            />
          </template>
        </el-table-column>
        <el-table-column label="创建时间" prop="createTime" width="180">
          <template #default="{ row }">
            {{ formatNstime(row.createTime, false) }}
          </template>
        </el-table-column>
        <el-table-column label="修改时间" prop="updateTime" width="180">
          <template #default="{ row }">
            {{ formatNstime(row.updateTime, false) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" prop="status" show-overflow-tooltip>
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              active-text="开启"
              :active-value="1"
              inactive-text="停用"
              :inactive-value="0"
              inline-prompt
              @change="handleFieldStatusChange($event, row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="备注" prop="remark" show-overflow-tooltip />
        <el-table-column label="操作" width="140">
          <template #default="{ row }">
            <el-button size="small" @click="handleShowFielInfo(row)">详情</el-button>
            <el-button
              :disabled="row.templateName === '系统默认模版'"
              size="small"
              @click="handleDeleteTemplate(row.id)"
            >
              删除
            </el-button>
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
        layout="total, sizes, prev, pager, next, jumper"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        @current-change="getAllTemplate"
        @size-change="getAllTemplate"
      />
    </template>
    <template v-else>
      <field-delete-history @back="hanldleShowHistory" @show-info="handleShowFielInfo" />
    </template>
    <add-field-module v-model="showAddFieldDialog" />
    <field-modules v-model:visible="fieldInfovVisible" :display="displayStatus" :template-data="templateData" />
  </div>
</template>

<style scoped lang="scss">
  .field-initialisation-container {
    padding: 14px !important;
    height: 100%;
    box-sizing: border-box;
    overflow-y: auto;
    h3 {
      margin-block: 0 0.5em;
    }
    &::-webkit-scrollbar {
      width: 0;
      height: 0;
    }
    h3 {
      margin-top: 0;
    }
  }
  .my-table {
    height: calc(100vh - 138px);
  }
</style>
