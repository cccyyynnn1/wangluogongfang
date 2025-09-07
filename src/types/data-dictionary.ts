export interface DictTypePageParame {
  pageNum: number
  pageSize: number
  dictName: string
  dictType: string
  /**正常  1； 停用 0 */
  enable: number | string
}
export interface DictTypePageItem {
  createTime?: number
  dictName: string
  dictType: string
  /**正常  1； 停用 0 */
  enable: number
  remark?: string
  id?: number
  isDefault?: number
}

export interface AssetLabelDict {
  createTime?: number
  transferValue: string
  sourceValue?: string
  remark?: string
  id?: number
}

export interface DictTypePageListItem {
  createTime?: number
  dictLabel: string
  dictValue: string
  dictTypeId?: string | number
  /**正常  1； 停用 0 */
  enable: number
  remark?: string
  id?: number
  isDefault?: number
}
