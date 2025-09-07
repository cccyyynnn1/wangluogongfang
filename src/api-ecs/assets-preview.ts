import request from '@/utils/request'
import {
  ResponseData,
  PreviewAssetsBodyType,
  PreviewAssetsItem,
  PreviewAssetsSummary,
  FingerprintItem,
  UpdateAssetInfo,
  pageChartsType,
  detailChartsType,
} from '@/types/index'

/**
 * @description 获取资产预览分页
 */
export const getAssetsPreviewListApi = (
  data: PreviewAssetsBodyType
): Promise<ResponseData<{ records: PreviewAssetsItem[]; total: number }>> => {
  return request({
    url: '/v3/ecsPlatform/assetNew/getAssetNewPage',
    method: 'post',
    data,
  })
}
/**
 * @description 获取资产-数据中心，网络分区统计
 */
export const getAssetsPreviewSummaryApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/assetNew/getSummary',
    method: 'post',
  })
}

/**
 * @description 资产标签模糊查询
 */
export const getAssetsLabelBySearchApi = (data: {
  searchStr?: string
  orderType: 0 | 1
  groupId: number
}): Promise<ResponseData<FingerprintItem[]>> => {
  return request({
    url: '/v3/ecsPlatform/assetLabel/getLabelBySearch',
    method: 'post',
    data,
  })
}

/**
 * 获取资产访问关系图
 */
export const getAssetNewRelatApi = (data: {
  ip: string
  startTime: string
  endTime: string
}): Promise<ResponseData<{ links: any[]; nodes: any[] }>> => {
  return request({
    url: '/v3/ecsPlatform/assetNew/getAssetNewRelat',
    method: 'post',
    data: data,
  })
}

/**
 * 获取防火墙策略数据
 */
export const getAssetFirewallPolicyApi = (params: {
  page: number
  limit: number
  query?: string
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/public/forward',
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/firewallPolicy/page?page=${params.page}&limit=${params.limit}&query=${params.query}`,
    },
  })
}

/**
 * 资产信息更新
 */
export const updateAssetInfoApi = (data: UpdateAssetInfo): Promise<any> => {
  return request({
    url: '/v3/ecsPlatform/assetNew/updateAssetNew',
    method: 'post',
    data,
  })
}

/**
 * 列表echarts数据
 */
export const pageChartsApi = (data: pageChartsType): Promise<any> => {
  return request({
    url: '/v3/ecsPlatform/assetNew/pageCharts',
    method: 'post',
    data,
  })
}

/**
 * 资产信息更新
 */
export const detailChartsApi = (data: detailChartsType): Promise<any> => {
  return request({
    url: '/v3/ecsPlatform/assetNew/detailCharts',
    method: 'post',
    data,
  })
}
