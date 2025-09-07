/**
 * 提取规则更新保存接口类型类型
 */
export interface normalizeSaveOrUpdateType {
  analysisFieldsStr: string
  analysisType: string
  createTime: number
  createType: string
  description: string
  field: string
  groupId: number | undefined
  id?: number
  indexType: string
  logSample: string
  name: string
  requestHost: string
  requestUrl: string
  status: string | number
  updateTime: number
  groupSplit: string
  kvSplit: string
  analysFields: any
  newFieldsStr: string | undefined
}

export interface normalizeGetGageType {
  pageNum: number
  pageSize: number
  groupId?: number | string | undefined
}
