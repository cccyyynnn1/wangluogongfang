<script setup lang="ts">
  import ThreatTrends from './components/workbench/threat-trends.vue'

  import RankingOfAsset from './components/workbench/ranking-of-asset.vue'

  import SystemDistribution from './components/workbench/system-distribution.vue'

  import AttackSource from './components/workbench/AttackSource.vue'

  import AttackType from './components/workbench/AttackType.vue'

  import AccessStatistics from './components/workbench/accessStatistics.vue'

  import Message from './components/workbench/message.vue'

  import { homeCountApi } from '~/src/api-ecs/dashboard'

  const router = useRouter()

  const data = [
    {
      name: '站点资产',
      redirect: '/site_index',
      icon: 'site-assets',
      desc: '跳转到资产栏目的站点资产页',
      color: '#1690FF',
    },

    {
      name: '路径追踪',
      redirect: '/retrieve,pathTracing',
      icon: 'path-tracing',
      desc: '跳转到检索栏目的路径追踪页',
      color: '#FFC87C',
    },

    {
      name: '攻击透视',
      redirect: '/alerts,isAttackPerspective',
      icon: 'attack-perspective',
      desc: '跳转到告警栏目的攻击透视页',
      color: '#6EDED6',
    },

    {
      name: '应用数据检索',
      redirect: '/retrieve,applicationLayer',
      icon: 'search',
      desc: '跳转到检索栏目的应用层会话',
      color: '#C39AEF',
    },

    {
      name: '网络数据检索',
      redirect: '/retrieve/network',
      icon: 'web-data-retrieval',
      desc: '跳转到检索栏目的网络层会话',
      color: '#FF9ACB',
    },
    {
      name: '资产访问',
      redirect: '/assets/asset-visits',
      icon: 'asset-access',
      desc: '跳转到资产栏目内的资产访问页',
      color: '#95DE64',
    },
    {
      name: '未知资产',
      redirect: '/assets/index,unKnown',
      icon: 'unknow-assets',
      desc: '跳转到资产栏目内的未知资产页',
      color: '#6BC1FF',
    },

    {
      name: '用户管理',
      redirect: '/managements/user_management/user',
      icon: 'user',
      desc: '跳转到系统栏目的用户管理',
      color: '#FED566',
    },
  ]

  // const line = ref() // 攻击趋势
  const threatLevel = ref() // 攻击趋势

  const attackIp = ref() // 攻击源

  const threatType = ref() // 攻击类型

  // const asset = ref() // 资产数据量

  const clientIp = ref() // 访问统计

  const victimIp = ref() // 资产威胁统计排行

  const sysType = ref() // 系统分布

  const initData = async () => {
    const { data } = await homeCountApi()
    // line.value = data.line
    threatLevel.value = data.threatLevel
    attackIp.value = data.attackIp
    threatType.value = data.threatType
    // asset.value = data.asset
    clientIp.value = data.clientIp
    victimIp.value = data.victimIp
    sysType.value = data.sysType
  }

  onMounted(() => {
    initData()
  })

  const skipToPage = (redirect: string) => {
    const arr = redirect.split(',')
    router.push({ path: arr[0], query: { params: arr[1] } })
  }

  provide('tiemFlag', null)
</script>

<script lang="ts">
  export default {
    name: 'Workbench',
  }
</script>
<template>
  <div class="workbench">
    <el-row :gutter="10">
      <el-col :span="18">
        <div class="navigation">
          <vab-card>
            <template #header>
              <span class="title">快捷导航</span>
            </template>
            <div v-for="(item, index) in data" :key="index" class="each-item" @click="() => skipToPage(item.redirect)">
              <el-image :src="require(`@/assets/message_image/${item.icon}.svg`)" />
              <span style="margin-top: 14px; color: #837e9e; font-weight: 700">{{ item.name }}</span>
            </div>
          </vab-card>
        </div>
        <el-row class="workbench_chart" :gutter="10">
          <el-col :span="12">
            <!-- <AttackTrend :line="line" /> -->
            <vab-card>
              <template #header>
                <span class="title">威胁趋势</span>
              </template>
              <ThreatTrends :threat-levels="threatLevel" />
            </vab-card>
          </el-col>
          <el-col :span="12">
            <vab-card>
              <template #header>
                <span class="title">资产威胁统计排行</span>
              </template>
              <RankingOfAsset :show-btn="false" :victim-ips="victimIp" />
            </vab-card>
          </el-col>
        </el-row>
        <el-row class="workbench_chart" :gutter="10">
          <el-col :span="8">
            <SystemDistribution :sys-types="sysType" />
          </el-col>
          <el-col :span="8">
            <AttackSource :attack-ip="attackIp" />
          </el-col>
          <el-col :span="8">
            <AccessStatistics :client-ip="clientIp" />
          </el-col>
        </el-row>
      </el-col>
      <el-col :span="6">
        <el-row class="workbench_chart" :gutter="12">
          <el-col :span="24">
            <Message />
          </el-col>
          <el-col :span="24">
            <AttackType :threat-type="threatType" />
          </el-col>
        </el-row>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped lang="scss">
  $defaultHeight: 245px;
  .workbench {
    padding: 10px;
    // background-color: #f6f8f9;
    .title {
      color: #303133;
      font-weight: 700;
    }
    .navigation {
      width: 100%;

      :deep() {
        .el-card {
          border-radius: 6px;
          height: 178px;
        }
        .el-card__header {
          // padding-top: 10px;
          border-bottom: none;
          height: 36px;
          // display: flex;
          // align-items: center;
        }
        .el-card__body {
          height: 130px;
          display: flex;
        }
      }
      .title {
        color: #303133;
        font-size: 15px;
        font-weight: 500;
      }
      .each-item {
        width: calc((100% - 84px) / 8);
        display: flex;
        align-items: center;
        justify-content: center;
        flex-direction: column;
        cursor: pointer;
        :deep(.el-image) {
          width: 60px;
          height: 60px;
          background: #efedff;
          border-radius: 22px;
        }
        &:not(:last-child) {
          margin-right: 12px;
        }
      }
    }
    .workbench_chart {
      :deep() {
        .ranking-of-asset,
        .system-distribution {
          height: calc((100vh - $defaultHeight) / 2 - 36px);
          padding-top: 0 !important;
          .echarts {
            height: calc(100% - 32px) !important;
          }
        }
        .system-distribution {
          .echarts {
            height: 100% !important;
            padding: 0 10px 10px;
          }
        }
        .threat-trends {
          height: calc((100vh - $defaultHeight) / 2 - 36px);
        }
        .title {
          color: #303133;
          font-size: 15px;
          font-weight: 500;
        }
        .el-card {
          height: calc((100vh - $defaultHeight) / 2);
          border-radius: 6px;
          // border: none;

          .el-card__header {
            border-bottom: none;
            height: 36px;
            // display: flex;
            // align-items: center;
            // justify-content: space-between;
          }
          .right {
            width: 90px;
            float: right;
          }
          .el-card__body {
            padding: 0;
            height: calc((100vh - $defaultHeight) / 2 - 36px);
            .echart-content {
              height: calc((100vh - $defaultHeight) / 2 - 36px);
              width: 100%;
            }
            .table,
            .data-node-state {
              height: 100%;
            }
          }
        }
        .message {
          height: calc((100vh - $defaultHeight) / 2 + 190px);
          margin-bottom: 12px;
          .el-card__body {
            padding: 0;
            height: calc((100vh - $defaultHeight) / 2 + 190px - 36px);
          }
        }
      }
    }
  }
</style>
