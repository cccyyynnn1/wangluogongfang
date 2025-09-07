<script lang="ts" setup>
  import { useRoutesStore } from '@/store/modules/routes'
  import { handleActivePath } from '@/utils/routes'
  import { VabRoute } from '/#/router'
  import { useSettingsStore } from '@/store/modules/settings'
  import { useUserStore } from '@/store/modules/user'
  import { proxyNet } from '@/config/index'
  import { GlobalMessageItem } from '~/types/store'
  import { ElNotification, messageEmits } from 'element-plus'
  import { usePubilcStore } from '@/store/modules/public'
  import { GlobalMsgSwitchIdApi } from '~/src/api-ecs/public'

  const publicStore = usePubilcStore()
  const $baseMessage: any = inject('$baseMessage')
  const { setGlobalMessageList, SetFlowProbes } = publicStore
  const { globalLoginMessageList, hasNewMessage } = storeToRefs(publicStore)
  const userStore = useUserStore()
  const {
    setTableColumns,
    getAllHightLight,
    highLightLastUpdateTime,
    indexFieldsLastUpdateTime,
    setIndexFieldsUpdateTime,
    sethighLightLastUpdateTime,
  } = userStore
  const route: VabRoute = useRoute()
  const router = useRouter()
  const routesStore: any = useRoutesStore()
  const { tab, activeMenu } = storeToRefs(routesStore)
  const settingsStore = useSettingsStore()
  const { foldSideBar, openSideBar, toggleCollapse } = settingsStore
  const { collapse } = storeToRefs(settingsStore)
  const showIcon = ref(false)

  onMounted(() => {
    if (route.matched[0].children.length <= 1) {
      foldSideBar()
    } else {
      setTimeout(() => {
        showIcon.value = true
      }, 100)
      openSideBar()
    }
    SetFlowProbes()
  })

  const setMessageTop = () => {
    const nodes = document.querySelectorAll('.global-msg-notification')
    nodes.forEach((_, index) => {
      nodes[index].setAttribute('style', `top: ${106 * index + 16}px !important;z-index: 3000`)
    })
  }

  const activeRef = ref(false)

  watch(
    route,
    () => {
      if (tab.value.data !== route.matched[0].name) tab.value.data = route.matched[0].name
      activeMenu.value.data = handleActivePath(route)
      activeRef.value =
        activeMenu.value.data.includes('/site_index') || activeMenu.value.data == '/alerts/threat-intelligence'
          ? true
          : false
      if (route.matched[0].children.length <= 1) {
        foldSideBar()
        showIcon.value = false
      } else {
        openSideBar()
        setTimeout(() => {
          showIcon.value = true
        }, 300)
      }
    },
    { immediate: true }
  )

  const url = process.env.NODE_ENV == 'development' ? proxyNet : window.location.hostname
  const globalMsgUrl = `wss://${url}/v3/ecsPlatform/websocket/globalMsg?EcsSessionId=${userStore.token}`
  const { data, close, send, status, open } = useWebSocket(globalMsgUrl, {
    autoReconnect: {
      retries: 3,
      delay: 6000,

      onFailed() {
        $baseMessage('WebSocket重连失败', 'error', 'vab-hey-message-error')
      },
    },
    onError(e) {
      $baseMessage('WebSocket链接失败', 'error', 'vab-hey-message-error')
    },
  })
  let timer = setInterval(() => {
    send('PING')
  }, 10000)
  const visibilitychange = () => {
    if (!document.hidden && (!status.value || status.value == 'CLOSED')) {
      open()
      clearInterval(timer)
      timer = setInterval(() => {
        send('PING')
      }, 10000)
    }
  }
  const updateOnlineStatus = (event: any) => {
    if (event.type == 'online' && (!status.value || status.value == 'CLOSED')) {
      open()
      clearInterval(timer)
      timer = setInterval(() => {
        send('PING')
      }, 10000)
    }
  }

  document.addEventListener('visibilitychange', visibilitychange)
  window.addEventListener('online', updateOnlineStatus)
  window.addEventListener('offline', updateOnlineStatus)

  onUnmounted(() => {
    close()
    clearInterval(timer)
    document.removeEventListener('visibilitychange', visibilitychange)
    window.removeEventListener('online', updateOnlineStatus)
    window.removeEventListener('offline', updateOnlineStatus)
  })

  class Task {
    private static instance: Task
    public task: GlobalMessageItem[]
    private falg: boolean
    constructor() {
      this.task = []
      this.falg = true
    }
    public static init(): Task {
      if (!Task.instance) {
        Task.instance = new Task()
      }
      return Task.instance
    }
    public next() {
      if (this.getFalg()) {
        this.setFlag(false)
        if (this.task.length > 0) {
          const item = this.task.shift()
          if (item) {
            // AddGlobalMessages(item)
            // item.msgType == 'sys' ? AddGlobalSysMessages(item) : AddGlobalUserMessages(item)
            router.currentRoute.value.fullPath == '/dashboard' &&
              (item.minorType == 'collectShare' || item.minorType == 'notice') &&
              setGlobalMessageList()
            ElNotification({
              dangerouslyUseHTMLString: true,
              customClass: 'global-msg-notification',
              message:
                `
                     <div class="item global-msg-notification-item global-msg-notification-td" data-id='${item.id}' style="z-index:${item.id} !important">
                        <div class="item_top" data-id='${item.id}'>
                           <div class="item_left" data-id='${item.id}'>` +
                `${
                  item.msgType == 'user'
                    ? `<div class="img user-img" data-id="${item.id}"></div>`
                    : `<div class="img sys-img" data-id="${item.id}"></div>`
                }` +
                `<div class="mark" data-id="${item.id}"></div>
                           </div>
                           <div class="item_right" data-id="${item.id}">
                              <div class="item_right-top" data-id="${item.id}">
                                <div class="item_right_name" data-id="${item.id}">[${item.title}]</div>
                                <div class="item_right_content" data-id="${item.id}">${item.notes}</div>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                           `,
              duration: 5000,
              onClick: handleCallBack,
            })
            nextTick(() => {
              setMessageTop()
            })
          }
        }
        setTimeout(() => {
          this.setFlag(true)
          if (this.task.length > 0) {
            this.next()
          }
        }, 2000)
      }
    }
    public add(arr: GlobalMessageItem[]) {
      arr.forEach((item) => {
        this.task.push(item)
      })
    }
    public delFirst() {
      this.task.shift()
    }
    public GetThisTask() {
      return this.task
    }

    private getFalg() {
      return this.falg
    }
    private setFlag(val: boolean) {
      this.falg = val
    }
  }

  watch(
    () => globalLoginMessageList.value,
    () => {
      const { newMsgs, unReadCount } = globalLoginMessageList.value
      const count = unReadCount
      if (count > 0) {
        setTimeout(() => {
          ElNotification({
            dangerouslyUseHTMLString: true,
            customClass: 'global-msg-notification global-login-msg-notification',
            message: `
                     <div class="item global-msg-notification-item" >
                        <div class="item_top">
                           <div class="item_left">
                             <div class="img sys-img"></div>
                             <div class="mark" ></div>
                           </div>
                           <div class="item_right">
                              <div class="item_right-top">
                                <div class="item_right_name">[提示]</div>
                                <div class="item_right_content">您有<span class="item_right_count">${count}</span>条未读新消息，快去看看</div>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                           `,
            duration: 6000,
            onClick: handleLoginCallBack,
          })
        }, 0)
      }
      if (newMsgs.length > 0) {
        newMsgs.forEach((item, index) => {
          setTimeout(() => {
            ElNotification({
              dangerouslyUseHTMLString: true,
              appendTo: '.vab-app-main',
              customClass: 'global-msg-notification-important',
              message: `
                             <div
                               class="item global-msg-notification-td"
                               data-id='${item.id}'
                               style="z-index:${item.id} !important"
                             >
                                <div class="item-left" data-id="${item.id}">
                                    <div class="item-left-img" data-id="${item.id}"></div>
                                    <span class="item-left-title" data-id="${item.id}">[${item.title}]</span>
                                    <span class="item-left-tip" data-id="${item.id}">${item.notes}</span>
                                </div>
                                <div class="global-msg-notification-item-btn" data-src="btn" data-type="${item.minorType}" data-id="${item.id}">
                                    <input type="checkbox" data-src="btn" data-type="${item.minorType}"  data-id="${item.id}" id="global-msg-notification-item-input" name="today" value="false">&nbsp;
                                    <sapn style="opacity:0.7" data-src="btn" data-type="${item.minorType}"  data-id="${item.id}">今日内，不在提示</span>
                                </div>
                             </div>
                             `,
              duration: 0,
              onClick: handleCallBack,
            })
          }, 2000 * index)
        })
      }
    },
    { immediate: true, deep: true }
  )

  watch(
    () => data.value,
    () => {
      if (data.value && data.value != 'ws服务连接成功!') {
        try {
          let newMsg: GlobalMessageItem[] | null = JSON.parse(data.value)?.newMsgs
          const GLOBAL_VALUE = JSON.parse(data.value)
          if (GLOBAL_VALUE?.indexFieldsLastUpdateTime !== indexFieldsLastUpdateTime) {
            setTableColumns()
            setIndexFieldsUpdateTime(GLOBAL_VALUE.indexFieldsLastUpdateTime)
          }
          if (GLOBAL_VALUE?.highLightLastUpdateTime !== highLightLastUpdateTime) {
            getAllHightLight()
            sethighLightLastUpdateTime(GLOBAL_VALUE.highLightLastUpdateTime)
          }
          hasNewMessage.value = GLOBAL_VALUE?.hasUnRead
          if (newMsg) {
            newMsg = newMsg.reverse()
            const messageList = newMsg.filter((item: GlobalMessageItem) => {
              return !item.important
            })
            const notificationList = newMsg.filter((item: GlobalMessageItem) => {
              return item.important
            })
            const task = Task.init()
            task.add(messageList)
            task.next()
            // messageList.forEach((item: GlobalMessageItem, index: number) => {
            //
            //   setTimeout(() => {

            //   }, 2000 * index)
            // })
            notificationList.forEach((item: GlobalMessageItem, index: number) => {
              // AddGlobalMessages(item)
              // AddGlobalSysMessages(item)
              // router.currentRoute.value.fullPath == '/dashboard' && setGlobalMessageList()
              setTimeout(() => {
                ElNotification({
                  dangerouslyUseHTMLString: true,
                  appendTo: '.vab-app-main',
                  customClass: 'global-msg-notification-important',
                  message: `
                             <div
                               class="item global-msg-notification-td"
                               data-id='${item.id}'
                               style="z-index:${item.id} !important"
                             >
                                <div class="item-left" data-id="${item.id}">
                                    <div class="item-left-img" data-id="${item.id}"></div>
                                    <span class="item-left-title" data-id="${item.id}">[${item.title}]</span>
                                    <span class="item-left-tip" data-id="${item.id}">${item.notes}</span>
                                </div>
                                <div class="global-msg-notification-item-btn" data-src="btn" data-type="${item.minorType}" data-id="${item.id}">
                                    <input type="checkbox" data-src="btn" data-id="${item.id}" data-type="${item.minorType}"  id="global-msg-notification-item-input" name="today" value="false">&nbsp;
                                    <sapn style="opacity:0.7" data-src="btn" data-type="${item.minorType}"  data-id="${item.id}">今日内，不在提示</span>
                                </div>
                             </div>
                             `,
                  duration: 0,
                  onClick: handleCallBack,
                })
              }, 2000 * index)
            })
          }
        } catch {
          console.error('数据不是Json')
        }
      }
    }
  )

  const handleLoginCallBack = () => {
    settingsStore.changeMessageVisible(true)
    const node = document.querySelector('.global-login-msg-notification')
    node?.parentNode?.removeChild(node)
  }

  const handleCallBack = async (e?: any) => {
    const id = e.target.dataset.id
    const src = e.target.dataset.src
    if (src != 'btn') {
      id && settingsStore.changeMessageVisible(true, id)
    } else {
      const minorType = e.target.dataset.type
      // todo get
      const { msg } = await GlobalMsgSwitchIdApi({
        minorType,
        enable: true,
      })
      $baseMessage(msg, 'success')
    }
    const nodes = document.querySelectorAll('.global-msg-notification')
    nodes.forEach((item, index) => {
      const node = nodes[index].querySelector('.global-msg-notification-td')
      if (node?.getAttribute('data-id') == id) {
        nodes[index]?.parentNode?.removeChild(nodes[index])
      }
    })
    const nodes1 = document.querySelectorAll('.global-msg-notification-important')
    nodes1.forEach((item, index) => {
      const node = nodes1[index].querySelector('.global-msg-notification-td')
      if (node?.getAttribute('data-id') == id) {
        nodes1[index]?.parentNode?.removeChild(nodes1[index])
      }
    })
  }
