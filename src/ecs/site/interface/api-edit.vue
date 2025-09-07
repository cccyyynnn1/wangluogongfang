<script lang="ts">
  export default {
    name: 'ApiEdit',
  }
</script>
<script setup lang="ts">
  import VueDraggable from 'vuedraggable'
  import { siteApiInterfaceItemType } from '@/types'
  import { useUserStore } from '@/store/modules/user'
  import { TableColumnItemType } from '~/types/store'
  import { updateApiInterfaceApi } from '~/src/api-ecs/site'
  import SaveModules from '@/ecs/config/init-configration/save-modules.vue'
  import { getHostApiFieldsApi } from '@/api-ecs/public'
  import { resetSiteDisPlaysFiledApi } from '~/src/api-ecs/custom-field'

  import { RefreshRight } from '@element-plus/icons-vue'
  const resetDisable = ref(false)
  const $baseMessage: any = inject('$baseMessage')
  const { getTableColumn } = useUserStore()
  const allField = getTableColumn(31)
  const props = defineProps<{
    showApiEdit: boolean
    apiData?: siteApiInterfaceItemType
    sessionId: number
  }>()

  const emits = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
    (e: 'on-reflash'): void
  }>()

  let defaultFields = [] as { id: number; fieldNameCn: string }[]
  let defaultFieldIds = [] as number[]
  const saveModulesVisible = ref(false)
  const visible = ref(false)
  const comparisonData = ref<{
    type?: number
    siteSessionId: number
    siteApiId: number
    tag: 'siteApi'
  }>()
  const replacementField = computed(() => fields.showField.map((i) => i?.fieldNameCn))
  onMounted(async () => {
    visible.value = props.showApiEdit
    const { data } = await getHostApiFieldsApi(31)
    defaultFields = (data.ids || []).map((id, index) => ({
      id: id,
      fieldNameCn: data.names[index],
    }))
    defaultFieldIds = data.ids
    if (!props.apiData) {
      displayFieldIds.value = data.ids
    }
  })

  // 表单数据
  const formData = reactive({
    apiName: '',
    apiUrl: '',
    searchSql: '',
    isDisplay: 0 as 0 | 1,
  })

  const fields = reactive({
    showField: [] as { id: number; fieldNameCn: string }[],
    choosableField: [] as { id: number; fieldNameCn: string }[],
  })

  const displayFieldIds = ref<number[]>([])

  const handleClose = () => {
    emits('on-closeEvent', false)
  }
  const handleSubmit = async () => {
    const id = props.apiData ? props.apiData.id : undefined
    const my_fields = fields.showField.map((i) => i.id)
    const { msg } = await updateApiInterfaceApi({
      ...formData,
      id,
      sessionId: props.sessionId,
      displayFields: my_fields.toString(),
    })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    emits('on-reflash')
    handleClose()
  }
  // 重置字段
  const handleReset = async () => {
    if (!props.apiData) return
    const { id, sessionId } = props.apiData
    const obj = {
      indexType: '',
      siteApiId: id,
      siteSessionId: '',
    }
    resetDisable.value = true
    try {
      const { msg, data } = await resetSiteDisPlaysFiledApi(obj)
      displayFieldIds.value = data || []
      // $baseMessage(msg, 'success', 'vab-hey-message-success')
      // emits('on-reflash')
    } finally {
      resetDisable.value = false
    }
  }

  const handle2save = () => {
    if (!props.apiData) return
    const { id, sessionId } = props.apiData
    comparisonData.value = {
      type: undefined,
      siteSessionId: sessionId,
      siteApiId: id,
      tag: 'siteApi',
    }
    saveModulesVisible.value = true
  }
  watchEffect(() => {
    if (props.apiData) {
      const { apiName, apiUrl, displayFieldsArr, searchSql, isDisplay } = props.apiData
      formData.apiName = apiName
      formData.apiUrl = apiUrl
      formData.searchSql = searchSql
      formData.isDisplay = isDisplay
      displayFieldIds.value = displayFieldsArr
    }
  })
  watchEffect(() => {
    fields.showField = displayFieldIds.value.map((i) => {
      const field = allField.find((item) => item.id === i)
      return field
    }) as TableColumnItemType[]
    fields.choosableField = allField.filter((i) => !displayFieldIds.value.includes(i.id))
  })
</script>

<template>
  <div class="api-edit">
    <el-dialog
      v-model="visible"
      :before-close="handleClose"
      :title="`${props.apiData ? '编辑' : '新增'}接口`"
      width="1200px"
    >
      <el-form ref="formRef" label-position="right" label-width="120px" :model="formData">
        <el-form-item label="接口名称：" prop="apiName">
          <el-input v-model="formData.apiName" />
        </el-form-item>
        <el-form-item label="接口地址：" prop="apiUrl">
          <el-input v-model="formData.apiUrl" />
        </el-form-item>
        <el-form-item label="接口是否展示：" prop="isDisplay">
          <el-switch
            v-model="formData.isDisplay"
            active-text="是"
            :active-value="1"
            inactive-text="否"
            :inactive-value="0"
            inline-prompt
            style="--el-switch-off-color: #ccc"
          />
        </el-form-item>
        <el-form-item label="过滤条件：" prop="searchSql">
          <el-input v-model="formData.searchSql" />
        </el-form-item>
        <el-form-item prop="roleName">
          <template #label>
            <div style="text-align: right">
              <span>展示字段：</span>
              <br />
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
              :list="fields.showField"
            >
              <template #item="{ element }">
                <el-tag>{{ element.fieldNameCn }}</el-tag>
              </template>
            </vue-draggable>
          </div>
        </el-form-item>
        <el-form-item label="可选字段：">
          <div class="space">
            <vue-draggable
              animation="300"
              chosen-class="chosenClass"
              ghost-class="ghost"
              group="my-group"
              item-key="id"
              :list="fields.choosableField"
              :sort="false"
            >
              <template #item="{ element }">
                <el-tag type="info">{{ element.fieldNameCn }}</el-tag>
              </template>
            </vue-draggable>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button
          v-if="props.apiData"
          :disabled="resetDisable"
          :icon="RefreshRight"
          link
          :loading="resetDisable"
          style="float: left; color: #6954f0; margin-left: 0px"
          @click="handleReset"
        >
          重置字段
        </el-button>
        <span class="dialog-footer">
          <span v-if="props.apiData" class="save" @click="handle2save">另存为</span>
          <el-button type="primary" @click="handleSubmit">确认</el-button>
          <el-button @click="handleClose">取消</el-button>
        </span>
      </template>
      <save-modules
        v-model="saveModulesVisible"
        :comparison-data="comparisonData"
        :replacement-field="replacementField"
      />
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .space {
    width: 100%;
    height: 150px;
    padding: 10px;
    border: 1px solid var(--el-border-color);
    border-radius: 2px;
    overflow-y: auto;
    div {
      height: calc(100% - 20px);
      min-height: 130px;
    }
  }
  :deep() {
    .el-tag {
      margin-left: 10px;
    }
    .save {
      cursor: pointer;
      margin-right: 10px;
      color: var(--el-color-primary);
      font-weight: 400;
      font-size: 14px;
      vertical-align: middle;
    }
  }
</style>
