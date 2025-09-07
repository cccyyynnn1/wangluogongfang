<script lang="ts">
  export default {
    name: 'NewReplayTask',
  }
</script>

<script setup lang="ts">
  import { getLocalPcapsApi, addLocalPcapsTaskApi } from '@/api-ecs/packet-replay'
  import type { ComponentSize, FormInstance, FormRules } from 'element-plus'
  // @ts-ignore
  import { default as Sortable, SortableEvent } from 'sortablejs'
  import { useScroll } from '@vueuse/core'
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')
  const props = defineProps<{
    modelValue: boolean
    flowProbes: { adapterId: string[]; adapterIdStr: string; name: string; id: string }[]
  }>()
  const emits = defineEmits<{
    (e: 'update:modelValue', visibility: boolean): void
    (e: 'reload'): void
  }>()
  const listArrivedState = ref()
  const replayTaskVisible = useVModel(props, 'modelValue', emits)
  const ruleFormRef = ref<FormInstance>()
  const loading = ref(true)
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
  const pcapTotal = ref(0)
  const handleGetLocalPcaps = async () => {
    loading.value = true
    const { data } = await getLocalPcapsApi(queryForm.queryList)
    const arr = data.records || []
    pcapList.value.push(...arr)
    pcapTotal.value = data.total || 0
    loading.value = false
  }
  const handleSelectionChange = (
    liss: {
      originName: string
      uniName: string
      remark: string
    }[]
  ) => {
    queryForm.fileNameList = liss.map(({ originName, remark, uniName }) => ({ originName, remark, uniName }))
  }
  const handleFlowProbeChange = (val: string) => {
    queryForm.flowProbeIdStr = props.flowProbes.find((item) => item.id === val)!.name
  }
  const handleSubmit = () => {
    ruleFormRef.value?.validate(async (valid) => {
      if (valid) {
        const { queryList, ...query } = queryForm
        const { msg } = await addLocalPcapsTaskApi(query)
        emits('reload')
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        replayTaskVisible.value = false
      }
    })
  }
  watch(
    () => replayTaskVisible.value,
    () => {
      if (replayTaskVisible.value) {
        handleGetLocalPcaps()
        nextTick(() => {
          const dom = document.querySelector('.packet-list-table .el-scrollbar__wrap') as HTMLElement
          const { arrivedState } = useScroll(dom)
          listArrivedState.value = arrivedState
          rowDrag()
        })
      } else {
        ruleFormRef.value?.resetFields()
      }
    }
  )
  watch(
    () => listArrivedState.value,
    () => {
      if (listArrivedState.value.bottom && pcapTotal.value > pcapList.value.length) {
        queryForm.queryList.pageNum++
        handleGetLocalPcaps()
      }
    },
    { deep: true }
  )
  const startX = ref(0)
  const startY = ref(0)
  const top = ref(0)
  const left = ref(0)
  // 行拖拽
  const rowDrag = function () {
    // 要拖拽元素的父容器
    const tbody = document.querySelector('.packet-list-table .el-table__body-wrapper tbody')
    if (!tbody) return
    Sortable.create(tbody as HTMLElement, {
      //  可被拖拽的子元素
      draggable: '.packet-list-table .el-table__row',
      // scroll: true,
      onEnd(event: SortableEvent) {
        if (event.oldIndex !== undefined && event.newIndex !== undefined) {
          const currRow = pcapList.value.splice(event.oldIndex, 1)[0]
          pcapList.value.splice(event.newIndex, 0, currRow)
        }
      },
      onStart: function (/**Event*/ event: any) {
        //       startX.value = event.clientX - draggable.value.offsetLeft;
        // startY.value = event.clientY - draggable.value.offsetTop;
        // // 添加鼠标移动事件监听
        // document.addEventListener('mousemove', doDrag);
        // // 添加鼠标释放事件监听
        // document.addEventListener('mouseup', stopDrag);
      },
    })
  }
</script>

<template>
  <div class="ioc-replay-task">
    <el-dialog v-model="replayTaskVisible" class="ioc-replay-task" title="新建任务" width="870px">
      <el-form
        ref="ruleFormRef"
        class="packet-replay-form"
        label-position="top"
        :model="queryForm"
        :rules="queryFormRules"
      >
        <el-form-item label="选择Pcap：" prop="fileNameList">
          <el-table
            v-loading="loading"
            border
            class="packet-list-table"
            :data="pcapList"
            row-key="id"
            style="height: 300px"
            @selection-change="handleSelectionChange"
          >
            <el-table-column type="selection" width="55" />
            <el-table-column label="序号" type="index" width="60" />
            <el-table-column label="Pcap名称" prop="originName" show-overflow-tooltip />
            <el-table-column label="备注" prop="remark" show-overflow-tooltip />
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
      <template #footer>
        <div class="dialog-footer">
          <el-button :auto-insert-space="false" type="primary" @click="handleSubmit">确定</el-button>
          <el-button :auto-insert-space="false" @click="replayTaskVisible = false">取消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .ioc-replay-task {
    :deep() {
      .packet-list-table {
        .el-table__row {
          transition: all 0.5s;
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
