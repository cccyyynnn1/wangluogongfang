<script setup lang="ts">
  import { overviewLoadBalanceListApi } from '~/src/api-ecs/equipment'

  const curId = ref('')

  const emit = defineEmits<{
    (e: 'on-change-menu', curId: string): void
  }>()

  const activeName = ref('')

  let listDate = reactive<{ id: string; name: string; loadBalanceViews: { id: string; hostName: string }[] }[]>([]) // 表格数据

  const queryPage = reactive({
    page: 1,
    limit: 10,
    moduleType: 'loadBalance',
    apiType: 'overviewloadbalance',
    query: {
      keyword: undefined,
    },
  })

  const listLoading = ref(false) // 是否加载

  // 得到设备列表
  const getBalanceList = async () => {
    try {
      listLoading.value = true
      // @ts-ignore
      const { data } = await overviewLoadBalanceListApi({ ...queryPage, query: JSON.stringify(queryPage.query) })
      listDate = data.list
      activeName.value = listDate[0].id
    } finally {
      listLoading.value = false
    }
  }

  // 点击某一项
  const handleClick = (id: string) => {
    curId.value = id
    emit('on-change-menu', curId.value)
  }

  onMounted(() => {
    getBalanceList()
  })
</script>

<script lang="ts">
  export default {
    name: 'LoadBalancingList',
  }
</script>
<template>
  <div class="load-balancing-list">
    <div class="top">
      <el-input v-model="queryPage.query.keyword" clearable placeholder="查询" @keyup.enter="getBalanceList" />
    </div>
    <div v-loading="listLoading" class="content">
      <el-collapse v-model="activeName" accordion>
        <el-collapse-item v-for="item in listDate" :key="item.id" :name="item.id" :title="item.name">
          <div
            v-for="td in item.loadBalanceViews"
            :key="td.id"
            class="item"
            :class="{ active: td.id === curId }"
            @click="handleClick(td.id)"
          >
            <span>{{ td.hostName }}</span>
          </div>
        </el-collapse-item>
      </el-collapse>
    </div>
  </div>
</template>

<style scoped lang="scss">
  :deep() {
    .el-collapse-item,
    .el-collapse-item__header {
      height: 30px !important;
    }
    .el-collapse-item__content {
      padding-bottom: 0;
    }
    .el-collapse {
      --el-collapse-border-color: none;
    }
  }
  .active {
    background-color: #567ef1;
    color: #fff;
  }
  .load-balancing-list {
    padding-right: 20px;

    width: 200px;
    overflow: hidden;
    .top {
      width: 100%;
      height: 40px;
    }
    .content {
      border-right: 1px solid var(--el-border-color);
      width: 100%;
      height: calc(100vh - 170px);
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 3px;
      }

      .item {
        height: 30px;
        padding-left: 8px;
        line-height: 30px;
        border-radius: 2px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        word-break: break-all;
        word-wrap: break-word;
        &:hover {
          cursor: pointer;
        }
      }
    }
  }
</style>
