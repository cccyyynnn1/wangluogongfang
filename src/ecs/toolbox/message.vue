<script lang="ts">
  export default {
    name: 'Message', // 系统信息
  }
</script>

<script setup lang="ts">
  import inputSearch from '~/library/components/VabInputSearch/index.vue'
  import { onClickOutside } from '@vueuse/core'
  import { usePubilcStore } from '@/store/modules/public'
  import { GlobalMessageItem } from '~/types/store'
  import { useSettingsStore } from '@/store/modules/settings'
  import { formatTime } from '@/utils/time'
  import {
    GetGlobalByIdApi,
    GetMsgTypeArrApi,
    GlobalMsgDelApi,
    GlobalMsgGetGlobalMsgPageApi,
    GlobalMsgUpdateApi,
  } from '~/src/api-ecs/public'
  import { MsgTypeArrType } from '~/src/types'

  const router = useRouter()
  const settingsStore = useSettingsStore()
  const publicStore = usePubilcStore()
  const { updateHasNewMessage, setGlobalMessageList } = publicStore

  const { msgId } = storeToRefs(settingsStore)

  const target = ref(null)
  const minorTypeList = ['sysStatusAlert', 'authExpire', 'upgrade', 'logPcapDownload', 'collectShare']

  interface Props {
    modelValue: boolean
  }
  const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
  })

  let curTotal = 0

  const toolsIndex = ref()
  const newsIndex = ref(1)
  const leftMsgType = ref<'user' | 'sys'>('user')
  const defaultActiveValue = ref('user-1')
  const queryData = reactive<{
    pageNum: number
    pageSize: number
    minorType?: string
    readStatus?: 0 | 1 // 1:已读 0:未读
    sendUserId?: number
  }>({
    pageNum: 1,
    pageSize: 30,
  })

  const emits = defineEmits<{
    (e: 'update:modelValue', visible: boolean): void
  }>()

  const visible = useVModel(props, 'modelValue', emits)

  const filterVisiable = ref(false)
  onClickOutside(target, (event) => (filterVisiable.value = false))

  const searchValue = ref()
  const handleSearch = () => {
    // if (searchValue.value?.trim()) {
    //   curList.value.forEach((_: GlobalMessageItem, index: number) => {
    //     if (
    //       curList.value[index].msgSource.includes(searchValue.value) ||
    //       curList.value[index].content.includes(searchValue.value)
    //     ) {
    //       curList.value[index].isHide = false
    //     } else {
    //       curList.value[index].isHide = true
    //     }
    //     const regExp = new RegExp(searchValue.value, 'g')
    //     const regExp1 = new RegExp('<span style="color: #1E1842;background:#FFED81">', 'g')
    //     const regExp2 = new RegExp('</span>', 'g')
    //     curList.value[index].msgSource = curList.value[index].msgSource.replaceAll(regExp1, '')
    //     curList.value[index].content = curList.value[index].content.replaceAll(regExp1, '')
    //     curList.value[index].msgSource = curList.value[index].msgSource.replaceAll(regExp2, '')
    //     curList.value[index].content = curList.value[index].content.replaceAll(regExp2, '')
    //     curList.value[index].msgSource = curList.value[index].msgSource.replaceAll(
    //       regExp,
    //       `<span style="color: #1E1842;background:#FFED81">${searchValue.value}</span>`
    //     )
    //     curList.value[index].content = curList.value[index].content.replaceAll(
    //       regExp,
    //       `<span style="color: #1E1842;background:#FFED81">${searchValue.value}</span>`
    //     )
    //   })
    // } else {
    //   curList.value.forEach((_: GlobalMessageItem, index: number) => {
    //     const regExp1 = new RegExp('<span style="color: #1E1842;background:#FFED81">', 'g')
    //     const regExp2 = new RegExp('</span>', 'g')
    //     curList.value[index].msgSource = curList.value[index].msgSource.replaceAll(regExp1, '')
    //     curList.value[index].content = curList.value[index].content.replaceAll(regExp1, '')
    //     curList.value[index].msgSource = curList.value[index].msgSource.replaceAll(regExp2, '')
    //     curList.value[index].content = curList.value[index].content.replaceAll(regExp2, '')
    //     curList.value[index].isHide = false
    //   })
    // }
    queryData.pageNum = 1
    curList.value = []
    getCurData()
  }

  const handleClose = () => {
    emits('update:modelValue', false)
  }

  const elData = ref()
  const el = ref<HTMLElement | null>(null)
  const curList = ref<GlobalMessageItem[]>([])
  onMounted(async () => {
    loading.value = true
    await GetMsgTypeArr()
    await getCurItem()
  })

  watch(
    () => elData.value,
    () => {
      if (elData.value.bottom) {
        if (curTotal > curList.value.length) {
          getCurData()
        }
      }
    },
    { deep: true }
  )

  const msgTypeArr = ref<{ sys: MsgTypeArrType[]; user: MsgTypeArrType[] }>({
    sys: [],
    user: [],
  })
  const GetMsgTypeArr = async () => {
    const { data } = await GetMsgTypeArrApi()
    msgTypeArr.value = data
  }

  const getCurItem = async () => {
    if (msgId.value) {
      const { data } = await GetGlobalByIdApi({ id: msgId.value })
      if (data) {
        const { minorType, msgType, sendUserId } = data
        let leftTd = undefined
        leftMsgType.value = msgType
        let index = 0
        if (sendUserId) {
          index = msgTypeArr.value[leftMsgType.value].findIndex((element) => {
            return element.sendUserId == sendUserId
          })
        } else {
          index = msgTypeArr.value[leftMsgType.value].findIndex((element) => {
            return element.minorType.includes(minorType)
          })
        }
        leftTd = msgTypeArr.value[leftMsgType.value][index]
        defaultActiveValue.value = `${msgType}-${index + 1}`
        queryData.pageNum = 1
        queryData.sendUserId = leftTd?.sendUserId
        queryData.minorType = leftTd?.minorType
        toolsIndex.value = leftTd?.groupName
        curList.value = []
        await getCurData()
        handleItemClick(msgId.value)
        nextTick(() => {
          roll(`msg-scoll-${msgId.value}`)
          el.value = document.querySelector('.global-msg-warp-content')
          const { arrivedState } = useScroll(el.value)
          elData.value = arrivedState
        })
      }
    } else {
      const { minorType, sendUserId, groupName } = msgTypeArr.value[leftMsgType.value][0]
      queryData.pageNum = 1
      queryData.sendUserId = sendUserId
      queryData.minorType = minorType
      toolsIndex.value = groupName
      curList.value = []
      await getCurData()
      nextTick(() => {
        el.value = document.querySelector('.global-msg-warp-content')
        const { arrivedState } = useScroll(el.value)
        elData.value = arrivedState
      })
    }
  }

  //点击便签滚动
  const roll = (id: string) => {
    if (id) {
      document.getElementById(id)
      document.getElementById(id)!.scrollIntoView({
        behavior: 'smooth', // 平滑过渡
        block: 'start', // 上边框与视窗顶部平齐。默认值
      })
    }
  }

  const loading = ref(false)
  const getCurData = async () => {
    try {
      loading.value = true
      const { data } = await GlobalMsgGetGlobalMsgPageApi({
        pageNum: queryData.pageNum++,
        pageSize: queryData.pageSize,
        sendUserId: queryData.sendUserId,
        minorType: queryData.minorType,
        readStatus: newsIndex.value == 0 ? 0 : undefined,
        searchStr: searchValue.value,
      })
      curTotal = data.total
      curList.value = [...curList.value, ...data.records]
    } finally {
      loading.value = false
    }
  }
  watch(
    () => newsIndex.value, // 全部消息/未读消息
    () => {
      queryData.pageNum = 1
      curList.value = []
      getCurData()
    }
  )

  const handleItemClick = async (id: number) => {
    let cur = -1
    cur = curList.value.findIndex((item) => {
      return item.id == id
    })
    if (curList.value[cur]?.isFold) return
    curList.value.forEach((item) => {
      if (item.id != id) {
        item.isFold = false
      }
    })
    curList.value[cur]['isFold'] = true
    if (!curList.value[cur].readStatus) {
      const { data } = await GlobalMsgUpdateApi({
        minorType: queryData.minorType,
        sendUserId: queryData.sendUserId,
        id: id,
      })
      msgTypeArr.value[leftMsgType.value].forEach((_: MsgTypeArrType, index: number) => {
        if (msgTypeArr.value[leftMsgType.value][index].groupName == toolsIndex.value) {
          msgTypeArr.value[leftMsgType.value][index].count = data.curTypeUnReadCount
        }
      })
      updateHasNewMessage(!!data.unReadCount)
      router.currentRoute.value.fullPath == '/dashboard' && setGlobalMessageList()
    }
    curList.value[cur].readStatus = true
  }

  const handleItemFold = (id: number) => {
    let cur = -1
    cur = curList.value.findIndex((item) => {
      return item.id == id
    })
    curList.value[cur].isFold = false
  }

  const formatNumber = (num: number) => {
    let res: number | null | string = null
    if (num > 99) {
      res = '99+'
    } else {
      res = num == 0 ? null : num
    }
    return res
  }

  const handleReadAll = async () => {
    const flag = curList.value.some((item) => {
      return item.readStatus == false
    })
    if (flag) {
      curList.value.forEach((_: GlobalMessageItem, index: number) => {
        curList.value[index].readStatus = true
      })
      const { data } = await GlobalMsgUpdateApi({
        minorType: queryData.minorType,
        sendUserId: queryData.sendUserId,
        containsAll: true,
      })
      msgTypeArr.value[leftMsgType.value].forEach((_: MsgTypeArrType, index: number) => {
        if (msgTypeArr.value[leftMsgType.value][index].groupName == toolsIndex.value) {
          msgTypeArr.value[leftMsgType.value][index].count = data.curTypeUnReadCount
        }
      })
      updateHasNewMessage(!!data.unReadCount)
      router.currentRoute.value.fullPath == '/dashboard' && setGlobalMessageList()
    }
  }
  const $baseConfirm: any = inject('$baseConfirm')
  const handleReset = () => {
    if (curList.value.length == 0) return
    $baseConfirm(`你确定要清空${toolsIndex.value}中的消息吗？`, null, async () => {
      curList.value = []
      curTotal = 0
      const { msg, data } = await GlobalMsgDelApi({
        minorType: queryData.minorType,
        sendUserId: queryData.sendUserId,
        containsAll: true,
      })
      updateHasNewMessage(!!data.unReadCount)
      msgTypeArr.value[leftMsgType.value].forEach((_: MsgTypeArrType, index: number) => {
        if (msgTypeArr.value[leftMsgType.value][index].groupName == toolsIndex.value) {
          msgTypeArr.value[leftMsgType.value][index].count = data.curTypeUnReadCount
        }
      })
      router.currentRoute.value.fullPath == '/dashboard' && setGlobalMessageList()
      ElMessage({ message: msg, type: 'success' })
    })
  }
  const handleDel = async (id: number) => {
    $baseConfirm(`你确定要删除此条信息吗？`, null, async () => {
      // curList.value = curList.value.filter((_: GlobalMessageItem, index: number) => {
      //   return curList.value[index].id !== id
      // })
      const { msg, data } = await GlobalMsgDelApi({
        minorType: queryData.minorType,
        sendUserId: queryData.sendUserId,
        id: id,
      })
      updateHasNewMessage(!!data.unReadCount)
      msgTypeArr.value[leftMsgType.value].forEach((_: MsgTypeArrType, index: number) => {
        if (msgTypeArr.value[leftMsgType.value][index].groupName == toolsIndex.value) {
          msgTypeArr.value[leftMsgType.value][index].count = data.curTypeUnReadCount
        }
      })
      ElMessage({ message: msg, type: 'success' })
      router.currentRoute.value.fullPath == '/dashboard' && setGlobalMessageList()
      queryData.pageNum = 1
      curList.value = []
      getCurData()
    })
  }
  const handleLeftItemClick = async (val: string, type: 'sys' | 'user') => {
    leftMsgType.value = type
    const item = msgTypeArr.value[leftMsgType.value].find((item: MsgTypeArrType) => {
      return val == item.groupName
    })
    toolsIndex.value = val
    queryData.pageNum = 1
    queryData.minorType = item?.minorType
    queryData.sendUserId = item?.sendUserId
    curList.value = []
    getCurData()
  }

  const logPcapDownloadEvent = (val: any) => {
    const { filePath } = JSON.parse(val)
    settingsStore.changeDownloadVisible(true, false, filePath)
  }

  const collectShareEvent = (val: any) => {
    try {
      const { workspaceId, indexType, searchStTime, searchEdTime, searchSql, inputSql, sqlRelat, filterSqlArr } =
        JSON.parse(val)
      const obj: {
        workspaceId: number
        indexType: number
        timeRanges: [searchStTime: string, searchEdTime: string]
        sqlRelat: string
        inputSql: string
        filterSqlArr: string
        searchSql: string
      } = {
        workspaceId,
        indexType,
        timeRanges: [searchStTime, searchEdTime],
        searchSql,
        inputSql,
        sqlRelat,
        filterSqlArr,
      }
      window.open(
        `${window.location.origin}/#/retrieve/index?params=applicationLayer&infos=${encodeURIComponent(
          JSON.stringify(obj)
        )}`
      )
    } catch (error) {
      console.log(error)
    }
  }

  const handleToDEtail = (item: GlobalMessageItem) => {
    switch (item.minorType) {
      case 'sysStatusAlert':
        router.push({ name: 'SystemCondition' })
        break
      case 'authExpire':
        router.push({ name: 'AuthorizationManagement' })
        break
      case 'upgrade':
        router.push({ name: 'FullFlowConfig' })
        break
      case 'logPcapDownload':
        logPcapDownloadEvent(item.content)
        break
      case 'collectShare':
        collectShareEvent(item.content)

        break

      default:
        break
    }
    settingsStore.changeMessageVisible(false)
  }

  const enterPressHadlerSearch = (e: any) => {
    if (e.keyCode === 13) {
      handleSearch()
    }
  }

  const formartContent = (val: string, minorType: string) => {
    if (val) {
      if (minorType === 'collectShare') {
        const str = JSON.parse(val)
        const { dataType, inputSql, filterSqlArr, searchStTime, searchEdTime, remarks } = str
        const res = `检索类型:  ${dataType}<br/>检索条件:  ${inputSql || '空'}<br/>过滤条件:  ${
          filterSqlArr || '空'
        }<br />检索时间: ${searchStTime} - ${searchEdTime}`
        return remarks ? `${res}<br />分享备注: ${remarks}` : res
      } else if (minorType === 'logPcapDownload') {
        let res
        const str = JSON.parse(val)
        const { filePath, indexType, downloadCnd } = str
        try {
          const { clientIp, clientPort, serverIp, serverPort, timeRange } = downloadCnd && JSON.parse(downloadCnd)
          res = `文件名:  ${filePath}<br/>检索类型:  ${indexType}<br/>查询条件:  源Ip = ${
            serverIp || '空'
          }、 目的Ip = ${clientIp || '空'}、 源端口 = ${serverPort || '空'}、  目的端口 = ${
            clientPort || '空'
          }<br/>时间范围: ${timeRange || '空'}`
        } catch (error) {
          res = `检索条件： ${downloadCnd}`
        }
        return res
      } else if (minorType === 'notice') {
        const regExp = new RegExp('《新手操作文档》', 'g')
        const res = val.replaceAll(regExp, `<a href="#/help" target="_blank">《新手操作文档》</a>`)
        return res
      } else {
        return val
      }
    }
  }

  const formartContentTitle = (val: string, minorType: string) => {
    if (val) {
      if (minorType === 'collectShare') {
        // const str = JSON.parse(val)
        // const { dataType } = str
        return ` 给您发送了一条新消息`
      } else if (minorType === 'logPcapDownload') {
        const str = JSON.parse(val) || val
        const { filePath } = str
        const res = `文件名:  ${filePath}`
        return res
      } else if (minorType === 'notice') {
        return ` 您好！非常高兴您的加入。`
      } else {
        return val
      }
    }
  }
