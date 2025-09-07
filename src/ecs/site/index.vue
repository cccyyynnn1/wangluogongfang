<script setup lang="ts">
  import EditSite from './edit-site.vue'
  import SiteSession from './site-session/index.vue'
  import ApiInterface from '@/ecs/site/interface/api-edit.vue'
  import TranspondData from './components/transpond-data.vue'
  // import ApiList from './api-list/index.vue'
  import ApiList from './api-list/api-list.vue'
  import SiteAppInteraction from './site-app-interaciton/index.vue'

  import { siteGetAllPageApi, getApiInterfaceApiByIdApi } from '~/src/api-ecs/site'
  import { siteCurrentType, siteApiInterfaceItemType } from '@/types'
  // @ts-ignore
  import inputSearch from '~/library/components/VabInputSearch'

  const router = useRouter()
  const { query } = useRoute()
  const activeName = ref('siteSession') // tabs选中项
  const activeData = reactive<{
    sessionId: string | number
    apiId: string | number
    hosts: string | null
    siteName: string | null
  }>({
    sessionId: '',
    apiId: '',
    hosts: null,
    siteName: '',
  })

  const title = ref('这里是站点名称')
  const sonTitle = ref('')
  let flag = false
  const siteSessionRef = ref()
  const mode = ref('') // 编辑站点还是新增站点
  const isFullscreen = ref(false) // 是否全屏
  const defaultPanel = ref<Element>()
  const heightStyle = { width: '' }
  const loading = ref(false) // 站点加载
  const showEditSite = ref(false) // 新增和编辑站点页面
  const showTranspond = ref(false) // 显示转发页面
  const siteList = ref<any[]>([])
  const siteSearch = ref()
  const moreNode = ref<Map<string, { status: boolean; list: siteApiInterfaceItemType[] }>>(new Map())
  const curSite = ref()
  const curApi = ref()
  const highlightId = ref('')
  const showApiInterface = ref(false)
  const apiInterface = ref()
  // 全屏设高
  const setHeight = () => {
    const res = defaultPanel.value?.getBoundingClientRect().width as number
    heightStyle.width = `${res}px`
  }

  // 获取站点列表
  const getSiteList = async () => {
    loading.value = true
    try {
      const { data } = await siteGetAllPageApi({ name: siteSearch.value })
      siteList.value = data
      if (!activeData.sessionId) {
        const curHost = activeData.hosts ? data.find((item: any) => item.id === curSite.value) : data[0]
        activeData.sessionId = curHost?.id
        curSite.value = curHost?.id
        activeData.hosts = curHost?.hosts
        activeData.siteName = curHost?.siteName
        currentItem.id = activeData.sessionId as string
        highlightId.value = activeData.apiId ? `api${activeData.apiId}` : `site${activeData.sessionId}`
      }
      if (!activeData.sessionId || flag) {
        const curHost = activeData.hosts ? data.find((item: any) => item.id === curSite.value) : data[0]
        activeData.hosts = curHost?.hosts
      }
      setTitle()
    } finally {
      loading.value = false
    }
  }

  const reflash = () => {
    flag = true
    getSiteList()
    loadNode(activeData.sessionId)
    siteSessionRef.value.siteSessionSearch()
  }
  // 全屏的开关
  const clickFullScreen = () => {
    isFullscreen.value = !isFullscreen.value
  }

  // 传给编辑页面的值
  const currentItem: siteCurrentType = reactive({
    id: 0,
    siteName: '',
    createTime: 0,
    createUser: 0,
    displayFields: '',
    hosts: '',
    module: '',
    phone: '',
    searchSql: '',
    updateTime: 0,
    updateUser: 0,
    user: '',
  })

  const targetNode = ref()

  // 设置标题
  const setTitle = () => {
    const { sessionId, apiId } = activeData
    if (!sessionId) return
    const session = siteList.value.find((i) => i.id == sessionId)
    const api = getSubData(sessionId)?.list.find((i) => i.id == apiId)
    title.value = `${session ? session.siteName : ''} ${api ? `- ${api.apiName} ` : ''}`
  }

  // 切换右侧站点
  const changeActive = (site: any) => {
    // if (site.id == activeData.sessionId) return false
    activeData.sessionId = site.id
    curSite.value = site.id
    currentItem.id = site.id
    activeData.hosts = site.hosts
    activeData.siteName = site.siteName
    activeData.apiId = ''
    moreNode.value.clear()
    highlightId.value = `site${site.id}`
    setTitle()
  }

  // 切换右侧站点
  const changeSubActive = (node: siteApiInterfaceItemType) => {
    if (node.id == activeData.apiId) return false
    activeData.apiId = node.id
    curApi.value = node.id
    activeData.sessionId = node.sessionId
    curSite.value = node.sessionId
    highlightId.value = `api${node.id}`
    setTitle()
  }

  const displayFieldsArr = ref()

  // 编辑和新增站点
  const addSite = async (val: boolean, res?: string) => {
    if (activeData.apiId) {
      const api = getSubData(activeData.sessionId)?.list.find((i) => i.id == activeData.apiId)
      apiInterface.value = api
      apiInterface.value.displayFieldsArr = displayFieldsArr.value
      return (showApiInterface.value = true)
    }
    mode.value = res || ''
    if (res === 'edit' && siteList.value.length === 0) return
    res === 'edit' ? (sonTitle.value = '编辑站点') : (sonTitle.value = '新增站点')
    showEditSite.value = val
  }

  const handleSetDisplayFieldsArr = (val: any) => {
    displayFieldsArr.value = val
  }

  // 跳转到资产站点
  const skipAssetsSite = () => {
    router.push({ path: '/assets/index', query: { params: 'site' } })
  }

  //子节点数据
  const getSubData = (key: number | string) => {
    return moreNode.value.get(key.toString())
  }

  // 转发
  const transpond = (val: boolean) => {
    showTranspond.value = val
  }

  // 滚动事件
  const scrollEvent = (id: number | number) => {
    const node = document.querySelector('#my-scroll-content')
    // @ts-ignore
    const nodes = Array.from(node?.children)?.filter((item: any) => {
      return item.className == 'text-item active'
    })
    targetNode.value = nodes[nodes.length - 1 || 0]

    if (targetNode.value?.offsetTop) {
      node!.scrollTop = targetNode.value?.offsetTop - 80
    }
  }

  const loadNode = async (id: number | string) => {
    const status = getSubData(id.toString())?.status
    if (status) return (getSubData(id.toString())!.status = false)
    loading.value = true
    const { data } = await getApiInterfaceApiByIdApi(id.toString())
    moreNode.value.clear()
    moreNode.value.set(id.toString(), { status: true, list: data })
    const crr = data?.[0]
    if (crr?.isDisplay === 0 && crr?.id === activeData?.apiId) {
      activeData.apiId = ''
      highlightId.value = `site${activeData.sessionId}`
    }
    loading.value = false
    setTitle()
  }

  const reflashApiInterface = () => {
    getSubData(activeData.sessionId.toString())!.status = false
    loadNode(activeData.sessionId)
    siteSessionRef.value.siteSessionSearch()
  }

  watch(
    () => query,
    () => {
      activeData.hosts = query.hosts?.toString() || ''
      activeData.apiId = query.apiId?.toString() || ''
      activeData.siteName = query.siteName?.toString() || ''
      curApi.value = query.apiId?.toString() || ''
      activeData.sessionId = query.sessionId?.toString() || ''
      curSite.value = activeData.sessionId
      highlightId.value = activeData.apiId ? `api${activeData.apiId}` : `site${activeData.sessionId}`
      if (query.hosts) activeName.value = 'apiList'
      if (query.sessionId) loadNode(query.sessionId as string)
    },
    {
      immediate: true,
    }
  )

  const unfold = async (id: number | string) => {
    const { data } = await getApiInterfaceApiByIdApi(id.toString())
    moreNode.value.set(id.toString(), { status: true, list: data })
    setTitle()
  }

  // 只修改右侧树高亮及滚动
  const handleChangeActive = async (data: any) => {
    curApi.value = data.apiId
    curSite.value = data.siteId
    await unfold(curSite.value)
    scrollEvent(curSite.value)
  }
  const route = useRoute()

  onMounted(() => {
    getSiteList()
    setHeight()
    if (route.query.activeName) {
      activeName.value = route.query.activeName as string
    }
  })

  const formatString = (str: string) => {
    // eslint-disable-next-line
    const len = str.replace(/[^\x00-\xff]/g, '  ').length
    let res = ''
    if (len > 24) {
      if (len - str.length > 10) {
        res = str.slice(0, 13)
      } else {
        res = str.slice(0, 21)
      }
    } else {
      res = str
    }
    return res
  }

  watch(
    () => isFullscreen.value,
    () => {
      const node = document.querySelector('.site .el-tabs__nav')
      if (!node) return
      isFullscreen.value ? node.classList.add('isHide') : node.classList.remove('isHide')
    }
  )
