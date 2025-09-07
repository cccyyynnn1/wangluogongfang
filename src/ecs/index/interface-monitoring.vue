<script setup lang="ts">
  import { getAllDisPlaysFiledApi } from '~/src/api-ecs/public'

  import { getSiteIdAppIdApi } from '~/src/api-ecs/site'

  import DesignMonitoring from './components/interface-monitoring/design-monitoring.vue'

  import DetailMonitoring from './components/interface-monitoring/detail-monitoring.vue'

  import { MonitoringType } from '@/types/index'

  import { deleteApiMonitorApi, getApiMonitorPageApi } from '~/src/api-ecs/dashboard'

  import type { TabPaneName, TabsPaneContext } from 'element-plus'

  import { usePubilcStore } from '@/store/modules/public'

  const publicStore = usePubilcStore()

  const { SetFlowProbes } = publicStore

  const $baseConfirm: any = inject('$baseConfirm')

  const $baseMessage: any = inject('$baseMessage')

  const pageName = ref('interface-monitoring') // 显示页面名字

  const model = ref('detail') // 显示页面名字

  const divRef = ref()
  const isBtn = ref(false)
  const siteList = ref<any[]>([])
  interface Obj {
    [key: string]: any
  }
  const curData = ref<Obj>({})

  // 页面数据
  const listData = ref<MonitoringType[]>([])

  const deleteQuery = reactive({
    visible: false,
    top: '0',
    left: '0',
  })

  // 是否加载
  const listLoading = ref(false)

  const deleteData = ref()

  // 表格数据
  const layout = ref('total, sizes, prev, pager, next, jumper')

  const queryPage = reactive({
    total: 0,
    pageNum: 1,
    pageSize: 20,
    title: '',
    loading: false,
  })

  const allField = ref()

  const currentItem = ref()

  let tabIndex = 0
  const editableTabsValue = ref('0')
  const editableTabs = ref<{ title: string; name: string; type: string; model: string; data: any }[]>([])

  const showDetail = (item: MonitoringType) => {
    currentItem.value = item
    pageName.value = 'detail-monitoring'
    model.value = 'detail'
    let flag = false
    if (editableTabs.value.length > 0) {
      flag = editableTabs.value.some((td: any) => {
        return td.title == currentItem.value.name
      })
    }
    if (!flag && editableTabs.value.length == 10) {
      editableTabs.value.shift()
    }
    if (!flag) {
      addTab(currentItem.value.name, 'detail-monitoring', 'detail', item)
    } else {
      const item = editableTabs.value.find((td: any) => {
        return td.title == currentItem.value.name
      })
      if (item) {
        editableTabsValue.value = item.name
        curData.value = item
      }
    }
  }

  const addItemEvent = () => {
    currentItem.value = {}
    isBtn.value = false
    model.value = 'design'
    let flag = false
    if (editableTabs.value.length > 0) {
      flag = editableTabs.value.some((td: any) => {
        return td.title == currentItem.value.name
      })
    }
    if (!flag && editableTabs.value.length == 10) {
      editableTabs.value.shift()
    }
    if (!flag) {
      addTab('新建窗口', 'design-monitoring', 'design')
    }
  }

  // 页容量改变
  const handleSizeChange = async (val: number) => {
    queryPage.pageSize = val
    await getData()
    verifyError()
  }

  // 页面改变
  const handleCurrentChange = async (val: number) => {
    queryPage.pageNum = val
    await getData()
    verifyError()
  }

  // 获取数据
  const getData = async () => {
    listLoading.value = true
    try {
      const { data } = await getApiMonitorPageApi({ pageNum: queryPage.pageNum, pageSize: queryPage.pageSize })
      listData.value = data.records
      queryPage.total = data.total
    } finally {
      listLoading.value = false
    }
  }

  // 改变tab签标题
  const changeTabName = (val: string) => {
    if (editableTabs.value[Number(editableTabsValue.value) - 1]) {
      editableTabs.value[Number(editableTabsValue.value) - 1].title = val
    }
  }

  // 返回
  const backEvent = async () => {
    if (currentItem.value.uuid) {
      pageName.value = 'detail-monitoring'
    }
    if (model.value == 'design') {
      removeTab(curData.value.name)
      editableTabsValue.value = '0'
      curData.value = {}
      tabIndex = editableTabs.value.length
    } else {
      pageName.value = 'detail-monitoring'
    }
    await getData()
    verifyError()
  }

  const detailBackEvent = async () => {
    // pageName.value = 'interface-monitoring'
    await getData()
    verifyError()
  }

  const verifyError = () => {
    // @ts-ignore
    const { siteIds, apiIds } = siteList.value
    listData.value.forEach((item: any) => {
      const list = JSON.parse(item.list)
      list.forEach((td: any) => {
        if (td.data.module == '0') {
          if (!siteIds.includes(td.data.siteSessionId)) {
            item['hasError'] = true
            td['hasError'] = true
          }
          if (td.data.siteApiId && !apiIds.includes(td.data.siteApiId)) {
            item['hasError'] = true
            td['hasError'] = true
          }
        }
      })
      item.list = JSON.stringify(list)
    })
  }

  // 获取站点列表
  const getSiteList = async () => {
    const { data } = await getSiteIdAppIdApi()
    siteList.value = data
    verifyError()
  }

  onMounted(async () => {
    await getData()
    getAllDisPlaysFiled()
    getSiteList()
    SetFlowProbes()
  })

  const toDesignPage = () => {
    isBtn.value = true
    pageName.value = 'design-monitoring'
    model.value = 'design'
  }

  const toDetailPage = async (val: MonitoringType) => {
    Object.keys(val).forEach((item) => {
      // @ts-ignore
      currentItem.value[item] = val[item]
    })
    pageName.value = 'detail-monitoring'
    await getData()
    verifyError()
  }

  // 获取用户所有字段
  const getAllDisPlaysFiled = async () => {
    const { data } = await getAllDisPlaysFiledApi()
    allField.value = data['1']
  }

  provide('allField', allField)

  useEventListener(document, 'click', (evt) => {
    deleteQuery.visible = false
  })

  const addTab = (targetName: string, type: string, model: string, item?: MonitoringType) => {
    const newTabName = `${++tabIndex}`
    editableTabs.value.push({
      title: targetName,
      name: newTabName,
      type,
      model,
      data: item,
    })
    curData.value = editableTabs.value[+newTabName - 1]
    editableTabsValue.value = newTabName
  }

  const removeTab = (name: TabPaneName) => {
    let tabs = [{ title: '仪表盘', name: '0', type: '', model: '', data: undefined }]
    tabs = [...tabs, ...editableTabs.value]
    let activeName = editableTabsValue.value
    if (activeName === name) {
      tabs.forEach((tab, index) => {
        if (tab.name === name) {
          const nextTab = tabs[index + 1] || tabs[index - 1]
          // console.log(nextTab, 'nextTab')
          if (nextTab) {
            activeName = nextTab.name
          }
        }
      })
      editableTabs.value = tabs.filter((tab) => tab.name !== name)
      editableTabs.value.shift()
      editableTabsValue.value = activeName
      curData.value['name'] = activeName
    } else {
      editableTabs.value = tabs.filter((tab) => tab.name !== name)
      editableTabs.value.shift()
    }
  }

  const handleClick = (pane: TabsPaneContext) => {
    if (pane.index != '0') {
      const num = Number(pane.index) - 1
      curData.value = editableTabs.value[num]
      pageName.value = curData.value.type
      model.value = curData.value.model
    }
  }

  // 右键事件
  const handleContextmenu = (data: any, event: any) => {
    const { clientX, clientY } = event
    deleteQuery.left = `${clientX + 20}px`
    deleteQuery.top = `${clientY}px`
    deleteQuery.visible = true
    deleteData.value = data
  }

  // 删除
  const handleDelete = (deleteAll: boolean) => {
    $baseConfirm('你确定要删除当前项吗', null, async () => {
      try {
        const ids: number[] = []
        listData.value.forEach((item) => {
          if (item.checked) ids.push(item.id!)
        })
        const { msg } = await deleteApiMonitorApi({ ids, deleteAll })
        $baseMessage(msg, 'success', 'vab-hey-message-success')
        editableTabs.value = editableTabs.value.filter((item: any) => {
          return !ids.includes(item.data.id)
        })
        listData.value.forEach((item: any, index: number) => {
          listData.value[index].checked = false
        })
        await getData()
        verifyError()
      } finally {
        deleteQuery.visible = false
      }
    })
  }
  const handleChange = () => {
    console.log(listData.value)
  }
