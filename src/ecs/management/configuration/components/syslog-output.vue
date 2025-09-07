<script lang="ts">
  export default {
    name: 'SyslogOutput',
  }
</script>

<script setup lang="ts">
  import { Plus, Search } from '@element-plus/icons-vue'
  import JsonPreview from '@/components/json-preview.vue'
  import SyslogOutputInfo from './syslog-putput-info.vue'
  import {
    getTemplateRuleListApi,
    deleteTemplateRulestApi,
    updateTemplateRulestStatusApi,
  } from '@/api-ecs/log-data-output'
  import { formatNstime } from '@/utils/time'
  const logFieldRef = ref()
  const logPopoverRef = ref()
  const logPopoverVisible = ref(false)
  const logPopoverFields = ref('{}')

  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const outputInfoRef = ref<InstanceType<typeof SyslogOutputInfo> | null>(null)
  const queryData = reactive({
    pageNum: 1,
    pageSize: 10,
    searchStr: '',
  })
  const list_total = ref(0)
  const listDate = ref<any[]>()
  const listLoading = ref(false)
  const multipleSelection = ref<any[]>([])
  // 多选项改变
  const setSelectRows = (val: any) => {
    multipleSelection.value = val
  }
  const getTemplateRuleListHandler = async () => {
    listLoading.value = true
    const {
      data: { records, total },
    } = await getTemplateRuleListApi(queryData)
    listDate.value = records
    list_total.value = total
    listLoading.value = false
  }

  const getFields = (configStr: string): string[] => {
    try {
      const config = JSON.parse(configStr) as {
        [key: string]: string[]
      }
      const keys = Object.keys(config) || []
      const valKyes = Object.values(config) || []
      return [...keys, ...valKyes].flat(2)
    } catch (error) {
      return []
    }
  }
  const handleDelete = async (rows: any[], deleteAll = false) => {
    $baseConfirm(!deleteAll ? '你确定要删除当前项吗' : '你确定要删除所有数据吗', null, async () => {
      const { msg } = await deleteTemplateRulestApi({ ids: rows.map((i) => i.id), deleteAll })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      getTemplateRuleListHandler()
    })
  }
  const handlerInfo = (row: any) => {
    const {
      id,
      ip,
      logType,
      attackStageIds,
      ouputSource,
      outputNull,
      port,
      threatTypeList,
      threatLevelIds,
      probeIps,
      status,
      protocol,
      templateId,
      configStr,
      fieldsEchoInfos,
    } = row
    outputInfoRef.value?.showDialog({
      title: '详情',
      infoValue: {
        id,
        ip,
        logType,
        attackStageIds,
        ouputSource,
        outputNull,
        port,
        threatTypeList,
        threatLevelIds,
        probeIps,
        status,
        protocol,
        templateId,
        configStr,
        fieldsEchoInfos,
      },
    })
  }
  function configStrMouseover(e: MouseEvent, val: any) {
    logFieldRef.value = e.currentTarget
    logPopoverFields.value = val.configStr
  }
  function configStrClick() {
    logPopoverVisible.value = !logPopoverVisible.value
  }
  function searchStrClick() {
    queryData.pageNum = 1
    getTemplateRuleListHandler()
  }

  const handleStatusChange = async (val: any, row: any) => {
    row.disabled = true
    const { msg } = await updateTemplateRulestStatusApi({ status: val, ids: [row.id] })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    row.disabled = false
  }
  onMounted(async () => {
    getTemplateRuleListHandler()
  })
</script>

