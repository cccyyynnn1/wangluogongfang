/**
 * @description 数据包回放查询参数
 */
export type PacketReplayQueryParams = {
  page: number
  limit: number
  query: {
    /**
     * 数据链路ID
     */
    flowProbeIdStr: string
    flowProbeId: string[]
    /**
     * 回放方式
     */
    type?: 'fast' | 'original'
    /**
     * 文件名
     */
    fileName: string
    /**
     * 指纹
     */
    fingerPoint: string
    /**
     * 创建时间
     */
    crtTimeArr: string[]
    crtTime: string
    /**
     * 最后回放时间
     */
    lastRunDateTimeArr: string[]
    lastRunDateTime: string
    /**
     * 最后回放结果
     */
    lastRunResult: string
    /**
     *  备注
     */
    note: string
  }
}
export type PacketReplayQuery = {
  reqType: 'GET' | 'POST' | 'DELETE' | 'PUT'
  url: string
  paramMap?: object
}

export type PacketReplayItem = {
  crtTime: string
  crtTimeStr: string
  crtUser: string
  crtUserStr: string
  fileName: string
  fileSize: 1169751
  fileSizeStr: string
  fingerPoint: string
  flowProbeId: string
  flowProbeIdStr: string
  id: string
  lastRunDateTime: string
  lastRunDateTimeStr: string
  lastRunResult: string
  note: string
  type: 'fast' | 'original'
  typeStr: string
}

export type RemoteFolderItem = {
  fileName: string
  filePath: string
  fileSize: string
  remark: string
  type: string
}