</script>
<div
  class="item"
  :class="{ itemActive: item.id == curIndex }"
  v-for="item in globalMessageList"
  :key="item.id"
  @click="handleItemClick(item.id)"
>

</div>
<template>
  <div class="vab-app-main" :class="{ active: activeRef }">
    <div v-if="showIcon" class="collapse-btn" :class="{ 'is-collapse': collapse }" @click="toggleCollapse">
      <el-icon>
        <ArrowLeftBold />
      </el-icon>
    </div>
    <section>
      <VabRouterView />
    </section>
    <!-- <vab-footer /> -->
  </div>
</template>
<style lang="scss">
  .global-msg-notification-important {
    width: calc(100vw);
    height: 52px;
    transition: all 0.5s;
    padding: 0px 40px 0px 20px !important;
    display: flex;
    align-items: center;
    border-radius: 0px !important;
    border: 0 solid #fff !important;
    background: #ffefef;
    top: 0 !important;
    right: 0 !important;
    box-shadow: none !important;
    .global-msg-notification-td {
      display: flex;
      align-items: center;
      justify-content: space-between;
    }
    .item-left {
      display: flex;
      align-items: center;
      justify-content: space-between;
    }
    .item-left-img {
      width: 16px;
      height: 20px;
      background-image: url('@/assets/message_image/lingdang.png');
      background-size: 100%;
    }
    .el-notification__group {
      margin: 0px;
      width: 100%;
    }
    .item-left-title {
      font-weight: 500;
      font-size: 16px;
      color: #2b2742;
      margin: 0 8px;
    }
    .item-left-tip {
      font-size: 15px;
      color: #2b2742;
    }
    .global-msg-notification-item-btn {
      display: flex;
      height: 18px;
      align-items: center;
      font-size: 13px;
      color: #4a4759;
      line-height: 18px;
      cursor: pointer;
    }
  }
  .global-msg-notification {
    padding: 20px !important;
    width: 330px;
    box-shadow: 0px 0px 4px 4px rgba(77, 81, 86, 0.05);
    border-radius: 16px;
    border: 1px solid #f1f0ff;
    .item_right_count {
      margin: 0 3px;
      font-weight: 700;
      font-size: 14px;
      color: #ff6767;
    }
    .el-notification__group {
      margin: 0px;
      width: 100%;
    }
    .global-msg-notification-item-btn {
      margin-top: 10px;
      display: flex;
      height: 18px;
      align-items: center;
      font-size: 13px;
      color: #4a4759;
      line-height: 18px;
    }
    .global-msg-notification-item {
      width: 330px;
      border-radius: 10px;
      // display: flex;
      // align-items: center;
      cursor: pointer;
      .item_top {
        display: flex;
        align-items: center;
      }
      .item_left {
        padding: 0px;
        display: flex;
        height: 48px;
        width: 58px;
        padding-right: 10px;
        align-items: center;
        justify-content: center;
        overflow: hidden;
        position: relative;
        .mark {
          position: absolute;
          width: 10px;
          height: 10px;
          background: #ff6767;
          border-radius: 50%;
          right: 9px;
          top: 0px;
        }
        .img {
          width: 48px;
          height: 48px;
          background: #ebeaff;
          border-radius: 16px;
        }
        .user-img {
          background-image: url('@/assets/message_image/user-message.svg');
        }
        .sys-img {
          background-image: url('@/assets/message_image//system-message.svg');
        }
      }
      .item_right {
        width: 220px;
        height: 48px;
        display: flex;
        // flex-direction: column;
        // padding-right: 12px;
        // justify-content: space-between;
        .item_right-top {
          display: flex;
          align-items: left;
          justify-content: space-between;
          flex-direction: column;
        }
        .item_right_content {
          font-size: 13px;
          color: #606266;
          width: 230px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          word-break: break-all;
          word-wrap: break-word;
        }
        .item_right_name {
          font-size: 14px;
          color: #4a4759;
          font-weight: 500;
          width: 200px;
          // width: 100%;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          word-break: break-all;
          word-wrap: break-word;
        }
        .item_right_time {
          width: calc(100% - 80px);
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          word-break: break-all;
          word-wrap: break-word;
          font-size: 12px;
          color: #b7b5bf;
          text-align: right;
        }
      }
      .item-bottom {
        margin-top: 10px;
        width: 290px;
        font-size: 14px;
        color: #4a4759;
        line-height: 24px;
        overflow: hidden;
        max-height: 48px;
        text-overflow: ellipsis;
        word-break: break-all;
        word-wrap: break-word;
        display: -webkit-box;
        // background-clip: text;
        appearance: none;
        line-clamp: 2;
        -webkit-line-clamp: 2;
        /* autoprefixer: ignore next */
        -webkit-box-orient: vertical;
      }
    }
  }
</style>
<style lang="scss" scoped>
  .vab-app-main {
    transition: all 0.4s cubic-bezier(0.645, 0.045, 0.355, 1), background 0s, background-color 0s, border 0s, color 0.1s,
      font-size 0s !important;
  }
  .vab-app-main.active {
    padding: 0 !important;
    transition: none;
  }
  .collapse-btn {
    position: fixed;
    top: 50%;
    color: #fff;
    margin-top: -25px;
    z-index: 10;
    left: 248px;
    width: 18px;
    height: 34px;
    background: #5b5290;
    border-radius: 4px 0px 0px 4px;
    transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1), border 0s, color 0.1s, font-size 0s;
    cursor: pointer;
    :deep() {
      .el-icon {
        margin-top: 10px;
        margin-left: 2px;
      }
    }
    &.is-collapse {
      left: 64px;
      transform: rotateY(180deg);
      background: #e1dcff;
      color: var(--el-color-primary);
    }
  }
</style>
