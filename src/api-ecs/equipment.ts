import request from '@/utils/request'
import { ResponseData, EquipmentListType, UpdateEquipmentType, FlowProbeListType } from '@/types/index'

/**
 * 设备列表接口
 */
export const equipmentListApi = (data: EquipmentListType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/public/forward',
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/flowDevice/page?page=${data.page}&limit=${data.limit}&query=${data.query}`,
    },
  })
}

/**
 * 更新设备接口
 */
export const updateEquipmentApi = (
  flowDeviceId: string,
  data: UpdateEquipmentType,
  params: { password: string }
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/public/forward',
    method: 'POST',
    data: {
      reqType: 'PUT',
      url: `/flowDevice/edit/${flowDeviceId}`,
      paramMap: data,
    },
    params,
  })
}

/**
 * 重启服务接口
 */
export const restartEquipmentApi = (flowDeviceId: string, params: { password: string }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'PUT',
      url: `/flowDevice/restart/${flowDeviceId}`,
    },
    params,
  })
}

/**
 * 链路列表接口
 */
export const flowProbeListApi = (params: FlowProbeListType): Promise<ResponseData<any>> => {
  // return request({
  //   url: '/v3/ecsPlatform/public/forward',
  //   method: 'get',
  //   params,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/flowProbe/page?page=${params.page}&limit=${params.limit}&query=${params.query}`,
    },
  })
}

/**
 * 获取⽹卡ID
 */
export const getAdapterByIdApi = (params: any): Promise<ResponseData<any>> => {
  // return request({
  //   url: `/v3/ecsPlatform/public/forward/adapter/page`,
  //   method: 'get',
  //   params,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/adapter/page?page=${params.page}&limit=${params.limit}&query=${params.query}&url=${params.url}`,
    },
  })
}

/**
 * 链路添加接口
 */
export const flowProbeAddApi = (data: any, params: { password: string }): Promise<ResponseData<any>> => {
  // return request({
  //   url: '/v3/ecsPlatform/public/forward/flowProbe/add',
  //   method: 'post',
  //   data,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'POST',
      url: `/flowProbe/add`,
      paramMap: data,
    },
    params,
  })
}

/**
 * 链路修改接口
 */
export const flowProbeUpdataApi = (id: string, data: any, params: { password: string }): Promise<ResponseData<any>> => {
  // return request({
  //   url: `/v3/ecsPlatform/public/forward/flowProbe/edit/${id}`,
  //   method: 'put',
  //   data,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'PUT',
      url: `/flowProbe/edit/${id}`,
      paramMap: data,
    },
    params,
  })
}

/**
 * 链路删除接口
 */
export const flowProbeDeleteApi = (data: string, params: { password: string }): Promise<ResponseData<any>> => {
  // return request({
  //   url: `/v3/ecsPlatform/public/forward/flowProbe/delete/${data}`,
  //   method: 'DELETE',
  //   data,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'DELETE',
      url: `/flowProbe/delete/${data}`,
    },
    params,
  })
}

/**
 * 设备⽹卡信息
 */
export const networkListApi = (params: FlowProbeListType): Promise<ResponseData<any>> => {
  // return request({
  //   url: '/v3/ecsPlatform/public/forward/adapter/page',
  //   method: 'get',
  //   params,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/adapter/page?page=${params.page}&limit=${params.limit}&query=${params.query}`,
    },
  })
}

/**
 * 启⽤⽹卡抓包
 */
export const adapterStartApi = (id: string, params: { password: string }): Promise<ResponseData<any>> => {
  // return request({
  //   url: `/v3/ecsPlatform/public/forward/adapter/start/${id}`,
  //   method: 'put',
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'PUT',
      url: `/adapter/start/${id}`,
    },
    params,
  })
}

/**
 * 停⽤⽹卡抓包
 */
export const adapterStopApi = (id: string, params: { password: string }): Promise<ResponseData<any>> => {
  // return request({
  //   url: `/v3/ecsPlatform/public/forward/adapter/stop/${id}`,
  //   method: 'put',
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'PUT',
      url: `/adapter/stop/${id}`,
    },
    params,
  })
}

/**
 * 设备列表接口
 */