</script>
<script lang="ts">
  export default {
    name: 'Sites',
  }
</script>
<template>
  <div class="site">
    <div class="left">
      <div class="left-bottom">
        <div class="box-card" style="width: 260px">
          <input-search
            v-model="siteSearch"
            clearable
            placeholder="请输入站点名"
            @keydown.enter="getSiteList"
            @onSearch="getSiteList"
          />
          <div id="my-scroll-content" v-loading="loading" class="content">
            <div
              v-for="site in siteList"
              :key="site.id"
              class="text-item"
              :class="{ active: highlightId == `site${site.id}` }"
              @click="changeActive(site)"
            >
              <div class="nodeIcon" :class="{ active: getSubData(site.id)?.status }">
                <el-icon v-if="site.hasChildren" @click.stop="loadNode(site.id)"><CaretRight /></el-icon>
              </div>
              <el-Popover
                class="content-item-tooltips"
                :content="site.siteName"
                effect="dark"
                :hide-after="0"
                placement="right"
                trigger="hover"
                :width="200"
              >
                <template #reference>
                  <div class="api-site">{{ formatString(site.siteName) }}</div>
                </template>
              </el-Popover>
              <template v-if="getSubData(site.id)?.status">
                <template v-for="item in getSubData(site.id)?.list">
                  <div
                    v-if="item.isDisplay === 1"
                    :key="item.id"
                    class="text-sub-item"
                    :class="{ active: highlightId == `api${item.id}` }"
                    @click.stop="changeSubActive(item)"
                  >
                    <el-Popover
                      class="content-item-tooltips"
                      :content="item.apiName"
                      effect="dark"
                      :hide-after="0"
                      placement="right"
                      trigger="hover"
                      :width="200"
                    >
                      <template #reference>
                        <div class="api-site">{{ item.apiName }}</div>
                      </template>
                    </el-Popover>
                  </div>
                </template>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="right">
      <!-- <div class="right-top">
       <div class="text">{{ title }}</div> 
      </div> -->
      <div ref="defaultPanel" class="default-panel" :class="{ 'vab-fullscreen': isFullscreen }" :style="heightStyle">
        <el-tabs v-model="activeName">
          <el-tab-pane label="会话" name="siteSession">
            <site-session
              v-if="activeName == 'siteSession'"
              ref="siteSessionRef"
              :api-id="activeData.apiId"
              :site-session-id="activeData.sessionId"
              @change-active="handleChangeActive"
              @set-display-fields-arr="handleSetDisplayFieldsArr"
            />
          </el-tab-pane>
          <el-tab-pane label="API列表" name="apiList">
            <api-list v-if="activeName == 'apiList'" :hosts="activeData.hosts" :site-name="activeData.siteName" />
          </el-tab-pane>
          <el-tab-pane label="站点应用交互" name="siteAppInteraction">
            <site-app-interaction v-if="activeName == 'siteAppInteraction'" :site-id="activeData.sessionId" />
          </el-tab-pane>
        </el-tabs>
        <div class="tool-bar">
          <div class="tool-bar-left"></div>

          <div class="tool-bar-right">
            <!-- <el-tooltip class="item" content="帮助中心" effect="dark" placement="top">
              <vab-icon class="icon" icon="question-line" />
            </el-tooltip> -->
            <!-- <el-tooltip class="item" content="回溯记录" effect="dark" placement="top">
              <vab-icon class="icon" icon="time-line" @click="siteSessionRef.showHistory = true" />
            </el-tooltip>
            <el-tooltip class="item" content="转发" effect="dark" placement="top">
              <vab-icon v-permissions="['Admin']" class="icon" icon="share-box-line" @click="transpond(true)" />
            </el-tooltip> -->
            <el-tooltip
              v-if="activeName == 'apiList'"
              class="item"
              :content="isFullscreen ? '关闭全屏' : '全屏'"
              effect="dark"
              placement="top"
            >
              <vab-icon
                class="icon"
                :icon="isFullscreen ? 'fullscreen-exit-fill' : 'fullscreen-fill'"
                @click="clickFullScreen"
              />
            </el-tooltip>
            <el-tooltip
              v-if="activeName == 'siteSession'"
              class="item"
              content="编辑站点"
              effect="dark"
              placement="top"
            >
              <vab-icon v-permissions="['Admin']" class="icon" icon="edit-2-line" @click="addSite(true, 'edit')" />
            </el-tooltip>
          </div>
        </div>
      </div>
    </div>

    <EditSite
      :current-item="currentItem"
      :mode="mode"
      :show-edit-site="showEditSite"
      source="site"
      :title="sonTitle"
      @on-closeEvent="addSite"
      @on-reflash="reflash"
    />
    <api-interface
      v-if="showApiInterface"
      :api-data="apiInterface"
      :session-id="+activeData.sessionId"
      :show-api-edit="showApiInterface"
      @on-closeEvent="showApiInterface = false"
      @on-reflash="reflashApiInterface"
    />

    <transpond-data
      v-if="showTranspond"
      :show-transpond="showTranspond"
      :site-session-id="activeData.sessionId"
      @on-close-event="transpond"
    />
  </div>
