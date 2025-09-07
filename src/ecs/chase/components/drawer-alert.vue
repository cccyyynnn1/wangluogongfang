<script lang="ts">
  export default {
    name: 'DrawerAlert', // 详情
  }
</script>

<script setup lang="ts">
  // const props = defineProps<{
  //   alldata: any
  // }>()
  const emits = defineEmits<{
    (e: 'on-close'): void
  }>()
</script>

<template>
  <div class="drawer-info sketchpad-drawer">
    <div class="sketchpad-drawer-top">
      <div>
        <span style="font-size: 18px; color: #1e1842; margin-right: 6px">告警列表</span>
        <span style="font-size: 14px; color: #4d3fa2">共计23条</span>
      </div>
      <el-icon @click="emits('on-close')"><Close /></el-icon>
    </div>
    <div class="row_1">
      <div class="row_1_inside">
        <img alt="" :src="require('@/assets/chase/aggressors.svg')" style="width: 16px; height: 16px" />
        <span style="color: #303133; font-weight: 500">172.168.12.203</span>
        <span style="font-weight: 500; font-size: 12px; color: #303133">内网</span>
      </div>
      <div class="row_1_arrow"></div>
      <div class="row_1_outside">
        <img alt="" :src="require('@/assets/chase/occupied.svg')" style="width: 16px; height: 16px" />
        <span style="color: #303133; font-weight: 500">172.168.12.203</span>
        <span style="font-weight: 500; font-size: 12px; color: #303133">内网</span>
      </div>
    </div>
    <div class="row_wrap">
      <div v-for="item in 7" :key="item" class="row_item">
        <div class="row_item_top">
          <div class="row_item_top_text">FatPipeg 任意文件上传</div>
        </div>
        <div class="row_item_bottom">
          <div style="text-align: center">
            <div class="info">23条</div>
            <div class="tip">告警数量</div>
          </div>
          <div style="text-align: center">
            <div class="info">404</div>
            <div class="tip">状态码</div>
          </div>
          <div style="text-align: center">
            <div class="info">内对外</div>
            <div class="tip">外网</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .row_wrap {
    width: 100%;
    height: calc(100% - 125px);
    overflow-y: auto;
    margin-top: 10px;
    .row_item_top {
      height: 34px;
      display: flex;
      align-items: center;
      .row_item_top_text {
        margin-left: 30px;
        font-weight: 500;
        font-size: 14px;
        color: #303133;
        position: relative;
        &::after {
          position: absolute;
          content: '';
          width: 12px;
          height: 12px;
          border-radius: 50%;
          background: #ff6767;
          border: 1px solid #ffffff;
          top: 5px;
          left: -15px;
        }
      }
    }
    .row_item_bottom {
      width: 100%;
      height: 66px;
      background: #fbfbff;
      border-radius: 0px 0px 8px 8px;
      border-top: 1px solid #e0ddf6;
      padding: 20px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      .info {
        font-weight: 500;
        font-size: 14px;
        color: #2b2742;
      }
      .tip {
        font-weight: 400;
        font-size: 13px;
        color: #918da5;
      }
    }
    &::-webkit-scrollbar {
      width: 0;
      height: 0;
    }
    .row_item {
      margin-bottom: 10px;
      width: 100%;
      height: 102px;
      background: #ffffff;
      border-radius: 8px;
      border: 1px solid #e0ddf6;
    }
  }
  .row_1 {
    height: 76px;
    width: 100%;
    margin-top: 10px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .row_1_inside,
    .row_1_outside {
      padding: 6px 0;
      width: 140px;
      height: 76px;
      background: #fff4f4;
      border-radius: 8px;
      border: 1px solid #ff6c6c;
      display: flex;
      align-items: center;
      justify-content: space-around;
      flex-direction: column;
    }
    .row_1_outside {
      border: 1px solid #ccc8f8;
      background: #fbfbff;
    }
    .row_1_arrow {
      width: 15px;
      height: 0px;
      border: 1px solid #8781b3;
      position: relative;
      margin-left: -3px;
      &::after {
        position: absolute;
        content: '';
        top: -6px;
        left: 12px;
        width: 0;
        height: 0;
        border-width: 6px;
        border-style: solid;
        border-color: transparent transparent transparent #8781b3;
      }
    }
  }
</style>
