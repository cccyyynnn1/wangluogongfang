<script lang="ts">
  export default {
    name: 'CustomRules', // 自定义告警
  }
</script>

<script setup lang="ts">
  import {
    getMaintenanceOptionsApi,
    getWarnCustomerRulePageApi,
    updateWarnCustomerRuleStatusApi,
    deleteWarnCustomerRuleApi,
    applyRulesApi,
  } from '@/api-ecs/forword-rule'
  import UpdateRules from './update-rules.vue'
  import { formatTime } from '@/utils/time'
  import { CheckboxValueType, Action } from 'element-plus'
  import { MaintenanceOptions, WarnCustomerRuleQueryType } from '@/types'
  import { debounce } from 'lodash'
  const route = useRoute()
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const fold = ref(false) // 筛选折叠
  const showPage = ref(false) // 弹框
  const total = ref(0) // 总条数
  const mode = ref()
  const listLoading = ref(false)
  const showRuleStatus = ref(false)
  const applyRulesStatusList = ref<{ msg: string; success: boolean; failReason: string }[]>([])
  let listDate = ref()
  // 检索条件
  const queryForm = reactive<WarnCustomerRuleQueryType>({
    pageSize: 10,
    pageNum: 1,
    ruleName: '',
    ruleTypes: [],
    severitys: [],
    enables: [],
    attackResults: [],
    confidences: [],
  })

  const maintenanceOptions = reactive<MaintenanceOptions & { enable: string[] }>({
    actionType: [],
    attackResult: [],
    attackType: [],
    killchain: [],
    attackMethod: [],
    ruleType: [],
    confidence: [],
    codeLanguage: [],
    severity: [],
    direction: [],
    enable: ['启用', '停用'],
  })

  const curIndex = computed(() => (queryForm.pageNum - 1) * queryForm.pageSize + 1)
  const multipleSelection = ref([])
  // 多选项改变
  const setSelectRows = (val: any) => {
    multipleSelection.value = val.map((i: any) => i.id)
  }
  const curRow = ref()
  const handleEdit = (row: any) => {
    mode.value = 'edit'
    curRow.value = row
    showPage.value = true
  }
  // 删除
  let delList: string[] = [] // 选中项的数组
  const handleDelete = (ids: number[], deleteAll = false) => {
    if (!deleteAll && !ids.length) return
    $baseConfirm('您确定要删除所选项吗', null, async () => {
      try {
        listLoading.value = true
        const { msg } = await deleteWarnCustomerRuleApi({ ids, deleteAll })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        handleReload()
      } catch {
        listLoading.value = false
      }
    })
  }
  const handleFold = () => {
    fold.value = !fold.value
  }
  const getData = async () => {
    listLoading.value = true
    const { data } = await getWarnCustomerRulePageApi({ ...queryForm })
    total.value = data.total
    listDate.value = data.records
    listLoading.value = false
  }

  const handleAdd = () => {
    mode.value = 'add'
    showPage.value = true
  }

  const handleReload = () => {
    handleReset()
  }

  const handleCheckAll = (mark: keyof WarnCustomerRuleQueryType, val: CheckboxValueType) => {
    const isMultiple = Array.isArray(queryForm[mark])
    const key = (isMultiple ? mark.slice(0, -1) : mark) as keyof MaintenanceOptions
    if (val) {
      Object.assign(queryForm, { [mark]: isMultiple ? maintenanceOptions[key] : '' })
    } else {
      Object.assign(queryForm, { [mark]: isMultiple ? [] : '' })
    }
  }
  const getConfidence = (code: string) => {
    const index = ['低', '中', '高', '可信'].findIndex((i) => code === i) + 1
    return index > 3 ? 3 : index
  }

  // 修改状态
  const handleUpdateStatus = async (enable: '启用' | '停用') => {
    listLoading.value = true
    const ids = multipleSelection.value
    try {
      const { msg } = await updateWarnCustomerRuleStatusApi({ enable, ids })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      handleReload()
    } catch {
      listLoading.value = false
    }
  }

  // 应用规则
  const handleUseRules = () => {
    $baseConfirm(
      '规则变更后执行该操作以应用最新规则配置，大约需要30秒重启服务，重启期间部分服务无法应用',
      null,
      async () => {
        try {
          listLoading.value = true
          const { data } = await applyRulesApi()
          applyRulesStatusList.value = data || []
          showRuleStatus.value = true
        } finally {
          listLoading.value = false
        }
      }
    )
  }
  const handleGetOptions = async () => {
    const { data } = await getMaintenanceOptionsApi()
    Object.assign(maintenanceOptions, data)
  }
  watch(
    () => [queryForm.attackResults, queryForm.confidences, queryForm.ruleTypes, queryForm.severitys, queryForm.enables],
    debounce(() => {
      getData()
    }, 900)
  )

  watch(
    () => showPage.value,
    () => {
      if (!showPage.value) {
        curRow.value = null
      }
    }
  )

  onMounted(() => {
    if (route.query.addRule) {
      const addRuleData = decodeURIComponent(route.query.addRule as string)
      const ruleData = JSON.parse(addRuleData) as {
        ruleVal: string
        ruleDirection: string
      }
      curRow.value = {
        ruleContent: JSON.stringify({ '': ruleData.ruleVal }),
        direction: ruleData.ruleDirection,
      }
      showPage.value = true
    }

    getData()
    handleGetOptions()
  })

  const handleReset = () => {
    multipleSelection.value = []
    Object.assign(queryForm, {
      pageSize: 10,
      pageNum: 1,
      ruleName: '',
      ruleTypes: [],
      severitys: [],
      enables: [],
      attackResults: [],
      confidences: [],
    })
  }
