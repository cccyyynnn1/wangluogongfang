import { isNull } from 'lodash'

/**
 * transpond更新保存接口类型类型
 */
export interface transpondSaveOrUpdateType {
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
  kafkaConfigId: number | any
  /**
   * 日志类型
   */
  matchTopic: number
  /**
   * 名称
   */
  name: string
  /**
   * 站点id
   */
  siteSessionId?: string | number | null
  id?: string | number | null
  url?: string
}

export interface transpondGetGageType {
  pageNum: number
  pageSize: number
}

export interface Filter {
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
