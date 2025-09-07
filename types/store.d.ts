import { AnyARecord } from 'node:dns'
import { VabRouteMeta, VabRouteRecord } from '/#/router'

declare interface AclModuleType {
  admin: boolean
  role: string[]
  permission: string[]
}
declare interface AiAnalyseModuleType {
  aiAnalyseHistory: Map<
    string,
    {
      history: { assistant: string; user: string; isAnalyse?: boolean }[]
      updateTime: number
    }
  >
}
declare interface ErrorLogModuleType {
  errorLogs: any[]
}

declare interface RoutesModuleType {
  tab: {
    data: string | undefined
  }
  tabMenu: string | undefined
  activeMenu: {
    data: string | undefined
  }
  routes: VabRouteRecord[]
}

declare type DeviceType = 'mobile' | 'desktop'
declare type LanguageType = 'zh' | 'en'
declare type ToolType =
  | 'all'
  | 'encoding-decoding'
  | 'unixtime'
  | 'hexconvert'
  | 'reset-password'
  | 'subnetmask'
  | 'regexp'
declare interface SettingsModuleType {
  theme: ThemeType
  device: DeviceType
  collapse: boolean
  language: LanguageType
  lock: boolean
  logo: string
  title: string
  echartsGraphic1: string[]
  echartsGraphic2: string[]
  SecConfig: any
  systemConfig: {
    value: {
      name: string
      nameAbb: string
      icp: string
      copyright: string
      logo: string
      loginBg: string
      logginTip: string
      loginLogo: string
    }
    id: number
    key: string
    uid: number
    name: string
  }
  toolboxVisible: boolean
  downloadVisible: boolean
  messageVisible: boolean
  toolType: ToolType
  dialogData: any
  filePath: string
  msgId: number | undefined
}

declare interface TabsModuleType {
  visitedRoutes: VabRouteRecord[]
}

declare interface OptionType {
  name?: string
  title?: string
  meta: VabRouteMeta
}
declare interface TableColumnType {
  [key: number]: TableColumnItemType[]
}

declare type TableFieldType = 'no_serach' | 'ip' | 'num' | 'sum' | 'average' | 'direction' | 'text' | 'isn_t'

declare interface TableColumnItemType {
  createTime: string | number | null
  fieldNameCn: string
  fieldNameEn: string
  type: string
  updateTime: string | number | null
  id: number
  isDeleted: 0 | 1 | 2
  isDisplay: 0 | 1
  supportAgg: boolean
  fieldType: TableFieldType
  supportAgg: boolean
}
declare interface UserModuleType {
  token: string | boolean
  expire: boolean
  username: string
  avatar: string
  tableColumns: TableColumnType | null
  userId: number
  indexTypeList?: []
  indexTypeINcouldsNetList?: []
  systemMenus: []
  highLightLastUpdateTime: number
  indexFieldsLastUpdateTime: number
  userDisPlaysFiled: { [key: number]: number[] }
  attackCharacterization: {
    highLightWhite: {
      content: string
      scope: 'requestHeader' | 'requestPayload' | 'responseHeader' | 'responsePayload' | 'all'
    }[]
    highLightConfig: {
      content: string
      scope: 'requestHeader' | 'requestPayload' | 'responseHeader' | 'responsePayload' | 'all'
    }[]
  }
}

export type GlobalMessageItem = {
  content: string
  id: number
  createTime: number
  minorType: string
  msgType: 'user' | 'sys'
  msgSource: string
  notes: string
  readStatus: boolean
  important: boolean
  title: string
  sendUserId?: number
  isFold?: boolean
  isHide?: boolean
}
export type GlobalLoginMessageType = {
  newMsgs: GlobalMessageItem[]
  unReadCount: number
}
declare interface PublicType {
  flowProbesList: any[]
  flowProbesIDList: any[]
  allDisPlaysFiledList: any[]
  globalMessageList: GlobalMessageItem[]
  globalLoginMessageList: GlobalLoginMessageType
  // globalSysMessageList: GlobalMessageItem[]
  // globalUserMessageList: GlobalMessageItem[]
  // globalSysIndex: number
  // globalUserIndex: number
  hasNewMessage: boolean
}
