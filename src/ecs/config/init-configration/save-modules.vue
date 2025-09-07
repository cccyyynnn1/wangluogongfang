<script lang="ts">
  export default {
    name: 'SaveModules',
  }
</script>

<script setup lang="ts">
  import { Plus } from '@element-plus/icons-vue'
  import AddFieldModule from '@/ecs/config/init-configration/add-field-module.vue'
  import VueDraggable from 'vuedraggable'
  import {
    getComparisonFieldsTemplateApi,
    getComparisonFieldsApi,
    deleteFieldsTemplateApi,
    contrastReplacementFieldsApi,
  } from '@/api-ecs/custom-field'
  import { InitialisationItem, tableSearch } from '@/types'
  import { formatTime } from '@/utils/time'
  const $baseMessage: any = inject('$baseMessage')
  const $baseConfirm: any = inject('$baseConfirm')
  const showAddFieldDialog = ref(false) // 打开编辑页
  const fieldInfovVisible = ref(false)

  const loading = ref(true)
  const virtualRef = ref()
  const replacementTooltipVisible = ref(false)

  const discrepancyFields = ref<string[]>([])
  const props = withDefaults(
    defineProps<{
      modelValue: boolean
      replacementField: string[]
      comparisonData?: {
        indexFiledsTemId?: number
        type?: number
        siteSessionId?: number
        siteApiId?: number
        tag: 'alarm' | 'survey' | 'siteSession' | 'siteApi'
      }
    }>(),
    {
      modelValue: false,
    }
  )
  const emits = defineEmits<{
    (e: 'update:modelValue', visible: boolean): void
    (e: 'show-info', showVal?: any): void
  }>()
  provide(tableSearch, () => {
    handleGetAllTemplate()
  })

  const templateComparisonFields = ref<string[]>([])
  const templateComparison = ref<InitialisationItem>()
  const checkedItem = ref<InitialisationItem>()
  const templateList = ref<InitialisationItem[]>([])
  const detailVisible = useVModel(props, 'modelValue', emits)
  const handleShowFielInfo = (moduleVal?: any) => {
    fieldInfovVisible.value = true
  }
  const handleGetAllTemplate = async () => {
    const { data } = await getComparisonFieldsTemplateApi()
    templateList.value = data || []
  }
  const checkboxChange = (value: boolean, data: InitialisationItem) => {
    checkedItem.value = value ? data : undefined
  }
  const handleReplacement = async () => {
    if (!checkedItem.value) return $baseMessage('请选择要替换的模版', 'error', 'vab-hey-message-error')
    const { msg } = await contrastReplacementFieldsApi({
      ...props.comparisonData,
      indexFiledsTemId: checkedItem.value?.id,
      names: props.replacementField.toString(),
    } as any)
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    detailVisible.value = false
  }
  const getReplacementFields = async (replacement: InitialisationItem, e: MouseEvent) => {
    const target = e.target as HTMLElement
    virtualRef.value = target.nodeName === 'DIV' ? target : target.parentElement
    templateComparison.value = replacement
    loading.value = true
    try {
      const { data } = await getComparisonFieldsApi({
        ...props.comparisonData,
        indexFiledsTemId: replacement.id,
      } as any)
      templateComparisonFields.value = data.namesList || []
      discrepancyFields.value = templateComparisonFields.value
        .filter((i) => !props.replacementField.includes(i))
        .concat(props.replacementField.filter((i) => !templateComparisonFields.value.includes(i)))
      replacementTooltipVisible.value = true
    } finally {
      loading.value = false
    }
  }

  const handleClickOutside = (e: any) => {
    const target = e.target as HTMLElement
    if (!target.classList.contains('compare') && target.nodeName !== 'svg') {
      replacementTooltipVisible.value = false
    }
  }

  const handleDeleteTemplate = async (ids: number[]) => {
    $baseConfirm('你确定要删除当前项吗？', null, async () => {
      const { msg } = await deleteFieldsTemplateApi(ids)
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      handleGetAllTemplate()
    })
  }
  watch(
    () => detailVisible.value,
    (val) => {
      if (val) {
        handleGetAllTemplate()
      } else {
        templateComparison.value = undefined
        checkedItem.value = undefined
      }
    }
  )
</script>

