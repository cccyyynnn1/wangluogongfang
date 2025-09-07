import request from '@/utils/request'
import {
  algorithmsContainersType,
  algorithmsListType,
  containersLogsType,
  dockerCopyType,
  ResponseData,
  uploadContaineType,
} from '@/types/index'

/**
 * 硬件信息
 */
export const getHardInfoApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/public/hard',
    method: 'get',
  })
}

/**
 * 算法镜像列表
 */
export const getAlgorithmsListApi = (data: algorithmsListType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/algorImage/getAlgorPage',
    method: 'post',
    data,
  })
}

/**
 * 容器列表
 */
export const getAlgorithmsContainersApi = (params: algorithmsContainersType): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/docker/api/containers/json',
    method: 'get',
    params,
  })
}

/**
 * 获取容器状态
 */
export const getContainersStatsApi = (params: {
  id: string
  stream?: boolean
  oneShot?: boolean
}): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/docker/api/containers/${params.id}/stats`,
    method: 'get',
    params: { stream: params.stream, oneShot: params.oneShot },
  })
}

/**
 * 获取容器详情
 */

export const getContainersInfoApi = (params: { id: string }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/docker/api/containers/${params.id}/json`,
    method: 'get',
  })
}

/**
 * 启动容器
 */
export const startContainersApi = (params: { id: string }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/admin/container/start',
    method: 'post',
    params,
  })
}

/**
 * 关闭容器
 */
export const stopContainersApi = (params: { id: string }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/admin/container/stop',
    method: 'post',
    params,
  })
}

/**
 * 删除镜像
 */
export const deleteContainersApi = (params: { id?: number }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/algorImage/delete',
    method: 'DELETE',
    params,
  })
}

/**
 * 镜像上传
 */
export const uploadContainersApi = (
  params: uploadContaineType,
  data: {
    file: File
  }
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/algorImage/importAlgorFile',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data: { ...params, ...data },
  })
}

/**
 * 镜像保存
 */
export const saveUpdateContainersApi = (
  data: uploadContaineType,
  fileUrl: string,
  imageSize: number
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/algorImage/saveUpdate',
    method: 'post',
    data: { ...data, fileUrl },
  })
}

/**
 * 创建容器
 */
export const createContainersApi = (params: { id: number; runParam: string }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/admin/container/run',
    method: 'post',
    params,
  })
}

/**
 * 容器进程
 */

export const getContainersProgressApi = (params: { id: string }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/docker/api/containers/${params.id}/top`,
    method: 'get',
  })
}

/**
 * 容器日志
 */

export const getContainersLogsApi = (params: containersLogsType): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/docker/api/containers/${params.id}/logs`,
    method: 'get',
    params: { ...params, id: null },
  })
}

/**
 * 上传文件
 */

export const dockerCopyApi = (
  params: dockerCopyType,
  data: {
    filedata: File | undefined
  }
): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/admin/docker/copy`,
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    params,
    data,
  })
}

/**
 * 通过id启动容器
 */

export const startContainersByIdApi = (params: { id: string }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/docker/api/containers/${params.id}/start`,
    method: 'post',
  })
}

/**
 * 通过id停止容器
 */

export const stopContainersByIdApi = (params: { id: string }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/docker/api/containers/${params.id}/stop`,
    method: 'post',
  })
}

/**
 * 通过id重启容器
 */

export const restartContainersByIdApi = (params: { id: string }): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/docker/api/containers/${params.id}/restart`,
    method: 'post',
  })
}

/**
 * 节点和链路数据来源接口
 */
export const flowDeviceCountApi = (): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/public/forward`,
    method: 'POST',
    data: {
      reqType: 'GET',
      url: `/flowDevice/flowDeviceCount`,
    },
  })
}
