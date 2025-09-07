import request from '@/utils/request'
import {
  ResponseData,
  customFieldSaveOrUpdateType,
  customFieldGetGageType,
  RetrieveIndexType,
  AssetLabelDict,
  InitialisationList,
  TemplateInfo,
} from '@/types/index'

/**
 * 自定义字段保存修改
 */
export const customFieldSaveOrUpdateApi = (data: customFieldSaveOrUpdateType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/indexField/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * 自定义字段分页
 */
export const customFieldGetPageApi = (data: customFieldGetGageType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/indexField/getPage',
    method: 'post',
    data,
  })
}

/**
 * 自定义字段删除数据
 */
export const customFieldDeleteApi = (params: { ids: string; deleteAll: boolean }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/indexField/delete',
    method: 'delete',
    params,
  })
}

/**
 * 通过类型获取自定义字段
 */
export const customFieldGetByTypeApi = (params: { type: string }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/indexField/getByType',
    method: 'get',
    params,
  })
}

/**
 * 更新用户字段
 */
export const updateDisplayApi = (data: {
  indexType: RetrieveIndexType
  displayIds: number[]
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/indexField/updateDisplay',
    method: 'post',
    data,
  })
}

/**
 * 获取所有字典
 */
export const getAllAssetDictApi = (params?: { searchStr: string }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/assetLabelDict/getAllAssetDict',
    method: 'get',
    params,
  })
}

/**
 * 更新字典信息
 */
export const updateAssetDictApi = (data: AssetLabelDict): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/assetLabelDict/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * 字段：默认配置（查询）
 */
export const indexFieldsConfigApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/indexFieldsConfig/queryAll',
    method: 'post',
  })
}

/**
 * 重置为默认初始化
 */
export const resetIndexFieldsConfigApi = (data: { indexType: number }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/indexFieldsConfig/restData',
    method: 'post',
    data,
  })
}

/**
 * 默认配置（修改） 只会修改当前登陆人的
 */
export const indexFieldsConfigUpdateApi = (data: {
  indexType: number
  displayIds: number[]
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/indexFieldsConfig/update',
    method: 'put',
    data,
  })
}

/**
 * 根据type重置展示字段
 */
export const resetDisPlaysFiledApi = (data: { indexType: number }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/indexFiledsUser/resetDisPlaysFiled',
    method: 'post',
    data,
  })
}

/**
 * 重置站点展示字段
 */
export const resetSiteDisPlaysFiledApi = (data: {
  indexType: string
  siteApiId: number | string
  siteSessionId: number | string
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/indexFiledsUser/resetDisPlaysFiled',
    method: 'post',
    data,
  })
}

/**
 * 根据ID获取模版字段详情
 */
export const getTemplateInfoApi = (data: { id: number }): Promise<ResponseData<TemplateInfo>> => {
  return request({
    url: '/v3/ecsPlatform/indexTemplate/findById',
    method: 'post',
    data,
  })
}

/**
 * 对比模版列表
 */
export const getComparisonFieldsTemplateApi = (): Promise<ResponseData<InitialisationList>> => {
  return request({
    url: '/v3/ecsPlatform/indexTemplate/getAllTemplate',
    method: 'post',
  })
}

/**
 * 字段初始化模版已删除列表
 */
export const getDeleteFieldsTemplateApi = (data: {
  templateName: string
  pageNum: number
  pageSize: number
}): Promise<
  ResponseData<{
    total: number
    records: InitialisationList
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/indexTemplate/findAllByIsDelete',
    method: 'post',
    data,
  })
}

/**
 * 字段初始化模版列表
 */
export const getFieldsTemplateApi = (data: {
  templateName: string
  pageNum: number
  pageSize: number
}): Promise<
  ResponseData<{
    total: number
    records: InitialisationList
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/indexTemplate/queryAll',
    method: 'post',
    data,
  })
}
/**
 * 启用字段初始化模版
 */
export const usingFieldsTemplateApi = (data: {
  id: number
  status: 0 | 1
}): Promise<
  ResponseData<{
    total: number
    records: InitialisationList
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/indexTemplate/updateStatus',
    method: 'post',
    data,
  })
}
/**
 * 删除字段初始化模版
 */
export const deleteFieldsTemplateApi = (
  ids: number[]
): Promise<
  ResponseData<{
    total: number
    records: InitialisationList
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/indexTemplate/logicallyDelete',
    method: 'delete',
    data: {
      ids,
    },
  })
}
/**
 * 字段重置
 */
export const fieldsTemplateRestApi = (data: {
  type: number
  siteSessionId: number | ''
  siteApiId: number | ''
  tag: 'alarm' | 'survey' | 'siteSession' | 'siteApi'
}): Promise<
  ResponseData<{
    namesList: string[]
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/indexTemplate/updateRest',
    method: 'post',
    data,
  })
}
/**
 * 更新字段初始化模版
 */
export const fieldsTemplateUpdateApi = (
  data: {
    id: number
    templateName: string
    remark: string
  } & TemplateInfo
): Promise<
  ResponseData<{
    namesList: string[]
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/indexTemplate/update',
    method: 'post',
    data,
  })
}
/**
 * 新增字段初始化模版
 */
export const fieldsTemplateAddApi = (
  data: {
    templateName: string
    remark: string
  } & TemplateInfo
): Promise<
  ResponseData<{
    namesList: string[]
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/indexTemplate/save',
    method: 'post',
    data,
  })
}

/**
 * 恢复已删除字段初始化模版
 */
export const fieldsTemplateRecoverApi = (data: {
  ids: number[]
}): Promise<
  ResponseData<{
    namesList: string[]
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/indexTemplate/recover',
    method: 'post',
    data,
  })
}

/**
 * 彻底删除字段初始化模版
 */
export const fieldsTemplateHistoryDeleteApi = (data: {
  ids: number[]
}): Promise<
  ResponseData<{
    namesList: string[]
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/indexTemplate/delete',
    method: 'delete',
    data,
  })
}

/**
 * 获取对比字段
 */
export const getComparisonFieldsApi = (data: {
  /**@对比模版ID*/
  indexFiledsTemId: number
  type: number
  siteSessionId: number | ''
  siteApiId: number | ''
  tag: 'alarm' | 'survey' | 'siteSession' | 'siteApi'
}): Promise<
  ResponseData<{
    namesList: string[]
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/indexTemplate/comparison',
    method: 'post',
    data,
  })
}

/**
 * 获取对比字段
 */
export const contrastReplacementFieldsApi = (data: {
  /**@对比模版ID*/
  indexFiledsTemId: number
  type: number
  siteSessionId: number | ''
  siteApiId: number | ''
  tag: 'alarm' | 'survey' | 'siteSession' | 'siteApi'
  names: string
}): Promise<
  ResponseData<{
    namesList: string[]
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/indexTemplate/contrastReplacement',
    method: 'post',
    data,
  })
}
