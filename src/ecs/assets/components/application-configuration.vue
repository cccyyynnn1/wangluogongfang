<script setup lang="ts">
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import VueDraggable from 'vuedraggable'
  import { TableColumnItemType } from '/#/store'
  import { RetrieveIndexType } from '@/types/index'
  import { useUserStore } from '@/store/modules/user'
  import { useCloned } from '@vueuse/core'
  import { RefreshRight } from '@element-plus/icons-vue'
  import SaveModules from '@/ecs/config/init-configration/save-modules.vue'
  import { resetSiteDisPlaysFiledApi } from '~/src/api-ecs/custom-field'

  const { getTableColumn } = useUserStore()
  interface ApplicationConfigurationProps {
    modelValue: boolean
    fields: TableColumnItemType[]
    retrieveIndexType: RetrieveIndexType
    showRestIcon?: boolean
  }
  const props = withDefaults(defineProps<ApplicationConfigurationProps>(), {
    modelValue: false,
    retrieveIndexType: 1,
    showRestIcon: false,
  })

  const saveModulesVisible = ref(false)
  const choosableField = ref<TableColumnItemType[]>([])
  const showField = ref<TableColumnItemType[]>([])
  const emits = defineEmits<{
    (e: 'update:modelValue', visible: boolean): void
    (e: 'handleok', showField: TableColumnItemType[]): void
    (e: 'handleCancel'): void
  }>()
  const resetDisable = ref(false)
  const dialogVisible = useVModel(props, 'modelValue', emits)
  const replacementField = computed(() => showField.value.map((i) => i.fieldNameCn))

  const comparisonData = ref<{
    type: number
    siteSessionId?: number
    siteApiId?: number
    tag: 'alarm' | 'survey' | 'siteSession' | 'siteApi'
  }>()

  function handleok() {
    dialogVisible.value = false
    emits('handleok', showField.value)
  }

  function handleCancel() {
    dialogVisible.value = false
    emits('handleCancel')
  }

  // 重置字段
  const handleReset = async () => {
    resetDisable.value = true
    const { retrieveIndexType } = props
    const obj = {
      indexType: retrieveIndexType.toString(),
      siteApiId: '',
      siteSessionId: '',
    }
    try {
      const { data, msg } = await resetSiteDisPlaysFiledApi(obj)
      const allColum = getTableColumn(retrieveIndexType)
      showField.value = data.map((_id: any) => allColum.find((item) => item.id === _id))
      const showFiels = showField.value.map((item) => item?.id)
      choosableField.value = allColum.filter((item) => !showFiels.includes(item?.id))
      // emits('handleok', showField.value)
    } finally {
      resetDisable.value = false
    }
  }
  const handle2save = () => {
    const { retrieveIndexType } = props
    // 只有在调查告警下使用，站点是单独的配置页面
    const isAlarm = [10, 29, 34].includes(retrieveIndexType)
    comparisonData.value = {
      type: retrieveIndexType,
      siteSessionId: undefined,
      siteApiId: undefined,
      // indexFiledsTemId: undefined,
      tag: isAlarm ? 'alarm' : 'survey',
    }
    saveModulesVisible.value = true
  }
  watch(
    () => dialogVisible.value,
    () => {
      if (dialogVisible.value) {
        const { retrieveIndexType, fields } = props
        const allColum = getTableColumn(retrieveIndexType)
        const { cloned } = useCloned(fields)
        showField.value = cloned.value

        const showFiels = fields.map((item) => item?.id)
        choosableField.value = allColum?.filter((item) => !showFiels.includes(item?.id))
      }
    },
    {
      immediate: true,
    }
  )
</script>

<script lang="ts">
  export default {
    name: 'ApplicationConfiguration',
  }
</script>

<template>
  <vab-dialog v-model="dialogVisible" destroy-on-close title="配置页" width="800px">
    <el-form label-position="right" label-width="90px">
      <el-form-item prop="roleName">
        <template #label>
          <div style="text-align: right">
            <span>展示字段：</span>
            <span style="color: #a9acb3; font-size: 13px">(拖动排序)&nbsp;&nbsp;&nbsp;</span>
          </div>
        </template>
        <div class="space">
          <vue-draggable
            animation="300"
            chosen-class="chosenClass"
            ghost-class="ghost"
            group="my-group"
            item-key="id"
            :list="showField"
          >
            <template #item="{ element }">
              <el-tag>{{ element.fieldNameCn }}</el-tag>
            </template>
          </vue-draggable>
        </div>
      </el-form-item>
      <el-form-item label="可选字段：" prop="roleName">
        <div class="space">
          <vue-draggable
            animation="300"
            chosen-class="chosenClass"
            ghost-class="ghost"
            group="my-group"
            item-key="id"
            :list="choosableField"
            :sort="false"
          >
            <template #item="{ element }">
              <el-tag v-if="element.fieldNameCn" type="info">{{ element.fieldNameCn }}</el-tag>
            </template>
          </vue-draggable>
        </div>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button
        v-if="retrieveIndexType !== 33"
        :disabled="resetDisable"
        :icon="RefreshRight"
        link
        :loading="resetDisable"
        style="float: left; color: #6954f0; margin-left: 85px"
        @click="handleReset"
      >
        重置字段
      </el-button>
      <span v-if="retrieveIndexType !== 33" class="save" :disabled="resetDisable" @click="handle2save">另存为</span>
      <el-button :disabled="resetDisable" type="primary" @click="handleok">保存</el-button>
      <el-button @click="handleCancel">取消</el-button>
    </template>
    <save-modules
      v-model="saveModulesVisible"
      :comparison-data="comparisonData"
      :replacement-field="replacementField"
    />
  </vab-dialog>
</template>

<style scoped lang="scss">
  .space {
    width: 100%;
    padding: 10px;
    border: 1px solid var(--el-border-color);
    border-radius: 2px;
    min-height: 100px;
    div {
      height: calc(100% - 20px);
      min-height: 80px;
    }
  }
  .save {
    cursor: pointer;
    margin-right: 10px;
    color: var(--el-color-primary);
    font-weight: 400;
    font-size: 14px;
    vertical-align: middle;
  }
  :deep(.el-tag) {
    margin-left: 10px;
    cursor: move;
  }
</style>
