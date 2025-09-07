<script lang="ts">
  export default {
    name: 'FieldModules',
  }
</script>

<script setup lang="ts">
  import UpdataField from './updata-field.vue'
  import { getTemplateInfoApi, fieldsTemplateUpdateApi, fieldsTemplateAddApi } from '@/api-ecs/custom-field'

  const $baseMessage: any = inject('$baseMessage')
  const moduleType = ref<'site' | 'alarm' | 'retrieve'>('site')
  const editFields = ref()
  const editFieldsTitle = ref('')
  const showFieldDialog = ref(false) // 打开编辑页

  const props = withDefaults(
    defineProps<{
      visible: boolean
      templateData: any
      display?: 'preview' | 'default' | 'add'
    }>(),
    {
      display: 'default',
    }
  )

  const emits = defineEmits<{
    (e: 'update:visible', val: boolean): void
    (e: 'closeHandle'): void
  }>()
  const fieldEditData = ref<GOD>()
  const curIndexType = ref<number | undefined>(undefined)
  const infoLoading = ref(false)
  const drawerVisible = useVModel(props, 'visible', emits)
  const baseTemplateFieldsInfo = ref()

  const hanldeModuleTypeChange = (type: 'site' | 'alarm' | 'retrieve') => {
    moduleType.value = type
  }

  type GOD = {
    type: number
    siteSessionId?: number
    siteApiId?: number
    tag: 'alarm' | 'survey' | 'siteSession' | 'siteApi'
  }
  const handleEdit2Fields = (fields: any, data: GOD) => {
    editFieldsTitle.value = fields.typeName
    editFields.value = fields
    curIndexType.value = data.type
    fieldEditData.value = data
    showFieldDialog.value = true
  }

  const handleGetTemplateInfo = async (id: number) => {
    infoLoading.value = true
    try {
      const { data } = await getTemplateInfoApi({ id })
      // data.session = data.session.splice(0, 1)
      baseTemplateFieldsInfo.value = data
    } finally {
      nextTick(() => {
        infoLoading.value = false
      })
    }
  }
  const handleFieldsEdit = async (fieldNames: string[]) => {
    editFields.value.namesList = fieldNames
    editFields.value.names = fieldNames.toString()
    // 是预览的时候不能直接保存
    if (props.display === 'default') {
      const { msg } = await fieldsTemplateUpdateApi({
        ...props.templateData,
        ...baseTemplateFieldsInfo.value,
      })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
    }
    showFieldDialog.value = false
  }

  const handleTemplateAdd = async () => {
    const { remark, templateName } = props.templateData
    infoLoading.value = true
    try {
      const { msg } = await fieldsTemplateAddApi({
        remark,
        templateName,
        ...baseTemplateFieldsInfo.value,
      })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      emits('closeHandle')
      drawerVisible.value = false
    } finally {
      infoLoading.value = false
    }
  }
  watch(
    () => props.visible,
    () => {
      if (props.visible) {
        editFields.value = undefined
        moduleType.value = 'site'
        props.templateData?.id && handleGetTemplateInfo(props.templateData?.id)
      }
    },
    { deep: true }
  )
</script>

