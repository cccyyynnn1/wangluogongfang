<script lang="ts">
  export default {
    name: 'EditDetail',
  }
</script>
<script setup lang="ts">
  import VueDraggable from 'vuedraggable'
  import { TableColumnItemType } from '/#/store'
  import { RetrieveIndexType } from '@/types/index'
  import { useUserStore } from '@/store/modules/user'
  import { useCloned } from '@vueuse/core'
  import { PublishRuleParams } from '@/types/index'

  import { publishRuleApi } from '@/api-ecs/retrieve'
  const $baseMessage: any = inject('$baseMessage')
  const { getTableColumn } = useUserStore()
  const props = defineProps<{
    showDialog: boolean
    fields: TableColumnItemType[]
    queryForm: PublishRuleParams
  }>()

  const emit = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
    (e: 'handleok', showField: TableColumnItemType[]): void
  }>()

  const visible = ref(false)
  const showField = ref<TableColumnItemType[]>([])
  const choosableField = ref<TableColumnItemType[]>([])

  // 表单数据
  const formData = reactive<{
    submit: PublishRuleParams
    customTiem: string | number
    timeType: string
    duration: number
  }>({
    submit: {
      searchSql: '',
      indexType: 1,
      pageNum: 1,
      pageSize: 10,
      orderType: 'desc',
      orderField: 'requestTimeNs',
      startTime: '',
      endTime: '',
      displayFields: [],
      ruleName: '',
      searchTime: 10,
    },
    customTiem: 1,
    duration: 0,
    timeType: 'min',
  })

  const timeDuratioOptions = [
    {
      value: 1,
      label: '最近1分钟',
    },
    {
      value: 10,
      label: '最近10分钟',
    },
    {
      value: 60,
      label: '最近60分钟',
    },
    {
      value: 'user-defined',
      label: '自定义时间',
    },
  ]

  const timeOptions = [
    {
      value: 'min',
      label: '分钟',
    },
    {
      value: 'hours',
      label: '小时',
    },
  ]

  const handleOk = async () => {
    const query = {
      ...formData.submit,
      displayFields: showField.value.map((i) => i.fieldNameEn),
      startTime: '',
      endTime: '',
    }
    if (!query.searchSql) return $baseMessage('请输入检索语句', 'error', 'vab-hey-message-error')
    try {
      const { msg } = await publishRuleApi(query)
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      emit('handleok', showField.value)
      handleClose()
    } catch (error) {
      $baseMessage('发布失败', 'error', 'vab-hey-message-error')
    }
  }

  const handleClose = () => {
    emit('on-closeEvent', false)
  }
  onMounted(() => {
    visible.value = props.showDialog
  })
  watchEffect(() => {
    const { queryForm } = props
    const { cloned } = useCloned(queryForm)
    formData.submit = cloned.value

    const isHours = queryForm.searchTime % 60 > 0 ? false : true
    formData.customTiem = queryForm.searchTime > 60 ? 'user-defined' : queryForm.searchTime
    formData.timeType = isHours ? 'hours' : 'min'
    formData.duration = isHours ? queryForm.searchTime / 60 : queryForm.searchTime
  })
  watchEffect(() => {
    const { fields } = props
    const allColum = getTableColumn(1)

    const { cloned } = useCloned(fields)
    showField.value = cloned.value

    const showFiels = fields.map((item) => item.id)
    choosableField.value = allColum.filter((item) => !showFiels.includes(item.id))
  })
  watchEffect(() => {
    if (typeof formData.customTiem === 'number') {
      formData.submit.searchTime = formData.customTiem
    }
  })
  watchEffect(() => {
    const unit = formData.timeType === 'min' ? 1 : 60
    formData.submit.searchTime = formData.duration * unit
  })
</script>

<template>
  <div class="edit-detail">
    <el-dialog v-model="visible" :before-close="handleClose" title="编辑" width="700px">
      <el-form ref="formRef" label-position="right" label-width="90px">
        <el-form-item label="时间：" prop="interfaceName">
          <el-select v-model="formData.customTiem" class="m-2" style="width: 100%">
            <el-option v-for="item in timeDuratioOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="formData.customTiem === 'user-defined'">
          <el-input v-model.number="formData.duration" placeholder="请输入">
            <template #append>
              <el-select v-model="formData.timeType" class="m-2" style="width: 115px; background: #fff">
                <el-option v-for="item in timeOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="规则名称：" prop="interfaceName">
          <el-input v-model="formData.submit.ruleName" />
        </el-form-item>
        <el-form-item label="查询语句：" prop="interfaceName">
          <el-input v-model="formData.submit.searchSql" />
        </el-form-item>
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
                <el-tag type="info">{{ element.fieldNameCn }}</el-tag>
              </template>
            </vue-draggable>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="handleOk">确认</el-button>
          <el-button @click="handleClose">取消</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .space {
    width: 100%;
    padding: 10px;
    border: 1px solid var(--el-border-color);
    border-radius: 2px;
  }

  :deep(.el-tag) {
    margin-left: 10px;
    cursor: move;
  }
</style>
