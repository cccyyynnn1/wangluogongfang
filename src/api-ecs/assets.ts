import request from '@/utils/request'
import {
  ResponseData,
  ServiceLink,
  ServiceNode,
  getAssetsType,
  ExportAssetsType,
  addAssetsType,
  NetworkPartitionType,
  ApplicationType,
  pageChartsType,
} from '@/types/index'

/**
 * 获取资产列表
 */
export const getAssetsListApi = (data: getAssetsType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/asset/getPage',
    method: 'post',
    data,
  })
}

/**
 * 资产更新
 */
export const updateAssetsApi = (data: any): Promise<ResponseData<string>> => {
  return request({
    url: '/v3/ecsPlatform/asset/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * 删除资产
 */
export const deleteAssetsApi = (params: {
  deleteAll: boolean
  knowOrNot: string
  ids: string
}): Promise<ResponseData<string>> => {
  return request({
    url: '/v3/ecsPlatform/asset/delete',
    method: 'DELETE',
    params,
  })
}

/**
 * 获取资产网段
 */
export const getAssetsSetApi = (): Promise<ResponseData<string>> => {
  return request({
    url: '/v3/ecsPlatform/assetRule/getRules',
    method: 'GET',
  })
}

/**
 * 修改资产网段
 */
export const updateAssetsSetApi = (params: { range: string }): Promise<ResponseData<string>> => {
  return request({
    url: '/v3/ecsPlatform/assetRule/saveUpdate',
    method: 'post',
    params,
  })
}

/**
 * 获取访问关系
 */
export const getServiceAccessApi = (params: {
  indexType: number
  pageNum: number
  pageSize: number
  searchSql: string
  startTime: string
  endTime: string
}): Promise<ResponseData<{ resList: any[]; total: number }>> => {
  return request({
    url: '/v3/ecsPlatform/public/getServiceAccessPage',
    method: 'post',
    data: { ...params, aggFieIds: 'clientIp,serverIp,serverPort' },
  })
}

/**
 * 获取访问关系
 */
export const getServiceLinksApi = (params: {
  indexType: number
  count: number
  searchSql?: string
  startTime: string
  endTime: string
  tid?: string
}): Promise<
  ResponseData<{
    links: ServiceLink[]
    nodes: ServiceNode[]
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/public/getServiceAccessForEvent',
    method: 'post',
    data: params,
  })
}

/** 更新溯源关系*/

/* content: {
  links: ServiceLink[]
  nodes: ServiceNode[]
} */
export const updateServiceLinksApi = (
  data: {
    tid?: string | number
    content: string
  }[]
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/traceSource/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * 导出资产
 */
export const exportAssetsApi = (data: ExportAssetsType): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/asset/exportAsset',
    method: 'post',
    data,
    responseType: 'arraybuffer',
  })
}

/**
 * 模版下载
 */
export const templateAssetsApi = (): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/asset/exportTemplate',
    method: 'post',
    responseType: 'arraybuffer',
  })
}

/**
 * 导入资产
 */
export const importAssetsApi = (data: { file?: File; type: number }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/asset/importAsset',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data,
  })
}

/**
 * 获取网路分区列表
 */
export const getNetworkPartitionListApi = (data: NetworkPartitionType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/netPartition/getNetPage',
    method: 'post',
    data,
  })
}

/**
 * 导出网路分区列表
 */
export const exportNetworkPartitionApi = (data: {
  ids?: string[]
  name?: string
  rule?: string
}): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/netPartition/exportNet',
    method: 'post',
    data,
    responseType: 'arraybuffer',
  })
}

/**
 * 网路分区列表模版下载
 */
export const templateNetworkPartitionApi = (): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/netPartition/exportTemplate',
    method: 'post',
    responseType: 'arraybuffer',
  })
}

/**
 * 导入网路分区列表
 */
export const importNetworkPartitionApi = (data: { file?: File }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/netPartition/importNet',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data,
  })
}

/**
 * 修改网路分区
 */
export const updateNetworkPartitionSetApi = (data: {
  name: string
  rule: string
  id?: number
}): Promise<ResponseData<string>> => {
  return request({
    url: '/v3/ecsPlatform/netPartition/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * 获取数据中心
 */
export const getDataSourceApi = (): Promise<
  ResponseData<{
    datasource: any[]
    region: {
      cn: string
      en: string
    }[]
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/netPartition/getDataSource',
    method: 'post',
  })
}

/**
 * 删除网路分区
 */
export const deleteNetworkPartitionApi = (params: {
  ids: string
  deleteAll: boolean
}): Promise<ResponseData<string>> => {
  return request({
    url: '/v3/ecsPlatform/netPartition/delete',
    method: 'DELETE',
    params,
  })
}

/**
 * 获取业务链列表
 */
export const getChainSsortListApi = (data: {
  pageNum?: number
  pageSize?: number
  site?: string
  status?: string
  limit?: number
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/chainSort/getPage',
    method: 'post',
    data,
  })
}

/**
 * 添加业务链
 */
export const addChainSsortApi = (data: { site: string }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/chainSort/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * 业务链删除选中
 */
export const deleteChainSsortApi = (data: { ids: number[]; deleteAll: boolean }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/chainSort/delete',
    method: 'post',
    data,
  })
}

/**
 * 业务链批量开启
 */
export const switchOpenChainSsortApi = (data: { ids: number[] }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/chainSort/switchOpen',
    method: 'post',
    data: data,
  })
}

/**
 * 业务链批量关闭
 */
export const switchCloseChainSsortApi = (data: { ids: number[] }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/chainSort/switchClose',
    method: 'post',
    data: data,
  })
}

/**
 * 保存/更新站点应用
 */
export const siteAppSaveUpdateApi = (data: ApplicationType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteApp/saveUpdate',
    method: 'post',
    data,
  })
}

/**
 * 获取站点应用列表
 */
export const getAppPageListApi = (data: NetworkPartitionType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteApp/getAppPage',
    method: 'post',
    data,
  })
}

/**
 * 导出站点应用列表
 */
export const exportApplicationApi = (data: { siteId: number; ids?: number[] }): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/siteApp/exportApp',
    method: 'post',
    data,
    responseType: 'arraybuffer',
  })
}

/**
 * 站点应用列表模版下载
 */
export const templateApplicationApi = (): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/siteApp/exportTemplate',
    method: 'post',
    responseType: 'arraybuffer',
  })
}

/**
 * 导入站点应用列表
 */
export const importApplicationApi = (data: { file?: File; siteId?: number }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteApp/importApp',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data,
  })
}

/**
 * 删除站点应用
 */
export const deleteApplicationApi = (data: {
  siteId: number
  deleteAll: boolean
  ids?: number[]
}): Promise<ResponseData<string>> => {
  return request({
    url: '/v3/ecsPlatform/siteApp/delete',
    method: 'DELETE',
    data,
  })
}

/**
 * 导出白名单
 */
export const exportAssetNewApi = (data: pageChartsType): Promise<BlobPart> => {
  return request({
    url: `/v3/ecsPlatform/assetNew/exportAssetNew`,
    method: 'post',
    data,
    responseType: 'arraybuffer',
  })
}
