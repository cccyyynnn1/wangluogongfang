/**
 * 获取算法镜像接口类型
 */
export interface algorithmsListType {
  pageNum: number // 页码
  pageSize: number // 每页数量
  algorName?: string // 算法名称
  imageName?: string // 算法镜像
}

/**
 * 获取容器列表接口类型
 */
export interface algorithmsContainersType {
  all?: boolean // bool值，默认false。默认仅返回运行中的容器
  limit?: number // 返回容器数量
  size?: boolean // bool值，默认false。返回容器大小字段，SizeRw 和 SizeRootFs
  filters?: string | any // 过滤器。json类型{"name":"镜像名称","status":"created|restarting|running|removing|paused|exited|dead"}
}

/**
 * 镜像上传类型
 */
export interface uploadContaineType {
  algorName: string // 算法名称，唯一
  imageName?: string // 容器名称，唯一
  description?: string // 描述
}

/**
 * 容器日志
 */

export interface containersLogsType {
  id: string
  follow?: boolean
  stdout?: boolean
  stderr?: boolean
  since?: number
  until?: number
  timestamps?: boolean
  tail?: string
}
/**
 * 上传文件类型
 */

export interface dockerCopyType {
  container: string
  permission?: number
  target: string
}
