import request from '@/utils/request'

import { MonitoringType, ResponseData } from '@/types/index'

/**
 * 数据趋势
 */
export const getDataTrendApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsChart/trend/getDataTrend',
    method: 'post',
  })
}

/**
 * 流量趋势
 */
export const getFlowTrendInfoApi = (userId: number): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsChart/trend/getFlowTrendInfo/${userId}`,
    method: 'GET',
  })
}

/**
 * 攻击源IP
 */
export const getAttackSourceApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsChart/trend/getAttackSource',
    method: 'post',
  })
}

/**
 * 攻击类型
 */
export const getAttackSourceTypeApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsChart/trend/getAttackSourceType',
    method: 'post',
  })
}

/**
 * 工作台数据量统计api
 */
export const getDataVolumeApi = (params: { startTime: string; endTime: string }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/event/indexTotal',
    method: 'get',
    params,
  })
}

/**
 * 各站点数据量总和
 */
export const getSiteDataVolume = (): Promise<ResponseData<{ site: any }>> => {
  return request({
    url: '/v3/ecsChart/trend/getSiteDataVolume',
    method: 'post',
  })
}

/**
 * 各站点数据量总和
 */
export const getAccessStatisticsApi = (): Promise<ResponseData<{ access: any }>> => {
  return request({
    url: '/v3/ecsChart/trend/getAccessStatistics',
    method: 'get',
  })
}

/**
 * 受攻击趋势
 */
export const getSiteAttackSourceApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsChart/trend/getSiteAttackSource',
    method: 'post',
  })
}

/**
 * 环比分析
 */
export const getStatMinDataApi = (data: { startTime: string; endTime: string }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsChart/trend/getStatMinData',
    method: 'post',
    data,
  })
}

/**
 * api监控分页
 */
export const getApiMonitorPageApi = (data: { pageNum: number; pageSize: number }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/apiMonitor/getPage',
    method: 'post',
    data,
  })
}

/**
 * api监控保存/更新
 */
export const saveOrUpdateApiMonitorApi = (data: MonitoringType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/apiMonitor/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * 删除站点api
 */
export const deleteApiMonitorApi = (data: { ids: number[]; deleteAll: boolean }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/apiMonitor/delete',
    method: 'delete',
    data,
  })
}

/**
 * 工作台统计api
 */
export const homeCountApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/event/homeCount',
    method: 'post',
  })
}

/**
 * 工作台统计api
 */
export const homeLineApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/event/homeLine',
    method: 'post',
  })
}
