<script setup lang="ts">
  import VabDialog from '@/plugins/VabDialog/index.vue'

  import ExamplesDetailOverview from './examples-detail-overview.vue'

  import ExamplesDetailProcess from './examples-detail-process.vue'

  import ExamplesDetailLogs from './examples-detail-logs.vue'

  import ExamplesDetailUpload from './examples-detail-upload.vue'

  import { Search } from '@element-plus/icons-vue'

  import { ExamplesDetailType } from '../type'

  import { getAlgorithmsContainersApi, getHardInfoApi, getContainersStatsApi } from '~/src/api-ecs/algorithms'

  import { algorithmsContainersType } from '~/src/types'

  import VueEvent from '@/data/event'

  const queryForm = reactive<algorithmsContainersType>({
    all: true,
    limit: undefined,
    size: false,
    filters: { name: undefined, status: undefined },
  })

  const listLoading = ref(false) // 是否加载

  const currentItem = ref() // 当前实例

  const instanceLsit = ref<any>([])

  // 硬件信息
  const hardInfo = reactive({
    cpuRate: { usage: 0, unused: 0 },
    memory: {
      free: '',
      rate: 0,
      total: '',
    },
  })

  // 状态选项
  const statusOptions = [
    { label: '运行中', value: 'running' },
    { label: '已创建', value: 'created' },
    { label: '启动中', value: 'restarting' },
    { label: '移除中', value: 'removing' },
    { label: '已暂停', value: 'paused' },
    { label: '退出的', value: 'exited' },
    { label: '已注销', value: 'dead' },
  ]

  // const allIds = ref<string[]>([])

  const examples_detail_visible = ref(false)

  const detailTabsValue = ref<ExamplesDetailType>('overview')

  // 关闭弹窗
  const closeEvent = () => {
    detailTabsValue.value = 'overview'
    examples_detail_visible.value = false
  }

  // 获取硬件信息
  const getHardInfo = async () => {
    const { data } = await getHardInfoApi()
    const obj = data.node.value && JSON.parse(data.node.value)
    Object.keys(hardInfo).forEach((item) => {
      // @ts-ignore
      hardInfo[item] = obj[item]
    })
  }

  // 得到运行状态
  const getStatus = (val: string) => {
    return statusOptions.find((item) => {
      return item.value == val
    })
  }

  // 获取容器列表
  const getAlgorithmsContainers = async () => {
    listLoading.value = true
    const res = await getAlgorithmsContainersApi({ ...queryForm })
    instanceLsit.value = res
    listLoading.value = false
    // allIds.value = []
    // instanceLsit.value.forEach((item: any) => {
    //   allIds.value.push(item.Id)
    // })
    // if (allIds.value.length > 0) {
    //   getAllContainersStats()
    // }
  }

  // 获取全部容器CPU使用率和内存使用率
  const getAllContainersStats = async (item: string) => {
    // todo 待优化
    // allIds.value.forEach(async (item: string) => {
    const res = await getContainersStatsApi({ id: item, stream: false, oneShot: true })
    // @ts-ignore
    const RAMUse = res.memory_stats.usage - (res.memory_stats.stats?.cache || 0)
    // @ts-ignore
    const RAMRate = ((RAMUse / res.memory_stats.limit) * 100.0).toFixed(2)
    // @ts-ignore
    const cpu_delta = res.cpu_stats.cpu_usage.total_usage - res.precpu_stats.cpu_usage.total_usage
    // @ts-ignore
    const system_cpu_delta = res.cpu_stats.system_cpu_usage - res.precpu_stats.system_cpu_usage
    // @ts-ignore
    const number_cpus = res.cpu_stats.online_cpus
    const CPURate = ((cpu_delta / system_cpu_delta) * number_cpus * 100.0).toFixed(2)
    instanceLsit.value.forEach((td: any) => {
      if (item == td.Id) {
        // @ts-ignore
        td.RAMRate = RAMRate
        // @ts-ignore
        td.CPURate = CPURate
        td.containersStats = res
      }
    })
    // })
  }

  function handel_examples_click(params: any) {
    currentItem.value = params
    examples_detail_visible.value = true
  }

  // 懒加载
  const vLazy = {
    mounted: (el: Element, binding: any) => {
      // 2. 创建一个观察对象，来观察当前使用指令的元素
      const observe = new IntersectionObserver(
        ([{ isIntersecting }]) => {
          if (isIntersecting) {
            // 停止观察
            observe.unobserve(el)
            getAllContainersStats(binding.value)
          }
        },
        {
          threshold: 0,
        }
      )
      // 开启观察
      observe.observe(el)
    },
  }

  onMounted(() => {
    getHardInfo()
    getAlgorithmsContainers()
  })

  // 刷新
  const containerRefresh = () => {
    getHardInfo()
    getAlgorithmsContainers()
  }

  // 监听事件
  VueEvent.on('containerRefresh', containerRefresh)
  // 事件销毁
  onUnmounted(() => {
    VueEvent.off('containerRefresh')
  })
