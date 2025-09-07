/**
 * customField更新保存接口类型类型
 */
export interface customFieldSaveOrUpdateType {
  fieldNameCn: string
  /**
   * 编辑时不为空
   */
  id?: string | number | null
  /**
   * 1：默认显示  0：不显示
   */
  isDisplay?: string | number | boolean
  /**
   * 传1，目前只存http类型的
   */
  type: number
}

export interface customFieldGetGageType {
  fieldNameCn: string
  indexType?: string | number
  pageNum: number
  pageSize: number
}

export type InitialisationItem = {
  content: string
  createTime: number
  id: number
  isDelete: 0 | 1
  remark: string
  status: 0 | 1
  templateName: string
  updateTime: number
  userId: number
}
export type InitialisationList = InitialisationItem[]

export type TemplateFeildItem = {
  namesList: string[]
  typeName: string
  type: string
  id?: number
  remark?: string
  apiUrl?: string
}
type TemplateSessionInfo = TemplateFeildItem & { apiList: TemplateFeildItem[] }
export type TemplateInfo = {
  api_default: TemplateFeildItem
  session_default: TemplateFeildItem
  session: TemplateSessionInfo[]
  alarm: TemplateFeildItem[]
  survey: TemplateFeildItem[]
}
