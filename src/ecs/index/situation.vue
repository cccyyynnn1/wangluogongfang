<script setup lang="ts">
  import { homeLineApi } from '~/src/api-ecs/dashboard'
  import AttackLine from './components/situation/attackLine.vue'
  import AbnormalWarn from './components/situation/abnormalWarn.vue'
  import Index from './components/situation/index.vue'
  import AttackIp from './components/situation/attackIp.vue'
  import ThreatType from './components/situation/threatType.vue'
  import WarnRequest from './components/situation/warnRequest.vue'
  import numberFormatter from '~/src/utils/number'

  const allData = reactive({
    index: undefined, // 应用数据量趋势
    attackLine: undefined, // 攻击趋势
    abnormalWarn: undefined, // 异常访问趋势
    warnRequest: undefined, // 异常访问占比
    attackIp: undefined, // 攻击源
    topCount: undefined, // 顶部数据
    threatType: undefined, // 攻击类型
  })

  const topCount = reactive({
    apiCount: {
      count: undefined,
      tradeNum: undefined,
    }, // 站点api总数
    warn_abonormal_count: {
      count: undefined,
      tradeNum: undefined,
    }, // 异常请求量
    assetCount: {
      count: undefined,
      tradeNum: undefined,
    }, // 资产总数
    warn_attack_count: {
      count: undefined,
      tradeNum: undefined,
    }, // 攻击总量
    siteCount: {
      count: undefined,
      tradeNum: undefined,
    }, // 站点总数
    assetNotKnowCount: {
      count: undefined,
      tradeNum: undefined,
    }, // 未知资产总数
  })

  const initData = async () => {
    const { data } = await homeLineApi()
    for (const key in allData) {
      allData[key as keyof typeof allData] = data[key]
    }
    if (data.topCount) {
      for (const key in topCount) {
        topCount[key as keyof typeof topCount] = data.topCount[key]
      }
    }
  }

  onMounted(() => {
    initData()
  })
</script>

<script lang="ts">
  export default {
    name: 'Situation', // 态势
  }
