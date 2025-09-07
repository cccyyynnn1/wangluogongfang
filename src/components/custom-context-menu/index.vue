<script lang="ts">
  import { defineComponent, ref, computed, nextTick, PropType, Ref, watch, onUnmounted, Fragment } from 'vue'
  import { cloneDeep, debounce } from 'lodash'
  export default defineComponent({
    name: 'CustomContextMenu',
    props: {
      appendToBody: {
        type: Boolean,
        default: true,
      },
      menuList: {
        type: Array as PropType<any[]>,
        required: true,
      },
      customClass: {
        type: String,
        default: '',
      },
      el: {
        type: Object as PropType<HTMLElement>,
        required: true,
      },
      params: {
        type: [String, Number, Array, Object] as PropType<any>,
        default: '',
      },
      disabled: {
        type: Function as PropType<() => boolean>,
        default: () => false,
      },
    },
    emits: ['open', 'close'],
    expose: ['show', 'close', 'showMenu'],
    setup(props, { emit }) {
      const showMenu = ref(false)
      const menuTop = ref(0)
      const menuLeft = ref(0)
      const contextmenuList = ref()
      const menuWrapper = ref()
      const clickDomEl = ref<HTMLElement>()
      const contextMenuShow = async (x = 0, y = 0) => {
        clickDomEl.value = document.elementFromPoint(x - 1, y - 1) as HTMLElement
        showMenu.value = true
        if (!showMenu.value) return
        contextmenuList.value = cloneDeep(props.menuList)
        await nextTick()
        menuLeft.value = x
        menuTop.value = y
      }
      const contextMenuClose = (e: Event) => {
        if (menuWrapper.value && !menuWrapper.value.contains(e.currentTarget)) {
          showMenu.value = false
          document.oncontextmenu = null
        }
      }

      const handleMenuItemClick = (item: any, $event: MouseEvent) => {
        const { callback, disabled, ...other } = item
        if (disabled) return
        if (callback && typeof callback === 'function') {
          const flag = callback(other, $event)
          if (flag === false) return
        }
        showMenu.value = false
      }

      onMounted(() => {
        document.addEventListener('mousedown', contextMenuClose)
      })
      onUnmounted(() => {
        document.removeEventListener('mousedown', contextMenuClose, true)
      })
      return {
        menuTop,
        menuLeft,
        showMenu,
        menuWrapper,
        handleMenuItemClick,
        show: contextMenuShow,
        close: contextMenuClose,
      }
    },
  })
</script>

<template>
  <div>
    <teleport :disabled="!appendToBody" to="body">
      <div
        v-if="showMenu"
        ref="menuWrapper"
        :class="['custom-context-menu', customClass]"
        :style="{ top: `${menuTop}px`, left: `${menuLeft}px` }"
      >
        <div
          v-for="menu of menuList"
          :key="menu.tips"
          class="itme-menu"
          @mousedown.stop="handleMenuItemClick(menu, $event)"
        >
          {{ menu.label }}
        </div>
      </div>
    </teleport>
  </div>
</template>

<style scoped lang="scss">
  .custom-context-menu {
    width: 156px;
    background: #ffffff;
    box-shadow: 0px 0px 4px 4px rgba(77, 81, 86, 0.05);
    border-radius: 8px;
    border: 1px solid #eceafb;
    border-radius: 6px;
    position: absolute;
    z-index: 999999;
    padding: 10px 5px;
    .itme-menu {
      border-radius: 4px;
      height: 36px;
      line-height: 36px;
      text-indent: 10px;
      color: #606266;
      cursor: pointer;
      &:hover {
        background: #f3f1fe;
        color: var(--el-color-primary);
      }
    }
  }
</style>