</script>

<template>
  <el-dialog v-model="visible" class="message-dialog" width="1080px">
    <template #header></template>
    <template #default>
      <div class="dialog_left">
        <div class="dialog_left_top">消息通知</div>
        <div class="dialog_left_bottom">
          <el-menu :default-active="defaultActiveValue" :default-openeds="['user', 'sys']">
            <el-sub-menu index="user">
              <template #title>
                <span>用户消息</span>
              </template>
              <el-menu-item
                v-for="(item, idx) in msgTypeArr.user"
                :key="idx"
                class="item"
                :index="`user-${idx + 1}`"
                @click="handleLeftItemClick(item.groupName, 'user')"
              >
                <div style="display: flex; align-items: center">
                  <img v-if="item.svg" alt="" :src="require(`@/assets/message_image/${item.svg}`)" />
                  <img v-else alt="" :src="require(`@/assets/message_image/user-message.svg`)" />
                  <span style="font-weight: 500; font-size: 15px; color: #4a4759; margin-left: 10px">
                    {{ item.groupName }}
                  </span>
                </div>
                <div v-if="item.count > 0" class="tips">{{ formatNumber(item.count) }}</div>
              </el-menu-item>
            </el-sub-menu>
            <el-sub-menu index="sys">
              <template #title>
                <span>系统消息</span>
              </template>
              <el-menu-item
                v-for="(item, idx) in msgTypeArr.sys"
                :key="idx"
                class="item"
                :index="`sys-${idx + 1}`"
                @click="handleLeftItemClick(item.groupName, 'sys')"
              >
                <div style="display: flex; align-items: center">
                  <img v-if="item.svg" alt="" :src="require(`@/assets/message_image/${item.svg}`)" />
                  <img v-else alt="" :src="require(`@/assets/message_image/user-message.svg`)" />
                  <span style="font-weight: 500; font-size: 15px; color: #4a4759; margin-left: 10px">
                    {{ item.groupName }}
                  </span>
                </div>
                <div v-if="item.count > 0" class="tips">{{ formatNumber(item.count) }}</div>
              </el-menu-item>
            </el-sub-menu>
          </el-menu>
          <div class="tools">
            <div class="tool_btn" @click="handleReset">一键清空</div>
            <div class="tool_btn" @click="handleReadAll">全部已读</div>
          </div>
        </div>
      </div>
      <div class="dialog_right">
        <div class="dialog_right_top">
          <div style="display: flex; align-items: center">
            <input-search
              v-model="searchValue"
              clearable
              placeholder="搜索"
              style="width: 288px"
              @keydown.enter="enterPressHadlerSearch"
              @on-search="handleSearch"
            />
            <el-popover
              ref="target"
              placement="bottom-start"
              popper-class="message_filter"
              trigger="click"
              :visible="filterVisiable"
              :width="100"
            >
              <template #reference>
                <div class="dialog_right_top_filter" @click="filterVisiable = true">
                  <img alt="" :src="require(`@/assets/message_image/screening.svg`)" />
                </div>
              </template>
              <template #default>
                <div class="filter_item" :class="{ active: newsIndex == 1 }" @click="newsIndex = 1">全部消息</div>
                <div class="filter_item" :class="{ active: newsIndex == 0 }" @click="newsIndex = 0">未读消息</div>
              </template>
            </el-popover>
          </div>
          <div class="dialog_close" @click="handleClose">
            <el-icon><Close /></el-icon>
          </div>
        </div>
        <div v-if="curList?.length > 0" v-loading="loading" class="dialog_right_content global-msg-warp-content">
          <div v-for="item in curList" :key="item.id">
            <div v-if="!item.isHide" :id="'msg-scoll-' + item.id" class="content_item">
              <div
                class="item-top"
                :style="{ cursor: item.isFold ? 'default' : 'pointer' }"
                @click="handleItemClick(item.id)"
              >
                <div class="content_item_left">
                  <div v-if="!item.readStatus" class="remark">&nbsp;&nbsp;</div>
                  <div v-else class="remark_white">&nbsp;&nbsp;</div>
                  <span class="user">
                    <span>[</span>
                    <span v-html="item.msgSource"></span>
                    <span>]</span>
                  </span>
                  <div
                    v-if="!item.isFold"
                    style="font-size: 14px; color: #4a4759"
                    v-html="formartContentTitle(item.content, item.minorType)"
                  ></div>
                </div>
                <div class="content_item_right">
                  <span style="font-size: 12px; color: #9d9baa">{{ formatTime(item.createTime) }}</span>
                  <el-icon v-if="!item.isFold" style="margin: 0 2px 0 6px; color: #6954f0">
                    <ArrowDown />
                  </el-icon>
                  <el-icon
                    v-if="item.isFold"
                    style="margin: 0 2px 0 6px; color: #6954f0; cursor: pointer"
                    @click.stop="handleItemFold(item.id)"
                  >
                    <ArrowUp />
                  </el-icon>
                  <span v-if="!item.isFold" style="font-weight: 500; font-size: 12px; color: #6954f0">展开</span>
                  <span
                    v-if="item.isFold"
                    style="font-weight: 500; font-size: 12px; color: #6954f0; cursor: pointer"
                    @click.stop="handleItemFold(item.id)"
                  >
                    收起
                  </span>
                </div>
              </div>
              <div
                v-if="item.isFold"
                class="content_detail"
                v-html="formartContent(item.content, item.minorType)"
              ></div>
              <div v-if="item.isFold" class="content_del">
                <el-button
                  v-if="minorTypeList.includes(item.minorType)"
                  link
                  style="margin-left: 9px"
                  type="primary"
                  @click="handleToDEtail(item)"
                >
                  查看详情
                </el-button>
                <span v-else></span>
                <div style="cursor: pointer" @click="handleDel(item.id)">
                  <el-icon><Delete /></el-icon>
                </div>
              </div>
            </div>
          </div>
          <el-empty v-if="curList?.length == 0" description="暂无数据" style="height: 80%; width: 100%" />
        </div>
        <div v-else v-loading="loading" class="dialog_right_content">
          <el-empty description="暂无数据" style="height: 80%; width: 100%" />
        </div>
      </div>
    </template>
  </el-dialog>