export const kafkaPushListApi = (params: EquipmentListType): Promise<ResponseData<any>> => {
  // return request({
  //   url: '/v3/ecsPlatform/public/forward/flowKafkaPush/page',
  //   method: 'get',
  //   params,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/flowKafkaPush/page?page=${params.page}&limit=${params.limit}&query=${params.query}`,
    },
  })
}

/**
 * kafka服务列表
 */
export const kafkaServerConfigApi = (params: any): Promise<ResponseData<any>> => {
  // return request({
  //   url: '/v3/ecsPlatform/public/forward/kafkaServerConfig/page',
  //   method: 'get',
  //   params,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/kafkaServerConfig/page?page=${params.page}&limit=${params.limit}&query=${params.query}`,
    },
  })
}

/**
 * kafka推送添加接口
 */
export const flowKafkaPushAddApi = (data: any, params: { password: string }): Promise<ResponseData<any>> => {
  // return request({
  //   url: '/v3/ecsPlatform/public/forward/flowKafkaPush/add',
  //   method: 'post',
  //   data,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'POST',
      url: `/flowKafkaPush/add`,
      paramMap: data,
    },
    params,
  })
}

/**
 * kafka推送修改接口
 */
export const flowKafkaPushUpdataApi = (
  id: string,
  data: any,
  params: { password: string }
): Promise<ResponseData<any>> => {
  // return request({
  //   url: `/v3/ecsPlatform/public/forward/flowKafkaPush/edit/${id}`,
  //   method: 'put',
  //   data,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'PUT',
      url: `/flowKafkaPush/edit/${id}`,
      paramMap: data,
    },
    params,
  })
}

/**
 * kafka推送删除接口
 */
export const flowKafkaPushDeleteApi = (data: string, params: { password: string }): Promise<ResponseData<any>> => {
  // return request({
  //   url: `/v3/ecsPlatform/public/forward/flowKafkaPush/delete/${data}`,
  //   method: 'DELETE',
  //   data,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'DELETE',
      url: `/flowKafkaPush/delete/${data}`,
    },
    params,
  })
}

/**
 * 系统服务列表接口
 */
export const sysServiceListApi = (params: FlowProbeListType): Promise<ResponseData<any>> => {
  // return request({
  //   url: '/v3/ecsPlatform/public/forward/sysService/page',
  //   method: 'get',
  //   params,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/sysService/page?page=${params.page}&limit=${params.limit}&query=${params.query}`,
    },
  })
}

/**
 * 启⽤服务
 */
export const sysServiceStartApi = (id: string, params: { password: string }): Promise<{ msg: string }> => {
  // return request({
  //   url: `/v3/ecsPlatform/public/forward/sysService/start/${id}`,
  //   method: 'put',
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'PUT',
      url: `/sysService/start/${id}`,
    },
    params,
  })
}

/**
 * 停⽤服务
 */
export const sysServiceStopApi = (id: string, params: { password: string }): Promise<{ msg: string }> => {
  // return request({
  //   url: `/v3/ecsPlatform/public/forward/sysService/stop/${id}`,
  //   method: 'put',
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'PUT',
      url: `/sysService/stop/${id}`,
    },
    params,
  })
}

/**
 * 负载均衡列表
 */
export const loadBalanceApi = (params: FlowProbeListType): Promise<ResponseData<any>> => {
  // return request({
  //   url: '/v3/ecsPlatform/public/forward/loadBalance/page?',
  //   method: 'get',
  //   params,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/loadBalance/page?page=${params.page}&limit=${params.limit}&query=${params.query}`,
    },
  })
}

/**
 * 负载均衡VS列表
 */
export const loadBalanceVsApi = (params: FlowProbeListType): Promise<ResponseData<any>> => {
  // return request({
  //   url: '/v3/ecsPlatform/public/forward/loadBalanceVs/page',
  //   method: 'get',
  //   params,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/loadBalanceVs/page?page=${params.page}&limit=${params.limit}&query=${params.query}`,
    },
  })
}

/**
 * 负载均衡POOL列表
 */
export const loadBalancePoolApi = (params: FlowProbeListType): Promise<ResponseData<any>> => {
  // return request({
  //   url: '/v3/ecsPlatform/public/forward/loadBalancePool/page',
  //   method: 'get',
  //   params,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/loadBalancePool/page?page=${params.page}&limit=${params.limit}&query=${params.query}`,
    },
  })
}

