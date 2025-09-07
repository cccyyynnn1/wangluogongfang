/**
 * kafka更新保存接口类型类型
 */
export interface kafkaSaveOrUpdateType {
  acount: string
  addressPorts: string
  description?: string
  id?: number | string
  name: string
  password: string
  host: string
}

export interface kafkaGetGageType {
  kafakName?: string
  pageNum: string | number
  pageSize: string | number
}

/**
 * syslog更新保存接口类型类型
 */
export interface syslogSaveOrUpdateType {
  /**
   * 字段数组
   */
  parseList: warnRuleAnalysFieldsListType[]
  /**
   * 客户端ip
   */
  clientIp: string
  /**
   * true：启用 false：禁用
   */
  enable: number
  id?: number
  /**
   * 样例日志
   */
  logSample: string
  /**
   * 名称
   */
  ruleName: string
  /**
   * 入库类型  0：告警  1：审计日志
   */
  saveIndex: number
  prefixRegex: string
}

/**
 * 字段数组
 */
// export interface AnalysFieldsList {
//   /**
//    * 解析器      0:正则 1：json
//    */
//   parseType: number
//   /**
//    * Key
//    */
//   matchRegex: string
//   matchLog: string
//   /**
//    * 映射字段
//    */
//   fieldMappings?: any[]
//   /**
//    * 正则
//    */
//   appendFields?: any[]
//   /**
//    * 提取值
//    */
//   matchJsonFields?: any[]
// }

export interface syslogGetGageType {
  ruleName?: string
  enable?: number
  pageNum: number
  pageSize: number
}

/**
 * 日志样本
 */
export interface warnRuleParseLogType {
  /**
   * 解析字段
   */
  parseList: warnRuleAnalysFieldsListType[]
  /**
   * 原始日志
   */
  logSample: string
  /**
   * 0:正则 1：json
   */
  // parseType: number
  /**
   * 1：告警  2：审计日志
   */
  saveIndex: number
  prefixRegex: string
}

export interface warnRuleAnalysFieldsListType {
  /**
   * 解析器      0:正则 1：json
   */
  parseType: number
  /**
   * Key
   */
  matchRegex: string
  matchLog: string
  /**
   * 映射字段
   */
  fieldMappings?: any[]
  // fieldMappingsVo?: any[]
  /**
   * 正则
   */
  appendFields?: any[]
  groupSplit?: ''
  kvSplit?: ''
  /**
   * 提取值
   */
  matchJsonFields?: any[]
}
