import request from '@/utils/request'
import {
  AlertItem,
  PacketDecodeItem,
  ServiceLink,
  ServiceNode,
  ResponseData,
  AlertWhiteList,
  WarnWhiteSaveOrUpdatesModel,
  SuricataRuleType,
  suricataRuleSaveOrUpdateType,
  InfoTotal,
  QueryInfoCloud,
  InfoCloudItem,
  InfoCustomItem,
  InfoWhiteItem,
  InfoCustomOptions,
  AddInfoCloud,
  AddInfoWhiteItem,
  MailAnalysisChartQuery,
  MailAnalysisItem,
  MailConfigItem,
  AbnormalLandingAlarmSaveEntityType,
  AbnormalLandingAlarmQueryByRegionalType,
  AbnormalLandingAlarmSrcIpOrAddressIpTopNType,
  AbnormalLandingAlarmDistributionOfTimeType,
  AbnormalLandingAlarmSearchType,
} from '@/types'

/**
 * 获取告警列表
 */
export const getAlertApi = (
  data: any
): Promise<
  ResponseData<{
    /** 表格数据 */
    resList: AlertItem[]
    /** 图表数据 */
    sumaryMap: string
    total: number
    scrollId: string
    aggOneCount: number
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/warnSearch/search',
    method: 'post',
    data,
  })
}
/**
 * 更新告警状态
 */
export const updateAlertStatusApi = (data: {
  startTime: string
  endTime: string
  indexType: number
  readStatus?: '0' | '1' | number
  ignoreStatus?: number | string
  color?: number
  ids: string[]
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/warnSearch/updateStatus',
    method: 'post',
    data,
  })
}
/**
 * 获取告警详情
 */
export const getAlertDetailApi = (
  data: {
    attackIp: string
    sourcePort: number
    victimIp: string
    targetPort: number
    startTimeNs: number
  }[]
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/warnSearch/getPcapDownloadInfo',
    method: 'post',
    data,
  })
}

/**
 * 获取告警详情
 */
export const getAggregationsApi = (
  data: {
    indexType: number
    startTime: string
    endTime: string
    aggregationFields: string
    topCount: number
    searchSql?: string
  }[]
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/event/aggregations',
    method: 'post',
    data,
  })
}
/**
 * 在线解码
 */
export const getPacketDecodeApi = (query: {
  timeRange: string
  ipAddr?: string
  serverIp?: string
  clientIp?: string
  serverPort?: string
  top: number
}): Promise<ResponseData<PacketDecodeItem[]>> => {
  return request({
    url: `/v3/ecsPlatform/public/flowSearch`,
    method: 'post',
    data: { searchSql: JSON.stringify(query) },
  })
}

/**
 * 获取数据包List
 */
export const getPacketDecodeListApi = (data: {
  top: number
  query: {
    objectList: {
      serverIp: string
      serverPort: string
      clientIp: string
      clientPort: string
      probeId: string
    }[]
    timeStep: string
    timeRange: string
  }
}): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/public/packetList`,
    method: 'post',
    data: { ...data, packetCut: '0' },
  })
}

/**
 * 获取数据包字节流
 */
export const getPacketDecodeDetailApi = (id: string): Promise<ResponseData<{ data: { hex: string } }>> => {
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/packetDecode/detail/${id}`,
    },
  })
}

/**
 * 数据流解码
 */
export const getPacketDecodeFlowDecodeApi = (data: {
  top: number
  flowType: string
  query: {
    objectList: {
      serverIp: string
      serverPort: string
      clientIp: string
      clientPort: string
      probeId: string
    }[]
    timeStep: string
    timeRange: string
  }
}): Promise<ResponseData<{ hexs: string[]; info: string; srcToDst: boolean }[]>> => {
  return request({
    url: `/v3/ecsPlatform/public/flowDecode`,
    method: 'post',
    data: { ...data, packetCut: '0' },
  })
}
/**
 * 获取攻击透视图标数据
 */
export const getAttackPerspectiveApi = (data: {
  ip: string
  startDate: string
  endDate: string
  limit: number
}): Promise<
  ResponseData<{
    links: ServiceLink[]
    nodes: ServiceNode[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/hugegraph/queryByAttackPerspec`,
    method: 'post',
    data,
  })
}
/**
 * 获取告警白名单
 */
export const getWarnWhiteApi = (data: {
  searchStr?: string
  pageSize: number
  pageNum: number
}): Promise<
  ResponseData<{
    records: AlertWhiteList[]
    total: number
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/warnWhite/getWarnWhitePage`,
    method: 'post',
    data,
  })
}

