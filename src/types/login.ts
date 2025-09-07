/**
 * 登录参数类型
 */
export type LoginPayload = {
  loginType: number
  username: string
  password: string
  verifyCode: string
  requestId: string
}

/**
 * 登录返回类型
 */
export type LoginRes = {
  token: string
}

/**
 * 用户信息类型
 */
export type UserInfoType = {
  username: string
  avatar: string
  roles: string[]
  permissions: any
  address: string
  createTime: null | string
  groupId: number
  id: number
  loginName: string
  mail: string
  nickName: string
  password: string
  status: number
  updateTime: number
  userType: number
  roleIds?: number[]
  theme: string
}

/**
 * 第三方登录
 */
export type LoginLogType = {
  addTime: number
  id: number
  ip: string
  loadTime: null | number
  msg: string
  name: string
  path: string
  proto: string
}

export type LoginOtherType = {
  /**
   * ldap:1,sso:2,radius:3
   */
  authType: number
  /**
   * 写死传1
   */
  loginType: string
  password: string
  /**
   * base64后的值
   */
  username: string
}
