<script lang="ts">
  export default {
    name: 'TrafficModel',
  }
</script>
<script setup lang="ts">
  import { getRuleApi, deleteTrafficApi } from '@/api-ecs/retrieve'
  import { RuleModelItem } from '@/types'
  const $baseMessage: any = inject('$baseMessage')
  const $baseConfirm: any = inject('$baseConfirm')
  const router = useRouter()

  // 表格数据
  const layout = ref('total, sizes, prev, pager, next, jumper')

  const queryPage = reactive({
    total: 0,
    pageNum: 1,
    pageSize: 30,
    title: '',
    loading: false,
  })

  const deleteQuery = reactive({
    visible: false,
    top: '0',
    left: '0',
  })

  const deleteData = ref()

  // 页面数据
  const listData = ref<RuleModelItem[]>([])

  // 页容量改变
  const handleSizeChange = (val: number) => {
    queryPage.pageSize = val
    getRuleList()
  }
  const trafficDeleteHandle = (data: any, event: any) => {
    const { clientX, clientY } = event
    deleteQuery.left = `${clientX + 20}px`
    deleteQuery.top = `${clientY}px`
    deleteQuery.visible = true
    deleteData.value = data
  }
  const deleteTraffic = async () => {
    $baseConfirm('你确定要删除当前项吗', null, async () => {
      const { msg } = await deleteTrafficApi({ id: deleteData.value.id })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      deleteQuery.visible = false
      getRuleList()
    })
  }

  // 页面改变
  const handleCurrentChange = () => {
    getRuleList()
  }

  const showDetail = (ruleId: number) => {
    router.push({ name: 'TrafficSceneModelDetail', params: { id: ruleId } })
  }
  async function getRuleList() {
    queryPage.loading = true
    const { pageNum, pageSize } = queryPage
    const {
      data: { total, records },
    } = await getRuleApi({ pageNum, pageSize })
    listData.value = records
    queryPage.total = total
    queryPage.loading = false
  }
  useEventListener(document, 'click', (evt) => {
    deleteQuery.visible = false
  })
  onMounted(() => {
    getRuleList()
  })
</script>

<template>
  <div v-loading="queryPage.loading" class="traffic-model-container">
    <!-- <el-tabs :model-value="'model'">
      <el-tab-pane label="流量场景模型" name="model"> -->
    <div class="content">
      <el-row :gutter="20">
        <el-col v-for="(item, index) in listData" :key="index" :span="4">
          <div
            class="item"
            @click="showDetail(item.id)"
            @contextmenu.prevent="($event) => trafficDeleteHandle(item, $event)"
          >
            <div class="warp">
              <el-image
                lazy
                :src="require('@/assets/theme_images/model-icon.png')"
                style="height: 40px; width: 40px; margin: 0 14px 0 20px"
              />
              {{ item.ruleName }}
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
    <!-- pagination -->
    <el-pagination
      v-model:current-page="queryPage.pageNum"
      background
      :layout="layout"
      :page-size="queryPage.pageSize"
      :page-sizes="[10, 20, 50, 100]"
      style="margin-bottom: 20px"
      :total="queryPage.total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />
    <!-- </el-tab-pane>
    </el-tabs> -->
    <div v-show="deleteQuery.visible" class="contextmenu" :style="{ left: deleteQuery.left, top: deleteQuery.top }">
      <el-button type="primary" @click="deleteTraffic()">删除</el-button>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .traffic-model-container {
    .content {
      width: 100%;
      min-height: 400px;
      margin-bottom: 20px;
    }

    .item {
      height: 120px;
      border-radius: 5px;
      background: url('~@/assets/theme_images/traffic-model.png');
      background-size: 100%;
      margin-top: 20px;
      border: 1px solid #f1f9ff;
      border-radius: 16px;
      display: flex;
      justify-content: center;
      align-items: center;
      flex-direction: column;
      height: 0;
      padding-bottom: 45%;
      position: relative;
      .warp {
        height: 60px;
        width: 100%;
        position: absolute;
        left: 50%;
        top: 50%;
        transform: translate(-50%, -50%);
        display: flex;
        // justify-content: center;
        align-items: center;
        // flex-direction: column;
      }

      &:hover {
        cursor: pointer;
        box-shadow: 0px 0px 8px 1px rgba(0, 0, 0, 0.08);
      }
    }
  }
  .contextmenu {
    position: fixed;
    top: 0;
    left: 0;
    margin: 0;
    background: #fff;
    z-index: 3000;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 400;
    color: #333;
    box-shadow: 2px 2px 3px 0 rgba(0, 0, 0, 0.3);
  }
</style>
