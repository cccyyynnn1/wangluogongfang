/**
 * 转发规则保存/更新接口类型
 */
export interface forwordRuleUpdateType {
  /**
   * 描述
   */
  description?: string
  /**
   * 过滤条件
   */
  filters?: Filter[]
  /**
   * 转发字段
   */
  forwardFields?: string
  /**
   * 转发服务topic
   */
  forwardTopic: string
  /**
   * 转发服务id
   */
  kafkaConfigId: number | string
  /**
   * 日志类型
   */
  matchTopic: number | string
  /**
   * 名称
   */
  name: string
  /**
   * 站点id
   */
  siteSessionId: string | number
  url?: string
}
interface Filter {
  /**
   * 字段
   */
  field?: string
  /**
   * 关系
   */
  relat?: string
  /**
   * 值
   */
  value?: string
}

export type MaintenanceOptions = {
  /**
   * 攻击方式
   */
  actionType: string[]
  /**
   * 攻击结果
   */
  attackResult: string[]
  /**
   * 攻击类型
   */
  attackType: string[]
  /**
   * 攻击阶段
   */
  killchain: string[]
  /**
   * 攻击手段
   */
  attackMethod: string[]
  /**
   * 规则类型
   */
  ruleType: string[]
  /**
   * 置信度
   */
  confidence: string[]
  /**
   * 代码语言
   */
  codeLanguage: string[]
  /**
   * 威胁等级
   */
  severity: {
    id: number
    dictValue: string
    dictLabel: string
  }[]
  /**
   * 攻击方向
   */
  direction: string[]
}

export type WarnCustomerRuleQueryType = {
  pageNum: number
  pageSize: number
  ruleName: string
  ruleTypes: string[]
  severitys: string[]
  enables: string[]
  attackResults: string[]
  confidences: string[]
}
