import request from '@/utils/request'

import { ResponseData, DowmloadType } from '@/types'

/**
 * 日志导出
 */
export const DownloadLogApi = (data: DowmloadType, params: number): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/downloadRecord/saveOrUpdate?dataType=${params}`,
    method: 'post',
    data,
  })
}

/**
 * 异步下载pcap
 */
export const DownloadPcapApi = (data: { uid: number; query: any }, params: number): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/downloadRecord/savePcap?dataType=${params}`,
    method: 'post',
    data,
  })
}

/**
 * 日志下载分页
 */
export const DownloadLogPageApi = (data: {
  pageNum: number
  pageSize: number
  dataType: number | null
  searchStr: string
}): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/downloadRecord/getDownLoadRecordPage',
    method: 'post',
    data,
  })
}

/**
 * 日志下载删除
 */
export const DownloadLogDelApi = (data: { ids?: Array<string>; deleteAll: boolean }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/downloadRecord/delete',
    method: 'DELETE',
    data,
  })
}
