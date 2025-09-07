<script lang="ts">
  export default {
    name: 'DrawerSelf', // 自定义连线
  }
</script>

<script setup lang="ts">
  import { node } from '~/src/types'
  import DrawerMemorise from './drawer-memorise.vue'
  import VueEvent from '@/data/event'
  import _ from 'lodash'
  import { default as Sortable, SortableEvent } from 'sortablejs'

  const emits = defineEmits<{
    (e: 'on-close'): void
  }>()

  const drawerMemorise = ref(false)
  const checkValue = ref(false)
  // let isSortable = false
  const clientIPList = ref<node[]>([])
  const serverIPList = ref<node[]>([])
  const remark = ref()

  const handleMouseenterSever = () => {
    remark.value = 'insideServerIsOutside'
    VueEvent.emit('insideIsOutside')
  }

  const handleMouseenterClient = () => {
    remark.value = 'insideClientIsOutside'
    VueEvent.emit('insideIsOutside')
  }

  const handleMouseleave = () => {
    remark.value = ''
    VueEvent.emit('insideIsOutside')
  }

  const curData = ref()
  const getResult = (val: any) => {
    const { data, isDrag } = val
    curData.value = data
    if (!isDrag) return
    const flag = clientIPList.value.some((item: any) => {
      return item.id == curData.value.id
    })
    const flag1 = serverIPList.value.some((item: any) => {
      return item.id == curData.value.id
    })
    if (flag && remark.value == 'insideServerIsOutside') {
      // 如果node存在于目的ip中且要放入源ip
      clientIPList.value = clientIPList.value.filter((item: node) => {
        return item.id !== curData.value.id
      })
      serverIPList.value.push(curData.value)
    }
    if (flag1 && remark.value == 'insideClientIsOutside') {
      // 如果node存在于源ip中且要放入目的ip
      serverIPList.value = serverIPList.value.filter((item: node) => {
        return item.id !== curData.value.id
      })
      clientIPList.value.push(curData.value)
    }
    if (!flag && !flag1) {
      if (remark.value == 'insideClientIsOutside') {
        clientIPList.value.push(curData.value)
      } else if (remark.value == 'insideServerIsOutside') {
        serverIPList.value.push(curData.value)
      }
    }
    // 还在拖拽中但已经离开拖拽区域
    if (!remark.value && isDrag) {
      serverIPList.value = serverIPList.value.filter((item: node) => {
        return item.id !== curData.value.id
      })
      clientIPList.value = clientIPList.value.filter((item: node) => {
        return item.id !== curData.value.id
      })
    }
  }

  const handelDelItem = (node: node) => {
    if (remark.value == 'insideClientIsOutside') {
      clientIPList.value = clientIPList.value.filter((item: node) => {
        return item.ip != node.ip
      })
    } else {
      serverIPList.value = serverIPList.value.filter((item: node) => {
        return item.ip != node.ip
      })
    }
  }

  onMounted(() => {
    initData()
  })

  let sortableServerInstance: Sortable
  let sortableClientInstance: Sortable
  nextTick(() => {
    const server = document.querySelector('.drag_content_server')
    const client = document.querySelector('.drag_content_client')
    if (!server) return
    sortableServerInstance = Sortable.create(server as HTMLElement, {
      group: { name: 'my_drag_item_name', pull: true, put: true },
      //  可被拖拽的子元素
      draggable: '.my_drag_item',
      delay: 0,
      sort: false,
      forceFallback: false, // 忽略 HTML5拖拽行为，强制回调进行

      // // 开始拖拽的时候
      // onStart(/**Event*/ evt) {
      //   curData.value = serverIPList.value[evt.oldIndex!]
      //   isSortable = true
      // },

      // onEnd(event: SortableEvent) {
      //   isSortable = false
      //   curData.value = undefined
      // },
    })
    if (!client) return
    sortableClientInstance = Sortable.create(client as HTMLElement, {
      group: { name: 'my_drag_item_name', pull: true, put: true },
      //  可被拖拽的子元素
      draggable: '.my_drag_item',
      delay: 0,
      disabled: false,
      sort: false,
      forceFallback: false, // 忽略 HTML5拖拽行为，强制回调进行
      // onDragLeave() {
      //   console.log(132)
      // },
      // // 开始拖拽的时候
      // onStart(/**Event*/ evt) {
      //   curData.value = clientIPList.value[evt.oldIndex!]
      //   isSortable = true
      // },
      // onEnd(event: SortableEvent) {
      //   isSortable = false
      //   curData.value = undefined
      // },
    })
  })

  const initData = () => {
    clientIPList.value = []
    serverIPList.value = []
  }

  // 自定义连线交换值
  const handleExchange = () => {
    const temp = _.cloneDeep(serverIPList.value)
    serverIPList.value = _.cloneDeep(clientIPList.value)
    clientIPList.value = temp
  }

  // 监听事件
  VueEvent.on('insideIsOutsideResult', getResult)

  onUnmounted(() => {
    sortableClientInstance && sortableClientInstance.destroy()
    sortableClientInstance && sortableServerInstance.destroy()
    // 事件销毁
    VueEvent.off('insideIsOutsideResult')
  })
</script>

