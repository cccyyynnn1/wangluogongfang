import type { InjectionKey, Ref } from 'vue'

interface AssetsIp_DetailKey {
  detailVisible: Ref<boolean>
  isEdit: Ref<boolean>
  currentRow: Ref<any>
}
export type AssetsType = 'assetVisits' | 'known' | 'unknow' | 'site' | 'nosite' | 'asset-overview'
export type IpDetailTabsValue = 'overview' | 'serviceVisits' | 'applicationData'
// export type IpDetailDataType = 'http' | 'dns' | 'icmp' | 'custom'
export type IpDetailDataType = 'HTTP' | 'DNS' | 'FTP' | 'SMB' | 'DB' | 'MAIL'

export const AssetsIp_Detail: InjectionKey<AssetsIp_DetailKey> = Symbol('AssetsIp_Detail')
