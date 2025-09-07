import request from '@/utils/request'
import { ResponseData, SyslogTemplate, TemplateItem, ApplyTemplate, TemplateRuleItem } from '@/types/index'
/**
 * 获取日志输出模版字段参数
 */
export const getAllSyslogTemplateApi = (): Promise<ResponseData<SyslogTemplate>> => {
  return request({
    url: '/v3/ecsPlatform/syslogTemplate/getAllSyslogTemplate',
    method: 'post',
  })
}

/**
 * 保存日志输出模版
 */
export const saveOrUpdateTemplateApi = (data: string): Promise<ResponseData<TemplateItem>> => {
  return request({
    url: '/v3/ecsPlatform/syslogTemplate/saveOrUpdate',
    method: 'post',
    data,
  })
}
/**
 * 应用日志模版
 */
export const saveOrUpdateApplyTemplateApi = (data: ApplyTemplate): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/syslogForwardNew/saveUpdate',
    method: 'post',
    data,
  })
}
/**
 * 应用日志模版
 */
export const getTemplateRuleListApi = (data: {
  pageNum: number
  pageSize: number
  searchStr: string
}): Promise<ResponseData<{ records: TemplateRuleItem[]; total: number }>> => {
  return request({
    url: '/v3/ecsPlatform/syslogForwardNew/getForwardNewPage',
    method: 'post',
    data,
  })
}

/**
 * 应用日志模版
 */
export const deleteTemplateRulestApi = (data: {
  deleteAll: boolean
  ids: number[]
}): Promise<ResponseData<{ records: TemplateRuleItem[]; total: number }>> => {
  return request({
    url: '/v3/ecsPlatform/syslogForwardNew/delete',
    method: 'delete',
    data,
  })
}
/**
 * 应用日志模版
 */
export const updateTemplateRulestStatusApi = (data: { status: number; ids: number[] }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/syslogForwardNew/updateStatus',
    method: 'post',
    data,
  })
}
