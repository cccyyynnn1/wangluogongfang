import { Component } from 'vue'
import { VIEWOBJ } from '../data/constant'

/**
 * 数据趋势每一项数据类型
 */
export type TrendseriesDataItemType = {
  name: string
  type: string
  data: number[]
  symbol: string
  smooth: boolean
  yAxisIndex: number
  showSymbol: boolean
  areaStyle: {
    opacity: number
    color: string
  }
}

/**
 * 监控的数据类型
 */
export interface MonitoringType {
  uuid: string
  name: string
  timeQuantum: string
  switchValue: boolean
  timeDuration: string
  hasError?: boolean
  id?: number
  timeDate?: string
  list: MonitoringItem[] | string
  checked?: boolean
}

export interface MonitoringItem {
  uuid: string
  // title: string
  style: UIStyle
  height?: number
  name: Component
  hasError?: boolean
  key: keyof typeof VIEWOBJ
  data: any
}

interface UIStyle {
  height?: string
  left?: number | string
  top?: string
  width?: number | string
  x?: number
  y?: number
  long?: number
}