<template>
  <vab-query-form>
    <vab-query-form-left-panel :span="12">
      <h3>日志数据输出</h3>
    </vab-query-form-left-panel>
    <vab-query-form-right-panel :span="12">
      <el-space>
        <el-input
          v-model="queryData.searchStr"
          class="search"
          placeholder="请输入检索条件"
          style="width: 220px; margin-right: -11px"
        />
        <el-button :icon="Search" type="primary" @click="searchStrClick">检索</el-button>
        <el-button :icon="Plus" type="primary" @click="() => outputInfoRef?.showDialog({})">新增</el-button>
        <el-dropdown style="margin-right: -8px !important">
          <span class="el-dropdown-link">
            <el-button type="danger">
              批量删除
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleDelete(multipleSelection)">删除选中</el-dropdown-item>
              <el-dropdown-item @click="() => handleDelete([], true)">删除所有</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-space>
    </vab-query-form-right-panel>
  </vab-query-form>
  <el-table
    v-loading="listLoading"
    :data="listDate"
    style="margin-top: 8px; height: calc(100vh - 120px)"
    @selection-change="setSelectRows"
  >
    <el-table-column type="selection" />
    <el-table-column :formatter="({ logType }) => `${logType === 0 ? '告警日志' : '审计日志'}`" label="输出日志类型" />
    <el-table-column label="日志模版组" prop="templateName" />
    <el-table-column label="输出日志字段" prop="outputlog" width="200">
      <template #default="{ row }">
        <el-popover
          placement="right-start"
          popper-class="no-arrow"
          popper-style="box-shadow: rgb(14 18 22 / 35%) 0px 10px 38px -10px, rgb(14 18 22 / 20%) 0px 10px 20px -15px; padding:0px;border-radius: 10px;"
          trigger="hover"
          :width="200"
        >
          <template #default>
            <ul class="allField">
              <li v-for="(field, index) in getFields(row.fieldsEchoInfos)" :key="index">{{ field }}</li>
            </ul>
          </template>
          <template #reference>
            <span>
              <span>全部共计</span>
              <el-link type="primary">{{ getFields(row.fieldsEchoInfos)?.length }}个字段</el-link>
            </span>
          </template>
        </el-popover>
      </template>
    </el-table-column>
    <el-table-column label="日志字段" width="120">
      <template #default="{ row }">
        <span class="configPreviewBtn" @click="configStrClick" @mouseover="(event) => configStrMouseover(event, row)">
          日志样例展示
        </span>
      </template>
    </el-table-column>
    <el-table-column :formatter="({ ip, port }) => `${ip}:${port}`" label="对端IP/PORT" prop="ip_prot" width="160" />
    <!-- <el-table-column  label="业务组别" prop="target" /> -->
    <el-table-column
      :formatter="({ createTime }) => formatNstime(createTime, false)"
      label="创建时间"
      prop="createTime"
      width="180"
    />
    <el-table-column
      :formatter="({ updateTime }) => formatNstime(updateTime, false)"
      label="更新时间"
      prop="updateTime"
      width="180"
    />
    <el-table-column label="状态" prop="status">
      <template #default="{ row }">
        <el-switch
          v-model="row.status"
          active-text="开启"
          :active-value="1"
          :disabled="row.disabled"
          inactive-text="关闭"
          :inactive-value="0"
          inline-prompt
          @change="(val) => handleStatusChange(val, row)"
        />
      </template>
    </el-table-column>
    <el-table-column fixed="right" label="操作" width="140">
      <template #default="{ row }">
        <el-button class="row_action" size="small" @click="handlerInfo(row)">详情</el-button>
        <el-button class="row_action" size="small" @click="handleDelete([row])">删除</el-button>
      </template>
    </el-table-column>
    <template #empty><el-empty /></template>
  </el-table>
  <el-pagination
    v-model:current-page="queryData.pageNum"
    v-model:page-size="queryData.pageSize"
    background
    layout="total, sizes, prev, pager, next, jumper"
    :page-sizes="[10, 20, 50, 100]"
    :total="list_total"
  />
  <syslog-output-info ref="outputInfoRef" :callback="getTemplateRuleListHandler" />
  <el-popover
    ref="logPopoverRef"
    placement="right-start"
    popper-style="box-shadow: rgb(14 18 22 / 35%) 0px 10px 38px -10px, rgb(14 18 22 / 20%) 0px
    10px 20px -15px;padding:24px 30px 30px"
    trigger="click"
    :virtual-ref="logFieldRef"
    virtual-triggering
    :width="640"
  >
    <template #default>
      <div style="display: flex; justify-content: space-between">
        <h3 style="margin: 0; font-size: 14px">解析日志示例</h3>
        <el-space>
          <el-icon style="cursor: pointer" @click="logFieldRef?.click?.()">
            <Close />
          </el-icon>
        </el-space>
      </div>
      <div class="mask">
        <json-preview :json-value="JSON.stringify(JSON.parse(logPopoverFields), null, 4)" />
      </div>
    </template>
  </el-popover>
</template>

<style scoped lang="scss">
  .mask {
    height: 560px;
    margin-top: 15px;
    overflow-y: auto;
    background-color: #f8f7ff;
    border: 1px solid rgb(220, 223, 230);
  }
  .configPreviewBtn {
    font-weight: 500;
    color: var(--el-color-primary);
    cursor: pointer;
    &:hover {
      text-decoration: underline;
    }
  }
  .allField {
    max-height: 360px;
    padding: 0;
    margin: 10px 0;
    overflow-y: auto;
    li {
      height: 36px;
      font-size: 14px;
      font-weight: 400;
      line-height: 36px;
      color: #303133;
      text-indent: 20px;
      &:nth-child(odd) {
        // background-color: #fff;
        background: #f8f7ff;
      }
    }
  }
</style>