<template>
  <div class="field-modules-info">
    <el-drawer v-model="drawerVisible" destroy-on-close size="85%" title="模版详情">
      <template #default>
        <div v-loading="infoLoading" class="modules-info-content">
          <div class="modules-head">
            <el-button-group class="modules-type">
              <el-button
                :type="moduleType === 'site' ? 'primary' : 'default'"
                @click="() => hanldeModuleTypeChange('site')"
              >
                站点
              </el-button>
              <el-button
                :type="moduleType === 'retrieve' ? 'primary' : 'default'"
                @click="() => hanldeModuleTypeChange('retrieve')"
              >
                调查
              </el-button>
              <el-button
                :type="moduleType === 'alarm' ? 'primary' : 'default'"
                @click="() => hanldeModuleTypeChange('alarm')"
              >
                告警
              </el-button>
            </el-button-group>
          </div>

          <div class="modules-configuration" :style="{ height: display !== 'add' && 'calc(100% - 50px)' }">
            <!-- 站点字段详情 -->
            <template v-if="moduleType === 'site'">
              <div class="modules-configuration-item title-item">
                <span class="tree-indent">
                  <span class="tree-indent-unit"></span>
                </span>
                <span class="fields-name fields-title">{{ baseTemplateFieldsInfo?.session_default?.typeName }}</span>
              </div>
              <!-- 默认模版 -->
              <div class="modules-configuration-item title-item">
                <span class="tree-indent">
                  <span class="tree-indent-unit"></span>
                  <span class="tree-indent-unit"></span>
                </span>
                <span class="fields-name">默认模版</span>
              </div>
              <div class="modules-configuration-item">
                <span class="tree-indent">
                  <span class="tree-indent-unit"></span>
                  <span class="tree-indent-unit"></span>
                </span>
                <div class="fields">
                  <span v-for="sessionName in baseTemplateFieldsInfo?.session_default?.namesList" :key="sessionName">
                    {{ sessionName }}
                  </span>
                </div>
              </div>
              <!-- 会话列表详情 -->
              <template v-for="(_session, index_session) in baseTemplateFieldsInfo?.session" :key="_session?.id">
                <!-- 当前会话名字 -->
                <div
                  class="modules-configuration-item title-item"
                  :class="{ 'end-modules': index_session === baseTemplateFieldsInfo?.session?.length - 1 }"
                >
                  <span class="tree-indent">
                    <span class="tree-indent-unit"></span>
                    <span class="tree-indent-unit"></span>
                  </span>
                  <span class="fields-name">
                    <el-icon
                      v-if="display !== 'preview'"
                      @click="
                        () =>
                          handleEdit2Fields(_session, {
                            type: 30,
                            tag: 'siteSession',
                            siteSessionId: _session.type,
                          })
                      "
                    >
                      <Edit />
                    </el-icon>
                    {{ _session?.typeName }}
                  </span>
                </div>
                <!-- 当前会话展示字段 -->
                <div
                  class="modules-configuration-item"
                  :class="{
                    'end-sub-modules':
                      index_session === baseTemplateFieldsInfo?.session?.length - 1 && index_session !== 0,
                  }"
                >
                  <span class="tree-indent">
                    <span class="tree-indent-unit"></span>
                    <span
                      class="tree-indent-unit"
                      :class="{ disable: index_session === baseTemplateFieldsInfo?.session?.length - 1 }"
                    ></span>
                  </span>
                  <div class="fields">
                    <span v-for="name in _session?.namesList" :key="name">{{ name }}</span>
                  </div>
                </div>
                <!-- 当前会话下的API -->
                <template v-if="_session?.apiList?.length">
                  <template v-for="(_api, index_api) in _session?.apiList" :key="_api?.id">
                    <!-- 会话API名字 -->
                    <div
                      class="modules-configuration-item title-item"
                      :class="{
                        'end-sub-modules':
                          index_session === baseTemplateFieldsInfo?.session?.length - 1 && index_session !== 0,
                        'last-field': index_api === _session?.apiList?.length - 1,
                      }"
                    >
                      <span class="tree-indent">
                        <span class="tree-indent-unit"></span>
                        <span class="tree-indent-unit"></span>
                        <span
                          class="tree-indent-unit"
                          :class="{ disable: index_session === baseTemplateFieldsInfo?.session?.length - 1 }"
                        ></span>
                      </span>
                      <span class="fields-name">
                        <el-icon
                          v-if="display !== 'preview'"
                          @click="
                            () =>
                              handleEdit2Fields(_api, {
                                type: 30,
                                tag: 'siteApi',
                                siteSessionId: _session.type,
                                siteApiId: _api.type,
                              })
                          "
                        >
                          <Edit />
                        </el-icon>
                        {{ _api.typeName }}
                        <span v-if="_api?.apiUrl" style="color: #666">（{{ _api?.apiUrl }}）</span>
                      </span>
                    </div>
                    <!-- 会话API展示字段 -->
                    <div
                      class="modules-configuration-item"
                      :class="{
                        'end-sub-modules':
                          index_session === baseTemplateFieldsInfo?.session?.length - 1 && index_session !== 0,
                        'last-field': index_api === _session?.apiList?.length - 1,
                      }"
                    >
                      <span class="tree-indent">
                        <span class="tree-indent-unit"></span>
                        <span class="tree-indent-unit"></span>
                        <span
                          class="tree-indent-unit"
                          :class="{ disable: index_api === _session?.apiList?.length - 1 }"
                        ></span>
                      </span>
                      <div class="fields">
                        <span v-for="name in _api?.namesList" :key="name">{{ name }}</span>
                      </div>
                    </div>
                  </template>
                </template>
              </template>
              <!-- 站点API列表 -->
              <div class="modules-configuration-item end-modules title-item">
                <span class="tree-indent">
                  <span class="tree-indent-unit"></span>
                </span>
                <span class="fields-name fields-title">{{ baseTemplateFieldsInfo?.api_default?.typeName }}</span>
              </div>
              <div class="modules-configuration-item">
                <span class="tree-indent">
                  <span class="tree-indent-unit disable"></span>
                </span>
                <div class="fields">
                  <span v-for="api in baseTemplateFieldsInfo?.api_default?.namesList" :key="api">{{ api }}</span>
                </div>
              </div>
            </template>
            <!-- 调查详情 -->
            <template v-else-if="moduleType === 'retrieve'">
              <div class="modules-configuration-item title-item">
                <span class="tree-indent">
                  <span class="tree-indent-unit"></span>
                </span>
                <span class="fields-name fields-title">调查</span>
              </div>
              <!-- 列表详情 -->
              <template v-for="(retrieve, index_retrieve) in baseTemplateFieldsInfo?.survey" :key="retrieve?.id">
                <!-- 当前名字 -->
                <div
                  class="modules-configuration-item title-item"
                  :class="{ 'end-modules': index_retrieve === baseTemplateFieldsInfo?.survey?.length - 1 }"
                >
                  <span class="tree-indent">
                    <span class="tree-indent-unit"></span>
                    <span class="tree-indent-unit"></span>
                  </span>
                  <span class="fields-name">
                    <el-icon
                      v-if="display !== 'preview'"
                      @click="
                        () =>
                          handleEdit2Fields(retrieve, {
                            type: +retrieve.type,
                            tag: 'survey',
                          })
                      "
                    >
                      <Edit />
                    </el-icon>
                    {{ retrieve?.typeName }}
                  </span>
                </div>
                <!-- 当前展示字段 -->
                <div
                  class="modules-configuration-item"
                  :class="{
                    'end-sub-modules':
                      index_retrieve === baseTemplateFieldsInfo?.survey?.length - 1 && index_retrieve !== 0,
                  }"
                >
                  <span class="tree-indent">
                    <span class="tree-indent-unit"></span>
                    <span
                      class="tree-indent-unit"
                      :class="{ disable: index_retrieve === baseTemplateFieldsInfo?.survey?.length - 1 }"
                    ></span>
                  </span>
                  <div class="fields">
                    <span v-for="name in retrieve?.namesList" :key="name">{{ name }}</span>
                  </div>
                </div>
              </template>
            </template>
            <!-- 告警详情 -->
            <template v-else-if="moduleType === 'alarm'">
              <div class="modules-configuration-item title-item">
                <span class="tree-indent">
                  <span class="tree-indent-unit"></span>
                </span>
                <span class="fields-name fields-title">调查</span>
              </div>
              <!-- 列表详情 -->
              <template v-for="(alarm, index_alarm) in baseTemplateFieldsInfo?.alarm" :key="alarm?.id">
                <!-- 当前名字 -->
                <div
                  class="modules-configuration-item title-item"
                  :class="{ 'end-modules': index_alarm === baseTemplateFieldsInfo?.alarm?.length - 1 }"
                >
                  <span class="tree-indent">
                    <span class="tree-indent-unit"></span>
                    <span class="tree-indent-unit"></span>
                  </span>
                  <span class="fields-name">
                    <el-icon
                      v-if="display !== 'preview'"
                      @click="
                        () =>
                          handleEdit2Fields(alarm, {
                            type: +alarm.type,
                            tag: 'alarm',
                          })
                      "
                    >
                      <Edit />
                    </el-icon>
                    {{ alarm?.typeName }}
                  </span>
                </div>
                <!-- 当前展示字段 -->
                <div
                  class="modules-configuration-item"
                  :class="{
                    'end-sub-modules': index_alarm === baseTemplateFieldsInfo?.alarm?.length - 1 && index_alarm !== 0,
                  }"
                >
                  <span class="tree-indent">
                    <!-- <span class="tree-indent-unit"></span> -->
                    <span class="tree-indent-unit"></span>
                    <span
                      class="tree-indent-unit"
                      :class="{
                        disable: index_alarm === baseTemplateFieldsInfo?.alarm?.length - 1,
                      }"
                    ></span>
                  </span>
                  <div class="fields">
                    <span v-for="name in alarm?.namesList" :key="name">{{ name }}</span>
                  </div>
                </div>
              </template>
            </template>
          </div>
          <div v-if="display === 'add'" class="footer">
            <el-button type="primary" @click="handleTemplateAdd">保存</el-button>
            <el-button @click="drawerVisible = false">取消</el-button>
          </div>
        </div>
      </template>
    </el-drawer>
    <updata-field
      v-model="showFieldDialog"
      :field-edit-data="fieldEditData"
      :fields="editFields?.namesList"
      :retrieve-index-type="curIndexType"
      :title="editFieldsTitle"
      @handleok="handleFieldsEdit"
    />
  </div>
