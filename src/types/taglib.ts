/**
 * tagLib更新保存接口类型类型
 */
export interface tagLibSaveOrUpdateType {
  field: string
  indexType: number
  name: string
  relat: string
  value: string
  id?: number | string
}

export interface tagLibGetGageType {
  payloadLabel: PayloadLabel
  query: Query
}

interface PayloadLabel {
  field: string
  indexType?: number
  name: string
}

interface Query {
  pageNum: number
  pageSize: number
}
