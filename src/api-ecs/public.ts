import request from '@/utils/request'
import { MsgTypeArrType, ResponseData } from '@/types/index'
import { GlobalLoginMessageType, GlobalMessageItem, TableColumnType } from '/#/store'
/**
 * 查询是否进行过检测过程
 */
export const getDeploymentCheckStatusApi = (): Promise<
  ResponseData<{
    id: number
    key: string
    name: string
    uid: number
    status: string
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/deploymentCheck/isCheck',
    method: 'post',
  })
}
/**
 * 查询是否进行过检测过程
 */
export const deploymentCheckApi = (status: number, data: any): Promise<ResponseData<any>> => {
  const URL = [
    null,
    '/v3/ecsPlatform/deploymentCheck/checkDb',
    '/v3/ecsPlatform/deploymentCheck/initializeData',
    '/v3/ecsPlatform/deploymentCheck/checkEs',
    '/v3/ecsPlatform/deploymentCheck/checkHugegraph',
    '/v3/ecsPlatform/deploymentCheck/checkKafka',
    '/v3/ecsPlatform/deploymentCheck/checkAlarmAddress',
    '/v3/ecsPlatform/deploymentCheck/restart',
  ]

  return request({
    url: URL[status]!,
    method: status === URL.length - 1 ? 'get' : 'post',
    data,
  })
}

/**
 * 是否激活
 */
export const getCheckActivateApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/eht/public/license',
    method: 'get',
  })
}

/**
 * 获取所有字段
 */
export const getAllFieldApi = (): Promise<ResponseData<TableColumnType>> => {
  return request({
    url: '/v3/ecsPlatform/indexField/getAll',
    method: 'get',
  })
}
/**
 * 获取所有字段
 */
export const getAuthorization = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/public/getActivateKey',
    method: 'get',
  })
}

/**
 * 获取用户绑定字段
 */
export const getAllDisPlaysFiledApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/indexFiledsUser/getAllDisPlaysFiled',
    method: 'get',
  })
}

/**
 * 获取所有索引类型
 */
export const getAllIndexTypeApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/indexField/getAllType',
    method: 'get',
  })
}
/**
 * 获取所有索引包含网络层会话类型
 */
export const getAllIndexTypeIncouldsNetApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/indexField/getAllTypeAndNetApp',
    method: 'get',
  })
}
/**
 * 获取机器信息（密钥激活）
 */
export const getActivateKeyApi = (): Promise<
  ResponseData<{
    customerName: string
    featureCode: string
    machineCode: string
  }>
> => {
  return request({
    url: '/v3/eht/public/getActivateKey',
    method: 'get',
  })
}
/**
 * 根据激活码激活
 */
export const toActivateApi = (data: {
  customerName: string
  machineCode: string
  featureCode: string
  productKey: string
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/eht/public/keyActivation',
    method: 'post',
    params: data,
  })
}

/**
 * 获取AI会话记录（全局AI）
 */
export const aiSessionChatsApi = (data: {
  pageNum: number
  pageSize: number
}): Promise<
  ResponseData<{
    records: {
      id: number
      title: string
    }[]
    total: number
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/aiChatHistory/getAiChatHistoryPage`,
    method: 'post',
    data,
  })
}

/**
 * 获取会话历史记录（全局AI）
 */
export const aiSessionContentApi = (
  historyId: number
): Promise<
  ResponseData<
    {
      content: string
      historyId: number
      id: number
      role: 'user' | 'assistant'
    }[]
  >
> => {
  return request({
    url: `/v3/ecsPlatform/aiChatContent/getContentByHistoryId?historyId=${historyId}`,
    method: 'get',
  })
}
/**
 * 删除会话历史记录（全局AI）
 */
export const deleteAiSessionContentApi = (data: {
  deleteAll: boolean
  ids: number[]
}): Promise<
  ResponseData<
    {
      content: string
      historyId: number
      id: number
      role: 'user' | 'assistant'
    }[]
  >
> => {
  return request({
    url: `/v3/ecsPlatform/aiChatHistory/deleteChatHistory`,
    method: 'delete',
    data,
  })
}

/**
 * 删除会话历史记录（全局AI）
 */
export const getAiQuickInputApi = (data: {
  attackIp: string
  victimIp: string
  host: string
  threatName: string
}): Promise<ResponseData<string[]>> => {
  return request({
    url: `/v3/ecsPlatform/gpt/quickInput`,
    method: 'post',
    data,
  })
}
/**
 * 获取站点会话的默认展示字段
 */
export const getHostApiFieldsApi = (
  indexType: number | string
): Promise<
  ResponseData<{
    ids: number[]
    names: string[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/indexTemplate/findDefaultFiledBySiteSessionAndApi?indexType=${indexType}`,
    method: 'get',
  })
}

