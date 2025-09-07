import request from '@/utils/request'
import { ResponseData, LoginPayload, LoginRes, UserInfoType, LoginLogType, LoginOtherType } from '@/types/index'
import { encryptedData } from '@/utils/encrypt'
import { loginRSA } from '@/config'

/**
 * 登录
 */
export const loginApi = async (data: LoginPayload): Promise<ResponseData<LoginRes>> => {
  if (loginRSA) {
    data = await encryptedData(data)
  }
  return request({
    url: '/v3/authority/auth/login',
    method: 'post',
    data,
  })
}

/**
 * 退出
 */
export const logoutApi = async () => {
  return request({
    url: '/v3/authority/auth/loginOut',
    method: 'get',
  })
}

/**
 * 获取用户信息
 */
export const getUserInfoApi = (): Promise<ResponseData<UserInfoType>> => {
  return request({
    url: '/v3/authority/auth/getCurUser',
    method: 'get',
  })
}

/**
 * 获取系统信息
 */
export const getSysInfoApi = (): Promise<ResponseData<{ value: string }>> => {
  return request({
    url: '/v3/ecsPlatform/public/getSysInfo',
    method: 'get',
  })
}

/**
 * 获取登录配置
 */
export const getOtherLoginStatusApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/settingVar/getOtherLoginStatus',
    method: 'post',
    // headers: {},
  })
}

/**
 * 获取验证码开关信息
 */
export const getVerifyImgStatusApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/authority/auth/getVerifyImgStatus',
    method: 'get',
  })
}

/**
 * 其他方式登录
 */
export const loginOtherTypeApi = (data: LoginOtherType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/authority/auth/loginOtherType',
    method: 'post',
    data,
    // headers: {},
  })
}

/**
 * 重置密码
 */
export const resetPasswordApi = (data: { password: string; oldPassword: string }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/user/changPwd',
    method: 'post',
    data,
  })
}
/**
 * 更新主题
 */
export const updateThemeApi = (data: { theme: string; id: number }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/user/updateTheme',
    method: 'post',
    data,
  })
}
/**
 * 获取登陆日志
 */
export const getAdminLogPageApi = (data: {
  pageNum: number
  pageSize: number
}): Promise<
  ResponseData<{
    records: LoginLogType[]
    total: number
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/adminLog/getAdminLogPage',
    method: 'post',
    data,
  })
}
