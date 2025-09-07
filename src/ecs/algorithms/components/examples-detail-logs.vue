<script setup lang="ts">
  import { Search, Calendar } from '@element-plus/icons-vue'

  import { CollapseModelValue } from 'element-plus'

  import { getContainersLogsApi } from '~/src/api-ecs/algorithms'

  import { containersLogsType } from '~/src/types'

  import dayjs from 'dayjs'

  const props = defineProps<{
    currentItem: any
  }>()

  const searchVal = ref('')

  const textarea = ref()

  let index = 0

  const LastThirtyDays = ref<any>([])

  const listLoading = ref(false) // 是否加载

  const playloadData = reactive<containersLogsType>({
    id: '',
    follow: false,
    stdout: true,
    stderr: true,
    since: undefined,
    until: undefined,
    timestamps: true,
    tail: 'all',
  })

  onMounted(() => {
    playloadData.id = props.currentItem.Id
    getLastThirtyDays()
    changeDate()
  })

  // 获取数据
  const getData = async () => {
    listLoading.value = true
    const res = await getContainersLogsApi({ ...playloadData })
    textarea.value = res
    listLoading.value = false
  }

  // 切换日期
  const changeDate = (date?: any, _index = 0) => {
    if (listLoading.value) {
      return false
    }
    index = _index
    if (date) {
      playloadData.until = date.time
      playloadData.since = date.preTime
    } else {
      playloadData.until = LastThirtyDays.value[0].time
      playloadData.since = LastThirtyDays.value[0].preTime
    }
    getData()
  }

  // 得到最近30天的列表
  const getLastThirtyDays = () => {
    LastThirtyDays.value = []
    const timeDate = dayjs()
    let startDate = timeDate.format('YYYY-MM-DD')
    let now = ~~(new Date().getTime() / 1000)
    for (let index = 1; index < 31; index++) {
      const obj = { date: '', preTime: 0, perDate: '', time: 0 }
      const endDate = timeDate.subtract(index, 'day').format('YYYY-MM-DD')
      obj.date = startDate
      obj.perDate = endDate
      // @ts-ignore
      obj.time = now
      // @ts-ignore
      obj.preTime = ~~(new Date(endDate).getTime(endDate) / 1000)
      LastThirtyDays.value.push(obj)
      now = obj.preTime
      startDate = endDate
    }
  }

  const handleChange = (activeNames: CollapseModelValue) => {
    console.log(activeNames)
  }
</script>

<script lang="ts">
  export default {
    name: 'ExamplesDetailLogs',
  }
</script>

<template>
  <div class="examplesDetailLogs">
    <!-- <div class="header">
      <el-button type="primary">导出</el-button>
       <el-space style="justify-content: space-between; margin-right: -8px">
        <el-input v-memo="searchVal" style="width: 280px" />
        <el-button :icon="Search" type="primary">检索</el-button>
      </el-space>
    </div> -->
    <div class="logs_container">
      <div class="logs_left">
        <div v-for="(item, index) in LastThirtyDays" :key="index" class="item" @click="changeDate(item, index)">
          <el-icon class="header-icon"><Calendar /></el-icon>
          {{ item.date }}
        </div>
      </div>
      <div v-loading="listLoading" class="logs_right">
        <span v-if="textarea">{{ textarea }}</span>
        <el-empty v-else description="description" />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .header {
    display: flex;
    justify-content: space-between;
    margin: 0 0 25px;
  }

  .logs_container {
    display: flex;

    .logs_left {
      width: 220px;
      border: 1px solid var(--el-border-color);
      height: 450px;
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
      .item {
        height: 40px;
        display: flex;
        align-items: center;
        background-color: #f8fbff;
        padding-left: 20px;
        border-bottom: 1px solid var(--el-border-color);
        font-size: 16px;
      }
      &:hover {
        cursor: pointer;
      }
      .header-icon {
        font-size: 16px;
        margin-right: 8px;
        margin-bottom: 2px;
      }
    }

    .logs_right {
      width: 920px;
      height: 450px;
      border: 1px solid var(--el-border-color);
      border-left: 0;
      padding: 35px 30px;
      line-height: 28px;
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
    }

    :deep() {
      .el-collapse-item__header {
        background-color: #f8fbff;
        color: #7a7a7a;
        text-indent: 0.5em;
        padding-left: 1em;
        font-weight: bold;
        user-select: none;
        .header-icon {
          font-size: 16px;
          margin-right: 6px;
        }
      }

      .el-collapse-item__content {
        padding-bottom: 0;
        user-select: none;

        ul {
          list-style: none;
          padding: 0;
          text-align: center;
          line-height: 40px;
          margin: 0;
          max-height: 400px;
          overflow-y: auto;

          li {
            color: #7a7a7ae2;
            cursor: pointer;

            &:hover {
              color: var(--el-color-primary);
              background: #f6fbff;
            }
          }
        }
      }
    }
  }
</style>
