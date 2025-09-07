import request from '@/utils/request'
import { ResponseData, forwordRuleUpdateType, MaintenanceOptions, WarnCustomerRuleQueryType } from '@/types'

/**
 * 获取所有kafka配置信息
 */
export const forwordRuleUpdateApi = (data: forwordRuleUpdateType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/forwardRules/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * 获取信息维护所有下拉
 */
export const getMaintenanceOptionsApi = (): Promise<ResponseData<MaintenanceOptions>> => {
  return request({
    url: '/v3/ecsPlatform/infoRule/getAllSelect',
    method: 'get',
  })
}

export const getMaintenancelistApi = (data: {
  pageNum: number
  pageSize: number
  searchStr: string
}): Promise<
  ResponseData<{
    records: object[]
    total: number
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/infoRule/getInfoRulePage',
    method: 'post',
    data,
  })
}

export const updateMaintenanceInfoApi = (
  data: any
): Promise<
  ResponseData<{
    records: object[]
    total: number
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/infoRule/updateInfoRule',
    method: 'post',
    data,
  })
}

export const updateMaintenanceStatusApi = (data: {
  ids: number[]
  enable: '启用' | '停用'
  selectAll: boolean
}): Promise<
  ResponseData<{
    records: object[]
    total: number
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/infoRule/updateStatus',
    method: 'post',
    data,
  })
}

export const getWarnCustomerRulePageApi = (
  data: WarnCustomerRuleQueryType
): Promise<
  ResponseData<{
    records: any[]
    total: number
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/warnCustomerRule/getWarnCustomerRulePage',
    method: 'post',
    data,
  })
}

export const updateWarnCustomerRulePApi = (
  data: any
): Promise<
  ResponseData<{
    records: any[]
    total: number
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/warnCustomerRule/saveOrUpdate',
    method: 'post',
    data,
  })
}
export const updateWarnCustomerRuleStatusApi = (data: {
  enable: string
  ids: number[]
}): Promise<
  ResponseData<{
    records: any[]
    total: number
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/warnCustomerRule/updateStatus',
    method: 'post',
    data,
  })
}

export const deleteWarnCustomerRuleApi = (data: {
  deleteAll: boolean
  ids: number[]
}): Promise<ResponseData<unknown>> => {
  return request({
    url: '/v3/ecsPlatform/warnCustomerRule/delete',
    method: 'delete',
    data,
  })
}

export const applyRulesApi = (): Promise<ResponseData<{ msg: string; success: boolean; failReason: string }[]>> => {
  return request({
    url: '/v3/ecsPlatform/warnCustomerRule/applyRules',
    method: 'post',
  })
}