/**
 * 删除告警白名单
 */
export const getDeleteWarnWhiteApi = (data: { ids: number[]; deleteAll: boolean }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/warnWhite/deleteWarnWhilte`,
    method: 'delete',
    data,
  })
}

/**
 * 更新告警白名单
 */
export const getUpdateWarnWhiteApi = (data: AlertWhiteList): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/warnWhite/saveOrUpdate`,
    method: 'post',
    data,
  })
}

/**
 * 下载白名单模版
 */
export const exportWarnWhiteTemplateApi = (): Promise<BlobPart> => {
  return request({
    url: `/v3/ecsPlatform/warnWhite/exportTemplate`,
    method: 'get',
    responseType: 'arraybuffer',
  })
}

/**
 * 导入白名单模版
 */
export const importWarnWhiteTemplateApi = (data: { file: File }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/warnWhite/importWarnWhite`,
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data,
  })
}

/**
 * 导出白名单
 */
export const exportWarnWhiteApi = (data: { searchStr?: string; ids?: string | number[] }): Promise<BlobPart> => {
  return request({
    url: `/v3/ecsPlatform/warnWhite/exportWarnWhite`,
    method: 'post',
    data,
    responseType: 'arraybuffer',
  })
}

/**
 * 删除告警白名单
 */
export const DeleteWarntimeApi = (data: { ids: number[]; deleteAll: boolean }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/warnWhiteExportBk/delete`,
    method: 'delete',
    data,
  })
}

/**
 * 获取下拉配置字段
 */
export const GetConfigFieldsApi = (): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/warnWhite/getConfigFields`,
    method: 'get',
  })
}

/**
 * 告警加白：保存/更新
 */
export const WarnWhiteSaveOrUpdatesApi = (data: WarnWhiteSaveOrUpdatesModel): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/warnWhite/saveOrUpdate`,
    method: 'post',
    data,
  })
}

/**
 * 获取告警白名单
 */
export const getWarnWhiteListApi = (data: {
  pageSize: number
  pageNum: number
}): Promise<
  ResponseData<{
    records: any
    total: number
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/warnWhiteExportBk/getWarnWhiteBkPage`,
    method: 'post',
    data,
  })
}

/**
 * suricataRule规则分页
 */
export const getSuricataRulePageApi = (data: SuricataRuleType): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/warnCustomerRule/getWarnCustomerRulePage`,
    method: 'post',
    data,
  })
}

/**
 * suricataRule获取所有下拉
 */
export const getAllSelectDataApi = (): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/warnCustomerRule/getAllSelectData`,
    method: 'post',
  })
}

/**
 * suricataRule新增或修改
 */
export const suricataRuleSaveOrUpdateApi = (data: suricataRuleSaveOrUpdateType): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/warnCustomerRule/saveOrUpdate`,
    method: 'post',
    data,
  })
}

/**
 * suricataRule修改状态
 */
export const suricataRuleUpdateStatusApi = (data: { status: number; ids: number[] }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/warnCustomerRule/updateStatus`,
    method: 'post',
    data,
  })
}

/**
 * suricataRule删除
 */
export const suricataRuleDelApi = (data: { ids?: Array<string>; deleteAll: boolean }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/warnCustomerRule/delete',
    method: 'DELETE',
    data,
  })
}

/**
 * 获取所有情报总数
 */
export const getIntelligenceInfoTotalApi = (): Promise<ResponseData<InfoTotal>> => {
  return request({
    url: `/v3/ecsPlatform/infoCloud/getInfoTotal`,
    method: 'get',
  })
}
/**
 * 获取云端情报列表
 */
export const getInfoCloudPageApi = (data: {
  searchStr: string
  pageNum: number
  pageSize: number
}): Promise<
  ResponseData<{
    total: number
    records: InfoCloudItem[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/infoIoc/search`,
    method: 'post',
    data,
  })
}
/**
 * 获取自定义情报列表
 */
export const getInfoCustomPageApi = (
  data: QueryInfoCloud
): Promise<
  ResponseData<{
    total: number
    records: InfoCustomItem[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/infoCustom/getInfoCustomPage`,
    method: 'post',
    data,
  })
}

/**
 * 获取自定义情报列表
 */
export const updateInfoCustomPageApi = (
  data: AddInfoCloud
): Promise<
  ResponseData<{
    total: number
    records: InfoCustomItem[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/infoCustom/saveOrUpdate`,
    method: 'post',
    data,
  })
}

