import request from '@/utils/request'

import { ResponseData } from '../types'

/**
 * @description '知识：查询'
 * @params
 */
export const knowledgeGetAllApi = (): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/knowledge/getAll',
    method: 'post',
  })
}

/**
 * @description '知识：上传文件'
 * @params data: { file: File }
 */
export const knowledgeUploadApi = (data: { file?: File }): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/knowledge/upload',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data,
  })
}
/**
 * @description '知识：删除'
 * @params data: { ids: number[]}
 */
export const knowledgeDeleteApi = (data: string): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/knowledge/delete',
    method: 'DELETE',
    data,
  })
}
/**
 * @description '临时知识删除'
 */
export const knowledgeFileDeleteApi = (path: string): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/knowledge/deleteByFilePath`,
    method: 'DELETE',
    data: {
      path,
    },
  })
}
/**
 * @description 保存知识库
 */
export const knowledgeFileSaveApi = (paths: string[]): Promise<ResponseData<any>> => {
  return request({
    url: `/v3/ecsPlatform/knowledge/save`,
    method: 'post',
    data: {
      paths,
    },
  })
}