</script>
<template>
  <div class="situation">
    <div class="navigation">
      <div class="each-item">
        <el-image class="img" :src="require('@/assets/situation/assent.png')" />
        <div class="show-data">
          <div class="show-data-top">资产总数</div>
          <div class="show-data-bottom">
            {{ (topCount.assetCount && numberFormatter.format(Number(topCount.assetCount.count))) || '' }}
            <el-image
              v-if="topCount.assetCount && topCount.assetCount.tradeNum == 1"
              class="small-img"
              :src="require('@/assets/situation/rise.png')"
            />
            <el-image
              v-if="topCount.assetCount && topCount.assetCount.tradeNum == -1"
              class="small-img"
              :src="require('@/assets/situation/decline.png')"
            />
          </div>
        </div>
      </div>
      <div class="each-item">
        <el-image class="img" :src="require('@/assets/situation/unknown.png')" />
        <div class="show-data">
          <div class="show-data-top">未知资产总数</div>
          <div class="show-data-bottom">
            {{ (topCount.assetNotKnowCount && numberFormatter.format(Number(topCount.assetNotKnowCount.count))) || '' }}
            <el-image
              v-if="topCount.assetNotKnowCount && topCount.assetNotKnowCount.tradeNum == 1"
              class="small-img"
              :src="require('@/assets/situation/rise.png')"
            />
            <el-image
              v-if="topCount.assetNotKnowCount && topCount.assetNotKnowCount.tradeNum == -1"
              class="small-img"
              :src="require('@/assets/situation/decline.png')"
            />
          </div>
        </div>
      </div>
      <div class="each-item">
        <el-image class="img" :src="require('@/assets/situation/konwn.png')" />
        <div class="show-data">
          <div class="show-data-top">站点总数</div>
          <div class="show-data-bottom">
            {{ (topCount.siteCount && numberFormatter.format(Number(topCount.siteCount.count))) || '' }}
            <el-image
              v-if="topCount.siteCount && topCount.siteCount.tradeNum == 1"
              class="small-img"
              :src="require('@/assets/situation/rise.png')"
            />
            <el-image
              v-if="topCount.siteCount && topCount.siteCount.tradeNum == -1"
              class="small-img"
              :src="require('@/assets/situation/decline.png')"
            />
          </div>
        </div>
      </div>
      <div class="each-item">
        <el-image class="img" :src="require('@/assets/situation/sum.png')" />
        <div class="show-data">
          <div class="show-data-top">站点API总数</div>
          <div class="show-data-bottom">
            {{ (topCount.apiCount && numberFormatter.format(Number(topCount.apiCount.count))) || '' }}
            <el-image
              v-if="topCount.apiCount && topCount.apiCount.tradeNum == 1"
              class="small-img"
              :src="require('@/assets/situation/rise.png')"
            />
            <el-image
              v-if="topCount.apiCount && topCount.apiCount.tradeNum == -1"
              class="small-img"
              :src="require('@/assets/situation/decline.png')"
            />
          </div>
        </div>
      </div>
      <div class="each-item">
        <el-image class="img" :src="require('@/assets/situation/attact.png')" />
        <div class="show-data">
          <div class="show-data-top">攻击总量(12h)</div>
          <div class="show-data-bottom">
            {{ (topCount.warn_attack_count && numberFormatter.format(Number(topCount.warn_attack_count.count))) || '' }}
            <el-image
              v-if="topCount.warn_attack_count && topCount.warn_attack_count.tradeNum == 1"
              class="small-img"
              :src="require('@/assets/situation/rise.png')"
            />
            <el-image
              v-if="topCount.warn_attack_count && topCount.warn_attack_count.tradeNum == -1"
              class="small-img"
              :src="require('@/assets/situation/decline.png')"
            />
          </div>
        </div>
      </div>
      <div class="each-item">
        <el-image class="img" :src="require('@/assets/situation/unusual.png')" />
        <div class="show-data">
          <div class="show-data-top">异常请求量(12h)</div>
          <div class="show-data-bottom">
            {{
              (topCount.warn_abonormal_count && numberFormatter.format(Number(topCount.warn_abonormal_count.count))) ||
              ''
            }}
            <el-image
              v-if="topCount.warn_abonormal_count && topCount.warn_abonormal_count.tradeNum == 1"
              class="small-img"
              :src="require('@/assets/situation/rise.png')"
            />
            <el-image
              v-if="topCount.warn_abonormal_count && topCount.warn_abonormal_count.tradeNum == -1"
              class="small-img"
              :src="require('@/assets/situation/decline.png')"
            />
          </div>
        </div>
      </div>
    </div>
    <el-row :gutter="12">
      <el-col :span="8">
        <Index :index="allData.index" />
      </el-col>
      <el-col :span="8">
        <AttackLine :attack-line="allData.attackLine" />
      </el-col>
      <el-col :span="8">
        <AbnormalWarn :abnormal-warn="allData.abnormalWarn" />
      </el-col>
    </el-row>
    <el-row :gutter="12">
      <el-col :span="8">
        <AttackIp :attack-ip="allData.attackIp" />
      </el-col>
      <el-col :span="8">
        <ThreatType :threat-type="allData.threatType" />
      </el-col>
      <el-col :span="8">
        <WarnRequest :warn-request="allData.warnRequest" />
      </el-col>
    </el-row>
  </div>
</template>

<style scoped lang="scss">
  .situation {
    background-color: #f6f8f9;
    .el-card {
      min-height: calc((100vh - 300px) / 2);
      border-radius: 6px;
    }
  }

  :deep() {
    .el-card {
      border: none;
      overflow: inherit;
    }

    .el-card__header {
      border-bottom: none;
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: space-between;
    }
  }
  :deep(.el-card__body) {
    min-height: calc((100vh - 310px) / 2 - 40px);
    padding: 0 !important;
  }
  :deep() {
    .title {
      position: relative;
      font-weight: 700;
      &::before {
        display: inline-block;
        width: 3px;
        height: 12px;
        vertical-align: bottom;
        margin: 0 6px 2px 0px;
        content: '';
        background: #1890ff;
      }
    }
    .echart-content {
      height: calc((100vh - 310px) / 2 - 40px);
      width: 100%;
    }
  }
  .navigation {
    margin-bottom: 12px;
    width: 100%;
    height: 152px;
    background: #ffffff;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: space-around;
    .each-item {
      display: flex;
      align-items: center;
      justify-content: space-around;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      word-break: break-all;
      word-wrap: break-word;
      .img {
        height: 64px;
        width: 64px;
      }
      .show-data {
        margin-left: 14px;
        font-size: 13px;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        .small-img {
          width: 16px;
          height: 10px;
        }
        .show-data-top {
          height: 18px;
          color: #a9acb3;
          line-height: 18px;
        }
        .show-data-bottom {
          margin-top: 10px;
          height: 24px;
          font-size: 18px;
          font-weight: 600;
          color: #303159;
          line-height: 24px;
        }
      }
    }
  }
</style>
