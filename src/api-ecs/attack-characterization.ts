import request from '@/utils/request'
import { AttackCharacterizationType, ResponseData } from '@/types/index'

/**
 * @description 获取特征库
 */
export const getAttackHighlightConfigApi = (data: {
  pageNum: number
  pageSize: number
}): Promise<
  ResponseData<{
    records: AttackCharacterizationType[]
    total: number
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/highLightConfig/getHighLightConfigPage',
    method: 'post',
    data,
  })
}
/**
 * @description 更新特征库
 */
export const updateAttackHighlightConfigApi = (
  data: AttackCharacterizationType & {
    type: 'highLightConfig' | 'highLightWhite'
  }
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/highLightConfig/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 *@description 获取特征白名单
 */
export const getAttackHighlightWhiteApi = (data: {
  pageNum: number
  pageSize: number
}): Promise<ResponseData<{ records: AttackCharacterizationType[]; total: number }>> => {
  return request({
    url: '/v3/ecsPlatform/highLightWhite/getHighLightWhitePage',
    method: 'post',
    data,
  })
}

/**
 * @description 更新特征白名单
 */
export const updateAttackHighlightWhiteApi = (
  data: AttackCharacterizationType & {
    type: 'highLightConfig' | 'highLightWhite'
  }
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/highLightWhite/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * @description 导出特征库和白名单
 */
export const exportHighLightApi = (): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/highLightConfig/exportHighLight',
    method: 'post',
    responseType: 'arraybuffer',
  })
}

/**
 * @description 导出特征库和白名单模板
 */
export const exportHighLightTemplateApi = (): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/highLightConfig/exportTemplate',
    method: 'post',
    responseType: 'arraybuffer',
  })
}

/**
 * @description 导入特征库和白名单
 */
export const importHighLightApi = (data: { file: File }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/highLightConfig/importHighLight',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data,
  })
}

/**
 * @description 删除高亮特征和特征白名单
 */
export const deleteAttackHighLightApi = (data: {
  deleteAll: boolean
  configIds: number[]
  whiteIds: number[]
  password: string
}): Promise<ResponseData<any>> => {
  const { password, ...other } = data
  return request({
    url: `/v3/ecsPlatform/highLightConfig/delete`,
    method: 'delete',
    data: other,
    params: { password },
  })
}
