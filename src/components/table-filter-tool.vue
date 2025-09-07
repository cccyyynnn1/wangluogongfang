<script lang="ts">
  export default {
    name: 'TableFilterTool',
  }
  type Tools = 'filter' | 'statistic'
  interface Props {
    column: any
    /**
     * @param  tools  filter 过滤  statistic 统计
     */
    tools: Tools[]
    filterValue?: any
    filterType?: 'text' | 'select' | 'date-range'
    filterOption?: { label: string; value: string | number }[]
  }
</script>

<script setup lang="ts">
  import { tableSearch } from '@/types'
  const popoverRef = ref()
  const props = withDefaults(defineProps<Props>(), {
    tools: () => ['filter', 'statistic'],
    filterType: 'text',
    filterOption: () => [],
    filterValue: undefined,
  })
  const handleTableSearch = inject(tableSearch)
  const emits = defineEmits<{
    (e: 'update:filter-value', val: any): void
  }>()

  const showTools = ref(true)
  const toolType = ref<Tools>('filter')
  const filterVal = ref(props.filterValue)
  const visible = ref(false)

  function clearFilterSearch() {
    filterVal.value = undefined
    emits('update:filter-value', undefined)
  }
  function handleFilterSearch() {
    emits('update:filter-value', filterVal.value)
    handleTableSearch?.()
    unref(popoverRef).popperRef?.triggerRef.click()
  }
  function handleToolsClick() {
    if (props.tools.length === 1) {
      showTools.value = false
    }
    visible.value = !visible.value
  }
  function handleToolsBtnClick(type: Tools) {
    toolType.value = type
    showTools.value = false
  }
  function popoverEnter() {
    filterVal.value = props.filterValue
  }
</script>

<template>
  <el-popover
    ref="popoverRef"
    :hide-after="0"
    placement="bottom-start"
    :show-arrow="false"
    trigger="click"
    width="auto"
    @before-enter="popoverEnter"
  >
    <template #reference>
      <div class="table-filter-tool" @click="handleToolsClick">
        <span class="title">{{ column.label }}</span>
        <vab-icon class="filter-icon" icon="filter-line" />
      </div>
    </template>
    <template #default>
      <div v-if="showTools" class="tool-btn">
        <div class="tool-btn-item" @click="() => handleToolsBtnClick('statistic')">字段统计</div>
        <div class="tool-btn-item" @click="() => handleToolsBtnClick('filter')">筛选</div>
      </div>
      <template v-else>
        <div v-if="toolType === 'filter'" class="tools-content filter">
          <p>筛选</p>
          <el-input v-if="filterType === 'text'" v-model="filterVal" placeholder="请输入">
            <template #prefix>
              <vab-icon icon="search-line" />
            </template>
          </el-input>
          <el-select
            v-if="filterType === 'select'"
            v-model="filterVal"
            placeholder="请选择"
            style="width: 100%"
            :teleported="false"
          >
            <el-option
              v-for="option in filterOption"
              :key="option.value"
              :label="option.label"
              placeholder="请选择"
              :value="option.value"
            />
          </el-select>
          <vab-date-time-picker v-if="filterType === 'date-range'" v-model="filterVal" />
          <div class="filter-footer">
            <el-button @click="clearFilterSearch">重置筛选条件</el-button>
            <el-button type="primary" @click="handleFilterSearch">检索</el-button>
          </div>
        </div>
        <div v-else-if="toolType === 'statistic'" class="tools-content statistic"></div>
      </template>
    </template>
  </el-popover>
</template>
<style lang="scss">
  .table-filter-tool {
    cursor: pointer;
  }
  .filter.tools-content {
    width: 372px;
    padding: 8px;
    p {
      margin-bottom: 1em;
    }
    .filter-footer {
      margin-top: 30px;
      display: flex;
      justify-content: space-between;
    }
  }
  .tool-btn-item {
    line-height: 22px;
    padding: 5px 16px;
    cursor: pointer;
    &:hover {
      background-color: var(--el-color-primary-light-9);
      color: var(--el-color-primary);
    }
  }
</style>
