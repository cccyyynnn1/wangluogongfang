<script setup lang="ts">
  import useContextMenu from '@/hooks/useContextMenu'

  const { x, y, isShow, list } = useContextMenu()

  // 点击每一项
  const handleClick = (item: any) => {
    isShow.value = false
    item.callback(item.value)
  }

  onMounted(() => {
    // console.log(x.value, y.value, isShow.value)
  })
</script>

<script lang="ts">
  export default {
    name: 'ContextMenu',
  }
</script>

<template>
  <Teleport to="body">
    <transition>
      <div v-if="isShow" class="contextmenu-content" :style="{ left: x + 'px', top: y + 'px' }">
        <div class="list">
          <div v-for="(item, index) in list" :key="index" class="item" @click="handleClick(item)">
            {{ item.label }}
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<style scoped lang="scss">
  .contextmenu-content {
    // height: 100px;
    position: fixed;
    left: 999999px;
    top: 999999px;
    z-index: 5000;
    user-select: none;
  }
  .list {
    // height: 100px;
    border: 1px solid var(--el-border-color);
    border-radius: 4px;
    min-width: 180px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    word-break: break-all;
    word-wrap: break-word;
  }
  .item {
    box-sizing: border-box;
    padding: 0 5px;
    height: 30px;
    line-height: 30px;
    word-break: keep-all;
    background-color: #fff;
  }
  .item:hover {
    background-color: var(--el-color-primary);
    color: #fff;
    cursor: pointer;
  }
  .v-enter-from {
    opacity: 0;
  }
  .v-enter-to {
    opacity: 1;
    transition: 0.5s;
  }
</style>
