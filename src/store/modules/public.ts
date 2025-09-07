/**
 * @description 异常捕获的状态拦截，请勿修改
 */
import { getFlowProbesAPI } from '~/src/api-ecs/retrieve'
import { GlobalMessageItem, PublicType } from '/#/store'
import {
  getAllDisPlaysFiledApi,
  GetNewSecondaryMsgsApi,
  GlobalMsgGetGlobalMsgPageApi,
  GlobalMsgUpdateApi,
} from '~/src/api-ecs/public'
import { gp } from '~/library/plugins/vab'
import _, { update } from 'lodash'

export const usePubilcStore = defineStore('public', {
  state: (): PublicType => ({
    flowProbesList: [],
    flowProbesIDList: [],
    allDisPlaysFiledList: [],
    globalMessageList: [],
    globalLoginMessageList: {
      newMsgs: [],
      unReadCount: 0,
    },
    hasNewMessage: false,
    // globalSysMessageList: [],
    // globalSysIndex: 1,
    // globalUserIndex: 1,
    // globalUserMessageList: [],
  }),
  getters: {
    GetFlowProbesList: (state) => state.flowProbesList,
    GetFlowProbesIDList: (state) => state.flowProbesIDList,
    GetAllDisPlaysFiledList: (state) => state.allDisPlaysFiledList,
    GetGlobalMessageLists: (state) => state.globalMessageList,
    GetGlobalLoginMessageList: (state) => state.globalLoginMessageList,
    // GetGlobalSysList: (state) => state.globalSysMessageList,
    // GetGlobalUserList: (state) => state.globalUserMessageList,
  },
  actions: {
    async SetFlowProbes() {
      const { data } = await getFlowProbesAPI()
      this.flowProbesList = []
      this.flowProbesList.push(data)
      this.flowProbesIDList = []
      this.setflowProbesIDList(this.flowProbesList)
    },
    setflowProbesIDList(arr: any[]) {
      arr?.forEach((item: any) => {
        const child = item.data?.children || item?.children
        if (item.realId) {
          this.flowProbesIDList.push({ id: item.id, realId: item.realId })
        }
        if (child && child.length > 0) {
          this.setflowProbesIDList(child)
        }
      })
    },
    getFlowProbesIDList() {
      return this.flowProbesIDList
    },
    getFlowProbesList() {
      return this.flowProbesList
    },
    async SetAllDisPlaysFiledList() {
      const { data } = await getAllDisPlaysFiledApi()
      this.allDisPlaysFiledList = data
    },
    async setGlobalMessageList() {
      const { data } = await GlobalMsgGetGlobalMsgPageApi({
        pageNum: 1,
        pageSize: 20,
        minorType: 'collectShare,notice',
      })
      this.globalMessageList = data.records
    },

    setHasNewMessage() {
      this.hasNewMessage = this.globalMessageList.some((item) => {
        return !item.readStatus
      })
    },
    async setGlobalLoginMessageList() {
      const { data } = await GetNewSecondaryMsgsApi()
      this.globalLoginMessageList = data
    },
    updateHasNewMessage(value: boolean) {
      this.hasNewMessage = value
    },
    AddGlobalMessages(item: GlobalMessageItem) {
      this.globalMessageList.unshift(item)
      // this.setHasNewMessage()
    },
    // AddGlobalSysMessages(item: GlobalMessageItem) {
    //   this.globalSysMessageList.unshift(item)
    // },
    // AddGlobalUserMessages(item: GlobalMessageItem) {
    //   this.globalUserMessageList.unshift(item)
    // },
    GetGlobalUserMessageList() {
      return this.globalMessageList
    },
    GetHasNewMessageList() {
      return this.hasNewMessage
    },
    cleanAll() {
      this.flowProbesList = []
      this.allDisPlaysFiledList = []
      this.globalMessageList = []
      // this.globalSysMessageList = []
      // this.globalUserMessageList = []
      // this.globalSysIndex = 0
      // this.globalUserIndex = 0
      this.hasNewMessage = false
      this.globalLoginMessageList = {
        newMsgs: [],
        unReadCount: 0,
      }
    },
  },
})
