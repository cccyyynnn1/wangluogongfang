<script lang="ts">
  export default {
    name: 'UpdateRules', //
  }
</script>

<script setup lang="ts">
  import type { FormInstance } from 'element-plus'
  import { updateWarnCustomerRulePApi } from '@/api-ecs/forword-rule'
  import { Back } from '@element-plus/icons-vue'
  import cloneDeep from 'lodash/cloneDeep'
  const $baseMessage: any = inject('$baseMessage')
  const props = defineProps<{
    showPages: boolean
    curData: any
    allSelectOption: any
  }>()

  const loading = ref(false)
  const formRef = ref<FormInstance>()
  // 表单数据
  const infoForm = ref<any>({
    module: 'ordinary',
    checkProtocol: 'http',
    enable: '启用',
    alarmThresholdS: 0,
    alarmIntervalS: 0,
    alarmThresholdTimes: 0,
  })
  // 普通模式规则Rule
  const ordinaryRule = {
    checkProtocol: [
      {
        required: true,
        message: '请选择检测协议',
        trigger: 'change',
      },
    ],
  }
  // 专家模式规则Rule
  const expertRule = {}
  //默认规则Rule
  const infoRules = {
    ruleName: [
      {
        required: true,
        message: '请输入规则名称',
        trigger: 'change',
      },
    ],
    ruleType: [
      {
        required: true,
        message: '请选择规则类型',
        trigger: 'change',
      },
    ],
    attackType: [
      {
        required: true,
        message: '请选择攻击类型',
        trigger: 'change',
      },
    ],
    severity: [
      {
        required: true,
        message: '请选择威胁等级',
        trigger: 'change',
      },
    ],
    direction: [
      {
        required: true,
        message: '请选择攻击方向',
        trigger: 'change',
      },
    ],
    actionType: [
      {
        required: true,
        message: '请选择攻击方式',
        trigger: 'change',
      },
    ],
    attackMethod: [
      {
        required: true,
        message: '请选择攻击手段',
        trigger: 'change',
      },
    ],
    killchain: [
      {
        required: true,
        message: '请选择攻击阶段',
        trigger: 'change',
      },
    ],
    attackResult: [
      {
        required: true,
        message: '请选择攻击结果',
        trigger: 'change',
      },
    ],
    confidence: [
      {
        required: true,
        message: '请选择置信度',
        trigger: 'change',
      },
    ],
    vulnHarm: [
      {
        required: true,
        message: '请输入风险危害',
        trigger: 'change',
      },
    ],
    detailInfo: [
      {
        required: true,
        message: '请输入威胁详情',
        trigger: 'change',
      },
    ],
    vulnDesc: [
      {
        required: true,
        message: '请输入威胁描述',
        trigger: 'change',
      },
    ],
    ruleContent: [
      {
        required: true,
        message: '请输入规则内容',
        trigger: 'change',
      },
    ],
  }
  //表单校验Rule
  const formRules = computed(() => {
    return infoForm.value.module === 'expert' ? infoRules : { ...infoRules, ...ordinaryRule }
  })

  const httpSelection = ref<{ label: string; type: 'common' | 'request' | 'response' }[]>([])
  const ruleConditionKeys = ref<string[]>([])
  const ruleConditionType = ref<'common' | 'request' | 'response'>('common')
  const ruleCondition = ref<{ field: string; value: string; error?: boolean }[]>([
    {
      field: '',
      value: '',
    },
  ])

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
  const emit = defineEmits<{
    (e: 'on-close-event'): void
    (e: 'on-reflash'): void
  }>()
  // 关闭回调
  const handleClose = () => {
    emit('on-close-event')
  }
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

  const formatHttpSelect = (httpSelect: { httpCommon: string[]; httpRequest: string[]; httpResponse: string[] }) => {
    const { httpCommon = [], httpRequest = [], httpResponse = [] } = httpSelect
    const _ruleConditionKeys = ruleCondition.value.map((i) => i.field)
    ruleConditionKeys.value = _ruleConditionKeys
    const handleFun = (data: string[], key: 'common' | 'request' | 'response') => {
      return data.map((i) => {
        if (_ruleConditionKeys.includes(i) && ruleConditionType.value === 'common') {
          ruleConditionType.value = key
        }
        return { label: i, type: key }
      })
    }
    return [
      handleFun(httpCommon, 'common'),
      handleFun(httpRequest, 'request'),
      handleFun(httpResponse, 'response'),
    ].flat(2)
  }

  const handelRuleConditionKeysClick = (
    select: { label: string; type: 'common' | 'request' | 'response' },
    index: number
  ) => {
    const { label, type } = select
    ruleConditionKeys.value[index] = label
    if (!['request', 'request'].includes(ruleConditionType.value)) {
      ruleConditionType.value = type
    }
  }

  const handleConditionValueBlur = (condition: { field: string; value: string; error?: boolean }) => {
    const { value } = condition
    condition.error = !value
  }

  const getSelectStatus = (select: { label: string; type: 'common' | 'request' | 'response' }) => {
    const { label, type } = select
    if (type === 'request' && ruleConditionType.value === 'response') return true
    if (type === 'response' && ruleConditionType.value === 'request') return true
    return false
  }

  // 提交表单
  const submitForm = async () => {
    await formRef.value?.validate(async (valid, fields) => {
      let condition_valid = false
      ruleCondition.value.forEach((condition, index) => {
        if (!condition.value || !condition.field || condition.value.length < 2 || condition.value.length > 64) {
          console.log(111111)
          ruleCondition.value[index].error = true
          condition_valid = true
        }
      })
      if (condition_valid && infoForm.value.module === 'ordinary') return
      if (infoForm.value.module === 'ordinary') {
        const _ruleContent = ruleCondition.value.reduce((obj, condition) => {
          const { field, value } = condition
          obj[field] = value
          return obj
        }, {} as { [key: string]: string })
        infoForm.value.ruleContent = JSON.stringify(_ruleContent)
      }
      if (valid) {
        loading.value = true
        infoForm.value.alarmThresholdS = formatTime2second(alarmThresholdS)
        infoForm.value.alarmIntervalS = formatTime2second(alarmIntervalS)
        try {
          const { msg } = await updateWarnCustomerRulePApi(infoForm.value)
          $baseMessage(msg, 'success', 'vab-hey-message-success')
          emit('on-reflash')
          handleClose()
        } finally {
          loading.value = false
        }
      }
    })
  }

  // 加一项
  const handleAddItem = () => {
    ruleCondition.value.push({ field: '', value: '' })
  }

  // 删除一行
  const handleDeleteItem = (index: number) => {
    ruleCondition.value.splice(index, 1)
  }
  watch(
    () => props.allSelectOption,
    () => {
      httpSelection.value = formatHttpSelect(props.allSelectOption?.httpSelect || {})
    },
    {
      deep: true,
      immediate: true,
    }
  )
  onMounted(() => {
    const propData = cloneDeep(props.curData)
    Object.assign(infoForm.value, propData)
    if (infoForm.value.module === 'ordinary') {
      if (infoForm.value.ruleContent) {
        const _ruleCondition = Object.entries(JSON.parse(infoForm.value.ruleContent)).map(([key, val]) => ({
          field: key,
          value: val as string,
        }))
        ruleCondition.value = infoForm.value.ruleContent ? _ruleCondition : []
      }
      infoForm.value.ruleContent = ''
    }
    Object.assign(alarmThresholdS, getTimes(infoForm.value.alarmThresholdS))
    Object.assign(alarmIntervalS, getTimes(infoForm.value.alarmIntervalS))
  })
