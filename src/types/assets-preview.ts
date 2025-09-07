/**
 * @description 资产预览分页查询参数
 * @param pageNum 页码
 * @param pageSize 页数
 * @param dataSourceId 数据中心
 * @param activeAssetType 活跃资产类型  0-全部 1-活跃 2-不活跃
 * @param labelIds  标签ids
 * @param lastTime 最后在线时间
 * @param searchStr 模糊查询
 */
export interface PreviewAssetsBodyType {
  pageNum: number
  pageSize: number
  dataSourceId: number[]
  activeAssetType?: number
  labelIds?: number[]
  lastTime?: string
  searchStr: string
  searchIps?: string[]
  otherDataSource: boolean
  labelRelat?: string
  serverIps?: string[]
  originImport?: 1 | null
  originFlow?: 1 | null
  know?: 0 | 1 | null
  conflict?: 0 | 1 | null
  netType?: 0 | 1 | null
  isOffline?: true | false | null
  containShandowAsset?: 0 | 1 | null
}

export interface PreviewAssetsItem {
  appName: string
  appType: string
  businessName: string
  businessType: string
  id?: number
  indexIp?: number[]
  ipv4?: string
  ipv4Num: string
  ipv6: string
  labels: CMDBfingerprintList[]
  middlewareVersion: string
  netPartation: string
  otherInfo: string
  portList: number[]
  hostVoList: {
    apis: string[]
    businessInfo: string[]
    host: string
    know: boolean
  }[]
  otherInfoObj: {
    apiList: {
      businessPort: number
      urls: string[]
    }[]
    serverPorts: number[]
  }
  phone: string
  protocol: number
  sysVersion: string
  dbVersion: string
  datasource: string
  host: string
  historyList: any[]
}
export type FingerprintItem = {
  groupId: number
  id: number
  labelName: string
  isNew: boolean
}
export type CMDBfingerprintList = {
  groupId: 1 | 2 | 3
  groupName: string
  labelList: FingerprintItem[]
  labelSize: number | null
}
export type NetPartitionSums = {
  count: number
  dataSourceId: number
  dataSourceName: string
  ipv4Ed: string
  ipv4St: string
  ipv6Ed: string
  ipv6St: string
}[]

/**
 * @description 资产预览分页查询参数
 * @param activityAssetNum 活跃数量
 * @param notActivityAssetNum  非活跃数量
 * @param labelVoList 指纹列表
 * @param netPartitionSums 数据中心
 */
export interface PreviewAssetsSummary {
  activityAssetNum: number
  labelVoList: CMDBfingerprintList[]
  netPartitionSums: NetPartitionSums
  notActivityAssetNum: number
  otherNetPartationCount: number
  topList: FingerprintItem[]
}

/**
 * @description 更新资产基本信息
 */
export interface UpdateAssetInfo {
  appName?: string
  appType?: string
  businessName?: string
  businessType?: string
  datasource?: string
  dbVersion?: string
  id: number
  middlewareVersion?: string
  netPartation?: string
  phone?: string
  protocol?: number
  sysVersion: string
}

export interface pageChartsType {
  /**
   * 活跃资产类型 0-全部 1-活跃  2-不活跃
   */
  activeAssetType?: number
  /**
   * 数据中心ids
   */
  dataSourceId?: number[]
  /**
   * 结束时间，yyyy-MM-dd HH:mm:ss
   */
  endTime: string
  /**
   * 标签ids
   */
  labelIds?: number[]
  /**
   * or and
   */
  labelRelat?: string
  /**
   * 最近在线时间 2023-12-24 10:31:09
   */
  lastTime?: string
  /**
   * 是否 非当前数据中心的
   */
  otherDataSource: boolean
  /**
   * 模糊检索ip
   */
  searchIps?: string[]
  /**
   * 名称模糊检索
   */
  searchStr?: string
  /**
   * 目的ip
   */
  serverIps?: string[]
  /**
   * 开始时间，yyyy-MM-dd HH:mm:ss
   */
  startTime: string
  /**
   * 威胁等级code
   */
  threatLevelCode?: number
  /**
   * 威胁名称
   */
  threatName?: string
  /**
   * 威胁类型
   */
  threatTypes?: string[]
  originImport?: 1 | null
  originFlow?: 1 | null
  know?: 0 | 1 | null
  conflict?: 0 | 1 | null
  netType?: 0 | 1 | null
  isOffline?: true | false | null
  containShandowAsset?: 0 | 1 | null
}

export interface detailChartsType {
  /**
   * 结束时间，yyyy-MM-dd HH:mm:ss
   */
  endTime: string
  /**
   * 目的ip
   */
  serverIps: string | number
  /**
   * 开始时间，yyyy-MM-dd HH:mm:ss
   */
  startTime: string
}
