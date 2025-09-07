/**
 * @description 登录、获取用户信息、退出登录、清除token逻辑，不建议修改
 */
import { useAclStore } from './acl'
import { useTabsStore } from './tabs'
import { useRoutesStore } from './routes'
import { usePubilcStore } from './public'
import { useSettingsStore } from './settings'
import { UserModuleType } from '/#/store'
import { socialLogin } from '@/api/user'
import { getAllFieldApi, getAllIndexTypeApi, getAllIndexTypeIncouldsNetApi } from '@/api-ecs/public'
import { getToken, removeToken, setToken } from '@/utils/token'
import router, { resetRouter } from '@/router'
import { isArray, isString } from '@/utils/validate'
import { tokenName } from '@/config'
import { gp } from '@gp'
import { loginApi, getUserInfoApi, logoutApi, loginOtherTypeApi } from '~/src/api-ecs/login'
import { LoginPayload, RetrieveIndexType, LoginOtherType } from '@/types/index'
import { updateSystemMenusApi, getSystemMenusApi } from '~/src/api-ecs/system'
import { getAllHightLightAPI, getAllDisPlaysFiledApi } from '~/src/api-ecs/public'

import { asyncRoutes } from '@/router'
export const useUserStore = defineStore('user', {
  state: (): UserModuleType => ({
    token: getToken() as string,
    username: '游客',
    userId: 0,
    avatar: '/download/upload/avatar.gif',
    tableColumns: null,
    systemMenus: [],
    indexTypeList: undefined,
    indexTypeINcouldsNetList: undefined,
    expire: false,
    indexFieldsLastUpdateTime: 0,
    highLightLastUpdateTime: 0,
    userDisPlaysFiled: {},
    attackCharacterization: {
      highLightConfig: [],
      highLightWhite: [],
    },
  }),
  getters: {
    getToken: (state) => state.token,
    getUsername: (state) => state.username,
    getAvatar: (state) => state.avatar,
  },
  actions: {
    /**
     * @description '获取getAll字段更新的最近时间'
     */
    getIndexFieldsUpdateTime() {
      return this.indexFieldsLastUpdateTime
    },
    /**
     * @description '设置getAll字段更新的最近时间'
     */
    setIndexFieldsUpdateTime(val: number) {
      this.indexFieldsLastUpdateTime = val
    },
    sethighLightLastUpdateTime(val: number) {
      this.highLightLastUpdateTime = val
    },
    /**
     * @description 设置token
     * @param {*} token
     */
    setToken(token: string) {
      this.token = token
      setToken(token)
    },
    /**
     * @description 设置token
     * @param {*} token
     */
    setExpire(expire: boolean) {
      this.expire = expire
    },
    /**
     * @description 设置用户名
     * @param {*} username
     */
    setUsername(username: string) {
      this.username = username
    },
    /**
     * @description 设置用户ID
     * @param {*} userId
     */
    setUserId(userId: number) {
      this.userId = userId
    },
    getUserId() {
      return this.userId
    },
    /**
     * @description 设置头像
     * @param {*} avatar
     */
    setAvatar(avatar: string) {
      this.avatar = avatar
    },
    /**
     * @description 登录拦截放行时，设置虚拟角色
     */
    setVirtualRoles() {
      const aclStore = useAclStore()
      aclStore.setFull(true)
      this.setUsername('admin(未开启登录拦截)')
      this.setAvatar('/download/upload/avatar.gif')
    },
    /**
     * @description 设置token并发送提醒
     * @param {string} token 更新令牌
     * @param {string} tokenName 令牌名称
     */
    async afterLogin(token: string, tokenName: string) {
      const settingsStore = useSettingsStore()
      if (token) {
        this.setToken(token)
        const myRout = asyncRoutes.slice(0, -1)
        await updateSystemMenusApi(myRout)
        const hour = new Date().getHours()
        const thisTime =
          hour < 8 ? '早上好' : hour <= 11 ? '上午好' : hour <= 13 ? '中午好' : hour < 18 ? '下午好' : '晚上好'
        gp.$baseNotify(`欢迎登录${settingsStore.title}`, `${thisTime}！`)
        router.push('/')
      } else {
        const err = `登录接口异常，未正确返回${tokenName}...`
        gp.$baseMessage(err, 'error', 'vab-hey-message-error')
        throw err
      }
    },
    /**
     * @description 登录
     * @param {*} userInfo
     */
    async login(userInfo: LoginPayload) {
      return await loginApi(userInfo)
    },
    /**
     * @description 登录
     * @param {*} userInfo
     */
    async loginOther(userInfo: LoginOtherType) {
      return await loginOtherTypeApi({ ...userInfo, username: window.btoa(userInfo.username) })
    },
    /**
     * @description 第三方登录
     * @param {*} tokenData
     */
    async socialLogin(tokenData: any) {
      const {
        data: { [tokenName]: token },
      } = await socialLogin(tokenData)
      this.afterLogin(token, tokenName)
    },
    /**
     * @description 获取用户信息接口 这个接口非常非常重要，如果没有明确底层前逻辑禁止修改此方法，错误的修改可能造成整个框架无法正常使用
     * @returns
     */
    async getUserInfo() {
      const {
        data: { nickName, avatar, roles, permissions, id, userType, theme },
      } = await getUserInfoApi()
      const settingsStore = useSettingsStore()
      /**
       * 检验返回数据是否正常，无对应参数，将使用默认用户名,头像,Roles和Permissions
       * username {String}
       * avatar {String}
       * roles {List}
       * ability {List}
       */
      if (
        (nickName && !isString(nickName)) ||
        (avatar && !isString(avatar)) ||
        (permissions && !isArray(permissions))
      ) {
        const err = 'getUserInfo核心接口异常，请检查返回JSON格式是否正确'
        gp.$baseMessage(err, 'error', 'vab-hey-message-error')
        throw err
      } else {
        const aclStore = useAclStore()
        // 如不使用username用户名,可删除以下代码
        if (nickName) this.setUsername(nickName)
        if (id) this.setUserId(id)
        // 如不使用avatar头像,可删除以下代码
        if (avatar) this.setAvatar(avatar)
        // 如不使用roles权限控制,可删除以下代码
        if (userType === 1) {
          aclStore.setRole(['Admin'])
        }
        if (roles) aclStore.setRole(roles)
        // 如不使用permissions权限控制,可删除以下代码
        if (permissions) aclStore.setPermission(permissions)
        if (theme) settingsStore.updateState({ theme: JSON.parse(theme) })
        // 获取表格头
        this.getTableColumns()
        this.setAllIndexType()
        this.setAllIndexINcouldsNetType()
        this.getSysMenus()
        this.getAllDisPlaysFiledList()
      }
    },
    getAllIndexType() {
      return this.indexTypeList || []
    },
    async setAllIndexType() {
      const { data } = await getAllIndexTypeApi()
      this.indexTypeList = data
    },
    getAllIndexINcouldsNetType() {
      return this.indexTypeINcouldsNetList || []
    },
    async setAllIndexINcouldsNetType() {
      const { data } = await getAllIndexTypeIncouldsNetApi()
      this.indexTypeINcouldsNetList = data
    },

    getAllField() {
      return this.tableColumns
    },
    async getTableColumns() {
      if (localStorage.getItem('allIndexTypeField')) {
        this.tableColumns = JSON.parse(localStorage.getItem('allIndexTypeField')!)
      } else {
        const { data } = await getAllFieldApi()
        this.tableColumns = data
        localStorage.setItem('allIndexTypeField', JSON.stringify(this.tableColumns))
      }
    },
    async setTableColumns() {
      const { data } = await getAllFieldApi()
      this.tableColumns = data
      localStorage.setItem('allIndexTypeField', JSON.stringify(this.tableColumns))
    },
    async getAllDisPlaysFiledList() {
      const { data } = await getAllDisPlaysFiledApi()
      this.userDisPlaysFiled = data
    },
    getDisPlaysFiled(tableKey: RetrieveIndexType) {
      return this.userDisPlaysFiled ? (this.userDisPlaysFiled[tableKey] as number[]) : []
    },
    getTableColumn(tableKey: RetrieveIndexType) {
      return this.tableColumns ? this.tableColumns[tableKey] : []
    },
    async getSysMenus() {
      const { data } = await getSystemMenusApi()
      this.systemMenus = data
    },
    async getAllHightLight() {
      const { data } = await getAllHightLightAPI()
      this.attackCharacterization = data
    },
    /**
     * @description 退出登录
     */
    async logout() {
      await logoutApi()
      await this.resetAll()
      // 解决横向布局退出登录显示不全的bug
      location.reload()
    },
    /**
     * @description 重置token、roles、permission、router、tabsBar等
     */
    async resetAll() {
      this.setToken('')
      this.setUsername('游客')
      this.setExpire(false)
      this.setAvatar('/download/upload/avatar.gif')
      this.setIndexFieldsUpdateTime(0)
      const aclStore = useAclStore()
      const routesStore = useRoutesStore()
      const publicStore = usePubilcStore()
      const tabsStore = useTabsStore()
      aclStore.setPermission([])
      aclStore.setFull(false)
      aclStore.setRole([])
      tabsStore.delAllVisitedRoutes()
      routesStore.clearRoutes()
      publicStore.cleanAll()
      resetRouter()
      removeToken()
      localStorage.clear()
    },
  },
})
