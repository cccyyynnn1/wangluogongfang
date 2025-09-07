import request from '@/utils/request'
import {
  ResponseData,
  QueryUserType,
  QuerySysDeptType,
  SysDeptData,
  EditSysDeptData,
  SysDeptItem,
  QueryUserParams,
  Role,
  AddUserParam,
  SysRole,
  EditSysRole,
  GetAuditLogQuery,
  AuditSystemLog,
  ParamsConfigKeys,
  SecConfigUpdateModel,
  BackupFileType,
  AuthToken,
  AuthTokenPage,
} from '@/types'
import type { VabRouteRecord } from '/#/router'
// 用户管理模块

/**
 * 获取用户列表
 */
export const getUserPageApi = (
  data: QueryUserParams
): Promise<
  ResponseData<{
    records: Role[]
    total: number
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/user/getPage',
    method: 'post',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
  })
}
/**
 * 获取用户列表
 */
export const getAllRoleOrDeptApi = (): Promise<
  ResponseData<{ role: { roleName: string; id: number }[]; dept: { deptName: string; id: number }[] }>
> => {
  return request({
    url: '/v3/ecsPlatform/user/getAllDeptRole',
    method: 'get',
  })
}
/**
 * 新增用户
 */
export const addUserApi = (data: AddUserParam): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/user/saveUser',
    method: 'post',
    data,
  })
}
/**
 * 编辑用户
 */
export const editUserApi = (data: any): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/user/updateUser',
    method: 'post',
    data,
  })
}

/**
 * 删除用户
 */
export const deleteUserApi = (params: { ids: string; deleteAll: boolean }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/user/delete`,
    method: 'delete',
    params,
  })
}
/**
 * 删除用户
 */
export const freeUserApi = (id: string): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/user/freeUser?id=${id}`,
    method: 'get',
  })
}

/**
 * 获取用户信息
 */
export const getUserInfoApi = (data: QueryUserType): Promise<ResponseData<QueryUserType>> => {
  return request({
    url: '/v3/ecsPlatform/user/getPage',
    method: 'post',
    data,
  })
}
/**
 * 获取部门信息
 */
export const getSysDeptApi = (data: QuerySysDeptType): Promise<ResponseData<SysDeptData>> => {
  return request({
    url: '/v3/ecsPlatform/sysDept/getPage',
    method: 'post',
    data,
  })
}
/**
 * 新增部门
 */
export const saveOrUpdateSysDeptApi = (data: EditSysDeptData): Promise<ResponseData<SysDeptData>> => {
  return request({
    url: '/v3/ecsPlatform/sysDept/saveOrUpdate',
    method: 'post',
    data,
  })
}

