<script lang="ts">
  export default {
    name: 'TopBar',
  }
</script>
<script setup lang="ts">
  import { Ref } from 'vue'

  const props = defineProps<{
    title: string
    mode: string
  }>()

  const showDialog = inject('activekey') as Ref<boolean>

  const modeNmae = inject('modekey') as Ref<string>

  const reflesh = () => {}

  const openDialog = () => {
    showDialog.value = true
    modeNmae.value = props.mode
  }
</script>

<template>
  <div class="top-bar">
    <div class="left">{{ title }}</div>
    <div class="right">
      <!-- <vab-icon class="btn" icon="refresh-line" @click="reflesh" /> -->
      <vab-icon class="btn" icon="fullscreen-fill" @click="openDialog" />
    </div>
  </div>
</template>

<style scoped lang="scss">
  .top-bar {
    padding: 10px;
    height: 53px;
    border-bottom: 1px solid #dcdfe6;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .left {
      font-size: 13px;
    }

    .btn {
      font-size: 13px;
      vertical-align: -webkit-baseline-middle;

      &:hover {
        cursor: pointer;
      }

      &:nth-child(1) {
        margin-right: 8px;
      }
    }
  }
</style>
