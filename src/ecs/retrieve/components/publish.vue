<script lang="ts">
  export default {
    name: 'Publish',
  }
</script>
<script setup lang="ts">
  const props = defineProps<{
    showPublish: boolean
  }>()

  const emits = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
    (e: 'on-submit', val: { searchTime: number; ruleName: string }): void
  }>()

  const visible = ref(false)

  onMounted(() => {
    visible.value = props.showPublish
  })

  const submitForm = reactive({
    searchTime: 0,
    ruleName: '',
  })

  // 表单数据
  const formData = reactive({
    monitoringName: undefined,
    timeType: 'min',
    duration: 0,
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

  const handleClose = () => {
    emits('on-closeEvent', false)
  }
  const handleOk = () => {
    emits('on-submit', submitForm)
  }

  watchEffect(() => {
    if (typeof formData.monitoringName === 'number') {
      submitForm.searchTime = formData.monitoringName
    }
  })
  watchEffect(() => {
    const unit = formData.timeType === 'min' ? 1 : 60
    submitForm.searchTime = formData.duration * unit
  })
</script>

<template>
  <div class="filed-config">
    <el-dialog
      v-model="visible"
      :before-close="handleClose"
      title="流量监控"
      width="600px"
    >
      <el-form label-position="right" label-width="90px">
        <el-form-item label="监控时长：" prop="monitoringName">
          <el-select
            v-model="formData.monitoringName"
            class="m-2"
            style="width: 100%"
          >
            <el-option
              v-for="item in timeDuratioOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="formData.monitoringName === 'user-defined'"
          prop="ruleName"
        >
          <el-input v-model.number="formData.duration" placeholder="请输入">
            <template #append>
              <el-select
                v-model="formData.timeType"
                class="m-2"
                style="width: 115px; background: #fff"
              >
                <el-option
                  v-for="item in timeOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="规则名称：" prop="ruleName">
          <el-input
            v-model="submitForm.ruleName"
            placeholder="请输入规则名称"
          />
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

<style scoped lang="scss"></style>
