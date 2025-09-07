<script lang="ts">
  export default {
    name: 'UpdateDialog',
  }
</script>
<script setup lang="ts">
  import UpdateConfig from '../update-config.vue'
  const props = defineProps<{
    isShowDialog: boolean
    mode: string
    description: string
    varVal: {
      updateSite: string
      interval: string
      timeUnit: number
      enable: boolean
    }
  }>()
  const emit = defineEmits<{
    (
      e: 'on-closeEvent',
      updateVal?: {
        updateSite: string
        interval: string
        timeUnit: number
        enable: boolean
      }
    ): void
  }>()
  const updateVal = reactive({
    updateSite: 'https://172.17.0.1',
    interval: 'day',
    timeUnit: 1,
    enable: false,
  })
  const modes = ref('') // 弹框类型

  const width = ref('35%') // 弹框类型

  const dialogVisible = ref(false) // 弹框显隐

  const title = ref('') // 弹框类型

  // 选项
  const roleOptions = [
    {
      value: 'hour',
      label: '小时',
    },
    {
      value: 'day',
      label: '天',
    },
    {
      value: 'week',
      label: '周',
    },
    {
      value: 'month',
      label: '月',
    },
    {
      value: 'year',
      label: '年',
    },
  ]

  // 表单数据
  const formData = reactive({
    source: '1',
    net: '',
    interval: '',
    data: 'day',
  })

  const handleClose = () => {
    emit('on-closeEvent')
  }

  const handelClick = () => {
    if (modes.value === 'updateConfig') return emit('on-closeEvent', updateVal)
    return emit('on-closeEvent')
  }

  const setTitle = () => {
    switch (modes.value) {
      case 'info':
        title.value = '版本说明'
        break
      case 'updateConfig':
        title.value = '升级配置'
        break
      case 'update':
        title.value = '升级'
        break
    }
  }

  watch(
    () => dialogVisible.value,
    () => {
      if (dialogVisible.value && modes.value === 'updateConfig') {
        const { enable, timeUnit, updateSite, interval } = props.varVal
        updateVal.enable = enable
        updateVal.timeUnit = timeUnit
        updateVal.updateSite = updateSite
        updateVal.interval = interval
      }
    }
  )

  onMounted(() => {
    modes.value = props.mode
    setTitle()
    modes.value == 'updateConfig' ? (width.value = '50%') : (width.value = '65%')
    dialogVisible.value = props.isShowDialog
  })
</script>

<template class="update-dialog">
  <el-dialog v-model="dialogVisible" :before-close="handleClose" :title="title" :width="width">
    <div v-if="modes === 'update'">
      <el-form ref="formRef" label-position="right" label-width="170px" :model="formData">
        <el-form-item label="升级包来源：" prop="source">
          <el-radio-group v-model="formData.source">
            <el-radio label="1">手动上传</el-radio>
            <el-radio disabled label="0">自动获取</el-radio>
          </el-radio-group>
        </el-form-item>
        <update-config v-if="formData.source === '1'" />
      </el-form>
    </div>
    <div v-else-if="modes === 'updateConfig'">
      <el-form ref="formRef" label-position="right" label-width="140px" :model="formData">
        <el-form-item label="升级网站">
          <el-input v-model="updateVal.updateSite" :style="{ width: '85%' }" />
        </el-form-item>
        <el-form-item label="升级时间间隔">
          <el-input v-model="updateVal.timeUnit" :style="{ width: '35%' }" />
          <el-select
            v-model="updateVal.interval"
            class="m-2"
            placeholder="选择类型"
            :style="{ width: '45%', marginLeft: '5%' }"
          >
            <el-option v-for="item in roleOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
      </el-form>
    </div>
    <div v-else-if="modes === 'info'">
      <div class="content-container">
        <div class="info-title">更新说明</div>
        <div class="particulars">
          <div class="sub-content">
            {{ props.description }}
          </div>
        </div>
      </div>
    </div>
    <template #footer>
      <span class="dialog-footer">
        <el-button type="primary" @click="handelClick">{{ modes === 'updateConfig' ? '保存' : '确定' }}</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<style scoped lang="scss">
  .info-title {
    margin-bottom: 20px;
  }

  .particulars {
    padding: 20px 20px;
    background: #fafafa;

    .subtitle {
      margin-bottom: 20px;
    }

    .sub-content {
      margin-left: 10px;
      line-height: 24px;
      white-space: break-spaces;
    }
  }
</style>