/**
 * 获取自定义情报下拉列表
 */
export const getInfoCustomOptionsApi = (): Promise<
  ResponseData<{
    alertThreatLevel_code_cn: InfoCustomOptions
    info_ioc_type: InfoCustomOptions
    info_reliable: InfoCustomOptions
    info_threat_type: InfoCustomOptions
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/infoCustom/getAllSelectData`,
    method: 'post',
  })
}

/**
 * 获取情报白名单列表
 */
export const getInfoWhitePageApi = (
  data: QueryInfoCloud
): Promise<
  ResponseData<{
    total: number
    records: InfoWhiteItem[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/infoWhite/getPage`,
    method: 'post',
    data,
  })
}

/**
 * 更新情报白名单
 */
export const updateInfoWhiteApi = (
  data: AddInfoWhiteItem
): Promise<
  ResponseData<{
    total: number
    records: InfoWhiteItem[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/infoWhite/saveUpdate`,
    method: 'post',
    data,
  })
}

/**
 * 更新情报白名单
 */
export const deleteInfoWhiteApi = (data: { ids: number[]; deleteAll: boolean }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/infoWhite/delete`,
    method: 'delete',
    data,
  })
}
/**
 * @description '白名单导出模版'
 * @params
 */
export const exportTemplateInfoWhiteApi = (): Promise<BlobPart> => {
  return request({
    url: `/v3/ecsPlatform/infoWhite/exportTemplate`,
    method: 'get',
    responseType: 'arraybuffer',
  })
}

/**
 * @description '白名单导入'
 * @params data: { file: File }
 */
export const infoWhiteImportInfoWhiteApi = (data: { file?: File }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/infoWhite/importInfoWhite',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data,
  })
}
/**
 * @description '白名单导出'
 * @params data: { ids: number[]; searchStr:string }
 */
export const infoWhiteExportInfoWhiteApi = (data: { ids: number[]; searchStr?: string }): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/infoWhite/exportInfoWhite',
    method: 'post',
    responseType: 'arraybuffer',
    data,
  })
}

/**
 * 删除自定义情报
 */
export const deleteInfoCustomApi = (data: { ids: number[]; deleteAll: boolean }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/infoCustom/delete`,
    method: 'delete',
    data,
  })
}
/**
 * 删除自定义情报
 * @param enable  1 开启 0 关闭
 */
export const updateInfoCustomStatusApi = (data: { ids: number[]; enable: 0 | 1 }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/infoCustom/updateStatus`,
    method: 'delete',
    data,
  })
}

/**
 * 获取告警描述
 */
export const getAlarmDescriptionApi = (
  data: any
): Promise<ResponseData<{ vulnDescript: string } & { [key: string]: any }>> => {
  return request({
    url: `/v3/ecsPlatform/vuln/getVuln`,
    method: 'post',
    data,
  })
}

/**
 * 获取告警描述
 */
export const getThreatTypeListApi = (): Promise<ResponseData<{ dictLabel: string; dictValue: string }[]>> => {
  return request({
    url: `/v3/ecsPlatform/warnSearch/getThreatTypeList`,
    method: 'post',
  })
}

/**
 * 导出自定义情报
 */
export const exportInfoCustomApi = (data: { searchStr: string }): Promise<BlobPart> => {
  return request({
    url: `/v3/ecsPlatform/infoCustom/exportInfoCustom`,
    method: 'post',
    responseType: 'arraybuffer',
    data,
  })
}
/**
 * 导入自定义情报
 */
export const importInfoCustomApi = (data: {
  file: File
}): Promise<ResponseData<{ dictLabel: string; dictValue: string }[]>> => {
  return request({
    url: `/v3/ecsPlatform/infoCustom/importInfoCustom`,
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data,
  })
}

/**
 * @description '白名单导出模版'
 * @params
 */
export const infoCustomExportTempApi = (): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/infoCustom/export/temp`,
    method: 'get',
    responseType: 'arraybuffer',
  })
}

/**
 * 获取邮件分析图表数据
 */
export const getMailAnalysisChartApi = (
  data: MailAnalysisChartQuery
): Promise<
  ResponseData<{
    domainPie: { [key: string]: number }
    line: {
      inIn: {
        count: number[]
        date: string[]
      }
      inOut: {
        count: number[]
        date: string[]
      }
      outIn: {
        count: number[]
        date: string[]
      }
      outOut: {
        count: number[]
        date: string[]
      }
    }
    senderBar: { [key: string]: number }
    label: {
      [key: string]: {
        count: number
        groupName: string
        id: number
        labelName: string
      }[]
    }
    threatLevel: {
      [key: number]: number
    }
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/aiMail/mailCharts`,
    method: 'post',
    data,
  })
}

