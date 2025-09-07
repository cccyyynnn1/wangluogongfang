<script lang="ts">
  export default {
    name: 'AssetsLabel',
  }
</script>

<script setup lang="ts">
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import { Search, Close } from '@element-plus/icons-vue'
  import { getAssetsLabelBySearchApi } from '@/api-ecs/assets-preview'
  import { CMDBfingerprintList, FingerprintItem } from '@/types/index'

  const props = defineProps<{
    selectedFieelds: FingerprintItem[]
  }>()

  let labelMap = new Map()
  const showInfo = ref(false)
  // 标签查询
  const searchStr = ref('')
  // 标签查询
  const orderType = ref<0 | 1>(0)

  // 标签分株类型
  const labelTypes = ref<Omit<CMDBfingerprintList, 'labelList' | 'labelSize'>[]>([
    {
      groupId: 1,
      groupName: '流量指纹',
    },
    {
      groupId: 2,
      groupName: '规则指纹',
    },
    {
      groupId: 3,
      groupName: '三方指纹',
    },
  ])
  // 当前分组标签数量
  const labelTypeNum = ref(0)
  // 当前标签类型
  const typeKey = ref(1)
  // 选中的标签Map
  const selectLabelsMap = ref<Map<number, FingerprintItem>>(new Map())
  // 当前归属下的标签
  const labelInfos = ref<FingerprintItem[]>([])
  // 选择的标签
  const selectLabels = computed<FingerprintItem[]>(() => [...selectLabelsMap.value.values()])

  const labelRelat = ref<'and' | 'or'>('or')

  const handLelabelItemClick = (labelitem: FingerprintItem) => {
    selectLabelsMap.value.has(labelitem.id)
      ? selectLabelsMap.value.delete(labelitem.id)
      : selectLabelsMap.value.set(labelitem.id, labelitem)
  }

  const handleGetAssetsLabelBySearch = async () => {
    const { data } = await getAssetsLabelBySearchApi({
      searchStr: searchStr.value,
      orderType: orderType.value,
      groupId: typeKey.value,
    })
    labelInfos.value = data || []
    labelTypeNum.value = data?.length ?? 0
  }
  const selectedToApply = () => {
    emits('confirm', { labels: selectLabels.value, labelRelat: labelRelat.value })
    showInfo.value = false
  }

  const resetSelect = () => {
    selectLabelsMap.value.clear()
  }
  const emits = defineEmits<{
    (e: 'confirm', val: { labels: FingerprintItem[]; labelRelat: 'and' | 'or' }): void
  }>()

  watch(
    () => [typeKey.value, orderType.value],
    () => {
      handleGetAssetsLabelBySearch()
    },
    {
      immediate: true,
    }
  )

  defineExpose({
    showDialog: (labelList: CMDBfingerprintList[], _labelRelat: 'and' | 'or') => {
      labelList.forEach(({ groupId, groupName, labelList }) => {
        labelMap.set(groupId, labelList)
      })
      selectLabelsMap.value.clear()
      props.selectedFieelds.forEach((field: FingerprintItem) => {
        selectLabelsMap.value.set(field.id, field)
      })
      typeKey.value = 1
      labelRelat.value = _labelRelat
      showInfo.value = true
    },
  })
</script>

