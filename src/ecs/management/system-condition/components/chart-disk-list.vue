<script lang="ts">
  export default {
    name: 'ChartDiskList',
  }
</script>

<script setup lang="ts">
  const props = defineProps<{
    total: any
  }>()
</script>

<template>
  <div class="disks">
    <el-row :gutter="10" style="margin-inline: 0">
      <el-col :span="6" v-for="item in total" :key="item.filesystem">
        <div class="disk">
          <div class="rate" :style="{ '--rate': `${item.usePercent}%` }"></div>
          <p>
            使用率：
            <span>{{ item.usePercent.toFixed(2) }}%</span>
          </p>
          <p>
            挂载点：
            <span>{{ item.mountPoint }}</span>
          </p>
          <div class="use">{{ item.used }}&nbsp;&nbsp;可用({{ item.size }})</div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped lang="scss">
  .disks {
    height: 600px;
    overflow: auto;
    .disk {
      border-radius: 10px;
      border: 1px solid #ebe9fa;
      margin-bottom: 20px;
      overflow: hidden;
      padding: 20px;
      padding-bottom: 0;
      p {
        font-weight: 400;
        font-size: 14px;
        color: #352e58;
        margin-block: 8px;
        span {
          display: inline-block;
          width: calc(100% - 62px);
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          text-align: right;
          vertical-align: bottom;
        }
      }
      .rate {
        height: 22px;
        border-radius: 4px;
        background: #f4f3fc;
        overflow: hidden;
        position: relative;
        margin-bottom: 7px;
        &::after {
          content: ' ';
          display: block;
          position: absolute;
          height: 22px;
          // inset: 0;
          // right: var(--rate);
          width: var(--rate);
          border-radius: 4px;
          background-color: var(--el-color-primary);
        }
      }
      .use {
        height: 42px;
        border-top: 1px solid #ebe9fa;
        text-align: center;
        font-weight: 400;
        font-size: 14px;
        margin-inline: -20px;
        color: #9d9baa;
        margin-top: 22px;
        line-height: 42px;
      }
    }
  }
</style>