// 删除部门组织
export const deleteSysDeptApi = (data: { ids: string; deleteAll: boolean }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/sysDept/delete`,
    method: 'delete',
    params: data,
  })
}
// 根据ID获取部门组织
export const getSysDeptByIdApi = (id: number): Promise<ResponseData<SysDeptItem[]>> => {
  return request({
    url: `/v3/ecsPlatform/sysDept/getByDeptId`,
    method: 'post',
    params: { id: id.toString() },
  })
}

// 更新系统配置
export function updateSystemDynamic(data: any, params: { password: string }): Promise<ResponseData<any>> {
  return request({
    url: '/v3/ecsPlatform/settingVar/updateVar',
    method: 'post',
    // headers: {
    //   'Content-Type': 'multipart/form-data',
    // },
    data: data,
    params: params,
  })
}

// 新增系统配置
export function addSystemDynamic(data: any): Promise<ResponseData<any>> {
  return request({
    url: '/v3/ecsPlatform/settingVar/addVar',
    method: 'post',
    data,
  })
}
// 新增系统配置
export function initialSystemApi(
  data: { ip: string; user: string; password: string },
  params: { password: string }
): Promise<ResponseData<any>> {
  return request({
    url: '/v3/ecsPlatform/public/ciInit',
    method: 'post',
    data: data,
    params,
  })
}
// 获取系统角色
export function getSystemRoleApi(data: {
  pageNum: number
  pageSize: number
  roleName: string
}): Promise<ResponseData<{ records: SysRole[]; total: number }>> {
  return request({
    url: 'v3/ecsPlatform/sysRole/getPage',
    method: 'post',
    data,
  })
}

// 获删除系统角色
export function deleteSystemRoleApi(params: { ids: string; deleteAll: boolean }): Promise<ResponseData<any>> {
  return request({
    url: 'v3/ecsPlatform/sysRole/delete',
    method: 'delete',
    params,
  })
}

// 编辑系统角色
export function editSystemRoleApi(data: EditSysRole): Promise<ResponseData<any>> {
  return request({
    url: '/v3/ecsPlatform/sysRole/saveOrUpdate',
    method: 'post',
    data,
  })
}

// 更新系统菜单
export function updateSystemMenusApi(data: VabRouteRecord[]): Promise<ResponseData<any>> {
  return request({
    url: 'v3/ecsPlatform/sysMenu/saveOrUpdate',
    method: 'post',
    data,
  })
}
// 更新系统菜单
export function getSystemMenusApi(): Promise<ResponseData<any>> {
  return request({
    url: '/v3/ecsPlatform/sysMenu/getAll',
    method: 'get',
  })
}
// 更新系统菜单
export function getSystemMenusByRoleApi(): Promise<ResponseData<any>> {
  return request({
    url: '/v3/ecsPlatform/sysMenu/getAllByRole',
    method: 'get',
  })
}

// 设置服务器时间
export function setSystemTimeApi(date: string): Promise<ResponseData<any>> {
  return request({
    url: '/v3/eht/date/setTime',
    method: 'post',
    data: { date },
  })
}
// 设置NTP服务器
export function setNtpdateApi(data: { server: string; enable: boolean }): Promise<ResponseData<any>> {
  return request({
    url: '/v3/eht/date/ntpdate',
    method: 'post',
    data,
  })
}
// 时间设置 升级包上传 需要单独调用接口
export function checkVerificationApi(params: { password: string }): Promise<ResponseData<any>> {
  return request({
    url: 'v3/ecsPlatform/customConfig/checkVerification',
    method: 'get',
    params,
  })
}

// 获取NTP服务器
export function getNtpdateApi(): Promise<ResponseData<{ enable: boolean; server: string }>> {
  return request({
    url: '/v3/eht/date/getNtp',
    method: 'get',
  })
}

// 查询系统日志
export function getAuditLogPageApi(
  data: GetAuditLogQuery
): Promise<ResponseData<{ resList: AuditSystemLog[]; scrollId: string; total: number }>> {
  return request({
    url: '/v3/ecsPlatform/auditLog/getAuditLogPage',
    method: 'post',
    data,
  })
}
// 缓存数据库Key列表
export function getKeysListApi(): Promise<ResponseData<ParamsConfigKeys[]>> {
  return request({
    url: '/v3/eht/cache/keysList',
    method: 'get',
  })
}
// 查询缓存数据库
export function getKeysInfoApi(key: string): Promise<ResponseData<ParamsConfigKeys>> {
  return request({
    url: `/v3/eht/cache/keys/${key}`,
    method: 'get',
  })
}
// 缓存数据库更新
export function postKeysInfoApi(data: string, parames: { key: string; ttl: number }): Promise<ResponseData<any>> {
  return request({
    url: `/v3/eht/cache/keys/${parames.key}?ttl=${parames.ttl}`,
    method: 'post',
    data,
    headers: {
      'Content-Type': 'text/plain',
    },
  })
}
// 缓存数据库删除
export function deleteKeysInfoApi(key: string): Promise<ResponseData<any>> {
  return request({
    url: `/v3/eht/cache/keys/${key}`,
    method: 'delete',
  })
}

// 系统配置
export function getSystemConfigApi(params?: { keys: string }): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/settingVar/getByKeys`,
    method: 'get',
    params,
  })
}

// 更新安全配置
export function SecConfigUpdateApi(
  data: SecConfigUpdateModel,
  params: { password: string }
): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/secConfig/update`,
    method: 'post',
    data,
    params,
  })
}

// 查询安全配置
export function getSecConfigApi(): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/secConfig/getSecConfig`,
    method: 'get',
  })
}

// 系统配置
export function setSystemConfigApi(
  data: { value: string; id?: number; key?: string },
  params: { password: string }
): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/settingVar/updateVar`,
    method: 'post',
    data,
    params: params,
  })
}

// 发送测试邮件
export function mailTestSendApi(params: { mailAddr: string }): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/mail/testSend`,
    method: 'get',
    params,
  })
}
// 发送测试邮件
export function systemUpgrade(): Promise<
  ResponseData<
    {
      date: string
      description: string
      latest: true
      version: string
      module: {
        auth: { filePath: string; md5: string }[]
        ecs: { filePath: string; md5: string }[]
        event: { filePath: string; md5: string }[]
        forward: { filePath: string; md5: string }[]
        forward_java: { filePath: string; md5: string }[]
        host: { filePath: string; md5: string }[]
        html: { filePath: string; md5: string }[]
        loges: { filePath: string; md5: string }[]
        mariadb: { filePath: string; md5: string }[]
        syslog: { filePath: string; md5: string }[]
        thu_wad: { filePath: string; md5: string }[]
        totp: { filePath: string; md5: string }[]
        warn_java: { filePath: string; md5: string }[]
      }[]
    }[]
  >
> {
  return request({
    url: `/v3/ecsPlatform/upgrade/getUpgradeList`,
    method: 'get',
  })
}

/**
 * @description '告警GPT：检查是否可用'
 * @params
 */
