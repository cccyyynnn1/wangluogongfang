import { TableColumnItemType } from '/#/store'

export type RetrieveIndexType = number
// export type IndexTypes = keyof typeof RetrieveIndexType
export type OrderType = 'asc' | 'desc'
/**
 * 用户信息类型
 */
export interface SearchBySqlParams {
  searchSql: string
  indexType: number
  pageNum: number
  pageSize: number
  orderType: OrderType | string
  orderField: string
  startTime: string
  endTime: string
  workspaceId?: number
  scrollId?: string
  inputSql?: string
  filterSqlArr?: string
  sqlRelat?: string
}

export interface NetworkLayerParams {
  searchSql: string
  indexType: number
  pageNum: number
  pageSize: number
  orderType: OrderType
  orderField: string
  startTime: string
  endTime: string
  scrollId?: string
}

export interface NewNetworkLayerParams {
  timeRange: string
  flowProbeIds: string[]
  searchTable: string
  searchModel: string
}

export interface PublishRuleParams extends SearchBySqlParams {
  displayFields: string[] | TableColumnItemType[]
  ruleName: string
  searchTime: number
  id?: number
}

export interface RuleModelItem {
  createTime: number
  displayFields: TableColumnItemType[]
  id: number
  indexType: string
  pageSize: number
  ruleName: string
  searchSql: string
  searchTime: number
  updateTime: number
}
export interface RuleModel {
  total: number
  records: RuleModelItem[]
}

export interface WorkerSpaceItem {
  id: number
  spaceName: string
  default: boolean
}
export interface RetrieveHistoryParams {
  pageNum: number
  pageSize: number
  /**
   * trur:收藏 false：检索历史
   */
  collectStatus: boolean
  workspaceId?: number
}
export interface FindHistoryParams {
  searchSql: string
  indexType: number
  collectStatus: boolean
  workspaceId: number
}
export interface ShareHistoryParams {
  pageNum: number
  pageSize: number
  sId?: number
  rName?: string
}
export interface RetrieveHistoryItem {
  /**
   * trur:收藏 false：检索历史
   */
  collectStatus: boolean
  dataType?: string
  id?: number
  indexType: number
  operateTime?: number
  remarks?: string
  searchEdTime: string
  searchStTime: string
  searchSql: string
  updateUser?: string
  workspaceId: number
  inputSql?: string
  filterSqlArr?: string
  sqlRelat?: string
}

export interface LevelRuleSaveUpdateModel {
  /**
   * 条件.json字符串
   */
  cnd: string
  /**
   * 更新时必填
   */
  id?: number
  /**
   * 索引类型
   */
  indexType: number
  /**
   * 0：普通 1：一般 2：严重
   */
  level: number | string
  /**
   * 名称
   */
  ruleName: string
}
