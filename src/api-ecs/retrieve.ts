import request from '@/utils/request'
import {
  ResponseData,
  SearchBySqlParams,
  PublishRuleParams,
  RuleModel,
  WorkerSpaceItem,
  RetrieveHistoryParams,
  RetrieveHistoryItem,
  ShareHistoryParams,
  ServiceLink,
  ServiceNode,
  NewNetworkLayerParams,
  FindHistoryParams,
  LevelRuleSaveUpdateModel,
} from '@/types/index'

/**
 * 获取资产列表
 */
export const searchBySqlApi = (data: SearchBySqlParams, signal?: AbortSignal): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/event/searchBySql',
    method: 'post',
    data,
    signal,
  })
}
/**
 * 发布流量场景模版
 */
export const publishRuleApi = (data: PublishRuleParams): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/eventRule/saveOrUpdate',
    method: 'post',
    data,
  })
}
/**
 * 获取流量场景模版
 */
export const getRuleApi = (data: { pageNum: number; pageSize: number }): Promise<ResponseData<RuleModel>> => {
  return request({
    url: '/v3/ecsPlatform/eventRule/getAll',
    method: 'post',
    data,
  })
}

/**
 * 根据id获取流量模版配置
 */
export const getRuleConfigsApi = (id: number): Promise<ResponseData<PublishRuleParams>> => {
  return request({
    url: '/v3/ecsPlatform/eventRule/getById',
    method: 'post',
    params: { id },
  })
}
/**
 * 获取平台地址
 */
export const getKeyVarApi = (key = 'yunche_ras'): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/settingVar/getByKeys?keys=${key}`,
    method: 'get',
  })
}

/**
 * 获取平台地址全部信息
 */
export const getKeyVarInfoApi = (key = 'yunche_ras,ecs_storage,ecs_kafka,ecs_neo4j'): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/settingVar/getByKeys?keys=${key}`,
    method: 'get',
  })
}

// 元数据
export function searchMetadata(data: { id: number; indexType: number; requestTimeNs: number }) {
  return request({
    url: '/v3/ecsPlatform/event/getEventByIddate',
    method: 'post',
    data,
  })
}
// 删除流量模型
export function deleteTrafficApi(params: { id: string }): Promise<ResponseData<any>> {
  return request({
    url: '/v3/ecsPlatform/eventRule/delete',
    method: 'delete',
    params,
  })
}
// 获取Top字段统计
export function getTopFieldApi(data: {
  indexType: number
  topCount: number
  startTime: string
  endTime: string
  aggregationFields: string
  searchSql: string
  whiteType?: number
}): Promise<
  ResponseData<{
    aggObj: { name: string; value: number }[]
  }>
> {
  return request({
    url: '/v3/ecsPlatform/event/aggregations',
    method: 'post',
    data,
  })
}

/**
 *  DNS数字编码表
 */
// export const getDnsDictApi = (): Promise<ResponseData<any>> => {
//   return request({
//     url: '/v3/ecsPlatform/event/getDnsDict',
//     method: 'get',
//   })
// }

/**
 *  应用层数字编码表
 */
// export const appProtocolDictApi = (): Promise<ResponseData<any>> => {
//   return request({
//     url: '/v3/ecsPlatform/event/getAppProtocolDict',
//     method: 'get',
//   })
// }

/**
 *  获取工作空间
 */
export const getAllSpaceApi = (): Promise<ResponseData<WorkerSpaceItem[]>> => {
  return request({
    url: '/v3/ecsPlatform/workSpace/getAllSpace',
    method: 'post',
  })
}

/**
 *  更新工作空间
 */
export const saveOrUpdateSpaceApi = (spaceName: string): Promise<ResponseData<WorkerSpaceItem[]>> => {
  return request({
    url: '/v3/ecsPlatform/workSpace/saveOrUpdate',
    method: 'post',
    data: { spaceName },
  })
}
/**
 *  清空工作空间
 */
export const emptySpaceApi = (id: number): Promise<ResponseData<WorkerSpaceItem[]>> => {
  return request({
    url: `/v3/ecsPlatform/workSpace/emptySpace?id=${id}`,
    method: 'get',
  })
}
/**
 *  删除工作空间
 */
export const deleteSpaceApi = (id: number[] | number): Promise<ResponseData<WorkerSpaceItem[]>> => {
  return request({
    url: '/v3/ecsPlatform/workSpace/delete',
    method: 'delete',
    params: { ids: id },
  })
}
/**
 * 获取所有用户名
 */
export const getAllUserNameAPI = (): Promise<any> => {
  return request({
    url: '/v3/ecsPlatform/user/getAllUserName',
    method: 'get',
  })
}
/**
 *  获取检索历史
 */