export function checkEnableApi(data: { value: string; key?: string }): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/gpt/checkEnable`,
    method: 'post',
    data,
  })
}

// 统一参数配置类型列表
export function getAllCustomConfigApi(): Promise<ResponseData<string[]>> {
  return request({
    url: `/v3/ecsPlatform/customConfig/getAllBelong`,
    method: 'get',
  })
}
// 根据当前类型获取所有参数
export function getParamsByBelongApi(
  type: string,
  search: string
): Promise<
  ResponseData<
    {
      belong: string
      createTime: number
      id: number
      param: string
      remark: string
      value: string
      desensitization: boolean
    }[]
  >
> {
  return request({
    url: `/v3/ecsPlatform/customConfig/findParams?belong=${type}&search=${search}`,
    method: 'post',
  })
}
// 编辑当前参数
export function updateParamsByBelongApi(params: { custom: string; password: string }): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/customConfig/update`,
    method: 'post',
    params,
  })
}

/**
 * @description '下载日历模版'
 * @params
 */
export const CalendartemplateApi = (): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/abnormal_landing_alarm/calendar/export',
    method: 'GET',
    responseType: 'arraybuffer',
  })
}

/**
 * @description '查询日历信息'
 * @params
 */
export function FindAllCalendarApi(params: { year: number; month?: number }): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/abnormal_landing_alarm/findAllCalendar`,
    method: 'post',
    params,
  })
}

/**
 * @description '日历导入'
 * @params
 */
export function CalendarImportApi(
  params: { annual: number; password: string },
  data: { file: File }
): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/abnormal_landing_alarm/calendar/import`,
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data,
    params,
  })
}

/**
 * @description 手动备份任务保存
 * @param {string} remark  备注
 */
export function saveManualBkApi({
  remark,
  password,
}: {
  remark: string
  password: string
}): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/databaseBk/saveManualBk`,
    method: 'post',
    data: { remark },
    params: { password },
  })
}

/**
 * @description 自动保存间隔
 * @param {string} intervalDays 自动保存间隔
 */
export function setIntervalDaysApi({
  intervalDays,
  password,
}: {
  intervalDays: number
  password: string
}): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/databaseBk/updateIntervalDays`,
    method: 'post',
    data: { intervalDays },
    params: { password },
  })
}
/**
 * @description 获取保存间隔
 */
export function getIntervalDaysApi(): Promise<ResponseData<number>> {
  return request({
    url: `/v3/ecsPlatform/databaseBk/getIntervalDays`,
    method: 'get',
  })
}

/**
 * @description 获取所有备份数据
 */
export function getAllBksApi(): Promise<ResponseData<BackupFileType[]>> {
  return request({
    url: `/v3/ecsPlatform/databaseBk/getAllBk`,
    method: 'get',
  })
}

/**
 * @description 恢复数据
 * @param {object} data 还原数据参数
 * @param {string} data.fileName 还原文件名
 * @param {string} data.password 敏感操作名密码
 */
export function rollBkApi({ password, fileName }: { fileName: string; password: string }): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/databaseBk/rollBk`,
    method: 'post',
    data: { fileName },
    params: { password },
  })
}
/**
 * @description 删除数据
 */
export function deleteBkApi({
  password,
  fileNames,
}: {
  password: string
  fileNames: string[]
}): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/databaseBk/deleteBkFile?`,
    method: 'delete',
    data: { fileNames },
    params: { password },
  })
}

/**
 * @description 获取授权信息列表
 */
export function getAuthTokenPageApi(data: {
  pageNum: number
  pageSize: number
}): Promise<ResponseData<{ records: AuthTokenPage[]; total: number }>> {
  return request({
    url: `/v3/ecsPlatform/authToken/getAuthTokenPage`,
    method: 'post',
    data,
  })
}

/**
 * @description 更新授权信息列表
 */
export function updateAuthTokenApi(data: AuthToken, params: { password: string }): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/authToken/saveOrUpdate`,
    method: 'post',
    data,
    params,
  })
}
/**
 * @description 删除授权信息列表
 */
export function deleteAuthTokenApi(
  data: { deleteAll: boolean; ids: number[] },
  params: { password: string }
): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/authToken/delete`,
    method: 'delete',
    data,
    params,
  })
}

/**
 * @description 获取API授权数据树
 */
export function getApiDataTreeApi(): Promise<ResponseData<any>> {
  return request({
    url: `/v3/ecsPlatform/authApi/getApiTree`,
    method: 'get',
  })
}

/**
 * @description 获取授权详情
 */
export function getAuthInfoApi(data: { id: number }): Promise<
  ResponseData<{
    authId: string
    secret: string
    timestamp: number
    token: string
  }>
> {
  return request({
    url: `/v3/ecsPlatform/authToken/getAuthInfo`,
    method: 'post',
    data,
  })
}
