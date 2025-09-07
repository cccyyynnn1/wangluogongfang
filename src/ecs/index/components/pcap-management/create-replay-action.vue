<script lang="ts">
  export default {
    name: 'CreateReplayAction', //
  }
</script>

<script setup lang="ts">
  import {
    getLocalPcapsApi,
    editLocalPcapsStatusApi,
    deleteLocalPcapsApi,
    addLocalPcapsTaskApi,
  } from '@/api-ecs/packet-replay'
  import type { ComponentSize, FormInstance, FormRules, TableInstance } from 'element-plus'
  import { default as Sortable, SortableEvent } from 'sortablejs'
  import { useScroll } from '@vueuse/core'
  const props = defineProps<{
    modelValue: boolean
    flowProbes: { adapterId: string[]; adapterIdStr: string; name: string; id: string }[]
  }>()
  const emits = defineEmits<{
    (e: 'update:modelValue', visibility: boolean): void
    (e: 'reload'): void
  }>()
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')

  // pack管理：选择pack
  const queryFormPageOne = reactive({
    pageNum: 1,
    pageSize: 20,
    searchStr: '',
  })
  const pcapManagementVisible = useVModel(props, 'modelValue', emits)
  const multipleTableRef = ref<TableInstance>()
  const pcapListPageOne = ref<any[]>([])
  const pcapTotalPageOne = ref(0)
  const loading = ref(true)

  const handleGetLocalAllPcaps = async (isClear = false) => {
    loading.value = true
    const { data } = await getLocalPcapsApi(queryFormPageOne)
    const arr = data.records || []
    if (isClear) {
      pcapListPageOne.value = arr
    } else {
      pcapListPageOne.value.push(...arr)
    }
    pcapTotalPageOne.value = data.total || 0
    loading.value = false
  }

  const handleEditNote = async (row: any) => {
    const { id, remark } = row
    const { msg } = await editLocalPcapsStatusApi({ id, remark })
    row.edit = false
    $baseMessage(msg, 'success', 'vab-hey-message-success')
  }
  const handleDelete = (row: any) => {
    $baseConfirm('确认删除当前Pcap文件么？', null, async () => {
      const { msg } = await deleteLocalPcapsApi({ ids: [row.id] })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      handleGetLocalAllPcaps(true)
    })
  }

  const toggleSelection = (rows?: any[], ignoreSelectable?: boolean) => {
    if (rows) {
      rows.forEach((row) => {
        multipleTableRef.value!.toggleRowSelection(row, true)
      })
    } else {
      multipleTableRef.value!.clearSelection()
    }
  }

  const initData = () => {
    // pcapListPageOne.value = []
    // queryFormPageOne.pageNum = 1
    // handleGetLocalAllPcaps()
    nextTick(() => {
      const dom = document.querySelector('.packet-management-table .el-scrollbar__wrap') as HTMLElement
      const { arrivedState } = useScroll(dom)
      listArrivedState.value = arrivedState
    })
    if (pcapList.value.length > 0) {
      nextTick(() => {
        toggleSelection(pcapList.value)
      })
    }
  }

  onMounted(() => {
    handleGetLocalAllPcaps()
    initData()
  })

  // 新建&排序
  const listArrivedState = ref()
  const replayTaskVisible = ref(false)
  const ruleFormRef = ref<FormInstance>()
  const orderChecked = ref(false)

  const queryForm = reactive({
    flowProbeId: '',
    flowProbeIdStr: '',
    type: 'original',
    typeStr: 'original',
    note: '',
    noteStr: '',
    fileNameList: [] as {
      originName: string
      uniName: string
      remark: string
    }[],
    queryList: {
      pageNum: 1,
      pageSize: 10,
      searchStr: '',
    },
  })
  const queryFormRules = {
    fileNameList: [{ required: true, message: '请选择Pcap文件', trigger: 'blur' }],
    flowProbeId: [{ required: true, message: '请选择链路', trigger: 'blur' }],
  }
  const pcapList = ref<any[]>([])
  const handleGetLocalPcaps = async (
    liss: {
      originName: string
      uniName: string
      remark: string
    }[]
  ) => {
    pcapList.value = liss
  }

  const handleFlowProbeChange = (val: string) => {
    queryForm.flowProbeIdStr = props.flowProbes.find((item) => item.id === val)!.name
  }
  const handleSubmit = () => {
    queryForm.fileNameList = pcapList.value.map(({ originName, remark, uniName }) => ({
      originName,
      remark,
      uniName,
    }))
    ruleFormRef.value?.validate(async (valid) => {
      if (valid) {
        const { queryList, ...query } = queryForm
        const { msg } = await addLocalPcapsTaskApi(query)
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        setTimeout(() => {
          emits('reload')
          pcapManagementVisible.value = false
        }, 200)
      }
    })
  }

  watch(
    () => replayTaskVisible.value,
    () => {
      if (replayTaskVisible.value) {
        nextTick(() => {
          rowDrag()
        })
      } else {
        ruleFormRef.value?.resetFields()
        // pcapList.value = []
        orderChecked.value = false
        initData()
      }
    }
  )

  watch(
    () => listArrivedState.value,
    () => {
      if (listArrivedState.value.bottom && pcapTotalPageOne.value > pcapListPageOne.value.length) {
        queryFormPageOne.pageNum++
        handleGetLocalAllPcaps()
      }
    },
    { deep: true }
  )
  const handleUp = (index: number) => {
    const temp = pcapList.value[index - 1]
    pcapList.value[index - 1] = pcapList.value[index]
    pcapList.value[index] = temp
  }

  const handleToTop = (item: any, index: number) => {
    pcapList.value.splice(index, 1)
    pcapList.value.unshift(item)
  }

  let sortableInstance: any

  watch(
    () => orderChecked.value,
    () => {
      const tbody = document.querySelector('.packet-list-table .el-table__body-wrapper tbody')
      if (orderChecked.value) {
        // 要拖拽元素的父容器
        if (!tbody) return
        tbody.setAttribute('style', 'cursor:move')
        sortableInstance = new Sortable(tbody as HTMLElement, {
          //  可被拖拽的子元素
          draggable: '.packet-list-table .el-table__row',
          direction: 'horizontal',
          delay: 0,
          disabled: !orderChecked.value,
          animation: 250,
          scrollSensitivity: 25,
          scrollSpeed: 50, // px
          forceFallback: true, // 忽略 HTML5拖拽行为，强制回调进行
          fallbackOnBody: true,
          bubbleScroll: false,
          dragoverBubble: true,
          fallbackClass: 'packet-list-table-sortable-fallback', // 当使用forceFallback的时候，被复制的dom的css类名
          onMove(e: any) {
            e.preventDefault()
            e.stop
          },
          // 开始拖拽的时候
          onEnd(event: SortableEvent) {
            if (event.oldIndex !== undefined && event.newIndex !== undefined) {
              const currRow = pcapList.value.splice(event.oldIndex, 1)[0]
              pcapList.value.splice(event.newIndex, 0, currRow)
            }
          },
        })
      } else {
        if (sortableInstance) {
          sortableInstance.destroy()
          if (!tbody) return
          tbody.setAttribute('style', 'cursor:default')
        }
      }
    }
  )

  // 行拖拽
  const rowDrag = function () {
    // 要拖拽元素的父容器
    const tbody = document.querySelector('.packet-list-table .el-table__body-wrapper tbody')
    if (!tbody) return
    tbody.setAttribute('style', 'cursor:default')
    tbody.addEventListener('selectstart', function (event) {
      event.preventDefault()
    })
  }