/**
 * @description '全局消息：获取所有消息'
 * @params
 */
export const GlobalMsgGetGlobalMsgPageApi = (data: {
  pageNum: number
  pageSize: number
  minorType?: string
  readStatus?: 0 | 1 // 1:已读 0:未读
  sendUserId?: number
  searchStr?: number
}): Promise<{ data: { records: GlobalMessageItem[]; total: number } }> => {
  return request({
    url: `/v3/ecsPlatform/globalMsg/getGlobalMsgPage`,
    method: 'post',
    data,
  })
}

/**
 * @description '全局消息：更新已读'
 * @params
 */
export const GlobalMsgUpdateApi = (data: {
  minorType?: string
  sendUserId?: number
  containsAll?: boolean
  id?: number
}): Promise<{ msg: string; data: { curTypeUnReadCount: number; unReadCount: number } }> => {
  return request({
    url: `/v3/ecsPlatform/globalMsg/updateGlobalMsg`,
    method: 'post',
    data,
  })
}

/**
 * @description '全局消息：删除'
 * @params
 */
export const GlobalMsgDelApi = (data: {
  minorType?: string
  sendUserId?: number
  containsAll?: boolean
  id?: number
}): Promise<{ msg: string; data: { curTypeUnReadCount: number; unReadCount: number } }> => {
  return request({
    url: `/v3/ecsPlatform/globalMsg/deleteGlobalMsg`,
    method: 'DELETE',
    data,
  })
}

/**
 * @description '全局消息：获取最新次要消息'
 * @params
 */
export const GetNewSecondaryMsgsApi = (): Promise<{ data: GlobalLoginMessageType }> => {
  return request({
    url: `/v3/ecsPlatform/globalMsg/getImportantMsgs`,
    method: 'get',
  })
}

/**
 * @description '全局消息-获取消息类型列表'
 * @params
 */
export const GetMsgTypeArrApi = (): Promise<{ data: { sys: MsgTypeArrType[]; user: MsgTypeArrType[] } }> => {
  return request({
    url: `/v3/ecsPlatform/globalMsg/getMsgTypeArr`,
    method: 'get',
  })
}

/**
 * @description '全局消息-根据id获取消息详情'
 * @params
 */
export const GetGlobalByIdApi = (params: { id: number }): Promise<{ data: GlobalMessageItem }> => {
  return request({
    url: `/v3/ecsPlatform/globalMsg/getGlobalById`,
    method: 'get',
    params,
  })
}

/**
 * @description '全局消息开关状态-更新开关'
 * @params
 */

export const GlobalMsgSwitchIdApi = (data: { minorType: string; enable: boolean }): Promise<{ msg: string }> => {
  return request({
    url: `/v3/ecsPlatform/globalMsgSwitch/updateSwitch`,
    method: 'POST',
    data,
  })
}
//获取检索历史
export const getSearchHistories = (data: {
  pageNum: number
  pagesize: number
  indexType?: number
  workspaceId?: number
  siteId?: number
  apiId?: number
}): Promise<{ data: { records: { id: number; searchSql: string }[] } }> => {
  return request({
    url: `/v3/ecsPlatform/collectHistory/getSearchHistoryPage`,
    method: 'post',
    data,
  })
}
//获取高亮特征
export const getAllHightLightAPI = (): Promise<
  ResponseData<{
    highLightWhite: {
      content: string
      scope: 'all' | 'requestHeader' | 'requestPayload' | 'responseHeader' | 'responsePayload'
    }[]
    highLightConfig: {
      content: string
      scope: 'all' | 'requestHeader' | 'requestPayload' | 'responseHeader' | 'responsePayload'
    }[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/highLightConfig/getAllHightLightData`,
    method: 'get',
  })
}

/**
 * @description '登录：29443'
 * @params
 */
export const getLogin29443Api = (): Promise<any> => {
  return request({
    url: '/v3/ecsPlatform/public/forwardLogin',
    method: 'POST',
    data: {
      reqType: 'POST',
      url: '/pub/login',
      paramMap: {
        captcha: '',
        password:
          'TAvSAnnhW1RKlCjkujETnPQhp2ephwRetG5ytBiPfGUXYGsXwpTeUdqg3+Wv6ODMytwW/GPA0qjqUISi8321sGMw7FiH+J6SFFEDl0ktV0h0GII0ImpnIiS0H6P11Asjgu2MKr1R/2uR/SZ3eRPaUNQEhGC2BUpK5uBF2bSZO5U=',
        userName: 'proadmin',
      },
    },
  })
}
