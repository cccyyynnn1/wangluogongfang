<script lang="ts">
  export default {
    name: 'SyslogOutputInfo',
  }
</script>

<script setup lang="ts">
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import SyslogOutputField from './syslog-output-field.vue'
  import JsonPreview from '@/components/json-preview.vue'
  import {
    getAllSyslogTemplateApi,
    saveOrUpdateTemplateApi,
    saveOrUpdateApplyTemplateApi,
  } from '@/api-ecs/log-data-output'
  import { SyslogTemplate } from '@/types/index'
  const $baseMessage: any = inject('$baseMessage')

  type BaseData = {
    id?: number | null
    templateId: number
    protocol: number
    logType: number
    status: number
    ip: string
    port: number
    outputNull: number
    ouputSource: number
    probeIps: number[]
    threatLevelIds: number[]
    threatTypeList: string[]
    attackStageIds: number[]
    fieldsEchoInfos: string
    configStr: string
  }
  const baseData: BaseData = {
    id: null,
    templateId: 1,
    protocol: 0,
    logType: 0,
    status: 1,
    ip: '',
    port: 80,
    outputNull: 0,
    ouputSource: 1,
    probeIps: [],
    threatLevelIds: [],
    threatTypeList: [],
    attackStageIds: [],
    fieldsEchoInfos: '{}',
    configStr: '{}',
  }

  const props = defineProps({
    callback: {
      type: Function,
      default: () => {},
    },
  })
  const showAlert = ref(false)
  const visible = ref(false)
  const fieldDrawerVisible = ref(false)
  const previewFieldsVisible = ref(false)
  const dialogTitle = ref('新增')
  const dialogValue = ref()
  const active = ref(0)
  const selectOptions = reactive<SyslogTemplate>({
    attackStageDict: [],
    templateList: [],
    threatLevelDict: [],
    threatTypeDict: [],
    flowDeviceList: [],
  })

  const templateName = ref('')
  const templateSaveVisible = ref(false)

  const fieldRef = ref<InstanceType<typeof SyslogOutputField>>()
  // 预览字段数据模版（为选择的字段填充value值预览）
  const previewObject = ref<{ [key: string]: any } | undefined>()
  // 选择的字段(根据此数据展示选择的字段)
  const selecFields = ref<{ [key: string]: string[] } | undefined>()
  const handleCheckedChange = (type: 'all' | 'invert') => {
    fieldRef.value?.filterFields(type)
  }
  const handleFieldSave = () => {
    previewObject.value = fieldRef.value?.getPreviewFieldData()
    selecFields.value = fieldRef.value?.getFields()
    templateDirtyChange()
    fieldDrawerVisible.value = false
  }
  // 获取输出日志内容下拉数据
  const getAllSyslogTemplateHandler = async () => {
    const { data } = await getAllSyslogTemplateApi()
    selectOptions.attackStageDict = data.attackStageDict
    selectOptions.templateList = data.templateList
    selectOptions.threatLevelDict = data.threatLevelDict
    selectOptions.threatTypeDict = data.threatTypeDict
    selectOptions.flowDeviceList = data.flowDeviceList
    if (dialogValue.value.templateId) templateChangeHandle(dialogValue.value.templateId)
  }
  const saveTemplateHandler = async (isToApply = false) => {
    const { outputNull, ouputSource, probeIps, threatLevelIds, threatTypeList, attackStageIds } = dialogValue.value
    const { msg, data } = await saveOrUpdateTemplateApi(
      JSON.stringify({
        templateName: templateName.value,
        configStr: {
          // 字段预览模版
          previewFields: { ...previewObject.value },
          outputNull,
          ouputSource,
          probeIps,
          threatLevelIds,
          threatTypeList,
          attackStageIds,
          // 选中的字段数据
          selecFields: selecFields.value,
        },
      })
    )
    selectOptions.templateList.push(data)
    dialogValue.value.templateId = data.id
    templateName.value = ''
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    templateSaveVisible.value = false
  }
  const templateChangeHandle = (val: any) => {
    try {
      const { configStr: curConfig } = selectOptions.templateList.find((item) => item.id === val)!
      const {
        threatTypeList,
        outputNull,
        ouputSource,
        previewFields,
        probeIps,
        threatLevelIds,
        attackStageIds,
        selecFields: _selecFields,
      } = JSON.parse(curConfig)
      dialogValue.value = {
        ...dialogValue.value,
        probeIps,
        threatLevelIds,
        threatTypeList,
        attackStageIds,
        outputNull,
        ouputSource,
      }
      selecFields.value = _selecFields
      previewObject.value = previewFields
      showAlert.value = false
    } catch (error) {
      console.error(error)
    }
  }
  const nextStep = () => {
    if (!dialogValue.value.ip || !dialogValue.value.port) {
      return $baseMessage(!dialogValue.value.ip ? 'IPV4不能为空' : '端口不能为空', 'error', 'vab-hey-message-error')
    }
    const ipv4Regex = /^(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})(\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})){3}$/
    if (dialogValue.value.ip && !ipv4Regex.test(dialogValue.value.ip)) {
      return $baseMessage('请输入合法IPV4地址', 'error', 'vab-hey-message-error')
    }
    active.value = 1
  }
  const templateToApply = async () => {
    const { probeIps, threatLevelIds, threatTypeList, attackStageIds } = dialogValue.value
    if (threatTypeList.length === 0) return $baseMessage('威胁类型不能为空', 'error', 'vab-hey-message-error')
    if (attackStageIds.length === 0) return $baseMessage('攻击阶段不能为空', 'error', 'vab-hey-message-error')
    if (threatLevelIds.length === 0) return $baseMessage('威胁等级不能为空', 'error', 'vab-hey-message-error')
    if (probeIps.length === 0) return $baseMessage('流量区域和探针不能为空', 'error', 'vab-hey-message-error')
    if (JSON.stringify(selecFields.value) === '{}')
      return $baseMessage('输出字段不能为空', 'error', 'vab-hey-message-error')
    const { msg, data } = await saveOrUpdateApplyTemplateApi({
      ...dialogValue.value,
      configStr: JSON.stringify(previewObject.value),
      fieldsEchoInfos: JSON.stringify(selecFields.value),
    })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    props.callback()
    visible.value = false
  }

  const templateDirtyChange = () => {
    dialogValue.value.templateId = null
    showAlert.value = true
  }

  watch(visible, () => {
    if (visible.value) getAllSyslogTemplateHandler()
  })

  defineExpose({
    showDialog: ({ title = '新增', infoValue = baseData }: { title?: string; infoValue?: BaseData }) => {
      dialogTitle.value = title
      dialogValue.value = infoValue
      selecFields.value = JSON.parse(infoValue.fieldsEchoInfos)
      previewObject.value = JSON.parse(infoValue.configStr)
      active.value = 0
      visible.value = true
    },
  })
