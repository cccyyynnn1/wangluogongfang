<script lang="ts">
  export default {
    name: 'EncodingDecoding',
  }
</script>
<script setup lang="ts">
  import Base64 from './base64.vue'
  import Uri from './uri.vue'
  import Unicode from './unicode.vue'
</script>

<template>
  <el-tabs class="demo-tabs" model-value="base64">
    <el-tab-pane label="Base64编解码" name="base64">
      <base64 />
    </el-tab-pane>
    <el-tab-pane label="URL编解码" name="url">
      <uri />
    </el-tab-pane>
    <el-tab-pane label="Unicode中文互转" name="unicode">
      <unicode />
    </el-tab-pane>
  </el-tabs>
</template>

<style scoped lang="scss"></style>