</script>

<script lang="ts">
  export default {
    name: 'InterfaceMonitoring',
  }
</script>
<template>
  <div class="interface-monitoring-container">
    <div class="monitoring-del">
      <el-dropdown placement="bottom-start" style="margin-left: 10px" trigger="click">
        <span class="el-dropdown-link">
          <el-button type="danger">
            批量删除
            <el-icon class="el-icon--right"><arrow-down /></el-icon>
          </el-button>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="handleDelete(false)">删除选中</el-dropdown-item>
            <el-dropdown-item @click="handleDelete(true)">删除所有</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
    <el-tabs
      v-model="editableTabsValue"
      class="demo-tabs"
      closable
      type="card"
      @tab-click="handleClick"
      @tab-remove="removeTab"
    >
      <!-- @keydown.stop="moveFalg = true" @keyup.stop="moveFalg = false" -->
      <!-- @mouseenter="
                      () => {
                        if (!moveFalg) {
                          item.checked = true
                        }
                      }
                    " -->
      <el-tab-pane key="0" closable label="任务仪表盘" name="0">
        <div style="height: 100%; width: 100%">
          <!-- 展示页 -->
          <div v-loading="listLoading" class="content">
            <div class="my-items">
              <el-row class="my_row" :gutter="20">
                <el-col v-for="item in listData" :key="item.uuid" :span="4">
                  <div ref="divRef" class="item" :class="{ hasError: item.hasError }" @click="showDetail(item)">
                    <div class="canChecked" @click.stop>
                      <el-checkbox v-model="item.checked" @change="handleChange" />
                    </div>
                    <!-- @contextmenu.prevent="($event) => handleContextmenu(item, $event)" -->
                    <div class="warp">
                      <el-image
                        :src="require('@/assets/theme_images/meter-icon.png')"
                        style="height: 48px; width: 48px"
                      />
                      <div class="text">{{ item.name }}</div>
                    </div>
                    <div v-if="item.hasError" class="mark"></div>
                    <div v-if="item.hasError" class="mark-boder"></div>
                  </div>
                </el-col>
                <el-col :span="4">
                  <div class="item" style="font-size: 16px; background: initial" @click="addItemEvent">
                    <div class="warp">
                      <el-icon style="font-size: 20px"><Plus /></el-icon>
                    </div>
                  </div>
                </el-col>
              </el-row>
            </div>
            <!-- pagination -->
            <el-pagination
              v-model:current-page="queryPage.pageNum"
              background
              :layout="layout"
              :page-size="queryPage.pageSize"
              :page-sizes="[10, 20, 30, 50]"
              style="margin-bottom: 20px"
              :total="queryPage.total"
              @current-change="handleCurrentChange"
              @size-change="handleSizeChange"
            />
          </div>
          <!-- <div
            v-show="deleteQuery.visible"
            class="contextmenu"
            :style="{ left: deleteQuery.left, top: deleteQuery.top }"
          >
            <el-button type="primary" @click="handleDelete">删除</el-button>
          </div> -->
        </div>
      </el-tab-pane>
      <el-tab-pane v-for="item in editableTabs" :key="item.name" :label="item.title" :name="item.name">
        <!-- <template v-if="condition"> -->
        <div v-if="item.type === 'design-monitoring' && item.model == 'design'">
          <DesignMonitoring :all-data="currentItem" :is-btn="isBtn" @on-back="backEvent" @toDetail="toDetailPage" />
        </div>
        <div v-else>
          <DesignMonitoring
            v-if="pageName === 'design-monitoring'"
            :all-data="item.data"
            :is-btn="isBtn"
            @changeName="changeTabName"
            @on-back="backEvent"
            @toDetail="toDetailPage"
          />
          <DetailMonitoring v-else :current-item="item.data" @on-back="detailBackEvent" @on-todesign="toDesignPage" />
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped lang="scss">
  .interface-monitoring-container {
    position: relative;
    height: calc(100vh - 20px);
    overflow: hidden;
    .monitoring-del {
      position: absolute;
      display: flex;
      align-items: center;
      justify-content: end;
      top: 10px;
      right: 20px;
      height: 40px;
      width: 120px;
      border-bottom: 1px solid var(--el-border-color-light);
    }
  }
  :deep() {
    .el-tabs__header {
      width: calc(100% - 120px);
    }
    .el-popper {
      max-height: 150px;
      max-width: 60vw;
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
    }
    .my-items {
      height: calc(100vh - 170px);
      overflow: hidden;
    }
    .my_row {
      overflow-y: auto;
      max-height: 100%;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }
    }
    .el-tabs__nav-wrap.is-top,
    .el-tabs__nav-scroll {
      overflow: inherit;
    }
    .el-tabs__header.is-top {
      overflow: hidden !important;
    }
    .el-tabs__nav.is-top {
      border: 0 solid #fff;

      height: 40px;

      .is-active {
        background-color: var(--el-color-primary) !important;
        color: #fff !important;
      }

      .el-tabs__item {
        border-radius: 6px 6px 0px 0px;
        margin-right: 10px;
        background: var(--el-color-primary-light-9);
        border: 1px solid var(--el-color-primary-light-9);
        color: #303133;
        position: relative;

        &:hover {
          background: var(--el-color-primary-light-9);
          border: 1px solid var(--el-color-primary);
          color: var(--el-color-primary);
          padding: 0 20px;
        }
        &:first-child {
          padding: 0 20px;
          &:hover {
            padding: 0 20px;
          }
          .el-icon.is-icon-close {
            display: none;
          }
        }
        &:nth-child(n + 2) {
          .el-icon.is-icon-close {
            position: absolute;
            top: 0px;
            right: 0px;
            z-index: 4;
            color: #fff;
            background-color: var(--el-color-primary);
            border-radius: 0 0 0 6px;
          }
        }
      }
    }

    .el-tabs__nav-wrap.is-top {
      height: 40px;
    }
    .el-tabs__new-tab {
      display: none;
    }
    .el-tabs__item {
      max-width: 400px;
      overflow: hidden;
      line-height: 40px !important;
      text-overflow: ellipsis;
      white-space: nowrap;
      word-break: break-all;
      word-wrap: break-word;
    }
  }

  .content {
    width: 100%;
    height: calc(100vh - 80px);
    overflow-y: auto;
    overflow-x: hidden;
    &::-webkit-scrollbar {
      width: 0;
      height: 0;
    }
    padding-bottom: 20px;
  }

  .item {
    margin-top: 20px;
    background: url('~@/assets/theme_images/meter.png');
    height: 0;
    padding-bottom: 71%;
    border: 1px solid #f1f9ff;
    background-size: 100%;
    border-radius: 16px;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    position: relative;
    .canChecked {
      position: absolute;
      right: 10px;
      top: 0px;
      width: 40px;
      height: 40px;
      display: flex;
      align-items: baseline;
      justify-content: flex-end;
    }
    &.hasError {
      background: url('~@/assets/theme_images/meter-err.png');
    }
    .mark,
    .mark-boder {
      border-radius: 16px;
      height: 100%;
      width: 100%;
      position: absolute;
      left: 50%;
      top: 50%;
      transform: translate(-50%, -50%);
      // background-color: #ff3636;
    }
    .mark-boder {
      border: 1px solid #ff3636;
      // background-color: inherit;
      // opacity: 0.5;
    }
    .warp {
      height: 80px;
      width: 100%;
      position: absolute;
      left: 50%;
      top: 50%;
      transform: translate(-50%, -50%);
      display: flex;
      justify-content: center;
      align-items: center;
      flex-direction: column;
    }

    .text {
      margin-top: 10px;
      width: 90%;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      vertical-align: middle;
      text-align: center;
    }

    &:hover {
      cursor: pointer;
      box-shadow: 0px 0px 8px 1px rgba(0, 0, 0, 0.08);
    }
  }

  .contextmenu {
    position: fixed;
    top: 0;
    left: 0;
    margin: 0;
    background: #fff;
    z-index: 3000;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 400;
    color: #333;
    box-shadow: 2px 2px 3px 0 rgba(0, 0, 0, 0.3);
  }
</style>
