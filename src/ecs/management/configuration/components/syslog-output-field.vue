<script lang="ts">
  export default {
    name: 'SyslogOutputField',
  }
</script>

<script setup lang="ts">
  import { CheckboxValueType } from 'element-plus/es/components/checkbox/src/checkbox'
  import { sysLogAlertField, LogFieldType } from '@/data/constant'
  import { PropType } from 'vue'
  import _ from 'lodash'
  let isAllChecked = false
  const props = defineProps({
    selectFields: {
      type: Object as PropType<{
        [key: string]: string[]
      }>,
      default: () => ({}),
    },
    showSelect: {
      type: Boolean,
      default: true,
    },
    // 展示选择字段模版，false是展示全量模版
    showSelectTemplate: {
      type: Boolean,
      default: false,
    },
  })
  // 展示字段 数据源
  const curPreviewFieldData = ref()

  // 告警模版所有字段
  const allTemplateFields = getAllTemplateFields()
  // 展示子属性
  const showItems = ref<Set<string>>(new Set())
  // 选中需要展示的字段
  const selectFields = ref<Map<string, string[]>>(new Map(Object.entries(_.cloneDeep(props.selectFields))))

  // 折叠子属性
  const handleFold = (field: LogFieldType) => {
    showItems.value.has(field.id) ? showItems.value.delete(field.id) : showItems.value.add(field.id)
  }
  const checkboxChange = (isChecked: CheckboxValueType, field: LogFieldType, subField?: LogFieldType) => {
    if (!props.showSelect) return
    const isParent = !subField
    const parentInfo = selectFields.value.has(field.id)
      ? selectFields.value.get(field.id)
      : selectFields.value.set(field.id, []).get(field.id)
    const isCheckedHandle = () => {
      const subFields = sysLogAlertField.find((item) => item.id === field.id)
      isParent
        ? selectFields.value.set(field.id, subFields?.children?.map((item) => item.id) || [])
        : parentInfo?.push(subField.id)
      isParent ? selectFields.value.set(field.id, subFields?.children?.map((item) => item.id) || []) : null
    }
    const unCheckedHandle = () => {
      isParent ? selectFields.value.delete(field.id) : parentInfo?.splice(parentInfo?.indexOf(subField.id), 1)
      if (parentInfo?.length === 0) selectFields.value.delete(field.id)
    }
    isChecked ? isCheckedHandle() : unCheckedHandle()
  }
  // 判断父节点是否选中
  const getParentFieldStatus = (field: LogFieldType) => {
    return selectFields.value.has(field.id)
  }
  // 判断父节点是否全选
  const getParentFieldIsIndeterminate = (field: LogFieldType) => {
    const checkedSize = selectFields.value.get(field.id)?.length || 0
    const hasChild = field.children
    const childrenSize = field.children?.length || 0
    return hasChild ? checkedSize > 0 && checkedSize < childrenSize : false
  }
  // 判断子节点是否选中
  const getSubFieldsStatus = (field: LogFieldType, subField: LogFieldType) => {
    return (selectFields.value.get(field.id) || []).includes(subField.id)
  }
  // 获取选择的字段
  const getSelectFields = () => {
    return [...selectFields.value.entries()].reduce(
      (obj: { [key: string]: string[] }, [key, val]) => ((obj[key] = val), obj),
      {}
    )
  }
  // 获取所有的告警模版字段
  function getAllTemplateFields() {
    return sysLogAlertField.reduce(
      (obj: { [key: string]: string[] }, { id, children }) => (
        (obj[id] = children ? children.map((i) => i.id) : []), obj
      ),
      {}
    )
  }
  // 获取模版展示数据
  function getSelectTemplateData(): LogFieldType[] {
    const { selectFields, showSelectTemplate } = props
    const templateData = Object.entries(selectFields).map(([key, value]) => {
      const _templateData = sysLogAlertField.find((i) => i.id === key)!
      return {
        ..._templateData,
        children: _templateData.children ? _templateData?.children?.filter((i) => value.includes(i.id)) : null,
      }
    })
    if (templateData.length === 1) showItems.value.add(templateData[0].id)
    return showSelectTemplate ? templateData : sysLogAlertField
  }
  // 获取模版预览数据
  function getPreviewFieldData(): {
    [key: string]: any
  } {
    const previewFields = getSelectFields()
    return Object.entries(previewFields).reduce((preview_obj: { [key: string]: any }, [key, value]) => {
      const _templateData = sysLogAlertField.find((i) => i.id === key)!
      const children = value.map((id) => _templateData.children?.find((i) => i.id === id)) as LogFieldType[]
      const childrenVal = children.reduce((child_obj: { [key: string]: any }, { label, value }) => {
        child_obj[label] = value
        return child_obj
      }, {})
      preview_obj[key] = value.length > 0 ? childrenVal : _templateData.value
      return preview_obj
    }, {})
  }
  function handleInvertOrAll(type: 'invert' | 'all') {
    if (type === 'all') {
      if (isAllChecked) return
      selectFields.value = new Map(Object.entries(allTemplateFields))
      isAllChecked = true
    } else {
      if (isAllChecked) {
        selectFields.value.clear()
        isAllChecked = false
      } else {
        if (selectFields.value.size === 0) {
          selectFields.value = new Map(Object.entries(allTemplateFields))
          return (isAllChecked = true)
        }
        Object.entries(allTemplateFields).forEach(([key, value]) => {
          const curFeild = selectFields.value.get(key)
          if (curFeild) {
            if (curFeild?.length) {
              const differenceKey = _.difference(value, curFeild)
              differenceKey.length
                ? selectFields.value.set(key, _.difference(value, curFeild))
                : selectFields.value.delete(key)
            } else {
              selectFields.value.delete(key)
            }
          } else {
            selectFields.value.set(key, value)
          }
        })
      }
    }
  }
  watchEffect(() => {
    selectFields.value = new Map(Object.entries(props.selectFields))
    nextTick(() => (curPreviewFieldData.value = getSelectTemplateData()))
  })
  onMounted(() => {
    curPreviewFieldData.value = getSelectTemplateData()
  })
  defineExpose({
    getFields: getSelectFields,
    getPreviewFieldData,
    filterFields: handleInvertOrAll,
  })
