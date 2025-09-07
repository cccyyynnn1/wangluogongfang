<script setup lang="ts">
  import { Search } from '@element-plus/icons-vue'
  import VabChart from '@/plugins/VabChart/index.vue'
  import { useTableCopy } from '@/utils'

  enum RiskStatus {
    ALL = 0,
    LOW,
    MID,
    HEIGHT,
    CRITICAL,
  }
  const gridData = [
    {
      date: '2016-05-02',
      name: 'John Smith',
      address: 'No.1518,  Jinshajiang Road, Putuo District',
      level: 2,
    },
    {
      date: '2016-05-04',
      name: 'John Smith',
      address: 'No.1518,  Jinshajiang Road, Putuo District',
      level: 1,
    },
    {
      date: '2016-05-01',
      name: 'John Smith',
      address: 'No.1518,  Jinshajiang Road, Putuo District',
      level: 4,
    },
    {
      date: '2016-05-03',
      name: 'John Smith',
      address: 'No.1518,  Jinshajiang Road, Putuo District',
      level: 3,
    },
  ]
  const levelKey: {
    [key: number]: {
      label: string
      class: string
    }
  } = {
    [RiskStatus.ALL]: {
      label: '全部',
      class: '',
    },
    [RiskStatus.LOW]: {
      label: '低危',
      class: 'low',
    },
    [RiskStatus.MID]: {
      label: '中危',
      class: 'mid',
    },
    [RiskStatus.HEIGHT]: {
      label: '高危',
      class: 'high',
    },
    [RiskStatus.CRITICAL]: {
      label: '危急',
      class: 'critical',
    },
  }
  const alertLevels = [
    {
      valeue: 4,
      label: '危急',
    },
    {
      valeue: 3,
      label: '高危',
    },
    {
      valeue: 2,
      label: '中危',
    },
    {
      valeue: 1,
      label: '低危',
    },
    {
      valeue: 0,
      label: '正常',
    },
  ]

  const dialogTableVisible = ref(false)

  const sourceIp = ref('')
  const sourceIpLevel = ref(4)
  const destinationsIp = ref('')
  const destinationsIpLevel = ref(4)
  const traceability_chart = ref<InstanceType<typeof VabChart>>()

  const obj = { 低危: 1, 中危: 2, 高危: 3, 危急: 4 }
  function getLevel(str: string) {
    // @ts-ignore
    const level = obj[str]
    return levelKey[level]
  }
</script>

<script lang="ts">
  export default {
    name: 'ExamplesTraceability',
  }
</script>

<template>
  <div class="traceability">
    <div class="hostChart">
      <vab-chart ref="traceability_chart" class="traceability-echart" theme="vab-echarts-theme" />
    </div>
    <div class="hostTable">
      <div class="hostTable_item">
        <div class="hostTable_title">源IP</div>
        <div class="hostTable_search">
          <el-input v-model="sourceIp" />
          <el-button :icon="Search" type="primary">检索</el-button>
        </div>
        <div class="alertLevel">
          <el-button-group>
            <el-button
              v-for="level in alertLevels"
              :key="level.valeue"
              :type="sourceIpLevel === level.valeue ? 'primary' : 'default'"
              @click="sourceIpLevel = level.valeue"
            >
              {{ level.label }}
            </el-button>
          </el-button-group>
        </div>
        <el-table :data="gridData" height="650px" style="width: 100%" @cell-contextmenu="useTableCopy">
          <el-table-column label="标签" prop="level" width="80">
            <template #default="{ row }">
              <el-tag :class="[getLevel(row.level).class]">
                {{ getLevel(row.level).label }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="时间" prop="date" width="140" />
          <el-table-column label="源IP" prop="address" width="140" />
          <el-table-column label="目的地IP" prop="address" width="140" />
          <el-table-column label="Host" prop="address" width="140" />
          <el-table-column label="XFF" prop="address" width="80" />
          <el-table-column label="URL" prop="address" width="140" />
        </el-table>
      </div>
      <div class="host">
        {{ '172.120.208.17' }}
      </div>
      <div class="hostTable_item">
        <div class="hostTable_title">目的地IP</div>
        <div class="hostTable_search">
          <el-input v-model="destinationsIp" />
          <el-button :icon="Search" type="primary">检索</el-button>
        </div>
        <div class="alertLevel">
          <el-button-group>
            <el-button
              v-for="level in alertLevels"
              :key="level.valeue"
              :type="destinationsIpLevel === level.valeue ? 'primary' : 'default'"
              @click="destinationsIpLevel = level.valeue"
            >
              {{ level.label }}
            </el-button>
          </el-button-group>
        </div>
        <el-table :data="gridData" height="650px" style="width: 100%" @cell-contextmenu="useTableCopy">
          <el-table-column label="标签" prop="level" width="80">
            <template #default="{ row }">
              <el-tag :class="[getLevel(row.level).class]">
                {{ getLevel(row.level).label }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="时间" prop="date" width="140" />
          <el-table-column label="源IP" prop="address" width="140" />
          <el-table-column label="目的地IP" prop="address" width="140" />
          <el-table-column label="Host" prop="address" width="140" />
          <el-table-column label="XFF" prop="address" width="80" />
          <el-table-column label="URL" prop="address" width="140" />
        </el-table>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .hostChart {
    width: 100%;
    height: 420px;
    border-radius: 4px;
    margin-bottom: 20px;
    background: #f3f9ff;
  }

  .hostTable {
    display: flex;

    .hostTable_item {
      flex: 1;
      overflow-x: auto;

      .hostTable_title {
        height: 14px;
        line-height: 12px;
        position: relative;
        text-indent: 0.8em;

        &::before {
          content: ' ';
          display: inline-block;
          position: absolute;
          width: 3px;
          height: 100%;
          left: 0;
          background: #0d88fe;
          margin-right: 5px;
        }
      }

      .hostTable_search {
        display: flex;
        justify-content: space-between;
        margin: 16px 0;

        .el-input {
          width: 80%;
        }
      }

      .alertLevel {
        margin-bottom: 20px;
      }
    }

    .host {
      width: 114px;
      height: 48px;
      margin: 20% 20px;
      background: #ecf4ff;
      border-radius: 4px;
      border: 1px solid #abcbff;
      line-height: 48px;
      text-align: center;
    }
  }

  .traceability-echart {
    width: 100%;
  }

  :deep() {
    .echarts {
      height: 100%;
    }
  }
</style>