<template>
  <div class="save-modules-content">
    <el-drawer v-model="detailVisible" destroy-on-close size="85%">
      <template #header="{ titleId, titleClass }">
        <h4 :id="titleId" :class="titleClass">另存为模版</h4>
        <el-button class="add2custom" :icon="Plus" @click="() => (showAddFieldDialog = true)">新建自定义</el-button>
      </template>
      <div class="template-list" @click="handleClickOutside" @scroll="() => (replacementTooltipVisible = false)">
        <el-row :gutter="20">
          <el-col v-for="template in templateList" :key="template.id" :span="6">
            <div class="modules-item">
              <div class="info">
                <div class="module-title">
                  <el-checkbox
                    :disabled="template.templateName === '系统默认模版'"
                    :model-value="checkedItem?.id === template.id"
                    style="vertical-align: middle"
                    @change="checkboxChange($event, template)"
                  />
                  <el-tooltip :content="template.templateName" effect="dark" placement="top" popper-class="diyTooltip">
                    <span class="module-name">{{ template.templateName }}</span>
                  </el-tooltip>
                  <span class="module-status" :class="{ ok: template.status === 1 }">
                    {{ template.status === 1 ? '已生效' : '未生效' }}
                  </span>
                </div>
                <div class="module-msg">
                  <p>创建时间：{{ formatTime(template.createTime) }}</p>
                  <p>修改时间：{{ formatTime(template.updateTime) }}</p>
                  <p>备&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;注：{{ template.remark }}</p>
                </div>
                <div class="module-footer">
                  <el-button
                    :disabled="template.templateName === '系统默认模版'"
                    plain
                    @click="handleDeleteTemplate([template.id])"
                  >
                    删除
                  </el-button>
                </div>
              </div>
              <div class="compare" @click="getReplacementFields(template, $event)">
                <svg
                  id="图层_1"
                  version="1.1"
                  viewBox="0 0 16 14"
                  x="0px"
                  xml:space="preserve"
                  xmlns="http://www.w3.org/2000/svg"
                  xmlns:xlink="http://www.w3.org/1999/xlink"
                  y="0px"
                >
                  <g id="配置-字段初始化">
                    <g id="配置-另存为3" transform="translate(-609, -168)">
                      <g id="编组-5" transform="translate(220, 0)">
                        <g id="资产数据量" transform="translate(30, 84)">
                          <g id="编组-6" transform="translate(353, 84)">
                            <g id="duibi" transform="translate(6, 0)">
                              <path
                                id="形状"
                                class="st0"
                                d="M7.2,0h1.6v14H7.2V0z M0,3.1v7.8c0,0.9,0.7,1.6,1.6,1.6h4.8v-1.6H1.6V3.1h4.8V1.6H1.6
							C0.7,1.6,0,2.3,0,3.1z M9.6,1.6h1.6v1.6H9.6V1.6z M9.6,10.9h1.6v1.6H9.6V10.9z M14.4,1.6h-1.6v1.6h1.6v1.6H16V3.1
							C16,2.3,15.3,1.6,14.4,1.6L14.4,1.6z M14.4,10.9h-1.6v1.6h1.6c0.9,0,1.6-0.7,1.6-1.6V9.3h-1.6V10.9z M14.4,6.2H16v1.6h-1.6
							V6.2z"
                              />
                            </g>
                          </g>
                        </g>
                      </g>
                    </g>
                  </g>
                </svg>
                对比
              </div>
            </div>
          </el-col>
        </el-row>
        <el-popover
          placement="right-start"
          popper-class="singleton-tooltip"
          :popper-options="{
            modifiers: [
              {
                name: 'computeStyles',
                options: {
                  adaptive: false,
                  enabled: false,
                },
              },
            ],
          }"
          :show-arrow="false"
          :teleported="false"
          trigger="focus"
          :virtual-ref="virtualRef"
          virtual-triggering
          :visible="replacementTooltipVisible"
          width="600"
        >
          <el-form label-position="top">
            <el-form-item prop="roleName">
              <template #label>
                <span>已选展示字段：</span>
              </template>
              <div class="space">
                <vue-draggable item-key="id" :list="replacementField" :sort="false">
                  <template #item="{ element }">
                    <el-tag :class="{ isdiscrepancyField: discrepancyFields.includes(element) }">{{ element }}</el-tag>
                  </template>
                </vue-draggable>
              </div>
            </el-form-item>
            <el-form-item v-loading="loading" :label="`${templateComparison?.templateName}字段：`" prop="roleName">
              <div class="space">
                <vue-draggable item-key="id" :list="templateComparisonFields" :sort="false">
                  <template #item="{ element }">
                    <el-tag :class="{ isdiscrepancyField: discrepancyFields.includes(element) }">{{ element }}</el-tag>
                  </template>
                </vue-draggable>
              </div>
            </el-form-item>
          </el-form>
        </el-popover>
      </div>
      <div class="footer">
        <el-button type="primary" @click="handleReplacement">替换</el-button>
        <el-button @click="detailVisible = false">取消</el-button>
      </div>
    </el-drawer>
    <add-field-module v-model="showAddFieldDialog" @show-info="handleShowFielInfo" />
  </div>