</script>

<template>
  <vab-dialog v-model="visible" destroy-on-close :title="dialogTitle" width="1200">
    <el-steps
      :active="active"
      align-center
      finish-status="wait"
      process-status="finish"
      :space="340"
      style="justify-content: center"
    >
      <el-step title="输出日志类型和目标" />
      <el-step title="请选择输出日志内容" />
    </el-steps>
    <div class="info-content">
      <!-- step 1 -->
      <el-form v-if="active === 0" inline label-width="110px" style="width: 450px">
        <el-form-item label="日志类型:" prop="logType">
          <el-select v-model="dialogValue.logType" style="width: 340px">
            <el-option label="告警日志" :value="0" />
            <el-option disabled label="审计日志" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item class="is-required" label="对端网络地址:">
          <el-input v-model="dialogValue.ip" clearable placeholder="IPV4" style="width: 230px" />
          <el-input v-model="dialogValue.port" placeholder="端口" style="width: 100px; margin-left: 10px" />
        </el-form-item>
        <el-form-item class="is-required" label="网络传输协议:" prop="protocol">
          <el-select v-model="dialogValue.protocol" style="width: 340px">
            <el-option label="TCP" :value="0" />
            <el-option label="UDP" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item class="is-required" label="状态:" prop="status">
          <el-switch
            v-model="dialogValue.status"
            active-text="开启"
            :active-value="1"
            inactive-text="关闭"
            :inactive-value="0"
            inline-prompt
          />
        </el-form-item>
      </el-form>
      <!-- step 2 -->
      <el-form v-else inline label-width="112px" style="padding-left: 80px">
        <el-form-item label="日志类型:" prop="logType">
          <el-select v-model="dialogValue.logType" disabled style="width: 340px">
            <el-option label="告警日志" :value="0" />
            <el-option disabled label="审计日志" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="日志模版组:" prop="templateId" style="width: 580px">
          <el-select
            v-model="dialogValue.templateId"
            style="width: 340px"
            value-key="id"
            @change="templateChangeHandle"
          >
            <el-option
              v-for="template in selectOptions.templateList"
              :key="template.id"
              :label="template.templateName"
              :value="template.id"
            />
          </el-select>
          <el-link style="margin-left: 10px" type="primary" :underline="false" @click="templateSaveVisible = true">
            另存为日志组模版
          </el-link>
          <el-link
            v-if="showAlert"
            style="
              position: absolute;
              top: 30px;
              right: 0;
              display: inline-block;
              width: 117px;
              margin-top: 20px;
              margin-left: 10px;
              line-height: 24px;
              color: #ffba00;
            "
            type="warning"
            :underline="false"
          >
            模版内容已更改,如需恢复请重选模版,也可保存为新模版
          </el-link>
        </el-form-item>
        <el-form-item label="威胁类型:" prop="threatTypeList">
          <el-select
            v-model="dialogValue.threatTypeList"
            collapse-tags
            collapse-tags-tooltip
            :max-collapse-tags="3"
            multiple
            placeholder="请选择威胁类型"
            style="width: 340px"
            @change="templateDirtyChange()"
          >
            <el-option v-for="type in selectOptions.threatTypeDict" :key="type" :label="type" :value="type" />
          </el-select>
        </el-form-item>
        <el-form-item label="攻击阶段:" prop="attackStageIds">
          <el-select
            v-model="dialogValue.attackStageIds"
            collapse-tags
            collapse-tags-tooltip
            :max-collapse-tags="3"
            multiple
            placeholder="请选择攻击阶段"
            style="width: 340px"
            @change="templateDirtyChange()"
          >
            <el-option
              v-for="attack in selectOptions.attackStageDict"
              :key="attack.id"
              :label="attack.dictValue"
              :value="attack.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="威胁等级:" prop="threatLevelIds">
          <el-select
            v-model="dialogValue.threatLevelIds"
            collapse-tags
            collapse-tags-tooltip
            :max-collapse-tags="3"
            multiple
            placeholder="请选择威胁等级"
            style="width: 340px"
            @change="templateDirtyChange()"
          >
            <el-option
              v-for="threat in selectOptions.threatLevelDict"
              :key="threat.id"
              :label="threat.dictValue"
              :value="threat.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="流量区域和探针:" prop="probeIps">
          <el-tree-select
            v-model="dialogValue.probeIps"
            check-on-click-node
            collapse-tags
            collapse-tags-tooltip
            :data="selectOptions.flowDeviceList"
            multiple
            placeholder="请选择流量区域和探针"
            :render-after-expand="false"
            show-checkbox
            style="width: 340px"
            @change="templateDirtyChange()"
          />
        </el-form-item>
        <el-form-item label="当字段为空:" prop="outputNull">
          <el-radio-group v-model="dialogValue.outputNull" @change="templateDirtyChange()">
            <el-radio :label="0">输出Null</el-radio>
            <el-radio :label="1">不处理</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="输出数据来源:" prop="ouputSource" style="width: 550px">
          <el-radio-group v-model="dialogValue.ouputSource" @change="templateDirtyChange()">
            <el-radio :label="1">是</el-radio>
            <el-radio :label="0">否</el-radio>
          </el-radio-group>
          <span style="margin-left: 20px; font-size: 12px; font-weight: 400; color: #a9acb3">
            例：TUH[vX.Y.Z] 2024/1/10/ 12:30:33
          </span>
        </el-form-item>
        <el-form-item label="&nbsp" style="width: 100%">
          <div style="display: flex; align-items: flex-start; margin-top: 12px">
            <div class="field-box">
              <syslog-output-field :select-fields="selecFields" :show-select="false" show-select-template />
            </div>
            <el-space :size="20" style="flex-direction: column; margin-left: 20px">
              <el-button
                style="margin-right: 14px; margin-bottom: 20px"
                type="primary"
                @click="fieldDrawerVisible = true"
              >
                字段筛选
              </el-button>
              <el-link type="primary" :underline="false" @click="previewFieldsVisible = true">预览编辑后示例</el-link>
            </el-space>
          </div>
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
      <el-checkbox
        v-if="active === 1"
        v-model="templateSaveVisible"
        label="保存为日志组模版"
        style="margin-right: 35px; color: #333; vertical-align: middle"
      />
      <el-button type="default" @click="visible = false">取消</el-button>
      <el-button v-if="active === 0" type="primary" @click="nextStep">下一步</el-button>
      <el-button v-else type="primary" @click="templateToApply">应用</el-button>
    </template>

    <!-- 模版字段选择器 -->
    <el-drawer v-model="fieldDrawerVisible" destroy-on-close direction="rtl" size="640px" title="字段筛选">
      <template #default>
        <syslog-output-field ref="fieldRef" :select-fields="selecFields" />
      </template>
      <template #footer>
        <el-space style="align-items: center; float: left; height: 32px">
          <el-link type="primary" @click="handleCheckedChange('all')">全选</el-link>
          <el-link style="margin-left: 20px" type="primary" @click="handleCheckedChange('invert')">反选</el-link>
        </el-space>
        <el-button type="default" @click="fieldDrawerVisible = false">取消</el-button>
        <el-button type="primary" @click="handleFieldSave">保存</el-button>
      </template>
    </el-drawer>
    <!-- 模版预览器 -->
    <vab-dialog v-model="previewFieldsVisible" class="previewDialog" destroy-on-close title="解析日志示例" width="860">
      <div style="height: 570px; margin-bottom: 30px; overflow: auto; background-color: #f8f7ff">
        <json-preview :json-value="JSON.stringify(previewObject, null, 4)" />
      </div>
    </vab-dialog>
    <!-- 另存日志模版 -->
    <vab-dialog
      v-model="templateSaveVisible"
      class="templateDialog"
      destroy-on-close
      title="另存为日志模版"
      width="650"
    >
      <el-form inline label-width="110px">
        <el-form-item label="模版名称:" required>
          <el-input v-model="templateName" placeholder="请输入模版名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button
          type="default"
          @click="
            () => {
              templateSaveVisible = false
              templateName = ''
            }
          "
        >
          取消
        </el-button>
        <el-button type="primary" @click="saveTemplateHandler()">保存</el-button>
      </template>
    </vab-dialog>
  </vab-dialog>
