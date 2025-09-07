<script lang="ts">
  export default {
    name: 'PacketReplay',
  }
</script>

<script setup lang="ts">
  import { DArrowRight, DArrowLeft } from '@element-plus/icons-vue'
  import type { FormInstance } from 'element-plus'
  import {
    getPacketReplayListApi,
    getFlowProbeApi,
    deleteFlowProbeApi,
    playBackFlowProbeApi,
  } from '@/api-ecs/packet-replay'
  import { PacketReplayQueryParams } from '@/types'
  import { PacketReplayItem } from '@/types/index'
  import CreateReplayAction from './components/pcap-management/create-replay-action.vue'
  import PcapUpload from './components/pcap-management/pcap-upload.vue'
  import ReplayOverview from './components/replay-overview.vue'
  const $baseConfirm: any = inject('$baseConfirm')
  const $baseMessage: any = inject('$baseMessage')

  const loading = ref(true)
  const showMore = ref(false)
  const tableListData = ref<PacketReplayItem[]>()
  const total = ref(0)
  // 获取表格序号
  const curIndex = computed(() => (replay_form.page - 1) * replay_form.limit + 1)
  // 链路数据
  const flowProbes = ref<{ adapterId: string[]; adapterIdStr: string; name: string; id: string }[]>([])
  const selectedData = ref<PacketReplayItem[]>([])

  const packetReplayForm = ref<FormInstance>()
  const replay_form = reactive<PacketReplayQueryParams>({
    page: 1,
    limit: 20,
    query: {
      flowProbeIdStr: '',
      flowProbeId: [],
      type: undefined,
      fileName: '',
      fingerPoint: '',
      crtTime: '',
      crtTimeArr: [],
      lastRunDateTime: '',
      lastRunDateTimeArr: [],
      lastRunResult: '',
      note: '',
    },
  })

  const uploadType = ref<'online' | 'directory'>('online')
  const uploadVisible = ref(false)
  const newReplayTaskVisible = ref(false)
  const managementVisible = ref(false)
  const replayVisible = ref(false)

  const onReset = () => {
    packetReplayForm.value?.resetFields()
    replay_form.query = {
      flowProbeIdStr: '',
      flowProbeId: [],
      type: undefined,
      fileName: '',
      fingerPoint: '',
      crtTime: '',
      crtTimeArr: [],
      lastRunDateTime: '',
      lastRunDateTimeArr: [],
      lastRunResult: '',
      note: '',
    }
  }

  //重放数据包
  const replayConfirmHandle = (id: string) => {
    $baseConfirm('是否重新回放？', null, async () => {
      const { data } = await playBackFlowProbeApi({
        reqType: 'PUT',
        url: `/packetPlayback/playback/${id}`,
      })
      $baseMessage(data.message, 'success', 'vab-hey-message-success')
      getPacketReplayList()
    })
  }
  //  表格筛选
  const handleSelectionChange = (val: PacketReplayItem[]) => {
    selectedData.value = val
  }

  // 删除数据包回放
  const handleDeleteFlowProbe = async (ids: string[] | 'all') => {
    $baseConfirm(ids === 'all' ? '确认删除所有数据么？' : '确认删除当前数据么？', null, async () => {
      const { data } = await deleteFlowProbeApi({
        reqType: 'DELETE',
        url: ids === 'all' ? `/packetPlayback/deleteAll` : `/packetPlayback/delete/${ids}`,
      })
      $baseMessage(data.message, 'success', 'vab-hey-message-success')
      getPacketReplayList()
    })
  }
  //  数据链路改变
  const handleFlowProbeChange = (vals: string[]) => {
    replay_form.query.flowProbeIdStr = vals
      .map((item) => flowProbes.value?.find((flow) => flow.id === item)?.name)
      .toString()
  }
  // 开始时间改变
  const handleCtrDateChange = (vals: string[]) => {
    replay_form.query.crtTime = vals ? vals.join(' - ') : ''
  }
  // 结束时间改变
  const handleEndDateChange = (vals: string[]) => {
    replay_form.query.lastRunDateTime = vals ? vals.join(' - ') : ''
  }
  // 或许数据链路数据
  const getFlowProbes = async () => {
    const paramsStr = new URLSearchParams({
      page: '1',
      limit: '300',
      query: JSON.stringify({ type: 3 }),
    })
    const { data } = await getFlowProbeApi({
      reqType: 'GET',
      url: `/flowProbe/page?${paramsStr}`,
    })
    flowProbes.value = data.list || []
  }
  // 获取数据包回放list
  const getPacketReplayList = async () => {
    loading.value = true
    try {
      const paramsData = Object.entries(replay_form.query).filter(([key, val]) => val && val.length)
      const _params = Object.fromEntries(paramsData) as Record<string, string>
      const paramsStr = new URLSearchParams({
        query: JSON.stringify(_params),
        page: replay_form.page.toString(),
        limit: replay_form.limit.toString(),
      })
      const { data } = await getPacketReplayListApi({
        reqType: 'GET',
        url: `/packetPlayback/page?${paramsStr}`,
      })
      tableListData.value = data.list || []
      total.value = data.total
    } finally {
      loading.value = false
    }
  }
  const handleShowUpload = (type: 'online' | 'directory') => {
    uploadType.value = type
    uploadVisible.value = true
  }
  onMounted(() => {
    getPacketReplayList()
    getFlowProbes()
  })

  const handleReload = () => {
    replay_form.page = 1
    getPacketReplayList()
  }
