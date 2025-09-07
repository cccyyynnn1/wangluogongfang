export interface DowmloadType {
  /**
   * 下载sql语句，日志：检索条件。数据包：下载地址
   */
  downloadCnd?: string
  /**
   * 检索结束时间，毫秒
   */
  edTime?: number
  /**
   * 索引类型，日志下载时  必填
   */
  indexType: number
  /**
   * 检索开始时间，毫秒
   */
  stTime?: number
  /**
   * 类型，1: 日志 0：数据包
   */
  type?: number
  id?: number
  status?: number
  dataType: number
}
