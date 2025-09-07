/**
 * 站点检索参数类型
 */
export type siteSearchType = {
  clientPort?: number | string | undefined
  serverPort?: number | string | undefined
  /**
   * 源IP
   */
  clientIp?: string
  /**
   * cookie
   */
  cookie?: string
  /**
   * 结束时间 yyyy-mm-dd hh:mm:ss
   */
  endTime: string
  /**
   * host
   */
  host?: string
  /**
   * 索引类型，暂时使用 1
   */
  indexType: string | number
  /**
   * 排序字段，暂时使用 requestTimeNs
   */
  orderField: string
  /**
   * 排序类型， desc / asc
   */
  orderType: 'desc' | 'asc'
  /**
   * 页数
   */
  pageNum: number
  /**
   * 条数
   */
  pageSize: number
  /**
   * 请求方式
   */
  requestMethod?: string
  /**
   * 请求payload
   */
  requestPayload?: string
  /**
   * 响应payload
   */
  responsePayload?: string
  /**
   * 状态码
   */
  responseStatusCode?: string
  /**
   * 输入框sql
   */
  searchSql?: string
  /**
   * 目的IP
   */
  serverIp?: string
  /**
   * 站点id
   */
  siteSessionId: string | number
  siteApiId?: string | number
  /**
   * 开始时间 yyyy-mm-dd hh:mm:ss
   */
  startTime: string
  /**
   * 标题
   */
  title?: string
  /**
   * URL
   */
  url?: string
  /**
   * User Agent
   */
  userAgent?: string
  /**
   * XFF
   */
  xff?: string
  labelIds?: number[]
  scrollId?: string
  inputSql?: string
  filterSqlArr?: string
  sqlRelat?: string
}

/**
 * 站点列表参数类型
 */
export type siteGetPageType = {
  pageNum: number
  pageSize: number
  hosts?: string
  /**
   * 站点名称
   */
  siteName?: string
  /**
   * 用户
   */
  user?: string
  module: number
}

export interface SiteUnknowPageType {
  host?: string
  ip?: string
  /**
   * 0：源ip  1：目的ip
   */
  ipType?: string
  pageNum: number
  pageSize: number
  position: string
  [property: string]: any
}

/**
 * 聚合检索类型
 */
export type siteAggregationsType = {
  /**
   * 聚合字段
   */
  aggregationFields: string[]
  /**
   * 客户端IP
   */
  clientIp: null | string
  /**
   * cookie信息
   */
  cookie: null | string
  /**
   * 结束时间，2001-01-01 01:00:00
   */
  endTime: string
  /**
   * 主机
   */
  host: null | string
  /**
   * 页码
   */
  pageNum: number
  /**
   * 每页数量
   */
  pageSize: number
  /**
   * 请求方式
   */
  requestMethod: null | string
  /**
   * 请求负载
   */
  requestPayload: null | string
  /**
   * 响应负载
   */
  responsePayload: null | string
  /**
   * 响应状态代码
   */
  responseStatusCode: null | string
  /**
   * 检索语句
   */
  searchSql: null | string
  /**
   * 服务端IP
   */
  serverIp: null | string
  /**
   * 会话ID
   */
  siteSessionId: null | string
  /**
   * 开始时间，2001-01-01 01:00:00
   */
  startTime: string
  /**
   * 标题
   */
  title: null | string
  topCount: number | null
  url: null | string
  /**
   * 用户代理信息
   */
  userAgent: null | string
  /**
   * xff
   */
  xff: number | null
}

/**
 * 删除站点会话类型
 */
export type siteSessionDeleteType = {
  id: number | string
  apiId?: number | string
}
export type filedDataType = {
  id: string | number
  fieldNameCn: string
  fieldNameEn: string
  isDisplay: number
  supportAgg: boolean
}
export type siteCurrentType = {
  createTime: number
  createUser: number
  displayFields: string
  hosts: string
  id: number | string
  module: string
  phone: string
  searchSql: string
  siteName: string
  updateTime: number
  updateUser: number
  user: string
}

/**
 * 新增修改会话类型
 */
export type siteSessionSaveUpdateType = {
  /**
   * 显示字段
   */
  displayFields: filedDataType[] | string | number[]
  /**
   * 主机
   */
  hosts: string
  /**
   * 类型
   */
  module: number | string
  /**
   * 电话
   */
  phone: string
  /**
   * 检索语句
   */
  searchSql: string
  /**
   * 站点名称
   */
  siteName: string
  /**
   * 用户
   */
  user: string
  /**
   * id
   */
  id: number | null | string

  /**
   * 要展示的API
   */
  displayApiIds: number[]
}

/**
 * 会话详情类型
 */
export type siteSessionInfoType = {
  clientIp: string
  clientPort: number
  serverPort: number
  serverIp: string
  requestTimeNs: number
  requestXForwardedFor: number
  responseTimeNs: string
  title: string
  responseStatusCode: number
  requestFullUrl: string
}
export type siteApiInterfaceItemType = {
  apiName: string
  apiUrl: string
  searchSql: string
  createTime: number
  createUser: number
  displayFieldsArr: number[]
  id: number
  sessionId: number
  updateTime: number
  updateUser: number
  /**
   *  0关闭 1 启用
   * **/
  isDisplay: 0 | 1
}
export type siteApiInterfaceType = {
  records: siteApiInterfaceItemType[]
  total: number
}

export type ShortcutListType = {
  key: string
  relation: string
  value?: string
  label: string
  enable: boolean
  id: string
}

export type SiteHistoryItem = {
  collectStatus: boolean
  createTime?: number
  createUser?: number
  dataType?: string
  filterSqlArr?: string
  sqlRelat?: string
  id?: number
  indexType: number | string
  inputSql?: string
  apiId?: number | string
  operateTime?: number
  remarks?: string
  searchEdTime: string
  searchSql: string
  searchStTime: string
  siteId: number | string
  updateTime?: number
  updateUser?: number
  workspaceId?: number
}

export interface siteHistoryParams {
  pageNum: number
  pageSize: number
}

export interface UpdateUnknowSiteResqust {
  /**
   * 展示字段ids
   */
  displayFields: string
  /**
   * hosts
   */
  hosts: string
  /**
   * 电话
   */
  phone?: string
  /**
   * 检索sql
   */
  searchSql?: string
  /**
   * 站点名称
   */
  siteName: string
  /**
   * 未知站点id
   */
  siteUnknowId: number
  /**
   * 负责人
   */
  user?: string
  [property: string]: any
}

export interface updateSiteUnknowGroupByRequest {
  /**
   * 源IP
   */
  clientIp?: string
  /**
   * 统计数量
   */
  groupByCount: number
  /**
   * 统计字段，host:host  目的IP:serverIp
   */
  groupByField: string
  /**
   * host
   */
  host?: string
  /**
   * 方向，内对内,内对外,外对内,外对外
   */
  position?: string
  /**
   * 目的ip
   */
  serverIp?: string
}