<template>
  <vab-dialog v-model="showInfo" destroy-on-close title="添加" width="1200">
    <div class="assets-label-content">
      <el-row>
        <el-col :span="12">
          <el-input
            v-model="searchStr"
            clearable
            placeholder="请输入"
            style="width: 300px"
            @clear="handleGetAssetsLabelBySearch"
          />
          <el-button
            :icon="Search"
            style="margin-left: -2px; letter-spacing: 0"
            type="primary"
            @click="handleGetAssetsLabelBySearch"
          >
            检索
          </el-button>
        </el-col>
        <el-col :span="12" style="text-align: right">
          <span style="vertical-align: sub">排序方式：</span>
          <el-button class="sortBtn" :class="{ selected: orderType === 0 }" @click="orderType = 0">点击热度</el-button>
          <el-button
            class="sortBtn"
            :class="{ selected: orderType === 1 }"
            style="margin-left: 10px"
            @click="orderType = 1"
          >
            引用次数
          </el-button>
        </el-col>
      </el-row>
      <div class="lable-list">
        <div class="lable-preview">
          <div class="lable-types">
            <div
              v-for="labelType in labelTypes"
              :key="labelType.groupId"
              class="leabl-type"
              :class="{ checked: typeKey === labelType.groupId }"
              @click="
                () => {
                  typeKey = labelType.groupId
                  labelTypeNum = 0
                }
              "
            >
              {{ labelType.groupName }}
              <template v-if="typeKey === labelType.groupId && labelTypeNum > 0">({{ labelTypeNum }})</template>
            </div>
          </div>
          <div class="lable-info">
            <div
              v-for="labelItem in labelInfos"
              :key="labelItem.id"
              class="lable-item"
              :class="{ checked: selectLabelsMap.has(labelItem.id) }"
              style="height: 32px; line-height: 18px"
              @click="() => handLelabelItemClick(labelItem)"
            >
              {{ labelItem.labelName }}
              <template v-if="labelItem.isNew">
                <span style="color: #ff6060; font-size: 10px; font-weight: bold">NEW</span>
              </template>
            </div>
          </div>
        </div>
      </div>
      <div class="lable-selected">
        <h5>已添加标签：</h5>
        <div class="lable-selected-list">
          <div
            v-for="label in selectLabels"
            :key="label.id"
            class="lable-select-item"
            style="height: 32px; line-height: 18px"
          >
            {{ label.labelName }}
            <div class="selectTag-closeable" @click="selectLabelsMap.delete(label.id)">
              <Close />
            </div>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <ul class="labelRelat" style="float: left; margin-top: 4px">
        <li :class="{ active: labelRelat == 'or' }" @click="labelRelat = 'or'">或</li>
        <li :class="{ active: labelRelat == 'and' }" @click="labelRelat = 'and'">且</li>
        <li class="tips">(ps:切换当前标签可筛选“或”和“且”的逻辑关系)</li>
      </ul>
      <el-button type="primary" @click="selectedToApply">应用</el-button>
      <el-button @click="resetSelect">重置</el-button>
    </template>
  </vab-dialog>
</template>

<style scoped lang="scss">
  .assets-label-content {
    font-family: PingFangSC, PingFang SC;
    font-size: 14px;
    .lable-list {
      height: 334px;
      margin-top: 20px;
      border: 1px solid #eeeef0;
      border-radius: 2px;
      .lable-preview {
        display: flex;
        height: 100%;
        .lable-types {
          width: 181px;
          height: 100%;
          padding: 20px;
          overflow-y: auto;
          border-right: 1px solid #eeeef0;
          .leabl-type {
            height: 40px;
            font-weight: 600;
            line-height: 40px;
            color: #666666;
            text-indent: 20px;
            cursor: pointer;
            border-radius: 4px;
            &.checked {
              color: var(--el-color-primary);
              background: var(--el-color-primary-light-9);
            }
          }
        }
        .lable-info {
          padding: 20px;
        }
      }
      .lable-search {
        height: 100%;
        padding: 20px;
        overflow-y: auto;
        .search-list {
          display: flex;
          &:not(:last-child) {
            margin-bottom: 10px;
          }
          .search-label-type {
            padding: 5px 10px;
            color: #888b91;
          }
        }
      }
      .lable-info {
        flex: 1;
        overflow-y: auto;
        .lable-item {
          display: inline-block;
          width: fit-content;
          padding: 7px 12px;
          margin: 0 10px 10px 0;
          font-weight: 400;
          color: #606266;
          cursor: pointer;
          background: #f3f3fe;
          border-radius: 4px;
          &.checked {
            color: #fff;
            background-color: var(--el-color-primary);
          }
        }
      }
    }

    .lable-selected {
      h5 {
        margin-block: 10px !important;
        font-weight: 400;
        color: #a9acb3;
      }
      .lable-selected-list {
        height: 180px;
        padding: 20px;
        overflow: auto;
        border: 1px solid #eeeef0;
        .lable-select-item {
          position: relative;
          display: inline-block;
          padding: 7px 20px;
          margin: 0 10px 10px 0;
          font-weight: 400;
          color: var(--el-color-primary);
          cursor: pointer;
          background: var(--el-color-primary-light-9);
          border-radius: 4px;

          // --el-color-primary-light-7

          .selectTag-closeable {
            position: absolute;
            top: -5px;
            right: -5px;
            width: 12px;
            height: 12px;
            line-height: 12px;
            color: #fff;
            text-align: center;
            background-color: #ff6060;
            border-radius: 100%;
            svg {
              width: 10px;
              height: 10px;
            }
          }
        }
      }
    }
    .sortBtn {
      border: none;

      &.selected {
        color: var(--el-color-primary);
        background-color: var(--el-color-primary-light-9);
      }
    }
  }
  :deep() {
    .el-dialog__footer:empty {
      padding: 0;
    }
    .el-button__text--expand {
      margin-right: 0 !important;
      letter-spacing: normal;
    }
  }
</style>
