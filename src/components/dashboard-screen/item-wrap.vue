<script setup lang="ts">
  import BorderBox from './border-box.vue'
  const props = withDefaults(
    defineProps<{
      // 标题
      title: number | string
    }>(),
    {
      title: '',
    }
  )
</script>

<template>
  <BorderBox>
    <div v-if="title" class="item_title">
      <span class="title-inner">{{ title }}</span>
    </div>
    <div class="item_content">
      <slot></slot>
    </div>
  </BorderBox>
</template>

<style scoped lang="scss">
  $item-title-height: 88px;
  $item_title_content-height: calc(100% - 100px);

  .item_title {
    height: $item-title-height;
    line-height: $item-title-height;
    width: 100%;
    font-size: 36px;
    font-weight: 600;
    color: #c2d6f3;
    text-indent: 40px;
    .title-inner {
      font-weight: 900;
      letter-spacing: 2px;
      // background: linear-gradient(92deg, #0072ff 0%, #00eaff 48.8525390625%, #01aaff 100%);
      // -webkit-background-clip: text;
      // -webkit-text-fill-color: transparent;
    }
  }

  .item_content {
    height: $item_title_content-height;
    :deep() {
      div[class*='-container'] {
        height: 100%;
        width: calc(100% - 8px);
        margin: 0 auto;
      }
    }
  }
</style>