</script>

<template>
  <div class="create-replay-action">
    <el-dialog v-model="pcapManagementVisible" title="新建任务" width="870px">
      <div v-if="!replayTaskVisible" style="font-size: 14px; line-height: 20px; margin-bottom: 5px">选择pacp：</div>
      <el-table
        v-if="!replayTaskVisible"
        ref="multipleTableRef"
        v-loading="loading"
        border
        class="packet-management-table"
        :data="pcapListPageOne"
        style="height: 550px"
        @selection-change="handleGetLocalPcaps"
      >
        <el-table-column align="center" type="selection" width="55" />
        <el-table-column label="序号" type="index" width="60" />
        <el-table-column label="Pcap名称" prop="originName" show-overflow-tooltip />
        <el-table-column label="备注" show-overflow-tooltip>
          <template #default="{ row }">
            <template v-if="row.edit">
              <el-input v-model="row.remark" clearable @keydown.enter="handleEditNote(row)" />
            </template>
            <template v-else>
              {{ row.remark }}
            </template>
          </template>
        </el-table-column>
        <el-table-column align="center" fixed="right" label="操作" width="140">
          <template #default="{ row }">
            <el-button v-if="row.edit" class="row_action" size="small" type="primary" @click="handleEditNote(row)">
              完成
            </el-button>
            <el-button v-else class="row_action" size="small" @click="row.edit = true">编辑</el-button>
            <el-button class="row_action" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-else class="ioc-replay-task">
        <el-form
          ref="ruleFormRef"
          class="packet-replay-form"
          label-position="top"
          :model="queryForm"
          :rules="queryFormRules"
        >
          <el-form-item label="已选Pcap：" prop="fileNameList" style="position: relative">
            <el-checkbox v-model="orderChecked" label="按顺序回放" style="position: absolute; top: -32px; right: 0" />
            <el-table
              v-loading="loading"
              border
              class="packet-list-table"
              :data="pcapList"
              row-key="id"
              style="height: 300px"
            >
              <el-table-column label="序号" type="index" width="60" />
              <el-table-column label="Pcap名称" prop="originName" show-overflow-tooltip width="350" />
              <el-table-column label="备注" prop="remark" show-overflow-tooltip />
              <el-table-column align="center" fixed="right" label="操作" width="200">
                <template #default="{ row, $index }">
                  <el-button
                    class="row_action"
                    :disabled="$index == 0 || !orderChecked"
                    size="small"
                    @click="handleToTop(row, $index)"
                  >
                    置顶
                  </el-button>
                  <el-button
                    class="row_action"
                    :disabled="$index == 0 || !orderChecked"
                    size="small"
                    @click="handleUp($index)"
                  >
                    上移一层
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-form-item>
          <el-form-item label="链路：" prop="flowProbeId">
            <el-select v-model="queryForm.flowProbeId" style="width: 100%" @change="handleFlowProbeChange">
              <el-option v-for="flow in flowProbes" :key="flow.id" :label="flow.name" :value="flow.id" />
            </el-select>
          </el-form-item>
          <el-form-item label="回放方式" prop="type">
            <el-radio-group v-model="queryForm.type" @change="queryForm.typeStr = $event">
              <el-radio label="original" size="large">原始速度回放</el-radio>
              <el-radio label="fast" size="large">快速回放</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="备注" prop="note">
            <el-input v-model="queryForm.note" @change="queryForm.noteStr = $event" />
          </el-form-item>
        </el-form>
      </div>
      <div v-if="!replayTaskVisible" style="display: flex; justify-content: space-between; width: 100%">
        <!-- <el-pagination
          v-model:current-page="queryFormPageOne.pageNum"
          v-model:page-size="queryFormPageOne.pageSize"
          background
          layout="sizes, prev, pager, next, jumper"
          :page-sizes="[10, 20, 30, 40, 50, 100]"
          :total="pcapTotalPageOne"
        /> -->
        <div></div>
        <el-button
          :disabled="pcapList.length == 0"
          style="margin-top: 12px"
          type="primary"
          @click="replayTaskVisible = true"
        >
          下一步
        </el-button>
      </div>
      <div v-else style="display: flex; justify-content: space-between; width: 100%; margin-top: 12px">
        <div></div>
        <div>
          <el-button :auto-insert-space="false" @click="replayTaskVisible = false">返回上一步</el-button>
          <el-button :auto-insert-space="false" type="primary" @click="handleSubmit">确定</el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<style>
  .packet-list-table-sortable-fallback {
    display: none;
  }
</style>

<style scoped lang="scss">
  :deep(.el-dialog) {
    .el-dialog__body {
      padding-top: 10px !important;
      .el-table__row {
        height: 50px;
      }
    }
  }
  .ioc-replay-task {
    :deep() {
      .packet-list-table {
        .hover-row {
          > .el-table__cell {
            background-color: #fff !important;
          }
        }
        // .el-table__row {
        //   &:hover {
        //     background-color: var(--el-table-row-hover-bg-color) !important;
        //   }
        // }
        .sortable-chosen {
          z-index: 9;
          > .el-table__cell {
            background-color: var(--el-table-row-hover-bg-color);
          }
        }
        .el-table__body {
          cursor: move;
        }
        .el-table__cell {
          padding-block: 8px;
        }
      }
    }
  }
</style>
