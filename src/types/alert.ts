export interface AlertItem {
  /** 攻击IP */
  attackIp: string
  host: string
  id: string
  /** 源端口 */
  clientPort: number
  startTimeNs: string
  status: string
  /** 威胁类型 */
  threatType: string
  url: string
  /** 受害ip  */
  victimIp: string
  sourceData?: string
}

export interface PacketItem {
  /**  时间 */
  dateTimeStr: string
  /**  时间差 */
  diffTime: string
  /**  信息 */
  info: string
  /**  大小 */
  lengthStr: string
  /**  序号 */
  num: string
  /**  协议 */
  protocol: string
  /**  相对时间 */
  relativeTime: string
  /**  源IP */
  source: string
  /**  源端口 */
  clientPort: string
  /**  目的IP */
  target: string
  /** 目的端口 */
  targetPort: string
}

export interface PacketDecodeItem {
  clientIp: string
  clientIpLocation: string
  clientPort: string
  probeId: string
  serverIp: string
  serverPort: string
  sessionInfo: string
  totalBytesStr: string
  id: string
}
export interface PacketDecodeQuery {
  clientIp: string
  clientPort: string
  probeId?: string
  probeIds?: number[]
  deviceId?: string
  serverIp: string
  serverPort: string
}

export interface AlertWhiteList {
  clientIp: string
  createTime?: number
  createUser?: number
  id?: number
  ruleId: string
  serverIp: string
  updateUser?: number
}

export interface WarnWhiteSaveOrUpdatesModel {
  cndList: CndList[]
  /**
   * 编辑时必填
   */
  id?: number
  /**
   * 规则id
   */
  ruleId?: string
  /**
   * 规则名称
   */
  ruleName: string
  /**
   * 1：源ip->目的ip 2：源ip->目的ip，目的端口 3：自定义 4：any->any
   */
  type: number
}

export interface CndList {
  /**
   * 字段名
   */
  field: string
  values: string[]
}

export interface SuricataRuleType {
  pageNum: number
  pageSize: number
  /**
   * 攻击结果，字典值
   */
  attackRes?: string[]
  /**
   * 启用状态，1;开启0:关闭
   */
  enable?: string[]
  /**
   * 置信度，字典值
   */
  reliable?: string[]
  /**
   * 规则类型，字典值
   */
  ruleType?: string[]
  /**
   * 规则名称或编号
   */
  searchStr?: string
  /**
   * 告警等级，字典值
   */
  threatLevel?: string[]
  /**
   * 告警间隔
   */
  threatInterval?: string[]
  /**
   * 告警阈值
   */
  threatThreshold?: string[]
}

export interface suricataRuleSaveOrUpdateType {
  tag?: string
  affectAppliance?: string
  /**
   * 规则编号，字典值
   */
  ruleNum?: string
  /**
   * 影响系统，字典值
   */
  affectSystem?: string
  /**
   * 攻击结果，字典值
   */
  attackRes?: string
  /**
   * 攻击方式，字典值
   */
  attackType?: string
  cnnvdNum?: string
  cnvdNum?: string
  /**
   * 代码语言，字典值
   */
  codeLanguage?: string
  cveNum?: string
  /**
   * 方向，字典值
   */
  direction?: string
  /**
   * 启用状态，1;开启0:关闭
   */
  enable?: number
  id?: number
  /**
   * 所属杀伤链，字典值
   */
  killChain?: string
  /**
   * 匹配条件，自定义json字符串
   */
  matchCnd?: string
  /**
   * 置信度，字典值
   */
  reliable?: string
  /**
   * 规则名称
   */
  ruleName: string
  /**
   * 规则类型，字典值
   */
  ruleType?: string
  /**
   * 解决方案
   */
  solution?: string
  /**
   * 告警描述
   */
  threatDesc: string
  /**
   * 告警危害信息
   */
  threatInfo?: string
  /**
   * 告警间隔，自定义json字符串
   */
  threatInterval?: string
  /**
   * 告警等级，字典值
   */
  threatLevel?: string
  /**
   * 告警阈值
   */
  threatThreshold?: string
}

