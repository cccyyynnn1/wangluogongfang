<script lang="ts">
  export default {
    name: 'ToolUri',
  }
</script>

<script setup lang="ts">
  import { Delete } from '@element-plus/icons-vue'
  import { URLDecode } from '@/assets/cyberChef/node/index.mjs'
  const leftVal = ref('')
  const rightVal = ref('')
  const clearHandle = () => {
    leftVal.value = ''
    rightVal.value = ''
  }
  const encodingHandle = () => {
    rightVal.value = encodeURIComponent(leftVal.value)
  }
  const decodingHandle = () => {
    const data = leftVal.value.replace(/\+/g, '%20')
    rightVal.value = decodeURIComponent(data)
    // rightVal.value = decodeURIComponent(leftVal.value)
  }
  const transformHandle = () => {
    const left = leftVal.value
    const right = rightVal.value
    leftVal.value = right
    rightVal.value = left
  }
</script>

<template>
  <div class="box">
    <div class="left">
      <h4>需要转换的内容</h4>
      <el-input v-model="leftVal" resize="none" type="textarea" />
    </div>
    <div class="center">
      <el-button type="primary" @click="decodingHandle">URL解码</el-button>
      <el-button type="primary" @click="encodingHandle">URL编码</el-button>
      <el-button type="primary" @click="transformHandle">交换内容</el-button>
      <el-button :icon="Delete" @click="clearHandle">清空</el-button>
    </div>
    <div class="right">
      <h4>转换结果</h4>
      <el-input v-model="rightVal" resize="none" type="textarea" />
    </div>
  </div>
</template>

<style scoped lang="scss">
  .box {
    display: flex;
    .left,
    .right {
      width: 446px;
      height: 550px;
      h4 {
        height: 40px;
        background: #f8fbff;
        border: 1px solid #e1e3ec;
        border-bottom: 0;
        text-indent: 18px;
        line-height: 40px;
        margin: 0;
      }
      :deep() {
        .el-textarea {
          border: 0;
        }
        .el-textarea__inner {
          height: 508px;
          border: 0;
          border-radius: 0;
          // box-shadow: none;
        }
      }
    }
    .center {
      width: 156px;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      margin: 0 auto;
      .el-button,
      .el-select {
        width: 100%;
        margin: 0 0 14px 0;
      }
    }
  }
</style>
