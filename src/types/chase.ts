export interface chasePathTracingResquest {
  aggregationFields: string[]
  count: number
  endTime: string
  indexType: number
  searchSql: string
  startTime: string
}

export interface chasePathTracingResponse {
  linksAll: link[]
  linksDefault: link[]
  nodesAll: node[]
  nodesDefault: node[]
}

export type link = {
  alarm: []
  clientIp: string
  count: string[]
  serverIp: string
  serverPort: { number: { protocolStr: string; appProtocol: string } }
  snat: string
  type: string
  tag: boolean
  source?: string
  target?: string
}

export type node = {
  IpLocation: string
  assetName: null | string
  attribute: string
  businessName: null | string
  color: string
  coordinate: string
  countryCode: string
  remark: string
  icon: string
  ip: string
  id?: string
  ipGroupIdStr: string
  name: string
  tag: boolean
  isUnfold: 0 | 1 | 2
  level?: number
  x?: null | number
  y?: null | number
  fx?: null | number
  vx?: null | number
  vy?: null | number
  fy?: null | number
  isPath?: string
}

export type DrawerSelfRes = {
  isDrag: boolean
  data?: node
}

export interface QueryAssetByIpRequest {
  /**
   * 2024-11-19 04:47:03
   */
  endTime: string
  /**
   * 192.178.1.1
   */
  ip: string
  /**
   * 2024-11-19 04:47:03
   */
  startTime: string
  /**
   * "路径追踪" 或者 "威胁狩猎"
   */
  status: string
}
