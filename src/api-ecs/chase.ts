import request from '@/utils/request'
import { chasePathTracingResponse, chasePathTracingResquest, QueryAssetByIpRequest, ResponseData } from '../types'

/**
 * @description '路径追踪检索功能'
 * @params chasePathTracingResquest
 */
export const HuntingPathTracingApi = (
  data: chasePathTracingResquest
): Promise<ResponseData<chasePathTracingResponse>> => {
  return request({
    url: '/v3/ecsPlatform/hunting_tracing/huntingPathTracing',
    method: 'post',
    data,
  })
}

/**
 * @description '根据节点IP查询左侧资产信息'
 * @params QueryAssetByIpRequest
 */
export const QueryAssetByIpApi = (
  data: QueryAssetByIpRequest
): Promise<{ msg: string; code: number; data: Record<string, unknown> }> => {
  return request({
    url: '/v3/ecsPlatform/hunting_tracing/queryAssetByIp',
    method: 'post',
    data,
  })
}

/**
 * @description '新增或修改节点颜色'
 * @params { ip: string, remark: string }
 */
export const UpdateOrSaveRemarkApi = (data: {
  ip: string
  remark: string
}): Promise<{ msg: string; code: number; data: Record<string, unknown> }> => {
  return request({
    url: '/v3/ecsPlatform/node_remark/updateOrSaveRemark',
    method: 'post',
    data,
  })
}

/**
 * @description '删除节点颜色
'
 * @params { ip: string }
 */
export const DeleteRemarkApi = (data: {
  ip: string
}): Promise<{ msg: string; code: number; data: Record<string, unknown> }> => {
  return request({
    url: '/v3/ecsPlatform/node_remark/deleteRemark',
    method: 'post',
    data,
  })
}
