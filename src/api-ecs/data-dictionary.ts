import request from '@/utils/request'
import { ResponseData, DictTypePageParame, DictTypePageItem, DictTypePageListItem } from '@/types'

/**
 * 获取字典类型列表
 */
export const getDictTypePageApi = (
  data: DictTypePageParame
): Promise<
  ResponseData<{
    records: DictTypePageItem[]
    total: number
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/sysDictType/getDictTypePage',
    method: 'post',
    data,
  })
}

/**
 * 编辑字典类型
 */
export const editDictTypePageApi = (
  data: DictTypePageItem,
  params: { password: string }
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/sysDictType/saveOrUpdate',
    method: 'post',
    data,
    params,
  })
}

/**
 * 删除字典类型
 */
export const deleteDictTypePageApi = (
  data: { ids: number[]; deleteAll: boolean },
  params: { password: string }
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/sysDictType/delete',
    method: 'delete',
    data,
    params,
  })
}

/**
 * 获取字典标签列表
 */
export const getDictTypePageListApi = (data: {
  pageNum: number
  pageSize: number
  dictTypeId: string | number | undefined
}): Promise<
  ResponseData<{
    records: DictTypePageListItem[]
    total: number
  }>
> => {
  return request({
    url: '/v3/ecsPlatform/sysDictData/getDictDataPage',
    method: 'post',
    data,
  })
}
/**
 * 编辑字典标签
 */
export const editDictTypePageListApi = (
  data: DictTypePageListItem,
  params: { password: string }
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/sysDictData/saveOrUpdate',
    method: 'post',
    data,
    params,
  })
}
/**
 * 删除字典标签
 */
export const deleteDictTypePageListApi = (
  data: {
    dictTypeId: string | number | undefined
    ids: number[]
    deleteAll: boolean
  },
  params: { password: string }
): Promise<ResponseData<any>> => {
  return request({
    url: '/v3/ecsPlatform/sysDictData/delete',
    method: 'delete',
    data,
    params,
  })
}
