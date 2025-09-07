<script setup lang="ts">
  import { equipmentListApi } from '@/api-ecs/equipment'

  const curId = ref('')

  const emit = defineEmits<{
    (e: 'on-change-menu', curId: string): void
  }>()

  const router = useRouter()

  let listDate = reactive<{ id: string; name: string }[]>([]) // 表格数据

  const queryPage = reactive({
    page: 1,
    limit: 10,
    query: {
      keyword: undefined,
    },
  })

  const listLoading = ref(false) // 是否加载

  // 得到设备列表
  const getDeviceList = async () => {
    try {
      listLoading.value = true
      const { data } = await equipmentListApi({ ...queryPage, query: JSON.stringify(queryPage.query) })
      listDate = data.list
    } finally {
      listLoading.value = false
    }
  }

  // 跳转到采集设备
  const skip = () => {
    router.push({ path: '/equipment/collecting-device' })
  }

  // 点击某一项
  const handleClick = (id: string) => {
    curId.value = id
    emit('on-change-menu', curId.value)
  }

  onMounted(() => {
    getDeviceList()
  })
</script>

<script lang="ts">
  export default {
    name: 'DeviceList',
  }
</script>
<template>
  <div class="device-list">
    <div class="top">
      <el-input v-model="queryPage.query.keyword" clearable placeholder="检索" @keyup.enter="getDeviceList" />
      <div class="tool">
        <el-tooltip class="tip" content="采集设备" effect="dark" placement="top">
          <vab-icon class="icon" icon="settings-3-line" style="margin: 0 4px" @click="skip" />
        </el-tooltip>
        <el-tooltip class="tip" content="刷新" effect="dark" placement="top">
          <vab-icon class="icon" icon="refresh-line" @click="getDeviceList" />
        </el-tooltip>
      </div>
    </div>
    <div v-loading="listLoading" class="content">
      <div
        v-for="item in listDate"
        :key="item.id"
        class="item"
        :class="{ active: item.id === curId }"
        @click="handleClick(item.id)"
      >
        <span>{{ item.name }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .active {
    background-color: #567ef1;
    color: #fff;
  }
  .device-list {
    padding-right: 20px;

    width: 200px;
    overflow: hidden;
    .top {
      margin-top: -4px;
      width: 100%;
      height: 40px;
      display: flex;
      align-items: center;
      .tool {
        width: 52px;
        color: #7a649c;

        .icon {
          &:hover {
            cursor: pointer;
          }
        }
      }
    }
    .content {
      background-color: #fff;
      border-right: 1px solid var(--el-border-color);
      width: 100%;
      height: calc(100vh - 170px);
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 3px;
      }

      .item {
        height: 30px;
        padding-left: 20px;
        line-height: 30px;
        border-radius: 2px;
        &:hover {
          cursor: pointer;
        }
      }
    }
  }
</style>
