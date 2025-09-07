<script lang="ts">
  export default {
    name: 'Message', // 系统消息按钮
  }
</script>

<script setup lang="ts">
  import { useSettingsStore } from '@/store/modules/settings'
  import { usePubilcStore } from '@/store/modules/public'
  const settingsStore = useSettingsStore()
  const publicStore = usePubilcStore()
  const { hasNewMessage } = storeToRefs(publicStore)
  const handleClick = () => {
    settingsStore.changeMessageVisible(true)
  }
</script>

<template>
  <div class="message" @click="handleClick">
    <!-- <vab-icon icon="download-2-line" /> -->
    <img alt="" class="img" src="@/assets/message_image/message.png" />
    <div v-if="hasNewMessage" class="reamrk"></div>
  </div>
</template>

<style scoped lang="scss">
  .img {
    height: 36px;
    width: 36px;
  }
  .message {
    // margin-bottom: 10px;
    position: relative;
    border-radius: 5px;
    &:hover {
      color: #7f71eb;
      // background-color: #f0eefd !important;
    }
    .reamrk {
      width: 9px;
      height: 9px;
      background: #ff6767;
      border-radius: 50%;
      position: absolute;
      border: 1px solid #d9d5ef;
      top: 7px;
      right: 7px;
    }
  }
</style>
