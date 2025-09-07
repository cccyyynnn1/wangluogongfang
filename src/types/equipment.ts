/**
 * 设备列表接口类型
 */
export interface EquipmentListType {
  query?: {}
  page: string | number
  limit: string | number
}

/**
 * 更新设备类型
 */
export interface UpdateEquipmentType {
  dataTransactionSecond: number // 数据传输合并间隔
  dbMemoryLimitGb: number // 数据库内存限制
  eventStatSecond: number // 事件统计间隔
  isNuma: boolean // 是否开启numa
  needRestartService: boolean // 是否需要重启服务
  systemTime: boolean // 是否同步系统时间
  zipStorage: boolean // 启⽤压缩存储
  logLevelStr: string // ⽇志级别
  dialogTitle: string
  moduleType: string
}

/**
 * 链路列表接口类型
 */
export interface FlowProbeListType {
  page: string | number
  limit: string | number
  query: any
}
