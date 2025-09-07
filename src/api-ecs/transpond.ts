import request from '@/utils/request'
import { ResponseData, transpondSaveOrUpdateType, transpondGetGageType } from '@/types/index'

/**
 * 转发规则保存修改
 */
export const transpondSaveOrUpdateApi = (data: transpondSaveOrUpdateType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/forwardRules/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * 转发规则分页
 */
export const transpondGetPageApi = (data: transpondGetGageType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/forwardRules/getPage',
    method: 'post',
    data,
  })
}

/**
 * 转发规则删除数据
 */
export const transpondDeleteApi = (params: { ids: string; deleteAll: boolean }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/forwardRules/delete',
    method: 'delete',
    params,
  })
}
