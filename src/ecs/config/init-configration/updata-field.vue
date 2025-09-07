<script setup lang="ts">
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import VueDraggable from 'vuedraggable'
  import { RetrieveIndexType } from '@/types/index'
  import { useUserStore } from '@/store/modules/user'
  import { fieldsTemplateRestApi } from '~/src/api-ecs/custom-field'

  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')

  interface ApplicationConfigurationProps {
    modelValue: boolean
    fields: string[]
    retrieveIndexType: RetrieveIndexType
    title?: string
    fieldEditData?: {
      type: number
      siteSessionId?: number
      siteApiId?: number
      tag: 'alarm' | 'survey' | 'siteSession' | 'siteApi'
    }
  }
  let indexType: undefined | number = undefined

  const { getTableColumn } = useUserStore()
  const props = withDefaults(defineProps<ApplicationConfigurationProps>(), {
    modelValue: false,
    retrieveIndexType: 1,
    fields: () => [],
    title: '字段配置',
  })

  const choosableField = ref<string[]>([])
  const showField = ref<string[]>([])
  const emits = defineEmits<{
    (e: 'update:modelValue', visible: boolean): void
    (e: 'handleok', fields: string[]): void
  }>()

  const dialogVisible = useVModel(props, 'modelValue', emits)

  const hanldeFieldsRest = async () => {
    if (props?.fieldEditData) {
      $baseConfirm('确定要重置为系统初始化字段?', null, async () => {
        const { type, siteSessionId = '', siteApiId = '', tag } = props.fieldEditData!
        const { data } = await fieldsTemplateRestApi({ type, siteSessionId, siteApiId, tag })
        $baseMessage('字段重置成功！', 'success', 'vab-hey-message-success')
        showField.value = data.namesList
        choosableField.value = getTableColumn(indexType!)
          .map((i) => i.fieldNameCn)
          .filter((i) => !data.namesList.includes(i) && Boolean(i))
        emits('handleok', data.namesList)
      })
    }
  }

  function handleok() {
    emits('handleok', showField.value)
  }

  function handleCancel() {
    dialogVisible.value = false
  }
  watch(
    () => props.modelValue,
    () => {
      if (props.retrieveIndexType && props.modelValue) {
        showField.value = [...props.fields]
        if (indexType !== props.retrieveIndexType) {
          indexType = props.retrieveIndexType
          choosableField.value = getTableColumn(indexType)
            .map((i) => i.fieldNameCn)
            .filter((i) => !showField.value.includes(i) && Boolean(i))
        }
      }
    }
  )
</script>

<script lang="ts">
  export default {
    name: 'UpdataField',
  }
</script>

<template>
  <vab-dialog v-model="dialogVisible" destroy-on-close :title="`${title}字段配置`" width="1000" v-bind="$attrs">
    <el-form label-position="right" label-width="90px">
      <el-form-item prop="roleName">
        <template #label>
          <div style="text-align: right">
            <span>展示字段：</span>
            <span style="color: #a9acb3; font-size: 13px">(拖动排序)&nbsp;&nbsp;&nbsp;</span>
          </div>
        </template>
        <vue-draggable
          animation="300"
          chosen-class="chosenClass"
          ghost-class="ghost"
          group="my-group"
          item-key="id"
          :list="showField"
        >
          <template #item="{ element }">
            <el-tag>{{ element }}</el-tag>
          </template>
        </vue-draggable>
      </el-form-item>
      <el-form-item label="可选字段：" prop="roleName">
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
            <el-tag type="info">{{ element }}</el-tag>
          </template>
        </vue-draggable>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button :auto-insert-space="false" style="float: left; margin-left: 90px" @click="hanldeFieldsRest">
        重置为系统默认
      </el-button>
      <el-button type="primary" @click="handleok">确定</el-button>
      <el-button @click="handleCancel">取消</el-button>
    </template>
  </vab-dialog>
</template>

<style scoped lang="scss">
  .el-form-item__content {
    & > div {
      width: 100%;
      padding: 10px;
      border: 1px solid var(--el-border-color);
      border-radius: 2px;
      min-height: 100px;
    }
  }
  :deep() {
    .el-tag {
      margin-left: 10px;
      cursor: move;
    }
    .el-dialog {
      border-radius: 10px;
      .el-button {
        --el-border-radius-base: 6px;
      }
    }
  }
</style>
