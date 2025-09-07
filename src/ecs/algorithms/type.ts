import type { InjectionKey, Ref } from 'vue'
interface AlgorithmsAlertsKey {
  isAdvanced: Ref<boolean>
  autoRefresh: Ref<boolean>
  autoRefreshDuration: Ref<number>
}

export type AlgorithmsType = 'examples' | 'alert' | 'image'

export type ExamplesDetailType = 'overview' | 'process' | 'logs' | 'upload'

export type ExamplesCheckType = 'port' | 'storage' | 'links' | 'internet'

export type OperationType = 'start' | 'stop' | 'restart'

export const Algorithms_alerts: InjectionKey<AlgorithmsAlertsKey> = Symbol('AlgorithmsAlerts')