</script>

<script lang="ts">
  export default {
    name: 'AlgorithmsExamples',
  }
</script>

<template>
  <!-- 查询表单 -->
  <el-form class="demo-form-inline" inline :model="queryForm" style="display: none">
    <!-- <el-form-item label="算法名称">
      <el-input v-model="queryForm.name" />
    </el-form-item> -->
    <el-form-item label="镜像名称">
      <el-input v-model="queryForm.filters.name" />
    </el-form-item>
    <el-form-item label="运行状态">
      <el-select v-model="queryForm.filters.status" clearable>
        <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
    </el-form-item>
    <el-form-item>
      <el-button :icon="Search" :loading="listLoading" type="primary" @click="getAlgorithmsContainers">检索</el-button>
    </el-form-item>
  </el-form>
  <!-- CPU & 内存 使用率 -->
  <el-row :gutter="20">
    <el-col :span="12">
      <vab-card skeleton>
        <div class="utilization-rate">
          <div class="utilization-rate-left">
            <div class="utilization-rate-title">CPU使用率</div>
            <div class="utilization-rate-icon">
              <el-image fit="cover" :src="require('@/assets/algorithms_images/cpu.png')" />
            </div>
          </div>
          <div class="utilization-rate-right">
            <div class="utilization-rate-cpu_info">{{ Number(hardInfo.cpuRate.usage).toFixed(2) }}%</div>
            <el-progress :percentage="+hardInfo.cpuRate.usage" :show-text="false" :stroke-width="8" />
          </div>
        </div>
      </vab-card>
    </el-col>
    <el-col :span="12">
      <vab-card skeleton>
        <div class="utilization-rate">
          <div class="utilization-rate-left">
            <div class="utilization-rate-title">内存使用率</div>
            <div class="utilization-rate-icon">
              <el-image fit="cover" :src="require('@/assets/algorithms_images/memory.png')" />
            </div>
          </div>
          <div class="utilization-rate-right">
            <div class="utilization-rate-memory_info">{{ hardInfo.memory.rate }}%</div>
            <el-progress :percentage="hardInfo.memory.rate" :show-text="false" status="success" :stroke-width="8" />
          </div>
        </div>
      </vab-card>
    </el-col>
  </el-row>
  <!-- 算法实例列表 -->
  <el-row :gutter="20">
    <el-col :span="24">
      <vab-card class="examples-container" skeleton>
        <div class="examples-container-title">算法实例</div>
        <el-row :gutter="20">
          <template v-for="item in instanceLsit" :key="item.Id">
            <el-col v-lazy="item.Id" :span="6" @click="handel_examples_click(item)">
              <vab-card shadow="hover">
                <template #header>
                  <div class="examples-content-title">{{ item.Names[0] }}</div>
                  <div class="examples-content-description">
                    <el-tooltip :content="item.Command" effect="dark" placement="bottom">
                      {{ item.Command }}
                    </el-tooltip>
                  </div>
                  <div class="examples-content-status" :class="{ success: getStatus(item.State)?.value === 'running' }">
                    {{ getStatus(item.State)?.label }}
                  </div>
                </template>
                <template #default>
                  <div class="utilization-rate-content">
                    <div class="utilization-rate-content-info">
                      <span>CPU</span>
                      <span>{{ item.CPURate }}%</span>
                    </div>
                    <el-progress :percentage="+item.CPURate" :show-text="false" :stroke-width="8" />
                  </div>
                  <div class="utilization-rate-content">
                    <div class="utilization-rate-content-info">
                      <span>RAM</span>
                      <span>{{ item.RAMRate }}%</span>
                    </div>
                    <el-progress :percentage="+item.RAMRate" :show-text="false" :stroke-width="8" />
                  </div>
                </template>
              </vab-card>
            </el-col>
          </template>
        </el-row>
      </vab-card>
    </el-col>
  </el-row>
  <vab-dialog
    v-model="examples_detail_visible"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    title="详情"
    width="1175px"
    @close="closeEvent"
  >
    <div class="ip-detail">
      <el-tabs v-model="detailTabsValue">
        <el-tab-pane label="总览" name="overview">
          <examples-detail-overview :current-item="currentItem" />
        </el-tab-pane>
        <el-tab-pane label="进程" name="process">
          <examples-detail-process :current-item="currentItem" />
        </el-tab-pane>
        <!-- todo  -->
        <!-- <el-tab-pane label="日志" name="logs">
          <examples-detail-logs :current-item="currentItem" />
        </el-tab-pane> -->
        <!-- <el-tab-pane label="上传文件" name="upload">
          <examples-detail-upload :current-item="currentItem" @closeEvent="closeEvent" />
        </el-tab-pane> -->
      </el-tabs>
    </div>
  </vab-dialog>