</template>
<style lang="scss">
  .el-overlay {
    background-color: rgba(0, 0, 0, 0.2) !important;
    backdrop-filter: none !important;
  }
  .message-dialog {
    .el-dialog__header {
      display: none;
    }
    .el-dialog__body {
      padding: 0 !important;
      height: 72vh;
      border-top: 0px !important;
      display: flex;
      // border: 1px solid #5e596d !important;
      border-radius: 10px;
    }
  }
  .message_filter {
    min-width: 100px !important;
    padding: 5px !important;
    .active {
      background: #f3f1fe;
      color: #5236ff;
    }
    .filter_item {
      width: 90px;
      height: 36px;
      border-radius: 4px;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      // &:hover {
      //   background: #f3f1fe;
      //   color: #5236ff;
      // }
    }
  }
</style>
<style scoped lang="scss">
  .dialog_left {
    padding: 0 20px;
    width: 280px;
    height: 100%;
    border-right: 1px solid #f1f0ff;
    :deep() {
      .dialog_left_bottom {
        height: calc(100% - 66px);
        padding-bottom: 20px;
        display: flex;
        justify-content: space-between;
        flex-direction: column;
        .el-sub-menu__title {
          margin-bottom: 1px;
          height: 42px;
          background: #f8f7ff;
          border-radius: 8px !important;
          font-weight: 500;
          font-size: 14px;
          color: #918da5;
        }
        .el-menu {
          height: calc(100% - 40px);
          overflow: hidden;
          overflow-y: auto;
          &::-webkit-scrollbar {
            width: 0;
            height: 0;
          }

          .el-sub-menu {
            margin-bottom: 6px;
          }
        }
        .active {
          color: #593eff;

          // ::after {
          //   position: absolute;
          //   content: '';
          //   width: 7px;
          //   height: 7px;
          //   background: #ff6767;
          //   border-radius: 50%;
          //   left: 46px;
          //   top: 12px;
          // }
        }
        .is-active.el-menu-item {
          span {
            color: #593eff !important;
          }
        }
        .item {
          position: relative;
          width: 100%;
          padding: 13px;
          height: 52px;
          border-radius: 8px;
          display: flex;
          align-items: center;
          justify-content: space-between;
          cursor: pointer;
          &:hover {
            background-color: rgba(240, 238, 253, 0.7);
          }
          img {
            width: 32px;
            height: 32px;
          }
          .tips {
            width: 30px;
            height: 20px;
            background: #ff6667;
            border-radius: 10px;
            color: #ffffff;
            line-height: 20px;
            text-align: center;
          }
        }

        .btns {
          position: relative;
          // height: calc(100% - 40px);
          width: 100%;
          overflow-y: auto;
          overflow-x: hidden;
          &::-webkit-scrollbar {
            width: 0;
            height: 0;
          }
        }
        .tools {
          display: flex;
          align-items: center;
          justify-content: space-between;
          .tool_btn {
            width: 112px;
            height: 34px;
            background: #ffffff;
            border-radius: 6px;
            border: 1px solid #eaeafb;
            display: flex;
            align-items: center;
            justify-content: center;
            cursor: pointer;
            color: #4a4759;
          }
        }
      }
    }
    .dialog_left_top {
      margin-left: 10px;
      font-weight: 500;
      font-size: 20px;
      color: #303133;
      text-align: left;
      height: 66px;
      line-height: 66px;
    }
  }
  .dialog_right {
    width: calc(100% - 280px);
    border-radius: 0 12px 12px 0;
    overflow: hidden;
    // padding: 0 20px;
    .dialog_right_top {
      width: 100%;
      width: 100%;
      height: 66px;
      background-color: #fff;
      border-bottom: 1px solid #f1f0ff;
      padding: 0 20px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      .dialog_right_top_filter {
        width: 32px;
        height: 32px;
        background: #ffffff;
        border-radius: 4px;
        border: 1px solid var(--el-border-color);
        margin-left: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        img {
          width: 13px;
          height: 13px;
          color: #9d9baa;
        }
      }

      .dialog_close {
        width: 26px;
        height: 26px;
        background: #eeedf9;
        display: flex;
        align-items: center;
        justify-content: center;
        float: right;
        border-radius: 50%;
        cursor: pointer;
      }
    }
    .dialog_right_content {
      width: 100%;
      height: calc(100% - 66px);
      background: #fcfcff;
      padding: 15px 20px;
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
      .content_item {
        background-color: #fff;
        border-radius: 8px;
        border: 1px solid #ebe9fd;
        padding: 15px;
        width: 100%;
        margin-bottom: 10px;
      }
      .item-top {
        display: flex;
        align-items: center;
        justify-content: space-between;
        width: 100%;
        .content_item_left {
          width: 70%;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          word-break: break-all;
          word-wrap: break-word;
          text-align: left;
          display: flex;
          align-items: center;
          .remark,
          .remark_white {
            width: 7px;
            height: 7px;
            background: #ff6767;
            border-radius: 50%;
          }
          .remark_white {
            background-color: #fff;
          }

          .user {
            font-weight: 500;
            font-size: 14px;
            color: #4a4759;
            margin: 0 6px;
          }
        }
        .content_item_right {
          display: flex;
          align-items: center;
        }
      }
      .content_detail {
        padding: 0 12px;
        font-size: 14px;
        color: #4a4759;
        line-height: 24px;
        overflow: hidden;
        // text-overflow: ellipsis;
        white-space: wrap;
        word-break: break-all;
        word-wrap: break-word;
        margin: 8px 0;
      }
      .content_del {
        display: flex;
        align-items: center;
        justify-content: space-between;
      }
    }
  }
</style>
