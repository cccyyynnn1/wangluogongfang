export interface AttackStageDict {
  createTime: number
  createUser: number
  dictLabel: string
  dictTypeId: number
  dictValue: string
  enable: number
  id: number
  isDefault: number
  remark: string
  updateTime: number
  updateUser: number
}
export interface TemplateItem {
  createTime: number
  createUser: number
  id: number
  templateName: string
  configStr: string
  updateTime: number
  updateUser: number
}
export type ThreatLevelDict = AttackStageDict
export type ThreatTypeDict = AttackStageDict
export interface FlowDeviceItem {
  label: string
  value: string
  children: FlowDeviceItem[]
}

export type SyslogTemplate = {
  /**
   * 攻击阶段字典
   */
  attackStageDict: AttackStageDict[]
  /**
   * 日志模版列表
   */
  templateList: TemplateItem[]
  /**
   * 攻击等级字典
   */
  threatLevelDict: ThreatLevelDict[]
  /**
   * 攻击类型字典
   */
  threatTypeDict: string[]
  flowDeviceList: FlowDeviceItem[]
}

export interface ApplyTemplate {
  id?: number
  logType: number
  ip: string
  status: 0 | 1 | 2
  port: number
  protocol: number
  templateId?: number
  ouputSource: number
  configStr: string
  outputNull: number
  threatTypeIds: number[]
  attackStageIds: number[]
  threatLevelIds: number[]
  probeIps: string[]
  fieldsEchoInfos: string
}

export type TemplateRuleItem = ApplyTemplate