</script>

<template>
  <div class="field-preview" :class="!showSelect && 'unchecked'">
    <template v-for="field in curPreviewFieldData" :key="field.id">
      <div class="field-item">
        <div class="field-key">
          <el-checkbox
            :indeterminate="getParentFieldIsIndeterminate(field)"
            :model-value="getParentFieldStatus(field)"
            @change="(val) => checkboxChange(val, field)"
          />
          {{ field.label }}
        </div>
        <div class="field-describe" :class="showItems.has(field.id) ? 'active' : ''" @click="() => handleFold(field)">
          {{ field.describe }}
          <span v-if="field.children" class="arrow"></span>
        </div>
      </div>
      <template v-if="field.children && showItems.has(field.id)">
        <div v-for="subfield in field.children" :key="subfield.id" class="field-item subfield">
          <div class="field-key">
            <el-checkbox
              :model-value="getSubFieldsStatus(field, subfield)"
              @change="(val) => checkboxChange(val, field, subfield)"
            />
            {{ subfield.label }}
          </div>
          <div class="field-describe">
            {{ subfield.describe }}
            <span v-if="subfield.children" class="arrow"></span>
          </div>
        </div>
      </template>
    </template>
  </div>
</template>

<style scoped lang="scss">
  .field-preview {
    width: 100%;
    height: 100%;
    border: 1px solid #ebeef5;
    overflow-y: auto;
    &.unchecked {
      :deep() {
        .el-checkbox__inner {
          display: none;
        }
      }
    }
    .field-item {
      display: flex;
      justify-content: space-between;
      line-height: 36px;
      font-size: 14px;
      font-weight: 400;
      color: #747983;
      height: 36px;
      text-indent: 20px;
      &:nth-of-type(even) {
        background: #f8f7ff;
      }
      .field-key {
        width: 250px;
        color: #303133;
        position: relative;
        cursor: pointer;
        :deep() {
          .el-checkbox {
            vertical-align: middle;
            margin-right: 8px;
            width: 14px;
          }
        }
      }
      .field-describe {
        position: relative;
        color: #747983;
        flex-grow: 1;
        pointer-events: none;

        &.active {
          .arrow {
            border-top-color: transparent;
            border-bottom-color: #d8dce6;
            top: 8px;
          }
        }
        .arrow {
          width: 0;
          height: 0;
          border: 8px solid transparent;
          border-top-color: #d8dce6;
          content: ' ';
          display: block;
          position: absolute;
          top: 16px;
          right: 15px;
          cursor: pointer;
          pointer-events: auto;
          transition: all 0.3s;
          transform-origin: 8px 0;
        }
      }
    }
    .subfield {
      .el-checkbox {
        margin-left: 26px;
      }
    }
  }
</style>