</script>

<template>
  <div v-loading="loading" class="packet-replay-container">
    <!-- 数据检索 -->
    <el-form
      ref="packetReplayForm"
      class="packet-replay-form"
      :inline="true"
      label-position="right"
      label-width="95px"
      :model="replay_form.query"
    >
      <el-form-item label="链路" prop="flowProbeId">
        <el-select
          v-model="replay_form.query.flowProbeId"
          clearable
          multiple
          :multiple-limit="2"
          @change="handleFlowProbeChange"
        >
          <el-option v-for="flow in flowProbes" :key="flow.id" :label="flow.name" :value="flow.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="回放方式" prop="type">
        <el-select v-model="replay_form.query.type" clearable>
          <el-option label="原始速度回放" value="original" />
          <el-option label="快速回放" value="fast" />
        </el-select>
      </el-form-item>
      <el-form-item label="文件名" prop="fileName">
        <el-input v-model="replay_form.query.fileName" clearable />
      </el-form-item>
      <el-form-item label="指纹" prop="fingerPoint">
        <el-input v-model="replay_form.query.fingerPoint" clearable />
      </el-form-item>
      <el-form-item label="最后回放结果" prop="lastRunResult">
        <el-input v-model="replay_form.query.lastRunResult" clearable />
      </el-form-item>
      <el-form-item label="创建时间" prop="crtTimeArr">
        <vab-date-time-picker v-model="replay_form.query.crtTimeArr" allow-clear @change="handleCtrDateChange" />
      </el-form-item>
      <el-form-item label="最后回放时间" prop="lastRunDateTimeArr">
        <vab-date-time-picker
          v-model="replay_form.query.lastRunDateTimeArr"
          allow-clear
          @change="handleEndDateChange"
        />
      </el-form-item>
      <el-form-item label="备注" prop="note">
        <el-input v-model="replay_form.query.note" clearable />
      </el-form-item>
    </el-form>
    <!-- <el-form-item class="search-btn" label="&nbsp;">
        <el-button :auto-insert-space="false" type="primary" @click="getPacketReplayList">检索</el-button>
        <el-button :auto-insert-space="false" plain @click="onReset">重置</el-button>
        <el-button :auto-insert-space="false" class="more-btn" plain @click="() => (showMore = !showMore)">
          {{ showMore ? '收起' : '展开' }}
          <el-icon>
            <DArrowRight v-if="showMore" />
            <DArrowLeft v-else />
          </el-icon>
        </el-button>
      </el-form-item> -->
    <el-space style="justify-content: space-between" wrap>
      <div>
        <el-button
          :auto-insert-space="false"
          style="margin-right: 10px"
          type="primary"
          @click="newReplayTaskVisible = true"
        >
          新建任务
        </el-button>
        <el-dropdown style="margin-right: 10px">
          <span class="el-dropdown-link">
            <el-button :auto-insert-space="false" type="primary">
              上传
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleShowUpload('online')">在线上传</el-dropdown-item>
              <el-dropdown-item @click="handleShowUpload('directory')">目录挂载</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-dropdown style="margin-right: 10px">
          <span class="el-dropdown-link">
            <el-button type="danger">
              批量删除
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleDeleteFlowProbe(selectedData.map((i) => i.id))">
                删除选中
              </el-dropdown-item>
              <el-dropdown-item @click="handleDeleteFlowProbe('all')">删除所有</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button :auto-insert-space="false" plain style="margin-right: 10px" @click="replayVisible = true">
          数据回放展示
        </el-button>
        <!-- <el-button :auto-insert-space="false" plain @click="managementVisible = true">Pcap文件管理</el-button> -->
        <!-- <el-button :auto-insert-space="false" plain @click="a2 = true">数据包调查</el-button>
      <el-button :auto-insert-space="false" plain @click="a1 = true">数据包告警</el-button> -->
      </div>
      <div style="margin-right: -8px">
        <el-button :auto-insert-space="false" type="primary" @click="getPacketReplayList">检索</el-button>
        <el-button :auto-insert-space="false" plain @click="onReset">重置</el-button>
      </div>
    </el-space>
    <el-table class="packet-replay-table" :data="tableListData" @selection-change="handleSelectionChange">
      <el-table-column align="center" type="selection" width="50" />
      <el-table-column align="center" :index="($index) => curIndex + $index" label="序号" type="index" width="60" />
      <el-table-column label="链路" prop="flowProbeIdStr" show-overflow-tooltip width="85" />
      <el-table-column label="回放方式" prop="typeStr" show-overflow-tooltip width="105" />
      <el-table-column label="文件名" prop="fileName" show-overflow-tooltip />
      <el-table-column label="文件大小" prop="fileSizeStr" show-overflow-tooltip width="85" />
      <el-table-column label="指纹" prop="fingerPoint" show-overflow-tooltip />
      <el-table-column label="创建时间" prop="crtTimeStr" show-overflow-tooltip width="165" />
      <el-table-column label="最后回放结果" prop="lastRunResult" show-overflow-tooltip width="140" />
      <el-table-column label="最后回放时间" prop="lastRunDateTimeStr" show-overflow-tooltip width="165" />
      <el-table-column label="备注" prop="note" show-overflow-tooltip width="100px" />
      <el-table-column align="center" label="操作" width="160px">
        <template #default="{ row }">
          <el-button plain size="small" @click="replayConfirmHandle(row.id)">重新回放</el-button>
          <el-button plain size="small" @click="handleDeleteFlowProbe([row.id])">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      v-model:current-page="replay_form.page"
      v-model:page-size="replay_form.limit"
      background
      class="pagination"
      layout="total, sizes, prev, pager, next, jumper"
      :page-sizes="[20, 30, 50]"
      :total="total"
      @current-change="getPacketReplayList"
      @size-change="getPacketReplayList"
    />
    <!-- 新建任务 -->
    <CreateReplayAction
      v-if="newReplayTaskVisible"
      v-model="newReplayTaskVisible"
      :flow-probes="flowProbes"
      @reload="handleReload"
    />
    <!-- <new-replay-task
      v-model="newReplayTaskVisible"
      :flow-probes="flowProbes"
      @reload="
        () => {
          replay_form.page = 1
          getPacketReplayList()
        }
      "
    /> -->
    <!-- <pcap-management v-model="managementVisible" /> -->
    <pcap-upload v-if="uploadVisible" v-model="uploadVisible" :upload-type="uploadType" />
    <replay-overview v-if="replayVisible" @close="replayVisible = false" />
    <!-- <el-dialog v-model="a1" class="alarm-packet-replay" fullscreen title="数据包告警" width="100%">
      <alarm-packet-replay />
    </el-dialog>
    <el-dialog v-model="a2" class="retrieve-packet-replay" fullscreen title="数据包调查" width="100%">
      <retrieve-packet-replay />
    </el-dialog> -->
  </div>
</template>

<style scoped lang="scss">
  .packet-replay-container {
    display: flex;
    flex-direction: column;
    .packet-replay-form {
      width: 100%;
      .el-form-item {
        width: 25%;
        margin-right: 0 !important;
        padding-right: 10px;
        &:nth-child(4n) {
          padding-right: 0;
        }
        &.search-btn {
          margin-right: 0;
          :deep(.el-form-item__content) {
            justify-content: end;
          }
        }
        .el-select {
          width: 100%;
        }
        .more-btn {
          .el-icon {
            transform: translateX(5px) rotate(-90deg);
          }
        }
      }
    }
    .packet-replay-table {
      width: 100%;
      margin-top: 20px;
      flex: 1;
    }
    .packet-replay-overview {
      position: absolute;
      inset: 0;
      z-index: 10;
    }
  }
</style>
