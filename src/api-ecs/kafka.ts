import request from '@/utils/request'

import {
  ResponseData,
  kafkaSaveOrUpdateType,
  kafkaGetGageType,
  syslogGetGageType,
  syslogSaveOrUpdateType,
  warnRuleParseLogType,
} from '@/types'

/**
 * 保存修改
 */
export const kafkaSaveOrUpdateApi = (data: kafkaSaveOrUpdateType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/kafkaConfig/saveOrUpdate',
    method: 'post',
    data,
  })
}

/**
 * 获取所有kafka配置信息
 */
export const kafkagetAllApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/kafkaConfig/getAll',
    method: 'post',
  })
}

/**
 * kafka分页
 */
export const kafkaGetPageApi = (data: kafkaGetGageType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/kafkaConfig/getPage',
    method: 'post',
    data,
  })
}

/**
 * kafka通过id获取数据
 */
export const kafkaByIdApi = (params: { id: string | number }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/kafkaConfig/getById',
    method: 'post',
    params,
  })
}

/**
 * @description '系统初始化检测:kafka检测'
 * @params
 */
export const CheckKafkaApi = (data: {
  ip: string
  post: string
  userName: string
  password: string
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/deploymentCheck/checkKafka',
    method: 'post',
    data,
  })
}

/**
 * kafka删除数据
 */
export const kafkaDeleteApi = (params: { ids: string; deleteAll: boolean }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/kafkaConfig/delete',
    method: 'delete',
    params,
  })
}

/**
 * syslog分页
 */
export const syslogGetPageApi = (data: syslogGetGageType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/warnRule/getWarnRulePage',
    method: 'post',
    data,
  })
}

/**
 * syslog删除数据
 */
export const syslogDeleteApi = (
  data: { ids: []; deleteAll: boolean },
  params: { password: string }
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/warnRule/delete',
    method: 'delete',
    data,
    params,
  })
}

/**
 * syslog保存修改
 */
export const syslogSaveOrUpdateApi = (
  data: syslogSaveOrUpdateType,
  params: { password: string }
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/warnRule/saveOrUpdate',
    method: 'post',
    data,
    params,
  })
}

/**
 * syslog保存修改
 */
export const parseLogApi = (data: warnRuleParseLogType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/warnRule/parseLog',
    method: 'post',
    data,
  })
}

/**
 * syslog获取字典
 */
export const getWarnRuleDictApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/warnRule/getWarnRuleDict',
    method: 'post',
  })
}

/**
 * syslog通过id获取数据
 */
export const getByWarnRuleIdApi = (params: { id: number }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/warnRule/getByWarnRuleId',
    method: 'get',
    params,
  })
}
