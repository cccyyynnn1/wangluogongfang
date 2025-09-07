<script setup lang="ts">
  import { ref, onMounted, reactive, nextTick } from 'vue'

  import EditWindow from '@/ecs/index/components/interface-monitoring/edit-window.vue'

  import type { FormInstance } from 'element-plus'

  import { MonitoringItem, filedDataType } from '../../../src/types/index'

  import { getTabelListAPI, searchBySqlApi } from '~/src/api-ecs/retrieve'

  import { useUserStore } from '@/store/modules/user'

  import { getBySiteIdApi, siteSearchApi } from '~/src/api-ecs/site'

  import _ from 'lodash'

  import dayjs from 'dayjs'

  const props = defineProps<{
    currentItem?: MonitoringItem
    timeSetting: any
  }>()

  const { getTableColumn } = useUserStore()

  const emit = defineEmits<{
    (e: 'on-delete-event', val: MonitoringItem): void
    (e: 'on-update-event', val: MonitoringItem): void
  }>()

  const formRef = ref<FormInstance>()

  const form = reactive({
    Confirm: '',
  })

  const showpage = ref(false)

  const mode = ref('add')

  const windowTitle = ref('添加窗口')

  // 检验规则
  const validatePass = (rule: any, value: any, callback: any) => {
    if (form.Confirm !== '确认删除') {
      callback(new Error('删除必须填写：确认删除'))
    } else {
      callback()
    }
  }

  const rules = {
    Confirm: [{ validator: validatePass, trigger: 'blur' }],
  }

  const timeSetting = ref()

  const currentItem = ref()

  const dialogFormVisible = ref(false)

  const className = ref('className')

  const resizeObserver = ref()

  const tableColumn = ref<filedDataType[]>([]) // 表头

  const element = ref()

  onMounted(() => {
    initData()
    // document.body.addEventListener('mousedown', observerUITable)

    document.body.addEventListener('mouseup', unObserverUITable)
  })
  onUnmounted(() => {
    unObserverUITable()
    // document.body.removeEventListener('mousedown', observerUITable)
    element.value?.removeEventListener('mousedown', observerUITable)
    document.body.removeEventListener('mouseup', unObserverUITable)
  })

  nextTick(() => {
    element.value = document.querySelector(`.${className.value}`) as Element
    element.value?.addEventListener('mousedown', observerUITable)
  })

  // 刷新数据
  const handleReflash = (val: MonitoringItem) => {
    currentItem.value = val
    emit('on-update-event', val)
  }

  const queryData = reactive({
    pageNum: 1,
    pageSize: 10,
  })

  const showData = async () => {
    listLoading.value = true
    const res = currentItem.value.data
    queryData.pageNum = currentItem.value.data.pageNum
    queryData.pageSize = currentItem.value.data.pageSize
    timeSetting.value = props.timeSetting
    let allField: filedDataType[] = []
    if (res.indexType == 22) {
      allField = getTableColumn(1)
    } else {
      allField = getTableColumn(res.indexType)
      if (res.indexType == 1 && res.module != '2') {
        allField = getTableColumn(30)
        res.displayFields = []
        if (currentItem.value.hasError) return (listLoading.value = false)
        const {
          data: { displayFieldsArr },
        } = await getBySiteIdApi({ id: res.siteSessionId, apiId: res.siteApiId })
        displayFieldsArr.forEach((item: number) => {
          allField.forEach((td) => {
            if (item == td.id) {
              res.displayFields.push(td)
            }
          })
        })
      }
    }
    tableColumn.value = res.displayFields
      ?.map((i: any) => {
        return allField.find((item) => item.id === i.id)
      })
      .filter(Boolean)
    try {
      if (res.module == '0') {
        if (currentItem.value.hasError) return (listLoading.value = false)
        const {
          data: { resList },
        } = await siteSearchApi({
          indexType: res.indexType,
          orderField: res.orderField,
          orderType: res.orderType,
          pageNum: res.pageNum,
          pageSize: res.pageSize,
          startTime: timeSetting.value.startTime,
          endTime: timeSetting.value.endTime,
          siteSessionId: res.siteSessionId,
          siteApiId: res.siteApiId,
        })
        // queryPage.listDate = []
        listDate.value = resList || []
      } else if (res.indexType == 24) {
        const query = {
          topField: 'topField',
          timeRange: '',
          flowProbeIds: ['0'],
          searchTable: 'eventStat',
          searchModel: 'all',
          eventFilterList: [],
          keyword: undefined,
          ip: '',
          clientIp: '',
          serverIp: '',
          serverPort: '',
        }
        query.timeRange = `${timeSetting.value.startTime} - ${timeSetting.value.endTime}`
        query.flowProbeIds = res.flowProbeIds
        query.searchModel = res.searchModel
        query.ip = res.ip
        query.clientIp = res.clientIp
        query.serverIp = res.serverIp
        query.serverPort = res.serverPort
        const keyword = res.ip || res.clientIp || res.serverIp || res.serverPort || undefined
        query.keyword = keyword
        // query.eventFilterList =
        // remark == 'node' ? nodelist.value : eventFilterList.value && JSON.parse(JSON.stringify(eventFilterList.value))
        const res1 = JSON.stringify({ ...query })
        const result = await getTabelListAPI({
          top: res.topValue || '100',
          moduleType: 'eventStat',
          query: res1,
          groupBy: 'clientIp,serverIp,serverPort',
        })
        listDate.value = result.list
      } else {
        const {
          data: { resList },
        } = await searchBySqlApi({
          indexType: res.indexType,
          orderField: res.orderField,
          orderType: res.orderType,
          pageNum: res.pageNum,
          pageSize: res.pageSize,
          searchSql: res.searchSql,
          startTime: timeSetting.value.startTime,
          endTime: timeSetting.value.endTime,
        })
        listDate.value = resList || []
      }
    } finally {
      listLoading.value = false
    }
  }

  // 获取表格序号
  const curIndex = computed(() => (queryData.pageNum - 1) * queryData.pageSize + 1)

  watch(
    () => currentItem.value,
    () => {
      if (currentItem.value && currentItem.value.data.indexType) {
        nextTick(() => {
          showData()
        })
      }
    },
    { immediate: true }
  )

  const listDate = ref<object[]>([]) // 表格数据

  const listLoading = ref(false) // 是否加载
  const observerUITable = _.debounce(() => {
    unObserverUITable()
    const parentNodeWidth = element.value.parentNode.getBoundingClientRect().width
    resizeObserver.value = new ResizeObserver((entries) => {
      let { width, height } = entries[0].contentRect
      height += 20
      let value = ((width + 30) / (parentNodeWidth - 2)) * 100
      value = value > 100 ? 100 : value
      currentItem.value.style.width = `${value}%`
      currentItem.value.style.height = `${height}px`
      currentItem.value.style.long = height + currentItem.value.style.y
      handleReflash(currentItem.value)
    })
    resizeObserver.value.observe(element.value)
  }, 50)

  // 计算可缩放的阈值
  const unObserverUITable = () => {
    resizeObserver.value?.unobserve(element.value)
  }

  const initData = () => {
    if (props.currentItem) {
      currentItem.value = JSON.parse(JSON.stringify(props.currentItem))
      className.value += currentItem.value.uuid
    }
  }

  const handleEdit = () => {
    showpage.value = true
  }

  // 确认删除
  const handleConfirm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl!.validate(async (valid) => {
      if (valid) {
        emit('on-delete-event', currentItem.value)
        dialogFormVisible.value = false
      } else {
        return false
      }
    })
  }

  // 打开删除弹框
  const handleDelete = () => {
    dialogFormVisible.value = true
    nextTick(() => {
      formRef.value!.clearValidate()
      formRef.value!.resetFields()
    })
  }

  // 关闭
  const handleClose = (val: boolean) => {
    showpage.value = val
  }

  function formatDate(row: any, key: string) {
    const time = row[key] / 1000000
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss.SSS')
  }

  defineExpose({
    unObserverUITable,
  })
