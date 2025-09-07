<script lang="ts" setup>
  import Download from './download.vue'
  import ToolBox from './toolbox.vue'
  import Message from './message.vue'
  import { VabRoute } from '/#/router'
  import { isExternal } from '@/utils/validate'
  import { translate } from '@/i18n'
  import { useRoutesStore } from '@/store/modules/routes'
  import { defaultOpeneds, openFirstMenu, uniqueOpened } from '@/config'
  import { useSettingsStore } from '@/store/modules/settings'
  import variables from '@vab/styles/variables/variables.module.scss'
  const route: VabRoute = useRoute()
  const router = useRouter()

  const settingsStore = useSettingsStore()
  const { theme, collapse } = storeToRefs(settingsStore)
  const routesStore = useRoutesStore()
  const {
    getTab: tab,
    getTabMenu: tabMenu,
    getActiveMenu: activeMenu,
    getRoutes: routes,
    getPartialRoutes: partialRoutes,
  }: any = storeToRefs(routesStore)

  const handleTabClick = () => {
    nextTick(() => {
      if (isExternal(tabMenu.value.path)) {
        window.open(tabMenu.value.path)
        setTimeout(() => {
          router.push('/')
        }, 1000)
      } else if (openFirstMenu) router.push(tabMenu.value.redirect || tabMenu.value)
    })
  }
</script>
<template>
  <el-scrollbar
    class="vab-column-bar-container"
    :class="{
      'is-collapse': collapse,
      ['vab-column-bar-container-' + theme.columnStyle]: true,
    }"
  >
    <vab-logo />
    <el-tabs v-model="tab.data" tab-position="left" @tab-click="handleTabClick">
      <template v-for="(item, index) in routes" :key="index + item.name">
        <el-tab-pane :name="item.name">
          <template #label>
            <div
              class="vab-column-grid"
              :class="{
                ['vab-column-grid-' + theme.columnStyle]: true,
              }"
              :title="translate(item.meta.title)"
            >
              <div>
                <vab-icon v-if="item.meta.icon" :icon="item.meta.icon" :is-custom-svg="item.meta.isCustomSvg" />
                <span>
                  {{ translate(item.meta.title) }}
                </span>
              </div>
            </div>
          </template>
        </el-tab-pane>
      </template>
    </el-tabs>
    <el-menu
      :background-color="variables['column-second-menu-background']"
      :default-active="activeMenu.data?.split('?')[0]"
      :default-openeds="defaultOpeneds"
      mode="vertical"
      :unique-opened="uniqueOpened"
    >
      <el-divider>
        {{ translate(tabMenu ? tabMenu.meta.title : tabMenu) }}
      </el-divider>
      <template v-for="item in partialRoutes" :key="item.path">
        <vab-menu v-if="!item.meta.hidden" :item="item" />
      </template>
    </el-menu>
    <div
      class="toolbox1"
      :class="{
        ['toolbox-' + theme.columnStyle]: true,
      }"
    >
      <Message />
      <Download :data="false" />
      <ToolBox />
      <div class="setting">
        <el-popover
          placement="right"
          popper-style="box-shadow: rgb(14 18 22 / 35%) 0px 10px 38px -10px, rgb(14 18 22 / 20%) 0px 10px 20px -15px; padding: 20px;"
          :width="200"
        >
          <template #reference>
            <img alt="" class="img" src="@/assets/tools/avatar.png" />
          </template>
          <template #default>
            <vab-avatar />
          </template>
        </el-popover>
      </div>
      <span class="text" style="font-size: 12px">v3.0.5 RC</span>
    </div>
  </el-scrollbar>
</template>