</template>

<style scoped lang="scss">
  .save-modules-content {
    .modules-item {
      height: 205px;
      border-radius: 16px;
      border: 1px solid #ebe9fa;
      display: flex;
      overflow: hidden;
      margin-bottom: 20px;
      .info {
        width: calc(100% - 54px);
        .module-title {
          margin-block: 14px 10px;
          padding-left: 20px;
          .module-name {
            display: inline-block;
            width: calc(100% - 95px);
            font-weight: 500;
            font-size: 18px;
            color: #342e58;
            text-overflow: ellipsis;
            white-space: nowrap;
            overflow: hidden;
            vertical-align: middle;
            margin-left: 8px;
          }
          .module-status {
            font-weight: 400;
            font-size: 13px;
            margin-left: 8px;
            &::before {
              content: ' ';
              display: inline-block;
              width: 10px;
              height: 10px;
              border-radius: 50%;
              background: #c5c2d6;
              margin-right: 5px;
            }
            &.ok {
              color: #332e55;
              &::before {
                background: #5ad85d;
              }
            }
          }
        }
        .module-msg {
          text-indent: 44px;
          font-weight: 400;
          font-size: 14px;
          color: #9d9baa;
          margin-bottom: 15px;
          p {
            margin-bottom: 6px;
            text-overflow: ellipsis;
            white-space: nowrap;
            overflow: hidden;
            vertical-align: middle;
          }
        }
        .module-footer {
          height: 56px;
          border-top: 1px solid #ebe9fa;
          text-indent: 44px;
          .el-button {
            margin-top: 12px;
            border-radius: 6px;
          }
        }
      }
      .compare {
        width: 54px;
        background: #f1eeff;
        display: flex;
        flex-direction: column;
        font-weight: 500;
        font-size: 14px;
        color: var(--el-color-primary);
        justify-content: center;
        align-items: center;
        cursor: pointer;
        svg {
          fill: var(--el-color-primary);
          width: 16px;
          margin-bottom: 6px;
        }
      }
    }
    :deep() {
      .el-drawer {
        border-radius: 30px 0px 0px 30px;
        background: #534b89;
        .el-drawer__header {
          padding: 12px 30px 12px;
          margin-bottom: 0;
          .add2custom {
            color: var(--el-color-primary);
            background-color: #fff;
            margin-right: 20px;
            border-radius: 6px;
            font-weight: 500;
            font-size: 14px;
          }
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
          padding: 20px;
          .modules-type.el-button-group {
            --el-border-radius-base: 6px;
            --el-border-color: var(--el-color-primary);
            .el-button.el-button--default {
              --el-button-text-color: var(--el-color-primary);
            }
          }
          .template-list {
            height: calc(100% - 42px);
            overflow-y: auto;
            overflow-x: hidden;
          }
          .footer {
            margin-top: 10px;
            text-align: right;
            display: block;
          }
          .singleton-tooltip {
            background: #ffffff;
            box-shadow: 0px 0px 4px 4px rgba(77, 81, 86, 0.05);
            border-radius: 10px;
            border: 1px solid #f1f0ff;
            transition: transform 0.3s cubic-bezier(0.23, 1, 0.32, 1);
          }
          .space {
            width: 100%;
            min-height: 160px;
            max-height: 280px;
            overflow-y: auto;
            padding: 10px 5px;
            background: linear-gradient(180deg, #ffffff 0%, #ffffff 100%);
            border-radius: 4px;
            border: 1px solid #ebe9fa;
            .el-tag {
              cursor: pointer;
              &.isdiscrepancyField {
                color: #ff3232;
                background: #ffe9e9;
                border-color: #ffe9e9;
              }
            }
          }
        }
      }
    }
  }
</style>
