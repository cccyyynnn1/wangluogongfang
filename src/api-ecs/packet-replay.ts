import request from '@/utils/request'
import { PacketReplayQuery, ResponseData, PacketReplayItem, RemoteFolderItem } from '@/types/index'
/**
 * 数据包回放列表
 */
export const getPacketReplayListApi = (
  data: PacketReplayQuery
): Promise<
  ResponseData<{
    total: number
    list: PacketReplayItem[]
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/public/forward',
    method: 'post',
    data,
  })
}
/**
 * 数据包回放-数据链路下拉
 */
export const getFlowProbeApi = (
  data: PacketReplayQuery
): Promise<
  ResponseData<{
    list: {
      adapterId: string[]
      adapterIdStr: string
      name: string
      id: string
    }[]
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/public/forward',
    method: 'post',
    data,
  })
}
/**
 * 数据包回放-删除数据包
 */
export const deleteFlowProbeApi = (
  data: PacketReplayQuery
): Promise<
  ResponseData<{
    message: string
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/public/forward',
    method: 'post',
    data,
  })
}

/**
 * 数据包重新回放
 */
export const playBackFlowProbeApi = (
  data: PacketReplayQuery
): Promise<
  ResponseData<{
    message: string
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/public/forward',
    method: 'post',
    data,
  })
}
/**
 * 数据包本地Pcap分页
 */
export const getLocalPcapsApi = (data: {
  pageNum: number
  pageSize: number
  searchStr: string
}): Promise<
  ResponseData<{
    records: any[]
    total: number
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/pcapUploadRecord/getPcapRecordPage',
    method: 'post',
    data,
  })
}
/**
 * 数据包本地Pcap上传
 */
export const uploadLocalPcapsApi = (data: {
  files: File[]
  remarks: string[]
}): Promise<
  ResponseData<{
    msg: string
  }>
> => {
  const formData = new FormData()
  data.files.forEach((file) => formData.append('files', file))
  formData.append('remarks', JSON.stringify(data.remarks))
  return request({
    url: '/v3/ecsPlatform/pcapUploadRecord/uploadLocalPcap',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data: formData,
  })
}

/**
 * 数据包本地Pcap上传
 */
export const addLocalPcapsTaskApi = (data: {
  type: string
  typeStr: string
  flowProbeId: string
  flowProbeIdStr: string
  note: string
  noteStr: string
  fileNameList: {
    originName: string
    uniName: string
    remark: string
  }[]
}): Promise<
  ResponseData<{
    msg: string
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/public/packetPlaybackSave',
    method: 'post',
    data: {
      ...data,
      dialogTitle: 'add',
      moduleType: 'packetPlayback',
    },
  })
}
/**
 * 数据包本地Pcap上传
 */
export const editLocalPcapsStatusApi = (data: {
  id: number
  remark: string
}): Promise<
  ResponseData<{
    msg: string
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/pcapUploadRecord/updatePcapRecord',
    method: 'post',
    data,
  })
}
/**
 * 数据包本地Pcap上传
 */
export const deleteLocalPcapsApi = (data: {
  ids: number[]
}): Promise<
  ResponseData<{
    msg: string
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/pcapUploadRecord/deletePcapRecord',
    method: 'delete',
    data,
  })
}

/**
 * 获取远程目录文件
 */
export const getRemoteFolderApi = (data: {
  readType: string
  ip?: string
  port?: number
  path?: string
  account: string
  password: string
}): Promise<ResponseData<RemoteFolderItem[]>> => {
  return request({
    url: '/v3/ecsPlatform/public/lookRemoteFolder',
    method: 'post',
    data,
  })
}

/**
 * 数据包本地Pcap上传
 */
export const saveRemoteFolderApi = (data: {
  readType: string
  ip: string
  port: number
  path: string
  account: string
  password: string
  requestId: number
  remoteFileList: {
    fileName: string
    filePath: string
    remark: string
  }[]
}): Promise<ResponseData<RemoteFolderItem[]>> => {
  return request({
    url: '/v3/ecsPlatform/public/saveRemotePcapToLocal',
    method: 'post',
    data,
  })
}