<style lang="scss" scoped>
  @use 'sass:math';
  @mixin active {
    &:hover {
      color: #7f71eb;
      // background-color: #f0eefd !important;

      i,
      svg {
        color: #7f71eb;
      }
    }

    &.is-active {
      color: #7f71eb;
      // background-color: #f0eefd !important;
    }
  }
  .img {
    height: 36px;
    width: 36px;
  }
  .download {
    margin-top: 5px;
    border-radius: 5px;
    &:hover {
      color: #7f71eb;
      // background-color: #f0eefd !important;
    }
  }
  .apps,
  .setting {
    margin-top: 5px;
    border-radius: 5px;
    &:hover {
      color: #7f71eb;
      // background-color: #f0eefd !important;
    }
  }

  .text {
    color: #fff;
    opacity: 0.8;
  }

  .toolbox1 {
    height: 185px;
    background-color: transparent;
    position: fixed;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
    bottom: 0px;
    width: $base-left-menu-width-min;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    // height: 0px;
    // border-radius: 5px;
    text-align: center;
    // color: var(--el-color-white);
    cursor: pointer;

    i {
      display: block;
      margin-top: 4px;
    }
    &:hover {
      // color: var(--el-color-white);
      // background: var(--el-color-white) !important;
      cursor: pointer;
    }
    &.toolbox-horizontal {
      width: 83px;
      height: 50px;
      line-height: 50px;
      border-radius: 0;
      text-align: left;
      left: 0;
      i {
        display: inline-block;
        margin-top: 0;
        margin-left: 10px;
        margin-right: 3px;
      }
    }
  }
  .vab-column-bar-container {
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    width: var(--el-left-menu-width);
    height: 100vh;
    // overflow: hidden;
    background: $base-column-second-menu-background;
    box-shadow: $base-box-shadow;
    :deep() {
      * {
        transition: $base-transition;
      }

      .el-tabs {
        box-shadow: $base-box-shadow;
        height: 100vh !important;
      }
    }

    &-vertical,
    &-card,
    &-arrow {
      :deep() {
        .el-tabs + .el-menu {
          left: calc($base-left-menu-width-min + 4px);
          width: calc(var(--el-left-menu-width) - #{$base-left-menu-width-min} - 4px);
          border: 0;
        }
      }
    }

    &-horizontal {
      :deep() {
        .logo-container-column {
          .logo {
            width: $base-left-menu-width-min * 1.3 !important;
          }

          .title {
            margin-left: $base-left-menu-width-min * 1.3 !important;
          }
        }

        .el-tabs + .el-menu {
          left: $base-left-menu-width-min * 1.3;
          width: calc(var(--el-left-menu-width) - #{$base-left-menu-width-min} * 1.3);
          border: 0;
        }
      }
    }

    &-card {
      :deep() {
        .el-tabs {
          .el-tabs__item {
            padding: 5px !important;

            .vab-column-grid {
              width: $base-left-menu-width-min - 10 !important;
              height: $base-left-menu-width-min - 10 !important;
              border-radius: 5px;
            }

            &.is-active {
              background: transparent !important;

              .vab-column-grid {
                background: #7f71eb;
              }
            }
          }
        }

        .el-tabs + .el-menu {
          left: $base-left-menu-width-min + 10;
          width: calc(var(--el-left-menu-width) - #{$base-left-menu-width-min} - 20px);
        }

        .el-sub-menu .el-sub-menu__title,
        .el-menu-item {
          min-width: 180px;
          margin-bottom: 5px;
          border-radius: 5px;
        }
      }
    }

    &-arrow {
      :deep() {
        .el-tabs {
          .el-tabs__item {
            &.is-active {
              background: transparent !important;

              .vab-column-grid {
                background: transparent !important;

                &:after {
                  position: absolute;
                  right: 0;
                  width: 0;
                  height: 0;
                  overflow: hidden;
                  content: '';
                  border-color: transparent #{var(--el-color-white)} transparent transparent;
                  border-style: solid dashed dashed;
                  border-width: 8px;
                }
              }
            }
          }
        }

        .el-tabs + .el-menu {
          left: $base-left-menu-width-min + 10;
          width: calc(var(--el-left-menu-width) - #{$base-left-menu-width-min} - 20px);
        }

        .el-sub-menu .el-sub-menu__title,
        .el-menu-item {
          min-width: 180px;
          margin-bottom: 5px;
          border-radius: 5px;
        }
      }
    }

    .vab-column-grid {
      display: flex;
      align-items: center;
      width: $base-left-menu-width-min;
      overflow: hidden;
      text-align: center;
      text-overflow: ellipsis;
      word-break: break-all;
      white-space: nowrap;

      &-vertical,
      &-card,
      &-arrow {
        justify-content: center;
        height: $base-left-menu-width-min;

        > div {
          svg {
            position: relative;
            top: 8px;
            display: block;
            width: $base-font-size-default + 4;
            height: $base-font-size-default + 4;
          }

          [class*='ri-'] {
            display: block;
            height: 20px;
          }
        }
      }

      &-horizontal {
        justify-content: left;
        width: $base-left-menu-width-min * 1.3;
        height: #{math.div($base-left-menu-width-min, 1.3)};
        padding-left: #{math.div($base-padding, 2)};
      }
    }

    :deep() {
      .el-scrollbar__wrap {
        overflow-x: hidden;
      }

      .el-tabs {
        position: fixed;

        .el-tabs__header.is-left {
          margin-right: 0 !important;

          .el-tabs__nav-wrap.is-left {
            margin-right: 0 !important;
            background: $base-column-first-menu-background;

            .el-tabs__nav-scroll {
              height: 100%;
            }
          }
        }

        .el-tabs__nav {
          height: calc(100vh - #{$base-logo-height} - 185px);
          background: $base-column-first-menu-background;
          overflow-y: auto;

          &::-webkit-scrollbar {
            width: 0;
            height: 0;
          }
        }

        .el-tabs__item {
          height: auto;
          padding: 0;
          height: 60px;
          font-size: 12px;
          color: var(--el-color-white);

          &.is-active {
            background: var(--el-color-primary);
          }
        }
      }

      .el-tabs__active-bar.is-left,
      .el-tabs--left .el-tabs__nav-wrap.is-left::after {
        display: none;
      }

      .el-menu {
        border: 0;

        .el-divider {
          margin: 0 0 $base-margin 0;
          background-color: $base-column-second-menu-background;

          &__text {
            color: $base-title-color;
            background-color: $base-column-second-menu-background;
          }
        }

        .el-menu-item,
        .el-sub-menu__title {
          height: $base-menu-item-height;

          overflow: hidden;
          line-height: $base-menu-item-height;
          text-overflow: ellipsis;
          white-space: nowrap;
          vertical-align: middle;

          // @include active;
        }
      }
    }

    &.is-collapse {
      :deep() {
        width: 64px;
      }
    }
  }
</style>