export type InfoTotal = {
  infoCloud: number
  infoCustom: number
  infoWhite: number
}
export type QueryInfoCloud = {
  pageNum: number
  pageSize: number
  searchStr?: string
}
export type InfoCloudItem = {
  id: number
  infoNum: string
  iocStr: string
  iocType: string
  organize: string
  reliable: string | number
  threatDesc: string
  threatLevel: string | number
  threatName: string
  threatType: string
}
export type AddInfoCloud = {
  id?: number
  iocType: string
  organize?: string
  reliable?: string | number
  enable?: 0 | 1
  origin?: string
  threatDesc?: string
  threatLevel?: string | number
  threatName?: string
  threatType?: string
  ip?: string
  domain?: string
  port?: string
  uri?: string
  md5?: string
}

export type InfoCustomOptions = {
  dictLabel: string
  dictValue: string
}[]
export type InfoCustomItem = InfoCloudItem & {
  createTime: number
  createUser: number
  enable: number
  origin: string
  updateTime: number
  updateUser: number
}
export type InfoWhiteItem = {
  createTime: number
  createUserStr: number
  enable: number
  id: number
  remark: string
  updateTime: number
  whiteName: string
  whiteType: string
  whiteTypeStr: string
  checked?: boolean
}
export type AddInfoWhiteItem = {
  whiteName: string
  whiteType: string
  remark?: string
  id?: number
}

export type MailAnalysisChartQuery = {
  indexType: number
  topCount: number
  startTime: string
  endTime: string
  labels: string[]
  labelRelat: 'or' | 'and'
  inOut: boolean
  inIn: boolean
  outIn: boolean
  senderMail: string
  domain: string
  searchStr: string
  threatLevel: number[]
}
export type MailAnalysisItem = {
  aiEmailSum: string
  attachment: string
  body: string
  clientIp: string
  clientPort: number
  emailTags: string[]
  kafkakey: string
  mailClients: string
  msgId: string
  receiverMail: string
  receiverName: string
  senderMail: string
  senderName: string
  serverIp: string
  serverPort: number
  startTimeNs: number
  subject: string
  threatLevel: string
}

export type MailConfigItem = {
  id?: number
  address: string
  name: string
  password: string
  remark: string
  enable: boolean
}

export interface AbnormalLandingAlarmSaveEntityType {
  /**
   * 市
   */
  city: string
  /**
   * 国家
   */
  country: string
  /**
   * 结束时间
   */
  endTime: string
  /**
   * 适配事件
   */
  incident: string
  /**
   * 规则名称
   */
  name: string
  /**
   * 省
   */
  province: string
  /**
   * 目的ip
   */
  serverIp: string
  /**
   * 开始时间
   */
  startTime: string
  /**
   * 时间类型 每天 工作日 等等
   */
  id?: number
  timeStatus: string
  assetArray: string
}

export interface AbnormalLandingAlarmQueryByRegionalType {
  searchSql: string
  startTime: string
  endTime: string
  indexType: number
  whiteType: number
  topCount: number
  aggregationFields: string[]
}

export interface AbnormalLandingAlarmSrcIpOrAddressIpTopNType {
  searchSql: string
  startTime: string
  endTime: string
  indexType: number
  whiteType: number
  topCount: number
  aggregationFields: 'attackIp' | ' victimIp'
}

export interface AbnormalLandingAlarmDistributionOfTimeType {
  attackIp?: string
  sourcePort?: number | string
  xff?: string
  url?: string
  host?: string
  victimIp?: string
  targetPort?: number | string
  threatType?: string //写死
  threatName?: string //写死
  orderField: string
  searchSql: string
  startTime: string
  attackResult?: any
  endTime: string
  orderType: string
  indexType: number
  pageNum: number
  pageSize: number
  whiteType: number
}

export interface AbnormalLandingAlarmSearchType {
  searchSql: string
  attackIp?: string
  sourcePort?: number | string
  xff?: string
  url?: string
  host?: string
  victimIp?: string
  targetPort?: number | string
  threatType?: string
  threatName?: string
  threatLevel?: string
  attackResult?: any
  readStatus?: number
  orderField: string
  orderType: string
  indexType: number
  pageNum: number
  pageSize: number
  startTime: string
  endTime: string
  whiteType: number
  total?: number
}
