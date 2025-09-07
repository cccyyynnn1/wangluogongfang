<script setup lang="ts">
  import AlgorithmsExamples from './components/algorithms-examples.vue'
  import AlgorithmsAlerts from './components/algorithms-alerts.vue'
  import AlgorithmsImages from './components/algorithms-images.vue'
  import { Algorithms_alerts, AlgorithmsType } from './type'
  import { timeDuratioOptions } from '@/ecs/alert/data/index'
  const activeName = ref<AlgorithmsType>('examples')
  // 高级筛选
  const isAdvanced = ref(false)
  // 自动刷新
  const autoRefresh = ref(false)
  const autoRefreshDuration = ref(30)
  provide(Algorithms_alerts, {
    isAdvanced: isAdvanced,
    autoRefresh: autoRefresh,
    autoRefreshDuration: autoRefreshDuration,
  })
</script>

<script lang="ts">
  export default {
    name: 'Algorithms',
  }
</script>

<template>
  <div class="algorithms-container">
    <el-tabs v-model="activeName">
      <el-tab-pane label="算法实例" name="examples">
        <algorithms-examples />
      </el-tab-pane>
      <el-tab-pane label="算法镜像" name="image">
        <algorithms-images />
      </el-tab-pane>
      <!-- <el-tab-pane label="算法告警" name="alert">
        <algorithms-alerts />
      </el-tab-pane> -->
    </el-tabs>
    <!-- 日志查询方式 -->
    <div v-if="activeName === 'alert'" class="table_action">
      <el-select v-model="autoRefreshDuration" class="m-2" :disabled="autoRefresh">
        <el-option v-for="item in timeDuratioOptions" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
      <div class="action">
        <label>自动刷新：</label>
        <el-switch v-model="autoRefresh" />
      </div>
      <div class="action">
        <el-button-group>
          <el-button :type="!isAdvanced ? 'primary' : 'default'" @click="isAdvanced = false">快速筛选</el-button>
          <el-button :type="isAdvanced ? 'primary' : 'default'" @click="isAdvanced = true">高级筛选</el-button>
        </el-button-group>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .algorithms-container {
    height: calc(100vh - 35px);
    position: relative;
    overflow-y: auto;
    &::-webkit-scrollbar {
      width: 0;
      height: 0;
    }
    .table_action {
      display: flex;
      position: absolute;
      align-items: center;
      top: 15px;
      right: 20px;

      .action {
        display: flex;
        align-items: center;
        margin-left: 30px;
      }
    }
  }
</style>