</template>
<style scoped lang="scss">
  .site {
    border-left: 1px solid rgba(238, 238, 240, 0.13);
    height: calc(100vh);
    background: #090123;
    overflow: hidden;
    width: 100%;
    // background: #f6f8f9;
    display: flex;
    :deep() {
      // 隐藏tab下划线
      .el-tabs__nav-wrap {
        &::after {
          display: none;
        }
      }
      .isHide {
        visibility: hidden;
      }
      .api-site {
        max-width: 180px;
        display: inline-block;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
        &:hover {
          color: rgba(102, 85, 231, 0.8) !important;
        }
      }
      .tooltip-base-box .box-item {
        width: 110px;
        margin-top: 10px;
      }
      .el-tabs__header {
        margin: 0 !important;
      }
    }

    .left {
      height: inherit;
      width: 260px;
      background: #090123;
      height: 100%;
      :deep() {
        .el-input__wrapper {
          background-color: transparent !important;
          background: #15112f;
          border-radius: 6px;
          box-shadow: 0 0 0 1px rgba(181, 178, 195, 1);
          svg {
            cursor: pointer;
          }
          .active {
            color: #6655e7 !important;
            font-size: 13px;
            g {
              stroke: #6655e7;
            }
          }
        }
        .is-focus {
          box-shadow: 0 0 0 1px #b9aff4 !important;
        }
      }
      .left-top {
        margin-bottom: 20px;
        background: #090123;
        font-size: 16px;
        font-weight: 500;
        display: flex;
        justify-content: space-between;
      }

      .left-bottom {
        height: calc(100% - 5px);
      }

      .box-card {
        height: 100%;
        background: #090123;
        padding: 20px;
      }

      .icon {
        &:hover {
          cursor: pointer;
        }
      }

      .content {
        height: calc(100vh - 80px);
        margin-top: 15px;
        overflow-y: auto;
        :deep() {
          .el-loading-mask {
            background: rgba(0, 0, 0, 0.1);
          }
        }
        color: #dcdae9;
        &::-webkit-scrollbar {
          width: 0px;
          height: 12px !important;
        }
        &::-webkit-scrollbar-thumb {
          background-color: rgba(255, 255, 255, 0.2);
        }
        &:hover {
          &::-webkit-scrollbar {
            width: 6px;
            height: 12px !important;
          }
        }
      }

      .text-item {
        position: relative;
        width: 210px;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
        .nodeIcon {
          width: 10px;
          // height: 34px;
          display: inline-block;
          margin-right: 10px;
          vertical-align: 3px;
          &.active {
            i {
              transform: rotate(90deg);
            }
          }
        }
        &.active {
          :deep() {
            & > .api-site {
              color: #6954f0 !important;
              font-weight: 500;
            }
          }
        }
        :deep() {
          i {
            margin-top: 7px;
            transform: rotate(0);
            transition: all cubic-bezier(0.19, 1, 0.22, 1);
          }
          span {
            text-indent: 10px;
          }
        }

        &:hover {
          cursor: pointer;
        }
      }
      .text-sub-item {
        padding-left: 35px;
        height: 28px;
        .api-site {
          max-width: 180px;
        }
        &.active {
          color: #6954f0 !important;
          font-weight: 500;
        }
      }
    }

    .text {
      color: #303338;
    }

    .right {
      height: calc(100%);
      min-width: calc(100% - 260px);
      padding-top: 13px;
      background-color: #fff;
      .default-panel {
        width: 100%;
        margin-bottom: 20px;
        min-width: 1000px;
        position: relative;
        padding: 0 20px;
        background: #ffffff;
        height: calc(100% - 30px);
        :deep(.el-tabs__nav) {
          min-width: 250px;
        }

        .m-2 {
          margin-right: 20px;
        }

        .icon {
          margin-right: 6px;
          color: #c3c3c3;

          &:hover {
            cursor: pointer;
          }
        }

        .tool-bar {
          position: absolute;
          right: 20px;
          top: 0px;
          width: 70%;
          height: 40px;
          display: flex;
          align-items: center;
          justify-content: space-between;
        }
      }

      .tool-bar-left {
        display: flex;
        align-items: center;
      }

      .right-top {
        margin-bottom: 20px;
        font-size: 16px;
        font-weight: 500;
      }
    }
  }
</style>
