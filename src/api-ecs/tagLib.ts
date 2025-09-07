import request from '@/utils/request'
import {
  ResponseData,
  tagLibSaveOrUpdateType,
  tagLibGetGageType,
} from '@/types/index'

/**
 * 保存修改
 */
export const tagLibSaveOrUpdateApi = (
  data: tagLibSaveOrUpdateType
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/payloadLabel/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * tagLib分页
 */
export const tagLibGetPageApi = (
  data: tagLibGetGageType
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/payloadLabel/getPage',
    method: 'post',
    data,
  })
}

/**
 * tagLib通过id获取数据
 */
export const tagLibByIdApi = (params: {
  id: string | number
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/payloadLabel/getLabelById',
    method: 'post',
    params,
  })
}

/**
 * tagLib删除数据
 */
export const tagLibDeleteApi = (params: {
  ids: string
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/payloadLabel/delete',
    method: 'delete',
    params,
  })
}

/**
 * 站点检索
 */
export const siteSearchApi = (data: any): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/siteSession/search',
    method: 'post',
    data,
  })
}

/**
 * 获取客户标签
 */
export const getAllServerIpLabelApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/payloadLabel/getAllServerIpLabel',
    method: 'post',
  })
}