</template>

<style scoped lang="scss">
  .info-content {
    display: flex;
    justify-content: center;
    padding-top: 38px;
  }
  .field-box {
    width: 688px;
    height: 250px;
  }
  .previewDialog {
    :deep() {
      .el-dialog__body {
        height: 630px;
        padding-inline: 30px !important;
      }
    }
  }
  .templateDialog {
    :deep() {
      .el-dialog__body {
        height: 165px;
        margin-top: 50px;
        border-top: 0 !important;
      }
      .el-dialog__header {
        padding-top: 35px;
      }
      .el-form-item__content {
        width: 460px;
      }
      .el-dialog__footer {
        padding-bottom: 40px !important;
      }
    }
  }
  :deep() {
    .el-radio {
      margin-right: 20px;
    }
    .el-checkbox__label,
    .el-radio__label {
      font-weight: normal;
    }
    .el-step__icon {
      width: 40px;
      height: 40px;
      font-size: 18px;
      font-weight: 600;
    }
    .el-step.is-horizontal .el-step__line {
      top: 19px;
    }
    .el-form-item {
      margin-right: 20px;
      margin-bottom: 20px;
    }
    .el-form-item__content {
      width: 340px;
    }
    .el-dialog__footer {
      &:empty {
        display: none;
      }
    }
    .el-dialog__body {
      height: 665px;
      padding-top: 30px !important;
      padding-bottom: 0 !important;
    }
    .el-drawer__header {
      height: 84px;
      padding: 0 40px;
      margin-bottom: 0;
      border-bottom: 1px solid var(--el-border-color);
      .el-drawer__title {
        font-size: 20px;
        font-weight: 500;
        color: #303133;
      }
    }
    .el-drawer__body,
    .el-drawer__footer {
      padding: 20px 55px 0;
    }
    .el-drawer__footer {
      padding-bottom: 20px;
    }
  }
</style>
