<script lang="ts" setup>
  import { useSettingsStore } from '@/store/modules/settings'
  import { useFavicon } from '@vueuse/core'
  // @ts-ignore
  const icon = useFavicon()

  const settingsStore = useSettingsStore()
  const { theme, logo, systemConfig } = storeToRefs(settingsStore)
  watchEffect(() => {
    icon.value = logo.value
  })
</script>

<template>
  <div
    class="logo-container"
    :class="{
      ['logo-container-' + theme.layout]: true,
    }"
  >
    <router-link to="/">
      <span class="logo">
        <!-- 使用自定义svg示例 -->
        <el-image :src="logo" />
        <!-- <el-image :src="logo1" v-else /> -->
      </span>
      <span class="title" :class="{ 'hidden-xs-only': theme.layout === 'horizontal' }">
        {{ systemConfig.value.nameAbb }}
      </span>
    </router-link>
  </div>
</template>

<style lang="scss" scoped>
  @mixin container {
    position: relative;
    height: $base-header-height;
    overflow: hidden;
    line-height: $base-header-height;
    background: transparent;
  }

  @mixin logo {
    display: inline-block;
    width: 32px;
    height: 32px;
    color: $base-title-color;
    vertical-align: middle;
    fill: currentColor;
  }

  @mixin title {
    display: inline-block;
    margin-left: 5px;
    overflow: hidden;
    font-size: 20px;
    line-height: 55px;
    color: $base-title-color;
    text-overflow: ellipsis;
    white-space: nowrap;
    // background: $base-column-second-menu-background !important;
    // vertical-align: middle;
  }

  .logo-container {
    &-horizontal,
    &-common {
      @include container;

      .logo {
        svg,
        img {
          @include logo;
        }
      }

      .title {
        @include title;
      }
    }

    &-vertical,
    &-column,
    &-comprehensive,
    &-float {
      @include container;

      height: $base-logo-height;
      line-height: $base-logo-height;
      text-align: center;

      .logo {
        svg,
        img {
          @include logo;
        }
      }

      .title {
        @include title;
        max-width: calc(var(--el-left-menu-width) - 60px);
      }
    }

    &-column {
      // background: $base-column-second-menu-background !important;

      .logo {
        position: fixed;
        top: 0;
        display: block;
        width: $base-left-menu-width-min;
        height: $base-logo-height;
        margin: 0;
        width: 64px;
        // background: $base-column-first-menu-background !important;
      }

      .title {
        // padding-right: 15px;
        // padding-left: 15px;
        margin-left: $base-left-menu-width-min !important;
        // color: var(--el-color-black) !important;
        background: $base-column-second-menu-background !important;
        @include title;
      }
    }

    :deep() {
      img.el-image__inner {
        @include logo;
      }
    }
  }
</style>
