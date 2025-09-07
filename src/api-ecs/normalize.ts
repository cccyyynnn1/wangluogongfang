import request from '@/utils/request'
import { ResponseData, normalizeSaveOrUpdateType, normalizeGetGageType } from '@/types/index'

/**
 * 提取规则保存修改
 */
export const normalizeSaveOrUpdateApi = (data: normalizeSaveOrUpdateType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/normalize/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * 提取规则分页
 */
export const normalizeGetPageApi = (data: normalizeGetGageType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/normalize/getPage',
    method: 'post',
    data,
  })
}

/**
 * 提取规则删除数据
 */
export const normalizeDeleteApi = (params: { ids: string; deleteAll: boolean }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/normalize/delete',
    method: 'delete',
    params,
  })
}

/**
 * 提取规则字段提取
 */
export const normalizeParseLogApi = (data: {
  logSample: string
  field: number | string
  kvSplit?: string
  groupSplit?: string
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/normalize/parseLog',
    method: 'post',
    data,
  })
}

/**
 * 提取规则得到所有分组
 */
export const normalizeGroupGetAllApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/normalizeGroup/getAll',
    method: 'GET',
  })
}

/**
 * 提取规则获取当前行数据
 */
export const normalizeGetRowByIdApi = (params: { id: string | number }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/normalize/getById',
    method: 'GET',
    params,
  })
}