<template>
  <div class="drawer-self sketchpad-drawer">
    <div class="sketchpad-drawer-top">
      <span style="font-size: 18px; color: #1e1842; margin-right: 6px">自定义连线</span>
      <div style="display: flex; align-items: center">
        <el-button plain size="small" style="margin-right: 6px" @click="drawerMemorise = true">
          <img
            alt=""
            :src="require('@/assets/chase/memorise.svg')"
            style="width: 16px; height: 14px; margin-right: 4px"
          />
          记忆库
        </el-button>
        <el-icon @click="emits('on-close')"><Close /></el-icon>
      </div>
    </div>
    <div class="drawer-self-warp">
      <div class="drawer-self-warp-content">
        <div class="drawer-self-warp-content-sourse" @mouseenter="handleMouseenterSever" @mouseleave="handleMouseleave">
          <div class="drag_top">源IP：</div>
          <div class="drag_content drag_content_server">
            <!-- <transition-group name="fade"> -->
            <div
              v-for="item in serverIPList"
              :key="item.id"
              class="my_drag_item"
              :style="{ opacity: item.id == curData?.id ? '0.5' : '1' }"
            >
              <div class="img">
                <img v-if="item.icon === 'user'" alt="" :src="require('@/assets/chase/default.svg')" />
                <img v-else alt="" :src="require('@/assets/chase/server.svg')" />
                <div class="img_del" @click="handelDelItem(item)">
                  <el-icon style="color: #fff; font-size: 8px"><Close /></el-icon>
                </div>
              </div>
              <div>{{ item.ip }}</div>
            </div>
            <!-- </transition-group> -->
          </div>
        </div>
        <div class="split" @click="handleExchange">
          <img
            alt=""
            :src="require('@/assets/chase/exchange.svg')"
            style="width: 18px; height: 12px; margin-right: 4px"
          />
          <span>交换源/目IP值</span>
        </div>
        <div
          class="drawer-self-warp-content-client"
          @mouseenter="handleMouseenterClient"
          @mouseleave="handleMouseleave"
        >
          <div class="drag_top">目的IP：</div>
          <div class="drag_content drag_content_client">
            <transition-group name="fade">
              <div
                v-for="item in clientIPList"
                :key="item.id"
                class="my_drag_item"
                :style="{ opacity: item.id == curData?.id ? '0.5' : '1' }"
              >
                <div class="img">
                  <img v-if="item.icon === 'user'" alt="" :src="require('@/assets/chase/default.svg')" />
                  <img v-else alt="" :src="require('@/assets/chase/server.svg')" />
                  <div class="img_del" @click="handelDelItem(item)">
                    <el-icon style="color: #fff; font-size: 8px"><Close /></el-icon>
                  </div>
                </div>
                <div>{{ item.ip }}</div>
              </div>
            </transition-group>
          </div>
        </div>
      </div>
      <div class="drawer-self-warp-tool">
        <el-checkbox v-model="checkValue" label="添加到记忆库，永久保存" />
        <el-button type="primary">保存连线</el-button>
      </div>
    </div>
    <!-- 记忆库的新增编辑 -->
    <DrawerMemorise v-if="drawerMemorise" @on-close="drawerMemorise = false" />
  </div>
</template>

<style scoped lang="scss">
  .drawer-self-warp {
    margin-top: 10px;
    height: calc(100% - 50px);
    width: 100%;

    .drawer-self-warp-content {
      border-radius: 4px;
      border: 1px solid #e0ddf6;
      height: calc(100% - 45px);
      position: relative;
      .drawer-self-warp-content-sourse,
      .drawer-self-warp-content-client {
        height: 50%;
        padding: 0 10px;
        overflow: hidden;
        .drag_top {
          height: 28px;
          line-height: 28px;
          align-items: center;
          font-weight: 500;
          font-size: 13px;
          color: #1e1842;
        }
        .drag_content {
          min-height: 175px;
          width: 297px;
          &::-webkit-scrollbar {
            width: 0;
            height: 0;
          }
          display: flex;
          flex-wrap: wrap;
          :deep() {
            .fade-enter-active,
            .fade-leave-active {
              transition: opacity 0.2s ease;
            }
            .fade-enter-from,
            .fade-leave-to {
              opacity: 0;
            }
            .my_drag_item {
              height: 88px;
              width: 99px;
              overflow: hidden;
              display: flex;
              align-items: center;
              justify-content: center;
              flex-direction: column;
              pointer-events: all;
            }
          }
          .img {
            width: 48px;
            height: 48px;
            border-radius: 50%;
            background: #ffffff;
            box-shadow: 0px 0px 4px 4px rgba(77, 81, 86, 0.05);
            border: 1px solid #eae7ff;
            display: flex;
            align-items: center;
            justify-content: center;
            position: relative;
            // cursor: pointer;
            margin: 0 auto;
            margin-bottom: 2px;
            img {
              height: 16px;
              width: 16px;
            }
            &:hover {
              .img_del {
                background: #f56c6c;
              }
            }
            .img_del {
              cursor: pointer;
              position: absolute;
              width: 16px;
              height: 16px;
              top: 0;
              right: 0;
              border-radius: 50%;
              display: flex;
              align-items: center;
              justify-content: center;
            }
          }
        }
      }
      .split {
        position: absolute;
        width: 124px;
        height: 30px;
        margin: 0 auto;
        background: #ffffff;
        border-radius: 4px;
        border: 1px solid #e5e4ef;
        display: flex;
        align-items: center;
        justify-content: center;
        left: 97px;
        top: 50%;
        transform: translateY(-50%);
        cursor: pointer;
        &::before,
        &::after {
          position: absolute;
          content: '';
          width: 99px;
          height: 1px;
          background: #e5e4ef;
          top: 14px;
        }
        &::before {
          left: -99px;
        }
        &::after {
          right: -99px;
        }
      }
    }
    .drawer-self-warp-tool {
      margin-top: 15px;
      height: 30px;
      width: 100%;
      display: flex;
      align-items: center;
      justify-content: space-between;
    }
  }
</style>
