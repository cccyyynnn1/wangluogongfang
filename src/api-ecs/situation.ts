import request from '@/utils/request'
import { ResponseData } from '@/types'
import { AnyKindOfDictionary } from 'lodash'

/**
 * 大屏相关:站点列表
 */
export const stationListApi = (data: { pageNum: number; pageSize: number }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/sceneHighwayTollStation/getStationPage',
    method: 'post',
    data,
  })
}

/**
 * 大屏相关:站点新增
 */
export const stationAddApi = (data: {
  assetIp: string
  defenseType: number
  eastLongitude: number
  firewallIps: string
  northernLatitude: number
  stationName: string
  xdrasUuid: string
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/sceneHighwayTollStation/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * 大屏相关:站点更新
 */
export const stationUpdate = (data: {
  id: number
  assetIp: string
  defenseType: number
  eastLongitude: number
  firewallIps: string
  northernLatitude: number
  stationName: string
  xdrasUuid: string
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/sceneHighwayTollStation/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * 大屏相关:站点删除
 */
export const stationDelApi = (data: { ids: number[] }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/sceneHighwayTollStation/delete',
    method: 'DELETE',
    data,
  })
}

/**
 * 大屏相关:导出
 */
export const stationExportApi = (): Promise<any> => {
  return request({
    url: '/v3/ecsPlatform/sceneHighwayTollStation/exportStation',
    method: 'get',
    responseType: 'arraybuffer',
  })
}

/**
 * 大屏相关:导入
 */
export const stationImportApi = (data: { file?: File }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/sceneHighwayTollStation/importStation',
    method: 'POST',
    data,
  })
}

/**
 * 大屏相关:查询探针列表API
 */
export const flowDeviceApi = (params: {
  page: number
  limit: number
  query: AnyKindOfDictionary
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/public/forward',
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/flowDevice/page?page=${params.page}&limit=${params.limit}&query=${params.query}`,
    },
  })
}

/**
 * 大屏相关:查询防火墙列表API
 */
export const firewallApi = (params: { page: number; limit: number; query: any }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/public/forward',
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/firewall/page?page=${params.page}&limit=${params.limit}&query=${params.query}`,
    },
  })
}