/**
 * 负载均衡添加接口
 */
export const loadBalanceAddApi = (data: any, params: { password: string }): Promise<ResponseData<any>> => {
  // return request({
  //   url: '/v3/ecsPlatform/public/forward/loadBalance/add',
  //   method: 'post',
  //   headers: {
  //     'Content-Type': 'multipart/form-data',
  //   },
  //   data,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data: {
      reqType: 'POST',
      url: `/loadBalance/add`,
      paramMap: data,
    },
    params,
  })
}

/**
 * 负载均删除接口
 */
export const loadBalanceDeleteApi = (data: string, params: { password: string }): Promise<ResponseData<any>> => {
  // return request({
  //   url: `/v3/ecsPlatform/public/forward/loadBalance/delete/${data}`,
  //   method: 'DELETE',
  //   data,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'DELETE',
      url: `/loadBalance/delete/${data}`,
    },
    params,
  })
}

/**
 * 负载均衡型号列表接口
 */
export const loadBalanceModelListApi = (params: FlowProbeListType): Promise<ResponseData<any>> => {
  // return request({
  //   url: '/v3/ecsPlatform/public/forward/loadBalanceModel/page',
  //   method: 'get',
  //   params,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/loadBalanceModel/page?page=${params.page}&limit=${params.limit}&query=${params.query}`,
    },
  })
}

/**
 * 流量设备更新列表接口
 */
export const flowDeviceUpdateListApi = (params: FlowProbeListType): Promise<ResponseData<any>> => {
  // return request({
  //   url: '/v3/ecsPlatform/public/forward/flowDeviceUpdate/page',
  //   method: 'get',
  //   params,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/flowDeviceUpdate/page?page=${params.page}&limit=${params.limit}&query=${params.query}`,
    },
  })
}

/**
 * 流量设备更新删除接口
 */
export const flowDeviceUpdateDeleteApi = (data: string): Promise<ResponseData<any>> => {
  // return request({
  //   url: `/v3/ecsPlatform/public/forward/flowDeviceUpdate/delete/${data}`,
  //   method: 'DELETE',
  //   data,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'DELETE',
      url: `/flowDeviceUpdate/delete/${data}`,
    },
  })
}

/**
 * 流量分析设备列表接口
 */
export const flowDeviceListApi = (params: {
  page: number
  limit: number
  query: string
}): Promise<ResponseData<any>> => {
  // return request({
  //   url: '/v3/ecsPlatform/public/forward/flowDevice/page',
  //   method: 'get',
  //   params,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/flowDevice/page?page=${params.page}&limit=${params.limit}&query=${params.query}`,
    },
  })
}

/**
 * 流量分析设备更新接口
 */
export const flowDeviceUpdateApi = (
  data: { id: string; flowDeviceIds: string },
  params: { password: string }
): Promise<{ message: string }> => {
  // return request({
  //   url: `/v3/ecsPlatform/public/forward/flowDeviceUpdate/update/${params.id}`,
  //   method: 'PUT',
  //   params: { flowDeviceIds: params.flowDeviceIds },
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'PUT',
      url: `/flowDeviceUpdate/update/${data.id}`,
      paramMap: {
        flowDeviceIds: data.flowDeviceIds,
      },
    },
    params: params,
  })
}

/**
 * 负载均衡概览列表接口
 */
export const overviewLoadBalanceListApi = (params: FlowProbeListType): Promise<ResponseData<any>> => {
  // return request({
  //   url: '/v3/ecsPlatform/public/forward/loadBalance/overviewloadbalance',
  //   method: 'get',
  //   params,
  // })
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/loadBalance/overviewloadbalance?page=${params.page}&limit=${params.limit}&query=${params.query}`,
    },
  })
}

/**
 * 基础配置更新接口
 */
export const settingVarUpdateVarApi = (
  data: { id: number; value: string },
  params: { password: string }
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/settingVar/updateVar',
    method: 'post',
    data,
    params,
  })
  // return request({
  //   url: `/v3/ecsPlatform/public/forward`,
  //   method: 'POST',
  //   data: {
  //     reqType: 'PUT',
  //     url: `/flowDeviceUpdate/update/${params.id}`,
  //     paramMap: {
  //       flowDeviceIds: params.flowDeviceIds,
  //     },
  //   },
  // })
}