</template>

<style scoped lang="scss">
  :deep() {
    .ip-detail {
      min-height: 600px;
      .el-tabs__header {
        margin-bottom: 30px;
      }
    }
  }

  .utilization-rate {
    display: flex;
    align-items: flex-end;
    color: #606266;

    .utilization-rate-left {
      width: 80px;
      margin-right: 40px;

      .utilization-rate-title {
        font-size: 14px;
        font-weight: 500;
        text-align: center;
      }

      .utilization-rate-icon {
        width: 60px;
        height: 60px;
        margin: 8px auto 0;
      }
    }

    .utilization-rate-right {
      flex: 1;
      padding-bottom: 10px;
      font-size: 16px;

      .utilization-rate-cpu_info {
        color: var(--el-color-primary);
        margin-bottom: 5px;
      }

      .utilization-rate-memory_info {
        color: var(--el-color-success);
        margin-bottom: 5px;
        display: inline-block;
      }
    }
  }

  .examples-container-title,
  .examples-content-title,
  .examples-content-description,
  .utilization-rate-content-info,
  .examples-content-status {
    color: #606266;
    margin: 0 0 20px;
  }
  .examples-content-title {
    font-weight: bold;
    font-size: 16px;
  }

  .examples-content-description {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    word-break: break-all;
    word-wrap: break-word;
    height: 24px;
    line-height: 24px;
    margin-bottom: 10px !important;
    margin-top: -10px !important;
  }

  .examples-content-status {
    margin-bottom: 0 !important;

    &.success {
      &::before {
        background-color: var(--el-color-success);
      }
    }

    &::before {
      content: ' ';
      display: inline-block;
      width: 14px;
      height: 14px;
      border-radius: 100%;
      background: var(--el-color-error);
      vertical-align: middle;
      margin-right: 5px;
      margin-bottom: 3px;
    }
  }

  .examples-container .vab-card {
    cursor: pointer;
  }

  .utilization-rate-content {
    margin-bottom: 0;

    &:not(:first-child) {
      margin-top: 20px;
    }

    .utilization-rate-content-info {
      margin-bottom: 5px;
      display: flex;
      justify-content: space-between;
    }
  }
</style>
