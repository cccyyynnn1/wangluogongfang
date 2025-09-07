/**
 * 添加资产类型
 */
export interface addAssetsType {
  assetsGroupId: string
  assetsGroupIdStr: string
  assetsIp: string
  assetsTagIds: string[]
  assetsTagIdsStr: string
  assetsTypeIds: string[]
  assetsTypeIdsStr: string
  coreVersion: string
  cpu: string
  dialogTitle: string
  disk: string
  hostName: string
  inManageIp: string
  macAddr: string
  moduleType: string
  name: string
  outManageIp: string
  phone: string
  ram: string
  responsible: string
  sys: string
  sysStr: string
}

/**
 * 添加资产类型
 */
export interface getAssetsType {
  pageNum: number
  pageSize: number
  // query: assetsQueryType | string
  knowOrNot: string
  name?: string | undefined
  assetsIp?: string | undefined
  lastTime?: string | undefined
  responsible?: string | undefined
}

export interface assetsQueryType {
  knowOrNot: string
  name?: string | undefined
  assetsIp?: string | number | undefined
  lastTime?: string | number | undefined
  responsible?: string | undefined
}
export interface ServiceLink {
  alarm: string[]
  clientIp: string
  serverIp: string
  serverPort: {
    [key: number]: {
      appProtocol: string
      protocolStr: string
    }
  }
  snat: string
  type: string
}

export interface ServiceNode {
  IpLocation: string
  coordinate: string
  ip: string
  x?: number
  y?: number
  fixed?: boolean
}

/**
 * 导出资产类型
 */
export interface ExportAssetsType {
  assetsIp?: string
  ids?: number[]
  knowOrNot: string
  lastTime?: string
  name?: string
  responsible?: string
}

/**
 * 网络分区获取页面数据类型
 */
export interface NetworkPartitionType {
  pageNum: number
  pageSize: number
  netTypes: string[]
  dataSources: string[]
  name?: string
  rule?: string
}

/**
 * 站点应用数据类型
 */
export interface ApplicationType {
  appName: string
  appIp: string
  siteId: number | undefined
  id?: string
}
