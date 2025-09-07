import request from '@/utils/request'
import { ResponseData, siteApiInterfaceType, siteApiInterfaceItemType } from '@/types/index'
import {
  siteGetPageType,
  siteSearchType,
  siteAggregationsType,
  siteSessionDeleteType,
  siteSessionSaveUpdateType,
  SiteHistoryItem,
  siteHistoryParams,
  SiteUnknowPageType,
  UpdateUnknowSiteResqust,
  updateSiteUnknowGroupByRequest,
} from '../types/site'

/**
 * 站点检索
 */
export const siteSearchApi = (data: siteSearchType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/search',
    method: 'post',
    data,
  })
}

/**
 * 站点所有列表
 */
export const siteGetAllPageApi = (params?: { name?: string }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/getAll',
    method: 'post',
    params,
  })
}

/**
 * 获取所有已存在站点id和apid
 */
export const getSiteIdAppIdApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/getSiteIdAppId',
    method: 'post',
  })
}

/**
 * 站点分页
 */
export const siteGetPageApi = (data: siteGetPageType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/getPage',
    method: 'post',
    data,
  })
}

/**
 * @description '未知站点分页'
 * @params
 */
export const getSiteUnknowPageApi = (data: SiteUnknowPageType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteUnknow/getSiteUnknowPage',
    method: 'post',
    data,
  })
}

/**
 * 聚合检索
 */
export const siteAggregationsApi = (data: siteAggregationsType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/aggregations',
    method: 'post',
    data,
  })
}

/**
 * 删除站点会话
 */
export const siteSessionDeleteApi = (params: {
  ids: string
  module: number
  deleteAll: boolean
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/delete',
    method: 'delete',
    params,
  })
}

/**
 * 新增修改会话
 */
export const siteSaveUpdateApi = (data: siteSessionSaveUpdateType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/saveUpdate',
    method: 'post',
    data,
  })
}

/**
 * 站点详情
 */
export const getBySiteIdApi = (params: siteSessionDeleteType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/getBySiteId',
    method: 'get',
    params,
  })
}

/**
 * @description '获取默认字段'
 * @params
 */
export const getAllDefaultFieldsApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/indexField/findAllIndexFieldsDefault',
    method: 'get',
  })
}

/**
 * 告警标签
 */
export const getAlretTagApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteLabel/getDataWarnLabels',
    method: 'get',
  })
}

/**
 * 所有标签
 */
export const getAllTagApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/payloadLabel/getAll',
    method: 'get',
  })
}

/**
 * 通过host获取ip
 */
export const getIpByHostApi = (data: { hosts: string }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/getIpByHost',
    method: 'post',
    data,
  })
}
/**
 * 获取API接口
 */
export const getApiInterfaceApi = (data: {
  pageNum: number
  pageSize: number
  apiName: string
  apiUrl: string
  sessionId: number
}): Promise<ResponseData<siteApiInterfaceType>> => {
  return request({
    url: '/v3/ecsPlatform/siteApi/getPage',
    method: 'post',
    data,
  })
}
/**
 * 更新API接口
 */
export const updateApiInterfaceApi = (data: {
  sessionId: number
  displayFields: string
  apiName: string
  apiUrl: string
  searchSql: string
  id?: number
  isDisplay: 0 | 1
}): Promise<ResponseData<siteApiInterfaceType>> => {
  return request({
    url: '/v3/ecsPlatform/siteApi/saveOrUpdate',
    method: 'post',
    data,
  })
}
/**
 * 删除API接口
 */
export const deleteApiInterfaceApi = (ids: string): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteApi/delete',
    method: 'delete',
    params: { ids },
  })
}
/**
 * 站点API接口
 */
export const getApiInterfaceApiByIdApi = (id: string): Promise<ResponseData<siteApiInterfaceItemType[]>> => {
  return request({
    url: '/v3/ecsPlatform/siteApi/getBySessionId',
    method: 'get',
    params: { id },
  })
}

/**
 * 导出站点资产
 */
export const exportSiteAssetsApi = (params: { ids: string }): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/exportSite',
    method: 'post',
    params,
    responseType: 'arraybuffer',
  })
}

/**
 * 导出未知站点资产
 */
export const exportUnknowAssetsApi = (data: {
  ids: number[]
  host?: string
  clientIp?: string
  serverIp?: string
  position?: string
}): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/siteUnknow/exportUnknowSite',
    method: 'post',
    data,
    responseType: 'arraybuffer',
  })
}

/**
 * 站点模版下载
 */
export const templateSiteAssetsApi = (): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/exportTemplate',
    method: 'post',
    responseType: 'arraybuffer',
  })
}

/**
 * 导入站点资产
 */
export const importSiteAssetsApi = (data: { file?: File }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/importSite',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data,
  })
}

/**
 * 站点根据host获取json
 */
export const getSiteApiListApi = (data: { hosts?: string[]; url?: string }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/getHostArr',
    method: 'post',
    data,
  })
}

/**
 * @description '站点根据host获取json'
 * @params
 */

export const getListBySiteApiApi = (data: {
  startTime: string
  endTime: string
  hosts?: string[]
  url?: string
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/getHostArrNew',
    method: 'post',
    data,
  })
}

/**
 * 导出站点资产
 */
export const getSiteAppInterActionApi = (id: number): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/getAppGraphs',
    method: 'get',
    params: { id },
  })
}

/**
 *  获取站点历史
 */
export const getSiteHistoryApi = (
  data: siteHistoryParams
): Promise<ResponseData<{ records: SiteHistoryItem[]; total: number }>> => {
  return request({
    url: '/v3/ecsPlatform/collectSiteHistory/getCollectSitePage',
    method: 'post',
    data,
  })
}

/**
 *  取消站点收藏
 */
export const UpdateHistoryApi = (data: SiteHistoryItem): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/collectSiteHistory/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * @description '导出'
 * @params
 */
export const exportResponseDataApi = (data: { hosts: string[] }): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/exportApiList',
    method: 'post',
    data,
    responseType: 'arraybuffer',
  })
}
/**
* @description '根据host搜索更多host'
* @params  {
    endTime: string;
    indexType: number;
    searchStr: string;
    startTime: string;
    topCount: number;
}
*/
export const getPortByHostApi = (data: {
  endTime: string
  indexType: number
  searchStr: string
  startTime: string
  topCount: number
  isUnknowSite?: boolean
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/getPortByHost',
    method: 'post',
    data,
  })
}

/**
 * @description '未知站点编辑'
 * @params
 */
export const updateUnknowSiteApi = (data: UpdateUnknowSiteResqust): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteUnknow/updateUnknowSite',
    method: 'post',
    data,
  })
}
/**
 * @description 获取API列表（带数量）
 * @params
 */
export const getApiListCountBySessionldApi = (id: string): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/siteApi/getApiListCountBySessionId?id=${id}`,
    method: 'get',
  })
}

/**
 * @description '未知站点字段统计'
 * @params updateSiteUnknowGroupByRequest
 */
export const updateSiteUnknowGroupByApi = (data: updateSiteUnknowGroupByRequest): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteUnknow/groupBy',
    method: 'post',
    data,
  })
}
