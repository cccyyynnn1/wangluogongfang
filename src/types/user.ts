export interface QueryUserType {
  query: {
    pageNum: number
    pageSize: number
  }
  user?: {
    nickName: string
  }
}
export interface QuerySysDeptType {
  pageNum: number
  pageSize: number
  deptName: string
}
export interface SysDeptItem {
  ancestors: string
  createTime: number | null
  createUser: number | null
  deptName: string
  id: number
  parentId: number
  updateTime: number | null
  updateUser: number | null
  hasChildren?: boolean
  children?: SysDeptItem[]
}

export interface SysDeptData {
  countId: number
  current: number
  maxLimit: string
  optimizeCountSql: boolean
  orders: string[]
  pages: number
  records: SysDeptItem[]
  size: number
  total: number
}

export interface EditSysDeptData {
  // 为空是新增 编辑为空
  id?: number
  parentId?: number
  deptName: string
}

export interface QueryUserParams {
  pageNum: number
  pageSize: number
  loginName?: string
  userType?: number
  deptIds?: number[]
  roleIds?: number[]
  createSt?: number
  createEd?: number
}
export interface Role {
  id: number
  createTime: number
  deptName: number
  deptId: number
  loginName: string
  mail?: string
  nickName?: string
  deptIds?: string
  phone?: string
  roleName?: number
  roleType?: number
  password: string
  status?: number
}
export interface AddUserParam {
  nickName: string
  loginName: string
  password: string
  deptId: number
  userType: number
  effectiveSt: string
  effectiveEd: string
  roleIds: number[]
  description?: string
  phone?: string
  mail?: string
  ipWhiteList?: string
  address?: string
  visibleEncryptData?: boolean
  twoFactorAuth?: boolean
  visibleBuiltinData?: boolean
  limitIp: string
  confirm?: boolean
}
export interface SysRole {
  createTime: number
  createUser: number
  createUserName: string
  description: number
  enable: number
  id?: number
  roleName: string
  roleType: number
  updateTime: number
  updateUser: number
  updateUserName: string
  menuIds: number[]
  userNames: string[]
}

export interface EditSysRole {
  enable: number
  id?: number
  menuIds: number[]
  roleName: string
  loginMenuIdList?: number[]
}