/**
 * 获取邮件分析图表数据
 */
export const getMailAnalysisListApi = (
  data: MailAnalysisChartQuery & {
    pageNum: number
    pageSize: number
    searchStr: string
    scrollId: string
  }
): Promise<
  ResponseData<{
    resList: MailAnalysisItem[]
    scrollId: string
    total: number
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/aiMail/getMailPage`,
    method: 'post',
    data,
  })
}

/**
 * 获取表格头的展示字段
 */
export const getMailFieldsApi = (): Promise<ResponseData<number[]>> => {
  return request({
    url: `/v3/ecsPlatform/mailAiField/getMailAiField`,
    method: 'post',
  })
}

/**
 * 设置邮件分析的表格头
 */
export const setMailFieldsApi = (data: number[]): Promise<ResponseData<number[]>> => {
  return request({
    url: `/v3/ecsPlatform/mailAiField/updateMailAiField`,
    method: 'post',
    data,
  })
}

/**
 * 获取内部邮箱
 */
export const getInnerMailApi = (data: {
  pageNum: number
  pageSize: number
}): Promise<
  ResponseData<{
    total: number
    records: {
      id: number
      mailAddress: string
    }[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/mailInner/getMailInnerPage`,
    method: 'post',
    data,
  })
}

/**
 * 删除内部邮箱
 */
export const deleteInnerMailApi = (data: {
  deleteAll: boolean
  ids: number[]
}): Promise<
  ResponseData<{
    total: number
    records: {
      id: number
      mailAddress: string
    }[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/mailInner/deleteMailInner`,
    method: 'delete',
    data,
  })
}
/**
 * 新增内部邮箱
 */
export const addInnerMailApi = (params: {
  address: string
}): Promise<
  ResponseData<{
    total: number
    records: {
      id: number
      mailAddress: string
    }[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/mailInner/saveMailInner`,
    method: 'post',
    params,
  })
}
/**
 * 编辑内部邮箱
 */
export const editInnerMailApi = (data: {
  id: number
  address: string
}): Promise<
  ResponseData<{
    total: number
    records: {
      id: number
      mailAddress: string
    }[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/mailInner/updateMailInner`,
    method: 'post',
    data,
  })
}

/**
 * 获取配置邮箱
 */
export const getMailConfigApi = (data: {
  pageNum: number
  pageSize: number
}): Promise<
  ResponseData<{
    total: number
    records: MailConfigItem[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/mailConfig/getMailConfigPage`,
    method: 'post',
    data,
  })
}
/**
 * 获取配置邮箱
 */
export const updateMailConfigApi = (
  data: MailConfigItem
): Promise<
  ResponseData<{
    total: number
    records: MailConfigItem[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/mailConfig/saveUpdate`,
    method: 'post',
    data,
  })
}

/**
 * 删除配置邮箱
 */
export const deleteMailConfigApi = (data: {
  deleteAll: boolean
  ids: number[]
}): Promise<
  ResponseData<{
    total: number
    records: MailConfigItem[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/mailConfig/deleteMailConfig`,
    method: 'delete',
    data,
  })
}

/**
 * 重置邮箱密码
 */
export const resetMailpasswordApi = (data: {
  id: number
  oldPassword: string
  password: string
}): Promise<
  ResponseData<{
    total: number
    records: MailConfigItem[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/mailConfig/changePassword`,
    method: 'post',
    data,
  })
}

/**
 * 获取算法配置
 */
export const getMailAlgoConfigPageApi = (): Promise<
  ResponseData<{
    total: number
    records: {
      id: number
      param: string
      type: 'default' | 'custom'
    }[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/mailAlgoConfig/getMailAlgoConfigData`,
    method: 'get',
  })
}

/**
 * 更新算法配置
 */
export const updateMailAlgoConfigPageApi = (data: {
  id?: number
  type: string
  param: string
}): Promise<
  ResponseData<{
    total: number
    records: MailConfigItem[]
  }>
> => {
  return request({
    url: `/v3/ecsPlatform/mailAlgoConfig/saveUpdate`,
    method: 'post',
    data,
  })
}

/**
 * 获取企业邮箱下拉
 */
export const getEnterpriseMailApi = (): Promise<
  ResponseData<
    {
      name: string
      used: boolean
    }[]
  >
