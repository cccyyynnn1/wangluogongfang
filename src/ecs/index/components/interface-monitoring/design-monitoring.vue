<script setup lang="ts">
  import UiTable from '~/library/components/ui/ui-table.vue'

  import { MonitoringItem, MonitoringType } from '@/types/index'

  import { uuid } from '~/src/utils'

  import type { FormInstance } from 'element-plus'

  import { requireRules } from '~/src/utils/rules'

  import VueDraggable from 'vuedraggable'

  import { saveOrUpdateApiMonitorApi } from '~/src/api-ecs/dashboard'

  import dayjs from 'dayjs'

  const formSaveRef = ref<FormInstance>()

  const draggableContent = ref()

  const disabled = ref(false)

  const props = defineProps<{
    allData: MonitoringType
    isBtn: boolean
  }>()

  const emit = defineEmits<{
    (e: 'on-back'): void
    (e: 'on-reflash', val: MonitoringType): void
    (e: 'change-name', val: string): void
    (e: 'toDetail', val: MonitoringType): void
  }>()

  const formSave = ref({
    id: undefined,
    uuid: uuid(),
    name: '',
    timeQuantum: '1hours',
    switchValue: false,
    timeDuration: '1.5min',
    timeDate: undefined,
  })

  const toolList = [{ id: 0, label: '添加窗口' }]

  const origin = {
    pageX: 0,
    pageY: 0,
  }

  const target = {
    pageX: 0,
    pageY: 0,
  }

  const styleObj = ref({
    left: 0,
    top: '0',
    height: '330px',
    width: '100%',
    x: 0,
    y: 0,
    long: 400,
  })

  const formSaveRules = {
    name: requireRules,
  }

  const UIRef = ref()

  const visible = ref(false)

  const timeQuantumOptions = [
    {
      value: '1hours',
      label: '1小时',
    },
    {
      value: '3hours',
      label: '3小时',
    },
    {
      value: '6hours',
      label: '6小时',
    },
    {
      value: '24hours',
      label: '24小时',
    },
    {
      value: 'today',
      label: '今日',
    },
    {
      value: 'yesterday',
      label: '昨天',
    },
    {
      value: 'last-three-days',
      label: '最近三天',
    },
    {
      value: 'last-week',
      label: '最近一周',
    },
    {
      value: 'last-two-week',
      label: '最近两周',
    },
    {
      value: 'user-defined',
      label: '自定义时间',
    },
  ]

  const timeSetting = ref({
    endTime: '',
    startTime: '',
  })

  const timeDuratioOptions = [
    // {
    //   value: '30s',
    //   label: '30秒',
    // },
    {
      value: '1.5min',
      label: '90秒',
    },
    {
      value: '5min',
      label: '5分钟',
    },
    {
      value: '10min',
      label: '10分钟',
    },
  ]

  const uiObj = {
    UiTable: markRaw(UiTable),
  }

  const allData = ref<MonitoringType>()

  const containerData = ref<MonitoringItem[]>([])

  const formatTimeDate = () => {
    if (formSave.value.timeDate) {
      // @ts-ignore
      const [startDate, endDate] = formSave.value.timeDate
      timeSetting.value.endTime = dayjs(endDate).format('YYYY-MM-DD HH:mm:ss')
      timeSetting.value.startTime = dayjs(startDate).format('YYYY-MM-DD HH:mm:ss')
    }
  }

  // 得到自定义时间
  watchEffect(() => {
    formatTimeDate()
  })

  const formatDayDate = () => {
    if (formSave.value.timeQuantum === 'user-defined') {
      formatTimeDate()
      return
    }
    const timeDate = dayjs()
    const endDate = timeDate.format('YYYY-MM-DD HH:mm:ss')
    const cur = timeDate.startOf('day')
    let startDate = null
    switch (formSave.value.timeQuantum) {
      case '1hours':
        startDate = timeDate.subtract(1, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case '3hours':
        startDate = timeDate.subtract(3, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case '6hours':
        startDate = timeDate.subtract(6, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case '24hours':
        startDate = timeDate.subtract(24, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'today':
        startDate = timeDate.startOf('day').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'yesterday':
        startDate = cur.subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'last-three-days':
        startDate = timeDate.subtract(3, 'day').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'last-week':
        startDate = timeDate.subtract(1, 'week').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'last-two-week':
        startDate = timeDate.subtract(2, 'week').format('YYYY-MM-DD HH:mm:ss')
        break
      default:
        startDate = timeDate.startOf('day').format('YYYY-MM-DD HH:mm:ss')
        break
    }
    timeSetting.value.endTime =
      formSave.value.timeQuantum === 'yesterday'
        ? timeDate.endOf('day').subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss')
        : endDate
    timeSetting.value.startTime = startDate
  }

  // 得到时间
  watchEffect(() => {
    formatDayDate()
  })

  // 添加
  const addComponent = () => {
    calcAfterAddItem(200)
    containerData.value.push({
      name: markRaw(UiTable),
      uuid: uuid(),
      key: 'UiTable',
      style: JSON.parse(JSON.stringify(styleObj.value)),
      data: {
        name: '名称',
      },
    })
    styleObj.value.top = '0'
    styleObj.value.y = 0
  }

  // 添加新的重新计算top和height
  const calcAfterAddItem = (num: number) => {
    let temp = 0
    containerData.value.forEach((item: MonitoringItem, index: number) => {
      const _index_1 =
        containerData.value[index].style.top!.indexOf('p') > 0
          ? containerData.value[index].style.top!.indexOf('p')
          : undefined
      const _index_2 =
        containerData.value[index].style.height!.indexOf('p') > 0
          ? containerData.value[index].style.height!.indexOf('p')
          : undefined
      const height =
        Number(containerData.value[index].style.top!.slice(0, _index_1)) +
        Number(containerData.value[index].style.height!.slice(0, _index_2))
      temp = height > temp ? height : temp
    })
    styleObj.value.top = `${temp}px`
    if (num != 0) {
      styleObj.value.y = temp
    }
    containerData.value.forEach((item: MonitoringItem) => {
      item.height = temp + num
      item.style.long = temp + num
    })
  }

  // 删除
  const handleDelete = (currentItem: MonitoringItem) => {
    const _index = containerData.value.findIndex((item: MonitoringItem) => {
      return item.uuid == currentItem.uuid
    })
    containerData.value.splice(_index, 1)
  }

  // 更新数据
  const handleUpdate = (value: MonitoringItem) => {
    let temp = 0
    containerData.value.forEach((item: MonitoringItem, index: number) => {
      if (containerData.value[index].uuid == value.uuid) {
        Object.keys(item).forEach((td: string) => {
          if (td !== 'name') {
            // @ts-ignore
            item[td] = value[td]
          }
        })
      }
      temp = containerData.value[index].style.long! > temp ? containerData.value[index].style.long! : temp
    })
    containerData.value.forEach((item: MonitoringItem) => {
      item.height = temp
      item.style.long = temp
    })
  }

  // 保存
  const btnDisabled = ref(false)
  const handleSave = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl!.validate(async (valid) => {
      if (valid) {
        try {
          btnDisabled.value = true
          const val = { ...toRaw(formSave.value), list: toRaw(containerData.value) }
          const payload = { ...formSave.value, list: JSON.stringify(containerData.value) }
          const { msg } = await saveOrUpdateApiMonitorApi({ ...val, list: JSON.stringify(val.list) })
          ElMessage({ message: msg, type: 'success' })
          visible.value = false
          emit('change-name', formSave.value.name)
          if (props.allData.uuid) {
            emit('toDetail', payload)
          } else {
            emit('on-back')
          }
        } finally {
          btnDisabled.value = false
        }
      } else {
        return false
      }
    })
  }

  // 打开保存弹窗
  const saveEvent = () => {
    if (containerData.value.length < 1)
      return ElMessageBox.alert('窗口不存在，保存失败', '提示', {
        confirmButtonText: '确定',
      })
    const flag = containerData.value.every((item: any) => {
      return item.data.indexType
    })
    if (!flag)
      return ElMessageBox.alert('存在未设定条件的窗口，保存失败', '提示', {
        confirmButtonText: '确定',
      })
    visible.value = true
    nextTick(() => {
      formSaveRef.value!.clearValidate()
      formSaveRef.value!.resetFields()
    })
  }

  const viewStart = (e: any) => {
    const { originalEvent } = e
    draggableContent.value.componentStructure.children.forEach((item: any) => {
      item.component.exposeProxy.unObserverUITable()
    })
    origin.pageX = 0
    origin.pageY = 0
    origin.pageX = originalEvent.pageX
    origin.pageY = originalEvent.pageY
  }

  const toolStart = (e: any) => {
    const { originalEvent } = e
    origin.pageX = 0
    origin.pageY = 0
    origin.pageX = originalEvent.pageX
    origin.pageY = originalEvent.pageY
  }

  const viewEnd = (e: any) => {
    const { to, originalEvent, newIndex } = e
    target.pageX = originalEvent.pageX
    target.pageY = originalEvent.pageY
    let x = target.pageX - origin.pageX + containerData.value[newIndex].style.x!
    let y = target.pageY - origin.pageY + containerData.value[newIndex].style.y!
    x = x <= 0 ? 0 : x
    y = y <= 0 ? 0 : y
    containerData.value[newIndex].style.x = x
    containerData.value[newIndex].style.y = y
    const width = to.clientWidth
    const left = `${(x / width) * 100}%`
    const top = `${y}px`
    containerData.value[newIndex].style.left = left
    containerData.value[newIndex].style.top = top
    calcAfterAddItem(0)
    styleObj.value.top = '0'
  }

  const toolEnd = (e: any) => {
    const { item, to, originalEvent, newIndex } = e
    const originRectX = item.getBoundingClientRect().x
    const originRectY = item.getBoundingClientRect().y
    const targetRectX = to.getBoundingClientRect().x
    const targetRectY = to.getBoundingClientRect().y
    target.pageX = originalEvent.pageX
    target.pageY = originalEvent.pageY
    const width = to.clientWidth
    let x = target.pageX - origin.pageX + (originRectX - targetRectX)
    let y = target.pageY - origin.pageY + (originRectY - targetRectY)
    x = x <= 0 ? 0 : x
    y = y <= 0 ? 0 : y
    const left = `${(x / width) * 100}%`
    const top = `${y}px`
    containerData.value[newIndex] = {
      name: markRaw(UiTable),
      uuid: uuid(),
      key: 'UiTable',
      style: JSON.parse(JSON.stringify(styleObj.value)),
      data: {
        name: '名称',
      },
    }
    containerData.value[newIndex].style.x = x
    containerData.value[newIndex].style.y = y
    containerData.value[newIndex].style.left = left
    containerData.value[newIndex].style.top = top
    calcAfterAddItem(0)
    styleObj.value.top = '0'
  }

  const goBack = () => {
    emit('on-back')
  }

  const handleClick = (e: any) => {
    disabled.value = e.srcElement.className == 'el-scrollbar__thumb' ? true : false
  }

  onMounted(() => {
    // @ts-ignore
    if (props.allData.uuid) {
      initData()
    }
  })

  const initData = () => {
    allData.value = JSON.parse(JSON.stringify(props.allData))
    const res = JSON.parse(allData.value?.list as string)
    containerData.value = []
    res?.forEach((element: any) => {
      // @ts-ignore
      element.name = uiObj[element.key]
      element.style = shallowRef(element.style).value
      containerData.value.push(element)
    })
    Object.keys(formSave.value).forEach((item) => {
      // @ts-ignore
      formSave.value[item] = allData.value[item]
    })
  }

  const index = ref(100)
</script>

<script lang="ts">
  export default {
    name: 'DesignMonitoring',
  }
</script>
<template>
  <div class="design-monitoring">
    <el-row class="toolbar" :gutter="10">
      <el-col v-if="isBtn" :span="2">
        <el-button type="primary" @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
      </el-col>
      <el-col :span="20">
        <vue-draggable
          animation="300"
          chosen-class="chosenClass"
          ghost-class="ghost"
          :group="{ name: 'people', pull: 'clone', put: false }"
          item-key="id"
          :list="toolList"
          :sort="false"
          @start="toolStart"
        >
          <template #item="{ element }">
            <el-button @click="addComponent">{{ element.label }}</el-button>
          </template>
        </vue-draggable>
      </el-col>
      <el-col :span="2">
        <el-button style="float: right" type="primary" @click="saveEvent">保存</el-button>
      </el-col>
    </el-row>
    <vue-draggable
      ref="draggableContent"
      chosen-class="myChosenClass"
      class="content myDraggableContent"
      :disabled="disabled"
      ghost-class="ghost"
      :group="{ name: 'people', pull: true, put: true }"
      item-key="uuid"
      :list="containerData"
      :sort="true"
      @add="toolEnd"
      @end="viewEnd"
      @mousemove="handleClick"
      @start="viewStart"
    >
      <template #item="{ element }">
        <component
          :is="element.name"
          ref="UIRef"
          :current-item="element"
          :style="element.style"
          :time-setting="timeSetting"
          @on-delete-event="handleDelete"
          @on-update-event="handleUpdate"
        />
      </template>
    </vue-draggable>
  </div>
  <el-dialog v-model="visible" title="保存接口监控视图" width="600">
    <el-form ref="formSaveRef" label-width="120px" :model="formSave" :rules="formSaveRules">
      <el-form-item label="监控视图名称：" prop="name">
        <el-input v-model="formSave.name" autocomplete="off" />
      </el-form-item>
      <el-form-item label="监控时长：">
        <el-select v-model="formSave.timeQuantum">
          <el-option v-for="item in timeQuantumOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <vab-date-time-picker v-if="formSave.timeQuantum === 'user-defined'" v-model="formSave.timeDate" />
      </el-form-item>
      <el-form-item label="自动刷新：">
        <el-select v-model="formSave.timeDuration">
          <el-option v-for="item in timeDuratioOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-switch v-model="formSave.switchValue" active-text="自动刷新" style="margin-left: 20px" />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button :disabled="btnDisabled" type="primary" @click="handleSave(formSaveRef)">保存</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<style scoped lang="scss">
  .myChosenClass {
    background-color: rgb(253, 249, 3) !important;
    opacity: 1 !important;
  }
  .toolbar {
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 15px;
  }
  .content {
    position: relative;
    z-index: 0;
    border: 1px solid #eee;
    height: 100%;
    width: 100%;
    height: calc(100vh - 150px);
    overflow-y: auto;
    &::-webkit-scrollbar {
      width: 0;
      height: 0;
    }
  }
</style>
