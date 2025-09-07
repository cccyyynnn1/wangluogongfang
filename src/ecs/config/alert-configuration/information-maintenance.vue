<script lang="ts">
  export default {
    name: 'InformationMaintenance',
  }
</script>

<script setup lang="ts">
  import {
    getMaintenanceOptionsApi,
    getMaintenancelistApi,
    updateMaintenanceInfoApi,
    updateMaintenanceStatusApi,
  } from '@/api-ecs/forword-rule'
  import cloneDeep from 'lodash/cloneDeep'
  import { tableSearch } from '@/types'
  import { Back } from '@element-plus/icons-vue'
  import TableFilterTool from '@/components/table-filter-tool.vue'
  import { formatNstime } from '@/utils/time'
  import type { FormInstance } from 'element-plus'
  import { MaintenanceOptions } from '@/types'
  const $baseMessage: any = inject('$baseMessage')
  const showInfo = ref(false)
  const severityType = ['危急', '高危', '中危', '低危']
  const queryPage = reactive({
    pageSize: 20,
    pageNum: 1,
    searchStr: '',
  })
  const total = ref(0)
  const list = ref()
  // 获取表格序号
  const curIndex = computed(() => (queryPage.pageNum - 1) * queryPage.pageSize + 1)
  const spacer = h(ElDivider, { direction: 'vertical' })
  const infoForm = ref()
  const multipleSelection = ref()
  const infoFormRef = ref<FormInstance>()
  const maintenanceOptions = reactive<MaintenanceOptions>({
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
  })
  const rules = {
    id: [
      {
        required: true,
        trigger: 'blur',
      },
    ],
    ruleName: [
      {
        required: true,
        message: '请输入规则名称',
        trigger: 'blur',
      },
    ],
    ruleType: [
      {
        required: true,
        message: '请选择规则类型',
        trigger: 'blur',
      },
    ],
    attackType: [
      {
        required: true,
        message: '请选择攻击类型',
        trigger: 'blur',
      },
    ],
    severity: [
      {
        required: true,
        message: '请选择威胁等级',
        trigger: 'blur',
      },
    ],
    direction: [
      {
        required: true,
        message: '请选择攻击方向',
        trigger: 'blur',
      },
    ],
    actionType: [
      {
        required: true,
        message: '请选择攻击方式',
        trigger: 'blur',
      },
    ],
    attackMethod: [
      {
        required: true,
        message: '请选择攻击手段',
        trigger: 'blur',
      },
    ],
    killchain: [
      {
        required: true,
        message: '请选择攻击阶段',
        trigger: 'blur',
      },
    ],
    attackResult: [
      {
        required: true,
        message: '请选择攻击结果',
        trigger: 'blur',
      },
    ],
    confidence: [
      {
        required: true,
        message: '请选择置信度',
        trigger: 'blur',
      },
    ],
    vulnHarm: [
      {
        required: true,
        message: '请输入风险危害',
        trigger: 'blur',
      },
    ],
    detailInfo: [
      {
        required: true,
        message: '请输入威胁详情',
        trigger: 'blur',
      },
    ],
    vulnDesc: [
      {
        required: true,
        message: '请输入威胁描述',
        trigger: 'blur',
      },
    ],
  }

  provide(tableSearch, () => {
    queryPage.pageNum = 1
    handleGetMaintenanceList()
  })
  const getConfidence = (code: string) => {
    return ['低', '中', '高'].findIndex((i) => code === i) + 1
  }
  // 告警阈值
  const alarmThresholdS = reactive({
    hour: 0,
    minutes: 0,
    second: 0,
  })
  // 告警间隔
  const alarmIntervalS = reactive({
    hour: 0,
    minutes: 0,
    second: 0,
  })

  const getTimes = (count: number) => {
    const hour = Math.floor(count / 3600)
    const minutes = Math.floor((count % 3600) / 60)
    const second = count % 60
    return {
      hour,
      minutes,
      second,
    }
  }
  const formatTime2second = (data: { hour: number; minutes: number; second: number }) => {
    return 3600 * data.hour + 60 * data.minutes + data.second
  }
  const hanldeUpdateStatus = async (status: '启用' | '停用') => {
    if (!multipleSelection.value?.length) return
    const arr = multipleSelection.value.length ? multipleSelection.value : list.value
    const { msg } = await updateMaintenanceStatusApi({
      ids: arr.map((i: any) => i.id),
      enable: status,
      selectAll: false,
    })
    queryPage.pageNum = 1
    handleGetMaintenanceList()
    $baseMessage(msg, 'success', 'vab-hey-message-success')
  }
  const handleSubmit = async () => {
    await infoFormRef.value?.validate(async (valid, fields) => {
      if (valid) {
        infoForm.value.alarmThresholdS = formatTime2second(alarmThresholdS)
        infoForm.value.alarmIntervalS = formatTime2second(alarmIntervalS)
        const { msg } = await updateMaintenanceInfoApi(infoForm.value)
        showInfo.value = false
        queryPage.pageNum = 1
        handleGetMaintenanceList()
        $baseMessage(msg, 'success', 'vab-hey-message-success')
      }
    })
  }
  const setSelectRows = (val: any) => {
    multipleSelection.value = val
  }
  const handleGetOptions = async () => {
    const { data } = await getMaintenanceOptionsApi()
    Object.assign(maintenanceOptions, data)
  }
  const handleGetMaintenanceList = async () => {
    const { data } = await getMaintenancelistApi(queryPage)
    list.value = data.records
    total.value = data.total
  }
  const handleEdit = (row: any) => {
    Object.assign(alarmThresholdS, getTimes(row.alarmThresholdS))
    Object.assign(alarmIntervalS, getTimes(row.alarmIntervalS))
    infoForm.value = cloneDeep(row)
    showInfo.value = true
  }
  onMounted(() => {
    handleGetOptions()
    handleGetMaintenanceList()
  })
