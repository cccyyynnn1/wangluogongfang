export interface GetAuditLogQuery {
  pageNum: number
  pageSize: number
  level?: string | number
  /**
   *  0  网络日志
   *  1  系统日志
   *  2  数据库日志
   */
  logType?: 0 | 1 | 2 | 3
  startTime: string
  endTime: string
  scrollId?: string
  orderField: 'timeStamp'
  clientIp?: string
  levelStr?: string
  apiPath?: string
  sysType?: number
  protocol?: number
  userName?: string
  searchStr?: string
}
export interface AuditSystemLog {
  clientIp: string
  logTypeStr: string
  logType: 'sysLog' | 'netLog' | 'dbLog'
  timeStamp: number
  sourceData: string
}
export interface ParamsConfigKeys {
  action: string
  node: {
    key: string
    timestamp: number
    ttl: number
    value: string
  }
}

export interface SecConfigUpdateModel {
  /**
   * 是否包含大写
   */
  capitalChar: boolean
  /**
   * 密码更换周期
   */
  changeDays: number
  /**
   * 修改密码配置
   */
  changePwd: boolean
  /**
   * 令牌失效时间
   */
  expireTime: number
  /**
   * ip白名单
   */
  ipWhilteList: string
  /**
   * 登录失败是否锁定
   */
  loginFailLock: boolean
  /**
   * 限制IP登录
   */
  loginIpLimit: boolean
  /**
   * 登录次数
   */
  loginTimes: number
  /**
   * 密码长度
   */
  minPwdLen: number
  maxPwdLen: number
  /**
   * 是否包含数字
   */
  numChar: boolean
  /**
   * 是否包含小写
   */
  ordinaryChar: boolean
  /**
   * 超时配置
   */
  overtimeConfig: boolean
  /**
   * 是否包含特殊字符
   */
  specialChar: boolean
  /**
   * 自动解锁时间
   */
  unlockTime: number
}

export type BackupFileType = {
  fileName: string
  remark: string
  createTime: number
  type: 'auto' | 'manual'
}

export type AuthToken = {
  name: string
  description: string
  id?: number
  apiIds: number[]
}

export type AuthTokenPage = AuthToken & {
  updateTime: number
  createTime: number
}