</script>

<script lang="ts">
  export default {
    name: 'UiTable',
  }
</script>
<template>
  <div class="ui-table" :class="className">
    <div class="top-bar">
      <div class="top-bar-left">{{ currentItem && currentItem.data.name }}</div>
      <div class="top-bar-right">
        <el-icon class="icon" style="margin-right: 6px; opacity: 0.7" @click="handleEdit"><EditPen /></el-icon>
        <el-icon class="icon" style="opacity: 0.7" @click="handleDelete"><Close /></el-icon>
      </div>
    </div>
    <div class="ui-table-content">
      <el-table v-loading="listLoading" :border="true" :data="listDate">
        <el-table-column v-if="listDate.length > 0" :align="'center'" label="序号" width="65">
          <template #default="{ $index }">
            <span>{{ curIndex + $index }}</span>
          </template>
        </el-table-column>
        <template v-for="item in tableColumn" :key="item?.id">
          <el-table-column
            v-if="item.fieldNameCn.includes('时间')"
            align="center"
            min-width="160"
            :prop="item.fieldNameEn"
            :resizable="true"
            show-overflow-tooltip
          >
            <template #header>
              {{ item.fieldNameCn }}
            </template>
            <template #default="{ row }">
              {{ formatDate(row, item.fieldNameEn) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else
            align="center"
            :min-width="120"
            :prop="item.fieldNameEn"
            :resizable="true"
            show-overflow-tooltip
          >
            <template #header>
              {{ item.fieldNameCn }}
            </template>
          </el-table-column>
          <!-- <el-table-column v-else align="center" :label="item.fieldNameCn" min-width="100" :prop="item.fieldNameEn"
                                                                                                                                                                                    :resizable="false" show-overflow-tooltip /> -->
        </template>
        <template #empty>暂无数据</template>
      </el-table>
    </div>
    <Teleport to="body">
      <el-dialog v-model="dialogFormVisible" title="窗口删除" width="500">
        <el-form ref="formRef" :model="form" :rules="rules">
          <el-form-item>
            如果确定删除
            <span
              style="
                color: #2f9bff;
                max-width: 220px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
                word-break: break-all;
                word-wrap: break-word;
              "
            >
              {{ currentItem && currentItem.title }}
            </span>
            窗口，请在下方输入确认删除
          </el-form-item>
          <el-form-item prop="Confirm">
            <el-input v-model="form.Confirm" autocomplete="off" placeholder="确认删除" />
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="dialogFormVisible = false">取消</el-button>
            <el-button type="primary" @click="handleConfirm(formRef)">确认</el-button>
          </span>
        </template>
      </el-dialog>
      <transition>
        <EditWindow
          v-if="showpage"
          :current-item="currentItem"
          :mode="mode"
          :showpage="showpage"
          :title="windowTitle"
          @on-reflash="handleReflash"
          @onCloseEvent="handleClose"
        />
      </transition>
    </Teleport>
  </div>
</template>

<style scoped lang="scss">
  .ui-table {
    position: absolute;
    z-index: 1;
    height: 200px;
    width: 100%;
    resize: both;
    max-width: 100%; /* 最大宽度为父元素的100% */
    max-height: 100%; /* 最大高度为父元素的100% */
    overflow: hidden;
    box-shadow: var(--el-box-shadow-light);
    border-radius: var(--el-card-border-radius);
    border: 1px solid var(--el-card-border-color);
    background-color: #fff;
    color: var(--el-text-color-primary);
    // transition: all 0.05s;
    padding: 5px 15px 15px 15px;
    .top-bar {
      height: 40px;
      padding: 10px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      // background-color: red;
      // border-bottom: 1px solid var(--el-border-color);

      .icon {
        &:hover {
          cursor: pointer;
        }
      }
    }
    .ui-table-content {
      height: calc(100% - 40px);
      :deep() {
        .el-table {
          height: 100%;
          border: 0px solid #fff;
          border-bottom: 1px solid var(--el-border-color);
        }
      }
    }
  }
</style>