export const getRetrieveHistoryApi = (
  data: RetrieveHistoryParams
): Promise<ResponseData<{ records: RetrieveHistoryItem[]; total: number }>> => {
  return request({
    url: '/v3/ecsPlatform/collectHistory/getCollectPage',
    method: 'post',
    data,
  })
}
/**
 *  获取检索历史
 */
export const findHistoryApi = (
  data: FindHistoryParams
): Promise<ResponseData<{ records: RetrieveHistoryItem[]; total: number }>> => {
  return request({
    url: '/v3/ecsPlatform/collectHistory/findHistory',
    method: 'post',
    data,
  })
}

/**
 *  获取分享历史
 */
export const getSkareHistoryApi = (
  data: ShareHistoryParams
): Promise<ResponseData<{ records: RetrieveHistoryItem[]; total: number }>> => {
  return request({
    url: '/v3/ecsPlatform/collectShare/getSharePage',
    method: 'post',
    data,
  })
}

/**
 *  取消收藏
 */
export const UpdateHistoryApi = (data: RetrieveHistoryItem): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/collectHistory/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 *  取消分享
 */
export const UpdateCollectShareApi = (id: number[] | number): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/collectShare/delete',
    method: 'delete',
    params: { ids: id },
  })
}
/**
 *  去分享
 */
export const toShareApi = (data: {
  rNames: string
  workspaceId?: number
  sName: string
  collectId?: number
  searchEdTime?: string
  searchStTime?: string
  searchSql?: string
  indexType?: number
  remarks?: string
  filterSqlArr?: string
  sqlRelat?: string
  inputSql?: string
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/collectShare/save',
    method: 'post',
    data,
  })
}
/**
 * 获取路径追踪数据
 */
export const getPathTracingApi = (data: {
  searchSql: string
  indexType: number
  startTime: string
  endTime: string
  aggregationFields: string[]
  count: number
  isPathTrack: boolean
  isTraceSource: boolean
  isAssetVisit: boolean
}): Promise<ResponseData<{ links: ServiceLink[]; nodes: ServiceNode[] }>> => {
  return request({
    url: '/v3/ecsPlatform/public/getServiceAccessForHttp',
    method: 'post',
    data,
  })
}

/**
 * 获取收藏夹数据类型
 */
export const getCollectLogDataTypeApi = (
  workspaceId: number
): Promise<ResponseData<{ name: string; value: number }[]>> => {
  return request({
    url: '/v3/ecsPlatform/collectLog/getAllDataType',
    method: 'get',
    params: { workspaceId },
  })
}
/**
 * 获取收藏夹类型日志
 */
export const getLogPageApi = (data: {
  pageNum: number
  pageSize: number
  indexType: number
  workspaceId: number
}): Promise<ResponseData<{ records: RetrieveHistoryItem[]; total: number }>> => {
  return request({
    url: '/v3/ecsPlatform/collectLog/getLogPage',
    method: 'post',
    data,
  })
}
/**
 * 获取收藏夹类型日志
 */
export const favoritesLogPageApi = (data: {
  dataLog: string
  indexType: number
  workspaceId: number
}): Promise<ResponseData<{ records: RetrieveHistoryItem[]; total: number }>> => {
  return request({
    url: '/v3/ecsPlatform/collectLog/saveOrUpdate',
    method: 'post',
    data,
  })
}
/**
 * 获取收藏夹类型日志
 */
export const deletefavoritesLogPageApi = (ids: number[]): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/collectLog/delete?ids=${ids.toString()}`,
    method: 'delete',
  })
}
/**
 * 添加日志到溯源图
 */
export const appendCollectTraceApi = (data: {
  dataLog: string
  indexType: number
  workspaceId: number
}): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/collectTrace/saveOrUpdate`,
    method: 'post',
    data,
  })
}
/**
 * 获取检索溯源图
 */
export const getTetrieveTraceApi = (data: {
  workspaceId: number
  isPathTrack: boolean
  isTraceSource: boolean
  isAssetVisit: boolean
}): Promise<ResponseData<{ links: ServiceLink[]; nodes: ServiceNode[] }>> => {
  return request({
    url: `/v3/ecsPlatform/public/getServiceAccessForHttp`,
    method: 'post',
    data,
  })
}
/**
 * 获取资产访问
 */
export const getAssetVisitApi = (data: {
  searchSql: string
  startTime: string
  endTime: string
  indexType: number
  isPathTrack: boolean
  isTraceSource: boolean
  isAssetVisit: boolean
  count: number
}): Promise<ResponseData<{ links: ServiceLink[]; nodes: ServiceNode[] }>> => {
  return request({
    url: `/v3/ecsPlatform/public/getServiceAccessForHttp`,
    method: 'post',
    data,
  })
}

/**
 * 获取检索溯源图深度挖掘类型
 */