</script>

<template>
  <div class="information-maintenance-container">
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>信息维护</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12">
        <el-input v-model="queryPage.searchStr" placeholder="模糊检索" style="width: 200px" />
        <el-button style="margin-left: -2px" type="primary" @click="handleGetMaintenanceList">
          <vab-icon icon="search-line" />
          检索
        </el-button>
        <el-button type="primary" @click="hanldeUpdateStatus('启用')">启用</el-button>
        <el-button style="margin-right: 0 !important" type="primary" @click="hanldeUpdateStatus('停用')">
          停用
        </el-button>
      </vab-query-form-right-panel>
    </vab-query-form>
    <el-table class="maintenance" :data="list" row-key="id" @selection-change="setSelectRows">
      <el-table-column fixed="left" type="selection" width="55" />
      <el-table-column fixed="left" label="序号" width="80">
        <template #default="{ $index }">
          {{ curIndex + $index }}
        </template>
      </el-table-column>
      <el-table-column label="规则ID" prop="id" show-overflow-tooltip width="120" />
      <el-table-column label="规则名称" prop="ruleName" show-overflow-tooltip width="480" />
      <el-table-column label="规则类型" prop="ruleType" show-overflow-tooltip width="180" />
      <el-table-column label="威胁等级" prop="severity" show-overflow-tooltip>
        <template #default="{ row }">
          <span
            v-if="row.severity"
            :class="['level', `level-${severityType.findIndex((type) => type === row?.severity)}`]"
          >
            <el-icon><WarnTriangleFilled /></el-icon>
            {{ row.severity }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="置信度" prop="confidence">
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
      <el-table-column label="攻击阶段" prop="killchain" show-overflow-tooltip />
      <el-table-column label="攻击结果" prop="attackResult" show-overflow-tooltip />
      <el-table-column label="CNNVD编号" prop="cnnvdId" show-overflow-tooltip width="180" />
      <el-table-column label="CVE编号" prop="cveId" show-overflow-tooltip width="180" />
      <el-table-column label="状态" show-overflow-tooltip>
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="queryPage.searchStr"
            :column="column"
            :filter-option="[
              { label: '启用', value: '启用' },
              { label: '停用', value: '停用' },
            ]"
            filter-type="select"
            :tools="['filter']"
          />
        </template>
        <template #default="{ row }">
          <span class="status" :class="{ on: row.enable === '启用' }">{{ row.enable }}</span>
        </template>
      </el-table-column>
      <el-table-column label="更新时间" prop="updateTime" show-overflow-tooltip sortable width="180">
        <template #default="{ row }">
          {{ formatNstime(row.updateTime, false) }}
        </template>
      </el-table-column>
      <el-table-column fixed="right" label="操作" width="80">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
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
      @current-change="handleGetMaintenanceList"
      @size-change="handleGetMaintenanceList"
    />
    <div v-if="showInfo" class="maintenance-info">
      <el-space :spacer="spacer">
        <el-button :auto-insert-space="false" :icon="Back" @click="showInfo = false">返回</el-button>
        <span class="maintenance-info-title">编辑信息</span>
      </el-space>
      <el-form ref="infoFormRef" class="infoForm" label-position="top" :model="infoForm" :rules="rules">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="规则名称" prop="ruleName">
              <el-input v-model="infoForm.ruleName" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="规则标签" prop="tags">
              <el-input v-model="infoForm.tags" />
            </el-form-item>
          </el-col>

          <el-col :span="6">
            <el-form-item label="规则ID" prop="id">
              <el-input v-model="infoForm.id" :disabled="!!infoForm.id" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="规则类型" prop="ruleType">
              <el-select v-model="infoForm.ruleType" style="width: 100%">
                <el-option v-for="item in maintenanceOptions?.ruleType" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="威胁等级" prop="severity">
              <el-select v-model="infoForm.severity" style="width: 100%">
                <el-option v-for="item in maintenanceOptions?.severity" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="攻击方向" prop="direction">
              <el-select v-model="infoForm.direction" style="width: 100%">
                <el-option v-for="item in maintenanceOptions?.direction" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
          </el-col>

          <el-col :span="6">
            <el-form-item label="攻击类型" prop="attackType">
              <el-select v-model="infoForm.attackType" style="width: 100%">
                <el-option v-for="item in maintenanceOptions?.attackType" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="攻击方式" prop="actionType">
              <el-select v-model="infoForm.actionType" style="width: 100%">
                <el-option v-for="item in maintenanceOptions?.actionType" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="攻击手段" prop="attackMethod">
              <el-select v-model="infoForm.attackMethod" style="width: 100%">
                <el-option v-for="item in maintenanceOptions?.attackMethod" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="攻击阶段" prop="killchain">
              <el-select v-model="infoForm.killchain" style="width: 100%">
                <el-option v-for="item in maintenanceOptions?.killchain" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
          </el-col>

          <el-col :span="6">
            <el-form-item label="攻击结果" prop="attackResult">
              <el-select v-model="infoForm.attackResult" style="width: 100%">
                <el-option v-for="item in maintenanceOptions?.attackResult" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="置信度" prop="confidence">
              <el-select v-model="infoForm.confidence" style="width: 100%">
                <el-option v-for="item in maintenanceOptions?.confidence" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="代码语言" prop="codeLanguage">
              <el-select v-model="infoForm.codeLanguage" style="width: 100%">
                <el-option v-for="item in maintenanceOptions?.codeLanguage" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="CNNVD编号" prop="cnnvdId">
              <el-input v-model="infoForm.cnnvdId" />
            </el-form-item>
          </el-col>

          <el-col :span="8">
            <el-form-item label="CVE编号" prop="cveId">
              <el-input v-model="infoForm.cveId" />
            </el-form-item>
          </el-col>

          <el-col :span="8">
            <el-form-item label="影响应用" prop="siteApp">
              <el-input v-model="infoForm.siteApp" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="影响系统" prop="serverType">
              <el-input v-model="infoForm.serverType" />
            </el-form-item>
          </el-col>

          <el-col :span="12">
            <el-form-item label="风险危害" prop="vulnHarm">
              <el-input v-model="infoForm.vulnHarm" resize="none" :rows="4" type="textarea" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="威胁详情" prop="detailInfo">
              <el-input v-model="infoForm.detailInfo" resize="none" :rows="4" type="textarea" />
            </el-form-item>
          </el-col>

          <el-col :span="12">
            <el-form-item label="解决方案" prop="bulletIn">
              <el-input v-model="infoForm.bulletIn" resize="none" :rows="4" type="textarea" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="威胁描述" prop="vulnDesc">
              <el-input v-model="infoForm.vulnDesc" resize="none" :rows="4" type="textarea" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="告警阈值">
              <el-space :size="5" wrap>
                每
                <el-input
                  v-model.number="alarmThresholdS.hour"
                  clearable
                  placeholder="请选择小时"
                  style="width: 100px"
                />
                小时
                <el-input
                  v-model.number="alarmThresholdS.minutes"
                  clearable
                  placeholder="请选择分钟"
                  style="width: 100px"
                />
                分钟
                <el-input
                  v-model.number="alarmThresholdS.second"
                  clearable
                  placeholder="请选择秒"
                  style="width: 100px"
                />
                秒内发生
                <el-input
                  v-model.number="infoForm.alarmThresholdTimes"
                  clearable
                  placeholder="请选择次数"
                  style="width: 80px"
                />
                次警告
              </el-space>
            </el-form-item>
          </el-col>

          <el-col :span="12">
            <el-form-item label="告警间隔">
              <el-space :size="5" wrap>
                每
                <el-input
                  v-model.number="alarmIntervalS.hour"
                  clearable
                  placeholder="请选择小时"
                  style="width: 100px"
                />
                小时
                <el-input
                  v-model.number="alarmIntervalS.minutes"
                  clearable
                  placeholder="请选择分钟"
                  style="width: 100px"
                />
                分钟
                <el-input
                  v-model.number="alarmIntervalS.second"
                  clearable
                  placeholder="请选择秒"
                  style="width: 100px"
                />
                秒告警一次
              </el-space>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="规则启用状态" prop="enable">
              <el-switch
                v-model="infoForm.enable"
                active-text="开启"
                active-value="启用"
                inactive-text="停用"
                inactive-value="停用"
                inline-prompt
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="&nbsp;">
              <div style="width: 100%; text-align: right">
                <el-button type="primary" @click="handleSubmit">确定</el-button>
                <el-button @click="showInfo = false">取消</el-button>
              </div>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .information-maintenance-container {
    h3 {
      margin-block: 0 0.5em;
    }
    .maintenance.el-table {
      height: calc(100vh - 120px);
      margin-top: 10px;
      .level {
        font-size: 14px;
        width: 58px;
        height: 24px;
        border-radius: 4px;
        display: block;
        line-height: 22px;
        text-align: center;
        font-size: 14px;
        background-color: #9d9aba;
        color: #fff;
        .el-icon {
          margin-right: -4px;
          vertical-align: -3px;
          font-size: 17px;
        }
        &-0 {
          background: #ca0a08;
        }
        &-1 {
          background: #ff2927;
        }
        &-2 {
          background: #ff7212;
        }
        &-3 {
          background: #ffbe36;
        }
      }
      .status {
        width: 40px;
        height: 24px;
        background: #9690b7;
        border-radius: 4px;
        display: block;
        font-weight: 500;
        font-size: 13px;
        color: #ffffff;
        text-align: center;
        line-height: 24px;
        &.on {
          background: #4abf57;
        }
      }
    }
    .maintenance-info {
      z-index: 999;
      position: absolute;
      inset: 0;
      padding: 20px 30px 0;
      background: #fff;
      &-title {
        font-weight: 500;
        font-size: 20px;
        color: #303133;
      }
      .infoForm {
        margin: 30px 110px 0;
      }
    }
  }
</style>
