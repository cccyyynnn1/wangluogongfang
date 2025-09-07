<script setup lang="ts">
  import {
    restartContainersByIdApi,
    startContainersByIdApi,
    stopContainersByIdApi,
    getContainersInfoApi,
  } from '~/src/api-ecs/algorithms'

  import { ExamplesCheckType, OperationType } from '../type'

  import { ElDivider } from 'element-plus'

  import dayjs from 'dayjs'

  import 'dayjs/locale/zh-cn'

  import { useTableCopy } from '@/utils'

  dayjs.locale('zh-cn')

  const relativeTime = require('dayjs/plugin/relativeTime')

  dayjs.extend(relativeTime)

  const props = defineProps<{
    currentItem: any
  }>()

  const currentItem = ref<any>()

  const listDate = ref<object[]>([]) // 表格数据

  const spacer = h(ElDivider, { direction: 'vertical' })

  const examples_check_type = ref<ExamplesCheckType>('port')

  const loading = ref(false)

  const btnDisabled = ref(false)

  const textarea = ref()

  const echo = reactive({
    name: '',
    versions: '',
    cmd: '--',
    time: '',
  })

  const checkType: {
    [key in ExamplesCheckType]: {
      label: string
      prop: string
    }[]
  } = {
    port: [
      {
        label: '本地端口',
        prop: 'PrivatePort',
      },
      {
        label: '容器端口',
        prop: 'PublicPort',
      },
      {
        label: '类型端口',
        prop: 'Type',
      },
    ],
    storage: [
      {
        label: '文件/文件夹',
        prop: 'Driver',
      },
      {
        label: '装载路径',
        prop: 'Destination',
      },
      {
        label: '类型',
        prop: 'Type',
      },
    ],
    links: [
      {
        label: '容器名称',
        prop: 'containersName',
      },
      {
        label: '别名',
        prop: 'name',
      },
    ],
    internet: [
      {
        label: '网络名称',
        prop: 'netName',
      },
      {
        label: '驱动程序 ',
        prop: 'drivers',
      },
    ],
  }

  onMounted(() => {
    initData()
  })

  // 初始化数据
  const initData = async () => {
    currentItem.value = JSON.parse(JSON.stringify(props.currentItem))
    echo.name = currentItem.value.Image?.split('/')[0] || ''
    echo.versions = currentItem.value.Image?.split('/')[1] || ''
    currentItem.value.containersInfo = await getContainersInfoApi({ id: currentItem.value.Id })
    textarea.value = currentItem.value.containersInfo.Config.Env
    echo.cmd = currentItem.value.containersInfo.Config.Cmd.join(' ')
    // @ts-ignore
    echo.time = dayjs(currentItem.value.Created * 1000).fromNow()
  }

  // 启动，停止，重启
  const operationEvent = async (val: OperationType) => {
    const res = undefined
    btnDisabled.value = true
    switch (val) {
      case 'start':
        await startContainersByIdApi({ id: currentItem.value.Id })
        ElMessage({ message: '启动成功', type: 'success' })
        btnDisabled.value = false
        break
      case 'stop':
        stopContainersByIdApi({ id: currentItem.value.Id })
        ElMessage({ message: '停止成功', type: 'success' })
        setTimeout(() => {
          btnDisabled.value = false
        }, 500)
        break
      case 'restart':
        restartContainersByIdApi({ id: currentItem.value.Id })
        ElMessage({ message: '重启成功', type: 'success' })
        setTimeout(() => {
          btnDisabled.value = false
        }, 500)
        break
    }
  }

  // 获取表格数据
  watch(
    () => examples_check_type.value,
    () => {
      loading.value = true
      if (examples_check_type.value == 'port') {
        setTimeout(() => {
          const res = JSON.parse(JSON.stringify(currentItem.value.Ports))
          const arr: string[] = []
          res.forEach((item: any, index: number) => {
            if (arr.includes(res[index].PrivatePort)) {
              delete res[index]
            } else {
              arr.push(res[index].PrivatePort)
            }
          })
          listDate.value = res.filter((item: any) => {
            return item
          })
        }, 0)
      } else if (examples_check_type.value == 'storage') {
        listDate.value = []
        currentItem.value.Mounts.forEach((item: any) => {
          const obj = { Driver: '', Destination: '', Type: '' }
          obj.Driver = item.Driver
          obj.Destination = item.Destination
          obj.Type = item.Type
          listDate.value.push(obj)
        })
      } else if (examples_check_type.value == 'links') {
        listDate.value = [{ containersName: currentItem.value.Image }]
      } else if (examples_check_type.value == 'internet') {
        const value4 = currentItem.value.containersInfo.NetworkSettings.Networks.bridge
        listDate.value = [{ netName: value4?.Gateway, drivers: value4?.DriverOpts }]
      }
      loading.value = false
    },
    { immediate: true }
  )

  function checkExamplesType(type: ExamplesCheckType) {
    loading.value = true
    examples_check_type.value = type
    loading.value = false
  }

  const columnList = computed(() => checkType[examples_check_type.value])
</script>

<script lang="ts">
  export default {
    name: 'ExamplesDetailOverview',
  }
</script>

