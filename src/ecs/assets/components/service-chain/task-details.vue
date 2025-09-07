<script setup lang="ts">
  import { switchOpenChainSsortApi, switchCloseChainSsortApi } from '~/src/api-ecs/assets'

  const $baseMessage: any = inject('$baseMessage')

  const props = defineProps<{
    showPage: boolean
    currentItem: any
  }>()

  const visible = ref(true) // 弹框显隐

  const currentItem = ref()

  const emit = defineEmits<{
    (e: 'on-close-event'): void
    (e: 'on-reflash'): void
  }>()

  const currentData = reactive({
    site: '',
    verifyCode: '',
    status: false,
  })

  // 关闭回调
  const handleClose = () => {
    emit('on-close-event')
    emit('on-reflash')
  }

  const btnDisable = ref(false)
  const handleChange = async () => {
    try {
      btnDisable.value = true
      const { msg } = currentData.status
        ? await switchOpenChainSsortApi({ ids: [currentItem.value.id] })
        : await switchCloseChainSsortApi({ ids: [currentItem.value.id] })
      // const { msg } = await updateChainSsortApi({ id: currentItem.value.id, status: currentData.status ? '0' : '1' })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
    } finally {
      btnDisable.value = false
    }
  }

  const initData = async () => {
    // const { data } = await getChainSsortInfoApi({ id: currentItem.value.id })
    currentData.site = props.currentItem.site
    currentData.verifyCode = props.currentItem.verifyCode
    currentData.status = props.currentItem.status === 0 ? true : false
  }

  onMounted(() => {
    visible.value = props.showPage
    currentItem.value = props.currentItem
    initData()
  })
</script>

<script lang="ts">
  export default {
    name: 'TaskDetails',
  }
</script>
<template>
  <div class="task-details">
    <el-dialog v-model="visible" :before-close="handleClose" title="任务详情" width="800px">
      <div class="content">
        <div class="item">
          <span>站点域名：</span>
          <span>{{ currentData.site }}</span>
        </div>
        <div class="item">
          <span>验证码：</span>
          <span>{{ currentData.verifyCode }}</span>
        </div>
        <div class="item">
          <span>任务状态：</span>
          <el-switch
            v-model="currentData.status"
            :disabled="btnDisable"
            style="margin-right: 10px"
            @change="handleChange"
          />
          <span class="text">任务默认执行12小时后自动关闭，如果需要重新开启，只需要打开开关，每次开启时间为12小时</span>
        </div>
        <div class="item">
          <span>tips：</span>
          <span class="text">梳理任务的域名不可更改，如果要修改，请将本任务删除并创建新的任务</span>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .text {
    color: #0eadff;
  }
</style>