export const getTetrieveTypeApi = (workspaceId: number): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/collectTrace/getAllDataType?workspaceId=${workspaceId}`,
    method: 'get',
  })
}
/**
 * 获取检索溯源图深度挖掘类型
 */
export const getTetrieveLogApi = (data: {
  pageNum: number
  pageSize: number
  indexType: number
  workspaceId: number
  nodeIp?: string
  clientIp?: string
  serverIp?: string
}): Promise<ResponseData<{ records: any[]; total: number }>> => {
  return request({
    url: `/v3/ecsPlatform/collectTrace/getLogPage`,
    method: 'post',
    data,
  })
}
/**
 * 获取检索溯源图深度挖掘类型
 */
export const getIpLabelApi = (data: string[]): Promise<ResponseData<string[]>> => {
  return request({
    url: `/v3/ecsPlatform/public/getXffLabel`,
    method: 'post',
    data,
  })
}
/**
 * 获取网络应用层右侧上方下拉树
 */
export const getFlowProbesAPI = (): Promise<{ data: any }> => {
  return request({
    url: '/v3/ecsPlatform/public/forward',
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/flowSearch/getFlowProbes`,
    },
  })
}
/**
 * 获取网络应用层右侧下方树
 */
export const getEventStatisticsAPI = (params: {
  top: string
  moduleType: string
  topField: string
  query: any
}): Promise<{ data: any }> => {
  return request({
    url: '/v3/ecsPlatform/public/forward',
    method: 'POST',
    data: {
      reqType: 'POST',
      url: `/flowSearch/eventStatistics?top=${params.top}&moduleType=${params.moduleType}&topField=${params.topField}&query=${params.query}`,
    },
  })
}

/**
 * 获取网络应用层表格数据
 */
export const getTabelListAPI = (params: {
  top: string
  moduleType: string
  query: any
  groupBy: string
}): Promise<any> => {
  return request({
    url: '/v3/ecsPlatform/public/forward',
    method: 'POST',
    data: {
      reqType: 'POST',
      url: `/flowSearch/query?top=${params.top}&moduleType=${params.moduleType}&groupBy=${params.groupBy}&query=${params.query}`,
    },
  })
}

/**
 *  删除分享
 */
export const deleteCollectApi = (id: number[] | number, name: string): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/collectShare/deleteByRname',
    method: 'delete',
    data: { ids: id, rName: name },
  })
}
/**
 * 获取网络应用层解码右侧list
 */
export const getPacketDecodeListAPI = (data: { top: string; packetCut: string; query: any }): Promise<any> => {
  return request({
    url: '/v3/ecsPlatform/public/forward',
    method: 'POST',
    data: {
      reqType: 'POST',
      url: `/packetDecode/list`,
      paramMap: data,
    },
  })
}

/**
 * 获取PCAP保存参数
 */
export const getDownloadQueryAPI = (data: { query: string }): Promise<any> => {
  return request({
    url: '/v3/ecsPlatform/public/netDownloadPcap',
    method: 'POST',
    data,
  })
}

/**
 * 更新规则
 */
export const LevelRuleSaveUpdateAPI = (data: LevelRuleSaveUpdateModel): Promise<any> => {
  return request({
    url: '/v3/ecsPlatform/levelRule/saveUpdate',
    method: 'post',
    data,
  })
}

/**
 * 规则分页
 */
export const GetLevelRulePageAPI = (data: { pageNum: number; pageSize: number }): Promise<any> => {
  return request({
    url: '/v3/ecsPlatform/levelRule/getLevelRulePage',
    method: 'post',
    data,
  })
}

/**
 * 规则删除
 */
export const DelLevelRulePageAPI = (data: { deleteAll: boolean; ids?: number[] }): Promise<any> => {
  return request({
    url: '/v3/ecsPlatform/levelRule/deleteLevelRule',
    method: 'DELETE',
    data,
  })
}

/**
 * 统计下载
 */
export const EventExportAggAPI = (data: { name: string; value: number }[]): Promise<any> => {
  return request({
    url: '/v3/ecsPlatform/event/exportAgg',
    method: 'post',
    data,
    responseType: 'arraybuffer',
  })
}

/**
 * 文件还原 rest/v1/flowSearch/showFile
 */
export const getFlowSearchShowFileAPI = (data: any): Promise<any> => {
  return request({
    url: '/v3/ecsPlatform/public/showFile',
    method: 'post',
    data,
  })
}

/**
 * 获取当前条件下的索引数量
 */
export const getIndexCountAPI = (data: SearchBySqlParams): Promise<ResponseData<{ [key: string]: string }[]>> => {
  return request({
    url: '/v3/ecsPlatform/event/getIndexCount',
    method: 'post',
    data,
  })
}