<template>
  <div v-loading="loading" class="ExamplesDetailOverview">
    <el-space alignment="center" :size="0">
      <div>
        <span class="title">{{ echo.name }}</span>
        <span>{{ echo.versions }}</span>
      </div>
      <el-space alignment="center" :size="8">
        <el-button :loading="btnDisabled" type="primary" @click="operationEvent('start')">启动</el-button>
        <el-button :loading="btnDisabled" type="primary" @click="operationEvent('stop')">停止</el-button>
        <el-button :loading="btnDisabled" type="primary" @click="operationEvent('restart')">重新启动</el-button>
        <!-- <el-button type="primary">强制停止</el-button> -->
      </el-space>
    </el-space>
    <el-row :gutter="20" style="margin-top: 20px; margin-bottom: 20px">
      <el-col :span="14">
        <el-descriptions border :column="1">
          <el-descriptions-item label="启用时间">{{ echo.time }}</el-descriptions-item>
          <el-descriptions-item label="CPU优先顺序">中</el-descriptions-item>
          <el-descriptions-item label="内存限制">自动</el-descriptions-item>
          <el-descriptions-item label="执行命令">{{ echo.cmd }}</el-descriptions-item>
        </el-descriptions>
      </el-col>
      <el-col :span="10">
        <el-row>
          <el-col :span="24">
            <div class="utilization-rate">
              <el-image fit="cover" :src="require('@/assets/algorithms_images/cpu.png')" />
              <div class="utilization-rate-right">
                <div class="utilization-rate-info">
                  <span>CPU使用率</span>
                  <span>{{ currentItem?.CPURate }}%</span>
                </div>
                <el-progress :percentage="Number(currentItem?.CPURate)" :show-text="false" :stroke-width="8" />
              </div>
            </div>
          </el-col>
          <el-col :span="24">
            <div class="utilization-rate">
              <el-image fit="cover" :src="require('@/assets/algorithms_images/memory.png')" />
              <div class="utilization-rate-right">
                <div class="utilization-rate-info">
                  <span>内存使用率</span>
                  <span>{{ currentItem?.RAMRate }}%</span>
                </div>
                <el-progress
                  :percentage="Number(currentItem?.RAMRate)"
                  :show-text="false"
                  status="success"
                  :stroke-width="8"
                />
              </div>
            </div>
          </el-col>
        </el-row>
      </el-col>
    </el-row>
    <el-space alignment="center" :size="4" :spacer="spacer" style="margin-bottom: 20px">
      <el-button link :type="examples_check_type === 'port' ? 'primary' : 'default'" @click="checkExamplesType('port')">
        端口设置
      </el-button>
      <el-button
        link
        :type="examples_check_type === 'storage' ? 'primary' : 'default'"
        @click="checkExamplesType('storage')"
      >
        储存空间
      </el-button>
      <el-button
        link
        :type="examples_check_type === 'links' ? 'primary' : 'default'"
        @click="checkExamplesType('links')"
      >
        链接
      </el-button>
      <el-button
        link
        :type="examples_check_type === 'internet' ? 'primary' : 'default'"
        @click="checkExamplesType('internet')"
      >
        网络
      </el-button>
    </el-space>
    <el-table
      v-loading="loading"
      align="center"
      border
      :data="listDate"
      style="width: 100%"
      @cell-contextmenu="useTableCopy"
    >
      <el-table-column
        v-for="column in columnList"
        :key="column.prop"
        align="center"
        :label="column.label"
        :prop="column.prop"
      />
    </el-table>
    <div class="title margin">环境变量</div>
    <el-input v-model="textarea" :autosize="{ minRows: 4 }" resize="none" type="textarea" />
  </div>
</template>

<style scoped lang="scss">
  .ExamplesDetailOverview {
    .title {
      font-size: 18px;
      color: #303133;
      margin-right: 25px;
      font-weight: 500;

      &.margin {
        margin-top: 30px;
      }

      & + span {
        vertical-align: top;
      }
    }

    & > .el-space--horizontal {
      &:first-child {
        width: 100%;
        justify-content: space-between;
      }

      :deep() {
        & > .el-space__item {
          &:nth-child(2) {
            margin-right: -8px !important;
          }
        }
      }
    }

    .utilization-rate {
      display: flex;
      align-items: top;
      margin-top: 10px;
      margin-bottom: 10px;

      .el-image {
        width: 60px;
      }

      .utilization-rate-right {
        flex: 1;
        margin-left: 15px;

        .utilization-rate-info {
          display: flex;
          justify-content: space-between;
          font-size: 18px;
          color: #303133;
          margin-bottom: 20px;
          margin-top: 4px;
          font-weight: 500;

          span {
            &:last-child {
              font-size: 14px;
            }
          }
        }
      }
    }

    :deep() {
      .el-descriptions__label.el-descriptions__cell.is-bordered-label {
        width: 168px;
      }

      .el-textarea {
        margin: 16px 0 26px;
        color: #303133;

        .el-textarea__inner {
          background: #f8fbff;
          padding: 20px;

          &:focus {
            box-shadow: 0 0 0 1px var(--el-border-color) inset;
          }
        }
      }
    }
  }
</style>