> => {
  return request({
    url: `/v3/ecsPlatform/mailConfig/getEnterpriseMail`,
    method: 'get',
  })
}

/**
 * @description '异常登录告警: 新增规则'
 * @params
 */
export const AbnormalLandingAlarmSaveEntityApi = (
  data: AbnormalLandingAlarmSaveEntityType
): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/abnormal_landing_alarm/saveEntity`,
    method: 'post',
    data,
  })
}

/**
 * @description '异常登录告警: 新增规则'
 * @params
 */
export const AbnormalLandingAlarmUpdateEntityApi = (
  data: AbnormalLandingAlarmSaveEntityType
): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/abnormal_landing_alarm/updateEntity`,
    method: 'post',
    data,
  })
}

/**
 * @description '异常登录告警: 查询规则列表'
 * @params
 */
export const AbnormalLandingAlarmQueryAllApi = (data: { name: string }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/abnormal_landing_alarm/queryAll`,
    method: 'post',
    data,
  })
}

/**
 * @description '异常登录告警: 删除规则'
 * @params
 */
export const AbnormalLandingAlarmDeleteByIdApi = (data: { ids: number[] }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/abnormal_landing_alarm/deleteById`,
    method: 'DELETE',
    data,
  })
}

/**
 * @description '异常登录告警: 规则禁用或启用'
 * @params
 */
export const AbnormalLandingAlarmUpdateStatusApi = (data: {
  ids: number[]
  status: number
}): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/abnormal_landing_alarm/rule/updateStatus`,
    method: 'POST',
    data,
  })
}

/**
 * @description '异常登录告警规则:导入'
 * @params
 */
export const AbnormalLandingAlarmImportApi = (data: { file: File }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/abnormal_landing_alarm/rule/import',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data,
  })
}

/**
 * @description '异常登录告警规则:模版'
 * @params
 */
export const AbnormalLandingAlarmExportTemplateApi = (): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/abnormal_landing_alarm/rule/exportTemplate',
    method: 'post',
    responseType: 'arraybuffer',
  })
}

/**
 * @description '异常登录告警规则:导出规则'
 * @params
 */
export const AbnormalLandingAlarmExportApi = (data: { ids: number[] }): Promise<BlobPart> => {
  return request({
    url: '/v3/ecsPlatform/abnormal_landing_alarm/rule/export',
    method: 'post',
    data,
    responseType: 'arraybuffer',
  })
}

/**
 * @description '异常登录告警规则：地区分布'
 * @params
 */
export const AbnormalLandingAlarmQueryByRegionalApi = (
  data: AbnormalLandingAlarmQueryByRegionalType
): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/abnormal_landing_alarm/queryByRegional`,
    method: 'POST',
    data,
  })
}

/**
 * @description '异常登录告警规则：根据源IP或目的IP获取TOP10'
 * @params
 */
export const AbnormalLandingAlarmSrcIpOrAddressIpTopNApi = (
  data: AbnormalLandingAlarmSrcIpOrAddressIpTopNType
): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/abnormal_landing_alarm/srcIpOrAddressIpTopN`,
    method: 'POST',
    data,
  })
}

/**
 * @description '异常登录告警规则：时间分布'
 * @params
 */
export const AbnormalLandingAlarmDistributionOfTimeNApi = (
  data: AbnormalLandingAlarmDistributionOfTimeType
): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/abnormal_landing_alarm/distributionOfTime`,
    method: 'POST',
    data,
  })
}

/**
 * @description '异常登录告警规则：异常登录告警列表'
 * @params
 */
export const AbnormalLandingAlarmSearchApi = (data: AbnormalLandingAlarmSearchType): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/abnormal_landing_alarm/search`,
    method: 'POST',
    data,
  })
}

/**
 * @description '异常登录告警规则：根据规则ID或者威胁名称查询规则信息'
 * @params
 */
export const AbnormalLandingAlarmQueryOneApi = (data: {
  name?: string
  id?: number
}): Promise<{
  data:
    | {
        createTime: string
        country: string
        province: string
        city: string
        name: string
        serverIp: string
        startTime: string
        endTime: string
        timeStatus: string
        statusName: string
        incident: string
      }
    | string
}> => {
  return request({
    url: `/v3/ecsPlatform/abnormal_landing_alarm/queryOne`,
    method: 'POST',
    data,
  })
}

/**
 * @description 应用自定义情报
 * @params
 */
export const applyRulesApi = (): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/infoCustom/applyRules`,
    method: 'POST',
  })
}