</script>

<template>
  <div v-loading="listLoading" class="custom-rules-container">
    <div v-if="!showPage" class="warp">
      <h3>自定义告警</h3>
      <div class="header">
        <div class="seacher-warp">
          <el-form inline label-position="right" label-width="auto" style="width: 100%">
            <el-form-item class="my-form-item" label="规则名称">
              <el-input v-model="queryForm.ruleName" clearable @keyup.enter="getData" />
            </el-form-item>
            <el-form-item class="my-form-item" label="规则类型">
              <el-select
                v-model="queryForm.ruleTypes"
                clearable
                collapse-tags
                :max-collapse-tags="1"
                multiple
                placeholder="请选择"
                popper-class="custom-header"
                style="width: 100%"
              >
                <template #default>
                  <el-checkbox style="margin-left: 20px" @change="(e) => handleCheckAll('ruleTypes', e)">
                    全部
                  </el-checkbox>
                  <el-divider style="margin-bottom: 0" />
                  <el-option v-for="item in maintenanceOptions?.ruleType" :key="item" :label="item" :value="item" />
                </template>
              </el-select>
            </el-form-item>
            <el-form-item class="my-form-item" label="威胁等级">
              <el-select
                v-model="queryForm.severitys"
                clearable
                collapse-tags
                :max-collapse-tags="1"
                multiple
                placeholder="请选择"
                popper-class="custom-header"
                style="width: 100%"
              >
                <template #default>
                  <el-checkbox style="margin-left: 20px" @change="(e) => handleCheckAll('severitys', e)">
                    全部
                  </el-checkbox>
                  <el-divider style="margin-bottom: 0" />
                  <el-option v-for="item in maintenanceOptions?.severity" :key="item" :label="item" :value="item" />
                </template>
              </el-select>
            </el-form-item>
            <el-form-item class="my-form-item" label="启用状态">
              <el-select
                v-model="queryForm.enables"
                clearable
                collapse-tags
                :max-collapse-tags="1"
                multiple
                placeholder="请选择"
                popper-class="custom-header"
                style="width: 100%"
              >
                <template #default>
                  <el-checkbox style="margin-left: 20px" @change="(e) => handleCheckAll('enables', e)">
                    全部
                  </el-checkbox>
                  <el-divider style="margin-bottom: 0" />
                  <el-option v-for="item in maintenanceOptions.enable" :key="item" :label="item" :value="item" />
                </template>
              </el-select>
            </el-form-item>
            <el-form-item class="my-form-item" label="攻击结果">
              <el-select
                v-model="queryForm.attackResults"
                clearable
                collapse-tags
                :max-collapse-tags="1"
                multiple
                placeholder="请选择"
                popper-class="custom-header"
                style="width: 100%"
              >
                <template #default>
                  <el-checkbox style="margin-left: 20px" @change="(e) => handleCheckAll('attackResults', e)">
                    全部
                  </el-checkbox>
                  <el-divider style="margin-bottom: 0" />
                  <el-option v-for="item in maintenanceOptions?.attackResult" :key="item" :label="item" :value="item" />
                </template>
              </el-select>
            </el-form-item>
            <el-form-item class="my-form-item" label="置信度" prop="reliable">
              <el-select
                v-model="queryForm.confidences"
                clearable
                collapse-tags
                :max-collapse-tags="1"
                multiple
                placeholder="请选择"
                popper-class="custom-header"
                style="width: 100%"
              >
                <template #default>
                  <el-checkbox style="margin-left: 20px" @change="(e) => handleCheckAll('confidences', e)">
                    全部
                  </el-checkbox>
                  <el-divider style="margin-bottom: 0" />
                  <el-option v-for="item in maintenanceOptions?.confidence" :key="item" :label="item" :value="item" />
                </template>
              </el-select>
            </el-form-item>
          </el-form>
        </div>
        <div class="operation-warp">
          <div class="my-left">
            <el-button :auto-insert-space="false" type="primary" @click="handleUseRules">应用规则</el-button>
            <el-button :auto-insert-space="false" type="primary" @click="handleAdd">新增</el-button>
            <el-button
              :auto-insert-space="false"
              :disabled="!multipleSelection.length"
              type="primary"
              @click="handleUpdateStatus('启用')"
            >
              启用
            </el-button>
            <el-button
              :auto-insert-space="false"
              :disabled="!multipleSelection.length"
              type="primary"
              @click="handleUpdateStatus('停用')"
            >
              停用
            </el-button>
            <el-dropdown style="margin-left: 10px">
              <span class="el-dropdown-link">
                <el-button type="danger">
                  批量删除
                  <el-icon class="el-icon--right"><arrow-down /></el-icon>
                </el-button>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleDelete(multipleSelection)">删除选中</el-dropdown-item>
                  <el-dropdown-item @click="handleDelete([], true)">删除所有</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
          <div class="my-right">
            <el-button :auto-insert-space="false" type="primary" @click="handleReset">重置</el-button>
          </div>
        </div>
      </div>
      <el-table :data="listDate" row-key="id" @selection-change="setSelectRows">
        <el-table-column show-overflow-tooltip type="selection" />
        <el-table-column label="序号" width="55">
          <template #default="{ $index }">
            {{ curIndex + $index }}
          </template>
        </el-table-column>
        <el-table-column label="规则编号" prop="ruleId" show-overflow-tooltip width="110" />
        <el-table-column align="center" label="规则名称" prop="ruleName" show-overflow-tooltip />
        <el-table-column align="center" label="规则类型" prop="ruleType" show-overflow-tooltip />
        <el-table-column align="center" label="威胁等级" prop="severity" width="90">
          <template #default="{ row }">
            <span :class="['alert_tag', row.severityStr]">
              <el-icon><WarnTriangleFilled /></el-icon>
              {{ row.severity }}
            </span>
          </template>
        </el-table-column>

        <el-table-column align="center" label="攻击结果" prop="attackResult" show-overflow-tooltip width="90" />
        <el-table-column align="center" label="启用状态" prop="enable" show-overflow-tooltip width="90" />
        <el-table-column align="center" label="添加时间" prop="createTime" show-overflow-tooltip width="180">
          <template #default="{ row }">
            {{ formatTime(row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="置信度" width="90">
          <template #default="{ row }">
            <el-rate
              :colors="['#67C23A', '#67C23A', '#67C23A']"
              disabled
              disabled-void-color="#C7C6D4"
              :max="3"
              :model-value="getConfidence(row.confidence)"
            />
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" width="140">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleDelete([row.id])">删除</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty class="vab-data-empty" description="暂无数据" />
        </template>
      </el-table>
      <el-pagination
        v-model:current-page="queryForm.pageNum"
        v-model:page-size="queryForm.pageSize"
        background
        :layout="'total, sizes, prev, pager, next, jumper'"
        :page-sizes="[10, 20, 50]"
        :total="total"
        @current-change="getData"
        @size-change="getData"
      />
    </div>
    <update-rules
      v-if="showPage"
      :all-select-option="maintenanceOptions"
      :cur-data="curRow"
      :remark="mode"
      :show-pages="showPage"
      @on-close-event="showPage = false"
      @on-reflash="handleReload"
    />
    <el-dialog v-model="showRuleStatus" destroy-on-close title="应用规则提示" width="520px">
      <div class="status-list">
        <div v-for="(item, index) of applyRulesStatusList" :key="index" class="status-item">
          <div class="status-icon success">
            <el-icon v-if="item.success" color="#55BF3A" size="18"><SuccessFilled /></el-icon>
            <el-icon v-else color="#FA6A66" size="18"><CircleCloseFilled /></el-icon>
          </div>
          <div class="status-msg">
            <p>{{ item.msg }}</p>
            <span v-if="item.failReason">{{ item.failReason }}</span>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  $criticalColor: #ff0202;
  $lowColor: #ffbe36;
  $midColor: #f1b04d;
  $highColor: #fa6d15;
  $defaultColor: #909399;
  .custom-rules-container {
    h3 {
      margin-block: 0 1em;
    }
    .warp {
      height: 100%;
      display: flex;
      flex-direction: column;
      .header {
        .seacher-warp {
          .my-form-item {
            width: 33.3333%;
            margin-right: 0;
            padding-right: 20px;
            &:nth-child(3n) {
              padding-right: 0;
            }
          }
        }
      }
      .el-table {
        .alert_tag {
          display: inline-block;
          width: 55px;
          padding: 2px 3px;
          border-radius: 5px;
          color: #fff;
          font-size: 14px;
          .el-icon {
            margin-right: -4px;
            vertical-align: -3px;
            font-size: 17px;
          }
          &.critical {
            background: #ca0a08;
          }
          &.high {
            background: #ff2927;
          }
          &.medium {
            background: #ff7212;
          }
          &.low {
            background: #ffbe36;
          }
        }
        flex: auto;
      }
    }
    .status-list {
      .status-item {
        display: flex;
        margin-bottom: 12px;
        .status-icon {
          width: 28px;
          display: flex;
          align-items: center;
          position: relative;
          &::after,
          &::before {
            content: ' ';
            display: block;
            position: absolute;
            width: 1px;
            background-color: #e8e6f6;
            left: 9px;
            z-index: 1;
          }
          &::after {
            top: 24px;
            bottom: 0;
          }
          &::before {
            top: -12px;
            bottom: calc(100% - 12px);
          }
          .el-icon {
            position: absolute;
            top: 11px;
            z-index: 2;
          }
        }
        .status-msg {
          flex: 1;
          background: #f6f4fd;
          border-radius: 4px;
          line-height: 26px;
          padding-block: 6px;
          padding-left: 14px;
          p {
            margin-block: 0;
            font-weight: 500;
            font-size: 14px;
            color: #2b2742;
          }
          span {
            font-weight: 400;
            font-size: 14px;
            color: #4a4759;
          }
        }
        &:last-child {
          .status-icon {
            &::after {
              width: 0;
            }
          }
        }
        &:first-child {
          .status-icon {
            &::before {
              width: 0;
            }
          }
        }
      }
    }
  }

  .operation-warp {
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  :deep() {
    .el-form-item__label-wrap {
      margin-right: 0 !important;
    }
    .custom-header {
      .el-checkbox {
        display: flex;
        height: unset;
      }
    }
  }
  .custom-header {
    .el-checkbox {
      display: flex;
      height: unset;
    }
  }
</style>