</template>

<style scoped lang="scss">
  .field-modules-info {
    :deep() {
      .el-drawer {
        border-radius: 30px 0px 0px 30px;
        background: #534b89;
        .el-drawer__header {
          padding: 12px 30px 12px;
          margin-bottom: 0;
          .el-drawer__title {
            font-size: 20px;
            font-weight: 500;
            font-size: 20px;
            line-height: 34px;
            color: #fff;
          }
          .el-drawer__close-btn {
            width: 26px;
            height: 26px;
            border-radius: 50%;
            background: #eeedf9;
            display: flex;
            justify-content: center;
            align-items: center;
          }
        }
        .el-drawer__body {
          border-radius: 30px 0px 0px 30px;
          background: #fff;
          padding-inline: 20px;
          .modules-info-content {
            height: 100%;
            .modules-head {
              position: sticky;
              top: 0;
              background-color: #fff;
              z-index: 99999;
              padding-block: 0 16px;
              .modules-type {
                &.el-button-group {
                  --el-border-radius-base: 6px;
                  --el-border-color: var(--el-color-primary);
                  .el-button.el-button--default {
                    --el-button-text-color: var(--el-color-primary);
                  }
                }
              }
            }
          }
          .footer {
            margin-top: 10px;
            text-align: right;
          }
        }
      }
    }
    .modules-configuration {
      height: calc(100% - 90px);
      overflow-y: auto;
      .modules-configuration-item {
        display: flex;
        align-items: flex-start;
        position: relative;
        --tree-title-height: 34px;
        --el-color-primary: #ccc;
        padding-bottom: 6px;
        .tree-indent {
          align-self: stretch;
          white-space: nowrap;
          user-select: none;
          margin-right: 4px;
          .tree-indent-unit {
            position: relative;
            height: 100%;
            display: inline-block;
            width: var(--tree-title-height);
            &::before {
              position: absolute;
              top: 0;
              inset-inline-end: 20px;
              bottom: -8px;
              border-inline-end: 1px solid var(--el-color-primary);
              content: '';
            }
            &.disable {
              --el-color-primary: #fff;
            }
          }
        }
        .fields {
          display: flex;
          border-radius: 4px;
          border: 1px solid #e6e6ed;
          flex: 1;
          padding: 8px 10px 0;
          flex-wrap: wrap;
          span {
            font-weight: 400;
            font-size: 13px;
            color: #6954f0;
            display: block;
            background: #f5f4ff;
            border-radius: 4px;
            padding: 4px 10px;
            margin-bottom: 8px;
            &:not(:last-child) {
              margin-right: 8px;
            }
          }
        }
        .fields-name {
          font-size: 14px;
          color: #3c394f;
          font-weight: 400;
          .el-icon {
            vertical-align: -2px;
            cursor: pointer;
          }
          &.fields-title {
            font-weight: 500;
          }
        }
        &.end-modules,
        &.end-sub-modules,
        &.last-field {
          .tree-indent-unit {
            &:last-child::before {
              bottom: 9px;
            }
          }
        }
        &.end-sub-modules {
          .tree-indent-unit {
            &:not(:first-child):not(:last-child) {
              --el-color-primary: #fff;
            }
          }
        }
        &.title-item .tree-indent-unit {
          &:last-child::after {
            position: absolute;
            width: 10px;
            height: calc(calc(var(--tree-title-height) / 2) - 4px);
            border-bottom: 1px solid var(--el-color-primary);
            content: '';
            left: 14px;
          }
        }
        &:first-child {
          .tree-indent-unit:last-child {
            &::before {
              top: 12px;
            }
          }
        }
        &:last-child {
          .tree-indent-unit:first-child {
            &::before {
              bottom: 90px;
            }
          }
          .tree-indent-unit:first-child::after {
            position: absolute;
            width: 10px;
            height: calc(calc(var(--tree-title-height) / 2) - 4px);
            border-bottom: 1px solid var(--el-color-primary);
            border-left: 1px solid var(--el-color-primary);
            content: '';
            left: 13px;
          }
        }
      }
    }
  }
</style>
