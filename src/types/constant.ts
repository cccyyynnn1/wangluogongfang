import { InjectionKey } from 'vue'
/**
 * 接口返回类型
 */
export type ResponseData<T> = {
  code: number
  data: T
  msg: string
}
/**
 * indexType类型
 */
export interface IndexTypeTpye {
  value: number
  label: string
  en: string
}
export const tableSearch = Symbol('tableSearch') as InjectionKey<() => void>