</script>

<template>
  <div v-loading="loading" class="my-contain">
    <el-space :size="0">
      <div>
        <el-button :auto-insert-space="false" :icon="Back" @click="handleClose">返回</el-button>
        <el-divider direction="vertical" />
        <span class="info-title">{{ infoForm?.id ? '编辑规则' : '新增规则' }}</span>
      </div>
      <el-button :auto-insert-space="false" type="primary" @click="submitForm">保存</el-button>
    </el-space>
    <div class="contain">
      <el-form
        ref="formRef"
        label-position="top"
        label-width="100px"
        :model="infoForm"
        :rules="formRules"
        :validate-on-rule-change="false"
      >
        <el-row :gutter="20">
          <el-col :span="16">
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
                <el-form-item label="规则ID" prop="ruleId">
                  <el-input v-model="infoForm.ruleId" disabled />
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="规则类型" prop="ruleType">
                  <el-select v-model="infoForm.ruleType" style="width: 100%">
                    <el-option v-for="item in allSelectOption?.ruleType" :key="item" :label="item" :value="item" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="威胁等级" prop="severity">
                  <el-select v-model="infoForm.severity" style="width: 100%">
                    <el-option v-for="item in allSelectOption?.severity" :key="item" :label="item" :value="item" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="攻击方向" prop="direction">
                  <el-select v-model="infoForm.direction" style="width: 100%">
                    <el-option v-for="item in allSelectOption?.direction" :key="item" :label="item" :value="item" />
                  </el-select>
                </el-form-item>
              </el-col>

              <el-col :span="6">
                <el-form-item label="攻击类型" prop="attackType">
                  <el-select v-model="infoForm.attackType" style="width: 100%">
                    <el-option v-for="item in allSelectOption?.attackType" :key="item" :label="item" :value="item" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="攻击方式" prop="actionType">
                  <el-select v-model="infoForm.actionType" style="width: 100%">
                    <el-option v-for="item in allSelectOption?.actionType" :key="item" :label="item" :value="item" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="攻击手段" prop="attackMethod">
                  <el-select v-model="infoForm.attackMethod" style="width: 100%">
                    <el-option v-for="item in allSelectOption?.attackMethod" :key="item" :label="item" :value="item" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="攻击阶段" prop="killchain">
                  <el-select v-model="infoForm.killchain" style="width: 100%">
                    <el-option v-for="item in allSelectOption?.killchain" :key="item" :label="item" :value="item" />
                  </el-select>
                </el-form-item>
              </el-col>

              <el-col :span="6">
                <el-form-item label="攻击结果" prop="attackResult">
                  <el-select v-model="infoForm.attackResult" style="width: 100%">
                    <el-option v-for="item in allSelectOption?.attackResult" :key="item" :label="item" :value="item" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="置信度" prop="confidence">
                  <el-select v-model="infoForm.confidence" style="width: 100%">
                    <el-option v-for="item in allSelectOption?.confidence" :key="item" :label="item" :value="item" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="代码语言" prop="codeLanguage">
                  <el-select v-model="infoForm.codeLanguage" style="width: 100%">
                    <el-option v-for="item in allSelectOption?.codeLanguage" :key="item" :label="item" :value="item" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="CNNVD编号" prop="cnnvdId">
                  <el-input v-model="infoForm.cnnvdId" />
                </el-form-item>
              </el-col>

              <el-col :span="6">
                <el-form-item label="CVE编号" prop="cveId">
                  <el-input v-model="infoForm.cveId" />
                </el-form-item>
              </el-col>

              <el-col :span="6">
                <el-form-item label="影响应用" prop="siteApp">
                  <el-input v-model="infoForm.siteApp" />
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="影响系统" prop="serverType">
                  <el-input v-model="infoForm.serverType" />
                </el-form-item>
              </el-col>
              <el-col :span="6">
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
              <el-col :span="24">
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
                      style="width: 100px"
                    />
                    次警告
                  </el-space>
                </el-form-item>
              </el-col>

              <el-col :span="24">
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
            </el-row>
          </el-col>
          <el-col :span="8">
            <el-form-item class="radio-item" label="规则内容" prop="module" style="position: relative">
              <span v-if="infoForm.module === 'expert'" class="ruleText">( 仅支持输入单条规则 )</span>
              <el-radio-group
                v-model="infoForm.module"
                style="position: absolute; right: 0; bottom: 0"
                @change="infoForm.ruleContent = ''"
              >
                <el-radio label="ordinary">普通模式</el-radio>
                <el-radio label="expert">专家模式</el-radio>
              </el-radio-group>
            </el-form-item>
            <!-- 普通模式 -->
            <div v-if="infoForm.module === 'ordinary'" class="config">
              <el-form-item label="检测协议：" prop="checkProtocol" style="position: relative">
                <el-input v-model="infoForm.checkProtocol" disabled />
              </el-form-item>
              <el-form-item label="检测条件：">
                <el-timeline style="width: 100%">
                  <el-timeline-item
                    v-for="(condition, index) in ruleCondition"
                    :key="index"
                    :hollow="true"
                    style="position: relative"
                    :type="'primary'"
                  >
                    <span v-if="index != 0" style="position: absolute; color: #6954f0; left: -20px; top: 0">且</span>
                    <el-row :gutter="3" style="flex: 1">
                      <el-col
                        class="el-form-item"
                        :class="{ 'is-error': condition.error && !condition.field }"
                        :span="10"
                        style="position: relative"
                      >
                        <el-select v-model="condition.field" style="width: 100%">
                          <el-option
                            v-for="select in httpSelection"
                            :key="select.label"
                            :disabled="getSelectStatus(select)"
                            :label="select.label"
                            :value="select.label"
                            @click="handelRuleConditionKeysClick(select, index)"
                          />
                        </el-select>
                        <span v-if="condition.error && !condition.field" class="error-msg">请选择检测条件</span>
                      </el-col>
                      <el-col
                        class="el-form-item"
                        :class="{ 'is-error': condition.error && condition.field }"
                        :span="11"
                        style="padding-left: 10px; position: relative"
                      >
                        <el-input v-model="condition.value" @blur="handleConditionValueBlur(condition)" />
                        <span v-if="condition.error && condition.field" class="error-msg">
                          检测条件值长度必须为2-64位
                        </span>
                      </el-col>
                      <el-col :span="3">
                        <div style="display: flex; justify-content: start; height: 100%; align-items: center">
                          <el-icon color="#6954f0" style="margin-inline: 10px" @click="handleAddItem">
                            <Plus />
                          </el-icon>
                          <el-icon v-if="index !== 0" color="#6954f0" @click="handleDeleteItem(index)">
                            <Delete />
                          </el-icon>
                        </div>
                      </el-col>
                    </el-row>
                  </el-timeline-item>
                </el-timeline>
              </el-form-item>
            </div>
            <!-- 专家模式 -->
            <el-form-item v-else class="diyItem" label=" " prop="ruleContent">
              <el-input
                v-model="infoForm.ruleContent"
                placeholder='支持Snort、Suricata语法，例如：alert http any any -> any any (msg:"HTTP响应体匹配示例"; flow:established,to_client; http.response_body; content:"Server response body"; classtype:bad-unknown; sid:120; rev:1;)'
                resize="none"
                :rows="30"
                type="textarea"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .my-contain {
    height: 100%;
    & > .el-space {
      width: 100%;
      justify-content: space-between;
      .info-title {
        font-weight: 500;
        font-size: 20px;
        vertical-align: middle;
        color: #303133;
      }
    }
    :deep() {
      .radio-item.el-form-item {
        margin-bottom: 0px !important;
      }
    }
    .contain {
      height: calc(100vh - 75px);
      overflow-y: auto;
      overflow-x: hidden;
      margin-top: 20px;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
      .ruleText {
        position: absolute;
        left: 70px;
        bottom: 2px;
        font-size: 12px;
      }
      .config {
        padding: 10px 16px;
        border: 1px solid var(--el-border-color);
        border-radius: 2.5px;
        height: 636px;
        overflow-y: auto;
        .el-icon {
          cursor: pointer;
        }
      }

      :deep() {
        .diyItem {
          .el-form-item__label {
            display: none;
          }
          .el-textarea__inner {
            height: 636px;
          }
        }
        .el-radio-group {
          &:not(:last-child).el-radio {
            margin-right: 15px;
          }
        }
        .error-msg {
          color: var(--el-color-danger);
          font-size: 12px;
          line-height: 1;
          padding-top: 2px;
          position: absolute;
          top: 100%;
          left: 10px;
        }
        .el-timeline {
          padding-left: 30px;
        }
        .el-timeline-item__tail {
          margin-top: 20px;
          border-left: 1px dashed #6954f0 !important;
        }
        .el-timeline-item__node {
          top: 7px;
        }
        .el-timeline-item__timestamp {
          display: none;
        }
      }
    }
  }
</style>
